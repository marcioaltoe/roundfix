---
task: task_03
spec: 0285-a-transient-profile-proof-timeout-does-not-park-the-delivery-queue
status: completed
type: backend
complexity: medium
---

# Task 03: Cleanup of an absent disposable Agent Session is silent

## Overview

On 2026-10-19 at 22:57 BRT a timed-out profile proof also logged `disposable
Agent Session cleanup found no Session to close`, with a recovery step that
did not apply (`_prd.md` → Sources 1). The setup deadline had stopped
`sessions ensure` before acpx recorded the Session, so the cleanup's `sessions
close` found none, which is what cleanup wants. This Task makes that case
silent. It is verifiable on its own through the compiled fake acpx.

## Requirements

1. MUST answer the 22:57 line of the Secondbrain Inbox Entry of 2026-10-20,
   [`references/2026-10-20-a-prova-de-perfil-do-preflight-estoura-o-tempo-em-serie.md`](references/2026-10-20-a-prova-de-perfil-do-preflight-estoura-o-tempo-em-serie.md)
   (`_prd.md` → Sources 1), by implementing API Contract 3 of `_techspec.md`
   in `closeDisposableSession` (ADR-0298): when `CloseSession` fails and
   `IsAgentSessionAbsent` holds, it returns `nil` and writes no warning.
2. MUST keep every other close failure as today: an
   `AgentSessionCleanupError`, joined with any setup error. The profile
   proof, the Fallback Chain probe and sealed prompts share the helper and
   all get the change; `EndSession`'s own warning does not change.
3. MUST change `TestMissingSessionIsRecognisedFromBothExitShapes` in
   `internal/agent/acpx_runner_test.go`, a declared break: for both exit
   shapes it asserts that no warning is written, instead of exactly one, and
   keeps its other assertions and its name.
4. MUST add `TestTimedOutProofSetupClosesAnAbsentSessionSilently` to the same
   file, its first statement `t.Parallel()`, waiting on no real time:
   `sessions ensure` blocks in the fake acpx, the test-deadline context
   expires after the fixture reports the command started, and `sessions
   close` exits `4`. The proof returns an error matching
   `context.DeadlineExceeded`, no `AgentSessionCleanupError` is in its chain,
   the last invocation is the disposable `sessions close`, and no warning is
   written.
5. MUST search the code, tests, guides and skills for `found no Session to
   close` and leave no occurrence.

## Subtasks

- [ ] Return without a warning when the disposable Session is absent.
- [ ] Update the exit-shape test and add the timed-out setup test.

## Acceptance Criteria

- [ ] A proof whose Session never opened returns its setup error alone and writes no cleanup warning.
- [ ] A close failure for a Session that exists is still reported and joined.

## Consumers

- `closeDisposableSession`'s warnings and errors — consumers: `./internal/agent`

## Context

- instruction: `docs/adr/0298-a-timed-out-profile-proof-gets-a-third-attempt-and-its-setup-deadline-is-configurable.md`
- instruction: `internal/agent/sessions.go`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/acpx_runner_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestTimedOutProofSetupClosesAnAbsentSessionSilently|TestMissingSessionIsRecognisedFromBothExitShapes|TestDisposableSessionCloseIsAppendedWhenAnOpenSessionWillNotClose|TestProveExactSelectionTimeoutCleanup|TestUnrelatedExitOneKeepsItsClassification)$' ./internal/agent 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTimedOutProofSetupClosesAnAbsentSessionSilently TestMissingSessionIsRecognisedFromBothExitShapes TestDisposableSessionCloseIsAppendedWhenAnOpenSessionWillNotClose TestProveExactSelectionTimeoutCleanup TestUnrelatedExitOneKeepsItsClassification; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; if grep -rq 'found no Session to close' internal docs/user-guide .agents/skills skills; then printf 'the misleading cleanup warning remains\n' >&2; exit 1; fi; GOOS=windows go build -buildvcs=false ./internal/agent || exit 1; GOOS=windows go vet ./internal/agent` — expected: exit 0. Before this Task the new test does not exist and the warning text remains, so the command fails at its first pass check. Afterwards the absent-Session tests pass with the existing cleanup tests, the warning text is gone, and the package and its tests compile for Windows.

## References

- `_prd.md` → Sources; Goals; Core Feature 3; Success Metric 4
- `_techspec.md` → API Contract 3; Invariant 5; Current behavior; Testing Approach; Build Order 3
- ADR-0298; ADR-0114; ADR-0259; ADR-0271; ADR-0285


## Result

Implemented the Task 03 slice; Task status and declared Verification remain
Daemon-owned.

- Acceptance criterion 1: `closeDisposableSession` now returns silently for
  an absent Session. `TestMissingSessionIsRecognisedFromBothExitShapes`
  retains its diagnosis and error-chain assertions and requires zero warnings
  for both exit 4 and the recognized exit-1 shape.
  `TestTimedOutProofSetupClosesAnAbsentSessionSilently` starts with
  `t.Parallel()`, blocks the compiled fake's `sessions ensure`, waits for its
  start marker, and explicitly expires the existing test-deadline context.
  It asserts `context.DeadlineExceeded`, no `AgentSessionCleanupError`, an
  exact final close of the same disposable Session with exit 4, and zero
  warnings. No elapsed duration triggers expiration.
- Acceptance criterion 2: all other close errors retain their existing
  `AgentSessionCleanupError` wrapping and joining behavior. The existing
  open-Session proof failure and Fallback Chain probe tests passed, including
  assertions that both selection and cleanup errors remain in the chain.
  The shared helper applies the change to proofs, probes and sealed prompts;
  `EndSession` is unchanged.

Focused evidence:

- Before the production edit,
  `go test ./internal/agent -run '^(TestTimedOutProofSetupClosesAnAbsentSessionSilently|TestMissingSessionIsRecognisedFromBothExitShapes)$' -count=1`
  exited 1: the new timeout test and both exit-shape subtests exposed the
  unwanted warning.
- After the edit and formatting,
  `go test ./internal/agent -run '^(TestTimedOutProofSetupClosesAnAbsentSessionSilently|TestMissingSessionIsRecognisedFromBothExitShapes|TestDisposableSessionCloseIsAppendedWhenAnOpenSessionWillNotClose|TestACPXProbeCleanupFailureJoinsSelectionError|TestProveExactSelectionCancelCleanup)$' -count=1`
  exited 0 (`ok roundfix/internal/agent`).
- Searched `internal`, `docs/user-guide`, `docs/agents`, `.agents/skills` and
  `skills` for the obsolete warning text: no matches. Historical Spec and
  ADR quotations remain unchanged.
- `gofmt` applied to the two changed Go files; `git diff --check` exited 0.

Environment: the configured writable root contains a symlink, so the sandbox
could not launch commands. Approved execution used the resolved checkout.
The only pre-existing worktree change was the Daemon's `status: in_progress`
line in this Task file, preserved without modification.

Not run: declared Task Verification (including Windows build/vet), repository
Verification and incremental Verification; these are reserved for the Daemon.
No follow-up implementation was required.
