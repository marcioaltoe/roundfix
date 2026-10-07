---
schema: roundfix/archive-record/v1
spec: 0045-context-driven-baseline-0-0-1-reset
title: Roundfix 0.0.1 Context-Driven Baseline reset
status: archived
created: "2026-07-22"
archived: "2026-07-24"
disposition: pass
source: docs/history/specs/0045-context-driven-baseline-0-0-1-reset
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: ""
qa_report: qa-report-2026-07-23.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0060
  - ADR-0061
  - ADR-0062
  - ADR-0063
sources: []
regeneration: []
promoted: []
pull_request: "35"
delivery_commit: 8122ceeab1cc8fa6d46616801489e375e1610e19
---

# Roundfix 0.0.1 Context-Driven Baseline reset

A real managed upgrade removed operational guidance while its generated corpus still passed internal coverage checks: 27 coarse entries stood in for a 573-line instruction source, and the generated agent-guidance set deleted 187 lines while adding 73. Roundfix needs a project-agnostic, independently auditable baseline that preserves every governed instruction and structured contract, gives the repository explicit control over non-portable policy, and restarts all Roundfix-owned version surfaces together at `0.0.1`.
