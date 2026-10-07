---
schema: roundfix/archive-record/v1
spec: 0228-queue-items-that-stay-current-with-main
title: Queue items that stay current with main
status: archived
created: "2026-10-04"
archived: "2026-10-05"
disposition: qa-override
source: docs/history/specs/0228-queue-items-that-stay-current-with-main
source_revision: a30dc847a037fe584812b4471aa7da1701a9ec59
qa_task: task_05
qa_report: qa-report-2026-10-05-01.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: 'Environment-only partial after corrective task_06: Q08 needs the historical 0225 queue item, which no longer exists (intervention log entries 157-159 record it), and Q14 has no Pull Request yet; the queue opens it. Every behavior row passed, including the real SKILL.md version collision.'
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: 3ef4d88faf1ed9ecc6edcddc61206707cfd3e85e
adrs:
  - ADR-0233
sources:
  - 2026-10-04-queued-specs-raise-the-same-skill-version.md
  - 2026-10-04-a-docs-fix-after-archive-forces-a-manual-pull-request.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "391"
delivery_commit: 00382172007d15ef7de178659bc84a4e92c156a4
---

# Queue items that stay current with main

This is a bug fix with two causes, both recorded on 2026-10-04 in the operator's intervention log and adopted from two Backlog Entries: [queued Specs raise the same skill version](references/2026-10-04-queued-specs-raise-the-same-skill-version.md) and [a docs fix after archive forces a manual Pull Request](references/2026-10-04-a-docs-fix-after-archive-forces-a-manual-pull-request.md).
