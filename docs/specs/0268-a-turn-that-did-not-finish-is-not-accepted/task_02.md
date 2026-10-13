---
task: task_02
spec: 0268-a-turn-that-did-not-finish-is-not-accepted
status: pending
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
