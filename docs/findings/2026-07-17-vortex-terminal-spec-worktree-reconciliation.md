---
status: pending
created_at: 2026-07-17
updated_at: 2026-07-17
---

# Vortex spec Runs — integrated work remained in terminal worktrees (2026-07-17)

Branch cleanup in `gesttione-solutions/vortex` found terminal spec Runs whose commits were
already present in the branch being preserved, but whose Run Worktrees and Run Branches still
required manual removal. No Active Run existed. The cleanup targeted
`ma/0004-alinhamento-oraculum-s2s` at `6d628739` after the implementation and review of specs
`0004-alinhamento-oraculum-s2s` and `0005-tokens-mcp-no-axis`.

Environment:

- Vortex checkout: `/Users/marcio/dev/vortex`.
- Preserved branch: `ma/0004-alinhamento-oraculum-s2s` at `6d628739`.
- Terminal Runs:
  - `run_20260717T013321Z_ebb6b69d19a998ee` — `Unresolved`, spec
    `0004-alinhamento-oraculum-s2s`, Run Branch tip `45933d18`.
  - `run_20260717T045329Z_58870573bec93323` — `Unresolved`, spec
    `0005-tokens-mcp-no-axis`, Run Branch tip `1b6d200b`.
  - `run_20260717T063323Z_64a0225939ed9b67` — `Unresolved`, spec
    `0005-tokens-mcp-no-axis`, Run Branch tip `765834e1`.

## 1. Terminal Run Worktrees survived after their commits reached the target branch

- **Symptom / evidence**: `roundfix runs list --state all --limit 0` reported all three Runs as
  terminal `Unresolved`. `git worktree list --porcelain` still reported Run Worktrees for
  `run_20260717T045329Z_58870573bec93323` and
  `run_20260717T063323Z_64a0225939ed9b67`; the older Run Branch for
  `run_20260717T013321Z_ebb6b69d19a998ee` also remained locally. Both retained worktrees were
  clean. For each Run Branch, this command returned exit code `0`:

  ```text
  git merge-base --is-ancestor <run-branch> ma/0004-alinhamento-oraculum-s2s
  ```

  The Run Branch commits were therefore already reachable from the branch the user chose to
  preserve.
- **Root cause**: `internal/worktree/worktree.go` calls `branchHasCommitsBeyondBase` during
  `PruneTerminalReport`. That check compares the Run Branch tip only with its creation base. Any
  changed tip prevents reaping, even when the complete Run Branch is already an ancestor of the
  target branch. Roundfix does not reconcile a terminal Run Worktree after its commits reach the
  target branch through a later Run or manual integration.
- **Action / suggestion**: extend terminal-worktree reconciliation with an ancestry check against
  the Run's recorded target branch. A terminal Run Worktree can be removed when its worktree is
  clean and its Run Branch tip is already reachable from that branch. Keep the current refusal
  when the worktree is dirty, the target branch is unknown, or ancestry cannot be proven.

## 2. The default Run listing hid retained workspaces behind an active-state filter

- **Symptom / evidence**: after the implementation Runs ended, the repository-scoped default
  listing printed:

  ```text
  No Runs found.
  (18 terminal Run(s) hidden; use --state all)
  ```

  At the same time, `git worktree list --porcelain` still contained two Roundfix-owned Run
  Worktrees. The output correctly stated that no Active Run existed, but it did not distinguish
  harmless terminal history from terminal Runs that still retained branches or worktrees.
- **Root cause**: `internal/cli/runs.go` filters rows only by Run state. Its hidden-state note
  counts terminal Runs but does not inspect whether a terminal Run still owns a Run Worktree or
  Run Branch.
- **Action / suggestion**: surface retained-workspace state in `runs list`. When the Active view is
  empty but terminal Runs still have Run Worktrees, the note can name that count and point to the
  terminal/reconciliation view. Terminal rows can also expose whether the Run Worktree is kept,
  missing, or safe to remove.

## 3. No first-party command reconciled the safe terminal workspace state

- **Symptom / evidence**: the cleanup required direct Git commands after proving ancestry and
  clean worktree state:

  ```text
  git worktree remove <run-worktree>
  git branch -d <run-branch>
  ```

  `roundfix gc` does not cover this state: its current contract removes eligible Run Event Journal
  rows and artifact directories only. `roundfix stop --force` is reserved for dead, stuck, or
  runaway Active Runs, while the preflight reaper skips Run Branches with commits beyond their
  creation base.
- **Root cause**: the CLI has cleanup mechanisms for empty terminal debris and retained Run
  storage, but no supported reconciliation path for a terminal Run Branch whose commits were
  integrated elsewhere after the Run ended.
- **Action / suggestion**: add a first-party dry-run and apply flow that reports terminal Run
  Worktrees as `safe`, `unintegrated`, `dirty`, or `unknown`. The apply path must remove only
  `safe` worktrees and branches, print the evidence used, and leave every ambiguous case intact.
  This keeps users and supervising Agents out of manual Git cleanup while preserving failed work.

## What worked — keep

- Terminal Runs preserved their work until integration could be proven.
- Run IDs, Run Branch names, worktree paths, and terminal outcomes made the residue traceable.
- Standard Git ancestry checks proved that cleanup would not discard unique commits.
- `git branch -d` refused unsafe branch deletion by default and completed only after ancestry was
  established.

## Planned resolution

Spec [0038 Terminal Run Worktree reconciliation](../specs/0038-terminal-run-worktree-reconciliation/_prd.md)
owns all three findings: proof-based classification, retained-worktree
visibility in Run discovery, and the dry-run-first Reconcile Command. It
depends on spec
[0037 Terminal outcome integrity](../specs/0037-terminal-outcome-integrity/_prd.md)
for guarded Integration Pending reconciliation.
