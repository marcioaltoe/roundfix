---
schema: roundfix/archive-record/v1
spec: 0039-review-source-evidence-and-detached-outcomes
title: Review Source evidence and Detached Run outcomes
status: archived
created: "2026-07-17"
archived: "2026-07-28"
disposition: qa-override
source: docs/history/specs/0039-review-source-evidence-and-detached-outcomes
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-28-04.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: ""
qa_override_reason: ""
qa_override_qa_outcome: ""
qa_override_qa_task_status: ""
qa_override_revision: ""
adrs:
  - ADR-0054
sources: []
regeneration: []
promoted: []
pull_request: "41"
delivery_commit: 397227ff8a1ecbc1f1fa5d180bea282a33429bee
---

# Review Source evidence and Detached Run outcomes

Roundfix currently treats several different CodeRabbit states as the same signal: a completed check can mean a real review or an explicit skip, an approval can be ignored when a check is missing, and one transient GitHub failure can end a watch before any Round. When failure happens in a Detached Run, zero issue counts and a context-free notification leave users and Supervisors without the evidence needed to recover. Prior dogfood evidence was absorbed into this Spec and remains in Git history; the still-open behavior is documented by the [Vortex detached-watch finding](../../findings/2026-07-16-vortex-pr87-detached-watch-notification.md).
