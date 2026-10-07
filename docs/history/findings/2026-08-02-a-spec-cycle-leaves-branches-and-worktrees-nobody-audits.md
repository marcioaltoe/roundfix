---
status: done
created_at: 2026-08-02
updated_at: 2026-09-08
absorbed_by: 0068-spec-close-audit
---

# A Spec cycle leaves branches and worktrees nobody audits (2026-08-02)

The per-Spec loop in `docs/agents/autonomous-work.md` ends at "squash merge and reconcile". Nothing at that boundary audits what the cycle created against what survived it, so debris accumulates silently and completed work stays invisible. Four distinct kinds surfaced in a single session, and the maintainer found all four by running `git branch -l` and `git worktree list` by hand.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-08-02-a-spec-cycle-leaves-branches-and-worktrees-nobody-audits.md`.
