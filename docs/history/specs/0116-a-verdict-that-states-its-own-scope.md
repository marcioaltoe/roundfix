---
schema: roundfix/archive-record/v1
spec: 0116-a-verdict-that-states-its-own-scope
title: A verdict that states its own scope
status: archived
created: "2026-08-26"
archived: "2026-08-30"
disposition: pass
source: docs/history/specs/0116-a-verdict-that-states-its-own-scope
source_revision: 1d1726f5d227731aa229fbc336dccbae9aca1398
qa_task: task_06
qa_report: qa-report-2026-08-30.md
qa_verdict: pass
unproven: []
adrs: []
sources:
  - 2026-08-14-a-clean-checker-left-eight-vacuous-commands-standing.md
regeneration: []
promoted:
  - docs/references/2026-08-14-a-clean-checker-left-eight-vacuous-commands-standing.md
pull_request: "170"
delivery_commit: fe5a038374f3deb8253832773f6eeabec339536a
---

# A verdict that states its own scope

A Spec author runs the Spec Consistency Check, reads `No findings.`, and treats the graph as checked. The check never ran a single authored Verification command, so a command that passes against the unchanged tree survives it, and the Daemon's pre-work probe refuses that Task at the starting line one Run later. A maintainer reads a QA Report and cannot tell which Roundfix produced it, so a finding raised by a binary older than the tree it audits is indistinguishable from a real one. Both are the same defect: a verdict that reports less than the reader takes it to mean, and does not say so where the reader is looking.
