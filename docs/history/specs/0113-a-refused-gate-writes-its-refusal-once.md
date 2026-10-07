---
schema: roundfix/archive-record/v1
spec: 0113-a-refused-gate-writes-its-refusal-once
title: A refused gate writes its refusal once
status: archived
created: "2026-08-25"
archived: "2026-08-26"
disposition: pass
source: docs/history/specs/0113-a-refused-gate-writes-its-refusal-once
source_revision: 1d1726f5d227731aa229fbc336dccbae9aca1398
qa_task: task_06
qa_report: qa-report-2026-08-26-01.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "167"
delivery_commit: 86c3ca440f661d9223727513c614df10e33af015
---

# A refused gate writes its refusal once

When the QA gate refuses at a precondition, it writes a report whose Results table is empty — correctly, because it stopped before building the matrix. The mechanical stage of every later run then reads that report and refuses:
