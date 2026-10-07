---
schema: roundfix/archive-record/v1
spec: 0042-verification-capacity-and-daemon-task-settlement
title: Verification Capacity and Daemon Task Settlement
status: archived
created: "2026-07-18"
archived: "2026-07-28"
disposition: pass
source: docs/history/specs/0042-verification-capacity-and-daemon-task-settlement
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: ""
qa_report: qa-report-2026-07-28-03.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0051
  - ADR-0056
  - ADR-0057
sources: []
regeneration: []
promoted: []
pull_request: "40"
delivery_commit: ed4abec3bcdd8d41aa47ac70a9a3d35012b838a1
---

# Verification Capacity and Daemon Task Settlement

Spec Runs currently use `worktree.concurrency` to limit an entire Task lifecycle, so increasing Task Capacity also lets multiple repository-wide Verification suites run at once. In Vortex this coupled two useful concurrent Agent implementations to two simultaneous integration suites, exhausting local listeners and setup capacity. The same Run also let Agents mark Tasks `failed` from their own gate runs, which bypassed Daemon Verification, Verification Feedback, and repair. This Spec absorbed the incident evidence; the original report remains available through Git history.
