---
spec: 0133-a-fixture-that-does-not-spawn-per-test
status: archived
created: 2026-09-13
surfaces: [backend]
archived: "2026-09-13"
source_slug: 0133-a-fixture-that-does-not-spawn-per-test
---


# The regressions this branch must not ship

Specs 0119 and 0132 each repaired a real defect and each introduced a new one.
Neither may reach the delivery target, so both are repaired here, on the same
branch, before any of the three Specs merges.

**The blocking one.** The operation enforcement Spec 0119 delivered requires the
`implement` operation from every Spec, and a record with an empty `paths` list
is invalid. Together: a Spec that touches no Governed Path can neither write a
valid record nor be implemented. Spec 0131 ran Clean in exactly that condition
before the enforcement landed, and would be undispatchable today. The canonical
clause says an absent record "grants nothing", and the implementation obeyed it
literally — conflating authority over protected tooling with authority to do the
work at all.

**The cost one.** Giving the Daemon's Task-cycle fixtures real Git provenance
fixed a correctness defect and introduced a cost one. The base fixture that every Task-cycle test
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
- Tooling authority: applicable — express maintainer authorization: "Aprovar os três caminhos", 2026-09-13, recorded in `docs/specs/0133-a-fixture-that-does-not-spawn-per-test/_authorization.md`; bounded files: `internal/baseline/assets/modules/core.json`, `docs/agents/agent-instructions.md`, `docs/agents/setup-context.json`. Sanctioned regeneration: `make baseline-digests` follows the approved module edit. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- A Spec that touches no Governed Path can be implemented without a maintainer
  record, while every change to a Governed Path still requires one.
- The per-test process setup this branch introduced is removed, and a named
  test keeps it from returning, while every journey keeps reading real Git
  provenance. Closing the remaining wall-clock gap to the delivery target is
  not this Spec's goal; see Decisions.
- The canonical rule and the implementation say the same thing about what an
  absent record withholds.

## User Stories

1. As the Supervisor, I want to implement a Spec that touches no protected
   tooling without waiting for a maintainer record, so that ordinary work does
   not need an approval it has no reason to need.
2. As the maintainer, I want a change to a Governed Path to keep refusing
   without an operative record, so narrowing the rule costs no authority.
3. As the maintainer, I want the complete Verification to finish within its
   budget, so that a correctness repair does not arrive as a wall-clock
   regression.
4. As the Supervisor, I want a later change that reintroduces per-test process
   setup to fail a named test, so the cost cannot creep back silently.

## Core Features

1. An absent, proposed, contradictory or withdrawn record withholds governed
   mutation and does not by itself refuse implementation, commit or push. The
   canonical clause states that division, and the implementation matches it.
2. Every change to a Governed Path still requires an operative record naming
   that exact path. Narrowing which actions an absent record blocks removes no
   authority from the tooling boundary.
3. A Spec that declares no protected tooling mutation dispatches with no record
   present, as Spec 0131 did before the enforcement landed.
4. The per-test base fixture performs no Git repository initialization. The
   committed seed exists once for the package and every test reuses it.
5. Every journey still reads committed bytes through the real authorization
   reader. No test passes because the reader was weakened, stubbed or given a
   non-Git shortcut.
6. The external and symlinked Spec Root journeys keep passing, since they
   exercise a different root shape and are why provenance was added.
7. A named test asserts the seed is created once and reused, so reintroducing
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

- The Daemon package's wall clock in one concurrent run is measurably lower
  than before the fixture repair, and the reduction is recorded with both
  measurements rather than asserted.
- A search of the per-test base fixture body finds no repository
  initialization, where it finds three Git calls today.

## Decisions

- Amended on 2026-09-13, after measurement falsified the original Goal. The
  Goal claimed the Daemon package's wall clock returns to the delivery target's
  range. Measured on the same host with warmed caches: the delivery target runs
  3.63s and the candidate 6.48s, and the immediate pre-repair control ran 7.70s.
  The fixture repair is therefore real and partial — it recovered about 1.2s of
  a 2.8s gap — and the remaining majority has a cause this Spec has not
  identified. Two earlier diagnoses of that gap were wrong: archive traversal,
  which the suite guard short-circuits when no violation exists, and fixture
  spawning alone, which accounts for the smaller part.
- The residual wall-clock gap is promoted to the repository's test-performance
  campaign, which owns suite wall clock, with both measurements and the partial
  decomposition above as its input. The maintainer took that decision on
  2026-09-13 rather than accept a third unmeasured diagnosis inside a branch
  already carrying four Specs. The impact the terminal QA gate assigned to the
  residual is friction, not blocked completion.
- Delivered as its own Spec rather than inside Spec 0132, because 0132's QA gate
  settled with a passing verdict and work added under a passing gate would make
  that pass certify what the gate never saw. The loader refuses that, correctly.
- Delivered on the same branch as Specs 0119 and 0132, before any of them
  reaches the target, so the regression never lands.

### Declared intentional breaks

1. An absent record stops refusing implementation, commit and push. A caller
   that relied on the absent record refusing everything sees those three
   proceed; every governed-path refusal is unchanged.
2. The per-test base fixture stops creating its own repository and reads a
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
