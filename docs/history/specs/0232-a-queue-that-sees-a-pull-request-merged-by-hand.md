---
schema: roundfix/archive-record/v1
spec: 0232-a-queue-that-sees-a-pull-request-merged-by-hand
title: A queue that sees a Pull Request merged by hand
status: archived
created: "2026-10-05"
archived: "2026-10-05"
disposition: qa-override
source: docs/history/specs/0232-a-queue-that-sees-a-pull-request-merged-by-hand
source_revision: a30dc847a037fe584812b4471aa7da1701a9ec59
qa_task: task_03
qa_report: qa-report-2026-10-05.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: 'Environment-only partial: the Run sandbox could not reach GitHub or the web (rows 7a-7c, 7f) and the original 0231 queue row is gone (7e); the operator confirmed with gh on 2026-10-05 that #404, #382 and #396 are MERGED from their item branches (merged 2026-10-05T20:20:08Z, 2026-10-04T21:03:19Z, 2026-10-05T15:52:46Z); row 10 has no Pull Request yet. Every behavior row passed.'
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 477a1ebe287fcd8887123fdd4542a36552d40200
adrs:
  - ADR-0237
  - ADR-0232
sources:
  - 2026-10-05-a-pull-request-merged-by-hand-stays-parked.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "407"
delivery_commit: f2bc1e05d5b0d1bc2da1d4b7aca118dd144feafa
---

# A queue that sees a Pull Request merged by hand

This is a bug fix adopted from one Backlog Entry, [a Pull Request merged by hand stays parked](references/2026-10-05-a-pull-request-merged-by-hand-stays-parked.md), recorded on 2026-10-05 from the operator's intervention log. When the operator merges a parked item's Pull Request by hand, the Delivery Queue never learns it, and no command answers the item's Pending Question.
