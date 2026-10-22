---
task: task_02
spec: 0295-the-archive-rewrites-the-satisfied-by-links-it-records
status: completed
type: backend
complexity: medium
---

# Task 02: The Record Link Repair and the strict check point a promoted target at its copy

## Overview

Records already written by an earlier archive keep `unproven` links that
resolve under `docs/history/specs/`. The Record Link Repair, the
`record-links` unit of `roundfix history sanitize` and of the history section
of `roundfix baseline update`, repairs a removed file's link with a Git
Locator, but points a file the archive promoted to `docs/references/` at the
same missing path or at its source, and `SC-HISTORY-RETAINED` suggests that
target. This Task makes both name the promoted copy, proves it through both
commands and the check, and states it in the `history`, `baseline` and `spec`
guides and the glossary.

## Requirements

1. MUST answer the Inbox Entry of 2026-10-20,
   [`references/2026-10-20-o-archive-copia-links-relativos-das-declaracoes-para-o-archive-record.md`](references/2026-10-20-o-archive-copia-links-relativos-das-declaracoes-para-o-archive-record.md)
   (`_prd.md` → Sources 1), by implementing API Contract 4 of `_techspec.md`
   in `internal/spec/record_links.go`: add `RecordLinkPromotion` with the
   signature `_techspec.md` → Interfaces sketches; in `PlanRecordLinks`, ask
   it first for each missing link of a record and, for a promotion, record a
   move, a `HistoryCitation` whose `Rewrite` is `docs/references/<file>`, and
   no locator; carry the moves per record in `RecordLinkPlan` beside the
   locators and apply both in `Apply`; give `rewriteRecordLinks` the moves
   parameter the sketch names and pass them to `CitationRewrite` as `Moves`,
   so a moved destination becomes the relative path from the record to the
   copy. `internal/spec/live_links.go` passes no moves. A link that is not a
   promotion keeps `ResolveRecordLink`'s locator.
2. MUST implement API Contract 5 in `internal/speccheck/history_retained.go`
   `detectHistoryRetained`: for a missing link in an Archive Record for which
   `RecordLinkPromotion` answers a copy, the fix line is `Replace the link
   with <relative path from the record to the copy>, the copy the archive
   promoted to docs/references/<file>.`, checked before Git; every other
   finding text, level, location and fix line is unchanged.
3. MUST state the hazard in the code's behavior: the record's `source`,
   `promoted` list and links were written by an earlier archive or sanitize;
   the repair writes into the record that later readers follow. A promoted
   copy is named only when the record lists it and the working tree holds it
   as a regular file (Invariant 4); nothing is fetched, and the `--apply` tag
   requirement is unchanged.
4. MUST keep everything else: `internal/cli/history.go`,
   `internal/cli/baseline_update_history.go`, the plan formats and the
   digest rules are not changed; every existing test name and assertion is
   unchanged, `TestPlanRecordLinksResolvesAgainstTheRecordThenItsSource` and
   the `live-links` tests included.
5. MUST add the tests `_techspec.md` → Testing Approach names, each top-level
   test calling `t.Parallel()` first, compiling for Windows, and using
   temporary repositories, homes and Git configuration:
   - in the new file `internal/spec/record_links_promoted_test.go`,
     `TestPlanRecordLinksPointsAPromotedTargetAtItsCopy`: a committed record
     shaped like the one `_techspec.md` → Current behavior step 1 shows, with
     `docs/references/promoted.md` in the tree; the plan's citation for the
     promoted link has `Rewrite` `docs/references/promoted.md`, the planned
     bytes hold `` [`references/promoted.md`](../../references/promoted.md) ``,
     the removed link keeps its Git Locator, planning writes nothing,
     `Apply` writes the planned bytes, and a second plan finds no link; and
     `TestRecordLinkPromotionRequiresTheListedCopy`, with subtests that answer
     no copy when the copy is absent from the tree, when the record does not
     list it under `promoted`, when the target lies outside the record's
     `source` folder, and for a self-link;
   - in the new file `internal/speccheck/history_retained_promoted_test.go`,
     `TestHistoryRetainedNamesThePromotedCopyForARecordLink`: the check on
     such a record reports the promoted link with the fix line of
     Requirement 2 and the removed link with its Git Locator fix;
   - in the new file `internal/cli/history_record_links_promoted_test.go`,
     `TestHistorySanitizePointsAPromotedRecordLinkAtItsCopyTranscript`, which
     asserts Surface Transcript 1 whole through `roundfix history sanitize`
     (every standard output line in order, an empty standard error and exit
     0, with `<commit>` and `<n>` the only variations) on the repository the
     transcript names, writing the record bytes the Delivery Base archive
     wrote; and `TestBaselineUpdatePlansAPromotedRecordLinkAtItsCopy`, which
     runs the update's history plan in JSON with the existing baseline update
     history fixtures and asserts the citation of the promoted link has
     `rewrite` `docs/references/promoted.md` and is bound by the Plan Digest.
6. MUST document the change (`_prd.md` → Core Feature 4):
   - in `docs/user-guide/commands/history.md` and
     `docs/user-guide/commands/baseline.md`, in the paragraph that begins
     "After removal units, `record-links` repairs", right after the sentence
     ending "holding the target wins.", add: "Before that, a link whose
     target, resolved against the `source` folder, is a file whose copy
     `docs/references/<file>` the record lists under `promoted` and the tree
     holds is rewritten to a relative link to that copy, and the plan names
     `docs/references/<file>` as its rewrite (ADR-0308).";
   - in `docs/user-guide/commands/spec.md`, right after the sentence ending
     "the deletion or replacement Git Locator.", add: "For a broken Archive
     Record link to a file the record promoted, the fix names the relative
     link to its `docs/references/` copy instead (ADR-0308).";
   - MUST re-record the Behavior Surface Record, whose three guide
     fingerprints change, by running
     `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`.
     The Agent runs it itself before handing back; the record is never
     hand-edited, and `docs/references/skill-coverage.json` does not change.
7. MUST add the glossary term **Record Link Repair** to `CONTEXT.md` through
   `domain-modeling` (`_prd.md` → Glossary), right after the **History
   Sanitize Command** entry, with this definition: "The `record-links` unit
   of the History Sanitize Command and of the history section of `roundfix
   baseline update`, which rewrites each broken relative link and self-link
   of a surviving Archive Record: a target the record promoted points at its
   `docs/references/` copy, and any other becomes a Git Locator tried against
   the record's directory, then its `source` folder" followed by a link to
   ADR-0308 relative to `CONTEXT.md`, and the line "_Avoid_: Live-links
   repair, record rewrite, Archive Relink".

## Subtasks

- [ ] Add `RecordLinkPromotion` and carry moves through the repair.
- [ ] Name the promoted copy in the strict check's fix line.
- [ ] Add the package, check and command tests.
- [ ] Update the three guides, re-record the Behavior Surface Record, and add the glossary term.

## Acceptance Criteria

- [ ] The repair rewrites a promoted target to its copy and names `docs/references/<file>` in both plans, and keeps every other target's locator.
- [ ] The strict check's fix line names the same copy.
- [ ] The guides and the glossary state the repair, and the Behavior Surface Record is current.

## Consumers

- `PlanRecordLinks`, `RecordLinkPlan` and the `record-links` plan citations — consumers: `./internal/spec`, `./internal/cli`
- The `SC-HISTORY-RETAINED` fix line — consumers: `./internal/speccheck`

## Context

- instruction: `docs/adr/0308-an-archive-rebases-the-unproven-links-it-records-and-the-record-link-repair-points-a-promoted-target-at-its-copy.md`
- instruction: `docs/adr/0299-the-history-conversion-rewrites-every-citation-it-lists-and-the-check-reports-a-removed-history-file-cited-without-a-git-locator.md`
- instruction: `docs/adr/0248-existing-history-is-sanitized-in-batches-after-a-history-full-tag.md`
- instruction: `internal/spec/git_locator.go`
- instruction: `internal/spec/history_entries.go`
- instruction: `internal/spec/record_links_test.go`
- instruction: `internal/speccheck/history_retained_test.go`
- instruction: `internal/cli/history.go`
- instruction: `internal/cli/history_test.go`
- instruction: `internal/cli/baseline_update_history.go`
- instruction: `internal/cli/baseline_update_citations_test.go`
- instruction: `internal/cli/baseline_update_history_test.go`
- instruction: `internal/docscontract/skill_coverage_test.go`
- creates: `internal/spec/record_links_promoted_test.go`
- creates: `internal/speccheck/history_retained_promoted_test.go`
- creates: `internal/cli/history_record_links_promoted_test.go`
- interface: `internal/spec/record_links.go`
- interface: `internal/spec/live_links.go`
- interface: `internal/speccheck/history_retained.go`
- interface: `docs/user-guide/commands/history.md`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `docs/user-guide/commands/spec.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestPlanRecordLinksPointsAPromotedTargetAtItsCopy|TestRecordLinkPromotionRequiresTheListedCopy|TestPlanRecordLinksResolvesAgainstTheRecordThenItsSource|TestMissingRelativeLinksFindsBrokenAndSelfLinks)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestPlanRecordLinksPointsAPromotedTargetAtItsCopy TestRecordLinkPromotionRequiresTheListedCopy TestPlanRecordLinksResolvesAgainstTheRecordThenItsSource; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^(TestHistoryRetainedNamesThePromotedCopyForARecordLink|TestHistoryRetainedReportsBrokenAndSelfLinksWithTheirGitLocator|TestHistoryRetainedPrintsTheLocatorOfTheFileABrokenLinkNames)$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestHistoryRetainedNamesThePromotedCopyForARecordLink TestHistoryRetainedReportsBrokenAndSelfLinksWithTheirGitLocator; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^(TestHistorySanitizePointsAPromotedRecordLinkAtItsCopyTranscript|TestBaselineUpdatePlansAPromotedRecordLinkAtItsCopy|TestHistorySanitizePlanWritesNothing|TestHistorySanitizeAppliesTheNextBatch|TestBaselineUpdateRewritesEveryListedCitation|TestBaselineUpdateHistoryDigestBindsCitations)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestHistorySanitizePointsAPromotedRecordLinkAtItsCopyTranscript TestBaselineUpdatePlansAPromotedRecordLinkAtItsCopy TestHistorySanitizeAppliesTheNextBatch; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for f in docs/user-guide/commands/history.md docs/user-guide/commands/baseline.md; do tr -s '[:space:]' ' ' < "$f" | grep -q -- 'the record lists under .promoted. and the tree holds is rewritten to a relative link to that copy' || { printf 'missing repair sentence in %s\n' "$f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/user-guide/commands/spec.md | grep -qF -- 'For a broken Archive Record link to a file the record promoted, the fix names the relative link' || { echo 'missing spec guide sentence' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Record Link Repair**: The' || { echo 'missing Record Link Repair term' >&2; exit 1; }; go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' || exit 1; GOOS=windows go vet ./internal/spec ./internal/speccheck ./internal/cli` — expected: exit 0. Before this Task the new tests do not exist, so the command fails at its first pass check. Afterwards the repair, both plans and the check name the promoted copy, the existing record-links, history, baseline citation and digest tests pass, the guides and the glossary carry their new text, the Behavior Surface Record is current, and the three packages' tests compile for Windows.

## References

- `_prd.md` → Sources; Goals; Core Features 2-4; Behavior changes; Success Metric 2; Success Metric 3; Glossary
- `_techspec.md` → Current behavior; Interfaces; API Contracts 4-5; Surface Transcript 1; Invariants 4-5; Testing Approach; Build Order 2
- ADR-0308; ADR-0248; ADR-0254; ADR-0299; ADR-0261; ADR-0184; ADR-0280; ADR-0259; ADR-0271; ADR-0285


## Result

Implemented this Task's slice; Task status and declared Verification remain
Daemon-owned.

- Repair and plans: `RecordLinkPromotion` resolves the historical destination
  against the record's `source`, requires containment under that folder, a
  matching `promoted` entry and a regular checkout copy (including no symlink
  components), and excludes self-links. `PlanRecordLinks` records per-record
  moves beside locators, names the copy's repository path in the citation,
  and applies both through the existing Citation Rewrite. Other links still
  use `ResolveRecordLink`; live-links passes no moves. The command handlers,
  plan formats and digest rules were not edited, and nothing is fetched.
- Strict check: `SC-HISTORY-RETAINED` names the relative promoted copy before
  Git Locator resolution, including outside Git; the existing finding text,
  severity, location and non-promotion fix remain covered by regression tests.
- Documentation: added the prescribed sentences to the history, baseline
  and spec guides and the Record Link Repair term after History Sanitize
  Command in `CONTEXT.md`, linking ADR-0308. Re-recorded the Behavior Surface
  Record with its owning generator; only those three guide fingerprints
  changed. `docs/references/skill-coverage.json` is unchanged.

Acceptance evidence:

| Criterion | Focused evidence |
| --- | --- |
| Promoted-copy repair and both plans, preserving other locators | `TestPlanRecordLinksPointsAPromotedTargetAtItsCopy` checks citation rewrites, relative planned bytes, planning without writes, Apply's exact bytes and an empty second plan. `TestRecordLinkPromotionRequiresTheListedCopy` covers absent/unlisted/directory copies, outside-source targets and self-links. Both new command tests passed: sanitize asserts the entire Surface Transcript 1, empty stderr and exit 0 using historical record bytes; baseline asserts JSON citations and proves changing only the promoted rewrite changes the digest. Existing record-link, live-link, history apply and baseline citation/digest regressions passed unchanged. |
| Strict fix names the same copy | `TestHistoryRetainedNamesThePromotedCopyForARecordLink` passed inside and outside Git, checking exact fixes and unchanged finding text, severity and location; the removed link retains its locator inside Git. Existing History Retained regressions passed. |
| Guides, glossary and current Behavior Surface Record | Required text is present, and the owning record command exited 0 after the final guide edits. Diff inspection confirms exactly the three expected guide fingerprints changed. |

Focused checks (all exited 0):

- `go test ./internal/spec ./internal/speccheck ./internal/cli -run 'Test(PlanRecordLinks|RecordLinkPromotion|MissingRelativeLinks|LiveLinks|HistoryRetained|HistorySanitize(RewritesBroken|PointsAPromoted|PlanWritesNothing|AppliesTheNextBatch)|BaselineUpdate(RewritesBroken|PlansAPromoted|RewritesEveryListedCitation|HistoryDigestBindsCitations))' -count=1`
- `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1` (required generation, not the Daemon's read-only Verification command).
- `GOOS=windows go test -c -o <temporary-directory>/spec.test.exe ./internal/spec`, then the equivalent compile commands for `./internal/speccheck` and `./internal/cli`; all three test executables compiled into a disposable directory under `/tmp`.
- `git diff --check`.

The initial focused spec-test attempt failed to compile because
`RecordLinkPromotion` did not yet exist. After implementation, the new
behavior tests and existing focused regressions passed. The sandbox rejected
the symlinked writable root before command execution; authorized escalated
execution accessed the assigned checkout without changing sandbox settings.

Not run: the authored `## Verification` command, repository Verification,
incremental Verification or the full suite; these remain for Daemon
settlement. No Task status, Task Graph manifest, other Task file, existing
test assertion or skill was edited; no commit, push or pull request was made.
No follow-up outside this slice was identified.
