---
schema: roundfix/archive-record/v1
spec: 0145-a-runner-that-owns-its-lock
title: A runner that owns its lock
status: archived
created: "2026-09-18"
archived: "2026-09-18"
disposition: pass
source: docs/history/specs/0145-a-runner-that-owns-its-lock
source_revision: 31c30719398d5c2e43af79d5868b1268fc495f10
qa_task: task_02
qa_report: qa-report-2026-09-18.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "215"
delivery_commit: bcfa076280df3b85d0272d018270e2b641a9df11
---

# A runner that owns its lock

The acpx runner keeps the state of every Agent Session it has ensured, warmed, started work on, or assigned a selection to, in maps guarded by a mutex it holds as a field. Sixteen of its methods take that struct by value.
