---
spec: 0261-a-cancelled-spec-leaves-a-record
prd: _prd.md
created: 2026-10-12
---

# A cancelled Spec leaves a record — Technical Spec

## Executive Summary

The Cancel Command is a sibling of the Archive Command. It reuses the record
builder, renderer and folder cut, and adds one disposition, four record fields
and Source Return. `internal/spec` owns the record, Source Return and the cut
(`Cancel`). `internal/cli` owns the refusals that need the Run Database, the
Delivery Queue and the other active Specs. Three existing readers that treat
an archived Spec as delivered work learn the disposition: the Supersede
Command, `deliver start` prerequisites and `deliver retry`. Every other reader
keeps its code and gains a test that proves it reads a cancelled record. The
primary trade-off is that a binary older than this Spec refuses a repository
holding a cancelled record ("unknown archive disposition"). We accept it
because the schema string stays `roundfix/archive-record/v1` and the
alternative, a second schema, would make every reader branch on two schemas
for one extra disposition (ADR-0262). The docs layout guide is
Baseline-owned, so the guidance Task edits one clause of the `spec-workflow`
module and regenerates its derived files.

## Project Constraints

- Identifier strategy: applicable — the record keeps
  `roundfix/archive-record/v1`, the path `<archive-root>/<slug>.md` and the
  disposition token `cancelled`. `completed_tasks` and `pending_tasks` hold
  Task identifiers from the Task Graph. `returned_sources` holds
  repository-relative paths. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the command reads the working
  tree, Git and the local Run Database, and writes only working-tree files. It
  opens no network connection and reads no credential. Tests use temporary
  repositories and homes. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0262 is implemented here and
  extends ADR-0247, ADR-0248 and ADR-0193. ADR-0057 keeps Task status
  Daemon-owned, so `Cancel` reads it and never writes a Task file;
  ADR-0057: "the Daemon is the sole writer of a Task's status". ADR-0154
  keeps the override for completed non-QA Tasks; ADR-0154: "Non-QA Tasks must
  still be complete". ADR-0223, ADR-0229 and
  ADR-0237 bind the Delivery Retry, which gains one refusal before any queue
  write. ADR-0251 binds `ParseArchiveRecord`'s closed disposition set, which
  gains `cancelled`. ADR-0058, ADR-0060 and ADR-0250 bind the module edit: the
  clause keeps its identity and enforcement, and its versions rise through the
  record command. ADR-0184 binds the Surface Transcripts. ADR-0187, ADR-0189
  and ADR-0233 bind the owned skills' versions. ADR-0256 binds the Skill
  Coverage Map. ADR-0259 binds parallel tests. ADR-0252 and ADR-0253 bind the
  contract run. ADR-0244 binds the glossary terms. ADR-0179 binds the Governed
  Paths. ADR-0261, from Spec 0256, makes Source Revision Recovery reach a
  cancelled record too; this Spec adds nothing to it. The rest of the PRD's
  row applies unchanged. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — task_02 alone changes Governed Paths, and
  the grant bounds exactly those. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0261-a-cancelled-spec-leaves-a-record/_authorization.md`;
  bounded files: `.agents/skills/archive-spec/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/archive.md`, `.agents/skills/roundfix/references/deliver.md`, `docs/agents/docs-layout.md`, `docs/agents/setup-context.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`, `internal/baseline/assets/modules/spec-workflow.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `skills/archive-spec/SKILL.md`, `skills/roundfix/SKILL.md`.

## System Architecture

No new package. `internal/spec` gains `cancel.go` beside `archive.go`;
`internal/cli` gains `cancel.go` beside `archive.go` and one dispatch case.

Absence evidence, measured on `7a8b7f70`:

- `grep -rn 'case "cancel"' internal/cli | wc -l` printed `0`: no command
  exists.
- `grep -rn 'cancelled' internal/spec/archive_record.go internal/spec/archive.go | wc -l`
  printed `0`: no disposition exists.
- `ls internal/spec/cancel*.go internal/cli/cancel*.go` matched no file.

### Reader inventory

Measured on `7a8b7f70` with a search for every caller of
`ParseArchiveRecord`, `ReadArchivedSpec`, `ReadArchivedSpecAt`,
`ArchivedTaskCompleted` and `ArchivedSpecSlugs`.

| reader | place | today with a `cancelled` record | after |
| --- | --- | --- | --- |
| record parser | `internal/spec/archive_record.go` `ParseArchiveRecord` | refused: unknown archive disposition | accepted with its four fields (task_01) |
| Task completion and merge evidence | `ArchivedTaskCompleted`, used by `internal/worktree/merged_head.go` | every Task completed | only `completed_tasks` (task_01) |
| Supersede Command | `internal/cli/supersede.go` `knownDelivererSpec` | a known deliverer | refused as `--by` (task_04) |
| queue prerequisites | `internal/cli/deliver.go` `validateDeliveryPrerequisites` | met | refused by `deliver start` (task_04) |
| Delivery Retry | `internal/cli/deliver.go` `runDeliverRetry` | retried | refused before the engine (task_04) |
| owner prerequisites | `internal/cli/deliver_workflow.go` `UnmetPrerequisites` | met by the record's presence | unchanged: the Cancel Command refuses while an active Spec requires the Spec, and `deliver start` refuses a new queue |
| review scope and Spec context | `internal/cli/review.go` | read as a record | unchanged, proven by a test (task_04) |
| Run causes | `internal/runcause/report.go` | `archived (cancelled)` | unchanged, proven by a test (task_04) |
| Spec check | `internal/speccheck` archived names and backlog rule | read as a record | unchanged, proven by a test (task_04) |
| History Sanitize Command | `internal/spec/history_sanitize.go` | never read: it plans only Legacy Archive Folders | unchanged, proven by a test (task_04) |
| delivery archive stage | `internal/cli/deliver_workflow.go` | not reached: the queue never archives a cancelled Spec | unchanged |

## Implementation Design

### Interfaces

```go
// internal/spec (task_01)
const ArchiveCancelled ArchiveDisposition = "cancelled"
// ArchiveRecord gains: CancelReason string; CompletedTasks, PendingTasks, ReturnedSources []string
type CancelRequest struct {
	SpecsRoot, Slug, Reason, SourceRevision, RepositoryRoot string
	BuiltInRoot                                             bool
	CancelledAt                                             time.Time
}
type ReturnedSource struct{ From, To, Status string } // repository-relative; Status open | pending
type CancelResult struct {
	RecordPath                string
	Completed, Pending        []string
	Returned                  []ReturnedSource
	RemovedFiles              int
	RemovedBytes              int64
	RemovedPaths              []string // repository-relative, for citation advice
}
func Cancel(req CancelRequest) (CancelResult, error)
```

`Cancel` runs in this order and writes nothing until step 6:

1. The Spec folder exists under the Spec Root with `_prd.md`, and no
   `_supersession.md`.
2. The reason is non-blank after whitespace collapses to single spaces, and
   `SourceRevision` is 40 hex characters.
3. When `_tasks.md` exists, `Load` reads the graph and fails on any manifest
   error. Tasks with status `completed` go to `Completed`, every other Task to
   `Pending`, each in graph order.
4. Neither `<archive-root>/<slug>.md` nor `<archive-root>/<slug>/` exists.
5. Source Return is planned from `references/_index.md`: each row whose type
   is `backlog` or `finding` is read. A missing or symbolic-link file, a
   missing `status` key in its front matter, or an occupied destination
   refuses the whole cancellation.
6. `BuildArchiveRecord` builds the record. `Cancel` then sets
   `Disposition: cancelled`, the four fields, the PRD's opening paragraph as
   `Outcome` and `archiveDate(CancelledAt)` as `Archived`. It writes the
   record and each returned source with `O_EXCL`, removes them all again if
   any write fails, and then removes the Spec folder.

### Data Models

A cancelled record renders these keys after `unproven` and before `adrs`, and
only for `cancelled`:

```yaml
cancel_reason: <reason, one line>
completed_tasks: [task_01, task_02]
pending_tasks: [task_03, task_04, task_05]
returned_sources: [docs/backlog/<name>.md, docs/findings/<name>.md]
```

Each list renders `[]` when empty. `ParseArchiveRecord` requires all four for
`cancelled` and refuses a record of any other disposition that carries one of
them with `archive record disposition <d> cannot carry <key>`. The QA fields
keep what `BuildArchiveRecord` observed: an empty `qa_report` and
`qa_verdict` without a report, or the newest report's name and verdict. The
record's `sources` keeps every indexed `path`, as today.

Source Return rewrites only front matter lines and, for a Finding, appends an
addendum:

- `backlog` → `docs/backlog/<path>`: the `status:` line becomes
  `status: open` and the `spec:` line becomes `spec: null`. A trailing
  comment on either line is kept, and a missing `spec:` line is inserted after
  `status:`.
- `finding` → `docs/findings/<path>`: the `status:` line becomes
  `status: pending`, `updated_at:` becomes the cancellation date, and an
  `absorbed_by:` line is removed. The body gains, after one blank line:

  ```markdown
  ## Addendum — <date> — Spec <slug> cancelled

  Spec <slug> was cancelled on <date> before it implemented this finding: <reason>
  ```

- Every other type stays in the folder, so it stays in Git at
  `source_revision`, and is named in the record's `sources`.

### API Contracts

1. API Contract: `roundfix cancel <slug> --reason <text>` cancels an active
   Spec as in Interfaces. It refuses with exit 2 and no write, in this order:
   a missing slug (`missing required Spec slug; pass roundfix cancel <slug>
   --reason <text>`); a missing or blank `--reason` (`missing required
   --reason; pass --reason <text>`); an unexpected argument; no Git working
   tree; an unknown active Spec (`unknown active Spec slug "<slug>"`); a
   non-Markdown file that names the active directory (the Archive Command's
   wording with `cancel`); uncommitted changes under the folder; an Active Run
   for the Spec after orphan reclaim (`Active Run <id> holds Spec "<slug>";
   stop it with: roundfix stop <id>`); a Delivery Queue item for the Spec at a
   stage other than `merged` or `parked`; an active Spec whose `requires`
   names it; and every refusal of `Cancel`. On success it exits 0. It exits 1
   only when a write or removal fails after validation.
2. API Contract: after a successful cancellation, the command writes one
   warning to standard error for a parked Delivery Queue item of the Spec. It
   then writes one warning per tracked Markdown line outside `docs/history/`
   that cites a removed path, found with `spec.HistoryCitations` over
   `RemovedPaths` and the folder path. Neither warning changes the exit code.
3. API Contract: `ParseArchiveRecord` accepts `cancelled` with its fields,
   and `RenderArchiveRecord(ParseArchiveRecord(b))` returns `b` for a record
   `Cancel` wrote. `ArchivedTaskCompleted(record, task)` for a cancelled
   record is true only for a Task in `completed_tasks`.
4. API Contract: `roundfix supersede --by <slug>` refuses a cancelled Spec
   with `superseding Spec "<slug>" is cancelled; a cancelled Spec delivers no
   content` and exit 2. `roundfix deliver start` refuses a queue whose Spec
   requires a cancelled Spec with `Delivery Queue Spec <a> requires cancelled
   Spec <b>` and exit 2. `roundfix deliver retry <slug>` refuses a Spec whose
   Archive Record in the checkout is cancelled with `Spec "<slug>" is
   cancelled in <record>; a cancelled Spec is never delivered` and exit 2,
   before it reads or writes the queue.

### Surface Transcripts

1. Surface Transcript: cancelling a Spec with 5 completed and 3 pending Tasks
   and two adopted sources.

   ```transcript
   $ roundfix cancel <slug> --reason "no consumer"
   stdout:
   cancelled <slug> -> docs/history/specs/<slug>.md; 5 completed Task(s), 3 pending; removed <count> file(s) (<bytes> bytes) kept in Git at <revision>
   returned docs/backlog/<entry>.md as open
   returned docs/findings/<finding>.md as pending
   stderr:
   exit: 0
   ```

2. Surface Transcript: a parked queue item and a live citation are reported.

   ```transcript
   $ roundfix cancel <slug> --reason "no consumer"
   stdout:
   cancelled <slug> -> docs/history/specs/<slug>.md; 0 completed Task(s), 2 pending; removed <count> file(s) (<bytes> bytes) kept in Git at <revision>
   stderr:
   roundfix: warning: Delivery Queue item <slug> stays parked; a Delivery Retry of a cancelled Spec is refused, so close its branch and pull request by hand
   roundfix: warning: <path>:<line> cites <target>, which this cancellation removed
   exit: 0
   ```

3. Surface Transcript: an unfinished Delivery Queue item refuses.

   ```transcript
   $ roundfix cancel <slug> --reason "no consumer"
   stdout:
   stderr:
   Preflight failed

   Reason:
     Spec "<slug>" cannot be cancelled while its Delivery Queue item is at stage <stage>; stop the queue with: roundfix deliver stop

   No side effects:
     Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.

   Usage:
     Run 'roundfix cancel --help' for usage.
   exit: 2
   ```

4. Surface Transcript: an Active Run refuses.

   ```transcript
   $ roundfix cancel <slug> --reason "no consumer"
   stdout:
   stderr:
   Preflight failed

   Reason:
     Active Run <run-id> holds Spec "<slug>"; stop it with: roundfix stop <run-id>

   No side effects:
     Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.

   Usage:
     Run 'roundfix cancel --help' for usage.
   exit: 2
   ```

5. Surface Transcript: an active Spec that requires it refuses.

   ```transcript
   $ roundfix cancel <slug> --reason "no consumer"
   stdout:
   stderr:
   Preflight failed

   Reason:
     Spec "<slug>" cannot be cancelled while active Spec <dependent> requires it; remove it from that Spec's requires first

   No side effects:
     Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.

   Usage:
     Run 'roundfix cancel --help' for usage.
   exit: 2
   ```

6. Surface Transcript: the help.

   ```transcript
   $ roundfix cancel --help
   stdout:
   Usage:
     roundfix cancel <slug> --reason <text>

   Cancels an active Spec that will not be delivered. Writes one Archive Record
   with disposition cancelled, the reason, the date and the Tasks completed and
   pending, returns each adopted Backlog Entry to docs/backlog/ (status: open)
   and each adopted Finding to docs/findings/ (status: pending), and removes the
   committed Spec folder. The removed bytes stay in Git at the record's
   source_revision. cancel creates no Run, never changes a Task status, and
   never commits or pushes.

   Refuses while an Active Run or an unfinished Delivery Queue item holds the
   Spec, while an active Spec requires it, and for a Spec that carries a
   supersession.

   Options:
     --reason  Why the Spec is cancelled, recorded in the Archive Record

   Exit codes:
     0  cancelled
     1  cancel write failed
     2  Preflight Validation failed
   stderr:
   exit: 0
   ```

In every transcript `<revision>` is the 12-hex prefix of `source_revision`,
and the built-in Spec Root is used. The root help gains the line
`  roundfix cancel <slug> --reason <text>` after the `archive` usage line and
`  cancel     Cancel an active Spec with an Archive Record` after the
`archive` command line.

## Coverage Map

- Goal 1 → `Cancel`, the Cancel Command and API Contract 1 (task_01, task_03)
- Goal 2 → Source Return in `Cancel` (task_01)
- Goal 3 → `ArchivedTaskCompleted` (task_01); Supersede Command, `deliver start` and `deliver retry` refusals (task_04)
- Goal 4 → the Cancel Command's refusals (task_03)
- Goal 5 → Baseline clause, skills, user guides and glossary (task_02)
- Story 1 → task_01, task_03
- Story 2 → task_01
- Story 3 → task_04
- Story 4 → task_03
- Story 5 → task_01, task_04
- Core Feature 1 → task_03
- Core Feature 2 → task_01
- Core Feature 3 → task_01
- Core Feature 4 → task_01, task_03
- Core Feature 5 → task_01, task_04
- Core Feature 6 → task_02, task_03
- Success Metric 1 → task_01, task_03
- Success Metric 2 → task_03
- Success Metric 3 → task_01, task_04
- Success Metric 4 → task_04

## Integration Points

- The Run Database: read-only `ActiveSpecRun` and `DeliveryQueue`, opened
  only when the database file exists, as `reopen` does. Orphaned active-run
  locks are reclaimed as `settle` does.
- Git: `rev-parse`, `status --porcelain` and `ls-files` through the existing
  helpers. Nothing is committed.

## Testing Approach

Every new top-level test in a Parallel Test Package calls `t.Parallel()`
first or names the process-wide state it changes (ADR-0259). Every test uses
temporary repositories and homes.

- `internal/spec` (task_01): the 0030-shaped fixture (8 Tasks, 5 completed);
  the record round-trip; refusal of the four keys on another disposition;
  `ArchivedTaskCompleted`; Source Return of a Backlog Entry and a Finding with
  byte comparison outside the changed lines; each `Cancel` refusal leaves the
  tree unchanged; an injected write failure rolls back.
- `internal/cli` (task_03): each Surface Transcript asserted whole in one
  named test; the queue refusal for every unfinished stage; a parked item
  that does not refuse; the root help lines.
- `internal/cli` and `internal/speccheck` (task_04): the three refusals, and
  review, Run causes, Spec check and History Sanitize characterizations over a
  repository holding a cancelled record.
- `internal/baseline` (task_02): the clause text, its rendering in this
  repository's guide and the Standard TypeScript golden, and a retained clause
  for a Source Baseline adopter.

## Build Order

1. The cancelled record, `ArchivedTaskCompleted` and `Cancel` with Source
   Return (task_01).
2. Guidance: the Baseline clause and its regeneration, the skills, the user
   guides including the cancel guide, and the glossary (task_02) (depends
   on: none; it runs beside step 1).
3. The Cancel Command, its refusals, warnings, help and coverage entry
   (task_03) (depends on: 1, 2).
4. The readers: Supersede Command, `deliver start`, `deliver retry`, and the
   characterization of review, Run causes, Spec check and History Sanitize
   (task_04) (depends on: 1, 2).
5. QA gate (task_05) (depends on: 3, 4).

Spec 0256 is delivered before this Spec and is named in `_tasks.md`
`requires`. It changes `ArchivedSpecSlugs`, the review's archive-root loop,
`internal/cli/archive.go`, the archive-spec skill, the Roundfix archive
reference and the **Archive Record** glossary entry. This Spec rebases on it.
It adds its glossary sentence at the end of the entry, its skill text under
its own headings, and no change to those readers. No Verification here runs a
test 0256 creates.

## Skill and document texts

task_02 writes these texts unless another Task is named. Each phrase must
appear verbatim on whitespace-normalized text.

- `CONTEXT.md`, new **Cancel Command**, placed immediately before **History
  Sanitize Command**: "The `roundfix cancel <slug> --reason <text>` command
  that retires an active Spec which will not be delivered: it writes an Archive
  Record with the disposition `cancelled`, performs Source Return and removes
  the Spec folder, which stays in Git at `source_revision`. It refuses while an
  Active Run or an unfinished Delivery Queue item holds the Spec or an active
  Spec requires it, never changes a Task status and never commits
  (ADR-0262)." `_Avoid_`: Delete command, abandon, archive override.
- `CONTEXT.md`, new **Source Return**, placed after **Cancel Command**: "The
  Cancel Command's hand-back of a cancelled Spec's adopted sources: a Backlog
  Entry returns to `docs/backlog/` with `status: open` and no Spec, a Finding
  returns to `docs/findings/` with `status: pending` and a dated addendum, and
  every other byte is kept. Other adopted sources stay in Git with the folder
  (ADR-0262)." `_Avoid_`: Unadoption, restore, re-triage.
- `CONTEXT.md`, revised **Archive Record**: add as the entry's last sentence
  "A cancelled Spec's record has the disposition `cancelled` and carries
  `cancel_reason`, `completed_tasks`, `pending_tasks` and `returned_sources`;
  only its listed completed Tasks count as completed (ADR-0262)."
- `CONTEXT.md`, revised **Supersession**: add "A cancelled Spec never delivers
  another Spec's content (ADR-0262)."
- `CONTEXT.md`, revised **Prerequisite Spec**: add "A cancelled Spec's record
  never meets a prerequisite (ADR-0262)."
- `CONTEXT.md`, revised **Delivery Retry**: add "It refuses a Spec whose
  Archive Record is `cancelled` (ADR-0262)."
- The archive reference of the Roundfix skill gains a `## Cancel Command`
  section carrying "`roundfix cancel <slug> --reason <text>` writes an Archive
  Record with the disposition `cancelled`" and "A cancelled Spec never meets a
  prerequisite, never supersedes another Spec, and a Delivery Retry refuses
  it". The Roundfix skill's routing row names `archive`, `cancel`, `history`,
  `supersede`.
- The archive-spec skill gains a `## Cancelling a Spec` section, never inside
  `### QA settlement`, carrying "A Spec that will not be delivered is cancelled
  with `roundfix cancel <slug> --reason <text>`, never deleted by hand" and
  "returns each adopted Backlog Entry to `docs/backlog/` with `status: open`".
- `docs/user-guide/commands/archive.md` carries "A Spec that will not be
  delivered is cancelled with `roundfix cancel`".
- `docs/user-guide/commands/supersede.md` carries "A cancelled Spec cannot be
  the superseding Spec".
- `docs/user-guide/commands/deliver.md` and the deliver reference of the
  Roundfix skill carry "`deliver start` refuses a
  queue whose Spec requires a cancelled Spec, and `deliver retry` refuses a
  cancelled Spec".
- task_02 writes `docs/user-guide/commands/cancel.md` with the help of
  Surface Transcript 6, the refusals of API Contract 1 and the warnings of API
  Contract 2, plus a `cancel` row in `docs/user-guide/commands.md` and in the
  command table of `docs/user-guide/usage.md`.
- The `guidance` of `clause.spec.keep-artifacts-in-spec-folder` in the
  `spec-workflow` module gains, after its last sentence: "A Spec dropped
  before delivery is cancelled, never deleted by hand: a runtime cancel
  command, such as `roundfix cancel <slug> --reason <text>`, writes an Archive
  Record with the disposition `cancelled` that keeps the reason, the date and
  the Tasks completed and pending at that moment, returns each adopted Backlog
  Entry to `docs/backlog/` with `status: open` and each adopted Finding to
  `docs/findings/` with `status: pending`, and removes the Spec folder in the
  same change. It refuses while a Run or a delivery queue item is active for
  the Spec. A cancelled Spec never meets a prerequisite and never supersedes
  another Spec. A runtime without a cancel command must report cancellation
  as unsupported rather than deleting the folder." Every earlier byte of the
  clause stays.

### Baseline versions and derived files

task_02 raises by one, from the value on its starting commit,
`rule.spec.docs-layout` and `guide.spec-docs-layout` in `spec-workflow`, then
runs the Module Version Record step, which chooses the module's own version
(ADR-0250). A prototype on a scratch clone of `7a8b7f70` applied the clause
and those two raises, then ran the record command, `make baseline-digests` and
`go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills
--yes --format text`. It changed exactly these files, a second refresh
reported `File changes: 0`, and `go test ./internal/baseline -count=1` passed:

- `internal/baseline/assets/modules/spec-workflow.json`
- `internal/baseline/module-versions.json`
- `internal/baseline/assets/profiles/standard-typescript-monorepo.json`
- `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`
- `internal/baseline/testdata/catalog.diagnostics.golden.json`
- `internal/baseline/testdata/catalog.digest`
- `internal/baseline/testdata/catalog.normalized.json`
- the four plan goldens under `internal/baseline/testdata/plan-characterization/`
- `docs/agents/docs-layout.md` and `docs/agents/setup-context.json`

The clause keeps its `id`, its `mandatory` enforcement and no `replaces` list,
so a Source Baseline adopter's refresh records it `retained`.

## Risks & Considerations

- An older binary refuses a cancelled record. The upgrade notice and the
  user guide say so, and the record is written only by a binary that reads
  it.
- A relative Markdown link inside a returned source that pointed into the
  Spec folder now dangles. The citation warning names live links into the
  folder; links inside the returned sources are out of scope.
- A parked queue item keeps its branch, Run Worktree and any open pull
  request. The warning says so, and `deliver retry` refuses the Spec.
- A History Finding retired earlier with `absorbed_by: <slug>` keeps that
  pointer to a cancelled Spec. The Spec check still resolves it, because the
  record exists. Repairing those pointers is a follow-up, recorded in the
  report.
- Spec 0256 touches the same glossary entry, skills and readers, so this Spec
  rebases on it and runs the strict check again.
- The Baseline edit reaches every adopter through the Managed Refresh. A
  runtime without the command is told to report cancellation as unsupported.

## Research record

Secondbrain, consulted on 2026-10-12 through `wiki/index.md` and
`qmd query "cancel a decomposed spec archive record cancelled disposition return adopted sources backlog" --all --files --min-score 0.3`:

- `inbox/roundfix/_triaged/2026-10-10-nao-ha-caminho-para-cancelar-uma-spec-ja-decomposta.md`
  is the adopter observation behind the adopted entry. It asked for a
  disposition of its own, the reason, the completed and pending Tasks and the
  return of adopted sources, which is the shape of the record.
- `projects/vortex/mirror/docs/references/2026-10-10-plano-de-specs-ate-o-release.md`,
  section "Specs fora do plano", shows the adopter's two outcomes: replaced
  ("Substituída") and dropped ("Deprecada"). That is why a cancellation does
  not name a delivering Spec and the Supersede Command stays separate.
- `projects/roundfix/mirror/docs/adr/0247-an-archive-leaves-an-archive-record-and-the-spec-folder-stays-in-git.md`
  confirmed that the record and the cut are reused unchanged.
- The pending `inbox/roundfix/` namespace was empty.

Exa, consulted the same day:

- PEP 1 (<https://peps.python.org/pep-0001/>): "It is still important to have
  a record of this fact." It keeps Withdrawn and Superseded as separate
  states, which supported a separate disposition and no `--by` option.
- GitHub Changelog, 2022-05-19
  (<https://github.blog/changelog/2022-05-19-the-new-github-issues-may-19th-update/>):
  closing as `complete` or `not-planned`, with a recorded reason. It
  supported recording the reason in the record rather than in a commit
  message.

## Glossary

- adds: **Cancel Command**
- adds: **Source Return**
- changes: **Archive Record**
- changes: **Supersession**
- changes: **Prerequisite Spec**
- changes: **Delivery Retry**

## Decisions

- One record schema with a new disposition and four disposition-only fields,
  rather than a second schema; see ADR-0262.
- `Cancel` lives in `internal/spec` beside `Archive` and reuses its builder,
  renderer and cut. The CLI owns only the checks that need the Run Database,
  the queue and other Specs.
- A Finding returns to the findings as `pending` with an addendum, never to
  the backlog; see ADR-0262.
- Only an unfinished queue item refuses; a parked one warns, and its retry is
  refused; see ADR-0262.
- `UnmetPrerequisites` stays unchanged, because the Cancel Command refuses
  while an active Spec requires the Spec and `deliver start` refuses a queue
  that requires a cancelled one.
- The docs layout guidance extends the existing clause rather than adding a
  clause, so no Source Baseline row is needed.
