---
spec: 0157-reconciliation-and-one-repository-identity
status: active
created: 2026-09-24
surfaces: [backend, cli]
---

# Reconciliation and one repository identity

## Executive Summary

Refuse a cleanup candidate that revalidates without evidence, compare QA Reports
by identity rather than by path root, and derive the repository identity from
the main worktree so every linked worktree shares its artifact directory and Run
records.

## Project Constraints

- Identifier strategy: applicable — repository identity moves from the checkout
  path to the main worktree's root; the main checkout's identity is unchanged.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local Git and filesystem only; no
  credential and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0080, ADR-0091, ADR-0093, ADR-0096,
  ADR-0097, ADR-0104, ADR-0117, ADR-0130, ADR-0155 and ADR-0156 hold. Source:
  `docs/agents/domain.md`.
- Tooling authority: not applicable — the intersection with
  `internal/speccheck/governed.go` is empty. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## The unproven candidate

`ApplyRunBranchCandidate` revalidates a candidate through `InspectTerminalRun`
and then calls `cleanupTerminalRun`, which dereferences the result's evidence.
When revalidation returns a non-releasable state with no evidence, the candidate
is refused with a reason naming the failed revalidation, before any Git command
that removes a worktree or branch.

## Report identity

`supersedingQAReport` compares the Run-side and default-side newest reports via
`NewestQAReportFromPaths`, which breaks ties on the full path. The comparison
instead keys a report by its file name within the Spec's QA directory, so
`docs/specs/<slug>/qa/<name>` and `docs/history/specs/<slug>/qa/<name>` are the
same report.

## Repository identity

`repoID` in `internal/config/config.go` hashes `filepath.Clean(gitRoot)`. It
becomes a hash of the main worktree's root, resolved through the repository's
common Git directory, so the main checkout keeps its current identity and every
linked worktree maps to it. The Run store's per-repository lookups use the same
identity, so a Run recorded from a linked worktree is found from the main
checkout. Artifact directories previously created under a worktree-derived
identity remain readable and are included when the repository's Runs are
listed.

## API Contracts

1. `reconcile --apply` refuses a candidate that revalidates without evidence and
   removes nothing for it.
2. A report and its archived copy are the same report for supersession.
3. Every worktree of a repository reports the same identity, equal to the main
   checkout's current identity.

## Coverage Map

- Goal 1 → The unproven candidate; API Contract 1.
- Goal 2 → Report identity; API Contract 2.
- Goal 3 → Repository identity; API Contract 3.
- Core Feature 1 → The unproven candidate.
- Core Feature 2 → Report identity.
- Core Feature 3 → Repository identity; API Contract 3.
- Core Feature 4 → Repository identity.
- Success Metric 1 → Testing Approach 1.
- Success Metric 2 → Testing Approach 2.
- Success Metric 3 → Testing Approach 3.
- API Contracts 1-3 → The unproven candidate, Report identity, Repository identity.

## Integration Points

- **Spec 0154.** Its fallback path is where the unproven candidate arises.
- **`internal/store`.** Per-repository Run lookups adopt the shared identity.
- **Spec 0156.** The delivery queue lists and resumes Runs by repository.

## Testing Approach

1. **Unproven candidate.** A candidate whose revalidation yields no evidence is
   refused, and no worktree or branch is removed. Fails on the tree as it
   stands, where cleanup dereferences nil evidence.
2. **Archived copy.** A Run-side report and a same-named archived report are
   treated as the same report.
3. **Shared identity.** A linked worktree and the main checkout report the same
   identity, equal to the main checkout's pre-change identity; a Run recorded
   from the linked worktree is listed from the main checkout.
4. **Earlier records.** An artifact directory created under a worktree-derived
   identity remains readable and listed.
5. **Repository gate.** The terminal QA Task records the Daemon's repository
   Verification result as a fact.

## Build Order

1. Refuse the unproven candidate (depends on: none).
2. Report identity independent of root (depends on: 1).
3. Shared repository identity and reachable earlier records (depends on: 2).
4. Terminal QA (depends on: 1, 2, 3, 5, 6).
5. Keep each Run's checkout, share only the identity (depends on: 3).
6. Same report means same content (depends on: 5).

## Risks & Considerations

- **Moving the main checkout's data.** The identity is anchored on the main
  worktree precisely so the main checkout's hash does not change; a test pins it.
- **Hiding a worktree's Runs.** Records created before this change under a
  worktree-derived identity must stay reachable, or the fix would lose the very
  Runs it exists to show.
