---
task: task_02
spec: 0222-a-review-through-claude-that-keeps-its-verdict
status: completed
type: backend
complexity: medium
---

# Task 02: The ACPX Runner keeps a read-only turn that ended after a refused permission

## Overview

acpx exits `5` when every permission request of a turn was denied, and a
read-only Agent Session denies every request that is not a read. The ACPX
Runner today reports that exit, after a parsed result, as a transport anomaly,
which makes the review discard a delivered verdict (finding of 2026-10-02 in
the adopted Backlog Entry
[2026-10-03](references/2026-10-03-the-review-through-claude-blocks-on-a-refused-permission.md)).
This Task keeps such a turn's result and marks the refusal. It is verifiable on
its own through the runner's fake-acpx tests.

## Requirements

1. MUST add the field `PermissionRefused bool` to `agent.ExecuteResult`.
2. MUST, in `ACPXRunner.RunPrompt`, when `cmd.Wait` reports exit code `5`, the
   stream parsed a prompt result whose stop reason is `end_turn`, the request's
   access is `SessionAccessReadOnly` and the request is not inert: publish the
   existing `acpxPermissionDeniedStatus` Run Event, set `PermissionRefused`,
   leave `TransportAnomaly` empty, and return the result with a nil error, as
   API Contract 1 states.
3. MUST keep today's behavior for every other case: a parsed result with any
   exit other than `130` becomes a transport anomaly (including exit `5` with a
   stop reason other than `end_turn`, with read-write access, or with an inert
   request), exit `130` stays a stop, and an exit without a parsed result is
   mapped as before.
4. MUST let `runFakeACPXPrompt` set the request's access and inert flag,
   defaulting to the values the existing cases use, so no existing case
   changes its outcome.
5. MUST add the tests named in Verification to
   `internal/agent/acpx_runner_test.go`: the first asserts API Contract 1 for a
   read-only `end_turn` turn with exit `5`, including the Run Event; the second
   has one case each for exit `1` read-only `end_turn`, exit `5` read-only
   `max_tokens`, exit `5` read-write `end_turn` and exit `5` inert `end_turn`,
   each asserting a non-empty transport anomaly naming its exit code and
   `PermissionRefused == false`.
6. MUST NOT change the daemon packages or the acpx arguments of any session.

## Subtasks

- [ ] Add the result field.
- [ ] Classify the read-only refused-permission exit before the anomaly case.
- [ ] Extend the fake-acpx harness with access and inert.
- [ ] Add the two tests.

## Acceptance Criteria

- [ ] A read-only `end_turn` turn with exit `5` returns its result, no error,
      no anomaly, the refusal marked and the permission-denied Run Event.
- [ ] Every other exit after a parsed result is still a transport anomaly, and
      the existing runner tests pass unchanged.

## Context

- interface: `internal/agent/agent.go`
- interface: `internal/agent/acpx_runner.go`
- interface: `internal/agent/acpx_runner_test.go`
- instruction: `docs/adr/0020-parsed-prompt-result-outranks-acpx-exit-code.md`
- instruction: `docs/adr/0227-a-review-through-claude-keeps-its-verdict-and-fits-its-window.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestACPXRunPromptKeepsAReadOnlyTurnThatRefusedAPermission|TestACPXRunPromptKeepsTheAnomalyOutsideAReadOnlyRefusal)$" ./internal/agent 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestACPXRunPromptKeepsAReadOnlyTurnThatRefusedAPermission TestACPXRunPromptKeepsTheAnomalyOutsideAReadOnlyRefusal; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the two tests do not exist, so the command fails.

## References

- `_prd.md` → Core Feature 1; User Story 1; Success Metric 1; Goal 1; Goal 4
- `_techspec.md` → The refused permission; API Contract 1; Testing Approach 1; Build Order 2
- ADR-0227
- ADR-0020

## Result

Implemented the runner slice. `ExecuteResult.PermissionRefused` marks a
non-inert read-only turn with a parsed `end_turn` result and acpx exit `5`.
The runner publishes `acpxPermissionDeniedStatus`, preserves the answer and
stream output, and returns no error or transport anomaly. Event publication
errors still propagate. Other exits retain the existing classification.

The fake-acpx harness now accepts access and inert settings; their zero values
preserve the existing read-write, non-inert defaults. Its argument parser also
recognizes the existing `--approve-reads` flag so read-only prompts reach the
fake process's prompt handler. Production session arguments are unchanged.

Focused implementation evidence:

- Before the runner change,
  `GOCACHE=/tmp/roundfix-task02-gocache rtk proxy go test -count=1 -run '^TestACPXRunPromptKeeps' ./internal/agent`
  exited `1`: the read-only refusal test observed the existing exit-5
  transport anomaly. The four neighboring cases passed.
- After the change,
  `GOCACHE=/tmp/roundfix-task02-gocache rtk proxy go test -count=1 -v -run '^TestACPX(RunPrompt|PromptExitClassificationMatrix|ExitCodeMapping)' ./internal/agent`
  exited `0`.
- Acceptance criterion 1: `TestACPXRunPromptKeepsAReadOnlyTurnThatRefusedAPermission`
  passed, asserting the retained answer, output and `end_turn`, nil error,
  empty anomaly, refusal flag and permission-denied Run Event.
- Acceptance criterion 2: `TestACPXRunPromptKeepsTheAnomalyOutsideAReadOnlyRefusal`
  passed all four required cases, asserting exit-bearing anomalies and a false
  refusal flag. Existing prompt, exit-classification and exit-mapping tests
  also passed, including parsed exit `130` and exit `5` without a result.
- `rtk proxy git -c core.fsmonitor=false diff --check` exited `0`.
- `GOCACHE=/tmp/roundfix-task02-gocache rtk make verify-incremental`
  initially exited `2`: the sandbox denied process-table reads and loopback
  binds, and editing this Result during the run triggered suiteguard's
  repository-change detector. The same command rerun with escalated access
  and no concurrent worktree edits exited `0`, covering formatting, vet,
  repository Go tests, skill checks and the CLI build. Existing agent tests
  passed in the first run and were reused from Go's cache in the second.

The authored Verification command was not run; Verification and Task settlement
remain Daemon-owned. No daemon package, other Task file or Task Graph changed.

## Carry-forward provenance

- Source Run: `run_20261004T145620Z_ee61ac9af5ab0ce2`
- Source commit: `b5fce724941c34472f27dd9972d64c6a3877a6cd`
