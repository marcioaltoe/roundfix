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
	prompt, err := BuildTaskPrompt(sampleTaskPromptRequest())
	if err != nil {
		t.Fatalf("BuildTaskPrompt returned error: %v", err)
	}
	if !strings.Contains(prompt, sampleTaskContent) {
		t.Fatalf("expected task prompt to embed the full task file content verbatim, got:\n%s", prompt)
	}
}

func TestBuildTaskPromptRendersSpecContextBundlePathOnly(t *testing.T) {
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

func TestBuildQAPromptStatesQAGateContract(t *testing.T) {
	prompt, err := BuildQAPrompt(
		"0001-implement-command",
		"/repo/docs/specs/0001-implement-command",
		"/repo/docs/specs/0001-implement-command/_prd.md",
	)
	if err != nil {
		t.Fatalf("BuildQAPrompt returned error: %v", err)
	}

	for _, expected := range []string{
		"Spec: 0001-implement-command",
		"Spec directory: /repo/docs/specs/0001-implement-command",
		"PRD: /repo/docs/specs/0001-implement-command/_prd.md",
		"Run the qa-gate process for this Spec",
		"Write the QA Report to the Spec's qa/ directory as qa-report-YYYY-MM-DD.md",
		"frontmatter must carry the verdict: pass, fail, or partial",
		"Use verdict: pass only when every criterion passes.",
		"Never commit, push, or open a pull request.",
	} {
		if !strings.Contains(prompt, expected) {
			t.Fatalf("expected QA prompt to contain %q, got:\n%s", expected, prompt)
		}
	}
}

func TestBuildQAPromptDeterministicForIdenticalInput(t *testing.T) {
	first, err := BuildQAPrompt("0001-implement-command", "/repo/docs/specs/0001-implement-command", "/repo/docs/specs/0001-implement-command/_prd.md")
	if err != nil {
		t.Fatalf("BuildQAPrompt returned error: %v", err)
	}
	second, err := BuildQAPrompt("0001-implement-command", "/repo/docs/specs/0001-implement-command", "/repo/docs/specs/0001-implement-command/_prd.md")
	if err != nil {
		t.Fatalf("BuildQAPrompt returned error: %v", err)
	}
	if first != second {
		t.Fatalf("expected identical input to yield identical output, got:\n--- first ---\n%s\n--- second ---\n%s", first, second)
	}
}

func TestBuildQAPromptValidatesRequiredFields(t *testing.T) {
	tests := []struct {
		name     string
		specSlug string
		specDir  string
		prdPath  string
	}{
		{name: "empty spec slug", specSlug: "", specDir: "/repo/docs/specs/0001", prdPath: "/repo/docs/specs/0001/_prd.md"},
		{name: "empty spec directory", specSlug: "0001", specDir: "", prdPath: "/repo/docs/specs/0001/_prd.md"},
		{name: "empty prd path", specSlug: "0001", specDir: "/repo/docs/specs/0001", prdPath: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt, err := BuildQAPrompt(tt.specSlug, tt.specDir, tt.prdPath)
			if err == nil {
				t.Fatalf("expected error for %s, got prompt:\n%s", tt.name, prompt)
			}
			if prompt != "" {
				t.Fatalf("expected empty prompt on error, got:\n%s", prompt)
			}
		})
	}
}
