package agent

import (
	"strings"
	"testing"
)

const sampleTaskContent = `---
task: task_02
spec: 0001-implement-command
status: pending
---

# Task 02: Sample slice

## Verification

- ` + "`rtk go test ./...`" + ` — expected: all tests pass.
`

func sampleTaskPromptRequest() TaskPromptRequest {
	return TaskPromptRequest{
		SpecSlug:    "0001-implement-command",
		TaskID:      "task_02",
		TaskPath:    "/repo/docs/specs/0001-implement-command/task_02.md",
		TaskContent: sampleTaskContent,
	}
}

func TestBuildTaskPromptStatesExecutionInvariants(t *testing.T) {
	t.Parallel()

	prompt, err := BuildTaskPrompt(sampleTaskPromptRequest())
	if err != nil {
		t.Fatalf("BuildTaskPrompt returned error: %v", err)
	}

	for _, expected := range []string{
		"Spec: 0001-implement-command",
		"Task: task_02",
		"Task file: /repo/docs/specs/0001-implement-command/task_02.md",
		"Implement only this Task's slice",
		"The Daemon owns Task status during Implement; do not edit the task file's status field.",
		"Do not run commands from the task file's ## Verification section",
		"Run focused implementation checks while working when useful",
		"Append or update a ## Result section in the task file with implementation and focused-check evidence",
		"Hand back implementation-ready work without claiming the Task is completed or failed",
		"Never commit, push, or open a pull request.",
		"Never edit the Task Graph manifest (_tasks.md) or any other task file.",
		"a prior Run died mid-task; start the Task fresh.",
		"The work-target lock guarantees no live owner.",
	} {
		if !strings.Contains(prompt, expected) {
			t.Fatalf("expected task prompt to contain %q, got:\n%s", expected, prompt)
		}
	}
	if strings.Contains(prompt, "Run the commands in the task file's ## Verification section while working; all must pass.") {
		t.Fatalf("expected task prompt to remove authoritative Verification requirement, got:\n%s", prompt)
	}
	for _, forbidden := range []string{
		"Set status: in_progress in the task file frontmatter when you start.",
		"Settle the task file frontmatter to status: completed or status: failed",
	} {
		if strings.Contains(prompt, forbidden) {
			t.Fatalf("expected task prompt to forbid Agent status authorship instead of containing %q, got:\n%s", forbidden, prompt)
		}
	}
}

func TestBuildTaskPromptEmbedsTaskContentVerbatim(t *testing.T) {
	t.Parallel()

	prompt, err := BuildTaskPrompt(sampleTaskPromptRequest())
	if err != nil {
		t.Fatalf("BuildTaskPrompt returned error: %v", err)
	}
	if !strings.Contains(prompt, sampleTaskContent) {
		t.Fatalf("expected task prompt to embed the full task file content verbatim, got:\n%s", prompt)
	}
}

func TestBuildTaskPromptRendersSpecContextBundlePathOnly(t *testing.T) {
	t.Parallel()

	req := sampleTaskPromptRequest()
	req.Context = SpecContextBundle{
		PRD:       "docs/specs/0001-implement-command/_prd.md",
		TechSpec:  "docs/specs/0001-implement-command/_techspec.md",
		TaskGraph: "docs/specs/0001-implement-command/_tasks.md",
		Instructions: []string{
			"AGENTS.md",
			".agents/skills/implement-task/SKILL.md",
			".agents/skills/golang-testing/SKILL.md",
		},
		Interfaces:        []string{"internal/spec/task.go"},
		PriorChangedFiles: []string{"internal/agent/spec_prompt.go", "internal/daemon/task_engine.go"},
		OmittedPriorFiles: 7,
	}

	prompt, err := BuildTaskPrompt(req)
	if err != nil {
		t.Fatalf("BuildTaskPrompt returned error: %v", err)
	}
	for _, expected := range []string{
		"Spec Context Bundle:",
		"- PRD: docs/specs/0001-implement-command/_prd.md",
		"- TechSpec: docs/specs/0001-implement-command/_techspec.md",
		"- Task Graph: docs/specs/0001-implement-command/_tasks.md",
		"- Instructions:",
		"  - AGENTS.md",
		"  - .agents/skills/implement-task/SKILL.md",
		"  - .agents/skills/golang-testing/SKILL.md",
		"- Interfaces:",
		"  - internal/spec/task.go",
		"- Prior changed files:",
		"  - internal/agent/spec_prompt.go",
		"  - internal/daemon/task_engine.go",
		"- Omitted prior files: 7",
	} {
		if !strings.Contains(prompt, expected) {
			t.Fatalf("expected prompt to contain %q, got:\n%s", expected, prompt)
		}
	}
	if got := strings.Count(prompt, sampleTaskContent); got != 1 {
		t.Fatalf("expected exactly one complete task file embedding, got %d", got)
	}
	for _, forbidden := range []string{
		"# Context-Efficient Runs",
		"package daemon",
		"diff --git",
		"@@ -",
		`{"kind":`,
	} {
		if strings.Contains(prompt, forbidden) {
			t.Fatalf("prompt exposed non-manifest content marker %q:\n%s", forbidden, prompt)
		}
	}
}

func TestBuildTaskPromptDeterministicForIdenticalInput(t *testing.T) {
	t.Parallel()

	first, err := BuildTaskPrompt(sampleTaskPromptRequest())
	if err != nil {
		t.Fatalf("BuildTaskPrompt returned error: %v", err)
	}
	second, err := BuildTaskPrompt(sampleTaskPromptRequest())
	if err != nil {
		t.Fatalf("BuildTaskPrompt returned error: %v", err)
	}
	if first != second {
		t.Fatalf("expected identical input to yield identical output, got:\n--- first ---\n%s\n--- second ---\n%s", first, second)
	}
}

func TestBuildTaskPromptValidatesRequiredFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		mutate func(req *TaskPromptRequest)
	}{
		{name: "empty spec slug", mutate: func(req *TaskPromptRequest) { req.SpecSlug = "" }},
		{name: "empty task id", mutate: func(req *TaskPromptRequest) { req.TaskID = "" }},
		{name: "empty task path", mutate: func(req *TaskPromptRequest) { req.TaskPath = "" }},
		{name: "empty task content", mutate: func(req *TaskPromptRequest) { req.TaskContent = "" }},
		{name: "whitespace-only task content", mutate: func(req *TaskPromptRequest) { req.TaskContent = "  \n\t" }},
		{name: "negative omitted prior file count", mutate: func(req *TaskPromptRequest) { req.Context.OmittedPriorFiles = -1 }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := sampleTaskPromptRequest()
			tt.mutate(&req)
			prompt, err := BuildTaskPrompt(req)
			if err == nil {
				t.Fatalf("expected error for %s, got prompt:\n%s", tt.name, prompt)
			}
			if prompt != "" {
				t.Fatalf("expected empty prompt on error, got:\n%s", prompt)
			}
		})
	}
}

func sampleQAPromptRequest() QAPromptRequest {
	return QAPromptRequest{
		SpecSlug:     "0001-implement-command",
		SpecDir:      "/repo/docs/specs/0001-implement-command",
		PRDPath:      "/repo/docs/specs/0001-implement-command/_prd.md",
		RunBranch:    "roundfix/run-run_20260728T041231Z_48d72c1b142ea37b",
		TargetBranch: "ma/widget-flow",
		UserCheckout: "/repo",
		PullRequest:  "#40 (owner/repo)",
	}
}

func TestBuildQAPromptCarriesTheSpecContextBundle(t *testing.T) {
	t.Parallel()

	t.Run("changed paths", func(t *testing.T) {
		t.Parallel()

		req := sampleQAPromptRequest()
		req.Context = SpecContextBundle{
			PRD:               "docs/specs/0001-implement-command/_prd.md",
			TechSpec:          "docs/specs/0001-implement-command/_techspec.md",
			TaskGraph:         "docs/specs/0001-implement-command/_tasks.md",
			Instructions:      []string{"AGENTS.md", ".agents/skills/implement-task/SKILL.md"},
			Interfaces:        []string{"internal/agent/spec_prompt.go"},
			PriorChangedFiles: []string{"internal/daemon/task_context.go"},
			OmittedPriorFiles: 3,
		}

		prompt, err := BuildQAPrompt(req)
		if err != nil {
			t.Fatalf("BuildQAPrompt returned error: %v", err)
		}
		for _, expected := range []string{
			"Spec Context Bundle:",
			"- PRD: docs/specs/0001-implement-command/_prd.md",
			"- TechSpec: docs/specs/0001-implement-command/_techspec.md",
			"- Task Graph: docs/specs/0001-implement-command/_tasks.md",
			"  - AGENTS.md",
			"  - .agents/skills/implement-task/SKILL.md",
			"  - internal/agent/spec_prompt.go",
			"- Prior changed files:\n  - internal/daemon/task_context.go",
			"- Omitted prior files: 3",
		} {
			if !strings.Contains(prompt, expected) {
				t.Fatalf("expected QA prompt to contain %q, got:\n%s", expected, prompt)
			}
		}
	})

	t.Run("no changed paths", func(t *testing.T) {
		t.Parallel()

		req := sampleQAPromptRequest()
		req.Context = SpecContextBundle{PRD: "docs/specs/0001-implement-command/_prd.md"}

		prompt, err := BuildQAPrompt(req)
		if err != nil {
			t.Fatalf("BuildQAPrompt returned error: %v", err)
		}
		if !strings.Contains(prompt, "- Prior changed files: none\n") {
			t.Fatalf("expected QA prompt to state that no paths changed, got:\n%s", prompt)
		}
	})
}

func TestBuildQAPromptCarriesThePreviousReportIdentity(t *testing.T) {
	t.Parallel()

	t.Run("previous report exists", func(t *testing.T) {
		t.Parallel()

		req := sampleQAPromptRequest()
		req.PreviousReportPath = "docs/specs/0001-implement-command/qa/qa-report-2026-08-10.md"
		req.PreviousReportHead = "0123456789abcdef"

		prompt, err := BuildQAPrompt(req)
		if err != nil {
			t.Fatalf("BuildQAPrompt returned error: %v", err)
		}
		for _, expected := range []string{
			"Previous QA Report: docs/specs/0001-implement-command/qa/qa-report-2026-08-10.md\n",
			"Previous QA Report head: 0123456789abcdef\n",
		} {
			if !strings.Contains(prompt, expected) {
				t.Fatalf("expected QA prompt to contain %q, got:\n%s", expected, prompt)
			}
		}
	})

	t.Run("no previous report", func(t *testing.T) {
		t.Parallel()

		prompt, err := BuildQAPrompt(sampleQAPromptRequest())
		if err != nil {
			t.Fatalf("BuildQAPrompt returned error: %v", err)
		}
		if !strings.Contains(prompt, "Previous QA Report: none exists for this Spec.\n") {
			t.Fatalf("expected QA prompt to state that no previous report exists, got:\n%s", prompt)
		}
		if strings.Contains(prompt, "Previous QA Report head:") {
			t.Fatalf("expected no previous-report head without a report, got:\n%s", prompt)
		}
	})
}

func TestBuildQAPromptStatesQAGateContract(t *testing.T) {
	t.Parallel()

	prompt, err := BuildQAPrompt(sampleQAPromptRequest())
	if err != nil {
		t.Fatalf("BuildQAPrompt returned error: %v", err)
	}

	for _, expected := range []string{
		"Spec: 0001-implement-command",
		"Spec directory: /repo/docs/specs/0001-implement-command",
		"PRD: /repo/docs/specs/0001-implement-command/_prd.md",
		"Pull Request: #40 (owner/repo)",
		"Run the qa-gate process for this Spec",
		"qa-report-YYYY-MM-DD.md for the day's first report",
		"qa-report-YYYY-MM-DD-NN.md with a numeric -NN suffix for same-day reruns",
		"frontmatter must carry the verdict: pass, fail, or partial",
		"rows_blocked_environment, rows_blocked_finding, and rows_blocked_declared",
		"no row is declared-blocked, finding-blocked, or skipped",
		"A nonzero rows_blocked_environment does not by itself prevent pass; a nonzero rows_blocked_finding or rows_blocked_declared does.",
		"Never commit, push, or open a pull request.",
	} {
		if !strings.Contains(prompt, expected) {
			t.Fatalf("expected QA prompt to contain %q, got:\n%s", expected, prompt)
		}
	}
	if strings.Contains(prompt, "pass only when every criterion passes") {
		t.Fatalf("the QA contract still imposes the permanent partial ceiling ADR-0080 removes, got:\n%s", prompt)
	}
	if strings.Contains(prompt, "-<scope-or-build>") {
		t.Fatalf("expected the QA prompt to allow only numeric same-day suffixes, got:\n%s", prompt)
	}
}

// The gate reasons about the user's branch from a checkout that can never
// be on it, so the prompt has to name both branches and separate them.
func TestBuildQAPromptStatesCheckoutFactsSeparatingRunBranchFromTarget(t *testing.T) {
	t.Parallel()

	prompt, err := BuildQAPrompt(sampleQAPromptRequest())
	if err != nil {
		t.Fatalf("BuildQAPrompt returned error: %v", err)
	}

	for _, expected := range []string{
		"Run Worktree branch: roundfix/run-run_20260728T041231Z_48d72c1b142ea37b (this checkout only — a per-Run branch that is never pushed and has no Pull Request of its own)\n",
		"Spec target branch: ma/widget-flow (the user branch this Spec's commits land on; any Pull Request for this Spec is open on this branch, never on the Run Worktree branch)\n",
		"User checkout: /repo (the user's repository root this Run Worktree was created from)\n",
		"Pull Request: #40 (owner/repo)\n",
	} {
		if !strings.Contains(prompt, expected) {
			t.Fatalf("expected QA prompt to contain %q, got:\n%s", expected, prompt)
		}
	}
	// Facts only: the QA contract keeps the rules it mirrors from the
	// protected qa-gate skill, and the checkout facts add none.
	if got := strings.Count(prompt, qaGateContract); got != 1 {
		t.Fatalf("expected the QA contract exactly once, got %d in:\n%s", got, prompt)
	}
	factsEnd := strings.Index(prompt, qaGateContract)
	if factsEnd < 0 {
		t.Fatalf("expected the QA contract in the prompt, got:\n%s", prompt)
	}
	if !strings.HasSuffix(prompt[:factsEnd], "Pull Request: #40 (owner/repo)\nPrevious QA Report: none exists for this Spec.\n\n") {
		t.Fatalf("expected the checkout facts to end the fact block before the QA contract, got:\n%s", prompt)
	}
}

func TestBuildQAPromptStatesPullRequestJourneysAreEnvironmentBlockedWhenNoneIsOpen(t *testing.T) {
	t.Parallel()

	req := sampleQAPromptRequest()
	req.PullRequest = ""
	req.PullRequestResolved = true

	prompt, err := BuildQAPrompt(req)
	if err != nil {
		t.Fatalf("BuildQAPrompt returned error: %v", err)
	}

	const expected = "Pull Request: none open; Pull Request journeys are environment-blocked.\n"
	if !strings.Contains(prompt, expected) {
		t.Fatalf("expected the QA prompt to contain %q, got:\n%s", expected, prompt)
	}
}

// A proven absence and an unknown one are both environment-blocked, but only
// the first may be reported as an absence. Collapsing them lets a gh failure
// reach the gate as evidence that no Pull Request exists.
func TestBuildQAPromptSeparatesUnresolvedPullRequestFromProvenAbsence(t *testing.T) {
	t.Parallel()

	resolved := sampleQAPromptRequest()
	resolved.PullRequest = ""
	resolved.PullRequestResolved = true

	unresolved := sampleQAPromptRequest()
	unresolved.PullRequest = ""
	unresolved.PullRequestResolved = false

	resolvedPrompt, err := BuildQAPrompt(resolved)
	if err != nil {
		t.Fatalf("BuildQAPrompt(resolved) returned error: %v", err)
	}
	unresolvedPrompt, err := BuildQAPrompt(unresolved)
	if err != nil {
		t.Fatalf("BuildQAPrompt(unresolved) returned error: %v", err)
	}

	if resolvedPrompt == unresolvedPrompt {
		t.Fatal("a proven absent Pull Request and an unresolvable one produced the same prompt")
	}
	if strings.Contains(unresolvedPrompt, "none open") {
		t.Fatalf("an unresolvable lookup reported a confirmed absence, got:\n%s", unresolvedPrompt)
	}
	const expected = "Pull Request: could not be resolved; Pull Request journeys are environment-blocked and their absence is unproven — do not record a confirmed absence.\n"
	if !strings.Contains(unresolvedPrompt, expected) {
		t.Fatalf("expected the QA prompt to contain %q, got:\n%s", expected, unresolvedPrompt)
	}
}

func TestBuildQAPromptOmitsUnrecordedCheckoutFacts(t *testing.T) {
	t.Parallel()

	req := sampleQAPromptRequest()
	req.RunBranch = ""
	req.TargetBranch = ""
	req.UserCheckout = ""

	prompt, err := BuildQAPrompt(req)
	if err != nil {
		t.Fatalf("BuildQAPrompt returned error: %v", err)
	}
	for _, forbidden := range []string{"Run Worktree branch:", "Spec target branch:", "User checkout:"} {
		if strings.Contains(prompt, forbidden) {
			t.Fatalf("expected no %q line for an unrecorded fact, got:\n%s", forbidden, prompt)
		}
	}
	if !strings.Contains(prompt, "PRD: /repo/docs/specs/0001-implement-command/_prd.md\nPull Request: #40 (owner/repo)\nPrevious QA Report: none exists for this Spec.\n\n"+qaGateContract) {
		t.Fatalf("expected a usable prompt with the Spec identity and the QA contract, got:\n%s", prompt)
	}
}

func TestBuildQAPromptOmitsIndividuallyUnrecordedCheckoutFacts(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		mutate    func(req *QAPromptRequest)
		forbidden string
	}{
		{name: "no run branch", mutate: func(req *QAPromptRequest) { req.RunBranch = "" }, forbidden: "Run Worktree branch:"},
		{name: "no target branch", mutate: func(req *QAPromptRequest) { req.TargetBranch = "" }, forbidden: "Spec target branch:"},
		{name: "no user checkout", mutate: func(req *QAPromptRequest) { req.UserCheckout = "" }, forbidden: "User checkout:"},
		{name: "whitespace target branch", mutate: func(req *QAPromptRequest) { req.TargetBranch = "  \t" }, forbidden: "Spec target branch:"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := sampleQAPromptRequest()
			tt.mutate(&req)
			prompt, err := BuildQAPrompt(req)
			if err != nil {
				t.Fatalf("BuildQAPrompt returned error: %v", err)
			}
			if strings.Contains(prompt, tt.forbidden) {
				t.Fatalf("expected no %q line, got:\n%s", tt.forbidden, prompt)
			}
			if !strings.Contains(prompt, qaGateContract) {
				t.Fatalf("expected the QA contract to survive a missing fact, got:\n%s", prompt)
			}
		})
	}
}

func TestBuildQAPromptDeterministicForIdenticalInput(t *testing.T) {
	t.Parallel()

	first, err := BuildQAPrompt(sampleQAPromptRequest())
	if err != nil {
		t.Fatalf("BuildQAPrompt returned error: %v", err)
	}
	second, err := BuildQAPrompt(sampleQAPromptRequest())
	if err != nil {
		t.Fatalf("BuildQAPrompt returned error: %v", err)
	}
	if first != second {
		t.Fatalf("expected identical input to yield identical output, got:\n--- first ---\n%s\n--- second ---\n%s", first, second)
	}
}

func TestBuildQAPromptValidatesRequiredFields(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		mutate func(req *QAPromptRequest)
	}{
		{name: "empty spec slug", mutate: func(req *QAPromptRequest) { req.SpecSlug = "" }},
		{name: "empty spec directory", mutate: func(req *QAPromptRequest) { req.SpecDir = "" }},
		{name: "empty prd path", mutate: func(req *QAPromptRequest) { req.PRDPath = "" }},
		{name: "negative omitted prior file count", mutate: func(req *QAPromptRequest) { req.Context.OmittedPriorFiles = -1 }},
		{name: "previous report without head", mutate: func(req *QAPromptRequest) { req.PreviousReportPath = "qa/qa-report-2026-08-10.md" }},
		{name: "previous report head without path", mutate: func(req *QAPromptRequest) { req.PreviousReportHead = "0123456789abcdef" }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := sampleQAPromptRequest()
			tt.mutate(&req)
			prompt, err := BuildQAPrompt(req)
			if err == nil {
				t.Fatalf("expected error for %s, got prompt:\n%s", tt.name, prompt)
			}
			if prompt != "" {
				t.Fatalf("expected empty prompt on error, got:\n%s", prompt)
			}
		})
	}
}
