---
schema: roundfix/archive-record/v1
spec: 0163-baseline-decisions-and-regeneration
title: Baseline decisions and regeneration
status: archived
created: "2026-09-24"
archived: "2026-09-25"
disposition: pass
source: docs/history/specs/0163-baseline-decisions-and-regeneration
source_revision: 31c30719398d5c2e43af79d5868b1268fc495f10
qa_task: task_07
qa_report: qa-report-2026-09-25.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration:
  - command: make baseline-digests
  - command: make skills-sync
promoted: []
pull_request: "252"
delivery_commit: 7dfda64d20e35bf6189502c6fdc67d04002f7140
---

# Baseline decisions and regeneration

Four defects in the Baseline and its tooling remain from portfolio Spec 0121.
