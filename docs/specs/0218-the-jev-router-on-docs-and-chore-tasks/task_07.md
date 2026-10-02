---
task: task_07
spec: 0218-the-jev-router-on-docs-and-chore-tasks
status: completed
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

## Result

Implemented the monthly key-limit safeguard for Daemon Verification.

- The bounded key reader now returns nullable `limit`, `limit_remaining`,
  and `limit_reset` alongside monthly usage. `KeyUsage` retains its existing
  usage-only interface for cost accounting and its malformed-usage refusals.
- The production prompt gate checks the reported key limit before allowing
  execution. A missing, lifetime, negative, or over-ceiling limit produces
  `jev_router_key_unbounded` with the remedy to set a monthly credit limit of
  at most the loaded ceiling at OpenRouter. A monthly key with numeric
  `limit_remaining <= 0` produces `jev_ceiling_reached`.
- The existing Selection Failure path handles both refusals: fallback before
  work and Work Item failure after work. The reason mapper now preserves
  `jev_router_key_unbounded` in the fallback notification.
- The runtime reference and configuration guide explain the monthly limit,
  remedy, exhausted-key refusal, and fallback behavior. Both Roundfix Skill
  version fields are `0.1.17`; `make skills-sync` regenerated the shipped
  mirror and the owned-version recording command recorded its digest.

Acceptance evidence: the routed prompt cannot pass the production `Before`
gate unless the key reports a numeric limit at or below the loaded ceiling
and exactly `monthly` reset. The five required tests in
`internal/jevrouter/key_limit_test.go` exercise the local fake key endpoint.
`TestJevRouterGateChecksTheReportedKeyLimit` exercises the production gate
against that HTTP boundary; the two new unbounded-key lifecycle tests prove
fallback before work and failure after work, including reason preservation.

Focused checks (with `GOCACHE=/tmp/roundfix-task07-gocache` for Go commands):

- `go test -count=1 ./internal/jevrouter`: exit 0, including all existing key,
  spend, and ledger tests and the five required key-limit tests.
- `go test -count=1 ./internal/daemon -run 'TestJevRouter|TestNonRoutedPrompt'`:
  exit 0, including the new real-gate and lifecycle tests.
- A temporary Go overlay removed only the production `CheckKeyLimit` call.
  Running `TestJevRouterGateChecksTheReportedKeyLimit` against that overlay
  exited 1 on unlimited, lifetime, over-ceiling, and exhausted keys: each
  incorrectly returned usage with no refusal. The overlay lives under `/tmp`
  and changes no repository file.
- The same real-gate regression test without the overlay exited 0 on the
  current implementation.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
  -record-skill-versions`: exit 0 after version `0.1.17` was synchronized.
- `make skills-sync`: exit 0. A Python byte comparison confirmed both skill
  files and runtime references match their shipped mirrors; phrase checks
  confirmed monthly reset, remaining limit, and the new reason in both docs.
- `make baseline-digests`: exit 0; no additional derived changes required.
- `git diff --check`: exit 0.
- First `make verify-incremental`: exit 2. Assertions passed, but the
  repository suite guard correctly refused edits made during its execution
  to this Result and `internal/daemon/jev_router_gate_test.go`. A subsequent
  run held the worktree unchanged while the check executed and exited 0:
  formatting, vet, package tests, skill sync/check, and build passed.

The sandbox initially refused local HTTP listeners; the focused tests were
rerun with the required permission scope and passed. Skill-source edits used
the same scope because `.agents` is protected, under the Spec's existing
explicit authorization. No live key was read and no live API call was made.
The live [OpenRouter key endpoint reference](https://openrouter.ai/docs/api/api-reference/api-keys/get-current-api-key)
confirmed the reported limit fields; the [TypeSafe documentation index](https://docs.typesafe.ai/llms.txt)
was consulted as required by the TypeSafe skill.

Scope postflight: the production gate wiring and reason preservation also
change `internal/daemon/engine.go`, `internal/daemon/agent_session_owner.go`,
and `internal/daemon/jev_router_gate_test.go`. These ordinary paths are
necessary to enforce and test this Task's safeguard through the existing
execution path. All governed edits stay within the Spec authorization; no
additional derived files remain changed.

Task status, checkboxes, and declared Verification commands remain
Daemon-owned. No commit, push, Pull Request, Task Graph edit, or other Task
edit was made.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/daemon/agent_session_owner.go`
- `internal/daemon/engine.go`
- `internal/daemon/jev_router_gate_test.go`

## Carry-forward provenance

- Source Run: `run_20261002T213602Z_376bbbe5b3ded83d`
- Source commit: `261a2164f3a271f7c67ed443c1e172632adf094a`
