---
schema: roundfix/archive-record/v1
spec: 0175-cleanup-after-a-squash-merge
title: Cleanup after a squash merge
status: archived
created: "2026-09-28"
archived: "2026-09-28"
disposition: pass
source: docs/history/specs/0175-cleanup-after-a-squash-merge
source_revision: 31c30719398d5c2e43af79d5868b1268fc495f10
qa_task: task_07
qa_report: qa-report-2026-09-28-01.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0115
sources:
  - 2026-09-25-squash-merge-leaves-runs-and-branches-behind.md
  - 2026-09-25-post-merge-cleanup-of-runs-and-branches.md
  - 2026-09-25-reconcile-proves-archived-deleted-target-runs.md
  - 2026-09-25-reconcile-apply-ignores-superseded-runs.md
  - 2026-09-25-legacy-run-repository-key-from-its-worktree.md
  - 2026-09-25-carry-forward-staging-worktree-leaks.md
  - 2026-09-25-half-removed-item-cleanup-ignores-nested-worktrees.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "268"
delivery_commit: 06afa8359886d319b648556021d42622fda551ee
---

# Cleanup after a squash merge

When a Spec's pull request is squash-merged, the Run Worktrees, Run Branches and staging worktrees its Runs left behind stay in the repository. The maintainer asked for every finished squash merge to be followed by a cleanup of completed Runs and branches. On 2026-09-25 this repository held 16 `roundfix/run-*` branches, 16 Run Worktrees and a locked carry-forward staging worktree, most of them for Specs already merged and archived. Roundfix released almost none of them. Six defects explain why:
