---
schema: roundfix/archive-record/v1
spec: 0134-a-governed-deletion-the-gate-can-see
title: A governed deletion the gate can see
status: archived
created: "2026-09-13"
archived: "2026-09-14"
disposition: pass
source: docs/history/specs/0134-a-governed-deletion-the-gate-can-see
source_revision: 1d1726f5d227731aa229fbc336dccbae9aca1398
qa_task: task_04
qa_report: qa-report-2026-09-14-01.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "188"
delivery_commit: 30069a624aa4adb01062311b43e2109b4915eaa8
---

# A governed deletion the gate can see

The operation authority this branch built refuses a Task that changes a Governed Path without the `commit` operation. It does not refuse a Task that *removes* one. The classifier compares the paths before and after the Agent turn and reports a governed mutation only for a path that appears in the second list and not the first, so every deletion is invisible to it — and a rename reads as a deletion plus an unrelated addition. A Task can drop `Makefile` without any operation being asked for.
