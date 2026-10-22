---
task: task_01
spec: 0308-recorded-paths-are-read-before-carry-forward-provenance
status: completed
type: backend
complexity: medium
---

# Task 01: A carried Task's recorded paths are read and replaced before its provenance

## Overview

The reader of a Task file's `## Recorded paths` section accepts it only when
it is the last level-2 section, and Task Carry-Forward appends `##
Carry-forward provenance` after it. Every carried Task therefore loses its
recorded paths: the corrective retry refuses a correction inside them,
reconcile keeps a Run Worktree dirty over them, and recording paths again
appends a second section. This Task makes the reader find the section among
the trailing Daemon-written settlement sections and makes the writer replace
it in place.

## Requirements

1. MUST answer the oraculum Inbox Entry of 2026-10-22,
   [`references/2026-10-22-retry-corrective-spec-required-ignora-recorded-paths.md`](references/2026-10-22-retry-corrective-spec-required-ignora-recorded-paths.md)
   (`_prd.md` → Sources 1), by implementing API Contracts 1-4 of
   `_techspec.md` in `internal/spec/recorded_paths.go`.
2. MUST replace `recordedPathsSectionOffset` with `recordedPathsSections`
   exactly as API Contract 1 states: unfenced level-2 sections with today's
   fence and heading rules, each starting at the newline before its heading,
   collected walking back from the end over `## Result` and `## Carry-forward
   provenance`, stopping at any other heading.
3. MUST make `RecordedTaskPaths` return the union of the collected sections'
   paths in file order without duplicates (API Contract 2), and keep its
   result for a file that ends with the recorded section.
4. MUST make `RecordTaskPaths` replace the first collected section in place,
   remove the others, keep every other byte, append at the end when none is
   collected, and skip the write when nothing changes (API Contract 3). Its
   input checks and the section's bytes stay as they are.
5. MUST keep API Contract 4: `DeclaredTaskPaths`, `ProveReviewCorrection` and
   its refusal texts, `dirtyPathsInArchivedSpec`, the Verification path
   provision, `commitTask`, `settleTaskAndCommit`, `RecordCarryForward` and
   the settlement heading classifiers keep their code. The corrective-scope
   rule is unchanged: a path no Task declares or records is still refused.
6. MUST state the hazard in the tests: the writer rewrites committed Task
   files during settlement, so a sabotaged writer that drops the provenance
   or duplicates the section must fail a named row.
7. MUST change the `Invariant:` header line of
   `internal/spec/recorded_paths_section_test.go` to "Only an exact unfenced
   Recorded paths heading followed by nothing but Daemon-written settlement
   sections is Daemon-owned." and add there, each top-level test calling
   `t.Parallel()` first, with synthetic fixtures:
   - `TestRecordedPathsAreReadBeforeTrailingSettlementSections` with the rows
     `carry-forward-provenance`, `result-then-provenance`,
     `duplicate-sections`, `authored-section-between` and
     `fenced-heading-before-provenance` of `_techspec.md` → Testing Approach,
     asserting `RecordedTaskPaths` and whether `DeclaredTaskPaths` holds the
     recorded path;
   - `TestRecordTaskPathsReplacesTheSectionBeforeCarryForwardProvenance` with
     the rows `same-paths` (bytes and file identity unchanged), `new-paths`,
     `no-paths` and `duplicate-sections`, each asserting the exact bytes
     written and the provenance kept byte for byte.
8. MUST add `TestARecordedPathBeforeCarryForwardProvenanceCountsAsDeclared` to
   `internal/cli/review_declared_correction_test.go`, with `archive-record`
   and `legacy-folder` subtests over `declaredCorrectionTask("", true)` plus a
   `## Carry-forward provenance` section: the real `ProveReviewCorrection`
   accepts a correction of `review.txt`, and refuses a correction that also
   changes `other.txt` with exactly `changed path other.txt is outside the
   named archived Specs and the paths their Tasks declared`.
9. MUST add `TestArchivedSpecLeftoverRecordedBeforeCarryForwardProvenance` to
   `internal/worktree/merged_spec_leftovers_test.go`: the archived Task of
   `newArchivedLeftoverFixture` rewritten to end with `## Recorded paths`
   listing `recorded.txt` then `## Carry-forward provenance`, `recorded.txt`
   dirty in the Run Worktree, reconciliation state `superseded`.
10. MUST NOT change any other existing test, golden or fixture, any guide,
    skill or glossary entry.

## Subtasks

- [ ] Replace the last-heading locator with the trailing settlement run.
- [ ] Read the union of the collected sections and replace them in place.
- [ ] Add the reader and writer rows in `internal/spec`.
- [ ] Add the corrective-retry and reconcile tests for the carried shape.

## Acceptance Criteria

- [ ] A carried Task's recorded path is accepted by the corrective retry for an Archive Record and a legacy folder, and a path in no Task is refused with today's text.
- [ ] Reconcile supersedes a Run Worktree whose only dirty path is a carried Task's recorded path.
- [ ] Recording paths on a carried Task leaves one section and the provenance bytes unchanged.
- [ ] An authored or fenced recorded heading before the provenance is not read, and every existing recorded-path test passes unchanged.

## Consumers

- `spec.RecordedTaskPaths` and `spec.DeclaredTaskPaths` results for a carried Task — consumers: `./internal/cli`, `./internal/worktree`, `./internal/speccheck`
- `spec.RecordTaskPaths` bytes for a carried Task — consumers: `./internal/daemon`, `./internal/cli`

## Context

- instruction: `docs/adr/0166-the-daemon-records-the-paths-a-task-changed-without-declaring-them.md`
- instruction: `internal/spec/spec.go`
- instruction: `internal/spec/declared_paths.go`
- instruction: `internal/spec/recorded_paths_test.go`
- instruction: `internal/spec/declared_paths_test.go`
- instruction: `internal/cli/deliver_review_correction.go`
- instruction: `internal/cli/review_resolved_checkout_test.go`
- instruction: `internal/cli/settle.go`
- instruction: `internal/cli/settle_recorded_paths_test.go`
- instruction: `internal/cli/review_conventions.go`
- instruction: `internal/worktree/merged_head.go`
- instruction: `internal/speccheck/verification_path.go`
- instruction: `internal/speccheck/verification_path_test.go`
- instruction: `internal/daemon/task_engine.go`
- instruction: `internal/daemon/recorded_paths_test.go`
- interface: `internal/spec/recorded_paths.go`
- interface: `internal/spec/recorded_paths_section_test.go`
- interface: `internal/cli/review_declared_correction_test.go`
- interface: `internal/worktree/merged_spec_leftovers_test.go`

## Verification

- `out="$(go test -timeout 10m -count=1 -v -run '^(TestRecordedPathsAreReadBeforeTrailingSettlementSections|TestRecordTaskPathsReplacesTheSectionBeforeCarryForwardProvenance|TestRecordedPathsKeepAFencedHeadingInAuthoredText|TestRecordedPathsKeepAnAuthoredHeadingFollowedBySections|TestRecordedPathsReplaceTheTrailingDaemonSection|TestRecordTaskPathsReplacesAnExistingSectionAndKeepsEveryOtherByte|TestRecordTaskPathsWithNoPathsRemovesTheSection|TestDeclaredTaskPathsReadsCreatesInterfaceAndRecordedPaths|TestARecordedPathBeforeCarryForwardProvenanceCountsAsDeclared|TestARecordedPathCountsAsDeclared|TestACodeCorrectionOutsideTheDeclaredPathsNamesTheFirstPath|TestDeclaredPathsAreReadAtTheRecordedRevision|TestSettleRecordsTheUndeclaredPathsOfItsTaskCommit|TestArchivedSpecLeftoverRecordedBeforeCarryForwardProvenance|TestArchivedSpecLeftoverScopeKinds|TestTaskCommitRecordsAnUndeclaredTestFile|TestTaskCommitRecordsNothingWhenEveryPathIsDeclared|TestVerificationPathUndeclaredIgnoresWhatThePlanProvides)$' ./internal/spec ./internal/cli ./internal/worktree ./internal/daemon ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRecordedPathsAreReadBeforeTrailingSettlementSections TestRecordTaskPathsReplacesTheSectionBeforeCarryForwardProvenance TestRecordedPathsKeepAFencedHeadingInAuthoredText TestRecordedPathsKeepAnAuthoredHeadingFollowedBySections TestRecordedPathsReplaceTheTrailingDaemonSection TestRecordTaskPathsReplacesAnExistingSectionAndKeepsEveryOtherByte TestRecordTaskPathsWithNoPathsRemovesTheSection TestDeclaredTaskPathsReadsCreatesInterfaceAndRecordedPaths TestARecordedPathBeforeCarryForwardProvenanceCountsAsDeclared TestARecordedPathCountsAsDeclared TestACodeCorrectionOutsideTheDeclaredPathsNamesTheFirstPath TestDeclaredPathsAreReadAtTheRecordedRevision TestSettleRecordsTheUndeclaredPathsOfItsTaskCommit TestArchivedSpecLeftoverRecordedBeforeCarryForwardProvenance TestArchivedSpecLeftoverScopeKinds TestTaskCommitRecordsAnUndeclaredTestFile TestTaskCommitRecordsNothingWhenEveryPathIsDeclared TestVerificationPathUndeclaredIgnoresWhatThePlanProvides; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go vet ./internal/spec ./internal/cli ./internal/worktree` — expected: exit 0. Before this Task the four new tests do not exist, so the command fails at its first pass check. Afterwards the new reader, writer, corrective-retry and reconcile tests pass with the unchanged recorded-path, settlement, Task commit and Verification path tests of the consuming packages, and the three packages' tests compile for Windows.

## References

- `_prd.md` → Sources; Goals; Core Features 1-3; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5
- `_techspec.md` → Current behavior; API Contracts 1-4; Invariants 1-4; Testing Approach; Build Order 1
- ADR-0166


## Result

Implemented the trailing settlement-section reader and in-place writer in
`internal/spec/recorded_paths.go`. The helper keeps the existing fence and
exact-heading rules, walks back over Daemon settlement sections, and returns
recorded section bounds in file order. Reading unions paths without duplicates;
writing replaces the first collected section, removes the rest, preserves all
other bytes, and avoids replacement when the bytes are unchanged. Input checks
and generated section bytes are unchanged. No consumer implementation changed.

Acceptance evidence from focused checks:

- Corrective retry: `TestARecordedPathBeforeCarryForwardProvenanceCountsAsDeclared`
  accepts `review.txt` through the real proof for `archive-record` and
  `legacy-folder`; each also refuses the added `other.txt` with exactly
  `changed path other.txt is outside the named archived Specs and the paths their Tasks declared`.
- Reconcile: `TestArchivedSpecLeftoverRecordedBeforeCarryForwardProvenance`
  classifies the synthetic Run Worktree with dirty `recorded.txt` as
  `superseded` when the archived Task records it before provenance.
- Writer preservation: all four rows of
  `TestRecordTaskPathsReplacesTheSectionBeforeCarryForwardProvenance` assert
  exact output and provenance bytes. `same-paths` also checks file identity;
  `duplicate-sections` collapses the two sections at the first one's position.
- Ownership: all five rows of
  `TestRecordedPathsAreReadBeforeTrailingSettlementSections` assert recorded
  paths and declarations, including file-order deduplication and exclusion of
  authored/fenced headings. Existing tests were left unchanged except the
  required suite Invariant line.

Focused commands and outcomes:

- Before changing production code,
  `go test ./internal/spec -run 'TestRecordedPathsAreReadBeforeTrailingSettlementSections|TestRecordTaskPathsReplacesTheSectionBeforeCarryForwardProvenance' -count=1`
  exited 1: carried reader rows lost paths, and writer rows retained or
  duplicated the old section.
- The four new tests together passed in `./internal/spec`, `./internal/cli`
  and `./internal/worktree` after the implementation change.
- Final focused regression check:
  `go test ./internal/spec ./internal/cli ./internal/worktree -run 'Recorded|RecordTaskPaths|DeclaredTaskPaths|DeclaredPaths|CodeCorrection|ArchivedSpecLeftover' -count=1`
  exited 0 in all three packages. This covers the existing recorded-path,
  declared-path, corrective-scope and archived-leftover tests selected by
  that expression, alongside the new tests.
- Temporary `go test -overlay` sabotage checks in `./internal/spec` each
  exited 1 at the intended named writer row: discarding the suffix failed
  `new-paths`; always appending instead of replacing failed `same-paths`.
  Overlays were created outside the repository and removed afterwards.
- `git diff --check` exited 0.

The Run checkout uses a symlinked writable root that prevents sandbox startup;
focused commands used the resolved checkout with execution escalation. Task
status remains Daemon-owned. Declared Verification, Windows vet, repository
Verification and the terminal QA gate were not run; the Daemon owns those
checks and settlement. No commit, push or pull request was made. No follow-up
implementation outside this slice was added.
