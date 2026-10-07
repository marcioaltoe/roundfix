---
schema: roundfix/archive-record/v1
spec: 0135-a-rename-the-gate-sees-from-both-sides
title: A rename the gate sees from both sides
status: archived
created: "2026-09-14"
archived: "2026-09-14"
disposition: pass
source: docs/history/specs/0135-a-rename-the-gate-sees-from-both-sides
source_revision: 1d1726f5d227731aa229fbc336dccbae9aca1398
qa_task: task_03
qa_report: qa-report-2026-09-14.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "188"
delivery_commit: 30069a624aa4adb01062311b43e2109b4915eaa8
---

# A rename the gate sees from both sides

Spec 0134 made governed-mutation classification compare both directions, so a Governed Path present before the Agent turn and absent after it is a mutation. The pre-Pull-Request review of that work then proved the repair is necessary but not reachable for a rename, and corrected a premise 0134 started from.
