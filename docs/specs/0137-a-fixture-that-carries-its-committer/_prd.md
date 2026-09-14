---
spec: 0137-a-fixture-that-carries-its-committer
status: active
created: 2026-09-14
surfaces: [backend]
---

# A fixture that carries its committer

Spec 0133 replaced the Daemon Task-cycle fixture's per-test repository setup
with one shared seed copied per test. The setup it replaced wrote a committer
identity into the repository's own config; the seed does not. Repository
hardening writes only maintenance keys, and the identity the test helper passes
is a per-invocation override that no copy inherits.

Every commit the Daemon makes in those fixtures therefore relies on whatever
identity Git can discover from the machine. macOS composes one from the local
user and hostname, so the suite passes locally. A Linux runner cannot, and
answers `Author identity unknown` with exit 128. The delivery's own pull request
gate is where this first appeared: three local `make verify` runs, the
repository docs gate, and a passing QA gate could not see it.

## Project Constraints

- Identifier strategy: not applicable — no identifier, refusal code or record field is introduced or changed. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, network or HTTP surface is touched. The committer identity is test fixture configuration, not a credential. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0057 keeps the Daemon the exclusive writer of Task status, which this Spec does not touch. Source: `docs/agents/spec-routing.md`, `docs/agents/domain.md`.
- Tooling authority: applicable — no protected tooling mutation proposed or authorized. The repair is ordinary test source in `internal/daemon`. Source: `docs/agents/agent-instructions.md`.

## Goals

- A Daemon Task-cycle fixture repository carries a committer identity in its own
  configuration, so a Daemon commit succeeds on a machine where Git can
  discover none.

## Core Features

1. The shared fixture seed writes a committer identity into the repository
   configuration, so every copy of the seed inherits it.
2. The identity is asserted as a value present in the copied repository's
   configuration, not as an assumption about the host.
3. The seed is still created once, so the cost 0133 removed stays removed.

## Non-Goals / Out of Scope

- Any change to the Daemon's commit behavior, to what it stages, or to what
  authority a governed change requires.
- Other fixtures and helpers that create repositories, unless they share this
  seed.
- The suite's wall-clock budget, which the test-performance campaign owns.

## Declared intentional breaks

None. This restores a property the replaced setup had, and adds an assertion
that it is present.

## Regression locks

- The identity is asserted by reading the copied repository's own resolved
  configuration, so the assertion fails on a machine that would otherwise mask
  the defect with a discoverable identity.
- The seed-created-once assertion that 0133 established keeps passing.
