---
schema: roundfix/archive-record/v1
spec: 0046-public-context-driven-baseline-command
title: Public Context-Driven Baseline Command
status: archived
created: "2026-07-23"
archived: "2026-07-24"
disposition: qa-override
source: docs/history/specs/0046-public-context-driven-baseline-command
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: ""
qa_report: qa-report-2026-07-24.md
qa_verdict: partial
unproven: []
qa_override: true
qa_override_approval: ""
qa_override_reason: ""
qa_override_qa_outcome: ""
qa_override_qa_task_status: ""
qa_override_revision: ""
adrs:
  - ADR-0066
  - ADR-0067
  - ADR-0068
  - ADR-0069
  - ADR-0070
sources: []
regeneration: []
promoted: []
pull_request: "35"
delivery_commit: 8122ceeab1cc8fa6d46616801489e375e1610e19
---

# Public Context-Driven Baseline Command

The Context-Driven Baseline can currently be configured only through a Python implementation distributed inside the `setup-context-driven` skill. The live Fluxus adoption recorded in [the setup adoption finding](../../findings/2026-07-23-setup-context-driven-adoption-process-improvements.md) completed safely, but required repository cleanup before source inventory, repeated large audit responses, a hand-built Decision File, manual HTTP evidence, an unexplained request for a new PostgreSQL contract document, file-level plan reconstruction, and correction of persisted verification commands. Roundfix needs one public Baseline Command that humans and automations can use without Codex, preserves the safety contracts delivered by Spec 0045, and makes the observed adoption friction part of the product contract.
