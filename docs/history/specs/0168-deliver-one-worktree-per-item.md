---
schema: roundfix/archive-record/v1
spec: 0168-deliver-one-worktree-per-item
title: Deliver, one worktree per item
status: archived
created: "2026-09-24"
archived: "2026-09-25"
disposition: pass
source: docs/history/specs/0168-deliver-one-worktree-per-item
source_revision: 31c30719398d5c2e43af79d5868b1268fc495f10
qa_task: task_05
qa_report: qa-report-2026-09-25-03.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "251"
delivery_commit: a6866da021ed71d62cc9d31158156c7e1c75ac0a
---

# Deliver, one worktree per item

Spec 0161 shipped with four recorded limits, all in park, and stated that no release may ship `roundfix deliver` until this Spec merges. This Spec inherits that gate: the command stays unreleased until it merges.
