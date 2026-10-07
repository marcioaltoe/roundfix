---
schema: roundfix/archive-record/v1
spec: 0097-a-wave-that-cannot-collide
title: A wave that cannot collide
status: archived
created: "2026-08-12"
archived: "2026-09-02"
disposition: pass
source: docs/history/specs/0097-a-wave-that-cannot-collide
source_revision: 1d1726f5d227731aa229fbc336dccbae9aca1398
qa_task: task_06
qa_report: qa-report-2026-09-02.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "173"
delivery_commit: 8235a8506d9a34efe9a0086538f1e9ada4d3b86c
---

# A wave that cannot collide

Raising Task Worktree concurrency made three latent failures visible at once, and all three cost whole Runs of finished Agent work. Sibling worktrees bootstrap against one shared Git directory and collide on its lock, reporting failure after having done every byte of the work. Tasks the graph declares independent edit the same file and die at integration, which never appeared while Tasks ran one at a time. And creating a second Task Worktree under suite load fails with a raw filesystem errno that names neither the Run, the Task, nor the concurrency that produced it. Concurrency is configured as a number today and verified as nothing; this Spec makes the safety a property the graph can prove before it dispatches.
