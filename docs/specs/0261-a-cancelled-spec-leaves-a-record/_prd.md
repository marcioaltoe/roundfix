---
spec: 0261-a-cancelled-spec-leaves-a-record
status: active
created: 2026-10-12
surfaces: [backend, cli, docs]
---

# A cancelled Spec leaves a record

A maintainer who drops a Spec that already has a Task Graph has no command to
record it. The Supersede Command needs a delivering Spec, the QA Archive
Override needs every non-QA Task completed, and Task status belongs to the
Daemon. On 2026-10-10 an adopter dropped eight such Specs, one with 5 Tasks
completed, and deleted the folders by hand. No Archive Record remains, and
what that Spec delivered has no history
([the adopted Backlog Entry](references/2026-10-12-a-decomposed-spec-cannot-be-cancelled.md)).

This Spec adds the **Cancel Command**, `roundfix cancel <slug> --reason
<text>`. It writes an Archive Record with the disposition `cancelled`, which
keeps the reason, the date, and the Tasks completed and pending at that
moment. It hands the adopted sources back through **Source Return**, removes
the Spec folder, and refuses while a Run or a Delivery Queue item is active
for the Spec. The readers of Archive Records, the Supersede Command and the
Delivery Queue learn what a cancelled record means. The archive-spec skill,
the Roundfix skill, the docs layout guide and the user guides describe it.

## Prerequisites

Spec 0256 (`0256-an-archive-record-that-survives-a-squash-merge`) is delivered
before this one and is named in the Task Graph's `requires`. It changes the
Archive Record entry of the glossary, the archive-root readers, the
archive-spec skill and the archive reference of the Roundfix skill, which this
Spec also changes. Its Source Revision Recovery and its record-only
`--pull-request` form also apply to a cancelled record. No Verification here
runs a test 0256 creates. Before this Spec enters a Run or the Delivery Queue
it is rebased on the default branch and checked again with
`roundfix spec check <slug> --strict --run-verification`.

Specs 0255, 0257, 0258, 0259 and 0260 are authored in the same cycle and may
also change `CONTEXT.md`, the Roundfix skill or Baseline modules. The operator
orders the queue.

## Project Constraints

- Identifier strategy: applicable — the cancelled record keeps the schema
  string `roundfix/archive-record/v1` and the record path
  `<archive-root>/<slug>.md`. The new disposition is the token `cancelled`.
  Tasks are named by their existing Task identifiers. Spec, Task and ADR
  identifiers keep their schemes. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the command reads the working
  tree, Git and the local Run Database and writes only working-tree files. It
  opens no network connection, reads no credential and calls no provider.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0262 (this Spec) decides the
  cancelled record, Source Return, what readers do with it and the refusals.
  It extends ADR-0247, ADR-0247: "Readers work from the record or from Git,
  never from a folder that may be gone". ADR-0248, ADR-0251 and ADR-0254 bind
  the History Sanitize Command, which converts only Legacy Archive Folders and
  leaves a cancelled record byte-identical. ADR-0193 binds Prerequisite Specs:
  a cancelled record meets none. ADR-0154 keeps the QA Archive Override for a
  Spec whose non-QA Tasks completed, which is why a cancellation is a separate
  command; ADR-0154: "Non-QA Tasks must
still be complete and adopted references self-contained." ADR-0215 binds archived evidence, and no record already archived
  changes. ADR-0223, ADR-0229 and ADR-0237 bind a Delivery Retry, which refuses
  a cancelled Spec. ADR-0057 keeps Task status Daemon-owned: the command reads
  status and never writes it; ADR-0057: "the Daemon is the sole writer of a Task's status". ADR-0165 and ADR-0169 bind the review's
  archived-Spec handling, which a cancelled record reuses unchanged. ADR-0121
  and ADR-0230 serve Legacy Archive Folders, which a cancellation never
  writes. ADR-0058, ADR-0060 and ADR-0250 bind the Baseline module edit,
  its retention and its recorded version. ADR-0257 binds the Daemon as
  the only full gate inside a Run. ADR-0256 binds a release to the skills that
  describe a changed surface. ADR-0187, ADR-0189 and ADR-0233 bind the owned
  skills' versions. ADR-0184 binds the Surface Transcripts. ADR-0258 binds
  authoring against the QA rerun classes. ADR-0259 binds parallel tests.
  ADR-0252 and ADR-0253 select and run the repository contracts. ADR-0244 holds
  the declared glossary terms. ADR-0179 bounds the Governed Paths. ADR-0182
  binds each Task's Verification to the facts its gate checks. The gate is
  bound by ADR-0080, ADR-0091, ADR-0096, ADR-0097, ADR-0104, ADR-0167 and
  ADR-0240. ADR-0261, from Spec 0256, binds Source Revision Recovery for every
  record, a cancelled one included. ADR-0093, ADR-0117, ADR-0156, ADR-0168,
  ADR-0176 and ADR-0183 check this Spec by citation and receipt. ADR-0099 binds
  retention accounting, which records the extended clause `retained`. ADR-0166
  binds the Daemon's recorded paths, which the QA scope audit counts. ADR-0160
  opens a red gate only through `precondition_repairs`, which this Spec does
  not declare, so it does not apply. ADR-0170 decides Task Carry-Forward,
  which this Spec does not change, so it does not apply. ADR-0194 and
  ADR-0195 cite ADR-0097 but decide what a QA row records and when it is
  observed again; this Spec changes neither, so neither applies. ADR-0196
  validates a pre-PR review finding before it parks, and ADR-0210 records
  Evidence Snapshots; this Spec changes neither, so neither applies. ADR-0249
  cites ADR-0194 for the QA directory's format and reopen's Late Dependency;
  the QA gate here uses both unchanged, so it applies only as the gate's
  existing contract. Source:
  `docs/agents/domain.md`. ADR-0260 (Spec 0255) makes an archive relink Live Documents before the
  folder leaves; a cancellation removes the folder too, so it applies the same
  Archive Relink and records `relinked:`.
- Tooling authority: applicable — task_02 alone changes Governed Paths: the
  `spec-workflow` Baseline module and the files its sanctioned regeneration and
  Managed Refresh rewrite, and the canonical Roundfix and archive-spec skill
  files with their `SKILL.md` mirrors. The maintainer authorized skill edits
  on 2026-09-30 ("considere autorizado a ajustar todas as skills se
  necessário") and Baseline sources and guides ("Autorizar os dois"). On
  2026-10-11 the maintainer granted the Governed Paths each cycle AE Spec
  declares ("Concedo"). No Makefile, workflow, lint, formatter or `go.mod`
  change is proposed. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0261-a-cancelled-spec-leaves-a-record/_authorization.md`;
  bounded files: `.agents/skills/archive-spec/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/archive.md`, `.agents/skills/roundfix/references/deliver.md`, `docs/agents/docs-layout.md`, `docs/agents/setup-context.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`, `internal/baseline/assets/modules/spec-workflow.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `skills/archive-spec/SKILL.md`, `skills/roundfix/SKILL.md`.

## Goals

1. A maintainer can drop a decomposed Spec with one command, and the
   repository keeps an Archive Record that says why, when, and which Tasks
   were completed and pending.
2. The Spec's adopted Backlog Entries and Findings return to their active
   directories in their open state, ready for another Spec to adopt.
3. Nothing reads a cancelled record as delivered work: it counts only its
   completed Tasks, meets no prerequisite, supersedes nothing and is never
   retried by the Delivery Queue.
4. A cancellation never races a Run or an active Delivery Queue item.
5. The skills, the docs layout guide, the user guides and the glossary
   describe the command, the disposition and Source Return.

## User Stories

1. As a maintainer, I want to cancel a Spec whose Tasks are partly completed,
   so that the reason and what it delivered stay in history instead of being
   lost with a deleted folder.
2. As a maintainer, I want the cancelled Spec's adopted sources back in the
   backlog and findings, so that the next Spec can adopt them under the usual
   contract.
3. As the Delivery Queue, I want to never treat a cancelled Spec as a met
   prerequisite or retry it, so that I never start work on a premise that was
   dropped.
4. As an operator, I want the command to refuse while a Run or an active queue
   item holds the Spec, and while another active Spec requires it, so that a
   cancellation never strands work in flight.
5. As a reader of history (a review, the Run causes report, the Spec check or
   the History Sanitize Command), I want a cancelled record to read like any
   other record, so that no existing reader breaks on it.

## Core Features

1. **The cancel command.** `roundfix cancel <slug> --reason <text>` cancels
   an active Spec committed at `HEAD`. It writes the Archive Record, performs
   Source Return, removes the Spec folder in the same working-tree change and
   prints one summary line plus one line per returned source. It never
   commits, pushes or creates a Run. It works with or without a Task Graph.
2. **The cancelled record.** The record's disposition is `cancelled`. Its
   archive date is the cancellation date. It carries the reason, the Tasks
   completed and the Tasks not completed in graph order, and the paths of the
   returned sources. Its QA fields record what was observed, and its outcome
   paragraph is the PRD's opening paragraph. Only a cancelled record carries
   these fields.
3. **Source Return.** Each adopted source indexed as a Backlog Entry goes
   back to the backlog with `status: open` and no Spec. Each one indexed as a
   Finding goes back to the findings with `status: pending`, a new
   `updated_at`, no absorption pointer and a dated addendum naming the
   cancellation. Every other byte of a returned source is kept. Other indexed
   sources stay in Git with the folder and are named in the record.
4. **Refusals.** The command refuses, with exit 2 and no write, a missing
   reason, an unknown or uncommitted Spec, a Spec that carries a Supersession,
   an occupied record or return destination, an Active Run for the Spec, an
   unfinished Delivery Queue item for the Spec, an active Spec that names it
   in `requires`, and a non-Markdown file that names its active directory. A
   parked queue item does not block, and the command warns about it.
5. **What a cancelled record means.** Only the Tasks it lists as completed
   count as completed. The Supersede Command refuses a cancelled Spec as the
   delivering Spec. `deliver start` refuses a queue whose Spec requires a
   cancelled Spec. `deliver retry` refuses a cancelled Spec. The review, the
   Run causes report, the Spec check and the History Sanitize Command read a
   cancelled record unchanged.
6. **Guidance.** The Baseline's docs layout guide, the archive-spec skill,
   the archive reference and routing table of the Roundfix skill, the cancel,
   archive, supersede and deliver user guides and the glossary describe the
   command. The glossary adds **Cancel Command** and **Source Return** and
   revises **Archive Record**, **Supersession**, **Prerequisite Spec** and
   **Delivery Retry**.

## User Experience

The operator runs `roundfix cancel <slug> --reason "<why>"` on a clean
checkout. On success, standard output carries one `cancelled …` line with the
record path, the counts of completed and pending Tasks, the removed file count
and the revision that keeps them, then one `returned …` line per returned
source with its new status. Standard error stays empty unless a parked queue
item exists for the Spec or a tracked Markdown line still cites a removed
path; each of those prints one `roundfix: warning:` line and the exit stays 0.
The operator reviews the change and commits it. A refusal prints the usual
Preflight failed block and changes nothing.

## Non-Goals / Out of Scope

- Cancelling an archived Spec, or reopening a cancelled one. A dropped idea
  that comes back is a new Spec that adopts the returned sources.
- Changing Task status. The Daemon owns it, and the record only reads it.
- Closing a pull request, deleting a branch or a Run Worktree, or removing a
  parked Delivery Queue item.
- Rewriting Markdown links that pointed into the cancelled Spec; the command
  reports them as warnings.
- Naming a delivering Spec. That is the Supersede Command's job.
- Changing any record already archived, the History Sanitize Command's
  conversion, or the squash recovery of Spec 0256.

## Success Metrics

1. Success Metric: a Spec shaped like the adopter's dropped Spec, with 5
   Tasks completed and 3 pending and one adopted Backlog Entry, cancels with
   exit 0. Its record lists those 5 and 3 Task identifiers, the reason and
   the date. The Backlog Entry is back in the backlog with `status: open` and
   `spec: null`, and every other byte equals the adopted copy.
2. Success Metric: each refusal of Core Feature 4 exits 2 and leaves the
   working tree byte-identical, measured with `git status --porcelain` and a
   hash of the Spec folder before and after.
3. Success Metric: a cancelled record round-trips through the record parser
   and renderer byte for byte, and the review scope, the Run causes report,
   the Spec check and a History Sanitize plan read a repository holding one
   without an error.
4. Success Metric: a queue that requires a cancelled Spec is refused by
   `deliver start`, a Delivery Retry of a cancelled Spec is refused, and the
   Supersede Command refuses a cancelled Spec as `--by`.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce, kept in
`docs/references/2026-10-12-a-decomposed-spec-dropped-without-a-record.md`:

- the adopter event of 2026-10-10: eight decomposed Specs dropped, one with 5
  Tasks `completed` and 3 `pending`, folders deleted with no record. Success
  Metric 1 reproduces that shape in a synthetic fixture.
- PEP 1: "It is still important to have a record of this fact." PEP 1 keeps
  "Withdrawn" apart from "Superseded".
- The GitHub Changelog of 2022-05-19: "When closing, you can opt to close as
  `complete` (default option) or `not-planned`."

## Glossary

- adds: **Cancel Command**
- adds: **Source Return**
- changes: **Archive Record**
- changes: **Supersession**
- changes: **Prerequisite Spec**
- changes: **Delivery Retry**

## Decisions

- A cancellation is its own command and disposition, never a widened QA
  Archive Override or Supersession; see ADR-0262.
- Source Return sends a Backlog Entry to the backlog as `open` and a Finding to
  the findings as `pending`, never a Finding to the backlog; see ADR-0262.
- A cancelled record counts only its completed Tasks, meets no prerequisite,
  delivers no supersession and is never retried; review, Run causes, the Spec
  check and the History Sanitize Command read it unchanged; see ADR-0262.
- Only an unfinished Delivery Queue item blocks a cancellation. A parked item
  warns, and its retry is refused; see ADR-0262.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the record fields, the Source Return
rules, the refusal order, the reader changes and the build order.
