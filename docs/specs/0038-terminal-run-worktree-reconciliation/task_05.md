---
task: task_05
spec: 0038-terminal-run-worktree-reconciliation
status: completed
type: backend
complexity: medium
---

# Task 05: Surface retained terminal Run Worktrees in Runs List

## Overview

Make retained terminal Run Worktrees discoverable without changing the stable
Runs List stdout rows. Each view reports the exact relevant count on stderr and
points to the Reconcile Command as the sole classification surface.

## Requirements

1. MUST preserve existing Runs List stdout row shape byte-for-byte.
2. MUST count terminal spec Runs that retain a recorded worktree path or Run
   Branch.
3. MUST report hidden retained residue when the default Active view has no
   matching terminal rows.
4. MUST report retained residue represented in terminal and all-state views.
5. MUST write the bounded note to stderr only.
6. MUST point to `roundfix reconcile` without classifying safety in Runs List.
7. MUST omit the note when no retained terminal surface exists.

## Subtasks

- [x] Derive repository-scoped retained-worktree counts.
- [x] Add the exact stderr guidance.
- [x] Preserve stdout row serialization.
- [x] Cover Active, terminal, and all-state views.
- [x] Cover zero-retention and mixed-repository cases.

## Acceptance Criteria

- [x] Active view reports the exact hidden terminal residue count on stderr.
- [x] Terminal and all-state views report only retained Runs relevant to their
      repository scope.
- [x] Existing stdout rows remain byte-identical.
- [x] The note contains the exact `roundfix reconcile` pointer.
- [x] Runs List never labels a retained entry safe or unsafe.
- [x] No note is printed when neither worktree nor Run Branch remains.

## Context

- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `.agents/skills/golang-cli/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- interface: `internal/cli/runs.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/store/store.go`
- interface: `internal/worktree/worktree.go`

## Verification

- `rtk go test ./internal/cli -run 'TestRun.*List.*Retained.*Worktree' -count=1`
  — expected: all views report exact stderr counts while stdout stays
  byte-stable.
- `rtk go test ./internal/cli -run 'TestRun.*List.*(Active|Terminal|All)' -count=1`
  — expected: existing Runs List view contracts remain passing.

## References

- `_prd.md` → Goal 3; User Story 3; Core Feature 8; Success Metrics.
- `_techspec.md` → API Contracts: Run listing; Testing Approach; Build Order 5.
- `CONTEXT.md` → Runs List and Reconcile Command vocabulary.

## Result

Runs List now derives one exact retained terminal Run Worktree count from its
unbounded, repository-scoped Run query. The shared proof-based inspector
distinguishes retained terminal spec Runs from `released` history, including
path-and-branch, path-only, and Run-Branch-only residue. When the count is
non-zero, stderr prints one bounded note that points to
`roundfix reconcile`; stdout row formatting and empty-result output are
unchanged. When no retained surface exists, the existing hidden-row guidance
remains the fallback.

Verification:

- `rtk proxy env GOCACHE=/private/tmp/roundfix-task05-gocache go test ./internal/cli -run 'TestRun.*List.*Retained.*Worktree' -count=1`
  — passed.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task05-gocache go test ./internal/cli -run 'TestRun.*List.*(Active|Terminal|All)' -count=1`
  — passed.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task05-gocache go test ./internal/cli -count=1`
  — passed.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task05-gocache make verify`
  — passed.
- The first focused invocation without the task-local `GOCACHE` did not reach
  compilation because the sandbox denied the host Go cache; the equivalent
  task-local-cache invocation above passed.

Acceptance evidence:

- `TestRunRunsListActiveReportsRetainedWorktreesWithoutChangingStdout` creates
  path-and-branch, branch-only, path-only, and released terminal spec Runs. The
  Active view reports exactly three retained Runs on stderr and matches the
  pre-existing stdout row serialization byte-for-byte.
- `TestRunRunsListTerminalAndAllReportRetainedWorktreesByRepository` proves
  terminal and all-state views count one retained Run in the current
  repository and two with `--all`, while excluding the other repository from
  repository-scoped stdout.
- The retained note is exactly
  `(N terminal Run Worktree(s) retained; run 'roundfix reconcile' to inspect)`
  with grammatical singular/plural rendering. Tests reject the words `safe`
  and `unsafe` and confirm the pointer never enters stdout.
- `TestRunRunsListRetainedWorktreeNoteOmittedWhenReleased` removes both the Run
  Worktree and Run Branch, keeps the terminal history row visible, and observes
  empty stderr.

Follow-ups: none.
