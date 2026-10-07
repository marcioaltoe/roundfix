---
schema: roundfix/archive-record/v1
spec: 0030-context-driven-agent-instructions
title: Context-driven agent instructions
status: archived
created: "2026-07-15"
archived: "2026-07-16"
disposition: pass
source: docs/history/specs/0030-context-driven-agent-instructions
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-16-post-spec0031-rerun-07.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0046
sources: []
regeneration: []
promoted: []
---

# Context-driven agent instructions

Repository maintainers need a fast, repeatable way to establish and maintain effective agent instructions without replacing repository-specific knowledge or repeatedly answering settled setup questions. The setup workflow must provide concise, modular guidance, validate both the root agent instructions and their supporting guides, and preserve enough durable configuration to make future updates safe and mostly automatic.
