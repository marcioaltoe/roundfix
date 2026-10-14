---
spec: 0271-history-sanitize-records-keep-their-summary-links-and-dates
prd: _prd.md
created: 2026-10-13
---

# History sanitize records keep their summary, links and dates — Technical Spec

## Executive Summary

Every defect lives in one generator: `BuildArchiveRecord` as
`PlanLegacyConversion` calls it with `Legacy: true`. This design adds four
readings behind that flag in `internal/spec/archive_record.go` (summary,
override fields, sources, TechSpec date), one Git lookup for the creation
commit in `internal/spec/history_sanitize.go`, and one new file,
`internal/spec/legacy_link_rebase.go`, that rebases the relative links of the
record's outcome and of each promoted copy. The two commands that plan
conversions, `history sanitize` and the history section of `baseline update`,
pass the creation date through. The help, the output lines and the record
schema do not change.

The trade-off accepted is ADR-0277's split: legacy conversions now rewrite
promoted copies and outcomes while the Archive Command keeps copying both byte
for byte, so the two writers of Archive Records differ until the sibling
Backlog Entry is built. In exchange, the Delivery Queue's exact-archive proof
and every record already written stay untouched.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or Archive Record
  field is added. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request, credential or network
  call; Git is read locally. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0277 governs every decision below.
  ADR-0251: "The rule applies wherever a Legacy Archive Folder becomes a
  record". ADR-0248: "Each folder becomes `<slug>.md` through
  `BuildArchiveRecord`, the builder the archive itself uses". ADR-0247:
  "Promoted copies must equal removed blobs". ADR-0254: "One plan, one
  digest". ADR-0260: "a promoted file keeps its bytes and its own outward
  relative links are not rebased", narrowed by ADR-0277 for conversions.
  ADR-0261: "records already archived stay byte-identical". ADR-0262: "A
  cancelled Spec leaves an Archive Record and returns its sources". ADR-0269:
  "Archive and cancel refuse ignored Spec files and finish a failed removal
  forward". ADR-0256: "A
  release waits for the skills that describe a changed surface". ADR-0244: "A
  Spec declares the domain terms it introduces". ADR-0257: "Inside a Run the
  Daemon is the only full gate". ADR-0258: "A Spec is authored against the QA
  rerun classes". ADR-0259: "A test runs in parallel unless it names why it
  cannot". ADR-0271: "A test bounds by real time only what it asserts".
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AG grant of 2026-10-13 under the cycle AE "Concedo" and the 2026-09-30
  skills grant ("considere autorizado a ajustar todas as skills se
  necessário"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0271-history-sanitize-records-keep-their-summary-links-and-dates/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`,
  `.agents/skills/roundfix/references/baseline.md`,
  `skills/roundfix/SKILL.md`.

## Current behavior

Characterized on 2026-10-13 with `bin/roundfix` built at `9963ffd3`, in a
temporary repository and home holding three synthetic Legacy Archive Folders
and an annotated `history-full` tag, by
`roundfix history sanitize --apply --batch 3 --promote <note> --promote <probe>`
(exit 0):

- `0001-demo`, whose `_prd.md` has `qa_override: true`, no `created`, and
  `# Demo` followed directly by `## Problem`, wrote `created: ""`,
  `qa_override_approval: ""`, `qa_override_reason: ""`, `sources: []` beside
  `promoted: [docs/references/2026-01-01-note.md]` (the note has Backlog Entry
  front matter, `status: promoted`, and no `_index.md` exists), and a body
  with the title and no summary. Its `_techspec.md` held `created: 2026-01-02`.
- `0002-quote`, whose summary is the two-line block quote
  `> Archived by maintainer decision` / `> without a QA verdict.`, wrote the
  body line `> Archived by maintainer decision > without a QA verdict.`
- `0003-link`, whose summary is `The measurement is in [the probe](references/probe.md).`,
  wrote that line verbatim, so from `docs/history/specs/` it resolves to a
  missing `docs/history/specs/references/probe.md` while the file sits at
  `docs/references/probe.md`.
- The promoted note was copied byte for byte: `[the PRD](../_prd.md)` now
  resolves to `docs/_prd.md`, `[ADR one](../../../../adr/0001-x.md)` to a path
  above the repository, and the text `docs/history/specs/0001-demo/` names the
  removed folder.

A read-only scratch clone of the adopter's repository at its `history-full`
tag showed the same six shapes in the records its batches wrote, matching the
Backlog Entry's identifiers. Five of its folders without `_index.md` hold
adopted Findings (front matter `status` and `created_at`) or Backlog Entries
(`type`, `status`, `created`) directly under `references/`.

## System Architecture

`internal/spec` owns the generator; `internal/cli` only plans and prints.

- `BuildArchiveRecord` (`archive_record.go`) gains the four legacy readings
  behind `in.Legacy`, plus an optional `Rebase` input applied to the outcome.
  The non-legacy path and `ParseArchiveRecord` are untouched.
- `PlanLegacyConversion` (`history_sanitize.go`) builds the `LegacyLinkRebase`
  from its folder, record path and promotions, rebases each promoted copy and
  passes the rebase to `BuildArchiveRecord`. It applies
  `LegacyConversionRequest.Created` when the folder supplied no date.
- `FindLegacyCreation` (`history_sanitize.go`) reads the creation commit as
  `FindLegacyDelivery` reads the delivery commit.
- `legacy_link_rebase.go` holds the rebase. It reuses the Archive Relink's
  scanner (`scanArchiveMarkdown`, `archiveLinkTarget`) and path pass
  (`ArchiveCitationMoves`, `RelinkContent`'s text pass), so a code span or
  fence keeps its destinations as the archive does.
- `planHistoryUnit` (`internal/cli/history.go`), shared by both commands,
  calls `FindLegacyCreation`; the two re-plans with promotions
  (`history.go` and `baseline_update_history.go`) carry `Created` from the
  first plan.

## Implementation Design

### Interfaces

```go
// internal/spec/archive_record.go
const LegacyNotRecorded = "legacy: not recorded"

type ArchiveRecordInput struct {
	// existing fields unchanged
	Rebase *LegacyLinkRebase // legacy only; nil leaves the outcome as read
}

// internal/spec/history_sanitize.go
type LegacyConversionRequest struct {
	// existing fields unchanged
	Created string // fallback YYYY-MM-DD when the folder holds none
}

func FindLegacyCreation(ctx context.Context, repoRoot, specRoot, archiveRoot, slug string) (string, error)

// internal/spec/legacy_link_rebase.go
type LegacyLinkRebase struct {
	Folder, RecordPath string            // repository-relative, slash-separated
	Files              []string          // repository-relative folder files
	Promoted           map[string]string // folder file -> docs/references/<base>
}

func (r LegacyLinkRebase) Rebase(content []byte, from, to string) []byte
```

### Data Models

No new record field. A legacy record's `created`, `sources`,
`qa_override_approval`, `qa_override_reason` and outcome take new values; the
other override fields and every other key keep today's values.

### API Contracts

1. API Contract: summary reading. Under `Legacy`, each line of the summary
   paragraph, read from the PRD after `# ` or from the newest QA Report's
   `## Outcome`, loses its block quote markers before the lines are joined:
   up to three spaces, `>`, and one optional space, repeated while the line
   still starts with a marker. When the PRD's title paragraph is empty, the
   summary is the first paragraph after the PRD's first `## ` heading, read
   the same way. A QA Report outcome and a supersession reason keep their
   precedence.
2. API Contract: override fields. Under `Legacy`, `qa_override: true` with an
   empty or blank `qa_override_approval` records `legacy: not recorded` as the
   approval, and likewise for `qa_override_reason`. A recorded value is kept
   byte for byte. `qa_override_qa_outcome`, `qa_override_qa_task_status` and
   `qa_override_revision` keep the legacy value, empty included.
3. API Contract: sources. Under `Legacy`, when `references/_index.md` is
   absent or lists no path, `sources` names, sorted, each regular `.md` file
   directly under `references/` whose front matter is a mapping holding
   `type`, `status` and `created` (Backlog Entry) or `status` and
   `created_at` (Finding). Each value is the file name, as an index's `path`
   column writes it. An index that lists a path keeps today's reading alone.
4. API Contract: created. Under `Legacy`, an empty PRD `created` takes the
   TechSpec's `created` when its first ten characters are a `YYYY-MM-DD`
   date, recorded as those ten characters. `FindLegacyCreation` returns the
   author date (`%as`) of the oldest commit reachable from `HEAD` that added
   `<specRoot>/<slug>/_prd.md` or `<archiveRoot>/<slug>/_prd.md`
   (`git log --no-renames --diff-filter=A --reverse`), or `""` when none; a
   Git failure is returned wrapped and refuses the command, as a delivery
   read does. `PlanLegacyConversion` records `req.Created` only when the
   folder supplied no date. Both commands pass it, including in the re-plan
   that adds promotions. The folder branch of `ReadArchivedSpec` reads no Git
   and stops at the TechSpec.
5. API Contract: Legacy Link Rebase. `Rebase(content, from, to)` rewrites
   each relative destination the archive scanner finds: it resolves the
   destination against `from`'s directory to a repository path `T`. When `T`
   is the folder or inside it, the new target is `Promoted[T]` or else
   `RecordPath`; otherwise it is `T`. The destination becomes the escaped
   path from `to`'s directory to that target, with its `?` or `#` suffix. A
   destination that leaves the repository, an absolute URL, a scheme, an
   anchor alone and a destination inside code are left as they are. The text
   pass rewrites path text with `ArchiveCitationMoves(Folder, RecordPath,
   Files, Promoted)` exactly as the Archive Relink does. Content with no such
   destination or path text is returned byte-identical. `PlanLegacyConversion`
   applies it to each promoted copy (`from` the source file, `to` its
   destination) and to the outcome (`from` the PRD or QA Report that supplied
   it, `to` the record path).
6. API Contract: everything else stays. Without `Legacy`, `BuildArchiveRecord`
   returns what it returns today; the Archive Command, `ParseArchiveRecord`,
   the Plan Digest format, the command help and every output line keep their
   bytes, and no code path rewrites an existing record.

### Surface Transcripts

None. No command's help, stdout or stderr changes shape: the plan and apply
lines keep their text, and only the byte counts of changed records move.

### Invariants

```text
1. A record from a folder without any of the six shapes has the bytes the delivery base writes.
2. A promoted copy without relative destinations or folder path text is copied byte for byte.
3. Rebase changes only destinations and path text; every other byte stays.
4. The Archive Command's records and every existing record keep their bytes.
5. A recorded override approval or reason is never replaced.
6. The plan and the apply of one conversion write the same record bytes.
```

## Coverage Map

Deviations from the PRD: None.

- User Story 1 → summary reading (API Contract 1)
- User Story 2 → `LegacyLinkRebase` (API Contract 5)
- User Story 3 → override fields, sources and created (API Contracts 2-4)
- Core Feature 1 → API Contract 1
- Core Feature 2 → API Contracts 2, 3 and 4, `FindLegacyCreation`
- Core Feature 3 → `LegacyLinkRebase`, API Contract 5
- Core Feature 4 → user guide, Roundfix skill archive and baseline references, `CONTEXT.md`
- Success Metric 1 → `internal/cli/history_legacy_fields_test.go`, `internal/cli/history_link_rebase_test.go` and the QA gate
- Success Metric 2 → Invariants 1-3, `internal/spec/legacy_link_rebase_test.go`, the QA gate
- Success Metric 3 → API Contract 6, the existing archive tests and the QA gate

## Absence Evidence

1. Block quote marker removal — command: `grep -n 'TrimPrefix(line, ">")\|"> "' internal/spec/archive_record.go`; result: no match at `9963ffd3`; `archiveParagraph` only trims spaces.
2. First-section summary fallback — command: `grep -n 'archiveParagraph(' internal/spec/archive_record.go`; result: three calls besides its definition at `9963ffd3`, after `# `, in `## Outcome`, and in `ParseArchiveRecord`; none reads a `## ` section of the PRD.
3. Legacy override marker — command: `grep -rn 'legacy: not recorded' internal --include='*.go'`; result: no match at `9963ffd3`; the only other "not recorded" text in production is the unrelated `QAPreconditionReasonUnrecorded`.
4. Adopted-source fallback — command: `grep -n 'references' internal/spec/archive_record.go`; result: only the `_index.md` read at `9963ffd3`.
5. `FindLegacyCreation` and the TechSpec date — command: `grep -rn '_techspec.md\|diff-filter=A' internal/spec/archive_record.go internal/spec/history_sanitize.go`; result: one `--diff-filter=A` in `FindLegacyDelivery`, which reads the archive root for delivery, and no `_techspec.md` read.
6. `LegacyLinkRebase` — command: `grep -rn 'filepath.Rel(filepath.Dir(citingPath)' internal/spec`; result: only `RelinkContent`, which resolves a destination against the file's current location and maps only moved targets; nothing rebases a file that moves.

## ADR Alternatives

- ADR-0277: this Spec's decision; it rejects rewriting existing records, changing the Archive Command, a review-convention exception and citation-only reporting.
- ADR-0251: no rejected alternative is taken; legacy manifests are still never rewritten, and the new readings tolerate rather than refuse.
- ADR-0248: no rejected alternative is taken; batches, the tag and the one-record rule stay.
- ADR-0247: no rejected alternative is taken; the Archive Command's proof is untouched.
- ADR-0254: no rejected alternative is taken; the update still converts every unit in the change it plans.
- ADR-0260: the rejected alternatives stay rejected; ADR-0277 narrows only its statement about promoted bytes, for conversions.
- ADR-0261: no rejected alternative is taken; no backfill rewrites a record.
- ADR-0272: no rejected alternative is taken; its correction proof reads Archive Records and this Spec changes none of its rules.
- ADR-0262: no rejected alternative is taken; the Cancel Command and its record are untouched.
- ADR-0269: no rejected alternative is taken; the archive and cancel refusals and their recovery are untouched.
- ADR-0256: no rejected alternative is taken; the docs Task updates both covering skill references.
- ADR-0244: no rejected alternative is taken; the two terms follow its declaration contract.
- ADR-0257: no rejected alternative is taken; agents run focused tests and the Daemon the declared Verification.
- ADR-0258: no rejected alternative is taken; every file and break is declared.
- ADR-0259: no rejected alternative is taken; every new test calls `t.Parallel()` first.
- ADR-0271: no rejected alternative is taken; no new test waits on real time.

## Integration Points

None outside the repository. Git is read through `git log` in the repository
that holds the Spec Root, as `FindLegacyDelivery` reads it.

## Testing Approach

New tests call `t.Parallel()` first and use temporary repositories built with
`gittest` and the existing fixtures (`legacyFixture`, `newHistoryFixture`,
`baselineCitationFixture`). No existing test is expected to change: the
existing legacy fixtures carry `created`, a title paragraph, no block quote,
no override without approval and promoted files without links. A Task that
finds otherwise names the test and the contract it moves to.

- `internal/spec/legacy_record_summary_test.go`: block quote removal, nested
  quotes, the first-section fallback, QA Report precedence, and
  `TestActiveArchiveRecordKeepsItsSummaryReading` for the non-legacy path.
- `internal/spec/legacy_record_fields_test.go`: the override marker and its
  preservation, the sources fallback and index precedence, the TechSpec date,
  `FindLegacyCreation`, and `TestActiveArchiveRecordKeepsItsFields`.
- `internal/cli/history_legacy_fields_test.go`: through `history sanitize
  --apply` and `baseline update --yes --promote`, a folder dated by neither
  PRD nor TechSpec records its creation commit's date.
- `internal/spec/legacy_link_rebase_test.go` and
  `internal/cli/history_link_rebase_test.go`: the rebase rule, its byte
  preservation, and both commands writing rebased outcomes and copies.

## Build Order

1. Summary reading and its tests (no dependency).
2. Override marker, sources fallback, TechSpec and creation-commit dates,
   command wiring, the `history` guide paragraph on those fields and tests
   (depends on: 1, which edits the same function).
3. `LegacyLinkRebase`, its use in `PlanLegacyConversion` and the outcome, and
   tests (depends on: 2, which edits the same request and function).
4. User guide, Roundfix skill references, mirrors, versions, Behavior Surface
   Record and glossary (depends on: 1, 2, 3).
5. QA gate (depends on: 4).

No other Spec is a prerequisite. Rebase on main after the other cycle AG
plans merge and run the strict check again.

## Risks & Considerations

- A promoted copy that links another Legacy Archive Folder converted in a
  later batch still breaks when that folder goes; the next plan lists it as a
  History Citation.
- The summary fallback can pick a long first section; `RenderArchiveRecord`
  still shortens it to the record's target size.
- An adopter whose records were written before this version keeps them as
  they are; the user guide says so.

## Glossary

- adds: **Legacy Link Rebase**
- changes: **Lenient Legacy Reading**

## Decisions

- The readings apply only under `Legacy`. See ADR-0277.
- Approval and reason take `legacy: not recorded`; the other override fields keep their legacy value. See ADR-0277.
- A folder link without a promoted copy points at the Archive Record. See ADR-0277.
