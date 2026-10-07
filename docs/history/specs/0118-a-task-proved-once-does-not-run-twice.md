---
schema: roundfix/archive-record/v1
spec: 0118-a-task-proved-once-does-not-run-twice
title: A Task proved once does not run twice
status: archived
created: "2026-08-27"
archived: "2026-08-27"
disposition: pass
source: docs/history/specs/0118-a-task-proved-once-does-not-run-twice
source_revision: 1d1726f5d227731aa229fbc336dccbae9aca1398
qa_task: task_07
qa_report: qa-report-2026-08-27.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0115
  - ADR-0023
sources:
  - 2026-08-12-five-unresolved-runs-to-deliver-one-spec.md
regeneration: []
promoted: []
pull_request: "169"
delivery_commit: b1d60aee619afc287edcbfa0fb1e29c3e200a938
---

# A Task proved once does not run twice

An Implement Run commits each settled Task inside its Run Worktree and moves those commits to the user's branch only when the Run reaches Clean. A Run that ends with an Unresolved Outcome therefore leaves every Task it completed behind on its Run Branch, and the checkout still reads `status: pending` for work that ran, passed its Verification, and was committed. The next Implement Run reads that checkout and re-executes it — not re-checks it: re-executes it, with a fresh Agent turn and a fresh Verification cycle. Roundfix already owns the command that hands proved work back, and that command refuses the one outcome that produces the problem. This Spec makes the remedy reachable, and makes the Implement Command stop before spending an Agent turn on work a prior Run already proved.
