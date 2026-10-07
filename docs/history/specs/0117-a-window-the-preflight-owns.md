---
schema: roundfix/archive-record/v1
spec: 0117-a-window-the-preflight-owns
title: A window the Preflight owns
status: archived
created: "2026-08-26"
archived: "2026-08-26"
disposition: pass
source: docs/history/specs/0117-a-window-the-preflight-owns
source_revision: 1d1726f5d227731aa229fbc336dccbae9aca1398
qa_task: task_06
qa_report: qa-report-2026-08-26.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0004
sources: []
regeneration: []
promoted: []
pull_request: "168"
delivery_commit: aba7a02060412303282adf78e85ae9a51bbd26ec
---

# A window the Preflight owns

A Supervisor running an unattended session bounds it by wall-clock time: work until the cutoff, then stop opening new work and let what is running finish. Today that bound lives in the Supervisor's own discipline — a script it must remember to consult before each Run. A bound the caller must remember to check is not a bound: skip the check once and Runs keep opening for the rest of the night, with the guard installed and working. The Preflight already refuses Run creation for a dozen reasons and is the one place a Run cannot get past. The cutoff belongs there.
