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
	Context     SpecContextBundle
}

// SpecContextBundle is the path-only Task-start manifest that orients an
// Agent without embedding larger Spec, instruction, source, or diff content.
type SpecContextBundle struct {
	PRD               string
	TechSpec          string
	TaskGraph         string
	Instructions      []string
	Interfaces        []string
	PriorChangedFiles []string
	OmittedPriorFiles int
}

// taskExecutionInvariants states the ADR-0057 handoff contract the Agent must
// follow for one Spec Task. Every invariant lives in this one constant; the
// dedicated authorial-Skill Task aligns the protected implement-task skill.
const taskExecutionInvariants = `Execution invariants:
- Implement only this Task's slice; work that belongs to another Task is a follow-up note, not part of this diff.
- The Daemon owns Task status during Implement; do not edit the task file's status field.
- Do not run commands from the task file's ## Verification section; the Daemon runs them after your turn.
- Run focused implementation checks while working when useful.
- Append or update a ## Result section in the task file with implementation and focused-check evidence for each acceptance criterion.
- Hand back implementation-ready work without claiming the Task is completed or failed.
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
	if req.Context.OmittedPriorFiles < 0 {
		return "", errors.New("omitted prior file count must not be negative")
	}
	var builder strings.Builder
	builder.WriteString("You are the Roundfix child Agent for one Spec Task.\n\n")
	builder.WriteString(fmt.Sprintf("Spec: %s\n", req.SpecSlug))
	builder.WriteString(fmt.Sprintf("Task: %s\n", req.TaskID))
	builder.WriteString(fmt.Sprintf("Task file: %s\n\n", req.TaskPath))
	builder.WriteString(taskExecutionInvariants)
	writeSpecContextBundle(&builder, req.Context)
	builder.WriteString("\nTask file content, verbatim:\n\n")
	builder.WriteString(req.TaskContent)
	return builder.String(), nil
}

func writeSpecContextBundle(builder *strings.Builder, bundle SpecContextBundle) {
	if !hasSpecContextBundle(bundle) {
		return
	}
	builder.WriteString("\nSpec Context Bundle:\n")
	writeNamedPath(builder, "PRD", bundle.PRD)
	writeNamedPath(builder, "TechSpec", bundle.TechSpec)
	writeNamedPath(builder, "Task Graph", bundle.TaskGraph)
	writePathList(builder, "Instructions", bundle.Instructions)
	writePathList(builder, "Interfaces", bundle.Interfaces)
	writePathList(builder, "Prior changed files", bundle.PriorChangedFiles)
	builder.WriteString(fmt.Sprintf("- Omitted prior files: %d\n", bundle.OmittedPriorFiles))
}

func hasSpecContextBundle(bundle SpecContextBundle) bool {
	return strings.TrimSpace(bundle.PRD) != "" ||
		strings.TrimSpace(bundle.TechSpec) != "" ||
		strings.TrimSpace(bundle.TaskGraph) != "" ||
		len(bundle.Instructions) > 0 ||
		len(bundle.Interfaces) > 0 ||
		len(bundle.PriorChangedFiles) > 0 ||
		bundle.OmittedPriorFiles > 0
}

func writeNamedPath(builder *strings.Builder, label string, path string) {
	if strings.TrimSpace(path) == "" {
		return
	}
	builder.WriteString(fmt.Sprintf("- %s: %s\n", label, path))
}

func writePathList(builder *strings.Builder, label string, paths []string) {
	if len(paths) == 0 {
		return
	}
	builder.WriteString(fmt.Sprintf("- %s:\n", label))
	for _, path := range paths {
		if strings.TrimSpace(path) == "" {
			continue
		}
		builder.WriteString(fmt.Sprintf("  - %s\n", path))
	}
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
