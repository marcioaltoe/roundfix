---
schema: roundfix/archive-record/v1
spec: 0231-checks-that-hold-in-delivery
title: Checks that hold in delivery
status: archived
created: "2026-10-05"
archived: "2026-10-05"
disposition: pass
source: docs/history/specs/0231-checks-that-hold-in-delivery
source_revision: a30dc847a037fe584812b4471aa7da1701a9ec59
qa_task: task_04
qa_report: qa-report-2026-10-05.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0236
sources:
  - 2026-10-05-the-survivor-test-meets-eperm-from-an-exiting-owner.md
  - 2026-10-05-a-pull-request-check-ran-on-a-stale-merge.md
regeneration: []
promoted:
  - docs/references/2026-10-05-a-pull-request-check-ran-on-a-stale-merge.md
pull_request: "404"
delivery_commit: 4b5ea48b9e7686d7b8de06e7904ceabba8d2bee8
---

# Checks that hold in delivery

This is a bug fix with two causes, both recorded on 2026-10-05 in the operator's intervention log and adopted from two Backlog Entries: [the survivor test meets EPERM from an exiting owner](references/2026-10-05-the-survivor-test-meets-eperm-from-an-exiting-owner.md) and [a Pull Request check ran on a stale merge](references/2026-10-05-a-pull-request-check-ran-on-a-stale-merge.md). Each is a check that gave a verdict about something other than the work it was checking.
