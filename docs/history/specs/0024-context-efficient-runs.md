---
schema: roundfix/archive-record/v1
spec: 0024-context-efficient-runs
title: Context-Efficient Runs
status: archived
created: "2026-07-10"
archived: "2026-07-15"
disposition: pass
source: docs/history/specs/0024-context-efficient-runs
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-10.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0014
  - ADR-0038
  - ADR-0008
sources: []
regeneration: []
promoted: []
pull_request: "27"
delivery_commit: 3000f5973d28af660032da2d2000d9f032f4c01e
---

# Context-Efficient Runs

A measured phase-09 single-Task Run produced a 400 KB Console Log, approximately 100,000 tokens of text, including 31 full unified-diff JSON payloads and 330 echoed file reads. A replacement Agent Session spent approximately 300,000 tokens re-establishing Spec and repository context. Supervisors also had to grep that verbose Console Log for state changes, which caused a false failure report from an unrelated `Failed:` test-runner line. Roundfix must keep lossless evidence while giving Agents, Supervisors, and Detached Run callers only the information each needs.
