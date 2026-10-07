---
schema: roundfix/archive-record/v1
spec: 0248-the-regeneration-contract-runs-when-its-inputs-change
title: The regeneration contract runs when its inputs change
status: archived
created: "2026-10-07"
archived: "2026-10-07"
disposition: pass
source: docs/specs/0248-the-regeneration-contract-runs-when-its-inputs-change
source_revision: 36b693fac941a733a937234aeb465921531ef9ff
qa_task: task_05
qa_report: qa-report-2026-10-07.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0253
sources: []
regeneration: []
promoted: []
---

# The regeneration contract runs when its inputs change

`TestRegenerationIsDeclared` proves that every path the declared regeneration commands rewrite is covered by a `derived_paths` declaration in `.roundfixrc.yml`. Spec 0245 added it after a module edit regenerated paths no declaration covered, which entry 233 of the operator's queue log records. Spec 0247 gave it the Contract Relevance `boundary`. The Archive Record is `docs/history/specs/0247-a-run-gate-that-runs-the-repository-contracts.md`.
