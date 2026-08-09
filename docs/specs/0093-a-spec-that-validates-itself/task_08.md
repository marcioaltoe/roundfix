---
task: task_08
spec: 0093-a-spec-that-validates-itself
status: completed
type: backend
complexity: medium
---

# Task 08: Run the citation detector in the authoring stages

## Overview

Corrective Task for F-001 of `qa/qa-report-2026-08-09-01.md`. The stage registry
never registered `SC-CITATION-UNSUPPORTED`, so `--stage prd` and `--stage
techspec` return `No findings` for an artifact the unscoped run correctly
rejects. The authoring skills instruct authors to run exactly those scopes, so
the shipped wiring under-checks while reporting clean — worse than not checking,
because the author believes the citation was read.

## Requirements

1. MUST register the citation detector in the stage registry so `--stage prd`
   and `--stage techspec` run it.
2. MUST assign it the earliest stage at which every artifact it reads exists,
   which is `prd`: the claim and the cited record are both present then.
3. MUST leave the unscoped run's findings identical, since it already detects
   the claim; this Task widens where the detector runs and changes no verdict.
4. MUST prove the gap is closed with the exact reproduction the finding names:
   the Spec 0090 fixture rejected under `--stage prd` and under `--stage
   techspec`, not only unscoped.
5. MUST NOT add or relax any detector rule.

## Subtasks

- [ ] Register the detector for the `prd` stage.
- [ ] Prove both scoped commands now reject the fixture.
- [ ] Prove the unscoped verdict is unchanged.

## Acceptance Criteria

- [ ] `--stage prd` reports `SC-CITATION-UNSUPPORTED` on the fixture and exits
      non-zero.
- [ ] `--stage techspec` does the same.
- [ ] The unscoped run's findings are unchanged across the corpus.
- [ ] No detector rule changed.

## Rehearsal Cases

- Case: the `internal/speccheck/testdata/citation/repo` fixture carrying Spec
  0090's original attribution, checked with `--stage prd`; Observation:
  `SC-CITATION-UNSUPPORTED` with both texts, non-zero exit.
- Case: the same fixture with `--stage techspec`; Observation: the same finding.
- Case: the same fixture unscoped; Observation: identical to before this Task.

## Bounded scope

This Task may create or modify only:

- `internal/speccheck/coherence.go`
- `internal/speccheck/coherence_test.go`
- `docs/specs/0093-a-spec-that-validates-itself/task_08.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/speccheck -run '^TestStageScope' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestStageScopeRunsTheCitationDetectorInAuthoringStages'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/speccheck -run '^TestStageScope' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestStageScopeDefaultSweepIsUnchanged'` — expected: exits 0, proving the unscoped verdict did not move.
- `grep -q 'CodeCitationUnsupported' internal/speccheck/coherence.go` — expected: exits 0. This string does not exist in the file before this Task.

## References

- `_prd.md` → Goals 1 and 2.
- `_techspec.md` → Build Order 3; System Architecture.
- `qa/qa-report-2026-08-09-01.md` → F-001.
- ADR-0117.

## Result

The stage registry now assigns the existing citation detector to `prd`, and
the authoring-stage path passes the PRD plus the TechSpec when present to that
detector. The detector's parsing, resolution, and verdict rules remain
unchanged.

Focused implementation evidence, separate from Daemon Verification:

- Acceptance criterion 1: `rtk proxy env GOCACHE="$PWD/.gocache" go test
  ./internal/speccheck -run
  '^TestStageScopeRunsTheCitationDetectorInAuthoringStages$' -count=1` exits
  `0`. Its `PRD` subtest checks the Spec 0090 fixture for one error-level
  `SC-CITATION-UNSUPPORTED` finding containing both the false attribution and
  ADR-0083's actual subject.
- Acceptance criterion 2: the same focused test's `TechSpec` subtest exercises
  `StageTechSpec` against the same fixture and requires the same finding,
  severity, and two texts.
- Acceptance criterion 3: `rtk proxy env GOCACHE="$PWD/.gocache" go test
  ./internal/speccheck -run '^TestStageScopeDefaultSweepIsUnchanged$' -count=1`
  exits `0` across the active and archived characterization corpus. The
  unscoped `Check` implementation was not edited.
- Acceptance criterion 4: the change only registers and invokes the existing
  detector from `coherence.go`; `internal/speccheck/citations.go` is unchanged.
  `rtk proxy env GOCACHE="$PWD/.gocache" go test ./internal/speccheck -run
  '^TestCitation' -count=1` exits `0` for the existing citation behavior suite.

The regression test first failed for both `PRD` and `TechSpec` with zero
citation findings, then passed after the registry and scoped invocation change.
The Task's authored `## Verification` commands were not run; the Daemon owns
those commands and Task settlement.
