---
task: task_04
spec: 0091-a-proof-that-can-refuse
status: pending
type: backend
complexity: low
---

# Task 04: Stop appending a close error for a session never opened

## Overview

When a proof fails, the real diagnosis is followed by a line about failing to
close a disposable Agent Session that was never created, ending an actionable
message in noise about something the maintainer did not do. Reproduced on
2026-08-09 while proving a nonexistent `codex` model.

## Requirements

1. MUST NOT append a session-close failure to a proof diagnosis when the session
   was never created.
2. MUST keep recording that close failure, so a genuine cleanup problem is not
   lost; recording it is not the same as putting it in the maintainer's message.
3. MUST keep appending a close failure for a session that was created and could
   not be closed, which is a real cleanup problem the maintainer can act on.
4. MUST leave the leading diagnosis — classification, adapter error, recovery,
   next step and fallback — byte-identical.

## Subtasks

- [ ] Distinguish a session that was never created from one that would not close.
- [ ] Record the first, append the second.

## Acceptance Criteria

- [ ] A proof that fails before the session exists reports the diagnosis with no
      trailing close error.
- [ ] A proof whose created session fails to close still reports it.
- [ ] The close failure is recorded in both cases.
- [ ] The leading diagnosis text is unchanged.

## Bounded scope

This Task may create or modify only:

- `internal/agent/acpx_runner.go`
- `internal/agent/acpx_runner_test.go`
- `docs/specs/0091-a-proof-that-can-refuse/task_04.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestDisposableSessionClose' -count=1 -v 2>&1 | grep -q '^--- PASS: TestDisposableSessionCloseIsNotAppendedWhenTheSessionNeverOpened'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestDisposableSessionClose' -count=1 -v 2>&1 | grep -q '^--- PASS: TestDisposableSessionCloseIsAppendedWhenAnOpenSessionWillNotClose'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/agent -count=1 2>&1 | grep -q '^ok'` — expected: exits 0.

## References

- `_prd.md` → Goal 4.
- `_techspec.md` → Build Order 4.
- `docs/backlog/2026-08-08-a-failed-proof-appends-a-cleanup-error-the-maintainer-cannot-act-on.md`
