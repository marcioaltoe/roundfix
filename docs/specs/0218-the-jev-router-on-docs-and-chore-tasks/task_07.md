---
task: task_07
spec: 0218-the-jev-router-on-docs-and-chore-tasks
status: pending
type: backend
complexity: medium
---

# Task 07: A routed prompt starts only under a monthly key limit within the ceiling

## Overview

Corrective Task from the 2026-10-02 measurement. The gate refuses a routed prompt when the month's spend has already reached the US$5 ceiling, but it cannot stop a prompt that is already running. One routed prompt cost US$6.81 in flight and took the month's Jev spend to US$7.44. Only OpenRouter can stop a running prompt, through the key's own credit limit. `GET /api/v1/key` reports `data.limit`, `data.limit_remaining` and `data.limit_reset` (measured 2026-10-02: `limit` 50, `limit_reset` null). The gate therefore also requires a server-side monthly limit no larger than the ceiling.

## Requirements

1. MUST extend the key read in `internal/jevrouter/key.go` to return `data.limit` (number or null), `data.limit_remaining` (number or null) and `data.limit_reset` (string or null) alongside `data.usage_monthly`, keeping the existing refusals for a malformed `usage_monthly`.
2. MUST refuse a routed prompt before it starts with the reason `jev_router_key_unbounded` unless the key reports `limit_reset` equal to `monthly` and a numeric `limit` no greater than the monthly ceiling. The refusal follows the existing gate path: fallback before work, and the Task fails after work began. Its message tells the maintainer to set a monthly credit limit of at most the ceiling on the key at OpenRouter.
3. MUST refuse with the existing ceiling reason when `limit_remaining` is a number at or below zero.
4. MUST add the tests named in Verification, with a fake key endpoint as `internal/jevrouter/key_test.go` already uses, covering: no limit, a lifetime limit, a monthly limit above the ceiling, an accepted monthly limit, and an exhausted `limit_remaining`. MUST leave existing tests passing.
5. MUST describe the requirement in the Roundfix Skill's runtime reference and the configuration guide where the router is described, raising the Roundfix Skill version, running `make skills-sync` and re-recording the version.

## Subtasks

- [ ] Read the key's limit fields.
- [ ] Refuse an unbounded or over-ceiling key before a routed prompt starts.
- [ ] Document the requirement.

## Acceptance Criteria

- [ ] A routed prompt never starts unless OpenRouter itself caps the key's monthly spend at or below the ceiling.

## Context

- interface: `internal/jevrouter/key.go`
- interface: `internal/jevrouter/spend.go`
- interface: `internal/jevrouter/key_test.go`
- creates: `internal/jevrouter/key_limit_test.go`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/runtime.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/roundfix/references/runtime.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/configuration.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestARoutedPromptNeedsAMonthlyKeyLimit|TestALifetimeKeyLimitIsRefused|TestAMonthlyKeyLimitAboveTheCeilingIsRefused|TestAMonthlyKeyLimitWithinTheCeilingIsAccepted|TestAnExhaustedKeyLimitIsRefused)$" ./internal/jevrouter 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestARoutedPromptNeedsAMonthlyKeyLimit TestALifetimeKeyLimitIsRefused TestAMonthlyKeyLimitAboveTheCeilingIsRefused TestAMonthlyKeyLimitWithinTheCeilingIsAccepted TestAnExhaustedKeyLimitIsRefused; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0
- `tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/runtime.md | grep -qF -- 'jev_router_key_unbounded' && cmp .agents/skills/roundfix/references/runtime.md skills/roundfix/references/runtime.md` — expected: exit 0

## References

- Measurement record → `measurement/jev-router.md`
- `_techspec.md` → The month's Jev spend; The gate and the record
