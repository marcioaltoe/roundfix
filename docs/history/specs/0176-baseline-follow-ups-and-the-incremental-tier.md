---
schema: roundfix/archive-record/v1
spec: 0176-baseline-follow-ups-and-the-incremental-tier
title: Baseline follow-ups and the incremental tier
status: archived
created: "2026-09-28"
archived: "2026-09-28"
disposition: pass
source: docs/history/specs/0176-baseline-follow-ups-and-the-incremental-tier
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_05
qa_report: qa-report-2026-09-28.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0118
sources:
  - 2026-09-25-the-incremental-verification-waiver-repeats-in-every-spec.md
  - 2026-09-25-docs-layout-still-describes-hand-stamped-overrides.md
  - 2026-09-25-skills-lock-preimage-taken-after-the-fetch.md
  - 2026-09-25-reconcile-required-removed-breaks-the-output-contract.md
regeneration:
  - command: make baseline-digests
  - command: make skills-sync
promoted: []
pull_request: "262"
delivery_commit: f45efbeef4ae091330f17276f915c60a1fac44ca
---

# Baseline follow-ups and the incremental tier

The Context-Driven Baseline still has four open gaps, left behind by Specs 0121, 0163 and 0169:
