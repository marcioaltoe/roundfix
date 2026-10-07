---
schema: roundfix/archive-record/v1
spec: 0187-a-queue-that-recovers-without-a-supervisor
title: A queue that recovers without a supervisor
status: archived
created: "2026-09-30"
archived: "2026-09-30"
disposition: pass
source: docs/history/specs/0187-a-queue-that-recovers-without-a-supervisor
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_05
qa_report: qa-report-2026-09-30.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0178
sources:
  - 2026-09-30-post-merge-cleanup-needs-the-merge-commit-locally.md
  - 2026-09-30-retry-refusal-does-not-say-carry-forward-first.md
  - 2026-09-30-a-queue-owner-older-than-its-starting-main-runs-silently.md
  - 2026-09-30-a-grant-widened-mid-delivery-needs-a-rebase.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "293"
delivery_commit: 853ddd103fff934d860057368072e3696f9fd7e4
---

# A queue that recovers without a supervisor

The first two full Delivery Queue trials, on 2026-09-29, still needed a Supervisor to finish their items. Four of the remaining interventions came from Roundfix reading a stale or partial fact:
