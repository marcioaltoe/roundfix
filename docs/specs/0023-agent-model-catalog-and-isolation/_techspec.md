---
spec: 0023-agent-model-catalog-and-isolation
prd: _prd.md
created: 2026-07-10
---

# Agent Model Catalog and Isolation - Technical Spec

## Executive Summary

Roundfix will resolve an explicit Agent Model and Default Reasoning Effort for
the selected ACP Runtime before any operational Run is created. Configuration
extends the existing typed overlay in `internal/config`; runtime catalogs and
ACP option names live at the Agent boundary; CLI resolution remains in
`internal/cli`. Availability is proved through a disposable Agent Session, not
by reading `~/.codex/config.toml`, inspecting a runtime cache, or treating the
built-in Model Catalog as an allowlist. The primary trade-off is one additional
adapter startup per operational invocation in exchange for preventing an
unsupported selection from creating a broken Run. See ADR-0039.

## System Architecture

The feature extends four existing modules and adds no cross-cutting framework:

- `internal/config` owns built-in, User Config, and Project Config values. Its
  current pointer-based overlay already provides per-key precedence; a typed
  `Runtimes` section applies that behavior independently to model and reasoning.
- `internal/agent` owns the Model Catalog, concrete `RuntimeSpec`, acpx option
  mapping, disposable Agent Session preflight, and real-session preparation.
- `internal/cli` owns invocation precedence, flags, Interactive Input, error
  routing, and the rule that selection preflight completes before `CreateRun`.
- `internal/store` persists the effective model and reasoning on the Run row so
  attach and the Live Run View do not reconstruct historical values from the
  current configuration.

Resolution flow:

```text
built-ins -> User Config -> Project Config -> invocation overrides
          -> EffectiveAgentSelection -> disposable ACP preflight
          -> Create Run -> prepare real Agent Session -> Agent work
```

`defaults.agent` remains the default ACP Runtime. The deprecated
`defaults.model` path is stripped before strict YAML decoding and emits the
existing ADR-0027 warning shape. No code reads or writes runtime-owned config.

## Implementation Design

### Interfaces

Configuration shape:

```go
type RuntimeDefaults struct {
    Model           string
    ReasoningEffort string
}

type Runtimes struct {
    Codex    RuntimeDefaults
    Claude   RuntimeDefaults
    OpenCode RuntimeDefaults
}
```

Selection and catalog boundary:

```go
type AgentSelection struct {
    Runtime         string
    Model           string
    ReasoningEffort string
}

type ModelChoice struct {
    Label, Value, Description string
}

type InvocationSelection struct {
    Model, ReasoningEffort string
    ModelSet, ReasoningEffortSet bool
}

func ResolveSelection(runtime string, defaults config.RuntimeDefaults,
    invocation InvocationSelection) (AgentSelection, error)
func ModelCatalog(runtime string) []ModelChoice
```

`ResolveSelection` lives in `internal/cli`; the Agent package never imports
configuration and therefore does not create a package cycle.

Agent Runner preflight:

```go
type ProbeRequest struct {
    Runtime RuntimeSpec
    WorkDir string
}

type Runner interface {
    Probe(context.Context, ProbeRequest) error
    Run(context.Context, ExecuteRequest, runevent.Sink) (ExecuteResult, error)
    EndSession(context.Context, RuntimeSpec, SessionRef) error
}
```

`RuntimeSpec` gains `ReasoningEffort`. `RuntimeFor` requires both resolved
values. The acpx adapter maps Codex reasoning to config id
`reasoning_effort`; Claude and OpenCode use `effort`. Session preparation is
always `sessions ensure --model <model>`, `set <reasoning-id> <value>`, then
the existing permission mode and Codex sandbox configuration. Reconnects use
the same desired options through acpx's persisted session configuration.

`ACPXRunner.Probe` first retains the pinned acpx/version checks. It then uses a
collision-resistant `roundfix-preflight-*` name, creates the disposable
session in `WorkDir` with the selected model, applies reasoning, and closes it.
Cleanup runs with a short bounded context derived from
`context.WithoutCancel`; cleanup errors are joined rather than swallowed. The
probe sends no prompt and creates no Roundfix Run Event.

### Data Models

Built-in configuration values are:

```yaml
defaults:
  agent: codex
runtimes:
  codex:    {model: gpt-5.5, reasoning_effort: xhigh}
  claude:   {model: opus, reasoning_effort: high}
  opencode: {model: "", reasoning_effort: ""}
```

Built-ins seed Codex and Claude with concrete values. User and Project Config
then overlay each key exactly as written, including an explicit empty string;
an empty effective model or reasoning value is a Preflight Validation error.
OpenCode intentionally starts with both values empty and therefore requires an
overlay or invocation override. Custom values are kept verbatim after
whitespace trimming and are validated by ACP, not by the catalog.

SQLite schema version 7 adds non-null `model` and `reasoning_effort` columns to
`runs`, both with an empty-string migration default. `store.Run` and
`CreateRunRequest` gain matching fields, and every insert/select/scan path is
updated together. Legacy rows display `-` for either missing value; every new
Agent Run must persist both concrete values. The Run's existing `agent` value
continues to identify the runtime, including the current `-custom` suffix for
an Agent Command override.

The Codex catalog is one ordered constant with the seven PRD entries and exact
identifiers. The Claude catalog is one ordered constant with five labels:
`Default`, `Opus`, `Fable`, `Sonnet`, and `Haiku`. `Default` resolves to the
effective configured Claude model and its prompt text includes that concrete
value; the remaining choices submit the adapter aliases `opus`, `fable`,
`sonnet`, and `haiku`. The labels and descriptions are presentation data, not
validation data.

### API Contracts

Operational commands that start Agent work (`resolve`, `watch`, and
`implement`) accept:

- `--model <value>`: existing one-Run override.
- `--reasoning-effort <value>`: new one-Run override.

The parser records whether each flag was explicitly visited. An unvisited value
resolves from `runtimes.<selected-runtime>` after `--agent` is known; an
explicit non-empty value outranks Project Config, User Config, and built-ins.
A visited empty value is invalid rather than a request to inherit a hidden
runtime default. `fetch` remains unchanged because it creates no Agent Session.

Interactive Input orders fields as Agent, Agent Model, Default Reasoning
Effort. It prints the selected runtime's ordered catalog and accepts either a
catalog number or a typed custom value. The effective configured value is the
default even when it is custom and therefore absent from the catalog. OpenCode
shows no fabricated catalog and requires typed/configured values. This keeps
the catalogs at exactly seven and five entries without adding a misleading
`Custom` model. Reasoning suggestions are `low`, `medium`, `high`, and `xhigh`
for Codex and `default`, `high`, and `maximum` for Claude; they are convenience
choices only, and ACP validates compatibility with the selected model.

Preflight rejection exits `2`, writes diagnostics only to stderr, and creates
zero Run rows and zero durable Agent Sessions. A typed selection error names
runtime, model, and reasoning effort, preserves the adapter rejection through
`Unwrap`, and ends with both recovery actions: update the runtime/adapter or
choose supported values. There is no fallback.

Initial progress and the Live Run View show:

```text
Agent: Codex
Agent Model: gpt-5.5
Default Reasoning Effort: xhigh
```

Inspection reads these values from the Run row. It never consults current
configuration for a historical Run.

An Agent Command override remains supported but must implement the same ACP
model/config-option contract. If it rejects or does not advertise the selected
options, the disposable preflight fails like any other runtime.

## Coverage Map

- Goal 1 and Story 1 -> `AgentSelection`, `RuntimeSpec`, real-session
  preparation
- Goal 2 and Story 2 -> `config.Runtimes`, pointer overlays, generated config
- Goal 3 and Story 3 -> Model Catalog constants and Interactive Input mapping
- Story 4 -> `--model`, `--reasoning-effort`, invocation precedence
- Story 5 -> typed custom values plus ACP validation instead of an allowlist
- Goal 4 and Story 6 -> disposable Agent Session preflight, selection error
- Story 7 -> usage guide, canonical Roundfix Skill, generated embedded Skill
- Core Feature 9 -> schema v7 Run fields and Live Run View
- Core Feature 10 -> ADR-0027 deprecated-key stripping and warning

## Integration Points

- **acpx 0.12.0**: existing stdio adapter boundary. Global `--model` selects
  the model; generic `set` applies the runtime-specific reasoning config id.
  The adapter performs authoritative availability validation.
- **codex-acp**: accepts model config and `reasoning_effort`; model changes can
  reset reasoning, which is why reasoning is always assigned after model.
- **claude-agent-acp**: accepts the user-facing aliases and config id `effort`.
- **OpenCode ACP**: accepts model plus model-variant `effort`; Roundfix supplies
  no built-in OpenCode values.

No provider API, pricing endpoint, or runtime config file becomes an
integration point.

## Testing Approach

Configuration table tests cover built-ins; User/Project per-key overlays;
mixed model/reasoning sources; strict unknown keys; generated YAML; and exactly
one warning for ignored `defaults.model`. Resolver tests cover all precedence
combinations and OpenCode's missing-value failure.

Agent unit tests pin catalog order, labels, identifiers, and descriptions.
The existing fake acpx harness asserts argument order for Codex, Claude, and
OpenCode; model-before-reasoning; custom commands; disposable session cleanup
on success, rejection, cancellation, and cleanup failure; and the absence of a
prompt. A rejection fixture must prove that neither a Run row nor a durable
Agent Session exists.

Store migration tests open schema v6 fixtures, migrate to v7, verify legacy
empty values, and round-trip new concrete values across every Run query. CLI
buffer tests pin flag parsing, stdout/stderr separation, exit `2`, Interactive
Input catalog numbering/custom entry, and effective-value output. Cockpit
model tests verify the persisted reasoning field without terminal emulation.

The final task runs `rtk make verify` and `make skills-sync`; `roundfix skills
check` must report no canonical/embedded drift.

## Build Order

1. Configuration and catalog foundations: typed per-runtime overlays,
   built-ins, deprecated-key migration, selection resolver, and catalog tests.
2. ACP selection contract: `RuntimeSpec` reasoning, runtime config-id mapping,
   deterministic real-session preparation, and adapter tests (depends on: 1).
3. Disposable selection preflight and actionable error mapping, wired before
   every operational `CreateRun` path (depends on: 2).
4. Invocation and Interactive Input surfaces: new flag, precedence resolution,
   catalogs, custom values, and command tests (depends on: 1, 3).
5. Run persistence and inspection: schema v7, store migration, progress output,
   attach, and Live Run View (depends on: 1, 4).
6. Shipped guidance: generated config examples, README/usage documentation,
   canonical Roundfix Skill updates, `make skills-sync`, and migration/recovery
   examples for all three runtimes (depends on: 3, 4, 5).

## Risks & Considerations

- Adapter vocabularies differ. Keeping the mapping in `internal/agent` avoids
  leaking `reasoning_effort`/`effort` conditionals through CLI code.
- A model can disappear after preflight. The real session intentionally repeats
  assignment and fails as infrastructure rather than silently changing values.
- Disposable sessions add startup latency and must not leak. Unique names,
  unconditional bounded cleanup, and harness coverage mitigate this.
- Catalog descriptions age faster than identifiers. They remain centralized
  presentation metadata and may be updated without changing availability rules.
- A custom Agent Command can be ACP-compatible yet lack configurable reasoning.
  Explicit ownership takes precedence; the command is rejected rather than
  inheriting a hidden default.

## Decisions

- Roundfix owns Agent Model and Default Reasoning Effort selection. See
  ADR-0037.
- Availability is validated through a disposable Agent Session. See ADR-0039.
- Model is assigned before reasoning because model changes may reset reasoning.
- The Model Catalog is a picker aid, never an allowlist.
- Effective values are persisted on the Run; legacy rows remain readable.
- OpenCode has no fabricated defaults or catalog.
- Canonical and embedded Roundfix Skills ship in sync with the CLI behavior.
