---
spec: 0133-a-fixture-that-does-not-spawn-per-test
status: active
created: 2026-09-13
surfaces: [backend]
---

# A fixture that does not spawn per test

Giving the Daemon's Task-cycle fixtures real Git provenance fixed a correctness
defect and introduced a cost one. The base fixture that every Task-cycle test
builds now initializes a repository, stages it, commits it, reads a revision and
resolves an authorization — roughly five process spawns per test, across about
forty tests. This suite is spawn-bound, so that setup is the entire slowdown.

Measured on 2026-09-13 against the delivery target, one concurrent package run
each: the Daemon package takes 3.2s at the target and 6.8s with the provenance
setup in place. Production is unaffected — it resolves the authorization once
per Implement and once per Settle — so the cost is test setup alone.

## Project Constraints

- Identifier strategy: applicable — preserve existing test, fixture and helper identities; no new identifier is introduced. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, network or HTTP surface is touched. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0056 has Spec Runs separate Task Capacity and Verification Capacity, and neither is widened here; ADR-0148 keeps every authored Verification command able to fail against the unchanged tree. Source: `docs/agents/spec-routing.md`, `docs/agents/domain.md`.
- Tooling authority: applicable — no protected tooling mutation proposed or authorized. The change is a test fixture in `internal/daemon`, which is not a Governed Path. Source: `docs/agents/agent-instructions.md`.

## Goals

- The Daemon package's wall clock returns to the delivery target's range while
  every journey keeps reading real Git provenance.
- The provenance a test needs is created once for the package, not once per
  test.

## User Stories

1. As the maintainer, I want the complete Verification to finish within its
   budget, so that a correctness repair does not arrive as a wall-clock
   regression.
2. As the Supervisor, I want a later change that reintroduces per-test process
   setup to fail a named test, so the cost cannot creep back silently.

## Core Features

1. The per-test base fixture performs no Git repository initialization. The
   committed seed exists once for the package and every test reuses it.
2. Every journey still reads committed bytes through the real authorization
   reader. No test passes because the reader was weakened, stubbed or given a
   non-Git shortcut.
3. The external and symlinked Spec Root journeys keep passing, since they
   exercise a different root shape and are why provenance was added.
4. A named test asserts the seed is created once and reused, so reintroducing
   per-test setup fails at that test rather than in a wall-clock budget.

## User Experience

Nobody sees a new flag or prompt. The complete Verification finishes faster and
every existing assertion still holds.

## Non-Goals / Out of Scope

- The repository-wide suite wall-clock goal, which the performance campaign
  owns. This Spec removes the regression it introduced, not the pre-existing
  budget gap.
- Any production change. Production already resolves the authorization once per
  Implement and once per Settle.
- Changing what any journey asserts about settlement, staging or commit content.

## Success Metrics

- The Daemon package's wall clock in one concurrent run is within the delivery
  target's range rather than roughly double it.
- A search of the per-test base fixture body finds no repository
  initialization, where it finds three Git calls today.

## Decisions

- Delivered as its own Spec rather than inside Spec 0132, because 0132's QA gate
  settled with a passing verdict and work added under a passing gate would make
  that pass certify what the gate never saw. The loader refuses that, correctly.
- Delivered on the same branch as Specs 0119 and 0132, before any of them
  reaches the target, so the regression never lands.

### Declared intentional breaks

1. The per-test base fixture stops creating its own repository and reads a
   shared one instead. A test that depended on having a private repository per
   invocation changes with it.

Everything else keeps behaving as it does: every journey's settlement, staging
and commit assertions are unchanged, the authorization reader is untouched, and
no deadline, budget, parallelism setting or skip moves.

### Regression locks

- The measured target and candidate package times are recorded above and are the
  reference the acceptance compares against.
- The new named test is the guard: it fails if per-test setup returns, so the
  cost cannot creep back without a test going red.

## Open Questions

None.
