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
		"Set status: in_progress in the task file frontmatter when you start.",
		"Run the commands in the task file's ## Verification section while working; all must pass.",
		"Append a ## Result section to the task file with evidence",
		"Settle the task file frontmatter to status: completed or status: failed",
		"Never commit, push, or open a pull request.",
		"Never edit the Task Graph manifest (_tasks.md) or any other task file.",
		"a prior Run died mid-task; start the Task fresh.",
		"The work-target lock guarantees no live owner.",
	} {
		if !strings.Contains(prompt, expected) {
			t.Fatalf("expected task prompt to contain %q, got:\n%s", expected, prompt)
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
