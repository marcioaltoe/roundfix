---
schema: roundfix/archive-record/v1
spec: 0212-gates-and-run-storage-that-let-a-correct-delivery-finish
title: Gates and Run storage that let a correct delivery finish
status: archived
created: "2026-10-01"
archived: "2026-10-02"
disposition: pass
source: docs/history/specs/0212-gates-and-run-storage-that-let-a-correct-delivery-finish
source_revision: ae593da46d28e20d4c905ea1ebbeca5cbe8d3b90
qa_task: task_05
qa_report: qa-report-2026-10-02-01.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0212
sources:
  - 2026-10-01-the-pre-pr-review-fails-on-a-large-vendored-diff.md
  - 2026-10-01-an-imported-qa-pass-can-carry-a-file-that-breaks-the-gate.md
  - 2026-10-01-task-context-has-no-kind-for-a-deleted-path.md
  - 2026-10-01-reconcile-releases-the-runs-of-a-squash-merged-spec.md
regeneration:
  - command: make skills-sync
promoted: []
pull_request: "345"
delivery_commit: 35c6714e92568adcded211ddbf1b043f081adcff
---

# Gates and Run storage that let a correct delivery finish

Between 2026-09-25 and 2026-10-01, four deliveries whose work was correct stopped on a gate or on the disk, and each needed an operator. The pre-PR review of Spec 0200 failed three times with `agent/protocol error`, twice in the Delivery Queue and once by hand, on a candidate diff of 1,295,055 bytes, and someone reviewed it by hand. In Spec 0203's delivery a QA Agent wrote its evidence as a Go test file; the next QA pass imported it byte for byte, the repository gate refused to format it before any row ran, and every retry would have imported it again. Spec 0200's task_04 deletes a skill file, but a Task's `## Context` has no kind for a deleted path, so the Spec Consistency Check refused the correct Task once it ran, and the operator relabelled the entry as a path the Task creates. And `roundfix reconcile` keeps every Run of a Spec merged by squash: on 2026-10-01 this repository held 32 Run Worktrees of merged Specs, 19.1 GB, which contributed to a disk-full stop of a delivery.
