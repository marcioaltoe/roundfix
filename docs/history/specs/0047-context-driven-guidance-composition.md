---
schema: roundfix/archive-record/v1
spec: 0047-context-driven-guidance-composition
title: Context-Driven Baseline guidance composition
status: archived
created: "2026-07-24"
archived: "2026-07-25"
disposition: pass
source: docs/history/specs/0047-context-driven-guidance-composition
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: ""
qa_report: qa-report-2026-07-25.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "35"
delivery_commit: 8122ceeab1cc8fa6d46616801489e375e1610e19
---

# Context-Driven Baseline guidance composition

The public Baseline Command can generate a compact root index and modular agent guides, but the accepted Fluxus greenfield result still required maintainers to redistribute rules from a catch-all carrier, infer instruction precedence, and complete the ADR and Findings guidance manually. This feature makes the generated guidance self-contained and gives every accepted rule one semantic owner while preserving the retention, confirmation, and rollback guarantees delivered by Specs 0045 and 0046. The [greenfield acceptance finding](../../findings/2026-07-24-greenfield-agent-guidance-acceptance-target.md) is the acceptance source.
