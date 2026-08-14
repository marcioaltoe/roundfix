//go:build repocontract

// Suite: repository suite-guard installation contract.
// Invariant: every internal package whose tests spawn a process installs suiteguard.Main.
// Boundary IN: internal Go test files and package-level TestMain wiring.
// Boundary OUT: subprocess lifetime and repository-write behavior inside each package.
package suiteguard_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/suiteguardcontract"
)

func TestEverySpawningPackageInstallsTheSuiteGuard(t *testing.T) {
	t.Run("names an unenumerated and unguarded spawning package", func(t *testing.T) {
		root := t.TempDir()
		writeContractFixture(t, filepath.Join(root, "internal", "uninstalled", "spawn_test.go"), `package uninstalled

import (
	"os/exec"
	"testing"
)

func TestSpawn(t *testing.T) {
	if err := exec.Command("fixture").Run(); err != nil {
		t.Fatal(err)
	}
}
`)

		findings, err := suiteguardcontract.AuditInstallations(root, nil)
		if err != nil {
			t.Fatal(err)
		}
		want := "internal/uninstalled spawns a process but is not enumerated as guarded"
		if !containsString(findings, want) {
			t.Fatalf("findings = %q, want %q", findings, want)
		}
	})

	t.Run("repository installations", func(t *testing.T) {
		root, err := filepath.Abs(filepath.Join("..", ".."))
		if err != nil {
			t.Fatalf("resolve repository root: %v", err)
		}
		findings, err := suiteguardcontract.AuditInstallations(
			root,
			suiteguardcontract.GuardedSpawningPackages(),
		)
		if err != nil {
			t.Fatal(err)
		}
		if len(findings) != 0 {
			t.Fatalf("suite guard installation contract failed:\n%s", strings.Join(findings, "\n"))
		}
	})
}

func writeContractFixture(t *testing.T, filePath, contents string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(filePath), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filePath, []byte(contents), 0o644); err != nil {
		t.Fatal(err)
	}
}

func containsString(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}
