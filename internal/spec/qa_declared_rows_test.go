package spec

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestUnreachableDeclarationNamesEveryCriterion(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name, fields string
		want         []string
		invalid      bool
	}{
		{"repeated and continued", "- criterion: Success Metric 10\n  criterion: Success Metric 11,\n    after release", []string{"Success Metric 10", "Success Metric 11, after release"}, false},
		{"empty first", "- criterion:\n  criterion: Success Metric 11", nil, true},
		{"empty last", "- criterion: Success Metric 10\n  criterion:", nil, true},
		{"missing", "- reason: unreachable", nil, true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			writeDeclaredRowsFile(t, filepath.Join(dir, "_prd.md"), "# PRD\n\n## Unreachable Acceptance\n\n"+tt.fields+"\n  reason: unreachable\n  satisfied-by: later evidence\n")
			got, err := Unreachable(dir)
			if tt.invalid {
				if err == nil || !strings.Contains(err.Error(), "missing criterion") {
					t.Fatalf("Unreachable = %+v, %v; want missing criterion", got, err)
				}
				return
			}
			if err != nil {
				t.Fatal(err)
			}
			if len(got) != 1 || !reflect.DeepEqual(got[0].Criteria, tt.want) || got[0].Criterion != tt.want[0] || got[0].Reason != "unreachable" || got[0].SatisfiedBy != "later evidence" {
				t.Fatalf("declaration = %+v; want criteria %q with shared fields", got, tt.want)
			}
		})
	}
}

func TestCriterionNamedMatchesTextOrLeadingLabel(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name, declared, row string
		want                bool
	}{
		{"label", "Success Metric 5, after release, evidence", "Success Metric 5", true},
		{"normalized", "`Success   Metric 5`. ", " success metric 5 ", true},
		{"colon", "Goal 1: evidence", "Goal 1", true},
		{"semicolon", "Goal 1; evidence", "Goal 1", true},
		{"period", "Goal 1. evidence", "Goal 1", true},
		{"parenthesis", "Goal 1 (evidence)", "Goal 1", true},
		{"em dash", "Goal 1 — evidence", "Goal 1", true},
		{"dash", "Goal 1 - evidence", "Goal 1", true},
		{"digit boundary", "Success Metric 10, evidence", "Success Metric 1", false},
		{"paraphrase", "no Run delivered after this Spec fails the test", "no later Run fails the test", false},
		{"empty", "Goal 1", "", false},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if got := CriterionNamed(tt.declared, tt.row); got != tt.want {
				t.Fatalf("CriterionNamed(%q,%q)=%v; want %v", tt.declared, tt.row, got, tt.want)
			}
		})
	}
	for _, tt := range []struct {
		status, criterion string
		valid             bool
	}{
		{"blocked (declared: Success Metric 5)", "Success Metric 5", true},
		{"BLOCKED (DECLARED:Success Metric 5 )", "Success Metric 5", true},
		{"blocked (declared: )", "", false},
		{"blocked (declared: criterion", "", false},
		{"blocked (environment: criterion)", "", false},
	} {
		t.Run(tt.status, func(t *testing.T) {
			got, ok := QADeclaredRowCriterion(tt.status)
			if got != tt.criterion || ok != tt.valid {
				t.Fatalf("classifier = %q,%v; want %q,%v", got, ok, tt.criterion, tt.valid)
			}
		})
	}
}

func TestQAReportEligibilityMatchesEachDeclaredRow(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name, id, criterion string
		count               int
		refusal             string
	}{
		{"two criteria", "Q11", "Success Metric 11", 2, ""},
		{"undeclared criterion", "Q12", "Success Metric 12", 2, "declared-blocked row Q12"},
		{"count exceeds listed rows", "Q11", "Success Metric 11", 3, "shortfall is 2"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			dir := t.TempDir()
			writeDeclaredRowsFile(t, filepath.Join(dir, "_prd.md"), "# PRD\n\n## Unreachable Acceptance\n\n- criterion: Success Metric 10, after release\n  criterion: Success Metric 11, after release\n  reason: later evidence\n  satisfied-by: operator\n")
			path := filepath.Join(dir, "report.md")
			writeDeclaredRowsFile(t, path, fmt.Sprintf("---\nverdict: partial\nrows_blocked_declared: %d\n---\n\n## Results\n\n| # | Status |\n| - | --- |\n| Q10 | blocked (declared: Success Metric 10) |\n| %s | blocked (declared: %s) |\n", tt.count, tt.id, tt.criterion))
			report, err := ReadQAReportFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if len(report.DeclaredRows) != 2 || report.DeclaredRows[1] != (QADeclaredRow{tt.id, tt.criterion}) {
				t.Fatalf("declared rows = %+v", report.DeclaredRows)
			}
			err = QAReportEligibility(dir, report)
			if tt.refusal == "" {
				if err != nil {
					t.Fatal(err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.refusal) {
				t.Fatalf("eligibility = %v; want %q", err, tt.refusal)
			}
			if tt.id == "Q12" && !strings.Contains(err.Error(), "shortfall is 1") {
				t.Fatalf("missing legacy shortfall: %v", err)
			}
		})
	}
}

func TestQAReportEligibilityKeepsTodaysCountForAParaphrasedRow(t *testing.T) {
	t.Parallel()
	for _, table := range []bool{true, false} {
		t.Run(fmt.Sprintf("table=%v", table), func(t *testing.T) {
			dir := t.TempDir()
			writeDeclaredRowsFile(t, filepath.Join(dir, "_prd.md"), "# PRD\n\n## Unreachable Acceptance\n\n- criterion: no Run delivered after this Spec fails the test\n  reason: future Run\n  satisfied-by: later evidence\n")
			content := "---\nverdict: partial\nrows_blocked_declared: 1\n---\n"
			if table {
				content += "\n## Results\n\n| # | Status |\n| - | --- |\n| Q09 | blocked (declared: no later Run fails the test) |\n"
			}
			path := filepath.Join(dir, "report.md")
			writeDeclaredRowsFile(t, path, content)
			report, err := ReadQAReportFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if err := QAReportEligibility(dir, report); err != nil {
				t.Fatal(err)
			}
		})
	}
}

func writeDeclaredRowsFile(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
}
