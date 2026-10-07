---
schema: roundfix/archive-record/v1
spec: 0138-a-matrix-the-spec-declares
title: A matrix the Spec declares
status: archived
created: "2026-09-14"
archived: "2026-09-16"
disposition: pass
source: docs/history/specs/0138-a-matrix-the-spec-declares
source_revision: 31c30719398d5c2e43af79d5868b1268fc495f10
qa_task: task_03
qa_report: qa-report-2026-09-16.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0155
sources: []
regeneration:
  - command: make baseline-digests
promoted: []
pull_request: "192"
delivery_commit: 84b9cae16e843a66f3a29d0881d2f1f01a813335
---

# A matrix the Spec declares

Today the QA gate decides its own matrix. The qa-gate skill and the QA contract the Daemon puts in the gate's prompt both tell the executor to cover every promise, story and criterion. As a result, two executors of the same Spec cover different sources with different inputs. Spec 0119 needed six QA Reports. Only two of its five failed reports found a defect in the Spec's own work. The rest carried noise the matrix produced:
