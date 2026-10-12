---
schema: roundfix/archive-record/v1
spec: 0261-a-cancelled-spec-leaves-a-record
title: A cancelled Spec leaves a record
status: archived
created: "2026-10-12"
archived: "2026-10-12"
disposition: pass
source: docs/specs/0261-a-cancelled-spec-leaves-a-record
source_revision: 7f22d85e5d4489efe76bba6c91740707e8a7b69a
qa_task: task_05
qa_report: qa-report-2026-10-12.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0262
sources:
  - 2026-10-12-a-decomposed-spec-cannot-be-cancelled.md
regeneration:
  - command: make baseline-digests
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/spec-workflow.json
      - internal/baseline/module-versions.json
  - command: make skills-sync
promoted: []
---

# A cancelled Spec leaves a record

QA observed that the Cancel Command records a dropped Spec’s reason, date and completed/pending Tasks, returns adopted Backlog Entries and Findings with their required bytes and states, and removes the active folder while retaining its original Git revision. Refusals preserved repository state, cancelled records remained readable, and supersession, delivery prerequisites and retry refused them. Guidance, glossary, Baseline retention and Task scope agreed with the implementation. The pre-PR Pull Request journey remains environment-blocked with every authored equivalent control recorded.
