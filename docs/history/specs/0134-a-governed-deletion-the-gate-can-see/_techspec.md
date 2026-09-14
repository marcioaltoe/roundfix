---
spec: 0134-a-governed-deletion-the-gate-can-see
prd: _prd.md
created: 2026-09-13
---

# A governed deletion the gate can see — Technical Spec

## Executive Summary

Make governed-mutation classification symmetric and make an enumerated
regeneration list authoritative. The trade-off is deliberate: both changes make
the gate refuse more, and one of them will refuse a commit that passes today, so
each is declared as an intentional break rather than slipped in as a fix. No new
component is introduced; both defects are one-sided comparisons inside existing
functions.

## Project Constraints

- Identifier strategy: applicable — preserve the existing refusal codes, operation vocabulary and record field names; no new identifier is introduced. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, network or HTTP surface is touched. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0130 keeps the Governed Path set monotonic and is not narrowed here; ADR-0057 keeps the Daemon the exclusive writer of Task status; ADR-0149 has the grant name the regeneration command while the ownership tree names its outputs, and that division is what the second defect violates. Source: `docs/agents/spec-routing.md`, `docs/agents/domain.md`.
- Tooling authority: applicable — no protected tooling mutation proposed or authorized. The repair is ordinary source in `internal/daemon`, `internal/cli` and `internal/speccheck`, none of which is a Governed Path. Source: `docs/agents/agent-instructions.md`.

## System Architecture

| Component | File | Responsibility |
| --- | --- | --- |
| Governed-mutation classifier | `internal/daemon/task_engine.go` | Compare the paths before and after an Agent turn and report whether a Governed Path changed, in either direction. |
| Settle classification | `internal/cli/settle.go` | Reach the same verdict as the Daemon path for the same change. |
| Sanctioned output resolution | `internal/speccheck/mechanical.go` | Treat an enumerated `outputs` list as the allowed set and resolve repository-owned outputs only for a command-only declaration. |

No component is new. Both defects are one-sided comparisons: one over a set of
paths, one over a set of outputs.

## Implementation Design

### Symmetric classification

The classifier builds a set from the paths before the turn and reports a
governed mutation for any path in the after list that is absent from that set.
A removal is therefore invisible, because a removed path never appears in the
after list at all. A rename compounds it: the source vanishes from the after
list and the destination arrives, so a governed source reads as unremarkable
and an ungoverned destination reads as no mutation.

Compare both directions. A Governed Path present before and absent after is a
mutation exactly as much as one that appears. The comparison stays over path
sets rather than over porcelain rename records, because the porcelain form
reports only the destination for a rename while the snapshot pair already
carries the source — using the snapshots keeps one source of truth and needs no
rename parsing.

The Settle path asks the same question through its own call site. Both must
consume one classifier so the two cannot drift; a second copy is how this class
of defect returns.

### An enumerated list is the whole list

ADR-0149 divides the record's job from the ownership tree's: the grant names the
regeneration command, and the repository-owned declarations name that command's
outputs. A record that also enumerates `outputs` has narrowed its own grant, and
the audit currently unions that narrowing with the full owner-derived set —
which restores exactly what the record excluded.

Resolve repository ownership only when the record supplies no list. When it
supplies one, that list is the allowed set, and an output outside it is out of
grant. This matches the suite guard, which already treats a nil list as the
signal to resolve ownership and a present list as final.

## Coverage Map

- PRD Goal 1 → Governed-mutation classifier, Settle classification.
- PRD Goal 2 → Sanctioned output resolution.
- Core Feature 1 → Governed-mutation classifier.
- Core Feature 2 → Governed-mutation classifier, Settle classification.
- Core Feature 3 → Sanctioned output resolution.
- Core Feature 4 → all three, by leaving every existing refusal path intact.

## Testing Approach

Focused tests at the two existing seams; no new seam is needed.

1. Today's classification answers for an addition, a removal and a rename of a
   Governed Path are recorded before the change.
2. A removal of a Governed Path without the required operation refuses, and a
   rename refuses from its source.
3. The Daemon path and the public Settle path reach the same verdict for the
   same change.
4. A record enumerating `outputs` refuses a consuming commit that changes an
   owner-derived output outside that list; a command-only record still resolves
   ownership and still passes.
5. Every refusal that works today still works with its existing token.

Observation 4 rests on evidence this Spec did not author: the repository's own
ownership declarations under the Baseline tree, written by earlier work, are the
owner-derived set the enumerated list must exclude. If those declarations cannot
be read, the row records blocked with that reason.

## Build Order

1. Characterize classification for an addition, a removal and a rename (depends on: none).
2. Make classification symmetric and share it with the Settle path (depends on: 1).
3. Make an enumerated regeneration list authoritative (depends on: none).
4. Terminal QA (depends on: 2, 3).

## Risks & Considerations

The risk is over-refusal. Symmetric classification will refuse work that
settled before, and an authoritative list will fail a commit the audit passed
before; both are declared breaks in the PRD rather than surprises. The
characterization holds the boundary: every refusal that exists today must keep
its token and its condition, so the change adds cases and moves none.

## Decisions

- Compare path snapshots rather than parse porcelain rename records: the
  snapshot pair already carries both sides, and one source of truth avoids the
  drift that produced this defect.
- One classifier shared by the Daemon and Settle paths, because the review found
  the same hole in both copies.

## Vocabulary Contract

No token is coined. Both refusals reuse existing tokens and their glossary
owners are unchanged.
