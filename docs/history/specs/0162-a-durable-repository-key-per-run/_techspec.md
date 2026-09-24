---
spec: 0162-a-durable-repository-key-per-run
status: active
created: 2026-09-24
surfaces: [backend, cli]
---

# A durable repository key per Run

## Executive Summary

Store each Run's repository key at creation, backfill existing rows where it
resolves, and make listing, reconcile, gc and Run Window lookups use it.

## Project Constraints

- Identifier strategy: applicable — a recorded repository key per Run. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local Git and the Run Database
  only. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0080, ADR-0091, ADR-0093, ADR-0096,
  ADR-0097, ADR-0104, ADR-0117, ADR-0130, ADR-0155 and ADR-0156 hold. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization recorded in
  [_authorization.md](_authorization.md); bounded files:
  `internal/cli/cli_test.go`. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`.

## The key

A schema migration after the current version adds `runs.repository_root`,
filled at creation from `config.RepositoryRoot(gitRoot)`. The migration
backfills rows whose `git_root` still resolves and leaves the rest empty. Fresh
and migrated schemas stay equivalent.

## Consumers

`ListRuns` and the settle and carry-forward lookups in `internal/store/store.go`
match `repository_root` equal to the current key, or, for rows with an empty
key, the existing `git_root` alias match. `sameRepository` in
`internal/cli/reconcile.go` compares repository keys. `classifyGCSanitationRoot`
in `internal/cli/gc.go` derives the default root from the recorded key. The Run
Window functions fall back to the exact checkout's row when none exists under
the key, and clear both.

## API Contracts

1. `runs.repository_root` is recorded for every new Run.
2. Listing, reconcile, gc and Run Window lookups agree on one repository key.

## Coverage Map

- Goal 1 → The key; Consumers; API Contract 1.
- Goal 2 → Consumers; API Contract 2.
- Core Feature 1 → The key.
- Core Features 2-5 → Consumers.
- Success Metric 1 → Testing Approach 1.
- Success Metrics 2-4 → Testing Approach 2.
- API Contracts 1-2 → The key, Consumers.

## Integration Points

- **Spec 0157.** Owns the query-time identity this Spec makes durable.
- **Spec 0156.** Its migrations precede this one.

## Testing Approach

1. **Key and migration.** A new Run records its key; migration from every
   earlier version backfills a resolvable row and leaves an unresolvable one
   empty; fresh and migrated schemas match. With real Git, a Run from a linked
   worktree is listed from the main checkout after `git worktree remove`.
2. **Consumers.** With real Git, `reconcile <run-id>` from the main checkout
   accepts a linked-worktree Run and still refuses another repository's Run;
   `gc --sanitize` does not classify the shared root `overridden` after the
   worktree is removed; a window keyed on the worktree path is shown and cleared.
3. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. The key and its migration (depends on: none).
2. Consumers (depends on: 1).
3. Terminal QA (depends on: 1, 2, 4, 5).
4. Every lookup, and only the recorded key (depends on: 2).
5. Removed worktrees and bare layouts (depends on: 4).

## Risks & Considerations

- **Claiming another repository's Run.** The negative cases prove a Run of an
  unrelated repository is neither listed, reconciled nor reclaimed.
