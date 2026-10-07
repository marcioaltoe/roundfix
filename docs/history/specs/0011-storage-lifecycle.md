---
schema: roundfix/archive-record/v1
spec: 0011-storage-lifecycle
title: Review Artifacts, Run Logs, and Spec Archiving
status: archived
created: "2026-07-06"
archived: "2026-07-06"
disposition: pass
source: docs/history/specs/0011-storage-lifecycle
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-06.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0029
  - ADR-0030
sources: []
regeneration: []
promoted: []
---

# Review Artifacts, Run Logs, and Spec Archiving

Roundfix scatters durable output into the user-scoped Roundfix Home: Round and Review Issue artifacts land in a loose `reviews/pr-<n>/` root, and every Batch writes an agent log file that duplicates what the Run Event Journal already holds — tens of megabytes a day during dogfooding. Neither survives review or travels with the feature it belongs to, and a completed Spec has no way to retire itself out of the active set. This Spec puts each kind of durable output where it belongs: review artifacts in the repository's spec tree, per-Batch agent logs behind an opt-in switch, and completed Specs into an archive.
