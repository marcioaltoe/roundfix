---
schema: roundfix/archive-record/v1
spec: 0202-a-qa-gate-that-reruns-only-stale-rows
title: A QA gate that reruns only stale rows
status: archived
created: "2026-09-30"
archived: "2026-10-01"
disposition: pass
source: docs/history/specs/0202-a-qa-gate-that-reruns-only-stale-rows
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_05
qa_report: qa-report-2026-10-01.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0194
  - ADR-0195
sources: []
regeneration:
  - command: make skills-sync
  - command: make baseline-digests
promoted: []
pull_request: "314"
delivery_commit: 22631b4bb74066b37c0388c9f483d4f95aa6974a
---

# A QA gate that reruns only stale rows

After a corrective Task, the authored QA gate runs again, and its Agent Session executes every row of the matrix. Most of those rows read inputs the correction never touched. ADR-0097 already allows a passing row to carry forward when its declared repository inputs are unmoved. In practice no row ever carries, for two reasons.
