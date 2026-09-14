---
spec: 0137-a-fixture-that-carries-its-committer
status: active
created: 2026-09-14
---

# TechSpec — A fixture that carries its committer

## Executive Summary

The shared Task-cycle fixture seed initialises a repository and commits into it
through the test helper, which supplies a committer identity as a
per-invocation override. Overrides are not written to disk, and repository
hardening appends only maintenance keys, so a copy of the seed has no identity
of its own. The Daemon's own committer runs plain Git and depends on host
discovery.

Write the identity into the seed's configuration, where the setup 0133 replaced
used to write it, and assert it by reading the copied repository's resolved
configuration rather than by trusting the host.

## Project Constraints

- Identifier strategy: not applicable — no identifier, refusal code or record field is introduced or changed. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no credential, network or HTTP surface is touched. Source: `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0057 keeps the Daemon the exclusive writer of Task status, untouched here. Source: `docs/agents/spec-routing.md`, `docs/agents/domain.md`.
- Tooling authority: applicable — no protected tooling mutation proposed or authorized. The repair is ordinary test source in `internal/daemon`. Source: `docs/agents/agent-instructions.md`.

## System Architecture

| Component | Responsibility | Change |
| --- | --- | --- |
| Task-cycle fixture seed | Build one repository once and copy it per test | Write a committer identity into the repository configuration |
| Repository hardening helper | Disable automatic maintenance in a repository | None; it owns maintenance keys, not identity |

## Implementation Design

### Write the identity where a copy can inherit it

The helper that runs Git in tests passes the identity as a command-line
override on each invocation, which is why the seed's own commits succeed and a
copy's do not. The helper that appends configuration is the one that persists
settings, and the replaced per-test setup used it for exactly this. Use it for
the seed, and the property travels with every copy.

Assert the result by asking the copied repository what identity it resolves,
with host and global configuration excluded from the answer. An assertion that
merely commits successfully passes on any machine whose Git can invent an
identity, which is the whole reason this reached a pull request gate.

## Coverage Map

- PRD Goal 1 → Task-cycle fixture seed.
- Core Features 1-2 → Task-cycle fixture seed.
- Core Feature 3 → Task-cycle fixture seed, by leaving the once-only creation
  intact.

## Testing Approach

1. A repository copied from the seed resolves a committer identity from its own
   configuration, with host and global configuration excluded.
2. The Daemon's real-repository Task-cycle journey still commits per Task and
   still excludes pre-existing dirt.
3. The seed is still created once.

Observation 1 rests on evidence this Spec did not author: Git's own identity
resolution, and the fact that a Linux runner answers `Author identity unknown`
where macOS composes an identity from the local user and hostname. The pull
request gate's recorded failure is the measurement. If that exclusion cannot be
expressed, the row records blocked with its reason.

## Build Order

1. Write the committer identity into the shared fixture seed and assert it from
   a copy (depends on: none).
2. Terminal QA (depends on: 1).

## Risks & Considerations

- Asserting through a successful commit would reproduce the blind spot. The
  assertion must read the resolved configuration with the host excluded.

## Decisions

- Restore the property in the seed rather than reintroducing per-test setup, so
  0133's removed cost stays removed.

## Vocabulary Contract

No token is coined.
