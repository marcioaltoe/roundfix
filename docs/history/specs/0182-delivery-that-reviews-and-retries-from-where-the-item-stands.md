---
schema: roundfix/archive-record/v1
spec: 0182-delivery-that-reviews-and-retries-from-where-the-item-stands
title: Delivery that reviews and retries from where the item stands
status: archived
created: "2026-09-29"
archived: "2026-09-29"
disposition: pass
source: docs/history/specs/0182-delivery-that-reviews-and-retries-from-where-the-item-stands
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_04
qa_report: qa-report-2026-09-29.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0169
  - ADR-0170
sources:
  - 2026-09-28-delivery-review-diffs-against-a-moved-main.md
  - 2026-09-28-deliver-retry-carries-forward-from-a-stale-run.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "278"
delivery_commit: a9a6b878356238fd71069074d55e680bc5cadda6
---

# Delivery that reviews and retries from where the item stands

`roundfix deliver` reviews a candidate and recovers a parked item from facts it reads at that moment. On 2026-09-28 two of those facts were read from the wrong place, and each cost a manual intervention:
