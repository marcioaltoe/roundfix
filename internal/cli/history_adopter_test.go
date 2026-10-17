package cli

// Suite: adopter history cleanup through public commands.
// Invariant: lineage, per-path refusals and live repairs share one reviewable plan.
// Boundary IN: adopted temporary repositories and command arguments.
// Boundary OUT: text, JSON, Plan Digest, applied bytes and spec-check findings.

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
	"roundfix/internal/speccheck"
)

func historyAdopterFixture(t *testing.T, baseline bool) (historyFixture, string) {
	t.Helper()
	var h historyFixture
	if baseline {
		h = historyFixture{repo: newBaselineUpdateRepository(t), home: t.TempDir()}
		h.env = commandEnvironmentForTest(t)
		h.env.homeDir, h.env.workDir = h.home, h.repo
		h.env.environ = withEnvValue(h.env.environ, "HOME", h.home)
	} else {
		h = newHistoryFixture(t, false)
		historyWrite(t, filepath.Join(h.repo, "docs/agents/setup-context.json"), `{"managedArtifacts":[{"id":"guide.docs-layout","path":"docs/agents/layout.md"}]}`)
		historyWrite(t, filepath.Join(h.repo, "docs/agents/layout.md"), "<!-- setup-context-driven:begin id=guide.docs-layout version=0.0.1 -->\n"+spec.RetirementRuleSentinel+"\n<!-- setup-context-driven:end id=guide.docs-layout -->\n")
	}
	historyWrite(t, filepath.Join(h.repo, "docs/specs/demo/_prd.md"), "# Demo\n")
	target := "docs/specs/old/_techspec.md"
	historyWrite(t, filepath.Join(h.repo, target), "# Original\n")
	original := historyCommit(t, h.repo, "docs: original techspec")
	gittest.Run(t, h.repo, "rm", "-r", "docs/specs/old")
	for path, content := range map[string]string{
		"docs/adr/0001-old.md":          "---\nstatus: superseded\nsuperseded_by: ADR-0002\n---\n# Old\n",
		"docs/adr/0002-new.md":          "---\nstatus: accepted\n---\n# New\n",
		"docs/adr/0003-rejected.md":     "---\nstatus: rejected\n---\n# Rejected\n",
		"docs/history/adr/0004-kept.md": "---\nstatus: superseded\n---\n# Kept\n",
		"docs/references/live.md":       "# Live\n\n[Current behavior](../specs/old/_techspec.md#current-behavior)\n",
		"docs/user-guide/live.md":       "# Guide\n\n[Current behavior](../specs/old/_techspec.md#current-behavior)\n",
	} {
		historyWrite(t, filepath.Join(h.repo, path), content)
	}
	h.revision = historyCommit(t, h.repo, "docs: adopter shapes")
	return h, original
}

func assertAdopterPlan(t *testing.T, report baselineUpdateHistory) {
	t.Helper()
	if len(report.Refused) != 1 || report.Refused[0].Unit != "docs/history/adr/0004-kept.md" || !strings.HasPrefix(report.Refused[0].Reason, "superseded ADR has no accepted successor;") {
		t.Fatalf("refused=%+v", report.Refused)
	}
	var live, lineage bool
	for _, unit := range report.Units {
		if unit.Unit == "live-links" {
			live = unit.Action == "rewrite" && unit.Files == 2 && unit.Links == 2 && strings.Join(unit.Paths, ",") == "docs/references/live.md,docs/user-guide/live.md"
		}
		if unit.Unit == "adr" && strings.Join(unit.Paths, ",") != "docs/adr/0001-old.md,docs/adr/0003-rejected.md" {
			t.Fatalf("ADR paths include successor or refusal: %+v", unit)
		}
	}
	for _, citation := range report.Citations {
		if citation.Path == "docs/adr/0002-new.md" {
			lineage = citation.Line == 6 && citation.Target == "docs/adr/0001-old.md" && strings.HasPrefix(citation.Rewrite, "Supersedes ADR-0001 (") && strings.Contains(citation.Rewrite, "; Git keeps its text at `docs/adr/0001-old.md@")
		}
	}
	if !live || !lineage {
		t.Fatalf("live=%v lineage=%v report=%+v", live, lineage, report)
	}
}

func TestHistorySanitizeAppliesTheAdopterShapes(t *testing.T) {
	t.Parallel()
	if !historyRemovalTestHome(t) {
		return
	}
	h, original := historyAdopterFixture(t, false)
	historyTag(t, h)
	before, refs, index := historySnapshot(t, h.repo), baselineHistoryGit(t, h, "for-each-ref"), baselineHistoryIndex(t, h)
	code, out, stderr := historyRun(t, h, "sanitize")
	if code != 0 || stderr != "" || !strings.Contains(out, "live-links: rewrites 2 link(s) in 2 file(s)\n") || !strings.Contains(out, "cites docs/adr/0002-new.md:6 names docs/adr/0001-old.md; rewrites to Supersedes ADR-0001 (") || !strings.Contains(out, "refused docs/history/adr/0004-kept.md: superseded ADR has no accepted successor;") {
		t.Fatalf("plan=%d %s %s", code, out, stderr)
	}
	code, out, stderr = historyRun(t, h, "sanitize", "--format", "json")
	var report baselineUpdateHistory
	if err := json.Unmarshal([]byte(out), &report); err != nil {
		t.Fatal(err)
	}
	if code != 0 || stderr != "" {
		t.Fatalf("JSON=%d %s %s", code, out, stderr)
	}
	assertAdopterPlan(t, report)
	baselineHistoryUnchanged(t, h, before, refs, index)
	code, out, stderr = historyRun(t, h, "sanitize", "--apply", "--batch", "10")
	if code != 0 || stderr != "" || !strings.Contains(out, "applied 2 unit(s)") {
		t.Fatalf("apply=%d %s %s", code, out, stderr)
	}
	for _, path := range []string{"docs/adr/0001-old.md", "docs/adr/0003-rejected.md"} {
		assertPathMissing(t, filepath.Join(h.repo, path))
	}
	next, err := os.ReadFile(filepath.Join(h.repo, "docs/adr/0002-new.md"))
	if err != nil || strings.Count(string(next), "Supersedes ADR-0001 (") != 1 {
		t.Fatalf("lineage=%s err=%v", next, err)
	}
	for _, path := range []string{"docs/references/live.md", "docs/user-guide/live.md"} {
		next, err := os.ReadFile(filepath.Join(h.repo, path))
		if err != nil || string(next) != strings.Replace(string(before[path]), "[Current behavior](../specs/old/_techspec.md#current-behavior)", "Current behavior (`docs/specs/old/_techspec.md@"+original+"`)", 1) {
			t.Fatalf("live page=%s err=%v", next, err)
		}
	}
	checked, err := speccheck.Check(filepath.Join(h.repo, "docs/specs"), h.repo, "demo")
	if err != nil {
		t.Fatal(err)
	}
	retained := 0
	for _, finding := range checked.Findings {
		if finding.Code == speccheck.CodeHistoryRetained {
			retained++
			if len(finding.Where) != 1 || finding.Where[0].Path != "docs/history/adr/0004-kept.md" {
				t.Fatalf("retained=%+v", finding)
			}
		}
	}
	if retained != 1 {
		t.Fatalf("history findings=%d: %+v", retained, checked)
	}
	code, out, stderr = historyRun(t, h, "sanitize", "--format", "json")
	if err := json.Unmarshal([]byte(out), &report); err != nil {
		t.Fatal(err)
	}
	if code != 0 || stderr != "" || len(report.Units) != 0 || len(report.Citations) != 0 || len(report.Refused) != 1 {
		t.Fatalf("second=%d %s %s", code, out, stderr)
	}
	// A clean tree with just a per-ADR refusal still has a successful apply.
	historyCommit(t, h.repo, "docs: applied cleanup")
	code, out, stderr = historyRun(t, h, "sanitize", "--apply", "--batch", "10", "--format", "json")
	if code != 0 || stderr != "" {
		t.Fatalf("refusal-only apply=%d %s %s", code, out, stderr)
	}
}

func TestBaselineUpdatePlansLineageAndLiveLinks(t *testing.T) {
	t.Parallel()
	if !historyRemovalTestHome(t) {
		return
	}
	h, _ := historyAdopterFixture(t, true)
	before, refs, index := historySnapshot(t, h.repo), baselineHistoryGit(t, h, "for-each-ref"), baselineHistoryIndex(t, h)
	preview, out, stderr, code := baselineHistoryRun(t, h, "json")
	if code != 3 || stderr != "" {
		t.Fatalf("plan=%d %s %s", code, out, stderr)
	}
	assertAdopterPlan(t, preview.History)
	_, text, stderr, code := baselineHistoryRun(t, h, "text")
	if code != 3 || stderr != "" || !strings.Contains(text, "- live-links: rewrites 2 link(s) in 2 file(s)\n  - docs/references/live.md\n  - docs/user-guide/live.md\n") || !strings.Contains(text, "- refused docs/history/adr/0004-kept.md:") || !strings.Contains(text, "Git keeps its text at") {
		t.Fatalf("text=%d %s %s", code, text, stderr)
	}
	baselineHistoryUnchanged(t, h, before, refs, index)
	plan := planBaselineUpdateHistory(t.Context(), h.repo, h.env)
	digest, err := plan.digest("same-baseline-digest")
	if err != nil {
		t.Fatal(err)
	}
	for i, c := range plan.Report.Citations {
		if c.Path == "docs/adr/0002-new.md" {
			plan.Report.Citations[i].Rewrite = strings.Replace(c.Rewrite, "Supersedes", "Supersedes changed", 1)
		}
	}
	changed, err := plan.digest("same-baseline-digest")
	if err != nil || changed == digest {
		t.Fatalf("digest did not bind exact lineage: %q %q %v", digest, changed, err)
	}
	applied, out, stderr, code := baselineHistoryRun(t, h, "json", "--confirm-plan", preview.PlanDigest)
	if code != 0 || stderr != "" || applied.History.Status != "applied" || applied.ApprovedPlanDigest != preview.PlanDigest {
		t.Fatalf("apply=%d %s %s", code, out, stderr)
	}
	after, out, stderr, code := baselineHistoryRun(t, h, "json")
	if code != 0 || stderr != "" || len(after.History.Units) != 0 || len(after.History.Refused) != 1 {
		t.Fatalf("second=%d %s %s", code, out, stderr)
	}
	t.Run("dirty successor refuses ADR unit", func(t *testing.T) {
		dirty, _ := historyAdopterFixture(t, true)
		historyWrite(t, filepath.Join(dirty.repo, "docs/adr/0002-new.md"), "---\nstatus: accepted\n---\n# Uncommitted successor\n")
		planned := planBaselineUpdateHistory(t.Context(), dirty.repo, dirty.env)
		var refused bool
		for _, r := range planned.Report.Refused {
			if r.Unit == "adr" && strings.Contains(r.Reason, "uncommitted changes") {
				refused = true
			}
		}
		if !refused {
			t.Fatalf("dirty successor accepted: %+v", planned.Report)
		}
	})
}
