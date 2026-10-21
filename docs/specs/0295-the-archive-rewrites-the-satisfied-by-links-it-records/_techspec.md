---
spec: 0295-the-archive-rewrites-the-satisfied-by-links-it-records
prd: _prd.md
created: 2026-10-21
---

# The archive rewrites the satisfied-by links it records as unproven — Technical Spec

## Executive Summary

`BuildArchiveRecord` in `internal/spec/archive_record.go` already applies the
Legacy Link Rebase to the record's outcome paragraph, but copies the
`unproven` entries, the `satisfied-by` text of the Spec's Unreachable
Acceptance declarations, unchanged. A relative link in an entry then resolves
against `docs/history/specs/` instead of the Spec folder and breaks. The
builder now applies the same rebase to each entry, read as text of the
folder's `_prd.md`, so every caller (the Archive Command, the in-Run and
Delivery Queue archives through `spec.Archive`, and the Legacy Archive Folder
conversion) gets the fix from one code path.

For records already written, the Record Link Repair (`PlanRecordLinks` in
`internal/spec/record_links.go`, the `record-links` unit of `roundfix history
sanitize` and of `roundfix baseline update`) first asks whether a broken
link's source-folder target was promoted: when the record lists
`docs/references/<file name>` under `promoted` and the tree holds it, the link
is rewritten to that copy through the Citation Rewrite's existing moves.
`SC-HISTORY-RETAINED` uses the same answer in its fix line.

The trade-off accepted: the repair matches a promoted copy by file name,
because a record keeps only the copy's path, in exchange for repairing every
record already written without reading the removed folder.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted; only the text of `unproven` entries changes. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; every test uses temporary repositories
  and homes. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — Build Order 1 follows ADR-0308,
  ADR-0247, ADR-0277, ADR-0278, ADR-0230, ADR-0260, ADR-0261, ADR-0251 and
  ADR-0240; Build Order 2 follows ADR-0308, ADR-0248, ADR-0254, ADR-0299 and
  ADR-0261; Build Order 3 follows ADR-0308 and ADR-0256; every step follows
  ADR-0259, ADR-0271, ADR-0285, ADR-0104, ADR-0276, ADR-0265, ADR-0183,
  ADR-0184 and ADR-0280, each named with its title in `_prd.md` → Project
  Constraints. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AN approval of 2026-10-21, authorization "the same rules as cycle AM",
  carrying the cycle AJ grant "Mexer em Governed Paths declarados" and the
  2026-09-30 skills grant. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0295-the-archive-rewrites-the-satisfied-by-links-it-records/_authorization.md`;
  bounded files: `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`.

## Current behavior

Characterized on 2026-10-21 at `956066a7` with `bin/roundfix` built by `make
build` from that commit, in temporary repositories under a temporary `HOME`
and `GIT_CONFIG_GLOBAL`. The Characterization Repository has a `main`
branch and holds:

- `docs/agents/agent-instructions.md`, `docs/agents/domain.md`, and
  `docs/agents/docs-layout.md` with a `guide.docs-layout` setup region
  carrying "A retirement deletes the file", recorded in
  `docs/agents/setup-context.json` as the `guide.docs-layout` managed
  artifact, so `SC-HISTORY-RETAINED` runs;
- the Spec `0001-demo`: a `_prd.md` with Project Constraints, Success
  Metrics and this section,

  ```markdown
  ## Unreachable Acceptance

  - criterion: the seven-day measurement
    reason: it needs seven days of production traffic
    satisfied-by: the operator records it in [`references/promoted.md`](references/promoted.md)
  - criterion: the alert record
    reason: it needs a live cluster
    satisfied-by: the operator follows [`references/removed.md`](references/removed.md)
  ```

  an approved `_authorization.md`, a `_tasks.md` with one completed
  `task_01.md`, `references/removed.md`, and
  `qa/qa-report-2026-10-21.md` with `verdict: partial`,
  `rows_blocked_declared: 2`, zero finding and environment rows, and an
  `## Outcome` paragraph;
- a second active Spec `0002-other` with one pending Task, used only as the
  strict check's subject.

That `main` commit is pushed to a bare `origin` with `origin/HEAD` set.
Then a branch `docs/0001-demo` adds `references/promoted.md` in one commit,
the delivery branch shape of the Source, and from it:

1. `bin/roundfix archive 0001-demo --promote references/promoted.md
   --pull-request 91` exited 0, warned that `source_revision` is not an
   ancestor of `origin/main`, and wrote a record whose list read:

   ```yaml
   unproven:
     - the operator records it in [`references/promoted.md`](references/promoted.md)
     - the operator follows [`references/removed.md`](references/removed.md)
   ```

   with `promoted: [docs/references/promoted.md]` and `pull_request: "91"`.
   The outcome paragraph, which had no link, was unchanged.
2. After committing the archive, `bin/roundfix spec check --strict
   0002-other` exited 1 with two errors, at lines 15 and 16 of the record:
   `SC-HISTORY-RETAINED: docs/history/specs/0001-demo.md links to a missing
   file or itself: docs/history/specs/references/promoted.md`, fix ``Replace
   the link with `docs/history/specs/references/promoted.md` (pull request
   #91).``, and the same for `removed.md` with fix ``Replace the link with
   `docs/specs/0001-demo/references/removed.md@<main commit>`.``
3. `bin/roundfix history sanitize` exited 0 with `record-links: rewrites 2
   link(s) in 1 Archive Record(s)` and `cites
   docs/history/specs/0001-demo.md:15 names
   docs/history/specs/references/promoted.md; rewrites to
   `docs/history/specs/references/promoted.md` (pull request #91)`, the
   Source's wrong target, and a correct Git Locator for `removed.md`.

With every file on `main` and the archive run there, step 1 wrote the same
list, step 2 suggested `@<commit>` locators for both, and step 3 planned the
promoted file's Git Locator at the source revision rather than its
`docs/references/` copy.

Where the text is copied without a rebase, at `956066a7`:

- `internal/spec/archive.go` `archiveUnprovenActions` returns each
  declaration's `SatisfiedBy` verbatim, and `Archive` passes it as
  `ArchiveRecordInput.Unproven`.
- `internal/spec/archive_record.go` `BuildArchiveRecord` (lines 290-301)
  applies `in.Rebase` only to `r.Outcome`; `r.Unproven` comes from the input
  or, for a Legacy Archive Folder, from the PRD front matter.
- `internal/spec/record_links.go` `PlanRecordLinks` resolves each broken
  record link through `ResolveRecordLink`, which knows only Git Locators,
  and `rewriteRecordLinks` passes only locators to `CitationRewrite`.
- `internal/speccheck/history_retained.go` `detectHistoryRetained` prints
  `ResolveRecordLink`'s locator as the fix of a record link.

`BuildArchiveRecord` has two callers, `spec.Archive` and the Legacy Archive
Folder conversion in `internal/spec/history_sanitize.go`; `PlanRecordLinks`
is called by the History Sanitize Command and, through `planLinksUnit`, by
the history section of `roundfix baseline update`. The Delivery Queue's
exact-archive proof (`archiveRetirementIsExact`) and Pull Request Record
Commit proof do not compare `unproven`. Spec 0293's containment logic in
`RepositoryPath` and the promotion and returned-source destinations is not
on any of these paths.

## System Architecture

- `internal/spec/archive_record.go`: `BuildArchiveRecord` rebases each
  `unproven` entry with `in.Rebase` after the outcome.
- `internal/spec/record_links.go`: new `RecordLinkPromotion`;
  `PlanRecordLinks`, `RecordLinkPlan.Apply` and `rewriteRecordLinks` carry
  per-record moves beside the locators.
- `internal/spec/live_links.go`: passes no moves to `rewriteRecordLinks`.
- `internal/speccheck/history_retained.go`: a record link with a promotion
  gets the promoted copy in its fix line.
- Guides, skill references and `CONTEXT.md` state the rewrites.
- No change to `internal/spec/archive.go`, `internal/cli/archive.go`,
  `internal/cli/history.go`, `internal/cli/baseline_update_history.go`, the
  record schema or the Delivery Queue proofs.

## Implementation Design

### Interfaces

```go
// RecordLinkPromotion names the docs/references/ copy an Archive Record
// promoted for a missing record link. The link's destination is resolved
// against record.Source; it qualifies when that target lies under the
// Source folder, "docs/references/"+base(target) is listed in
// record.Promoted, and the working tree holds that path as a regular file.
func RecordLinkPromotion(root string, link MissingLink, record ArchiveRecord) (string, bool)

// rewriteRecordLinks rewrites only link spans whose record-relative target
// has a locator or a move; a move becomes the relative path to the copy.
func rewriteRecordLinks(content []byte, path string, locators map[string]GitLocator, moves map[string]string) ([]byte, int, error)
```

### Data Models

None. The Archive Record schema and fields are unchanged. `RecordLinkPlan`
gains an unexported per-record map of moves beside its locators.

### API Contracts

1. API Contract: `BuildArchiveRecord` with a non-nil `Rebase`. After the
   outcome, each `unproven` entry is passed to `Rebase.Rebase(entry,
   <Rebase.Folder>/_prd.md, Rebase.RecordPath)`, the same call the outcome
   uses with its own source file. The result replaces the entry in a new
   slice; the caller's slice is never written. Entry order and count are
   kept; an entry without a relative link keeps its bytes; an error is
   returned as `rebase archive unproven entry <n>: <error>` before anything
   is written. A nil `Rebase` keeps today's behavior.
2. API Contract: the Archive Command and every archive through
   `spec.Archive`. A folder destination in an entry points at its promoted
   copy when the archive promotes it (`[label](../../references/<file>)`
   for the built-in Spec Root) and otherwise becomes `label
   (`<folder path>@<commit>`)`, or the locator with `(pull request #<n>)` or
   `(not reachable from any branch or tag)` exactly as the outcome's
   locators read today. A destination outside the folder keeps naming the
   same repository path. Output, exit codes and refusals are unchanged.
3. API Contract: the Legacy Archive Folder conversion. The same rebase
   applies to the legacy `unproven` lines, after Lenient Legacy Reading has
   turned legacy maps into text.
4. API Contract: `PlanRecordLinks` and the `record-links` unit. For each
   broken link of a record, `RecordLinkPromotion` is asked first. A
   promotion records a move, a `HistoryCitation` whose `Rewrite` is the
   repository path `docs/references/<file>`, and rewrites the link
   destination to the relative path from the record to the copy; every other
   link keeps `ResolveRecordLink`'s locator. `Apply` re-reads each record and
   applies both. The text and JSON plan of `roundfix history sanitize` and of
   `roundfix baseline update` print that `Rewrite`, and the update's Plan
   Digest binds it as it binds every citation. The `--apply` tag requirement
   is unchanged.
5. API Contract: `SC-HISTORY-RETAINED`. For a broken link in an Archive
   Record with a promotion, the fix line is `Replace the link with
   <relative path>, the copy the archive promoted to docs/references/<file>.`;
   the finding text, level, location and every other fix line are unchanged.
   The promotion is checked before Git, so the fix also names the copy
   outside Git.

### Surface Transcripts

1. Surface Transcript: `roundfix history sanitize` on the Characterization
   Repository in its main shape (every file committed on `main`, the archive
   run there) after the Delivery Base binary wrote and the operator committed
   its archive.

   ```transcript
   $ roundfix history sanitize
   stdout:
   history sanitize plan: 1 unit(s) pending; 0 file(s) (0 bytes) leave docs/history
   record-links: rewrites 2 link(s) in 1 Archive Record(s)
   cites docs/history/specs/0001-demo.md:15 names docs/history/specs/references/promoted.md; rewrites to docs/references/promoted.md
   cites docs/history/specs/0001-demo.md:16 names docs/history/specs/references/removed.md; rewrites to `docs/specs/0001-demo/references/removed.md@<commit>`
   apply with: roundfix history sanitize --apply --batch <n> (needs the annotated tag history-full at or before HEAD)
   stderr:
   exit: 0
   ```

### Invariants

```text
1. A record the archive writes holds no relative link from an unproven entry to a file under the removed folder.
2. An unproven entry without a relative link, and the order and count of entries, keep their bytes.
3. The unproven entries and the outcome are rewritten by the same LegacyLinkRebase value.
4. The Record Link Repair names a promoted copy only when the record lists it and the tree holds it.
5. The repair's plan and the strict check's fix line name the same promoted copy.
```

## Coverage Map

- Goal 1, Core Feature 1, Success Metric 1 → `BuildArchiveRecord`
  (API Contracts 1-3), Build Order 1.
- Goal 2 → API Contract 2, Invariants 1 and 3.
- Goal 3, Core Feature 2, Success Metric 2 → `RecordLinkPromotion`,
  `PlanRecordLinks` (API Contract 4), Surface Transcript 1, Build Order 2.
- Goal 4, Core Feature 3 → `detectHistoryRetained` (API Contract 5),
  Invariant 5, Build Order 2.
- Success Metric 3 → Invariant 2, Invariant 4, Build Orders 1-2.
- Core Feature 4 → the guides and `CONTEXT.md` in Build Orders 1-2 and the
  skill references in Build Order 3.

Deviations from the PRD: None.

## Absence Evidence

1. `RecordLinkPromotion` — command: `rg -n "Promoted|promoted" internal/spec/record_links.go internal/spec/live_links.go internal/speccheck/history_retained.go`; result: no match at `956066a7` (exit 1); neither the repair nor the check reads the record's `promoted` list.
2. The rebase of `unproven` entries — command: `rg -n "Rebase" internal/spec/archive_record.go`; result: only the `Rebase` field and the outcome block (lines 65, 290, 295 and 296); no line rebases `r.Unproven`.

## ADR Alternatives

- ADR-0308: this design is that decision; it takes none of the alternatives
  ADR-0308 rejects (leaving repair to the sanitize, quoting entries as code,
  a second rewrite for `unproven`, a source path with the pull request for a
  promoted file).
- ADR-0247: keeps the record's `unproven` list and its 2 KB target; only the
  outcome is shortened, and entries are never dropped.
- ADR-0277: extends its Legacy Link Rebase from the outcome and promoted
  copies to the `unproven` entries; it does not take its rejected
  alternative of rewriting records an earlier sanitize wrote, except through
  the operator-applied repair it already allows.
- ADR-0278: the History Root still holds only Archive Records; no reference
  stays beside the record.
- ADR-0230: the archived text keeps its relative links' targets, the
  property that ADR rejected leaving to a later checker; here the archive
  rewrites them itself.
- ADR-0260: promotion stays the archive's own decision (`--promote` or a
  Live Document citation); the repair never promotes, it only links a copy
  the record already lists.
- ADR-0261: a locator without a reachable commit still names the record's
  pull request; no ref is pushed and no record is backfilled.
- ADR-0251: Lenient Legacy Reading still produces the legacy `unproven`
  text; the rebase runs after it.
- ADR-0248: the `history-full` tag and the dry run by default are
  unchanged.
- ADR-0254: the update plans the same `record-links` unit with its digest;
  no unit is applied without approval.
- ADR-0299: every listed citation is still rewritten, and a promoted copy is
  the "surviving record or promoted copy" target it already names.
- ADR-0240: a qualifying `partial` still archives under the one partial
  policy; only the text recorded changes.
- ADR-0259: every new top-level test calls `t.Parallel()` first.
- ADR-0271: no test waits on real time.
- ADR-0285: the new tests compile for Windows; none uses a Unix-only call.
- ADR-0104: acceptance rests on `go doc`, published documentation and the
  Delivery Base binary.
- ADR-0276: the network-prohibited documentation row qualifies offline;
  byte invariants compare with the Delivery Base.
- ADR-0265: absence evidence, ADR relations and the coverage cross-read are
  recorded here.
- ADR-0183: each ADR attribution in the PRD quotes its title.
- ADR-0184: the changed plan output is stated as Surface Transcript 1.
- ADR-0280: each guide Task records its own surface fingerprint.
- ADR-0256: the skill references change in the same Spec as the behavior.
- ADR-0092: unrelated to its rejected alternatives; this design changes nothing about triage intent, so the PRD sets no obligation from it.
- ADR-0120: unrelated to its rejected alternatives; this design changes nothing about the single history root, so the PRD sets no obligation from it.
- ADR-0167: unrelated to its rejected alternatives; this design changes nothing about the pre-PR Pull Request row, so the PRD sets no obligation from it.
- ADR-0229: unrelated to its rejected alternatives; this design changes nothing about operator archive resumption, so the PRD sets no obligation from it.
- ADR-0250: unrelated to its rejected alternatives; this design changes nothing about module versions and the coverage record, so the PRD sets no obligation from it.
- ADR-0258: unrelated to its rejected alternatives; this design changes nothing about QA rerun classes and reduced history, so the PRD sets no obligation from it.
- ADR-0262: unrelated to its rejected alternatives; this design changes nothing about cancellation records, so the PRD sets no obligation from it.
- ADR-0269: unrelated to its rejected alternatives; this design changes nothing about ignored Spec files and forward removal recovery, so the PRD sets no obligation from it.
- ADR-0272: unrelated to its rejected alternatives; this design changes nothing about review records across a fix, so the PRD sets no obligation from it.
- ADR-0281: unrelated to its rejected alternatives; this design changes nothing about rebuilt test executables, so the PRD sets no obligation from it.
- ADR-0282: unrelated to its rejected alternatives; this design changes nothing about superseded ADR lineage lines, so the PRD sets no obligation from it.
- ADR-0286: unrelated to its rejected alternatives; this design changes nothing about Correction Reviews, so the PRD sets no obligation from it.
- ADR-0295: unrelated to its rejected alternatives; this design changes nothing about default-branch reachability, so the PRD sets no obligation from it.
- ADR-0298: unrelated to its rejected alternatives; this design changes nothing about profile proof attempts, so the PRD sets no obligation from it.
- ADR-0165: unrelated to its rejected alternatives; this design changes nothing about post-archive blocking reviews, so the PRD sets no obligation from it.
- ADR-0197: unrelated to its rejected alternatives; this design changes nothing about pre-PR reviewer lineage, so the PRD sets no obligation from it.
- ADR-0237: unrelated to its rejected alternatives; this design changes nothing about Delivery Retry merges, so the PRD sets no obligation from it.
- ADR-0279: unrelated to its rejected alternatives; this design changes nothing about Queue Owner status, so the PRD sets no obligation from it.
- ADR-0290: unrelated to its rejected alternatives; this design changes nothing about Item Binaries and Run Worktrees, so the PRD sets no obligation from it.
- ADR-0169: unrelated to its rejected alternatives; this design changes nothing about pre-PR review diffs, so the PRD sets no obligation from it.
- ADR-0233: unrelated to its rejected alternatives; this design changes nothing about skill version regeneration, so the PRD sets no obligation from it.
- ADR-0196: unrelated to its rejected alternatives; this design changes nothing about pre-PR review findings, so the PRD sets no obligation from it.
- ADR-0306: unrelated to its rejected alternatives; its containment of archive and cancel paths is untouched, because no resolver, destination or refusal changes here.

## Integration Points

None external. The commands read the checkout and its local Git history;
nothing is fetched or sent.

## Testing Approach

All tests use temporary repositories, homes and Git configuration, call
`t.Parallel()` first, and compile for Windows.

- `internal/spec/archive_unproven_links_test.go` (new):
  `TestArchiveRebasesTheUnprovenLinksItRecords` drives `spec.Archive` on a
  Characterization-shaped Spec, with subtests for a removed file a default
  branch commit holds, a removed file only the delivery branch holds with a
  pull request, a promoted file, an entry without a link, a link outside the
  folder, and an external Spec Root; `TestLegacyConversionRebasesTheUnprovenLinksItRecords`
  drives the Legacy Archive Folder conversion.
- `internal/cli/archive_unproven_links_test.go` (new):
  `TestArchiveOfAQualifyingPartialPassesTheStrictHistoryCheck` runs `archive
  --promote` and then `spec check --strict` through the command.
- `internal/spec/record_links_promoted_test.go` (new):
  `TestPlanRecordLinksPointsAPromotedTargetAtItsCopy` and
  `TestRecordLinkPromotionRequiresTheListedCopy`.
- `internal/speccheck/history_retained_promoted_test.go` (new):
  `TestHistoryRetainedNamesThePromotedCopyForARecordLink`.
- `internal/cli/history_record_links_promoted_test.go` (new):
  `TestHistorySanitizePointsAPromotedRecordLinkAtItsCopyTranscript` asserts
  Surface Transcript 1; `TestBaselineUpdatePlansAPromotedRecordLinkAtItsCopy`
  asserts the update's JSON citation.
- Existing consumers run unchanged: `TestPlanRecordLinksResolvesAgainstTheRecordThenItsSource`,
  the `TestLegacyUnproven*` tests, the archive and Spec 0058 replay tests,
  and the history and baseline update history tests.

## Build Order

1. The archive rebases each `unproven` entry; the `archive` guide and the
   glossary say so (API Contracts 1-3).
2. The Record Link Repair and the strict check point a promoted target at
   its copy; the `history`, `baseline` and `spec` guides and the glossary
   say so (API Contracts 4-5) (depends on: 1, because both edit `CONTEXT.md`
   and the Behavior Surface Record).
3. The Roundfix skill's `archive`, `baseline` and `spec` references state
   both, with the skill version raised (depends on: 1, 2).
4. The QA gate (depends on: 3).

No Spec is listed in `requires`: Specs 0291-0293 deliver first and touch none
of these code paths.

## Risks & Considerations

- A file name match can point a link to an unpromoted file at a promoted copy
  of the same name in another directory of the Spec (ADR-0308,
  Consequences). The archive refuses two promotions with one name, so no
  link can match two copies.
- The repair still names the record-relative path with the pull request for
  a target neither promoted nor held by a reachable commit; this Spec keeps
  that existing behavior and its test.
- Locator text makes an entry longer; the record's 2 KB target shortens only
  the outcome, so a record with many entries can exceed it, as a long
  outcome-free record could before.

## Glossary

- changes: **Legacy Link Rebase**
- changes: **Lenient Legacy Reading**
- adds: **Record Link Repair**
- not a term: **Characterization Repository** — the name of this TechSpec's fixture, local to the Spec.

## Decisions

- Rebase `unproven` in `BuildArchiveRecord`, not in `archiveUnprovenActions`,
  so the legacy conversion shares it. See ADR-0308.
- Read each entry as text of `<folder>/_prd.md`, the file that holds the
  declarations.
- Represent a promotion as a Citation Rewrite move, the mechanism the
  history plan already uses for promoted copies.
- Name the promoted copy's repository path as the plan's `Rewrite`, as the
  history entries plan already does for moves.
