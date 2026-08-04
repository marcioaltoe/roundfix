---
task: task_01
spec: 0076-force-stop-exit-proof
status: pending
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
