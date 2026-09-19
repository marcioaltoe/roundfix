---
spec: 0148-a-profile-that-declares-both-tiers
prd: _prd.md
created: 2026-09-19
---

# A profile that declares both tiers — Technical Spec

## Executive Summary

The Profile gains one verification decision, the guide template renders it, and
the repository's generated guides are regenerated so they agree with the sources
that own them. Nothing derives one tier from the other.

The trade-off this design accepts is that a repository which has a fast command
but never declares it keeps publishing nothing for it. Inferring the value would
fill the silence the clause exists to expose.

## Project Constraints

- Identifier strategy: applicable — a decision's identifier is how the guide template and the planner name it, so the new decision takes an identifier in the existing scheme and renames none. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the change is a declared decision and its rendering; no credential, transport or HTTP policy is touched. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the Baseline's derived artifacts and their regeneration are governed by accepted decisions this Spec must preserve. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0155 applies: the `qa` Task declares the gate's matrix, so this Spec's terminal gate covers what its Requirements name.
  ADR-0156 applies: this Spec declares its Success Metrics and API Contracts as numbered units and names each in a Task.
  ADR-0104 applies: acceptance rests on evidence this Spec did not author, which here is the shipped clause that asks for the missing value and the Profile that lacks it.
- Tooling authority: applicable — the Profile, the guide template and the generated repository guides are Baseline-owned and governed. Express maintainer authorization: granted 2026-09-19, recorded in [_authorization.md](_authorization.md); bounded files: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`. The record was narrowed on 2026-09-19 when the maintainer split the Spec; narrowing removes authority and needs no further approval. The last two ride on the standing authorization of 2026-09-18. Sanctioned regeneration: `make baseline-digests` and `make skills-sync`. The generated guides change only through regeneration. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## System Architecture

| Component | Existing package or proposed file | Responsibility |
| --- | --- | --- |
| Profile decisions | `internal/baseline/assets/profiles` | Carry the incremental verification decision beside the complete gate. |
| Profile reader | `internal/baseline` | Expose both decisions, each with its own source, without deriving either. |
| Derived artifacts | `internal/baseline/testdata` | Follow the declaration through the sanctioned regeneration. |

No new package, profile or configuration file is proposed.

## Implementation Design

### The decision

The Profile's verification decisions gain one whose identifier names the
incremental tier, in the shape the others already use: an identifier, a kind, a
tool where one applies, and a command. It sits beside the complete gate decision
rather than inside it, so a reader asking "what runs fast here" and a reader
asking "what runs over the assembled tree" consult different fields.

A Profile that omits the decision produces nothing for it. The reader reports
absence as absence — not an empty string, and not the other tier's value.

### The derived artifacts

A Profile with one more decision is a different Profile, so the catalog digest,
the normalized catalog and the plan characterizations move with it. They are
ordinary source and are rewritten only by the sanctioned regeneration command.

### What this design no longer does

Publishing the decision into the generated guides was part of this design until
delivery measured it: it reaches the guide template, the index that declares
which tokens a template may render, this repository's decision record and the
formatter's golden copy — four governed paths beyond the Profile, found one at a
time by the Agent refusing and the gate refusing. The maintainer split the Spec
on 2026-09-19. Publication returns to Spec 0121, to be authored with its blast
radius measured by running the sanctioned command first.

### Interfaces

No exported signature changes. The decision set gains one member, and the
reader's existing accessor pattern serves it.

### Data Models

No stored record changes. The Profile asset gains one entry; the derived digests
that cover Baseline assets move with it, through the sanctioned command.

## API Contracts

1. A Profile that declares the incremental verification command exposes it
   through the reader beside the complete gate, each with its own source.
2. A Profile that declares none reports absence, and every other field of its
   output is unchanged.
3. Changing either tier's command leaves the other tier's value and source
   unchanged, and no decision is derived from the other.

## Coverage Map

- Goal 1 → The decision.
- Goal 2 → The decision; API Contract 3.
- Goal 3 → The decision (absent case); API Contract 2.
- Goal 4 → The derived artifacts.
- User Story 1 → The decision.
- User Story 2 → The decision.
- User Story 3 → The decision (source preserved).
- User Story 4 → The decision (absent case).
- Core Feature 1 → The decision.
- Core Feature 2 → The decision; API Contract 3.
- Core Feature 3 → The decision (absent case).
- Core Feature 4 → The derived artifacts.
- Success Metric 1 → Testing Approach 2.
- Success Metric 2 → Testing Approach 3.
- Success Metric 3 → Testing Approach 1.
- API Contracts 1-3 → The decision, The derived artifacts.

## Integration Points

- **Baseline digests.** Any derived pin the asset change moves is rewritten by
  the sanctioned command, never by hand.
- **The repository's own guides.** Regenerated so their content matches the
  modules and templates that own them.
- **Spec 0121.** The HTTP decision, the Greenfield refusal, skill regeneration
  ownership and skill-lock reconciliation stay there.

## Testing Approach

1. **Independence, at the profile reader's unit seam.** Changing one tier's
   command leaves the other's value and source untouched, and neither is derived
   from the other.
2. **Declaration, at the profile reader's seam.** A Profile declaring both
   exposes both; the case fails on the tree as it stands today, where the
   Profile cannot declare the incremental tier at all.
3. **Absence, at the same seam.** A Profile declaring none reports absence and
   changes no other field.
4. **Regeneration is clean.** The derived catalog and plan artifacts match what
   the sanctioned command produces, asserted by the existing contracts.
5. **Outside evidence.** The two shipped clauses are read as they stand and shown
   to name a value the Profile could not produce before this change. Where a
   clause cannot be read, the row records that and does not block.
6. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. The Profile decision and the reader that exposes it independently, with unit
   tests (depends on: none).
2. Regenerate the derived catalog and plan artifacts through the sanctioned
   command (depends on: 1).
3. Terminal QA (depends on: 1, 2).

## Risks & Considerations

- **A derived-artifact diff.** Regeneration moves generated guides and possibly
  digests. That is the sanctioned path, and the regeneration contract test is
  what proves the output matches its source.
- **A tempting inference.** Deriving the fast tier from the complete one would
  make the change smaller and the contract weaker; the clause exists precisely to
  expose repositories that have not decided.

## Decisions

- **Two decisions, no derivation.** See the PRD's Decisions; the independence is
  what the third API contract pins.
- **Absence stays absence.** An empty value and a missing decision are different
  statements, and the guidance already has words for the second.
- **Declare first, publish later.** Publication proved to be its own slice, and
  the maintainer split the Spec rather than let one slice keep growing.

## Vocabulary Contract

No token is coined. Baseline Profile, verification decision and repository
Verification are existing terms used with their existing meanings.
