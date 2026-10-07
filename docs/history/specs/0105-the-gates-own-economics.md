---
schema: roundfix/archive-record/v1
spec: 0105-the-gates-own-economics
title: The gate's own economics
status: archived
created: "2026-08-12"
archived: "2026-08-31"
disposition: pass
source: docs/history/specs/0105-the-gates-own-economics
source_revision: 1d1726f5d227731aa229fbc336dccbae9aca1398
qa_task: task_06
qa_report: qa-report-2026-08-31.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "172"
delivery_commit: 6d45b261ef3e13bc363b7916c6fca671e87a2ece
---

# The gate's own economics

Of 201 failed Tasks measured across five repositories, 123 are the QA gate returning a verdict rather than code breaking. The gate is not too strict — it finds defects no suite catches — but its expensive tail has one cause: the Spec assumes the world behaves like its fakes, so a design premise survives authoring, implementation and every unit test and dies at the gate several Runs later. Two structural costs compound it. The Pull Request row is unreachable by design in every Spec, because the authored gate runs before a Pull Request exists, and one Spec paid six of its eight gate executions for that row alone. And the gate's own Verification is authored by hand over a verdict the Daemon already derived, which in one measured case produced a gate that passed itself having failed.
