---
status: accepted
created_at: 2026-07-10T20:01:58Z
updated_at: 2026-08-11T21:35:56Z
deprecated_at: null
superseded_by: null
---

# Roundfix owns Agent Model and reasoning selection

Inheriting an ACP Runtime's local model configuration made Runs depend on hidden machine state and allowed a model without metadata to fail after launch. Roundfix therefore resolves a concrete Agent Model and Default Reasoning Effort from its layered runtime configuration, passes both explicitly for every Agent Session, and fails Preflight Validation when the runtime does not support them. Silent fallback and mutation of runtime-owned configuration were rejected because either would make identical Roundfix inputs produce different Agent behavior.
