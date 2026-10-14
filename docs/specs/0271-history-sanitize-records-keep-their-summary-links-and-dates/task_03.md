---
task: task_03
spec: 0271-history-sanitize-records-keep-their-summary-links-and-dates
status: completed
type: backend
complexity: high
---

# Task 03: Legacy Link Rebase points the outcome and promoted copies at surviving files

## Overview

A conversion moves two kinds of text out of the Legacy Archive Folder: the
outcome into the record one level up, and each promoted file into
`docs/references/`. Both keep destinations written relative to the folder,
so `references/<file>.md` in a summary and `../_prd.md` in a copy stop
resolving (`_techspec.md` → Current behavior). This Task adds the Legacy Link
Rebase and applies it to both in `PlanLegacyConversion`. It is verifiable on
its own through `PlanLegacyConversion` and through both commands on temporary
repositories.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-13,
   [history sanitize records lose their summary, links and dates](references/2026-10-13-history-sanitize-records-lose-their-summary-links-and-dates.md),
   for its fourth defect, both the summary links and the links inside
   promoted references.
2. MUST add `LegacyLinkRebase` and its `Rebase` method in the new
   `internal/spec/legacy_link_rebase.go` following `_techspec.md` →
   Interfaces and API Contract 5, reusing `scanArchiveMarkdown`,
   `archiveLinkTarget`, `ArchiveCitationMoves` and the text pass of
   `RelinkContent` rather than a second scanner. A destination inside the
   folder resolves to its promoted copy or else to the record; any other
   relative destination keeps naming the same repository path; a destination
   leaving the repository, a scheme, an anchor alone and code are untouched.
3. MUST apply it in `PlanLegacyConversion`: to each promoted copy, from its
   source path to its destination, and to the outcome, through
   `ArchiveRecordInput.Rebase`, from the PRD or QA Report that supplied it to
   the record path. The plan and the apply MUST write the same bytes
   (Invariant 6), and `ReadArchivedSpec`'s folder branch, which writes no
   record, passes no rebase.
4. MUST keep every other byte (Invariants 2 and 3): a promoted file without
   relative destinations or folder path text stays byte-identical, which
   `TestBaselineUpdateHistoryPromotionsRepeatAndRefuseSharedBasenames` keeps
   asserting unchanged.
5. MUST NOT change `RelinkContent`, `PlanArchiveRelink` or the Archive
   Command (API Contract 6).
6. MUST add the tests named in Verification, each calling `t.Parallel()`
   first. In the new `internal/spec/legacy_link_rebase_test.go`:
   `TestLegacyLinkRebasePointsFolderLinksAtTheirSurvivingFile` (a promoted
   target, an unpromoted folder file and a core artifact),
   `TestLegacyLinkRebaseKeepsOutsideTargets` (an ADR outside the folder
   before and after the move), `TestLegacyLinkRebaseLeavesOtherBytes` (no
   links, an absolute URL, an anchor, a code span and a fenced block keep
   their bytes) and `TestLegacyConversionRebasesTheOutcomeAndPromotedCopies`
   (asserting `RenderArchiveRecord` and the copy content). In the new
   `internal/cli/history_link_rebase_test.go`:
   `TestHistorySanitizeRebasesPromotedLinks` and
   `TestBaselineUpdatePromotionRebasesLinks`, each resolving every rewritten
   destination from its written file to an existing path. Every existing test
   MUST pass unchanged.

## Subtasks

- [ ] Add `LegacyLinkRebase` on the archive scanner and path pass.
- [ ] Rebase the promoted copies in `PlanLegacyConversion`.
- [ ] Rebase the outcome through `ArchiveRecordInput.Rebase`.
- [ ] Add the six tests.

## Acceptance Criteria

- [ ] A summary link to a promoted reference resolves from the record to the promoted copy.
- [ ] A link in a promoted copy to a core artifact resolves to the record, and one to a file outside the folder still resolves to it.
- [ ] A copy or record without such links keeps its bytes.
- [ ] `history sanitize --apply` and `baseline update --yes --promote` write the same rebased bytes their plans report.

## Consumers

- Promoted copy bytes and Archive Record outcome of a conversion — consumers: `./internal/spec`, `./internal/cli`

## Context

- instruction: `docs/adr/0277-a-legacy-archive-folder-conversion-keeps-the-summary-links-and-dates-its-folder-held.md`
- instruction: `docs/adr/0260-an-archive-relinks-live-citations-and-a-baseline-update-lists-the-citations-it-breaks.md`
- instruction: `internal/spec/live_citations.go`
- creates: `internal/spec/legacy_link_rebase.go`
- creates: `internal/spec/legacy_link_rebase_test.go`
- creates: `internal/cli/history_link_rebase_test.go`
- interface: `internal/spec/history_sanitize.go`
- interface: `internal/spec/archive_record.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestLegacyLinkRebasePointsFolderLinksAtTheirSurvivingFile|TestLegacyLinkRebaseKeepsOutsideTargets|TestLegacyLinkRebaseLeavesOtherBytes|TestLegacyConversionRebasesTheOutcomeAndPromotedCopies)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestLegacyLinkRebasePointsFolderLinksAtTheirSurvivingFile TestLegacyLinkRebaseKeepsOutsideTargets TestLegacyLinkRebaseLeavesOtherBytes TestLegacyConversionRebasesTheOutcomeAndPromotedCopies; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run 'Legacy|Archive|History|Relink' ./internal/spec || exit 1` — expected: exit 0. Before this Task the four tests do not exist, so the PASS check fails.
- `out="$(go test -count=1 -v -run '^(TestHistorySanitizeRebasesPromotedLinks|TestBaselineUpdatePromotionRebasesLinks)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestHistorySanitizeRebasesPromotedLinks TestBaselineUpdatePromotionRebasesLinks; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run 'History|Legacy|Archive|BaselineUpdate' ./internal/cli || exit 1` — expected: exit 0. Before this Task the two command tests do not exist, so the PASS check fails; afterwards they, the archive and Archive Relink tests and the baseline update history tests pass.

## References

- `_prd.md` → Goals; User Story 2; Core Feature 3; Success Metric 1; Success Metric 2
- `_techspec.md` → API Contract 5; API Contract 6; Interfaces; Invariants 2, 3, 4 and 6; Current behavior; Testing Approach; Build Order 3
- ADR-0277; ADR-0260

## Result

Implemented the Task 03 slice for the adopted Backlog Entry's fourth defect.
`LegacyLinkRebase` reuses the Archive Markdown scanner, destination resolver,
citation moves and bounded path-text pass. It rebases relative destinations
from their source file to their written file, retaining query and fragment
suffixes. Targets inside the removed folder name their promoted copy or the
record; outside targets retain their repository path. URLs, schemes, anchors,
repository escapes and relative destinations inside code retain their bytes.
The shared path-text pass keeps its existing treatment of code.

`PlanLegacyConversion` collects all promotion destinations before rebasing
copies, so links between promoted files resolve regardless of promotion order.
The planned copy bytes are the bytes `ApplyLegacyConversion` writes.
`ArchiveRecordInput.Rebase` applies only to legacy outcomes, using the PRD,
selected QA Report or supersession artifact that supplied the outcome. The
folder reader passes no rebase. `RelinkContent`, `PlanArchiveRelink`, the
Archive Command and existing tests were not edited.

Added all six required tests with `t.Parallel()` first. Acceptance evidence:

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| Summary links resolve from the record to promoted references | `TestLegacyLinkRebasePointsFolderLinksAtTheirSurvivingFile` and `TestLegacyConversionRebasesTheOutcomeAndPromotedCopies` exercise promoted targets, PRD summaries, first-section summaries and QA outcomes. Both command tests check the written summary and resolve its destinations to existing files. |
| Copy links to core artifacts resolve to the record; outside links keep their target | The rebase tests cover core, unpromoted and unlisted folder targets, folder links, encoded paths, reference definitions and outside ADRs. Both command tests resolve every written destination, including core, unpromoted, peer-promotion and ADR links with suffixes. |
| Copies and records without such links keep their bytes | `TestLegacyLinkRebaseLeavesOtherBytes` checks exact bytes, including binary, CRLF, URL, scheme, anchor, code and repository-escape cases. The conversion test compares a link-free record with rendering without rebase and compares link-free promoted bytes. The existing `TestBaselineUpdateHistoryPromotionsRepeatAndRefuseSharedBasenames` passes unchanged. |
| Both commands write their planned rebased bytes | The conversion test compares rendering, planned copies and applied files. `TestHistorySanitizeRebasesPromotedLinks` and `TestBaselineUpdatePromotionRebasesLinks` compare written records to `PlanLegacyConversion.Rendered`, assert the exact expected copy bytes and resolve destinations after folder removal. |

Focused checks run in this turn:

- Before implementation: `go test ./internal/cli -run 'RebasesPromotedLinks|PromotionRebasesLinks' -count=1` exited 1. Both new command tests reproduced the unchanged summary destinations from the removed folder.
- After implementation: `go test ./internal/spec -run 'LegacyLinkRebase|LegacyConversionRebases' -count=1` exited 0.
- `go test ./internal/cli -run 'RebasesPromotedLinks|PromotionRebasesLinks|TestBaselineUpdateHistoryPromotionsRepeatAndRefuseSharedBasenames' -count=1` exited 0.
- `go test ./internal/spec -run '^Test(RelinkContentRewritesOnlyDestinationsAndPaths|PlanArchiveRelinkPromotesLiveCitedFilesAndTargetsTheRecord|ArchiveWritesTheArchiveRecordAndRemovesTheSpecFolder|LegacyConversionWritesARecordGitStillHolds|LegacyConversionPreservesTheRendererOutcomeBudget|ActiveArchiveRecordKeepsItsSummaryReading|ActiveArchiveRecordKeepsItsFields)$' -count=1` exited 0.
- `go test ./internal/cli -run '^Test(ArchiveRelinksLiveCitationsTranscript|ArchiveWithoutInboundCitationsIsUnchanged|ArchiveRelinkCombinesManualPromotionsAndPreservesHistory)$' -count=1` exited 0.
- `gofmt` applied to the five implementation/test files; `git diff --check` exited 0.

The declared Verification commands and repository gates were not run; the
Daemon owns them and Task settlement. The pre-existing status edit remains
`in_progress`. No other Task file or manifest was edited, and no commit, push
or pull request was made. No follow-up work was identified in this slice.

## Carry-forward provenance

- Source Run: `run_20261014T193608Z_0f6bbae56906b671`
- Source commit: `39045f444e2df3ca717c414b052bb91c517ea86f`
