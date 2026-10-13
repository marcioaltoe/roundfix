---
task: task_04
spec: 0266-the-review-record-survives-a-fix
status: completed
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

## Result

Implemented the declared-path correction slice. `spec.DeclaredTaskPaths`
uses the existing Task parser to collect `creates:` and `interface:` entries
and `RecordedTaskPaths`, returning clean, sorted, unique repository paths.
Instruction and deletion entries do not grant correction paths; invalid Task
content or unsafe recorded paths return an error.

The proof first applies the existing archive-path rule. When a changed path
lies outside those archives, it reads each named archive's form at the parked
candidate and collects the declared set once. Legacy Task files are read at
that candidate; Archive Record Task files are read under `source` at
`source_revision`, using `spec.EnsureSourceRevision` and the configured
delivery remote for missing objects. All Task listings and contents come
from Git. Exact declared-path matches join the existing archive-path rule.
The first other changed path is refused in Git's listing order with the
required extended reason. Unavailable declarations become a refusal with
the outside-archive path context, their cause and no returned error. A
correction confined to archive paths needs no Task declarations and retains
its existing acceptance even when the record's source revision is absent.

Acceptance evidence from focused implementation checks:

| Acceptance criterion | Evidence |
| --- | --- |
| A disposed fix to a declared file returns to review, for a record and for a folder. | `TestACodeCorrectionInsideTheDeclaredPathsReturnsToReview`, subtests `archive-record` and `legacy-folder`, exercises public review/dispose commands with a fake reviewer in real temporary symlinked Git worktrees and asserts `Accepted: true`. `TestARecordedPathCountsAsDeclared` also admits a recorded-only path in both forms. |
| Any other changed path is refused by name, and the declared set cannot be widened by the correction. | `TestACodeCorrectionOutsideTheDeclaredPathsNamesTheFirstPath` asserts the exact refusal naming `other.txt` when the correction also changes `review.txt`. `TestDeclaredPathsAreReadAtTheRecordedRevision` refuses `other.txt` for both archive forms with committed and uncommitted Task edits, including an uncommitted source-folder copy for records. |
| Unreadable declared paths refuse without an error. | `TestUnreadableDeclaredPathsRefuseTheCorrection` covers a missing source revision with no remote and unparsable Task content in both archive forms. Its proof helper requires no returned error, and each case asserts the unavailable reason with a cause. |

Focused checks run:

- Before implementation, `rtk proxy go test -count=1 -run '^TestACodeCorrectionInsideTheDeclaredPathsReturnsToReview$' ./internal/cli`
  exited 1: both archive forms refused `review.txt` with the original
  outside-archive reason, reproducing the missing behavior.
- After implementation, `rtk proxy go test -count=1 -race -run 'Test(DeclaredTaskPathsReads|ACodeCorrection|ARecordedPathCounts|DeclaredPathsAreRead|UnreadableDeclaredPaths)' ./internal/spec ./internal/cli`
  exited 0: `internal/spec` and `internal/cli` passed. This selection runs all
  six new top-level tests, including the parser checks for deduplication,
  context-kind selection, path cleaning and malformed declarations.
- `rtk proxy go test -count=1 -run 'Test(AReviewOnlyCorrectionIsProvedFromTheReviewRecord|ACorrectionOutsideTheArchivedSpecIsRefused|ACorrectionWithAnUndisposedFindingIsRefused|ACorrectionThatDoesNotDescendIsRefused|ReviewCorrectionRequiresTheParkedReviewRecord|TheCorrectionProof.*)' ./internal/cli`
  exited 0. Existing correction, record, disposition, ancestry and symlinked
  checkout tests passed unchanged, including both required refusal regressions.
- `rtk proxy git diff --check` exited 0.

The new tests use isolated local repositories, temporary homes and a fake
reviewer, have `t.Parallel()` first, configure no remote and add no real-time
wait budget. No network or real reviewer was reached. The implementation
leaves the delivery state machine, round-2/gate/check sequence, corrective-Spec
output line, flags and exit codes unchanged.

The configured sandbox could not start because the writable root contains a
symlink; approved escalated commands were used for worktree access and focused
checks. Task status remains Daemon-owned. Authored Verification, repository
Verification and settlement are left to the Daemon. No commit, push, pull
request, Task Graph edit or other Task-file edit was made. No follow-up work
was identified for this slice.

### Verification Feedback repair — attempt 1

Inspected the Daemon's diagnostic artifact at
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261013T214348Z_2ba6216d89c5b1cf/verification/batch-002-attempt-1.log`
and the existing `TestReviewCorrectionAllowsTheArchiveRecordPath`.
The initial implementation eagerly recovered Task declarations even for a
correction confined to an Archive Record, introducing a source-availability
prerequisite into the existing archive-path proof. The initial focused
selection did not include that record-path regression.

Repaired only the Task's correction proof: archive paths are admitted first;
Task declarations are collected once on the first path outside the named
archives. If declarations cannot be read, the required unavailable reason
also names that outside-archive path. This preserves the existing sibling-path
refusal while retaining conservative refusal for code paths with unavailable
declarations. Existing tests were left unchanged.

Focused feedback checks:

- Before the repair, `rtk proxy go test -count=1 -run '^TestReviewCorrectionAllowsTheArchiveRecordPath$' ./internal/cli`
  exited 1 and reproduced the missing-source archive-record regression.
- After the repair, `rtk proxy go test -count=1 -race -run 'Test(ReviewCorrectionAllowsTheArchiveRecordPath|DeclaredTaskPathsReads|ACodeCorrection|ARecordedPathCounts|DeclaredPathsAreRead|UnreadableDeclaredPaths|ACorrectionOutsideTheArchivedSpecIsRefused|TheCorrectionProofRefusesACodeFixUnderASymlinkedWorktreeRoot)' ./internal/cli ./internal/spec`
  exited 0 for both packages. The existing record-path test covers source
  present/absent acceptance and sibling-path refusal. All six new top-level
  tests and both required outside-path refusal regressions passed under the
  race detector. This refreshes evidence for all three acceptance criteria
  after the proof repair.
- Final `rtk proxy git diff --check` exited 0. Read-only scope inspection
  confirmed exactly the same five task_04 paths, no staged changes, and
  byte-identical authored Task content apart from the pre-existing
  Daemon-owned `status: in_progress` value; only `## Result` was revised.

No Task-declared Verification command was rerun. The Daemon retains status,
full Verification and settlement ownership. No network or real reviewer was
reached, and no commits, pushes or pull requests were created.
