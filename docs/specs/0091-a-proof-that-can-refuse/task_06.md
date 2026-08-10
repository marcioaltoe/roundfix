---
task: task_06
spec: 0091-a-proof-that-can-refuse
status: pending
type: test
complexity: medium
---

# Task 06: Teach the sibling tests the catalogue read

## Overview

Corrective Task. Task 02 made proof read what a runtime advertises before
requesting a selection, which inserts one `sessions ensure` and one
`sessions show` into the acpx invocation sequence. Sixteen sibling tests
assert that sequence literally and now fail against the new behaviour. Task 02
should have carried them — a Task's declared scope includes the tests its own
change invalidates — so this Task closes that gap before the graph continues.

The tests are wrong about the sequence, not about the contract: each still
means what it meant, and each must keep meaning it after the catalogue read is
accounted for. Do not weaken an assertion to make it pass, and do not delete a
case.

## Requirements

1. MUST update every failing test so it accounts for the catalogue read that
   proof now performs before requesting a selection.
2. MUST keep each test's original subject intact: a test that proved a
   disposable session closes still proves that, and a test that proved a
   rejection diagnosis still proves that same diagnosis.
3. MUST NOT delete a test case, relax an assertion to a substring of what it
   checked before, or replace an exact invocation sequence with a length check.
4. MUST leave the tests that already pass byte-identical.

## Subtasks

- [ ] Account for the catalogue read in the asserted acpx invocation sequences.
- [ ] Update the diagnosis expectations that now surface capability evidence
      before the selection error.
- [ ] Confirm no other package asserts the old sequence.

## Acceptance Criteria

1. These sixteen tests pass, and each still asserts the behaviour its name
   claims:
   `TestACPXProbeSelectionRejectionClosesDisposableSession`,
   `TestACPXProbeSelectionSetupUsesBoundedContext`,
   `TestACPXProbeSkipsEmptyReasoningEffort`,
   `TestACPXProbeValidatesSelectionWithDisposableSession`,
   `TestAgentSelectionProfilesMacro`,
   `TestApplySessionSelectionDisposableAndLiveOrder`,
   `TestDisposableSessionCloseIsAppendedWhenAnOpenSessionWillNotClose`,
   `TestDisposableSessionCloseIsNotAppendedWhenTheSessionNeverOpened`,
   `TestProfileProofAppliesExactReasoningAndClosesDisposableSession`,
   `TestProfileProofClosesDisposableSessionOnSelectionFailure`,
   `TestProveExactSelectionCancelCleanup`,
   `TestProveExactSelectionCleanupJoinedFailure`,
   `TestProveExactSelectionEffectiveMismatchCleanup`,
   `TestProveExactSelectionOfficialFixturesNoPrompt`,
   `TestRunImplementDetachPrintsReportAndCompletesRun`,
   `TestRunImplementDetachSurvivesCallerProcessGroupKill`.
2. No test file outside `internal/agent/acpx_runner_test.go` and
   `internal/cli/implement_test.go` changes. The Daemon audits the Task commit
   against this boundary; a Verification command cannot prove it, because a
   check that nothing else changed passes most easily when no work happened at
   all.

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestACPXProbeSkipsEmptyReasoningEffort$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestACPXProbeSkipsEmptyReasoningEffort'` — expected: exits 0. Fails today against the unchanged tree, because the asserted sequence omits the catalogue read.
- `GOCACHE="$PWD/.gocache" go test ./internal/agent -run '^TestDisposableSessionCloseIsAppendedWhenAnOpenSessionWillNotClose$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestDisposableSessionCloseIsAppendedWhenAnOpenSessionWillNotClose'` — expected: exits 0. Fails today, because the diagnosis now surfaces capability evidence first.
- `GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^TestRunImplementDetachSurvivesCallerProcessGroupKill$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestRunImplementDetachSurvivesCallerProcessGroupKill'` — expected: exits 0. Fails today, and proves the CLI half was carried too.

## References

- PRD: `_prd.md` — the proof that can refuse.
- TechSpec: `_techspec.md` — the catalogue read introduced by Task 02.
