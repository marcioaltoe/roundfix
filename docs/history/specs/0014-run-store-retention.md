---
schema: roundfix/archive-record/v1
spec: 0014-run-store-retention
title: Run Store Retention
status: archived
created: "2026-07-06"
archived: "2026-07-06"
disposition: pass
source: docs/history/specs/0014-run-store-retention
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-06.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0033
sources: []
regeneration: []
promoted: []
---

# Run Store Retention

The Run Database's Run Event Journal is append-only and never pruned, so `run_events` grows without bound — a dogfood database reached ~220 MB, almost entirely journal rows carrying every agent payload of every Run, forever. The redundant per-Batch agent log files are already going opt-in (spec 0011), but the journal itself and its on-disk artifact directories keep accumulating. This Spec bounds the Run store: terminal Runs older than a configured retention window have their journal and artifact directory pruned — on demand through a GC Command and best-effort during the preflight sweep — while the load-bearing Run state (rows and active-run locks) is never touched.
