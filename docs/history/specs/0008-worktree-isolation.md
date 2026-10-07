---
schema: roundfix/archive-record/v1
spec: 0008-worktree-isolation
title: Worktree Isolation
status: archived
created: "2026-07-05"
archived: "2026-07-06"
disposition: pass
source: docs/history/specs/0008-worktree-isolation
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-05.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0023
  - ADR-0024
sources: []
regeneration: []
promoted: []
pull_request: "17"
delivery_commit: 5afc9d6c88bd9090cd56613b3a4e0b6b1b9ff5b6
---

# Worktree Isolation

Every Run today executes inside the user's own checkout. Three chronic findings share that single root: concurrent user work gets swept into task commits, the implement preflight must reject any dirty tree (making failed- Run recovery a chore), and the fragile snapshot-diff dance exists only to tell the Agent's changes apart from everyone else's. This Spec moves execution into an isolated per-Run git worktree on a named Run Branch, and brings the commits back to the user's branch through a porcelain-only integration protocol proven safe by experiment — the user keeps working in their checkout while Runs work in theirs.
