---
spec: 0023-agent-model-catalog-and-isolation
status: archived
created: 2026-07-10
surfaces: [frontend, cli, data, infra, docs]
archived: "2026-07-15"
source_slug: 0023-agent-model-catalog-and-isolation
---


# Agent Model Catalog and Isolation

A 2026-07-10 dogfood Run failed before useful Agent work because Codex inherited `gpt-5.6-sol` from local configuration while the installed client had no metadata for that model. Changing the local selection to `gpt-5.5` restored normal Task settlement, proving that model selection is a Run prerequisite rather than an Agent failure. Roundfix must make the Agent Model and reasoning choice explicit, reproducible, and visible without depending on runtime-owned configuration.

## Goals

- Every Agent Session starts with an explicit Agent Model and reasoning effort selected by Roundfix.
- Project Config can pin separate model and reasoning defaults for Codex, Claude, and OpenCode without changing runtime-owned configuration.
- Interactive Input presents the known Codex and Claude Model Catalogs while non-interactive callers retain custom model values.
- An unavailable model or reasoning value fails Preflight Validation before a Run is created and names the next useful action.

## User Stories

1. As a developer starting a Run, I want Roundfix to select the Agent Model and reasoning explicitly, so that changes in local Codex, Claude, or OpenCode configuration cannot break or alter the Run.
2. As a repository owner, I want Project Config defaults for each ACP Runtime, so that the repository chooses reproducible Agent behavior while retaining User Config as the lower-precedence layer.
3. As a developer using Interactive Input, I want an ordered Model Catalog for the selected ACP Runtime, so that I can choose a known model without memorizing identifiers.
4. As an automation author, I want model and reasoning overrides for one Run, so that an exceptional workload does not require editing shared configuration.
5. As a developer with access to a newly released model, I want to pass a custom model value, so that the built-in Model Catalog does not block forward compatibility.
6. As a developer whose selected model lacks runtime metadata, I want an actionable Preflight Validation failure before the Run exists, so that I can update the runtime or select an available model without recovering a failed Run.
7. As a Supervisor or developer adopting the new configuration, I want the Roundfix Skill and usage documentation to explain defaults, precedence, overrides, and recovery, so that I can operate the feature without reading source code.

## Core Features

1. **P0 - Runtime defaults.** User Config and Project Config accept `runtimes.<runtime>.model` and `runtimes.<runtime>.reasoning_effort` independently for `codex`, `claude`, and `opencode`. Existing per-key overlay precedence remains built-in, then User Config, then Project Config. See ADR-0037.
2. **P0 - Built-in selections.** Codex defaults to `gpt-5.5` with `xhigh` reasoning. Claude defaults to `opus` with `high` reasoning. OpenCode has no built-in model or reasoning; selecting it without both configured fails Preflight Validation.
3. **P0 - Explicit launch selection.** Every Run that starts Agent work resolves a concrete model and reasoning value and assigns both to its Agent Session. Runtime-owned model and reasoning defaults never participate in that selection. See ADR-0037.
4. **P0 - Availability preflight.** Before creating a Run, Roundfix checks that the selected ACP Runtime accepts the effective Agent Model and reasoning effort. Failure identifies the rejected values and tells the user to update the runtime or select supported values. Roundfix never substitutes another model silently. See ADR-0037.
5. **P0 - Codex Model Catalog.** Interactive Input presents these entries in order: `gpt-5.6-sol` (latest frontier agentic coding model), `gpt-5.6-terra` (balanced everyday agentic coding), `gpt-5.6-luna` (fast and affordable agentic coding), `gpt-5.5` (initial Default Agent Model), `gpt-5.4` (everyday coding), `gpt-5.4-mini` (small and cost-efficient), and `gpt-5.3-codex-spark` (ultra-fast coding).
6. **P0 - Claude Model Catalog.** Interactive Input presents `Default`, `Opus`, `Fable`, `Sonnet`, and `Haiku` in that order. The built-in Default resolves to Opus 4.8 with a 1M context window; the other labels represent Opus 4.8, Fable 5, Sonnet 5, and Haiku 4.5 respectively. A Project Config override changes what Default resolves to for that repository.
7. **P1 - Per-Run overrides.** The model override remains available and a reasoning-effort override is added to non-interactive commands and Interactive Input. Explicit invocation values outrank Project Config, User Config, and built-ins.
8. **P1 - Forward-compatible custom values.** Non-interactive model and reasoning values are not restricted to the built-in catalogs, but they must pass the same availability preflight as catalog entries.
9. **P1 - Effective selection visibility.** Run output and existing Run inspection surfaces show the concrete Agent Model and reasoning effort, never an ambiguous `auto` value for Agent work.
10. **P1 - Config migration.** The global `defaults.model` key is deprecated and ignored with one stderr warning that points to the per-runtime replacement, following ADR-0027. Generated User Config and Project Config examples use only the per-runtime structure.
11. **P0 - Shipped guidance.** The Roundfix Skill and user documentation ship with the behavior change. They document the per-runtime configuration shape, built-in defaults, overlay and invocation precedence, both catalogs, custom values, Interactive Input, non-interactive overrides, effective-selection inspection, deprecated-key migration, and recovery from unsupported model or reasoning failures. Copy-paste examples cover Codex, Claude, and the explicit-configuration requirement for OpenCode.

## User Experience

Interactive Input first selects the ACP Runtime, then offers its Model Catalog and reasoning choices with the effective Project Config or User Config values selected. The Claude `Default` label shows the concrete model it resolves to. Non-interactive commands accept the existing model override and the new reasoning-effort override; diagnostics remain on stderr.

If the selected model or reasoning effort is unavailable, the command exits during Preflight Validation before creating a Run. The error names the runtime, model, and reasoning effort and gives two recovery paths: update the runtime or choose supported values. No flow edits the user's Codex, Claude, or OpenCode configuration.

## Non-Goals / Out of Scope

- Automatically installing or upgrading Codex, Claude, OpenCode, acpx, or their model metadata.
- Falling back to another Agent Model or reasoning effort when the selected values are unavailable.
- Isolating credentials, permissions, provider endpoints, or other runtime-owned settings unrelated to model and reasoning selection.
- Adding a built-in OpenCode Model Catalog before its required choices are defined.
- Model pricing, quota-aware routing, task-by-task automatic model choice, or model benchmarking.
- Providing LLM credentials required by a target repository's Verification commands.

## Success Metrics

- 100% of Agent Sessions created by the test suite carry a non-empty concrete Agent Model and reasoning effort.
- A selected model with missing metadata creates zero Runs and starts zero Agent Sessions.
- The Codex picker exposes exactly seven ordered catalog entries and the Claude picker exposes exactly five.
- Changing runtime-owned model or reasoning configuration does not change the effective values of a Run with identical Roundfix inputs.
- Project Config overrides User Config independently for each runtime model and reasoning key; explicit invocation values override both.
- A deprecated `defaults.model` key produces exactly one stderr warning and never appears on stdout.
- The canonical Roundfix Skill matches every new flag, configuration key, default, Preflight Validation outcome, and recovery path; the embedded skill copy has zero sync drift.
- Usage documentation contains working examples for Project Config defaults and one-Run model/reasoning overrides for Codex and Claude.

## Decisions

- Roundfix owns the effective Agent Model and reasoning selection for every Agent Session. See ADR-0037.
- Codex initially defaults to `gpt-5.5` because current dogfood metadata accepts it; `gpt-5.6-sol` remains selectable and can become the default only through a deliberate validated change.
- Per-runtime configuration uses `runtimes.<runtime>.model` and `runtimes.<runtime>.reasoning_effort`.
- Codex defaults to `xhigh`, Claude defaults to `high`, and OpenCode requires explicit configuration.
- Custom values remain available but must pass Preflight Validation; no silent fallback is allowed.
- Bare `Fable` no longer names the supervising role; Supervisor is the canonical role and Fable remains a Claude Agent Model label.
- The Roundfix Skill and usage documentation are part of the feature's completion contract, not follow-up work.

## Open Questions

None. The TechSpec must verify the concrete Claude launch identifiers and each runtime's reasoning vocabulary against the supported adapter versions without changing these user-facing labels.
