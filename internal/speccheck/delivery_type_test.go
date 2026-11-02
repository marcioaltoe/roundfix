package speccheck

// Suite: PRD Delivery Type findings through the Spec Check.
// Boundary IN: disposable repositories and PRD frontmatter.
// Boundary OUT: archiving and publication belong to their package tests.
// Invariant: only unsupported declarations add SC-PRD-TYPE.

import (
	"path/filepath"
	"reflect"
	"testing"
)

func TestPRDTypeFindingNamesTheValue(t *testing.T) {
	t.Parallel()
	for _, tc := range []struct {
		name, front string
		invalid     bool
	}{
		{name: "absent"},
		{name: "no frontmatter"},
		{name: "feat", front: "type: feat\n"},
		{name: "fix", front: "type: fix\n"},
		{name: "perf", front: "type: perf\n"},
		{name: "refactor", front: "type: refactor\n"},
		{name: "docs", front: "type: docs\n", invalid: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			prd := "# Example\n\n## Decisions\n\nAutonomous route: no person present; decisions taken from the test fixture.\n"
			if tc.name != "no frontmatter" {
				prd = "---\nspec: " + clarificationFixtureSlug + "\n" + tc.front + "---\n\n" + prd
			}
			root := clarificationRepo(t, prd)
			result, err := Check(filepath.Join(root, "docs/specs"), root, clarificationFixtureSlug)
			if err != nil {
				t.Fatal(err)
			}
			findings := receiptCodeFindings(result, CodePRDType)
			if !tc.invalid {
				if len(findings) != 0 {
					t.Fatalf("findings=%+v", findings)
				}
				return
			}
			want := []Finding{{
				Code: CodePRDType, Severity: SeverityError,
				Summary: "PRD type \"docs\" is not one of feat, fix, perf, refactor",
				Where:   []Location{{Path: clarificationFixturePRD}},
				Fix:     "Set the PRD's type to feat, fix, perf or refactor, or remove it to deliver as feat.",
			}}
			if !reflect.DeepEqual(findings, want) {
				t.Fatalf("findings=%+v, want %+v", findings, want)
			}
			// Removing the unsupported declaration leaves every other finding unchanged.
			writeReceiptFixture(t, root, clarificationFixturePRD, "---\nspec: "+clarificationFixtureSlug+"\n---\n\n# Example\n\n## Decisions\n\nAutonomous route: no person present; decisions taken from the test fixture.\n")
			baseline, err := Check(filepath.Join(root, "docs/specs"), root, clarificationFixtureSlug)
			if err != nil {
				t.Fatal(err)
			}
			var other []Finding
			for _, finding := range result.Findings {
				if finding.Code != CodePRDType {
					other = append(other, finding)
				}
			}
			if !reflect.DeepEqual(other, baseline.Findings) {
				t.Fatalf("other findings changed: got %+v, baseline %+v", other, baseline.Findings)
			}
		})
	}
}
