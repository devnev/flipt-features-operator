/*
Copyright (C) 2024.

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package e2e

import (
	"fmt"
	"os/exec"
	"slices"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/devnev/flipt-features-operator/test/utils"
)

const operatorNamespace = "flipt-features-operator-system"
const targetsNamespace = "e2e-test-flipt"
const appNamespace = "e2e-test-app"

var _ = Describe("controller", Ordered, func() {
	BeforeAll(func() {
		By("installing prometheus operator")
		Expect(utils.InstallPrometheusOperator()).To(Succeed())

		By("installing the cert-manager")
		Expect(utils.InstallCertManager()).To(Succeed())

		By("creating operator namespace")
		cmd := exec.Command("kubectl", "delete", "ns", operatorNamespace)
		_, _ = utils.Run(cmd)
		cmd = exec.Command("kubectl", "create", "ns", operatorNamespace)
		_, _ = utils.Run(cmd)
	})

	AfterAll(func() {
		By("uninstalling the Prometheus manager bundle")
		utils.UninstallPrometheusOperator()

		By("uninstalling the cert-manager bundle")
		utils.UninstallCertManager()

		By("deleting testing namespaces")
		for _, ns := range []string{operatorNamespace, targetsNamespace, appNamespace} {
			cmd := exec.Command("kubectl", "delete", "ns", ns)
			_, _ = utils.Run(cmd)
		}
	})

	Context("Operator on clean cluster", func() {
		const projectimage = "devnev/flipt-features-operator:e2e-test"
		It("should run successfully", func() {
			By("building the manager(Operator) image")
			cmd := exec.Command("make", "docker-build", fmt.Sprintf("IMG=%s", projectimage))
			Expect(utils.Run(cmd)).Error().NotTo(HaveOccurred())

			By("loading the the manager(Operator) image on Kind")
			err := utils.LoadImageToKindClusterWithName(projectimage)
			Expect(err).NotTo(HaveOccurred())

			By("deploying the controller-manager")
			cmd = exec.Command("make", "deploy", fmt.Sprintf("IMG=%s", projectimage))
			Expect(utils.Run(cmd)).Error().NotTo(HaveOccurred())

			By("validating that the controller-manager pod is running as expected")
			expectControllerToComeup()
		})

		It("should gather newly created Features and Targets", func() {
			By("deploying testdata manifests")
			cmd := exec.Command("kubectl", "apply", "-k", "test/e2e/testdata/targets_and_features")
			Expect(utils.Run(cmd)).Error().NotTo(HaveOccurred())

			By("validating that configmaps are created")
			expectConfigMapWithKeys(targetsNamespace, "flipt-all-features-cm", []string{"e2e-test-app_features-appns.yml", "e2e-test-app_features-nons.yml", "e2e-test-app_features-testonly.yml"})
			expectConfigMapWithKeys(targetsNamespace, "flipt-testonly-features-cm", []string{"e2e-test-app_features-testonly.yml"})
		})

		It("should shut down", func() {
			By("undeploying the controller-manager")
			cmd := exec.Command("make", "undeploy", fmt.Sprintf("IMG=%s", projectimage))
			Expect(utils.Run(cmd)).Error().NotTo(HaveOccurred())
		})
	})

	Context("Operator with existing resources", func() {
		const projectimage = "devnev/flipt-features-operator:e2e-test"
		It("should run successfully", func() {
			By("building the manager(Operator) image")
			cmd := exec.Command("make", "docker-build", fmt.Sprintf("IMG=%s", projectimage))
			Expect(utils.Run(cmd)).Error().NotTo(HaveOccurred())

			By("loading the the manager(Operator) image on Kind")
			err := utils.LoadImageToKindClusterWithName(projectimage)
			Expect(err).NotTo(HaveOccurred())

			By("installing CRDs")
			cmd = exec.Command("make", "install")
			Expect(utils.Run(cmd)).Error().NotTo(HaveOccurred())

			By("deploying testdata manifests")
			cmd = exec.Command("kubectl", "apply", "-k", "test/e2e/testdata/targets_and_features")
			Expect(utils.Run(cmd)).Error().NotTo(HaveOccurred())

			By("clearing any test configmaps")
			cmd = exec.Command("kubectl", "delete", "cm", "-n", targetsNamespace, "--all")
			Expect(utils.Run(cmd)).Error().NotTo(HaveOccurred())

			By("deploying the controller-manager")
			cmd = exec.Command("make", "deploy", fmt.Sprintf("IMG=%s", projectimage))
			Expect(utils.Run(cmd)).Error().NotTo(HaveOccurred())

			By("validating that the controller-manager pod is running as expected")
			expectControllerToComeup()
		})

		It("should gather existing Features and Targets", func() {
			By("validating that configmaps are created")
			expectConfigMapWithKeys(targetsNamespace, "flipt-all-features-cm", []string{"e2e-test-app_features-appns.yml", "e2e-test-app_features-nons.yml", "e2e-test-app_features-testonly.yml"})
			expectConfigMapWithKeys(targetsNamespace, "flipt-testonly-features-cm", []string{"e2e-test-app_features-testonly.yml"})
		})

		It("should shut down", func() {
			By("undeploying the controller-manager")
			cmd := exec.Command("make", "undeploy", fmt.Sprintf("IMG=%s", projectimage))
			Expect(utils.Run(cmd)).Error().NotTo(HaveOccurred())
		})
	})
})

func expectControllerToComeup() {
	verifyControllerUp := func() error {
		// Get pod name
		podOutput, err := utils.Run(exec.Command("kubectl", "get",
			"pods", "-l", "control-plane=controller-manager",
			"-o", "go-template={{ range .items }}"+
				"{{ if not .metadata.deletionTimestamp }}"+
				"{{ .metadata.name }}"+
				"{{ \"\\n\" }}{{ end }}{{ end }}",
			"-n", operatorNamespace,
		))
		Expect(err).NotTo(HaveOccurred())
		podNames := utils.GetNonEmptyLines(string(podOutput))
		if len(podNames) != 1 {
			return fmt.Errorf("expect 1 controller pods running, but got %d", len(podNames))
		}
		controllerPodName := podNames[0]
		Expect(controllerPodName).Should(ContainSubstring("controller-manager"))

		// Validate pod status
		status, err := utils.Run(exec.Command("kubectl", "get",
			"pods", controllerPodName, "-o", "jsonpath={.status.phase}",
			"-n", operatorNamespace,
		))
		Expect(err).NotTo(HaveOccurred())
		if string(status) != "Running" {
			return fmt.Errorf("controller pod in %s status", status)
		}
		return nil
	}
	Eventually(verifyControllerUp, time.Minute, time.Second).Should(Succeed())
}

func expectConfigMapWithKeys(ns string, name string, keys []string) {
	slices.Sort(keys)
	verifyConfigMapHasKeys := func() error {
		cmdOutput, err := utils.Run(exec.Command(
			"kubectl", "get", "cm", "-n", ns, name,
			"-o", `go-template={{ range $key, $item := .data }}{{ printf "%s\n" $key }}{{ end }}`,
		))
		if err != nil && strings.HasPrefix(string(cmdOutput), "Error from server (NotFound):") {
			return err
		}
		Expect(err).NotTo(HaveOccurred())
		cmKeys := utils.GetNonEmptyLines(string(cmdOutput))
		slices.Sort(cmKeys)
		if !slices.Equal(keys, cmKeys) {
			return fmt.Errorf("expected configmap keys %q, but got %q", keys, cmKeys)
		}

		return nil
	}
	Eventually(verifyConfigMapHasKeys, time.Minute, time.Second).Should(Succeed())
}
