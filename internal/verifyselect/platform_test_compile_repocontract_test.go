//go:build repocontract

// Suite: Linux and Windows test compilation.
// Invariant: every repository test compiles for the release platforms no gate tests.
// Boundary IN: repository Go sources with both contract tags and a temporary module.
// Boundary OUT: executing cross-platform test binaries and mutating the checkout.
package verifyselect_test

import (
	"context"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// compiledPlatforms are the release platforms no gate runs a test on.
var compiledPlatforms = []string{"linux", "windows"}

// compileEveryTest runs go vet over every package of dir for goos.
func compileEveryTest(ctx context.Context, dir, goos string) ([]byte, error) {
	command := exec.CommandContext(ctx, "go", "vet", "-tags", "repocontract docscontract", "./...")
	command.Dir = dir
	command.Env = append(os.Environ(), "GOOS="+goos)
	return command.CombinedOutput()
}

func TestEveryTestCompilesForLinuxAndWindows(t *testing.T) {
	t.Parallel()
	for _, platform := range compiledPlatforms {
		t.Run(platform, func(t *testing.T) {
			if output, err := compileEveryTest(t.Context(), findRepositoryRoot(t), platform); err != nil {
				t.Fatalf("compile every test for %s: %v\n%s", platform, err, output)
			}
		})
	}
	t.Run("a_file_only_linux_fails_to_compile", func(t *testing.T) {
		dir := t.TempDir()
		writeFile(t, dir, "go.mod", "module compilefixture\n\ngo 1.26\n")
		writeFile(t, dir, "fixture.go", "package compilefixture\n")
		writeFile(t, dir, "broken_linux_test.go", "//go:build linux\n\npackage compilefixture\nimport \"testing\"\nfunc TestBroken(t *testing.T) { undefinedLinuxSymbol() }\n")
		output, err := compileEveryTest(t.Context(), dir, "linux")
		if err == nil || !strings.Contains(string(output), "broken_linux_test.go") {
			t.Fatalf("Linux must reject broken_linux_test.go: %v\n%s", err, output)
		}
		if output, err := compileEveryTest(t.Context(), dir, "windows"); err != nil {
			t.Fatalf("Windows must compile without the Linux-only file: %v\n%s", err, output)
		}
	})
}
