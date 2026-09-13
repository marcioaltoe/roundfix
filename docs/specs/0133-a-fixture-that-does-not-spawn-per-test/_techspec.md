---
spec: 0133-a-fixture-that-does-not-spawn-per-test
prd: _prd.md
created: 2026-09-13
---

# A fixture that does not spawn per test — Technical Spec

## Executive Summary

Move the Task-cycle fixture's committed seed from per-test creation to
per-package creation. The trade-off is deliberate: one shared repository read by
every test, against a private repository per test that costs five process spawns
in a spawn-bound suite. Provenance stays real; only its lifetime changes.

## Project Constraints

- Identifier strategy: applicable — preserve existing test, fixture and helper identities; no new identifier is introduced. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, network or HTTP surface is touched. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0056 has Spec Runs separate Task Capacity and Verification Capacity, and neither is widened here; ADR-0148 keeps every authored Verification command able to fail against the unchanged tree. Source: `docs/agents/spec-routing.md`, `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: "Aprovar os três caminhos", 2026-09-13, recorded in `docs/specs/0133-a-fixture-that-does-not-spawn-per-test/_authorization.md`; bounded files: `internal/baseline/assets/modules/core.json`, `docs/agents/agent-instructions.md`, `docs/agents/setup-context.json`. Sanctioned regeneration: `make baseline-digests` follows the approved module edit. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## System Architecture

| Component | File | Responsibility |
| --- | --- | --- |
| Absent-record rule | `internal/baseline/assets/modules/core.json` and its rendered guides | State that an absent record withholds governed mutation, not the work itself. |
| Operation gate | `internal/authorization/authorization.go`, `internal/cli/implement.go`, `internal/cli/settle.go`, `internal/daemon` | Ask for an operation only where a governed mutation is at stake, and keep every governed-path refusal. |
| Task-cycle fixture | `internal/daemon/task_engine_test.go` | Build a Task-cycle world per test, reading a committed seed created once for the package. |
| Shared Git helper | `internal/gittest/gittest.go` | Create and populate the seed; read only, unchanged by this Spec. |

## Implementation Design

### Two authorities, stated apart

The canonical clause says an absent record "grants nothing", and the gate obeyed
it by refusing the `implement` operation from every Spec. That reads one
sentence as answering two questions: may this change a protected path, and may
this work happen at all. Tooling authority has always answered the first. The
clause now answers only that, and the gate asks for an operation where a
governed mutation is at stake rather than at every dispatch.

Nothing about the governed boundary moves. A change to a Governed Path still
needs an operative record naming that exact path, an absent record still refuses
it, and the changed-path audit is untouched. What changes is that a Spec with no
protected mutation stops needing a maintainer signature to run — which is the
condition Spec 0131 delivered under, before the enforcement existed.

The empty-`paths` rule stays as it is. A record that exists must still name at
least one path, because a record's purpose is to bound something; the repair is
that a Spec with nothing to bound needs no record, not that it writes an empty
one.

### One seed per package

The seed's content is identical for every Task-cycle test: an initialized
repository with the Task source committed. Nothing in the journeys mutates that
seed in a way another test observes, which is why one copy can serve all of
them. Create it once for the package and have each fixture derive its working
state from that copy rather than repeating `init`, `add` and `commit`.

Two properties must survive the change. The authorization reader must still be
handed a real repository root and a real revision, because reading committed
bytes is the behavior the provenance was added for; and a test that needs to
mutate repository state must get its own copy rather than writing into the
shared seed, or tests would observe each other. Where a journey needs private
state, it derives a copy from the seed instead of building one from nothing.

The external and symlinked Spec Root journeys keep their own provenance, because
they exercise a different root shape. There are three of them, so their setup
cost is bounded and stays where it is.

## Coverage Map

- PRD Goal 1 → Absent-record rule, Operation gate.
- PRD Goal 2 → Task-cycle fixture.
- PRD Goal 3 → Absent-record rule.
- User Story 1 → Operation gate.
- User Story 2 → Operation gate.
- User Story 3 → Task-cycle fixture.
- User Story 4 → Task-cycle fixture, through the named reuse test.
- Core Feature 1 → Absent-record rule, Operation gate.
- Core Feature 2 → Operation gate.
- Core Feature 3 → Operation gate.
- Core Feature 4 → Task-cycle fixture.
- Core Feature 5 → Task-cycle fixture, by leaving the reader untouched.
- Core Feature 6 → Task-cycle fixture.
- Core Feature 7 → Task-cycle fixture, through the named reuse test.

## Testing Approach

Focused tests in the Daemon package. Required observations:

1. A Spec that declares no protected tooling mutation and carries no record
   dispatches, while a change to a Governed Path without an operative record
   still refuses. Both observed through the public commands.
2. The canonical clause and the rendered guides state the division, and the
   managed refresh converges.
3. The per-test base fixture body contains no repository initialization.
2. A named test asserts the seed is created once and reused across fixtures.
3. The three external and symlinked Spec Root journeys pass unchanged.
4. Every journey's settlement, staging and commit assertion still passes.
5. The package's wall clock in one concurrent run is measurably lower than the
   immediate pre-repair control, with both numbers recorded.

Observation 5 rests on evidence this Spec did not author: the delivery target's
own package time, measured from a checkout of `main` this Spec did not build,
3.63s on 2026-09-13 against a 6.48s candidate and a 7.70s pre-repair control.
If that checkout cannot be obtained, the row records blocked with that reason.
The observation deliberately asks for a recorded reduction rather than parity
with the target: measurement on 2026-09-13 showed the fixture repair recovers
the smaller part of the gap, and the PRD promotes the residual to the
test-performance campaign.

## Build Order

1. Narrow the absent-record rule in the canonical module and its rendered guides, and ask for an operation only where a governed mutation is at stake (depends on: none).
2. Create the seed once for the package and read it from the fixture (depends on: 1).
3. Terminal QA (depends on: 1, 2).

## Risks & Considerations

The risk is cross-test interference: a shared seed that a journey mutates would
make tests observe each other, which is worse than the cost being removed. The
design answers it by deriving a private copy wherever a journey writes, and the
package's existing assertions are the detector — they pass only if isolation
held.

## Decisions

- One shared seed per package, private copies where a journey writes.
- No production change; the cost is test setup alone.

## Vocabulary Contract

No token is coined.
