---
task: task_04
spec: 0089-an-effort-the-runtime-actually-receives
status: completed
type: backend
complexity: high
---

# Task 04: Warm the session, apply the effort, and publish the receipt

## Overview

The `effort` option exists only once a queue-owner process holds the selected
model, and the adapter raises that owner on the session's first prompt. This
Task sends one minimal warm-up prompt before any work turn, applies the
requested effort, observes the effective value, and publishes that observation
as a Run Event — because a transcript records what an Agent said and a receipt
records what the system executed.

## Requirements

1. MUST, for a `runtime_deferred` assignment, warm the Agent Session with one
   minimal prompt before applying the effort.
2. MUST apply the requested effort after the warm-up and observe the effective
   value, failing the session when the observed value differs.
3. MUST publish the observed effective effort as a Run Event so the evidence
   outlives the call that checked it. Run Event kinds are a closed enum in
   `internal/runevent`; if the receipt needs a kind that does not exist, add it
   there rather than overloading an existing one, and keep the existing
   `agent-selection` stream category.
4. MUST be idempotent per Agent Session: a session already warmed is not warmed
   again.
5. MUST order the acpx commands ensure, warm-up prompt, effort set, then work
   prompt — the ordering is the fix and MUST be asserted, not assumed.
6. MUST use a constant warm-up prompt that reads as setup rather than as an
   instruction the Agent could act on, and MUST record its exact bytes in this
   Task's Result so a reviewer can judge them.
7. MUST NOT warm a session whose assignment is not `runtime_deferred`.
8. MUST re-record the coverage record in this Task's own commit if any test is
   renamed or removed.

## Subtasks

- [ ] Add the warm-up step to the session lifecycle behind the encoding check.
- [ ] Apply and observe the effort, failing on a mismatch.
- [ ] Publish the observed value as a Run Event.
- [ ] Make the step idempotent per session.
- [ ] Assert the full command ordering.
- [ ] Edit the break-half characterization test and declare the break.

## Acceptance Criteria

- [ ] A `runtime_deferred` session records the acpx sequence ensure, prompt,
      effort set, prompt — asserted in order.
- [ ] The observed effective effort appears as a Run Event with the requested
      value.
- [ ] A second work turn on the same session issues no second warm-up.
- [ ] A session whose observed effort differs from the requested one fails
      before any work turn.
- [ ] A Codex session's command sequence is unchanged and contains no warm-up.
- [ ] An `opencode` session with an empty effort contains no warm-up and no
      effort set.

## Context

- interface: `internal/agent/acpx_runner.go`
- interface: `internal/runevent`

## Bounded scope

This Task may create or modify only:

- `internal/agent/acpx_runner.go`
- `internal/agent/acpx_runner_test.go`
- `internal/agent/acpx_session_effort_characterization_test.go`
- `internal/agent/selection_assignment.go`
- `internal/runevent/event.go`
- `internal/runevent/selection.go`
- `internal/runevent/stream.go`
- `internal/runevent/stream_test.go`
- `docs/references/coverage-record.json`
- `docs/specs/0089-an-effort-the-runtime-actually-receives/task_04.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go test ./internal/agent -count=1` — expected: exits 0.
- `go test ./internal/agent -run 'WarmSession' -count=1 -v` — expected: exits 0 and names at least one test; `no tests to run` fails this Task.
- `go test ./internal/agent -run 'WarmSessionIsIdempotent' -count=1 -v` — expected: exits 0, proving a second turn adds no second warm-up.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -count=1` — expected: exits 0.

## References

- `_prd.md` → Goals 1 and 3; Core Features: session warm-up.
- `_techspec.md` → Implementation Design: Interfaces; Consultation; Build Order 4.
- ADR-0108.

## Result

The Agent Session lifecycle now recognizes `runtime_deferred` assignments,
sends one setup prompt before applying the advertised effort, validates the
effective value returned by acpx, and stops before work when that evidence is
contradictory. The per-session warmed state is retained independently of the
fully ensured state, so a retry after a later setup failure does not send a
second setup prompt; ending the Agent Session clears both states.

The observed value is published as the new `agent.selection_receipt` Run Event.
Its payload carries the Agent Session, runtime, model, requested effort,
observed effort, and `applied` status. The public Run Event Stream projects the
receipt through the existing `agent-selection` category with an
`agent_session` scope.

Warm-up prompt exact bytes (14 ASCII/UTF-8 bytes, with no trailing newline):
`Session setup.` Hex: `53 65 73 73 69 6f 6e 20 73 65 74 75 70 2e`.

Pre-change signal:

- `GOCACHE="$PWD/.gocache" rtk go test ./internal/agent -run '^TestACPXSessionEffortCharacterizationDeclaredBreakWarmSessionOrdersEffortBeforeWork$' -count=1` — failed to compile because `acpxDeferredEffortWarmupPrompt` did not exist. The replaced characterization test had asserted the prior behavior: ensure followed by the work prompt with no effort set in between.

Focused checks run after the last Go edit:

- `GOCACHE="$PWD/.gocache" rtk go test ./internal/agent -run '^(TestACPXSessionEffortCharacterizationDeclaredBreakWarmSessionOrdersEffortBeforeWork|TestACPXRunWarmSessionPublishesEffectiveEffortReceipt|TestACPXRunWarmSessionIsIdempotent|TestACPXRunWarmSessionMismatchStopsBeforeWork|TestACPXRunAppliesSelectionBeforePrompt)$' -count=1` — 8 tests/subtests passed.
- `GOCACHE="$PWD/.gocache" rtk go test ./internal/runevent -count=1` — 46 tests passed, including the new receipt projection through `agent-selection`.
- `GOCACHE="$PWD/.gocache" rtk go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record` — exited 0 and regenerated `docs/references/coverage-record.json` after the characterization test rename; the record was generated rather than hand-edited.
- `rtk git diff --check` — exited 0.

Acceptance evidence:

- Criterion 1: `TestACPXSessionEffortCharacterizationDeclaredBreakWarmSessionOrdersEffortBeforeWork` passed with the ordered command keys `sessions ensure`, `prompt`, `set effort`, `prompt`, and asserted the setup and work prompt bytes in that order.
- Criterion 2: `TestACPXRunWarmSessionPublishesEffectiveEffortReceipt` passed with exactly one `agent.selection_receipt` whose requested and observed values were both `xhigh`; `TestProjectStreamEventProjectsSelectionReceiptInExistingCategory` passed with category `agent-selection`.
- Criterion 3: `TestACPXRunWarmSessionIsIdempotent` passed across two work turns on the same Agent Session with three total prompts: one setup prompt and two work prompts.
- Criterion 4: `TestACPXRunWarmSessionMismatchStopsBeforeWork` passed with acpx returning `high` as the set value but `low` as the effective current value. The capability boundary rejected `contradictory_response`, and the prompt ledger contained only the setup prompt.
- Criterion 5: the `codex reasoning_effort` subtest of `TestACPXRunAppliesSelectionBeforePrompt` passed with its unchanged ensure, `set reasoning_effort`, work-prompt sequence and no setup prompt.
- Criterion 6: the `opencode model-managed reasoning issues no effort set` subtest passed with only ensure and the work prompt; it contained neither a setup prompt nor an effort set.

No follow-up work was found inside this Task's slice. The commands authored
under `## Verification` were not run; Daemon Verification remains the
settlement boundary.
