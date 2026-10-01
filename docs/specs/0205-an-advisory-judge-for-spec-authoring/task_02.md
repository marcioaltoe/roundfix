---
task: task_02
spec: 0205-an-advisory-judge-for-spec-authoring
status: pending
type: backend
complexity: high
---

# Task 02: The judge asks, records every call and stops at the monthly ceiling

## Overview

task_01 plans a Spec's judgments and nothing asks them. This Task adds the
transport selection, OpenRouter's System One API first and the direct
TypeSafe endpoint second, the request, the model pin, the monthly Judge Log in
Roundfix Home, the spending ceiling read from it, and `Run`, which asks a plan
in order, compares each answer from the pinned model with its threshold and
fails open on every error (ADR-0200, ADR-0201). It is verifiable on its own
through an injected `http.RoundTripper`: a fixture Spec yields the advisory,
clear and skipped outcomes the TechSpec names on both transports, one Judge
Log line per request, and no request at all when no Jev key is set or the
ceiling is reached.

## Requirements

1. MUST implement `selectTransport` and the transport rule of `_techspec.md` → Invariants, Invariant 2: `ROUNDFIX_JEV_OPENROUTER_API_KEY` non-empty selects `openrouter`, else `ROUNDFIX_TYPESAFE_API_KEY` non-empty selects `typesafe`, else the run is skipped with `ROUNDFIX_JEV_OPENROUTER_API_KEY is not set (nor ROUNDFIX_TYPESAFE_API_KEY)`; the generic `OPENROUTER_API_KEY` and `TYPESAFE_API_KEY` are never read (`TestRunIgnoresTheGenericOpenRouterKey` MUST cover an environment holding only `OPENROUTER_API_KEY`, only `TYPESAFE_API_KEY`, and both: each is a skip that sends nothing); a key is sent only to its own transport's endpoint; and the transport is chosen once per run. Then MUST implement the request of `_techspec.md` → Asking step 3 and API Contract 3: `POST` to the selected transport's `endpoint` from the question file, the headers `Authorization: Bearer <key>` and `Content-Type: application/json` and no other header the code sets, the body `{"state", "model", "questions"}` with the selected transport's `request_model` (`jev-1.13` on OpenRouter, `jev-1.13.0` on TypeSafe) and one question under its `question_id`, 30 seconds per attempt, and at most two retries of a `429` or `529` after `Retry-After` (capped at 10 seconds) or 1.5 and 3 seconds. Every wait MUST end when the context is cancelled.
2. MUST read the answer as API Contract 3 states, including OpenRouter's optional `id`, `provider` and `usage.cost`, and apply `_techspec.md` → Asking step 4 and Invariant 1 through `Questions.pinned`: an answer whose reported model does not match `accepted_model_pattern` (`jev-1.13.<patch>`, or `typesafe/jev-1.13` with an optional `-YYYYMMDD` suffix) is skipped with its reason, logged, and never raised or cleared; a malformed answer is skipped as `unreadable answer`; otherwise the thresholds of the question file decide `advisory` or `clear`, compared exactly as written (`>=` for the confidence floor, `<` for the probability ceiling).
3. MUST implement every skip and stop reason of `_techspec.md` → Asking steps 1, 2, 5 and 6 with the exact reason texts, including `402` among the key refusals, and send no request after a stop, when no Jev key is set, when the Judge Log is unreadable, or when the month's cost plus the run's cost is at or above `monthly_ceiling_usd`.
4. MUST implement the Judge Log of `_techspec.md` → Data Models and API Contract 4: `<HomeDir>/.roundfix/judge/<YYYY-MM>.jsonl` for the UTC month of `Request.Now`, directory `0700`, file `0600`, one appended line per request with every field of the example line, including `transport`, `response_id`, `provider` and `cost_source`, the cost of `_techspec.md` → Invariant 3 (the reported `usage.cost`, else input tokens at the question file's price), and the month's cost summed from the current month's file across both transports. No key MUST appear in any line.
5. MUST implement `Run` as `_techspec.md` → Interfaces sketches: it asks the plan task_01 builds, asks identical states once, fills the `Report` fields that section names, and returns an error only when the Spec's PRD cannot be read.
6. MUST put the tests in `internal/judge/client_test.go`, `internal/judge/log_test.go` and `internal/judge/judge_test.go`, with a fake `http.RoundTripper`, a fixed clock and `t.TempDir()` as Roundfix Home. `TestRequestsCarryOnlySpecArtifactText` MUST follow `_techspec.md` → Testing Approach 3: its fixture repository also holds a Go source file and a findings file with sentinel strings, and it requires, once per transport, that every state string, split at `the cited decision`, occurs in the PRD, the TechSpec or an ADR, that no sentinel and no key appears in any request body, and that no header beyond the two named is set. `TestAskSendsEachTransportsOwnModelID` MUST require that a request on the `openrouter` transport carries `"model":"jev-1.13"` and one on the `typesafe` transport carries `"model":"jev-1.13.0"`. `TestRunSelectsTheTransportByKey` MUST cover each key alone and both keys together (every request goes to the OpenRouter endpoint with the OpenRouter key, and the TypeSafe key appears in no request); `TestRunIgnoresTheGenericOpenRouterKey` MUST set only `OPENROUTER_API_KEY` and require the named skip and no request; `TestModelPinAcceptsOnlyJev113` MUST accept `jev-1.13.0`, `typesafe/jev-1.13` and `typesafe/jev-1.13-20260917` and refuse `jev-1.14.0`, `typesafe/jev-1.14-20261101`, `~typesafe/jev-latest` and `""`; `TestRunRecordsTheReportedCost` MUST use an OpenRouter fixture response that answers the requested `jev-1.13` as `typesafe/jev-1.13-20260917` with `id`, `provider` and `usage.cost`, and require that its Judge Log line holds that cost as `reported`, while a direct fixture without `usage.cost` is logged as `computed` and the ceiling sums both.
7. MUST NOT open a real network connection, read the process's `ROUNDFIX_JEV_OPENROUTER_API_KEY`, `ROUNDFIX_TYPESAFE_API_KEY` or `OPENROUTER_API_KEY`, or write outside `t.TempDir()`.
8. MUST prove each new gate can fail. The Result MUST record one sabotage of the ceiling (for example `>` instead of `>=`), one of the model pin (for example accepting any `typesafe/` model) and one of the transport selection (for example falling back to `OPENROUTER_API_KEY`), each with the test that failed, and that the code was restored.

## Subtasks

- [ ] Select the transport from the keys.
- [ ] Send one request and read its answer.
- [ ] Record each call in the Judge Log and sum the month.
- [ ] Ask a plan with thresholds, skips and stops.
- [ ] Prove the data boundary on captured requests.
- [ ] Record one sabotage per gate in the Result.

## Acceptance Criteria

- [ ] A `says_nothing` answer at confidence 0.8 is `advisory`, at 0.79 `clear`; a `noul` of 0.29 is `advisory`, 0.3 `clear`.
- [ ] An answer reported as `jev-1.13.0` or `typesafe/jev-1.13-20260917` is compared with its threshold; one reported as `jev-1.14.0` is skipped as `answered by jev-1.14.0, thresholds belong to jev-1.13`, and logged.
- [ ] The OpenRouter request names `jev-1.13` and the TypeSafe request names `jev-1.13.0`.
- [ ] With both keys set, every request goes to `https://openrouter.ai/api/v1/systemone` with `ROUNDFIX_JEV_OPENROUTER_API_KEY`; with only `ROUNDFIX_TYPESAFE_API_KEY`, to `https://api.typesafe.ai/v1/systemone`; with only `OPENROUTER_API_KEY`, nowhere.
- [ ] An OpenRouter answer's `usage.cost` is the logged cost (`reported`); a direct answer's cost is input tokens at the pinned price (`computed`).
- [ ] With no Jev key, an unreadable Judge Log, or a month already at US$5.00, the fake transport receives no request.
- [ ] A `503` stops the run with `service unavailable (HTTP 503)` and no later request; a `422` skips one judgment and the next is asked; two `429` answers with `Retry-After: 0` followed by a `200` give one answer after three attempts.
- [ ] Each request leaves one Judge Log line with every field, the file is `0600`, and no line holds the key.

## Context

- creates: `internal/judge/client.go`
- creates: `internal/judge/log.go`
- creates: `internal/judge/judge.go`
- creates: `internal/judge/client_test.go`
- creates: `internal/judge/log_test.go`
- creates: `internal/judge/judge_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestAskSendsTheMeasuredRequest|TestAskRetriesARateLimitThenAnswers|TestRunRaisesAtTheMeasuredThresholds|TestRunNeverComparesAnotherModelsAnswer|TestRunSendsNothingWithoutAKey|TestRunStopsAtTheMonthlyCeiling|TestRunSendsNothingWithAnUnreadableLog|TestRunStopsWhenTheServiceFails|TestRunSkipsARefusedRequestAndContinues|TestJudgeLogRecordsEveryCall|TestRequestsCarryOnlySpecArtifactText|TestAskSendsEachTransportsOwnModelID|TestRunSelectsTheTransportByKey|TestRunIgnoresTheGenericOpenRouterKey|TestModelPinAcceptsOnlyJev113|TestRunRecordsTheReportedCost)$" ./internal/judge 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestAskSendsTheMeasuredRequest TestAskRetriesARateLimitThenAnswers TestRunRaisesAtTheMeasuredThresholds TestRunNeverComparesAnotherModelsAnswer TestRunSendsNothingWithoutAKey TestRunStopsAtTheMonthlyCeiling TestRunSendsNothingWithAnUnreadableLog TestRunStopsWhenTheServiceFails TestRunSkipsARefusedRequestAndContinues TestJudgeLogRecordsEveryCall TestRequestsCarryOnlySpecArtifactText TestAskSendsEachTransportsOwnModelID TestRunSelectsTheTransportByKey TestRunIgnoresTheGenericOpenRouterKey TestModelPinAcceptsOnlyJev113 TestRunRecordsTheReportedCost; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the sixteen tests exists, so the command fails.

## References

- [_prd.md](_prd.md) — Goals 2, 3 and 4; User Stories 3, 4, 5 and 6; Core Features 5, 6, 8, 9 and 10; Success Metrics 2, 3, 4 and 7; Recorded limits
- [_techspec.md](_techspec.md) — Interfaces; Invariants; Asking; Data Models; API Contract 3; API Contract 4; Testing Approach 3; Testing Approach 4; Build Order 2
- ADR-0200; ADR-0201; ADR-0089
