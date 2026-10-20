package spec

// Suite: ADR lineage identity and supersession dates.
// Invariant: numbered references select only numbered ADRs; dates belong to the retired ADR.
// Boundary IN: isolated temporary Git histories and ADR frontmatter.
// Boundary OUT: public history plans, exact lineage text and per-path refusals.

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/gittest"
)

func historyIdentityCommit(t *testing.T, root, date string) string {
	t.Helper()
	gittest.Run(t, root, "add", "--all")
	ctx, cancel := context.WithCancel(t.Context())
	args := append([]string{"-C", root}, gittest.ConfigArgs()...)
	cmd := exec.CommandContext(ctx, "git", append(args, "commit", "--allow-empty", "-m", "docs: identity fixture")...)
	cmd.Env = append(gittest.IsolatedEnv(), "GIT_AUTHOR_DATE="+date, "GIT_COMMITTER_DATE="+date)
	var output bytes.Buffer
	cmd.Stdout, cmd.Stderr = &output, &output
	waited := false
	t.Cleanup(func() {
		cancel()
		if cmd.Process != nil && !waited {
			_ = cmd.Wait()
		}
	})
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	err := cmd.Wait()
	waited = true
	if err != nil {
		t.Fatalf("commit: %v\n%s", err, output.String())
	}
	return strings.TrimSpace(gittest.Run(t, root, "rev-parse", "HEAD"))
}

func TestHistoryADRLineageNamesADateNamedADRByItsFileName(t *testing.T) {
	t.Parallel()
	if !citationSourcesTestHome(t) {
		return
	}
	const removed = "docs/adr/2026-09-18-old-cluster.md"
	const successor = "docs/adr/2026-10-01-new-cluster.md"
	for _, existing := range []string{"", "\nSupersedes 2026-09-18-old-cluster (2026-10-01).\n", "supersedes: 2026-09-18-old-cluster\n"} {
		t.Run(fmt.Sprintf("existing=%q", existing), func(t *testing.T) {
			root, _ := gitLocatorRepo(t)
			next := "---\nstatus: accepted\n"
			if strings.HasPrefix(existing, "supersedes:") {
				next += existing
			}
			next += "---\n# New\n"
			if strings.HasPrefix(existing, "\n") {
				next += existing
			}
			writeFile(t, filepath.Join(root, removed), "---\nstatus: superseded\nsuperseded_by: 2026-10-01-new-cluster\nupdated_at: 2026-10-01T00:00:00Z\n---\n# Old\n")
			writeFile(t, filepath.Join(root, successor), next)
			revision := historyIdentityCommit(t, root, "2026-10-19T12:00:00Z")
			plan, err := PlanHistoryKind(t.Context(), root, ArchiveKindADR)
			if err != nil || !reflect.DeepEqual(plan.Files, []string{removed}) || len(plan.Refused) != 0 {
				t.Fatalf("plan=%+v err=%v", plan, err)
			}
			if existing != "" {
				if len(plan.Lineage) != 0 {
					t.Fatalf("duplicate lineage: %+v", plan.Lineage)
				}
				return
			}
			want := []HistoryLineage{{Successor: successor, Removed: removed, Line: fmt.Sprintf("Supersedes 2026-09-18-old-cluster (2026-10-01); Git keeps its text at `%s@%s`.", removed, revision)}}
			if !reflect.DeepEqual(plan.Lineage, want) {
				t.Fatalf("lineage=%+v want=%+v", plan.Lineage, want)
			}
		})
	}
}

func TestHistoryADRLineageDatesTheSupersession(t *testing.T) {
	t.Parallel()
	if !citationSourcesTestHome(t) {
		return
	}
	const removed = "docs/adr/0001-old.md"
	const successor = "docs/adr/0002-new.md"
	for _, tc := range []struct {
		name, updated, want string
		untracked           bool
	}{
		{name: "timestamp converted to UTC", updated: "2026-10-05T22:30:00-03:00", want: "2026-10-06"},
		{name: "YAML date", updated: "2026-10-05", want: "2026-10-05"},
		{name: "quoted date", updated: "\"2026-10-05\"", want: "2026-10-05"},
		{name: "YAML timestamp", updated: "2026-10-05 22:30:00", want: "2026-10-05"},
		{name: "missing updated_at", want: "2026-10-03"},
		{name: "invalid updated_at", updated: "not-a-date", want: "2026-10-03"},
		{name: "null updated_at", updated: "null", want: "2026-10-03"},
		{name: "structured updated_at", updated: "[2026-10-05]", want: "2026-10-03"},
		{name: "HEAD fallback", want: "2026-10-19", untracked: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, _ := gitLocatorRepo(t)
			old := "---\nstatus: superseded\nsuperseded_by: ADR-0002\n"
			if tc.updated != "" {
				old += "updated_at: " + tc.updated + "\n"
			}
			old += "---\n# Old\n"
			writeFile(t, filepath.Join(root, successor), "---\nstatus: accepted\n---\n# New\n")
			if !tc.untracked {
				writeFile(t, filepath.Join(root, removed), old)
			}
			revision := historyIdentityCommit(t, root, "2026-10-02T22:30:00-03:00")
			historyIdentityCommit(t, root, "2026-10-19T12:00:00Z")
			if tc.untracked {
				writeFile(t, filepath.Join(root, removed), old)
			}
			plan, err := PlanHistoryKind(t.Context(), root, ArchiveKindADR)
			if err != nil || !reflect.DeepEqual(plan.Files, []string{removed}) || len(plan.Lineage) != 1 || len(plan.Refused) != 0 {
				t.Fatalf("plan=%+v err=%v", plan, err)
			}
			locator := fmt.Sprintf("`%s@%s`", removed, revision)
			if tc.untracked {
				locator = "`" + removed + "` (not reachable from any branch or tag)"
			}
			want := fmt.Sprintf("Supersedes ADR-0001 (%s); Git keeps its text at %s.", tc.want, locator)
			if plan.Lineage[0].Line != want {
				t.Fatalf("line=%q want=%q", plan.Lineage[0].Line, want)
			}
		})
	}
}

func TestHistoryADRNumericReferencesSelectOnlyNumberedADRs(t *testing.T) {
	t.Parallel()
	if !citationSourcesTestHome(t) {
		return
	}
	const removed = "docs/adr/0001-old.md"
	for _, tc := range []struct {
		name, successor, reference string
		accepted                   bool
	}{
		{"year", "2026-10-01-new-cluster.md", "2026", false},
		{"prefixed year", "2026-10-01-new-cluster.md", "ADR-2026", false},
		{"numbered", "0064-new.md", "ADR-0064", true},
		{"date stem", "2026-10-01-new-cluster.md", "2026-10-01-new-cluster", true},
		{"date filename", "2026-10-01-new-cluster.md", "2026-10-01-new-cluster.md", true},
		{"date path", "2026-10-01-new-cluster.md", "../adr/2026-10-01-new-cluster.md", true},
		{"named stem", "new-cluster.md", "new-cluster", true},
		{"named filename", "new-cluster.md", "new-cluster.md", true},
		{"named path", "new-cluster.md", "docs/adr/new-cluster.md", true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			root, _ := gitLocatorRepo(t)
			writeFile(t, filepath.Join(root, removed), "---\nstatus: superseded\nsuperseded_by: "+tc.reference+"\n---\n# Old\n")
			writeFile(t, filepath.Join(root, "docs/adr", tc.successor), "---\nstatus: accepted\n---\n# New\n\nSupersedes ADR-0001.\n")
			legacyCommit(t, root, "docs: successor identity")
			plan, err := PlanHistoryKind(t.Context(), root, ArchiveKindADR)
			if err != nil {
				t.Fatal(err)
			}
			if tc.accepted {
				if !reflect.DeepEqual(plan.Files, []string{removed}) || len(plan.Refused) != 0 || len(plan.Lineage) != 0 {
					t.Fatalf("plan=%+v", plan)
				}
			} else {
				want := []HistoryRefusal{{Path: removed, Reason: historyADRRefusal}}
				if len(plan.Files) != 0 || !reflect.DeepEqual(plan.Refused, want) || len(plan.Lineage) != 0 {
					t.Fatalf("plan=%+v want refusal=%+v", plan, want)
				}
			}
		})
	}
}
