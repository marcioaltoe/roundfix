package cli

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"roundfix/internal/gittest"
	"roundfix/internal/spec"
	"roundfix/internal/testfixture"
)

func historyRemovalFixture(t *testing.T, baseline bool) (historyFixture, string, string, string, int) {
	t.Helper()
	var h historyFixture
	if baseline {
		h = historyFixture{repo: newBaselineUpdateRepository(t), home: t.TempDir()}
		h.env = commandEnvironmentForTest(t)
		h.env.homeDir, h.env.workDir = h.home, h.repo
		h.env.environ = withEnvValue(h.env.environ, "HOME", h.home)
		mustMkdir(t, filepath.Join(h.repo, "docs/specs"))
	} else {
		h = newHistoryFixture(t, false)
	}
	active := "docs/findings/original.md"
	historyWrite(t, filepath.Join(h.repo, active), "# Original\n\nFull finding.\n")
	provenance := historyCommit(t, h.repo, "docs: original finding")
	gittest.Run(t, h.repo, "rm", active)
	retired := "docs/history/findings/reduced.md"
	content := "---\nstatus: done\n---\n\n# Reduced\n\nSummary.\n\nFull text in Git at `" + provenance + "`: `" + active + "`.\n"
	historyWrite(t, filepath.Join(h.repo, retired), content)
	historyWrite(t, filepath.Join(h.repo, "docs/user-guide/a.md"), "# Guide\n\n[finding](../history/findings/reduced.md)\n")
	h.revision = historyCommit(t, h.repo, "docs: retire finding")
	return h, retired, active, provenance, len(content)
}

func TestHistorySanitizePlanTranscript(t *testing.T) {
	t.Parallel()
	if !historyRemovalTestHome(t) {
		return
	}
	h, retired, active, provenance, size := historyRemovalFixture(t, false)
	historyTag(t, h)
	before, refs, index := historySnapshot(t, h.repo), baselineHistoryGit(t, h, "for-each-ref"), baselineHistoryIndex(t, h)
	code, out, stderr := historyRemovalBinary(t, h, "sanitize")
	want := fmt.Sprintf("history sanitize plan: 1 unit(s) pending; 1 file(s) (%d bytes) leave docs/history\nfindings: removes 1 file(s) (%d bytes)\ncites docs/user-guide/a.md:3 names %s; rewrites to `%s@%s`\napply with: roundfix history sanitize --apply --batch <n> (needs the annotated tag history-full at or before HEAD)\n", size, size, retired, active, provenance)
	if code != 0 || stderr != "" || out != want {
		t.Fatalf("exit=%d stdout=%q stderr=%q want=%q", code, out, stderr, want)
	}
	baselineHistoryUnchanged(t, h, before, refs, index)
}

func TestHistorySanitizeApplyTranscript(t *testing.T) {
	t.Parallel()
	if !historyRemovalTestHome(t) {
		return
	}
	h, retired, active, provenance, size := historyRemovalFixture(t, false)
	historyTag(t, h)
	code, out, stderr := historyRemovalBinary(t, h, "sanitize", "--apply", "--batch", "1")
	want := fmt.Sprintf("history sanitize applied 1 unit(s): wrote 0 Archive Record(s), removed 1 file(s) (%d bytes) kept in Git at %.12s and tag history-full; rewrote 1 citation(s) in 1 file(s); promoted 0 file(s) to docs/references/; 0 unit(s) remain\n", size, h.revision)
	if code != 0 || stderr != "" || out != want {
		t.Fatalf("exit=%d stdout=%q stderr=%q want=%q", code, out, stderr, want)
	}
	assertPathMissing(t, filepath.Join(h.repo, retired))
	got, err := os.ReadFile(filepath.Join(h.repo, "docs/user-guide/a.md"))
	if err != nil || string(got) != "# Guide\n\nfinding (`"+active+"@"+provenance+"`)\n" {
		t.Fatalf("guide=%s err=%v", got, err)
	}
	code, out, stderr = historyRun(t, h, "sanitize")
	if code != 0 || stderr != "" || out != "history sanitize plan: 0 unit(s) pending; nothing pending\n" {
		t.Fatalf("second: %d %s %s", code, out, stderr)
	}
}

func TestHistorySanitizeRefusesAPathNoReachableCommitHolds(t *testing.T) {
	t.Parallel()
	if !historyRemovalTestHome(t) {
		return
	}
	h := historyFixture{repo: newBaselineUpdateRepository(t), home: t.TempDir()}
	h.env = commandEnvironmentForTest(t)
	h.env.homeDir, h.env.workDir = h.home, h.repo
	mustMkdir(t, filepath.Join(h.repo, "docs/specs"))
	historyWrite(t, filepath.Join(h.repo, "docs/history/backlog/covered.md"), "# Covered\n")
	historyCommit(t, h.repo, "docs: covered backlog")
	historyTag(t, h)
	gittest.Run(t, h.repo, "checkout", "-b", "feat/unreachable")
	path := "docs/history/findings/unreachable.md"
	historyWrite(t, filepath.Join(h.repo, path), "# Branch only\n")
	historyCommit(t, h.repo, "docs: branch finding")
	before := historySnapshot(t, h.repo)
	code, out, stderr := historyRun(t, h, "sanitize", "--apply", "--batch", "1")
	reason := "history-full does not hold " + path + " and no commit the default branch reaches holds it"
	if code != 2 || out != "" || !strings.Contains(stderr, reason) {
		t.Fatalf("%d %s %s", code, out, stderr)
	}
	if !reflect.DeepEqual(before, historySnapshot(t, h.repo)) {
		t.Fatal("refusal changed files")
	}
	p := planBaselineUpdateHistory(t.Context(), h.repo, h.env)
	if len(p.Report.Refused) != 1 || p.Report.Refused[0].Unit != "findings" || p.Report.Refused[0].Reason != reason || len(p.selected) != 1 || p.selected[0].name != "backlog" {
		t.Fatalf("history=%+v", p.Report)
	}
	r, out, stderr, code := baselineHistoryRun(t, h, "json", "--yes")
	if code != 0 || stderr != "" || len(r.History.Refused) != 1 || r.History.Applied == nil || r.History.Applied.Removed != 1 {
		t.Fatalf("%d %s %s", code, out, stderr)
	}
	assertPathMissing(t, filepath.Join(h.repo, "docs/history/backlog/covered.md"))
	if _, err := os.Stat(filepath.Join(h.repo, path)); err != nil {
		t.Fatal("refused path removed")
	}

}

func TestBaselineUpdateRemovesRetiredHistoryAndRewritesCitations(t *testing.T) {
	t.Parallel()
	if !historyRemovalTestHome(t) {
		return
	}
	h, retired, active, provenance, _ := historyRemovalFixture(t, true)
	// An Archive Record is a rewrite source, even though it is under history.
	record := spec.ArchiveRecord{Spec: "0001-kept", Title: "Kept", Disposition: spec.ArchiveNoQA, Source: "docs/specs/0001-kept", SourceRevision: h.revision, Outcome: "See [finding](../findings/reduced.md)."}
	content, err := spec.RenderArchiveRecord(record)
	if err != nil {
		t.Fatal(err)
	}
	historyWrite(t, filepath.Join(h.repo, "docs/history/specs/0001-kept.md"), string(content))
	superseded := record
	superseded.Spec, superseded.Source, superseded.Disposition, superseded.SupersededBy, superseded.Outcome = "0002-uncited", "docs/specs/0002-uncited", spec.ArchiveSuperseded, "0001-kept", "Retired."
	retiredRecord, err := spec.RenderArchiveRecord(superseded)
	if err != nil {
		t.Fatal(err)
	}
	historyWrite(t, filepath.Join(h.repo, "docs/history/specs/0002-uncited.md"), string(retiredRecord))
	h.revision = historyCommit(t, h.repo, "docs: cite finding in record")
	before, refs, index := historySnapshot(t, h.repo), baselineHistoryGit(t, h, "for-each-ref"), baselineHistoryIndex(t, h)
	preview, out, stderr, code := baselineHistoryRun(t, h, "json")
	if code != 3 || stderr != "" || len(preview.History.Units) != 2 || preview.History.Units[0].Action != "remove" || strings.Contains(out, "bytesAfter") {
		t.Fatalf("%d %s %s", code, out, stderr)
	}
	if len(preview.History.Citations) != 2 {
		t.Fatalf("citations=%+v", preview.History.Citations)
	}
	for _, c := range preview.History.Citations {
		if c.Rewrite != "`"+active+"@"+provenance+"`" {
			t.Fatalf("citation=%+v", c)
		}
	}
	baselineHistoryUnchanged(t, h, before, refs, index)
	// Digest binds the rewrite text and has no retired bytesAfter field.
	p := planBaselineUpdateHistory(t.Context(), h.repo, h.env)
	digest, err := p.digest("baseline")
	if err != nil {
		t.Fatal(err)
	}
	// This is the API Contract 12 payload: remove actions, rewrite destinations,
	// and no bytesAfter field, including a superseded removal unit.
	locator := "`" + active + "@" + provenance + "`"
	payload := fmt.Sprintf(`{"revision":%q,"action":"create","commit":%q,"units":[{"unit":"findings","action":"remove","paths":[%q]},{"unit":"superseded","action":"remove","paths":["docs/history/specs/0002-uncited.md"]}],"citations":[{"path":"docs/history/specs/0001-kept.md","line":%d,"target":%q,"rewrite":%q},{"path":"docs/user-guide/a.md","line":3,"target":%q,"rewrite":%q}]}`, h.revision, h.revision, retired, p.Report.Citations[0].Line, retired, locator, retired, locator)
	sum := sha256.Sum256([]byte("roundfix/baseline-update-history/v1\x00baseline\x00" + payload))
	if digest != "sha256:"+hex.EncodeToString(sum[:]) {
		t.Fatalf("digest=%s payload=%s", digest, payload)
	}
	p.Report.Citations[0].Rewrite = "changed"
	changed, err := p.digest("baseline")
	if err != nil || changed == digest {
		t.Fatal("digest did not bind rewrite")
	}
	result, out, stderr, code := baselineHistoryRun(t, h, "json", "--confirm-plan", preview.PlanDigest)
	if code != 0 || stderr != "" || result.History.Applied == nil || result.History.Applied.Reduced != 0 || result.History.Applied.Removed != 2 || result.History.Applied.RewrittenCitations != 2 || result.History.Applied.RewrittenFiles != 2 {
		t.Fatalf("%d %s %s", code, out, stderr)
	}
	var raw map[string]any
	if err := json.Unmarshal([]byte(out), &raw); err != nil {
		t.Fatal(err)
	}
	assertPathMissing(t, filepath.Join(h.repo, retired))
	assertPathMissing(t, filepath.Join(h.repo, "docs/history/specs/0002-uncited.md"))
	again, out, stderr, code := baselineHistoryRun(t, h, "json")
	if code != 0 || stderr != "" || again.History.Status != "current" {
		t.Fatalf("second: %d %s %s", code, out, stderr)
	}
}

func TestBaselineUpdateRemovesWhatTheBaselineUsedToRelocate(t *testing.T) {
	t.Parallel()
	if !historyRemovalTestHome(t) {
		return
	}
	h := historyFixture{repo: newBaselineUpdateRepository(t), home: t.TempDir()}
	h.env = commandEnvironmentForTest(t)
	h.env.homeDir, h.env.workDir = h.home, h.repo
	files := []string{"docs/findings/_archived/a.md", "_archived/findings/a.md", "_archived/backlog/a.md", "_archived/adr/a.md", "docs/adr/0001-old.md", "docs/backlog/done.md", "docs/specs/reviews/pr-1/outcome.md"}
	for _, path := range files {
		historyWrite(t, filepath.Join(h.repo, path), "---\nstatus: deprecated\n---\n# Old\n")
	}
	historyWrite(t, filepath.Join(h.repo, files[6]), "---\npull_request_state: closed\nrecorded_at: 2026-10-16\n---\n")
	h.revision = historyCommit(t, h.repo, "docs: legacy layout")
	r, out, stderr, code := baselineHistoryRun(t, h, "json", "--yes")
	if code != 0 || stderr != "" || len(r.HistoryMoves) != 0 || r.History.Applied == nil || r.History.Applied.Removed != len(files) {
		t.Fatalf("%d %s %s", code, out, stderr)
	}
	for _, path := range files {
		assertPathMissing(t, filepath.Join(h.repo, path))
	}
	again, out, stderr, code := baselineHistoryRun(t, h, "json")
	if code != 0 || stderr != "" || again.History.Status != "current" {
		t.Fatalf("second: %d %s %s", code, out, stderr)
	}
}

func TestUpgradeNoticeNamesReducedEntriesAsPending(t *testing.T) {
	t.Parallel()
	if !historyRemovalTestHome(t) {
		return
	}
	prepareNoticeUpgrade(t, "current")
	env := commandEnvironmentForTest(t)
	gittest.InitRepo(t, env.workDir, "-b", "main")
	mustMkdir(t, filepath.Join(env.workDir, "docs/specs"))
	historyWrite(t, filepath.Join(env.workDir, "docs/history/findings/reduced.md"), "# Reduced\n\nFull text in Git at `"+strings.Repeat("a", 40)+"`: `docs/findings/a.md`.\n")
	historyWrite(t, filepath.Join(env.workDir, "docs/adr/0001-old.md"), "---\nstatus: rejected\n---\n# Old\n")
	r := spec.ArchiveRecord{Spec: "0001-superseded", Title: "Superseded", Disposition: spec.ArchiveSuperseded, SupersededBy: "0002-successor", Source: "docs/specs/0001-superseded", SourceRevision: strings.Repeat("a", 40), Outcome: "Retired."}
	content, err := spec.RenderArchiveRecord(r)
	if err != nil {
		t.Fatal(err)
	}
	historyWrite(t, filepath.Join(env.workDir, "docs/history/specs/0001-superseded.md"), string(content))
	historyCommit(t, env.workDir, "docs: pending units")
	var out, stderr bytes.Buffer
	code := runWithContext(context.Background(), []string{"upgrade"}, &out, &stderr, env)
	want := upgradeFixtureNotice() + "roundfix: history: 3 unit(s) pending sanitize (findings, adr, superseded); run roundfix baseline update to plan them\n"
	if code != 0 || out.String() != "already current 1.0.0\n" || stderr.String() != want {
		t.Fatalf("%d %s %s want=%s", code, &out, &stderr, want)
	}
}

func historyRemovalBinary(t *testing.T, h historyFixture, args ...string) (int, string, string) {
	t.Helper()
	binary := filepath.Join(t.TempDir(), "roundfix")
	build := exec.CommandContext(t.Context(), "go", "build", "-buildvcs=false", "-o", binary, "./cmd/roundfix")
	build.Dir = filepath.Join("..", "..")
	if output, err := build.CombinedOutput(); err != nil {
		t.Fatalf("build: %v\n%s", err, output)
	}
	command := exec.CommandContext(t.Context(), binary, append([]string{"history"}, args...)...)
	command.Dir = h.repo
	command.Env = withEnvValue(gittest.IsolatedEnv(), "HOME", h.home)
	var out, stderr bytes.Buffer
	command.Stdout, command.Stderr = &out, &stderr
	err := command.Run()
	code, note := testfixture.ChildExitCode(err)
	if note != "" {
		t.Fatal(note)
	}
	return code, out.String(), stderr.String()
}

// Preserve the toolchain caches when isolating HOME, so executable transcripts
// build offline and never populate a temporary home with read-only modules.
func historyRemovalTestHome(t *testing.T) bool {
	t.Helper()
	const key = "ROUNDFIX_HISTORY_REMOVAL_TEST"
	if os.Getenv(key) == t.Name() {
		return true
	}
	output, err := exec.CommandContext(t.Context(), "go", "env", "-json", "GOCACHE", "GOMODCACHE").Output()
	if err != nil {
		t.Fatalf("read tool caches: %v", err)
	}
	var caches map[string]string
	if err := json.Unmarshal(output, &caches); err != nil {
		t.Fatal(err)
	}
	command := exec.CommandContext(t.Context(), os.Args[0], "-test.run=^"+regexp.QuoteMeta(t.Name())+"$", "-test.count=1")
	command.Env = withEnvValue(gittest.IsolatedEnv(), "HOME", t.TempDir())
	for k, v := range caches {
		command.Env = withEnvValue(command.Env, k, v)
	}
	command.Env = withEnvValue(command.Env, "GOPROXY", "off")
	command.Env = withEnvValue(command.Env, "GOSUMDB", "off")
	command.Env = withEnvValue(command.Env, key, t.Name())
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("isolated history test: %v\n%s", err, output)
	}
	return false
}
