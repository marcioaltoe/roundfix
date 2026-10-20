---
task: task_03
spec: 0285-a-transient-profile-proof-timeout-does-not-park-the-delivery-queue
status: pending
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
