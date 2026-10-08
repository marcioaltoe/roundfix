package cli

// Suite: Baseline update Pending History.
// Invariant: approval binds selected history; previews and refusals preserve bytes and refs.
// Boundary IN: adopted temporary Git repositories, temporary homes, command arguments.
// Boundary OUT: stdout, JSON, working tree, index, local refs and commit count.

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"roundfix/internal/baseline"
	"roundfix/internal/gittest"
	"roundfix/internal/spec"
)

func baselineHistoryFixture(t *testing.T) historyFixture {
	t.Helper()
	h := historyFixture{repo: newBaselineUpdateRepository(t), home: t.TempDir()}
	if err := os.MkdirAll(filepath.Join(h.repo, "docs/specs"), 0755); err != nil {
		t.Fatal(err)
	}
	h.env = commandEnvironmentForTest(t)
	h.env.homeDir, h.env.homeDirErr, h.env.workDir, h.env.workDirErr = h.home, nil, h.repo, nil
	h.env.environ = withEnvValue(h.env.environ, "HOME", h.home)
	for _, slug := range []string{"0001-maps-unproven", "0002-plain", "0003-broken"} {
		historyLegacyFolder(t, h, slug)
	}
	historyWrite(t, filepath.Join(h.repo, "docs/history/specs/0001-maps-unproven/_prd.md"), "---\nspec: 0001-maps-unproven\ncreated: 2026-09-01\nunproven:\n  - row: 03\n    claim: Retained outcome\n    reason: No standalone evidence\n---\n\n# Legacy\n\nPreserved outcome.\n")
	historyWrite(t, filepath.Join(h.repo, "docs/history/specs/0003-broken/_prd.md"), "---\nspec: [broken\n---\n\n# Broken\n")
	historyWrite(t, filepath.Join(h.repo, "docs/history/findings/entry.md"), "---\nstatus: closed\nreason: done\n---\n\n# Retired\n\nFirst paragraph.\n\n## Detail\n\n"+strings.Repeat("Full detail.\n", 50))
	h.revision = historyCommit(t, h.repo, "Legacy history (#42)")
	return h
}

func baselineHistoryRun(t *testing.T, h historyFixture, format string, args ...string) (baselineUpdateResult, string, string, int) {
	t.Helper()
	var out, stderr bytes.Buffer
	command := append([]string{"--repo", h.repo, "--no-skills", "--format", format}, args...)
	code := runBaselineUpdateCommandWithSkillsStage(context.Background(), command, &out, &stderr, h.env, func(context.Context, baselineUpdateSkillsRequest) (baselineUpdateSkillsResult, error) {
		return baselineUpdateSkillsResult{Status: baselineUpdateSkillsVerified}, nil
	})
	var result baselineUpdateResult
	if format == "json" {
		if err := json.Unmarshal(out.Bytes(), &result); err != nil {
			t.Fatalf("decode %s: %v", out.String(), err)
		}
	}
	return result, out.String(), stderr.String(), code
}

func baselineHistoryDigest(t *testing.T, h historyFixture) string {
	t.Helper()
	catalog, err := baseline.LoadEmbeddedCatalog()
	if err != nil {
		t.Fatal(err)
	}
	input, err := baseline.ResolveManifestInput(h.repo, catalog)
	if err != nil {
		t.Fatal(err)
	}
	dirs, err := h.env.executableDirectories("test")
	if err != nil {
		t.Fatal(err)
	}
	outcome, err := baseline.BuildPlan(context.Background(), baseline.PlanRequest{Repository: h.repo, ProfileID: input.ProfileID, Decisions: input.Decisions, Preservation: baseline.RootPreservationRequest{Mode: baseline.PreservationModeManagedRefresh}, ExecutableDirectories: dirs})
	if err != nil || outcome.Plan == nil {
		t.Fatalf("plan: %+v %v", outcome, err)
	}
	return outcome.Plan.PlanDigest
}

func baselineHistoryGit(t *testing.T, h historyFixture, args ...string) string {
	t.Helper()
	return gittest.Run(t, h.repo, args...)
}
func baselineHistoryIndex(t *testing.T, h historyFixture) string {
	t.Helper()
	path := strings.TrimSpace(baselineHistoryGit(t, h, "rev-parse", "--git-path", "index"))
	if !filepath.IsAbs(path) {
		path = filepath.Join(h.repo, path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func baselineHistoryUnchanged(t *testing.T, h historyFixture, before map[string][]byte, refs, index string) {
	t.Helper()
	if !reflect.DeepEqual(before, historySnapshot(t, h.repo)) || refs != baselineHistoryGit(t, h, "for-each-ref") || index != baselineHistoryIndex(t, h) {
		t.Fatal("tree, index or refs changed")
	}
}

func TestBaselineUpdatePlansPendingHistoryWithoutWrites(t *testing.T) {
	h := baselineHistoryFixture(t)
	stamp := time.Unix(1, 0)
	if err := os.Chtimes(filepath.Join(h.repo, "docs/history/specs/0002-plain/_prd.md"), stamp, stamp); err != nil {
		t.Fatal(err)
	}
	before, refs, index := historySnapshot(t, h.repo), baselineHistoryGit(t, h, "for-each-ref"), baselineHistoryIndex(t, h)
	r, out, stderr, code := baselineHistoryRun(t, h, "json")
	if code != 3 || stderr != "" || r.State != "plan_ready" || len(r.History.Units) != 3 || len(r.History.Refused) != 1 || r.History.Status != "pending" || r.History.Tag.Action != "create" || r.History.Tag.Commit != h.revision || r.PlanDigest == baselineHistoryDigest(t, h) {
		t.Fatalf("%d %s %s", code, out, stderr)
	}
	_, text, stderr, code := baselineHistoryRun(t, h, "text")
	if code != 3 || stderr != "" {
		t.Fatalf("%d %s %s", code, text, stderr)
	}
	var files int
	var total int64
	for _, u := range r.History.Units {
		files += u.Files
		total += u.Bytes
	}
	a, b, c := r.History.Units[0], r.History.Units[1], r.History.Units[2]
	want := fmt.Sprintf("History: pending\nHistory units: 3 (%d file(s), %d bytes leave docs/history)\n- folder %s: removes %d file(s) (%d bytes) and writes %s (no-qa)\n- folder %s: removes %d file(s) (%d bytes) and writes %s (no-qa)\n- findings: reduces 1 file(s) from %d to %d bytes\nHistory refused: 1\n- refused docs/history/specs/0003-broken: %s\nHistory tag: creates annotated history-full at %.12s; push it with git push origin history-full\nPlan Digest: %s\nNext action: review the plan and rerun with --confirm-plan %s, or use --yes to approve the digest computed in one invocation\n", files, total, a.Unit, a.Files, a.Bytes, a.Record, b.Unit, b.Files, b.Bytes, b.Record, c.Bytes, c.BytesAfter, r.History.Refused[0].Reason, h.revision, r.PlanDigest, r.PlanDigest)
	if !strings.HasSuffix(text, want) || !strings.HasPrefix(text, "Baseline update: plan ready\nCategory: approval\nResult: guidance matches the current Baseline catalog; 3 history unit(s) pending sanitize\n") {
		t.Fatalf("transcript mismatch:\n%s\nwant suffix:\n%s", text, want)
	}
	baselineHistoryUnchanged(t, h, before, refs, index)
}

func TestBaselineUpdateAppliesPendingHistoryAndCreatesTheTag(t *testing.T) {
	h := baselineHistoryFixture(t)
	broken := snapshotDirectoryFiles(t, filepath.Join(h.repo, "docs/history/specs/0003-broken"))
	refs, count := baselineHistoryGit(t, h, "for-each-ref"), baselineHistoryGit(t, h, "rev-list", "--count", "HEAD")
	preview, _, _, _ := baselineHistoryRun(t, h, "json")
	r, out, stderr, code := baselineHistoryRun(t, h, "json", "--confirm-plan", preview.PlanDigest)
	if code != 0 || stderr != "" || r.State != "verified" || r.History.Status != "applied" || r.History.Applied.Units != 3 || r.History.Applied.Records != 2 || r.History.Applied.Reduced != 1 || r.ApprovedPlanDigest != preview.PlanDigest {
		t.Fatalf("%d %s %s", code, out, stderr)
	}
	for _, slug := range []string{"0001-maps-unproven", "0002-plain"} {
		assertPathMissing(t, filepath.Join(h.repo, "docs/history/specs", slug))
		data, err := os.ReadFile(filepath.Join(h.repo, "docs/history/specs", slug+".md"))
		if err != nil {
			t.Fatal(err)
		}
		record, err := spec.ParseArchiveRecord(data)
		if err != nil || record.Disposition != spec.ArchiveNoQA {
			t.Fatalf("record: %+v %v", record, err)
		}
		if slug == "0001-maps-unproven" && len(record.Unproven) != 1 {
			t.Fatalf("maps: %+v", record)
		}
	}
	data, err := os.ReadFile(filepath.Join(h.repo, "docs/history/findings/entry.md"))
	if err != nil || !spec.IsReducedHistoryEntry(data) {
		t.Fatalf("finding %s %v", data, err)
	}
	if !reflect.DeepEqual(broken, snapshotDirectoryFiles(t, filepath.Join(h.repo, "docs/history/specs/0003-broken"))) {
		t.Fatal("broken folder changed")
	}
	if strings.TrimSpace(baselineHistoryGit(t, h, "cat-file", "-t", "history-full")) != "tag" || strings.TrimSpace(baselineHistoryGit(t, h, "rev-parse", "history-full^{commit}")) != h.revision || count != baselineHistoryGit(t, h, "rev-list", "--count", "HEAD") {
		t.Fatal("tag or commits differ")
	}
	var otherRefs []string
	for _, line := range strings.Split(baselineHistoryGit(t, h, "for-each-ref"), "\n") {
		if !strings.Contains(line, "refs/tags/history-full") {
			otherRefs = append(otherRefs, line)
		}
	}
	if strings.Join(otherRefs, "\n") != refs {
		t.Fatal("other ref changed")
	}
	r, out, stderr, code = baselineHistoryRun(t, h, "json")
	if code != 0 || stderr != "" || r.State != "current" || r.History.Status != "current" || len(r.History.Refused) != 1 {
		t.Fatalf("second: %d %s %s", code, out, stderr)
	}
	_, text, _, _ := baselineHistoryRun(t, h, "text")
	if !strings.Contains(text, "History: current\nHistory refused: 1\n- refused docs/history/specs/0003-broken:") {
		t.Fatal(text)
	}
}

func TestBaselineUpdateWithoutPendingHistoryKeepsTheBaselineDigest(t *testing.T) {
	h := baselineHistoryFixture(t)
	if err := os.RemoveAll(filepath.Join(h.repo, "docs/history")); err != nil {
		t.Fatal(err)
	}
	historyCommit(t, h.repo, "No history")
	r, out, stderr, code := baselineHistoryRun(t, h, "json")
	if code != 0 || stderr != "" || r.History.Status != "current" || r.History.Units == nil || r.History.Refused == nil || r.PlanDigest != baselineHistoryDigest(t, h) {
		t.Fatalf("%d %s %s", code, out, stderr)
	}
	_, text, _, _ := baselineHistoryRun(t, h, "text")
	if strings.Contains(text, "\nHistory") {
		t.Fatal(text)
	}
	const wrong = "sha256:stale"
	rejected, out, _, code := baselineHistoryRun(t, h, "json", "--confirm-plan", wrong)
	wantMessage := fmt.Sprintf("confirmed Plan Digest %q does not match supplied plan %q", wrong, r.PlanDigest)
	if code != 3 || rejected.State != "action_required" || rejected.Category != "approval" || rejected.Message != wantMessage || rejected.NextAction != "review the supplied plan and rerun with --confirm-plan "+r.PlanDigest {
		t.Fatalf("legacy approval: %d %s", code, out)
	}

}

func TestBaselineUpdateNoHistorySkipsTheSection(t *testing.T) {
	h := baselineHistoryFixture(t)
	before, refs := historySnapshot(t, h.repo), baselineHistoryGit(t, h, "for-each-ref")
	for _, args := range [][]string{{"--no-history"}, {"--no-history", "--yes"}} {
		r, out, stderr, code := baselineHistoryRun(t, h, "json", args...)
		if code != 0 || stderr != "" || r.History.Status != "skipped" || len(r.History.Units) != 0 {
			t.Fatalf("%d %s %s", code, out, stderr)
		}
		_, text, _, _ := baselineHistoryRun(t, h, "text", args...)
		if strings.Contains(text, "\nHistory") {
			t.Fatal(text)
		}
	}
	if !reflect.DeepEqual(before, historySnapshot(t, h.repo)) || refs != baselineHistoryGit(t, h, "for-each-ref") {
		t.Fatal("opt out changed history or refs")
	}
	_, text, _, code := baselineHistoryRun(t, h, "text", "--help")
	if code != 0 || !strings.Contains(text, "--no-history") {
		t.Fatal(text)
	}
}

func TestBaselineUpdateRefusesUnitsWithUncommittedChanges(t *testing.T) {
	h := baselineHistoryFixture(t)
	for _, p := range []string{"docs/history/specs/0002-plain/_prd.md", "docs/history/specs/0001-maps-unproven/untracked.md"} {
		f := filepath.Join(h.repo, p)
		if strings.Contains(p, "0002") {
			data, err := os.ReadFile(f)
			if err != nil {
				t.Fatal(err)
			}
			historyWrite(t, f, string(data)+"\nDirty.\n")
		} else {
			historyWrite(t, f, "Untracked.\n")
		}
	}
	before := historySnapshot(t, h.repo)
	r, out, stderr, code := baselineHistoryRun(t, h, "json")
	if code != 3 || stderr != "" || len(r.History.Units) != 1 || len(r.History.Refused) != 3 {
		t.Fatalf("%d %s %s", code, out, stderr)
	}
	for _, ref := range r.History.Refused[:2] {
		if ref.Reason != "has uncommitted changes under "+ref.Unit+"; commit or restore them and rerun" {
			t.Fatalf("refusal: %+v", ref)
		}
	}
	_, out, stderr, code = baselineHistoryRun(t, h, "json", "--yes")
	if code != 0 {
		t.Fatalf("%d %s %s", code, out, stderr)
	}
	after := historySnapshot(t, h.repo)
	for p, data := range before {
		if strings.HasPrefix(p, "docs/history/specs/") && !bytes.Equal(data, after[p]) {
			t.Fatalf("refused bytes changed: %s", p)
		}
	}
}

func TestBaselineUpdateHistoryHonorsAnExistingTag(t *testing.T) {
	t.Run("annotated", func(t *testing.T) {
		h := baselineHistoryFixture(t)
		historyTag(t, h)
		tag := baselineHistoryGit(t, h, "rev-parse", "history-full")
		historyLegacyFolder(t, h, "0004-later")
		h.revision = historyCommit(t, h.repo, "Later history")
		r, out, stderr, code := baselineHistoryRun(t, h, "json")
		if code != 3 || stderr != "" || r.History.Tag.Action != "present" || len(r.History.Units) != 3 || len(r.History.Refused) != 2 || r.History.Refused[1].Reason != "history-full does not hold docs/history/specs/0004-later/_prd.md" {
			t.Fatalf("%d %s %s", code, out, stderr)
		}
		before := snapshotDirectoryFiles(t, filepath.Join(h.repo, "docs/history/specs/0004-later"))
		_, out, stderr, code = baselineHistoryRun(t, h, "json", "--yes")
		if code != 0 || tag != baselineHistoryGit(t, h, "rev-parse", "history-full") || !reflect.DeepEqual(before, snapshotDirectoryFiles(t, filepath.Join(h.repo, "docs/history/specs/0004-later"))) {
			t.Fatalf("%d %s %s", code, out, stderr)
		}
	})
	t.Run("lightweight", func(t *testing.T) {
		h := baselineHistoryFixture(t)
		baselineHistoryGit(t, h, "tag", "history-full")
		r, out, stderr, code := baselineHistoryRun(t, h, "json")
		if code != 0 || stderr != "" || r.State != "current" || r.History.Status != "blocked" || len(r.History.Units) != 0 || r.PlanDigest != baselineHistoryDigest(t, h) {
			t.Fatalf("%d %s %s", code, out, stderr)
		}
	})
}

func TestBaselineUpdateRejectsAStaleHistoryDigest(t *testing.T) {
	h := baselineHistoryFixture(t)
	before, refs, index := historySnapshot(t, h.repo), baselineHistoryGit(t, h, "for-each-ref"), baselineHistoryIndex(t, h)
	r, out, _, code := baselineHistoryRun(t, h, "json", "--confirm-plan", baselineHistoryDigest(t, h))
	if code != 3 || r.State != "action_required" || r.Category != "approval" || r.PlanDigest == baselineHistoryDigest(t, h) {
		t.Fatalf("%d %s", code, out)
	}
	baselineHistoryUnchanged(t, h, before, refs, index)
}

func TestBaselineUpdateReportsATagFailureBeforeAnyConversion(t *testing.T) {
	h := baselineHistoryFixture(t)
	before, refs := historySnapshot(t, h.repo), baselineHistoryGit(t, h, "for-each-ref")
	baselineHistoryGit(t, h, "config", "tag.gpgSign", "true")
	baselineHistoryGit(t, h, "config", "gpg.program", "false")
	r, out, stderr, code := baselineHistoryRun(t, h, "json", "--yes")
	if code != 1 || r.State != "failed" || r.Category != "history" || !strings.Contains(r.Message, "history-full") || r.CurrentCatalog.Digest == "" || r.StatusMatrix == nil || r.ApprovedPlanDigest == "" {
		t.Fatalf("%d %s %s", code, out, stderr)
	}
	if !reflect.DeepEqual(before, historySnapshot(t, h.repo)) || refs != baselineHistoryGit(t, h, "for-each-ref") {
		t.Fatal("conversion or tag happened after tag failure")
	}
}

func TestBaselineUpdateHistoryPlanningBoundaries(t *testing.T) {
	t.Run("similarly named tag", func(t *testing.T) {
		h := baselineHistoryFixture(t)
		baselineHistoryGit(t, h, "tag", "-a", "history-full-backup", "-m", "backup")
		r, out, stderr, code := baselineHistoryRun(t, h, "json")
		if code != 3 || stderr != "" || r.History.Status != "pending" || r.History.Tag.Action != "create" {
			t.Fatalf("%d %s %s", code, out, stderr)
		}
	})
	t.Run("non ancestor tag", func(t *testing.T) {
		h := baselineHistoryFixture(t)
		future := strings.TrimSpace(baselineHistoryGit(t, h, "commit-tree", "HEAD^{tree}", "-p", "HEAD", "-m", "Future commit"))
		baselineHistoryGit(t, h, "tag", "-a", "history-full", "-m", "future", future)
		r, out, stderr, code := baselineHistoryRun(t, h, "json")
		if code != 0 || stderr != "" || r.State != "current" || r.History.Status != "blocked" || len(r.History.Units) != 0 || !strings.Contains(r.History.Message, "not an ancestor") {
			t.Fatalf("%d %s %s", code, out, stderr)
		}
	})
	t.Run("configuration failure and opt out", func(t *testing.T) {
		h := baselineHistoryFixture(t)
		h.env.homeDirErr = errors.New("home unavailable\nretry")
		r, out, stderr, code := baselineHistoryRun(t, h, "json")
		if code != 0 || stderr != "" || r.State != "current" || r.History.Status != "blocked" || strings.Contains(r.History.Message, "\n") {
			t.Fatalf("%d %s %s", code, out, stderr)
		}
		r, out, stderr, code = baselineHistoryRun(t, h, "json", "--no-history", "--yes")
		if code != 0 || stderr != "" || r.History.Status != "skipped" || r.History.Units == nil || r.History.Refused == nil {
			t.Fatalf("%d %s %s", code, out, stderr)
		}
	})
	t.Run("dirty kind", func(t *testing.T) {
		h := baselineHistoryFixture(t)
		file := filepath.Join(h.repo, "docs/history/findings/entry.md")
		data, err := os.ReadFile(file)
		if err != nil {
			t.Fatal(err)
		}
		data = append(data, []byte("Dirty.\n")...)
		historyWrite(t, file, string(data))
		r, out, stderr, code := baselineHistoryRun(t, h, "json", "--yes")
		if code != 0 || stderr != "" || len(r.History.Refused) != 2 || r.History.Refused[1].Reason != "has uncommitted changes under docs/history/findings; commit or restore them and rerun" {
			t.Fatalf("%d %s %s", code, out, stderr)
		}
		after, err := os.ReadFile(file)
		if err != nil || !bytes.Equal(data, after) {
			t.Fatal("refused Finding changed")
		}
	})
	t.Run("apply text", func(t *testing.T) {
		h := baselineHistoryFixture(t)
		preview, _, _, _ := baselineHistoryRun(t, h, "json")
		_, out, stderr, code := baselineHistoryRun(t, h, "text", "--confirm-plan", preview.PlanDigest)
		var removed int
		var bytes int64
		for _, u := range preview.History.Units[:2] {
			removed += u.Files
			bytes += u.Bytes
		}
		want := fmt.Sprintf("History: applied\nHistory applied: 3 unit(s): wrote 2 Archive Record(s), reduced 1 file(s), removed %d file(s) (%d bytes) kept in Git at %.12s and tag history-full\nHistory refused: 1\n- refused docs/history/specs/0003-broken: %s\nHistory tag: created annotated history-full at %.12s; push it with git push origin history-full\n", removed, bytes, h.revision, preview.History.Refused[0].Reason, h.revision)
		if code != 0 || stderr != "" || !strings.HasPrefix(out, "Baseline update: verified\n") || !strings.Contains(out, want) {
			t.Fatalf("%d %s %s want %s", code, out, stderr, want)
		}
	})
}
