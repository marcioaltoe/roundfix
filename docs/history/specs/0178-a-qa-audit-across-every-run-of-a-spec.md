---
schema: roundfix/archive-record/v1
spec: 0178-a-qa-audit-across-every-run-of-a-spec
title: A QA audit across every Run of a Spec
status: archived
created: "2026-09-28"
archived: "2026-09-28"
disposition: pass
source: docs/history/specs/0178-a-qa-audit-across-every-run-of-a-spec
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_05
qa_report: qa-report-2026-09-28-01.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0138
sources:
  - 2026-09-14-the-mechanical-stage-audits-only-the-current-runs-task-commits.md
  - 2026-09-25-the-auditor-staleness-signal-is-recorded-but-never-acted-on.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "270"
delivery_commit: 08fd9b53edbdee77212742c9a30469c2bf34b9c5
---

# A QA audit across every Run of a Spec

The terminal QA gate must audit what a Spec delivers, not what its last Run happened to commit, and the auditor identity a QA Report records must mean one thing and be acted on. Two defects break that today:
