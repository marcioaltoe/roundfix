---
schema: roundfix/archive-record/v1
spec: 0162-a-durable-repository-key-per-run
title: A durable repository key per Run
status: archived
created: "2026-09-24"
archived: "2026-09-24"
disposition: pass
source: docs/history/specs/0162-a-durable-repository-key-per-run
source_revision: 31c30719398d5c2e43af79d5868b1268fc495f10
qa_task: task_03
qa_report: qa-report-2026-09-24-02.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "243"
delivery_commit: 2dde97041d592f0ba17b2569bbab4a9e07c7235b
---

# A durable repository key per Run

Spec 0157 gave a repository one identity across its worktrees, but found that identity at query time, through the `.git/worktrees/*/gitdir` entries that `git worktree remove` deletes. Its second pre-PR review recorded four limits that follow from that choice:
