---
schema: roundfix/archive-record/v1
spec: 0089-an-effort-the-runtime-actually-receives
title: An effort the runtime actually receives
status: archived
created: "2026-08-09"
archived: "2026-08-09"
disposition: qa-override
source: docs/history/specs/0089-an-effort-the-runtime-actually-receives
source_revision: 1d1726f5d227731aa229fbc336dccbae9aca1398
qa_task: task_08
qa_report: qa-report-2026-08-09-02.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: ""
qa_override_reason: ""
qa_override_qa_outcome: ""
qa_override_qa_task_status: ""
qa_override_revision: ""
adrs:
  - ADR-0108
  - ADR-0106
sources:
  - 2026-08-09-the-opencode-runtime-hands-back-the-floor-of-every-range.md
regeneration: []
promoted: []
---

# An effort the runtime actually receives

Spec 0088 made the OpenCode runtime reachable and, to keep Exact Agent Selection Proof token-free, made it model-managed: Roundfix refuses any non-empty `reasoning_effort` there and accepts whatever the model opens at. ADR-0106 recorded that trade honestly. Measured a day later, the trade is worse than it looked.
