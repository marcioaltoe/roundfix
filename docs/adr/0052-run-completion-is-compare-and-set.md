---
status: accepted
created_at: 2026-07-17T16:00:35Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# Run completion is compare-and-set

Ordinary Run completion may move only a non-terminal Run to a terminal outcome: an identical replay is idempotent, while a competing terminal outcome is rejected without changing the Run, lock, event history, or notification. Force stop completes Stopped only after the recorded owner process is proven absent; the sole terminal-state exception is an explicit, evidence-backed Integration Pending reconciliation to Clean, recorded with its prior outcome.
