---
task: task_02
spec: 0076-force-stop-exit-proof
status: completed
type: test
complexity: medium
---

# Task 02: Prove the escalation caused the exit

## Overview

With the helper alive, two things remain. The parent starts its `cmd.Wait`
goroutine before scanning the readiness line, so an early exit can close the
stdout pipe first and surface as `read |0: file already closed`. And
`TestOwnerProcessControllerForceKillExitProof` asserts only that the process is
gone — a condition a process that crashed on its own also satisfies, which is
exactly how this defect survived green.

This slice reorders the handshake and makes the proof assert causation.

## Requirements

1. MUST scan the helper's readiness line before starting the goroutine that
   observes process exit, so neither can close a pipe the other still needs.
2. MUST fix the ordering rather than widen a timeout. A budget that hides a
   race is the same defect one step later.
3. MUST make the force-kill proof assert that the controller's escalation ended
   the process, not merely that the process is gone.
4. MUST add a regression that fails when the helper exits before the controller
   acts, naming the premature exit.
5. MUST prove that regression can fail, by making the helper exit prematurely,
   observing the failure, and reverting within the same check.
6. MUST leave `TerminateAndWait` and every production symbol unchanged.

## Subtasks

- [ ] Move readiness consumption ahead of the exit-observation goroutine.
- [ ] Assert the escalation caused the exit in the force-kill proof.
- [ ] Add the premature-exit regression.
- [ ] Run the negative probe and revert it.

## Acceptance Criteria

- [ ] The parent never fails with `file already closed`, across the repeated
      run below.
- [ ] `TestOwnerProcessControllerForceKillExitProof` asserts the controller's
      escalation ended the process, not only that it is gone.
- [ ] A helper that exits before the controller acts fails the regression with
      a message naming the premature exit.
- [ ] That regression was observed failing under a deliberate premature exit
      and passing after the revert, both recorded in the Result.
- [ ] Parent and helper together pass at `-count=50` with no child emitting a
      `fatal error`.
- [ ] No production file under `internal/store` is modified.

## Context

- interface: `internal/store/process_unix_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/store -count=1 -run 'OwnerProcess' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the owner-process suite ran and passed.
- `go test ./internal/store -count=50 -run 'OwnerProcessControllerForceKillExitProof' 2>&1 | grep -qE "FAIL|fatal error|file already closed" && exit 1 || exit 0`
  — expected: exit 0; fifty repetitions with no failure, no runtime fatal error,
  and no closed-pipe error.
- `go test ./internal/store -count=1` — expected: exit 0.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0; the full parallel sweep is green, which is where the
  original failure surfaced.
- `if git diff --name-only HEAD | grep "^internal/store/" | grep -v "_test.go" | grep -q .; then exit 1; fi`
  — expected: exit 0; no production file changed.

## References

- `_prd.md` → Core Features 2, 3, and 4; Decisions.
- `_techspec.md` → Interfaces; Build Order 2 and 3; Testing Approach.
- ADR-0089.

## Result

### Implementation

- `startOwnerProcessHelper` now consumes and validates the helper's `ready`
  line before starting the `cmd.Wait` goroutine. Cleanup waits synchronously
  when readiness fails before that goroutine starts, so every child is still
  reaped without letting `Wait` close a pipe the scanner needs.
- `TestOwnerProcessControllerForceKillExitProof` now consumes the helper's wait
  result and requires a signaled `syscall.WaitStatus` whose signal is
  `SIGKILL`. A clean early exit fails with
  `exited prematurely before controller force-kill escalation` instead of
  satisfying the proof.
- `TerminateAndWait` and all production symbols remain unchanged.

### Focused checks

- Before the assertion change, a deliberate `return` immediately after the
  helper's `ready` line left
  `rtk proxy env GOCACHE=<worktree>/.gocache GOFLAGS=-buildvcs=false go test ./internal/store -run '^TestOwnerProcessControllerForceKillExitProof$' -count=1 -v`
  green. This reproduced the false-positive defect; the mutation was reverted
  before implementation.
- With the new assertion, the same deliberate premature return made that
  focused command exit 1 at
  `owner process ... exited prematurely before controller force-kill escalation: <nil>`.
  The mutation was reverted in the same check, and the identical command then
  exited 0 with the force-kill proof passing.
- `rtk proxy env GOCACHE=<worktree>/.gocache GOFLAGS=-buildvcs=false go test ./internal/store -run '^TestOwnerProcessController(GracefulExitProof|ForceKillExitProof)$' -count=20`
  exited 0. All 20 parent/helper repetitions completed without `file already
  closed`, a child runtime fatal exit, or another failure.
- `rtk proxy env GOCACHE=<worktree>/.gocache GOFLAGS=-buildvcs=false go test -race ./internal/store -run '^TestOwnerProcessControllerForceKillExitProof$' -count=1`
  exited 0 after the goroutine-ordering change.
- `rtk gofmt -w internal/store/process_unix_test.go` completed without output.

### Acceptance evidence

- Closed-pipe ordering: diff inspection places readiness scanning before the
  `cmd.Wait` goroutine, and the focused 20-repeat parent/helper check passed
  without `file already closed`. The authored 50-repeat command remains for
  Daemon Verification.
- Escalation causation: the force-kill proof now requires the child process's
  observed exit signal to equal `SIGKILL`; its focused positive and race checks
  passed.
- Premature-exit diagnostic: the deliberate clean helper exit failed the proof
  with a diagnostic containing `exited prematurely`.
- Mutation sensitivity: the deliberate premature exit was observed passing
  before the regression, failing after it, and passing again after the helper
  mutation was reverted.
- Repeated parent/helper behavior: 20 focused repetitions passed without a
  closed-pipe failure or child fatal exit. The acceptance criterion's exact
  `-count=50` evidence remains Daemon-owned Verification and is not claimed
  here.
- Production scope: diff inspection shows the only `internal/store` change is
  `internal/store/process_unix_test.go`; no production file was modified.

The Daemon-owned `## Verification` commands were not run in this Agent turn.
