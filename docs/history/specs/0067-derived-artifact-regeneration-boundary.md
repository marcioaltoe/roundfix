---
schema: roundfix/archive-record/v1
spec: 0067-derived-artifact-regeneration-boundary
title: Derived artifact regeneration boundary
status: archived
created: "2026-08-01"
archived: "2026-08-05"
disposition: pass
source: docs/history/specs/0067-derived-artifact-regeneration-boundary
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: task_05
qa_report: qa-report-2026-08-05-04.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "123"
delivery_commit: 0bfa6c4e6dbad114f76e90e4b07c7fb696655688
---

# Derived artifact regeneration boundary

`internal/baseline/testdata` is listed in `DERIVED_DIGEST_PATHS`, so the natural reading is that `make baseline-digests` owns everything under it. It does not, and the gap has now produced the same failure three times.
