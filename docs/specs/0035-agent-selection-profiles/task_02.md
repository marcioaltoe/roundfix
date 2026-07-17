---
task: task_02
spec: 0035-agent-selection-profiles
status: pending
type: backend
complexity: medium
---

# Task 02: Enforce author-declared Task Types

## Overview

Make Task Type a closed authoring and runtime-routing contract before any profile or Run side effect can occur. The slice updates the repo-owned authoring skill and the Spec parser together so generated Tasks and operational validation cannot drift.

## Requirements

1. MUST parse exactly `backend`, `frontend`, `data`, `infra`, `docs`, `test`, and `chore` as Task Types and reject empty, mixed-case, whitespace-padded, or unknown values.
2. MUST include the task file, invalid value, complete allowed set, and frontmatter correction action in every Task Type error.
3. MUST validate all requested Task files before profile resolution, disposable proof, Run persistence, branch/worktree creation, or Agent Session start.
4. MUST keep dependencies only in `_tasks.md`, status and type only in task frontmatter, and the manifest projection type identical to the file type.
5. MUST keep type selection based on the dominant delivered outcome and contain no runtime, model, recommendation, or profile configuration policy in `write-tasks`.
6. MUST synchronize the canonical and embedded repo-owned `write-tasks` skill copies.

## Subtasks

- [ ] Introduce the closed Task Type parser and actionable error contract.
- [ ] Expose the parsed type on loaded Tasks.
- [ ] Validate the complete Task Graph before operational side effects.
- [ ] Pin all accepted and rejected authoring cases.
- [ ] Align the `write-tasks` approval and template contracts.
- [ ] Regenerate and verify the embedded skill copy.

## Acceptance Criteria

- [ ] All seven canonical values load and remain stable through the Task Graph.
- [ ] Missing, padded, mixed-case, and unknown values fail with the exact file and allowed values.
- [ ] Invalid Task Type produces zero model probes, Run rows, branches, worktrees, and Agent invocations.
- [ ] Generated task frontmatter and the `_tasks.md` projection contain one identical allowed type.
- [ ] The `write-tasks` skill contains no recommendation, runtime id, model id, or profile configuration logic.
- [ ] Canonical and embedded `write-tasks` copies are byte-identical.

## Context

- instruction: `.agents/skills/write-tasks/SKILL.md`
- instruction: `.agents/skills/write-tasks/references/task-template.md`
- interface: `internal/spec/task.go`
- interface: `internal/spec/spec.go`
- interface: `internal/cli/implement.go`

## Verification

- `rtk go test ./internal/spec ./internal/cli -run 'Test(TaskType|ImplementRejectsInvalidTaskType)' -count=1` — expected: accepted values and zero-side-effect rejection cases pass.
- `cmp .agents/skills/write-tasks/SKILL.md skills/write-tasks/SKILL.md && cmp .agents/skills/write-tasks/references/task-template.md skills/write-tasks/references/task-template.md` — expected: canonical and embedded authoring contracts are byte-identical.
- `rtk make skills-sync-check` — expected: every repo-owned skill copy remains synchronized.

## References

- `_prd.md` → Goals 7; User Stories 1 and 5; Core Feature 7; Success Metrics.
- `_techspec.md` → Domain types: TaskType; `write-tasks` contract; Operational preflight steps 1-3; Build Order 2.
