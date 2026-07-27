---
status: pending
created_at: 2026-07-26
updated_at: 2026-07-26
---

# Agent Selection — the configured Claude ACP adapter is deprecated and advertises no capability evidence (2026-07-26)

`claude/claude-opus-5/xhigh` fails Agent Selection preflight with
`capability_evidence_invalid: missing_config_options`. The failure is not
about the model: `claude-opus-5` is advertised and gets selected successfully.
Roundfix is pointed at `@zed-industries/claude-code-acp`, which npm deprecated
in favour of `@agentclientprotocol/claude-agent-acp` — the same rename Roundfix
already followed for Codex. The replacement adapter advertises the exact two
controls Roundfix requires, so the fix is an adapter migration, not a parser
change.

Nothing has broken in production yet because `roundfix watch` only routes the
`review` category, which is Codex. The first Task routed to the `frontend`
profile in `roundfix implement` will hit this at preflight.

## Session evidence

- Repository `/Users/marcio/dev/roundfix` at `d4b679e`, macOS 25.5.0.
- `acpx 0.12.1`; `claude-code-acp` on PATH resolves to
  `@zed-industries/claude-code-acp@0.16.2`; `codex-acp` resolves to
  `@agentclientprotocol/codex-acp@1.1.7`; `claude-agent-acp` is not installed.
- `.roundfixrc.yml` `profiles.frontend.preferred` is
  `runtime: claude / model: claude-opus-5 / reasoning_effort: xhigh`.
- Roundfix's own disposable preflight Session from this date is still on disk:
  `~/.acpx/sessions/0f4c9e8f-99c9-43b9-b195-950a0814d64a.json`, name
  `roundfix-preflight-ce640e44d0c56b10`, created `2026-07-26T16:32:43.316Z`,
  cwd `/Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260726T151545Z_50417ecaca9bab8d`.
  50 Session records for `claude-code-acp` share the same shape.
- The Claude adapter could not be spawned live from this session: Claude Code
  refuses to start nested (`Claude Code cannot be launched inside another
  Claude Code session`), and the documented bypass (`unset CLAUDECODE`) crashes
  active sessions. Every claim below is proven from stored Session records,
  read-only `acpx` commands that spawn no adapter, and package sources on disk.

## 1. Preflight rejects the Session before any model is evaluated

- Symptom / evidence: `acpx sessions show` on the stored preflight Session —
  a read-only command that reconnects to nothing — returns the exact payload
  Roundfix parses:

  ```console
  $ acpx --cwd /Users/marcio/.roundfix/worktrees/roundfix-339f8dac/run_20260726T151545Z_50417ecaca9bab8d \
      --format json --json-strict claude sessions show roundfix-preflight-ce640e44d0c56b10
  {"schema":"acpx.session.v1", ... ,"acpx":{"current_model_id":"claude-opus-5",
   "available_models":["default","sonnet","haiku","claude-opus-5"],
   "model_control":"legacy_set_model","session_options":{"model":"claude-opus-5"}}}
  ```

  `acpx.config_options` is absent. `current_model_id` is `claude-opus-5`, so
  the model was advertised, requested, and applied.
- Root cause: `ParseSessionCapabilitySnapshot`
  ([`selection_capabilities.go:287`](../../internal/agent/selection_capabilities.go))
  returns `missing_config_options` when `acpx.config_options` is nil, and
  `startSessionSelectionWithEnsure`
  ([`selection_assignment.go:162`](../../internal/agent/selection_assignment.go))
  propagates that `CapabilityEvidenceError` straight out of preflight.
  `PlanSelectionAssignment` never runs, so neither the model name nor the
  reasoning effort is ever considered. `capability_evidence_invalid` is
  therefore a true statement about the evidence and a misleading one about the
  Agent Selection.
- Action / suggestion: Fix the adapter (finding 2). No change to
  `selection_capabilities.go` is needed for the supported path; see finding 4
  for the diagnostic gap that remains.

## 2. The configured adapter was renamed upstream and no longer receives updates

- Symptom / evidence: the npm registry carries a deprecation notice on the
  package Roundfix configures:

  ```console
  $ curl -s https://registry.npmjs.org/@zed-industries/claude-code-acp | jq -r '.versions[."dist-tags".latest].deprecated'
  This package has been renamed to @agentclientprotocol/claude-agent-acp. Please migrate to continue receiving updates.
  ```

  Same wording, same date, and same scope move as the Codex notice Roundfix
  already acted on (`@zed-industries/codex-acp` → `@agentclientprotocol/codex-acp`).

  | | configured | replacement |
  | --- | --- | --- |
  | package | `@zed-industries/claude-code-acp` | `@agentclientprotocol/claude-agent-acp` |
  | latest | `0.16.2`, frozen since 2026-03-26 | `0.62.0`, published 2026-07-24 |
  | releases | — | 54 versions, 1.14M downloads/week |

  Reading `dist/` of the installed `0.16.2`: zero occurrences of
  `configOptions`. `session/new` returns ACP's legacy `models`
  (`availableModels` / `currentModelId`) plus `unstable_setSessionModel`, and
  `MAX_THINKING_TOKENS` is its only thinking control — an environment variable,
  not an ACP control. acpx classifies this correctly as
  `model_control: "legacy_set_model"` and omits `config_options` entirely.
- Root cause: three places pin the deprecated lineage.
  1. `~/.acpx/config.json` sets `agents.claude.command` to `claude-code-acp`.
     This is what acpx actually spawns, and Roundfix reads the same file first
     in `configuredAdapterInvocation`.
  2. [`acpx_runner.go:54`](../../internal/agent/acpx_runner.go) defaults
     `defaultAdapterCommands["claude"]` to `claude-code-acp`.
  3. [`acpx_runner.go:59-60`](../../internal/agent/acpx_runner.go) points both
     `adapterInstallCommands` entries at the deprecated `@zed-industries` scope.

  acpx 0.12.1 already defaults `claude` to
  `npx -y @agentclientprotocol/claude-agent-acp@^0.60.0`
  (`ACP_ADAPTER_PACKAGE_RANGES.claude`), so both overrides are actively
  choosing the stale adapter.
- Action / suggestion: Migrate all three, following the Codex pattern. See
  "Implementation instructions".

## 3. The replacement adapter advertises both controls Roundfix requires

- Symptom / evidence: `@agentclientprotocol/claude-agent-acp@0.62.0`,
  `dist/acp-agent.js`, has 29 references to `configOptions`. `buildConfigOptions()`
  at line 4657 emits:

  ```js
  { id: "model",  category: "model",         type: "select", currentValue: models.currentModelId, ... }
  { id: "effort", category: "thought_level", type: "select", currentValue: validEffort, ... }
  ```

  The option ids are exported constants (`MODE_CONFIG_ID`, `MODEL_CONFIG_ID`,
  `EFFORT_CONFIG_ID` at lines 4593-4595). Effort values come from the model's
  `supportedEffortLevels`, typed in `@anthropic-ai/claude-agent-sdk@0.3.219`:

  ```ts
  // package/sdk.d.ts:553
  export declare type EffortLevel = 'low' | 'medium' | 'high' | 'xhigh' | 'max';
  ```

  In the Claude Code 2.1.220 binary, the `xhigh` eligibility check excludes
  `claude-3-*`, `claude-opus-4-0`, `claude-opus-4-1`, `claude-opus-4-5`,
  `claude-opus-4-6`, `claude-sonnet-4-0`, `claude-sonnet-4-5`,
  `claude-sonnet-4-6`, and `claude-haiku-4-5`. `claude-opus-5` is not in the
  exclusion list.
- Root cause: not a defect — this is the capability Roundfix's contract was
  written against, and it exists in the replacement adapter.
- Action / suggestion: Three pieces of Roundfix already line up and need no
  change.
  - `id: "model"` with `category: "model"` satisfies the model-option detection
    in `ParseSessionConfigOptions`.
  - `id: "effort"` matches `acpxGenericReasoningEffortKey`
    ([`acpx_runner.go:42`](../../internal/agent/acpx_runner.go)), which
    `acpxReasoningEffortConfigKey`
    ([`acpx_runner.go:1477`](../../internal/agent/acpx_runner.go)) already maps
    for the `claude` runtime. acpx's `resolveCompatibleConfigId` remaps only
    legacy Zed Codex invocations, so `effort` passes through unchanged.
  - Plain model ids with no `[effort]` suffix plus a separate effort select
    resolve to `SelectionEncodingIndependent`, which is the encoding
    `PlanSelectionAssignment` already implements.

  `.roundfixrc.yml` can therefore keep `reasoning_effort: xhigh` on the
  `frontend` profile.

## 4. Roundfix verifies adapter lineage for Codex only

- Symptom / evidence: `CheckAdapter` returns bare evidence for every
  non-Codex runtime ([`acpx_runner.go:508`](../../internal/agent/acpx_runner.go)):

  ```go
  if runtimeID != "codex" {
      return evidence, nil
  }
  return inspectCodexAdapter(ctx, invocation)
  ```

  For Codex, `inspectCodexAdapter` reads `--version`, rejects a foreign package
  with `adapter_lineage_unknown`, and rejects an old version with
  `adapter_version_unsupported` — each carrying an install command. For Claude,
  any binary on PATH named `claude-code-acp` passes, and the mismatch surfaces
  much later as `capability_evidence_invalid`.
- Root cause: the lineage and version gate was built for the Codex migration
  and never generalized.
- Action / suggestion: Extend the same gate to `claude` so a machine still
  holding the deprecated adapter fails at the adapter check with a legible
  message and a `npm install -g` next action, instead of failing mid-preflight
  with a capability-evidence error that names the wrong thing.

## 5. The Codex adapter pin no longer describes the installed reality

- Symptom / evidence: `PinnedCodexAdapterVersion = "1.1.4"`
  ([`acpx_runner.go:28`](../../internal/agent/acpx_runner.go)); the installed
  `codex-acp` reports `@agentclientprotocol/codex-acp 1.1.7`; acpx 0.12.1
  requires `^1.1.5`. `roundfix setup` writes
  `npx -y @agentclientprotocol/codex-acp@1.1.4`
  ([`setup.go:360`](../../internal/cli/setup.go)).
- Root cause: the pin is a floor (`compareAdapterVersions(...) < 0`), so 1.1.7
  passes and nothing has failed. The stale constant only misleads: a fresh
  `roundfix setup` pins an exact version acpx considers below its supported
  range.
- Action / suggestion: Raise the pin to at least acpx's `^1.1.5` floor when the
  Claude migration lands. Separable from findings 1-4 and safe to defer.

## Implementation instructions

Ordered. Steps 1 and 2 are independently verifiable; do not start step 3 until
step 2 passes.

### 1. Prove the model × effort matrix

The one link this investigation could not close: whether `claude-opus-5`
reports `supportsEffort` with `xhigh` for this account at runtime. It is
resolved by the SDK at session creation, so it needs a live adapter. Run
outside a Claude Code session:

```bash
npm install -g @agentclientprotocol/claude-agent-acp
acpx --cwd /tmp --model claude-opus-5 --format json --json-strict \
  --agent claude-agent-acp sessions ensure --name rf-probe
acpx --cwd /tmp --format json --json-strict \
  --agent claude-agent-acp sessions show rf-probe
```

Pass condition: the `sessions show` payload contains `acpx.config_options`
with an option `{"id":"effort", ...}` whose `options` include `xhigh`, and an
option `{"id":"model", "category":"model", ...}` whose `options` include
`claude-opus-5`. Record the observed payload and the adapter version as an
addendum to this finding.

If `effort` is absent, `claude-opus-5` does not support effort levels on this
account. In that case the profile must drop to `reasoning_effort: ""`, which
`PlanSelectionAssignment` resolves as `SelectionEncodingModelManaged`, and the
migration is still worth doing for finding 2.

### 2. Migrate the adapter

- `~/.acpx/config.json`: set `agents.claude.command`. Either remove the
  override entirely and let acpx use its own default, or pin explicitly the way
  Codex is pinned. Pin to the version proven in step 1.
- [`acpx_runner.go:54`](../../internal/agent/acpx_runner.go): change
  `defaultAdapterCommands["claude"]`. Mirror the Codex constants — add
  `ClaudeAdapterPackage`, `PinnedClaudeAdapterVersion`,
  `ClaudeAdapterCommand()`, and `ClaudeAdapterInstallCommand()` next to their
  Codex counterparts at lines 27-28 and 76-85.
- [`acpx_runner.go:59-60`](../../internal/agent/acpx_runner.go): repoint both
  `adapterInstallCommands` entries at `@agentclientprotocol/claude-agent-acp`.
- [`setup.go:356`](../../internal/cli/setup.go): add an
  `officialClaudeACPXOverride()` beside `officialCodexACPXOverride()` so
  `roundfix setup` writes and migrates the Claude override the same way it
  already migrates a stale Codex override.

### 3. Generalize the adapter lineage gate

Turn `inspectCodexAdapter` into a runtime-parameterized check covering `codex`
and `claude`, keyed off the per-runtime package and pinned version constants,
and drop the `runtimeID != "codex"` early return at
[`acpx_runner.go:508`](../../internal/agent/acpx_runner.go). `AdapterLineageError`
and `AdapterVersionError` already carry `Command`, `Package`, `Version`, and an
install command; only the Codex-specific message construction at
`acpx_runner.go:342-350` needs to become runtime-aware.

### 4. Confirm the profile end to end

Run a `roundfix implement` preflight against a Task routed to the `frontend`
profile and confirm a `SelectionProof` with `Status: proven`,
`Encoding: independent`, `ReasoningKey: effort`, and `ReasoningValue: xhigh`.

### Operational note

acpx carries a built-in diagnostic for a known session-creation stall with
some Claude Code and `@agentclientprotocol/claude-agent-acp` combinations; its
recommended mitigation is `--approve-all` with
`nonInteractivePermissions=deny`. The current `~/.acpx/config.json` already
sets `defaultPermissions: approve-all` and `nonInteractivePermissions: deny`.
If preflight starts timing out after the migration, this is the first thing to
check.

## Suggested acceptance checks

- A Claude Agent Session created through the new adapter yields
  `SelectionCapabilities` with a non-nil `ReasoningOption` whose `ID` is
  `effort`.
- `claude/claude-opus-5/xhigh` produces a `SelectionProof` with
  `Status: proven` and `Encoding: independent`.
- A machine whose `claude` adapter resolves to `@zed-industries/claude-code-acp`
  fails `CheckAdapter` with `adapter_lineage_unknown` and a next action naming
  `@agentclientprotocol/claude-agent-acp` — not with
  `capability_evidence_invalid`.
- A machine whose Claude adapter is below the pinned version fails with
  `adapter_version_unsupported` and an install command.
- `roundfix setup` on a config holding the deprecated Claude override offers a
  migration, matching the existing Codex override behaviour.
- `roundfix doctor` reports the Claude adapter package and version as evidence,
  the way it already does for Codex.

## What worked — keep

- The disposable preflight Session survives on disk after cleanup with its
  complete capability state. Diagnosing this without spawning a single adapter
  was possible only because of that record.
- `acpx sessions show` reads the stored record without reconnecting to the
  agent, which makes it a safe read-only probe inside a constrained harness.
- `ParseSessionCapabilitySnapshot` failed closed. It refused to guess a model
  or an effort from `available_models` when the evidence contract was not met,
  which is why a deprecated adapter surfaced as a preflight error instead of a
  silently downgraded Agent Selection.
- The Codex adapter migration left a complete, reusable pattern: pinned package
  constants, a lineage probe, a version floor, typed errors carrying install
  commands, and a `roundfix setup` migration path. Finding 2 is the same
  problem, and the Codex work is the template.
