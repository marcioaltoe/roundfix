---
schema: roundfix/archive-record/v1
spec: 0254-a-faster-make-test
title: A faster make test
status: archived
created: "2026-10-08"
archived: "2026-10-09"
disposition: qa-override
source: docs/specs/0254-a-faster-make-test
source_revision: c74b91f92628c31480650b05d4b812685afe2909
qa_task: task_05
qa_report: qa-report-2026-10-09-01.md
qa_verdict: fail
unproven: []
qa_override: true
qa_override_approval: 'maintainer decision 2026-10-09: ''Publicar e investigar depois (Recommended)'' — publish 0254 with F6 recorded as a known failure and investigate its cause afterwards'
qa_override_reason: 'Every coverage row passed; the timing targets are met (suite at most 70 % of the start, internal/cli sequential phase at most 25 %); the seven packages pass -race -short and -shuffle=on. Known failure F6, not environmental: in the full suite the item-binary child of TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate is sometimes killed by an external signal (exit -1); the sender is unidentified and is a follow-up Backlog Entry. Remaining blocks are environmental (api.github.com, no Pull Request) and the declared post-merge CI budget.'
qa_override_qa_outcome: fail
qa_override_qa_task_status: failed
qa_override_revision: c74b91f92628c31480650b05d4b812685afe2909
adrs:
  - ADR-0259
sources: []
regeneration: []
promoted: []
---

# A faster make test

The Archive Record must retain failed QA; this child does not archive or settle Task state.
