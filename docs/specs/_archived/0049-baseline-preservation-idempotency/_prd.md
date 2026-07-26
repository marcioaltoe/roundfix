---
spec: 0049-baseline-preservation-idempotency
status: archived
created: 2026-07-26
archived: 2026-07-26
qa_override: true
surfaces: [cli, infra]
---

# Baseline Preservation idempotency

## Problem

After an approved Preservation semantically redistributes existing root
instructions, a fresh Baseline Plan inventories the original root bytes and
recognized repository-rule carriers again. It can therefore propose another
backup, duplicate retained rules, and update generated guidance even though the
repository has not received new instructions.

## Project Constraints

- Identifier strategy: not applicable — this correction creates no Internal
  Identifier or application identity. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the correction is confined to
  local Baseline planning and changes no authentication provider, HTTP
  contract, or route. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — byte-exhaustive accounting,
  fail-closed retention, bounded carrier mutation, portable planning, and
  recoverable apply remain binding; ADR-0078 replaces only the live-root
  retention outcome after confirmed redistribution. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — no protected tooling mutation is proposed or
  authorized. Source: `docs/agents/agent-instructions.md`.

## Goals

- Make the first approved Preservation the sole owner-transfer operation for
  the root bytes it accounts for.
- Make an unchanged subsequent Preservation produce zero file changes and no
  new backup.
- Preserve strict accounting, confirmation, and rollback guarantees when new
  unmarked instructions are later added to `AGENTS.md`.

## User Stories

1. As a maintainer, I want redistributed root rules removed from the live root
   carrier, so that every accepted rule has one operative owner.
2. As a maintainer, I want an unchanged Preservation rerun to be empty, so that
   updating the Baseline is safe and predictable.
3. As a maintainer adding new root instructions, I want only those new bytes
   reviewed and migrated, with an immutable backup of that source version.

## Core Features

1. Preservation classification excludes setup-managed blocks and already
   retained repository-owned rule carriers.
2. An approved Preservation postimage rebuilds `AGENTS.md` from active managed
   root blocks and confirmed root-level repository-owned content, excluding
   source bytes whose dispositions were accepted.
3. The immutable backup contains the complete pre-migration root carrier and
   remains addressed by its content digest.
4. A new unmarked addition produces one new review and one backup for its new
   complete root content identity; after apply, the live root converges again.
5. Apply and recovery continue to validate exact preimages, postimages,
   dispositions, retention entries, and backup identities.

## Non-Goals

- Rewriting or deleting arbitrary nested instruction carriers.
- Changing the public Baseline Plan, Result, confirmation, or exit-code
  schemas.
- Reclassifying repository-owned semantic blocks that already have an accepted
  owner.
- Changing Baseline Profile content or protected repository tooling.
