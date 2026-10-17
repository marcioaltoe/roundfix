---
schema: roundfix/archive-record/v1
spec: 0275-the-f6-sender-is-named-and-a-rebuilt-item-binary-is-replaced
title: The F6 sender is named, and a rebuilt item binary is replaced
status: archived
created: "2026-10-16"
archived: "2026-10-17"
disposition: pass
source: docs/specs/0275-the-f6-sender-is-named-and-a-rebuilt-item-binary-is-replaced
source_revision: 791bc063c271249df9165d9615a642d69e7739a6
qa_task: task_03
qa_report: qa-report-2026-10-17-01.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0281
  - ADR-0266
sources:
  - 2026-10-13-trace-the-sender-of-the-f6-signal.md
regeneration: []
promoted: []
relinked:
  - docs/backlog/2026-10-16-an-item-binary-built-in-place-can-be-killed-by-code-signing.md
---

# The F6 sender is named, and a rebuilt item binary is replaced

The Archive Record may carry verdict pass: ten consecutive full suites passed, the held-binary sabotage failed with a matching CODESIGNING diagnostic and passed after restoration, and the scope and terminology audits passed. The sole environment block is the pre-PR Pull Request row. The conditional process-sender declaration remains untriggered, not proven. No new reusable knowledge beyond ADR-0281 was identified. Archive planning was not invoked because task_03 Requirement 12 forbids provider requests and prior F2 establishes that this plan can invoke one. QA performed no promotion, archive, workflow, commit, push or Pull Request operation; status and settlement remain Daemon-owned.
