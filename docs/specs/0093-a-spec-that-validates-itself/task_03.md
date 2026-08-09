---
task: task_03
spec: 0093-a-spec-that-validates-itself
status: completed
type: backend
complexity: medium
---

# Task 03: Let a caller ask for one authoring stage

## Overview

The checker runs every detector over every artifact. An author finishing a PRD
wants the rules decidable from a PRD, and wants them in under a second. This
Task adds the scope, assigns each detector the earliest stage at which every
artifact it reads exists, and leaves the default sweep byte-identical so
`make verify` is unchanged.

## Requirements

1. MUST add a stage scope with `prd`, `techspec`, and `tasks` values, plus a
   default that runs every detector as today.
2. MUST assign every existing detector the earliest stage at which every
   artifact it reads exists.
3. MUST leave the default sweep's findings identical to today's for every Spec
   in the corpus; a rule may be checked earlier, never dropped.
4. MUST report, for a scoped run, which detectors it skipped, so a narrow run is
   never mistaken for a clean full sweep.
5. MUST NOT change any detector's rule in this Task.

## Subtasks

- [ ] Add the stage type and its default.
- [ ] Assign each detector its earliest stage.
- [ ] Report skipped detectors on a scoped run.

## Acceptance Criteria

- [ ] A `prd`-scoped run executes the detectors that read only a PRD.
- [ ] A `tasks`-scoped run executes every detector, because the whole Spec
      exists by then.
- [ ] The default sweep produces the same findings as before this Task across
      the corpus.
- [ ] A scoped run names the detectors it did not run.

## Bounded scope

This Task may create or modify only:

- `internal/speccheck/coherence.go`
- `internal/speccheck/coherence_test.go`
- `docs/specs/0093-a-spec-that-validates-itself/task_03.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/speccheck -run '^TestStageScope' -count=1 -v 2>&1 | grep -q '^--- PASS: TestStageScopeRunsOnlyDetectorsTheStageCanDecide'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/speccheck -run '^TestStageScope' -count=1 -v 2>&1 | grep -q '^--- PASS: TestStageScopeDefaultSweepIsUnchanged'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/speccheck -run '^TestStageScope' -count=1 -v 2>&1 | grep -q '^--- PASS: TestStageScopeNamesTheDetectorsItSkipped'` — expected: exits 0.
- `grep -q 'StagePRD' internal/speccheck/coherence.go` — expected: exits 0. This string does not exist before this Task.

## References

- `_prd.md` → Goal 2.
- `_techspec.md` → Build Order 3; System Architecture.
- ADR-0117.

## Result

Implemented cumulative authoring-stage scope in `internal/speccheck`. `StageAll`
keeps the existing `Check` sweep, `StageTasks` delegates to that same full
sweep, and PRD and TechSpec scopes run only the detector groups whose inputs
exist by that stage. Scoped results add one `SkippedDetector` entry per later
detector, named by its existing finding code.

Focused checks:

- `rtk gofmt -w internal/speccheck/coherence.go internal/speccheck/coherence_test.go`
  — exited `0`.
- `rtk proxy env GOCACHE="$PWD/.gocache" go test ./internal/speccheck -run '^TestStageScope' -count=1`
  — exited `0`; the focused stage suite reported `ok` in 2.334 seconds.
- `rtk proxy env GOCACHE="$PWD/.gocache" go test ./internal/speccheck -count=1`
  — exited `0`; the package reported `ok` in 4.817 seconds.

Acceptance evidence:

- PRD scope executes PRD-only detectors:
  `TestStageScopeRunsOnlyDetectorsTheStageCanDecide` observes
  `SC-CONSTRAINT-MISSING` findings from an incomplete PRD and no findings from
  TechSpec- or Task-stage detectors.
- Tasks scope executes every detector: the same test proves the complete
  `StageTasks` result is deeply equal to the existing unscoped `Check` result.
- The default sweep is unchanged: `TestStageScopeDefaultSweepIsUnchanged`
  compares `StageAll` with `Check` for every active and archived Spec in the
  repository corpus and requires deeply equal findings and skipped-detector
  records.
- Scoped runs name omitted detectors:
  `TestStageScopeNamesTheDetectorsItSkipped` requires all nine TechSpec- and
  Task-stage detector codes in the PRD result's `Skipped` records.

The Daemon-owned `## Verification` commands were not run in this Agent turn.
