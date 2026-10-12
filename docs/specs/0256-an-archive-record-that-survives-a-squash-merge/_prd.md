---
spec: 0256-an-archive-record-that-survives-a-squash-merge
status: active
created: 2026-10-12
surfaces: [backend, cli, docs]
---

# An Archive Record that survives a squash merge

An Archive Record says the removed Spec folder is "kept in Git at
`source_revision`". That revision is a commit of the delivery branch, because
the archive runs there before the pull request opens. A squash merge leaves
the commit out of the default branch. Once the branch is deleted, the folder
survives only under the forge's pull request head ref, which a plain clone
does not fetch. On 2026-10-12 the records of this repository's Specs 0252,
0253 and 0254 all pointed at such commits
([the adopted squash entry](references/2026-10-12-an-archive-record-points-at-a-branch-commit-a-squash-merge-drops.md)).
On 2026-10-09 an adopter's `roundfix review` also stopped before calling its
provider. The reason was a prose `README.md` in the archive root, which the
review read as a malformed Archive Record
([the adopted review entry](references/2026-10-12-review-reads-any-markdown-in-the-archive-root-as-an-archive-record.md)).

This Spec makes the record name its delivery pull request. It makes every
reader that opens the folder at `source_revision` recover a missing commit
from the pull request head ref, and makes `roundfix archive` warn when the
revision is outside the default branch. It also makes the archive-root
readers treat only real records as records.

## Prerequisites

Specs 0255, 0257, 0258, 0259 and 0260 are authored in the same cycle. Spec
0255 changes `roundfix archive` and its skill reference. Others may raise the
Roundfix skill's version or change `CONTEXT.md`. No Verification here depends
on another active Spec's artifacts. The queue does not enforce an order, so
the operator orders it. Before this Spec enters a Run or the Delivery Queue it
is rebased on the default branch and checked again with
`roundfix spec check <slug> --strict --run-verification`.

The cancellation item of the same maintainer decision,
`docs/backlog/2026-10-12-a-decomposed-spec-cannot-be-cancelled.md`, stays
open. It does not fit this Spec's four implementation Tasks and is reserved
for a separate Spec.

## Project Constraints

- Identifier strategy: applicable — the record keeps the schema string
  `roundfix/archive-record/v1` and its existing optional `pull_request` field,
  a decimal pull request number. An Archive Record is a file directly under
  the archive root named `NNNN-<slug>.md` whose front matter declares that
  schema. Spec, Task and ADR identifiers keep their schemes. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: applicable — Source Revision Recovery runs
  `git fetch` against the configured remote, so a reader can reach the
  network. The fetch is non-interactive, never prompts for credentials, has a
  time bound, writes no ref and runs only when the commit is missing and a pull
  request number is known. Roundfix adds no credential handling and reads no
  secret. Every test, Verification command and QA row uses a local bare
  repository as the remote and opens no network connection. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0261 (this Spec) decides that the
  record names its pull request through a Pull Request Record Commit, that
  readers fetch `refs/pull/<N>/head` when the commit is missing, and that the
  archive warns. It extends ADR-0247, whose record and cut stay as decided:
  ADR-0247: "Readers work from the record or from Git, never from a folder
  that may be gone". ADR-0248 and ADR-0251 bind the History Sanitize Command's
  records, which already write `pull_request`. ADR-0254 binds sanitize inside
  `baseline update`. ADR-0260 (Spec 0255, merged first) makes the owner archive relink Live
  Documents and record them under `relinked:`; the Pull Request Record Commit
  changes only `pull_request`, never `relinked` or a relinked file, and the
  queue's exact-archive proof keeps accepting the relinked files it already
  proves. ADR-0215 binds archived evidence, and no archived record
  is rewritten. ADR-0223, ADR-0229 and ADR-0237 bind a delivery's work across retries. The new commit must survive them, and ADR-0232's merge evidence
  must still read the merged record. ADR-0165 and ADR-0169 bind the pre-PR
  review's archived-Spec handling and its merge-base diff. ADR-0196 validates
  a review finding before it parks. ADR-0236 judges a failed check on a merge
  with the current default branch. ADR-0193 counts an archived prerequisite.
  ADR-0121 and ADR-0230 serve Legacy Archive Folders, which have no record and
  need no recovery. ADR-0256 binds a release to the skills that describe a
  changed surface. ADR-0187, ADR-0189 and ADR-0233 bind the owned skills'
  versions. ADR-0184 binds the Surface Transcripts. ADR-0258 binds authoring
  against the QA rerun classes. ADR-0257 binds the Daemon as the only full gate
  inside a Run. ADR-0259 binds parallel tests. ADR-0252 and ADR-0253 select
  and run the repository contracts. ADR-0244 holds the declared glossary
  terms. ADR-0179 bounds the Governed Paths. ADR-0182 binds each Task's
  Verification to the facts its gate checks. The gate is bound by ADR-0080,
  ADR-0091, ADR-0096, ADR-0097, ADR-0104, ADR-0167 and ADR-0240. ADR-0194 and
  ADR-0195 cite ADR-0097 but decide what a QA row records and when it is
  observed again; this Spec changes neither, so neither applies. ADR-0210
  cites ADR-0097 for Evidence Snapshots, which this Spec does not change, so
  it does not apply. ADR-0250 cites ADR-0189 for Baseline module versions and
  the coverage record; this Spec changes no module and does not re-record the
  coverage record, so it does not apply. ADR-0249 cites ADR-0194 for the QA
  directory's format and reopen's Late Dependency; the QA gate here uses both
  unchanged, so it applies only as the gate's existing contract. ADR-0093,
  ADR-0117, ADR-0156, ADR-0168, ADR-0176 and ADR-0183 check this Spec by
  citation and receipt. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — only task_01 changes Governed Paths: the
  canonical Roundfix and archive-spec skill files and the two `SKILL.md`
  mirrors. The maintainer authorized skill edits on 2026-09-30 ("considere
  autorizado a ajustar todas as skills se necessário"). On 2026-10-11 the
  maintainer granted the Governed Paths each cycle AE Spec declares
  ("Concedo"). No Makefile, workflow, lint, formatter or `go.mod` change is
  proposed. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`; Spec-contained authorization record:
  `docs/specs/0256-an-archive-record-that-survives-a-squash-merge/_authorization.md`;
  bounded files: `.agents/skills/archive-spec/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/archive.md`, `.agents/skills/roundfix/references/deliver.md`, `.agents/skills/roundfix/references/review.md`, `.agents/skills/roundfix/references/runs.md`, `.agents/skills/roundfix/references/spec-delivery.md`, `skills/archive-spec/SKILL.md`, `skills/roundfix/SKILL.md`.

## Goals

1. A Spec delivered by the Delivery Queue and squash-merged leaves a record
   whose `pull_request` names that pull request.
2. A clone that lacks a record's `source_revision` reads the removed folder
   after one fetch of the pull request head ref, including for records
   archived before this Spec. The fetch writes no ref.
3. An operator who archives on a branch learns from the archive's warning that
   a squash merge would drop the revision, and how to record the pull request.
4. A review, the Run causes report, reconcile and the Spec check never treat
   a non-record Markdown file in the archive root as an Archive Record.
5. The skills, user guides and glossary describe the field, the recovery and
   the warning.

## User Stories

1. As the maintainer, I want an archived Spec's folder to stay readable after
   a squash merge deletes its branch, so that QA reports and Task Results
   remain history without changing my merge policy.
2. As an operator archiving by hand, I want the archive to tell me when its
   revision is not on the default branch and how to record the pull request,
   so that I can make the record recoverable.
3. As the Delivery Queue, I want to record the pull request number in the
   record after I open the pull request, and to prove that commit changes
   nothing else, so that I can publish it without a person.
4. As a reader of Archive Records (a review, Run causes, reconcile or the Spec
   check), I want a missing `source_revision` recovered from the pull request
   head ref, so that I report it unavailable only when the forge no longer
   has it.
5. As a maintainer with a prose index in the archive root, I want reviews to
   ignore it, so that a pull request touching it is still reviewed.

## Core Features

1. **The record names its pull request.** `roundfix archive <slug>
   --pull-request <n>` archives an active Spec and stores the number in the
   record. With a record the current branch added and the default branch does
   not yet hold, it stores the number in that record and changes nothing else.
   It refuses a record already on the default branch, a different number
   already stored, and a value that is not a positive integer.
2. **The Pull Request Record Commit.** After the Delivery Queue opens a
   Spec's pull request, and before it judges checks, it records the number in
   the record. It does this through the Archive Command, commits the change as
   its own commit and pushes it. It accepts the commit only when it changes the
   record alone, and in the record only `pull_request` and a shortened outcome.
   A retry finds that commit and adopts it instead of making a second one. A
   failed recording parks the item and never merges without the field.
3. **Source Revision Recovery.** A reader that needs the folder at
   `source_revision` and finds the commit missing takes the pull request
   number from the record. Without one, it takes the number from the trailing
   `(#N)` of the subject of the commit that added the record. It then fetches
   that pull request's head ref from the remote without writing a ref, and
   reads again. This covers the review's archived-Spec context, the Run causes
   report, reconcile's leftovers and QA readers, and the Spec check's archived
   Task Context. When recovery fails, each reader behaves as it does today.
4. **The archive warning.** When `roundfix archive` writes a record whose
   `source_revision` is not an ancestor of the default branch, it warns on
   standard error and still exits 0. With a pull request recorded, the warning
   names the head ref readers will use. Without one, it says how to record it.
5. **Only real records are records.** Every reader of the archive root treats
   as an Archive Record only a file named `NNNN-<slug>.md` whose front matter
   declares `schema: roundfix/archive-record/v1`, and ignores every other
   Markdown file there. A file that declares the schema but is malformed is
   still an error.
6. **Guidance.** This covers the archive, review, Run causes, delivery and
   hand-delivery references of the Roundfix skill; the archive-spec skill; the
   archive, review, Run causes and delivery user guides; and the glossary
   terms **Source Revision Recovery** and **Pull Request Record Commit**, with
   the revised **Archive Record** and **Archive Command** entries.

## User Experience

An operator who archives on a delivery branch sees the usual `archived …`
line and one `roundfix: warning:` line naming the revision and the default
branch. The warning tells the operator to record the pull request with
`--pull-request` once it opens, or, when the number was given, names the head
ref readers will use. Recording a number prints one `recorded pull request …`
line. A reader that recovers a revision prints nothing new. Its output is the
output it would give if the commit had been local.

## Non-Goals / Out of Scope

- `roundfix cancel` and the `cancelled` disposition. They stay in the open
  Backlog Entry and get a separate Spec.
- Rewriting, backfilling or re-rendering any record already archived, or any
  Legacy Archive Folder.
- Pushing any new ref, keeping a local ref for a fetched pull request, or
  changing a repository's merge policy.
- Changing the Baseline delivery order or any Baseline module.
- Forges other than those that expose `refs/pull/<N>/head`. A failed fetch
  keeps the existing unavailable behavior.
- `roundfix history sanitize`, which already fills `pull_request` for the
  folders it converts.

## Success Metrics

1. Success Metric: a Delivery Queue run against a fake forge ends with a
   merged head whose record carries the opened pull request's number. A retry
   after an interruption at any point between opening the pull request and
   recording it makes exactly one Pull Request Record Commit.
2. Success Metric: in a clone without the commit, with a local bare remote
   that holds `refs/pull/7/head`, each of the four reader classes reads the
   folder. This holds for a record with `pull_request: "7"` and for a record
   without it whose adding commit's subject ends with `(#7)`. Afterwards
   `git for-each-ref` lists the same refs as before. With no reachable
   remote, each reader keeps today's unavailable outcome.
3. Success Metric: archiving on a branch whose `source_revision` is outside
   the default branch prints the warning and exits 0. Archiving on the
   default branch prints no warning.
4. Success Metric: a review whose diff changes a prose `README.md` in the
   archive root reaches its provider. Before this Spec it stopped with
   `read review archive record: missing YAML frontmatter opening marker`.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- The measurement of 2026-10-12 in
  `docs/references/2026-10-12-squash-merged-archive-source-revisions.md`. In
  a fresh repository that fetched only `main`, the `source_revision` of the
  0252, 0253 and 0254 records was absent and not an ancestor of `main`. A
  fetch of `refs/pull/481/head`, `refs/pull/484/head` and `refs/pull/487/head`
  recovered each one and wrote no ref. Each delivery commit's subject ended
  with that `(#N)`.
- GitHub Docs, "Checking out pull requests locally": "Commits in a pull
  request are available in a repository even before the pull request is
  merged." The page also says the `refs/pull/` namespace is read-only, which
  rules out pushing an archive ref there.
- The Git documentation of `git fetch`: "Passing `--no-write-fetch-head` from
  the command line tells Git not to write the file."
- The adopter evidence in the two adopted Backlog Entries: on 2026-10-09
  Spec 0043's record pointed at `b5390aab`, a branch-only commit, and a
  review stopped on `docs/history/specs/README.md`.

## Unreachable Acceptance

- criterion: a reader recovers a real squash-dropped `source_revision` from
  `github.com` inside the Run.
  reason: the Run sandbox denies network access, and Verification and QA
  open no network connection.
  satisfied-by: the authoring measurement in
  `docs/references/2026-10-12-squash-merged-archive-source-revisions.md` and
  the hermetic local-remote tests of Success Metric 2.

## Glossary

- adds: **Source Revision Recovery**
- adds: **Pull Request Record Commit**
- changes: **Archive Record**
- changes: **Archive Command**

## Decisions

- The record names its pull request through a Pull Request Record Commit made
  after the pull request opens, and the operator records it with
  `--pull-request`. Readers recover a missing revision from
  `refs/pull/<N>/head`. The archive warns. No ref is pushed and no archived
  record is rewritten; see ADR-0261.
- A record without `pull_request` is recovered through the trailing `(#N)` of
  its adding commit's subject, so records archived before this Spec are
  reached without a backfill; see ADR-0261.
- An Archive Record is a file named `NNNN-<slug>.md` that declares the record
  schema, and every other Markdown file in the archive root is ignored.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the reader inventory, the recovery
contract, the commit proof and the build order.
