---
task: task_04
spec: 0266-the-review-record-survives-a-fix
status: pending
type: backend
complexity: medium
---

# Task 04: A code correction inside the archived Spec's declared paths returns to review

## Overview

ADR-0233's third retry rule admits a changed path only under an archived Spec
the blocker names, so a disposed code fix is refused even when the archived
Spec's own Task declared the file (`_techspec.md` → Current behavior, second
characterization). On 2026-10-13 the maintainer answered "Relaxar, com
limites": the retry also admits a path the archived Spec's Tasks declared.
This Task adds `spec.DeclaredTaskPaths` and makes the correction proof collect
each named Spec's declared paths at the recorded revision. It is verifiable on
its own: a disposed fix to a declared file is accepted, for a Spec archived as
a record and as a folder, and any other file is refused by name.

## Requirements

1. MUST implement `_techspec.md` → API Contract 6, refining ADR-0233 as
   ADR-0272 records: add `DeclaredTaskPaths(content []byte, taskFile string) ([]string, error)`
   in `internal/spec/declared_paths.go`, returning the clean, deduplicated
   `creates:` and `interface:` paths of one Task file's `## Context` and the
   paths of its `## Recorded paths` (`RecordedTaskPaths`), through the
   package's existing Task parser; and add the unexported
   `archivedSpecDeclaredPaths` in `internal/cli/deliver_review_correction.go`.
2. MUST read the declared paths from Git, never from the working tree: for a
   Legacy Archive Folder, its `task_*.md` files at the parked candidate; for an
   Archive Record, the `task_*.md` files under its `source` at its
   `source_revision`, recovered with `spec.EnsureSourceRevision` when the
   commit is missing locally (ADR-0261). Unreadable or unparsable Task files
   refuse with the reason
   `declared paths of archived Spec <slug> are unavailable: <cause>`.
3. MUST admit a changed path that lies under the Spec's archive path, as
   today, or equals one of its declared paths exactly, and refuse the first
   other path in the order Git lists it with the reason
   `changed path <path> is outside the named archived Specs and the paths their Tasks declared`.
   The descent, disposition and record checks and their reasons MUST stay as
   task_02 and task_03 leave them, and the accepted item MUST still return to
   `reviewing` (round 2, archive, repository gate and required checks follow,
   unchanged in `internal/delivery`).
4. MUST add `internal/spec/declared_paths_test.go` with
   `TestDeclaredTaskPathsReadsCreatesInterfaceAndRecordedPaths`, and
   `internal/cli/review_declared_correction_test.go` with the five tests
   `_techspec.md` → Testing Approach names for task_04, each calling
   `t.Parallel()` first, on temporary repositories with a fake reviewer:
   - `TestACodeCorrectionInsideTheDeclaredPathsReturnsToReview`, subtests
     `archive-record` (the Spec's folder committed with a Task declaring
     `interface: review.txt`, then archived as a record whose
     `source_revision` is that commit) and `legacy-folder` (the Task file
     inside `docs/history/specs/0225-example/` at the candidate): review with
     one finding on `review.txt`, a fix to `review.txt`, dispose
     `--fixed-by`, and the proof returns `Accepted: true`.
   - `TestARecordedPathCountsAsDeclared`: a path only under the Task's
     `## Recorded paths` is admitted.
   - `TestACodeCorrectionOutsideTheDeclaredPathsNamesTheFirstPath`: a fix that
     changes `review.txt` and `other.txt` is refused with the reason naming
     `other.txt`.
   - `TestDeclaredPathsAreReadAtTheRecordedRevision`: a correction that also
     edits the archived Task file, or a working-tree copy, to declare
     `other.txt` does not admit `other.txt` (Invariant 5).
   - `TestUnreadableDeclaredPathsRefuseTheCorrection`: a record whose
     `source_revision` cannot be read, with no remote, is refused with the
     unavailable reason and no error.
5. MUST keep every existing test passing unchanged, including
   `TestACorrectionOutsideTheArchivedSpecIsRefused` and task_02's
   `TestTheCorrectionProofRefusesACodeFixUnderASymlinkedWorktreeRoot`, whose
   asserted reasons stay prefixes of the new one.
6. MUST run the new tests once with `-race` (ADR-0259), give no wait a real
   time bound (ADR-0271), reach no network, and MUST NOT change
   `internal/delivery`, the review command's corrective-Spec line, or any
   command's flag or exit code.

## Subtasks

- [ ] Add `spec.DeclaredTaskPaths` and its test.
- [ ] Collect the declared paths at the recorded revision in the proof and widen the path rule.
- [ ] Add the five correction tests.

## Acceptance Criteria

- [ ] A disposed fix to a declared file returns to review, for a record and for a folder.
- [ ] Any other changed path is refused by name, and the declared set cannot be widened by the correction.
- [ ] Unreadable declared paths refuse without an error.

## Consumers

- `spec.DeclaredTaskPaths` and the review correction proof's path rule — consumers: `./internal/cli`, `./internal/spec`

## Context

- instruction: `docs/adr/0272-a-review-record-survives-a-fix-and-a-correction-inside-the-archived-specs-declared-paths-returns-to-review.md`
- instruction: `docs/adr/0233-a-skill-version-raise-is-regenerated-at-merge-and-a-review-only-correction-returns-to-review.md`
- instruction: `docs/adr/0261-an-archive-record-names-its-pull-request-and-readers-fetch-its-head.md`
- instruction: `docs/user-guide/commands/deliver.md`
- instruction: `internal/spec/task.go`
- instruction: `internal/spec/recorded_paths.go`
- instruction: `internal/spec/source_revision.go`
- instruction: `internal/cli/archive_reader.go`
- instruction: `internal/cli/deliver_review_correction_test.go`
- interface: `internal/cli/deliver_review_correction.go`
- creates: `internal/spec/declared_paths.go`
- creates: `internal/spec/declared_paths_test.go`
- creates: `internal/cli/review_declared_correction_test.go`

## Verification

- `out="$(go test -count=1 -race -v -run '^(TestACodeCorrectionInsideTheDeclaredPathsReturnsToReview|TestARecordedPathCountsAsDeclared|TestACodeCorrectionOutsideTheDeclaredPathsNamesTheFirstPath|TestDeclaredPathsAreReadAtTheRecordedRevision|TestUnreadableDeclaredPathsRefuseTheCorrection)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestACodeCorrectionInsideTheDeclaredPathsReturnsToReview TestACodeCorrectionInsideTheDeclaredPathsReturnsToReview/archive-record TestACodeCorrectionInsideTheDeclaredPathsReturnsToReview/legacy-folder TestARecordedPathCountsAsDeclared TestACodeCorrectionOutsideTheDeclaredPathsNamesTheFirstPath TestDeclaredPathsAreReadAtTheRecordedRevision TestUnreadableDeclaredPathsRefuseTheCorrection; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the tests do not exist, so their PASS lines are missing and the command fails; afterwards they pass under the race detector.
- `out="$(go test -count=1 -race -v -run '^TestDeclaredTaskPathsReadsCreatesInterfaceAndRecordedPaths$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestDeclaredTaskPathsReadsCreatesInterfaceAndRecordedPaths ' || { printf 'missing passing spec test\n' >&2; exit 1; }; grep -q 'archivedSpecDeclaredPaths(' internal/cli/deliver_review_correction.go || { printf 'the proof does not read declared paths\n' >&2; exit 1; }; go test -count=1 -run 'Review|Correction|Dispos|Lineage' ./internal/cli || exit 1; go test -count=1 -run 'Task|Recorded|Archive|Source' ./internal/spec || exit 1` — expected: exit 0. Before this Task the test and the helper are absent, so the command fails; afterwards the consuming CLI and spec suites pass.

## References

- `_prd.md` → Goals; User Story 3; Core Feature 4; Success Metric 4
- `_techspec.md` → Current behavior; API Contract 6; Invariant 4; Invariant 5; Testing Approach; Build Order 4
- ADR-0272; ADR-0233; ADR-0165; ADR-0261; ADR-0259; ADR-0271
