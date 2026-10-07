---
status: done
absorbed_by: 0103-a-suite-that-leaks-nothing
created_at: 2026-08-06
updated_at: 2026-09-08
---

# The detach tests leak the process they prove survives

**Date:** 2026-08-06 **Found by:** `ps` on a developer machine, while looking for stale Runs before starting unrelated work. Nothing in Roundfix reported these.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-08-06-the-detach-tests-leak-the-process-they-prove-survives.md`.
