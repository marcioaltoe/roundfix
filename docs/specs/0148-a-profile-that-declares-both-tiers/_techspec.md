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
- Tooling authority: applicable — the Profile, the guide template and the generated repository guides are Baseline-owned and governed. Express maintainer authorization: granted 2026-09-19, recorded in [_authorization.md](_authorization.md); bounded files: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `internal/baseline/assets/templates/guides/agent-instructions.md`, `internal/baseline/assets/templates/index.json`, `docs/agents/spec-routing.md`, `docs/agents/agent-instructions.md`, `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. The last two ride on the standing authorization of 2026-09-18. Sanctioned regeneration: `make baseline-digests` and `make skills-sync`. The generated guides change only through regeneration. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## System Architecture

| Component | Existing package or proposed file | Responsibility |
| --- | --- | --- |
| Profile decisions | `internal/baseline/assets/profiles` | Carry the incremental verification decision beside the complete gate. |
| Profile reader | `internal/baseline` | Expose both decisions, each with its own source, without deriving either. |
| Guide template | `internal/baseline/assets/templates/guides` | Publish both commands where it publishes the gate today. |
| Template index | `internal/baseline/assets/templates/index.json` | Declare the new token for that template, so it may be rendered at all. |
| Generated guides | `docs/agents` | Agree with their sources, by regeneration only. |

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

### Rendering

Where the guide template renders the complete gate, it renders the incremental
command too. A Profile that declares none renders the sentence the clauses
already use for an unmet two-tier contract, so silence keeps meaning what it
means today.

A template may render only the tokens the template index declares for it. The
new token is declared there beside the gate's, and that template's version rises
as the index requires — without it, the rendering cannot exist, which is where
this Spec's first Run correctly stopped.

### Regeneration, not editing

The repository's own guides live inside setup-context markers and are generated
from Baseline sources. They are bounded here so the delivery can run the
sanctioned regeneration and commit its output; no line inside a marker is typed
by hand.

### Interfaces

No exported signature changes. The decision set gains one member, and the
reader's existing accessor pattern serves it.

### Data Models

No stored record changes. The Profile asset gains one entry; the derived digests
that cover Baseline assets move with it, through the sanctioned command.

## API Contracts

1. A Profile that declares the incremental verification command publishes it in
   the generated guides beside the complete gate.
2. A Profile that declares none publishes none, and every other field of its
   output is unchanged.
3. Changing either tier's command leaves the other tier's value and source
   unchanged, and no decision is derived from the other.

## Coverage Map

- Goal 1 → The decision.
- Goal 2 → The decision; API Contract 3.
- Goal 3 → Rendering.
- Goal 4 → Rendering (absent case); API Contract 2.
- User Story 1 → Rendering.
- User Story 2 → The decision.
- User Story 3 → The decision (source preserved).
- User Story 4 → Rendering (absent case).
- Core Feature 1 → The decision.
- Core Feature 2 → The decision; API Contract 3.
- Core Feature 3 → Rendering; Regeneration, not editing.
- Core Feature 4 → Rendering (absent case).
- Success Metric 1 → Testing Approach 2.
- Success Metric 2 → Testing Approach 3.
- Success Metric 3 → Testing Approach 1.
- API Contracts 1-3 → The decision, Rendering.

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
2. **Rendering, at the guide generation seam.** A Profile declaring both
   publishes both; the case fails on the tree as it stands today, where the
   Profile cannot declare the incremental tier at all.
3. **Absence, at the same seam.** A Profile declaring none publishes none and
   changes no other field.
4. **Regeneration is clean.** The repository's generated guides match what their
   sources produce, asserted by the existing regeneration contract.
5. **Outside evidence.** The two shipped clauses are read as they stand and shown
   to name a value the Profile could not produce before this change. Where a
   clause cannot be read, the row records that and does not block.
6. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. The Profile decision and the reader that exposes it independently, with unit
   tests (depends on: none).
2. Rendering in the guide template, including the absent case, with generation
   tests (depends on: 1).
3. Regenerate the repository's guides and any derived pin through the sanctioned
   commands (depends on: 1, 2).
4. Terminal QA (depends on: 1, 2, 3).

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
- **Regeneration only.** A hand edit inside a setup-context marker is undone by
  the next Baseline update.

## Vocabulary Contract

No token is coined. Baseline Profile, verification decision and repository
Verification are existing terms used with their existing meanings.
