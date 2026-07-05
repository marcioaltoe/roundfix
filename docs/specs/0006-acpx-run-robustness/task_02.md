---
task: task_02
spec: 0006-acpx-run-robustness
status: pending
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

- [ ] Command skeleton, flags, usage, dispatch
- [ ] Preflight Validation with per-failure messages
- [ ] Verbatim verification loop with streamed stdout lines
- [ ] Settlement: status write, stage-all commit, deterministic report
- [ ] Buffer-captured tests over temp repo + store

## Acceptance Criteria

- [ ] A fixture reproducing the real incident (failed status, completed work
      in tree) settles: commit message and trailers byte-equal to the
      Daemon's for the same task, staged content includes the worktree
      changes and task file, exit 0.
- [ ] Every preflight refusal from Requirement 2 has a test asserting exit 2
      and its message; nothing is written in any refusal or failure path.
- [ ] Verification-failure path proves status and tree untouched, exit 1.
- [ ] An Active Run on the same working tree blocks settle (exit 2 naming
      the run id).
- [ ] `settle --help` documents flags, exit codes, and the stage-everything
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
