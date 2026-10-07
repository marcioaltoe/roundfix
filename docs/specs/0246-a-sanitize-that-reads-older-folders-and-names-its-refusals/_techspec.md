---
spec: 0246-a-sanitize-that-reads-older-folders-and-names-its-refusals
prd: _prd.md
created: 2026-10-07
---

# A sanitize that reads older folders and names its refusals — Technical Spec

## Executive Summary

The History Sanitize Command plans each Legacy Archive Folder through
`PlanLegacyConversion` (`internal/spec/history_sanitize.go:181`). That calls
`BuildArchiveRecord` (`internal/spec/archive_record.go:60`), which reads the
QA Task through `ReadCauseGraph` (`archive_record.go:143`). `ReadCauseGraph`
validates the projection table with today's rules (`internal/spec/spec.go:640`,
`spec.go:822`). The builder also copies the newest QA Report's verdict into the
disposition (`archive_record.go:157`), and `ParseArchiveRecord` accepts only
five values (`archive_record.go:371`). `runHistoryCommand` turns any planning
error into a Preflight refusal of the whole command
(`internal/cli/history.go:254`).

The fix has three parts. A Legacy Archive Folder's graph is read with a
lenient manifest reader that tolerates and names projection rows outside the
graph and retired Task types. A failed QA without an override becomes the
disposition `failed-qa`. A unit that cannot be planned becomes a Refused Unit,
listed with its reason, while planning continues and a batch fills itself with
the next units it can convert.

The trade-off accepted: the plan for `--batch <n>` now depends on planning, so
the batch's tag coverage is checked after planning, not before. In exchange,
Refused Units never consume a batch, and an adopter's history converts without
editing archived bytes (ADR-0251).

## Project Constraints

- Identifier strategy: not applicable. No identifier scheme changes. The
  disposition value `failed-qa` is kebab-case under the existing `disposition`
  key, and the new stdout lines start with the lowercase words `refused` and
  `tolerates`. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable. No request or credential is added.
  The command reads local files and Git, and tests use temporary repositories
  and a temporary home. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0251 (this Spec) governs the Lenient
  Legacy Reading, `failed-qa` and Refused Units. ADR-0248: "Each folder becomes
  `<slug>.md` through `BuildArchiveRecord`, the builder the archive itself
  uses". The builder stays shared, and only the `Legacy` input changes how it
  reads the graph. ADR-0248's "A folder with neither a QA Report, an override
  nor a supersession gets the new disposition `no-qa`" is extended by
  `failed-qa`. ADR-0154: "A QA archive override records user authority, not a
  pass". A failed QA without an override is never `qa-override`. ADR-0247: "The
  folder's bytes stay in Git at the commit the record names as
  `source_revision`", and the record's shape is unchanged. ADR-0184: "A TechSpec states a command surface
  as a transcript", answered below in Surface Transcripts. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable. task_01 edits the Roundfix Skill, whose
  canonical files and `SKILL.md` mirror are Governed Paths. Express maintainer
  authorization: "considere autorizado a ajustar todas as skills se
  necessário", the standing grant "Concedo" for the Governed Paths each Spec
  declares, and "Sim, as três frentes" (2026-10-07) for this Spec. Bounded
  files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`, `skills/roundfix/SKILL.md`.
  No other Governed Path changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`. The
  Spec-contained authorization record is
  `docs/specs/0246-a-sanitize-that-reads-older-folders-and-names-its-refusals/_authorization.md`.

## Current behavior

The current behavior was characterized at `405274c6` with the built binary in
temporary repositories holding one synthetic folder each:

- A projection row for a Task commented out of the graph:
  `build legacy record: Task Graph manifest "<path>/_tasks.md": projection table row names unknown Task "task_02"; remove it or add a matching graph node`,
  exit 2.
- A `refactor` type row:
  `build legacy record: Task Graph manifest "<path>/_tasks.md": projection row for Task "task_01" has invalid type "refactor" (allowed: backend, frontend, data, infra, docs, test, chore, qa); ...`,
  exit 2.
- A `verdict: fail` QA Report without an override:
  `parse legacy record: unknown archive disposition "fail"`, exit 2.
- A clean folder plans with `no-qa`. With all four folders together, the plan
  refuses at the first one.

## System Architecture

No new package or command. Three components change.

- **Graph reading** (`internal/spec/spec.go`, `cause_graph.go`). A legacy
  variant of `parseManifestNodes` (`spec.go:598`) keeps every front-matter
  check. It differs only in the projection rules at `spec.go:640-653` and in
  `parseTaskTypeProjections` (`spec.go:793`, type check at 822).
  `ReadLegacyCauseGraph` exposes it with the tolerances it found.
  `ReadCauseGraph`, `ReadCauseGraphAt` and `spec.Load` are unchanged.
- **Archive Record** (`internal/spec/archive_record.go`, `archive_reader.go`).
  `ArchiveRecordInput` gains `Legacy`. `BuildArchiveRecord` reads the graph
  leniently when it is set and maps a `fail` verdict without an override to
  `failed-qa`. `ParseArchiveRecord` accepts `failed-qa` under guards.
  `ReadArchivedSpec`'s folder branch (`archive_record.go:436`) and
  `PlanLegacyConversion` (`history_sanitize.go:271`) set `Legacy`, and the
  Archive Command (`archive.go:266`) does not. `ArchivedTaskCompleted`
  (`archive_reader.go:66`) excludes the QA Task of a `failed-qa` record.
- **History Sanitize Command** (`internal/cli/history.go`,
  `internal/spec/history_entries.go`). `historyInventory` (`history.go:80`)
  keeps an inventory error as that unit's refusal. `PlanHistoryKinds`
  (`history_entries.go:111`) is split into a per-kind `PlanHistoryKind`, and
  the existing function loops over it. The planning loop (`history.go:247`)
  records a Refused Unit instead of returning, stops once it has `n`
  convertible units, and checks tag coverage over them afterwards. Today that
  check runs at `history.go:216`. The plan and apply printers add the
  `tolerates` and `refused` lines.

### Readers that branch on disposition

Each reader was located with a search for the disposition constants and for
`.Disposition`, `QAOverride` and `ArchivedTaskCompleted`:

1. `spec.ParseArchiveRecord`, behind `ReadArchivedSpec` and
   `ReadArchivedSpecAt`, accepts `failed-qa` under Invariant 6.
2. Reconcile reads the merged head through `taskCompletedAtMergedHead` and
   `archivedTaskCompletedAtHead` (`internal/worktree/merged_head.go:492`, 774),
   and both call `ArchivedTaskCompleted`. The QA Task of a `failed-qa` record
   is not completed (Invariant 7), and every other Task is, as for `no-qa`.
3. Review: `internal/cli/review.go:429` and 1471 read `Source` and
   `SourceRevision` only, and are unchanged. `reviewSpecQAOverride`
   (`internal/cli/review_conventions.go:109`) returns false, because a
   `failed-qa` record has no override.
4. Delivery: `deliver_workflow.go:170` sets `QAOverride` false, so a Retry
   treats the record as an archive without operator override. The prerequisite
   check in `deliver.go:725` accepts the record as an archived Spec.
   `archiveRetirementIsExact` (`deliver_workflow.go:1610`) never meets one,
   because the Archive Command does not write `failed-qa`. All three are
   unchanged.
5. `SC-ARCHIVE-LICENSE` (`internal/speccheck/citations.go:470`) reads archived
   Spec names only, so a `failed-qa` record resolves an `absorbed_by` like any
   other. `archivedTaskContextExists` (`citations.go:1779`) reads `Source` only.
   Both are unchanged.
6. The cause report and `runs causes` (`internal/runcause/report.go:305`,
   `internal/cli/runs_causes.go:143`) print the value as recorded and are
   unchanged.
7. The sanitize plan line (`history.go:362`) prints the value as recorded and
   is unchanged.
8. `supersede.go:166` and `spec_check.go:669` test only that a record exists,
   and are unchanged.

## Implementation Design

### Interfaces

```go
// internal/spec
const ArchiveFailedQA ArchiveDisposition = "failed-qa"
type ArchiveRecordInput struct { /* existing */ Legacy bool }
type LegacyConversion struct { /* existing */ Tolerated []string }
// ReadLegacyCauseGraph reads a Legacy Archive Folder's graph leniently.
func ReadLegacyCauseGraph(root, slug string) (CauseGraph, []string, error)
func PlanHistoryKind(repositoryRoot, revision string, kind ArchiveKind) (HistoryKindPlan, error)

// internal/cli
type historyUnit struct { /* existing */ refusal string }
```

### Invariants

1. The Lenient Legacy Reading applies only when `ArchiveRecordInput.Legacy` is
   true. `PlanLegacyConversion` and `ReadArchivedSpec`'s folder branch set it,
   and the Archive Command never does.
2. The legacy graph keeps every front-matter check of `parseManifestNodes`:
   schema, `requires`, the `qa:` declaration, node ids and files, duplicate ids
   and unknown `needs`. In the projection table, a row whose id is not a graph
   node is skipped and tolerated as
   `projection row <id> names a Task outside the graph`. A row whose type is
   outside `AllowedTaskTypes` is tolerated as
   `projection row <id> has retired Task type "<type>"`. A row that is both is
   tolerated once, as outside the graph. Malformed and duplicate rows refuse
   with today's messages. Tolerances keep table order.
3. `spec.Load`, `ReadCauseGraph`, `ReadCauseGraphAt` and `BuildArchiveRecord`
   without `Legacy` keep today's errors byte for byte.
4. `PlanLegacyConversion` fills `Tolerated` from `ReadLegacyCauseGraph` and
   never writes it into the record.
5. In `BuildArchiveRecord`, a newest QA Report with verdict `fail` and no
   override gives `failed-qa`, keeping `qa_report` and `qa_verdict: fail`. An
   override still gives `qa-override`, and supersession and `no-qa` are
   unchanged. Any other verdict is copied as today and refused by the parser.
6. `ParseArchiveRecord` accepts `failed-qa` only with `qa_verdict: fail`, a
   non-empty `qa_report`, and no `qa_override` key. Otherwise it returns
   `archive record disposition failed-qa requires qa_verdict fail`,
   `archive record disposition failed-qa requires qa_report` or
   `archive record disposition failed-qa cannot carry qa_override`.
7. `ArchivedTaskCompleted(record, record.QATask)` is false for a `failed-qa`
   record. Every other answer is unchanged.
8. The Archive Command still refuses an active Spec with a failing QA and no
   override, with today's message.
9. Units keep today's order. The command plans each one in that order. An
   inventory, conversion or kind-planning error makes that unit a Refused Unit
   whose reason is the error text, and planning continues. With `--batch <n>`,
   planning stops after the `n`th convertible unit. A unit beyond it is not
   examined. Refused Units do not count toward `n`.
10. These still refuse the whole command with exit 2 and write nothing:
    - a usage or configuration error, an external Spec Root, or no Git
      repository;
    - a dirty tree;
    - a missing, lightweight or non-ancestor `history-full` tag, checked before
      planning;
    - a Git error while reading delivery history;
    - a promotion outside the selected convertible units
      (`promotion "<p>" is outside the batch`);
    - a promotion inside a Refused Unit
      (`promotion "<p>" is in refused unit <unit>: <reason>`);
    - a duplicate promotion destination or a promotion error;
    - a tag that misses a path of a selected convertible unit, checked after
      planning.
11. Plan output: when Refused Units exist, the first line appends
    `; <r> unit(s) refused`. Otherwise it is unchanged. Each folder line is
    followed by one `tolerates <folder>: <tolerance>` line per tolerance, and
    then by its candidate lines. After all unit lines and before the `cites`
    lines comes one `refused <unit>: <reason>` line per Refused Unit, in
    order. `<unit>` is the folder's repository-relative path or the kind name.
    The plan exits 0.
12. Apply output: the `refused` lines come first on stdout, followed by the
    confirmation line. When Refused Units exist, that line inserts
    `; <r> unit(s) refused` before `; <remaining> unit(s) remain`. Remaining
    counts every inventoried unit not applied, refused ones included. Refused
    Units are byte-identical after apply.
13. When every examined unit is refused, apply prints the `refused` lines on
    stdout, prints a Preflight refusal on stderr whose reason is
    `history sanitize --apply found no convertible unit; <r> unit(s) refused`,
    exits 2 and writes nothing. When nothing is inventoried, both outputs stay
    as today.

### Data Models

The Archive Record schema stays `roundfix/archive-record/v1` and gains the
disposition value `failed-qa`. No other field, journal or store changes.

### API Contracts

1. API Contract: Lenient Legacy Reading. Per Invariants 1 to 4, a Legacy
   Archive Folder converts with tolerated projection rows, and an active Spec
   refuses exactly as before.
2. API Contract: `failed-qa`. Per Invariants 5 to 8, the record keeps the
   verdict and the report, and every reader handles it as listed above.
3. API Contract: Refused Units. Per Invariants 9 to 13, the plan and the batch
   list each Refused Unit and continue, and the documented stdout lines are
   the ones in Surface Transcripts 1 to 3.

### Surface Transcripts

The QA gate runs each command in a temporary repository it builds with
synthetic Legacy Archive Folders, an annotated `history-full` tag where apply
needs one, and a temporary home.

1. Surface Transcript: the plan names a tolerance, a `failed-qa` record and a
   Refused Unit. The fixture holds a folder with a projection row outside the
   graph, a folder with a failing QA Report and no override, and a folder with
   an unparsable PRD.

   ```transcript
   $ roundfix history sanitize
   stdout:
   history sanitize plan: <u> unit(s) pending; <f> file(s) (<b> bytes) leave docs/history; 1 unit(s) refused
   ...
   tolerates docs/history/specs/<slug>: projection row <id> names a Task outside the graph
   ...
   folder docs/history/specs/<slug>: removes <n> file(s) (<b> bytes) and writes docs/history/specs/<slug>.md (<b> bytes, failed-qa, delivery <d>)
   ...
   refused docs/history/specs/<slug>: <reason>
   ...
   apply with: roundfix history sanitize --apply --batch <n> (needs the annotated tag history-full at or before HEAD)
   stderr:
   exit: 0
   ```

2. Surface Transcript: a batch skips a Refused Unit and fills itself. Three
   folders, with the first one refused and the tag in place.

   ```transcript
   $ roundfix history sanitize --apply --batch 2
   stdout:
   refused docs/history/specs/<slug>: <reason>
   history sanitize applied 2 unit(s): wrote 2 Archive Record(s), reduced 0 file(s), removed <d> file(s) (<b> bytes) kept in Git at <rev> and tag history-full; promoted 0 file(s) to docs/references/; 1 unit(s) refused; <r> unit(s) remain
   stderr:
   exit: 0
   ```

3. Surface Transcript: a batch whose every unit is refused writes nothing. One
   refused folder, with the tag in place.

   ```transcript
   $ roundfix history sanitize --apply --batch 1
   stdout:
   refused docs/history/specs/<slug>: <reason>
   stderr:
   Preflight failed
   ...
   exit: 2
   ```

## Coverage Map

- Goal 1 → Invariants 1 to 4; API Contract 1.
- Goal 2 → Invariants 5 to 8; API Contract 2.
- Goal 3 → Invariants 9 and 11; API Contract 3.
- Goal 4 → Invariants 9, 12 and 13; API Contract 3.
- Core Feature 1 → API Contract 1.
- Core Feature 2 → API Contract 2.
- Core Feature 3 → API Contract 3.
- Core Feature 4 → Build Order 1.
- Success Metric 1 → Testing Approach 3.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 3.
- Success Metric 4 → Testing Approach 1 and 2.
- Success Metric 5 → Testing Approach 4.

## Integration Points

- Git, read-only, as today: delivery lookup and the `history-full` tag checks.
- The adopter's four folders, read once during authoring for their structural
  shape only. Tests build synthetic folders and never read the adopter.

Research record. The Secondbrain was read first: `wiki/index.md`, then `qmd`
queries on the sanitize refusal and on archiving with a failed QA. They
returned the triaged adopter report
(`inbox/roundfix/_triaged/2026-10-07-history-sanitize-recusa-specs-arquivadas-antigas.md`),
this repository's history command reference, and ADR-0154 and ADR-0247 in the
mirror. Exa returned Martin Fowler's "Tolerant Reader" (read only what you
need) and pgloader's batch documentation (rejected rows are logged with their
reason, and the load continues). Those two shaped the Lenient Legacy Reading
and the Refused Unit.

## Testing Approach

1. Graph reading: new `internal/spec/history_sanitize_legacy_test.go`. It
   covers a projection row outside the graph and a retired type, each converted
   and named; a malformed row still refused; and an active Spec with the same
   rows still refused by `spec.Load` and `ReadCauseGraph`, and by
   `BuildArchiveRecord` without `Legacy`. It also covers `ReadArchivedSpec`
   reading such a legacy folder.
2. `failed-qa`: new `internal/spec/archive_failed_qa_test.go`. It covers the
   legacy conversion, the round trip, the parser guards,
   `ArchivedTaskCompleted`, and the Archive Command's unchanged refusal.
3. CLI: new `internal/cli/history_refusal_test.go`. It covers the plan listing
   every Refused Unit, apply filling the batch past one, apply refusing when
   every examined unit is refused, a promotion in a Refused Unit, a refused
   kind unit, and the four adopter shapes converted in one apply.
4. The existing `./internal/spec` and `./internal/cli` history and archive
   tests pass. The one exception is `TestHistorySanitizePreflightsWholeBatch`
   in `internal/cli/history_test.go`, whose expectation changes on purpose: a
   batch of 2 over a malformed second folder now writes the first and third
   records and lists the second.

## Build Order

1. The glossary, the history command reference and the Roundfix Skill's
   archive reference describe the Lenient Legacy Reading, `failed-qa` and
   Refused Units. The skill version is raised and recorded.
2. A Legacy Archive Folder's graph is read leniently (API Contract 1).
3. `failed-qa` is built, parsed and read (API Contract 2) (depends on: 2).
4. The History Sanitize Command lists Refused Units and fills the batch (API
   Contract 3) (depends on: 2, 3).
5. The final QA gate (depends on: 1, 4).

## Risks & Considerations

- A refusal reason carries the underlying error text, which may hold an
  absolute path. The plan is operator output, so this is accepted.
- `--batch <n>` examines units until it has `n` it can convert. A long run of
  Refused Units at the head of the order lengthens planning but never the
  batch.
- An adopter whose Legacy Archive Folder was read by an older Roundfix may now
  see a record where it saw an error. Readers only become more permissive for
  legacy folders.

## Vocabulary Contract

None. The `refused`, `tolerates` and confirmation lines are documented in
`docs/user-guide/commands/history.md` by task_01 and asserted by task_04's
tests.

Glossary terms adopted or revised in `CONTEXT.md` by task_01: **Refused Unit**,
**Lenient Legacy Reading**, **Archive Record** and **Sanitize Batch**.

## Glossary

None.

## Decisions

- A Legacy Archive Folder's graph is read leniently, wherever it becomes a
  record, and an active Spec's graph never is; see ADR-0251.
- A failed QA without an override is `failed-qa`, never a pass or an override;
  see ADR-0251.
- Refused Units are listed and do not count toward the batch; see ADR-0251.
