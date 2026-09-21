---
spec: 0152-one-declared-acceptance-policy
status: active
created: 2026-09-21
surfaces: [backend, cli, docs]
---

# One declared-acceptance policy

## Executive Summary

The rule that decides whether the newest QA Report is acceptable becomes one
exported function. Archive calls it instead of judging inline, and the derived
Verification command delegates to it instead of re-deciding in awk.

## Project Constraints

- Identifier strategy: applicable — report names and newest-report selection are
  unchanged. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local file reads only. Source:
  `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0080, ADR-0091, ADR-0093, ADR-0096,
  ADR-0097, ADR-0104, ADR-0117, ADR-0130, ADR-0155 and ADR-0156 all hold; this
  Spec unifies an implementation under them. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the derived Verification command and the
  shipped skill are public contract surfaces. Express maintainer authorization:
  recorded in [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `internal/spec/archive.go`. Sanctioned regeneration: `make skills-sync`.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`,
  `docs/agents/specific-repository.md`.

## The two implementations today

`internal/spec/archive.go:143-156` reads the newest report and accepts a `pass`
outright; for a `partial` it requires the blocked rows to be declared
unreachable, and refuses anything else with the verdict it found.

`DerivedQAVerification` in `internal/spec/task.go:95` renders a shell command
whose awk tail reads `exit(closed && verdicts == 1 && verdict == "pass" ? 0 : 1)`.
It never looks at a blocked-row declaration, because it cannot: the fields it
would need are the ones archive parses.

The Daemon executes the rendered command to settle the `qa` Task. That ordering
is what makes the divergence bite — the stricter copy runs first, so the
documented rule never gets a turn.

## The one decision

A single exported function in `internal/spec/qa.go` takes a parsed report and
answers whether it is acceptable, returning the reason when it is not. It states
exactly the rule archive states today, so nothing that is acceptable now stops
being acceptable.

Archive calls it in place of its inline judgement.

## How the derived command reaches it

The rendered command keeps finding the newest report by name — that selection is
not in dispute and stays where a reader can see it. What changes is the
judgement: instead of an awk expression, the command invokes the binary to apply
the one decision and exits on its status.

This is the point of the Spec. A rendered copy of a rule is a copy that drifts;
a rendered call to the rule cannot.

## API Contracts

1. A newest report with verdict `pass` and no blocked rows settles and archives,
   exactly as today.
2. A newest report with verdict `partial` whose blocked rows are declared
   unreachable settles, and archives as today.
3. `fail`, an undeclared `partial`, a missing report, an unparseable report and
   a `pass` carrying blocked rows each refuse, with today's reason.

## Coverage Map

- Goal 1 → The one decision.
- Goal 2 → How the derived command reaches it; API Contract 2.
- Goal 3 → API Contracts 1 and 3.
- Core Feature 1 → The one decision.
- Core Feature 2 → How the derived command reaches it.
- Core Feature 3 → API Contract 2; Testing Approach 2.
- Core Feature 4 → API Contract 3; Testing Approach 3.
- Success Metric 1 → Testing Approach 2.
- Success Metric 2 → Testing Approach 4.
- Success Metric 3 → Testing Approach 3.
- API Contracts 1-3 → The one decision, How the derived command reaches it.

## Integration Points

- **`internal/spec/archive.go`.** Loses its inline judgement, keeps its
  behavior.
- **`DerivedQAVerification`.** Keeps its newest-report selection, delegates its
  verdict judgement.
- **Spec 0122.** Core Feature 5's settlement-semantics documentation sweep stays
  there.

## Testing Approach

1. **The decision itself.** Each accepted and refused shape is asserted directly
   against the one function, including the reason it gives.
2. **A qualifying partial settles.** The derived command accepts a newest report
   whose verdict is `partial` with declared-unreachable blocked rows. Fails on
   the tree as it stands, where the rendered awk requires `pass`.
3. **Nothing else loosens.** `fail`, an undeclared partial, a missing report, an
   unparseable report and a `pass` carrying blocked rows are each still refused
   by the derived command.
4. **Archive is unchanged.** The reports archive accepts and refuses today are
   accepted and refused after the change, for the same reasons.
5. **The skill is true.** The shipped skill and its mirror state the one policy
   and what settles under it, and the mirror matches the canonical file.
6. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. The one eligibility decision, with unit tests for every shape (depends on:
   none).
2. Archive calls it, with its behavior pinned by tests (depends on: 1).
3. The derived command delegates to it, with tests for the qualifying partial
   and every refusal (depends on: 1).
4. The shipped skill (depends on: 2, 3).
5. Terminal QA (depends on: 1, 2, 3, 4).

## Risks & Considerations

- **A gate that stops gating.** Delegating the judgement to a command means a
  command that cannot run would make the check fail, which is the safe
  direction, but a command that exits 0 on an unreadable report would be the
  unsafe one. Testing Approach 3 exercises the missing and unparseable cases for
  exactly this reason.
- **This Spec's own gate.** The change alters the command that settles every
  `qa` Task, including this Spec's. Its terminal gate runs under the new rule,
  so a mistake that loosened settlement would settle this Spec too. Testing
  Approach 3 is the control that does not depend on the gate being honest.
- **A rendered command readers no longer recognise.** The Task file will show a
  call rather than the whole rule. That is the trade this Spec makes: a reader
  loses the inline text and gains the guarantee that what they read is what
  runs.
