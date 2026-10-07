---
schema: roundfix/archive-record/v1
spec: 0023-agent-model-catalog-and-isolation
title: Agent Model Catalog and Isolation
status: archived
created: "2026-07-10"
archived: "2026-07-15"
disposition: pass
source: docs/history/specs/0023-agent-model-catalog-and-isolation
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-15.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0037
sources: []
regeneration: []
promoted: []
pull_request: "27"
delivery_commit: 3000f5973d28af660032da2d2000d9f032f4c01e
---

# Agent Model Catalog and Isolation

A 2026-07-10 dogfood Run failed before useful Agent work because Codex inherited `gpt-5.6-sol` from local configuration while the installed client had no metadata for that model. Changing the local selection to `gpt-5.5` restored normal Task settlement, proving that model selection is a Run prerequisite rather than an Agent failure. Roundfix must make the Agent Model and reasoning choice explicit, reproducible, and visible without depending on runtime-owned configuration.
