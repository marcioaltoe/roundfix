---
spec: 0255-an-archive-that-relinks-and-an-update-that-lists-citations
prd: _prd.md
created: 2026-10-12
---

# An archive that relinks and an update that lists citations — Technical Spec

## Executive Summary

One new file in `internal/spec` holds the shared parts: the Live Document
classifier, a citation rewrite that changes only link destinations and path
text, and the planning of an Archive Relink. `spec.Archive` uses them to
promote the files Live Documents cite and rewrite those citations before it
removes the folder, and the record gains an optional `relinked` list. The
Delivery Queue stages those files and proves each rewrite by recomputing it
from the parent commit. The update's history section reuses the existing
detector, `spec.HistoryCitations`, to list the History Citations its selected
units break, binds them into the Plan Digest, accepts `--promote`, and reports
Stale History Citations from the History Full Tag.

The trade-off accepted: an archive now changes files outside the Spec folder,
so the exact-archive proof grows a recomputation instead of staying a pure
path check. The alternative, refusing the archive while a citation exists,
would stall the normal queue order. The update keeps its subject narrow: it
lists and copies, and never rewrites a citing file.

## Project Constraints

- Identifier strategy: not applicable. No identifier scheme changes. The
  record gains the optional lowercase key `relinked`; the update result gains
  `citations` and `staleCitations` with lowercase keys under `history`; the
  update gains the lowercase repeatable flag `--promote`. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable. No request, credential or network
  call is added. Both commands read files and run local Git; tests use
  temporary repositories built with Git and a temporary home. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0260 (this Spec) governs the Archive
  Relink, its refusals and warnings, the `relinked` key, the queue proof, the
  History Citation listing, `--promote` on the update and the Stale History
  Citation report: "Each Live Document citation of a promoted file or of the
  folder is rewritten in place". One owner per source stays, ADR-0083:
  "secondary Specs link the owner's copy". Promoted bytes stay, ADR-0247:
  "Promoted copies must equal removed blobs". The queue precedent is
  ADR-0230: "whose only content changes besides the archive stamp are these
  rewrites". The digest precedent is ADR-0173: "The warnings are part of the
  plan, so the Plan Digest binds them". ADR-0254 governs the history section
  and ADR-0260 narrows "The history section changes only files under the
  History Root" for `--promote`. The tag stays, ADR-0248: "holds every path the
  batch removes or rewrites". ADR-0184: "A TechSpec states a command surface as
  a transcript", answered in Surface Transcripts. ADR-0100, ADR-0121, ADR-0177
  and ADR-0215 are cited by those records and do not apply: no instruction
  carrier, legacy archive move, Relocation Citation scan or archived evidence
  changes. ADR-0116, ADR-0168, ADR-0176 , ADR-0183, ADR-0251 and ADR-0258 cite listed records and do not
  apply: this Spec changes no Spec check citation rule, no Lenient Legacy
  Reading or failed-QA disposition, and no retirement of a Finding or Backlog
  Entry. Source: `docs/agents/domain.md`.
- Tooling authority: applicable. task_01 edits owned skills and task_02 edits
  `internal/spec/archive.go`. Express maintainer authorization: "considere
  autorizado a ajustar todas as skills se necessário" and the 2026-10-11 cycle
  grant "Concedo" for the Governed Paths each Spec declares. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.
  Spec-contained authorization record:
  `docs/specs/0255-an-archive-that-relinks-and-an-update-that-lists-citations/_authorization.md`.
  Bounded files: `.agents/skills/archive-spec/SKILL.md`,
  `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `.agents/skills/write-idea/SKILL.md`, `.agents/skills/write-prd/SKILL.md`,
  `.agents/skills/write-tasks/SKILL.md`,
  `.agents/skills/write-techspec/SKILL.md`, `internal/spec/archive.go`,
  `skills/archive-spec/SKILL.md`, `skills/roundfix/SKILL.md`,
  `skills/write-idea/SKILL.md`, `skills/write-prd/SKILL.md`,
  `skills/write-tasks/SKILL.md`, `skills/write-techspec/SKILL.md`.

## Current behavior

Read on 2026-10-12 at `7a8b7f70` (v0.65.0):

- `spec.Archive` checks the Spec's outward links with `prepareArchiveLinks`,
  copies each `--promote` file to `docs/references/<basename>`, writes the
  record and removes the folder. It reads no file outside the folder, so an
  inbound citation is neither found nor reported. Its doc comment says "no PRD
  or link is rewritten".
- `speccheck.ActiveSpecPathPins` refuses an archive while a non-Markdown file
  names the active folder; Markdown is excluded by design.
- The Delivery Queue's `archiveRetirementIsExact` accepts a change set equal to
  the folder, the record and the promoted copies, each promoted blob equal to a
  removed blob. Any other changed path fails the proof.
- `roundfix history sanitize` computes citations with `spec.HistoryCitations`
  over the batch's files and its folders and prints
  `cites <path>:<line> names <target>`. The update's
  `planBaselineUpdateHistory` never calls the detector, and its digest payload
  holds only the revision, the tag and the units.
- `spec.HistoryCitations` lists tracked Markdown with `git ls-files`, skips
  `docs/history/`, and matches a relative link destination or a bounded
  repository-relative path. It errors outside a Git work tree.
- Existing `internal/spec` archive tests call `spec.Archive` with a
  `RepositoryRoot` that is not a Git work tree.
- Measured outside this repository: the characterization in
  `docs/references/2026-10-12-live-citations-characterization.md`.

Research record (2026-10-12). Secondbrain: `wiki/index.md`, then
`qmd query "roundfix archive promote references inbound links broken citations history sanitize"`,
which pointed at the mirrored Archive Record of the Spec that made archived
Specs keep their outward links and at the two triaged adopter reports,
`inbox/roundfix/_triaged/2026-10-09-o-baseline-update-sanitiza-o-historico-sem-relatar-os-links-que-quebra.md`
and `inbox/roundfix/_triaged/2026-10-10-spec-secundaria-linka-a-references-de-outra-spec-e-quebra-no-archive.md`.
The reports fixed the target classes and asked that adopters learn what to do
with Specs already written; the adopter note answers that nothing is needed.
Exa: the MD Links extension for VS Code
(marketplace.visualstudio.com/items?itemName=TarasErmolenko.fix-refs-vscode)
supports rewriting only the destination and warning about links left pointing
at a removed file; Zensical Studio
(zensical.org/docs/studio/workflows/link-intelligence/) also rebases a moved
file's own relative links, which this design declines, ADR-0247: "Promoted copies
must equal removed blobs".

## System Architecture

```text
internal/spec/live_citations.go (new)
  IsLiveDocument, ArchiveCitationMoves, RelinkContent, PlanArchiveRelink
spec.Archive
  eligibility -> outward links -> PlanArchiveRelink (Git work tree only)
  -> refusals (collision, duplicate basename, dirty Live Document)
  -> record (promoted + relinked) -> write record, copies, rewrites -> remove folder
cli archive: confirmation suffix, stderr warnings, --plan lines
cli deliver: stage record.Relinked; archiveRetirementIsExact recomputes RelinkContent
cli baseline update history section
  selected units -> spec.HistoryCitations -> live, offer, promotedTo
  --promote -> PlanLegacyConversion with Promote
  digest += citations, promotions (omitted when empty)
  tag present -> spec.StaleHistoryCitations (new, internal/spec/stale_citations.go)
```

## Implementation Design

### Interfaces

```go
// internal/spec/live_citations.go
func IsLiveDocument(path, specRootRel, archiveRootRel string) bool
// ArchiveCitationMoves maps each removed repository path to its new home:
// promoted sources to their destination; the folder and its core artifacts to the record.
func ArchiveCitationMoves(folderRel, recordRel string, folderFiles []string, promotions map[string]string) map[string]string
// RelinkContent rewrites link destinations and bounded repository paths naming a move.
func RelinkContent(content []byte, citingPath string, moves map[string]string) ([]byte, int)
type ArchiveRelinkPlan struct {
	Promote    []string          // folder-relative files promoted automatically, sorted
	Relinked   map[string][]byte // repository-relative Live Document -> rewritten bytes
	Unrelinked []HistoryCitation // citations still present after the rewrite
}
func PlanArchiveRelink(ctx context.Context, repositoryRoot, specRootRel, archiveRootRel, folderRel, recordRel string, manual []string) (ArchiveRelinkPlan, error)

// internal/spec/stale_citations.go
func StaleHistoryCitations(ctx context.Context, repositoryRoot, historyRoot string, held map[string]bool) ([]HistoryCitation, error)
```

### Data Models

- `ArchiveRecord.Relinked []string`, rendered as `relinked:` after `promoted`
  only when non-empty, read when present. `ArchiveResult` gains `Relinked
  []string` and `Unrelinked []HistoryCitation`.
- `baselineUpdateHistory` gains `Citations` (`json:"citations"`) and
  `StaleCitations` (`json:"staleCitations"`), both always present, never null,
  of `baselineUpdateHistoryCitation{Path, Line, Target, Live, Offer,
  PromotedTo}` with keys `path`, `line`, `target`, `live`, `offer` and
  `promotedTo`, the last two omitted when empty.
- `baselineUpdateRequest.promote []string`.

### API Contracts

1. API Contract: Live Documents. `IsLiveDocument` is true for a repository
   path ending in `.md` under the Spec Root and not under its archive root,
   under `docs/backlog/`, under `docs/findings/`, or under `docs/references/`.
   Nothing under the History Root is live.
2. API Contract: the rewrite. `RelinkContent` reads the destinations of inline
   links and reference definitions with the archive's Markdown scanner, code
   spans masked. A relative destination whose decoded path, resolved against
   the citing file's directory, equals a move's source, or the folder with a
   trailing slash, becomes the relative path from that directory to the move's
   destination, URL-escaped, with its `?` or `#` suffix kept. Outside those
   destinations, each occurrence of a move's source as a bounded
   repository-relative path (the boundary rule `spec.HistoryCitations` uses),
   in prose, code spans or fences, becomes the destination; the longest source
   matches first. No other byte changes. The count is the number of
   replacements.
3. API Contract: Archive Relink planning. When `RepositoryRoot` is inside a
   Git work tree, `PlanArchiveRelink` lists every citation of the folder's
   files and of the folder with `spec.HistoryCitations`, excluding the folder
   itself and the archive root. A cited file that is not a core artifact (the
   set `--promote` refuses) and that a Live Document cites joins `Promote`,
   together with each `manual` path. `ArchiveCitationMoves` maps each promotion
   to `docs/references/<basename>`, and the folder and each core artifact to
   the record. Every Live Document holding a citation is rewritten with
   `RelinkContent`, and a document whose bytes change joins `Relinked`. The
   citations the detector still finds in the rewritten bytes, and every
   citation outside Live Documents, are `Unrelinked`. Outside a Git work tree
   the plan is empty.
4. API Contract: archive refusals. Before writing, `spec.Archive` refuses, with
   exit 2 from the command and no file changed, when a promotion destination
   exists, when two promotions share a basename, or when
   `git status --porcelain --untracked-files=all` over the relinked paths is not
   empty: `Spec "<slug>" cannot archive: <path> has changes not committed at HEAD; commit them so the archive can relink it`.
   A Git failure inside a work tree also refuses. An existing refusal keeps its
   wording.
5. API Contract: archive writes. The record lists every promotion, sorted,
   under `promoted`, and the relinked paths, sorted, under `relinked`. The
   archive writes the record, the promoted copies and the relinked bytes, and
   then removes the folder; a write failure restores every relinked file and
   removes every file it created. The confirmation appends
   `; relinked <m> file(s)` after the promotion suffix when `<m>` is positive.
   Each Unrelinked citation prints one stderr line:
   `roundfix: warning: archive leaves <path>:<line> naming removed <target> unrelinked`.
   `--plan` prints `promotes <file> to docs/references/<basename>`,
   `relinks <path>: <n> citation(s)` and
   `leaves <path>:<line> naming <target> unrelinked` before its
   `promote with:` line. Without inbound citations every output is today's.
6. API Contract: the queue. The archive stage stages `record.Relinked` with
   the record and the promotions. `archiveRetirementIsExact` accepts each
   relinked path only when it is outside the folder and the History Root,
   existed at the parent, and its committed (or working-tree) bytes equal
   `RelinkContent` of its parent bytes under `ArchiveCitationMoves`. The
   promotions map pairs each promoted destination with the one folder file of
   the same basename and blob; zero or several candidates fail the proof.
7. API Contract: History Citations in the update. When the history section
   selects a unit, the update runs `spec.HistoryCitations` over the files of
   each selected folder and each selected `remove` kind, then over each
   selected folder path, dropping a folder citation whose file and line already
   hold one, as the History Sanitize Command does. `reduce` kinds are left out,
   because their files stay. Each citation records whether its file is a Live
   Document. A Live Document citation of a non-core file of a selected folder
   carries `offer` equal to its target, unless the target is promoted, which
   sets `promotedTo`. The Plan Digest payload gains `citations` (path, line,
   target, in report order) and, per folder unit, `promoted`; both are omitted
   when empty, so a digest without them is today's.
8. API Contract: `--promote` on the update. The flag is repeatable and refused
   with `--no-history` as invalid input. Each value must be a clean
   repository-relative file inside a selected folder, outside every Refused
   Unit, not a core artifact, with a basename no other promotion uses and no
   existing `docs/references/` file. Otherwise the update exits 2 with category
   `invalid`, names the value and writes nothing. Valid values re-plan their
   folder with `PlanLegacyConversion` and `Promote`, so apply copies them and
   the record lists them. On a history write failure the next action also names
   `docs/references`.
9. API Contract: text output. After `History refused`, the update prints, when
   present:

   ```text
   History citations: <n>
   - cites <path>:<line> names <target>
   - cites <path>:<line> names <target> (live)
   - cites <path>:<line> names <target> (live; promote with --promote <target>)
   - cites <path>:<line> names <target> (promoted to docs/references/<basename>)
   History stale citations: <n>
   - cites <path>:<line> names <target>; history-full holds it: git show history-full:<target>
   ```

   The history block also prints when only citations or stale citations exist.
10. API Contract: Stale History Citations. When `history-full` exists and the
    section is not `blocked` or `skipped`, `StaleHistoryCitations` reads every
    tracked Markdown file outside the History Root and reports each relative
    link destination or bounded repository path that names a file of the tag
    tree under the History Root which is absent from the working tree. A
    directory never counts. Each line is looked up in a set, so the scan stays
    linear in the lines read. The report never enters the digest, state,
    category or exit code.

### Invariants

```text
1. An archive with no inbound citation writes, records and prints exactly what it does today.
2. Outside a Git work tree, spec.Archive behaves exactly as today.
3. A relinked file differs from its HEAD bytes only in link destinations and bounded path text naming a move.
4. A refused archive changes no file.
5. Every promoted copy equals a removed blob.
6. The queue proof accepts no changed path it cannot recompute.
7. With no History Citation and no promotion, the update's digest and text equal today's.
8. The update never changes a citing file.
9. A Stale History Citation never changes the update's state, category, exit code or digest.
10. Nothing under the History Root is a Live Document or is rewritten.
```

### Surface Transcripts

The archive transcript runs the binary built from the audited head, as
`roundfix`, from the root of a temporary repository: an owner Spec
`0001-owner`, completed with a passing QA Report and an adopted
`references/2026-01-01-source.md`; an active `0002-secondary` whose `_prd.md`
links that source relatively and whose `task_01.md` names it by repository
path; an open Backlog Entry `docs/backlog/2026-01-02-follow-up.md` linking the
owner's `_prd.md`; and a handoff `docs/handoffs/2026-01-03-notes.md` linking the
source. All committed, under a temporary home.

The update transcripts run with `go run -buildvcs=false ./cmd/roundfix`
against a temporary adopted repository whose Baseline is current, with no
`history-full`, holding the committed Legacy Archive Folder
`docs/history/specs/0001-legacy` (with `references/2026-01-01-lesson.md`), an
unresolved Finding `docs/findings/2026-01-04-gap.md` linking that lesson, and a
handoff `docs/handoffs/2026-01-03-notes.md` linking the folder's `_prd.md`.

1. Surface Transcript: the owner's archive promotes, relinks and warns.

   ```transcript
   $ roundfix archive 0001-owner
   stdout:
   archived 0001-owner -> docs/history/specs/0001-owner.md; removed <n> file(s) (<bytes> bytes) kept in Git at <hex>; promoted 1 file(s) to docs/references/; relinked 3 file(s)
   stderr:
   roundfix: warning: archive leaves docs/handoffs/2026-01-03-notes.md:<line> naming removed docs/specs/0001-owner/references/2026-01-01-source.md unrelinked
   exit: 0
   ```

2. Surface Transcript: the update preview lists the History Citations.

   ```transcript
   $ go run -buildvcs=false ./cmd/roundfix baseline update --repo <repo> --no-skills --format text
   stdout:
   Baseline update: plan ready
   Category: approval
   Result: guidance matches the current Baseline catalog; 1 history unit(s) pending sanitize
   ...
   History: pending
   History units: 1 (<files> file(s), <bytes> bytes leave docs/history)
   - folder docs/history/specs/0001-legacy: removes <n> file(s) (<b> bytes) and writes docs/history/specs/0001-legacy.md (<disposition>)
   History citations: 2
   - cites docs/findings/2026-01-04-gap.md:<line> names docs/history/specs/0001-legacy/references/2026-01-01-lesson.md (live; promote with --promote docs/history/specs/0001-legacy/references/2026-01-01-lesson.md)
   - cites docs/handoffs/2026-01-03-notes.md:<line> names docs/history/specs/0001-legacy/_prd.md
   History tag: creates annotated history-full at <hex>; push it with git push origin history-full
   Plan Digest: <digest>
   Next action: <next>
   stderr:
   exit status 3
   exit: 1
   ```

3. Surface Transcript: the confirmed update promotes on request only.

   ```transcript
   $ go run -buildvcs=false ./cmd/roundfix baseline update --repo <repo> --no-skills --yes --promote docs/history/specs/0001-legacy/references/2026-01-01-lesson.md --format text
   stdout:
   Baseline update: verified
   ...
   History: applied
   History applied: 1 unit(s): wrote 1 Archive Record(s), reduced 0 file(s), removed <n> file(s) (<b> bytes) kept in Git at <hex> and tag history-full
   History citations: 2
   - cites docs/findings/2026-01-04-gap.md:<line> names docs/history/specs/0001-legacy/references/2026-01-01-lesson.md (promoted to docs/references/2026-01-01-lesson.md)
   - cites docs/handoffs/2026-01-03-notes.md:<line> names docs/history/specs/0001-legacy/_prd.md
   History tag: created annotated history-full at <hex>; push it with git push origin history-full
   ...
   stderr:
   exit: 0
   ```

4. Surface Transcript: the second update names the Stale History Citations.

   ```transcript
   $ go run -buildvcs=false ./cmd/roundfix baseline update --repo <repo> --no-skills --format text
   stdout:
   Baseline update: current
   Result: the repository already matches the current Baseline catalog
   ...
   History: current
   History stale citations: 2
   - cites docs/findings/2026-01-04-gap.md:<line> names docs/history/specs/0001-legacy/references/2026-01-01-lesson.md; history-full holds it: git show history-full:docs/history/specs/0001-legacy/references/2026-01-01-lesson.md
   - cites docs/handoffs/2026-01-03-notes.md:<line> names docs/history/specs/0001-legacy/_prd.md; history-full holds it: git show history-full:docs/history/specs/0001-legacy/_prd.md
   History tag: history-full at <hex> holds every planned path
   ...
   stderr:
   exit: 0
   ```

## Coverage Map

- Goal "Live Document citations keep resolving" → API Contracts 1, 2, 3 and 5,
  Invariants 3, 5 and 10, Surface Transcript 1.
- Goal "names what it does not rewrite; refuses collisions and dirty files" →
  API Contracts 4 and 5, Invariant 4, Surface Transcript 1.
- Goal "the queue commits and proves a relinking archive" → API Contract 6,
  Invariant 6.
- Goal "the plan lists History Citations; `--promote` copies only" → API
  Contracts 7, 8 and 9, Invariants 7 and 8, Surface Transcripts 2 and 3.
- Goal "Stale History Citations with the restore command" → API Contract 10,
  Invariant 9, Surface Transcript 4.
- User Story 1 → API Contracts 3, 5 and 6.
- User Story 2 → API Contract 5.
- User Story 3 → API Contracts 7 and 9.
- User Story 4 → API Contract 8.
- User Story 5 → API Contract 10.
- User Story 6 → Build Order 1.
- Core Feature 1 → API Contracts 1, 2 and 3.
- Core Feature 2 → API Contracts 4 and 5, Invariants 1 and 2.
- Core Feature 3 → API Contracts 5 and 6.
- Core Feature 4 → API Contracts 7 and 9.
- Core Feature 5 → API Contract 8, Invariant 8.
- Core Feature 6 → API Contract 10, Invariant 9.
- Core Feature 7 → Build Order 1.
- Success Metric 1 → API Contracts 3 and 5, Surface Transcript 1.
- Success Metric 2 → API Contract 4, Invariant 4.
- Success Metric 3 → Invariants 1 and 2.
- Success Metric 4 → API Contract 6.
- Success Metric 5 → API Contracts 7, 8 and 9, Surface Transcripts 2 and 3.
- Success Metric 6 → API Contract 10, Invariant 7, Surface Transcript 4.

## Integration Points

- Git, through the existing helpers: `git ls-files` (inside
  `spec.HistoryCitations`), `rev-parse --is-inside-work-tree`, `status
  --porcelain` over relinked paths, and in the queue `show <parent>:<path>`,
  `ls-tree` and `diff --name-only`.
- The update keeps its existing tag and tree reads; the tag tree it already
  lists feeds `StaleHistoryCitations`.
- No network, provider or Secondbrain access.

## Testing Approach

- `internal/spec/live_citations_test.go` (task_02, new): table tests of
  `IsLiveDocument` and `RelinkContent` on synthetic bytes (inline link, angle
  link, reference definition, fragment, code span and fence path, prose path,
  a near-miss prefix that must not change, a folder link), and
  `PlanArchiveRelink` on a temporary Git repository.
- `internal/cli/archive_relink_test.go` (task_02, new): command tests on the
  existing archive fixtures (`newImplementWorkspace`, `writeArchiveQAReport`,
  `commitArchiveFixture`), asserting Surface Transcript 1 whole, the refusals
  and the `--plan` lines. Existing archive tests are characterization and pass
  unchanged.
- `internal/cli/deliver_archive_relink_test.go` (task_03, new): extends the
  fixtures of `deliver_archive_record_test.go` with an inbound citation and
  asserts the staged commit and the proof, including a tampered relinked file.
- `internal/spec/stale_citations_test.go` and
  `internal/cli/baseline_update_citations_test.go` (task_04, new): the update
  tests reuse `newBaselineUpdateRepository` with legacy history committed under
  `docs/history/`, a temporary home and `--no-skills`, and assert Surface
  Transcripts 2 to 4 whole. Existing `baseline update` and `history` tests pass
  unchanged.
- `internal/cli` already installs `suiteguard.Main`; no new spawning package.

## Build Order

1. `CONTEXT.md` terms, the archive, baseline and deliver command references,
   the Roundfix Skill, the archive skill and the four authoring skills. They
   describe ADR-0260 and this TechSpec, so every code step names its guide.
2. The Archive Relink in `internal/spec` and the archive command's output
   (depends on: 1, which documents it).
3. The Delivery Queue's staging and proof of relinked files (depends on: 2,
   because it calls `RelinkContent` and `ArchiveCitationMoves` and reads
   `relinked`).
4. History Citations, `--promote` and Stale History Citations in the update
   (depends on: 2, because it calls `IsLiveDocument` from the file step 2
   creates).
5. QA gate (depends on: 3, 4).

No other active Spec is a prerequisite.

## Risks & Considerations

- A promoted file's own relative links are not rebased, ADR-0247: "Promoted
  copies must equal removed blobs". An adopted source rarely carries them; a rebase is a follow-up.
- An archive inside a Run worktree rewrites active Specs that another Run may
  be editing. The Delivery Queue archives one item at a time on its branch, and
  the dirty-file refusal keeps uncommitted edits out of the archive change.
- `RelinkContent` replaces path text in fences, Verification commands
  included, which is what keeps a moved path valid; the bounded match keeps a
  longer path that shares the prefix intact.
- The stale scan reads every tracked Markdown file once per update when the
  tag exists; a set lookup per candidate keeps it linear.
- Basename collisions in `docs/references/` refuse rather than rename, so an
  operator decides the name.

## Glossary

- adds: **Live Document**
- adds: **Archive Relink**
- adds: **History Citation**
- adds: **Stale History Citation**
- changes: **Archive Record**
- changes: **Archive Command**
- changes: **Managed Refresh**

## Decisions

- The shared classifier and rewrite live in one new `internal/spec` file; see
  ADR-0260.
- Folder and core-artifact citations target the Archive Record; other cited
  files are promoted; see ADR-0260.
- The queue proves each relinked file by recomputation from the parent; see
  ADR-0260.
- Outside a Git work tree the archive relinks nothing, which keeps every
  existing library test; see ADR-0260.
- The update's new digest fields are omitted when empty, so existing digests
  hold; see ADR-0260.
