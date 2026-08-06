---
status: done
created_at: 2026-07-25
updated_at: 2026-08-06
absorbed_by: 0060-spec-owned-reference-lifecycle
---

# Spec workflow — adopted source documents must travel with their owning Spec (2026-07-25)

The repository describes inbox-to-finding-to-Spec promotion as the normal
document lifecycle, but promotion currently creates links without transferring
the source document. A Spec can therefore be complete and archived while the
evidence that shaped it remains in a separate lifecycle with links that no
longer resolve.

Session evidence:

- [`docs/agents/docs-layout.md`](../agents/docs-layout.md) says
  `inbox → findings → spec` is the normal flow and that a file must move when
  its job changes.
- The same guide treats `docs/findings/` as immutable history with addenda and
  does not define how a promoted Finding becomes Spec-owned evidence.
- [`2026-07-23-setup-context-driven-adoption-process-improvements.md`](2026-07-23-setup-context-driven-adoption-process-improvements.md)
  still links Specs 0045 and 0046 under `docs/specs/`, although both Specs now
  live under `docs/specs/_archived/`.
- [`2026-07-24-spec-implementation-sequence.md`](../_inbox/2026-07-24-spec-implementation-sequence.md)
  remained in the inbox after the Specs it routed were implemented and
  archived.
- Archived Specs 0035, 0040, and 0041 already contain a `references/`
  directory, but no authorial workflow defines when or how source documents
  enter it.

## 1. Linked source documents remain outside the Spec lifecycle

- **Symptom / evidence**: PRDs and Findings can link each other while living
  under independent paths. Archiving moves the entire Spec directory but does
  not update external Finding or inbox links, so the traceability chain
  degrades after successful completion.
- **Root cause**: The authorial skills create and archive Spec artifacts, but
  no step assigns ownership of an adopted inbox note or Finding, moves it into
  the Spec, or validates links after the archive.
- **Action / suggestion**: Make source adoption an explicit Spec lifecycle
  transition. Once a Spec accepts a document as an implementation source, the
  Spec becomes its sole owner and the document moves with Git history to
  `docs/specs/<spec>/references/`.

## 2. The move must preserve provenance without duplicating content

- **Symptom / evidence**: Copying a Finding into a Spec would leave two
  authoritative documents that can diverge. Leaving a stub at the original
  path would preserve a second artifact whose lifecycle and archive behavior
  remain ambiguous.
- **Root cause**: The current workflow has link semantics but no ownership,
  provenance, or deduplication contract for Spec references.
- **Action / suggestion**: Use one move, not a copy:

  1. During PRD creation or promotion, inventory every adopted source from
     `docs/_inbox/` and `docs/findings/`.
  2. Assign exactly one owning Spec to each source.
  3. Move the file to `docs/specs/<spec>/references/` without changing its
     basename or historical content. A Finding records `status: done` and its
     Spec link in the same change before it moves.
  4. Record the original path, source type, owning Spec, adoption date, and
     current relative path in a Spec-local reference index.
  5. Rewrite repository links to the new location and fail the authoring gate
     when an adopted source remains at its original path or a link is broken.
  6. Archive the Spec directory as one unit; its PRD, TechSpec, Tasks, QA, and
     source references then retain stable relative links.

## 3. Authorial and archive workflows need one shared contract

- **Symptom / evidence**: `write-idea`, `write-prd`, `write-techspec`,
  `write-tasks`, `qa-gate`, and `archive-spec` consume adjacent parts of the
  traceability chain, but none owns the source-document transfer.
- **Root cause**: Reference collection is treated as document linking rather
  than a state transition with an owner and validation gate.
- **Action / suggestion**: A future Spec must update the repository-owned
  authorial skills, `docs/agents/docs-layout.md`,
  `docs/agents/spec-routing.md`, and the Spec validation contract together.
  `write-prd` is the recommended owner of the move because that is the first
  stage where the repository commits to implementing the accepted scope.
  `archive-spec` must validate that every declared source is already inside
  the Spec and that no source link points back to `docs/_inbox/` or
  `docs/findings/`.

## 4. Shared and historical sources require explicit decisions

- **Symptom / evidence**: One source can influence multiple Specs, and existing
  active or archived Specs already link documents outside their folders.
- **Root cause**: A physical move can have only one destination, while
  cross-Spec evidence can have more than one consumer.
- **Action / suggestion**: The future Spec must decide:

  - how one primary owning Spec is selected when multiple Specs consume the
    same source;
  - how secondary Specs link the owning Spec without creating a duplicate;
  - whether already archived Specs and historical Findings are migrated or
    only new promotions use the workflow;
  - whether raw inbox content must first become a Finding or can move directly
    into a Spec after its reliability and scope are reviewed; and
  - how repository-wide discovery finds Findings after they move into archived
    Spec references.

## What worked — keep

- Keep one purpose per `docs/` folder.
- Keep Findings evidence-first and append-only; moving a file must not rewrite
  the observations it records.
- Keep Spec directories portable so archiving remains one directory move.
- Keep downstream documents linked to upstream sources rather than duplicating
  their content.

## Routing

This finding proposes a new Spec-authoring workflow and remains `pending` until
an implementation Spec defines its source inventory, ownership rules,
reference-index schema, migration boundary, validation errors, skill changes,
and archive behavior.

## Addendum — 2026-07-28 — Routed to Spec 0060

Owned by
[Spec 0060 — Spec-owned reference lifecycle](../specs/0060-spec-owned-reference-lifecycle/_prd.md),
whose PRD resolves this report's open decisions: one primary owning Spec,
new promotions only, and triaged inbox source material may move directly.
The authorial-skill mutations were expressly authorized in that Spec's
Tooling authority entry.
