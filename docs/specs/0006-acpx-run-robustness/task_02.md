---
task: task_02
spec: 0006-acpx-run-robustness
status: completed
type: backend
complexity: high
---

# Task 02: Build the Settle Command

## Overview

One-command recovery for a failed Task whose work is preserved in the
working tree: `roundfix settle --spec <slug> --task <task_id>` re-runs the
Task's Verification verbatim and, only on pass, settles `completed` and
creates the standard Task commit. Replaces the hand-played Daemon recovery
used twice this week. Verifiable through buffer-captured CLI tests over real
temp repos.

## Requirements

1. MUST add the `settle` command: `--spec <slug>` and `--task <task_id>`
   both required (no interactive picker — recovery is deliberate); house
   exit codes (0 settled, 1 verification failed, 2 Preflight Validation).
2. MUST run Preflight Validation with one actionable message per failure:
   repository resolves; Spec loads valid; the task id exists in the Task
   Graph; task status is exactly `failed` (other statuses name the status
   and the right path: pending/in_progress → implement, completed → nothing
   to do); no Active Run holds the spec target or the working tree.
3. MUST run the Task's Verification commands verbatim through the existing
   verifier in the repository root, streaming `verify <command> — ok|failed`
   lines on stdout in order, stopping at the first failure.
4. MUST, on all-pass: set the task status to `completed`, stage **all**
   current worktree changes plus the task file, and commit with the standard
   Task derivation (message and trailers byte-identical to a Daemon
   settlement); final stdout line `settled <task_id> completed — <short
   sha>`.
5. MUST, on verification failure: write nothing (status untouched, no
   commit), end stdout with `<task_id> stays failed — verification failed`,
   exit 1.
6. MUST create no Run, write no journal events, never push; help text and
   top-level usage updated truthfully. The stage-everything contract is
   stated in the help text's one-line description of what the commit
   contains.

## Subtasks

- [x] Command skeleton, flags, usage, dispatch
- [x] Preflight Validation with per-failure messages
- [x] Verbatim verification loop with streamed stdout lines
- [x] Settlement: status write, stage-all commit, deterministic report
- [x] Buffer-captured tests over temp repo + store

## Acceptance Criteria

- [x] A fixture reproducing the real incident (failed status, completed work
      in tree) settles: commit message and trailers byte-equal to the
      Daemon's for the same task, staged content includes the worktree
      changes and task file, exit 0.
- [x] Every preflight refusal from Requirement 2 has a test asserting exit 2
      and its message; nothing is written in any refusal or failure path.
- [x] Verification-failure path proves status and tree untouched, exit 1.
- [x] An Active Run on the same working tree blocks settle (exit 2 naming
      the run id).
- [x] `settle --help` documents flags, exit codes, and the stage-everything
      contract; full suite passes.

## Verification

- `rtk go test ./internal/cli/` — expected: all tests pass.
- `rtk go run ./cmd/roundfix settle --help` — expected: usage with both
  flags, exit 0.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 2; Core Feature 2; Decisions (stage-everything, no
Run). `_techspec.md` → Interfaces (settle), Build Order 2, Risks. ADR-0014,
ADR-0016 (target keys). Round-2 dogfood finding 3.

## Result

- Implemented `roundfix settle --spec <slug> --task <task_id>` with required
  flags, exit codes, top-level dispatch, and help text documenting that the
  recovery commit contains all current worktree changes plus the task file.
- Added settle Preflight Validation for repository resolution, valid Spec
  loading, Task Graph membership, exact `failed` status, and Active Run locks
  for both the Spec target and current working tree.
- Added the verbatim Task Verification loop through the existing verifier,
  reporting `verify <command> — ok|failed` on stdout and stopping on the first
  failure.
- Added settlement that rewrites the Task status to `completed`, stages all
  dirty paths plus the task file, commits with `daemon.TaskCommitMessage`, and
  reports `settled <task_id> completed — <short sha>`.
- Evidence: `TestRunSettleCommitsFailedTaskWorktreeWithDaemonMessage` covers
  the incident-shaped failed Task with preserved work, asserts the Daemon
  commit message/trailers byte-equal, asserts the worktree file and task file
  are in the commit, and exits 0.
- Evidence: `TestRunSettlePreflightRefusalsWriteNothing` and
  `TestRunSettleRequiresSpecAndTask` assert exit 2 and messages for every
  Preflight refusal and required flag refusal, with task files, git status,
  commit count, and Run Database state unchanged as applicable.
- Evidence: `TestRunSettleVerificationFailureLeavesTaskAndTreeUntouched`
  proves status and tree unchanged, no later Verification command runs, and
  exit 1 with the required final stdout line.
- Evidence: `TestRunSettleActiveRunOnSameWorkingTreeBlocks` creates an Active
  Run on the same working tree, asserts exit 2 names the Run id and stop
  command, and proves no Run Event Journal entry was written.
- Evidence: `TestRunSettleHelpDocumentsContract` asserts `settle --help`
  documents both flags, exit codes, and the stage-everything contract; it also
  asserts top-level usage lists `settle`.
- Verification passed:
  - `rtk go test ./internal/cli/` — 196 passed in 1 package.
  - `rtk go run ./cmd/roundfix settle --help` — printed usage with both flags
    and exited 0.
  - `rtk go test ./...` — 605 passed in 16 packages.
  - `rtk make verify` — passed (`go test ./...`, `skills check`, and build).
