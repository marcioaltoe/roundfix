package daemon

// IN: real Spec Consistency settlement checker, temporary Git repository.
// OUT: attempt baseline comparison, covered by the settlement checks suite.
// Invariant: a file pinning the checked Spec is a refusing settlement finding.

import (
	"path/filepath"
	"testing"

	"roundfix/internal/speccheck"
)

func TestSettlementRefusesATaskThatPinsItsSpecPath(t *testing.T) {
	t.Parallel()
	fixture, plan, _, _, _ := settlementRepositoryFixture(t, []string{"declared check"})
	checker := SpecCheckSettlementChecker{}
	before, err := checker.RefusingFindings(plan.SpecsRoot, plan.WorkDir, plan.Spec.Slug)
	if err != nil {
		t.Fatal(err)
	}
	for _, finding := range before {
		if finding.Code == speccheck.CodeSpecPathPinned {
			t.Fatalf("unexpected baseline pin: %+v", finding)
		}
	}
	relative, err := filepath.Rel(plan.WorkDir, filepath.Join(plan.SpecsRoot, plan.Spec.Slug))
	if err != nil {
		t.Fatal(err)
	}
	mustWriteForTest(t, filepath.Join(fixture.gitRoot, "pin_test.go"), "package pin\nvar path = \""+filepath.ToSlash(relative)+"/_techspec.md\"\n")
	findings, err := checker.RefusingFindings(plan.SpecsRoot, plan.WorkDir, plan.Spec.Slug)
	if err != nil {
		t.Fatal(err)
	}
	var pins []speccheck.Finding
	for _, finding := range findings {
		if finding.Code == speccheck.CodeSpecPathPinned {
			pins = append(pins, finding)
		}
	}
	if len(pins) != 1 || pins[0].Severity != speccheck.SeverityError || len(pins[0].Where) != 1 || pins[0].Where[0].Path != "pin_test.go" || pins[0].Where[0].Line != 2 {
		t.Fatalf("refusing pins = %+v", pins)
	}
}
