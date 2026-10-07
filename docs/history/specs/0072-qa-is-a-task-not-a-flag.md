---
schema: roundfix/archive-record/v1
spec: 0072-qa-is-a-task-not-a-flag
title: QA is a Task, not a flag
status: archived
created: "2026-08-02"
archived: "2026-08-03"
disposition: qa-override
source: docs/history/specs/0072-qa-is-a-task-not-a-flag
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: task_07
qa_report: qa-report-2026-08-03-03.md
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
pull_request: "88"
delivery_commit: 86c8d8c21183039f4173cdc67c08bcaccf0a1305
---

# QA is a Task, not a flag

The QA gate is requested with `--qa` on the Implement Command. The Daemon withholds it correctly — no gate begins while any Task is unsettled, and Spec 0057's first Run proves it by ending with one Task failed and no gate at all. Nothing runs early.
