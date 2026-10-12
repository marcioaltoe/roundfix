package cli

// Suite: baseline update History Citations, promotions and stale advice.
// Invariant: approval binds citations and promotions; citing files never change.
// Boundary IN: adopted disposable Git repositories, temporary homes and CLI flags.
// Boundary OUT: stdout, stderr, JSON, exit codes, file bytes, index and refs.

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"roundfix/internal/spec"
)

const baselineCitationFolder = "docs/history/specs/0001-legacy"
const baselineCitationLesson = baselineCitationFolder + "/references/2026-01-01-lesson.md"
const baselineCitationFinding = "docs/findings/2026-01-04-gap.md"
const baselineCitationHandoff = "docs/handoffs/2026-01-03-notes.md"
const baselineCitationDestination = "docs/references/2026-01-01-lesson.md"

func baselineCitationFixture(t *testing.T, citing bool) historyFixture {
	t.Helper()
	h := historyFixture{repo: newBaselineUpdateRepository(t), home: t.TempDir()}
	h.env = commandEnvironmentForTest(t)
	h.env.homeDir, h.env.homeDirErr, h.env.workDir, h.env.workDirErr = h.home, nil, h.repo, nil
	h.env.environ = withEnvValue(h.env.environ, "HOME", h.home)
	if err := os.MkdirAll(filepath.Join(h.repo, "docs/specs"), 0755); err != nil {
		t.Fatal(err)
	}
	historyLegacyFolder(t, h, "0001-legacy")
	historyWrite(t, filepath.Join(h.repo, baselineCitationLesson), "# Lesson\n\nReusable source bytes.\n")
	if citing {
		historyWrite(t, filepath.Join(h.repo, baselineCitationFinding), "---\nstatus: pending\n---\n\n# Gap\n\n[lesson](../history/specs/0001-legacy/references/2026-01-01-lesson.md)\n")
		historyWrite(t, filepath.Join(h.repo, baselineCitationHandoff), "# Notes\n\n[PRD](../history/specs/0001-legacy/_prd.md)\n")
	}
	h.revision = historyCommit(t, h.repo, "Legacy history (#42)")
	return h
}

// Only the transcript's declared ... blocks vary; match all other lines,
// their order, stderr and exit.
func assertBaselineCitationTranscript(t *testing.T, out, stderr string, code int, want string, wantCode int, wantError ...string) {
	t.Helper()
	pattern := regexp.QuoteMeta(want)
	pattern = strings.ReplaceAll(pattern, regexp.QuoteMeta("...\n"), "(?:[^\\n]*\\n)*")
	expectedError := ""
	if len(wantError) > 0 {
		expectedError = wantError[0]
	}
	if code != wantCode || stderr != expectedError || !regexp.MustCompile("\\A"+pattern+"\\z").MatchString(out) {
		t.Fatalf("exit=%d stderr=%q\nstdout:\n%s\nwant exit=%d stdout:\n%s", code, stderr, out, wantCode, want)
	}
}

// Exercise the executable and go run wrapper that the Surface Transcripts name.
// Keep the adopter's home private while reusing only the Go tool's build/module
// caches, so running this local command never needs a dependency download.
func baselineCitationGoRun(t *testing.T, h historyFixture, args ...string) (string, string, int) {
	t.Helper()
	var caches map[string]string
	cacheOutput, err := exec.CommandContext(t.Context(), "go", "env", "-json", "GOCACHE", "GOMODCACHE").Output()
	if err != nil {
		t.Fatalf("resolve Go caches: %v", err)
	}
	if err := json.Unmarshal(cacheOutput, &caches); err != nil {
		t.Fatal(err)
	}
	commandArgs := append([]string{"run", "-buildvcs=false", "./cmd/roundfix", "baseline", "update", "--repo", h.repo, "--no-skills", "--format", "text"}, args...)
	command := exec.CommandContext(t.Context(), "go", commandArgs...)
	command.Dir, err = filepath.Abs(filepath.Join("..", ".."))
	if err != nil {
		t.Fatal(err)
	}
	command.Env = h.env.environ
	for key, value := range caches {
		command.Env = withEnvValue(command.Env, key, value)
	}
	command.Env = withEnvValue(command.Env, "GOPROXY", "off")
	var stdout, stderr bytes.Buffer
	command.Stdout, command.Stderr = &stdout, &stderr
	err = command.Run()
	code := 0
	if err != nil {
		var exit *exec.ExitError
		if !errors.As(err, &exit) {
			t.Fatalf("run transcript: %v", err)
		}
		code = exit.ExitCode()
	}
	return stdout.String(), stderr.String(), code
}

func baselineCitationLines(promoted bool) string {
	suffix := " (live; promote with --promote " + baselineCitationLesson + ")"
	if promoted {
		suffix = " (promoted to " + baselineCitationDestination + ")"
	}
	return "History citations: 2\n- cites " + baselineCitationFinding + ":7 names " + baselineCitationLesson + suffix + "\n- cites " + baselineCitationHandoff + ":3 names " + baselineCitationFolder + "/_prd.md\n"
}

func baselineCitationExpected(promoted bool) []baselineUpdateHistoryCitation {
	c := baselineUpdateHistoryCitation{Path: baselineCitationFinding, Line: 7, Target: baselineCitationLesson, Live: true, Offer: baselineCitationLesson}
	if promoted {
		c.Offer, c.PromotedTo = "", baselineCitationDestination
	}
	return []baselineUpdateHistoryCitation{c, {Path: baselineCitationHandoff, Line: 3, Target: baselineCitationFolder + "/_prd.md", Live: false}}
}

func assertBaselineCitationJSON(t *testing.T, output, field string, want []baselineUpdateHistoryCitation) {
	t.Helper()
	var result struct {
		History map[string]json.RawMessage `json:"history"`
	}
	if err := json.Unmarshal([]byte(output), &result); err != nil {
		t.Fatal(err)
	}
	var raw []map[string]json.RawMessage
	if err := json.Unmarshal(result.History[field], &raw); err != nil {
		t.Fatal(err)
	}
	if raw == nil || len(raw) != len(want) {
		t.Fatalf("%s=%s", field, result.History[field])
	}
	for i, c := range want {
		keys := map[string]bool{"path": true, "line": true, "target": true, "live": true}
		if c.Offer != "" {
			keys["offer"] = true
		}
		if c.PromotedTo != "" {
			keys["promotedTo"] = true
		}
		if len(raw[i]) != len(keys) {
			t.Fatalf("citation keys=%v want=%v", raw[i], keys)
		}
		for k := range raw[i] {
			if !keys[k] {
				t.Fatalf("unexpected key %q", k)
			}
		}
	}
	var got []baselineUpdateHistoryCitation
	if err := json.Unmarshal(result.History[field], &got); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("%s=%+v want=%+v", field, got, want)
	}
}

func TestBaselineUpdateListsHistoryCitationsTranscript(t *testing.T) {
	t.Parallel()
	h := baselineCitationFixture(t, true)
	before, refs, index := historySnapshot(t, h.repo), baselineHistoryGit(t, h, "for-each-ref"), baselineHistoryIndex(t, h)
	r, output, stderr, code := baselineHistoryRun(t, h, "json")
	if code != 3 || stderr != "" || r.State != "plan_ready" || r.History.Status != "pending" || len(r.History.Units) != 1 {
		t.Fatalf("%d %s %s", code, output, stderr)
	}
	assertBaselineCitationJSON(t, output, "citations", baselineCitationExpected(false))
	assertBaselineCitationJSON(t, output, "staleCitations", []baselineUpdateHistoryCitation{})
	u := r.History.Units[0]
	output, stderr, code = baselineCitationGoRun(t, h)
	want := fmt.Sprintf("Baseline update: plan ready\nCategory: approval\nResult: guidance matches the current Baseline catalog; 1 history unit(s) pending sanitize\n...\nHistory: pending\nHistory units: 1 (%d file(s), %d bytes leave docs/history)\n- folder %s: removes %d file(s) (%d bytes) and writes %s (%s)\n", u.Files, u.Bytes, u.Unit, u.Files, u.Bytes, u.Record, u.Disposition)
	want += baselineCitationLines(false) + fmt.Sprintf("History tag: creates annotated history-full at %.12s; push it with git push origin history-full\nPlan Digest: %s\nNext action: review the plan and rerun with --confirm-plan %s, or use --yes to approve the digest computed in one invocation\n", h.revision, r.PlanDigest, r.PlanDigest)
	assertBaselineCitationTranscript(t, output, stderr, code, want, 1, "exit status 3\n")
	baselineHistoryUnchanged(t, h, before, refs, index)
}

func TestBaselineUpdatePromotesOnlyOnRequestTranscript(t *testing.T) {
	t.Parallel()
	h := baselineCitationFixture(t, true)
	before := historySnapshot(t, h.repo)
	preview, _, _, _ := baselineHistoryRun(t, h, "json")
	promoted, output, stderr, code := baselineHistoryRun(t, h, "json", "--promote", baselineCitationLesson)
	if code != 3 || stderr != "" || preview.PlanDigest == promoted.PlanDigest {
		t.Fatalf("promotion preview: %d %s %s", code, output, stderr)
	}
	assertBaselineCitationJSON(t, output, "citations", baselineCitationExpected(true))
	if !reflect.DeepEqual(before, historySnapshot(t, h.repo)) {
		t.Fatal("promotion preview wrote files")
	}
	output, stderr, code = baselineCitationGoRun(t, h, "--yes", "--promote", baselineCitationLesson)
	u := preview.History.Units[0]
	want := fmt.Sprintf("Baseline update: verified\n...\nHistory: applied\nHistory applied: 1 unit(s): wrote 1 Archive Record(s), reduced 0 file(s), removed %d file(s) (%d bytes) kept in Git at %.12s and tag history-full\n", u.Files, u.Bytes, h.revision)
	want += baselineCitationLines(true) + fmt.Sprintf("History tag: created annotated history-full at %.12s; push it with git push origin history-full\n...\n", h.revision)
	assertBaselineCitationTranscript(t, output, stderr, code, want, 0)
	after := historySnapshot(t, h.repo)
	if !bytes.Equal(after[baselineCitationDestination], before[baselineCitationLesson]) {
		t.Fatal("copy differs from removed source")
	}
	for _, path := range []string{baselineCitationFinding, baselineCitationHandoff} {
		if !bytes.Equal(before[path], after[path]) {
			t.Fatalf("citing file changed: %s", path)
		}
	}
	record, err := spec.ParseArchiveRecord(after[baselineCitationFolder+".md"])
	if err != nil || !reflect.DeepEqual(record.Promoted, []string{baselineCitationDestination}) || len(record.Relinked) != 0 {
		t.Fatalf("record=%+v err=%v", record, err)
	}
	assertPathMissing(t, filepath.Join(h.repo, baselineCitationFolder))
	if blob := baselineHistoryGit(t, h, "show", h.revision+":"+baselineCitationLesson); blob != string(after[baselineCitationDestination]) {
		t.Fatal("copy differs from committed blob")
	}
	t.Run("without request", func(t *testing.T) {
		h := baselineCitationFixture(t, true)
		before := historySnapshot(t, h.repo)
		r, output, stderr, code := baselineHistoryRun(t, h, "json", "--yes")
		if code != 0 || stderr != "" || r.History.Status != "applied" {
			t.Fatalf("%d %s %s", code, output, stderr)
		}
		after := historySnapshot(t, h.repo)
		if _, exists := after[baselineCitationDestination]; exists {
			t.Fatal("promotion without --promote")
		}
		for _, path := range []string{baselineCitationFinding, baselineCitationHandoff} {
			if !bytes.Equal(before[path], after[path]) {
				t.Fatalf("citing file changed: %s", path)
			}
		}
		record, err := spec.ParseArchiveRecord(after[baselineCitationFolder+".md"])
		if err != nil || len(record.Promoted) != 0 {
			t.Fatalf("record=%+v err=%v", record, err)
		}
	})
}

func TestBaselineUpdateReportsStaleHistoryCitationsTranscript(t *testing.T) {
	t.Parallel()
	h := baselineCitationFixture(t, true)
	_, output, stderr, code := baselineHistoryRun(t, h, "json", "--yes", "--promote", baselineCitationLesson)
	if code != 0 {
		t.Fatalf("apply: %d %s %s", code, output, stderr)
	}
	before, refs, index := historySnapshot(t, h.repo), baselineHistoryGit(t, h, "for-each-ref"), baselineHistoryIndex(t, h)
	r, output, stderr, code := baselineHistoryRun(t, h, "json")
	if code != 0 || stderr != "" || r.State != "current" || r.Category != "" || r.History.Status != "current" || r.PlanDigest != baselineHistoryDigest(t, h) {
		t.Fatalf("%d %s %s", code, output, stderr)
	}
	wantCitations := baselineCitationExpected(false)
	wantCitations[0].Offer = ""
	assertBaselineCitationJSON(t, output, "staleCitations", wantCitations)
	assertBaselineCitationJSON(t, output, "citations", []baselineUpdateHistoryCitation{})
	output, stderr, code = baselineCitationGoRun(t, h)
	want := "Baseline update: current\nResult: the repository already matches the current Baseline catalog\n...\nHistory: current\nHistory stale citations: 2\n"
	for _, c := range wantCitations {
		want += fmt.Sprintf("- cites %s:%d names %s; history-full holds it: git show history-full:%s\n", c.Path, c.Line, c.Target, c.Target)
	}
	want += fmt.Sprintf("History tag: history-full at %.12s holds every planned path\n...\n", h.revision)
	assertBaselineCitationTranscript(t, output, stderr, code, want, 0)
	baselineHistoryUnchanged(t, h, before, refs, index)
	historyWrite(t, filepath.Join(h.repo, baselineCitationFinding), "# No stale citation\n")
	next, output, stderr, code := baselineHistoryRun(t, h, "json")
	if code != 0 || stderr != "" || next.State != r.State || next.Category != r.Category || next.PlanDigest != r.PlanDigest || len(next.History.StaleCitations) != 1 {
		t.Fatalf("stale advice: %d %s %s", code, output, stderr)
	}
}

func TestBaselineUpdateHistoryDigestBindsCitations(t *testing.T) {
	t.Parallel()
	h := baselineCitationFixture(t, true)
	prior, _, _, _ := baselineHistoryRun(t, h, "json")
	path := filepath.Join(h.repo, baselineCitationFinding)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	historyWrite(t, path, "\n"+string(data))
	before, refs, index := historySnapshot(t, h.repo), baselineHistoryGit(t, h, "for-each-ref"), baselineHistoryIndex(t, h)
	next, output, stderr, code := baselineHistoryRun(t, h, "json")
	if code != 3 || stderr != "" || next.PlanDigest == prior.PlanDigest || next.History.Citations[0].Line != 8 {
		t.Fatalf("changed citation: %d %s %s", code, output, stderr)
	}
	rejected, output, stderr, code := baselineHistoryRun(t, h, "json", "--confirm-plan", prior.PlanDigest)
	if code != 3 || rejected.Category != "approval" || rejected.PlanDigest != next.PlanDigest {
		t.Fatalf("stale digest: %d %s %s", code, output, stderr)
	}
	baselineHistoryUnchanged(t, h, before, refs, index)
}

func TestBaselineUpdateRefusesAnInvalidPromotion(t *testing.T) {
	t.Parallel()
	for _, tt := range []struct {
		name, path string
		extra      []string
		prepare    func(*testing.T, historyFixture)
	}{
		{name: "core", path: baselineCitationFolder + "/_prd.md"},
		{name: "valid then core", path: baselineCitationFolder + "/_prd.md", extra: []string{"--promote", baselineCitationLesson}},
		{name: "outside units", path: baselineCitationFinding},
		{name: "occupied destination", path: baselineCitationLesson, prepare: func(t *testing.T, h historyFixture) {
			historyWrite(t, filepath.Join(h.repo, baselineCitationDestination), "Do not overwrite\n")
		}},
		{name: "no history", path: baselineCitationLesson, extra: []string{"--no-history"}},
		{name: "unclean", path: "./" + baselineCitationLesson},
		{name: "absolute", path: "/" + baselineCitationLesson},
		{name: "directory", path: baselineCitationFolder + "/references"},
		{name: "missing", path: baselineCitationFolder + "/references/missing.md"},
		{name: "duplicate", path: baselineCitationLesson, extra: []string{"--promote", baselineCitationLesson}},
		{name: "refused unit", path: baselineCitationLesson, prepare: func(t *testing.T, h historyFixture) {
			historyWrite(t, filepath.Join(h.repo, baselineCitationLesson), "Dirty history\n")
		}},
	} {
		t.Run(tt.name, func(t *testing.T) {
			h := baselineCitationFixture(t, true)
			if tt.prepare != nil {
				tt.prepare(t, h)
			}
			// Give guidance its own pending write: invalid promotions must stop
			// before either half of the combined plan mutates the repository.
			manifest := ReadBaselineSetupManifest(t, h.repo)
			manifest.CatalogDigest = "sha256:" + strings.Repeat("0", 64)
			WriteBaselineSetupManifest(t, h.repo, manifest)
			if tt.name == "core" {
				preview, output, stderr, code := baselineHistoryRun(t, h, "json", "--no-history")
				if code != 3 || stderr != "" || len(preview.FileChanges) == 0 {
					t.Fatalf("pending guidance: %d %s %s", code, output, stderr)
				}
			}
			before, refs, index := historySnapshot(t, h.repo), baselineHistoryGit(t, h, "for-each-ref"), baselineHistoryIndex(t, h)
			args := append([]string{"--yes", "--promote", tt.path}, tt.extra...)
			r, output, stderr, code := baselineHistoryRun(t, h, "json", args...)
			if code != 2 || r.Category != "invalid" || !strings.Contains(r.Message, tt.path) || !strings.Contains(stderr, tt.path) || r.History.Citations == nil || r.History.StaleCitations == nil {
				t.Fatalf("%d %s %s", code, output, stderr)
			}
			baselineHistoryUnchanged(t, h, before, refs, index)
		})
	}
}

func TestBaselineUpdateWithoutCitationsKeepsDigestAndText(t *testing.T) {
	t.Parallel()
	h := baselineCitationFixture(t, false)
	p := planBaselineUpdateHistory(context.Background(), h.repo, h.env)
	// Pre-ADR-0260 payload and field order are the compatibility oracle.
	type legacyUnit struct {
		Unit         string   `json:"unit"`
		Action       string   `json:"action"`
		Paths        []string `json:"paths"`
		Record       string   `json:"record,omitempty"`
		RecordDigest string   `json:"recordDigest,omitempty"`
		BytesAfter   *int64   `json:"bytesAfter,omitempty"`
	}
	c := p.selected[0].conversion
	recordHash := sha256.Sum256(c.Rendered)
	payload := struct {
		Revision string       `json:"revision"`
		Action   string       `json:"action"`
		Commit   string       `json:"commit"`
		Units    []legacyUnit `json:"units"`
	}{h.revision, "create", h.revision, []legacyUnit{{Unit: baselineCitationFolder, Action: "convert", Paths: p.selected[0].files, Record: c.RecordPath, RecordDigest: hex.EncodeToString(recordHash[:])}}}
	data, err := json.Marshal(payload)
	if err != nil {
		t.Fatal(err)
	}
	hash := sha256.Sum256(append([]byte("roundfix/baseline-update-history/v1\x00"+baselineHistoryDigest(t, h)+"\x00"), data...))
	digest := "sha256:" + hex.EncodeToString(hash[:])
	r, output, stderr, code := baselineHistoryRun(t, h, "json")
	if code != 3 || stderr != "" || r.PlanDigest != digest {
		t.Fatalf("%d %s %s want digest=%s", code, output, stderr, digest)
	}
	assertBaselineCitationJSON(t, output, "citations", []baselineUpdateHistoryCitation{})
	assertBaselineCitationJSON(t, output, "staleCitations", []baselineUpdateHistoryCitation{})
	u := r.History.Units[0]
	_, output, stderr, code = baselineHistoryRun(t, h, "text")
	want := fmt.Sprintf("Baseline update: plan ready\nCategory: approval\nResult: guidance matches the current Baseline catalog; 1 history unit(s) pending sanitize\n...\nHistory: pending\nHistory units: 1 (%d file(s), %d bytes leave docs/history)\n- folder %s: removes %d file(s) (%d bytes) and writes %s (%s)\nHistory tag: creates annotated history-full at %.12s; push it with git push origin history-full\nPlan Digest: %s\nNext action: review the plan and rerun with --confirm-plan %s, or use --yes to approve the digest computed in one invocation\n", u.Files, u.Bytes, u.Unit, u.Files, u.Bytes, u.Record, u.Disposition, h.revision, digest, digest)
	assertBaselineCitationTranscript(t, output, stderr, code, want, 3)
}

func TestBaselineUpdateHistoryCitationBoundaries(t *testing.T) {
	t.Parallel()
	h := baselineCitationFixture(t, true)
	historyWrite(t, filepath.Join(h.repo, baselineCitationFinding), "[lesson](../history/specs/0001-legacy/references/2026-01-01-lesson.md)\n[core](../history/specs/0001-legacy/_prd.md)\n[folder](../history/specs/0001-legacy/)\n[reduced](../history/findings/retired.md)\n[removed](../history/reviews/pr-1/review.md)\n")
	historyWrite(t, filepath.Join(h.repo, "docs/history/findings/retired.md"), "---\nstatus: closed\n---\n\n# Retired\n\nFirst paragraph.\n\n## Detail\n\nFull detail.\n")
	historyWrite(t, filepath.Join(h.repo, "docs/history/reviews/pr-1/review.md"), "# Review\n")
	h.revision = historyCommit(t, h.repo, "Mixed history")
	r, output, stderr, code := baselineHistoryRun(t, h, "json")
	if code != 3 || stderr != "" || r.History.Status != "pending" || len(r.History.Units) != 3 {
		t.Fatalf("%d %s %s", code, output, stderr)
	}
	want := []baselineUpdateHistoryCitation{
		{Path: baselineCitationFinding, Line: 1, Target: baselineCitationLesson, Live: true, Offer: baselineCitationLesson},
		{Path: baselineCitationFinding, Line: 2, Target: baselineCitationFolder + "/_prd.md", Live: true},
		{Path: baselineCitationFinding, Line: 5, Target: "docs/history/reviews/pr-1/review.md", Live: true},
		{Path: baselineCitationHandoff, Line: 3, Target: baselineCitationFolder + "/_prd.md"},
		{Path: baselineCitationFinding, Line: 3, Target: baselineCitationFolder + "/", Live: true},
	}
	assertBaselineCitationJSON(t, output, "citations", want)
	_, output, stderr, code = baselineHistoryRun(t, h, "text")
	for _, line := range []string{
		"- cites " + baselineCitationFinding + ":2 names " + baselineCitationFolder + "/_prd.md (live)\n",
		"- cites " + baselineCitationFinding + ":3 names " + baselineCitationFolder + "/ (live)\n",
	} {
		if !strings.Contains(output, line) {
			t.Fatalf("missing %q\n%s", line, output)
		}
	}
}

func TestBaselineUpdateHistoryCitationDetectorFailureIsLocal(t *testing.T) {
	t.Parallel()
	h := baselineCitationFixture(t, true)
	if err := os.Remove(filepath.Join(h.repo, baselineCitationHandoff)); err != nil {
		t.Fatal(err)
	}
	before := historySnapshot(t, h.repo)
	r, output, stderr, code := baselineHistoryRun(t, h, "json", "--yes")
	if code != 0 || stderr != "" || r.State != "verified" || r.History.Status != "blocked" || r.History.Tag != nil || len(r.History.Units) != 0 || r.History.Message == "" || strings.ContainsAny(r.History.Message, "\r\n") || r.PlanDigest != baselineHistoryDigest(t, h) {
		t.Fatalf("%d %s %s", code, output, stderr)
	}
	if !reflect.DeepEqual(before, historySnapshot(t, h.repo)) {
		t.Fatal("blocked history changed")
	}
	assertBaselineCitationJSON(t, output, "citations", []baselineUpdateHistoryCitation{})
	assertBaselineCitationJSON(t, output, "staleCitations", []baselineUpdateHistoryCitation{})
}

func TestBaselineUpdateHistoryStaleDetectorFailureKeepsOutcome(t *testing.T) {
	t.Parallel()
	h := baselineCitationFixture(t, true)
	_, output, stderr, code := baselineHistoryRun(t, h, "json", "--yes")
	if code != 0 {
		t.Fatalf("%d %s %s", code, output, stderr)
	}
	if err := os.Remove(filepath.Join(h.repo, baselineCitationHandoff)); err != nil {
		t.Fatal(err)
	}
	r, output, stderr, code := baselineHistoryRun(t, h, "json")
	if code != 0 || stderr != "" || r.State != "current" || r.Category != "" || r.History.Status != "current" || r.History.Message == "" || r.PlanDigest != baselineHistoryDigest(t, h) {
		t.Fatalf("%d %s %s", code, output, stderr)
	}
	assertBaselineCitationJSON(t, output, "staleCitations", []baselineUpdateHistoryCitation{})
}

func TestBaselineUpdateHistoryPromotionsRepeatAndRefuseSharedBasenames(t *testing.T) {
	t.Parallel()
	for _, collision := range []bool{false, true} {
		name := "distinct files"
		if collision {
			name = "shared basename"
		}
		t.Run(name, func(t *testing.T) {
			h := baselineCitationFixture(t, true)
			file := "other.md"
			if collision {
				file = filepath.Base(baselineCitationLesson)
			}
			second := baselineCitationFolder + "/other/" + file
			historyWrite(t, filepath.Join(h.repo, second), "Second source\n")
			h.revision = historyCommit(t, h.repo, "Second source")
			before, refs, index := historySnapshot(t, h.repo), baselineHistoryGit(t, h, "for-each-ref"), baselineHistoryIndex(t, h)
			r, output, stderr, code := baselineHistoryRun(t, h, "json", "--yes", "--promote", baselineCitationLesson, "--promote", second)
			if collision {
				if code != 2 || r.Category != "invalid" || !strings.Contains(r.Message, second) {
					t.Fatalf("%d %s %s", code, output, stderr)
				}
				baselineHistoryUnchanged(t, h, before, refs, index)
				return
			}
			if code != 0 || stderr != "" {
				t.Fatalf("%d %s %s", code, output, stderr)
			}
			after := historySnapshot(t, h.repo)
			if !bytes.Equal(before[second], after["docs/references/"+file]) {
				t.Fatal("second copy differs")
			}
			record, err := spec.ParseArchiveRecord(after[baselineCitationFolder+".md"])
			if err != nil || !reflect.DeepEqual(record.Promoted, []string{baselineCitationDestination, "docs/references/other.md"}) {
				t.Fatalf("record=%+v err=%v", record, err)
			}
		})
	}
}
