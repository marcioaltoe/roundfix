---
task: task_02
spec: 0237-review-selection-failures-that-say-why
status: completed
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

## Result

The runner now retains the first JSON-RPC error's step, code, bounded message,
and whether the prompt had been sent at that failure. It correlates echoed
client requests and responses by id and derives the documented fallback step
when an error has a null or unmatched id, or when stdout contains no error.
Only a non-zero exit without a parsed result attaches the detail to a mapped
Batch failure; conversion to a Selection failure copies it without changing
its reason or wrapped cause. Both error texts name the step and adapter
message. `DescribeProtocolFailure` also places the documented preparation
errors and bounds their infrastructure stderr tails.

Acceptance evidence:

- Different step errors and prompt position: the new runner tests exercise
  captured acpx 0.19.4 stdout for `initialize`, `session/set_model`,
  `session/prompt`, and a null-id disconnect. They assert the three distinct
  error texts, JSON-RPC code, prompt-sent flags, first-error retention,
  string/unmatched/null ids, answered requests, and empty stdout.
- Bounded messages without request content: the message test checks the exact
  rune-safe truncation length, collapsed whitespace, description and both
  error texts, and exclusion of error `data`, model parameters and prompt
  content. The preparation test checks bounded, valid UTF-8 stderr tails,
  every specified preparation type, wrapped errors, protocol precedence and
  refusal to place plain errors.
- Existing classification: the new mapping test checks the unchanged reason,
  stderr-derived wrapped error, Batch exit code/stderr, infrastructure and
  stop paths, exit-zero errors without protocol detail, and parsed-result
  transport anomalies. The unchanged exit-classification matrix also passed.

Focused checks (all with `GOCACHE=/private/tmp/roundfix-task02-gocache`):

- Red starting point: `rtk proxy go test ./internal/agent -run
  '^TestACPXPromptNamesTheFailingProtocolStep$' -count=1` failed to compile
  because the protocol interface did not yet exist.
- `rtk proxy go test ./internal/agent -run
  'TestACPXPromptNames|TestACPXPromptRecords|TestACPXPromptBounds|TestACPXPromptProtocolDetail|TestDescribeProtocol'
  -count=1`: exit 0 after the final fixture and assertion edits.
- `rtk proxy go test ./internal/agent -run
  '^TestACPXPromptExitClassificationMatrix$' -count=1`: exit 0; the existing
  runner test file was not edited.
- `rtk proxy git -c core.fsmonitor=false diff --check`: exit 0.

Fixture provenance: installed `acpx --version` returned `0.19.4`. A disposable
local fake adapter produced the captured `acpx prompt` stdout; only its
working-directory path is normalized to `/tmp` in the test constants. The
three named-step failures and disconnect each exited 1. The live named-session
command printed its normal session banner to stderr; runner fixtures use
empty stderr to isolate the stdout diagnostic. The sandbox denied the local
queue-owner socket, so recording was repeated with host permission. Recording
used no remote adapter, credentials or network calls. Raw captures remain in
`/var/folders/_7/68y3l_1s55jcsdmmcmmm4dhh0000gn/T/roundfix-task02-acpx-i_mb57cz/`
and
`/var/folders/_7/68y3l_1s55jcsdmmcmmm4dhh0000gn/T/roundfix-task02-acpx-cq0x54ay/`.

The first `rtk make verify-incremental` attempt exited 2: two CLI owner-process
integration tests could not enumerate the host process table, and suiteguard
correctly detected the fixture edit made while that check was running. A new
attempt with host permission and a stable worktree exited 0:
`GOCACHE=/private/tmp/roundfix-task02-gocache rtk make verify-incremental`
passed formatting, vet, the package suite (including unchanged runner tests),
skill synchronization/checks and the CLI build. The log is
`/private/tmp/roundfix-task02-incremental-host.log`; the hidden final lines were
recovered with `rtk recall 918dfbd0df8a` without rerunning the check. The declared
Task Verification remains for the Daemon. Status, the Task Graph, other Task
files, daemon code, command arguments and existing runner tests were left
unchanged; no commit, push or Pull Request was made. Review retry behavior
belongs to task_03 and is outside this diff.
