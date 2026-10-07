---
schema: roundfix/archive-record/v1
spec: 0203-a-reviewer-that-validates-its-findings-and-remembers-its-rounds
title: A reviewer that validates its findings and remembers its rounds
status: archived
created: "2026-09-30"
archived: "2026-10-01"
disposition: pass
source: docs/history/specs/0203-a-reviewer-that-validates-its-findings-and-remembers-its-rounds
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_05
qa_report: qa-report-2026-10-01.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0196
  - ADR-0197
sources: []
regeneration:
  - command: make skills-sync
  - command: make baseline-digests
promoted: []
pull_request: "317"
delivery_commit: b62f5aa62491c9a09b99f3e9f33a5325d5577fb3
---

# A reviewer that validates its findings and remembers its rounds

The Pre-PR Review Command runs the configured reviewer over the candidate before a Pull Request exists. With the Delivery Queue, each finding parks the item until an operator acts. Three problems make that costly.
