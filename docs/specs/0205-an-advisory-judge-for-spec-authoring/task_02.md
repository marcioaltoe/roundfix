---
task: task_02
spec: 0205-an-advisory-judge-for-spec-authoring
status: completed
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

1. MUST implement `selectTransport` and the transport rule of `_techspec.md` → Invariants, Invariant 2: `ROUNDFIX_OPENROUTER_API_KEY` non-empty selects `openrouter`, else `ROUNDFIX_TYPESAFE_API_KEY` non-empty selects `typesafe`, else the run is skipped with `ROUNDFIX_OPENROUTER_API_KEY is not set (nor ROUNDFIX_TYPESAFE_API_KEY)`; the generic `OPENROUTER_API_KEY` and `TYPESAFE_API_KEY` are never read (`TestRunIgnoresTheGenericOpenRouterKey` MUST cover an environment holding only `OPENROUTER_API_KEY`, only `TYPESAFE_API_KEY`, and both: each is a skip that sends nothing); a key is sent only to its own transport's endpoint; and the transport is chosen once per run. Then MUST implement the request of `_techspec.md` → Asking step 3 and API Contract 3: `POST` to the selected transport's `endpoint` from the question file, the headers `Authorization: Bearer <key>` and `Content-Type: application/json` and no other header the code sets, the body `{"state", "model", "questions"}` with the selected transport's `request_model` (`jev-1.13` on OpenRouter, `jev-1.13.0` on TypeSafe) and one question under its `question_id`, 30 seconds per attempt, and at most two retries of a `429` or `529` after `Retry-After` (capped at 10 seconds) or 1.5 and 3 seconds. Every wait MUST end when the context is cancelled.
2. MUST read the answer as API Contract 3 states, including OpenRouter's optional `id`, `provider` and `usage.cost`, and apply `_techspec.md` → Asking step 4 and Invariant 1 through `Questions.pinned`: an answer whose reported model does not match `accepted_model_pattern` (`jev-1.13.<patch>`, or `typesafe/jev-1.13` with an optional `-YYYYMMDD` suffix) is skipped with its reason, logged, and never raised or cleared; a malformed answer is skipped as `unreadable answer`; otherwise the thresholds of the question file decide `advisory` or `clear`, compared exactly as written (`>=` for the confidence floor, `<` for the probability ceiling).
3. MUST implement every skip and stop reason of `_techspec.md` → Asking steps 1, 2, 5 and 6 with the exact reason texts, including `402` among the key refusals, and send no request after a stop, when no Jev key is set, when the Judge Log is unreadable, or when the month's cost plus the run's cost is at or above `monthly_ceiling_usd`.
4. MUST implement the Judge Log of `_techspec.md` → Data Models and API Contract 4: `<HomeDir>/.roundfix/judge/<YYYY-MM>.jsonl` for the UTC month of `Request.Now`, directory `0700`, file `0600`, one appended line per request with every field of the example line, including `transport`, `response_id`, `provider` and `cost_source`, the cost of `_techspec.md` → Invariant 3 (the reported `usage.cost`, else input tokens at the question file's price), and the month's cost summed from the current month's file across both transports. No key MUST appear in any line.
5. MUST implement `Run` as `_techspec.md` → Interfaces sketches: it asks the plan task_01 builds, asks identical states once, fills the `Report` fields that section names, and returns an error only when the Spec's PRD cannot be read.
6. MUST put the tests in `internal/judge/client_test.go`, `internal/judge/log_test.go` and `internal/judge/judge_test.go`, with a fake `http.RoundTripper`, a fixed clock and `t.TempDir()` as Roundfix Home. `TestRequestsCarryOnlySpecArtifactText` MUST follow `_techspec.md` → Testing Approach 3: its fixture repository also holds a Go source file and a findings file with sentinel strings, and it requires, once per transport, that every state string, split at `the cited decision`, occurs in the PRD, the TechSpec or an ADR, that no sentinel and no key appears in any request body, and that no header beyond the two named is set. `TestAskSendsEachTransportsOwnModelID` MUST require that a request on the `openrouter` transport carries `"model":"jev-1.13"` and one on the `typesafe` transport carries `"model":"jev-1.13.0"`. `TestRunSelectsTheTransportByKey` MUST cover each key alone and both keys together (every request goes to the OpenRouter endpoint with the OpenRouter key, and the TypeSafe key appears in no request); `TestRunIgnoresTheGenericOpenRouterKey` MUST set only `OPENROUTER_API_KEY` and require the named skip and no request; `TestModelPinAcceptsOnlyJev113` MUST accept `jev-1.13.0`, `typesafe/jev-1.13` and `typesafe/jev-1.13-20260917` and refuse `jev-1.14.0`, `typesafe/jev-1.14-20261101`, `~typesafe/jev-latest` and `""`; `TestRunRecordsTheReportedCost` MUST use an OpenRouter fixture response that answers the requested `jev-1.13` as `typesafe/jev-1.13-20260917` with `id`, `provider` and `usage.cost`, and require that its Judge Log line holds that cost as `reported`, while a direct fixture without `usage.cost` is logged as `computed` and the ceiling sums both.
7. MUST NOT open a real network connection, read the process's `ROUNDFIX_OPENROUTER_API_KEY`, `ROUNDFIX_TYPESAFE_API_KEY` or `OPENROUTER_API_KEY`, or write outside `t.TempDir()`.
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
- [ ] With both keys set, every request goes to `https://openrouter.ai/api/v1/systemone` with `ROUNDFIX_OPENROUTER_API_KEY`; with only `ROUNDFIX_TYPESAFE_API_KEY`, to `https://api.typesafe.ai/v1/systemone`; with only `OPENROUTER_API_KEY`, nowhere.
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

## Result

Implemented task_02's transport client, Judge Log, and fail-open `Run` in the
six declared Go files. Transport selection uses only the injected Roundfix
keys, once per run; redirects cannot carry a key to another endpoint. The
question file supplies each endpoint, request model, model pin, threshold,
price, and ceiling. Each HTTP attempt gets a 30-second context, is recorded
before any retry or later judgment, and counts toward calls and spend. The
current UTC month's log is reread before each attempt, including retries;
therefore recorded spend from either transport and another caller is included.
The report retains planning skips, unasked judgments, answers, and stop reasons.
No CLI, skill, tooling, Task Graph, or other Task changes belong to this diff.

Starting evidence: `client.go`, `log.go`, and `judge.go` did not exist (`ls`
exited 1). The initial worktree already had the Daemon's task_02 status edit;
that edit was preserved. No Task status or authored Verification was changed.

Focused checks on restored implementation:

- `GOCACHE=/private/tmp/roundfix-0205-task02-go-cache rtk proxy go test ./internal/judge -count=1 -run 'Test(Ask|Run|ModelPin|JudgeLog|Requests|Retry)'`
  — exit 0 after correcting null-field preservation in log redaction.
- `GOCACHE=/private/tmp/roundfix-0205-task02-go-cache rtk proxy go test ./internal/judge -count=1 -race -run 'Test(Ask|Run|ModelPin|JudgeLog|Requests|Retry)'`
  — exit 0 on the restored code and final tests, `ok roundfix/internal/judge`.
- `rtk proxy gofmt -l` on the six declared Go files — no output.
- `rtk proxy git -c core.fsmonitor=false diff --check` — exit 0.

Acceptance evidence from those focused checks:

| Criterion | Implemented behavior and test evidence |
| --- | --- |
| Exact confidence/probability boundaries | `TestRunRaisesAtTheMeasuredThresholds` observes advisory at confidence 0.8 and noul 0.29, clear at confidence 0.79 and noul 0.3. |
| Answering model pin | `TestRunNeverComparesAnotherModelsAnswer` compares both accepted model forms, skips Jev 1.14 with the exact required reason, and checks every logged model. `TestModelPinAcceptsOnlyJev113` covers all seven required accepted/refused IDs, including the empty model. |
| Endpoint-specific request models | `TestAskSendsEachTransportsOwnModelID` captures `jev-1.13` on OpenRouter and `jev-1.13.0` on TypeSafe. `TestAskSendsTheMeasuredRequest` checks the POST, endpoint, state, one question, two headers, and deadline. |
| Key selection and isolation | `TestRunSelectsTheTransportByKey` checks each dedicated key alone and both together, including absence of the direct key in OpenRouter requests. `TestRunIgnoresTheGenericOpenRouterKey` checks only the generic OpenRouter key, only the generic TypeSafe key, and both generic keys: the exact named skip and no request. `TestRunAsksIdenticalStatesOnceAndKeepsTheSelectedTransport` verifies deduplication and selection persistence. `TestRunNeverFollowsARedirectWithAKey` verifies no redirected request. |
| Reported/computed cost and combined ceiling | `TestRunRecordsTheReportedCost` captures OpenRouter's snapshot, response ID, provider, and reported cost; then appends a direct response with computed token cost to the same file and verifies the combined ceiling prevents another request. `TestJudgeLogUsesTheUTCMonthAndRejectsBadCosts` checks UTC month selection and fallback for null, negative, or nonnumeric reported cost. |
| No requests without budget evidence | `TestRunSendsNothingWithoutAKey`, `TestRunSendsNothingWithAnUnreadableLog`, and `TestRunStopsAtTheMonthlyCeiling` verify zero requests for the named conditions and stop after a call reaches the ceiling. `TestRunRechecksSpendBeforeRetryAndHonorsCancellation` also verifies spend rechecks before retries. |
| Stop/continue/retry rules | `TestRunStopsWhenTheServiceFails` checks exact reasons and no later judgment for 503, 401, 402, 403, exhausted 429, and exhausted 529. `TestRunSkipsARefusedRequestAndContinues` checks 422 followed by an answered judgment. `TestAskRetriesARateLimitThenAnswers` observes three attempts after two 429 responses with `Retry-After: 0`, one advisory answer, and three log lines. Cancellation, capped waits, default waits, and per-attempt timeout are covered by `TestRetryWaitsAndAttemptsHonorCancellation` and `TestRunRechecksSpendBeforeRetryAndHonorsCancellation`. |
| Complete private Judge Log | `TestJudgeLogRecordsEveryCall` checks all 29 fields, one line per attempt, directory 0700, file 0600, answer nullability, state hash length, and absence of the fake key. `TestRunHandlesMalformedAnswersAndWriteFailures` checks unreadable answers, the sole caller error for an unreadable PRD, a failed append preserving the answer in hand while stopping later calls, and credential redaction from transport errors and logs. |

`TestRequestsCarryOnlySpecArtifactText` exercises both transports with a real
temporary fixture repository holding a Go source sentinel and a findings
sentinel. Every state string, split at `the cited decision`, is asserted to
occur in its PRD, TechSpec, or accepted ADR; neither sentinel nor the fake key
appears in any body, and only Authorization and Content-Type are set. All new
tests inject their RoundTripper, keys, fixed Request clock, and temporary Home;
they open no socket and do not consult process credential variables.
`TestAskIgnoresUnknownFieldsAndRejectsMalformedValues` additionally protects
HTTP metadata from unknown response fields and rejects missing or mistyped
answers and negative token counts.
`TestRunReportsSkippedArtifactPathsWithoutSendingText` verifies that a
non-English PRD is reported with its repository-relative path and sends nothing.

Sabotage evidence (each mutation tested alone, each restored byte-for-byte):

| Gate sabotaged | Mutation and focused command | Observed failure |
| --- | --- | --- |
| Ceiling | Changed both `>= q.MonthlyCeilingUSD` guards to `>`. `GOCACHE=/private/tmp/roundfix-0205-task02-go-cache go test ./internal/judge -count=1 -run '^TestRunStopsAtTheMonthlyCeiling$'` | Exit 1: `already_exactly_at_ceiling` reported `unexpected request`; `run_reaches_ceiling` observed two calls and US$10 instead of one call and a stop. |
| Model pin | Accepted any `typesafe/` prefix in addition to the measured pattern. `GOCACHE=/private/tmp/roundfix-0205-task02-go-cache go test ./internal/judge -count=1 -run '^TestModelPinAcceptsOnlyJev113$'` | Exit 1: `typesafe/jev-1.14-20261101` reported `pin accepted=true want=false`. |
| Transport selection | Fell back to `OPENROUTER_API_KEY` for OpenRouter. `GOCACHE=/private/tmp/roundfix-0205-task02-go-cache go test ./internal/judge -count=1 -run '^TestRunIgnoresTheGenericOpenRouterKey$'` | Exit 1: `only_OPENROUTER_API_KEY` and `both_generic_keys` reported `unexpected request`. |

The subsequent race-enabled focused check ran after all three restorations.
The Task's declared Verification command and repository settlement gate were
not run; they remain Daemon-owned. No commit, push, or PR was made.

Environment limits: the first focused check with the default Go cache reported
standard-library packages as missing although their source directories existed;
the task-scoped GOCACHE allowed compilation and test execution. Fetching
`https://docs.typesafe.ai/api.md` through `rtk curl-cffi` was blocked by the
sandbox domain allowlist. Implementation follows this Spec's API Contract 3
and its maintainer-recorded endpoint measurements; no live API call was made.

## Carry-forward provenance

- Source Run: `run_20261001T205638Z_5e1457a22648dd9b`
- Source commit: `8e90c5863d62e54c1bb8cc04aa02f1db880bb402`
