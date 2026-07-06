---
task: task_03
spec: 0001-implement-command
status: completed
type: backend
complexity: low
---

# Task 03: Build the task and QA prompt builders

## Overview

Add `BuildTaskPrompt` and `BuildQAPrompt` beside the existing review prompt builder: minimal string builders that carry the execution invariants mirroring the implement-task and qa-gate contracts, so agent-side and daemon-side expectations cannot drift. Verifiable on its own through unit tests asserting the contract lines.

## Requirements

1. MUST build the task prompt from the Spec slug, Task id, task file path, and the full task file content embedded in the prompt.
2. MUST state the execution invariants: implement only this Task's slice; set `status: in_progress` on start; run the Verification commands while working; append a `## Result` section with evidence; settle `status: completed` or `status: failed`; never commit, push, or open a pull request; never edit the Task Graph manifest or other task files.
3. MUST state that a stale `in_progress` status means a dead prior Run and the Agent starts the Task fresh — the work-target lock guarantees no live owner.
4. MUST build the QA prompt from the Spec slug, spec directory, and PRD path: run the qa-gate process for the Spec, write the QA Report with its verdict frontmatter into the Spec's `qa/` directory, and never commit.
5. MUST produce deterministic output for identical input, with the invariants kept in one commented constant so a future templating pass (work-plan item 5) replaces them in one place.
6. SHOULD keep both builders free of any dependency on the Spec parser — inputs are plain strings.

## Subtasks

- [x] Task prompt builder with embedded task content
- [x] Execution invariants constant mirroring implement-task
- [x] QA prompt builder mirroring qa-gate
- [x] Unit tests asserting each invariant line and determinism

## Acceptance Criteria

- [x] Task prompt tests assert every invariant from Requirements 2–3 appears, the task content is embedded verbatim, and identical input yields identical output.
- [x] QA prompt tests assert the report destination, the verdict frontmatter requirement, and the never-commit rule appear.
- [x] Neither builder imports the Spec parser or any store/daemon package.

## Verification

- `rtk go test ./internal/agent/` — expected: all tests pass.
- `rtk go build ./...` — expected: builds cleanly.

## References

`_prd.md` → Core Feature 4; Decisions (minimal dedicated builder). `_techspec.md` → Interfaces (prompt builders), Build Order 5, Risks (prompt-contract drift). ADR-0013, ADR-0014.

## Result

`internal/agent` gains two prompt builders beside the review `BuildPrompt`:
`BuildTaskPrompt(TaskPromptRequest) (string, error)` builds the Spec Task
prompt from plain strings (SpecSlug, TaskID, TaskPath, TaskContent) with the
full task file embedded verbatim, and `BuildQAPrompt(specSlug, specDir,
prdPath string) (string, error)` builds the QA gate prompt. Both validate
required fields (empty or whitespace-only input → error, empty prompt) and
are pure functions of their inputs, so identical input yields identical
output. The execution invariants mirroring implement-task live in the single
commented constant `taskExecutionInvariants` and the qa-gate rules in the
single commented constant `qaGateContract` (both in
`internal/agent/spec_prompt.go`), so the work-plan item 5 templating pass
replaces each contract in one place.

Commands run:

- `rtk go test ./internal/agent/` — pass (37 tests, including the 7 new
  prompt-builder tests).
- `rtk go build ./...` — builds cleanly.
- `make verify` — pass (fmt-check, 318 tests in 16 packages,
  `roundfix skills check`, build).

Evidence per acceptance criterion:

- Invariants, verbatim embedding, determinism —
  `TestBuildTaskPromptStatesExecutionInvariants` asserts every invariant
  line from Requirements 2–3 (slice-only scope, `status: in_progress` on
  start, run `## Verification` commands, append `## Result` with evidence,
  settle `completed`/`failed`, never commit/push/PR, never edit `_tasks.md`
  or other task files, stale `in_progress` → dead prior Run → start fresh
  under the work-target lock); `TestBuildTaskPromptEmbedsTaskContentVerbatim`
  asserts the full task file content appears unchanged;
  `TestBuildTaskPromptDeterministicForIdenticalInput` asserts byte-equal
  output for identical input.
- QA prompt contract — `TestBuildQAPromptStatesQAGateContract` asserts the
  `qa/` report destination (`qa-report-YYYY-MM-DD.md`), the
  `verdict: pass, fail, or partial` frontmatter requirement, and the
  never-commit rule.
- Import scope — `internal/agent/spec_prompt.go` imports only `errors`,
  `fmt`, and `strings`; an `rg` sweep of `internal/agent/` for
  `roundfix/internal/spec`, `roundfix/internal/store`, and
  `roundfix/internal/daemon` returns zero matches.

Follow-ups: none new — the prompt-contract drift risk stays open by design
until the work-plan item 5 templating pass (already tracked in
`_techspec.md` Risks).
