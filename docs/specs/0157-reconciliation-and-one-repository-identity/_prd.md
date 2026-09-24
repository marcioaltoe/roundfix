---
spec: 0157-reconciliation-and-one-repository-identity
status: active
created: 2026-09-24
surfaces: [backend, cli]
---

# Reconciliation and one repository identity

Two defects shipped with Spec 0154 on the maintainer's decision to record them
rather than take a third corrective Task, and a third surfaced the moment the
restructured queue started running Specs in parallel.

**`reconcile --apply` can panic.** Spec 0154's absent-target fallback can make a
clean Run Branch a cleanup candidate. On revalidation `InspectTerminalRun`
returns an unintegrated result without evidence, and `cleanupTerminalRun` reads
`fresh.evidence.worktreePresent` through a nil pointer. The path is new since
0154, so the crash is a regression in a destructive command.

**An archived copy of the same QA Report is not recognised.** Archive moves a
report from `docs/specs/<slug>/qa` to `docs/history/specs/<slug>/qa` without
renaming it. Run-side and default-side reports then share date and sequence, the
tie-break is by full path, and the active path wins — so valid archived evidence
is missed and the Run stays preserved.

**Each worktree is a different repository.** `repoID` in
`internal/config/config.go` hashes the checkout path, and the Run store keys Runs
by the Git root. Running Specs 0155 and 0156 from linked worktrees on
2026-09-24 produced artifact directories `8917c132…` and another, apart from the
main checkout's `339f8dac…`, and Spec 0155's Run does not appear in
`roundfix runs list` from the main checkout. Parallel delivery and the durable
queue both depend on one repository having one identity.

This Spec is carved from delivery 3 of the restructured queue: the two limits
Spec 0154 recorded, and Spec 0125 Core Feature 2.

## Project Constraints

- Identifier strategy: applicable — the repository identity becomes that of the
  main worktree, shared by every linked worktree; the main checkout's existing
  identity is unchanged, so its records need no migration. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local Git and filesystem only; no
  credential and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0093 checks Spec consistency by
  citation, ADR-0104 accepts on evidence a Spec did not author, ADR-0130 keeps a
  path governed once bounded, ADR-0155 makes the `qa` Task declare the matrix and
  ADR-0156 makes a declared promise name a consuming Task. This Spec's gate is
  bound by ADR-0080, ADR-0091, ADR-0096, ADR-0097 and ADR-0117. All hold.
  Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — the intersection of this Spec's changed
  paths with `internal/speccheck/governed.go` is empty. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- `reconcile --apply` cannot crash on a candidate it cannot prove.
- Archived evidence is recognised as the evidence it is.
- A repository has one identity, whichever of its worktrees a command runs in.

## Core Features

1. **A candidate without evidence is refused, not cleaned.** When revalidation
   of a cleanup candidate returns no evidence, the candidate is refused with a
   named reason before any cleanup runs.
2. **Report identity survives archiving.** Comparing QA Reports for supersession
   treats a report and its archived copy as the same report, independent of the
   active or archived root.
3. **One identity per repository.** The repository identity derives from the
   main worktree, so linked worktrees share its artifact directory and its Run
   records; a Run started from a linked worktree is listed from the main
   checkout.
4. **Earlier worktree-derived records stay reachable.** Artifact directories
   created under a worktree-derived identity remain readable and are listed
   with the repository's Runs.

## Non-Goals / Out of Scope

- The verification and settlement features of delivery 3 (Spec 0122 Core
  Features 5 to 8, Spec 0129 Core Features 1, 4 and 5, Spec 0123 Core Feature 3),
  which follow in their own Spec.
- Changing the branch-naming policy or Run namespace.
- Changing what reconciliation releases when evidence is positive.

## Success Metrics

1. A cleanup candidate whose revalidation returns no evidence is refused and no
   worktree or branch is removed, where today the command panics.
2. A Run-side report and its archived copy with the same date and sequence are
   recognised as the same report.
3. A Run started from a linked worktree appears in `roundfix runs list` from the
   main checkout, and the main checkout's identity is unchanged.

## Recorded limits

- A Run recorded before this Spec from a linked worktree that has since been
  removed stays unlisted: its checkout no longer resolves to the repository, and
  the maintainer's corrective ceiling of two Tasks was spent on the defects the
  pre-PR review of 2026-09-24 found. Reproduction: start a Run from a linked
  worktree on the previous release, remove the worktree, run `roundfix runs list`
  from the main checkout. A durable repository key per Run is the carried fix.

## Decisions

- **Anchor on the main worktree.** Deriving identity from the main worktree's
  root keeps the main checkout's existing identity, so its records need no
  migration; only records created from linked worktrees need to stay reachable.
- **Refuse rather than guess.** A cleanup candidate that cannot be proven at the
  moment of cleanup is refused, matching the rule that nothing is released on an
  absence of evidence.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in the
Task Graph. The negative cases carry the weight: an identity change that moved
the main checkout's artifacts, or a refusal that still removed a worktree, would
pass the happy path while doing damage.

## Research basis

The two reconciliation defects were reported by independent pre-PR review of
Spec 0154 and recorded in its PRD. The identity defect was observed directly on
2026-09-24: `repoID` hashes `filepath.Clean(gitRoot)`, runs launched from
linked worktrees wrote to separate artifact directories, and one of them was
absent from `roundfix runs list` in the main checkout.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
