---
task: task_01
spec: 0272-the-history-holds-only-archive-records
status: pending
type: backend
complexity: high
---

# Task 01: A citation of a removed folder file becomes a Git Locator in the record and its copies

## Overview

The Archive Command copies the outcome and each promoted file byte for byte,
so a link into the removed Spec folder resolves to nothing, and the Legacy
Link Rebase points such a link at the record itself (`_techspec.md` → Current
behavior). This Task adds the Git Locator and the Citation Rewrite in
`internal/spec/git_locator.go`, makes the Archive Command and every Legacy
Archive Folder conversion rewrite the outcome and each promoted copy through
them, and teaches the Delivery Queue's exact-archive proof to recompute the
copies. It is verifiable on its own through `roundfix archive` on a temporary
repository and through the proof.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-13,
   [an Archive Record's outcome keeps links relative to the removed Spec folder](references/2026-10-13-an-archive-record-outcome-keeps-links-relative-to-the-removed-spec-folder.md),
   including its revision on the self-link, and the Backlog Entry of
   2026-10-14,
   [promoted copies keep prose paths to unpromoted folder files](references/2026-10-14-promoted-copies-keep-prose-paths-to-unpromoted-folder-files.md),
   whose open rule this Task settles as the Git Locator.
2. MUST implement API Contract 1 in `ResolveGitLocator`: the preferred commit
   when it holds the path and the default branch head or a tag reaches it,
   else the newest commit the default branch head reaches that touches the
   path (its first parent when that commit deleted it), else no commit. The
   default branch comes from `preflight.DetectDefaultBranch`, read from
   `refs/heads/<name>` or else `refs/remotes/origin/<name>`. Every Git command
   runs under the context, and a failure other than an absent object is
   returned wrapped with `%w`. Data comes only from the local repository; no
   fetch is made.
3. MUST implement API Contract 2 (`GitLocator.Text`) and API Contract 3
   (`CitationRewrite.Apply`), including the link, image and reference
   definition forms, the in-code form, the rebase of destinations outside the
   removed set when `from` differs from `to`, and byte preservation of every
   other byte. `RelinkContent` MUST delegate to `Apply` and return the bytes it
   returns at the delivery base for every input (Invariant 3).
4. MUST implement API Contract 4: `LegacyLinkRebase` gains `Locate` and
   `Rebase` returns an error; the Archive Command builds one with the folder's
   promotions and a resolver preferring the record's `source_revision`, and
   `BuildArchiveRecord` applies an input `Rebase` to the outcome whether or not
   the record is legacy. Each promoted copy is written through the same
   rebase. `PlanLegacyConversion` passes the same resolver, so a folder
   destination that is not promoted becomes a Git Locator, never the record.
   Live Documents keep the Archive Relink's moves unchanged.
5. MUST implement API Contract 5: the outcome shortening of
   `RenderArchiveRecord` never cuts inside a code span.
6. MUST implement API Contract 6 in `archiveRetirementIsExact`: a promoted
   copy is accepted when it equals the rewrite of the one removed file with its
   basename, recomputed with the parent preferred; every other check of the
   proof, including its refusal of a relinked path under the History Root,
   stays.
7. MUST NOT write a Git Locator that names a commit no default-branch head or
   tag reaches (Invariant 1), and MUST NOT write a link from an Archive Record
   to itself (Invariant 2).
8. MUST change these existing tests to the new contract and name the change in
   the Result: `TestLegacyLinkRebasePointsFolderLinksAtTheirSurvivingFile` and
   `TestLegacyConversionRebasesTheOutcomeAndPromotedCopies`
   (`internal/spec/legacy_link_rebase_test.go`: a folder destination without
   a promoted copy becomes a Git Locator) and the cases of
   `internal/cli/history_link_rebase_test.go`,
   `internal/cli/archive_relink_test.go` and
   `internal/cli/deliver_archive_relink_test.go` whose promoted copies or
   outcomes carry a folder link. Every other existing test MUST pass
   unchanged.
9. MUST add the tests named in Verification, each calling `t.Parallel()` first,
   on temporary repositories with `HOME` set to a temporary directory; a squash
   merge is simulated by a commit on a deleted branch. In
   `internal/spec/git_locator_test.go`:
   `TestResolveGitLocatorPrefersAReachableProvenanceCommit`,
   `TestResolveGitLocatorFallsBackToTheNewestDefaultBranchCommit`,
   `TestResolveGitLocatorNamesNoUnreachableCommit`,
   `TestResolveGitLocatorReadsTheRemoteDefaultBranch`,
   `TestGitLocatorTextForms`,
   `TestCitationRewriteReplacesLinksAndPathsWithGitLocators`,
   `TestCitationRewriteKeepsRelinkContentBytes`,
   `TestLegacyLinkRebaseNeverLinksTheRecordToItself` and
   `TestArchiveRecordShorteningKeepsGitLocatorsWhole`. In
   `internal/cli/archive_git_locator_test.go`:
   `TestArchiveOutcomeLinksBecomeGitLocators` (an outcome link to an adopted
   reference names the planning commit on the default branch, and one to a QA
   evidence file reads `(not reachable from any branch or tag)`),
   `TestArchivePromotedCopyKeepsItsTargetsThroughGitLocators` and
   `TestDeliveryAcceptsAnArchiveWithRewrittenPromotedCopies`.
10. MUST describe in `docs/user-guide/commands/archive.md` that the outcome
    and each promoted copy name a removed folder file by Git Locator, with the
    phrase `Git Locator`, and in `docs/user-guide/commands/deliver.md` that the
    exact-archive proof recomputes the promoted copies, with the phrase
    `Git Locator`, then re-record the Behavior Surface Record with
    `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`.

## Subtasks

- [ ] Add `GitLocator`, `ResolveGitLocator` and `CitationRewrite`.
- [ ] Make `RelinkContent` and `LegacyLinkRebase` delegate to the rewrite.
- [ ] Rewrite the Archive Command's outcome and promoted copies.
- [ ] Keep Git Locators whole when the outcome is shortened.
- [ ] Recompute promoted copies in the exact-archive proof.
- [ ] Add the tests and describe the archive and delivery behavior in the guide.

## Acceptance Criteria

- [ ] Every Git Locator written names a reachable commit that holds the path, or reads the pull request or `not reachable` form.
- [ ] No outcome or promoted copy links into the removed folder or to the record itself.
- [ ] `RelinkContent` output is unchanged for every existing Archive Relink case.
- [ ] The Delivery Queue accepts an archive whose promoted copy was rewritten and still refuses a tampered relink.

## Consumers

- `CitationRewrite`, `RelinkContent` and `LegacyLinkRebase.Rebase` — consumers: `./internal/spec`, `./internal/cli`
- Archive Record outcome and promoted copy bytes — consumers: `./internal/spec`, `./internal/cli`

## Context

- instruction: `docs/adr/0278-a-retirement-deletes-and-the-history-root-holds-only-archive-records.md`
- instruction: `docs/adr/0260-an-archive-relinks-live-citations-and-a-baseline-update-lists-the-citations-it-breaks.md`
- instruction: `docs/adr/0247-an-archive-leaves-an-archive-record-and-the-spec-folder-stays-in-git.md`
- instruction: `docs/adr/0277-a-legacy-archive-folder-conversion-keeps-the-summary-links-and-dates-its-folder-held.md`
- instruction: `internal/preflight/default_branch.go`
- creates: `internal/spec/git_locator.go`
- creates: `internal/spec/git_locator_test.go`
- creates: `internal/cli/archive_git_locator_test.go`
- interface: `internal/spec/live_citations.go`
- interface: `internal/spec/legacy_link_rebase.go`
- interface: `internal/spec/legacy_link_rebase_test.go`
- interface: `internal/spec/archive.go`
- interface: `internal/spec/archive_record.go`
- interface: `internal/spec/history_sanitize.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver_archive_relink_test.go`
- interface: `internal/cli/archive_relink_test.go`
- interface: `internal/cli/history_link_rebase_test.go`
- interface: `docs/user-guide/commands/archive.md`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestResolveGitLocatorPrefersAReachableProvenanceCommit|TestResolveGitLocatorFallsBackToTheNewestDefaultBranchCommit|TestResolveGitLocatorNamesNoUnreachableCommit|TestResolveGitLocatorReadsTheRemoteDefaultBranch|TestGitLocatorTextForms|TestCitationRewriteReplacesLinksAndPathsWithGitLocators|TestCitationRewriteKeepsRelinkContentBytes|TestLegacyLinkRebaseNeverLinksTheRecordToItself|TestArchiveRecordShorteningKeepsGitLocatorsWhole)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestResolveGitLocatorPrefersAReachableProvenanceCommit TestResolveGitLocatorFallsBackToTheNewestDefaultBranchCommit TestResolveGitLocatorNamesNoUnreachableCommit TestResolveGitLocatorReadsTheRemoteDefaultBranch TestGitLocatorTextForms TestCitationRewriteReplacesLinksAndPathsWithGitLocators TestCitationRewriteKeepsRelinkContentBytes TestLegacyLinkRebaseNeverLinksTheRecordToItself TestArchiveRecordShorteningKeepsGitLocatorsWhole; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run 'Archive|Legacy|Relink|Locator|History' ./internal/spec || exit 1` — expected: exit 0. Before this Task the nine tests do not exist, so the PASS check fails.
- `out="$(go test -count=1 -v -run '^(TestArchiveOutcomeLinksBecomeGitLocators|TestArchivePromotedCopyKeepsItsTargetsThroughGitLocators|TestDeliveryAcceptsAnArchiveWithRewrittenPromotedCopies)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestArchiveOutcomeLinksBecomeGitLocators TestArchivePromotedCopyKeepsItsTargetsThroughGitLocators TestDeliveryAcceptsAnArchiveWithRewrittenPromotedCopies; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -run 'Archive|Relink|Legacy|History|Delivery|Deliver|Cancel|Review.*Override|ConventionC5' ./internal/cli || exit 1; for f in docs/user-guide/commands/archive.md docs/user-guide/commands/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Git Locator' || { printf 'missing Git Locator in %s\n' "$f" >&2; exit 1; }; done; go test -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract || exit 1` — expected: exit 0. Before this Task the three command tests do not exist and neither guide names the Git Locator, so the command fails.

## References

- `_prd.md` → Goals; User Story 2; Core Feature 3; Success Metric 3
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; Invariants 1-4; Current behavior; Testing Approach; Build Order 1
- ADR-0278; ADR-0260; ADR-0247; ADR-0277; ADR-0261
