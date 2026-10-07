---
status: done
created_at: 2026-07-28
updated_at: 2026-09-08
absorbed_by: 0053-qa-gate-reachability-and-verdict-semantics
---

# Failed QA Runs strand Run Branches that block the next review Run, and nothing can classify them as superseded (2026-07-28)

Every QA gate attempt creates a Run Branch and commits its QA Report there. When the gate fails, the outcome is Unresolved, so nothing integrates and the branch is kept. After a few attempts the repository accumulates Run Branches that each hold one commit, and Branch Integrity Preflight then refuses to start a review Run until a human resolves every one of them.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-07-28-failed-qa-runs-strand-branches-that-block-review-runs.md`.
