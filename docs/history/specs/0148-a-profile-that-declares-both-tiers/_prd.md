---
spec: 0148-a-profile-that-declares-both-tiers
status: archived
created: 2026-09-19
surfaces: [backend, docs]
archived: "2026-09-19"
source_slug: 0148-a-profile-that-declares-both-tiers
---


# A profile that declares both tiers

Two mandatory clauses in the shipped guidance tell an Agent the same thing: use
the active Baseline Profile's **declared incremental verification command** for
fast local checks, and keep the complete command for the assembled tree. One of
them goes further and says that a missing incremental command leaves the
Profile's two-tier contract unmet, and never authorizes skipping the local tier.

The Profile has no such field. Its verification decisions are format, lint,
test, build and workspace; the complete gate is published to the generated
guides as one value, and the incremental tier is published nowhere, because
nothing declares it.

So the clause asks an Agent to read a decision that cannot exist. In this
repository the gap is filled by knowing that `make verify-incremental` is the
fast tier — knowledge that lives in a Makefile an Agent is not pointed at, and
that no other repository adopting the Baseline would share.

This Spec is the first slice carved from Spec 0121, which keeps the HTTP
decision, the Greenfield refusal, skill regeneration ownership and the
skill-lock reconciliation.

## Project Constraints

- Identifier strategy: applicable — a decision's identifier is how the guide template and the planner name it, so the new decision takes an identifier in the existing scheme and renames none. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the change is a declared decision and its rendering; no credential, transport or HTTP policy is touched. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — the Baseline's derived artifacts and their regeneration are governed by accepted decisions this Spec must preserve. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0155 applies: the `qa` Task declares the gate's matrix, so this Spec's terminal gate covers what its Requirements name.
  ADR-0156 applies: this Spec declares its Success Metrics and API Contracts as numbered units and names each in a Task.
  ADR-0104 applies: acceptance rests on evidence this Spec did not author, which here is the shipped clause that asks for the missing value and the Profile that lacks it.
- Tooling authority: applicable — the Profile is Baseline-owned and governed. Express maintainer authorization: granted 2026-09-19, recorded in [_authorization.md](_authorization.md); bounded files: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`. The record was narrowed on 2026-09-19 when the maintainer split the Spec, and the guide template, the template index, this repository's generated guides and the formatter golden left this Spec's authority with the publication that needed them. Sanctioned regeneration: `make baseline-digests`. Its derived output is ordinary source and changes only through that command. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- The incremental verification tier is a declared decision, so the clause that
  asks for it can be obeyed.
- The two tiers stay distinguishable: neither is inferred from the other, and a
  catalog default stays distinguishable from a repository's choice.
- A repository that declares no incremental command keeps the contract unmet,
  and the absence is reported as absence rather than silently filled.
- The derived artifacts the declaration moves are regenerated through the
  sanctioned command.

## User Stories

1. As an Agent finishing a Task, I want the fast local command named by the
   Profile, so that I can answer whether the slice still holds without guessing
   a Makefile target.
2. As a maintainer adopting the Baseline, I want both tiers declared in one
   place, so that the two commands do not drift apart in guidance.
3. As a maintainer reading a Profile, I want a repository's own choice
   distinguishable from a catalog default, as every other decision already is.
4. As a maintainer of a repository that has no fast tier, I want that absence
   visible, so that nobody reads silence as permission to skip the local check.

## Core Features

1. **The Profile declares the incremental tier.** A verification decision for the
   fast local command exists beside the complete gate, with the same shape every
   other decision has.
2. **The two tiers are separate decisions.** Neither is derived from the other.
   Changing one leaves the other exactly as it was, and a repository's explicit
   choice stays distinguishable from a catalog default.
3. **An undeclared tier stays undeclared.** A Profile without the incremental
   decision exposes no value for it: absence is reported as absence, never as an
   empty string and never as the other tier's value.
4. **The derived artifacts follow.** The catalog digests and plan
   characterizations that the declaration moves are regenerated through the
   sanctioned command, never by hand.

## User Experience

A maintainer reading the Profile sees two verification decisions where there was
one, each with its own source, and a reader asking which command is the fast
tier gets an answer instead of a guess.

Publishing that answer into the generated guides is deliberately not part of
this Spec; until it lands, an Agent still reads only the complete gate there.

## Non-Goals / Out of Scope

- Changing what either command is in this repository, or changing any gate's
  composition.
- Rewriting the clauses that ask for the incremental command; this Spec makes the
  value they ask for exist.
- Publishing the decision into the generated guides, the guide template, the
  template index, this repository's decision record and the formatter golden.
  Delivery proved that a larger change than this slice can carry, and it returns
  to Spec 0121 with its blast radius measured first.
- The HTTP decision and its unused default, the Greenfield refusal on retained
  managed source, skill regeneration ownership, and skill-lock reconciliation.
  Spec 0121 keeps those.
- Adding a profile, or changing another profile's decisions.
- Hand-editing any generated guide: they change by regeneration only.

## Declared intentional breaks

- The catalog digests and plan characterizations move, because a Profile with one
  more decision is a different Profile. That is what the sanctioned regeneration
  is for, and the regeneration contract test is what proves the output matches
  its source.

## Regression locks

- The complete gate decision keeps its identifier, its value and its source.
- Every other verification decision keeps its identifier, kind, tool and command.
- A Profile that declares no incremental command produces the same output it
  produces today for every other field.
- The generated guides are untouched: this Spec publishes nothing to them.

## Acceptance evidence

At least one acceptance row rests on evidence this Spec did not author: the two
shipped clauses that require a declared incremental command, and the Profile
that carries no such field. The replay reads both, shows that the value the
clauses name cannot be produced today, and shows it produced after the change.
Where a clause or the Profile cannot be read, the row records that reason and
does not block.

## Success Metrics

1. A Profile that declares the incremental tier exposes it through the reader
   beside the complete gate, where today the reader has no such value to expose.
2. A Profile that declares no incremental tier reports absence, and its other
   output is unchanged.
3. Changing one tier's command leaves the other's value and source unchanged.

## Decisions

- **A separate decision, not a derived one.** Deriving the fast tier from the
  complete one would make a repository with an unusual pairing undeclarable, and
  would hide whether a maintainer chose it.
- **Declare first, publish later.** Delivery proved publication is its own
  slice: it reaches the guide template, the template index, this repository's
  decision record and the formatter golden, and it cost five Runs and three
  widenings before the maintainer split the Spec on 2026-09-19.
- **The clauses stay as they are.** They already say the right thing; what was
  missing is the value they name.
- **Slice, not portfolio.** Spec 0121 carries seven Core Features; this Spec
  takes the two-tier declaration and leaves the rest there.

## Research basis

**Secondbrain.** Consulted before authoring, index first, then a query on
two-tier verification and on declaring commands in a project profile. The
results were dominated by mirrors of this repository, which are references
rather than independent knowledge, and no source changed the design.

The repository's pending Inbox Entries were read before authoring; none is a
source for this Spec.

**Exa MCP.** Consultation was attempted, and no Exa MCP tool was available in
this session. No external source was read, so no external validation is claimed.

**Local measurement.** The two clauses that require the declared incremental
command, the Profile's five verification decisions, the guide template's single
published gate, and the absence of any incremental field were each read in this
repository before this PRD was written.

## Open Questions

None.
