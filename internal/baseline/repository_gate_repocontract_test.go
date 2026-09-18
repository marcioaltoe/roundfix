// Suite: repository gate composition
// Invariant: both verification tiers run an analyzer that rejects a known diagnostic.
// Boundary IN: the repository Makefile and the Go toolchain's analyzer.
// Boundary OUT: the analyzer target implementation, which the Makefile owns.

//go:build repocontract

package baseline

import (
	"bytes"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestRepositoryGateRunsTheAnalyzer(t *testing.T) {
	repository := filepath.Clean(filepath.Join("..", ".."))
	command := exec.CommandContext(
		t.Context(),
		"go", "vet", "./internal/baseline/analyzer/testdata/printfdiagnostic",
	)
	command.Dir = repository
	output, err := command.CombinedOutput()
	var exitError *exec.ExitError
	if !errors.As(err, &exitError) || exitError.ExitCode() == 0 {
		t.Fatalf("go vet diagnostic fixture error = %v, want a non-zero exit; output:\n%s", err, output)
	}
	if !bytes.Contains(output, []byte("fmt.Printf format %d")) {
		t.Fatalf("go vet diagnostic fixture output:\n%s\nwant diagnostic naming fmt.Printf format %%d", output)
	}

	makefile, err := os.ReadFile(filepath.Join(repository, "Makefile"))
	if err != nil {
		t.Fatalf("read repository gate composition: %v", err)
	}
	for _, target := range []string{"verify", "verify-incremental"} {
		dependencies := makeTargetDependencies(t, makefile, target)
		if !containsMakeDependency(dependencies, "vet") {
			t.Errorf("Makefile target %q dependencies = %v, want analyzer target vet", target, dependencies)
		}
	}
}

func makeTargetDependencies(t *testing.T, makefile []byte, target string) []string {
	t.Helper()
	prefix := target + ":"
	for _, line := range strings.Split(string(makefile), "\n") {
		if strings.HasPrefix(line, prefix) {
			return strings.Fields(strings.TrimPrefix(line, prefix))
		}
	}
	t.Fatalf("Makefile has no %q target", target)
	return nil
}

func containsMakeDependency(values []string, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
