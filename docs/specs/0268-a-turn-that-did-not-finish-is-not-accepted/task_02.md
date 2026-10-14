---
task: task_02
spec: 0268-a-turn-that-did-not-finish-is-not-accepted
status: completed
type: backend
complexity: medium
---

# Task 02: The ACPX Runner never accepts a turn that the provider failed

## Overview

Today `RunPrompt` returns success for any non-empty stop reason, and the
stream reader drops `session_info_update` metadata, so a Codex turn that
streamed `threadStatus` `systemError` and the text "Selected model is at
capacity. Please try a different model." before `end_turn` is accepted as a
finished turn. This Task marks that turn as a Provider System Error, keeps the
provider text as its bounded detail, and refuses it at the sealed prompt
boundary too. It answers the Backlog Entry
[a model-at-capacity reply is accepted as a finished turn](references/2026-10-13-a-model-at-capacity-reply-is-accepted-as-a-finished-turn.md)
of 2026-10-13. It is verifiable on its own at the runner seam with synthetic
stdout in the recorded shape.

## Requirements

1. MUST add `ProviderSystemErrorStatus`, `ProviderSystemErrorReason`,
   `ProviderSystemError` and `DescribeProviderSystemError` with the shapes of
   `_techspec.md` → Interfaces in a new `internal/agent/provider_system_error.go`,
   and the field `ProviderSystemError bool` on `ProtocolFailure`.
2. MUST track, in the prompt stream reader, the turn's last
   `_meta.codex.threadStatus.type` and the `agent_message_chunk` text after a
   `systemError`, per Invariants 1 and 2, and make `RunPrompt` return the
   failure of Invariant 4 with the text of API Contract 1, only on the path of
   Invariant 3.
3. MUST append `: <detail>` after the step in the protocol text of a
   `ProviderSystemError` protocol, so the text matches API Contract 1, and
   keep every other error text, `Reason`, `Err`, `Stderr` and exit mapping as
   they are.
4. MUST make `DescribeProviderSystemError` behave per Invariant 5.
5. MUST make the sealed prompt parser refuse a Provider System Error per
   Invariant 6.
6. MUST add `internal/agent/provider_system_error_test.go` with the tests
   named in Verification, driving `RunPrompt` through `runFakeACPXPrompt` with
   synthetic stdout shaped like the recorded stream (fake ids, a
   `session/prompt` request, `threadStatus` `active`, then `systemError`, then
   a capacity chunk, then `{"stopReason":"end_turn","usage":null}`):
   - `TestProviderSystemErrorIsNotAFinishedTurn`: with an Agent message before
     the status, the error is a `*BatchFailureError` whose text is
     `Agent Batch failed: provider system error at session/prompt: Selected model is at capacity. Please try a different model.`,
     and the result still carries the stop reason and messages.
   - `TestProviderSystemErrorWithoutOutputIsASelectionFailure`: with no
     Agent output at all (the status, then the result), the error is a
     `*SelectionFailureError` described with the runtime `codex` and the
     detail `threadStatus systemError`. The provider's text is Agent output
     to the existing classifier, so a turn that streamed it stays a Batch
     failure.
   - `TestRecoveredThreadStatusIsAFinishedTurn`: `systemError` followed by
     `active` and `idle` before `end_turn` is a successful turn.
   - `TestProviderSystemErrorDetailIsBounded`: no chunk after the status gives
     `threadStatus systemError`; a multi-line chunk longer than
     `ProtocolFailureMessageLimit` is cut to one line within the limit without
     a broken rune.
   - `TestSealedPromptRefusesAProviderSystemError`: the sealed parser returns
     the error of Invariant 6 and no output.
   - `TestDescribeProviderSystemErrorTraversesFailureChains`: wrapped and
     joined chains are described; `StopError`, `context.Canceled` and a Lost
     Rollout are not.
7. MUST add `internal/cli/review_provider_system_error_test.go` with
   `TestReviewRecordsAProviderSystemErrorAsARuntimeFailure`: given a Provider
   System Error and a result whose answer is the capacity text, the review
   classifier records a `review runtime failure:` reason naming the detail and
   never parses the answer as review output. It calls `t.Parallel()` first
   (ADR-0259).
8. MUST NOT change the daemon, delivery or `internal/cli` production code, the
   acpx arguments of any command, or any existing test in `internal/agent`.

## Subtasks

- [ ] Add the marker, the description function and the protocol field.
- [ ] Track the last thread status and the provider text in the stream reader.
- [ ] Return the failure from `RunPrompt` and refuse it in the sealed parser.
- [ ] Add the runner and review classifier tests.

## Acceptance Criteria

- [ ] The capacity turn fails with the API Contract 1 text, with and without
      prior Agent output.
- [ ] A turn whose status recovers before its result is still a finished turn.
- [ ] The sealed boundary refuses the turn and returns no proposal.
- [ ] Every existing `internal/agent` test, including the Lost Rollout and
      protocol failure tests, passes unchanged.

## Context

- creates: `internal/agent/provider_system_error.go`
- creates: `internal/agent/provider_system_error_test.go`
- creates: `internal/cli/review_provider_system_error_test.go`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/protocol_failure.go`
- interface: `internal/agent/sealed.go`
- instruction: `internal/agent/lost_rollout.go`
- instruction: `internal/agent/lost_rollout_test.go`
- instruction: `internal/agent/acpx_runner_test.go`
- instruction: `internal/cli/review.go`
- instruction: `docs/adr/0274-a-turn-that-did-not-finish-is-not-accepted.md`
- instruction: `docs/adr/0020-parsed-prompt-result-outranks-acpx-exit-code.md`

## Consumers

- `BatchFailureError` and `SelectionFailureError` text and `DescribeProtocolFailure` for a Provider System Error — consumers: `./internal/agent`, `./internal/cli`

## Verification

- `out="$(go test -count=1 -v -run '^(TestProviderSystemErrorIsNotAFinishedTurn|TestProviderSystemErrorWithoutOutputIsASelectionFailure|TestRecoveredThreadStatusIsAFinishedTurn|TestProviderSystemErrorDetailIsBounded|TestSealedPromptRefusesAProviderSystemError|TestDescribeProviderSystemErrorTraversesFailureChains|TestLostRolloutIsNamedMidTurn|TestOtherErrorsAreNotLostRollouts|TestDescribeLostRolloutTraversesFailureChains)$' ./internal/agent 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestProviderSystemErrorIsNotAFinishedTurn TestProviderSystemErrorWithoutOutputIsASelectionFailure TestRecoveredThreadStatusIsAFinishedTurn TestProviderSystemErrorDetailIsBounded TestSealedPromptRefusesAProviderSystemError TestDescribeProviderSystemErrorTraversesFailureChains TestLostRolloutIsNamedMidTurn; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the six new tests do not exist, so their PASS lines are missing and the command fails.
- `out="$(go test -count=1 -v -run '^(TestReviewRecordsAProviderSystemErrorAsARuntimeFailure)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestReviewRecordsAProviderSystemErrorAsARuntimeFailure' || { printf 'missing pass: TestReviewRecordsAProviderSystemErrorAsARuntimeFailure\n' >&2; exit 1; }` — expected: exit 0; before this Task the test does not exist, so its PASS line is missing.

## References

- [_prd.md](_prd.md) — Goal 1; User Story 1; Core Feature 1; Success Metric 1; Acceptance evidence
- [_techspec.md](_techspec.md) — Interfaces; Invariants 1-7; API Contract 1; Testing Approach 1; Testing Approach 2; Build Order 2
- [2026-10-13-a-model-at-capacity-reply-is-accepted-as-a-finished-turn.md](references/2026-10-13-a-model-at-capacity-reply-is-accepted-as-a-finished-turn.md)
- ADR-0274; ADR-0020

## Result

Implemented the Task 02 runner slice for Daemon Verification. The prompt stream
now retains the last Codex thread status and the text after a `systemError`
until the parsed prompt result. A zero-exit turn with that final status returns
a marked Provider System Error through the existing Batch/Selection classifier,
with bounded detail and the API Contract 1 text. A later status clears the
failure. The sealed parser uses the same tracking and returns the specified
error with no output.

`DescribeProviderSystemError` traverses wrapped and joined Batch/Selection
failures, retains the selection runtime, and excludes cancellation, deadlines,
StopError values and pointers, and unmarked Lost Rollouts. Protocol text gains
the provider detail only for the new marker. Existing nonzero-exit mapping,
parsed-result anomalies, read-only refusals, command arguments, and CLI
production code are unchanged.

### Focused implementation evidence

- Before implementation,
  `rtk proxy go test -count=1 -run '^TestProviderSystemErrorIsNotAFinishedTurn$' ./internal/agent`
  exited 1: the recorded-shape capacity stream returned a nil error where the
  regression required a Batch failure. The same command exited 0 after the
  runner change.
- After the final code changes,
  `rtk proxy go test -count=1 -run 'TestProviderSystemError|TestRecoveredThreadStatus|TestSealedPromptRefusesAProvider|TestDescribeProviderSystem|TestLostRollout|TestOtherErrorsAreNotLostRollouts|TestDescribeLostRollout|TestACPXPromptNamesTheFailingProtocolStep|TestACPXPromptProtocolDetailKeepsTheExitMapping|TestDescribeProtocolFailure|TestSealedPromptStreamIgnoresNonTerminalJSONRPCResults|TestACPXRunPromptKeeps' ./internal/agent`
  exited 0 (`ok roundfix/internal/agent`).
- After the final code changes,
  `rtk proxy go test -count=1 -run 'TestReviewRecordsAProvider|TestReviewCommandBlocksOnRuntimeFailure' ./internal/cli`
  exited 0 (`ok roundfix/internal/cli`).
- `rtk proxy git diff --check` exited 0.

### Acceptance criterion evidence

1. Capacity failures: `TestProviderSystemErrorIsNotAFinishedTurn` asserts the
   exact Batch error with prior Agent output, preserved raw stdout, messages,
   stop reason, and null usage. `TestProviderSystemErrorWithoutOutputIsASelectionFailure`
   asserts the exact Selection error, runtime and default detail when no Agent
   output exists; its provider-text case asserts that the capacity chunk alone
   still makes a Batch failure. `TestProviderSystemErrorDetailIsBounded` checks
   chunk concatenation, exclusion of earlier text, absent status metadata,
   single-line detail, byte bounds and intact UTF-8.
2. Recovery: `TestRecoveredThreadStatusIsAFinishedTurn` checks `active`, `idle`,
   and `active` then `idle` recovery through both parsers.
   `TestProviderSystemErrorUsesTheStatusAtTheResult` checks that updates after
   the result do not change its classification, capacity text alone does not
   mark a failure, and a new failure discards an older failure's text.
3. Sealed refusal: `TestSealedPromptRefusesAProviderSystemError` asserts the
   specified error and empty output for capacity text and absent text.
4. Existing-test compatibility: the focused agent command above passed the
   selected existing Lost Rollout, protocol failure, sealed nonterminal-result,
   parsed-result anomaly and read-only refusal tests unchanged.
   `TestProviderSystemErrorKeepsExistingExitPaths` also checks the new status
   against nonzero-exit anomalies, read-only refusal and exit 130. Every
   pre-existing `internal/agent` test file is unchanged. The complete agent
   suite and repository gates remain for Daemon settlement; no full-suite
   claim is made here.

The new CLI test calls `t.Parallel()` first and asserts a blocked review with
`review runtime failure:` plus the provider detail, no findings, and the
`session/prompt` failed step, even for an apparent clean verdict. The new
chain-description test covers wrapped/joined failures, cancellation, deadlines,
StopError values and pointers, nil failures and Lost Rollouts.

The Task's declared Verification commands and selected repository Verification
were not run by this child Agent. Task status, the Task Graph and other Task files remain
Daemon-owned. No commit, push or pull request was made. No follow-up changes
were added to this slice.


### Verification feedback — attempt 1

Inspected the Daemon diagnostic artifact:
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261014T125708Z_601c103ef8eee9d3/verification/batch-002-attempt-1.log`.
The Daemon's `make verify-changed` exited 2. Its output records the complete
`internal/agent` package as passing, but the existing CLI delivery test
`TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate/stdout` observed a
child exit code of -1 and no output rather than the fixture's expected exit 2.
This is recorded Daemon evidence, not a child rerun of the full gate.

Inspected `deliver_item_binary_test.go`, its fixture builder, and
`stepExecutable`/`runRoundfix` in `deliver_workflow.go`. This test compiles a
standalone fixture binary, copies it into a disposable repository for each
step, and probes it with `migrate --check`. The exit code -1 comes from the
child's `exec.ExitError.ExitCode()`, indicating signal termination. The new
review classifier test changes no process environment, starts no child, and
uses no shared fixture state. No Task 02 defect was identified; the cause of
the child termination remains unproven. Delivery production code and its
existing tests are outside this Task's slice and were not changed.

Focused investigation and checks:

- `rtk proxy go test -count=1 -v -run 'TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate' ./internal/cli`
  exited 0, including both stream cases.
- A temporary Go overlay restored `acpx_runner.go`, `protocol_failure.go` and
  `sealed.go` from baseline revision
  `e1e376dd9518e19a279093d37bbc219821da029a` and replaced this Task's three new
  Go files with package-only files. Repository source files were untouched.
  `go test -overlay /private/tmp/roundfix-task02-baseline-6u9hxnbg/overlay.json -count=1 -v -run TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate ./internal/cli`
  exited 0.
- `rtk proxy go test -overlay /private/tmp/roundfix-task02-baseline-6u9hxnbg/overlay.json -count=10 -run 'TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate' ./internal/cli`
  exited 0. This bounded baseline probe did not reproduce the signal
  termination; it does not establish the failure's cause or resolution.
- `rtk proxy go test -count=1 -run 'TestReviewRecordsAProvider|TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate' ./internal/cli`
  exited 0 with both the Task-owned classifier test and the reported delivery
  test selected together.
- `rtk proxy go test -count=1 -run 'TestProviderSystemError|TestRecoveredThreadStatus|TestSealedPromptRefusesAProvider|TestDescribeProviderSystem' ./internal/agent`
  exited 0, refreshing the focused evidence for capacity failure,
  recovery, sealed refusal, bounded detail and failure-chain traversal.
- `rtk proxy git diff --check` exited 0 after this Result update.

Only Result evidence was changed during this feedback turn. The Daemon's full
gate failure remains unresolved; a passing isolated check is not a passing
repository gate. If signal termination recurs, it needs a separate delivery
fixture/process investigation. No declared Verification command or selected
repository gate was rerun by this child, and Task status was not edited.

## Carry-forward provenance

- Source Run: `run_20261014T125708Z_601c103ef8eee9d3`
- Source commit: `a4555eb17560521207367ff3f02541834e8da828`
