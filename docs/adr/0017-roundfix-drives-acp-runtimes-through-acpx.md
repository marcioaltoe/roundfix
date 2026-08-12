---
status: accepted
created_at: 2026-07-05T22:17:04Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# Roundfix drives ACP Runtimes through acpx

Roundfix replaces its hand-rolled ACP client layer (`coder/acp-go-sdk`) with process orchestration of the acpx CLI as the only agent layer — a hard cutover covering review and spec Runs alike, with the acpx version pinned and verified before Runs start. This accepts Node.js as a required dependency of the agent layer (already a de-facto requirement, since the default runtimes launch via npx) and acpx's declared alpha churn; the pinned version and the `agent.Runner` seam are the containment, so rollback is a git revert, not a runtime toggle. A dual SDK/acpx backend was rejected because it doubles the test matrix of the most fragile subsystem indefinitely.
