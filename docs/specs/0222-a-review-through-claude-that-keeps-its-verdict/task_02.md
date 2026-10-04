---
task: task_02
spec: 0222-a-review-through-claude-that-keeps-its-verdict
status: pending
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
