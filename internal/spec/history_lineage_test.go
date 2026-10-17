package spec

// Suite: superseded ADR retirement.
// Invariant: accepted successors retain exact dated lineage; refusals are per path.
// Boundary IN: temporary Git histories and ADR frontmatter.
// Boundary OUT: public history plans and applied Markdown bytes.

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

func TestHistoryADRPlanWritesTheLineageASuccessorLacks(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	for _, ending := range []string{"", "\n", "\n\n"} {
		t.Run(strings.ReplaceAll(ending, "\n", "newline"), func(t *testing.T) {
			root, _ := gitLocatorRepo(t)
			successor := "docs/adr/0003-new.md"
			original := "---\nstatus: accepted\n---\n# New\n\nKeep every byte.  " + ending
			writeFile(t, filepath.Join(root, successor), original)
			removed := []string{"docs/adr/0001-old.md", "docs/history/adr/0002-old.md"}
			for _, path := range removed {
				writeFile(t, filepath.Join(root, path), "---\nstatus: superseded\nsuperseded_by: ADR-0003\n---\n# Old\n")
			}
			gittest.PersistIdentity(t, root)
			gittest.Run(t, root, "add", "--all")
			cmd := exec.CommandContext(t.Context(), "git", "-C", root, "commit", "-m", "docs: lineage fixture")
			cmd.Env = append(gittest.IsolatedEnv(), "GIT_COMMITTER_DATE=2026-10-17T23:30:00-03:00")
			if out, err := cmd.CombinedOutput(); err != nil {
				t.Fatalf("commit: %v %s", err, out)
			}
			revision := strings.TrimSpace(gittest.Run(t, root, "rev-parse", "HEAD"))
			plan, err := PlanHistoryKind(t.Context(), root, ArchiveKindADR)
			if err != nil || !reflect.DeepEqual(plan.Files, removed) || len(plan.Lineage) != 2 || len(plan.Refused) != 0 {
				t.Fatalf("plan=%+v err=%v", plan, err)
			}
			want := original
			if !strings.HasSuffix(want, "\n") {
				want += "\n"
			}
			for i, path := range removed {
				line := fmt.Sprintf("Supersedes ADR-%04d (2026-10-18); Git keeps its text at `%s@%s`.", i+1, path, revision)
				if got := plan.Lineage[i]; got.Successor != successor || got.Removed != path || got.Line != line {
					t.Fatalf("lineage=%+v want=%q", got, line)
				}
				want += "\n" + line + "\n"
			}
			citations, err := plan.LineageCitations(root)
			if err != nil || len(citations) != 2 {
				t.Fatalf("citations=%+v err=%v", citations, err)
			}
			for i, c := range citations {
				if strings.Split(want, "\n")[c.Line-1] != plan.Lineage[i].Line || c.Target != removed[i] || c.Path != successor || c.Rewrite != plan.Lineage[i].Line {
					t.Fatalf("citation=%+v", c)
				}
			}
			before, err := os.ReadFile(filepath.Join(root, successor))
			if err != nil || string(before) != original {
				t.Fatalf("plan wrote successor: %q %v", before, err)
			}
			if err := ApplyHistoryKind(root, plan); err != nil {
				t.Fatal(err)
			}
			after, err := os.ReadFile(filepath.Join(root, successor))
			if err != nil || string(after) != want {
				t.Fatalf("successor=%q want=%q err=%v", after, want, err)
			}
			again, err := PlanHistoryKind(t.Context(), root, ArchiveKindADR)
			if err != nil || len(again.Files) != 0 || len(again.Lineage) != 0 {
				t.Fatalf("second plan=%+v err=%v", again, err)
			}
			if err := ApplyHistoryKind(root, again); err != nil {
				t.Fatal(err)
			}
			afterAgain, err := os.ReadFile(filepath.Join(root, successor))
			if err != nil || string(afterAgain) != want {
				t.Fatalf("second apply=%q err=%v", afterAgain, err)
			}
		})
	}
}

func TestHistoryADRPlanRefusesOnlyAnADRWithoutAnAcceptedSuccessor(t *testing.T) {
	t.Parallel()
	if !gitLocatorTestHome(t) {
		return
	}
	root, _ := gitLocatorRepo(t)
	for path, content := range map[string]string{
		"docs/adr/0001-no-successor.md":       "---\nstatus: superseded\n---\n# Missing\n",
		"docs/adr/0002-proposed-successor.md": "---\nstatus: superseded\nsuperseded_by: ADR-0003\n---\n# Proposed\n",
		"docs/adr/0003-proposed.md":           "---\nstatus: proposed\n---\n# Proposed\n",
		"docs/adr/0004-rejected.md":           "---\nstatus: rejected\n---\n# Rejected\n",
	} {
		writeFile(t, filepath.Join(root, path), content)
	}
	legacyCommit(t, root, "docs: refusal fixture")
	plan, err := PlanHistoryKind(t.Context(), root, ArchiveKindADR)
	reason := "superseded ADR has no accepted successor; name the accepted ADR that replaces it in superseded_by, or set its status to deprecated"
	want := []HistoryRefusal{{Path: "docs/adr/0001-no-successor.md", Reason: reason}, {Path: "docs/adr/0002-proposed-successor.md", Reason: reason}}
	if err != nil || !reflect.DeepEqual(plan.Refused, want) || !reflect.DeepEqual(plan.Files, []string{"docs/adr/0004-rejected.md"}) || len(plan.Lineage) != 0 {
		t.Fatalf("plan=%+v err=%v", plan, err)
	}
	if err := ApplyHistoryKind(root, plan); err != nil {
		t.Fatal(err)
	}
	for _, r := range want {
		if _, err := os.Stat(filepath.Join(root, r.Path)); err != nil {
			t.Fatal(err)
		}
	}
	if _, err := os.Stat(filepath.Join(root, plan.Files[0])); !os.IsNotExist(err) {
		t.Fatalf("retired ADR remains: %v", err)
	}
}
