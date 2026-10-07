---
spec: 0243-a-history-that-holds-only-records
prd: _prd.md
created: 2026-10-07
---

# A history that holds only records — Technical Spec

## Executive Summary

The conversion lives in `internal/spec`, beside the archive and its record.
A Legacy Archive Folder is planned and converted by new functions in
`internal/spec/history_sanitize.go`, which call the same
`BuildArchiveRecord` and `RenderArchiveRecord` that `roundfix archive` uses.
The other history kinds are planned and reduced or removed by
`internal/spec/history_entries.go`. The command, `roundfix history sanitize`
in `internal/cli/history.go`, orders the units, checks the History Full Tag
and the working tree, prints the plan and applies one batch. It reuses
`judge.AdviseArchive` for `--advise`.

Spec 0242 already moved every reader to the record or to Git. An authoring
ablation that converted the whole history of a clone left both repository
gates at exit 0, so this Spec fixes no reader. It changes the builder only to
accept folders that predate the archive stamp and to name a folder without QA
`no-qa`.

The primary trade-off is reviewability against speed. Six Pull Requests
replace one, and each stays a diff of deletions plus about forty small
records that a person can read and revert (ADR-0248).

## Project Constraints

- Identifier strategy: applicable — new names are the command `history
  sanitize`, its flags `--apply`, `--batch`, `--advise` and `--promote`, the
  tag `history-full`, the disposition `no-qa`, the provenance line of a
  Reduced History Entry, and Go identifiers. The record schema stays
  `roundfix/archive-record/v1`. Spec, Task and ADR identifiers keep their
  schemes. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — only `--advise` can send a request,
  through `internal/judge`'s existing client, key variables, monthly ceiling
  and judge log, by calling `judge.AdviseArchive` once per folder of the
  batch. With no key or a reached ceiling it sends nothing and prints the
  skip reason. Every test injects an `http.RoundTripper` and a temporary
  home. Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0248 (this Spec) decides the
  batched sanitize after a `history-full` tag, the conversion through the
  archive's builder with `no-qa`, the reduction of retired Findings and
  Backlog Entries, the removal of retired Review Artifacts and handoffs, and
  that legacy readers stay. ADR-0247 introduced the record and expects this
  migration, ADR-0247: "Folders archived before this decision stay readable
  until a separate, batched migration replaces them"; its advice "never
  refuses, never moves a file and never decides an archive", and `--advise`
  keeps that. Each batch's Pull Request is the later change whose rule is
  ADR-0215: "A removal is its own later change, and the maintainer's
  explicit approval of that removal is recorded in it", and it carries the
  approval recorded here. ADR-0120 keeps the History Root,
  ADR-0120: "retired documentation is documentation". ADR-0230's link pass
  and ADR-0121's ledger stay for adopters. ADR-0154 binds an override's
  fields, which the converted record carries. ADR-0232, ADR-0227 and
  ADR-0193 read archived Specs through the record, unchanged. ADR-0165,
  ADR-0169, ADR-0229 and ADR-0237 govern the Delivery Queue's parks, review
  diff, operator archives and retries around an archive; a batch runs
  outside the queue and changes none of them, so none applies. ADR-0179
  bounds the Governed Paths, and ADR-0187, ADR-0189 and ADR-0233 bind the
  owned skill's version. ADR-0184 binds the Surface Transcripts, ADR-0182
  each Task's Verification, and ADR-0210 hashes declared inputs only. The
  gate is bound by ADR-0080, ADR-0091, ADR-0096, ADR-0097, ADR-0104 and
  ADR-0167, ADR-0196 validates a pre-PR finding, and ADR-0240 decides a QA
  partial. ADR-0194 and ADR-0195 cite ADR-0097 but decide what a QA row
  records and when it is observed again; this Spec changes neither, so
  neither applies. ADR-0093, ADR-0117, ADR-0156, ADR-0168, ADR-0176 and
  ADR-0183 check consistency by citation and receipt. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — 8 of the declared files are Governed
  Paths. They were measured with `GovernedPath` through a `go test -overlay`
  probe that wrote nothing, and each is bounded in `_authorization.md`. The
  express authorization comes from the maintainer: "Concedo" on 2026-10-06,
  beside the standing grants of 2026-09-30 ("considere autorizado a ajustar
  todas as skills se necessário"; "Autorizar os dois"). Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0243-a-history-that-holds-only-records/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/archive.md`, `docs/agents/docs-layout.md`, `docs/agents/setup-context.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`, `internal/baseline/assets/modules/context-workflow.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `skills/roundfix/SKILL.md`.

## Measured inventory

Measured on 2026-10-07 at `b07a080e` with `git ls-files` per family:

| Family under `docs/history` | Files | Bytes | Disposition |
| --- | ---: | ---: | --- |
| `specs/`: 223 Legacy Archive Folders | 5,790 | 58,211,598 | one Archive Record each |
| `specs/`: Archive Records | 1 | 1,892 | kept |
| `reviews/`: 50 `pr-<n>` folders | 501 | 1,619,231 | removed |
| `handoffs/` | 11 | 83,508 | removed |
| `backlog/` (terminal) | 31 | 84,742 | Reduced History Entries |
| `findings/` (terminal) | 98 | 609,053 | Reduced History Entries |
| `adr/` (retired ADRs) | 30 | 43,109 | kept whole |

The builder was run over every folder through a `go test -overlay` probe in
a scratch clone:

- 221 folders built, rendered and round-tripped: 166 `pass`, 6 `partial`,
  39 `qa-override`, 10 `superseded`.
- Two folders, `0003-dogfood-polish` and `0010-run-robustness`, have no QA
  Report, override or supersession and failed to parse with an empty
  disposition. They get `no-qa`.
- 13 folders (0001-0010, 0033-0035) predate the archive stamp. Their
  `_prd.md` still says `status: active`, and today's `ReadArchivedSpec`
  refuses them. The conversion does not require the stamp.
- The records total 263,305 bytes. The largest is 2,041 bytes, so no record
  exceeds `ArchiveRecordTargetBytes`.
- The delivery rule of Invariant 5 finds a delivery commit for 193 folders,
  188 of them with a `(#<n>)` Pull Request. 30 folders, delivered inside
  older multi-Spec commits and later relocated, have none.
- 97 of the 98 terminal Findings carry `absorbed_by`, naming 38 distinct
  Spec slugs, and the other carries closure fields. A slug resolves through
  its Archive Record after conversion, so no `absorbed_by` needs
  repointing.
- Reducing the 129 Findings and Backlog Entries took them from 691,114 to
  94,873 bytes.
- Markdown outside the History Root cites folder paths in two places:
  `CHANGELOG.md` lines 188 and 198 name the 0217 and 0218 measurements,
  which now live in `docs/references/`.

The ablation then committed the whole conversion in the clone. `make -k
verify` and `make -k verify-docs` both exited 0 there.

## System Architecture

| Component | Where | Change |
| --- | --- | --- |
| Record builder | `internal/spec/archive_record.go` | `no-qa` disposition; parse accepts it |
| Legacy conversion | new `internal/spec/history_sanitize.go` | Lists folders, finds the delivery, plans and applies one folder |
| Other kinds | new `internal/spec/history_entries.go` | Plans, reduces and removes; scans citations |
| History Sanitize Command | new `internal/cli/history.go`, `internal/cli/cli.go` | Unit order, tag and tree checks, plan, advice, apply |
| Guidance | command guide, Roundfix skill, one Baseline clause, `CONTEXT.md`, `CHANGELOG.md` | Describe the command and the procedure |
| Export contract | new `internal/docscontract/secondbrain_export_test.go` | Ties `.secondbrain-export` to Legacy Archive Folders |

## Implementation Design

### Interfaces

```go
// internal/spec/archive_record.go
const ArchiveNoQA ArchiveDisposition = "no-qa"

// internal/spec/history_sanitize.go
type LegacyDelivery struct{ Commit, PullRequest, Date string } // Date is YYYY-MM-DD
// FindLegacyDelivery reads first-parent history of HEAD in repoRoot.
// specRoot and archiveRoot are repository-relative slash paths.
func FindLegacyDelivery(ctx context.Context, repoRoot, specRoot, archiveRoot, slug string) (LegacyDelivery, error)
// LegacyArchiveFolders lists the directories under archiveRoot that hold _prd.md, sorted.
func LegacyArchiveFolders(archiveRoot string) ([]string, error)

type LegacyConversionRequest struct {
	RepositoryRoot string         // absolute
	ArchiveRoot    string         // absolute, inside RepositoryRoot
	Slug           string
	SourceRevision string         // 40- or 64-hex HEAD
	Delivery       LegacyDelivery
	Promote        []string       // repository-relative files inside the folder
}
type LegacyConversion struct {
	Slug, Folder, RecordPath string // repository-relative slash paths
	Record                   ArchiveRecord
	Rendered                 []byte
	Files                    []string // folder-relative, sorted
	Bytes                    int64
	Promoted                 []string // docs/references/<basename>
}
func PlanLegacyConversion(req LegacyConversionRequest) (LegacyConversion, error) // writes nothing
func ApplyLegacyConversion(repositoryRoot string, c LegacyConversion) error

// internal/spec/history_entries.go
type HistoryKindAction string // "reduce", "remove"
type HistoryKindPlan struct {
	Kind                    ArchiveKind // findings, backlog, reviews, handoffs
	Action                  HistoryKindAction
	Files                   []string // repository-relative slash paths, sorted
	BytesBefore, BytesAfter int64
	reduced                 map[string][]byte
}
func PlanHistoryKinds(repositoryRoot, revision string) ([]HistoryKindPlan, error)
func ApplyHistoryKind(repositoryRoot string, plan HistoryKindPlan) error
func ReduceHistoryEntry(content []byte, revision, path string) ([]byte, error)
func IsReducedHistoryEntry(content []byte) bool
type HistoryCitation struct{ Path string; Line int; Target string }
func HistoryCitations(ctx context.Context, repositoryRoot string, removed []string) ([]HistoryCitation, error)
```

### Invariants

1. A Legacy Archive Folder is a directory directly under the resolved
   archive root that holds `_prd.md`. An Archive Record `<slug>.md` is never
   one. The PRD's `status` does not matter.
2. Units are ordered: every Legacy Archive Folder by slug, then the kinds
   `findings`, `backlog`, `reviews` and `handoffs`. A kind is a unit only
   while it has a pending file. A file under `findings/` or `backlog/` is
   pending unless `IsReducedHistoryEntry` holds. Every file under
   `reviews/` and `handoffs/` is pending. `adr/` is never a unit. `--batch
   <n>` selects the first `n` pending units.
3. `BuildArchiveRecord` sets `no-qa` when no other rule set a disposition:
   no QA Report, no override and no supersession. `ParseArchiveRecord`
   accepts `no-qa`. `ArchivedTaskCompleted` treats it like `pass`. No other
   disposition rule changes.
4. `PlanLegacyConversion` calls `BuildArchiveRecord` with `SpecDir` the
   folder, `Source` the folder's repository-relative path (for example
   `docs/history/specs/<slug>`) and `SourceRevision` the request's. It then
   sets `PullRequest` and `DeliveryCommit` from `Delivery`, and `Archived`
   from `Delivery.Date` when the PRD carries no `archived` stamp. It renders
   and parses the record back and refuses when the round trip differs. It
   refuses a promotion that is outside the folder, not a regular file, a
   core artifact (`_*.md` at the folder root, `task_*.md`, a QA Report), a
   duplicate destination, or a destination that already exists under
   `docs/references/`. It writes nothing.
5. `FindLegacyDelivery` returns the newest first-parent commit of `HEAD`
   that deleted `<specRoot>/<slug>/_prd.md`, with rename detection off.
   Without one, it returns the newest first-parent commit that added
   `<archiveRoot>/<slug>/_prd.md`, unless that same commit deleted a
   `_prd.md` of the same slug elsewhere, which makes it a relocation. Without
   either, every field is empty. `PullRequest` is the decimal `<n>` of a
   subject that ends with `(#<n>)`. `Date` is the commit's author date.
6. `ApplyLegacyConversion` writes the record, copies each promotion
   byte-identically, and only then removes the folder. A failure before the
   removal removes what it wrote and leaves the folder whole.
7. `ReduceHistoryEntry` keeps the front matter byte for byte, including both
   `---` lines, then writes a blank line, the first `# ` line, a blank line,
   the first paragraph after that title that is not a heading with its
   whitespace collapsed, a blank line, and
   ``Full text in Git at `<revision>`: `<path>`.``. It refuses a file without
   front matter or without a title. `IsReducedHistoryEntry` holds exactly for
   a file whose last non-empty line has that form with a 40- or 64-hex
   revision. Reducing a reduced entry is refused.
8. `ApplyHistoryKind` rewrites each `reduce` file with its planned bytes and
   deletes each `remove` file, then removes directories left empty under
   that kind. It touches no other path.
9. `HistoryCitations` reads the tracked `*.md` files outside the History
   Root (`git ls-files`) and reports each line that names a removed path or
   a removed folder prefix. It is reported and never refuses.
10. The command reads only the repository's own archive root. It refuses,
    exit 2, when the resolved Spec Root is external.
11. Without `--apply`, the command writes nothing under the repository.
    With `--advise`, its only write is the judge log under the Roundfix
    Home. It calls `judge.AdviseArchive` once per folder of the batch with
    that folder's candidate files: every file outside the core artifacts and
    `qa/evidence/`, as `roundfix archive <slug> --plan` groups them.
12. `--apply` runs these checks before any write, in order, and each refusal
    exits 2 through Preflight Validation:
    - `--batch <n>` is given and `n` is at least 1;
    - `git status --porcelain --untracked-files=all` is empty;
    - `refs/tags/history-full` exists and `git cat-file -t` names a `tag`;
    - `git merge-base --is-ancestor history-full^{commit} HEAD` exits 0;
    - `git ls-tree -r --name-only history-full` holds every path the batch
      removes or rewrites;
    - every unit of the batch plans without error.
13. `--apply` then applies the units in order with `SourceRevision` and the
    reduction revision set to `HEAD`. A write failure exits 1 and names the
    unit; the operator restores the tree with `git restore` and `git clean`.
    The command never commits, tags, pushes or opens a Pull Request.
14. Usage errors exit 2: `--advise` without `--batch` or with `--apply`,
    `--promote` without `--apply`, a `--batch` that is not a positive
    integer, an unknown flag or an extra argument.

### Data Models

A converted record, as `--apply` writes it:

```markdown
---
schema: roundfix/archive-record/v1
spec: 0160-a-review-that-reaches-a-verdict
title: A review that reaches a verdict
status: archived
created: "2026-09-18"
archived: "2026-09-19"
disposition: pass
source: docs/history/specs/0160-a-review-that-reaches-a-verdict
source_revision: <HEAD at apply, 40-hex>
qa_task: task_05
qa_report: qa-report-2026-09-19.md
qa_verdict: pass
unproven: []
adrs: [...]
sources: [...]
regeneration: []
promoted: []
pull_request: "244"
delivery_commit: 51232b1a<...40-hex>
---

# A review that reaches a verdict

<one paragraph>
```

A Reduced History Entry:

```markdown
---
<front matter, unchanged>
---

# <title, unchanged>

<first paragraph>

Full text in Git at `<HEAD at apply>`: `docs/history/findings/<name>.md`.
```

No Run Database field changes.

### API Contracts

1. API Contract: `roundfix history sanitize [--batch <n>]` prints the plan
   to stdout and exits 0. It changes no file in the repository.
2. API Contract: `roundfix history sanitize --batch <n> --advise` adds Jev's
   advice per candidate file of the batch's folders and exits 0. It fails
   open.
3. API Contract: `roundfix history sanitize --apply --batch <n> [--promote
   <path> ...]` applies the next `n` units, prints one confirmation line and
   exits 0. It exits 2 for each refusal of Invariant 12 and 1 for a write
   failure.
4. API Contract: when no unit is pending, the plan says so and `--apply`
   exits 0 without writing.

### Surface Transcripts

1. Surface Transcript: the plan.

   ```transcript
   $ roundfix history sanitize
   stdout:
   history sanitize plan: <u> unit(s) pending; <f> file(s) (<b> bytes) leave docs/history
   folder docs/history/specs/<slug>: removes <n> file(s) (<b> bytes) and writes docs/history/specs/<slug>.md (<r> bytes, <disposition>, delivery <12-hex> #<pr>)
   candidate docs/history/specs/<slug>/references/<file>.md <bytes> bytes
   findings: reduces <n> file(s) from <b> to <a> bytes
   reviews: removes <n> file(s) (<b> bytes)
   cites CHANGELOG.md:<line> names docs/history/specs/<slug>/<path>
   apply with: roundfix history sanitize --apply --batch <n> (needs the annotated tag history-full at or before HEAD)
   stderr:
   exit: 0
   ```

2. Surface Transcript: a batch.

   ```transcript
   $ roundfix history sanitize --apply --batch 2 --promote docs/history/specs/<slug>/references/<file>.md
   stdout:
   history sanitize applied 2 unit(s): wrote 2 Archive Record(s), reduced 0 file(s), removed <n> file(s) (<b> bytes) kept in Git at <12-hex> and tag history-full; promoted 1 file(s) to docs/references/; <r> unit(s) remain
   stderr:
   exit: 0
   ```

3. Surface Transcript: a batch without the tag.

   ```transcript
   $ roundfix history sanitize --apply --batch 2
   stdout:
   stderr:
   Preflight failed
   ...
   Reason:
     history sanitize --apply needs the annotated tag history-full at or before HEAD; tag the last commit before the first batch with git tag -a history-full
   ...
   exit: 2
   ```

4. Surface Transcript: advice without a key.

   ```transcript
   $ roundfix history sanitize --batch 1 --advise
   stdout:
   history sanitize plan: 1 unit(s) pending; <f> file(s) (<b> bytes) leave docs/history
   folder docs/history/specs/<slug>: removes <n> file(s) (<b> bytes) and writes docs/history/specs/<slug>.md (<r> bytes, <disposition>, delivery unknown)
   candidate docs/history/specs/<slug>/references/<file>.md <bytes> bytes: no advice (<KEY_VARIABLE> is not set)
   apply with: roundfix history sanitize --apply --batch <n> (needs the annotated tag history-full at or before HEAD)
   stderr:
   exit: 0
   ```

## Coverage Map

- Goal 1 → plan (Invariants 1, 2, 9, 11; API Contracts 1, 4; Surface Transcript 1)
- Goal 2 → batch and tag (Invariants 2, 12, 13; API Contract 3; Surface Transcripts 2, 3)
- Goal 3 → legacy conversion (Invariants 3-6)
- Goal 4 → other kinds (Invariants 7, 8)
- Goal 5 → guidance and export contract (Exact texts; Operator batch procedure)
- User Story 1 → conversion and reduction
- User Story 2 → plan, `--advise`, `--promote` (Surface Transcript 4)
- User Story 3 → Invariant 12
- User Story 4 → no reader removed (ADR-0248)
- Core Feature 1 → Invariants 2, 9, 11
- Core Feature 2 → Invariants 12, 13
- Core Feature 3 → Invariants 3-5
- Core Feature 4 → Invariants 7, 8
- Core Feature 5 → Exact texts; Operator batch procedure
- Success Metric 1 → task_03 tests
- Success Metric 2 → task_01 and task_03 tests
- Success Metric 3 → task_03 tests
- Success Metric 4 → QA ablation row
- Success Metric 5 → task_04 tests

## Integration Points

- Git, through the local `git` binary: `rev-parse`, `status --porcelain`,
  `cat-file -t`, `merge-base --is-ancestor`, `ls-tree -r`, `log
  --first-parent`, `diff-tree`, `ls-files`.
- The Spec judge's existing transports, for `--advise` only.

## Testing Approach

- `internal/spec/history_sanitize_test.go` builds temporary repositories
  with Legacy Archive Folders: a passing Spec, an override, a superseded
  Spec, a pre-stamp Spec and one without a QA Report. It covers the size
  bound, the round trip, `git show <source_revision>:<source>/<path>`, each
  delivery rule including a relocation, each promotion refusal, and a
  failure before removal.
- `internal/spec/history_entries_test.go` covers the reduction, its
  idempotence and refusals, the kind plan's order and byte counts, the
  removal of reviews and handoffs, and the citation scan.
  `internal/speccheck/reduced_history_entry_test.go` runs the Spec check over
  reduced Findings with `absorbed_by` and a Rollup member, and over a reduced
  entry whose license is broken, which must still fail.
- `internal/cli/history_test.go` runs the command in temporary repositories
  with a fake judge transport and a temporary home, and asserts Surface
  Transcripts 1-4, every refusal of Invariants 12 and 14, and that a plan
  leaves the tree byte-identical.
- `internal/docscontract/secondbrain_export_test.go` asserts the export rule
  over this repository and over fixtures in both directions.
- `internal/baseline/history_sanitize_clause_test.go` asserts the appended
  sentence in the embedded module, the formatter golden and this
  repository's generated guide.

## Build Order

1. The `no-qa` disposition and the legacy conversion with its delivery rule
   and promotions.
2. The other kinds: reduction, removal and the citation scan. This step
   shares no file with step 1.
3. The History Sanitize Command, its command guide and its index row
   (depends on: 1, 2).
4. The Roundfix skill reference and versions, the Baseline sentence with
   `make baseline-digests` and the Managed Refresh, the glossary, the
   `CHANGELOG.md` paths and the export contract. This step shares no file
   with steps 1-3.
5. Final QA gate (depends on: 3, 4).

## Operator batch procedure

This runs after the Spec merges, never inside a Task. The command guide
carries it.

1. Build `bin/roundfix` from `main`. On a clean `main`, create and push the
   tag: `git tag -a history-full -m "docs/history before the sanitize
   batches (ADR-0248)"` and `git push origin history-full`.
2. Run `roundfix history sanitize` and keep its plan. Run `--batch 40
   --advise` for each batch when Jev's monthly ceiling allows. Its ceiling
   was reached on 2026-10-02, so the advice is skipped until the month
   turns, and the Jev judgment of 2026-10-06 stands in.
3. For each batch, on a branch `chore/history-sanitize-<k>` from `main`, run
   `roundfix history sanitize --apply --batch 40`, then `make verify` and
   `make verify-docs`. Commit `chore: sanitize history batch <k>`, open one
   Pull Request, merge it, and update `main`. A batch is reverted by
   reverting its commit.
4. Promote in the batch that holds the file:
   - with 0129, `--promote docs/history/specs/0129-spec-authoring-and-gate-recovery/references/2026-08-12-a-queue-of-eight-specs-shows-where-the-loop-breaks.md`;
   - with 0231, `--promote docs/history/specs/0231-checks-that-hold-in-delivery/references/2026-10-05-a-pull-request-check-ran-on-a-stale-merge.md`.

   Send a cross-project lesson to the Secondbrain inbox under
   `docs/agents/secondbrain.md`. The 0035 skill analysis and the 0079 pilot
   report are already there.
5. Six batches of 40 cover the 223 folders. The sixth also takes the four
   kind units. That last batch also drops the `docs/history` exclusions from
   `.secondbrain-export`, which the export contract then requires.
6. `roundfix history sanitize` reports `0 unit(s) pending` at the end.

## Risks & Considerations

- A removed file is recoverable only from Git. The tag and each record's
  `source_revision` name where. A batch branch starts from `main`, so `HEAD`
  at apply is a `main` commit that a squash merge keeps reachable.
- The Secondbrain mirror drops the removed paths on the next sync. Wiki
  pages that cite `projects/roundfix/mirror/docs/history/...` paths then
  point at files the mirror no longer holds. The Secondbrain's own lint owns
  that.
- The 30 records without a delivery commit carry no `pull_request`. No
  reader requires it; reconcile names a merged head without it.
- A batch that fails mid-write leaves a partial tree. The command names the
  unit, and the clean-tree refusal stops the next run until the operator
  restores the tree.
- Thirteen folders have `status: active` in a PRD under the History Root.
  Converting them ends the refusal that `ReadArchivedSpec` gives today.
- The `history-full` tag must be pushed, or a clone without it cannot
  apply. The command checks the local tag only.

## Exact texts

The Baseline sentence is appended after the last sentence of
`clause.context.docs-one-job-per-directory`. Every existing sentence stays
byte-identical:

"When the runtime offers a history sanitize, such as `roundfix history
sanitize --apply --batch <n>`, it replaces each Spec folder an earlier
archive left under the history root with its Archive Record, reduces each
retired Finding and Backlog Entry to its front matter, title, first paragraph
and the revision that holds its full text, and removes retired Review
Artifacts and handoffs, one reviewed batch at a time after a tag marks the
full history; retired ADRs stay whole."

The Roundfix skill's archive reference gains `## History Sanitize Command`
after `## Archive Command`. It names `roundfix history sanitize`, `--apply
--batch <n>`, `--advise`, `--promote`, the `history-full` tag and the
`no-qa` disposition. The skill index row for `archive` adds `history` to its
commands.

## Glossary

- adds: **History Sanitize Command**
- adds: **Legacy Archive Folder**
- adds: **Sanitize Batch**
- adds: **Reduced History Entry**
- adds: **History Full Tag**
- changes: **Archive Record**
- changes: **History Root**

## Research basis

- Secondbrain: `wiki/index.md` was read, then
  `qmd query "roundfix history archive record sanitize existing folders" --all --files --min-score 0.3`.
  It returned the adopted Backlog Entry mirrored with Spec 0242 (score
  1.00), the 0226 Finding on relative links in archived Specs, and
  `wiki/sources/digest-roundfix-autorizacoes-historicas-e-propriedade-da-regeneracion-2026-09-09.md`.
  A grep of the inbox found the 0035 skill analysis and the 0079 pilot
  report triaged on 2026-10-06 under `inbox/secondbrain/_triaged/`, so the
  procedure promotes only the 0129 and 0231 files. The mirror's history
  paths are cited by wiki pages, which is the second risk above.
- Exa: Git's `git tag` page
  (<https://git-scm.com/docs/git-tag>) separates annotated tags, "meant for
  release", from lightweight ones "for private or temporary object labels",
  which is why `--apply` refuses a lightweight tag. Git's `git merge-base`
  page (<https://git-scm.com/docs/git-merge-base>) defines
  `--is-ancestor`, which the ancestry refusal uses.
- Spec 0242's Archive Record, PRD and TechSpec at
  `3afb60ad06ae829ed51380a0a42265e04698fffa` supplied the list this Spec
  answers and the builder it reuses.

## Decisions

- Batched sanitize after an annotated `history-full` tag, dry run by
  default, refusing a dirty tree. See ADR-0248.
- Conversion through `BuildArchiveRecord` with `source` the folder's own
  path and `source_revision` `HEAD`, and the `no-qa` disposition. See
  ADR-0248.
- Findings and Backlog Entries reduced, Review Artifacts and handoffs
  removed, ADRs kept, and legacy readers kept for adopters. See ADR-0248.
- Advice through ADR-0247's `AdviseArchive`, opt-in per batch. See ADR-0247.
