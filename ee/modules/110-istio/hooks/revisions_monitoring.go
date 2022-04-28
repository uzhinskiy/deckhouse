/*
Copyright 2022 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package hooks

import (
	"context"
	"encoding/json"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook/metrics"

	v1 "k8s.io/api/core/v1"

	"github.com/flant/addon-operator/pkg/module_manager/go_hook"
	"github.com/flant/addon-operator/sdk"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/deckhouse/deckhouse/ee/modules/110-istio/hooks/internal"
	"github.com/deckhouse/deckhouse/go_lib/dependency"
)

var (
	revisionsMonitoringMetricsGroup              = "revisions"
	revisionsMonitoringActualNEDesiredMetricName = "d8_istio_revision_actual_ne_desired"
	revisionsMonitoringDesiredIsNotConfigured    = "d8_istio_desired_revision_is_not_configured"
	revisionsMonitoringNoDesired                 = "d8_istio_no_desired_revision"
)

var _ = sdk.RegisterFunc(&go_hook.HookConfig{
	Queue: internal.Queue("revisions-discovery-monitoring"),
	Kubernetes: []go_hook.KubernetesConfig{
		{
			Name:          "namespaces_global_revision",
			ApiVersion:    "v1",
			Kind:          "Namespace",
			FilterFunc:    applyNamespaceFilter, // from revisions_discovery.go
			LabelSelector: &metav1.LabelSelector{MatchLabels: map[string]string{"istio-injection": "enabled"}},
		},
		{
			Name:       "namespaces_definite_revision",
			ApiVersion: "v1",
			Kind:       "Namespace",
			FilterFunc: applyNamespaceFilter, // from revisions_discovery.go
			LabelSelector: &metav1.LabelSelector{
				MatchExpressions: []metav1.LabelSelectorRequirement{
					{
						Key:      "istio.io/rev",
						Operator: "Exists",
					},
				},
			},
		},
	},
	Schedule: []go_hook.ScheduleConfig{ // Due to we are afraid of subscribing to all Pods in the cluster,
		{Name: "cron", Crontab: "10 * * * *"}, // we run the hook every 10 minutes to discover orphan Pods.
	},
}, dependency.WithExternalDependencies(revisionsMonitoring))

type IstioDrivenPod v1.Pod
type IstioPodStatus struct {
	Revision string `json:"revision"`
	// ... we aren't interested in the other columns
}

func (p *IstioDrivenPod) getIstioRevision() string {
	var istioPodStatus IstioPodStatus

	istioStatusJSON := p.Annotations["sidecar.istio.io/status"]
	json.Unmarshal([]byte(istioStatusJSON), &istioPodStatus)

	return istioPodStatus.Revision
}

func revisionsMonitoring(input *go_hook.HookInput, dc dependency.Container) error {
	// isn't discovered yet
	if !input.Values.Get("istio.internal.globalRevision").Exists() {
		return nil
	}
	if !input.Values.Get("istio.internal.revisionsToInstall").Exists() {
		return nil
	}

	input.MetricsCollector.Expire(revisionsMonitoringMetricsGroup)

	var globalRevision = input.Values.Get("istio.internal.globalRevision").String()
	var revisionsToInstall = make([]string, 0)
	var revisionsToInstallResult = input.Values.Get("istio.internal.revisionsToInstall").Array()
	for _, revisionResult := range revisionsToInstallResult {
		revisionsToInstall = append(revisionsToInstall, revisionResult.String())
	}

	var namespaceRevisionMap = map[string]string{}
	for _, ns := range append(input.Snapshots["namespaces_definite_revision"], input.Snapshots["namespaces_global_revision"]...) {
		nsInfo := ns.(NamespaceInfo)
		if nsInfo.Revision == "global" {
			namespaceRevisionMap[nsInfo.Name] = globalRevision
		} else {
			namespaceRevisionMap[nsInfo.Name] = nsInfo.Revision
		}
	}

	k8sClient, err := dc.GetK8sClient()
	if err != nil {
		return err
	}

	podList, err := k8sClient.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{LabelSelector: "service.istio.io/canonical-name,sidecar.istio.io/inject!=false"})
	if err != nil {
		return err
	}
	for _, pod := range podList.Items {
		istioPod := IstioDrivenPod(pod)

		// if the Pod has label 'istio.io/rev' to configure definite revision
		if definiteRevision, ok := istioPod.Labels["istio.io/rev"]; ok {
			if !internal.Contains(revisionsToInstall, definiteRevision) {
				// ALARM! Desired revision isn't configured to install
				labels := map[string]string{
					"namespace":        istioPod.GetNamespace(),
					"desired_revision": definiteRevision,
				}
				input.MetricsCollector.Set(revisionsMonitoringDesiredIsNotConfigured, 1, labels, metrics.WithGroup(revisionsMonitoringMetricsGroup))
			}
			continue
		}

		// the Pod doesn't have definite revision, let's check its Namespace
		// if the Pod's Namespace has istio injection configured
		if revision, ok := namespaceRevisionMap[istioPod.GetNamespace()]; ok {
			if istioPod.getIstioRevision() != revision {
				// ALARM! The Pod's revision isn't equal the desired one, after Pod recreating, the actual revision will be changed
				labels := map[string]string{
					"namespace":        istioPod.GetNamespace(),
					"desired_revision": revision,
					"actual_revision":  istioPod.getIstioRevision(),
				}
				input.MetricsCollector.Set(revisionsMonitoringActualNEDesiredMetricName, 1, labels, metrics.WithGroup(revisionsMonitoringMetricsGroup))
			}
			continue
		} else {
			istioPod.getIstioRevision()
			// ALARM! The pod is orphaned, pod restarting will lead to sidecar miss
			labels := map[string]string{
				"namespace":       istioPod.GetNamespace(),
				"actual_revision": istioPod.getIstioRevision(),
			}
			input.MetricsCollector.Set(revisionsMonitoringNoDesired, 1, labels, metrics.WithGroup(revisionsMonitoringMetricsGroup))
			continue
		}
	}

	return nil
}
