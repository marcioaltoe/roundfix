---
spec: 0038-terminal-run-worktree-reconciliation
status: active
created: 2026-07-17
surfaces: [backend, cli, data, docs]
---

# Terminal Run Worktree reconciliation

Terminal spec Runs can retain clean Run Worktrees and Run Branches after their commits have already reached the user's target branch. Roundfix currently compares the Run Branch only with its creation base, hides the residue behind the default Active Run listing, and offers no supported cleanup command. Prior dogfood evidence, now absorbed into this Spec and retained in Git history, showed users having to prove ancestry and run Git commands manually.

## Project Constraints

- Identifier strategy: not applicable — reconciliation reuses Run IDs, recorded
  branches, and Git object identities and creates no project-owned Internal
  Identifier. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — reconciliation is confined to the
  local Run Database and Git repository and adds no authentication or HTTP
  contract. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0023 and ADR-0024 govern Run
  Worktrees and porcelain integration, ADR-0052 owns the guarded terminal
  transition, and ADR-0053 requires proof-based reconciliation. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — on 2026-07-26, the maintainer expressly
  authorizes changes to exactly `.agents/skills/roundfix/SKILL.md` and
  `skills/roundfix/SKILL.md`; no other protected tooling mutation is
  authorized. Source:
  `docs/agents/agent-instructions.md`.

## Goals

- Users and Supervisors can distinguish safely integrated terminal Run Worktrees from work that must be preserved.
- The Reconcile Command is read-only by default and removes only work whose safety is positively proven.
- Run discovery exposes retained Run Worktrees without requiring direct Git inspection.
- An Integration Pending Run becomes Clean only through recorded integration evidence, while every other historical terminal outcome remains unchanged.

## User Stories

1. As a user finishing work that was integrated manually or by a later Run, I want Roundfix to prove whether the retained Run Worktree is safe to remove, so that I do not repeat fragile Git archaeology.
2. As an Agent, I want a machine-readable reconciliation report before applying cleanup, so that I can act only on entries classified safe.
3. As a user running `runs list`, I want retained terminal Run Worktrees called out even when the default view contains no Active Runs, so that storage and pending work are not hidden.
4. As a user applying repository-wide reconciliation, I want dirty, unintegrated, and unknown entries preserved, so that bulk cleanup cannot discard work.
5. As a user who completed an Integration Pending command, I want the Run reconciled to Clean with the proof recorded, so that its current state reflects completed integration without losing its original outcome history.

## Core Features

1. Every terminal spec Run with a recorded Run Worktree or Run Branch is classified as `safe`, `unintegrated`, `dirty`, `unknown`, or `released`.
2. `safe` requires a clean retained worktree, a resolvable Run Branch tip and recorded target branch, and proof that the Run Branch tip is an ancestor of the current target branch tip.
3. `unintegrated` means the worktree is clean and the relevant refs resolve, but the Run Branch tip is not an ancestor of the target branch. `dirty` means the retained worktree has tracked or untracked changes. `unknown` covers missing metadata, Git failures, ambiguous refs, or any case where safety cannot be proven. `released` means neither retained worktree nor Run Branch remains.
4. `roundfix reconcile [run-id]` prints a dry-run report by default. `--apply` removes only `safe` Run Worktrees and Run Branches. Without a Run ID, the command scans terminal spec Runs for the current repository.
5. Text and JSON output report the Run ID, terminal outcome, classification, Run Worktree, Run Branch, target branch, relevant heads, evidence, planned or applied actions, and any refusal reason.
6. Applying a safe Integration Pending entry uses the guarded reconciliation transition from spec 0037 to record Integration Pending → Clean before cleanup. Applying another terminal outcome keeps that outcome and records only the reconciliation evidence and cleanup event.
7. The command is idempotent. Re-running dry-run or apply after successful cleanup reports `released` and makes no Git or Run Database change.
8. `runs list` preserves its existing row shape. When a view hides retained terminal Run Worktrees, or includes terminal Runs that retain them, stderr reports their exact count and points to `roundfix reconcile` for classification.
9. The existing automatic terminal reaper uses the same classifier and may remove only `safe` or already `released` entries. It never uses the old creation-base comparison as sufficient evidence.

## User Experience

```text
roundfix reconcile
run_...  IntegrationPending  safe         roundfix/run-... -> ma/feature
run_...  Unresolved          unintegrated roundfix/run-... -> ma/feature

Dry run: 1 Run Worktree can be released; 1 was preserved.
Apply with: roundfix reconcile --apply
```

`--apply` prints one action per safe Run and one preserved line per non-safe Run. A scan with only preserved entries is successful because no unsafe mutation was attempted. Operational Git or database failures produce a non-zero exit and name the next safe action.

## Non-Goals / Out of Scope

- Integrating unique Run Branch commits into the target branch.
- Repairing dirty Run Worktrees or choosing which local changes to keep.
- Deleting Task Worktrees whose failed Tasks still require settlement.
- Reclaiming Run Event Journal or Artifact Directory storage; the GC Command retains that responsibility.
- Automatic cleanup based only on age, terminal outcome, or a missing path.
- Rewriting terminal outcomes other than the guarded Integration Pending → Clean transition.
- Cross-repository reconciliation through one invocation; the no-ID form is scoped to the current repository.

## Success Metrics

- Every retained terminal Run Worktree fixture produces exactly one of the five documented classifications.
- `--apply` removes all and only `safe` fixtures; zero dirty, unintegrated, or unknown fixtures lose a path or ref.
- A second `--apply` against an already reconciled Run performs zero mutations and reports `released`.
- Run listing reports the exact count of retained terminal Run Worktrees without changing its existing stdout row contract.
- A safe Integration Pending fixture ends Clean with a durable reconciliation event containing both heads and the prior outcome.

## Decisions

- Reconciliation is a dedicated command, dry-run by default, rather than an extension of GC or an automatic-only Preflight side effect.
- Positive cleanliness and ancestry proof are mandatory; ambiguity always preserves work.
- Repository-wide reconciliation is the no-ID default; `--apply` is the only mutation switch.
- See [ADR-0053](../../adr/0053-terminal-run-worktree-reconciliation-is-proof-based.md).

## Open Questions

None.
