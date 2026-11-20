package cli

// Suite: the documentation part of Baseline update.
// Invariant: one digest approves exact documentation changes; repeated updates converge.
// Boundary IN: adopted temporary repositories, isolated command environments and arguments.
// Boundary OUT: JSON, text, exits, fileChanges and byte-preserving repository writes.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"
	"testing"

	"roundfix/internal/baseline"
)

type documentationFixture struct {
	root string
	env  commandEnvironment
}

func newDocumentationFixture(t *testing.T, on bool) documentationFixture {
	t.Helper()
	f := documentationFixture{root: newBaselineApplyTestRepository(t), env: commandEnvironmentForTest(t)}
	f.env.homeDir = t.TempDir()
	f.env.homeDirErr = nil
	f.env.workDir = f.root
	f.env.workDirErr = nil
	f.env.environ = withEnvValue(f.env.environ, "HOME", f.env.homeDir)
	decisions := []baseline.DecisionValue{
		{ID: "language.generated", Value: "English"},
		{ID: "verification.gate", Value: "make verify"}, {ID: "verification.incremental", Value: "make verify-incremental"},
		{ID: "branch.prefix", Value: "<type>/"}, {ID: "spec.scaffold", Value: true}, {ID: "domain.layout", Value: "single-context"},
		{ID: "triage.external", Value: false}, {ID: "autonomous.enabled", Value: true}, {ID: "runtime.backend", Value: "codex gpt-5.5 xhigh"},
		{ID: "runtime.design", Value: "claude opus xhigh"}, {ID: "secondbrain.enabled", Value: on}, {ID: "repository.extension.enabled", Value: false},
	}
	dirs, err := f.env.executableDirectories("documentation fixture")
	if err != nil {
		t.Fatal(err)
	}
	outcome, err := baseline.BuildPlan(t.Context(), baseline.PlanRequest{Repository: f.root, ProfileID: "go-cli-tui", Decisions: decisions, Preservation: baseline.RootPreservationRequest{Mode: baseline.PreservationModeGreenfield}, ExecutableDirectories: dirs})
	if err != nil || outcome.Plan == nil {
		t.Fatalf("adoption plan: %+v %v", outcome, err)
	}
	if _, err := baseline.ApplyPlan(t.Context(), f.root, *outcome.Plan, outcome.Plan.PlanDigest); err != nil {
		t.Fatal(err)
	}
	for path, content := range map[string]string{
		"README.md": "# Repository\n", "CONTEXT.md": "# Context\n",
		"docs/adr/0001-rule.md":     "---\nstatus: accepted\nknowledge: true\n---\n# Rule\n",
		"docs/adr/0002-unmarked.md": "---\nstatus: accepted\n---\n# Unmarked\n",
		"docs/adr/0003-tool.md":     "---\nstatus: accepted\nknowledge: false\n---\n# Tool\n",
		"docs/private/notes.md":     "# Private\n",
	} {
		writeBaselinePlanTestFile(t, f.root, path, content)
	}
	if err := os.MkdirAll(filepath.Join(f.root, "docs/specs"), 0755); err != nil {
		t.Fatal(err)
	}
	commitBaselinePlanTestRepository(t, f.root)
	return f
}

func documentationRun(t *testing.T, f documentationFixture, format string, args ...string) (baselineUpdateResult, string, int) {
	t.Helper()
	command := append([]string{"--repo", f.root, "--no-skills", "--no-history", "--format", format}, args...)
	var out, stderr bytes.Buffer
	code := runBaselineUpdateCommand(t.Context(), command, &out, &stderr, f.env)
	if stderr.Len() != 0 {
		t.Fatalf("exit %d: %s\n%s", code, out.String(), stderr.String())
	}
	var result baselineUpdateResult
	if format == "json" {
		if err := json.Unmarshal(out.Bytes(), &result); err != nil {
			t.Fatal(err)
		}
	}
	return result, out.String(), code
}
func documentationRead(t *testing.T, f documentationFixture, path string) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(f.root, path))
	if err != nil {
		t.Fatal(err)
	}
	return b
}
func documentationSeedExport(t *testing.T, f documentationFixture) string {
	t.Helper()
	own := "# Own bytes  \r\nnotes.md\r\n"
	writeBaselinePlanTestFile(t, f.root, baseline.ExportListPath, "docs/\nAGENTS.md\n"+own)
	return own
}
func documentationRequireCurrent(t *testing.T, f documentationFixture) {
	t.Helper()
	r, out, code := documentationRun(t, f, "json")
	if code != exitOK || r.State != "current" || r.Documentation.Status != "current" || len(r.FileChanges) != 0 {
		t.Fatalf("%d %s", code, out)
	}
}

func TestBaselineUpdateDocumentationPart(t *testing.T) {
	t.Parallel()
	t.Run("plan_names_the_list_and_the_losses", func(t *testing.T) {
		f := newDocumentationFixture(t, true)
		documentationSeedExport(t, f)
		before := documentationRead(t, f, baseline.ExportListPath)
		r, out, code := documentationRun(t, f, "json")
		if code != exitUnverified || r.Documentation.Status != "pending" || r.Documentation.ExportList.Action != "update" || r.Documentation.ADRs.Unmarked != 1 || r.Documentation.ADRs.Total != 3 {
			t.Fatalf("%d %s", code, out)
		}
		want := []string{"README.md", "CONTEXT.md", "docs/user-guide/", "docs/references/", "docs/design/", "docs/lessons/", "docs/adr/0001-rule.md"}
		if !reflect.DeepEqual(r.Documentation.ExportList.Standard, want) {
			t.Fatalf("standard: %+v", r.Documentation.ExportList)
		}
		losses := map[string]string{}
		for _, loss := range r.Documentation.ExportList.StopsExporting {
			losses[loss.Path] = loss.Kind
		}
		for path, kind := range map[string]string{"AGENTS.md": "process", "docs/agents/": "process", "docs/specs/": "process", "docs/adr/": "other", "docs/private/": "other"} {
			if losses[path] != kind {
				t.Errorf("%s loss=%s, want %s", path, losses[path], kind)
			}
		}
		if !bytes.Equal(before, documentationRead(t, f, baseline.ExportListPath)) {
			t.Fatal("preview wrote export list")
		}
	})
	t.Run("text_section", func(t *testing.T) {
		f := newDocumentationFixture(t, true)
		documentationSeedExport(t, f)
		writeBaselinePlanTestFile(t, f.root, "docs/adr/0002-unmarked.md", "---\nstatus: accepted\n---\n# Rule\nAmended by ADR-0004.\n")
		r, _, _ := documentationRun(t, f, "json")
		_, out, code := documentationRun(t, f, "text")
		want := "Documentation: pending\nExport list: update .secondbrain-export (7 line(s) in the Standard Export List)\n"
		for _, line := range r.Documentation.ExportList.Added {
			want += "- adds " + line + "\n"
		}
		want += "- removes docs/\n- removes AGENTS.md\n"
		want += "- stops exporting AGENTS.md (process)\n- stops exporting docs/agents/ (process)\n- stops exporting docs/specs/ (process)\n- stops exporting docs/adr/ (not in the standard list)\n- stops exporting docs/private/ (not in the standard list)\n"

		want += "ADRs without a Knowledge Mark: 1 of 3; they are not exported\nADR headers: 1\n- docs/adr/0002-unmarked.md gains amended_by: [ADR-0004]\n"
		if code != exitUnverified || !strings.Contains(out, want+"Plan Digest:") {
			t.Fatalf("%d %s\nwant %s", code, out, want)
		}
		writeBaselinePlanTestFile(t, f.root, baseline.ExportListPath, baseline.ExportListBeginMarker+"\n")
		_, out, _ = documentationRun(t, f, "text")
		refused := "Documentation: pending\nExport list refused: # roundfix:standard-export:end is missing; nothing is written to .secondbrain-export\nADRs without a Knowledge Mark: 1 of 3; they are not exported\nADR headers: 1\n"
		if !strings.Contains(out, refused) {
			t.Fatal(out)
		}
	})
	t.Run("json_member", func(t *testing.T) {
		f := newDocumentationFixture(t, true)
		_, out, _ := documentationRun(t, f, "json")
		if !strings.HasPrefix(out, `{"history":`) || !strings.Contains(out, `,"documentation":{"status":"pending","exportList":{"path":".secondbrain-export","action":"create","standard":[`) {
			t.Fatal(out)
		}
		var raw map[string]json.RawMessage
		if err := json.Unmarshal([]byte(out), &raw); err != nil {
			t.Fatal(err)
		}
		var d map[string]json.RawMessage
		if err := json.Unmarshal(raw["documentation"], &d); err != nil {
			t.Fatal(err)
		}
		var e map[string]json.RawMessage
		if err := json.Unmarshal(d["exportList"], &e); err != nil {
			t.Fatal(err)
		}
		for _, field := range []string{"standard", "added", "removed", "converted", "stopsExporting"} {
			if len(e[field]) == 0 || e[field][0] != '[' {
				t.Fatalf("%s missing array: %s", field, out)
			}
		}
		if string(d["adrHeaders"]) != "[]" || string(d["adrs"]) != `{"total":3,"unmarked":1}` || string(e["refusal"]) != `""` {
			t.Fatal(out)
		}
	})
	t.Run("digest_moves_only_when_pending", func(t *testing.T) {
		f := newDocumentationFixture(t, true)
		a, _, _ := documentationRun(t, f, "json")
		writeBaselinePlanTestFile(t, f.root, "docs/adr/0002-unmarked.md", "---\nstatus: accepted\nknowledge: true\n---\n# Rule\n")
		b, _, _ := documentationRun(t, f, "json")
		if a.PlanDigest == b.PlanDigest {
			t.Fatal("pending list does not bind ADR membership")
		}
		documentationRun(t, f, "json", "--yes")
		current, _, _ := documentationRun(t, f, "json")
		if current.Documentation.Status != "current" {
			t.Fatal(current.Documentation)
		}
		// An advisory missing mark changes neither the guidance nor the current digest.
		writeBaselinePlanTestFile(t, f.root, "docs/adr/0003-tool.md", "---\nstatus: accepted\n---\n# Tool\n")
		unmarked, _, _ := documentationRun(t, f, "json")
		if unmarked.PlanDigest != current.PlanDigest || unmarked.Documentation.ADRs.Unmarked == current.Documentation.ADRs.Unmarked {
			t.Fatal("current advisory moved digest")
		}
		writeBaselinePlanTestFile(t, f.root, baseline.ExportListPath, baseline.ExportListBeginMarker+"\n")
		refused, _, _ := documentationRun(t, f, "json")
		if refused.PlanDigest != current.PlanDigest || refused.Documentation.ExportList.Action != "refused" {
			t.Fatal("current refusal moved digest")
		}
		writeBaselinePlanTestFile(t, f.root, "docs/adr/0002-unmarked.md", "---\nstatus: accepted\nknowledge: true\n---\nAmended by ADR-0042.\n")
		pending, _, _ := documentationRun(t, f, "json")
		if pending.Documentation.Status != "pending" || pending.PlanDigest == refused.PlanDigest {
			t.Fatal("header fill did not bind pending plan")
		}
		writeBaselinePlanTestFile(t, f.root, baseline.ExportListPath, baseline.ExportListEndMarker+"\n")
		changedRefusal, _, _ := documentationRun(t, f, "json")
		if changedRefusal.PlanDigest == pending.PlanDigest {
			t.Fatal("pending plan did not bind refusal")
		}
		writeBaselinePlanTestFile(t, f.root, "docs/adr/0002-unmarked.md", "---\nstatus: accepted\nknowledge: true\n---\nAmended by ADR-0042 and ADR-0043.\n")
		changedHeader, _, _ := documentationRun(t, f, "json")
		if changedHeader.PlanDigest == changedRefusal.PlanDigest {
			t.Fatal("pending plan did not bind header identifiers")
		}

	})
	t.Run("plan_ready_message", func(t *testing.T) {
		f := newDocumentationFixture(t, true)
		r, out, code := documentationRun(t, f, "json")
		if code != exitUnverified || r.State != "plan_ready" || r.Message != "guidance matches the current Baseline catalog; 1 documentation change(s) pending" || !strings.Contains(r.NextAction, r.PlanDigest) {
			t.Fatalf("%d %s", code, out)
		}
	})
	t.Run("file_changes", func(t *testing.T) {
		f := newDocumentationFixture(t, true)
		writeBaselinePlanTestFile(t, f.root, "docs/adr/0002-unmarked.md", "# Legacy\nAmended by ADR-0004.\n")
		r, _, _ := documentationRun(t, f, "json")
		want := map[string]string{baseline.ExportListPath: "export.standard-list", "docs/adr/0002-unmarked.md": "adr.amendment-list"}
		if len(r.FileChanges) != len(want) {
			t.Fatalf("%+v", r.FileChanges)
		}
		for _, change := range r.FileChanges {
			if !slices.Equal(change.ManagedEntries, []string{want[change.Path]}) || change.AfterIdentity == "" {
				t.Fatalf("%+v", change)
			}
		}
	})
	t.Run("confirmed_apply", func(t *testing.T) {
		f := newDocumentationFixture(t, true)
		own := documentationSeedExport(t, f)
		r, _, _ := documentationRun(t, f, "json")
		applied, out, code := documentationRun(t, f, "json", "--confirm-plan", r.PlanDigest)
		data := documentationRead(t, f, baseline.ExportListPath)
		if code != exitOK || applied.Documentation.Status != "applied" || applied.ApprovedPlanDigest != r.PlanDigest || !bytes.HasPrefix(data, []byte(baseline.ExportListBeginMarker+"\n")) || !bytes.HasSuffix(data, []byte("\n\n"+own)) || bytes.Contains(data, []byte("\nAGENTS.md\n")) || bytes.Contains(data, []byte("\ndocs/\n")) {
			t.Fatalf("%d %s\n%q", code, out, data)
		}
	})
	t.Run("second_run_current", func(t *testing.T) {
		f := newDocumentationFixture(t, true)
		documentationRun(t, f, "json", "--yes")
		documentationRequireCurrent(t, f)
	})
	t.Run("secondbrain_off", func(t *testing.T) {
		f := newDocumentationFixture(t, false)
		r, out, code := documentationRun(t, f, "json", "--yes")
		if code != exitOK || r.Documentation.ExportList.Action != "none" {
			t.Fatalf("%d %s", code, out)
		}
		if _, err := os.Stat(filepath.Join(f.root, baseline.ExportListPath)); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("export file exists: %v", err)
		}
		// Even an unreadable-as-a-file entry is ignored when the module is off.
		if err := os.Mkdir(filepath.Join(f.root, baseline.ExportListPath), 0755); err != nil {
			t.Fatal(err)
		}
		r, out, code = documentationRun(t, f, "json")
		if code != exitOK || r.Documentation.ExportList.Action != "none" {
			t.Fatalf("%d %s", code, out)
		}
	})
	t.Run("marker_refused_rest_applies", func(t *testing.T) {
		f := newDocumentationFixture(t, true)
		broken := baseline.ExportListBeginMarker + "\nREADME.md\n"
		writeBaselinePlanTestFile(t, f.root, baseline.ExportListPath, broken)
		writeBaselinePlanTestFile(t, f.root, "docs/adr/0002-unmarked.md", "# Legacy\nAmended by ADR-0004.\n")
		r, out, code := documentationRun(t, f, "json", "--yes")
		if code != exitOK || r.Documentation.Status != "applied" || r.Documentation.ExportList.Action != "refused" || string(documentationRead(t, f, baseline.ExportListPath)) != broken || !bytes.Contains(documentationRead(t, f, "docs/adr/0002-unmarked.md"), []byte("amended_by: [ADR-0004]")) {
			t.Fatalf("%d %s", code, out)
		}
	})
	t.Run("fills_amendment_lists", func(t *testing.T) {
		f := newDocumentationFixture(t, true)
		before := map[string]string{}
		forms := []string{"Amended by", "Narrowed by", "Extended by", "Superseded in part by"}
		for i, form := range forms {
			for j, ids := range []string{"ADR-0042", "[ADR-0043](forty-three.md) and ADR-0042"} {
				path := fmt.Sprintf("docs/adr/%04d-form.md", 10+i*2+j)
				original := "---\r\nstatus: accepted\r\nknowledge: false\r\nupdated_at: 2026-01-01T00:00:00Z # original\r\n---\r\n# Rule\r\n**" + form + " " + ids + ".**\r\n"
				before[path] = original
				writeBaselinePlanTestFile(t, f.root, path, original)
			}
		}
		before["docs/adr/domain/0054-nested.md"] = "---\r\nstatus: accepted\r\nknowledge: false\r\n---\r\n# Nested\r\nAmended by ADR-0042.\r\n"
		before["docs/adr/0050-legacy.md"] = "# Legacy\nAmended by ADR-0042.\n"
		before["docs/adr/0051-code.md"] = "# Code\n```md\nAmended by ADR-0042.\n```\n`Extended by ADR-0043`\n"
		before["docs/adr/0052-self.md"] = "# Self\nAmended by ADR-0052.\n"
		before["docs/adr/0053-retired.md"] = "---\nstatus: rejected\nknowledge: true\n---\nAmended by ADR-0042.\n"
		for path, content := range before {
			writeBaselinePlanTestFile(t, f.root, path, content)
		}
		r, out, code := documentationRun(t, f, "json", "--yes")
		if code != exitOK || len(r.Documentation.ADRHeaders) != 10 || slices.Contains(r.Documentation.ExportList.Standard, "docs/adr/0053-retired.md") {
			t.Fatalf("%d %s", code, out)
		}
		for path, original := range before {
			after := string(documentationRead(t, f, path))
			if strings.Contains(path, "-code") || strings.Contains(path, "-self") || strings.Contains(path, "-retired") {
				if after != original {
					t.Fatalf("changed %s", path)
				}
				continue
			}
			list := "amended_by: [ADR-0042]"
			if strings.Contains(original, "ADR-0043") {
				list = "amended_by: [ADR-0042, ADR-0043]"
			}
			want := ""
			if strings.Contains(path, "legacy") {
				want = "---\n" + list + "\n---\n" + original
			} else {
				want = strings.Replace(original, "knowledge: false\r\n", "knowledge: false\r\n"+list+"\r\n", 1)
			}
			if after != want {
				t.Fatalf("%s\ngot %q\nwant %q", path, after, want)
			}
		}
		documentationRequireCurrent(t, f)
	})
	t.Run("drops_its_own_identifier", func(t *testing.T) {
		f := newDocumentationFixture(t, true)
		original := "---\nstatus: accepted\nknowledge: false\namended_by: [ADR-42, ADR-55]\n---\n# Rule\nAmended by ADR-43.\n"
		want := "---\nstatus: accepted\nknowledge: false\namended_by: [ADR-42, ADR-43]\n---\n# Rule\nAmended by ADR-43.\n"
		writeBaselinePlanTestFile(t, f.root, "docs/adr/0055-lists-itself.md", original)
		_, out, code := documentationRun(t, f, "json", "--yes")
		if code != exitOK || string(documentationRead(t, f, "docs/adr/0055-lists-itself.md")) != want {
			t.Fatalf("%d %s", code, out)
		}
		documentationRequireCurrent(t, f)
	})
	t.Run("fill_beside_history_unit", func(t *testing.T) {
		f := newDocumentationFixture(t, true)
		writeBaselinePlanTestFile(t, f.root, "docs/adr/0004-old.md", "---\nstatus: superseded\nsuperseded_by: ADR-0001\n---\n# Old\n")
		original := "---\nstatus: accepted\nknowledge: true\n---\n# Rule\nAmended by ADR-0042.\n"
		writeBaselinePlanTestFile(t, f.root, "docs/adr/0001-rule.md", original)
		commitBaselinePlanTestRepository(t, f.root)
		var out, stderr bytes.Buffer
		code := runBaselineUpdateCommand(t.Context(), []string{"--repo", f.root, "--no-skills", "--yes", "--format", "json"}, &out, &stderr, f.env)
		after := string(documentationRead(t, f, "docs/adr/0001-rule.md"))
		if code != exitOK || stderr.Len() != 0 || !strings.HasPrefix(after, strings.Replace(original, "knowledge: true\n", "knowledge: true\namended_by: [ADR-0042]\n", 1)) || !strings.Contains(after, "Supersedes ADR-0004 (") || !strings.Contains(after, "Git keeps its text at") {
			t.Fatalf("%d %s %s\n%s", code, out.String(), stderr.String(), after)
		}
		if _, err := os.Stat(filepath.Join(f.root, "docs/adr/0004-old.md")); !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("old ADR remains: %v", err)
		}
	})
	t.Run("no_history_does_not_skip", func(t *testing.T) {
		f := newDocumentationFixture(t, true)
		r, out, code := documentationRun(t, f, "json", "--yes")
		if code != exitOK || r.History.Status != "skipped" || r.Skills.Status != baselineUpdateSkillsSkipped || r.Documentation.Status != "applied" {
			t.Fatalf("%d %s", code, out)
		}
	})
	t.Run("apply_failure", func(t *testing.T) {
		f := newDocumentationFixture(t, true)
		var out, stderr bytes.Buffer
		attempted := ""
		stages := baselineUpdateStages{skills: runBaselineUpdateSkillsStage, documentation: func(ctx context.Context, root string, p baselineDocumentationPlan) (baselineUpdateDocumentation, error) {
			anchored, err := os.OpenRoot(root)
			if err != nil {
				t.Fatal(err)
			}
			defer anchored.Close()
			return p.apply(ctx, root, anchored, func(path string, content []byte, mode os.FileMode) error {
				attempted = path
				return errors.New("documentation output write failed")
			})
		}}

		code := runBaselineUpdateCommandWithStages(t.Context(), []string{"--repo", f.root, "--no-skills", "--no-history", "--yes", "--format", "json"}, &out, &stderr, f.env, stages)
		var r baselineUpdateResult
		if err := json.Unmarshal(out.Bytes(), &r); err != nil {
			t.Fatal(err)
		}
		if attempted != baseline.ExportListPath || code != exitRunFailed || r.Category != "documentation" || r.State != "failed" || r.NextAction != "restore .secondbrain-export and docs/adr with git restore, then rerun roundfix baseline update" || !strings.Contains(stderr.String(), "baseline update failed") {
			t.Fatalf("%d %s %s", code, out.String(), stderr.String())
		}
	})
}
