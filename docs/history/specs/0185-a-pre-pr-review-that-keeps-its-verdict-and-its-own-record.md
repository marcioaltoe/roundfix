---
schema: roundfix/archive-record/v1
spec: 0185-a-pre-pr-review-that-keeps-its-verdict-and-its-own-record
title: A pre-PR review that keeps its verdict and its own record
status: archived
created: "2026-09-29"
archived: "2026-09-29"
disposition: pass
source: docs/history/specs/0185-a-pre-pr-review-that-keeps-its-verdict-and-its-own-record
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_03
qa_report: qa-report-2026-09-29.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0174
sources:
  - 2026-09-29-a-review-verdict-after-progress-text-is-unclassifiable.md
  - 2026-09-29-one-review-record-is-shared-by-every-checkout.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "283"
delivery_commit: ec608ba5070a777c9d9da69bb7d931667ba6df05
---

# A pre-PR review that keeps its verdict and its own record

`roundfix review` (Specs 0153, 0160 and 0179) runs the configured reviewer over the candidate and records the result. On 2026-09-29 it lost a real verdict and let one checkout's review replace another's.
