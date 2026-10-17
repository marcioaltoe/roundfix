---
task: task_03
spec: 0272-the-history-holds-only-archive-records
status: pending
type: backend
complexity: high
---

# Task 03: The check holds the History Root to Archive Records and links to files that exist

## Overview

Nothing reports a file left in a retired family, a terminal Finding left among
unresolved ones or a relative link to a missing file, and 40 links in Archive
Records are already broken (`_techspec.md` → Current behavior). This Task adds
`MissingRelativeLinks` and the `record-links` unit that fixes broken record
links through both history commands, the gated `SC-HISTORY-RETAINED` error
with its Git Locator diagnostic, the Git fallbacks of Claim Receipts and Rollup
members, and the delete wording of `SC-BACKLOG-UNMOVED`. It then brings the
Roundfix and archive-spec skills and the glossary up to date with what
task_01 to task_03 deliver. It is verifiable on its own through planted
fixtures, both commands and the skill and glossary checks.

## Requirements

1. MUST implement `MissingRelativeLinks`, the exported
   `RetirementRuleSentinel` constant and API Contract 11 in
   `internal/spec/record_links.go`, and the `record-links` unit in both
   history commands after the removal units, skipping a record the same plan
   removes, with the line and JSON shapes of API Contract 12.
2. MUST implement API Contract 13 as `SC-HISTORY-RETAINED` in
   `internal/speccheck/history_retained.go`, wired beside the repository-level
   detectors in `constraints.go`, registered at stage `prd` in `coherence.go`,
   gated on `spec.RetirementRuleSentinel` in the recorded `guide.docs-layout`
   region, and listed as skipped with its reason when the gate is closed. Its
   inputs are this repository's files and local Git only; its diagnostic
   prints the replacement `ResolveGitLocator` returns.
3. MUST implement API Contract 14: the Git fallback for Claim Receipts
   (removing `docs/history/adr` from the filesystem loop), the Git fallback
   for Rollup members, and the Fix of a terminal Backlog Entry under
   `SC-BACKLOG-UNMOVED` saying to delete it. `SC-ARCHIVE-LICENSE` is
   unchanged.
4. MUST add `SC-HISTORY-RETAINED` to the active corpus code set at 0:
   `internal/docscontract/corpus_test.go`'s `corpusFindingCodes`,
   `internal/docscontract/testdata/corpus-golden.json` with a sentence in its
   `update` text, and the pin in
   `TestArchiveLayoutCharacterizationPinsCorpusGoldenAfterSpec0095`.
5. MUST change these existing tests to the new contract and name each change
   in the Result: `TestReceiptsResolveArchivedAndInactiveRecords` and
   `TestReceiptSourcesRefuseAmbiguousRecordsAndADRDirectoryLinks`
   (`internal/speccheck/receipts_test.go`), `TestCheckRollupMember`
   (`internal/speccheck/citations_test.go`), `TestCheckBacklogUnmoved` and
   `TestTerminalBacklogEntryLeftActiveIsUnmoved`
   (`internal/speccheck/backlog_test.go`) and
   `TestADeferredBacklogEntryLeftActiveIsReported`
   (`internal/speccheck/backlog_deferred_test.go`). Every other existing test
   MUST pass unchanged.
6. MUST add the tests named in Verification, each calling `t.Parallel()`
   first. In `internal/spec/record_links_test.go`:
   `TestMissingRelativeLinksFindsBrokenAndSelfLinks` (code spans, fences, URLs
   and Git Locators are not links) and
   `TestPlanRecordLinksResolvesAgainstTheRecordThenItsSource`. In
   `internal/speccheck/history_retained_test.go`:
   `TestHistoryRetainedIsSilentBeforeTheBaselineCarriesTheRule`,
   `TestHistoryRetainedReportsRetainedFamiliesAndTerminalFindings` (a planted
   Reduced History Entry, a file in each retired family and a `done` Finding),
   `TestHistoryRetainedReportsBrokenAndSelfLinksWithTheirGitLocator`,
   `TestHistoryRetainedReportsAGitLocatorNoReachableCommitHolds`,
   `TestReceiptsResolveARemovedADRThroughGit`,
   `TestRollupMembersResolveThroughGit` and
   `TestTerminalBacklogEntryFixSaysDelete`. In
   `internal/cli/history_record_links_test.go`:
   `TestHistorySanitizeRewritesBrokenRecordLinks` and
   `TestBaselineUpdateRewritesBrokenRecordLinks`.
7. MUST document `SC-HISTORY-RETAINED` in `docs/user-guide/commands/spec.md`
   and `.agents/skills/roundfix/references/spec.md`, and the `record-links`
   unit in `docs/user-guide/commands/history.md` and
   `docs/user-guide/commands/baseline.md`, with the phrase `record-links`.
8. MUST describe deletion and Git Locators in the Roundfix skill references
   `.agents/skills/roundfix/references/archive.md` (Archive Command and
   History Sanitize Command, without `reduces retired Findings`),
   `baseline.md`, `deliver.md` and `setup.md` (Pending History includes
   Reduced History Entries), each with the phrase `Git Locator` except
   `setup.md`, which carries `Reduced History Entry`, and in the Archive Relink
   section of `.agents/skills/archive-spec/SKILL.md`, with the phrase
   `Git Locator`. No text inside any `### QA settlement` section may change.
   Run `make skills-sync`, record both skills' versions with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   and re-record the Behavior Surface Record with
   `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`.
9. MUST, through `domain-modeling`, add **Git Locator** to `CONTEXT.md` near
   **History Citation**: the inline code `<path>@<commit>` that names a removed
   file at a commit that holds it and that the default branch or a tag
   reaches, written in place of a link to it, or its path and the record's
   pull request when no such commit exists, citing ADR-0278. It MUST revise,
   each citing ADR-0278: **History Root** (only Archive Records), **History
   Sanitize Command** (removes retired families and rewrites their citations
   to Git Locators), **Reduced History Entry** (a legacy shape the command
   removes; no retirement writes one), **Rollup** (retires by deletion; its
   members resolve as unresolved Findings or in Git), **Backlog Entry** (a
   terminal entry is deleted), **Legacy Link Rebase** (a folder destination
   without a promoted copy becomes a Git Locator, never the record),
   **History Relocation** (never into the Findings, Backlog, ADR or Review
   families) and **History Citation** (rewritten to a Git Locator in Archive
   Records, Live Documents and the user guide).
10. MUST, when this Task removes or renames a test, regenerate the test
    inventory with `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`
    after its last test change, so `TestCoverageEquivalence` reports no
    regression.

## Subtasks

- [ ] Add `MissingRelativeLinks` and the `record-links` unit to both commands.
- [ ] Add the gated `SC-HISTORY-RETAINED` detector and its diagnostic.
- [ ] Add the Git fallbacks for receipts and Rollup members and the delete Fix.
- [ ] Add the code to the corpus golden and its pin.
- [ ] Update the skills, their mirrors, versions and the Behavior Surface Record.
- [ ] Add and revise the glossary entries.

## Acceptance Criteria

- [ ] A planted file in each retired family, a planted Reduced History Entry, a `done` Finding and a broken or self record link each fail the check once the guide carries the rule, and the same tree passes without that sentence.
- [ ] Each diagnostic names the Git Locator or the deletion to make.
- [ ] Both commands rewrite a broken record link to its Git Locator.
- [ ] The skills and the glossary describe deletion and Git Locators.

## Consumers

- Spec Consistency Check findings and the active corpus counts — consumers: `./internal/speccheck`, `./internal/docscontract`, `./internal/spec`
- `record-links` unit in both history commands — consumers: `./internal/cli`, `./internal/spec`
- Owned skill text and versions — consumers: `./skills`

## Context

- instruction: `docs/adr/0278-a-retirement-deletes-and-the-history-root-holds-only-archive-records.md`
- instruction: `docs/adr/0183-an-attributed-claim-carries-a-receipt-the-check-proves.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- creates: `internal/spec/record_links.go`
- creates: `internal/spec/record_links_test.go`
- creates: `internal/speccheck/history_retained.go`
- creates: `internal/speccheck/history_retained_test.go`
- creates: `internal/cli/history_record_links_test.go`
- interface: `internal/speccheck/constraints.go`
- interface: `internal/speccheck/coherence.go`
- interface: `internal/speccheck/receipts.go`
- interface: `internal/speccheck/receipts_test.go`
- interface: `internal/speccheck/citations.go`
- interface: `internal/speccheck/citations_test.go`
- interface: `internal/speccheck/backlog.go`
- interface: `internal/speccheck/backlog_test.go`
- interface: `internal/speccheck/backlog_deferred_test.go`
- interface: `internal/docscontract/corpus_test.go`
- interface: `internal/docscontract/testdata/corpus-golden.json`
- interface: `internal/spec/archive_layout_characterization_test.go`
- interface: `internal/cli/history.go`
- interface: `internal/cli/baseline_update_history.go`
- interface: `docs/user-guide/commands/spec.md`
- interface: `docs/user-guide/commands/history.md`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `.agents/skills/roundfix/references/archive.md`
- interface: `skills/roundfix/references/archive.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/references/spec.md`
- interface: `skills/roundfix/references/spec.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/archive-spec/SKILL.md`
- interface: `skills/archive-spec/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`
- interface: `docs/references/coverage-record.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestMissingRelativeLinksFindsBrokenAndSelfLinks|TestPlanRecordLinksResolvesAgainstTheRecordThenItsSource)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestMissingRelativeLinksFindsBrokenAndSelfLinks TestPlanRecordLinksResolvesAgainstTheRecordThenItsSource; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^(TestHistoryRetainedIsSilentBeforeTheBaselineCarriesTheRule|TestHistoryRetainedReportsRetainedFamiliesAndTerminalFindings|TestHistoryRetainedReportsBrokenAndSelfLinksWithTheirGitLocator|TestHistoryRetainedReportsAGitLocatorNoReachableCommitHolds|TestReceiptsResolveARemovedADRThroughGit|TestRollupMembersResolveThroughGit|TestTerminalBacklogEntryFixSaysDelete)$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestHistoryRetainedIsSilentBeforeTheBaselineCarriesTheRule TestHistoryRetainedReportsRetainedFamiliesAndTerminalFindings TestHistoryRetainedReportsBrokenAndSelfLinksWithTheirGitLocator TestHistoryRetainedReportsAGitLocatorNoReachableCommitHolds TestReceiptsResolveARemovedADRThroughGit TestRollupMembersResolveThroughGit TestTerminalBacklogEntryFixSaysDelete; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run 'Receipt|Rollup|Backlog|Citation|Archive|History|Finding' ./internal/speccheck || exit 1; go test -count=1 -run 'ArchiveLayoutCharacterization|Archive|History' ./internal/spec || exit 1; go test -count=1 -tags docscontract -run '^TestCheckCorpusGolden$' ./internal/docscontract || exit 1` — expected: exit 0. Before this Task the nine tests do not exist, so the PASS checks fail.
- `out="$(go test -count=1 -v -run '^(TestHistorySanitizeRewritesBrokenRecordLinks|TestBaselineUpdateRewritesBrokenRecordLinks)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestHistorySanitizeRewritesBrokenRecordLinks TestBaselineUpdateRewritesBrokenRecordLinks; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run 'History|BaselineUpdate|Upgrade|SpecCheck' ./internal/cli || exit 1; for f in docs/user-guide/commands/spec.md .agents/skills/roundfix/references/spec.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'SC-HISTORY-RETAINED' || { printf 'missing SC-HISTORY-RETAINED in %s\n' "$f" >&2; exit 1; }; done; for f in docs/user-guide/commands/history.md docs/user-guide/commands/baseline.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'record-links' || { printf 'missing record-links in %s\n' "$f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/archive.md .agents/skills/roundfix/references/baseline.md .agents/skills/roundfix/references/deliver.md .agents/skills/archive-spec/SKILL.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Git Locator' || { printf 'missing Git Locator in %s\n' "$f" >&2; exit 1; }; done; if tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/archive.md | grep -qF -- 'reduces retired Findings'; then printf 'archive reference still reduces\n' >&2; exit 1; fi; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/setup.md | grep -qF -- 'Reduced History Entry' || { printf 'setup reference does not name reduced entries\n' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills || exit 1; go test -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract || exit 1` — expected: exit 0. Before this Task the two command tests do not exist and neither the guides nor the skills name the new code, unit or Git Locator, so the command fails.
- `tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Git Locator**:' || { printf 'CONTEXT.md does not define **Git Locator**\n' >&2; exit 1; }; for term in '**Git Locator**' '**History Root**' '**History Sanitize Command**' '**Reduced History Entry**' '**Rollup**' '**Backlog Entry**' '**Legacy Link Rebase**' '**History Relocation**' '**History Citation**'; do grep -A1 -F -x -- "$term:" CONTEXT.md | grep -qF -- '0278-a-retirement-deletes' || { printf 'CONTEXT.md entry %s does not cite ADR-0278\n' "$term" >&2; exit 1; }; done` — expected: exit 0. Before this Task the glossary lacks **Git Locator** and no entry cites ADR-0278, so the command fails.

## References

- `_prd.md` → Goals; User Story 2; Core Feature 2; Core Feature 7; Success Metric 2; Success Metric 3; Glossary
- `_techspec.md` → API Contract 11; API Contract 12; API Contract 13; API Contract 14; Invariants 1, 2 and 5; Testing Approach; Build Order 3; Glossary
- ADR-0278; ADR-0183; ADR-0244; ADR-0256; ADR-0261
