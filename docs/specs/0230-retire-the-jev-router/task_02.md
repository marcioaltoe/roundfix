---
task: task_02
spec: 0230-retire-the-jev-router
status: pending
type: backend
complexity: medium
---

# Task 02: Configuration refuses the models the subscription rule reserves, and the router's floor key degrades to a warning

## Overview

After task_01 the router is gone from the runtime, but configuration still
carries its rules: only Project Config may select the router, its effort
must be empty, and the User Config floor `jev.router_min_credit_usd` is
validated and stored. This Task replaces those rules with the subscription
rule in every scope that names an Agent Selection, so a configuration that
would reach an OpenAI or Anthropic model through OpenRouter fails to load and
`roundfix profiles configure` refuses it, and turns the floor into a
deprecated key that loads with one warning. It is verifiable on its own
through configuration loading and the configure command.

## Requirements

1. MUST call `CheckSubscriptionRule` from selection normalization with field
   `<path>.model`, after the model is known non-empty and before any
   reasoning-effort rule (TechSpec Invariant 6), so Project Config, User
   Config, the legacy `runtimes` defaults once they become profiles, a
   one-Run override and `roundfix profiles configure` all refuse with API
   Contract 1 or 2 (API Contract 3).
2. MUST remove the router's configuration rules: the Project-Config-only
   source check, the one-Run override refusal, the router's empty-effort
   rule and the legacy-runtimes router check; and MUST delete the config
   package's router file and its two router test files.
3. MUST register `jev.router_min_credit_usd` as a deprecated key with no
   replacement, removed before validation in User Config and Project Config
   whatever its value, with the one warning of API Contract 5; MUST remove
   the Project Config ignore warning and the value validation for that key,
   the `RouterMinCreditUSD` field and its overlay; and MUST keep the existing
   deprecated-key warnings byte-identical.
4. MUST keep `jev.monthly_ceiling_usd` behavior unchanged: User Config only,
   the Project Config warning, and its validation.
5. MUST add tests: through `Load` with disposable Homes and repositories,
   each scope refuses the four models of Success Metric 1 with the field, the
   model and the rule in the message, while
   `openrouter/deepseek/deepseek-v4-pro` on `opencode` and `gpt-6.1-sol` on
   `codex` load, and `ResolveProfile` refuses a refused override; through
   `Load`, `jev.router_min_credit_usd` values `20` and `0` in User Config and
   in Project Config each load with exactly one API Contract 5 warning; and
   through the CLI harness, `roundfix profiles configure --scope project`
   with a refused fragment exits 2 with Surface Transcript 1's stderr and
   leaves Project Config byte-identical.
6. MUST NOT change the runner, the daemon, the judge, any guide or skill, or
   any built-in or Recommended Profile; no test may reach the network or the
   real `~/.roundfix`.

## Subtasks

- [ ] Call the rule from selection normalization.
- [ ] Remove the router's configuration rules and files.
- [ ] Deprecate the floor key and remove its field.
- [ ] Add the scope, deprecated-key and configure tests.

## Acceptance Criteria

- [ ] Every configuration scope refuses each refused model with the field,
      the model and the rule, and still accepts the allowed selections.
- [ ] `roundfix profiles configure` refuses a refused fragment, exits 2 with
      Surface Transcript 1's stderr and writes nothing.
- [ ] A configuration carrying `jev.router_min_credit_usd`, with any value,
      loads with exactly one API Contract 5 warning.
- [ ] No source in the repository names the router's configuration rules or
      the floor field any more.

## Context

- instruction: `docs/adr/0027-removed-config-keys-degrade-to-warnings.md`
- deletes: `internal/config/jev_router.go`
- deletes: `internal/config/jev_router_test.go`
- deletes: `internal/config/jev_router_credit_test.go`
- interface: `internal/config/profiles.go`
- interface: `internal/config/config.go`
- creates: `internal/config/subscription_rule_config_test.go`
- creates: `internal/cli/profiles_configure_subscription_test.go`

## Verification

- `test ! -e internal/config/jev_router.go && test ! -e internal/config/jev_router_test.go && test ! -e internal/config/jev_router_credit_test.go && ! grep -rqE 'JevRouter|RouterMinCreditUSD|router_min_credit_usd must be|which only Project Config may select|cannot be a one-Run override' internal cmd && go build -buildvcs=false ./... && go vet ./internal/config ./internal/cli` — expected: exit 0; before this Task the config router files, the floor field and the router rules exist, so the command fails.
- `out="$(go test -count=1 -v -run '^(TestSubscriptionRuleRefusesEveryConfigScope|TestRetiredRouterCreditFloorIsDeprecated)$' ./internal/config 2>&1)" || { printf '%s\n' "$out"; exit 1; }; cli="$(go test -count=1 -v -run '^TestProfilesConfigureRefusesSubscriptionModels$' ./internal/cli 2>&1)" || { printf '%s\n' "$cli"; exit 1; }; for name in TestSubscriptionRuleRefusesEveryConfigScope TestRetiredRouterCreditFloorIsDeprecated TestProfilesConfigureRefusesSubscriptionModels; do printf '%s\n%s\n' "$out" "$cli" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the three tests exists, so the command fails.

## References

- `_prd.md` → Goals; User Stories 1, 2 and 3; Core Features 2 and 4; Success Metrics 1 and 4
- `_techspec.md` → Interfaces; API Contract 1; API Contract 2; API Contract 3; API Contract 5; Surface Transcript 1; Surface Transcript 2; Testing Approach; Build Order 3
- ADR-0235; ADR-0027; ADR-0049; ADR-0231
