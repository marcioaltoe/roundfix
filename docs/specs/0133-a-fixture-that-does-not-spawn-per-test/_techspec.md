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
- Tooling authority: applicable — no protected tooling mutation proposed or authorized. The change is a test fixture in `internal/daemon`, which is not a Governed Path. Source: `docs/agents/agent-instructions.md`.

## System Architecture

| Component | File | Responsibility |
| --- | --- | --- |
| Task-cycle fixture | `internal/daemon/task_engine_test.go` | Build a Task-cycle world per test, reading a committed seed created once for the package. |
| Shared Git helper | `internal/gittest/gittest.go` | Create and populate the seed; read only, unchanged by this Spec. |

## Implementation Design

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

- PRD Goal 1 → Task-cycle fixture.
- PRD Goal 2 → Task-cycle fixture.
- User Story 1 → Task-cycle fixture.
- User Story 2 → Task-cycle fixture, through the named reuse test.
- Core Feature 1 → Task-cycle fixture.
- Core Feature 2 → Task-cycle fixture, by leaving the reader untouched.
- Core Feature 3 → Task-cycle fixture.
- Core Feature 4 → Task-cycle fixture, through the named reuse test.

## Testing Approach

Focused tests in the Daemon package. Required observations:

1. The per-test base fixture body contains no repository initialization.
2. A named test asserts the seed is created once and reused across fixtures.
3. The three external and symlinked Spec Root journeys pass unchanged.
4. Every journey's settlement, staging and commit assertion still passes.
5. The package's wall clock in one concurrent run is within the delivery
   target's range.

Observation 5 rests on evidence this Spec did not author: the delivery target's
own package time, measured from a checkout of `main` this Spec did not build,
3.2s on 2026-09-13. If that checkout cannot be obtained, the row records blocked
with that reason.

## Build Order

1. Create the seed once for the package and read it from the fixture (depends on: none).
2. Terminal QA (depends on: 1).

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
