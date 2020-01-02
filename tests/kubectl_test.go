package tests_test

import (
	"strings"

	. "github.com/onsi/ginkgo"
	"github.com/onsi/ginkgo/extensions/table"
	. "github.com/onsi/gomega"

	"kubevirt.io/client-go/kubecli"
	"kubevirt.io/kubevirt/tests"
)

var _ = Describe("[rfe_id:3423][vendor:cnv-qe@redhat.com][level:component]oc/kubectl get vm/vmi tests", func() {
	tests.FlagParse()

	var k8sClient string

	virtCli, err := kubecli.GetKubevirtClient()
	tests.PanicOnError(err)

	BeforeEach(func() {
		k8sClient = tests.GetK8sCmdClient()
		tests.SkipIfNoCmd(k8sClient)
		tests.BeforeTestCleanup()
	})

	table.DescribeTable("should verify set of columns for", func(verb, resource string, expectedHeader []string) {
		vm := tests.NewRandomVirtualMachine(tests.NewRandomVMI(), false)
		vm, err := virtCli.VirtualMachine(tests.NamespaceTestDefault).Create(vm)
		Expect(err).NotTo(HaveOccurred())

		tests.StartVirtualMachine(vm)

		result, _, _ := tests.RunCommand(k8sClient, verb, resource)
		Expect(result).ToNot(BeNil())
		resultFields := strings.Fields(result)
		columnHeaders := resultFields[:len(expectedHeader)]
		// Verify the generated header is same as expected
		Expect(columnHeaders).To(Equal(expectedHeader))
		// Name will be there in all the cases, so verify name
		Expect(resultFields[len(expectedHeader)]).To(Equal(vm.Name))
	},
		table.Entry("[test_id:3464]virtualmachine", "get", "vm", []string{"NAME", "AGE", "RUNNING", "VOLUME"}),
		table.Entry("[test_id:3465]virtualmachineinstance", "get", "vmi", []string{"NAME", "AGE", "PHASE", "IP", "NODENAME"}),
	)

	table.DescribeTable("should verify set of wide columns for", func(verb, resource, option string, expectedHeader []string, verifyPos int, expectedData string) {
		vm := tests.NewRandomVirtualMachine(tests.NewRandomVMI(), false)
		vm, err := virtCli.VirtualMachine(tests.NamespaceTestDefault).Create(vm)
		Expect(err).NotTo(HaveOccurred())

		tests.StartVirtualMachine(vm)

		result, _, _ := tests.RunCommand(k8sClient, verb, resource, "-o", option)

		Expect(result).ToNot(BeNil())
		resultFields := strings.Fields(result)
		columnHeaders := resultFields[:len(expectedHeader)]
		// Verify the generated header is same as expected
		Expect(columnHeaders).To(Equal(expectedHeader))
		// Name will be there in all the cases, so verify name
		Expect(resultFields[len(expectedHeader)]).To(Equal(vm.Name))
		// Verify one of the wide column output field
		Expect(resultFields[len(resultFields)-verifyPos]).To(Equal(expectedData))
	},
		table.Entry("[test_id:3421]virtualmachine", "get", "vm", "wide", []string{"NAME", "AGE", "RUNNING", "VOLUME", "CREATED"}, 1, "true"),
		table.Entry("[test_id:3422]virtualmachineinstance", "get", "vmi", "wide", []string{"NAME", "AGE", "PHASE", "IP", "NODENAME", "LIVE-MIGRATABLE"}, 1, "True"),
	)
})
