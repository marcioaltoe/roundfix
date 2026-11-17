---
task: task_03
spec: 0400-planning-references-and-release-preparation-need-no-manual-classification
status: pending
type: backend
complexity: medium
---

# Task 03: The release plan command reports a Release-Preparation Commit and reaches the same decision after it

## Overview

`roundfix release plan --format json` marks each Release-Preparation Commit
with `releasePreparation: true`. Fixture repositories committed in
temporary directories prove, through the public command, that a planning
commit with a `docs/references/` file and a release-preparation commit
leave the decision, the proposed version and the Release Ceiling's answer
equal to the range without them, and that every other shape keeps today's
answer. One named test reproduces Surface Transcript 1 on a shared clone of
this repository's own history.

## Requirements

1. MUST add `ReleasePreparation bool` with the JSON tag
   `releasePreparation,omitempty` to `releasePlanChangeJSON`, right after
   `CrossesMaintenanceOnlyBoundary`, and fill it from
   `ChangeEvidence.ReleasePreparation` (`_techspec.md` → API Contract 2).
   The text output MUST NOT change shape (API Contract 3).
2. MUST add `internal/cli/releaseplan_release_preparation_test.go` with
   `TestReleasePlanNeedsNoClassificationForPlanningReferencesOrReleasePreparation`,
   calling `t.Parallel()` first, one subtest per case, each building its
   repository with `newReleasePlanCommandRepo`, `writeReleasePlanFile`,
   `gitReleasePlan` and `commitReleasePlan`, with a delivery shaped as in
   `TestReleasePlanAnswersApprovalFromTheReleaseCeiling` (the record is
   deleted and an Archive Record added in the same commit) whose record
   grants `release` with `release_ceiling: minor` and
   `release_paths: [VERSION, CHANGELOG.md]`, and each running the public
   command in text and JSON:
   - `planning_reference_is_answered`: a planning commit that adds the Spec
     folder and `docs/references/evidence.md`, then the `feat:` delivery:
     decision `ready`, `Release authorization: answered`, no
     `Blocking commits:` line;
   - `without_the_reference_the_same_decision`: the same range without the
     reference file reaches the same state, proposed version and
     authorization answer;
   - `preparation_keeps_the_decision`: the first range plus a
     `chore: release v0.2.0` commit changing only `VERSION` and
     `CHANGELOG.md` reaches the same state, proposed version and
     authorization answer, and its JSON change has
     `releasePreparation: true` and `crossesMaintenanceOnlyBoundary: false`
     while every other change omits `releasePreparation`;
   - `preparation_with_an_unlisted_path`: the preparation also changes
     `README.md`: `manual_classification_required`, naming it;
   - `record_without_a_ceiling`: the record has no ceiling and no paths: the
     preparation is `manual_classification_required`, as today;
   - `fix_on_release_paths`: a `fix:` commit changing only `CHANGELOG.md`
     keeps `[impact=patch, type=fix]` and the plan its impact;
   - `classified_by_hand`: `--impact minor --reason x` on the answered range
     reports `Release authorization: not answered — the plan was classified by hand`.

   Each subtest MUST record `snapshotReleasePlanRepo` before the command and
   call `assertReleasePlanRepoUnchanged` after it.
3. MUST add `TestTheV0143ReleaseRangeIsReadyWithItsPreparationCommit` in
   the same file. It MUST find this repository's top level with
   `git rev-parse --show-toplevel`, make a shared clone of it in
   `t.TempDir()`, check out `v0.143.4` detached there, run
   `release plan --from v0.143.3 --to v0.143.4` in the clone, and assert
   Surface Transcript 1 of `_techspec.md` whole: the eight lines from
   `Decision: ready` to the `Next action:` line in order, then any lines,
   then `Determining commits:` and its one delivery line as the last two
   lines, empty standard error and exit 0. It MUST NOT fetch, push, tag or
   write in this repository.
4. MUST keep every existing test of `./internal/cli` passing, including
   `TestReleasePlanAnswersApprovalFromTheReleaseCeiling` and the skill
   coverage tests, without changing them.
5. MUST change no Governed Path and no file outside the two this Task
   declares.

## Subtasks

- [ ] Add the JSON evidence field.
- [ ] Add the fixture-repository cases through the public command.
- [ ] Add the transcript test on a shared clone of this repository's history.

## Acceptance Criteria

- [ ] A planning reference and a release-preparation commit leave the decision, the proposed version and the ceiling's answer of a range unchanged.
- [ ] JSON marks exactly the Release-Preparation Commit, and the text output keeps its shape.
- [ ] An unlisted path, a record without a ceiling, a `fix:` commit and a classification by hand keep today's answers.
- [ ] Surface Transcript 1 holds on this repository's own `v0.143.3..v0.143.4` range, and no fixture repository changes.

## Context

- instruction: `.agents/skills/roundfix/references/release.md`
- instruction: `internal/cli/releaseplan_release_ceiling_test.go`
- instruction: `internal/cli/releaseplan_git_source_test.go`
- interface: `internal/cli/releaseplan_command.go`
- creates: `internal/cli/releaseplan_release_preparation_test.go`

## Consumers

- `roundfix.release-plan/0.0.1` JSON change items — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestReleasePlanNeedsNoClassificationForPlanningReferencesOrReleasePreparation|TestTheV0143ReleaseRangeIsReadyWithItsPreparationCommit)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestReleasePlanNeedsNoClassificationForPlanningReferencesOrReleasePreparation/planning_reference_is_answered TestReleasePlanNeedsNoClassificationForPlanningReferencesOrReleasePreparation/without_the_reference_the_same_decision TestReleasePlanNeedsNoClassificationForPlanningReferencesOrReleasePreparation/preparation_keeps_the_decision TestReleasePlanNeedsNoClassificationForPlanningReferencesOrReleasePreparation/preparation_with_an_unlisted_path TestReleasePlanNeedsNoClassificationForPlanningReferencesOrReleasePreparation/record_without_a_ceiling TestReleasePlanNeedsNoClassificationForPlanningReferencesOrReleasePreparation/fix_on_release_paths TestReleasePlanNeedsNoClassificationForPlanningReferencesOrReleasePreparation/classified_by_hand TestTheV0143ReleaseRangeIsReadyWithItsPreparationCommit; do printf '%s\n' "$out" | grep -qF -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m ./internal/cli` — expected: exit 0. Before this Task the tests do not exist, so the command fails on the first missing pass; afterwards every case and the transcript pass by name and the whole command package passes.

## References

- `_prd.md` → Goals (second, third and fourth bullets); Core Feature 3; User Experience; Success Metric 2; Success Metric 3
- `_techspec.md` → Command surface; API Contracts 1-3; Surface Transcript 1; Testing Approach (Task 3); Build Order 3
- ADR-0412, ADR-0349, ADR-0143
