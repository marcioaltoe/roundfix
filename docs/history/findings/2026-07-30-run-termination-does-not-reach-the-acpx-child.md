---
status: done
created_at: 2026-07-30
updated_at: 2026-09-08
absorbed_by: 0066-run-teardown-reclaims-what-it-created
---

# Run termination does not reach the acpx child (2026-07-30)

Checking whether any Run was still active turned up four `acpx` processes that had been spinning for **three days and six hours**, since 2026-07-27 11:26. They belonged to Spec 0037's live QA fixture. Their parent was gone, the worktrees they pointed at were gone, and nothing had ever told them to stop.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-07-30-run-termination-does-not-reach-the-acpx-child.md`.
