// Suite: review scope and diagnostics through disposable Git repositories.
// Boundary IN: committed candidate paths and locks, public review command.
// Boundary OUT: captured provider prompts and persisted checkout review records.
package cli

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

	"roundfix/internal/agent"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
)

func TestReviewOmitsQAEvidenceAndUpstreamSkills(t *testing.T) {
	t.Parallel()
	runner := &reviewCommandRunner{results: []reviewCommandRunResult{{result: agent.ExecuteResult{Message: "No findings."}}}}
	fixture := newReviewCommandFixture(t, "codex", runner)
	commitLineageFile(t, fixture, "skills-lock.json", `{"skills":{"external":{}}}`)
	for _, path := range []string{"docs/specs/example/qa/evidence/raw.txt", "docs/history/specs/old/qa/evidence/raw.txt", ".agents/skills/external/SKILL.md"} {
		if err := os.MkdirAll(filepath.Dir(filepath.Join(fixture.repository, path)), 0o755); err != nil {
			t.Fatal(err)
		}
		commitLineageFile(t, fixture, path, strings.Repeat("excluded payload\n", reviewDiffBound/16))
	}
	commitLineageFile(t, fixture, "source.go", "package example\n")
	code, record, stderr := runLineageReview(t, fixture)
	if code != exitOK {
		t.Fatalf("exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	if strings.Contains(runner.request.Prompt, "excluded payload") {
		t.Fatal("omitted payload reached reviewer")
	}
	if !strings.Contains(runner.request.Prompt, "diff --git a/source.go b/source.go") {
		t.Fatal("source missing from reviewed diff")
	}
	assertReviewScope(t, fixture, runner, record, []reviewOmittedPath{{Path: ".agents/skills/external/SKILL.md", Reason: "upstream-skill"}, {Path: "docs/history/specs/old/qa/evidence/raw.txt", Reason: "qa-evidence"}, {Path: "docs/specs/example/qa/evidence/raw.txt", Reason: "qa-evidence"}}, "--- END CANDIDATE DIFF ---")
}

func assertReviewScope(t *testing.T, fixture reviewCommandFixture, runner *reviewCommandRunner, record reviewRecord, want []reviewOmittedPath, marker string) {
	t.Helper()
	if !reflect.DeepEqual(record.OmittedPaths, want) {
		t.Fatalf("omitted paths=%+v want=%+v", record.OmittedPaths, want)
	}
	_, after, ok := strings.Cut(runner.request.Prompt, marker+"\n")
	if !ok {
		t.Fatalf("missing diff end marker %q", marker)
	}
	qa, upstream := 0, 0
	for _, path := range want {
		if path.Reason == "qa-evidence" {
			qa++
		} else {
			upstream++
		}
	}
	summary := fmt.Sprintf("Omitted from this diff (not reviewed): %d path(s) of QA evidence, %d of upstream-managed skills.\n", qa, upstream)
	if !strings.HasPrefix(after, summary) {
		t.Fatalf("missing summary %q", summary)
	}
	last := -1
	for _, path := range want {
		index := strings.Index(after, path.Path+" ("+path.Reason+")")
		if index <= last {
			t.Fatalf("missing or unsorted omission %s in %q", path.Path, after)
		}
		last = index
	}
	startMarker := "--- BEGIN CANDIDATE DIFF ---\n"
	if marker == "--- END ROUND 2 DELTA DIFF ---" {
		startMarker = "--- BEGIN ROUND 2 DELTA DIFF ---\n"
	}
	_, body, _ := strings.Cut(runner.request.Prompt, startMarker)
	body, _, _ = strings.Cut(body, marker+"\n")
	for _, path := range want {
		if strings.Contains(body, "diff --git a/"+path.Path+" b/") {
			t.Fatalf("omitted path in diff: %s", path.Path)
		}
	}
	// ExecGitRunner removes the Git output's trailing newline; the prompt
	// adds a separator before the end marker, which is not part of the diff.
	body = strings.TrimSuffix(body, "\n")
	if record.DiffBytes != len(body) {
		t.Fatalf("diffBytes=%d prompt diff=%d", record.DiffBytes, len(body))
	}
	persisted, err := readReviewRecord(filepath.Join(reviewCheckoutDir(fixture.artifactDir, fixture.repository), reviewRecordFileName))
	if err != nil {
		t.Fatal(err)
	}
	if persisted.DiffBytes != record.DiffBytes || !reflect.DeepEqual(persisted.OmittedPaths, want) {
		t.Fatalf("persisted scope=%+v", persisted)
	}
}

func commitScopeFile(t *testing.T, fixture reviewCommandFixture, path, content string) string {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(filepath.Join(fixture.repository, path)), 0o755); err != nil {
		t.Fatal(err)
	}
	return commitLineageFile(t, fixture, path, content)
}

func TestReviewOmitsASkillTheLockDroppedAtTheHead(t *testing.T) {
	t.Parallel()
	runner := &reviewCommandRunner{results: []reviewCommandRunResult{{result: agent.ExecuteResult{Message: "No findings."}}}}
	fixture := newReviewCommandFixture(t, "codex", runner)
	base := commitScopeFile(t, fixture, "skills-lock.json", `{"skills":{"dropped":{}}}`)
	fixture.baseCommit = base
	gittest.Run(t, fixture.repository, "branch", "-f", "main", base)
	commitScopeFile(t, fixture, ".agents/skills/dropped/SKILL.md", "dropped payload\n")
	commitScopeFile(t, fixture, "skills-lock.json", `{"skills":{}}`)
	code, record, stderr := runLineageReview(t, fixture)
	if code != exitOK || strings.Contains(runner.request.Prompt, "dropped payload") {
		t.Fatalf("exit=%d record=%+v stderr=%q prompt=%q", code, record, stderr, runner.request.Prompt)
	}
	assertReviewScope(t, fixture, runner, record, []reviewOmittedPath{{Path: ".agents/skills/dropped/SKILL.md", Reason: "upstream-skill"}}, "--- END CANDIDATE DIFF ---")
}

func TestReviewKeepsTheQAReportInTheDiff(t *testing.T) {
	t.Parallel()
	runner := &reviewCommandRunner{results: []reviewCommandRunResult{{result: agent.ExecuteResult{Message: "No findings."}}}}
	fixture := newReviewCommandFixture(t, "codex", runner)
	report := "docs/specs/example/qa/qa-report-2026-10-02.md"
	commitScopeFile(t, fixture, report, "QA Report stays reviewed\n")
	// Without a lock, locally owned skills and lookalike evidence paths stay reviewed.
	commitScopeFile(t, fixture, ".agents/skills/local/SKILL.md", "local skill stays reviewed\n")
	commitScopeFile(t, fixture, "docs/specs/example/qa/evidence-lookalike/file", "lookalike stays reviewed\n")
	code, record, stderr := runLineageReview(t, fixture)
	if code != exitOK {
		t.Fatalf("exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	for _, payload := range []string{"+QA Report stays reviewed", "+local skill stays reviewed", "+lookalike stays reviewed"} {
		if !strings.Contains(runner.request.Prompt, payload) {
			t.Fatalf("missing %q", payload)
		}
	}
	assertReviewScope(t, fixture, runner, record, []reviewOmittedPath{}, "--- END CANDIDATE DIFF ---")
}

func TestReviewBlocksOnAMalformedSkillsLock(t *testing.T) {
	t.Parallel()
	for _, lock := range []string{"{", `{"skills":[]}`, `{"skills":null}`, `null`} {
		t.Run(lock, func(t *testing.T) {
			runner := &reviewCommandRunner{}
			fixture := newReviewCommandFixture(t, "codex", runner)
			commitScopeFile(t, fixture, "skills-lock.json", lock)
			code, record, stderr := runLineageReview(t, fixture)
			assertBlockedReviewCommand(t, code, record, stderr, "review scope: read skills-lock.json:")
			if runner.probeCalls != 0 || runner.prepareCalls != 0 || runner.preparedCalls != 0 {
				t.Fatalf("provider reached: %+v", runner)
			}
		})
	}
}

func TestReviewBlocksAboveTheBoundWithoutAProviderCall(t *testing.T) {
	t.Parallel()
	for _, round := range []int{1, 2} {
		t.Run(fmt.Sprint(round), func(t *testing.T) {
			runner := &reviewCommandRunner{results: []reviewCommandRunResult{{result: agent.ExecuteResult{Message: sessionReviewFinding}}}}
			fixture := newReviewCommandFixture(t, "codex", runner)
			writeReviewCommandProfileConfig(t, fixture.repository, "codex", fixture.artifactDir, "codex", "codex")
			if round == 2 {
				runLineageReview(t, fixture)
			}
			probes, prepares, prompts := runner.probeCalls, runner.prepareCalls, runner.preparedCalls
			commitScopeFile(t, fixture, "huge.txt", strings.Repeat("x", reviewDiffBound)+"\n")
			code, record, stderr := runLineageReview(t, fixture)
			want := fmt.Sprintf("review diff too large: %d bytes after omitting 0 path(s) exceeds the review bound of %d bytes", record.DiffBytes, reviewDiffBound)
			assertBlockedReviewCommand(t, code, record, stderr, want)
			if record.Reason != want || record.DiffBytes <= reviewDiffBound {
				t.Fatalf("record=%+v", record)
			}
			if runner.probeCalls != probes || runner.prepareCalls != prepares || runner.preparedCalls != prompts {
				t.Fatalf("provider reached: %+v", runner)
			}
		})
	}
}

func TestReviewRecordsTheRuntimeStderrTail(t *testing.T) {
	t.Parallel()
	for _, width := range []int{10, 200} {
		t.Run(fmt.Sprint(width), func(t *testing.T) {
			lines := make([]string, 30)
			for i := range lines {
				lines[i] = fmt.Sprintf("line %02d %s", i+1, strings.Repeat("x", width))
			}
			stderrText := strings.Join(lines, "\n") + "\n"
			failure := fmt.Errorf("wrapped: %w", &agent.BatchFailureError{Stderr: stderrText, Err: errors.New("runtime stopped")})
			runner := &reviewCommandRunner{results: []reviewCommandRunResult{{err: failure}}}
			fixture := newReviewCommandFixture(t, "codex", runner)
			code, record, stderr := runLineageReview(t, fixture)
			assertBlockedReviewCommand(t, code, record, stderr, "review runtime failure: "+failure.Error())
			want := strings.Join(lines[20:], "\n")
			if len(want) > 1024 {
				want = want[len(want)-1024:]
			}
			if record.RuntimeStderrTail != want || len(record.RuntimeStderrTail) > 1024 || len(strings.Split(record.RuntimeStderrTail, "\n")) > 10 {
				t.Fatalf("tail=%q want=%q", record.RuntimeStderrTail, want)
			}
			persisted, err := readReviewRecord(filepath.Join(reviewCheckoutDir(fixture.artifactDir, fixture.repository), reviewRecordFileName))
			if err != nil || persisted.RuntimeStderrTail != want {
				t.Fatalf("persisted=%+v err=%v", persisted, err)
			}
		})
	}
}

func TestReviewReadsARecordWithoutTheNewFields(t *testing.T) {
	t.Parallel()
	path := filepath.Join(t.TempDir(), "old.json")
	mustWrite(t, path, `{"repository":"/repo","baseCommit":"base","headCommit":"head","provider":"codex","source":"built-in","outcome":"reviewed","specs":[],"skippedSpecs":[],"archivedSpecs":[],"specContextTruncated":false}`)
	record, err := readReviewRecord(path)
	if err != nil {
		t.Fatal(err)
	}
	if record.DiffBytes != 0 || len(record.OmittedPaths) != 0 || record.RuntimeStderrTail != "" {
		t.Fatalf("legacy record=%+v", record)
	}
	var output bytes.Buffer
	if err := writeReviewRecord(&output, record); err != nil {
		t.Fatal(err)
	}
	var fields map[string]json.RawMessage
	if err := json.Unmarshal(output.Bytes(), &fields); err != nil {
		t.Fatal(err)
	}
	if string(fields["omittedPaths"]) != "[]" || string(fields["diffBytes"]) != "0" {
		t.Fatalf("new fields=%s", output.String())
	}
}

func TestReviewScopesRoundTwoDelta(t *testing.T) {
	t.Parallel()
	runner := &reviewCommandRunner{results: []reviewCommandRunResult{{result: agent.ExecuteResult{Message: sessionReviewFinding}}, {result: agent.ExecuteResult{Message: "No findings."}}}}
	fixture := newReviewCommandFixture(t, "codex", runner)
	commitScopeFile(t, fixture, "skills-lock.json", `{"skills":{"external":{}}}`)
	runLineageReview(t, fixture)
	commitScopeFile(t, fixture, ".agents/skills/external/SKILL.md", "round two excluded payload\n")
	commitScopeFile(t, fixture, "docs/specs/example/qa/evidence/raw.txt", "round two excluded payload\n")
	commitScopeFile(t, fixture, "delta.txt", "round two reviewed payload\n")
	code, record, stderr := runLineageReview(t, fixture)
	if code != exitOK || record.Lineage.Round != 2 || strings.Contains(runner.request.Prompt, "round two excluded payload") {
		t.Fatalf("exit=%d record=%+v stderr=%q prompt=%q", code, record, stderr, runner.request.Prompt)
	}
	if !strings.Contains(runner.request.Prompt, "+round two reviewed payload") {
		t.Fatal("missing reviewed delta")
	}
	assertReviewScope(t, fixture, runner, record, []reviewOmittedPath{{Path: ".agents/skills/external/SKILL.md", Reason: "upstream-skill"}, {Path: "docs/specs/example/qa/evidence/raw.txt", Reason: "qa-evidence"}}, "--- END ROUND 2 DELTA DIFF ---")
}

func TestReviewScopesConfiguredSpecRoots(t *testing.T) {
	t.Parallel()
	runner := &reviewCommandRunner{results: []reviewCommandRunResult{{result: agent.ExecuteResult{Message: "No findings."}}}}
	fixture := newReviewCommandFixture(t, "codex", runner)
	writeReviewCommandConfigWithSpecsRoot(t, fixture.repository, "codex", fixture.artifactDir, "planning")
	for _, path := range []string{"planning/example/qa/evidence/raw.txt", "planning/_archived/old/qa/evidence/raw.txt", "docs/specs/other/qa/evidence/raw.txt"} {
		commitScopeFile(t, fixture, path, path+" payload\n")
	}
	code, record, stderr := runLineageReview(t, fixture)
	if code != exitOK {
		t.Fatalf("exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	if !strings.Contains(runner.request.Prompt, "+docs/specs/other/qa/evidence/raw.txt payload") {
		t.Fatal("unconfigured root wrongly omitted")
	}
	assertReviewScope(t, fixture, runner, record, []reviewOmittedPath{{Path: "planning/_archived/old/qa/evidence/raw.txt", Reason: "qa-evidence"}, {Path: "planning/example/qa/evidence/raw.txt", Reason: "qa-evidence"}}, "--- END CANDIDATE DIFF ---")
}

func TestReviewBlocksOnAMalformedBaseLock(t *testing.T) {
	t.Parallel()
	runner := &reviewCommandRunner{}
	fixture := newReviewCommandFixture(t, "codex", runner)
	fixture.baseCommit = commitScopeFile(t, fixture, "skills-lock.json", "{")
	commitScopeFile(t, fixture, "skills-lock.json", `{"skills":{}}`)
	code, record, stderr := runLineageReview(t, fixture)
	assertBlockedReviewCommand(t, code, record, stderr, "review scope: read skills-lock.json:")
	if runner.probeCalls != 0 || runner.prepareCalls != 0 {
		t.Fatal("malformed base lock reached provider")
	}
}

type unreadableReviewLockGit struct{ preflight.ExecGitRunner }

func (runner unreadableReviewLockGit) RunGit(ctx context.Context, root string, args ...string) (string, error) {
	if len(args) == 2 && args[0] == "show" && strings.HasSuffix(args[1], ":skills-lock.json") {
		return "", errors.New("lock object unreadable")
	}
	return runner.ExecGitRunner.RunGit(ctx, root, args...)
}

func TestReviewBlocksOnAnUnreadableSkillsLock(t *testing.T) {
	t.Parallel()
	runner := &reviewCommandRunner{}
	fixture := newReviewCommandFixture(t, "codex", runner)
	commitScopeFile(t, fixture, "skills-lock.json", `{"skills":{}}`)
	withReviewSpecGitRunner(t, unreadableReviewLockGit{})
	code, record, stderr := runLineageReview(t, fixture)
	assertBlockedReviewCommand(t, code, record, stderr, "review scope: read skills-lock.json: lock object unreadable")
	if runner.probeCalls != 0 || runner.prepareCalls != 0 || runner.preparedCalls != 0 {
		t.Fatal("unreadable lock reached provider")
	}
}

func TestReviewAdmitsADiffExactlyAtTheBound(t *testing.T) {
	t.Parallel()
	runner := &reviewCommandRunner{results: []reviewCommandRunResult{{result: agent.ExecuteResult{Message: "No findings."}}}}
	fixture := newReviewCommandFixture(t, "codex", runner)
	head := commitScopeFile(t, fixture, "bound.txt", "x\n")
	diff, err := (preflight.ExecGitRunner{}).RunGit(t.Context(), fixture.repository, "diff", "--no-ext-diff", "--no-textconv", "--no-color", fixture.baseCommit, head, "--")
	if err != nil {
		t.Fatal(err)
	}
	// Only the added line's payload grows, so Git's fixed headers stay equal.
	commitScopeFile(t, fixture, "bound.txt", strings.Repeat("x", reviewDiffBound-len(diff)+1)+"\n")
	code, record, stderr := runLineageReview(t, fixture)
	if code != exitOK || record.DiffBytes != reviewDiffBound || runner.preparedCalls != 1 {
		t.Fatalf("exit=%d diffBytes=%d stderr=%q calls=%d", code, record.DiffBytes, stderr, runner.preparedCalls)
	}
}
