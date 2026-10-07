---
schema: roundfix/archive-record/v1
spec: 0224-an-archived-retry-that-needs-no-recorded-candidate
title: An archived retry that needs no recorded candidate
status: archived
created: "2026-10-04"
archived: "2026-10-04"
disposition: qa-override
source: docs/history/specs/0224-an-archived-retry-that-needs-no-recorded-candidate
source_revision: a30dc847a037fe584812b4471aa7da1701a9ec59
qa_task: task_04
qa_report: qa-report-2026-10-04-01.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: 'maintainer standing authorization 2026-09-30 (unattended program): qa_override only for environment-only partials'
qa_override_reason: 'Environment-only partial after the corrective task_05: R7 needs the original 0220 queue item, which no longer exists (the intervention log entries 140-142 record it), and R13 has no Pull Request yet; the queue opens it. Every behavior row passed.'
qa_override_qa_outcome: partial
qa_override_qa_task_status: failed
qa_override_revision: fe4baadb3398ce95040a59ad297442bdb4a33178
adrs:
  - ADR-0229
sources: []
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "381"
delivery_commit: f0353479ee550b1fc846b3649542e56c3c83538b
---

# An archived retry that needs no recorded candidate

On 2026-10-04 the Delivery Queue parked Spec 0220 `run-unresolved`. Its newest QA Report was `partial` with no finding-blocked row, and its two environment-blocked rows both waited for an open Pull Request: one was the Pull Request row, the other the Linux CI run that only a Pull Request starts. The operator archived the Spec with the QA Archive Override, as the standing authorization allows, and `roundfix deliver retry` refused with `candidate head is missing`. No review had run, so no candidate head was recorded, and only a `qa-environment-partial` park may use the Implement start head of its Run in place of a candidate. The operator opened Pull Request #367 and ran the pre-PR review by hand. Spec 0211 ended the same way on 2026-10-02, before that exception existed.
