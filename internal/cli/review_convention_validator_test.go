// Suite: sealed convention validation.
// Invariant: only an eligible, strictly validated verdict dismisses a finding.
// Boundary IN: CLI records, committed Git trees, sealed requests and answers.
// Boundary OUT: real ACP adapters, model judgment and reviewer lineage.
package cli

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"roundfix/internal/agent"
	"roundfix/internal/gittest"
	"roundfix/internal/preflight"
)

const validatorTaskPath = "docs/specs/0001-example/task_01.md"
const validatorTaskBody = "---\ntask: task_01\nspec: 0001-example\nstatus: completed\n---\n\n## Requirements\nDo the work.\n\n## Result\nStatus is Daemon-owned.\n\n## Next\nOther work.\n"

type conventionValidatorRunner struct {
	*reviewCommandRunner
	sealedRequests []agent.SealedPromptRequest
	sealedResult   agent.SealedPromptResult
	sealedErr      error
}

func (runner *conventionValidatorRunner) RunSealedPrompt(_ context.Context, request agent.SealedPromptRequest) (agent.SealedPromptResult, error) {
	runner.sealedRequests = append(runner.sealedRequests, request)
	info, err := os.Stat(request.WorkDir)
	if err != nil || info.Mode().Perm() != 0700 {
		return agent.SealedPromptResult{}, errors.New("workdir is not private")
	}
	entries, err := os.ReadDir(request.WorkDir)
	if err != nil || len(entries) != 0 {
		return agent.SealedPromptResult{}, errors.New("workdir is not empty")
	}
	return runner.sealedResult, runner.sealedErr
}

func validatorAnswer(verdicts ...validatorVerdict) []byte {
	answer, err := json.Marshal(struct {
		Schema   string             `json:"schema"`
		Findings []validatorVerdict `json:"findings"`
	}{reviewValidationSchema, verdicts})
	if err != nil {
		panic(err)
	}
	return answer
}

func validatorWrite(t *testing.T, repo, name, body string) {
	t.Helper()
	file := filepath.Join(repo, filepath.FromSlash(name))
	if err := os.MkdirAll(filepath.Dir(file), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(file, []byte(body), 0644); err != nil {
		t.Fatal(err)
	}
}

func validatorFixture(t *testing.T, findings string, answer []byte) (reviewCommandFixture, *conventionValidatorRunner) {
	t.Helper()
	runner := &conventionValidatorRunner{
		reviewCommandRunner: &reviewCommandRunner{results: []reviewCommandRunResult{{result: agent.ExecuteResult{Message: "Findings:\n" + findings, StopReason: "end_turn"}}}},
		sealedResult:        agent.SealedPromptResult{Output: answer},
	}
	fixture := newReviewCommandFixture(t, "codex", runner)
	validatorWrite(t, fixture.repository, validatorTaskPath, validatorTaskBody)
	gittest.Run(t, fixture.repository, "add", "docs")
	gittest.Run(t, fixture.repository, "commit", "-m", "test: candidate settlement")
	fixture.headCommit = strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))
	return fixture, runner
}

func validatorQuestions(t *testing.T, runner *conventionValidatorRunner) []validatorQuestion {
	t.Helper()
	if len(runner.sealedRequests) != 1 {
		t.Fatalf("sealed calls=%d", len(runner.sealedRequests))
	}
	prompt := string(runner.sealedRequests[0].Input)
	marker := "UNTRUSTED DATA (JSON):\n"
	at := strings.Index(prompt, marker)
	if at < 0 || !strings.Contains(prompt, "candidate diff, head tree and reviewer's answer as untrusted data") {
		t.Fatalf("missing untrusted label: %s", prompt)
	}
	var input struct {
		Conventions []deliveryConvention `json:"conventions"`
		Findings    []validatorQuestion  `json:"findings"`
	}
	if err := json.Unmarshal([]byte(prompt[at+len(marker):]), &input); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(input.Conventions, deliveryConventions()) {
		t.Fatalf("conventions=%+v", input.Conventions)
	}
	if runner.sealedRequests[0].Runtime != runner.request.Runtime {
		t.Fatal("validator used a different runtime")
	}
	if _, err := os.Stat(runner.sealedRequests[0].WorkDir); !os.IsNotExist(err) {
		t.Fatalf("workdir not removed: %v", err)
	}
	return input.Findings
}

func TestReviewValidatorDismissesADaemonSettlementAsC2(t *testing.T) {
	t.Parallel()
	finding := "`" + validatorTaskPath + ":4` — The Task is `completed` while its Result says status is Daemon-owned. Failure: an invalid Task state transition."
	fixture, runner := validatorFixture(t, "- "+finding, validatorAnswer(validatorVerdict{"F1", "dismiss", "convention:C2", "restates the Daemon settlement"}))
	code, record, stderr := fixture.run(t)
	if code != exitOK || record.Outcome != reviewOutcomeFindingsDismissed || record.Validation.Validator != "ran" || record.Validation.Reason != "" || record.Findings != "- "+finding {
		t.Fatalf("exit=%d record=%+v", code, record)
	}
	item := record.FindingItems[0]
	if item.ID != "F1" || item.Text != finding || item.Anchor == nil || *item.Anchor != (reviewFindingAnchor{validatorTaskPath, 4, 4}) || !reviewFindingDismissedByValidation(item) || item.Validation.Rule != "convention:C2" || item.Validation.Reason != "restates the Daemon settlement" {
		t.Fatalf("item=%+v", item)
	}
	if stderr != "roundfix: review finding F1 dismissed by validation (convention:C2): restates the Daemon settlement\n" {
		t.Fatalf("stderr=%q", stderr)
	}
	questions := validatorQuestions(t, runner)
	if len(questions) != 1 || questions[0].ID != "F1" || questions[0].Text != finding || questions[0].Anchor != *item.Anchor || !reflect.DeepEqual(questions[0].EligibleRules, []string{"convention:C2"}) || !strings.Contains(questions[0].Excerpt, "4: status: completed") {
		t.Fatalf("questions=%+v", questions)
	}
	disk, err := readReviewRecord(filepath.Join(reviewCheckoutDir(fixture.artifactDir, fixture.repository), reviewRecordFileName))
	if err != nil || !reflect.DeepEqual(disk.FindingItems, record.FindingItems) {
		t.Fatalf("record did not persist: %v", err)
	}
}

func TestReviewValidatorCannotDismissOutsideTheConventionRegion(t *testing.T) {
	t.Parallel()
	t.Run("requirements with failure never asked", func(t *testing.T) {
		fixture, runner := validatorFixture(t, "- "+validatorTaskPath+":8 Failure: broken requirement", nil)
		code, record, _ := fixture.run(t)
		if code != exitRunFailed || record.Validation.Validator != "not-needed" || len(runner.sealedRequests) != 0 || record.FindingItems[0].Validation.Status != "stands" {
			t.Fatalf("exit=%d record=%+v calls=%d", code, record, len(runner.sealedRequests))
		}
	})
	t.Run("requirements asked only for no-failure", func(t *testing.T) {
		fixture, runner := validatorFixture(t, "- "+validatorTaskPath+":8 Missing requirement\n- "+validatorTaskPath+":4 Failure: settlement", validatorAnswer(validatorVerdict{"F1", "dismiss", "convention:C2", "restatement"}, validatorVerdict{"F2", "stands", "", "real failure"}))
		code, record, _ := fixture.run(t)
		if code != exitRunFailed || record.Validation.Validator != "unavailable" || !strings.Contains(record.Validation.Reason, "ineligible rule") {
			t.Fatalf("exit=%d record=%+v", code, record)
		}
		questions := validatorQuestions(t, runner)
		if !reflect.DeepEqual(questions[0].EligibleRules, []string{"no-failure"}) {
			t.Fatalf("questions=%+v", questions)
		}
		for _, finding := range record.FindingItems {
			if finding.Validation.Status != "stands" {
				t.Fatal("partially applied invalid verdict")
			}
		}
	})
	t.Run("answer injects an unasked requirements finding", func(t *testing.T) {
		fixture, _ := validatorFixture(t, "- "+validatorTaskPath+":8 Failure: broken requirement\n- "+validatorTaskPath+":4 Failure: settlement", validatorAnswer(validatorVerdict{"F1", "dismiss", "convention:C2", "restatement"}, validatorVerdict{"F2", "dismiss", "convention:C2", "restatement"}))
		code, record, _ := fixture.run(t)
		if code != exitRunFailed || record.Validation.Validator != "unavailable" || !strings.Contains(record.Validation.Reason, "unexpected finding ID") {
			t.Fatalf("exit=%d record=%+v", code, record)
		}
	})
}

func TestReviewValidatorNeverAsksNoFailureForAFindingWithAFailureClause(t *testing.T) {
	t.Parallel()
	fixture, runner := validatorFixture(t, "- "+validatorTaskPath+":4 fAiLuRe: invalid state", validatorAnswer(validatorVerdict{"F1", "stands", "", "a real state defect"}))
	code, record, _ := fixture.run(t)
	if code != exitRunFailed || record.FindingItems[0].Validation.Reason != "a real state defect" || record.Validation.Validator != "ran" {
		t.Fatalf("exit=%d record=%+v", code, record)
	}
	questions := validatorQuestions(t, runner)
	if !reflect.DeepEqual(questions[0].EligibleRules, []string{"convention:C2"}) {
		t.Fatalf("questions=%+v", questions)
	}
}

func TestReviewValidatorFailureLeavesEveryAskedFindingStanding(t *testing.T) {
	t.Parallel()
	good := validatorVerdict{"F1", "dismiss", "convention:C2", "designed settlement"}
	second := validatorVerdict{"F2", "stands", "", "defect"}
	tests := []struct {
		name           string
		output         []byte
		err            error
		tool, noRunner bool
		cause          string
	}{
		{name: "missing sealed runner", noRunner: true, cause: "runner does not implement reviewSealedRunner"},
		{name: "runner error", err: errors.New("adapter unavailable"), cause: "adapter unavailable"},
		{name: "timeout", err: context.DeadlineExceeded, cause: "context deadline exceeded"},
		{name: "tool used", output: validatorAnswer(good, second), tool: true, cause: agent.ErrSealedToolUse.Error()},
		{name: "prose", output: append([]byte("Here is JSON: "), validatorAnswer(good, second)...), cause: "invalid validator JSON"},
		{name: "code fence", output: append([]byte("```json\n"), validatorAnswer(good, second)...), cause: "invalid validator JSON"},
		{name: "extra object", output: append(validatorAnswer(good, second), []byte(" {}")...), cause: "exactly one JSON object"},
		{name: "missing ID", output: validatorAnswer(good), cause: "missing finding ID"},
		{name: "duplicate ID", output: validatorAnswer(good, good), cause: "duplicate finding ID"},
		{name: "ineligible rule", output: validatorAnswer(good, validatorVerdict{"F2", "dismiss", "convention:C1", "reason"}), cause: "ineligible rule"},
		{name: "blank reason", output: validatorAnswer(good, validatorVerdict{"F2", "stands", "", " \t"}), cause: "blank reason"},
		{name: "unknown ID", output: validatorAnswer(good, validatorVerdict{"F3", "stands", "", "reason"}), cause: "unexpected finding ID"},
		{name: "case-insensitive top-level fields", output: []byte(`{"Schema":"roundfix/review-validation/v1","Findings":[]}`), cause: "invalid validator object fields"},
		{name: "JSON array", output: []byte(`[]`), cause: "invalid validator object fields"},
		{name: "JSON null", output: []byte(`null`), cause: "invalid validator object fields"},
		{name: "wrong schema", output: []byte(`{"schema":"other","findings":[]}`), cause: "invalid validator schema"},
		{name: "unknown field", output: []byte(`{"schema":"roundfix/review-validation/v1","findings":[],"extra":true}`), cause: "invalid validator object"},
		{name: "missing rule", output: []byte(`{"schema":"roundfix/review-validation/v1","findings":[{"id":"F1","verdict":"stands","reason":"reason"}]}`), cause: "invalid verdict fields"},
		{name: "null rule", output: []byte(`{"schema":"roundfix/review-validation/v1","findings":[{"id":"F1","verdict":"stands","rule":null,"reason":"reason"}]}`), cause: "missing verdict field"},
		{name: "duplicate JSON key", output: []byte(`{"schema":"bad","schema":"roundfix/review-validation/v1","findings":[]}`), cause: "duplicate or invalid JSON key"},
		{name: "invalid verdict", output: validatorAnswer(good, validatorVerdict{"F2", "pass", "", "reason"}), cause: "invalid verdict"},
		{name: "oversized output", output: []byte(strings.Repeat("x", agent.SealedPromptMaxOutputBytes+1)), cause: agent.ErrSealedOutputTooLarge.Error()},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			fixture, runner := validatorFixture(t, "- "+validatorTaskPath+":4 Failure: invalid status\n- "+validatorTaskPath+":11 Failure: invalid result", test.output)
			runner.sealedErr, runner.sealedResult.ToolUsed = test.err, test.tool
			if test.noRunner {
				withAgentRunner(t, runner.reviewCommandRunner)
			}
			code, record, stderr := fixture.run(t)
			if code != exitRunFailed || record.Outcome != reviewOutcomeFindings || record.Validation.Validator != "unavailable" || !strings.HasPrefix(record.Validation.Reason, "validator unavailable: ") || !strings.Contains(record.Validation.Reason, test.cause) || stderr != "" {
				t.Fatalf("exit=%d record=%+v stderr=%q", code, record, stderr)
			}
			if len(record.FindingItems) != 2 {
				t.Fatalf("items=%+v", record.FindingItems)
			}
			for _, finding := range record.FindingItems {
				if finding.Validation.Status != "stands" || finding.Validation.Rule != "" || finding.Validation.Reason != "anchored in the candidate diff" {
					t.Fatalf("finding=%+v", finding)
				}
			}
			if !test.noRunner {
				if questions := validatorQuestions(t, runner); len(questions) != 2 {
					t.Fatalf("questions=%+v", questions)
				}
			}
		})
	}
}

func TestReviewConventionRegionsAdmitAndRefuseAtTheirEdges(t *testing.T) {
	t.Parallel()
	fixture, _ := validatorFixture(t, "- review.txt:2 Failure: test", nil)
	archived := "docs/history/specs/0002-archived"
	validatorWrite(t, fixture.repository, archived+"/qa/report.md", "QA\n")
	validatorWrite(t, fixture.repository, "docs/specs/0001-example/qa/report.md", "QA\n")
	planning := "docs/specs/0003-plan"
	validatorWrite(t, fixture.repository, planning+"/task_01.md", strings.Replace(validatorTaskBody, "status: completed", "status: pending", 1))
	validatorWrite(t, fixture.repository, planning+"/task_02.md", strings.Replace(validatorTaskBody, "status: completed", "status: pending", 1))
	validatorWrite(t, fixture.repository, planning+"/_tasks.md", "graph\n")
	validatorWrite(t, fixture.repository, "docs/specs/0004-plan/_tasks.md", "graph\n")
	gittest.Run(t, fixture.repository, "add", "docs")
	gittest.Run(t, fixture.repository, "commit", "-m", "test: regions")
	roots, err := reviewCandidateSpecRoots("docs/specs", fixture.repository)
	if err != nil {
		t.Fatal(err)
	}
	repo := reviewRepository{fixture.repository, strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD")), fixture.baseCommit, roots, preflight.ExecGitRunner{}}
	tests := []struct {
		name, path, rule string
		start, end       int
		want             bool
	}{
		{"active QA", "docs/specs/0001-example/qa/report.md", "C1", 1, 20, true},
		{"archived QA", archived + "/qa/report.md", "C1", 1, 20, true},
		{"PRD excluded", "docs/specs/0001-example/_prd.md", "C1", 1, 1, false},
		{"front matter", validatorTaskPath, "C2", 1, 5, true},
		{"front matter edge", validatorTaskPath, "C2", 4, 7, false},
		{"result", validatorTaskPath, "C2", 10, 12, true},
		{"result crosses next", validatorTaskPath, "C2", 11, 13, false},
		{"requirements", validatorTaskPath, "C2", 7, 8, false},
		{"archive root", archived + "/_prd.md", "C3", 1, 50, true},
		{"active not archive", validatorTaskPath, "C3", 1, 1, false},
		{"outside archive prefix", "docs/history/specs-other/file.md", "C3", 1, 1, false},
		{"planning manifest", planning + "/_tasks.md", "C4", 1, 20, true},
		{"planning without Task files", "docs/specs/0004-plan/_tasks.md", "C4", 1, 1, true},
		{"planning task", planning + "/task_01.md", "C4", 1, 5, true},
		{"planning requirements", planning + "/task_01.md", "C4", 5, 7, false},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			rules, err := conventionRegions(context.Background(), repo, reviewFindingAnchor{test.path, test.start, test.end})
			if err != nil {
				t.Fatal(err)
			}
			if got := containsString(rules, test.rule); got != test.want {
				t.Fatalf("rules=%v want %s=%v", rules, test.rule, test.want)
			}
		})
	}
	for _, state := range []string{"completed", "QA report"} {
		t.Run(state, func(t *testing.T) {
			if state == "completed" {
				validatorWrite(t, fixture.repository, planning+"/task_02.md", validatorTaskBody)
			} else {
				validatorWrite(t, fixture.repository, planning+"/task_02.md", strings.Replace(validatorTaskBody, "status: completed", "status: pending", 1))
				validatorWrite(t, fixture.repository, planning+"/qa/report.md", "QA\n")
			}
			gittest.Run(t, fixture.repository, "add", "docs")
			gittest.Run(t, fixture.repository, "commit", "-m", "test: planning state")
			changed := repo
			changed.Head = strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))
			for _, file := range []string{"_tasks.md", "task_01.md"} {
				rules, err := conventionRegions(context.Background(), changed, reviewFindingAnchor{planning + "/" + file, 1, 4})
				if err != nil || containsString(rules, "C4") {
					t.Fatalf("rules=%v err=%v", rules, err)
				}
				// The original head remains eligible even after checkout contents changed.
				original, err := conventionRegions(context.Background(), repo, reviewFindingAnchor{planning + "/" + file, 1, 4})
				if err != nil || !containsString(original, "C4") {
					t.Fatalf("head tree drifted: %v %v", original, err)
				}
			}
		})
	}
}

func TestReviewDisposeRefusesAConventionDismissal(t *testing.T) {
	t.Parallel()
	fixture, _ := validatorFixture(t, "- "+validatorTaskPath+":4 Failure: invalid status", validatorAnswer(validatorVerdict{"F1", "dismiss", "convention:C2", "restates settlement"}))
	if code, record, stderr := fixture.run(t); code != exitOK {
		t.Fatalf("exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	code, stdout, stderr := runReviewDispose(t, "F1", "--dismiss", "--evidence", "restates the Daemon settlement")
	if code != exitPreflight || stdout != "" || stderr != "roundfix: review dispose refused: finding \"F1\" was dismissed by validation (convention:C2)\n" {
		t.Fatalf("exit=%d stdout=%q stderr=%q", code, stdout, stderr)
	}
	if _, err := os.Stat(filepath.Join(fixture.artifactDir, reviewDispositionLedgerFileName)); !os.IsNotExist(err) {
		t.Fatalf("ledger written: %v", err)
	}
}

func TestReviewValidatorUsesTheSuccessfulFallbackSelection(t *testing.T) {
	t.Parallel()
	fixture, runner := validatorFixture(t, "- "+validatorTaskPath+":4 Failure: invalid status", validatorAnswer(validatorVerdict{"F1", "dismiss", "convention:C2", "restates settlement"}))
	writeReviewCommandProfileConfig(t, fixture.repository, "codex", fixture.artifactDir, "codex", "codex")
	runner.prepareErrors = []error{&agent.SelectionFailureError{Runtime: "codex", Reason: "adapter startup"}, &agent.SelectionFailureError{Runtime: "codex", Reason: "adapter startup"}}
	code, record, stderr := fixture.run(t)
	if code != exitOK || record.Validation.Validator != "ran" {
		t.Fatalf("exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	validatorQuestions(t, runner)
	if runner.sealedRequests[0].Runtime.Model != "fallback-model" || runner.prepareCalls != 3 {
		t.Fatalf("selection=%+v prepares=%d", runner.sealedRequests[0].Runtime, runner.prepareCalls)
	}
}

func TestReviewValidatorNeverRunsForReviewedOrBlockedAnswers(t *testing.T) {
	t.Parallel()
	for _, answer := range []string{"No findings.", "Findings:\n- no anchor", "malformed verdict"} {
		t.Run(answer, func(t *testing.T) {
			fixture, runner := validatorFixture(t, "unused", nil)
			runner.results[0].result.Message = answer
			_, record, _ := fixture.run(t)
			if record.Outcome != reviewOutcomeReviewed && record.Outcome != reviewOutcomeBlocked {
				t.Fatalf("record=%+v", record)
			}
			if len(runner.sealedRequests) != 0 {
				t.Fatal("validator ran")
			}
		})
	}
}

func TestReviewValidatorExcerptsAreBoundedAndUseTheMergeBaseForDeletedFiles(t *testing.T) {
	t.Parallel()
	fixture, _ := validatorFixture(t, "unused", nil)
	var body strings.Builder
	for i := 0; i < 100; i++ {
		body.WriteString("source line\n")
	}
	validatorWrite(t, fixture.repository, "source.txt", body.String())
	gittest.Run(t, fixture.repository, "add", "source.txt")
	gittest.Run(t, fixture.repository, "commit", "-m", "test: excerpt base")
	base := strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))
	repo := reviewRepository{Root: fixture.repository, Head: base, Base: base, Git: preflight.ExecGitRunner{}}
	for _, test := range []struct {
		name              string
		start, end, count int
		first, last       string
	}{
		{"five lines of context", 10, 12, 13, "5: source line", "17: source line"},
		{"sixty line cap", 10, 90, 60, "5: source line", "64: source line"},
		{"file start", 1, 1, 6, "1: source line", "6: source line"},
		{"file end", 100, 100, 6, "95: source line", "100: source line"},
	} {
		t.Run(test.name, func(t *testing.T) {
			excerpt, err := reviewValidatorExcerpt(context.Background(), repo, reviewFindingAnchor{"source.txt", test.start, test.end})
			if err != nil || strings.Count(excerpt, "\n") != test.count || !strings.HasPrefix(excerpt, test.first+"\n") || !strings.HasSuffix(excerpt, test.last+"\n") {
				t.Fatalf("excerpt=%q err=%v", excerpt, err)
			}
		})
	}
	gittest.Run(t, fixture.repository, "rm", "source.txt")
	gittest.Run(t, fixture.repository, "commit", "-m", "test: deleted source")
	repo.Head = strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD"))
	excerpt, err := reviewValidatorExcerpt(context.Background(), repo, reviewFindingAnchor{"source.txt", 10, 12})
	if err != nil || !strings.Contains(excerpt, "10: source line") {
		t.Fatalf("deleted excerpt=%q err=%v", excerpt, err)
	}
}

func TestReviewConventionC2IncludesEverySettlementSection(t *testing.T) {
	t.Parallel()
	fixture, _ := validatorFixture(t, "unused", nil)
	body := "---\nstatus: completed\n---\n## Recorded paths\nfile.go\n## Carry-forward provenance\nsource\n## Result\nevidence\n## Requirements\nrequirement\n"
	validatorWrite(t, fixture.repository, validatorTaskPath, body)
	gittest.Run(t, fixture.repository, "add", "docs")
	gittest.Run(t, fixture.repository, "commit", "-m", "test: settlement sections")
	roots, err := reviewCandidateSpecRoots("docs/specs", fixture.repository)
	if err != nil {
		t.Fatal(err)
	}
	repo := reviewRepository{fixture.repository, strings.TrimSpace(gittest.Run(t, fixture.repository, "rev-parse", "HEAD")), fixture.baseCommit, roots, preflight.ExecGitRunner{}}
	for _, test := range []struct {
		name       string
		start, end int
		want       bool
	}{
		{"recorded paths", 4, 5, true}, {"carry-forward", 6, 7, true}, {"adjacent settlement sections", 4, 9, true}, {"crosses requirement", 9, 10, false}, {"invalid zero", 0, 1, false}, {"past end", 8, 20, false},
	} {
		t.Run(test.name, func(t *testing.T) {
			rules, err := conventionRegions(context.Background(), repo, reviewFindingAnchor{validatorTaskPath, test.start, test.end})
			if err != nil || containsString(rules, "C2") != test.want {
				t.Fatalf("rules=%v err=%v", rules, err)
			}
		})
	}
}

func TestReviewValidatorCanDismissNoFailureOnlyWhenEligible(t *testing.T) {
	t.Parallel()
	fixture, runner := validatorFixture(t, "- review.txt:2 A concern without a failure", validatorAnswer(validatorVerdict{"F1", "dismiss", "no-failure", "states no failure"}))
	code, record, stderr := fixture.run(t)
	if code != exitOK || record.Outcome != reviewOutcomeFindingsDismissed || record.FindingItems[0].Validation.Rule != "no-failure" || stderr != "roundfix: review finding F1 dismissed by validation (no-failure): states no failure\n" {
		t.Fatalf("exit=%d record=%+v stderr=%q", code, record, stderr)
	}
	questions := validatorQuestions(t, runner)
	if len(questions) != 1 || !reflect.DeepEqual(questions[0].EligibleRules, []string{"no-failure"}) {
		t.Fatalf("questions=%+v", questions)
	}
}

func TestReviewValidatorRejectsOversizedInputBeforeCallingTheRunner(t *testing.T) {
	t.Parallel()
	fixture, runner := validatorFixture(t, "unused", nil)
	workDir := t.TempDir()
	if err := os.Chmod(workDir, 0700); err != nil {
		t.Fatal(err)
	}
	_, err := runConventionValidator(context.Background(), runner, agent.RuntimeSpec{}, workDir, []validatorQuestion{{ID: "F1", Anchor: reviewFindingAnchor{validatorTaskPath, 4, 4}, Text: strings.Repeat("x", agent.SealedPromptMaxInputBytes), EligibleRules: []string{"convention:C2"}}})
	if err == nil || !strings.Contains(err.Error(), "input exceeds sealed limit") || len(runner.sealedRequests) != 0 {
		t.Fatalf("err=%v calls=%d repo=%s", err, len(runner.sealedRequests), fixture.repository)
	}
}
