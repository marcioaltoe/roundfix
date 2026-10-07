---
schema: roundfix/archive-record/v1
spec: 0033-console-log-tool-summary-deduplication
title: Console Log tool-summary deduplication
status: archived
created: "2026-07-16"
archived: "2026-07-17"
disposition: pass
source: docs/history/specs/0033-console-log-tool-summary-deduplication
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-17.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0008
  - ADR-0009
  - ADR-0030
sources: []
regeneration: []
promoted: []
pull_request: "32"
delivery_commit: d7fdd8cf3f986bda0cbd3683e239c799dcf6b41a
---

# Console Log tool-summary deduplication

A Detached Run can write the same compact read or edit summary multiple times when one ACP Runtime repeats identical content across a tool call's lifecycle events. The Run Event Journal is correctly lossless, but the caller-facing Console Log becomes noisy and can make one edit look like many edits. Console rendering must suppress only proven same-tool duplicates while preserving the journal and every distinct operation.
