---
schema: roundfix/archive-record/v1
spec: 0210-evidence-snapshots-that-stay-small
title: Evidence snapshots that stay small
status: archived
created: "2026-10-01"
archived: "2026-10-01"
disposition: pass
source: docs/history/specs/0210-evidence-snapshots-that-stay-small
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_03
qa_report: qa-report-2026-10-01.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0210
sources: []
regeneration: []
promoted: []
pull_request: "324"
delivery_commit: 5aa87d2bdfd222170036aacd93f1d0a812dd987c
---

# Evidence snapshots that stay small

Spec 0202 (merged, not yet released) has the Daemon record an Evidence Snapshot when a QA pass closes: for each carriable passing row, every file a declared input matches, with its SHA-256 digest. A glob input matches every file under it. On 2026-10-01, Spec 0203's QA Report declared `internal/**`, `cmd/roundfix/**`, `go.mod`, `go.sum`, `Makefile`, `**` and the Spec directory on each of 11 rows. The report grew to 157,212 lines, 156,696 of them in the `evidence_snapshots` block, and the pre-PR review's agent session then failed with `agent/protocol error` on the diff. The operator removed the block by hand to deliver Spec 0203 (commit `a1fc8402`). Every later QA Report whose rows declare a broad glob would grow the same way, and v0.24.0 would ship that behavior.
