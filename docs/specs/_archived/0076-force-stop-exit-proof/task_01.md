---
task: task_01
spec: 0076-force-stop-exit-proof
status: completed
type: test
complexity: medium
---

# Task 01: Keep the signal-ignoring helper alive

## Overview

The `ignore` helper mode exists so a process refuses `SIGTERM` and forces the
controller to escalate. It blocks on `select {}`, which Go's runtime reports as
`all goroutines are asleep - deadlock!` and kills with exit 2, so the process
dies of its own accord before any escalation happens. This slice makes it stay
alive, and adds the direct observation that would have caught it.

Verifiable on its own: run the helper binary alone and watch it live.

## Requirements

1. MUST make `ignore` mode block in a way the Go runtime does not classify as a
   deadlock, so the process survives until something kills it.
2. MUST keep `ignore` mode ignoring `SIGTERM`; the mode's purpose is unchanged.
3. MUST NOT block on a sleep, a deadline, or any wall-clock duration. This
   repository has paid for load-sensitive timing three times in one session, and
   a timer would defeat the deadlock detector while reintroducing that class.
4. MUST add a test that runs the helper alone, observes readiness, and asserts
   the process is still alive afterwards and emitted no `fatal error`.
5. MUST leave the `terminate` mode and every production symbol in
   `internal/store` unchanged.

## Subtasks

- [ ] Replace the `ignore` mode block with a registered signal source.
- [ ] Add the direct liveness observation of the helper binary.
- [ ] Confirm the helper emits no runtime fatal error.

## Acceptance Criteria

- [ ] Run alone in `ignore` mode, the helper prints readiness and is still
      alive afterwards.
- [ ] The helper's output contains no `fatal error` and no
      `all goroutines are asleep`.
- [ ] The helper still ignores `SIGTERM`: sending it does not end the process.
- [ ] No production file under `internal/store` is modified.
- [ ] No sleep, timer, or deadline is introduced into the helper's block.

## Context

- interface: `internal/store/process_unix_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/store -count=1 -run 'Helper|Alive|Liveness' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the liveness observation ran and passed.
- `go test -c -buildvcs=false -o /tmp/rf-0076-store.test ./internal/store && ROUNDFIX_OWNER_PROCESS_HELPER=ignore /tmp/rf-0076-store.test -test.run='^TestOwnerProcessHelper$' -test.timeout=5s 2>&1 | grep -q "all goroutines are asleep" && exit 1 || exit 0`
  — expected: exit 0; the helper no longer dies of a runtime deadlock.
- `if grep -n "select {}" internal/store/process_unix_test.go | grep -q .; then exit 1; fi`
  — expected: exit 0; the unbounded empty select is gone.
- `if git diff --name-only HEAD | grep "^internal/store/" | grep -v "_test.go" | grep -q .; then exit 1; fi`
  — expected: exit 0; no production file changed.

## References

- `_prd.md` → Core Feature 1; Goals.
- `_techspec.md` → Interfaces; Build Order 1; Risks & Considerations.
- ADR-0089.

## Result

### Implementation

- The `ignore` helper now registers `SIGUSR1` with `signal.Notify`, keeps
  ignoring `SIGTERM`, and blocks on the registered signal channel without a
  sleep, timer, or deadline. A liveness probe emits `alive` and returns to the
  same signal-backed block.
- `TestOwnerProcessHelperIgnoreModeStaysAlive` launches the helper binary by
  itself, consumes `ready`, sends `SIGTERM`, waits for the signal-backed
  liveness acknowledgement, asserts that the process has not exited, and
  rejects `fatal error` or `all goroutines are asleep` in captured stderr.

### Focused checks

- Before the helper change,
  `rtk proxy env GOCACHE=<worktree>/.gocache GOFLAGS=-buildvcs=false go test ./internal/store -count=1 -run '^TestOwnerProcessHelperIgnoreModeStaysAlive$' -v`
  failed at `owner process helper did not acknowledge liveness`, establishing
  the regression signal.
- After the helper change, the same focused test passed.
- `rtk proxy env GOCACHE=<worktree>/.gocache GOFLAGS=-buildvcs=false go test ./internal/store -count=1 -run '^TestOwnerProcessControllerForceKillExitProof$' -v`
  passed, exercising the existing controller escalation against the
  signal-ignoring helper.
- The first focused-test attempt used the default Go cache and stopped before
  compilation because the sandbox denied access to
  `/Users/marcio/Library/Caches/go-build`; the repository-local cache resolved
  that environment-only failure.

### Acceptance evidence

- Readiness and continued liveness: the focused liveness test consumed
  `ready`, received `alive` after its probe, observed no exit, and passed.
- No runtime fatal output: the focused liveness test captured stderr after
  reaping the helper and rejected both forbidden diagnostics; it passed.
- `SIGTERM` remains ignored: the liveness acknowledgement arrived only after
  the test sent `SIGTERM`, and the existing force-kill proof passed.
- Production scope: final status inspection listed only this Task file and
  `internal/store/process_unix_test.go`; no production file under
  `internal/store` changed.
- No wall-clock block: diff inspection shows the `ignore` branch uses
  `signal.Notify` plus a channel range and introduces no sleep, timer, or
  deadline.

The Daemon-owned `## Verification` commands were not run in this Agent turn.
