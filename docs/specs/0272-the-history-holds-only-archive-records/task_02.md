---
task: task_02
spec: 0272-the-history-holds-only-archive-records
status: pending
type: backend
complexity: high
---

# Task 02: The history commands remove every retired family and rewrite what cites it

## Overview

The History Sanitize Command and the history section of `roundfix baseline
update` reduce Findings and Backlog Entries, never plan retired ADRs, leave
superseded records in place, only list the citations they break, and the
Baseline relocates retired files into families a later run must clean
(`_techspec.md` → Current behavior). This Task makes both commands remove every
retired family wherever an earlier layout left it, including the `superseded`
records nothing cites, prove each removed path is reachable in Git, rewrite
each citation of a removed path to its Git Locator, and stops the Baseline
from relocating into retired families. It is verifiable on its own through
both commands, their text and JSON, a second run, and the upgrade notice.

## Requirements

1. MUST implement API Contract 7: `PlanHistoryKinds(ctx, repositoryRoot)`
   plans `findings`, `backlog`, `adr`, `reviews` and `handoffs`, each with the
   action `remove`, from the History Root family, its legacy trees and the
   active-directory leftovers the contract lists; a Reduced History Entry is
   pending; a superseded ADR whose successor does not name it refuses the `adr`
   unit by name; `ApplyHistoryKind` accepts every kind, validates every path
   before the first removal and keeps `docs/adr`, `docs/backlog`,
   `docs/findings` and the Spec Root. `ReduceHistoryEntry` and its reduce
   branch are retired.
2. MUST implement API Contract 8: a `superseded` unit after the kinds, holding
   the uncited `superseded` Archive Records without a regeneration grant, with
   citation measured against every tracked Markdown file that survives the
   plan, outside skill trees and `testdata`.
3. MUST implement API Contract 9 in both commands: apply refuses, by name, a
   path neither `history-full` nor a commit `ResolveGitLocator` finds holds;
   `baseline update` refuses that unit and plans the rest.
4. MUST implement API Contract 10: after the removals, both commands rewrite
   the citations of removed paths in Archive Records, Live Documents and
   `docs/user-guide/` with the Citation Rewrite of task_01, locating a Reduced
   History Entry at its provenance pair when that commit qualifies, and list
   every other citation as today. `HistoryCitations` reads Archive Records as
   sources. The rewrite reads only the repository and writes only the citing
   files the plan lists.
5. MUST implement API Contract 12 for the new units and lines, the JSON fields
   and the Plan Digest, keeping `applied.reduced` at 0, and MUST make
   `historyInventory` take a context, so `roundfix upgrade` names reduced
   entries, retired ADRs and uncited superseded records as Pending History.
6. MUST make `DiscoverHistoryLayout` stop producing relocations into the
   Findings, Backlog, ADR and Review families (`docs/findings/_archived/`,
   `_archived/{findings,adr,backlog}/`, retired `docs/adr/` files, terminal
   `docs/backlog/` entries and finished Review Artifacts), which the history
   section now removes in the same plan; Legacy Archive Folder relocations and
   the live Review Artifact root rename stay.
7. MUST implement Surface Transcripts 1 and 2 exactly, each asserted in one
   named test: every stdout and stderr line in order and the exit code, with
   only the transcript's angle-bracket variations.
8. MUST NOT write anything in plan mode (Invariant 5), and a second plan after
   an applied plan MUST report nothing pending for the applied units
   (Invariant 6).
9. MUST change these existing tests to the new contract and name each change
   in the Result: in `internal/spec/history_entries_test.go`, remove
   `TestReduceHistoryEntryKeepsFrontMatterTitleAndFirstParagraph`,
   `TestReduceHistoryEntryRefusesWhatItCannotReduce` and
   `TestReducedHistoryEntryIsNotPendingAgain`, and change
   `TestPlanHistoryKindsOrdersAndMeasures`,
   `TestApplyHistoryKindRemovesReviewsAndHandoffsOnly` and
   `TestHistoryCitationsNameMarkdownThatCitesARemovedPath`; in
   `internal/cli`, `TestHistorySanitizePlanWritesNothing`,
   `TestHistorySanitizeAppliesTheNextBatch`,
   `TestHistorySanitizeFolderCitationsAndKindNamedSlugs`,
   `TestHistorySanitizeRefusesAKindUnitAndPlansTheRest`,
   `TestBaselineUpdatePlansPendingHistoryWithoutWrites`,
   `TestBaselineUpdateAppliesPendingHistoryAndCreatesTheTag`,
   `TestBaselineUpdateHistoryPlanningBoundaries`,
   `TestBaselineUpdateHistoryCitationBoundaries` and
   `TestUpgradeNoticeIsSilentWithoutPendingHistory`; in `internal/baseline`,
   `TestDiscoverHistoryLayoutLegacyArchiveShapes`,
   `TestDiscoverHistoryLayoutClassifiesActiveDocuments`,
   `TestHistoryLayoutRelocatesLegacyReviewRootWhateverItsLiveness` and
   `TestHistoryMoveApply`. Any other test that changes is named with the
   contract it moves to.
10. MUST add the tests named in Verification, each calling `t.Parallel()`
    first, on temporary repositories and homes. In
    `internal/spec/history_removal_test.go`:
    `TestPlanHistoryKindsRemovesEveryRetiredFamily`,
    `TestPlanHistoryKindsRemovesWhatTheOldLayoutLeft`,
    `TestPlanHistoryKindsRefusesASupersededADRItsSuccessorDoesNotName`,
    `TestApplyHistoryKindKeepsTheActiveDirectories`,
    `TestPlanSupersededRecordsSkipsCitedAndGrantedRecords` and
    `TestRemovalRewriteLocatesAReducedEntryAtItsProvenance`. In
    `internal/cli/history_removal_test.go`:
    `TestHistorySanitizePlanTranscript`, `TestHistorySanitizeApplyTranscript`,
    `TestHistorySanitizeRefusesAPathNoReachableCommitHolds`,
    `TestBaselineUpdateRemovesRetiredHistoryAndRewritesCitations` (JSON
    actions, `applied.reduced` 0, the rewritten counts, a Plan Digest without
    `bytesAfter`, and a second run whose history is `current`),
    `TestBaselineUpdateRemovesWhatTheBaselineUsedToRelocate` and
    `TestUpgradeNoticeNamesReducedEntriesAsPending`. In
    `internal/baseline/history_layout_test.go`:
    `TestDiscoverHistoryLayoutRelocatesNothingIntoRetiredFamilies`.
11. MUST describe the removal, the coverage rule, the rewrite and the new
    lines in `docs/user-guide/commands/history.md` (with the phrase
    `Git Locator`, and without `reduces retired Findings`) and
    `docs/user-guide/commands/baseline.md` (with the phrase `Git Locator`, and
    without `It never rewrites citing files`), and in
    `docs/user-guide/commands/upgrade.md` that a Reduced History Entry counts
    as Pending History, with the phrase `Reduced History Entry`; then
    re-record the Behavior Surface Record with
    `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`.
12. MUST, because this Task removes and renames tests, regenerate the test
    inventory with `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`
    after its last test change, so `TestCoverageEquivalence` reports no
    regression; MUST replace every remaining caller of the retired reduction
    writer, including the legacy fixture in
    `internal/speccheck/reduced_history_entry_test.go`, with literal legacy
    bytes, and keep `TestReducedFindingsKeepTheirLicensesAndMembers` and
    `TestReducedFindingWithABrokenLicenseStillFails` asserting what they
    assert today.

## Subtasks

- [ ] Plan removal for every family and leftover, and retire the reduce path.
- [ ] Add the `superseded` unit and the coverage rule.
- [ ] Rewrite the citations of removed paths after apply.
- [ ] Print and publish the new units, lines and fields, and pass a context to the inventory.
- [ ] Stop relocating into retired families.
- [ ] Add the tests, update the changed ones and document the commands.

## Acceptance Criteria

- [ ] A plan lists every retired file, including reduced entries, retired ADRs and uncited superseded records, and writes nothing.
- [ ] Apply removes them, rewrites each citation in scope to its Git Locator and refuses a path no reachable commit holds.
- [ ] A second run after an apply reports the history current.
- [ ] Both transcripts match the built product line for line.

## Consumers

- `PlanHistoryKinds`, `ApplyHistoryKind` and `HistoryCitations` — consumers: `./internal/spec`, `./internal/cli`
- `baseline update` history JSON and Plan Digest — consumers: `./internal/cli`
- `DiscoverHistoryLayout` relocations — consumers: `./internal/baseline`, `./internal/cli`

## Context

- instruction: `docs/adr/0278-a-retirement-deletes-and-the-history-root-holds-only-archive-records.md`
- instruction: `docs/adr/0248-existing-history-is-sanitized-in-batches-after-a-history-full-tag.md`
- instruction: `docs/adr/0254-baseline-update-sanitizes-pending-history-in-the-change-it-plans.md`
- interface: `internal/spec/history_entries.go`
- interface: `internal/spec/history_entries_test.go`
- creates: `internal/spec/history_removal_test.go`
- interface: `internal/cli/history.go`
- interface: `internal/cli/history_test.go`
- interface: `internal/cli/history_refusal_test.go`
- interface: `internal/cli/baseline_update_history.go`
- interface: `internal/cli/baseline_update_history_test.go`
- interface: `internal/cli/baseline_update_citations_test.go`
- interface: `internal/cli/baseline_history_citation_test.go`
- interface: `internal/cli/upgrade.go`
- interface: `internal/cli/upgrade_notice_test.go`
- creates: `internal/cli/history_removal_test.go`
- interface: `internal/baseline/history_layout.go`
- interface: `internal/baseline/history_layout_test.go`
- interface: `internal/baseline/transaction_test.go`
- interface: `internal/baseline/apply_test.go`
- interface: `internal/baseline/history_citations_test.go`
- interface: `internal/baseline/history_citation_plan_test.go`
- interface: `docs/user-guide/commands/history.md`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `docs/user-guide/commands/upgrade.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `docs/references/coverage-record.json`
- interface: `internal/speccheck/reduced_history_entry_test.go`
- interface: `internal/baseline/reduced_retirement_clauses_test.go`
- interface: `internal/cli/baseline_update.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestPlanHistoryKindsRemovesEveryRetiredFamily|TestPlanHistoryKindsRemovesWhatTheOldLayoutLeft|TestPlanHistoryKindsRefusesASupersededADRItsSuccessorDoesNotName|TestApplyHistoryKindKeepsTheActiveDirectories|TestPlanSupersededRecordsSkipsCitedAndGrantedRecords|TestRemovalRewriteLocatesAReducedEntryAtItsProvenance)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestPlanHistoryKindsRemovesEveryRetiredFamily TestPlanHistoryKindsRemovesWhatTheOldLayoutLeft TestPlanHistoryKindsRefusesASupersededADRItsSuccessorDoesNotName TestApplyHistoryKindKeepsTheActiveDirectories TestPlanSupersededRecordsSkipsCitedAndGrantedRecords TestRemovalRewriteLocatesAReducedEntryAtItsProvenance; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run 'History|Retire|Reduced|Archive|Locator' ./internal/spec || exit 1; out="$(go test -count=1 -v -run '^TestDiscoverHistoryLayoutRelocatesNothingIntoRetiredFamilies$' ./internal/baseline 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- "--- PASS: TestDiscoverHistoryLayoutRelocatesNothingIntoRetiredFamilies " || { printf 'missing passing relocation test\n' >&2; exit 1; }; go test -count=1 -run 'History|Relocat|Transaction|Apply|Plan' ./internal/baseline || exit 1` — expected: exit 0. Before this Task the seven tests do not exist, so the PASS checks fail.
- `out="$(go test -count=1 -v -run '^(TestHistorySanitizePlanTranscript|TestHistorySanitizeApplyTranscript|TestHistorySanitizeRefusesAPathNoReachableCommitHolds|TestBaselineUpdateRemovesRetiredHistoryAndRewritesCitations|TestBaselineUpdateRemovesWhatTheBaselineUsedToRelocate|TestUpgradeNoticeNamesReducedEntriesAsPending)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestHistorySanitizePlanTranscript TestHistorySanitizeApplyTranscript TestHistorySanitizeRefusesAPathNoReachableCommitHolds TestBaselineUpdateRemovesRetiredHistoryAndRewritesCitations TestBaselineUpdateRemovesWhatTheBaselineUsedToRelocate TestUpgradeNoticeNamesReducedEntriesAsPending; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run 'History|BaselineUpdate|BaselinePlan|BaselineApply|Upgrade|Archive|Legacy|ReleasePlan' ./internal/cli || exit 1; for f in docs/user-guide/commands/history.md docs/user-guide/commands/baseline.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Git Locator' || { printf 'missing Git Locator in %s\n' "$f" >&2; exit 1; }; done; if tr -s '[:space:]' ' ' < docs/user-guide/commands/history.md | grep -qF -- 'reduces retired Findings'; then printf 'history guide still reduces\n' >&2; exit 1; fi; if tr -s '[:space:]' ' ' < docs/user-guide/commands/baseline.md | grep -qF -- 'It never rewrites citing files'; then printf 'baseline guide still refuses to rewrite\n' >&2; exit 1; fi; tr -s '[:space:]' ' ' < docs/user-guide/commands/upgrade.md | grep -qF -- 'Reduced History Entry' || { printf 'upgrade guide does not name reduced entries\n' >&2; exit 1; }; go test -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract || exit 1` — expected: exit 0. Before this Task the six command tests do not exist and the guides still describe reduction, so the command fails.

## References

- `_prd.md` → Goals; User Stories 1 and 3; Core Feature 1; Core Feature 4; Success Metric 5
- `_techspec.md` → API Contract 7; API Contract 8; API Contract 9; API Contract 10; API Contract 12; Surface Transcript 1; Surface Transcript 2; Invariants 1, 4, 5 and 6; Current behavior; Testing Approach; Build Order 2
- ADR-0278; ADR-0248; ADR-0254; ADR-0120; ADR-0258; ADR-0092
