---
task: task_02
spec: 0035-agent-selection-profiles
status: completed
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

- [x] Introduce the closed Task Type parser and actionable error contract.
- [x] Expose the parsed type on loaded Tasks.
- [x] Validate the complete Task Graph before operational side effects.
- [x] Pin all accepted and rejected authoring cases.
- [x] Align the `write-tasks` approval and template contracts.
- [x] Regenerate and verify the embedded skill copy.

## Acceptance Criteria

- [x] All seven canonical values load and remain stable through the Task Graph.
- [x] Missing, padded, mixed-case, and unknown values fail with the exact file and allowed values.
- [x] Invalid Task Type produces zero model probes, Run rows, branches, worktrees, and Agent invocations.
- [x] Generated task frontmatter and the `_tasks.md` projection contain one identical allowed type.
- [x] The `write-tasks` skill contains no recommendation, runtime id, model id, or profile configuration logic.
- [x] Canonical and embedded `write-tasks` copies are byte-identical.

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

## Result

Implemented the closed Task Type contract in `internal/spec`: loaded tasks now expose a parsed `TaskType`, task file frontmatter rejects missing, empty, padded, mixed-case, and unknown values with the file path, invalid value, full allowed set, and frontmatter correction action, and `_tasks.md` projection rows are validated against task file types when present. Implement preflight now fails invalid Task Types before selection probes, Run persistence, worktree creation, or Agent execution because `spec.Load` rejects the graph before those operational steps.

Updated the repo-owned `write-tasks` skill and task template so Task Type remains based on the dominant delivered outcome, dependencies remain only in `_tasks.md`, status/type remain task-frontmatter owned, and the projection table must carry the identical allowed value. Synchronized `.agents/skills/write-tasks` with `skills/write-tasks`, and refreshed setup-context-driven skill snapshot digests so the full skill audit accepts the updated authoring contract.

Evidence:

- Canonical values: `TestTaskTypeCanonicalValuesLoadThroughTaskGraph` covers `backend`, `frontend`, `data`, `infra`, `docs`, `test`, and `chore` through graph loading.
- Invalid values: `TestTaskTypeRejectsInvalidFrontmatterValues` covers missing, empty, padded, mixed-case, and unknown values with actionable diagnostics.
- Zero side effects: `TestImplementRejectsInvalidTaskTypeBeforeSideEffects` asserts no probes, fallback probes, Run Database, Run Worktree root, or Agent calls.
- Projection identity: `TestTaskTypeProjectionMustMatchTaskFile` rejects mismatched `_tasks.md` type projection versus task frontmatter.
- Skill policy: `rtk rg -n "codex|claude|gpt|opus|model|runtime|recommend|profile" .agents/skills/write-tasks skills/write-tasks` returned no matches.
- Verification passed: `GOCACHE=/private/tmp/roundfix-gocache rtk go test ./internal/spec ./internal/cli -run 'Test(TaskType|ImplementRejectsInvalidTaskType)' -count=1` → 10 tests passed.
- Verification passed: `cmp .agents/skills/write-tasks/SKILL.md skills/write-tasks/SKILL.md && cmp .agents/skills/write-tasks/references/task-template.md skills/write-tasks/references/task-template.md` → no diff.
- Verification passed: `rtk make skills-sync-check` → no drift output.
- Additional package check passed: `GOCACHE=/private/tmp/roundfix-gocache rtk go test ./internal/spec ./internal/cli -count=1` → 578 tests passed.
- Full gate passed: `GOCACHE=/private/tmp/roundfix-gocache rtk make verify` → Go tests, setup-context checks, skills check, and build passed.
