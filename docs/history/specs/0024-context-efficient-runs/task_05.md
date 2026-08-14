---
task: task_05
spec: 0024-context-efficient-runs
status: completed
type: backend
complexity: high
---

# Task 05: Build bounded Spec Context Bundles

## Overview

Give each Task a deterministic path manifest plus exactly one embedded full
Task contract, without embedding larger Spec, instruction, source, or diff
content. The slice is verifiable against sequential and parallel Git histories,
external Spec Roots, explicit context validation, and the 200-path ceiling.

## Requirements

1. MUST parse optional labeled Task context entries for instruction and interface paths without changing frontmatter or Task Graph ownership.
2. MUST reject external, unclean, unknown-kind, or more than 50 unique explicit context entries as typed Task-file validation errors.
3. MUST include standard Spec paths, root Agent instructions, the implement-task Skill, explicit Task context, and sorted prior changed files by path only.
4. MUST derive prior files from the persisted Run HEAD to the current Task worktree HEAD so only already integrated Task changes appear.
5. MUST cap the complete manifest at 200 paths, reserve standard/explicit paths first, and report the omitted prior-file count.
6. MUST embed the complete assigned Task exactly once and embed no full PRD, TechSpec, Skill, source file, or prior diff.
7. MUST regenerate the same bounded bundle for the next Task after session resume or replacement.

## Subtasks

- [ ] Add typed Task Context parsing and validation.
- [ ] Add the Git prior-changed-path resolver.
- [ ] Assemble deterministic bounded bundle categories.
- [ ] Extend Task prompts with the manifest and one full Task.
- [ ] Handle sequential, parallel, external-Spec, and replacement-session paths.
- [ ] Add parser, Git, bound, and prompt regression tests.

## Acceptance Criteria

- [ ] Missing `## Context` remains valid, while labeled valid entries load into their correct categories.
- [ ] Invalid paths, labels, or a 51st unique entry fail with a Task-specific validation error.
- [ ] Standard and explicit paths are retained before sorted prior files fill the remaining 200-path capacity.
- [ ] The bundle reports the exact number of omitted prior files when the ceiling is exceeded.
- [ ] A parallel Task sees files integrated before its Task Worktree base and no unintegrated sibling diff.
- [ ] The prompt contains one complete assigned Task and zero complete larger documents, source files, or prior diffs.
- [ ] A resumed/replacement session's next Task receives the same deterministic bundle for the same repository state.

## Verification

- `rtk go test ./internal/spec ./internal/agent ./internal/daemon` - expected: context parsing, path validation, Git history, bounds, concurrency, and prompt tests pass.
- `rtk go test -race ./internal/daemon` - expected: parallel Task Context Bundle construction passes under the race detector.
- `rtk make verify` - expected: formatting, all tests, Skill sync checks, Skill checks, and build pass.

## Context

- instruction: `.agents/skills/golang-testing/SKILL.md`
- instruction: `.agents/skills/golang-concurrency/SKILL.md`
- interface: `internal/spec/task.go`
- interface: `internal/spec/spec.go`
- interface: `internal/agent/spec_prompt.go`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/worktree/worktree.go`

## References

`_prd.md` -> User Story 5; Core Features 10-11; Success Metrics. `_techspec.md` -> Data Models: Task Context references; API Contracts: Spec Context Bundle; Build Order 5. ADR-0035.

## Result

- Missing `## Context` remains valid and labeled entries load into instruction/interface categories: covered by `TestLoadParsesOptionalTaskContext`.
- Invalid labels, external/escaping or unclean paths, and the 51st unique entry fail as typed Task-file validation errors: covered by `TestLoadRejectsInvalidTaskContext`.
- Standard and explicit paths are retained before sorted prior files fill the 200-path capacity: covered by `TestAssembleTaskContextBundleReservesExplicitPathsAndCountsOmittedPriorFiles`.
- The omitted prior-file count is exact when the ceiling is exceeded: covered by `TestAssembleTaskContextBundleReservesExplicitPathsAndCountsOmittedPriorFiles`.
- A parallel Task sees files integrated before its Task Worktree base and excludes an unintegrated sibling diff: covered by `TestPriorChangedFilesUseCurrentWorktreeHeadAndIgnoreSiblingBranch` and `TestTaskCycleParallelTaskPromptUsesTaskWorktreeContextBase`.
- The prompt contains one complete assigned Task and no complete larger documents, source files, or prior diffs: covered by `TestBuildTaskPromptRendersSpecContextBundlePathOnly` and `TestTaskCyclePromptContainsBundleWithoutReferencedBodies`.
- A resumed or replacement session receives the same deterministic bundle for the same repository state: covered by `TestAssembleTaskContextBundleIsDeterministic`.

Verification:

- `rtk go test ./internal/spec ./internal/agent ./internal/daemon` passed: 251 tests in 3 packages.
- `rtk go test -race ./internal/daemon` passed: 81 tests in 1 package.
- `rtk make verify` passed: `go test ./...` reported 1097 tests in 19 packages, the Roundfix skill check passed, and the build completed.
