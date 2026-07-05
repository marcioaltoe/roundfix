package agent

import (
	"errors"
	"fmt"
	"strings"
)

// TaskPromptRequest carries the plain-string inputs for one Spec Task
// prompt. Inputs stay plain strings so the builders never depend on the
// Spec parser or any store/daemon package.
type TaskPromptRequest struct {
	SpecSlug    string
	TaskID      string
	TaskPath    string
	TaskContent string
}

// taskExecutionInvariants states the execution contract the Agent must
// follow for one Spec Task, mirroring the implement-task skill. Every
// invariant lives in this one constant so the templating pass (work-plan
// item 5) replaces the contract in a single place.
const taskExecutionInvariants = `Execution invariants:
- Implement only this Task's slice; work that belongs to another Task is a follow-up note, not part of this diff.
- Set status: in_progress in the task file frontmatter when you start.
- Run the commands in the task file's ## Verification section while working; all must pass.
- Append a ## Result section to the task file with evidence for each acceptance criterion.
- Settle the task file frontmatter to status: completed or status: failed before you finish.
- Never commit, push, or open a pull request.
- Never edit the Task Graph manifest (_tasks.md) or any other task file.
- If the task file arrives with status: in_progress, a prior Run died mid-task; start the Task fresh. The work-target lock guarantees no live owner.
`

// qaGateContract states the QA contract the Agent must follow for one
// Spec, mirroring the qa-gate skill. Every rule lives in this one constant
// so the templating pass (work-plan item 5) replaces the contract in a
// single place.
const qaGateContract = `QA contract:
- Run the qa-gate process for this Spec: validate every user story and acceptance criterion from the PRD and the task files against the real application.
- Write the QA Report to the Spec's qa/ directory as qa-report-YYYY-MM-DD.md, using today's date.
- The QA Report frontmatter must carry the verdict: pass, fail, or partial. Use verdict: pass only when every criterion passes.
- Never commit, push, or open a pull request.
`

// BuildTaskPrompt builds the prompt for one Spec Task: the Task identity,
// the execution invariants mirroring implement-task, and the full task
// file content embedded verbatim.
func BuildTaskPrompt(req TaskPromptRequest) (string, error) {
	if strings.TrimSpace(req.SpecSlug) == "" {
		return "", errors.New("Spec slug is required")
	}
	if strings.TrimSpace(req.TaskID) == "" {
		return "", errors.New("Task id is required")
	}
	if strings.TrimSpace(req.TaskPath) == "" {
		return "", errors.New("task file path is required")
	}
	if strings.TrimSpace(req.TaskContent) == "" {
		return "", errors.New("task file content is required")
	}
	var builder strings.Builder
	builder.WriteString("You are the Roundfix child Agent for one Spec Task.\n\n")
	builder.WriteString(fmt.Sprintf("Spec: %s\n", req.SpecSlug))
	builder.WriteString(fmt.Sprintf("Task: %s\n", req.TaskID))
	builder.WriteString(fmt.Sprintf("Task file: %s\n\n", req.TaskPath))
	builder.WriteString(taskExecutionInvariants)
	builder.WriteString("\nTask file content, verbatim:\n\n")
	builder.WriteString(req.TaskContent)
	return builder.String(), nil
}

// BuildQAPrompt builds the prompt for one Spec QA gate: the Spec identity
// and the QA contract mirroring qa-gate.
func BuildQAPrompt(specSlug, specDir, prdPath string) (string, error) {
	if strings.TrimSpace(specSlug) == "" {
		return "", errors.New("Spec slug is required")
	}
	if strings.TrimSpace(specDir) == "" {
		return "", errors.New("Spec directory is required")
	}
	if strings.TrimSpace(prdPath) == "" {
		return "", errors.New("PRD path is required")
	}
	var builder strings.Builder
	builder.WriteString("You are the Roundfix child Agent for one Spec QA gate.\n\n")
	builder.WriteString(fmt.Sprintf("Spec: %s\n", specSlug))
	builder.WriteString(fmt.Sprintf("Spec directory: %s\n", specDir))
	builder.WriteString(fmt.Sprintf("PRD: %s\n\n", prdPath))
	builder.WriteString(qaGateContract)
	return builder.String(), nil
}
