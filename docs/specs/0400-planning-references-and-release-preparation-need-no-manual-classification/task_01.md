---
task: task_01
spec: 0400-planning-references-and-release-preparation-need-no-manual-classification
status: pending
type: backend
complexity: medium
---

# Task 01: The Release Plan counts planning references and a Release-Preparation Commit as maintenance-only

## Overview

The release plan domain stops asking for a classification by hand for two
commits that ship nothing. A path under `docs/references/` becomes planning
evidence, and a commit with no automatic impact whose every changed path is
a `release_paths` entry of a granted record the range delivers becomes a
Release-Preparation Commit inside the maintenance boundary. `Build` reads
the grants once, before classifying, when a commit may be one. Table tests
in `internal/releaseplan`, with no Git, prove each rule and its limits. No
command output changes yet; task 03 adds the JSON evidence.

## Requirements

1. MUST add the prefix `docs/references/` to the planning-evidence paths of
   `MaintenanceOnly`, and keep every earlier answer of `MaintenanceOnly` for
   every other path (`_techspec.md` → Interfaces, invariant 1).
2. MUST add `ReleasePaths []string` to `ClassifyRequest` and
   `ReleasePreparation bool` to `ChangeEvidence`, with the names of
   `_techspec.md` → Interfaces.
3. MUST mark a commit as a Release-Preparation Commit in `ClassifyChanges`
   exactly when invariant 2 holds: its automatic impact is `none`, it
   changes at least one path, and each changed path, cleaned as
   `MaintenanceOnly` cleans it, equals an entry of `ReleasePaths`. Such a
   commit MUST have `ReleasePreparation` true and
   `CrossesMaintenanceOnlyBoundary` false; no other commit has
   `ReleasePreparation` true (invariant 3). A commit with a `feat`, `fix`,
   `perf` or breaking impact MUST keep its impact and boundary even when it
   changes only release paths.
4. MUST leave `ClassifyChanges` returning the same result as before for
   every request with empty `ReleasePaths` and no path under
   `docs/references/` (invariant 4).
5. MUST make `Build` call `ReleaseGrants` at most once per plan: before
   `ClassifyChanges` when the source implements `ReleaseGrantSource` and a
   commit has automatic impact `none` and crosses the boundary, passing the
   `Paths` of every grant with `Granted` true as `ReleasePaths`; otherwise
   after the proposal, only when it is `approval_required`, as today.
   `AnswerApproval` MUST receive the grants already loaded and stay
   unchanged (invariants 5 and 6).
6. MUST extend `TestMaintenanceOnly` with the case
   `docs/references/evidence.md` true, keeping its other cases.
7. MUST add `internal/releaseplan/release_preparation_test.go` with
   `TestReleasePreparationCommitSetsNoImpact`, calling `t.Parallel()`
   first, one subtest per case, each asserting the whole `ChangeEvidence`
   of every commit, the result state, the classification source and
   `BlockingCommits`, never a substring:
   - `inside_release_paths`: a `chore: release v0.2.0` commit changing
     `CHANGELOG.md` and `VERSION`, with both in `ReleasePaths`, after a
     `fix:` commit: state `ready`, no blocking commit;
   - `cleaned_paths_match`: the same commit with the path written
     `./CHANGELOG.md`;
   - `one_unlisted_path`: the same commit also changing `README.md`: it
     blocks, as today;
   - `no_release_paths`: empty `ReleasePaths`: it blocks, as today;
   - `no_changed_path`: a commit with no path: it blocks, as today;
   - `fix_on_release_paths`: a `fix:` commit changing only `CHANGELOG.md`
     keeps impact `patch` and crosses the boundary;
   - `breaking_on_release_paths`: a `chore!:` commit changing only
     `CHANGELOG.md` keeps impact `major`;
   - `docs_references_planning_commit`: a `docs:` commit changing
     `docs/specs/0001-x/_prd.md` and `docs/references/evidence.md` is
     maintenance-only with no `ReleasePaths`.
8. MUST rename `TestBuildLooksForGrantsOnlyWhenApprovalIsRequired` to
   `TestBuildLoadsReleaseGrantsOnceWhenACommitMayPrepareTheRelease` and
   assert the call count of each subtest: `patch` 0, `maintenance` 0,
   `manual_classification` (a `refactor:` commit) 1, `minor` 1,
   `candidate_in_a_minor_range` (a `feat:` delivery and a `chore:` commit) 1,
   and add `preparation_is_answered`: a `feat:` delivery and a `chore:`
   commit changing only a path the fake's granted grant lists end `ready`
   with the authorization answered. The fake grant gains `Paths`. Keep the
   `source without the optional capability` and `read error` subtests.
   This is a declared change of that test's expectation (`_techspec.md` →
   Testing Approach, Task 1).
9. MUST keep every test of `TestAnswerApprovalFollowsEveryCeiling` passing
   unchanged.
10. MUST change no Governed Path and no file outside the ones this Task
    declares.

## Subtasks

- [ ] Add `docs/references/` to the planning-evidence paths.
- [ ] Add the request and evidence fields and the Release-Preparation Commit rule.
- [ ] Load the grants once in `Build` before classifying when a candidate exists.
- [ ] Add the table test and extend the boundary and Build tests.

## Acceptance Criteria

- [ ] A commit changing only planning paths and `docs/references/` files is maintenance-only, and every earlier boundary answer holds.
- [ ] A commit with no impact that changes only granted release paths is a Release-Preparation Commit; one unlisted path, no paths, an empty commit or a commit with an impact keeps today's classification.
- [ ] `Build` reads the grants once, and only when a candidate exists or approval is required.
- [ ] A Release Ceiling still never answers a plan classified by hand.

## Context

- instruction: `docs/adr/0412-planning-references-and-the-release-preparation-commit-set-no-release-impact.md`
- instruction: `docs/adr/0349-a-spec-release-ceiling-answers-the-release-plan-approval.md`
- interface: `internal/releaseplan/classifier.go`
- interface: `internal/releaseplan/model.go`
- interface: `internal/releaseplan/build.go`
- interface: `internal/releaseplan/classifier_test.go`
- interface: `internal/releaseplan/release_grant_test.go`
- creates: `internal/releaseplan/release_preparation_test.go`

## Consumers

- `releaseplan.ClassifyRequest`, `releaseplan.ChangeEvidence` and `releaseplan.Build` — consumers: `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestReleasePreparationCommitSetsNoImpact|TestBuildLoadsReleaseGrantsOnceWhenACommitMayPrepareTheRelease|TestMaintenanceOnly)$' ./internal/releaseplan 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestReleasePreparationCommitSetsNoImpact/inside_release_paths TestReleasePreparationCommitSetsNoImpact/cleaned_paths_match TestReleasePreparationCommitSetsNoImpact/one_unlisted_path TestReleasePreparationCommitSetsNoImpact/no_release_paths TestReleasePreparationCommitSetsNoImpact/no_changed_path TestReleasePreparationCommitSetsNoImpact/fix_on_release_paths TestReleasePreparationCommitSetsNoImpact/breaking_on_release_paths TestReleasePreparationCommitSetsNoImpact/docs_references_planning_commit TestBuildLoadsReleaseGrantsOnceWhenACommitMayPrepareTheRelease/manual_classification TestBuildLoadsReleaseGrantsOnceWhenACommitMayPrepareTheRelease/candidate_in_a_minor_range TestBuildLoadsReleaseGrantsOnceWhenACommitMayPrepareTheRelease/preparation_is_answered TestMaintenanceOnly/docs/references/evidence.md; do printf '%s\n' "$out" | grep -qF -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m ./internal/releaseplan ./internal/cli` — expected: exit 0. Before this Task the new tests and subtests do not exist, so the command fails on the first missing pass; afterwards each case passes by name, and the whole release plan domain and the command package that consumes it pass.

## References

- `_prd.md` → Goals (first, second and fourth bullets); Core Features 1 and 2; Success Metric 1; Non-Goals
- `_techspec.md` → Classifier; Build; Interfaces (invariants 1 to 6); API Contract 1; Testing Approach (Task 1); Build Order 1
- ADR-0412, ADR-0349, ADR-0143
