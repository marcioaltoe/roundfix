---
schema: roundfix/archive-record/v1
spec: 0265-tests-that-wait-on-a-child-do-not-race-the-machine
title: Tests that wait on a child do not race the machine
status: archived
created: "2026-10-13"
archived: "2026-10-13"
disposition: partial
source: docs/specs/0265-tests-that-wait-on-a-child-do-not-race-the-machine
source_revision: 68f079ce82e6e221e3dc2b076da8caadc8f13f69
qa_task: task_05
qa_report: qa-report-2026-10-13.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0271
sources: []
regeneration: []
promoted: []
pull_request: "509"
---

# Tests that wait on a child do not race the machine

All 31 inventory tests passed 50 consecutive race-detector repetitions under niced CPU load, including the real local ACPX round trip without a skip; three consecutive full suites passed under their authored pre-0259 exclusions with no killed or negative-exit diagnostic. The new rule rejected a 500 ms timeout in a disposable repository copy with its exact file and line, the two production dependency defaults were proved, and all eight detach assertion bodies and protected regions were preserved. No production budget value changed. The Archive Record should carry this qualifying partial and the absence of pre-PR Merge-Ready/review-artifact evidence. Durable knowledge already lives in ADR-0271, CONTEXT.md and the repository guide; task_05 settlement and subsequent delivery remain the Daemon's responsibility.
