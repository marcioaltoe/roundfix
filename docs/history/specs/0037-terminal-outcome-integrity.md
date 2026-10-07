---
schema: roundfix/archive-record/v1
spec: 0037-terminal-outcome-integrity
title: Terminal outcome integrity
status: archived
created: "2026-07-17"
archived: "2026-07-27"
disposition: pass
source: docs/history/specs/0037-terminal-outcome-integrity
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-27.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0052
sources: []
regeneration: []
promoted: []
pull_request: "38"
delivery_commit: 8ec92ad968bb8113128d5ca9b17465fb4a1144a3
---

# Terminal outcome integrity

A force-stopped Run can currently be completed again by its still-running owner, a Stop Request can remain unnoticed throughout a Review Source wait, and cleanup can target an Agent Session that never reached the active lifecycle. The resulting state is unsafe for users and Supervisors: the Run Database can contradict the Stop Command, the released lock can coexist with live work, and secondary cleanup noise can obscure the primary failure. Prior dogfood evidence was absorbed into this Spec and remains in Git history; the still-open behavior is reproduced by the [Vortex detached-watch finding](../../findings/2026-07-16-vortex-pr87-detached-watch-notification.md).
