---
spec: 0255-an-archive-that-relinks-and-an-update-that-lists-citations
status: active
created: 2026-10-12
surfaces: [backend, cli, docs]
---

# An archive that relinks and an update that lists citations

Two retirements cut files that live documents still cite, and neither says so.
On 2026-10-09 an adopter's `roundfix baseline update` sanitized its Pending
History; the plan named no citation, and 46 citations in tracked Markdown
outside the History Root stopped resolving
([the adopted update entry](references/2026-10-12-baseline-update-sanitizes-history-without-reporting-the-citations-it-breaks.md)).
On 2026-10-10 the same adopter found open Specs linking into the
`references/` of an owner Spec, as the authoring skills ask, and every one of
those links breaks when the owner archives first, which is the normal queue
order ([the adopted archive entry](references/2026-10-12-an-archive-breaks-links-other-specs-make-into-its-references.md)).

The maintainer decided both outcomes. The owner Spec's archive promotes a
cited file to `docs/references/` and rewrites the inbound citations from each
**Live Document** ("Archive promove e reescreve", 2026-10-11): an **Archive Relink**. The update lists every **History Citation** its sanitization breaks,
in text and JSON and bound by the Plan Digest, and copies a cited file
upstream only when the operator passes `--promote` (2026-10-12). When the
History Full Tag exists, the update also names each **Stale History Citation**, so an adopter finds what an earlier update already broke.

## Prerequisites

None. No other Spec is active in this worktree. Specs 0256 to 0259 of the same
cycle are authored in parallel. None of their artifacts is read by this Spec's
Verification. Spec 0258 also changes the `write-techspec` and `write-tasks`
skills, and Spec 0256 also changes the archive and its readers; each raises a
skill to the next free version when it records it, and whichever lands second
rebases and reruns `roundfix spec check --strict --run-verification`.

## Project Constraints

- Identifier strategy: not applicable. No identifier scheme changes. The
  Archive Record gains the optional lowercase key `relinked`, the update result
  gains lowercase JSON keys under its existing `history` object, and the update
  gains the lowercase flag `--promote`. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable. No request, credential or network
  call is added. The archive and the update read local files and run local
  Git; every test uses temporary repositories and a temporary home. Source:
  `docs/agents/cli.md`, `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable. ADR-0260 (this Spec) decides the Archive
  Relink, the History Citation listing, `--promote` on the update and the Stale
  History Citation report: "The archive refuses, writing nothing, when a
  promotion destination already exists". One owner per source stays, ADR-0083:
  "secondary Specs link the owner's copy"; the archive now carries that copy
  upstream. Promoted bytes stay, ADR-0247: "Promoted copies must equal removed
  blobs". ADR-0230 is the precedent for archive rewrites the Delivery Queue
  accepts: "whose only content changes besides the archive stamp are these
  rewrites". ADR-0173 is the precedent for citations a plan reports and binds:
  "The warnings are part of the plan, so the Plan Digest binds them". ADR-0254
  governs the update's history section, whose rule "The history section changes
  only files under the History Root" ADR-0260 narrows for `--promote`.
  The History Full Tag stays, ADR-0248: "holds every path the batch removes or
  rewrites". ADR-0184: "A TechSpec states a command surface as a transcript",
  answered by the TechSpec's Surface Transcripts. ADR-0100, ADR-0121, ADR-0177
  and ADR-0215 are cited by those records and do not apply: this Spec changes no
  instruction carrier, no legacy archive move, no Relocation Citation scan and
  removes no archived evidence. ADR-0116, ADR-0168, ADR-0176 , ADR-0183, ADR-0251 and ADR-0258 cite listed records and do not
  apply: this Spec changes no Spec check citation rule, no Lenient Legacy
  Reading or failed-QA disposition, and no retirement of a Finding or Backlog
  Entry. Source: `docs/agents/domain.md`.
- Tooling authority: applicable. task_01 edits owned skills and task_02 edits
  the archive source, which are Governed Paths. Express maintainer
  authorization covers this: "considere autorizado a ajustar todas as skills se
  necessário" and the 2026-10-11 cycle grant "Concedo" for the Governed Paths
  each Spec declares. No Baseline module, Baseline guide, Makefile,
  `.roundfixrc.yml`, CI workflow, formatter or lint configuration, or `go.mod`
  changes. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`. The Spec-contained authorization record is
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

## Goals

- Archiving an owner Spec leaves every Live Document citation into its folder
  resolving: a cited file outside the core artifacts is promoted and the
  citation points at the promoted copy, and a citation of the folder or a core
  artifact points at the Archive Record.
- The archive names every citation it does not rewrite, and refuses, writing
  nothing, when a promotion would collide or a Live Document to rewrite holds
  uncommitted changes.
- The Delivery Queue commits and proves an archive that relinked files.
- Before approval, the update's plan lists every History Citation its
  sanitization breaks, marks the Live Document ones with a `--promote` offer,
  and its Plan Digest binds them; `--promote` copies the named file upstream on
  apply and nothing else changes in the citing files.
- With the History Full Tag present, the update names each Stale History
  Citation and how to restore its target.

## User Stories

1. As a maintainer delivering a queue in which secondary Specs link the
   owner's adopted source, I want the owner's archive to promote the source and
   rewrite those links, so that the secondary Specs keep resolving.
2. As an operator archiving a Spec, I want the citations the archive cannot
   rewrite named on stderr, so that I can fix them myself.
3. As an adopter whose update plans Pending History, I want every citation the
   sanitization breaks listed in the plan I approve, so that nothing breaks
   unseen.
4. As that adopter, I want to promote the cited files I keep with `--promote`,
   so that they stay in the tree under `docs/references/`.
5. As an adopter whose earlier update already broke citations, I want the
   update to name them and the tag that holds their targets, so that I can
   restore them.
6. As a Spec author, I want the authoring skills to say that the archive
   carries the owner's copy upstream, so that I keep linking it with one
   owner.

## Core Features

1. **Archive Relink.** Before writing, the archive finds every citation of the
   Spec folder from tracked Markdown outside that folder and outside the History
   Root. A file outside the core artifacts that a Live Document cites is
   promoted as `--promote` promotes it; a citation of the folder or a core
   artifact targets the Archive Record. Each Live Document citation of a
   promoted file, the folder or a core artifact is rewritten; only the link
   destination or the path text changes (ADR-0260).
2. **Archive refusals and warnings.** The archive refuses, writing nothing,
   when a promotion destination exists, when two promotions share a basename,
   or when a Live Document it must rewrite differs from `HEAD`. Every other
   citation of the folder is a warning on stderr naming file, line and target.
   An archive with no inbound citation behaves exactly as today (ADR-0260).
3. **The record and the queue.** The Archive Record lists automatic and manual
   promotions under `promoted` and rewritten files under an optional
   `relinked`, absent when empty. The Delivery Queue stages both and accepts a
   relinked file only when its committed bytes equal the archive's rewrite of
   the parent's bytes (ADR-0260, ADR-0230).
4. **History Citations in the update.** The update's history section lists
   every citation from tracked Markdown outside the History Root of a path a
   selected unit removes, in text and JSON, and the Plan Digest binds them. A
   Live Document citation of a promotable file carries the offer
   `--promote <path>` (ADR-0260, ADR-0173).
5. **Promotion on request.** `--promote <path>`, repeatable, copies a file of a
   selected Legacy Archive Folder to `docs/references/` on apply and lists it
   in that folder's Archive Record. An invalid promotion refuses before any
   write. The update never rewrites a citing file (ADR-0260).
6. **Stale History Citations.** When the History Full Tag exists, the update
   reports each citation that already fails to resolve while the tag holds its
   file, with the restore command. The report is advice outside the Plan Digest
   and never changes the state or the exit code (ADR-0260).
7. **Docs and glossary.** `CONTEXT.md` gains **Live Document**, **Archive Relink**, **History Citation** and **Stale History Citation** and revises  **Archive Record**, **Archive Command** and **Managed Refresh**. The command
   references, the Roundfix Skill, the archive skill and the four authoring
   skills describe the behavior, and the update reference carries the adopter
   note on finding citations an earlier update broke.

## Non-Goals / Out of Scope

- Changing the History Sanitize Command: its plan, `cites` lines and manual
  `--promote` stay as they are.
- Automatic promotion or any rewrite of citing files by the update.
- Rebasing the outward relative links of a promoted file; its bytes stay equal
  to the removed blob.
- Rewriting citations in files that are not Live Documents, such as handoffs,
  inbox notes, guides or the README; they are reported.
- A Spec check finding for links into another Spec's `references/`.
- Restoring Stale History Citations automatically, or reading a remote.
- Editing archived bytes, existing Archive Records or any adopter repository.

## Success Metrics

1. Success Metric: in a temporary repository whose owner Spec has one adopted
   source, where an active secondary Spec cites that source by a relative link
   in its PRD and by a repository path in a Task file, an open Backlog Entry
   links the owner's PRD, and a handoff links the source, the owner's archive
   exits 0. The promoted copy equals the removed blob, the three Live Document
   citations resolve afterwards with every other byte unchanged, the handoff is
   named on stderr and unchanged, and the record lists the promotion and the
   three relinked files.
2. Success Metric: the same archive refuses with exit 2 and changes no file
   when the promotion destination already exists, and again when the secondary
   Spec has an uncommitted change.
3. Success Metric: an archive with no inbound citation prints, writes and
   records exactly what it does today, and every existing archive test passes
   unchanged.
4. Success Metric: the Delivery Queue's archive stage commits the record, the
   promoted copy and the relinked files and proves the archive exact; the same
   commit with one extra byte in a relinked file is refused.
5. Success Metric: in a temporary adopted repository with Pending History whose
   Live Documents cite files of a convertible Legacy Archive Folder, the update
   preview lists each History Citation, marks the Live ones with their offer,
   and its digest changes when one citing line moves. The confirmed update with
   `--promote` copies that file, leaves every citing file byte-identical, and a
   `--promote` outside the selected units exits 2 and writes nothing.
6. Success Metric: a second update over that repository reports each citation
   whose target is gone as a Stale History Citation naming the tag, while a
   repository without Pending History and without stale citations prints and
   digests exactly what it does today.

## Acceptance evidence

The outside-evidence row rests on sources this Spec did not produce:

- The read-only measurement of the adopter repository, committed by this
  planning change as `docs/references/2026-10-12-live-citations-characterization.md`:
  46 citations broken by the 2026-10-09 update, 27 of them from Live
  Documents, with the classes of their targets, and ten secondary-Spec links
  into two owners' `references/`.
- The adopter reports triaged in the Secondbrain on 2026-10-09 and 2026-10-10,
  named in that document.
- The MD Links extension for VS Code, which on a move "rewrites only the URL
  token of broken links. Everything else in the document stays byte-for-byte
  unchanged", and warns on delete "if any Markdown links still point to" the
  file: the published practice behind relinking only the destination and
  reporting what is not rewritten.

## Glossary

- adds: **Live Document**
- adds: **Archive Relink**
- adds: **History Citation**
- adds: **Stale History Citation**
- changes: **Archive Record**
- changes: **Archive Command**
- changes: **Managed Refresh**

## Decisions

- The owner's archive promotes a cited file and rewrites Live Document
  citations; folder and core-artifact citations target the Archive Record; see
  ADR-0260.
- Only Live Documents are rewritten and trigger a promotion; other citations
  are warned; see ADR-0260.
- The archive refuses a colliding promotion or a dirty Live Document instead of
  guessing; see ADR-0260.
- The update lists History Citations bound by its digest and promotes only on
  `--promote`, as the History Sanitize Command does; see ADR-0260.
- The update reports Stale History Citations from the History Full Tag as
  advice outside the digest; see ADR-0260.

## Open Questions

None.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
