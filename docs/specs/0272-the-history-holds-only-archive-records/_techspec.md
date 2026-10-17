---
spec: 0272-the-history-holds-only-archive-records
prd: _prd.md
created: 2026-10-16
---

# The history holds only Archive Records — Technical Spec

## Executive Summary

One new file in `internal/spec`, `git_locator.go`, owns the Git Locator: it
chooses the commit that names a removed file and rewrites a citation into
`` `<path>@<commit>` ``. Every writer that removes a file calls it: the Archive
Command for the outcome it copies and the copies it promotes, the Legacy Link
Rebase of a Legacy Archive Folder conversion, and the History Sanitize Command
and the history section of `roundfix baseline update`, which stop reducing and
start removing every retired family, including what the Baseline used to
relocate into them. The Spec Consistency Check gains one error code,
`SC-HISTORY-RETAINED`, which holds the History Root to Archive Records and
Archive Records and Live Documents to links that resolve; it stays silent until
the repository's docs-layout guide carries the new retirement clause. The
Baseline states retirement in one clause, and the delivery cleans this
repository through its own Managed Refresh.

The trade-off accepted is that a citation of a removed file now depends on Git
reachability. A locator names only a commit that the default branch or a tag
reaches, so a file that lived only on a squash-merged branch is named by its
path and the record's pull request instead of a commit. In exchange no
relative link points at a missing file, and the history keeps no copy an Agent
would read as current.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, numbering or
  Archive Record field is added; the record schema stays
  `roundfix/archive-record/v1`, and a Git Locator reuses a repository path and
  a 40- or 64-hex commit already in use. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request, credential or network
  call is added; reachability is proved against local refs, and tests run on
  temporary repositories and homes. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0278 decides the design: "A
  retirement deletes, and the History Root holds only Archive Records", and
  "Git keeps every text." It narrows ADR-0248, whose sanitize stays a "Dry run
  by default."; ADR-0254, which keeps "One plan, one digest"; ADR-0260, where
  "a promoted file keeps its bytes and its own outward relative links are not
  rebased" no longer holds for a promoted copy; ADR-0247, whose proof reads
  "Promoted copies must equal removed blobs"; ADR-0261, where "records already
  archived stay byte-identical" now excepts link destinations to removed
  files; ADR-0277, whose "No other byte changes" gives way to Git Locators for
  removed folder files; ADR-0120 and ADR-0258 for the families and reduced
  entries; and ADR-0092, where "promotion moves the entry into the adopting
  Spec's `references/`" is unchanged. ADR-0186: "A Baseline clause now states
  only behavior the shipped product has." ADR-0202: "A preference is written
  as an obligation about its exception." ADR-0257: "the Baseline states each
  rule once". ADR-0183: "A claim that attributes behavior to a decision record
  now carries a Claim Receipt". ADR-0244: "A Spec declares the domain terms it
  introduces". ADR-0256: "A release waits for the skills that describe a
  changed surface". ADR-0259: "A test runs in parallel unless it names why it
  cannot". ADR-0271: "A test bounds by real time only what it asserts".
  ADR-0095: "The secondbrain inbox is the single capture door". Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AH approval of 2026-10-16 ("Governed Paths declarados") and the 2026-09-30
  skills grant ("considere autorizado a ajustar todas as skills se
  necessário"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0272-the-history-holds-only-archive-records/_authorization.md`;
  bounded files: `.agents/skills/archive-spec/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/roundfix/references/setup.md`,
  `.agents/skills/roundfix/references/spec.md`,
  `docs/agents/docs-layout.md`, `docs/agents/secondbrain.md`,
  `docs/agents/setup-context.json`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`,
  `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/secondbrain.md`,
  `internal/baseline/assets/modules/context-workflow.json`,
  `internal/baseline/assets/modules/secondbrain.json`,
  `internal/baseline/assets/profiles/standard-typescript-monorepo.json`,
  `internal/docscontract/publicdocs_test.go`,
  `internal/docscontract/testdata/corpus-golden.json`,
  `internal/spec/archive.go`,
  `internal/spec/archive_layout_characterization_test.go`,
  `internal/speccheck/backlog.go`, `internal/speccheck/backlog_test.go`,
  `internal/speccheck/coherence.go`, `internal/speccheck/constraints.go`,
  `skills/archive-spec/SKILL.md`, `skills/baseline_skill_contract_test.go`,
  `skills/roundfix/SKILL.md`, and, by the amendments of 2026-10-17,
  `docs/references/coverage-record.json`.

## Current behavior

Characterized on 2026-10-16 on the authoring branch, reading the code and the
tree without writing:

- `PlanHistoryKind` (`internal/spec/history_entries.go`) gives Findings and
  Backlog Entries the action `reduce` and skips an entry that is already
  reduced; `PlanHistoryKinds` never plans the ADR kind, and `ApplyHistoryKind`
  refuses it. `HistoryCitations` skips every source under `docs/history/`, so
  it never reads an Archive Record.
- `historyInventory` (`internal/cli/history.go`) lists Legacy Archive Folders,
  then Findings, Backlog, Review and Handoff units, and skips reduced entries;
  `historyTagCoverage` refuses a batch path the `history-full` tree lacks. The
  plan prints `cites <path>:<line> names <target>` and nothing rewrites it. The
  apply line is `reduced %d file(s), removed %d file(s)`.
- `baseline update`'s history section (`internal/cli/baseline_update_history.go`)
  publishes `applied.reduced`, the `reduce` action and `bytesAfter`, and binds
  `bytesAfter` into the Plan Digest. The history plan is computed before the
  Baseline's History Moves are applied (`internal/cli/baseline_update.go`).
- `DiscoverHistoryLayout` (`internal/baseline/history_layout.go`) relocates
  `docs/findings/_archived/` and `_archived/{findings,adr,backlog}/` into the
  matching History Root family, retired `docs/adr/` files into
  `docs/history/adr/`, terminal `docs/backlog/` entries into
  `docs/history/backlog/`, and finished Review Artifacts into
  `docs/history/reviews/`, so those files survive the run that moved them.
- `BuildArchiveRecord` rebases the outcome only under `Legacy`; the Archive
  Command copies the outcome and every promoted file byte for byte. The record
  of Spec 0262 links `references/<file>.md`, which resolves to nothing.
  `LegacyLinkRebase.Rebase` points a folder destination that is not promoted
  at the record itself, which produced self-links in an adopter's records.
- The Delivery Queue's exact-archive proof (`archiveRetirementIsExact` in
  `internal/cli/deliver_workflow.go`) matches each promoted copy to a removed
  blob by hash, recomputes Live Document relinks with `ArchiveCitationMoves`
  and `RelinkContent`, never compares the outcome, and refuses a relinked
  path under the History Root.
- The Spec Consistency Check resolves Claim Receipts in `docs/adr` and then
  `docs/history/adr`, resolves Rollup members in active or archived Findings,
  and tells a terminal Backlog Entry to "Move" to `docs/history/backlog/`. No
  detector reports a terminal Finding left in `docs/findings/`, a file under a
  retired family, or a broken relative link.
- The tree holds 383 files under `docs/history/`: 253 Archive Records, 98
  reduced Findings and 32 reduced Backlog Entries, every entry naming a
  provenance commit that `main` reaches. Ten `superseded` records (Specs 0120
  to 0129) hold no regeneration grant; nine are cited only by Findings removed
  in the same change, and the 0129 record is also cited by slug in a promoted
  Finding under `docs/references/`. Five links from Live Documents and the
  user guide name reduced Findings; 40 relative links in 32 Archive Records
  are already broken. The 30 records written by the Archive Command (0242 to
  0271) name a `source_revision` that no branch or tag reaches. `history-full`
  holds 129 of the 140 paths this delivery removes.
- The rendered `guide.docs-layout` region is 14,494 bytes. Its clauses
  `docs-one-job-per-directory`, `findings-05-done`, `findings-08-rollup`,
  `findings-09-archive`, `findings-10-live-work-health`,
  `findings-11-rollup-closure` and the last paragraph of
  `backlog-01-operational-contract` describe reduced entries and history
  families; `docs-one-job-per-directory` keeps retired ADRs "whole" and calls
  handoffs the capture door.

## System Architecture

`internal/spec` owns the Git Locator, the citation rewrite and every removal
plan; `internal/cli` selects units, prints and applies; `internal/speccheck`
reports; `internal/baseline` stops relocating into retired families and
carries the clause.

- `git_locator.go` (new) resolves a `GitLocator` and applies a
  `CitationRewrite`. `RelinkContent` and `LegacyLinkRebase.Rebase` delegate
  to `CitationRewrite.Apply`, so the Archive Relink, the Legacy Link Rebase,
  the removal rewrite and the Delivery Queue's recomputation share one
  scanner and one destination rule. The default branch is the one
  `preflight.DetectDefaultBranch` names, the same detection that
  `resolveDefaultBranchHead` in `internal/worktree` wraps; `internal/worktree`
  imports `internal/spec`, so the function is not imported from there.
- The Archive Command (`internal/spec/archive.go`) builds one
  `LegacyLinkRebase` for the folder it removes, with its promotions and a
  locator resolver preferring the record's `source_revision`, and applies it
  to the outcome and to each promoted copy. `BuildArchiveRecord` applies an
  input `Rebase` whether or not the record is legacy.
- `history_entries.go` plans removal for every retired family, from the
  History Root and from the places an earlier layout left them, plus the
  `superseded` records nothing cites, and plans the citation rewrite of the
  removed paths. `record_links.go` (new) finds relative links to missing files
  and self-links, for the check and for the `record-links` unit.
- `historyInventory`, `planHistoryUnit` and the two apply paths in
  `internal/cli` add the `adr`, `superseded` and `record-links` units, apply
  the citation rewrite after the removals and print the new lines.
- `DiscoverHistoryLayout` keeps relocating Legacy Archive Folders and the live
  Review Artifact root, and no longer produces a relocation into the Findings,
  Backlog, ADR or Review families.
- `history_retained.go` (new, `internal/speccheck`) holds the repository-level
  detector wired beside `detectFindingsConsistency`.

## Implementation Design

### Interfaces

```go
// internal/spec/git_locator.go
type GitLocator struct {
	Path        string // repository-relative, slash-separated
	Commit      string // empty when no reachable commit holds Path
	PullRequest string // the citing record's pull request, used when Commit is empty
}

// Text renders the locator; inCode omits the surrounding backticks.
func (l GitLocator) Text(inCode bool) string

func ResolveGitLocator(ctx context.Context, repositoryRoot, path, preferred, pullRequest string) (GitLocator, error)
func GitLocatorHolds(ctx context.Context, repositoryRoot, path, commit string) (bool, error)

type CitationRewrite struct {
	Moves  map[string]string                         // removed path -> surviving repository path
	Locate func(path string) (GitLocator, bool, error) // removed path without a surviving home
}

func (rw CitationRewrite) Apply(content []byte, from, to string) ([]byte, int, error)
```

```go
// internal/spec/legacy_link_rebase.go (changed)
type LegacyLinkRebase struct {
	Folder, RecordPath string
	Files              []string
	Promoted           map[string]string
	Locate             func(path string) (GitLocator, bool, error)
}

func (r LegacyLinkRebase) Rebase(content []byte, from, to string) ([]byte, error)

// internal/spec/record_links.go
type MissingLink struct{ Path string; Line int; Destination, Target string; Self bool }

func MissingRelativeLinks(repositoryRoot, path string, content []byte) []MissingLink
func PlanRecordLinks(ctx context.Context, repositoryRoot, archiveRootRel string, skip map[string]bool) (RecordLinkPlan, error)

const RetirementRuleSentinel = "A retirement deletes the file"
```

### Data Models

No schema changes. `HistoryKindPlan.Action` is always `remove` and its
`BytesAfter` field and the unexported reduced bytes go away. A removal plan for
the `superseded` unit is a `HistoryKindPlan` of kind `specs` whose files are
records directly under the archive root. `RecordLinkPlan` holds the rewritten
bytes per record and the number of links. The `baseline update` JSON keeps
every field it has; `applied.reduced` stays and is always 0.

### API Contracts

1. API Contract: Git Locator resolution. `ResolveGitLocator` returns `Commit =
   preferred` when `preferred` is non-empty, `git cat-file -e
   <preferred>:<path>` succeeds, and the default branch head or a tag reaches
   `preferred` (`git merge-base --is-ancestor`, `git tag --contains`).
   Otherwise it takes the newest commit reachable from the default branch head
   that touches the path (`git log -1 --format=%H <head> -- <path>`), or that
   commit's first parent when that commit deleted it, provided the result holds
   the path. The default branch is the name `preflight.DetectDefaultBranch`
   gives, read from `refs/heads/<name>` or, when that ref is absent, from
   `refs/remotes/origin/<name>`; when neither resolves, only tags count. When
   no commit qualifies, `Commit` is empty. A Git failure other than an absent
   object is returned wrapped with `%w`, and every command runs under the
   context. A path may name a directory.
2. API Contract: Git Locator text. `Text(false)` is `` `<path>@<commit>` ``;
   with an empty commit it is `` `<path>` (pull request #<n>) `` when
   `PullRequest` is set, else `` `<path>` (not reachable from any branch or
   tag) ``. `Text(true)` writes the same without the backticks, for a position
   already inside a code span or fence. Text inside a Git Locator is never a
   citation.
3. API Contract: Citation Rewrite. `Apply(content, from, to)` scans
   destinations with `scanArchiveMarkdown` and resolves each relative
   destination against `from`'s directory to a repository path `T`. When
   `Moves[T]` exists, the destination becomes the escaped path from `to`'s
   directory to it, with its `?` or `#` suffix. Otherwise, when `Locate(T)`
   reports a removal, the whole inline link `[text](dest "title")` becomes
   `text (` + `Text(false)` + `)`, an image drops its `!` the same way, and a
   reference definition `[label]: dest` becomes `\[label]: ` + `Text(false)`.
   Otherwise, when `from` differs from `to`, the destination is rebased to the
   same `T`; else it keeps its bytes. The path text pass replaces a bounded
   occurrence of a moved path with its move, and of a located path with
   `Text(inCode)`, where `inCode` is true inside a code span or fence. Link
   text and every other byte stay. `RelinkContent(content, path, moves)`
   equals `Apply` with `Moves: moves`, `Locate: nil` and `from == to == path`,
   byte for byte.
4. API Contract: Archive Command outcome and promoted copies. The Archive
   Command and every Legacy Archive Folder conversion build a
   `LegacyLinkRebase` whose `Promoted` holds the folder's promotions and whose
   `Locate` reports every path that is the folder or inside it and not
   promoted, resolved with the record's `source_revision` preferred and its
   `pull_request`. It rewrites the outcome (`from` the PRD, QA Report or
   supersession file that supplied it, `to` the record path) and each promoted
   copy (`from` the folder file, `to` its `docs/references/` destination). A
   destination inside the folder therefore points at its promoted copy or
   becomes a Git Locator, never the record itself; a destination outside the
   folder keeps naming the same repository path; folder path text becomes a
   Git Locator or the promoted path. Content without such a destination or
   path is byte-identical. Live Documents keep the Archive Relink unchanged:
   the folder and core artifacts point at the record and cited non-core files
   at their promoted copy.
5. API Contract: outcome shortening. `RenderArchiveRecord` never cuts inside a
   code span: the sentence end it looks for and the cut it falls back to lie
   outside any span, so a shortened outcome keeps each Git Locator whole.
6. API Contract: exact-archive proof. `archiveRetirementIsExact` accepts a
   promoted copy when its basename names exactly one removed folder file and
   the copy's bytes equal that file's bytes after the Archive Command's
   rewrite of API Contract 4, recomputed with the parent as the preferred
   commit. It keeps every other check, including its refusal of a relinked
   path under the History Root.
7. API Contract: removal planning. `PlanHistoryKinds` plans, in this order,
   `findings`, `backlog`, `adr`, `reviews` and `handoffs`, each with the
   action `remove`. A kind's files are its History Root family, every file in
   its legacy trees (`docs/findings/_archived/` and `_archived/findings/` for
   findings, `_archived/backlog/` for backlog, `_archived/adr/` for adr), and
   what an earlier layout left in an active directory: a Finding or Rollup
   directly under `docs/findings/` whose status is `deferred`, `done`,
   `deprecated`, `superseded`, `closed` or `cancelled`; a Backlog Entry
   directly under `docs/backlog/` with a terminal status; an ADR directly
   under `docs/adr/` that is `rejected` or `deprecated`, or `superseded` with a
   `superseded_by` naming an accepted ADR in `docs/adr/` whose text names it;
   and a Review Artifact under the Spec Root's review roots that
   `ClassifyReview` finds finished. A Reduced History Entry is pending like any
   other file. A superseded ADR whose successor does not name it refuses the
   `adr` unit with a reason naming it. An unknown or missing status is not
   planned. `ApplyHistoryKind` validates every path before the first removal,
   removes the files, and prunes the directories they leave empty without
   removing `docs/adr`, `docs/backlog`, `docs/findings` or the Spec Root.
   `ReduceHistoryEntry` is retired; `IsReducedHistoryEntry` stays for
   diagnostics.
8. API Contract: superseded records. The `superseded` unit holds each Archive
   Record directly under the archive root whose `disposition` is
   `superseded`, whose `regeneration` is empty, and that no file surviving the
   plan cites. A citation is a relative link resolving to the record, or the
   record's path or slug as a bounded token outside a Git Locator, in tracked
   Markdown outside skill trees (`.agents/skills/`, `skills/`) and `testdata`
   directories; front matter `requires`, `absorbed_by` and `spec` values are
   text and count.
9. API Contract: removal coverage. Apply refuses, by name, a path to remove
   that neither the `history-full` tree holds nor a commit
   `ResolveGitLocator` finds: `history-full does not hold <path> and no commit
   the default branch reaches holds it`. The annotated tag remains required
   at or before `HEAD`.
10. API Contract: removal rewrite. After a plan's removals, every citation of
    a removed path in an Archive Record, a Live Document or a page under
    `docs/user-guide/` is rewritten with `CitationRewrite`: a Legacy Archive
    Folder's folder and core artifacts move to its record, its promoted files
    to their copies, and every other removed path is located; a Reduced
    History Entry is located at the path and commit its provenance line names
    when that commit qualifies under API Contract 1. A citation elsewhere is
    listed and left. The rewrite never touches a file the same plan removes.
11. API Contract: record links. The `record-links` unit holds each Archive
    Record that survives the plan and has a relative link to a missing file
    or to itself. A missing target is located first as resolved against the
    record, then as resolved against the record's `source` folder; the first
    with a qualifying commit wins, else the record-relative path is located
    with the record's `pull_request`. A self-link is replaced by the Git
    Locator of the record's `source` folder.
12. API Contract: command output. Both commands keep every existing line
    except: the `reduces` lines are gone; a removal unit prints `<unit>:
    removes <n> file(s) (<b> bytes)`; `record-links: rewrites <l> link(s) in
    <n> Archive Record(s)` (`- ` prefixed in `baseline update`); a citation in
    rewrite scope appends `; rewrites to <replacement>` to its `cites` line;
    the apply lines drop `reduced <f> file(s), ` and gain `; rewrote <c>
    citation(s) in <m> file(s)` after the History Full Tag clause. In
    `baseline update` JSON a unit's `action` is `convert`, `remove` or
    `rewrite`, `bytesAfter` is omitted, a `rewrite` unit carries `links`,
    every pending unit carries `paths`, the sorted paths it removes or the
    records it rewrites, which the text plan prints as `  - <path>` lines
    under the unit (task_06, added 2026-10-17 for User Story 3), a
    rewritten citation carries `rewrite`, and `applied` gains
    `rewrittenCitations` and `rewrittenFiles` while `reduced` stays 0. The Plan
    Digest's unit carries `action: remove` and no `bytesAfter`, and its
    citation carries `rewrite`. `History units` counts only files that leave
    the tree. The upgrade notice names every pending unit, reduced entries
    included.
13. API Contract: retained history check. `SC-HISTORY-RETAINED` is an error
    of stage `prd`. It applies only when the repository's `guide.docs-layout`
    managed region, at the path `docs/agents/setup-context.json` records for
    it, contains `RetirementRuleSentinel` after whitespace collapses; otherwise
    the detector is listed as skipped. It reports a file under the Findings,
    Backlog, ADR, Review or Handoff family of the History Root; a terminal
    Finding directly under `docs/findings/`; a relative link from an Archive
    Record, a Live Document or a page under `docs/user-guide/` to a missing
    file, or from a record to itself, outside code; and a Git Locator in those
    files whose commit does not hold its path or that no default-branch head
    or tag reaches. Each finding names the path and line; its Fix names the
    deletion or the replacement `Text(false)` that `ResolveGitLocator`
    returns. Outside Git the locator parts are skipped by name.
14. API Contract: readers. A Claim Receipt for `ADR-NNNN` resolves in
    `docs/adr/` and otherwise in Git: the newest commit reachable from `HEAD`
    that added or modified `docs/adr/NNNN-*.md` or `docs/history/adr/NNNN-*.md`,
    read with `git show`. A Rollup member resolves as an unresolved Finding, an
    archived file, or a basename that Git shows deleted under `docs/findings/`
    or `docs/history/findings/`. The Fix of `SC-BACKLOG-UNMOVED` for a
    terminal entry says to delete it, Git keeping its text.
    `SC-ARCHIVE-LICENSE` is unchanged.
15. API Contract: the Baseline. The `context-workflow` module's
    `docs-one-job-per-directory` states the directory jobs only, with the
    history root holding Archive Records, `docs/_inbox/` only for a repository
    without a Secondbrain and handoffs for session continuity.
    `findings-09-archive` becomes the one retirement clause, mandatory, opening
    with `RetirementRuleSentinel`, naming what retires, the Git Locator that
    `roundfix spec check` prints and the runtime's `roundfix history
    sanitize` for what an earlier layout left, and listing
    `clause.context.findings-10-live-work-health` and
    `clause.context.findings-11-rollup-closure` in `replaces`.
    `findings-05-done`, `findings-08-rollup` and the last paragraph of
    `backlog-01-operational-contract` point at that clause. The `secondbrain`
    module's `clause.secondbrain.capture-trigger` adds the conditional lesson
    capture. Both modules take the next recorded version, and the rendered
    `guide.docs-layout` region is at most 11,000 bytes.

### Surface Transcripts

1. Surface Transcript: `history sanitize` plans the removal of a Reduced
   History Entry and the rewrite of the user guide link that cites it. The
   fixture is a temporary repository whose `history-full` tag holds one
   reduced Finding under the History Root, whose provenance commit on the
   default branch holds the Finding at its active path, and whose user guide
   page `docs/user-guide/a.md` links it on line 3.

   ```transcript
   $ roundfix history sanitize
   stdout:
   history sanitize plan: 1 unit(s) pending; 1 file(s) (<bytes> bytes) leave docs/history
   findings: removes 1 file(s) (<bytes> bytes)
   cites docs/user-guide/a.md:3 names <reduced entry>; rewrites to <locator>
   apply with: roundfix history sanitize --apply --batch <n> (needs the annotated tag history-full at or before HEAD)
   stderr:
   exit: 0
   ```

2. Surface Transcript: applying that batch removes the entry and rewrites the
   citation.

   ```transcript
   $ roundfix history sanitize --apply --batch 1
   stdout:
   history sanitize applied 1 unit(s): wrote 0 Archive Record(s), removed 1 file(s) (<bytes> bytes) kept in Git at <revision> and tag history-full; rewrote 1 citation(s) in 1 file(s); promoted 0 file(s) to docs/references/; 0 unit(s) remain
   stderr:
   exit: 0
   ```

### Invariants

```text
1. A Git Locator a writer produces names a commit that holds its path and that the default branch or a tag reaches.
2. No writer produces a link from an Archive Record to itself.
3. RelinkContent keeps the bytes it produces at the delivery base for every input.
4. Content with no destination or path text to rewrite keeps its bytes in every writer.
5. A plan writes nothing; apply removes only planned paths and rewrites only planned citations.
6. A second plan after an applied plan reports no pending history for the units it applied.
```

## Coverage Map

Deviations from the PRD:

- Core Feature 5: the delivery runs the cleanup through the history section of
  the Managed Refresh, `roundfix baseline update`, which runs the History
  Sanitize Command's planning and applies it in one plan; the History Sanitize
  Command's apply requires a clean working tree, which the Task's own edits
  break. One promoted Finding's routing pointer moves by hand from the 0129
  record to its successor so that record is no longer cited.
- Core Feature 6: the retirement rule is one clause whose sentences name the
  deletion, the Git Locator and the runtime's command. A separate clause per
  obligation would need a new Source Baseline row or a reused clause identity
  that misstates the Source Baseline transition.

- Goal 1, User Story 1, Core Feature 1 → API Contracts 7 and 15, the
  `context-workflow` clause
- Goal 2, User Story 2, Core Feature 3 → `ResolveGitLocator`,
  `CitationRewrite` (API Contracts 1-6, 10, 11)
- Goal 3, User Story 3, Core Feature 4 → `PlanHistoryKinds`, the `superseded`
  and `record-links` units, removal coverage (API Contracts 7-12)
- Goal 4, User Stories 4-5, Core Feature 6 → API Contract 15
- Core Feature 2 → `SC-HISTORY-RETAINED`, the readers (API Contracts 13-14)
- Core Feature 5 → the dogfood Task's Managed Refresh
- Core Feature 7 → the user guide pages, the Roundfix and archive-spec skills,
  `CONTEXT.md`, the ADR notes
- Success Metric 1 → the dogfood Task and the QA gate
- Success Metric 2 → `internal/speccheck/history_retained_test.go` and the QA
  sabotage rows
- Success Metric 3 → Invariants 1-2, `internal/spec/git_locator_test.go`,
  API Contract 13's locator check
- Success Metric 4 → `internal/baseline/retirement_rule_clauses_test.go`
- Success Metric 5 → `internal/cli/history_removal_test.go` and the QA
  adopter dry run

## Absence Evidence

1. `GitLocator` and `CitationRewrite` — command: `grep -rn '@<commit>\|GitLocator\|merge-base.*--is-ancestor' internal/spec --include=*.go`; result: no locator type or reachability helper in `internal/spec`; the only ancestry checks are in `internal/cli` and `internal/worktree`.
2. Rewriting a whole link into text and code — command: `grep -rn 'func rewriteArchiveLinks\|func RelinkContent\|func (r LegacyLinkRebase) Rebase' internal/spec`; result: all three replace only the destination bytes inside `](…)`.
3. Removal of reduced entries and ADRs — command: `grep -n 'ArchiveKindADR\|plan.Action = "reduce"' internal/spec/history_entries.go`; result: `reduce` for Findings and Backlog Entries, and no ADR kind planned.
4. `superseded` and `record-links` units — command: `grep -rn 'ArchiveSuperseded\|record-links' internal/cli/history.go internal/cli/baseline_update_history.go`; result: no match.
5. `MissingRelativeLinks` — command: `grep -rn 'does not resolve\|missing file' internal/spec internal/speccheck --include=*.go`; result: only `prepareArchiveLinks`, which checks links leaving a Spec folder that a legacy move relocates.
6. `SC-HISTORY-RETAINED` and the sentinel — command: `grep -rn 'SC-HISTORY-RETAINED\|A retirement deletes the file' internal docs/agents`; result: no match.

## ADR Alternatives

- ADR-0278: this Spec's decision; its rejected alternatives (keeping superseded ADRs and completed change folders) stay rejected, and the lineage lives in the Archive Record, the successor ADR and the Git Locator.
- ADR-0248: no rejected alternative is taken; the dry run, the batches and the History Full Tag stay.
- ADR-0254: no rejected alternative is taken; one plan and one digest still cover the Baseline and the history.
- ADR-0260: the rewrite of citing files it rejected for the update is now taken for removed files, as ADR-0278 decides; its other rejected alternatives stay rejected.
- ADR-0247: no rejected alternative is taken; the proof keeps matching promoted copies to removed files, through the rewrite.
- ADR-0261: no rejected alternative is taken; no record is backfilled beyond link destinations to missing files.
- ADR-0232: no rejected alternative is taken; Merge Evidence still reads the delivered Spec's Archive Record, which no unit removes.
- ADR-0277: its rejection of citation-only reporting stays; a folder file without a promoted copy now gets a Git Locator instead of the record.
- ADR-0120: no rejected alternative is taken; the History Root stays the single root.
- ADR-0258: no rejected alternative is taken; its authoring rules stand.
- ADR-0092: no rejected alternative is taken; promotion is unchanged.
- ADR-0186: no rejected alternative is taken; each clause states what the product does.
- ADR-0202: no rejected alternative is taken; openness is an obligation about its exception.
- ADR-0257: no rejected alternative is taken; the retirement rule is stated once.
- ADR-0183: no rejected alternative is taken; every attribution carries a receipt.
- ADR-0244: no rejected alternative is taken; the terms follow its declaration contract.
- ADR-0256: no rejected alternative is taken; the skills Task updates every covering reference.
- ADR-0259: no rejected alternative is taken; every new test calls `t.Parallel()` first.
- ADR-0271: no rejected alternative is taken; no new test waits on real time.
- ADR-0095: no rejected alternative is taken; the Secondbrain inbox stays the single capture door.
- ADR-0203: no rejected alternative is taken; the retirement clause is product behavior every adopter runs.
- ADR-0184: no rejected alternative is taken; the two changed outputs are stated as Surface Transcripts.
- ADR-0251: no rejected alternative is taken; Lenient Legacy Reading still reads the folders a conversion removes.
- ADR-0262: no rejected alternative is taken; the Cancel Command and its record are untouched.
- ADR-0269: no rejected alternative is taken; the archive and cancel refusals and their recovery are untouched.
- ADR-0250: no rejected alternative is taken; each module version is chosen when the record step runs.
- ADR-0280: no rejected alternative is taken; every Task that edits a user guide page re-records the Behavior Surface Record itself.
- ADR-0281: no rejected alternative is taken; a Verification builds its binary into a new temporary directory.
- ADR-0165: no rejected alternative is taken; review parking after archive is unchanged.
- ADR-0169: no rejected alternative is taken; the pre-PR review diff is unchanged.
- ADR-0196: no rejected alternative is taken; review finding validation is unchanged.
- ADR-0233: no rejected alternative is taken; skill version raises follow its merge rule.
- ADR-0237: no rejected alternative is taken; Delivery Retry is unchanged.
- ADR-0272: no rejected alternative is taken; review records are unchanged.
- ADR-0279: no rejected alternative is taken; queue status is unchanged.

## Integration Points

None outside the repository. Git is read locally through `git` commands in
the repository that holds the Spec Root, as the existing archive and history
readers do.

## Testing Approach

New tests call `t.Parallel()` first, build temporary repositories with
`gittest` and the existing fixtures (`newHistoryFixture`,
`baselineHistoryFixture`, `legacyFixture`), and set `HOME` to a temporary
directory. A squash merge is simulated by a commit on a branch that is
deleted, so no ref reaches it. Tests that change are named in the Task that
changes them, with the contract they move to.

- `internal/spec/git_locator_test.go`: commit choice, the remote-tracking
  fallback, an unreachable commit, the text forms, the rewrite of links,
  images, reference definitions, code spans and fences, and byte equality of
  `RelinkContent` with the delivery base on a corpus of the existing Archive
  Relink fixtures.
- `internal/cli/archive_git_locator_test.go`: `roundfix archive` on a
  repository whose planning commit is on the default branch, and the Delivery
  Queue's proof of that archive.
- `internal/spec/history_removal_test.go` and
  `internal/cli/history_removal_test.go`: every family and legacy location,
  the superseded unit, coverage, the rewrite, both transcripts, the JSON and a
  second run.
- `internal/spec/record_links_test.go`, `internal/speccheck/history_retained_test.go`
  and `internal/cli/history_record_links_test.go`: the gate, each reported
  shape, the readers and the `record-links` unit through both commands.
- `internal/baseline/retirement_rule_clauses_test.go`: the clause text, the
  region size, the sentinel and the `replaces` transition.

## Build Order

1. Git Locator, Citation Rewrite, the Archive Command's outcome and promoted
   copies, the Legacy Link Rebase change, the exact-archive proof, the archive
   and deliver guide pages, and tests (no dependency).
2. Removal of every retired family, the `superseded` unit, removal coverage
   and rewrite, the Baseline's relocation change, command output, the upgrade
   notice, the history, baseline and upgrade guide pages, and tests (depends
   on: 1, which provides the locator and the rewrite).
3. `record-links` unit, `SC-HISTORY-RETAINED` and the readers, the corpus
   golden, the spec guide page, the Roundfix and archive-spec skills and the
   glossary, and tests (depends on: 2, which edits the same commands and
   removal plan).
4. Baseline modules, version record, derived files, the Managed Refresh that
   applies this repository's history, the context-driven development guide and
   the clause tests (depends on: 3, whose check the refreshed guide turns on).
5. QA gate (depends on: 4).

No other Spec is a prerequisite. Spec 0273 changes Delivery Queue code; if it
merges first and touches `internal/cli/deliver_workflow.go`, rebase and run the
strict check again before the Run starts.

## Risks & Considerations

- A plan amended on the default branch between an archive and its proof can
  move the newest commit that holds a folder file, so the proof's
  recomputation differs and the item parks; `roundfix deliver retry`
  re-archives.
- The Archive Command cannot name the delivery pull request when it writes
  the record, so a link to a file that only the Run branch holds reads
  `(not reachable from any branch or tag)`; the record's `pull_request`,
  added later, names the pull request whose head keeps it.
- An adopter whose binary is newer than its Baseline gets no refusal until
  `baseline update` refreshes the guide; an adopter on an older binary keeps
  reducing until it upgrades, and the newer binary then removes those entries.
- A repository whose Legacy Archive Folders still sit under `_archived/specs`
  needs one run to relocate them and a second to convert them, as before.
- An active Spec authored beside this one with a relative link to a missing
  file fails the check once this repository's guide is refreshed; its author
  rebases and runs the strict check again, as the autonomous-work guide asks.
- The check reads Git for each Git Locator; it batches commits and paths so
  the cost grows with distinct commits, not with citations.

## Glossary

- changes: **History Relocation**
- changes: **History Citation**

## Decisions

- One function chooses every Git Locator commit, preferring the provenance
  commit the file already names. See ADR-0278.
- A whole link becomes its text followed by the locator in code, so no
  relative destination names a missing file. See ADR-0278.
- The Archive Command rewrites the outcome and each promoted copy with the
  Legacy Link Rebase, and the exact-archive proof recomputes the copies. See
  ADR-0278, which narrows ADR-0247 and ADR-0260.
- Relocations into retired families become removals in the same plan, Review
  Artifacts included, so a second run reports the history current.
- The new check is gated by a sentence of the refreshed guide, not by a
  module version, because no adopter file records module versions.
- Public JSON keeps `applied.reduced` at 0 and only gains fields.
- The retirement rule is one clause with `replaces`, because a new clause
  needs a Source Baseline row.
