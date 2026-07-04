---
task: task_03
spec: 0001-implement-command
status: pending
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

- [ ] Task prompt builder with embedded task content
- [ ] Execution invariants constant mirroring implement-task
- [ ] QA prompt builder mirroring qa-gate
- [ ] Unit tests asserting each invariant line and determinism

## Acceptance Criteria

- [ ] Task prompt tests assert every invariant from Requirements 2–3 appears, the task content is embedded verbatim, and identical input yields identical output.
- [ ] QA prompt tests assert the report destination, the verdict frontmatter requirement, and the never-commit rule appear.
- [ ] Neither builder imports the Spec parser or any store/daemon package.

## Verification

- `rtk go test ./internal/agent/` — expected: all tests pass.
- `rtk go build ./...` — expected: builds cleanly.

## References

`_prd.md` → Core Feature 4; Decisions (minimal dedicated builder). `_techspec.md` → Interfaces (prompt builders), Build Order 5, Risks (prompt-contract drift). ADR-0013, ADR-0014.
