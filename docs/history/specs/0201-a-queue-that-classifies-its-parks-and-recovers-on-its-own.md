---
schema: roundfix/archive-record/v1
spec: 0201-a-queue-that-classifies-its-parks-and-recovers-on-its-own
title: A queue that classifies its parks and recovers on its own
status: archived
created: "2026-09-30"
archived: "2026-10-01"
disposition: pass
source: docs/history/specs/0201-a-queue-that-classifies-its-parks-and-recovers-on-its-own
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_05
qa_report: qa-report-2026-10-01.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0192
  - ADR-0193
sources: []
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "311"
delivery_commit: 7a2ed622b7e4a11343e7b1c52ae8b023a3b473d7
---

# A queue that classifies its parks and recovers on its own

The v0.22.0 Delivery Queue ran from 11:41 to 18:26 on 2026-09-30 and needed twelve manual interventions. Only Spec 0199 went from start to merge through the queue, and it still needed two manual steps. The finding [2026-09-30-the-v0-22-0-queue-needed-twelve-manual-interventions.md](../../history/findings/2026-09-30-the-v0-22-0-queue-needed-twelve-manual-interventions.md) groups them by class. Three classes are the queue's own:
