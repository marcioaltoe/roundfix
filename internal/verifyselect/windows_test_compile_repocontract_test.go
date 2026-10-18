//go:build repocontract

// Suite: Windows test compilation.
// Invariant: every repository test compiles for the Windows release platform.
// Boundary IN: repository Go sources with both contract tags and the Go toolchain.
// Boundary OUT: executing Windows test binaries and mutating the checkout.
package verifyselect_test

import (
	"os"
	"os/exec"
	"testing"
)

func TestEveryTestCompilesForWindows(t *testing.T) {
	t.Parallel()
	command := exec.CommandContext(t.Context(), "go", "vet", "-tags", "repocontract docscontract", "./...")
	command.Dir = findRepositoryRoot(t)
	command.Env = append(os.Environ(), "GOOS=windows")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("compile every test for Windows: %v\n%s", err, output)
	}
}
