---
task: task_02
spec: 0218-the-jev-router-on-docs-and-chore-tasks
status: completed
type: backend
complexity: medium
---

# Task 02: The month's Jev spend counts the judge, the router and the key's own usage

## Overview

Spec 0205's judge stops at US$5 of Jev spend a month, counted from its Judge
Log. The router spends on the same OpenRouter key, so one ceiling must count
both, and must not trust local lines alone. This Task adds
`internal/jevrouter`: a reader of the key's `usage_monthly`, the month's
spend from the Judge Log and that usage, and a writer of the
`router-prompt` Judge Log line (ADR-0218). It requires Spec 0205's judge
package in the tree.

## Requirements

1. MUST add `KeyUsage` in `internal/jevrouter/key.go`: one
   `GET <endpoint>/key` with `Authorization: Bearer <key>` and no other
   header, 10 seconds, no retry, returning `data.usage_monthly` and an error
   for any non-`200`, a timeout, a network failure, or a body without a
   non-negative number there. The error MUST NOT contain the key.
2. MUST add `MonthSpend` and `Deps` in `internal/jevrouter/spend.go` as the
   TechSpec "The month's Jev spend" states: Judge Log lines of the UTC month
   read through the judge package's Judge Log reader, `cost_usd` summed by
   `transport`, and `Total = TypeSafeLogged + max(OpenRouterLogged, KeyUsageMonthly)`.
   The ceiling MUST come from the judge package's loaded questions, never a
   second literal.
3. MUST add `Ledger.Append` in `internal/jevrouter/ledger.go`, writing the
   `router-prompt` line with exactly the fields and values of the TechSpec
   "The gate and the record", through the judge package's Judge Log writer or
   file mode (`0700` directory, `0600` file, append). If Spec 0205 left that
   reader or writer unexported, the Task MAY export it; that file is then a
   recorded path.
4. MUST add `TestRouterLineIsReadByTheJudgeLog`, which appends a line and
   reads the month back through the judge package, counting its cost.
5. MUST test the key reader against an `httptest.Server` only, in a temporary
   Home, and MUST NOT read a real `ROUNDFIX_OPENROUTER_API_KEY` or
   `OPENROUTER_API_KEY`.

## Subtasks

- [ ] Add the key reader.
- [ ] Add the month's spend over the Judge Log and the key.
- [ ] Add the router line and prove the judge reads it.

## Acceptance Criteria

- [ ] With Judge Log lines of US$0.10 TypeSafe and US$0.20 OpenRouter and a
      key usage of US$0.50, the spend is US$0.60; with a key usage of
      US$0.05 it is US$0.30.
- [ ] A failing, non-`200` or malformed key endpoint is an error without the
      key in it.
- [ ] A router line is read by the judge package and counted.

## Context

- creates: `internal/jevrouter/key.go`
- creates: `internal/jevrouter/spend.go`
- creates: `internal/jevrouter/ledger.go`
- creates: `internal/jevrouter/key_test.go`
- creates: `internal/jevrouter/spend_test.go`
- creates: `internal/jevrouter/ledger_test.go`
- instruction: `docs/history/specs/0205-an-advisory-judge-for-spec-authoring/_techspec.md`
- instruction: `docs/adr/0201-the-judge-sends-only-spec-artifacts-and-spends-under-a-monthly-ceiling.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestKeyUsageSendsOnlyTheBearerHeader|TestKeyUsageRefusesAMalformedAnswer|TestMonthSpendAddsTypeSafeToTheLargerOpenRouterFigure|TestMonthSpendFailsOnAnUnreadableKeyEndpoint|TestRouterLineIsReadByTheJudgeLog)$' ./internal/jevrouter 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestKeyUsageSendsOnlyTheBearerHeader TestKeyUsageRefusesAMalformedAnswer TestMonthSpendAddsTypeSafeToTheLargerOpenRouterFigure TestMonthSpendFailsOnAnUnreadableKeyEndpoint TestRouterLineIsReadByTheJudgeLog; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task `internal/jevrouter` does not exist, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 3; User Story 3; Core Features 4-5; Success Metric 4; Prerequisites
- [_techspec.md](_techspec.md) — Interfaces; The month's Jev spend; The gate and the record; API Contracts 2 and 4; Testing Approach 2; Build Order 2
- ADR-0218; ADR-0201

## Result

Implemented the Task 02 slice for Daemon Verification:

- `KeyUsage` sends one bearer-only GET, disables implicit User-Agent and
  compression headers, cookies and redirects, honors cancellation, and caps
  the request at ten seconds. Non-200 responses, invalid or missing usage,
  body-read failures and network errors return errors with the supplied key
  redacted; wrapped causes remain matchable.
- `MonthSpend` reads the UTC month's rows through `judge.ReadMonth`, sums
  TypeSafe and OpenRouter costs separately, and adds TypeSafe to the larger
  OpenRouter figure. A missing log is empty; malformed logs and key endpoint
  failures are errors. The default ceiling comes from `judge.Load`; an
  explicit `Deps.Ceiling` is reserved for the caller's loaded judge settings.
  Dependencies carry an explicit `Env`, `HomeDir`, `Client` and `Endpoint`, so
  tests never read either real OpenRouter key variable.
- `Ledger.Append` uses the judge's schema and append-only writer, preserving
  its 0700 directory and 0600 file modes. It records the documented identity,
  tokens, latency, cost delta and outcomes, clamps negative deltas, and records
  unreadable post-prompt usage as zero cost with an error.
- Additional recorded path: `internal/judge/log.go`. Exported `LogLine`,
  `ReadMonth` and `AppendLogLine` share the existing log boundary; the judge's
  private monthly-cost reader now uses the same parser as the router.

Focused-check evidence, 2026-10-02:

- `GOCACHE=/private/tmp/roundfix-0218-task02-gocache rtk proxy go test -count=1 -v ./internal/jevrouter ./internal/judge`
  exited 0 with loopback-server permission. Router package: 10.360s; judge
  package: 0.512s. The judge's opt-in live measurement remained skipped.
- Spend criterion: `TestMonthSpendAddsTypeSafeToTheLargerOpenRouterFigure`
  passed for usage 0.50, 0.05 and zero, demonstrating totals 0.60, 0.30 and
  0.30, excluding the prior month and using the loaded judge ceiling.
- Endpoint criterion: `TestKeyUsageRefusesAMalformedAnswer`,
  `TestKeyUsageRefusesHTTPFailuresAndRedirects`,
  `TestKeyUsageHonorsCancellationAndNetworkFailure`,
  `TestKeyUsageCapsTheRequestAtTenSeconds` and
  `TestMonthSpendFailsOnAnUnreadableKeyEndpoint` passed. The network test
  includes the sentinel in the failed URL to exercise redaction. The stalled
  endpoint stopped at the reader's own ten-second deadline.
- Log criterion: `TestRouterLineIsReadByTheJudgeLog` passed, checking all
  record values, null judgment fields, private file modes and two appended
  rows whose costs sum to 1.00 through the shared judge reader.
  `TestLedgerRecordsUnreadableUsageAndClampsNegativeDelta` also passed.
- `TestKeyUsageSendsOnlyTheBearerHeader` and
  `TestMonthSpendRejectsUnreadableLogsAndMissingKeys` passed against local
  servers and temporary Homes only. No live API usage occurred.
- The initial checks hit sandbox restrictions on the shared Go cache and
  loopback listeners; a task-local cache and loopback permission resolved
  those environment restrictions.

The live [OpenRouter key reference](https://openrouter.ai/docs/api/api-reference/api-keys/get-current-api-key)
was checked for the request and response contract. Before implementation,
`internal/jevrouter` was absent. Task status, authored Verification, the Task
Graph and other Task files were left unchanged. No commit, push or pull
request was made; the Daemon owns Verification and settlement.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/judge/log.go`

## Carry-forward provenance

- Source Run: `run_20261002T194325Z_56cb22f8a792bdff`
- Source commit: `4ed4b70394c344d34b265b8ed232fbc5a351b992`
