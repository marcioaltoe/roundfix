---
schema: roundfix/archive-record/v1
spec: 0364-a-rerun-finishes-a-partial-npm-publication
title: A rerun finishes a partial npm publication
status: archived
created: "2026-11-06"
archived: "2026-11-06"
disposition: partial
source: docs/specs/0364-a-rerun-finishes-a-partial-npm-publication
source_revision: 634a646e7decc9fa33e86dcab0588b95030dcad6
qa_task: task_04
qa_report: qa-report-2026-11-06-01.md
qa_verdict: partial
unproven:
  - the QA proves the same path on the workflow's own step scripts with a stub registry and records this as one blocked row; the first real release after delivery is the live confirmation, recorded session evidence.
adrs:
  - ADR-0374
  - ADR-0146
sources:
  - 2026-11-05-the-release-workflow-survives-a-transparency-log-409.md
regeneration:
  - command: go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage
    outputs:
      - docs/references/behavior-surfaces.json
promoted: []
type: fix
---

# A rerun finishes a partial npm publication

Offline QA proves that a rerun compares published tar content and publishes only missing coordinates, refuses different or unreadable content before publication, preserves the time-only refusal and transparency-log failure behavior, and documents the recovery consistently. The corrected fielded declaration is accepted by the Run-built public CLI, resolving the prior settlement finding. The offline archive plan (`docs/specs/0364-a-rerun-finishes-a-partial-npm-publication/qa/evidence/2026-11-06-rerun-01/archive-plan.txt` (not reachable from any branch or tag)) exits 0, records `judge: skipped (offline)` and identifies the adopted backlog candidate with `no advice (offline)`; reusable recovery knowledge already lives in the glossary, runbook and accepted ADR, so no promotion or archive is performed.
