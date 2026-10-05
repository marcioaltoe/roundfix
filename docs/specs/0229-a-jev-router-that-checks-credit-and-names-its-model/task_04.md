---
task: task_04
spec: 0229-a-jev-router-that-checks-credit-and-names-its-model
status: completed
type: backend
complexity: medium
---

# Task 04: An OpenRouter credit refusal has its own name

## Overview

In the 2026-10-04 Jev Router measurement OpenRouter refused a routed repair
prompt with HTTP 402 ("would exceed your available credits"), and Roundfix
reported `agent/protocol error`. With task_03's relay noting a 402, this Task
names that failure `openrouter_credit_refused` (API Contract 5), as the
TechSpec's "The runner" and "The Judge Log line" state: before Agent work it
is a failed selection that the Fallback Chain takes; after it the Work Item
fails with that reason, as `jev_ceiling_reached` does (ADR-0234, ADR-0114).

## Requirements

1. MUST, when a routed prompt fails and its observation carries a `Refusal`,
   replace the prompt's failure reason with API Contract 5, naming the
   `limit_source` or `unspecified`: without Agent output, return a
   `SelectionFailureError` for the runtime and report it as a failed
   selection; with Agent output, return a `BatchFailureError` whose `Reason`
   is that text in place of `agent/protocol error`, keeping its exit code and
   stderr. A routed prompt without a `Refusal`, and every non-routed prompt,
   keeps its current classification.
2. MUST make `selectionReasonCode` return `openrouter_credit_refused` for that
   failure, as it returns `jev_ceiling_reached`, so the fallback receipt
   names it.
3. MUST make `Ledger.Append` start the line's `error` with
   `openrouter_credit_refused: <limit_source>` when the record's observation
   carries a `Refusal`, followed by `"; "` and the key-usage error when that is
   present, and leave the outcome `skipped` for a failed prompt.
4. MUST add the tests named in Verification as Testing Approach 5 and 6
   describe: in `internal/agent/jev_router_test.go`, through the compiled fake
   acpx and a local upstream answering 402 with `limit_source`
   `openrouter_credits`, a failed prompt with no Agent output returns a
   `SelectionFailureError` with API Contract 5 and one with Agent output a
   `BatchFailureError` whose message carries API Contract 5 and not
   `agent/protocol error`; in `internal/jevrouter/ledger_test.go` the refusal
   written first in `error`; and in `internal/daemon/jev_router_gate_test.go`
   a fake runner's `openrouter_credit_refused` selection failure falling back
   before work with a receipt naming it, and its batch failure after work
   failing the Task with that reason.
5. MUST NOT retry the refused request, switch models after Agent work began,
   change the gate's checks, the relay's forwarding, or any default or
   Recommended Profile; no test reaches beyond the loopback interface.

## Subtasks

- [ ] Name the refusal in the runner before and after Agent work.
- [ ] Return the reason code from the session owner.
- [ ] Record the refusal on the Judge Log line.
- [ ] Add the tests.

## Acceptance Criteria

- [ ] A routed prompt refused with 402 before Agent output is a failed
      selection carrying
      `openrouter_credit_refused: OpenRouter refused a routed request for credit (openrouter_credits)`,
      and the fallback receipt names `openrouter_credit_refused`.
- [ ] After Agent output the Task fails with that reason and never with
      `agent/protocol error`.
- [ ] The prompt's `router-prompt` line carries the refusal in `error`.

## Context

- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/acpx_runner_test.go`
- interface: `internal/agent/jev_router_test.go`
- interface: `internal/jevrouter/ledger.go`
- interface: `internal/jevrouter/ledger_test.go`
- interface: `internal/daemon/agent_session_owner.go`
- interface: `internal/daemon/jev_router_gate_test.go`

## Verification

- `grep -q openrouter_credit_refused internal/agent/acpx_runner.go || { printf 'runner does not name the refusal\n' >&2; exit 1; }; grep -q openrouter_credit_refused internal/daemon/agent_session_owner.go || { printf 'session owner has no reason code\n' >&2; exit 1; }; out="$(go test -count=1 -v -run "^(TestJevRouterCreditRefusalBeforeWorkIsAFailedSelection|TestJevRouterCreditRefusalAfterWorkNamesItsReason|TestRouterLineRecordsACreditRefusal|TestJevRouterCreditRefusalFallsBackBeforeWork|TestJevRouterCreditRefusalFailsTheTaskAfterWork|TestLedgerRecordsUnreadableUsageAndClampsNegativeDelta)$" ./internal/agent ./internal/jevrouter ./internal/daemon 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestJevRouterCreditRefusalBeforeWorkIsAFailedSelection TestJevRouterCreditRefusalAfterWorkNamesItsReason TestRouterLineRecordsACreditRefusal TestJevRouterCreditRefusalFallsBackBeforeWork TestJevRouterCreditRefusalFailsTheTaskAfterWork TestLedgerRecordsUnreadableUsageAndClampsNegativeDelta; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task neither the runner nor the session owner names the refusal and the new tests do not exist, so the command fails.

## References

- `_prd.md` → User Story 3; Core Feature 4; Success Metric 2
- `_techspec.md` → The runner; The Judge Log line; API Contract 5; Testing Approach 5; Testing Approach 6; Build Order 4
- ADR-0234; ADR-0114; ADR-0050; ADR-0218


## Result

Implemented the task_04 refusal classification slice. The runner consumes the
routed prompt's relay observation before classifying its failure, replaces the
Batch failure reason with API Contract 5, and uses the existing no-output
selection classifier. It retains exit code and stderr after Agent output.
The session owner recognizes `openrouter_credit_refused` in fallback receipts.
The Judge Log starts `error` with the refusal and limit source, appends an
unreadable key-usage error after `; `, and keeps failed prompts `skipped`.
An absent limit source is `unspecified`. Gate checks, forwarding, profiles and
retry policy were not changed.

Acceptance evidence:

- Before work: `TestJevRouterCreditRefusalBeforeWorkIsAFailedSelection` uses
  the compiled acpx fixture and a loopback 402 upstream, checks the exact
  selection reason, runtime and failed-selection status, and asserts one
  upstream request. `TestJevRouterCreditRefusalFallsBackBeforeWork` checks
  routed-to-fallback execution and notification ordering with the receipt's
  `openrouter_credit_refused` code and full reason.
- After work: `TestJevRouterCreditRefusalAfterWorkNamesItsReason` checks the
  Batch failure reason, exit code and stderr after an Agent message, and
  excludes `agent/protocol error`. The daemon test
  `TestJevRouterCreditRefusalFailsTheTaskAfterWork` emits an Agent message,
  checks the work-start status and persisted failed Task outcome, and
  excludes retry or fallback activation.
- Judge Log: `TestRouterLineRecordsACreditRefusal` reads the appended line
  through the shared judge reader and checks refusal-first error ordering,
  both limit-source cases, optional usage errors, cost and `skipped` outcome.

Focused checks:

- Red reproduction: `rtk proxy go test ./internal/agent -run
  'TestJevRouterCreditRefusal' -count=1` exited 1 before production edits,
  showing both before/after failures still named `agent/protocol error`.
  The initial three-package focused run also exposed the missing ledger
  refusal and fallback reason code; its agent tests encountered the sandbox
  loopback-listener restriction. Subsequent local listener/process checks
  ran with approved sandbox escalation.
- `rtk proxy go test ./internal/agent ./internal/jevrouter ./internal/daemon
  -run 'Test(JevRouterCreditRefusal|RouterLineRecordsACreditRefusal|JevRouterPromptReturnsTheRelayObservation|LedgerRecordsUnreadableUsage)'
  -count=1` exited 0 in all three packages after the production edits.
- After strengthening the daemon post-work test to emit Agent output,
  `rtk proxy go test ./internal/daemon -run 'TestJevRouterCreditRefusal'
  -count=1` exited 0.
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0.
- Initial `rtk make verify-incremental` exited 2: the CLI tests reported
  `PASS`, but suiteguard detected this Agent's concurrent Result edit to
  `task_04.md`. All other packages passed. The worktree was then frozen for
  a fresh incremental run; this was not a production or assertion failure.
- Fresh `rtk make verify-incremental` exited 0 with the worktree held unchanged
  during testing: formatting, vet, Go tests, skill checks and CLI build passed.
  Only this Result evidence was appended afterward.

The authored Verification commands were not run. Status remains Daemon-owned;
no other Task, Task Graph, commit, push or Pull Request was changed.

## Carry-forward provenance

- Source Run: `run_20261005T141012Z_a1b532ab932de7ab`
- Source commit: `807d6896a642308a52ef76b1b6c556240b588037`
