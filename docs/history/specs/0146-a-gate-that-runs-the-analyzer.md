---
schema: roundfix/archive-record/v1
spec: 0146-a-gate-that-runs-the-analyzer
title: A gate that runs the analyzer
status: archived
created: "2026-09-18"
archived: "2026-09-18"
disposition: pass
source: docs/history/specs/0146-a-gate-that-runs-the-analyzer
source_revision: 31c30719398d5c2e43af79d5868b1268fc495f10
qa_task: task_03
qa_report: qa-report-2026-09-18-01.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "220"
delivery_commit: 8b012f8eef0d142cacb62a93a04f9d5349e775c9
---

# A gate that runs the analyzer

The Go toolchain ships an analyzer that the repository's Verification gate does not run. For a long time it could not: the agent package reported thirty-one diagnostics about a mutex copied by value, so composing the analyzer into the gate would have meant composing a red gate.
