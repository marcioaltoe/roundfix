---
status: accepted
created_at: 2026-07-06T21:05:00Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# Run Worktrees can bootstrap their environment

A Run or Task Worktree is created from committed Git state, so it lacks the
untracked, machine-local environment a project's Verification needs — installed
dependencies (`node_modules`), a migrated and seeded database, warmed build
caches. `worktree.copy` already places untracked files (like `.env`) into the
worktree, but it only copies files; it runs nothing. Roundfix therefore gains a
configured `worktree.bootstrap` command that runs once in each newly created
worktree, after `worktree.copy` and before any Agent work or Verification, to
prepare that environment (for example `bun install && bun run db:migrate && bun
run db:seed`). A bootstrap failure ends the Run or Task with a distinct
bootstrap-failed outcome, because Verification would fail anyway. Roundfix owns
running the command, not what it does: database provisioning and dependency
strategy live in the command, and shared-state projects (one database across
Tasks) run at `worktree.concurrency: 1` so bootstrap runs once on the reused Run
Worktree. This unblocks Roundfix for stateful monorepos (TypeScript + database +
package manager) without abandoning Worktree isolation.
