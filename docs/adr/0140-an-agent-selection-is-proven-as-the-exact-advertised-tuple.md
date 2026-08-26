---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-26T00:00:00Z
updated_at: 2026-08-26T00:00:00Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# An Agent Selection is proven as the exact advertised tuple

Before a Run exists, Roundfix opens a uniquely named disposable Agent Session, discovers the controls the adapter actually advertises, derives an adapter-specific assignment plan, applies it, and proves the session represents the exact requested runtime/model/reasoning tuple — rejecting static family assumptions and unproven overrides, and failing Preflight without fallback on any rejection. An empty Default Reasoning Effort is itself a valid selection meaning the model manages reasoning: the set call is skipped on both the disposable and the live session, and the effective selection records the model-managed state. A static catalog or runtime-owned config inspection cannot prove that the exact selection will start; the extra adapter startup buys that proof.

Consolidates ADR-0039, ADR-0040, and ADR-0055 (2026-08-26); all are archived under docs/history/adr/.
