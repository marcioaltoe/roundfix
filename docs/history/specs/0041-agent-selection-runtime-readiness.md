---
schema: roundfix/archive-record/v1
spec: 0041-agent-selection-runtime-readiness
title: Agent Selection Runtime Readiness
status: archived
created: "2026-07-17"
archived: "2026-07-21"
disposition: pass
source: docs/history/specs/0041-agent-selection-runtime-readiness
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-18.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0050
  - ADR-0055
sources: []
regeneration: []
promoted: []
pull_request: "33"
delivery_commit: f62ec2b7e550df7ea439747796354b1191896f86
---

# Agent Selection Runtime Readiness

Roundfix now owns Agent Selection Profiles, but the generated defaults and the effective ACP adapter can disagree about how an exact selection is expressed. Dogfooding with Codex CLI `0.144.5`, ACPX `0.12.0`, and a setup-generated bare `codex-acp` override that resolved to the legacy `@zed-industries/codex-acp` `0.16.0` exposed four connected failures:
