package cli

// Suite: history sanitize command.
// Invariant: plans and refusals preserve bytes; a tagged clean batch touches only selected units.
// Boundary IN: CLI dispatch, disposable gittest repositories, temporary homes and fake HTTP.
// Boundary OUT: conversion/reduction internals and provider services.

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"roundfix/internal/gittest"
	"roundfix/internal/judge"
	"roundfix/internal/spec"
)

type historyTransport func(*http.Request) (*http.Response, error)

func (f historyTransport) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

type historyFixture struct {
	repo, home, revision string
	env                  commandEnvironment
}

func newHistoryFixture(t *testing.T, full bool) historyFixture {
	t.Helper()
	h := historyFixture{repo: t.TempDir(), home: t.TempDir()}
	gittest.InitRepo(t, h.repo, "-b", "main")
	if err := os.MkdirAll(filepath.Join(h.repo, "docs/specs"), 0755); err != nil {
		t.Fatal(err)
	}
	historyWrite(t, filepath.Join(h.repo, "README.md"), "# Fixture\n\nCites docs/history/specs/aaa/references/lesson.md\n")
	if full {
		for _, slug := range []string{"aaa", "bbb", "ccc"} {
			folder := filepath.Join(h.repo, "docs/history/specs", slug)
			historyWrite(t, filepath.Join(folder, "_prd.md"), "---\nspec: "+slug+"\ncreated: 2026-09-01\n---\n\n# "+slug+"\n\nA preserved outcome.\n")
			historyWrite(t, filepath.Join(folder, "references/lesson.md"), "Reusable knowledge.\n")
			historyWrite(t, filepath.Join(folder, "qa/evidence/out.txt"), "Evidence\n")
		}
		for _, kind := range []string{"findings", "backlog"} {
			historyWrite(t, filepath.Join(h.repo, "docs/history", kind, "entry.md"), "---\nstatus: closed\nreason: done\n---\n\n# Retired\n\nFirst paragraph.\n\n## Detail\n\nFull detail.\n")
		}
		historyWrite(t, filepath.Join(h.repo, "docs/history/reviews/pr-1/review.md"), "Review\n")
		historyWrite(t, filepath.Join(h.repo, "docs/history/handoffs/handoff.md"), "Handoff\n")
		historyWrite(t, filepath.Join(h.repo, "docs/history/adr/retired.md"), "Keep whole\n")
		historyWrite(t, filepath.Join(h.repo, "docs/history/specs/record.md"), "Existing record\n")
	}
	h.revision = historyCommit(t, h.repo, "Archive fixtures (#42)")
	h.env = commandEnvironment{workDir: h.repo, homeDir: h.home, environ: []string{"PATH=" + os.Getenv("PATH")}}
	h.env.dependencies.judgeTransport = historyTransport(func(*http.Request) (*http.Response, error) {
		t.Fatal("unexpected provider request")
		return nil, fmt.Errorf("unexpected request")
	})
	h.env.dependencies.judgeNow = func() time.Time { return time.Date(2026, 10, 7, 0, 0, 0, 0, time.UTC) }
	return h
}
func historyWrite(t *testing.T, name, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(name), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(name, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
}
func historySnapshot(t *testing.T, root string) map[string][]byte {
	t.Helper()
	result := make(map[string][]byte)
	err := filepath.WalkDir(root, func(name string, entry fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if name == filepath.Join(root, ".git") {
			if entry.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}
		if entry.IsDir() {
			return nil
		}
		data, err := os.ReadFile(name)
		if err != nil {
			return err
		}
		relative, err := filepath.Rel(root, name)
		if err != nil {
			return err
		}
		result[filepath.ToSlash(relative)] = data
		return nil
	})
	if err != nil {
		t.Fatalf("snapshot working tree %s: %v", root, err)
	}
	return result
}
func historyCommit(t *testing.T, repo, subject string) string {
	t.Helper()
	gittest.Run(t, repo, "add", "--all")
	gittest.Run(t, repo, "commit", "-m", subject)
	return strings.TrimSpace(gittest.Run(t, repo, "rev-parse", "HEAD"))
}
func historyTag(t *testing.T, h historyFixture) {
	t.Helper()
	gittest.Run(t, h.repo, "tag", "-a", "history-full", "-m", "full")
}
func historyRun(t *testing.T, h historyFixture, args ...string) (int, string, string) {
	t.Helper()
	var out, err bytes.Buffer
	code := runWithContext(context.Background(), append([]string{"history"}, args...), &out, &err, h.env)
	return code, out.String(), err.String()
}
func historyUnchanged(t *testing.T, h historyFixture, tree map[string][]byte, status string) {
	t.Helper()
	if got := historySnapshot(t, h.repo); !reflect.DeepEqual(got, tree) {
		t.Fatal("repository file bytes changed")
	}
	if got := gittest.Run(t, h.repo, "status", "--porcelain", "--untracked-files=all"); got != status {
		t.Fatalf("status changed: %q -> %q", status, got)
	}
}
func historyRefusal(t *testing.T, h historyFixture, args ...string) {
	t.Helper()
	before := historySnapshot(t, h.repo)
	status := gittest.Run(t, h.repo, "status", "--porcelain", "--untracked-files=all")
	code, out, err := historyRun(t, h, args...)
	if code != 2 || out != "" || !strings.Contains(err, "Preflight failed") {
		t.Fatalf("code=%d stdout=%s stderr=%s", code, out, err)
	}
	historyUnchanged(t, h, before, status)
}

func TestHistorySanitizePlanWritesNothing(t *testing.T) {
	t.Parallel()
	h := newHistoryFixture(t, true)
	before := historySnapshot(t, h.repo)
	status := gittest.Run(t, h.repo, "status", "--porcelain", "--untracked-files=all")
	code, out, err := historyRun(t, h, "sanitize")
	if code != 0 || err != "" {
		t.Fatalf("code=%d stderr=%s", code, err)
	}
	var expected strings.Builder
	var conversions []spec.LegacyConversion
	var files int
	var size int64
	for _, slug := range []string{"aaa", "bbb", "ccc"} {
		c, e := spec.PlanLegacyConversion(spec.LegacyConversionRequest{RepositoryRoot: h.repo, ArchiveRoot: filepath.Join(h.repo, "docs/history/specs"), Slug: slug, SourceRevision: h.revision, Delivery: spec.LegacyDelivery{Commit: h.revision, PullRequest: "42", Date: strings.TrimSpace(gittest.Run(t, h.repo, "show", "-s", "--format=%as", "HEAD"))}})
		if e != nil {
			t.Fatal(e)
		}
		conversions = append(conversions, c)
		files += len(c.Files)
		size += c.Bytes
	}
	kinds, e := spec.PlanHistoryKinds(h.repo, h.revision)
	if e != nil {
		t.Fatal(e)
	}
	for _, k := range kinds {
		files += len(k.Files)
		size += k.BytesBefore
	}
	fmt.Fprintf(&expected, "history sanitize plan: 7 unit(s) pending; %d file(s) (%d bytes) leave docs/history\n", files, size)
	for _, c := range conversions {
		fmt.Fprintf(&expected, "folder %s: removes %d file(s) (%d bytes) and writes %s (%d bytes, no-qa, delivery %.12s #42)\ncandidate %s/references/lesson.md 20 bytes\n", c.Folder, len(c.Files), c.Bytes, c.RecordPath, len(c.Rendered), h.revision, c.Folder)
	}
	for _, k := range kinds {
		if k.Action == "reduce" {
			fmt.Fprintf(&expected, "%s: reduces %d file(s) from %d to %d bytes\n", k.Kind, len(k.Files), k.BytesBefore, k.BytesAfter)
		} else {
			fmt.Fprintf(&expected, "%s: removes %d file(s) (%d bytes)\n", k.Kind, len(k.Files), k.BytesBefore)
		}
	}
	expected.WriteString("cites README.md:3 names docs/history/specs/aaa/references/lesson.md\napply with: roundfix history sanitize --apply --batch <n> (needs the annotated tag history-full at or before HEAD)\n")
	if out != expected.String() {
		t.Fatalf("plan:\n%s\nwant:\n%s", out, expected.String())
	}
	historyUnchanged(t, h, before, status)
}
func TestHistorySanitizeAppliesTheNextBatch(t *testing.T) {
	t.Parallel()
	h := newHistoryFixture(t, true)
	historyTag(t, h)
	original := snapshotDirectoryFiles(t, filepath.Join(h.repo, "docs/history/specs/ccc"))
	before := historySnapshot(t, h.repo)
	code, out, err := historyRun(t, h, "sanitize", "--apply", "--batch", "2")
	if code != 0 || err != "" || !strings.Contains(out, "applied 2 unit(s): wrote 2 Archive Record(s)") || !strings.Contains(out, "5 unit(s) remain") {
		t.Fatalf("%d %s %s", code, out, err)
	}
	for _, slug := range []string{"aaa", "bbb"} {
		assertPathMissing(t, filepath.Join(h.repo, "docs/history/specs", slug))
		data, e := os.ReadFile(filepath.Join(h.repo, "docs/history/specs", slug+".md"))
		if e != nil {
			t.Fatal(e)
		}
		record, e := spec.ParseArchiveRecord(data)
		if e != nil || record.SourceRevision != h.revision {
			t.Fatalf("record=%+v err=%v", record, e)
		}
		for _, name := range []string{"_prd.md", "references/lesson.md", "qa/evidence/out.txt"} {
			got := gittest.Run(t, h.repo, "show", "history-full:docs/history/specs/"+slug+"/"+name)
			want := before["docs/history/specs/"+slug+"/"+name]
			if got != string(want) {
				t.Fatalf("tag recovery %s", name)
			}
		}
	}
	if !reflect.DeepEqual(original, snapshotDirectoryFiles(t, filepath.Join(h.repo, "docs/history/specs/ccc"))) {
		t.Fatal("third folder changed")
	}
	revision := historyCommit(t, h.repo, "Batch one")
	code, out, err = historyRun(t, h, "sanitize", "--apply", "--batch", "1")
	if code != 0 || !strings.Contains(out, "4 unit(s) remain") {
		t.Fatalf("%d %s %s", code, out, err)
	}
	data, e := os.ReadFile(filepath.Join(h.repo, "docs/history/specs/ccc.md"))
	if e != nil {
		t.Fatal(e)
	}
	record, e := spec.ParseArchiveRecord(data)
	if e != nil || record.SourceRevision != revision {
		t.Fatalf("record=%+v err=%v", record, e)
	}
	historyCommit(t, h.repo, "Batch two")
	code, out, err = historyRun(t, h, "sanitize", "--apply", "--batch", "4")
	if code != 0 || !strings.Contains(out, "reduced 2 file(s), removed 2 file(s)") {
		t.Fatalf("%d %s %s", code, out, err)
	}
	for _, kind := range []string{"findings", "backlog"} {
		data, e := os.ReadFile(filepath.Join(h.repo, "docs/history", kind, "entry.md"))
		if e != nil || !spec.IsReducedHistoryEntry(data) {
			t.Fatalf("kind %s not reduced: %v", kind, e)
		}
	}
	assertPathMissing(t, filepath.Join(h.repo, "docs/history/reviews"))
	assertPathMissing(t, filepath.Join(h.repo, "docs/history/handoffs"))
	if data, e := os.ReadFile(filepath.Join(h.repo, "docs/history/adr/retired.md")); e != nil || string(data) != "Keep whole\n" {
		t.Fatal("ADR changed")
	}
}
func TestHistorySanitizeRefusesWithoutAnAnnotatedAncestorTag(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"missing", "lightweight", "not ancestor"} {
		t.Run(kind, func(t *testing.T) {
			h := newHistoryFixture(t, true)
			switch kind {
			case "lightweight":
				gittest.Run(t, h.repo, "tag", "history-full")
			case "not ancestor":
				gittest.Run(t, h.repo, "checkout", "-b", "future")
				historyWrite(t, filepath.Join(h.repo, "future.md"), "future")
				historyCommit(t, h.repo, "Future")
				historyTag(t, h)
				gittest.Run(t, h.repo, "checkout", "main")
			}
			historyRefusal(t, h, "sanitize", "--apply", "--batch", "2")
		})
	}
}
func TestHistorySanitizeRefusesATagWithoutTheBatchPaths(t *testing.T) {
	t.Parallel()
	h := newHistoryFixture(t, false)
	historyTag(t, h)
	historyWrite(t, filepath.Join(h.repo, "docs/history/reviews/pr-1/r.md"), "review")
	historyCommit(t, h.repo, "Later history")
	historyRefusal(t, h, "sanitize", "--apply", "--batch", "1")
}
func TestHistorySanitizeRefusesADirtyTree(t *testing.T) {
	t.Parallel()
	for _, kind := range []string{"tracked", "untracked", "staged"} {
		t.Run(kind, func(t *testing.T) {
			h := newHistoryFixture(t, true)
			historyTag(t, h)
			name := "README.md"
			if kind == "untracked" {
				name = "new.md"
			}
			historyWrite(t, filepath.Join(h.repo, name), "dirty")
			if kind == "staged" {
				gittest.Run(t, h.repo, "add", name)
			}
			historyRefusal(t, h, "sanitize", "--apply", "--batch", "2")
		})
	}
}
func TestHistorySanitizeAdviseFailsOpenWithoutAKey(t *testing.T) {
	t.Parallel()
	h := newHistoryFixture(t, true)
	before := historySnapshot(t, h.repo)
	status := gittest.Run(t, h.repo, "status", "--porcelain", "--untracked-files=all")
	q, e := judge.Load()
	if e != nil {
		t.Fatal(e)
	}
	code, out, err := historyRun(t, h, "sanitize", "--batch", "1", "--advise")
	if code != 0 || err != "" || !strings.HasPrefix(out, "history sanitize plan: 1 unit(s) pending;") || !strings.Contains(out, "candidate docs/history/specs/aaa/references/lesson.md 20 bytes: no advice ("+q.KeyVariables()[0]+" is not set)") {
		t.Fatalf("%d %s %s", code, out, err)
	}
	historyUnchanged(t, h, before, status)
}
func TestHistorySanitizeAdviseWithAFakeJudge(t *testing.T) {
	t.Parallel()
	for _, choice := range []string{"reusable_knowledge", "repository_record", "transient_evidence"} {
		t.Run(choice, func(t *testing.T) {
			h := newHistoryFixture(t, true)
			before := historySnapshot(t, h.repo)
			status := gittest.Run(t, h.repo, "status", "--porcelain", "--untracked-files=all")
			q, e := judge.Load()
			if e != nil {
				t.Fatal(e)
			}
			h.env.environ = append(h.env.environ, "ROUNDFIX_OPENROUTER_API_KEY=fake-key")
			calls := 0
			h.env.dependencies.judgeTransport = historyTransport(func(r *http.Request) (*http.Response, error) {
				calls++
				body, e := io.ReadAll(r.Body)
				if e != nil || !bytes.Contains(body, []byte("references/lesson.md")) {
					t.Fatalf("request=%s err=%v", body, e)
				}
				probabilities := map[string]float64{}
				for key := range q.Archive.Question.Criteria {
					probabilities[key] = 0
				}
				probabilities[choice] = 1
				distribution, e := json.Marshal(probabilities)
				if e != nil {
					t.Fatal(e)
				}
				answer := fmt.Sprintf(`{"model":"jev-1.13.0","answers":{"%s":{"type":"choice","choice":"%s","probabilities":%s,"confidence":1}},"usage":{"input_tokens":10,"output_tokens":2}}`, q.Archive.QuestionID, choice, distribution)
				return &http.Response{StatusCode: 200, Header: make(http.Header), Body: io.NopCloser(strings.NewReader(answer))}, nil
			})
			code, out, err := historyRun(t, h, "sanitize", "--batch", "2", "--advise")
			if code != 0 || err != "" || calls != 2 || strings.Count(out, ": "+choice+" (p=1.00)") != 2 {
				t.Fatalf("%d calls=%d %s %s", code, calls, out, err)
			}
			historyUnchanged(t, h, before, status)
			data, e := os.ReadFile(filepath.Join(h.home, ".roundfix/judge/2026-10.jsonl"))
			if e != nil || !bytes.Contains(data, []byte(`"model":"jev-1.13.0"`)) {
				t.Fatalf("judge log=%s err=%v", data, e)
			}
		})
	}
}
func TestHistorySanitizePromotesIntoReferences(t *testing.T) {
	t.Parallel()
	h := newHistoryFixture(t, true)
	historyTag(t, h)
	historyRefusal(t, h, "sanitize", "--apply", "--batch", "1", "--promote", "docs/history/specs/bbb/references/lesson.md")
	historyRefusal(t, h, "sanitize", "--apply", "--batch", "1", "--promote", "docs/history/specs/aaa/_prd.md")
	historyRefusal(t, h, "sanitize", "--apply", "--batch", "1", "--promote", filepath.Join(h.repo, "docs/history/specs/aaa/references/lesson.md"))
	historyRefusal(t, h, "sanitize", "--apply", "--batch", "1", "--promote", "docs/history/specs/aaa/../bbb/references/lesson.md")
	historyRefusal(t, h, "sanitize", "--apply", "--batch", "2", "--promote", "docs/history/specs/aaa/references/lesson.md", "--promote", "docs/history/specs/bbb/references/lesson.md")
	code, out, err := historyRun(t, h, "sanitize", "--apply", "--batch", "1", "--promote", "docs/history/specs/aaa/references/lesson.md")
	if code != 0 || !strings.Contains(out, "promoted 1 file(s)") {
		t.Fatalf("%d %s %s", code, out, err)
	}
	data, e := os.ReadFile(filepath.Join(h.repo, "docs/references/lesson.md"))
	if e != nil || string(data) != "Reusable knowledge.\n" {
		t.Fatalf("promotion: %s %v", data, e)
	}
	data, e = os.ReadFile(filepath.Join(h.repo, "docs/history/specs/aaa.md"))
	if e != nil {
		t.Fatal(e)
	}
	record, e := spec.ParseArchiveRecord(data)
	if e != nil || !reflect.DeepEqual(record.Promoted, []string{"docs/references/lesson.md"}) {
		t.Fatalf("record=%+v err=%v", record, e)
	}
}
func TestHistorySanitizeUsageErrors(t *testing.T) {
	t.Parallel()
	for _, args := range [][]string{{"sanitize", "--apply"}, {"sanitize", "--advise"}, {"sanitize", "--advise", "--apply", "--batch", "1"}, {"sanitize", "--promote", "x"}, {"sanitize", "--batch", "0"}, {"sanitize", "--batch", "-1"}, {"sanitize", "--batch", "x"}, {"sanitize", "--batch"}, {"sanitize", "--unknown"}, {"sanitize", "extra"}, {"unknown"}} {
		t.Run(strings.Join(args, " "), func(t *testing.T) { h := newHistoryFixture(t, true); historyRefusal(t, h, args...) })
	}
	h := newHistoryFixture(t, false)
	for _, args := range [][]string{nil, {"--help"}, {"sanitize", "--help"}} {
		code, out, err := historyRun(t, h, args...)
		if code != 0 || out != historyUsage || err != "" {
			t.Fatalf("help: %d %s %s", code, out, err)
		}
	}
}
func TestHistorySanitizeRefusesAnExternalSpecRoot(t *testing.T) {
	t.Parallel()
	h := newHistoryFixture(t, true)
	external := t.TempDir()
	gittest.InitRepo(t, external, "-b", "main")
	historyWrite(t, filepath.Join(h.repo, ".roundfixrc.yml"), "specs:\n  root: "+external+"\n")
	historyCommit(t, h.repo, "External root")
	historyRefusal(t, h, "sanitize")
}
func TestHistorySanitizeWithNothingPending(t *testing.T) {
	t.Parallel()
	h := newHistoryFixture(t, false)
	before := historySnapshot(t, h.repo)
	for _, args := range [][]string{{"sanitize"}, {"sanitize", "--apply", "--batch", "2"}} {
		code, out, err := historyRun(t, h, args...)
		if code != 0 || out != "history sanitize plan: 0 unit(s) pending; nothing pending\n" || err != "" {
			t.Fatalf("%d %s %s", code, out, err)
		}
		historyUnchanged(t, h, before, "")
	}
}

func TestHistorySanitizePreflightsWholeBatch(t *testing.T) {
	t.Parallel()
	h := newHistoryFixture(t, true)
	historyWrite(t, filepath.Join(h.repo, "docs/history/specs/bbb/_prd.md"), "invalid PRD")
	historyCommit(t, h.repo, "Malformed second unit")
	historyTag(t, h)

	code, out, err := historyRun(t, h, "sanitize", "--apply", "--batch", "1")
	if code != 0 || !strings.Contains(out, "wrote 1 Archive Record(s)") {
		t.Fatalf("unselected invalid folder blocked batch: %d %s %s", code, out, err)
	}
	h = newHistoryFixture(t, true)
	historyWrite(t, filepath.Join(h.repo, "docs/history/specs/bbb/_prd.md"), "invalid PRD")
	historyCommit(t, h.repo, "Malformed second unit")
	historyTag(t, h)
	before := snapshotDirectoryFiles(t, filepath.Join(h.repo, "docs/history/specs/bbb"))
	code, out, err = historyRun(t, h, "sanitize", "--apply", "--batch", "2")
	if code != 0 || err != "" || !strings.Contains(out, "refused docs/history/specs/bbb: ") {
		t.Fatalf("%d %s %s", code, out, err)
	}
	for _, slug := range []string{"aaa", "ccc"} {
		if _, e := os.Stat(filepath.Join(h.repo, "docs/history/specs", slug+".md")); e != nil {
			t.Fatal(e)
		}
	}
	if !reflect.DeepEqual(before, snapshotDirectoryFiles(t, filepath.Join(h.repo, "docs/history/specs/bbb"))) {
		t.Fatal("refused folder changed")
	}

}
func TestHistorySanitizeFolderCitationsAndKindNamedSlugs(t *testing.T) {
	t.Parallel()
	h := newHistoryFixture(t, false)
	historyWrite(t, filepath.Join(h.repo, "docs/history/specs/findings/_prd.md"), "---\nspec: findings\ncreated: 2026-09-01\n---\n\n# Findings\n\nA result.\n")
	historyWrite(t, filepath.Join(h.repo, "docs/history/findings/entry.md"), "---\nstatus: closed\n---\n\n# Retired\n\nResult.\n")
	historyWrite(t, filepath.Join(h.repo, "README.md"), "# Fixture\n\nFolder docs/history/specs/findings/\n")
	historyCommit(t, h.repo, "History")
	code, out, err := historyRun(t, h, "sanitize")
	if code != 0 || !strings.Contains(out, "cites README.md:3 names docs/history/specs/findings/") || !strings.Contains(out, "findings: reduces 1 file(s)") {
		t.Fatalf("%d %s %s", code, out, err)
	}
	historyTag(t, h)
	code, out, err = historyRun(t, h, "sanitize", "--apply", "--batch", "9")
	if code != 0 || !strings.Contains(out, "wrote 1 Archive Record(s), reduced 1 file(s)") {
		t.Fatalf("%d %s %s", code, out, err)
	}
}
func TestHistorySanitizeAdviceIsAdvisory(t *testing.T) {
	t.Parallel()
	for _, mode := range []string{"service failure", "monthly ceiling"} {
		t.Run(mode, func(t *testing.T) {
			h := newHistoryFixture(t, true)
			q, e := judge.Load()
			if e != nil {
				t.Fatal(e)
			}
			h.env.environ = append(h.env.environ, q.KeyVariables()[0]+"=fake")
			if mode == "monthly ceiling" {
				historyWrite(t, filepath.Join(h.home, ".roundfix/config.yml"), "jev:\n  monthly_ceiling_usd: 1\n")
				historyWrite(t, filepath.Join(h.home, ".roundfix/judge/2026-10.jsonl"), "{\"cost_usd\":1}\n")
			} else {
				h.env.dependencies.judgeTransport = historyTransport(func(*http.Request) (*http.Response, error) { return nil, fmt.Errorf("offline fixture") })
			}
			before := historySnapshot(t, h.repo)
			status := gittest.Run(t, h.repo, "status", "--porcelain", "--untracked-files=all")
			code, out, err := historyRun(t, h, "sanitize", "--batch", "1", "--advise")
			reason := "service unavailable"
			if mode == "monthly ceiling" {
				reason = "monthly ceiling reached"
			}
			if code != 0 || err != "" || !strings.Contains(out, reason) {
				t.Fatalf("%d %s %s", code, out, err)
			}
			historyUnchanged(t, h, before, status)
		})
	}
}

func TestHistorySanitizeAdviceWithoutAKnownDelivery(t *testing.T) {
	t.Parallel()
	h := newHistoryFixture(t, false)
	prd := "---\nspec: demo\ncreated: 2026-09-01\n---\n\n# Demo\n\nOutcome.\n"
	historyWrite(t, filepath.Join(h.repo, "docs/old/demo/_prd.md"), prd)
	historyCommit(t, h.repo, "Old archive")
	if err := os.RemoveAll(filepath.Join(h.repo, "docs/old")); err != nil {
		t.Fatal(err)
	}
	historyWrite(t, filepath.Join(h.repo, "docs/history/specs/demo/_prd.md"), prd)
	historyWrite(t, filepath.Join(h.repo, "docs/history/specs/demo/references/lesson.md"), "Lesson\n")
	historyCommit(t, h.repo, "Relocate archived folder")
	before := historySnapshot(t, h.repo)
	code, out, err := historyRun(t, h, "sanitize", "--batch", "1", "--advise")
	if code != 0 || err != "" || !strings.Contains(out, "no-qa, delivery unknown)") || !strings.Contains(out, "candidate docs/history/specs/demo/references/lesson.md 7 bytes: no advice (") {
		t.Fatalf("%d %s %s", code, out, err)
	}
	historyUnchanged(t, h, before, "")
}

func TestHistorySnapshotSkipsGitMetadataBeforeReading(t *testing.T) {
	t.Parallel()
	h := newHistoryFixture(t, false)
	// Like the fsmonitor socket, this metadata node cannot be read as a file.
	// A dangling symlink exercises the boundary without requiring socket privileges.
	metadata := filepath.Join(h.repo, ".git/unreadable-metadata")
	if err := os.Symlink("missing-metadata-target", metadata); err != nil {
		t.Fatal(err)
	}
	if _, err := os.ReadFile(metadata); err == nil {
		t.Fatal("regression fixture must be unreadable")
	}
	got := historySnapshot(t, h.repo)
	want := map[string][]byte{"README.md": []byte("# Fixture\n\nCites docs/history/specs/aaa/references/lesson.md\n")}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("snapshot=%v want=%v", got, want)
	}
}
