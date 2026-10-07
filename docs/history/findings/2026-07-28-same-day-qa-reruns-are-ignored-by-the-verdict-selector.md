---
status: done
created_at: 2026-07-28
updated_at: 2026-09-08
absorbed_by: 0053-qa-gate-reachability-and-verdict-semantics
---

# QA verdict — every same-day rerun was ignored, so a passing Spec still reported fail (2026-07-28)

`NewestQAReport` picked the day's **first** QA Report, not its latest. A Spec that failed its gate and then passed on a rerun the same day kept reporting the stale failure, and `roundfix archive` refused it.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-07-28-same-day-qa-reruns-are-ignored-by-the-verdict-selector.md`.
