---
schema: roundfix/archive-record/v1
spec: 0154-evidence-when-the-target-branch-is-gone
title: Evidence when the target branch is gone
status: archived
created: "2026-09-23"
archived: "2026-09-24"
disposition: pass
source: docs/history/specs/0154-evidence-when-the-target-branch-is-gone
source_revision: 31c30719398d5c2e43af79d5868b1268fc495f10
qa_task: task_04
qa_report: qa-report-2026-09-23-02.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration:
  - command: make skills-sync
  - command: make baseline-digests
promoted: []
pull_request: "237"
delivery_commit: 786abcf2f5fc80cc06a86d5aaf9e7d145481afce
---

# Evidence when the target branch is gone

Reconciliation asks one question about a terminal Run: did its work reach the target branch? It answers with `merge-base --is-ancestor`, and when that misses it looks for content evidence — a superseding QA Report on the target for the same Spec.
