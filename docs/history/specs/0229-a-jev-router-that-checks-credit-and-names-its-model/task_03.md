---
task: task_03
spec: 0229-a-jev-router-that-checks-credit-and-names-its-model
status: completed
type: backend
complexity: high
---

# Task 03: A loopback relay records the routed model on the Judge Log

## Overview

Both `router-prompt` Judge Log lines of the 2026-10-04 measurement recorded an
empty model, because neither the ACP stream nor OpenCode's records carry the
model the router chose; only OpenRouter's responses do. This Task adds the
relay of the TechSpec's "The relay": the ACPX runner points each routed
session's inline provider at a loopback relay under a per-session token, the
relay forwards requests and responses unchanged and notes each response's
`model`, `provider`, `id` and any HTTP 402, the runner returns that
observation with the prompt's result, and the Judge Log line records it
(API Contract 6, ADR-0234). Naming a credit refusal is task_04's.

## Requirements

1. MUST add, in a new file of the `jevrouter` package, `OpenRouterAPI`,
   `Relay`, `StartRelay`, `Open`, `Take`, `Release`, `Close`, `Observation`
   and `Refusal` with the TechSpec's Interfaces, built on the standard
   library's `httputil.ReverseProxy` as "The relay" describes: listen on
   `127.0.0.1` only; map `/<token>/<rest>` to `<upstream>/<rest>` with the
   query unchanged and no forwarding headers; flush streamed answers at once;
   answer 404 without any upstream request for a path whose first segment is
   not an open token; answer 502 on an upstream failure without logging.
2. MUST observe a copy of each response body without changing a byte: a
   `text/event-stream` body line by line, reading the JSON after `data: `, any
   other body whole up to 1 MiB; from each JSON object read only `id`,
   `model`, `provider` and, on a 402 status, `error.metadata.limit_source`;
   keep distinct models and providers in first-seen order, the last id, and
   the first 402 as a `Refusal`. MUST hold every invariant of "The relay":
   it never logs, stores or writes a header, a body or the key, and keeps no
   header value after a request ends.
3. MUST make the inline provider configuration a function of the base URL,
   keeping the key the `{env:ROUNDFIX_OPENROUTER_API_KEY}` placeholder, and
   add `ACPXRunner.RouterEndpoint` (empty means `jevrouter.OpenRouterAPI`).
   For a routed session, `codexEnvForSession` MUST start the runner's relay on
   first use, open one token per session name, reuse it for every later
   command of that session, and pass the relay's base URL; `clearSessionState`
   and the close of a disposable session MUST release the token, and the
   runner MUST close the relay when no token remains. The relay's serving
   goroutine is owned and awaited by `Close`.
4. MUST add `Router jevrouter.Observation` to `agent.ExecuteResult` and make
   `RunPrompt` take the session's observation after a routed prompt ends, on
   success and on failure, and return it there.
5. MUST add `Reported jevrouter.Observation` to `jevrouter.PromptRecord`, make
   `runPrepared` copy `result.Router` into it, and make `Ledger.Append` write
   `model` and `provider` as the distinct values joined by `", "` and
   `response_id` as the last id (API Contract 6); each stays empty when
   nothing was reported. The Judge Log schema and every other field of the
   line are unchanged.
6. MUST add the tests named in Verification as Testing Approach 4, 5 and 6
   describe: `internal/jevrouter/relay_test.go` (new) against a local upstream
   for byte-identical streamed and whole answers with the request's
   `Authorization` header and body reaching the upstream unchanged, the
   observation's order and last id, a noted 402 with its `limit_source`, a 404
   with no upstream request for an unknown token, and `Close`; in
   `internal/agent/jev_router_test.go`, through the compiled fake acpx, a
   loopback base URL with the key placeholder that is the same across two
   prompts of one session, a fake prompt that posts through that URL to a
   local upstream and returns its observation in `ExecuteResult.Router`, and a
   cleared session releasing its token; in `internal/jevrouter/ledger_test.go`
   a record reporting `a/one` and `b/two` written as `a/one, b/two` with its
   providers and last id; and in `internal/daemon/jev_router_gate_test.go` a
   routed prompt whose fake result carries an observation producing that line.
7. MUST change `TestJevRouterSessionCarriesThePlaceholderConfig` so it
   compares the routed child's configuration with the one built for the
   relay's base URL it received, and still proves the key is in no
   configuration, argument or acpx file.
8. MUST NOT send any test request beyond the loopback interface, change the
   router's selection rules, the gate's checks, the Judge Log schema or any
   default or Recommended Profile.

## Subtasks

- [ ] Add the relay with its observation and its token lifecycle.
- [ ] Point the routed session's inline provider at the relay and own the relay in the runner.
- [ ] Return the observation with the prompt's result and record it on the Judge Log line.
- [ ] Add the tests and update the placeholder test.

## Acceptance Criteria

- [ ] A streamed and a whole answer pass through the relay byte-identical,
      and an unknown token reaches no upstream.
- [ ] A routed session's OpenCode configuration names a loopback base URL and
      the key placeholder only, the same for every prompt of the session.
- [ ] A routed prompt whose responses report `a/one` then `b/two` writes one
      `router-prompt` line with `model` `a/one, b/two`, its providers and the
      last response id; the key appears in no log, file or Run Event.

## Context

- creates: `internal/jevrouter/relay.go`
- creates: `internal/jevrouter/relay_test.go`
- interface: `internal/jevrouter/ledger.go`
- interface: `internal/jevrouter/ledger_test.go`
- interface: `internal/agent/jev_router.go`
- interface: `internal/agent/jev_router_test.go`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/acpx_runner_test.go`
- interface: `internal/agent/agent.go`
- interface: `internal/daemon/agent_session_owner.go`
- interface: `internal/daemon/jev_router_gate_test.go`

## Verification

- `grep -q RouterEndpoint internal/agent/acpx_runner.go || { printf 'runner has no relay endpoint\n' >&2; exit 1; }; out="$(go test -count=1 -v -run "^(TestRelayForwardsStreamedAndWholeAnswersUnchanged|TestRelayNotesModelsProvidersAndTheLastResponseID|TestRelayNotesACreditRefusal|TestRelayRefusesAPathWithoutAnOpenToken|TestRelayCloseStopsServing|TestRouterLineRecordsTheReportedModels|TestRouterLineIsReadByTheJudgeLog|TestJevRouterSessionCarriesThePlaceholderConfig|TestJevRouterSessionPointsOpenCodeAtTheRelay|TestJevRouterPromptReturnsTheRelayObservation|TestJevRouterClearedSessionReleasesItsRelayToken|TestJevRouterPromptRecordsTheReportedModels|TestJevRouterPromptAppendsOneLine)$" ./internal/jevrouter ./internal/agent ./internal/daemon 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRelayForwardsStreamedAndWholeAnswersUnchanged TestRelayNotesModelsProvidersAndTheLastResponseID TestRelayNotesACreditRefusal TestRelayRefusesAPathWithoutAnOpenToken TestRelayCloseStopsServing TestRouterLineRecordsTheReportedModels TestRouterLineIsReadByTheJudgeLog TestJevRouterSessionCarriesThePlaceholderConfig TestJevRouterSessionPointsOpenCodeAtTheRelay TestJevRouterPromptReturnsTheRelayObservation TestJevRouterClearedSessionReleasesItsRelayToken TestJevRouterPromptRecordsTheReportedModels TestJevRouterPromptAppendsOneLine; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the runner has no relay endpoint and the relay and model tests do not exist, so the command fails.

## References

- `_prd.md` → User Story 4; Core Feature 5; Success Metric 3
- `_techspec.md` → The relay; The runner; The Judge Log line; Interfaces; API Contract 6; Testing Approach 4-6; Build Order 3
- ADR-0234; ADR-0218; ADR-0201

## Result

Implemented the Task 03 slice for Daemon Verification. The relay binds only
`127.0.0.1`, maps live session tokens to the upstream API through
`httputil.ReverseProxy`, forwards bodies unchanged, flushes immediately, and
retains only models, providers, the last response id and the first HTTP 402.
Whole-body parsing stops at 1 MiB; malformed or oversized bodies still pass
unchanged. Unknown and released tokens return 404 without contacting upstream;
upstream failures return 502 without diagnostics. `Close` awaits the serving
goroutine.

Routed ACPX commands receive one stable loopback base URL per session and the
existing environment-key placeholder. Normal and disposable session cleanup
release tokens through `clearSessionState`; releasing the last token closes
the runner's relay. `RunPrompt` takes the observation on success and failure,
`runPrepared` copies it into the prompt record, and the existing Judge Log
fields carry models and providers joined by `", "` and the last response id.
Credit-refusal naming remains task_04's slice.

Acceptance evidence from focused checks:

- Byte-identical streamed and whole responses, unchanged upstream request body
  and Authorization, unchanged query, stripped forwarding headers, unknown and
  released token rejection, model/provider ordering, last id, first 402,
  malformed/oversized passthrough and serving shutdown are exercised by the
  `TestRelay*` checks. `TestRelayFlushesBeforeUpstreamFinishes` receives the
  first streamed bytes before allowing the upstream to finish;
  `TestRelayKeepsTokensAndTakesIndependent` also checks escaped paths and token
  isolation.
- `TestJevRouterSessionPointsOpenCodeAtTheRelay` crosses the compiled fake
  acpx boundary twice and checks the same loopback URL and key placeholder.
  The updated `TestJevRouterSessionCarriesThePlaceholderConfig` compares
  against the configuration built for that observed URL and excludes the
  sentinel key from configuration, acpx arguments and acpx configuration.
  `TestJevRouterClearedSessionReleasesItsRelayToken` checks 404 while a second
  session remains live and listener shutdown after the last token is released.
- `TestJevRouterPromptReturnsTheRelayObservation` posts through the configured
  relay to a local upstream on both successful and failed prompts and checks
  the returned observation and its clearing. `TestRouterLineRecordsTheReportedModels`
  and `TestJevRouterPromptRecordsTheReportedModels` check one Judge Log line
  with `a/one, b/two`, `p, q` and `last`. Existing empty-report/schema checks
  remain in the focused selection; daemon tests exclude the sentinel key from
  the Judge Log and Run Events and exclude private prompt/output text.

Focused commands:

- `GOCACHE=/tmp/roundfix-task03-gocache rtk proxy go test -count=1 -run 'TestRelay|TestRouterLine|TestJevRouter(Session|Prompt|Cleared)' ./internal/jevrouter ./internal/agent ./internal/daemon`
  — exit 0 for all three packages against the final source.
- `GOCACHE=/tmp/roundfix-task03-gocache rtk proxy go test -race -count=1 -run 'TestRelay|TestRouterLine|TestJevRouter(Session|Prompt|Cleared)' ./internal/jevrouter ./internal/agent ./internal/daemon`
  — exit 0 for all three packages; no race reported.
- `rtk proxy git -c core.fsmonitor=false diff --check` — exit 0.
- `GOCACHE=/tmp/roundfix-task03-gocache rtk make verify-incremental`
  — final frozen-tree run exited 0: formatting, vet, repository tests, skill
  checks and build passed. The first run exited 2 because I edited
  `internal/agent/jev_router.go` and `internal/jevrouter/relay_test.go` during
  the tests; the suite guard reported those concurrent edits. The rerun kept
  the repository unchanged until the gate exited.

The default Go cache initially reported the installed standard-library
`net/http/httputil` package missing; the isolated task cache resolved it. The
sandbox denied local test listeners, so the successful checks used approved
execution with loopback sockets available. Every test HTTP request stays on
loopback. No live OpenRouter or TypeSafe API request was made.

The authored Verification command was not run. Task status, the Task Graph,
other Task files and selection/gate/default/Profile behavior were not edited.
No commit, push or Pull Request was created.

## Carry-forward provenance

- Source Run: `run_20261005T141012Z_a1b532ab932de7ab`
- Source commit: `f9328e75126d5f42b75ee5be4c489aeaeb7a89bc`
