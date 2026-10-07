---
schema: roundfix/archive-record/v1
spec: 0001-implement-command
title: Implement Command
status: archived
created: "2026-07-04"
archived: "2026-07-06"
disposition: partial
source: docs/history/specs/0001-implement-command
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-04.md
qa_verdict: partial
unproven: []
adrs:
  - ADR-0012
  - ADR-0013
  - ADR-0014
  - ADR-0015
sources: []
regeneration: []
promoted: []
pull_request: "17"
delivery_commit: 5afc9d6c88bd9090cd56613b3a4e0b6b1b9ff5b6
---

# Implement Command

Roundfix today resolves review feedback on Open Pull Requests, but the Specs produced by the planning workflow — PRD, Task Graph, task files — are still executed by hand, Task by Task. The Implement Command makes Roundfix execute a Spec's Task Graph: it runs Agents over the Tasks in dependency order, verifies and commits each one through the Daemon, and optionally ends with a QA gate — so a completed Spec becomes a verified local branch with one command.
