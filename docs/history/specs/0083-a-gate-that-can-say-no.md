---
schema: roundfix/archive-record/v1
spec: 0083-a-gate-that-can-say-no
title: A gate that can say no
status: archived
created: "2026-08-07"
archived: "2026-08-09"
disposition: pass
source: docs/history/specs/0083-a-gate-that-can-say-no
source_revision: 1d1726f5d227731aa229fbc336dccbae9aca1398
qa_task: task_07
qa_report: qa-report-2026-08-07.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
---

# A gate that can say no

The repository's own rule says the local gate is the only gate. On 2026-08-07 that gate was observed exiting `0` on a working tree whose Go suite exits `1`, with the failing package omitted from its summary. A gate that cannot say no is not a gate, and every completion claim routed through it is unverified.
