/*
Copyright 2022 Flant JSC
Licensed under the Deckhouse Platform Enterprise Edition (EE) license. See https://github.com/deckhouse/deckhouse/blob/main/ee/LICENSE
*/

package hooks

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/deckhouse/deckhouse/testing/hooks"
)

var _ = FDescribe("Istio hooks :: revisions_discovery ::", func() {
	f := HookExecutionConfigInit(`{"istio":{}}`, "")
	f.RegisterCRD("install.istio.io", "v1alpha1", "IstioOperator", true)

	Context("Empty cluster and minimal settings", func() {
		BeforeEach(func() {
			values := `
internal:
  supportedVersions: ["1.2.3-beta.45"]
`
			f.ValuesSetFromYaml("istio", []byte(values))

			f.BindingContexts.Set(f.KubeStateSet(``))
			f.RunHook()
		})

		It("Hook must execute successfully", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(f.LogrusOutput.Contents()).To(HaveLen(0))

			Expect(f.ConfigValuesGet("istio.globalVersion").String()).To(Equal("1.2.3-beta.45"))
			Expect(f.ValuesGet("istio.internal.applicationNamespaces").Array()).To(BeEmpty())
			Expect(f.ValuesGet("istio.internal.revisionsToInstall").String()).To(MatchJSON(`["v1x2x3beta45"]`))
			Expect(f.ValuesGet("istio.internal.operatorRevisionsToInstall").String()).To(MatchJSON(`["v1x2x3beta45"]`))
			Expect(f.ValuesGet("istio.internal.globalRevision").String()).To(Equal("v1x2x3beta45"))
		})
	})

	Context("Application namespaces with labels and IstioOperator", func() {
		BeforeEach(func() {
			values := `
globalVersion: "1.1.0"
additionalVersions: ["1.4","1.3"]
internal:
  supportedVersions: ["1.1.0", "1.8.0-alpha.2", "1.3", "1.4"]
`
			f.ValuesSetFromYaml("istio", []byte(values))
			f.BindingContexts.Set(f.KubeStateSet(`
---
# regular ns
apiVersion: v1
kind: Namespace
metadata:
  name: ns0
---
# ns with global revision
apiVersion: v1
kind: Namespace
metadata:
  name: ns1
  labels:
    istio-injection: enabled
---
# ns with global revision
apiVersion: v1
kind: Namespace
metadata:
  name: ns2
  labels:
    istio-injection: enabled
---
# ns with definite revision
apiVersion: v1
kind: Namespace
metadata:
  name: ns3
  labels:
    istio.io/rev: v1x7x4
---
# ns with definite revision
apiVersion: v1
kind: Namespace
metadata:
  name: ns4
  labels:
    istio.io/rev: v1x5x0
---
# ns with definite revision
apiVersion: v1
kind: Namespace
metadata:
  name: ns5
  labels:
    istio.io/rev: v1x7x4
---
# ns with definite revision with d8 prefix
apiVersion: v1
kind: Namespace
metadata:
  name: d8-ns6
  labels:
    istio.io/rev: v1x8x0
---
# ns with global revision with d8 prefix
apiVersion: v1
kind: Namespace
metadata:
  name: d8-ns7
  labels:
    istio-injection: enabled
---
# ns with definitee revision with kube prefix
apiVersion: v1
kind: Namespace
metadata:
  name: kube-ns8
  labels:
    istio.io/rev: v1x9x0
---
# ns with global revision with kube prefix
apiVersion: v1
kind: Namespace
metadata:
  name: kube-ns9
  labels:
    istio-injection: enabled
---
apiVersion: install.istio.io/v1alpha1
kind: IstioOperator
metadata:
  name: v1x8x0alpha2
  namespace: d8-istio
spec:
  revision: v1x8x0alpha2
`))

			f.RunHook()
		})
		It("Should count all namespaces and revisions properly", func() {
			Expect(f).To(ExecuteSuccessfully())
			Expect(f.ValuesGet("istio.internal.applicationNamespaces").AsStringSlice()).To(Equal([]string{"d8-ns6", "d8-ns7", "kube-ns8", "kube-ns9", "ns1", "ns2", "ns3", "ns4", "ns5"}))
			Expect(f.ValuesGet("istio.internal.revisionsToInstall").AsStringSlice()).To(Equal([]string{"v1x1x0", "v1x3", "v1x4"}))
			Expect(f.ValuesGet("istio.internal.operatorRevisionsToInstall").AsStringSlice()).To(Equal([]string{"v1x1x0", "v1x3", "v1x4", "v1x8x0alpha2"}))
			Expect(f.ValuesGet("istio.internal.globalRevision").String()).To(Equal("v1x1x0"))
			Expect(f.ConfigValuesGet("istio.globalVersion").String()).To(Equal("1.1.0"))
		})
	})

	Context("Unsupported versions", func() {
		BeforeEach(func() {
			values := `
globalVersion: "1.2.3-beta.45"
additionalVersions: ["1.7.4", "1.1.0", "1.8.0-alpha.2", "1.3.1", "1.9.0"]
internal:
  supportedVersions: ["1.1.0","1.2.3-beta.45","1.3.1"]
`
			f.ValuesSetFromYaml("istio", []byte(values))
			f.RunHook()
		})
		It("Should return errors", func() {
			Expect(f).ToNot(ExecuteSuccessfully())

			Expect(f.GoHookError).To(MatchError("unsupported revisions: [v1x7x4,v1x8x0alpha2,v1x9x0]"))
		})
	})
})
