---
spec: 0162-a-durable-repository-key-per-run
status: active
created: 2026-09-24
surfaces: [backend, cli]
---

# A durable repository key per Run

Spec 0157 gave a repository one identity across its worktrees, but found that
identity at query time, through the `.git/worktrees/*/gitdir` entries that
`git worktree remove` deletes. Its second pre-PR review recorded four limits
that follow from that choice:

- A Run recorded from a linked worktree drops out of `roundfix runs list`,
  `reconcile` and settle lookups from the main checkout once its worktree is
  removed, which is the normal end of a parallel delivery.
- `reconcile <run-id>` from the main checkout refuses a linked-worktree Run
  that `implement` there proposes for carry-forward, because its repository
  check compares checkout paths.
- `gc --sanitize` classifies the shared artifact root `overridden` after such a
  worktree is removed, and then reclaims nothing from it.
- A Run Window set from a linked worktree before Spec 0157 is keyed on that
  worktree's path and is no longer found.

## Project Constraints

- Identifier strategy: applicable — each Run records its repository key, the
  main worktree root at creation, beside its own checkout; existing rows are
  backfilled where the key resolves. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local Git and the Run Database only;
  no credential and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0093 checks Spec consistency by
  citation, ADR-0104 accepts on evidence a Spec did not author, ADR-0130 keeps a
  path governed once bounded, ADR-0155 makes the `qa` Task declare the matrix
  and ADR-0156 makes a declared promise name a consuming Task. This Spec's gate
  is bound by ADR-0080, ADR-0091, ADR-0096, ADR-0097 and ADR-0117. All hold.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the standing grant of 2026-09-21 for governed
  paths a slice needs, recorded in [_authorization.md](_authorization.md);
  bounded files: `internal/cli/cli_test.go`. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- A Run stays with its repository after the worktree it ran in is gone.
- Every command that asks "same repository?" gets one answer.

## Core Features

1. **A recorded repository key.** A new Run records its repository key at
   creation. A migration adds the column and backfills it for existing rows
   whose checkout still resolves; a row that cannot be resolved keeps an empty
   key and its current behaviour.
2. **Listing by key.** Per-repository listing, reconcile selection and settle
   lookups match a Run by its recorded key, falling back to today's checkout
   match only for rows without one.
3. **One repository check.** `reconcile <run-id>` accepts a Run whose recorded
   key, or resolved checkout, equals the current repository's key.
4. **Garbage collection by key.** `gc --sanitize` compares an artifact root with
   the default derived from the Run's recorded key, so a removed worktree does
   not make the shared root look overridden.
5. **Run Windows found from any checkout.** A window lookup that finds none under
   the repository key falls back to a row keyed on the exact checkout, and
   `window clear` removes either.

## Non-Goals / Out of Scope

- Moving artifact directories.
- Changing the per-checkout Active Run lock.

## Success Metrics

1. A Run recorded from a linked worktree is listed from the main checkout after
   that worktree is removed.
2. `reconcile <run-id>` from the main checkout accepts a linked-worktree Run.
3. `gc --sanitize` classifies the shared root as it did before the worktree was
   removed.
4. A Run Window keyed on a linked worktree's path is shown and cleared from it.

## Recorded limits

The corrective ceiling of two Tasks was spent on the defects of the first two
pre-PR reviews of 2026-09-24. The third review found one minor gap, carried to a
future slice:

- In a bare-repository layout, or a linked worktree of a `--separate-git-dir`
  repository, the default Artifact Root moved from the worktree path to the
  common Git directory. `gc --sanitize` then classifies a pre-upgrade
  per-worktree root `overridden` and preserves it instead of reclaiming it; once
  retention prunes those rows, the old directories are no longer found. Nothing
  is deleted. Reproduction: bare clone, `git worktree add proj/main`, record a
  terminal Run on the previous release, upgrade, run `gc --sanitize`. Carried
  fix: accept a recorded root that equals the default derived from either the
  recorded key or the recorded checkout.

## Decisions

- **Record, don't rediscover.** The repository a Run belongs to is a fact at
  creation; rediscovering it later depends on Git metadata that is deleted on
  purpose.
- **Backfill only what resolves.** A row whose checkout no longer resolves keeps
  its current behaviour rather than a guessed key.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in the
Task Graph, with real Git worktrees that are created and removed. The negative
cases carry the weight: a Run of another repository listed or reclaimed under
this one would pass every happy path.

## Research basis

The four limits and their reproductions are recorded in the archived PRD of
Spec 0157, from its second pre-PR review of 2026-09-24.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
