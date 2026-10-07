---
schema: roundfix/archive-record/v1
spec: 0066-run-teardown-reclaims-what-it-created
title: Run teardown reclaims what it created
status: archived
created: "2026-08-01"
archived: "2026-08-05"
disposition: pass
source: docs/history/specs/0066-run-teardown-reclaims-what-it-created
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: task_06
qa_report: qa-report-2026-08-05.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "118"
delivery_commit: edf3d1d843704f185cd6f25f63d0c0fbefde48ca
---

# Run teardown reclaims what it created

A Run that ends without passing leaves two kinds of debris behind, and nothing reclaims either.
