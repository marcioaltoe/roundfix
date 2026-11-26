//go:build docscontract

// Suite: One Build repository rule.
// Invariant: tests build binaries only through the package fixture helpers.
// Boundary IN: Go test sources and exec calls.
// Boundary OUT: testfixture's implementation and the frozen governed exception.
//
//verify:always
package testfixture

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Frozen inventory: entries may only shrink.
var oneBuildInventory = map[string]string{
	"internal/cli/baseline_release_gate_test.go": "Governed Path outside this Spec's grant",
}

func TestEveryTestBuildsABinaryOnlyThroughItsPackagesOnePlace(t *testing.T) {
	t.Parallel()
	reports, err := goToolReports("../..", true, oneBuildInventory)
	if err != nil {
		t.Fatal(err)
	}
	for _, report := range reports {
		t.Error(report)
	}
}

func TestTheOneBuildRuleReportsSeededViolations(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, call string
		bad        bool
	}{
		{"build", `exec.Command("go", "build", "./cmd/app")`, true},
		{"run", `exec.CommandContext(nil, "go", "run", "app.go")`, true},
		{"compile", `exec.Command("go", "test", "-c")`, true},
		{"slice", `args := []string{"test", "-c"}; exec.Command("go", args...)`, true},
		{"list", `exec.Command("go", "list")`, false},
		{"test", `exec.Command("go", "test")`, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root := t.TempDir()
			if err := os.WriteFile(filepath.Join(root, "seed_test.go"), []byte("package seeded\nimport \"os/exec\"\nfunc fixture(){"+tc.call+"}"), 0o600); err != nil {
				t.Fatal(err)
			}
			reports, err := goToolReports(root, true, nil)
			if err != nil {
				t.Fatal(err)
			}
			if (len(reports) > 0) != tc.bad {
				t.Fatalf("reports = %q, violation = %v", reports, tc.bad)
			}
			if tc.bad && !strings.Contains(reports[0], "seed_test.go:") {
				t.Fatalf("missing filename: %q", reports)
			}
		})
	}
}
