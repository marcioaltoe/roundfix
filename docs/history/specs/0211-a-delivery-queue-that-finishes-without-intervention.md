---
schema: roundfix/archive-record/v1
spec: 0211-a-delivery-queue-that-finishes-without-intervention
title: A delivery queue that finishes without intervention
status: archived
created: "2026-10-01"
archived: "2026-10-02"
disposition: qa-override
source: docs/history/specs/0211-a-delivery-queue-that-finishes-without-intervention
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_05
qa_report: qa-report-2026-10-02-01.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: QA partial with zero findings after the QA-0211-01 fix; blocked rows need the operator's intervention log outside the repository (07), GitHub API access the sandbox denies (08) and an open Pull Request (13)
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 649bfeedc6de4cfad3740f25ff3eb9159e3084d3
adrs:
  - ADR-0211
sources:
  - 2026-10-01-an-archived-retry-cannot-find-the-runs-implement-start-head.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "343"
delivery_commit: 7c1bc29a654ae36ed88b81753955a3aba3266a23
---

# A delivery queue that finishes without intervention

On 2026-10-01 the Delivery Queue needed an operator four times for items it should have finished on its own, and every time the queue's own answer did not work. Spec 0204 parked `qa-environment-partial`; the operator followed the printed answer, and the Delivery Retry refused because it could not find the Implement start head of a Run the queue owner had started. Spec 0203 parked `checks-failed`; after the operator re-ran the failed job and retried, the queue merged while the re-run was still pending, GitHub refused, and the item parked `delivery-error`. A retry of that park refused even though the Pull Request had turned green, and the operator merged Pull Request #317 by hand.
