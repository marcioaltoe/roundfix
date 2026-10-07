---
spec: 0246-a-sanitize-that-reads-older-folders-and-names-its-refusals
status: active
created: 2026-10-07
surfaces: [backend, cli, docs]
---

# A sanitize that reads older folders and names its refusals

On 2026-10-07 an adopter repository running Roundfix 0.55.0 reported that
`roundfix history sanitize` printed no plan at all
([the adopted Backlog Entry](references/2026-10-07-history-sanitize-refuses-older-archived-specs.md)).
`BuildArchiveRecord` checks each Legacy Archive Folder against today's schema,
and the command stops at the first folder that fails. Four of the adopter's
folders fail for reasons the adopter cannot repair:

- Two Task Graph projection tables name Tasks the author had taken out of the
  graph.
- One table uses the Task type `refactor`, which today's closed set no longer
  holds.
- One Spec was archived by maintainer decision with a failing QA Report and no
  override stamp. Its record takes the disposition `fail`, which
  `ParseArchiveRecord` refuses.

Archived metadata is never hand-edited, and `archive --qa-override` acts only
on an active Spec.

This is a bug fix. The Archive Record schema, the Archive Command and the
loading of active Specs keep their contracts. The History Sanitize Command now
reads a Legacy Archive Folder leniently. It records a failed QA without an
override as `failed-qa`, and it names every unit it cannot convert while a
batch converts the others. This minimal PRD exists for the downstream artifact
contract; the design lives in the [_techspec.md](_techspec.md) and ADR-0251.

## Prerequisites

None. No other Spec is active.

## Project Constraints

- Identifier strategy: not applicable. No identifier scheme changes. The Archive
  Record gains the kebab-case disposition value `failed-qa` under the existing
  `disposition` key, beside `no-qa`. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable. No request, credential or network
  call is added. The command reads local files and Git, and every test uses a
  temporary repository and a temporary home. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0251 (this Spec) decides the Lenient
  Legacy Reading, the `failed-qa` disposition and Refused Units. ADR-0248:
  "Each folder becomes `<slug>.md` through `BuildArchiveRecord`, the builder
  the archive itself uses". That holds: the builder is shared, and only a
  Legacy Archive Folder reads its graph leniently. ADR-0248 also says "A folder
  with neither a QA Report, an override nor a supersession gets the new
  disposition `no-qa`", which this Spec extends with `failed-qa`. ADR-0154: "A QA
  archive override records user authority, not a pass", so a failed QA without
  an override is never written as `qa-override`. ADR-0247: "The folder's bytes
  stay in Git at the commit the record names as `source_revision`", and the
  record's shape is unchanged. ADR-0184: "A TechSpec states a command surface as
  a transcript", answered by the TechSpec's Surface Transcripts. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable. task_01 edits the Roundfix Skill, whose
  canonical files and `SKILL.md` mirror are Governed Paths. Express maintainer
  authorization covers this: "considere autorizado a ajustar todas as skills
  se necessário", the standing grant "Concedo" for the Governed Paths each Spec
  declares, and "Sim, as três frentes" (2026-10-07) for this Spec. No other
  Governed Path changes, and no Makefile, formatter or lint configuration
  changes. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`. The Spec-contained authorization record is
  `docs/specs/0246-a-sanitize-that-reads-older-folders-and-names-its-refusals/_authorization.md`.
  Bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/archive.md`, `skills/roundfix/SKILL.md`.

## Goals

- A Legacy Archive Folder converts when its projection table names a Task
  outside the graph or a type outside today's closed set. The plan names what
  was tolerated, and an active Spec still refuses both.
- A Legacy Archive Folder whose newest QA Report says `fail`, with no override,
  becomes a valid Archive Record. It keeps the report name and the `fail`
  verdict, is never a pass or an override, and its QA Task never reads as
  completed.
- The plan lists every unit it cannot convert, each with its reason, and keeps
  planning the others.
- `--apply --batch <n>` converts the next `n` units it can convert. It leaves
  each Refused Unit untouched and listed, and does not count it toward `n`.

## Core Features

1. **Lenient Legacy Reading.** The graph of a Legacy Archive Folder tolerates
   projection rows outside the graph and retired Task types (ADR-0251).
2. **The `failed-qa` disposition.** A record for a failed QA without an
   override, accepted by `ParseArchiveRecord` and read correctly by every
   reader that branches on disposition (ADR-0251).
3. **Refused Unit.** The plan and the batch name each unit they cannot
   convert, with its reason, and continue (ADR-0251).
4. **Docs and glossary.** `CONTEXT.md` gains **Refused Unit** and
   **Lenient Legacy Reading** and revises **Archive Record** and
   **Sanitize Batch**. The history command reference and the Roundfix Skill's
   archive reference describe the behavior.

## Non-Goals / Out of Scope

- Relaxing any check for an active Spec: `spec.Load`, the Implement Command,
  the Spec Consistency Check and the Archive Command keep the current schema.
- Writing `failed-qa` from the Archive Command, or giving an archived Spec an
  override after the fact.
- Tolerating malformed or duplicate projection rows, a manifest with another
  schema, or a QA Report verdict other than `pass`, `partial` or `fail`. These
  stay refusals, now listed per unit.
- Editing archived bytes, existing Archive Records or the adopter repository.

## Success Metrics

1. Success Metric: in a temporary repository with synthetic folders of the four
   adopter shapes, `roundfix history sanitize` exits 0 and plans every folder,
   and `--apply` writes four Archive Records that `ParseArchiveRecord` reads.
2. Success Metric: the failed-QA folder's record reads `disposition: failed-qa`,
   `qa_verdict: fail` and its report name, with no `qa_override` field, and
   `ArchivedTaskCompleted` is false for its QA Task.
3. Success Metric: with one folder that cannot be converted among three,
   `--apply --batch 2` writes the other two records, leaves that folder
   byte-identical, prints `refused <folder>: <reason>` and exits 0.
4. Success Metric: an active Spec whose projection table names a Task outside
   the graph still fails to load with today's message, and the Archive Command
   still refuses an active Spec with a failing QA and no override.
5. Success Metric: every existing sanitize, archive record and spec test
   passes. The one exception is the whole-batch refusal test, which this Spec
   changes on purpose.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- The adopter's report, triaged in the Secondbrain at
  `inbox/roundfix/_triaged/2026-10-07-history-sanitize-recusa-specs-arquivadas-antigas.md`,
  with the exact refusal and the 71-unit plan measured under local patches.
- The adopter's four folders, read only under the maintainer's grant, for
  their structural shape: commented-out graph nodes still in the projection
  table, a `refactor` type row, and three `fail` QA Reports with no
  `qa_override` field.
- Martin Fowler's "Tolerant Reader" (martinfowler.com/bliki/TolerantReader.html,
  2011): read only what you need and ignore what you do not.
- pgloader's batch documentation (pgloader.readthedocs.io, "Batch Processing"):
  rejected rows are logged with their reason, and with `on error resume next`
  the load continues with the rows that are accepted.

## Glossary

- adds: **Refused Unit**
- adds: **Lenient Legacy Reading**
- changes: **Archive Record**
- changes: **Sanitize Batch**

## Decisions

- Read a Legacy Archive Folder's graph leniently, never an active Spec's; see
  ADR-0251.
- Record a failed QA without an override as `failed-qa`, keeping the verdict
  and the report; see ADR-0251.
- List every Refused Unit, and do not count it toward the batch; see ADR-0251.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
