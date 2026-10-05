---
task: task_02
spec: 0229-a-jev-router-that-checks-credit-and-names-its-model
status: pending
type: backend
complexity: medium
---

# Task 02: The Jev Router gate refuses below the account credit floor

## Overview

On 2026-10-04 (the Jev Router measurement addendum of that date) the gate
sent a routed repair prompt that the OpenRouter account, holding about
US$5.11, could not cover, and OpenRouter cut it off. This Task adds the User
Config value `jev.router_min_credit_usd` (US$15 when unset), the credit read
`GET /api/v1/credits`, and the gate's refusal `openrouter_credit_low` when the
lower of the account balance and the key's remaining limit is below the
floor, as the TechSpec's "The configured floor" and "The credit read and the
gate" state (ADR-0234). Before Agent work the refusal hands the Task to the
Fallback Chain; after it the Work Item fails, as for `jev_ceiling_reached`.

## Requirements

1. MUST add `RouterMinCreditUSD` to `config.Jev` and `router_min_credit_usd`
   to `jevOverlay`, decode it from User Config, and leave it zero in
   `Builtin()` and when unset; MUST remove it from a Project Config document
   before decoding and print API Contract 2 through the existing
   ignored-setting warning, as `jev.monthly_ceiling_usd` is handled; and MUST
   refuse a User Config value that is not a finite number greater than zero
   with API Contract 3 inside the existing `parse config "<path>": ` wrapping.
2. MUST add, in a new file of the `jevrouter` package, `DefaultMinCreditUSD`
   (15), `Credits` with `Balance`, `ReadCredits`, `MinCredit` and `CreditLeft`
   with the TechSpec's Interfaces. `ReadCredits` MUST follow `ReadKey`'s
   transport rules (ten-second bound, no redirect, no cookie jar, no
   `User-Agent`, at most 1 MiB read) and MUST return an error, with the key
   redacted, for a non-200 status, a malformed body, or a missing, negative
   or non-finite `data.total_credits` or `data.total_usage`.
3. MUST add `MinCreditUSD` to `jevrouter.Deps`, and make
   `jevRouterGate.Before`, after its existing spend, key-limit and ceiling
   checks, read the credits and return a `SelectionFailureError` for runtime
   `opencode` with `jev_spend_unreadable: <error>` when the read fails and
   with API Contract 4 when `CreditLeft` is below `MinCredit(deps.MinCreditUSD)`.
   The prompt is not sent in either case. MUST NOT change `MonthSpend` or
   `Spend.CheckKeyLimit`.
4. MUST add `JevRouterMinCreditUSD` to `daemon.Dependencies`, make `NewEngine`
   pass it to the default gate's `MinCreditUSD`, and set it from the loaded
   configuration in the Implement Run engine in `internal/cli/implement.go`
   and the resolve engine in `internal/cli/cli.go`, beside
   `JevMonthlyCeilingUSD`.
5. MUST make `selectionReasonCode` return `openrouter_credit_low` for that
   refusal, as it returns `jev_ceiling_reached`, so the fallback receipt
   names it.
6. MUST add the tests named in Verification as Testing Approach 1 to 3
   describe: `internal/config/jev_router_credit_test.go` (new) for a User
   Config floor of 20, the unset default, the Project Config warning and the
   refusals of 0, -1 and `.inf`; `internal/jevrouter/credits_test.go` (new) for
   the bearer header, every malformed answer, HTTP failure and redirect with
   no key in the error text, `CreditLeft` and `MinCredit`; and in
   `internal/daemon/jev_router_gate_test.go` a local stand-in serving `/key`
   and `/credits` that proves the balance 5.11 with a key remaining limit of
   40.31 is refused naming US$15, a balance of 20 is sent, a key remaining
   limit of 3 with a balance of 200 is refused, a 403 credits answer is
   `jev_spend_unreadable`, a gate built with `JevRouterMinCreditUSD` 4 sends
   the 5.11 balance, and a fake-gate `openrouter_credit_low` refusal falls back
   before Agent work and fails the Task after it.
7. MUST change `TestJevRouterGateChecksTheReportedKeyLimit` and
   `TestJevRouterGateAcceptsAKeyLimitAtTheConfiguredCeiling`, whose local
   servers answer every path with the key's body, so that they also serve a
   valid `/credits` answer and keep asserting what they asserted.
8. MUST NOT change `.roundfixrc.yml`, the `roundfix config init` templates,
   any other configuration key, the router's selection rules or any default
   or Recommended Profile, or write any file outside the repository; every
   test uses a disposable home and a local endpoint, and no test reaches the
   network.

## Subtasks

- [ ] Add the configuration value, its Project Config removal and its validation.
- [ ] Add the credit read, the floor default and the credit-left rule.
- [ ] Refuse below the floor in the gate and pass the configured floor from both engines.
- [ ] Name the reason code and add the tests.

## Acceptance Criteria

- [ ] `jev.router_min_credit_usd: 20` in User Config loads as 20; unset loads
      as zero; a Project Config value is ignored with API Contract 2; 0, -1
      and `.inf` fail with API Contract 3.
- [ ] With the account at 5.11 and the key at 40.31 remaining, the gate
      refuses with `openrouter_credit_low: OpenRouter credit US$5.1100 is below the US$15.0000 floor`
      and sends no prompt; at 20 it sends.
- [ ] Before Agent work the refusal falls back with a receipt naming
      `openrouter_credit_low`; after it the Work Item fails with that reason.

## Context

- interface: `internal/config/config.go`
- creates: `internal/config/jev_router_credit_test.go`
- creates: `internal/jevrouter/credits.go`
- creates: `internal/jevrouter/credits_test.go`
- interface: `internal/jevrouter/spend.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/daemon/agent_session_owner.go`
- interface: `internal/daemon/jev_router_gate_test.go`
- interface: `internal/cli/implement.go`
- interface: `internal/cli/cli.go`
- instruction: `docs/user-guide/configuration.md`

## Verification

- `grep -q JevRouterMinCreditUSD internal/cli/implement.go || { printf 'implement engine does not pass the floor\n' >&2; exit 1; }; grep -q JevRouterMinCreditUSD internal/cli/cli.go || { printf 'resolve engine does not pass the floor\n' >&2; exit 1; }; out="$(go test -count=1 -v -run "^(TestJevRouterMinCreditIsReadFromUserConfig|TestJevRouterMinCreditIsUnsetByDefault|TestProjectConfigCannotSetTheJevRouterMinCredit|TestJevRouterMinCreditRefusesANonPositiveOrInfiniteValue|TestReadCreditsSendsOnlyTheBearerHeader|TestReadCreditsRefusesMalformedAnswersAndHTTPFailures|TestCreditLeftTakesTheLowerOfBalanceAndKeyLimit|TestMinCreditDefaultsToFifteenDollars|TestJevRouterGateRefusesCreditBelowTheFloor|TestJevRouterGateRefusesUnreadableCredits|TestJevRouterDefaultGateUsesTheConfiguredCreditFloor|TestJevRouterCreditLowFallsBackBeforeWork|TestJevRouterCreditLowFailsTheTaskAfterWork|TestJevRouterGateChecksTheReportedKeyLimit|TestJevRouterGateAcceptsAKeyLimitAtTheConfiguredCeiling)$" ./internal/config ./internal/jevrouter ./internal/daemon 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestJevRouterMinCreditIsReadFromUserConfig TestJevRouterMinCreditIsUnsetByDefault TestProjectConfigCannotSetTheJevRouterMinCredit TestJevRouterMinCreditRefusesANonPositiveOrInfiniteValue TestReadCreditsSendsOnlyTheBearerHeader TestReadCreditsRefusesMalformedAnswersAndHTTPFailures TestCreditLeftTakesTheLowerOfBalanceAndKeyLimit TestMinCreditDefaultsToFifteenDollars TestJevRouterGateRefusesCreditBelowTheFloor TestJevRouterGateRefusesUnreadableCredits TestJevRouterDefaultGateUsesTheConfiguredCreditFloor TestJevRouterCreditLowFallsBackBeforeWork TestJevRouterCreditLowFailsTheTaskAfterWork TestJevRouterGateChecksTheReportedKeyLimit TestJevRouterGateAcceptsAKeyLimitAtTheConfiguredCeiling; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task neither engine passes the floor and the new tests do not exist, so the command fails.

## References

- `_prd.md` → User Stories 1, 2 and 5; Core Features 1-3; Success Metric 1; Success Metric 4
- `_techspec.md` → The configured floor; The credit read and the gate; Interfaces; API Contract 1; API Contract 2; API Contract 3; API Contract 4; Testing Approach 1-3; Build Order 2
- ADR-0234; ADR-0231; ADR-0218; ADR-0114; ADR-0050; ADR-0027
