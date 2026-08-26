---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-26T00:00:00Z
updated_at: 2026-08-26T00:00:00Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# The adapter's refusal fires first, and the catalogue is the net

Whichever unadvertised-model refusal fires first stands: all three live runtimes reach requested-model application and return `selection_rejected` from the adapter before Roundfix's own membership check can speak, so the adapter's message reaches the operator unchanged. Roundfix still reads the runtime's honest advertised catalogue — through a session ensured without the requested override, keeping preflight token-free — and keeps its membership verdict as the net beneath an adapter that declines to refuse (measured 2026-08-09: codex-acp refuses an unknown model, claude-agent-acp accepts it and echoes it back as current). Membership is decided against evidence, not against a stderr string match.

Consolidates ADR-0112 and ADR-0119 (2026-08-26); both are archived under docs/history/adr/.
