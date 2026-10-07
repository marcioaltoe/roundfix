---
schema: roundfix/archive-record/v1
spec: 0173-a-delivery-queue-that-recovers
title: A delivery queue that recovers
status: archived
created: "2026-09-28"
archived: "2026-09-28"
disposition: pass
source: docs/history/specs/0173-a-delivery-queue-that-recovers
source_revision: 31c30719398d5c2e43af79d5868b1268fc495f10
qa_task: task_07
qa_report: qa-report-2026-09-28.md
qa_verdict: pass
unproven: []
adrs: []
sources:
  - 2026-09-25-a-parked-delivery-item-cannot-be-resumed.md
  - 2026-09-25-carry-forward-refuses-inputs-moved-by-a-sibling-task.md
  - 2026-09-25-project-config-dropped-silently-from-task-commits.md
  - 2026-09-25-fallback-proof-timeout-blocks-dispatch.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "263"
delivery_commit: 6fac37ea45c845761e851e78b4713af7a71d272c
---

# A delivery queue that recovers

`roundfix deliver` exists so an operator does not orchestrate a Spec from Run to merge by hand. In the first efficiency wave all three queued Specs parked `run-unresolved`, and the queue had no way back: the owner loop skips every parked item, `deliver resume` only restarts the owner, and `deliver start` records a new queue. Each Spec was finished by hand: fast-forward the item branch to the Run Branch, add a corrective Task, implement, review with an explicit base, archive, gate, open the pull request and merge. Four defects make that recovery necessary or make dispatch fail for reasons the operator cannot see:
