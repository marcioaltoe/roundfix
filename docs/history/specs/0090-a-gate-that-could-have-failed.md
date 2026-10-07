---
schema: roundfix/archive-record/v1
spec: 0090-a-gate-that-could-have-failed
title: A gate that could have failed
status: archived
created: "2026-08-09"
archived: "2026-08-09"
disposition: qa-override
source: docs/history/specs/0090-a-gate-that-could-have-failed
source_revision: 1d1726f5d227731aa229fbc336dccbae9aca1398
qa_task: task_07
qa_report: qa-report-2026-08-09-01.md
qa_verdict: fail
unproven: []
qa_override: true
qa_override_approval: ""
qa_override_reason: ""
qa_override_qa_outcome: ""
qa_override_qa_task_status: ""
qa_override_revision: ""
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "146"
delivery_commit: d5b96ca96bd520a580a0622ad91c1902fb8e546b
---

# A gate that could have failed

A Task's `## Verification` commands are the only thing standing between an Agent turn and a Task marked `completed`. Spec 0089 proved that this barrier can be zero-height without anyone noticing.
