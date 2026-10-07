---
schema: roundfix/archive-record/v1
spec: 0025-optional-reasoning-effort
title: Optional Reasoning Effort
status: archived
created: "2026-07-11"
archived: "2026-07-15"
disposition: pass
source: docs/history/specs/0025-optional-reasoning-effort
source_revision: b785bc9ad6165874b2cce480f8da6ef81b584c00
qa_task: ""
qa_report: qa-report-2026-07-11.md
qa_verdict: pass
unproven: []
adrs: []
sources: []
regeneration: []
promoted: []
pull_request: "27"
delivery_commit: 3000f5973d28af660032da2d2000d9f032f4c01e
---

# Optional Reasoning Effort

Roundfix cannot drive the codex gpt-5.6 Agent Model family. codex-acp exposes the `reasoning_effort` session config option only for model presets that support more than one effort; the gpt-5.6 family manages reasoning itself, so every `session/set_config_option "reasoning_effort"` value is rejected with ACP -32602. claude-code-acp does not implement `session/set_config_option` at all, so every effort value is rejected for every Claude model. Roundfix hard-requires a non-empty Default Reasoning Effort and unconditionally issues the set call during Agent selection, so selection fails for any gpt-5.6 model regardless of configuration, and the claude runtime's built-in defaults (`opus` + `high`) fail the same way.
