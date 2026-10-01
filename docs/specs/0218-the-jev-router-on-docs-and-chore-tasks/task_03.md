---
task: task_03
spec: 0218-the-jev-router-on-docs-and-chore-tasks
status: pending
type: backend
complexity: high
---

# Task 03: Every routed prompt runs under the shared ceiling and leaves a Judge Log line

## Overview

Tasks 01 and 02 give Roundfix a routed selection and a way to read the
month's Jev spend. Nothing yet stops a routed prompt at the ceiling or
records what it cost. This Task adds a `JevRouterGate` to the Daemon's
dependencies and calls it around each routed prompt in
`agentSessionOwner.runPrepared`: a refusal before Agent work begins moves the
Fallback Chain on, a refusal after it fails the Work Item (ADR-0050,
ADR-0114), and each routed prompt appends one `router-prompt` line
(ADR-0218). It also describes the router in the Roundfix Skill.

## Requirements

1. MUST add `JevRouterGate` and `Dependencies.JevRouter` in
   `internal/daemon/engine.go`, defaulting when nil to a gate built from the
   process environment, Roundfix Home, the `https://openrouter.ai/api/v1`
   endpoint and the judge's loaded ceiling.
2. MUST call `Before` in `agentSessionOwner.runPrepared` only when the active
   runtime and model satisfy `agent.IsJevRouterSelection`, and MUST turn its
   refusal into the `SelectionFailureError` reasons the TechSpec "The gate
   and the record" states (`jev_spend_unreadable: <error>`,
   `jev_ceiling_reached: month's Jev spend US$<total> of US$<ceiling>`), so
   the existing `Run` loop falls back before work and fails the Work Item
   after it. A refused prompt MUST NOT reach the runner.
3. MUST call `After` once per routed prompt that ran, with the Run, Spec,
   scope, category, attempt, latency, tokens, the usage `Before` returned and
   whether the prompt failed. A failed `After` MUST be written to the Run's
   progress stream and MUST NOT change the prompt's result.
4. MUST leave non-routed prompts byte-for-byte as today, with no key
   endpoint call.
5. MUST add a `### Jev Router` section to
   `.agents/skills/roundfix/references/runtime.md` stating the routed
   selection, the Project Config rule, the key variable, the shared US$5
   ceiling, and the refusal classifications `jev_router_key_missing`,
   `jev_spend_unreadable` and `jev_ceiling_reached`; raise both version
   fields of `.agents/skills/roundfix/SKILL.md`; run `make skills-sync`; and
   re-record the version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
6. MUST test with a fake gate and a fake runner only, with a sentinel key
   that appears in no line, and MUST NOT change the `### QA settlement`
   section of any skill.

## Subtasks

- [ ] Add the gate dependency and its default.
- [ ] Call it before and after routed prompts in the session owner.
- [ ] Prove fallback before work, failure after work, and one line per prompt.
- [ ] Describe the router in the Roundfix Skill and raise its version.

## Acceptance Criteria

- [ ] At the ceiling, the first routed prompt never runs and the fallback
      selection starts after its notification; after work began, the Work
      Item fails with `jev_ceiling_reached`.
- [ ] An unreadable spend falls back with `jev_spend_unreadable`.
- [ ] Below the ceiling, one routed prompt appends one `router-prompt` line
      whose cost is the usage change, and a non-routed prompt calls no gate.

## Context

- creates: `internal/daemon/jev_router_gate_test.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/daemon/agent_session_owner.go`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/runtime.md`
- interface: `skills/roundfix/references/runtime.md`
- interface: `skills/testdata/owned-skill-versions.json`
- instruction: `docs/adr/0218-the-jev-router-is-a-project-selected-opencode-model-under-the-shared-jev-ceiling.md`
- instruction: `docs/adr/0050-configured-fallbacks-activate-after-notification.md`
- instruction: `docs/adr/0114-opening-an-agent-session-is-not-agent-work.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestJevRouterCeilingFallsBackBeforeWork|TestJevRouterCeilingFailsTheTaskAfterWork|TestJevRouterUnreadableSpendFallsBack|TestJevRouterPromptAppendsOneLine|TestNonRoutedPromptCallsNoGate)$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestJevRouterCeilingFallsBackBeforeWork TestJevRouterCeilingFailsTheTaskAfterWork TestJevRouterUnreadableSpendFallsBack TestJevRouterPromptAppendsOneLine TestNonRoutedPromptCallsNoGate; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for phrase in '### Jev Router' 'roundfix-openrouter/typesafe/jev-router' 'jev_ceiling_reached' 'jev_router_key_missing' 'ROUNDFIX_OPENROUTER_API_KEY'; do tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/runtime.md | grep -qF -- "$phrase" || { printf 'missing phrase in the runtime reference: %s\n' "$phrase" >&2; exit 1; }; done; diff -r .agents/skills/roundfix skills/roundfix >/dev/null || { printf 'mirror differs: skills/roundfix\n' >&2; exit 1; }; vout="$(go test -count=1 -v -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills 2>&1)" || { printf '%s\n' "$vout"; exit 1; }; printf '%s\n' "$vout" | grep -q -- '--- PASS: TestEveryOwnedSkillVersionIsRecorded' || { printf '%s\n' "$vout"; exit 1; }; make skills-sync-check` — expected: exit 0; before this Task the named Daemon tests do not exist, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 3; User Story 3; Core Features 4, 5 and 7; Success Metric 3
- [_techspec.md](_techspec.md) — Interfaces; The gate and the record; API Contracts 3-4; Testing Approach 3-4; Build Order 3
- ADR-0218; ADR-0050; ADR-0114; ADR-0189
