---
schema: roundfix/archive-record/v1
spec: 0016-worktree-bootstrap
title: Worktree Bootstrap
status: archived
created: "2026-07-06"
archived: "2026-07-06"
disposition: pass
source: docs/history/specs/0016-worktree-bootstrap
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-06.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0034
sources: []
regeneration: []
promoted: []
---

# Worktree Bootstrap

Roundfix executes each agent Run in an isolated worktree created from committed Git state. That worktree has no installed dependencies, no migrated or seeded database, and no warmed caches — so for a stateful project (a TypeScript monorepo with a package manager, a database, and a build cache) the Daemon's Verification fails not because the work is wrong but because the environment was never prepared. `worktree.copy` already places untracked files like `.env` into the worktree, but it only copies files. This Spec adds a configured `worktree.bootstrap` command that Roundfix runs once in each new worktree — after copying files, before any Agent work and before Verification — so the worktree is a working environment. Together with `worktree.copy` and sequential execution, this unblocks Roundfix for stateful monorepos without giving up worktree isolation.
