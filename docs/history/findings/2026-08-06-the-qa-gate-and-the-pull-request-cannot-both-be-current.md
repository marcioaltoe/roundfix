---
status: done
created_at: 2026-08-06
updated_at: 2026-09-08
absorbed_by: 0127-durable-unattended-spec-workflow
---

# 2026-08-06 — The QA gate and the pull request cannot both be current

An authored QA gate that inspects the pull request can never pass on its first run after any corrective Task, because the artifact it verifies and the artifact the pull request shows are different by construction.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-08-06-the-qa-gate-and-the-pull-request-cannot-both-be-current.md`.
