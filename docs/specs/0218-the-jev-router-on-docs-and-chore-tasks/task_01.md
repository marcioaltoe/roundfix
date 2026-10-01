---
task: task_01
spec: 0218-the-jev-router-on-docs-and-chore-tasks
status: pending
type: backend
complexity: medium
---

# Task 01: The Jev Router is one OpenCode selection that only Project Config can name

## Overview

The Jev Router can only be reached through OpenCode, configured with a
provider for OpenRouter. This Task makes
`roundfix-openrouter/typesafe/jev-router` with an empty effort the one routed
OpenCode selection, accepts it only from Project Config, gives each of its
sessions a fixed inline provider whose key is the placeholder
`{env:ROUNDFIX_OPENROUTER_API_KEY}`, and refuses it with
`jev_router_key_missing` when the variable is empty (ADR-0218). It also
documents the selection in the configuration guide.

## Requirements

1. MUST add `internal/agent/jev_router.go` with `JevRouterModel`,
   `JevRouterKeyEnv`, `JevRouterKeyMissing`, `jevRouterProviderConfig` (the
   exact JSON of the TechSpec "The inline provider and the key") and
   `IsJevRouterSelection`.
2. MUST add, at the top of `codexEnvForSession`, the routed branch the
   TechSpec states: a `SelectionFailureError` with reason
   `jev_router_key_missing: ROUNDFIX_OPENROUTER_API_KEY is not set` when the
   runner's base environment has no non-empty `ROUNDFIX_OPENROUTER_API_KEY`,
   and otherwise the single override
   `OPENCODE_CONFIG_CONTENT=<jevRouterProviderConfig>`. Every other runtime
   and OpenCode model MUST keep today's overrides byte for byte.
3. MUST make `normalizeSelection`, `validateProfiles` and `ResolveProfile`
   in `internal/config` refuse a routed selection with a non-empty effort,
   from any source other than `project`, and as a one-Run override, with the
   messages of the TechSpec "The routed selection", reading
   `agent.IsJevRouterSelection` rather than a second literal.
4. MUST add to `docs/user-guide/configuration.md`, under
   `## Agent selection profiles`, a `### Jev Router` subsection with a
   Project Config example naming the router as the `docs` preferred selection
   and the current default as its fallback, and stating: only Project Config
   may select it; its effort is `""`; Roundfix gives OpenCode the provider
   through `OPENCODE_CONFIG_CONTENT`, replacing an inherited value for routed
   sessions only; the key is `ROUNDFIX_OPENROUTER_API_KEY` and never
   `OPENROUTER_API_KEY`; and every routed prompt runs under the US$5 monthly
   Jev ceiling shared with `roundfix spec judge`.
5. MUST prove, with a sentinel key value, that the sentinel reaches no acpx
   argument and no configuration value, and MUST provide every fake OpenCode
   or acpx through the compiled test binary (ADR-0125).
6. MUST NOT read `OPENROUTER_API_KEY` anywhere, or change the Recommended
   Profile, the built-in profiles or `.roundfixrc.yml`.

## Subtasks

- [ ] Add the routed model, its provider and its key name.
- [ ] Give routed sessions the inline provider, or refuse without the key.
- [ ] Refuse the router outside Project Config and with an effort.
- [ ] Document the selection in the configuration guide.

## Acceptance Criteria

- [ ] A Project Config naming the router loads; User Config, a non-empty
      effort and a one-Run override are refused with the TechSpec messages.
- [ ] A routed session's acpx environment holds the fixed inline provider and
      no sentinel; without the key it is refused before any acpx call; other
      OpenCode models get no inline provider.
- [ ] No production file names `"OPENROUTER_API_KEY"`.

## Context

- creates: `internal/agent/jev_router.go`
- creates: `internal/agent/jev_router_test.go`
- creates: `internal/config/jev_router_test.go`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/config/profiles.go`
- interface: `docs/user-guide/configuration.md`
- instruction: `docs/adr/0218-the-jev-router-is-a-project-selected-opencode-model-under-the-shared-jev-ceiling.md`
- instruction: `internal/agent/codex_spawn.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestJevRouterLoadsFromProjectConfig|TestJevRouterRefusedFromUserConfig|TestJevRouterRefusesAReasoningEffort|TestJevRouterRefusedAsOneRunOverride|TestJevRouterSessionCarriesThePlaceholderConfig|TestJevRouterRefusedWithoutTheKey|TestOtherOpenCodeModelsGetNoInlineConfig)$' ./internal/config ./internal/agent 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestJevRouterLoadsFromProjectConfig TestJevRouterRefusedFromUserConfig TestJevRouterRefusesAReasoningEffort TestJevRouterRefusedAsOneRunOverride TestJevRouterSessionCarriesThePlaceholderConfig TestJevRouterRefusedWithoutTheKey TestOtherOpenCodeModelsGetNoInlineConfig; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for phrase in '### Jev Router' 'roundfix-openrouter/typesafe/jev-router' 'only Project Config' 'OPENCODE_CONFIG_CONTENT' 'ROUNDFIX_OPENROUTER_API_KEY'; do tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- "$phrase" || { printf 'missing phrase in configuration guide: %s\n' "$phrase" >&2; exit 1; }; done; if grep -rn --include='*.go' '"OPENROUTER_API_KEY"' internal cmd | grep -v '_test.go:'; then printf 'a production file reads the generic OpenRouter key\n' >&2; exit 1; fi` — expected: exit 0; before this Task the named tests do not exist, so the command fails.

## References

- [_prd.md](_prd.md) — Goals 1-2; User Stories 1-2; Core Features 1-3 and 7; Success Metrics 1-2; Declared breaks
- [_techspec.md](_techspec.md) — Interfaces; The routed selection; The inline provider and the key; API Contract 1; Surface Transcript 1; Testing Approach 1; Build Order 1
- ADR-0218; ADR-0037; ADR-0049; ADR-0125
