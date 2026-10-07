---
schema: roundfix/archive-record/v1
spec: 0028-settlement-and-reporting
title: Settlement and Reporting Robustness
status: archived
created: "2026-07-14"
archived: "2026-07-15"
disposition: pass
source: docs/history/specs/0028-settlement-and-reporting
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-15.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0044
sources: []
regeneration: []
promoted: []
---

# Settlement and Reporting Robustness

Field sessions running spec Runs (specs 0004/0005, 0009/0010 dogfood cycles) lost work or misled operators through settlement and reporting gaps that are each small but compound badly. A killed run process leaves an orphaned Active-Run lock that blocks every relaunch until a manual force stop. An Agent that writes an obvious status synonym (`done` instead of `completed`) voids an entire batch of finished work. A Task whose Verification was already green on the baseline settles `completed` with a commit containing nothing but the status flip. The Settle Command sweeps every change in a shared worktree into one task's commit without saying so. The final report names failed items but never says why they failed, forcing journal archaeology. And a missing ACP adapter binary surfaces late, as an opaque spawn failure with no install hint. This spec closes those gaps.
