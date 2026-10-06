---
task: task_02
spec: 0240-a-lost-rollout-is-infrastructure
status: completed
type: backend
complexity: medium
---

# Task 02: The ACPX Runner names a Lost Rollout by its code and phrase

## Overview

Today a lost Codex rollout reaches Roundfix as `agent/protocol error at
session setup: Internal error (JSON-RPC -32603)` or `... at session/prompt:
Internal error`. The runner discards `data.details`, which holds `no rollout
found for thread id <uuid>`, and does not track `session/resume`. This Task
marks the failure as a Lost Rollout, names its step and keeps the bounded
detail. It answers the Backlog Entry "A lost Codex rollout fails the Task and
spends a retry" of 2026-10-06. It is verifiable on its own at the runner seam
with stdout recorded from acpx 0.19.4.

## Requirements

1. MUST add `LostRolloutPhrase`, `LostRollout` and `DescribeLostRollout` with
   the shapes of `_techspec.md` → Interfaces in a new
   `internal/agent/lost_rollout.go`. `ProtocolFailure` gains the fields
   `LostRollout bool` and `Detail string`.
2. MUST read the JSON-RPC error's `data.details` only when it is a string, and
   only to decide and keep a Lost Rollout, per Invariants 1 and 3. No other
   field of `data` is read, and a non-matching error keeps an empty `Detail`.
3. MUST add `session/resume` to the tracked protocol requests, per
   Invariant 2.
4. MUST append `; lost rollout: <detail>` after `(JSON-RPC <code>)` in the
   protocol text of `BatchFailureError` and `SelectionFailureError` for a Lost
   Rollout, per API Contract 1. `Reason`, `Err`, `Stderr`, the exit mapping and
   the batch or selection classification stay as they are (Invariant 3).
5. MUST make `DescribeLostRollout` behave per Invariant 4. It takes `Runtime`
   from a `SelectionFailureError` in the chain when one is present.
6. MUST add the tests named in Verification to a new
   `internal/agent/lost_rollout_test.go`. They drive `RunPrompt` through
   `runFakeACPXPrompt` with the stdout in the appendix of
   `references/2026-10-06-lost-rollout-investigation.md`.
   - A loss during `session/resume` is a `SelectionFailureError` described as
     step `session/resume`, `PromptSent` false, with the phrase in its detail
     and error text.
   - A loss during `session/prompt` is a `BatchFailureError` described as step
     `session/prompt`, `PromptSent` true.
   - An `Internal error` whose `data.details` lacks the phrase, the phrase
     under code `-32602`, and a `data.details` that is an object are not Lost
     Rollouts. Their text is byte-identical to today's.
   - The phrase in `message` with code `-32600` is a Lost Rollout.
   - A detail longer than `ProtocolFailureMessageLimit` and spread over lines
     is cut to one line within the limit without a broken rune.
   - `session/resume` answered with a non-lost error names `session/resume`.
7. MUST NOT change the daemon, delivery or CLI packages, the acpx arguments of
   any command, or the existing tests in `internal/agent`.

## Subtasks

- [ ] Add the matcher, the description function and the two fields.
- [ ] Read `data.details` for the match and track `session/resume`.
- [ ] Extend the protocol text for a Lost Rollout.
- [ ] Add the runner tests with the recorded stdout.

## Acceptance Criteria

- [ ] Both recorded losses are described as Lost Rollouts with their step and
      the phrase.
- [ ] No other error line changes text or classification.
- [ ] Every existing `internal/agent` test passes unchanged.

## Context

- creates: `internal/agent/lost_rollout.go`
- creates: `internal/agent/lost_rollout_test.go`
- interface: `internal/agent/protocol_failure.go`
- interface: `internal/agent/acpx_runner.go`
- instruction: `internal/agent/acpx_runner_test.go`
- instruction: `internal/agent/protocol_failure_test.go`
- instruction: `docs/specs/0240-a-lost-rollout-is-infrastructure/references/2026-10-06-lost-rollout-investigation.md`
- instruction: `docs/adr/0245-a-lost-rollout-is-runtime-infrastructure.md`
- instruction: `docs/adr/0242-a-blocked-review-names-the-protocol-step-and-retries-once-before-the-prompt.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestLostRolloutIsNamedDuringResume|TestLostRolloutIsNamedMidTurn|TestOtherErrorsAreNotLostRollouts|TestLostRolloutDetailIsBounded|TestSessionResumeIsATrackedStep)$' ./internal/agent 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestLostRolloutIsNamedDuringResume TestLostRolloutIsNamedMidTurn TestOtherErrorsAreNotLostRollouts TestLostRolloutDetailIsBounded TestSessionResumeIsATrackedStep; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the five tests do not exist, so the command fails.

## References

- `_prd.md` → Core Feature 1; Success Metric 1; Goal 1; Acceptance evidence
- `_techspec.md` → Interfaces; Invariants 1 to 4; API Contract 1; Testing Approach 1; Build Order 2
- ADR-0245
- ADR-0242

## Result

Implemented the runner slice for Daemon Verification. The Task status remains
Daemon-owned; no declared Verification command was run, and no commit, push or
Pull Request was made.

- Both recorded losses: the new runner tests replay the investigation
  appendix's exact stdout through `runFakeACPXPrompt`. Resume retains a
  `SelectionFailureError`, step `session/resume`, Runtime `codex`, and
  `PromptSent: false`; mid-turn retains a `BatchFailureError`, step
  `session/prompt`, and `PromptSent: true`. Both expose the phrase through
  `DescribeLostRollout` and assert the exact error text including the suffix.
- Other error lines: negative controls assert byte-identical diagnostics for
  both selection and batch paths when details lack the phrase, the code is
  `-32602`, details are an object, the phrase is in another data field, or
  data itself is a string. Only `data.details` is decoded, and only matching
  losses keep `Detail`. The existing Reason, Err, Stderr and exit mapping are
  preserved. A non-lost resume now names its tracked step as required.
- Existing agent tests: left unchanged. The focused package check below
  exercised all existing tests alongside the new tests. Additional new
  coverage checks the phrase in `message` under `-32600`, multiline Unicode
  detail truncation, wrapped/joined failure chains, Runtime propagation, and
  cancellation/StopError exclusion.

Focused-check evidence:

- Before implementation,
  `GOCACHE=/tmp/roundfix-task02-gocache rtk proxy go test -count=1 -run TestLostRolloutIsNamedDuringResume ./internal/agent`
  exited 1: the new API and ProtocolFailure fields were absent.
- After the final code/test edits,
  `GOCACHE=/tmp/roundfix-task02-gocache rtk proxy go test -count=1 ./internal/agent`
  exited 0 (`ok roundfix/internal/agent 2.752s`).
- `rtk proxy git -c core.fsmonitor=false diff --check` exited 0.
- The first
  `GOCACHE=/tmp/roundfix-task02-gocache rtk make verify-incremental`
  exited 2. Formatting and vet passed; agent tests passed. Two CLI
  process-owner tests could not read the host process table (`operation not
  permitted`). Several suite guards also detected this Agent's Result edit
  while that check was running.
- A subsequent
  `GOCACHE=/tmp/roundfix-task02-gocache rtk make verify-incremental`
  with host process-table access and no concurrent worktree edits exited 0.
  Formatting, vet, the repository Go suite, shipped-skill checks and CLI build
  passed. This is incremental-check evidence, not Daemon Verification or a
  Task settlement verdict.

The diff is bounded to the two new agent files, the two declared agent
interfaces, and this Result. No daemon, delivery, CLI, existing test, Task
Graph, other Task file, or acpx command arguments changed. No follow-up was
identified within this runner slice. Declared Verification and Task settlement
remain with the Daemon.
