package baseline

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestBaselinePlansTheClaudeInstructionLink(t *testing.T) {
	t.Parallel()
	const rule = "# Local rules\n- Keep this rule.\n"
	const unlinkedWarning = "this checkout holds the symbolic link CLAUDE.md -> AGENTS.md as a regular file; run git config core.symlinks true, then check CLAUDE.md out again"
	const withheldSuffix = "; move each into docs/agents/specific-repository.md, then run roundfix baseline update to replace it with a symbolic link to AGENTS.md"
	tests := []struct {
		name, shape, content, target, warning string
		refresh, preservation, convert        bool
	}{
		{name: "adoption_greenfield_absent", shape: "absent", convert: true},
		{name: "adoption_greenfield_file", shape: "file", content: rule, convert: true},
		{name: "adoption_preservation_file", shape: "file", content: "Keep the repository release name stable.\n", preservation: true, convert: true},
		{name: "adoption_import_line", shape: "file", content: "@AGENTS.md\n", convert: true},
		{name: "adoption_link_elsewhere", shape: "link", target: "docs/claude-notes.md", content: rule, convert: true},
		{name: "adoption_unlinked_checkout_left", shape: "unlinked", warning: unlinkedWarning},
		{name: "refresh_absent_created", shape: "absent", refresh: true, convert: true},
		{name: "refresh_empty_file_converted", shape: "file", refresh: true, convert: true},
		{name: "refresh_import_line_converted", shape: "file", content: "@AGENTS.md\n", refresh: true, convert: true},
		{name: "refresh_managed_regions_only_converted", shape: "file", content: managedRegionFixture("root.core", "Managed guidance."), refresh: true, convert: true},
		{name: "refresh_copy_of_agents_converted", shape: "copy", refresh: true, convert: true},
		{name: "refresh_link_to_clean_file_converted", shape: "link", target: "docs/claude-notes.md", content: "@AGENTS.md\n", refresh: true, convert: true},
		{name: "refresh_file_with_a_rule_withheld", shape: "file", content: rule, refresh: true, warning: "CLAUDE.md holds 1 passage of its own (CLAUDE.md:1-2)" + withheldSuffix},
		{name: "refresh_link_to_file_with_a_rule_withheld", shape: "link", target: "docs/claude-notes.md", content: rule, refresh: true, warning: "CLAUDE.md links to docs/claude-notes.md, which holds 1 passage of its own (docs/claude-notes.md:1-2)" + withheldSuffix},
		{name: "refresh_unlinked_checkout_left", shape: "unlinked", refresh: true, warning: unlinkedWarning},
		{name: "refresh_link_kept", shape: "link", target: ClaudeInstructionLinkTarget, refresh: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var repo string
			var request PlanRequest
			if tt.refresh {
				repo, request = newManagedRefreshPlanFixture(t)
				if err := os.Remove(filepath.Join(repo, ClaudeInstructionLinkPath)); err != nil {
					t.Fatal(err)
				}
			} else {
				repo = newPlanRepository(t)
				request = PlanRequest{Repository: repo, ProfileID: "go-cli-tui", Decisions: planTestDecisions(), Preservation: RootPreservationRequest{Mode: PreservationModeGreenfield}}
			}
			content := tt.content
			if tt.shape == "copy" {
				agents, err := os.ReadFile(filepath.Join(repo, ClaudeInstructionLinkTarget))
				if err != nil {
					t.Fatal(err)
				}
				content = string(agents)
			}
			switch tt.shape {
			case "file", "copy":
				writeInspectionFile(t, repo, ClaudeInstructionLinkPath, content)
				if err := os.Chmod(filepath.Join(repo, ClaudeInstructionLinkPath), 0o640); err != nil {
					t.Fatal(err)
				}
			case "link":
				if tt.target != ClaudeInstructionLinkTarget {
					writeInspectionFile(t, repo, tt.target, content)
				}
				mustInspectionSymlink(t, tt.target, filepath.Join(repo, ClaudeInstructionLinkPath))
			case "unlinked":
				if !tt.refresh {
					writeInspectionFile(t, repo, ClaudeInstructionLinkTarget, "")
				}
				mustInspectionSymlink(t, ClaudeInstructionLinkTarget, filepath.Join(repo, ClaudeInstructionLinkPath))
				commitInspectionRepository(t, repo, "record instruction link")
				if err := os.Remove(filepath.Join(repo, ClaudeInstructionLinkPath)); err != nil {
					t.Fatal(err)
				}
				runPlanGit(t, repo, "-c", "core.symlinks=false", "checkout", "--", ClaudeInstructionLinkPath)
			}
			before := snapshotVisibleTree(t, repo)
			var plan PlanDocument
			if tt.preservation {
				plan = buildRootPreservationPlan(t, repo, "repository-rules", specificRepositoryPath)
			} else {
				outcome, err := BuildPlan(t.Context(), request)
				if err != nil {
					t.Fatal(err)
				}
				if outcome.Plan == nil {
					t.Fatalf("no plan: %+v", outcome.Result)
				}
				plan = *outcome.Plan
			}
			var posts []Postimage
			var entries []ManagedEntry
			var changes []FileChange
			var warnings []Finding
			for _, p := range plan.Postimages {
				if p.Path == ClaudeInstructionLinkPath {
					posts = append(posts, p)
				}
			}
			for _, e := range plan.ManagedEntries {
				if e.Path == ClaudeInstructionLinkPath {
					entries = append(entries, e)
				}
			}
			for _, c := range plan.FileChanges {
				if c.Path == ClaudeInstructionLinkPath {
					changes = append(changes, c)
				}
			}
			for _, w := range plan.Warnings {
				if strings.HasPrefix(w.Code, "baseline.root.claude-link.") {
					warnings = append(warnings, w)
				}
			}
			if tt.convert {
				want := Postimage{Path: ClaudeInstructionLinkPath, Kind: PreimageSymlink, Content: []byte(ClaudeInstructionLinkTarget)}
				if !reflect.DeepEqual(posts, []Postimage{want}) {
					t.Fatalf("link postimages = %+v", posts)
				}
				action := "update"
				if tt.shape == "absent" {
					action = "create"
				}
				preimage := preimagesByPath(plan.Preimages)[ClaudeInstructionLinkPath]
				wantEntry := ManagedEntry{ID: "root.claude-instruction-link", Path: ClaudeInstructionLinkPath, Kind: "link", Action: action, BeforeIdentity: preimage.ContentIdentity}
				if len(entries) != 1 {
					t.Fatalf("link entries = %+v", entries)
				}
				wantEntry.Ordinal = entries[0].Ordinal
				if !reflect.DeepEqual(entries[0], wantEntry) {
					t.Fatalf("entry=%+v want=%+v", entries[0], wantEntry)
				}
				wantChange := FileChange{Path: ClaudeInstructionLinkPath, Action: action, BeforeIdentity: preimage.ContentIdentity, ManagedEntries: []string{"root.claude-instruction-link"}}
				if !reflect.DeepEqual(changes, []FileChange{wantChange}) {
					t.Fatalf("changes = %+v", changes)
				}
			} else if len(posts)+len(entries)+len(changes) != 0 {
				t.Fatalf("unexpected conversion: %+v %+v %+v", posts, entries, changes)
			}
			var wantWarnings []Finding
			if tt.warning != "" {
				code := "baseline.root.claude-link.withheld"
				if tt.shape == "unlinked" {
					code = "baseline.root.claude-link.unlinked-checkout"
				}
				wantWarnings = []Finding{{Code: code, Path: ClaudeInstructionLinkPath, Message: tt.warning}}
			}
			if !reflect.DeepEqual(warnings, wantWarnings) {
				t.Fatalf("warnings=%+v want=%+v", warnings, wantWarnings)
			}
			var claudeBackups []ManagedEntry
			for _, e := range plan.ManagedEntries {
				if strings.HasPrefix(e.Path, "CLAUDE.") && e.Kind == "backup" {
					claudeBackups = append(claudeBackups, e)
				}
			}
			wantBackup := tt.convert && (tt.shape == "file" || tt.shape == "copy")
			if wantBackup {
				identity := planContentIdentity([]byte(content))
				backupPath := "CLAUDE." + strings.TrimPrefix(identity, "sha256:") + ".md"
				if len(claudeBackups) != 1 || claudeBackups[0].ID != "backup:CLAUDE.md" || claudeBackups[0].Path != backupPath || claudeBackups[0].ContentIdentity != identity {
					t.Fatalf("backups=%+v", claudeBackups)
				}
				if p := planPostimage(t, plan, backupPath); p.Kind != PreimageRegular || !bytes.Equal(p.Content, []byte(content)) {
					t.Fatalf("backup postimage=%+v", p)
				}
			} else if len(claudeBackups) != 0 {
				t.Fatalf("unexpected backup=%+v", claudeBackups)
			}
			if _, err := ApplyPlan(t.Context(), repo, plan, plan.PlanDigest); err != nil {
				t.Fatal(err)
			}
			if tt.convert || tt.shape == "link" && tt.target == ClaudeInstructionLinkTarget {
				assertAppliedPostimage(t, repo, Postimage{Path: ClaudeInstructionLinkPath, Kind: PreimageSymlink, Content: []byte(ClaudeInstructionLinkTarget)})
			} else {
				after := snapshotVisibleTree(t, repo)
				if !reflect.DeepEqual(claudePlanTreeEntry(t, before), claudePlanTreeEntry(t, after)) {
					t.Fatalf("CLAUDE.md changed: before=%+v after=%+v", claudePlanTreeEntry(t, before), claudePlanTreeEntry(t, after))
				}
			}
			if tt.shape == "link" && tt.target != ClaudeInstructionLinkTarget {
				target, err := os.ReadFile(filepath.Join(repo, tt.target))
				if err != nil || string(target) != content {
					t.Fatalf("target changed: %q %v", target, err)
				}
			}
			if wantBackup {
				backup, err := os.ReadFile(filepath.Join(repo, claudeBackups[0].Path))
				if err != nil || string(backup) != content {
					t.Fatalf("backup=%q %v", backup, err)
				}
			}
			if tt.convert {
				request.Decisions = plan.Decisions
				request.Preservation = RootPreservationRequest{Mode: PreservationModeManagedRefresh}
				next := buildManagedRefreshPlan(t, request)
				if len(next.FileChanges) != 0 {
					t.Fatalf("next plan is not current: %+v", next.FileChanges)
				}
				if tt.refresh {
					if _, err := ApplyPlan(t.Context(), repo, plan, plan.PlanDigest); err != nil {
						t.Fatalf("exact reapply: %v", err)
					}
				}
			}
		})
	}
	t.Run("stale_claude_md_refuses_apply", func(t *testing.T) {
		t.Parallel()
		repo, request := newManagedRefreshPlanFixture(t)
		if err := os.Remove(filepath.Join(repo, ClaudeInstructionLinkPath)); err != nil {
			t.Fatal(err)
		}
		writeInspectionFile(t, repo, ClaudeInstructionLinkPath, "@AGENTS.md\n")
		plan := buildManagedRefreshPlan(t, request)
		writeInspectionFile(t, repo, ClaudeInstructionLinkPath, "Keep a newly written rule.\n")
		before := snapshotVisibleTree(t, repo)
		_, err := ApplyPlan(t.Context(), repo, plan, plan.PlanDigest)
		var applyErr *ApplyError
		if !errors.As(err, &applyErr) || applyErr.Kind != ApplyErrorStale {
			t.Fatalf("stale apply=%v", err)
		}
		if after := snapshotVisibleTree(t, repo); !reflect.DeepEqual(before, after) {
			t.Fatal("stale apply changed repository")
		}
	})
	t.Run("guide_states_the_rule", func(t *testing.T) {
		t.Parallel()
		content, err := os.ReadFile("assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/agent-instructions.md")
		if err != nil {
			t.Fatal(err)
		}
		const sentence = "The root instruction file is `AGENTS.md`; `CLAUDE.md` is a symbolic link to it, which `roundfix baseline` creates and keeps, never a file of its own."
		if got := strings.Count(strings.Join(strings.Fields(string(content)), " "), sentence); got != 1 {
			t.Fatalf("guide sentence count=%d", got)
		}
	})
}

func TestClaudeInstructionLinkRollbackRestoresTheFile(t *testing.T) {
	t.Parallel()
	for _, shape := range []string{"regular_file", "link_elsewhere"} {
		t.Run(shape, func(t *testing.T) {
			t.Parallel()
			repo, request := newManagedRefreshPlanFixture(t)
			if err := os.Remove(filepath.Join(repo, ClaudeInstructionLinkPath)); err != nil {
				t.Fatal(err)
			}
			if shape == "regular_file" {
				writeInspectionFile(t, repo, ClaudeInstructionLinkPath, "@AGENTS.md\n")
				if err := os.Chmod(filepath.Join(repo, ClaudeInstructionLinkPath), 0o640); err != nil {
					t.Fatal(err)
				}
			} else {
				writeInspectionFile(t, repo, "docs/claude-notes.md", "@AGENTS.md\n")
				mustInspectionSymlink(t, "docs/claude-notes.md", filepath.Join(repo, ClaudeInstructionLinkPath))
			}
			plan := buildManagedRefreshPlan(t, request)
			before := snapshotVisibleTree(t, repo)
			tx := beginTestTransaction(t, repo, plan)
			injected := errors.New("injected failure after instruction-link replacement")
			replaced := false
			tx.phaseHook = func(point transactionFaultPoint) error {
				if point.Phase == transactionPhaseReplaced && point.Path == ClaudeInstructionLinkPath {
					replaced = true
					assertAppliedPostimage(t, repo, Postimage{Path: ClaudeInstructionLinkPath, Kind: PreimageSymlink, Content: []byte(ClaudeInstructionLinkTarget)})
					return injected
				}
				return nil
			}
			if _, err := tx.Apply(context.Background()); !errors.Is(err, injected) {
				t.Fatalf("fault error=%v", err)
			}
			if !replaced {
				t.Fatal("instruction-link replacement hook was not reached")
			}
			if err := tx.Close(); err != nil {
				t.Fatal(err)
			}
			after := snapshotVisibleTree(t, repo)
			if !reflect.DeepEqual(before, after) {
				for _, original := range before {
					found := false
					for _, current := range after {
						if current.Path == original.Path {
							found = true
							if !reflect.DeepEqual(current, original) {
								t.Errorf("rollback changed %s: before=%+v after=%+v", original.Path, original, current)
							}
							break
						}
					}
					if !found {
						t.Errorf("rollback removed %s", original.Path)
					}
				}
				if len(before) != len(after) {
					t.Errorf("rollback tree size: before=%d after=%d", len(before), len(after))
				}
			}
		})
	}
}

func claudePlanTreeEntry(t *testing.T, entries []visibleTreeEntry) visibleTreeEntry {
	t.Helper()
	for _, entry := range entries {
		if entry.Path == ClaudeInstructionLinkPath {
			return entry
		}
	}
	t.Fatal("CLAUDE.md is missing from the fixture snapshot")
	return visibleTreeEntry{}
}
