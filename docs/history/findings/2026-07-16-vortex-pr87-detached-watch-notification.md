---
status: done
created_at: 2026-07-16
updated_at: 2026-09-08
absorbed_by: 0039-review-source-evidence-and-detached-outcomes
---

# Detached watch — terminal failure notification lacked actionable context (2026-07-16)

A detached watch for Vortex PR #87 failed after a GitHub API timeout during a brief network interruption. The user later asked whether the review was still running; only then did the Supervisor inspect the Run Database and Console Log, discover the terminal failure, and start a replacement Run. This report extends the shipped [Run Outcome Notifications spec](../specs/_archived/0019-run-outcome-notifications/_prd.md) with evidence from an unattended review Run.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-07-16-vortex-pr87-detached-watch-notification.md`.
