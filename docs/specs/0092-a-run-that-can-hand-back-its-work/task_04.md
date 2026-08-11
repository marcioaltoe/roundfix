---
task: task_04
spec: 0092-a-run-that-can-hand-back-its-work
status: completed
type: backend
complexity: high
---

# Task 04: Derive the Run outcome from unresolved work

## Overview

The Task that makes Task 03 safe. "Did this Batch finish?" and "is there
unresolved work left?" are different questions sharing one answer today, which is
what forces a failed Batch to overwrite the outcomes it recorded. This Task
computes the Run outcome from the Review Issues that remain unresolved, and
rewrites the six tests Task 01 enumerated to the new contract.

## Requirements

1. MUST compute a Run's outcome from the count of Review Issues still unresolved,
   not from whether a Batch reported failure.
2. MUST keep a Run reporting `Unresolved` with exit 1 while any issue remains
   unresolved, including after a Batch failure.
3. MUST let a Run reach `Clean` only when no unresolved issue remains, regardless
   of whether a Batch failed on the way.
4. MUST rewrite each of the six outcome-contract tests Task 01 enumerated,
   stating in each what the new contract asserts and why it is better rather
   than merely different.
5. MUST NOT report `Clean` for a Run whose Agent failed while unresolved work
   remains; that state is what a reverted 2026-08-09 attempt produced.

## Subtasks

- [ ] Compute the outcome from unresolved issues.
- [ ] Rewrite the six outcome-contract tests.
- [ ] Prove a crashed Agent with remaining work still ends Unresolved.

## Acceptance Criteria

- [ ] A Batch failure with issues still unresolved ends the Run `Unresolved`.
- [ ] A Batch failure with every issue resolved ends the Run `Clean`.
- [ ] All six rewritten tests pass and each states its new contract.
- [ ] No path reports `Clean` while an unresolved issue exists.

## Rehearsal Cases

- Case: a Batch of two issues where the Agent resolves one and then crashes;
  Observation: the resolved issue keeps `resolved`, the other is `failed`, and
  the Run ends `Unresolved` with exit 1.
- Case: a Batch of one issue the Agent resolves before the runtime crashes on
  teardown; Observation: the issue keeps `resolved` and the Run ends `Clean`.
- Case: a Batch that fails before resolving anything; Observation: unchanged from
  today — every issue `failed`, Run `Unresolved`.

## Bounded scope

This Task may create or modify only:

- `internal/daemon/engine.go`
- `internal/daemon/engine_test.go`
- `internal/rounds/rounds.go`
- `internal/rounds/rounds_test.go`
- `internal/cli/cli_test.go`
- `docs/specs/0092-a-run-that-can-hand-back-its-work/task_04.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestRunOutcomeDerived' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestRunOutcomeDerivedFromUnresolvedIssues'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestRunOutcomeDerived' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestRunOutcomeDerivedStaysUnresolvedAfterAFailedBatch'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/cli ./internal/daemon -count=1 2>&1 | tee /dev/stderr | grep -c '^ok' | grep -q '^2$'` — expected: exits 0, proving both packages pass rather than one being skipped. The `tee /dev/stderr` is load-bearing: without it the counting pipeline swallows every failing test name, and the Daemon records an empty diagnostic that names no failure to repair.
- `test "$(grep -c 'New contract:' internal/cli/cli_test.go internal/daemon/engine_test.go | awk -F: '{s+=$2} END {print s}')" -ge 6` — expected: exits 0, proving all six rewritten tests state their new contract.

## References

- `_prd.md` → Goals 2 and 5.
- `_techspec.md` → Build Order 4; Risks.
- ADR-0010, ADR-0113.

## Result

Implemented the derived Run-outcome contract on top of Task 03's preserved
Review Issue outcomes:

- `ResolveCycle` documents that a failed Batch preserves Terminal Review Issue
  outcomes, marks only unfinished issues failed, and exposes the recomputed
  unresolved count as the caller's sole Run-outcome input.
- Rewrote the four enumerated CLI contracts and two enumerated daemon contracts.
  Each carries a `New contract:` statement explaining both the asserted outcome
  and why unresolved-work derivation is safer than Batch-failure derivation.
  The daemon contracts retain their established discovery names and also expose
  the two `TestRunOutcomeDerived...` entry points authored by this Task.
- Extended the CLI fake Agent only enough to model one Batch where the Agent
  resolves one Review Issue, leaves the other unfinished, then crashes. The
  resulting public flow preserves `resolved`, marks the unfinished issue
  `failed`, reports `Unresolved`, and exits 1.

Pre-change signal:

- `GOCACHE="$PWD/.gocache" rtk go test ./internal/cli -run '^TestRunResolveAgentFailureMarksBatchFailed$' -count=1`
  reached the existing test and failed with `expected Run failure exit 1, got
  0`, proving Task 03 had exposed the old outcome-contract assertion.

Focused checks after the last implementation edit:

- `GOCACHE="$PWD/.gocache" rtk go test ./internal/daemon -run '^(TestRunOutcomeDerivedFromUnresolvedIssues|TestRunOutcomeDerivedStaysUnresolvedAfterAFailedBatch|TestResolveCycleVerificationFailureFailsBatchAndContinues|TestResolveCycleContinuesToNextBatchAfterFailedBatch)$' -count=1`
  passed the two daemon outcome-count contracts through both their established
  discovery names and the Task-authored entry points (`4 passed`).
- `GOCACHE="$PWD/.gocache" rtk go test ./internal/cli -run '^(TestRunResolveVerificationFailureDoesNotCommit|TestRunResolveAgentFailureMarksBatchFailed|TestRunResolveAgentFailureContinuesWithLaterBatches|TestRunResolveClosesAgentSessionForTerminalOutcomes)$' -count=1`
  passed all selected CLI functions and their subtests (`8 passed`). The CLI
  test environment required authorized access after the sandbox blocked
  `api.github.com`.
- `rtk grep -c 'New contract:' internal/cli/cli_test.go internal/daemon/engine_test.go`
  reported `4` CLI statements and `2` daemon statements.
- `rtk git diff --check` passed.

Acceptance evidence:

- A failed Batch with unresolved work stays `Unresolved`: the daemon case keeps
  the crashed first Batch's issue `failed` while the second Batch resolves, and
  returns `Remaining == 1`; the CLI mixed-status crash case reports Unresolved
  and exits 1.
- A failed Batch with every issue resolved reaches `Clean`: the daemon case
  returns `Remaining == 0` after failed Verification preserves `resolved`; the
  CLI Agent-crash and Verification-failure cases report Clean and exit 0.
- All six rewritten contracts passed in the focused package selections, and the
  source contains exactly six `New contract:` statements describing the safer
  behavior.
- No exercised path reports Clean with unresolved work: both negative cases
  retain a positive unresolved count, while every Clean case asserts a resolved
  issue set and zero unresolved work.

The authored `## Verification` commands were not run; the Daemon owns them and
terminal settlement.

### Verification Feedback repair — attempt 1

The Daemon diagnostic identified one additional CLI event test whose first
assertion still derived the Run outcome from failed Verification. Its Review
Issue was already `resolved`, so the implementation correctly computed zero
unresolved work and returned Clean under ADR-0113.

The repair keeps the test's event invariant: failed Verification must be
journaled and must not create a Batch commit. It now also requires the derived
Clean outcome and exactly one Final Push because no unresolved Review Issue
remains.

Focused checks after the repair:

- `GOCACHE="$PWD/.gocache" rtk go test ./internal/cli -run '^TestFailedVerificationJournalsFailureWithoutCommitEvents$' -count=1`
  passed the diagnostic's failing case (`1 passed`).
- `GOCACHE="$PWD/.gocache" rtk go test ./internal/cli -run '^(TestRunResolveVerificationFailureDoesNotCommit|TestRunResolveAgentFailureMarksBatchFailed|TestRunResolveAgentFailureContinuesWithLaterBatches|TestRunResolveClosesAgentSessionForTerminalOutcomes|TestFailedVerificationJournalsFailureWithoutCommitEvents)$' -count=1`
  passed the affected CLI outcome and event cases (`9 passed`).
- `GOCACHE="$PWD/.gocache" rtk go test ./internal/daemon -run '^(TestRunOutcomeDerivedFromUnresolvedIssues|TestRunOutcomeDerivedStaysUnresolvedAfterAFailedBatch|TestResolveCycleVerificationFailureFailsBatchAndContinues|TestResolveCycleContinuesToNextBatchAfterFailedBatch)$' -count=1`
  passed the daemon outcome-count contracts (`4 passed`).

The Daemon's declared Verification command was not rerun during this repair
turn; the Daemon owns the single post-repair attempt and terminal settlement.
