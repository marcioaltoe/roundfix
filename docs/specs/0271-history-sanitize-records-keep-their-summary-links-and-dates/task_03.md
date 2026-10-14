---
task: task_03
spec: 0271-history-sanitize-records-keep-their-summary-links-and-dates
status: pending
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
