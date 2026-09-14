---
spec: 0134-a-governed-deletion-the-gate-can-see
status: archived
created: 2026-09-13
surfaces: [backend]
archived: "2026-09-14"
source_slug: 0134-a-governed-deletion-the-gate-can-see
---


# A governed deletion the gate can see

The operation authority this branch built refuses a Task that changes a
Governed Path without the `commit` operation. It does not refuse a Task that
*removes* one. The classifier compares the paths before and after the Agent
turn and reports a governed mutation only for a path that appears in the second
list and not the first, so every deletion is invisible to it — and a rename
reads as a deletion plus an unrelated addition. A Task can drop `Makefile`
without any operation being asked for.

Separately, an authorization record that lists its regeneration `outputs`
explicitly does not get the last word: the audit unions those with every output
the repository's ownership records derive for the same command, so a grant
naming output A lets a consuming commit change output B pass as sanctioned.

Both were found by the pre-Pull-Request review of this branch on 2026-09-13.
The second had already been reported one round earlier and was not acted on.

## Project Constraints

- Identifier strategy: applicable — preserve the existing refusal codes, operation vocabulary and record field names; no new identifier is introduced. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, network or HTTP surface is touched. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0130 keeps the Governed Path set monotonic and is not narrowed here; ADR-0057 keeps the Daemon the exclusive writer of Task status; ADR-0149 has the grant name the regeneration command while the ownership tree names its outputs, and that division is what the second defect violates. Source: `docs/agents/spec-routing.md`, `docs/agents/domain.md`.
- Tooling authority: applicable — no protected tooling mutation proposed or authorized. The repair is ordinary source in `internal/daemon`, `internal/cli` and `internal/speccheck`, none of which is a Governed Path. Source: `docs/agents/agent-instructions.md`.

## Goals

- Removing a Governed Path requires the same operation authority as changing
  one, so the gate cannot be stepped around by deleting instead of editing.
- An authorization record that enumerates its regeneration outputs is
  authoritative for that command; nothing widens it silently.

## Core Features

1. Governed-mutation classification reports a mutation when a Governed Path is
   removed, not only when one appears. A rename of a Governed Path is therefore
   classified from its source as well as its destination.
2. The same classification applies wherever it is asked, so the public Settle
   path and the Daemon's commit path reach the same verdict for the same change.
3. When a record supplies explicit regeneration `outputs`, that list is the
   allowed set. Repository-owned outputs are resolved only for a
   command-only declaration, which is the division ADR-0149 draws.
4. Every refusal that works today keeps working, with its existing token: a
   governed change outside the bounded set, a grant edited in the commit that
   consumes it, and a Spec that declares no authorization.

## Non-Goals / Out of Scope

- External and symlinked Spec Root identity. The same review round produced
  four findings in that family, and three consecutive rounds have now produced
  findings there, so it is promoted to its own Spec with its own threat
  enumeration rather than patched at the end of this branch.
- Any change to the Governed Path set, to what a grant permits, or to the
  operation vocabulary.
- The residual Daemon package wall clock, which the test-performance campaign
  owns.

## Declared intentional breaks

1. A Task that removes a Governed Path without the required operation starts
   refusing, where it settled silently before.
2. A consuming commit that changes a repository-owned output the record does
   not enumerate starts failing the changed-path audit, where an explicit
   `outputs` list previously let it pass.

Everything else keeps behaving as it does: refusal tokens, the operation
vocabulary, the Governed Path set, command-only regeneration resolution, and
Daemon ownership of Task status.

## Regression locks

- Today's classification answers for an addition, a removal and a rename are
  captured as characterization before the change, so the two intended moves are
  visible in the diff and nothing else moves.
- Both new refusals fire on a positive observation — an observed removal, an
  output outside an enumerated list — never on missing or unreadable evidence.
