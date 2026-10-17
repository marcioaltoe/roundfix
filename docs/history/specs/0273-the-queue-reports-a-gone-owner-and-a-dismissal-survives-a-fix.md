---
schema: roundfix/archive-record/v1
spec: 0273-the-queue-reports-a-gone-owner-and-a-dismissal-survives-a-fix
title: The queue reports a gone owner, and a dismissal survives a fix
status: archived
created: "2026-10-16"
archived: "2026-10-16"
disposition: partial
source: docs/specs/0273-the-queue-reports-a-gone-owner-and-a-dismissal-survives-a-fix
source_revision: 61ea55013b6e42a0bf83259a94632d318b7bf733
qa_task: task_05
qa_report: qa-report-2026-10-16-01.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0279
sources:
  - 2026-10-14-a-dismissal-after-a-committed-fix-needs-a-reset.md
  - 2026-10-14-deliver-status-does-not-say-the-queue-owner-is-gone.md
regeneration:
  - command: make skills-sync
promoted: []
---

# The queue reports a gone owner, and a dismissal survives a fix

The current gate closes as a qualifying partial: the owner report, owning-checkout refusals and dismissal-after-fix flows pass through the current binary, and all scope, documentation and focused checks pass. The only unmet row is the authored pre-PR Pull Request row. The Daemon owns QA Task settlement and subsequent delivery. No archive, status, commit, push or Pull Request mutation was performed.
