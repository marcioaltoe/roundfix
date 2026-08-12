---
task: task_04
spec: 0083-a-gate-that-can-say-no
status: completed
type: test
complexity: medium
---

# Task 04: Measure the corpus sweep in a unit the machine cannot move

## Overview

The corpus budget asserts the sweep finishes under one second of wall clock. Its
author already guarded against in-suite parallelism with a dedicated-run check,
but wall clock cannot be guarded against a busy machine — and on 2026-08-07 a
loaded laptop produced 2.5 seconds and a red build with nothing wrong. This task
keeps the performance signal and expresses it in a unit competing load cannot
change.

## Requirements

1. MUST stop failing on elapsed time. Wall clock may still be logged for a
   human, but it MUST NOT gate.
2. MUST keep a performance signal that a real regression would trip — work
   performed by the sweep, such as operations or file reads, rather than time
   taken.
3. MUST choose a unit that stays meaningful as the corpus grows, so adding
   Specs does not silently consume the budget the way a fixed count would.
4. MUST honor the batched-read constraint the repository already accepted: the
   sweep's cost is bounded by operations, and re-expressing the assertion MUST
   NOT introduce per-file work the repository removed on purpose.
5. MUST keep the existing dedicated-run guard or state why it is no longer
   needed once the unit is load-independent.
6. MUST change only these repository-relative paths plus this Task file:
   `internal/speccheck/constraints_characterization_test.go`. Any other changed
   path fails this Task.

## Subtasks

- [x] Replace the wall-clock assertion with a load-independent unit.
- [x] Keep wall clock as a logged observation.
- [x] Confirm the unit trips on an induced regression.
- [x] Settle the dedicated-run guard's fate and record the reason.
- [x] Confirm the changed-file set matches the declared boundary.

## Acceptance Criteria

- [x] The check passes on a deliberately loaded machine, proven by running it
      under induced load rather than asserted.
- [x] An induced inefficiency in the sweep still fails the check, proven by
      observation.
- [x] Wall clock still appears in the check's output.
- [x] No assertion in the check compares a duration against a constant.

## Context

- instruction: `docs/workflow/authorizations/2026-08-07-make-the-gate-honest.md`
- interface: `internal/speccheck/constraints_characterization_test.go`

## Verification

- `go test ./internal/speccheck -run '^TestCheckCorpusBudget$' -count=1 -parallel=1 -v > /tmp/task_04-1.log 2>&1 && grep -q '^--- PASS: TestCheckCorpusBudget' /tmp/task_04-1.log` — expected: exits 0.
- `grep -n -E '(time\.(Since|Until)|\.(After|Before)\()' internal/speccheck/constraints_characterization_test.go | grep -E '(Errorf|Fatalf|Fatal|assert|require)' ; test $? -eq 1` — expected: exits 0, proving no time-comparison function feeds a failing assertion.
- `go test ./internal/speccheck -count=1` — expected: exits 0.
- `(git diff --name-only HEAD; git ls-files --others --exclude-standard) | grep -v -E '^(internal/speccheck/constraints_characterization_test\.go|docs/specs/0083-a-gate-that-can-say-no/task_04\.md)$' | grep . ; test $? -eq 1` — expected: exits 0, proving no path outside the declared boundary changed.

## References

- `_techspec.md` → Build Order 6; Risks: any wall-clock threshold on a shared machine is a future false alarm.
- `_prd.md` → Core Feature 4; Goal 2.
- ADR-0090.

## Result

### Implementation

- Replaced the one-second wall-clock assertion with a budget of at most one
  `speccheck.Check` operation per eligible Spec. The allowance is derived from
  the measured Spec count, so corpus growth expands the budget proportionally.
- Counted the existing expensive operation at its sweep call site without
  adding a directory walk or per-file read. The corpus-golden check uses the
  same sweep with accounting disabled.
- Kept elapsed wall clock in the verbose log beside the operation count and
  normalized allowance.
- Removed the dedicated-run guard because package concurrency and host load
  cannot alter the operation count; the test records that reason directly.

### Focused checks

- Before the change, eight concurrent normal-priority `yes` processes plus
  `rtk proxy env GOCACHE="$PWD/.gocache" go test ./internal/speccheck -run Budget -count=3 -parallel=1 -v`
  failed all three runs at `1.123577625s`, `1.188922791s`, and `1.0320415s`
  against the unchanged one-second wall-clock budget.
- After the change, the same loaded three-run probe exited zero. Each run
  reported `82 Check operations across 82 Specs`; one logged `1.042782125s`,
  proving elapsed time no longer controlled its verdict.
- After formatting and removing the temporary mutation, a fresh eight-burner
  run exited zero while logging `1.13787625s` and `82 Check operations across
  82 Specs`.
- With a temporary second real `checkCorpusSpec` call for every eligible Spec,
  the focused check exited non-zero at `164 Check operations across 82 Specs,
  want at most 82`. Removing the mutation restored the focused check.
- `rtk proxy env GOCACHE="$PWD/.gocache" go test ./internal/speccheck -run 'TestCheckCorpus(Golden|Budget)$' -count=1 -parallel=1 -v`
  exited zero: both corpus checks passed, and the budget check logged
  `642.721833ms` plus `82 Check operations across 82 Specs`.
- `rtk rg -n "time\\.Since|elapsed\\s*(>=|>)|time\\.Second|corpusBudget" internal/speccheck/constraints_characterization_test.go`
  found only the `time.Since` observation. No duration budget or comparison
  remains.
- `rtk git diff --name-only HEAD` and
  `rtk git -c core.fsmonitor=false status --short --untracked-files=all` listed
  only the characterization test and this Task file.

### Acceptance evidence

1. Eight CPU burners left the final operation-budget run passing while its
   logged wall clock exceeded the retired one-second threshold.
2. Doubling the expensive sweep operation made the check fail at two operations
   per Spec, so the normalized allowance detects the induced inefficiency.
3. Every focused budget run logged `full Spec corpus sweep completed in
   <duration>` before its work accounting.
4. The only remaining duration expression feeds `t.Logf`; the failing assertion
   compares counted `Check` operations with the Spec-normalized allowance.

### Daemon verification

Not run in this Agent turn. The Daemon owns every command in `## Verification`
and the Task's terminal status.
