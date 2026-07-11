---
spec: 0025-optional-reasoning-effort
status: active
created: 2026-07-11
surfaces: [cli, backend, docs]
---

# Optional Reasoning Effort

Roundfix cannot drive the codex gpt-5.6 Agent Model family. codex-acp exposes
the `reasoning_effort` session config option only for model presets that
support more than one effort; the gpt-5.6 family manages reasoning itself, so
every `session/set_config_option "reasoning_effort"` value is rejected with
ACP -32602. Roundfix hard-requires a non-empty Default Reasoning Effort and
unconditionally issues the set call during Agent selection, so selection fails
for any gpt-5.6 model regardless of configuration.

## Goals

- A runtime configured with an empty Default Reasoning Effort selects the
  Agent Model without issuing any reasoning option, and the Run records that
  the model manages reasoning.
- An explicitly configured Default Reasoning Effort keeps today's contract: a
  runtime rejection fails Preflight Validation without fallback.
- This repository's Project Config drives codex with `gpt-5.6-sol` as the
  Default Agent Model.

## Core Features

1. **P0 - Model-managed reasoning.** Empty `runtimes.<runtime>.reasoning_effort`
   is a valid selection: Agent selection, the disposable preflight session, and
   the live Agent Session skip the reasoning set call entirely. See ADR-0040.
2. **P0 - Explicit values keep failing loudly.** A non-empty configured or
   flag-passed effort that the runtime rejects still fails Preflight
   Validation, and the failure names the model-managed remediation.
3. **P0 - Honest surfaces.** Run headers, the Run row, Attach, and Interactive
   Input represent the model-managed state instead of erroring on the missing
   value.
4. **P0 - Shipped configuration and guidance.** `.roundfixrc.yml` moves to
   `gpt-5.6-sol` with model-managed reasoning; README, CONTEXT.md, and both
   Roundfix skill copies document the optional semantics with zero sync drift.

## Non-Goals / Out of Scope

- Changing the built-in Default Agent Model or Default Reasoning Effort
  constants for any runtime.
- Per-model effort metadata in the Model Catalog or any static claim about
  which models accept which efforts.
- Tolerating or retrying a rejected non-empty reasoning value.
- Upgrading or replacing the codex-acp adapter.
