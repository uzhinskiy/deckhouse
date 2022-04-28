/*
Copyright 2022 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package hooks

import (
	"context"

	"k8s.io/utils/pointer"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/yaml"

	"github.com/flant/shell-operator/pkg/metric_storage/operation"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/deckhouse/deckhouse/go_lib/dependency"
	. "github.com/deckhouse/deckhouse/testing/hooks"
)

var _ = Describe("Istio hooks :: revisions_monitoring ::", func() {
	f := HookExecutionConfigInit(`{"istio":{"internal":{}}}`, "")

	Context("Empty cluster and minimal settings", func() {
		BeforeEach(func() {
			f.RunHook()
		})

		It("Hook must execute successfully", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(string(f.LogrusOutput.Contents())).To(HaveLen(0))

			m := f.MetricsCollector.CollectedMetrics()
			Expect(m).To(HaveLen(0))

		})
	})

	Context("Empty cluster and revisions are discovered", func() {
		BeforeEach(func() {
			f.ValuesSet("istio.internal.globalRevision", "v1x42")
			f.ValuesSet("istio.internal.revisionsToInstall", []string{})
			f.RunHook()
		})

		It("Hook must execute successfully", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(string(f.LogrusOutput.Contents())).To(HaveLen(0))

			m := f.MetricsCollector.CollectedMetrics()
			Expect(m).To(HaveLen(1))
			Expect(m[0]).To(BeEquivalentTo(operation.MetricOperation{
				Group:  revisionsMonitoringMetricsGroup,
				Action: "expire",
			}))
		})
	})

	Context("There are different desired and actual revisions", func() {
		BeforeEach(func() {
			f.ValuesSet("istio.internal.globalRevision", "v1x42")
			f.ValuesSet("istio.internal.revisionsToInstall", []string{"v1x15", "v1x42"})

			namespacesYAML := `
---
apiVersion: v1
kind: Namespace
metadata:
  labels:
    istio-injection: enabled
  name: ns-global
---
apiVersion: v1
kind: Namespace
metadata:
  labels:
    istio.io/rev: v1x15
  name: ns-rev1x15
---
apiVersion: v1
kind: Namespace
metadata:
  name: ns-nodesired
`
			podsYAML := []string{
				`---
apiVersion: v1
kind: Pod
metadata:
  name: pod-to-ignore
  namespace: ns-global
  labels:
    sidecar.istio.io/inject: false
    service.istio.io/canonical-name: qqq
  annotations:
    sidecar.istio.io/status: '{"a":"b","revision":"v1x13"}'
`,
				`---
apiVersion: v1
kind: Pod
metadata:
  name: pod-definite-revision-installed
  namespace: ns-nodesired
  labels:
    istio.io/rev: v1x15
    service.istio.io/canonical-name: qqq
`,
				`---
apiVersion: v1
kind: Pod
metadata:
  name: pod-definite-revision-not-installed
  namespace: ns-nodesired
  labels:
    istio.io/rev: v1xexotic
    service.istio.io/canonical-name: qqq
`,
				`---
apiVersion: v1
kind: Pod
metadata:
  name: pod-global-revision-actual
  namespace: ns-global
  labels:
    service.istio.io/canonical-name: qqq
  annotations:
    sidecar.istio.io/status: '{"a":"b","revision":"v1x42"}'
`,
				`---
apiVersion: v1
kind: Pod
metadata:
  name: pod-global-revision-not-actual
  namespace: ns-global
  labels:
    service.istio.io/canonical-name: qqq
  annotations:
    sidecar.istio.io/status: '{"a":"b","revision":"v1x15"}'
`,
				`---
apiVersion: v1
kind: Pod
metadata:
  name: pod-definite-ns-revision-actual
  namespace: ns-rev1x15
  labels:
    service.istio.io/canonical-name: qqq
  annotations:
    sidecar.istio.io/status: '{"a":"b","revision":"v1x15"}'
`,
				`---
apiVersion: v1
kind: Pod
metadata:
  name: pod-definite-ns-revision-not-actual
  namespace: ns-rev1x15
  labels:
    service.istio.io/canonical-name: qqq
  annotations:
    sidecar.istio.io/status: '{"a":"b","revision":"v1x42"}'
`,
				`---
apiVersion: v1
kind: Pod
metadata:
  name: pod-orphan
  namespace: ns-nodesired
  labels:
    service.istio.io/canonical-name: qqq
  annotations:
    sidecar.istio.io/status: '{"a":"b","revision":"v1x15"}'
`,
			}

			f.BindingContexts.Set(f.KubeStateSet(namespacesYAML))

			for _, podYAML := range podsYAML {
				var pod v1.Pod
				var err error
				err = yaml.Unmarshal([]byte(podYAML), &pod)
				Expect(err).To(BeNil())

				_, err = dependency.TestDC.MustGetK8sClient().
					CoreV1().
					Pods(pod.GetNamespace()).
					Create(context.TODO(), &pod, metav1.CreateOptions{})
				Expect(err).To(BeNil())
			}

			f.RunHook()
		})

		It("Hook must execute successfully", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(string(f.LogrusOutput.Contents())).To(HaveLen(0))

			m := f.MetricsCollector.CollectedMetrics()
			Expect(m).To(HaveLen(5))
			Expect(m[0]).To(BeEquivalentTo(operation.MetricOperation{
				Group:  revisionsMonitoringMetricsGroup,
				Action: "expire",
			}))
			Expect(m[1]).To(BeEquivalentTo(operation.MetricOperation{
				Name:   "d8_istio_revision_actual_ne_desired",
				Group:  revisionsMonitoringMetricsGroup,
				Action: "set",
				Value:  pointer.Float64Ptr(1.0),
				Labels: map[string]string{
					"actual_revision":  "v1x15",
					"namespace":        "ns-global",
					"desired_revision": "v1x42",
				},
			}))
			Expect(m[2]).To(BeEquivalentTo(operation.MetricOperation{
				Name:   "d8_istio_desired_revision_is_not_configured",
				Group:  revisionsMonitoringMetricsGroup,
				Action: "set",
				Value:  pointer.Float64Ptr(1.0),
				Labels: map[string]string{
					"desired_revision": "v1xexotic",
					"namespace":        "ns-nodesired",
				},
			}))
			Expect(m[3]).To(BeEquivalentTo(operation.MetricOperation{
				Name:   "d8_istio_no_desired_revision",
				Group:  revisionsMonitoringMetricsGroup,
				Action: "set",
				Value:  pointer.Float64Ptr(1.0),
				Labels: map[string]string{
					"actual_revision": "v1x15",
					"namespace":       "ns-nodesired",
				},
			}))
			Expect(m[4]).To(BeEquivalentTo(operation.MetricOperation{
				Name:   "d8_istio_revision_actual_ne_desired",
				Group:  revisionsMonitoringMetricsGroup,
				Action: "set",
				Value:  pointer.Float64Ptr(1.0),
				Labels: map[string]string{
					"actual_revision":  "v1x42",
					"desired_revision": "v1x15",
					"namespace":        "ns-rev1x15",
				},
			}))
		})
	})
})
