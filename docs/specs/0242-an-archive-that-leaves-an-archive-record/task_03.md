---
task: task_03
spec: 0242-an-archive-that-leaves-an-archive-record
status: completed
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


## Result

Implemented the Task 03 reader and repository-test slice. Task status and the
authored shell Verification remain Daemon-owned; neither declared Verification
command was run in this turn. No commit, push or Pull Request was made.

### Implementation

- `gittest.PinnedHistory` materializes only the requested paths with `git archive`
  and `archive/tar`, handles Git's global PAX header, preserves file bytes and
  modes, and uses the test's context and temporary directory. The pin is the
  shared `PinnedHistoryRevision` constant; an unavailable commit names the pin
  in its skip reason. The helper never writes the source repository.
- The Spec check reads archived slugs through `spec.ArchivedSpecSlugs`.
  Record-only Finding absorption resolves, and a promoted Backlog Entry's fix
  names the Archive Record rather than a removed references directory.
  Task Context resolves both former active and legacy archive spellings through
  the record's `source` and `source_revision`. Missing records, revisions and
  unresolved historical files remain unresolved.
- The audit claims the Archive Record, collapses retired active Spec file
  claims onto that record, and retains legacy-folder claims when no record is
  claimed. The CLI accepts a valid record-only slug in built-in and custom
  archive roots. Legacy audit fixtures retain their original subject.
- The suite guard reads record frontmatter with its own YAML projection and
  does not import `internal/spec`; active and legacy authorization declarations
  still participate in the same output resolution and ordering.
- Corpus, authorization, receipt, transcript, judge and 0058 replay tests read
  their historical subjects at the pin. Regeneration copies and archived-byte
  checks accept an absent archive root and retain records when present. Archive
  layout characterization accepts an absent or record-only Spec archive root.
- Coverage excludes packages under `docs/`. The record was regenerated by its
  existing command, removing the three formerly recorded docs packages and
  retaining every previously recorded test in the remaining packages. Its
  other additions reflect the current repository's test inventory. The seven
  user-guide links now name their historical Specs in prose.

### Focused checks

Checks used `GOCACHE=/private/tmp/roundfix-task03-cache` and real temporary Git
repositories. The disposable folder-free tree is
`/private/tmp/roundfix-task03-no-folders-ezu5ukvn/repository`; it contains the
working files, omits every Spec folder under `docs/history/specs`, and reaches
`PinnedHistoryRevision` through Git alternates. Only its temporary Git index
was populated for repository-copy fixtures; the work branch was not committed.

| Acceptance criterion | Implementation and observed evidence |
| --- | --- |
| All 22 named ablation tests pass without Spec folders, without a top-level skip | The corpus selection, regeneration contracts, governed-path contract and guide-link selection below produced a top-level PASS for every one of the 22 names. No named test was skipped. JSON event logs retain the individual results. |
| Archived-slug, Task Context, audit and regeneration reads find a record-only Spec | All seven newly required tests passed in the final folder-free copy. They exercise absorption and promotion, built-in and custom Task Context roots and path spellings, unavailable Git provenance, delivered and branch-held audit records, malformed CLI records, and record declarations beside active and legacy authorizations. |
| Nothing under `docs/history` changes | `git -c core.fsmonitor=false diff --exit-code -- docs/history` exited 0, and `git -c core.fsmonitor=false status --short -- docs/history` returned no paths. All historical materialization and folder removal happened in temporary directories. |

Focused commands and outcomes:

- `go test -count=1 -json ./internal/spec ./internal/authorization ./internal/judge ./internal/speccheck -run 'Archived|RepositorySpecCorpus|Spec0058Replay|QAGateLegacy|ArchiveLayout|QAReportReader|Coverage|CurrentRecord|QuestionFileLoads|GroupingQuestionIs|ReceiptCharacterization|ThisSpecs'`
  exited 0 in the folder-free copy. The 17 named corpus tests passed.
  Events: `/private/tmp/roundfix-task03-corpus.jsonl`.
- `go test -count=1 -json -tags repocontract ./internal/baseline ./internal/speccheck ./skills -run 'MeasuredSanctionedOwnership|DeclaredStepRegeneration|EveryBoundedPath|OwnedSkillEdit'`
  recorded PASS for both Baseline regeneration contracts and the owned-skill
  byte-identity contract. The first governed-path attempt exposed that merging
  historical records must restore sorted order; that was corrected, then
  `go test -count=1 -json -tags repocontract ./internal/speccheck -run 'EveryBoundedPath'`
  exited 0. Events: `/private/tmp/roundfix-task03-contracts.jsonl` and
  `/private/tmp/roundfix-task03-governed-final.jsonl`.
- `go test -count=1 -json -tags docscontract ./internal/docscontract -run 'UserGuideLinks'`
  exited 0 in the folder-free copy. Events:
  `/private/tmp/roundfix-task03-links.jsonl`.
- `go test -count=1 ./internal/gittest ./internal/speccheck ./internal/specaudit ./internal/suiteguardcontract`
  exited 0 in the final folder-free copy, as did
  `go test -count=1 ./internal/cli -run 'SpecAudit|ValidateSpecAudit'`.
- `go test -count=1 -json ./internal/gittest ./internal/speccheck ./internal/specaudit ./internal/cli ./internal/suiteguardcontract ./internal/spec -run 'ArchiveRecord|PinnedHistory|ArchivedSlugs|AuditClaims|CoverageRecordCounts'`
  exited 0 in the final folder-free copy, with PASS events for all seven new
  required test names. Events: `/private/tmp/roundfix-task03-seams-final.jsonl`.
- A read-only `go test -overlay` probe restored the five original reader
  implementations from `HEAD` while retaining the new regression tests. The
  absorption/context, audit, CLI slug and regeneration tests all failed on
  their intended missing-record behavior. This records the red signal without
  changing repository files. Output: `/private/tmp/roundfix-task03-red.txt`.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record`
  generated the record. On its first run the test passed but the suite guard
  rejected the changed output (exit 1); the deterministic second run wrote
  identical bytes and exited 0. No generated value was hand-edited.
- `git -c core.fsmonitor=false diff --check` exited 0.

### Incremental check and follow-up

The sandboxed `GOCACHE=/private/tmp/roundfix-task03-cache rtk make verify-incremental`
exited 2. The analyzer and all Task 03 packages passed; the only failures were
`TestRunForceStopLegacyRunWithoutOwnerIdentityStillStopsOwner` and
`TestRunForceStopOwnerProcessIntegrationProvesExitBeforeStoreCompletion`, both
reporting `read process table for non-session owner: operation not permitted`.
Output: `/private/tmp/roundfix-task03-incremental.log`.

The same incremental command rerun with host process access exited 0. It
passed the CLI integration tests, the test suite, formatting and analyzer
checks, skill sync/checks and the build. Existing cache results were reused
where the inputs were unchanged. Output:
`/private/tmp/roundfix-task03-incremental-host.log`.

Follow-up outside this slice: the existing coverage generator needs a
sanctioned-regeneration declaration for a changed output. The consuming
`_authorization.md` and the suite guard's declared-command registration were
not widened here. The first generator attempt and its boundary diagnostic
remain recorded above.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/specaudit/audit_test.go`

## Carry-forward provenance

- Source Run: `run_20261006T232220Z_7c2e4a372c1dd72e`
- Source commit: `2b7ec9ab36c0a19fc1aa68acf93c7c989fa4962c`
