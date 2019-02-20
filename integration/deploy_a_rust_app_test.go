package integration_test

import (
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/cloudfoundry/libbuildpack/cutlass"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("V3 Wrapped CF Rust Buildpack", func() {
	var app *cutlass.App
	AfterEach(func() {
		DestroyApp(app)
	})

	Context("When pushing a simple Rust app", func() {
		BeforeEach(func() {
			app = cutlass.New(filepath.Join(bpDir, "integration", "testdata", "simple_app"))
			app.Disk = "1G"
		})

		It("uses the requested php version and runs successfully", func() {
			Expect(app.Push()).To(Succeed())
			Eventually(func() ([]string, error) { return app.InstanceStates() }, 120*time.Second).Should(Equal([]string{"RUNNING"}))
			Eventually(app.Stdout.ANSIStrippedString).Should(MatchRegexp(`Rust \d+\.\d+\.\d+: Contributing`))
			Eventually(app.Stdout.String).Should(ContainSubstring("OUT SUCCESS"))
		})
	})

	Context("Unbuilt buildpack (eg github)", func() {
		var bpName string

		BeforeEach(func() {
			if cutlass.Cached {
				Skip("skipping cached buildpack test")
			}

			tmpDir, err := ioutil.TempDir("", "")
			Expect(err).NotTo(HaveOccurred())
			defer os.RemoveAll(tmpDir)

			bpName = "unbuilt-v3-php"
			bpZip := filepath.Join(tmpDir, bpName+".zip")

			app = cutlass.New(filepath.Join(bpDir, "integration", "testdata", "simple_app"))
			app.Buildpacks = []string{bpName + "_buildpack"}
			app.Disk = "1G"
			app.HealthCheck = "process"

			cmd := exec.Command("git", "archive", "-o", bpZip, "HEAD")
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			cmd.Dir = bpDir
			Expect(cmd.Run()).To(Succeed())

			Expect(cutlass.CreateOrUpdateBuildpack(bpName, bpZip, "")).To(Succeed())
		})

		AfterEach(func() {
			Expect(cutlass.DeleteBuildpack(bpName)).To(Succeed())
		})

		It("runs", func() {
			Expect(app.Push()).To(Succeed())
			Eventually(func() ([]string, error) { return app.InstanceStates() }, 120*time.Second).Should(Equal([]string{"RUNNING"}))

			Eventually(app.Stdout.ANSIStrippedString).Should(MatchRegexp(`Rust \d+\.\d+\.\d+: Contributing`))
			Eventually(app.Stdout.String).Should(ContainSubstring("OUT SUCCESS"))
		})
	})
})
