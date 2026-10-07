---
schema: roundfix/archive-record/v1
spec: 0006-acpx-run-robustness
title: ACPX Run Robustness
status: archived
created: "2026-07-05"
archived: "2026-07-06"
disposition: pass
source: docs/history/specs/0006-acpx-run-robustness
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-05.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0020
sources: []
regeneration: []
promoted: []
pull_request: "17"
delivery_commit: 5afc9d6c88bd9090cd56613b3a4e0b6b1b9ff5b6
---

# ACPX Run Robustness

Two dogfood Runs in one day lost a finished Task to the same transport failure: the Agent completed and verified its work, acpx's 10 MiB message buffer blew on a large adapter message at end of turn, the process exited 1, and Roundfix — classifying Batches by exit code alone — settled the Task failed and ended the Run Unresolved with correct work stranded uncommitted. Recovery meant hand-playing the Daemon: re-verify, rewrite the status, craft the commit. This Spec makes the agent layer classify honestly (the parsed result outranks the exit code), gives failed-but-done Tasks a one-command recovery, and applies whatever mitigation exists for the buffer limit itself.
