---
task: task_01
spec: 0016-worktree-bootstrap
status: pending
type: backend
complexity: medium
---

# Task 01: Bootstrap config and Run Worktree bootstrap with failure mapping

## Overview

Add the bootstrap primitive and wire it into Run Worktree creation: a
`worktree.bootstrap` command Roundfix runs once in a new worktree after the
`worktree.copy` step and before Agent work, bounded by a timeout, with a
bootstrap failure ending the Run with a clear message. This covers sequential
Runs and review Runs (which use the Run Worktree).

## Requirements

1. MUST add `worktree.bootstrap` (string; empty = no bootstrap, today's
   behavior) and `worktree.bootstrap_timeout` (duration, default `10m`) to
   config, strict-decoded with Project > User > builtin precedence.
2. MUST run the bootstrap command in the new worktree root immediately after the
   `worktree.copy` placement and before the worktree is used for Agent work,
   under the configured timeout.
3. MUST, on a non-zero exit, start failure, or timeout, return a distinct
   bootstrap error and end the owning Run before any Batch is assigned, with
   stderr carrying `worktree bootstrap failed: <command>: <reason>`; bootstrap
   output MUST be journaled/streamed like other Run diagnostics.
4. MUST leave Run behavior byte-stable when `worktree.bootstrap` is empty.

## Subtasks

- [ ] `worktree.bootstrap` + `worktree.bootstrap_timeout` config keys
- [ ] `runBootstrap` in internal/worktree: run command in worktree root, timeout, BootstrapError
- [ ] Call after CopyList placement in Run Worktree creation
- [ ] Map bootstrap failure to a Run-ending outcome with the message
- [ ] Tests: success, non-zero exit, timeout, empty-bootstrap byte-stable

## Acceptance Criteria

- [ ] With `worktree.bootstrap` set, the command runs in the new Run Worktree root after copy and before Agent work.
- [ ] A non-zero bootstrap exit ends the Run before any Batch, with `worktree bootstrap failed: <command>: <reason>` on stderr.
- [ ] A bootstrap command exceeding `bootstrap_timeout` fails with a timeout reason.
- [ ] With an empty `worktree.bootstrap`, Run behavior is unchanged (existing tests pass).
- [ ] Config loads both keys with correct precedence; an unknown `worktree` key still fails strict validation.

## Verification

- `rtk go test ./internal/worktree/ ./internal/config/ ./internal/cli/` — expected: bootstrap and config tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 1-4; Core Features 1-2, 4. `_techspec.md` → Where
bootstrap runs, Failure handling, Build Order 1, Interfaces: `runBootstrap`.
ADR-0034. Work-plan bootstrap finding.
