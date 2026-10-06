---
task: task_02
spec: 0237-review-selection-failures-that-say-why
status: pending
type: backend
complexity: medium
---

# Task 02: The ACPX Runner keeps the failing protocol step, the adapter's message and whether the prompt was sent

## Overview

An `acpx prompt` process that exits non-zero without a parsed result today
loses the JSON-RPC error line that carried the adapter's message, so a failure
at `initialize`, `session/set_model` or `session/prompt` reads the same
`agent/protocol error`. This Task keeps the first error line, names the step
by the request it answers, records whether `session/prompt` had been sent, and
adds `DescribeProtocolFailure`, which also places a session-preparation
failure. It answers the Backlog Entry "A pre-PR review agent selection failure
says nothing actionable" of 2026-10-06. Verifiable on its own at the runner
seam with stdout recorded from acpx 0.19.4.

## Requirements

1. MUST add `ProtocolFailure`, `ProtocolDescription`,
   `ProtocolFailureMessageLimit` and `DescribeProtocolFailure` with the shapes
   of `_techspec.md` → Interfaces, and the field `Protocol *ProtocolFailure`
   on `BatchFailureError` and `SelectionFailureError`.
2. MUST track client requests, their responses, `PromptSent` and the first
   JSON-RPC error while `readPromptStream` parses stdout, per Invariants 1 to
   3, reading the JSON-RPC `id` without changing how any other line is
   handled.
3. MUST set `Protocol` only on a `BatchFailureError` mapped from a non-zero
   exit without a parsed prompt result, and copy it in
   `classifyNoOutputFailure`, per Invariant 6; `Reason`, `Err`, `Stderr`, the
   exit-code mapping and every other path stay as they are.
4. MUST print the step and the message in both error types' text as
   `_techspec.md` → Error text states, and bound every message per
   Invariant 5: only the JSON-RPC `message`, never `data`, request parameters
   or the prompt; one line; at most `ProtocolFailureMessageLimit` bytes on a
   rune boundary.
5. MUST make `DescribeProtocolFailure` place the errors of Invariant 4,
   reading an `InfrastructureError`'s stderr tail as the message of a
   preparation failure, and return `false` for any other error.
6. MUST add the tests named in Verification to a new
   `internal/agent/protocol_failure_test.go`, driving `RunPrompt` through
   `runFakeACPXPrompt` with stdout recorded from acpx 0.19.4 (`_prd.md` →
   Acceptance evidence): failures at `initialize`, `session/set_model` and
   `session/prompt` give three different errors naming their step, code and
   message, with `PromptSent` false, false and true; an error line with a null
   id after `session/prompt` names `session/prompt` with `PromptSent` true; an
   exit `1` with no stdout names `adapter startup`; a message longer than the
   limit and spread over lines is cut to one line within the limit without a
   broken rune; `Reason` stays `agent/protocol error` and `Err` keeps today's
   stderr-derived value; and `DescribeProtocolFailure` places each
   preparation type of Invariant 4 and refuses a plain error.
7. MUST NOT change the daemon packages, the acpx arguments of any command, or
   the existing tests in `internal/agent/acpx_runner_test.go`.

## Subtasks

- [ ] Add the protocol failure types and the description function.
- [ ] Track requests, responses and the first error in the prompt stream.
- [ ] Attach the trace on a non-zero exit without a parsed result.
- [ ] Print the step and the bounded message in both error types.
- [ ] Add the runner and description tests.

## Acceptance Criteria

- [ ] The three recorded step failures produce three different errors, each
      naming its step and the adapter's message, and report whether the
      prompt was sent.
- [ ] Every message is one line within the limit and carries no `data` or
      request content.
- [ ] Reason, wrapped error, exit mapping and the existing runner tests are
      unchanged.

## Context

- creates: `internal/agent/protocol_failure.go`
- creates: `internal/agent/protocol_failure_test.go`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/agent.go`
- instruction: `internal/agent/acpx_runner_test.go`
- instruction: `docs/adr/0242-a-blocked-review-names-the-protocol-step-and-retries-once-before-the-prompt.md`
- instruction: `docs/adr/0020-parsed-prompt-result-outranks-acpx-exit-code.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestACPXPromptNamesTheFailingProtocolStep|TestACPXPromptRecordsWhetherThePromptWasSent|TestACPXPromptBoundsTheAdapterMessage|TestACPXPromptProtocolDetailKeepsTheExitMapping|TestDescribeProtocolFailurePlacesSessionPreparation)$' ./internal/agent 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestACPXPromptNamesTheFailingProtocolStep TestACPXPromptRecordsWhetherThePromptWasSent TestACPXPromptBoundsTheAdapterMessage TestACPXPromptProtocolDetailKeepsTheExitMapping TestDescribeProtocolFailurePlacesSessionPreparation; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the five tests do not exist, so the command fails.

## References

- `_prd.md` → Core Feature 1; User Story 1; Success Metric 1; Success Metric 2; Goal 1; Acceptance evidence
- `_techspec.md` → Interfaces; Invariants 1 to 6; Error text; API Contract 1; API Contract 2; Testing Approach 1; Build Order 2
- ADR-0242
- ADR-0020
