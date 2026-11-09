package speccheck_test

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strconv"
	"strings"
	"testing"

	"roundfix/internal/speccheck"
)

const unreachableSectionHeader = "# PRD\n\n## Unreachable Acceptance\n\n"

func TestAnUnreachableSectionReadAsWrittenReportsNothing(t *testing.T) {
	t.Parallel()

	for _, tt := range []struct{ name, content string }{
		{"fielded", unreachableSectionHeader + "- criterion: release\n  reason: network\n  satisfied-by: operator\n- criterion: first\n  criterion: second\n  wrapped criterion\n  reason: credentials\n  satisfied-by: maintainer\n"},
		{"none", unreachableSectionHeader + "None. All criteria are reachable\nby the local gate.\n"},
		{"empty", unreachableSectionHeader},
		{"absent", "# PRD\n"},
		{"template comment", unreachableSectionHeader + `<!-- Optional. One Unreachable Acceptance declaration may name several criteria, one criterion line each, and each criterion becomes its own QA row; acceptance matches every declared-blocked row to a declaration that names its criterion, by its text or its leading label.
     A report whose declared-blocked rows do not exceed its declarations is accepted as before.
     Declare unreachable criteria with these fields:
     - criterion: <criterion>
       criterion: <another criterion>
       reason: <why it is unreachable>
       satisfied-by: <evidence or follow-up that satisfies it>
     Omit this section when every criterion is reachable. -->
`},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			checkUnreachableFixture(t, tt.content, func(t *testing.T, result speccheck.Result, _ string) {
				if got := findingsWithCode(result, speccheck.CodeUnreachableUnread); len(got) != 0 {
					t.Fatalf("unread findings = %+v", got)
				}
				for _, skip := range result.Skipped {
					if skip.Code == speccheck.CodeUnreachableUnread {
						t.Fatalf("present PRD skipped: %+v", skip)
					}
				}
			})
		})
	}
}

func TestANumberedUnreachableEntryIsReportedAtItsFirstLine(t *testing.T) {
	t.Parallel()

	const first = "1. A real rerun of a failed release job on GitHub Actions against the npm"
	checkUnreachableFixture(t, unreachableSectionHeader+first+"\n   registry needs credentials\n   and an operator.\n", func(t *testing.T, result speccheck.Result, display string) {
		got := findingsWithCode(result, speccheck.CodeUnreachableUnread)
		wantSummary := fmt.Sprintf("%s line 5 under ## Unreachable Acceptance is not part of a fielded declaration, so the QA gate does not read it: %s", display, strconv.Quote(first))
		if len(got) != 1 {
			t.Fatalf("findings = %+v, want one", got)
		}
		if got[0].Severity != speccheck.SeverityError || got[0].Summary != wantSummary || !reflect.DeepEqual(got[0].Where, []speccheck.Location{{Path: display, Line: 5}}) {
			t.Errorf("finding = %+v, want error at line 5 with summary %q", got[0], wantSummary)
		}
		if !strings.Contains(got[0].Fix, `"- criterion: <text the QA row repeats>"`) {
			t.Errorf("fix = %q, want fielded form", got[0].Fix)
		}
		const wantFix = `Write each entry in the fielded form: a top-level "- criterion: <text the QA row repeats>" line, then indented "reason: <why no Verification or QA row can reach it>" and "satisfied-by: <evidence or follow-up that satisfies it>" lines; start the section with None. when no criterion is unreachable.`
		if got[0].Fix != wantFix {
			t.Errorf("fix = %q, want %q", got[0].Fix, wantFix)
		}
	})
}

func TestUnreachableUnreadSkipsOnlyAMissingPRD(t *testing.T) {
	t.Parallel()

	repo := t.TempDir()
	root := filepath.Join(repo, "docs", "specs")
	const slug = "missing-prd"
	if err := os.MkdirAll(filepath.Join(root, slug), 0o755); err != nil {
		t.Fatal(err)
	}
	for _, stage := range []speccheck.Stage{speccheck.StageAll, speccheck.StagePRD, speccheck.StageTechSpec} {
		t.Run("stage="+string(stage), func(t *testing.T) {
			t.Parallel()
			result, err := speccheck.CheckStage(root, repo, slug, stage)
			if err != nil {
				t.Fatal(err)
			}
			if got := findingsWithCode(result, speccheck.CodeUnreachableUnread); len(got) != 0 {
				t.Fatalf("missing PRD findings = %+v", got)
			}
			var skips []speccheck.SkippedDetector
			for _, skip := range result.Skipped {
				if skip.Code == speccheck.CodeUnreachableUnread {
					skips = append(skips, skip)
				}
			}
			if len(skips) != 1 || skips[0].Missing != "docs/specs/"+slug+"/_prd.md" {
				t.Fatalf("skips = %+v, want one naming the missing PRD", skips)
			}
		})
	}
}

func TestUnreadRunsKeepTheirFirstLines(t *testing.T) {
	t.Parallel()

	checkUnreachableFixture(t, unreachableSectionHeader+"first\n second\n\nthird\n<!-- comment -->\nfourth\n", func(t *testing.T, result speccheck.Result, display string) {
		got := findingsWithCode(result, speccheck.CodeUnreachableUnread)
		if len(got) != 3 {
			t.Fatalf("findings = %+v, want three runs", got)
		}
		for index, line := range []int{5, 8, 10} {
			if !reflect.DeepEqual(got[index].Where, []speccheck.Location{{Path: display, Line: line}}) {
				t.Errorf("finding %d = %+v, want line %d", index, got[index], line)
			}
		}
	})
}

func TestProseAfterADeclarationAndAMissingFieldAreReported(t *testing.T) {
	t.Parallel()

	checkUnreachableFixture(t, unreachableSectionHeader+"- criterion: release\n  reason: network\nProse the gate drops\n", func(t *testing.T, result speccheck.Result, display string) {
		got := findingsWithCode(result, speccheck.CodeUnreachableUnread)
		if len(got) != 2 {
			t.Fatalf("findings = %+v, want unread prose then missing field", got)
		}
		if got[0].Where[0] != (speccheck.Location{Path: display, Line: 7}) || !strings.Contains(got[0].Summary, strconv.Quote("Prose the gate drops")) {
			t.Errorf("prose finding = %+v", got[0])
		}
		want := display + " line 5 starts an Unreachable Acceptance declaration without its satisfied-by: field"
		if got[1].Summary != want || got[1].Where[0].Line != 5 || got[1].Severity != speccheck.SeverityError || got[1].Fix != got[0].Fix {
			t.Errorf("missing-field finding = %+v, want %q with same fix", got[1], want)
		}
	})
}

func TestNoneBesideADeclarationAndAPlainBulletAreReported(t *testing.T) {
	t.Parallel()

	t.Run("none beside declaration", func(t *testing.T) {
		t.Parallel()
		checkUnreachableFixture(t, unreachableSectionHeader+"None.\n\n- criterion: release\n  reason: network\n  satisfied-by: operator\n", func(t *testing.T, result speccheck.Result, _ string) {
			got := findingsWithCode(result, speccheck.CodeUnreachableUnread)
			if len(got) != 1 || got[0].Where[0].Line != 5 || !strings.Contains(got[0].Summary, strconv.Quote("None.")) {
				t.Fatalf("findings = %+v, want None. at line 5", got)
			}
		})
	})
	t.Run("plain bullet", func(t *testing.T) {
		t.Parallel()
		checkUnreachableFixture(t, unreachableSectionHeader+"- A real rerun of a failed release job.\n  An indented second line.\n", func(t *testing.T, result speccheck.Result, _ string) {
			got := findingsWithCode(result, speccheck.CodeUnreachableUnread)
			if len(got) != 2 || got[0].Where[0].Line != 6 || got[1].Where[0].Line != 5 || !strings.Contains(got[1].Summary, "without its criterion: field") {
				t.Fatalf("findings = %+v, want unread continuation then missing criterion", got)
			}
		})
	})
}

func checkUnreachableFixture(t *testing.T, content string, check func(*testing.T, speccheck.Result, string)) {
	t.Helper()
	repo := t.TempDir()
	const slug = "unreachable"
	root := filepath.Join(repo, "docs", "specs")
	dir := filepath.Join(root, slug)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "_prd.md"), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	for _, stage := range []speccheck.Stage{speccheck.StageAll, speccheck.StagePRD, speccheck.StageTechSpec} {
		t.Run("stage="+string(stage), func(t *testing.T) {
			var result speccheck.Result
			var err error
			if stage == speccheck.StageAll {
				result, err = speccheck.Check(root, repo, slug)
			} else {
				result, err = speccheck.CheckStage(root, repo, slug, stage)
			}
			if err != nil {
				t.Fatal(err)
			}
			check(t, result, "docs/specs/"+slug+"/_prd.md")
		})
	}
}
