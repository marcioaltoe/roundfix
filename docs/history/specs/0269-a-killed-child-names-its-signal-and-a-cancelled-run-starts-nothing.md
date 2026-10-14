---
schema: roundfix/archive-record/v1
spec: 0269-a-killed-child-names-its-signal-and-a-cancelled-run-starts-nothing
title: A killed child names its signal, and a cancelled Run starts nothing
status: archived
created: "2026-10-13"
archived: "2026-10-14"
disposition: partial
source: docs/specs/0269-a-killed-child-names-its-signal-and-a-cancelled-run-starts-nothing
source_revision: 4595dea9ce8874f54dbc893c6e538b7ecfcd5c17
qa_task: task_05
qa_report: qa-report-2026-10-14.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0056
  - ADR-0266
sources:
  - 2026-10-13-a-test-child-is-killed-by-a-signal-from-an-unknown-sender.md
  - 2026-10-13-the-verification-capacity-cancellation-test-starts-a-second-command.md
regeneration:
  - command: make skills-sync
promoted: []
---

# A killed child names its signal, and a cancelled Run starts nothing

The Delivery Queue and all nine child-Roundfix test-helper sites name a Signal Exit, preserving ordinary-exit wording. Verification Capacity rejects cancellation after an immediate grant, releases the permit and starts no new Verification or pre-work probe; the controlled cancellation negative test exposes the former second start in 20/20 runs. Documentation, the glossary and the open sender-tracing Backlog are consistent. The Archive Record can carry a qualifying partial QA result with only the pre-PR Pull Request row unproven; this does not claim Merge-Ready acceptance or delivery.
