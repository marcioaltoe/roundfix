---
task: task_04
spec: 0091-a-proof-that-can-refuse
status: completed
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

- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestDisposableSessionClose' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestDisposableSessionCloseIsNotAppendedWhenTheSessionNeverOpened'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestDisposableSessionClose' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestDisposableSessionCloseIsAppendedWhenAnOpenSessionWillNotClose'` — expected: exits 0.

## References

- `_prd.md` → Goal 4.
- `_techspec.md` → Build Order 4.
- `docs/backlog/2026-08-08-a-failed-proof-appends-a-cleanup-error-the-maintainer-cannot-act-on.md`

## Result

### Implementation

- A close command that exits with acpx's missing-session status now preserves
  that typed reason through `CloseSession`.
- Disposable proof cleanup records an absent Session through the runner's
  warning sink and omits it from the returned proof diagnosis. Every other
  close failure remains an `AgentSessionCleanupError` in the returned chain.

### Focused checks

- Red signal before the production change:
  `GOCACHE=/private/tmp/roundfix-task04-gocache rtk proxy go test ./internal/agent -run '^TestDisposableSessionCloseIsNotAppendedWhenTheSessionNeverOpened$' -count=1 -v`
  failed because the exact leading selection diagnosis was followed by the
  missing-session close error.
- `GOCACHE=/private/tmp/roundfix-task04-gocache rtk go test ./internal/agent -run '^TestDisposableSessionCloseIsNotAppendedWhenTheSessionNeverOpened$' -count=1`
  passed (1 test).
- `GOCACHE=/private/tmp/roundfix-task04-gocache rtk go test ./internal/agent -run '^TestDisposableSessionCloseIsAppendedWhenAnOpenSessionWillNotClose$' -count=1`
  passed (1 test).
- `GOCACHE=/private/tmp/roundfix-task04-gocache rtk go test ./internal/agent -run '^(TestProveExactSelectionCleanupJoinedFailure|TestACPXProbeCleanupFailureJoinsSelectionError|TestACPXProbeFallbackReportsCleanupInfrastructureError|TestACPXCloseSessionReturnsCloseFailure|TestDisposableSessionCloseIsNotAppendedWhenTheSessionNeverOpened|TestDisposableSessionCloseIsAppendedWhenAnOpenSessionWillNotClose)$' -count=1`
  passed (6 tests).
- `GOCACHE=/private/tmp/roundfix-task04-gocache rtk go test ./internal/agent -count=1`
  passed (302 tests).
- `rtk git diff --check` passed. `rtk git diff --name-only` listed only
  `internal/agent/acpx_runner.go`, `internal/agent/acpx_runner_test.go`, and this
  Task file.

### Acceptance criteria evidence

- **Failure before the Session exists:**
  `TestDisposableSessionCloseIsNotAppendedWhenTheSessionNeverOpened` asserts
  the complete returned diagnosis equals the pre-change leading diagnosis and
  contains no `AgentSessionCleanupError`.
- **Open Session will not close:**
  `TestDisposableSessionCloseIsAppendedWhenAnOpenSessionWillNotClose` proves
  ensure and selection ran before close, then asserts the returned chain still
  contains `AgentSessionCleanupError` after the unchanged leading diagnosis.
- **Close failure recorded in both cases:** the missing-session case captures
  the warning sink record, while the genuine close failure is observable via
  `errors.As` in the returned chain.
- **Leading diagnosis unchanged:** both regression tests compare the exact
  existing selection-diagnosis bytes before checking the cleanup suffix.

The authored `## Verification` commands were not run; the Daemon owns them.
