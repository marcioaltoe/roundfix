---
schema: roundfix/archive-record/v1
spec: 0064-spec-artifact-consistency-gate
title: Spec artifact consistency gate
status: archived
created: "2026-08-01"
archived: "2026-08-04"
disposition: pass
source: docs/history/specs/0064-spec-artifact-consistency-gate
source_revision: 436b2d919ede38fc45ec654c236e3f863f304756
qa_task: task_09
qa_report: qa-report-2026-08-04-02.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "108"
delivery_commit: 9c744cf498635ffdf52d6adb0009142c4d80d174
---

# Spec artifact consistency gate

Half of the QA findings this repository paid twenty-minute cycles for were not code defects — they were contradictions between a Spec's own artifacts, each detectable by reading files. Spec 0058's QA-001 found a PRD promising a preflight check npm makes impossible while its own ADR acknowledged the limit; QA-004 found the workflow emitting five failure prefixes while the runbook documented four. Spec 0056 repeated both shapes: F-001 found Project Constraints omitting an ADR that governs the changed behavior, and F-002 found a Core Feature contradicting the TechSpec decision that superseded it. Four findings, four full gate cycles, zero of them needing a running binary.
