---
spec: 0008-worktree-isolation
status: active
created: 2026-07-05
surfaces: [cli, infra]
---

# Worktree Isolation

Every Run today executes inside the user's own checkout. Three chronic
findings share that single root: concurrent user work gets swept into task
commits, the implement preflight must reject any dirty tree (making failed-
Run recovery a chore), and the fragile snapshot-diff dance exists only to
tell the Agent's changes apart from everyone else's. This Spec moves
execution into an isolated per-Run git worktree on a named Run Branch, and
brings the commits back to the user's branch through a porcelain-only
integration protocol proven safe by experiment — the user keeps working in
their checkout while Runs work in theirs.

## Goals

- A Run never reads or writes the user's checkout during execution: Agents,
  Verification, snapshots, and commits all happen in the Run Worktree. See
  ADR-0023.
- Concurrent user edits and commits become invisible to Run commits — the
  multi-writer worktree stance is satisfied structurally, not by warnings.
- The Run's commits reach the user's branch only through safe porcelain:
  fast-forward in their checkout, ancestry-checked branch move when they
  switched away, and an explicit Integration Pending outcome when neither is
  safe. See ADR-0024.
- Failed work lives in exactly one place — the kept Run Worktree — where
  inspection and the Settle Command find it.
- The Live Run View, Attach, and reports read Task state from where
  execution actually happens.

## User Stories

1. As a developer, I want to keep editing and committing in my checkout
   while a Run executes, so that my work and the Run's commits never
   contaminate each other.
2. As a developer starting an implement Run on a dirty tree, I want the Run
   to proceed (with a note about integration implications) instead of being
   rejected, so that a failed attempt's leftovers no longer block progress.
3. As a developer whose Run ends Clean, I want its commits fast-forwarded
   onto my branch and the Run Worktree cleaned up, so that success looks
   exactly like it does today — commits on my branch, nothing extra lying
   around.
4. As a developer whose checkout had overlapping changes or new commits, I
   want the Run to end Integration Pending with the exact command to
   integrate when I am ready, so that nothing is forced over my work and
   nothing is silently lost.
5. As a developer with a failed Task, I want its preserved work in the kept
   Run Worktree and `roundfix settle` pointed there, so that recovery has
   one obvious surface.
6. As a developer watching or attaching to a Run, I want Task statuses read
   from the Run Worktree, so that the cockpit shows execution truth, not my
   stale checkout.
7. As a developer whose verification needs untracked files (env files,
   local certs), I want a config list copied into the Run Worktree at start,
   so that isolation does not break my gate.

## Core Features

1. **Run Worktree lifecycle.** Run start creates the worktree under
   Roundfix Home on the Run Branch at the Run's head commit; the Run row
   records its path; Clean-and-integrated outcomes remove worktree and
   branch; every other outcome keeps both and reports the path. A preflight
   sweep prunes worktree debris of terminal Clean Runs. See ADR-0023.
2. **Execution retarget.** Both engines, the Agent Session cwd, prompts,
   Verification, snapshots, commits, task-file reads and writes, and the QA
   step operate on the Run Worktree; the before-snapshot rule survives
   unchanged inside it (failed-Task dirt still never enters later commits).
3. **Porcelain-only integration.** At Run end (and after settle): user still
   on the Run's branch → `merge --ff-only` in their checkout (succeeds on
   clean and non-overlapping-dirty trees; refuses safely otherwise); user
   elsewhere → ancestry-checked `branch -f`; refusal → Integration Pending
   with the commits on the Run Branch and the command in the report. Final
   Push happens only after successful integration. See ADR-0024.
4. **Preflight recalibration.** The implement dirty-tree rejection becomes a
   stderr note about integration implications; the Active Run lock stays
   keyed on the user checkout; branch detection works in the worktree via
   the named Run Branch.
5. **Readers follow execution.** Live Run View, Attach, the settle command,
   and end-of-run reports read from the recorded Run Worktree (falling back
   to the user root for terminal Runs whose worktree is gone).
6. **Worktree provisioning.** A config list of untracked files to copy into
   the Run Worktree at start (default empty); the Artifact Directory's
   builtin default moves out of the repository tree so sandboxed runtimes
   and snapshots never see it (explicit configs are untouched).

## User Experience

Same commands and flags. New surfaces: the Run header names the Run
Worktree; non-Clean outcomes print its path; Integration Pending prints the
one-line integration command; the implement dirty-tree error becomes a
note. Success is visually identical to today — commits land on the branch,
worktree gone.

## Non-Goals / Out of Scope

- Parallel Task/Batch execution and ready-set scheduling — this Spec builds
  the isolation layer that will host it, nothing more.
- Review artifacts relocation into the spec tree (next Spec; ADR-0003
  untouched beyond the builtin default's location).
- Any change to commit messages, verification semantics, or ownership
  (ADRs 0010/0013/0014 behavior is preserved inside the worktree).
- Merge-conflict resolution — Integration Pending hands divergence to the
  developer by design.

## Success Metrics

- A Run executes while the user's checkout is dirty and receiving commits;
  the Run's commits contain zero user files, and the user's work is
  untouched — proven in tests for both engines.
- Clean outcome: branch fast-forwarded, worktree and Run Branch gone.
- Overlap and divergence fixtures end Integration Pending with the branch
  unmoved and the documented command working when run manually.
- A failed Task's work is inspectable in the kept worktree and settled from
  it by the Settle Command.
- The full existing suite passes with contract changes limited to the
  documented preflight demotion and new report lines.

## Decisions

- Worktree-per-Run on a named Run Branch under Roundfix Home; per-Task
  worktrees wait for parallelism. See ADR-0023.
- Integration is porcelain-only with the two-case protocol and Integration
  Pending as a new terminal state (exit code 1 family — the Run is not
  fully delivered). See ADR-0024.
- The glossary gains Run Worktree, Run Branch, and Integration Pending.
- Untracked provisioning is an explicit config list, default empty — magic
  copying was rejected.
- The Artifact Directory builtin default leaves the repo tree; explicit
  configuration keeps working unchanged.

## Open Questions

None.
