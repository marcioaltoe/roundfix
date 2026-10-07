---
schema: roundfix/archive-record/v1
spec: 0172-a-qa-gate-that-tells-the-truth
title: A QA gate that tells the truth
status: archived
created: "2026-09-25"
archived: "2026-09-25"
disposition: pass
source: docs/history/specs/0172-a-qa-gate-that-tells-the-truth
source_revision: 31c30719398d5c2e43af79d5868b1268fc495f10
qa_task: task_06
qa_report: qa-report-2026-09-25-01.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0056
sources:
  - 2026-09-25-hollow-qa-report-passes.md
  - 2026-09-25-qa-report-sequence-reuses-a-gap.md
  - 2026-09-25-temporary-retry-hides-a-deterministic-failure.md
  - 2026-09-25-events-stream-aborts-on-vacuous-probe.md
  - 2026-09-16-an-unobserved-verification-is-published-without-its-classification.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "259"
delivery_commit: 0160f70a3ce53e198ade50b4c5c0817de7a4f506
---

# A QA gate that tells the truth

A QA gate's verdict and the Verification signal on the Run Event Stream must mean exactly what they say. Five defects make them say more, or less, than was measured:
