package integration_test

import (
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

		It("uses the requested rustup version and runs successfully", func() {
			Expect(app.Push()).To(Succeed())
			Eventually(func() ([]string, error) { return app.InstanceStates() }, 180*time.Second).Should(Equal([]string{"RUNNING"}))
			Eventually(app.Stdout.ANSIStrippedString).Should(MatchRegexp(`RustUp \d+\.\d+\.\d+: Contributing`))
			Eventually(app.Stdout.ANSIStrippedString).Should(MatchRegexp(`Rocket has launched from`))
		})
	})
})
