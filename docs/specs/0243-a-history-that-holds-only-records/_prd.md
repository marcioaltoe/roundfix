---
spec: 0243-a-history-that-holds-only-records
status: active
created: 2026-10-07
surfaces: [backend, cli, docs]
---

# A history that holds only records

Spec 0242 made every new archive leave a small **Archive Record** and remove
the Spec folder. It deliberately left the existing history alone. On
2026-10-07 `docs/history` still held 223 Spec folders archived before that
change, with 5,790 files and 58.2 MB, plus 501 Review Artifact files (1.6 MB),
11 handoffs, 31 terminal Backlog Entries and 98 terminal Findings. On
2026-10-06 the maintainer decided: "Não quero nada no histórico que não seja
relevante para o secondbrain", and the clean-up "deve atuar também no que já
existe no diretório".

This Spec is the second of the two. It adds the **History Sanitize
Command**, `roundfix history sanitize`. Without `--apply` it only reports.
With `--apply --batch <n>` it converts the next `n` units of the existing
history: each **Legacy Archive Folder** becomes an Archive Record, retired
Findings and Backlog Entries become a **Reduced History Entry**, and retired
Review Artifacts and handoffs are removed. Every removed byte stays in Git at
the **History Full Tag** and at the revision each record names. The operator
runs each **Sanitize Batch** after this Spec merges, one Pull Request per
batch. This Spec's Tasks never apply a batch to this repository.

## Prerequisites

Spec 0242 is merged (`b07a080e`). This Spec builds on its
`BuildArchiveRecord`, `RenderArchiveRecord`, `ParseArchiveRecord`,
`ReadArchivedSpec` and `judge.AdviseArchive`. No other Spec is active, and no
Verification here depends on another active Spec's artifacts.

## Project Constraints

- Identifier strategy: applicable — the command is `roundfix history
  sanitize` with the flags `--apply`, `--batch`, `--advise` and `--promote`.
  The tag is `history-full`. A record keeps the schema
  `roundfix/archive-record/v1`, and its `disposition` gains `no-qa`. A
  Reduced History Entry keeps its file name and front matter. Spec, Task and
  ADR identifiers keep their schemes. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — only `--advise` can send a request,
  through the Spec judge's existing transport, stage key, monthly ceiling and
  judge log, exactly as `roundfix archive <slug> --plan` does. It sends
  nothing without a key and fails open. Every test, Verification command and
  QA row uses a fake transport and a temporary home and reaches no provider.
  Source: `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0248 (this Spec) decides that the
  existing history is sanitized in batches after a `history-full` tag, that a
  Legacy Archive Folder becomes an Archive Record through the archive's own
  builder, that retired Findings and Backlog Entries are reduced and retired
  Review Artifacts and handoffs removed, and that the legacy readers stay for
  adopters. ADR-0247 introduced the record, ADR-0247: "Folders archived
  before this decision stay readable until a separate, batched migration
  replaces them", and this Spec builds that migration on its builder and its
  advice, which "never refuses, never moves a file and never decides an
  archive". Each batch's Pull Request is the later change whose rule is
  ADR-0215: "A removal is its own later change, and the maintainer's
  explicit approval of that removal is recorded in it", and it carries the
  approval this Spec records.
  ADR-0120 keeps the single History Root, ADR-0120: "retired documentation
  is documentation", and the records and reduced entries stay under it.
  ADR-0230's link pass and ADR-0121's relocation ledger keep serving adopters
  with legacy folders and are not removed. ADR-0154 binds an override's
  fields, and a converted record keeps every one. ADR-0232 grants merge
  evidence from the archived Spec, and ADR-0227's reconcile and ADR-0193's
  archived prerequisite read a record the same way after conversion.
  ADR-0165, ADR-0169, ADR-0229 and ADR-0237 govern the Delivery Queue's
  parks, review diff, operator archives and retries around an archive; a
  batch runs outside the queue and changes none of them, so none applies.
  ADR-0179 bounds the Governed Paths, and ADR-0187, ADR-0189 and ADR-0233
  bind the owned skill's version. ADR-0184 binds the Surface Transcripts and
  ADR-0182 each Task's Verification. ADR-0210 hashes declared inputs and
  reads no history. The QA gate is bound by ADR-0080, ADR-0091, ADR-0096,
  ADR-0097, ADR-0104 and ADR-0167, ADR-0196 validates a pre-PR finding
  before it parks, and ADR-0240 decides when a QA partial qualifies.
  ADR-0194 and ADR-0195 cite ADR-0097 but decide what a QA row records and
  when it is observed again; this Spec changes neither, so neither applies.
  ADR-0093, ADR-0117, ADR-0156, ADR-0168, ADR-0176 and ADR-0183 check this
  Spec by citation and receipt. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the Tasks change Governed Paths: the
  Roundfix skill entry, its archive reference and its mirror, one Baseline
  module and its derived files. The maintainer authorized skill edits on
  2026-09-30 ("considere autorizado a ajustar todas as skills se
  necessário") and Baseline sources and guides ("Autorizar os dois"). On
  2026-10-06 the maintainer granted the Governed Paths of the history
  clean-up ("Concedo"). Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0243-a-history-that-holds-only-records/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/archive.md`, `docs/agents/docs-layout.md`, `docs/agents/setup-context.json`, `internal/baseline/assets/formatter-fixtures/standard-typescript-monorepo/golden/docs/agents/docs-layout.md`, `internal/baseline/assets/modules/context-workflow.json`, `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `skills/roundfix/SKILL.md`.

## Goals

1. A dry run lists every pending unit of the existing history and changes no
   file in the repository.
2. `--apply --batch <n>` converts exactly the next `n` units, only after the
   History Full Tag, on a clean tree, so each batch is one reviewable and
   revertible Pull Request.
3. A converted Legacy Archive Folder leaves an Archive Record of at most
   2,048 bytes, built by the archive's own builder, whose removed bytes read
   back from Git at its `source_revision` and at the tag.
4. Retired Findings and Backlog Entries keep what the Spec check reads, and
   retired Review Artifacts and handoffs leave the tree.
5. The command guide, the Roundfix skill, a Baseline clause and the glossary
   describe the sanitize, and the Secondbrain export follows the history's
   form.

## User Stories

1. As the maintainer, I want the existing history reduced to records, so
   that the repository and its Secondbrain mirror keep only what the
   Secondbrain and the repository's readers need.
2. As the operator, I want to see what each batch removes and which files
   Jev thinks are reusable before I apply it, so that I can promote them in
   the same Pull Request.
3. As the operator, I want the command to refuse a batch that a tag does not
   cover, so that no byte leaves the tree before a tag keeps it reachable.
4. As an adopter whose history was never sanitized, I want every reader of
   legacy folders to keep working, so that sanitizing stays my choice.

## Core Features

1. **The plan.** `roundfix history sanitize` lists, in the fixed unit
   order, every Legacy Archive Folder with its files, bytes, the record it
   would write and its candidate files, and then each other history kind
   with its file count and bytes before and after. It names Markdown files
   outside the History Root that cite a path the plan removes. `--batch
   <n>` limits the plan to the next `n` units. `--advise` adds Jev's advice
   on the batch's candidate files.
2. **The batch.** `--apply --batch <n>` writes each record, reduces or
   removes each kind's files, copies each `--promote` file to
   `docs/references/` and removes each converted folder. It refuses a dirty
   working tree, a missing or lightweight `history-full` tag, a tag that is
   not an ancestor of `HEAD`, and a tag that does not hold a path the batch
   removes or rewrites. It writes nothing when any unit of the batch fails to
   build.
3. **The record of a legacy folder.** `BuildArchiveRecord` builds it with
   `source` set to the folder's path and `source_revision` set to `HEAD`. It
   fills `pull_request`, `delivery_commit` and a missing `archived` date
   from the commit that delivered the Spec, when Git shows one. A folder
   with no QA Report, override or supersession gets the disposition `no-qa`.
4. **The other kinds.** A Reduced History Entry keeps the front matter
   byte-identical, the title, the first paragraph and one line naming the
   revision and path of the full text. Review Artifacts and handoffs are
   removed. Retired ADRs are not touched.
5. **Guidance.** A command guide page, the Roundfix skill's archive
   reference, one Baseline sentence, the glossary terms of this Spec, the
   operator's batch procedure, and a documentation contract that ties the
   Secondbrain export to the presence of Legacy Archive Folders.

## Non-Goals / Out of Scope

- Applying any batch to this repository, creating or pushing the
  `history-full` tag, and editing `.secondbrain-export`. Those are the
  operator's work after this Spec merges, following the procedure in
  `docs/user-guide/commands/history.md`.
- Retiring `ReadArchivedSpec`'s folder branch, the exact-move proof,
  ADR-0230's link pass, ADR-0121's ledger or the pinned-history test reads.
  Adopters with unsanitized history still need them (ADR-0248).
- Changing how future Review Artifacts, handoffs, Findings or Backlog
  Entries retire. A later run of the command reduces whatever accumulates.
- Rewriting Git history or shrinking the pack ("Não agora", 2026-10-06).
- Writing into the Secondbrain. The operator sends cross-project lessons to
  its inbox under `docs/agents/secondbrain.md`.
- A Spec Root outside this repository. The command refuses it.

## Success Metrics

1. Success Metric: on a fixture repository with three Legacy Archive
   Folders (pass, `qa-override`, and no QA Report), a dry run exits 0 and
   leaves every tracked and untracked file byte-identical.
2. Success Metric: `--apply --batch 2` on that fixture converts exactly the
   first two folders. Each record is at most 2,048 bytes, round-trips through
   `ParseArchiveRecord`, and every removed file reads back byte-identical
   from `git show <source_revision>:<source>/<path>` and from the tag.
3. Success Metric: `--apply` exits 2 and changes no file for a dirty tree, a
   missing tag, a lightweight tag, a tag that is not an ancestor of `HEAD`
   and a tag that lacks a batch path.
4. Success Metric: on a disposable clone of this repository, converting
   every unit leaves `docs/history` under 1 MB from 60.6 MB. The authoring
   ablation measured 263,305 bytes of records, 94,873 bytes of reduced
   entries and 43,109 bytes of ADRs, and both `make verify` and
   `make verify-docs` exited 0. With this Spec's export contract, they still
   exit 0 once the clone's `.secondbrain-export` drops its history
   exclusions.
5. Success Metric: the documentation contract fails for an export that keeps
   history exclusions without Legacy Archive Folders, and for one that drops
   them while a folder remains.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- The Jev judgment of 2026-10-06
  (`~/.roundfix-operator/history-value-2026-10-06.md`, model
  `typesafe/jev-1.13-20260917`, 123 sampled artifacts, 9 of 10 hand labels
  matching). It rated QA reports, core artifacts and Task files as records
  and `qa/evidence/` as transient, and it named the reusable files:
  the 0217 and 0218 measurements, already in `docs/references/`; the 0035
  skill analysis and the 0079 pilot report, already triaged in the
  Secondbrain inbox on 2026-10-06; and the 0129 "queue of eight Specs"
  Finding and the 0231 stale-merge Finding, not yet upstream. That is why the
  plan lists candidates and the procedure names those two promotions.
- Git's documentation for `git tag`
  (<https://git-scm.com/docs/git-tag>): `-a` makes a tag object with a
  tagger, date and message, "meant for release", while a lightweight tag is
  "for private or temporary object labels". Git's `git merge-base
  --is-ancestor` documentation
  (<https://git-scm.com/docs/git-merge-base>) exits 0 only when the first
  commit is an ancestor of the second. That is the refusal for the tag.
- The Git book on maintenance and data recovery
  (<https://git-scm.com/book/en/v2/Git-Internals-Maintenance-and-Data-Recovery>):
  removing a file from the tree keeps every reachable version, so the tag
  keeps every byte the batches remove.
- An ablation during authoring on 2026-10-07. A `git clone --no-local` of
  `b07a080e` in a scratch directory converted all 223 folders with
  `BuildArchiveRecord`, reduced the 129 Findings and Backlog Entries and
  removed reviews and handoffs, in one commit (`_techspec.md` → Measured
  inventory). `make -k verify` and `make -k verify-docs` both exited 0 on
  that clone, so no reader needs a fix before the batches. Spec 0242 had
  already moved every reader to the record or to Git.

## Glossary

- adds: **History Sanitize Command**
- adds: **Legacy Archive Folder**
- adds: **Sanitize Batch**
- adds: **Reduced History Entry**
- adds: **History Full Tag**
- changes: **Archive Record**
- changes: **History Root**

## Decisions

- The existing history is sanitized in batches through Pull Requests after a
  `history-full` tag, by an explicit command that is a dry run by default;
  see ADR-0248.
- A Legacy Archive Folder becomes an Archive Record through the archive's
  builder, with the new disposition `no-qa`; see ADR-0248.
- Retired Findings and Backlog Entries are reduced, retired Review Artifacts
  and handoffs removed, and retired ADRs kept; see ADR-0248.
- Jev advises on candidate files through ADR-0247's advice and never gates;
  see ADR-0247.

## Open Questions

None. The operator's batch procedure is in `_techspec.md` → Operator batch
procedure.

## Technical candidate

The [_techspec.md](_techspec.md) records the measured inventory, the unit
order, the command contract, the implementation map and the build order.
