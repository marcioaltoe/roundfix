---
schema: roundfix/archive-record/v1
spec: 0253-authoring-rules-that-stop-qa-reruns
title: Authoring rules that stop QA reruns
status: archived
created: "2026-10-08"
archived: "2026-10-08"
disposition: partial
source: docs/specs/0253-authoring-rules-that-stop-qa-reruns
source_revision: 1740b715c99db7a1ae7bd78e6bcd3a0b71f51ac7
qa_task: task_05
qa_report: qa-report-2026-10-08.md
qa_verdict: partial
unproven:
  - 'the operator counts, over the first five Specs delivered after this release, the QA reruns whose only defect was a transcript line and the archives with `qa_override: true`, and compares them with the audit''s five reruns in 0219–0247 and four overrides in 0236–0248'
adrs:
  - ADR-0258
sources: []
regeneration:
  - command: make baseline-digests
  - command: go test ./internal/baseline -run '^TestEveryBaselineModuleVersionIsRecorded$' -record-module-versions -count=1
    outputs:
      - internal/baseline/assets/modules/autonomous-work.json
      - internal/baseline/assets/modules/context-workflow.json
      - internal/baseline/assets/modules/go.json
      - internal/baseline/assets/modules/spec-workflow.json
      - internal/baseline/module-versions.json
  - command: make skills-sync
promoted: []
---

# Authoring rules that stop QA reruns

The authored rules, canonical/mirrored skills, embedded catalog, generated guides and reduced-retirement form are verified at the audited head. This gate does not establish reduced QA rerun/override rates after release or Merge-Ready acceptance. The Archive Record may carry this qualifying declared partial and the future metric's satisfied-by follow-up as unproven; the Daemon owns settlement and the QA commit.
