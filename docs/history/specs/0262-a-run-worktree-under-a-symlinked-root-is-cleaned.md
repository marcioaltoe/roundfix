---
schema: roundfix/archive-record/v1
spec: 0262-a-run-worktree-under-a-symlinked-root-is-cleaned
title: A Run Worktree under a symlinked root is cleaned
status: archived
created: "2026-10-12"
archived: "2026-10-13"
disposition: pass
source: docs/specs/0262-a-run-worktree-under-a-symlinked-root-is-cleaned
source_revision: 107b3255599852dc9a05abeeb37c7fd47624ed2b
qa_task: task_04
qa_report: qa-report-2026-10-13.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0268
sources: []
regeneration: []
promoted: []
pull_request: "517"
---

# A Run Worktree under a symlinked root is cleaned

Since 2026-10-10 the maintainer's `~/.roundfix/worktrees` has been a symbolic link to another volume, and no Run Worktree under it has been cleaned since. Every Delivery Queue merge ends with `warning: cleanup failed: keep Run "<run-id>": recorded Run Worktree path is unsafe`, `roundfix reconcile` printed that reason 550 times on 2026-10-12, and the worktree root held 24 GB that nothing reclaims ([the symlinked root Backlog Entry](references/2026-10-12-a-run-worktree-under-a-symlinked-root-is-never-cleaned.md)). Roundfix records each worktree in the configured spelling, Git registers the resolved path, and reconciliation refuses any recorded path that crosses a symbolic link. The cycle's other worktree source, the timing-dependent admin-lock test, is fixed by Spec 0259's corrective task_06 and closed in history, so this Spec does not own it.
