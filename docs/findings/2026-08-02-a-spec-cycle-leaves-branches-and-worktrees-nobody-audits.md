---
status: open
created_at: 2026-08-02
updated_at: 2026-08-02
---

# A Spec cycle leaves branches and worktrees nobody audits (2026-08-02)

The per-Spec loop in `docs/agents/autonomous-work.md` ends at "squash merge and
reconcile". Nothing at that boundary audits what the cycle created against what
survived it, so debris accumulates silently and completed work stays invisible.
Four distinct kinds surfaced in a single session, and the maintainer found all
four by running `git branch -l` and `git worktree list` by hand.

## Symptom / evidence

At one point on 2026-08-02, with three Specs shipped that session:

```text
$ git worktree list
/Users/marcio/dev/roundfix                     [ma/profiles-configure-merge-semantics]
/private/tmp/.../scratchpad/wt-arch            [ma/archive-0058]
/private/tmp/.../scratchpad/wt-queue           [ma/spec-queue-from-findings]
/Users/marcio/.roundfix/worktrees/...144653Z   [roundfix/run-run_20260801T144653Z_...]

$ git branch -l
* ma/profiles-configure-merge-semantics
  ma/spec-queue-from-findings
  ma/archive-0058
  main
```

Four kinds of residue, each with a different cause:

1. **Supervisor scratch worktrees.** `wt-arch` and `wt-queue` were created to
   work off the main checkout while a Run was Active. After the push only the
   branches were needed; the working copies were never removed.
2. **A Run Worktree orphaned by its own merge.** `reconcile` reported
   `classification: unknown`, `action: preserve`, `refusal-reason: target branch
   could not be resolved unambiguously`. The target branch
   `ma/npm-trusted-publishing-and-release-preflight` no longer existed because
   the squash merge used `--delete-branch`. The behavior is correct — preserving
   beats guessing — but **every squash merge with `--delete-branch` orphans any
   retained Run Worktree for that branch**.
3. **A stale remote backup branch.** `roundfix/run-run_20260731T195234Z_...` was
   pushed as a backup before a travel pause. Its purpose ended when the work
   merged, and nothing reclaimed it. It sat behind `main` by 19 lines of
   `release.yml`, including a security hardening the branch predated.
4. **Branches held open by unmerged Pull Requests.** Two PRs were opened and
   left unreviewed while the Supervisor moved to the next task. The visible
   consequence is worse than the branch count: Spec 0058's **archive** and the
   five newly queued Specs existed only on those branches, so `main` still
   showed 0058 as active and showed none of the new Specs. Work reported as
   delivered was not where a reader would look for it.

## Root cause

Two gaps, one mechanical and one procedural.

**Mechanical:** `reconcile` resolves a Run to its target branch **by name**. A
squash merge that deletes the branch destroys the only handle it has, so it
degrades to `unknown` and preserves forever. The content is provably in `main` —
`git diff --name-status origin/main <run-branch>` showed no file unique to the
branch and identical implementation files — but nothing performs that check.

**Procedural:** the loop treats "open the Pull Request" as the last Supervisor
action of a Spec and has no closing audit. There is no point at which the loop
asks: which branches did this Spec create, which still exist, which worktrees
are still checked out, and is anything reported as delivered absent from the
default branch?

## Action / suggestion

Add a closing audit to the per-Spec cycle, run at merge-and-sync:

- Enumerate branches and worktrees created by the Spec and report which survive
  and why. A branch backing an open Pull Request is a legitimate survivor; a
  branch whose Pull Request merged is not.
- Resolve a Run's integration **by content** when its target branch no longer
  exists, so `--delete-branch` stops orphaning Run Worktrees. Absence of any
  unique file plus identical implementation files is the proof already used by
  hand.
- Reclaim Supervisor scratch worktrees at push time; the branch is the artifact,
  the working copy is not.
- Report anything the Spec claims to have delivered that is not present on the
  synced default branch. This is the check that would have caught an archive and
  five queued Specs living only on unmerged branches.

The audit must preserve on ambiguity, matching the posture `reconcile` and GC
already take. Its value is naming residue, not deleting aggressively.

## What worked — keep

- `reconcile` refused to guess. `classification: unknown` with an explicit
  `refusal-reason` is what made the orphaned worktree diagnosable in one read.
- The manual proof was cheap and decisive: `git diff --name-status` against
  `main` settled in one command whether the branch carried anything unique.

## Routing — 2026-08-02

Routed to [Spec 0068](../specs/0068-spec-close-audit/_prd.md) on 2026-08-02.
