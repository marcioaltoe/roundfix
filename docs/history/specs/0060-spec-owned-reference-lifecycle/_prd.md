---
spec: 0060-spec-owned-reference-lifecycle
status: archived
created: 2026-07-28
surfaces: [docs]
archived: "2026-07-31"
source_slug: 0060-spec-owned-reference-lifecycle
---


# Spec-owned reference lifecycle

Promotion from inbox to finding to Spec currently creates links without
transferring the source document, so a Spec can be complete and archived
while the evidence that shaped it stays behind in a separate lifecycle with
links that no longer resolve — archived Specs are linked from findings at
their pre-archive paths, and triaged inbox notes outlive the Specs they
routed. The 2026-07-28 triage session did all of this reconciliation by
hand: removing stale inbox notes, flipping finding statuses, and repointing
archived-Spec links, exactly the manual work a lifecycle contract should
own. Evidence:
[adopted source documents must travel with their owning Spec](../../findings/2026-07-25-spec-owned-reference-lifecycle.md).

## Project Constraints

- Identifier strategy: not applicable — documents keep their basenames and
  Git history; no project-owned Internal Identifier is created. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — documentation lifecycle only;
  no authentication or HTTP surface. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0083 makes adopted sources move
  to their one owning Spec with Git history, new promotions only; no other
  active ADR governs the Spec-artifact documentation lifecycle, whose
  remaining operative contracts are the docs-layout and spec-routing
  guides. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — on 2026-07-28, the maintainer expressly
  authorizes changes to exactly `.agents/skills/write-idea/SKILL.md`,
  `.agents/skills/write-prd/SKILL.md`,
  `.agents/skills/write-techspec/SKILL.md`,
  `.agents/skills/write-tasks/SKILL.md`,
  `.agents/skills/archive-spec/SKILL.md`, their embedded counterparts
  `skills/write-idea/SKILL.md`, `skills/write-prd/SKILL.md`,
  `skills/write-techspec/SKILL.md`, `skills/write-tasks/SKILL.md`, and
  `skills/archive-spec/SKILL.md`, plus the deterministic Skill-digest
  fallout in exactly `internal/baseline/assets/setups/go-cli.json`,
  `internal/baseline/assets/setups/rust-cli.json`,
  `internal/baseline/assets/setups/typescript-bun.json`,
  `internal/baseline/testdata/catalog.digest`,
  `internal/baseline/testdata/catalog.normalized.json`,
  `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`,
  and `internal/baseline/testdata/parity-corpus/v1/manifest.json` (a
  non-roundfix Skill edit cascades into all three setup snapshots;
  extension granted 2026-07-28). The repository-authored sections of
  `docs/agents/docs-layout.md` and `docs/agents/spec-routing.md` are
  documentation, not protected tooling. No other protected tooling
  mutation is authorized. Source: `docs/agents/agent-instructions.md`.

## Goals

- A Spec that adopts a source document becomes its sole owner: the document
  moves with Git history into the Spec and travels with it through archive.
- Traceability survives completion — an archived Spec's PRD, evidence, and
  references resolve as one portable unit.
- The authoring and archive gates enforce the lifecycle mechanically, so
  the manual triage this session performed by hand stops recurring.

## User Stories

1. As a Spec author, I want adopted inbox notes and findings inventoried
   and moved into the Spec's references during PRD creation, so that the
   Spec carries its own evidence.
2. As a maintainer archiving a Spec, I want the archive gate to fail when a
   declared source still sits at its original path or a link points back to
   the inbox or findings tree, so that archives stay self-contained.
3. As a reader of an archived Spec, I want a reference index recording each
   source's original path, type, owning Spec, and adoption date, so that
   provenance survives the move.
4. As an author of a Spec that shares a source with another Spec, I want
   one primary owner and explicit secondary links, so that no document has
   two authoritative homes.
5. As a future investigator, I want moved findings discoverable from the
   findings tree's history and status trail, so that "didn't we already
   see this?" still gets answered.

## Core Features

1. Source adoption is an explicit lifecycle transition owned by PRD
   creation: adopted sources from the inbox and findings move — one move,
   never a copy — into the Spec's references with basename and historical
   content preserved; a finding records `status: done` and its Spec link
   in the same change before it moves.
2. A Spec-local reference index records original path, source type, owning
   Spec, adoption date, and current relative path for every adopted
   source.
3. Repository links to the moved document are rewritten, and the authoring
   gate fails while an adopted source remains at its original path or a
   link is broken.
4. The archive gate validates that every declared source lives inside the
   Spec and no reference resolves to the inbox or findings tree; the Spec
   directory then archives as one unit with stable relative links.
5. Shared-source rules: exactly one primary owning Spec (the first Spec
   that commits to implementation); secondary Specs link the owner's copy
   and never duplicate it.
6. Migration boundary: the workflow applies to new promotions only;
   existing archived Specs and historical findings are not migrated, and a
   finding that already routed to a shipped Spec simply keeps its `done`
   status and links.
7. Triaged raw inbox notes may move directly into a Spec's references when
   they are source material rather than field evidence; anything that
   records observed behavior becomes a finding first.

## User Experience

- Authoring a PRD with adopted sources produces the move, the index, and
  the rewritten links in one reviewable change.
- The archive gate's failure names the offending source or link and the
  transition that fixes it.
- Findings that moved leave a legible trail: `done` status, Spec link, and
  Git history at the old path.

## Non-Goals / Out of Scope

- Rewriting the observations inside moved findings — evidence stays
  append-only and byte-preserved through the move.
- Migrating existing archived Specs or historical findings.
- A generated cross-repository knowledge index; discovery stays Git-native.
- Changing the one-purpose-per-folder docs layout.

## Success Metrics

- A new Spec adopting one finding and one inbox note carries both in its
  references with the index populated, zero stale links, and the finding's
  status flipped in the same change.
- Archiving that Spec succeeds and every reference resolves inside the
  archived directory; an injected stale source path makes the archive gate
  fail with the source named.
- After adoption, no repository link points at the document's original
  path.

## Decisions

- `write-prd` owns the move — the first stage where the repository commits
  to implementing the accepted scope; `archive-spec` owns the validation.
- One move, never a copy or a stub: a document has exactly one
  authoritative home. See
  [ADR-0083](../../adr/0083-adopted-sources-move-to-their-owning-spec.md).
- New promotions only; history is left as history.
- The `CONTEXT.md` glossary entries this Spec's behavior touches — at least
  `Spec`, whose artifact set gains the references directory and index —
  update in this Spec's documentation task, never ahead of implementation.

## Open Questions

None.
