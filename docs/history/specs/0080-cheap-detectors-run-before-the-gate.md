---
schema: roundfix/archive-record/v1
spec: 0080-cheap-detectors-run-before-the-gate
title: Cheap detectors run before the gate
status: archived
created: "2026-08-06"
archived: "2026-08-11"
disposition: pass
source: docs/history/specs/0080-cheap-detectors-run-before-the-gate
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: task_08
qa_report: qa-report-2026-08-11-07.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0091
  - ADR-0088
sources:
  - 2026-08-03-verification-performance-contract.md
  - 2026-08-06-two-stage-qa-gate-economics.md
regeneration: []
promoted: []
pull_request: "155"
delivery_commit: a2a4c86b7570e4ef782ccc2ff390033f586d47fd
---

# Cheap detectors run before the gate

The QA gate is the slowest step in the loop and the one every fleet project pays for, because Roundfix is the shared harness. It is a single agent audit that rebuilds its whole criterion matrix from scratch every round, so a one-line correction costs a full audit — measured on Spec 0079 as 92 minutes for the first Run and then 29 and 30 minutes for rounds whose only work was moving one declared constant. Inside those rounds the authoritative `make verify` runs cold and complete at about 90 seconds, when the same gate on an unchanged tree would cost 4.9 seconds.
