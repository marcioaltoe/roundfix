---
spec: 0148-a-profile-that-declares-both-tiers
status: active
created: 2026-09-19
surfaces: [backend, docs]
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
- Tooling authority: applicable — the Profile, the guide template and the generated repository guides are Baseline-owned and governed. Express maintainer authorization: granted 2026-09-19, recorded in [_authorization.md](_authorization.md); bounded files: `internal/baseline/assets/profiles/standard-typescript-monorepo.json`, `internal/baseline/assets/templates/guides/agent-instructions.md`, `internal/baseline/assets/templates/index.json`, `docs/agents/spec-routing.md`, `docs/agents/agent-instructions.md`, `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. The last two ride on the standing authorization of 2026-09-18. Sanctioned regeneration: `make baseline-digests` and `make skills-sync`. The generated guides change only through regeneration. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- The incremental verification tier is a declared decision, so the clause that
  asks for it can be obeyed.
- The two tiers stay distinguishable: neither is inferred from the other, and a
  catalog default stays distinguishable from a repository's choice.
- The generated guides publish both, so an Agent reads them where it already
  reads the gate.
- A repository that declares no incremental command keeps the contract unmet,
  and the unmet state is visible rather than silently filled.

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
3. **The generated guides publish both.** Where the guide template publishes the
   complete gate today, it publishes the incremental command too, and the
   repository's own generated guides are regenerated so they agree.
4. **An undeclared tier stays undeclared.** A Profile without the incremental
   decision publishes no value for it, and the guidance keeps saying the
   two-tier contract is unmet.

## User Experience

An Agent reads the generated instructions and finds both commands named: the
fast one to run after its slice, the complete one that CI runs over the
assembled tree. A maintainer reading the Profile sees two decisions where there
was one, each with its own source.

## Non-Goals / Out of Scope

- Changing what either command is in this repository, or changing any gate's
  composition.
- Rewriting the clauses that ask for the incremental command; this Spec makes the
  value they ask for exist.
- The HTTP decision and its unused default, the Greenfield refusal on retained
  managed source, skill regeneration ownership, and skill-lock reconciliation.
  Spec 0121 keeps those.
- Adding a profile, or changing another profile's decisions.
- Hand-editing any generated guide: they change by regeneration only.

## Declared intentional breaks

- The generated guides gain a line. A repository that pinned their exact bytes
  sees a diff after regeneration, which is what regeneration is for.

## Regression locks

- The complete gate decision keeps its identifier, its value and its source.
- Every other verification decision keeps its identifier, kind, tool and command.
- A Profile that declares no incremental command produces the same output it
  produces today for every other field.
- The generated guides stay byte-identical to what their sources produce.

## Acceptance evidence

At least one acceptance row rests on evidence this Spec did not author: the two
shipped clauses that require a declared incremental command, and the Profile
that carries no such field. The replay reads both, shows that the value the
clauses name cannot be produced today, and shows it produced after the change.
Where a clause or the Profile cannot be read, the row records that reason and
does not block.

## Success Metrics

1. A Profile that declares the incremental tier publishes both commands in the
   generated guides, where today only the complete gate appears.
2. A Profile that declares no incremental tier publishes none, and its other
   output is unchanged.
3. Changing one tier's command leaves the other's value and source unchanged.

## Decisions

- **A separate decision, not a derived one.** Deriving the fast tier from the
  complete one would make a repository with an unusual pairing undeclarable, and
  would hide whether a maintainer chose it.
- **Regenerate the guides, never hand-edit them.** They live inside
  setup-context markers, which the next Baseline update rewrites.
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
