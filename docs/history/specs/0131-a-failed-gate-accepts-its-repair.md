---
schema: roundfix/archive-record/v1
spec: 0131-a-failed-gate-accepts-its-repair
title: A failed gate accepts its repair
status: archived
created: "2026-09-10"
archived: "2026-09-10"
disposition: pass
source: docs/history/specs/0131-a-failed-gate-accepts-its-repair
source_revision: 1d1726f5d227731aa229fbc336dccbae9aca1398
qa_task: task_03
qa_report: qa-report-2026-09-10.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "183"
delivery_commit: 113766686eb04daf9f9fb2ada0834b1ec279c2f4
---

# A failed gate accepts its repair

When a Spec's terminal QA gate fails, the documented recovery is to author corrective Tasks and let the gate run again. Today that recovery cannot be executed: adding a Task under a gate that has already settled makes the Task Graph refuse to load, so the Supervisor can neither check the Spec nor dispatch the corrective work. The gate that found the defect is the reason the defect cannot be fixed.
