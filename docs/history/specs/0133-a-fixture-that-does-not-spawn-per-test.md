---
schema: roundfix/archive-record/v1
spec: 0133-a-fixture-that-does-not-spawn-per-test
title: The regressions this branch must not ship
status: archived
created: "2026-09-13"
archived: "2026-09-13"
disposition: pass
source: docs/history/specs/0133-a-fixture-that-does-not-spawn-per-test
source_revision: 1d1726f5d227731aa229fbc336dccbae9aca1398
qa_task: task_03
qa_report: qa-report-2026-09-13-01.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration:
  - command: make baseline-digests
promoted: []
pull_request: "188"
delivery_commit: 30069a624aa4adb01062311b43e2109b4915eaa8
---

# The regressions this branch must not ship

Specs 0119 and 0132 each repaired a real defect and each introduced a new one. Neither may reach the delivery target, so both are repaired here, on the same branch, before any of the three Specs merges.
