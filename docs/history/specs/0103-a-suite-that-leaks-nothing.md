---
schema: roundfix/archive-record/v1
spec: 0103-a-suite-that-leaks-nothing
title: A suite that leaks nothing
status: archived
created: "2026-08-12"
archived: "2026-08-15"
disposition: pass
source: docs/history/specs/0103-a-suite-that-leaks-nothing
source_revision: 1d1726f5d227731aa229fbc336dccbae9aca1398
qa_task: task_09
qa_report: qa-report-2026-08-15-01.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "163"
delivery_commit: cec521b703e1f515a4ddb5662d778efb9a4d670d
---

# A suite that leaks nothing

The test suite writes into the repository it is reading, leaves processes running for days, and carries a flake family whose common condition is spawn density — and no Roundfix command can see any of it. A test that applies a baseline against the live tree deleted tracked files while another test was copying them, so the authoritative gate can both report a false failure and destroy the operator's working tree. Four processes from the detach tests survived up to three days and burned two hours and forty minutes of CPU doing nothing; `runs list` reported no Runs, because they never registered one. A tool whose central promise is detached execution has no command that answers what it detached that is still running.
