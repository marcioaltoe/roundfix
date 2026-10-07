---
status: done
created_at: 2026-07-17
updated_at: 2026-09-08
absorbed_by: 0059-run-storage-compaction-and-global-sanitation
---

# Run storage — sanitation is repository-scoped and does not compact SQLite (2026-07-17)

An inspection of Roundfix Home followed the cleanup of terminal Run Worktrees and evaluated whether Roundfix can reclaim its remaining Run storage without direct filesystem or SQLite operations. The existing GC Command provides automated retention cleanup, but no active Spec covers global Artifact Directory sanitation or physical Run Database compaction.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-07-17-global-run-storage-sanitation-and-compaction.md`.
