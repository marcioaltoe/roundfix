---
status: done
created_at: 2026-07-29
updated_at: 2026-09-08
absorbed_by: 0124-verification-capacity-and-measured-economics
---

# QA cycles — the cost is a cold Run Worktree and the Agent's turn count, not the gate itself (2026-07-29)

A QA cycle on Spec 0061 took roughly twenty minutes of wall-clock. Measuring where that time goes contradicted the obvious hypothesis: the repository gate accounts for about an eighth of it, and compilation for a fraction of that. This report separates what is specific to this repository from what every repository using Roundfix pays, because the two need different fixes.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-07-29-qa-cycle-cost-is-cold-environments-and-agent-turns.md`.
