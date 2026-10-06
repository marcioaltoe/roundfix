---
task: task_03
spec: 0242-an-archive-that-leaves-an-archive-record
status: pending
type: backend
complexity: high
---

# Task 03: The Spec check, the audit, the suite guard and the repository tests need no archived Spec folder

## Overview

On 2026-10-06 an ablation removed every Spec folder under
`docs/history/specs` from a copy of this repository. `make verify` failed 17
tests and `make verify-docs` failed 5, including seven user-guide links
(`_techspec.md` → Readers). The Spec check's archived-slug set, the Spec
audit and the suite guard's sanctioned regenerations also read folders only.
This Task makes all of them work from Archive Records, from fixtures they
own, or from Git at the pinned commit
`40a7893d872c8a6705f6d7745e6efe430ae9deeb`. Spec 0243 can then migrate the
history without breaking a gate. It answers the Backlog Entry "History keeps
only what the Secondbrain needs" of 2026-10-06, and it carries forward the
reader analysis of the parked draft Spec 0238.

## Requirements

1. MUST add `internal/gittest/pinned_history.go` with
   `PinnedHistoryRevision` and `PinnedHistory` per Invariant 13. It must
   not write the repository.
2. MUST make the Spec check's archived-slug set
   (`repositoryDirectoryNames` as used by `detectFindingsConsistency` and
   `detectBacklogPromotion`) come from `spec.ArchivedSpecSlugs`, so an
   `absorbed_by` or promoted `spec` naming a record-only Spec resolves. The
   backlog fix text points at the record. `detectTaskContextReferences`
   resolves a Task Context path under an archived Spec through its record and
   `source_revision`, per Invariant 9.
3. MUST make the Spec audit claim the record for an archived Spec
   (`claimedArtifacts`, `archivedSpecArtifactPath`). It must not report a
   removed Spec file as undelivered. `validateSpecAuditSlug` accepts a slug
   whose record exists.
4. MUST make `readSanctionedRegenerations` also read the `regeneration` list
   of every record under `docs/history/specs`, per Invariant 12, without
   importing `internal/spec`. Active and legacy archived authorization
   records keep being read.
5. MUST make each of the 22 tests that failed in the ablation pass on a tree
   without Spec folders, keeping the test's subject:
   - the corpus tests in `internal/spec` and the 0058 replays read the
     corpus through `gittest.PinnedHistory`;
   - `TestCurrentRecordPermitsItsDeclaredOperations` and
     `TestEveryBoundedPathIsGoverned` read the 0119 and 0130 authorization
     records at the pinned revision;
   - `TestReceiptCharacterizationKeepsEveryParsedClaim` compares its
     `testdata` copies with the pinned blobs;
   - the transcript and judge tests read the archived TechSpecs of Specs
     0191, 0205 and 0209 through `gittest.PinnedHistory`;
   - `TestArchiveLayoutCharacterizationRecordsEveryRetiredFamily` accepts an
     absent `docs/history/specs` or one holding only records;
   - the two derived-regeneration contract tests and
     `TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical` tolerate an
     absent `docs/history/specs` and include records;
   - `TestCoverageEquivalence` counts no package under `docs/`, and
     `docs/references/coverage-record.json` is re-recorded with
     `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record`;
   - the seven user-guide links that reach `docs/history/specs` name their
     Spec in prose or link an ADR or a `docs/references` document instead.
6. MUST add the tests named in Verification:
   - `TestPinnedHistoryMaterializesThePinnedPaths`;
   - `TestArchivedSlugsIncludeArchiveRecords`, covering a Finding's
     `absorbed_by` and a promoted Backlog Entry that name a record-only
     Spec;
   - `TestTaskContextResolvesThroughTheArchiveRecord`;
   - `TestAuditClaimsTheArchiveRecord`;
   - `TestSpecAuditAcceptsAnArchiveRecordSlug`;
   - `TestSanctionedRegenerationsReadArchiveRecords`;
   - `TestCoverageRecordCountsNoPackageUnderDocs`.
7. MUST NOT delete, move or rewrite any file under `docs/history`, and MUST
   NOT change `spec.Archive`, the Delivery Queue, reconcile or the review.

## Subtasks

- [ ] Add the pinned-history helper.
- [ ] Make the Spec check, the audit and the suite guard read records.
- [ ] Move the 22 repository tests off the archived folders.
- [ ] Re-record the coverage record and fix the seven user-guide links.
- [ ] Prove the gates' tests on a copy without Spec folders.

## Acceptance Criteria

- [ ] On a copy of the tree with no Spec folder under `docs/history/specs`,
      all 22 tests pass. None of them skips.
- [ ] An archived-slug, Task Context, audit or regeneration read finds a
      record-only Spec.
- [ ] Nothing under `docs/history` changes.

## Context

- creates: `internal/gittest/pinned_history.go`
- creates: `internal/gittest/pinned_history_test.go`
- creates: `internal/speccheck/archive_record_test.go`
- creates: `internal/specaudit/archive_record_test.go`
- creates: `internal/cli/spec_audit_record_test.go`
- creates: `internal/suiteguardcontract/archive_record_test.go`
- interface: `internal/speccheck/citations.go`
- interface: `internal/speccheck/backlog.go`
- interface: `internal/specaudit/audit.go`
- interface: `internal/cli/spec_check.go`
- interface: `internal/suiteguardcontract/regeneration.go`
- interface: `internal/authorization/authorization_test.go`
- interface: `internal/judge/questions_test.go`
- interface: `internal/judge/grouping_test.go`
- interface: `internal/spec/archive_layout_characterization_test.go`
- interface: `internal/spec/spec_test.go`
- interface: `internal/spec/archive_test.go`
- interface: `internal/spec/qa_test.go`
- interface: `internal/spec/coverage_test.go`
- interface: `internal/spec/qa_frontmatter_test.go`
- interface: `internal/speccheck/receipt_characterization_test.go`
- interface: `internal/speccheck/transcripts_test.go`
- interface: `internal/speccheck/governed_repocontract_test.go`
- interface: `internal/baseline/derived_regeneration_repocontract_test.go`
- interface: `internal/baseline/derived_ownership_test.go`
- interface: `skills/owned_skill_edit_repocontract_test.go`
- interface: `docs/references/coverage-record.json`
- interface: `docs/user-guide/usage.md`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/user-guide/commands/detached-runs.md`
- interface: `docs/user-guide/commands/reconcile.md`
- interface: `docs/user-guide/commands/stop.md`
- instruction: `docs/adr/0247-an-archive-leaves-an-archive-record-and-the-spec-folder-stays-in-git.md`
- instruction: `docs/adr/0215-a-measurement-lands-as-a-record-and-archived-evidence-leaves-only-with-approval.md`
- instruction: `docs/references/archived-evidence-measurement.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestPinnedHistoryMaterializesThePinnedPaths|TestArchivedSlugsIncludeArchiveRecords|TestTaskContextResolvesThroughTheArchiveRecord|TestAuditClaimsTheArchiveRecord|TestSpecAuditAcceptsAnArchiveRecordSlug|TestSanctionedRegenerationsReadArchiveRecords|TestCoverageRecordCountsNoPackageUnderDocs)$' ./internal/gittest ./internal/speccheck ./internal/specaudit ./internal/cli ./internal/suiteguardcontract ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestPinnedHistoryMaterializesThePinnedPaths TestArchivedSlugsIncludeArchiveRecords TestTaskContextResolvesThroughTheArchiveRecord TestAuditClaimsTheArchiveRecord TestSpecAuditAcceptsAnArchiveRecordSlug TestSanctionedRegenerationsReadArchiveRecords TestCoverageRecordCountsNoPackageUnderDocs; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name (" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the seven tests do not exist, so the command fails.
- `p="" && tmp="$(mktemp -d)" && objects="$(cd "$(git rev-parse --git-common-dir)" && pwd)/objects" && mkdir "$tmp/r" && git ls-files -co --exclude-standard | grep -v '^docs/history/specs/[^/]*/' | while IFS= read -r p; do if [ -f "$p" ]; then printf '%s\n' "$p"; fi; done | tar -cf - -T - | tar -xf - -C "$tmp/r" && cd "$tmp/r" && git init -q && printf '%s\n' "$objects" > .git/objects/info/alternates && git update-ref refs/pinned 40a7893d872c8a6705f6d7745e6efe430ae9deeb && git add -A && git -c user.name=v -c user.email=v@example.invalid commit -q -m copy && test ! -e docs/history/specs/0001-implement-command && out="$(go test -count=1 -v -run '^(TestCurrentRecordPermitsItsDeclaredOperations|TestQuestionFileLoads|TestGroupingQuestionIsTheMeasuredOne|TestArchiveLayoutCharacterizationRecordsEveryRetiredFamily|TestRepositorySpecCorpusStillLoads|TestQAGateLegacyArchivedManifestsLoadUnchanged|TestArchivedPassCorpusRemainsArchiveEligible|TestArchivedQAOverrideCorpusIncludesFailedSpec|TestSpec0058ReplayArchivesDeclaredUnreachableRelease|TestSpec0058ReplayRefusesUnmatchedBlockedRow|TestSpec0058ReplayReportsWronglyDeclaredReachableRow|TestArchivedQAReportCorpusRemainsReadable|TestCoverageEquivalence|TestQAReportReaderAgreesWithTheDerivedVerificationOnTheArchive|TestReceiptCharacterizationKeepsEveryParsedClaim|TestThisSpecsSurfaceTranscriptsAreWellFormed|TestThisSpecsTranscriptsHaveImplementationAndGateReferences)$' ./internal/authorization ./internal/judge ./internal/spec ./internal/speccheck 2>&1; go test -count=1 -v -tags repocontract -run '^(TestMeasuredSanctionedOwnershipMatchesRecords|TestDeclaredStepRegenerationAndFrozenBoundaries|TestEveryBoundedPathIsGoverned|TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical)$' ./internal/baseline ./internal/speccheck ./skills 2>&1; go test -count=1 -v -tags docscontract -run '^TestUserGuideLinksResolve$' ./internal/docscontract 2>&1)"; printf '%s\n' "$out" | grep -E -- '^(--- FAIL|FAIL)' >&2; for name in TestCurrentRecordPermitsItsDeclaredOperations TestQuestionFileLoads TestGroupingQuestionIsTheMeasuredOne TestArchiveLayoutCharacterizationRecordsEveryRetiredFamily TestRepositorySpecCorpusStillLoads TestQAGateLegacyArchivedManifestsLoadUnchanged TestArchivedPassCorpusRemainsArchiveEligible TestArchivedQAOverrideCorpusIncludesFailedSpec TestSpec0058ReplayArchivesDeclaredUnreachableRelease TestSpec0058ReplayRefusesUnmatchedBlockedRow TestSpec0058ReplayReportsWronglyDeclaredReachableRow TestArchivedQAReportCorpusRemainsReadable TestCoverageEquivalence TestQAReportReaderAgreesWithTheDerivedVerificationOnTheArchive TestReceiptCharacterizationKeepsEveryParsedClaim TestThisSpecsSurfaceTranscriptsAreWellFormed TestThisSpecsTranscriptsHaveImplementationAndGateReferences TestMeasuredSanctionedOwnershipMatchesRecords TestDeclaredStepRegenerationAndFrozenBoundaries TestEveryBoundedPathIsGoverned TestOwnedSkillEditLeavesDerivedArtifactsByteIdentical TestUserGuideLinksResolve; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name (" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0. It copies the working tree without any Spec folder under `docs/history/specs` into a temporary repository that reaches the pinned commit through Git alternates. It then requires a PASS, not a skip, for each of the 22 tests. Measured on 2026-10-06 at `40a7893d`: all 22 fail, so the command fails.

## References

- `_prd.md` → Goal 4; User Story 4; Core Feature 4; Core Feature 5; Success Metric 4; Success Metric 5; Acceptance evidence
- `_techspec.md` → Readers; Invariants 9, 12 and 13; Build Order 3
- ADR-0247; ADR-0215
