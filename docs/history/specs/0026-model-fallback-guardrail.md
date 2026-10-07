---
schema: roundfix/archive-record/v1
spec: 0026-model-fallback-guardrail
title: Model Fallback Guardrail
status: archived
created: "2026-07-11"
archived: "2026-07-15"
disposition: pass
source: docs/history/specs/0026-model-fallback-guardrail
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-11.md
qa_verdict: pass
unproven: []
adrs:
  - ADR-0041
sources: []
regeneration: []
promoted: []
pull_request: "27"
delivery_commit: 3000f5973d28af660032da2d2000d9f032f4c01e
---

# Model Fallback Guardrail

When a configured Agent Model fails Agent selection Preflight Validation — a model the adapter does not advertise, a rejected Default Reasoning Effort, or an adapter gap like fable before claude-code-acp shipped effort support — the Run simply fails, and the developer must diagnose adapters by hand before any work can start. Roundfix knows the runtime's Model Catalog and can prove which Agent Model actually works, so it should offer the recovery itself: propose each provider's most recent functional Agent Model at its highest functional reasoning effort. Because switching models can consume tokens very differently than the user planned, the fallback must never be silent or autonomous — Roundfix reports the problem and starts work only after explicit human confirmation.
