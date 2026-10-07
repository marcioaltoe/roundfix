---
schema: roundfix/archive-record/v1
spec: 0004-watch-merge-readiness
title: Watch Merge Readiness
status: archived
created: "2026-07-05"
archived: "2026-07-06"
disposition: pass
source: docs/history/specs/0004-watch-merge-readiness
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-05.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0019
  - ADR-0008
sources: []
regeneration: []
promoted: []
pull_request: "17"
delivery_commit: 5afc9d6c88bd9090cd56613b3a4e0b6b1b9ff5b6
---

# Watch Merge Readiness

The review dogfood proved the watch loop works end to end — and exposed where it wastes time and stops short. It sleeps before its first useful check even when the review finished long ago; it declares Clean the moment the local Review Issue queue empties, while the pull request still shows the Review Source's status check re-reviewing the just-pushed commit; its stdout is empty, leaving the exit code as the only machine-readable result; and its stderr buries the Daemon's milestones under the full Agent console. This Spec makes watch end exactly when the pull request is truly ready for the developer's merge decision, and makes both ends of the pipe honest.
