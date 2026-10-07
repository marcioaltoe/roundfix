---
schema: roundfix/archive-record/v1
spec: 0081-a-journal-cheap-to-write-and-keep
title: A journal cheap to write and cheap to keep
status: archived
created: "2026-08-06"
archived: "2026-08-11"
disposition: pass
source: docs/history/specs/0081-a-journal-cheap-to-write-and-keep
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: task_09
qa_report: qa-report-2026-08-11-02.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0008
sources:
  - 2026-08-06-event-journal-payload-economics.md
regeneration: []
promoted: []
pull_request: "156"
delivery_commit: 4b9f0d35eff3522f5ccb58de7cbf9875798fa076
---

# A journal cheap to write and cheap to keep

The Run Event Journal is 2.9 GB of a 3.1 GB Run Database, and every byte of it is inside the fourteen-day retention window: 1,645,457 events, one Run holding 42,000, roughly 1.8 KB of payload each. Retention is working exactly as configured — `roundfix gc --dry-run` finds 302 eligible Runs and zero prunable rows, and no event predates the boundary.
