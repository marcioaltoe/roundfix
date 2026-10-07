---
status: done
created_at: 2026-07-28
updated_at: 2026-09-08
absorbed_by: 0054-tooling-task-and-verification-hygiene
---

# Protected-tooling Tasks demand a green repository and an undocumented commit choreography (2026-07-28)

A Task that edits protected tooling has to satisfy two conditions nothing in the repository states, and failing either costs a Run plus a QA cycle. Both were hit repeatedly on 2026-07-27 and 2026-07-28 across Specs 0037, 0038, and 0039 — by Agents behaving correctly and by the supervisor, who got the choreography wrong twice in a row while knowing the rule existed.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-07-28-tooling-tasks-need-a-green-repo-and-an-undocumented-commit-choreography.md`.
