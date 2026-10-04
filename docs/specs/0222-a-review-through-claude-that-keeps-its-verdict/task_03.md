---
task: task_03
spec: 0222-a-review-through-claude-that-keeps-its-verdict
status: pending
type: backend
complexity: medium
---

# Task 03: The review classifies a verdict delivered after a refused permission and names a prompt that is too long

## Overview

With the runner keeping a read-only turn that ended after a refused
permission, `roundfix review` classifies that answer by its verdict, so a clean
answer passes and findings get ids that `review dispose` accepts, and the
record says the refusal happened. The review also names a runtime's
`Prompt is too long` answer in its reason instead of a protocol error. Both
answer the reports of 2026-10-02 in the adopted Backlog Entry
[2026-10-03](references/2026-10-03-the-review-through-claude-blocks-on-a-refused-permission.md).
This Task is verifiable on its own through the review command's tests with a
fake runner.

## Requirements

1. MUST add the record field `PermissionRefused bool` with the JSON name
   `permissionRefused`, omitted when false, and set it from
   `result.PermissionRefused` in `classifyReviewCommandResult`.
2. MUST classify the final answer of a result with `PermissionRefused` exactly
   as a clean result's answer, so findings get ids `F1..Fn`, pass anchor
   validation as today, and exit `1`, and a `No findings.` answer exits `0`;
   and MUST append ` (after the read-only session refused a permission request)`
   to the reason of a record that still ends `blocked`, as API Contract 2
   states.
3. MUST, before the runtime-failure and transport-anomaly cases and before
   verdict classification, block with the reason
   `review prompt too long: <line>` when a line of `result.Answer()`, trimmed,
   starts with `Prompt is too long`, cutting the line to at most 512 bytes on a
   rune boundary; admission, timeout and Spec-read failures keep their
   precedence, as API Contract 3 states.
4. MUST add the tests named in Verification to the new file
   `internal/cli/review_permission_test.go`, using `newReviewCommandFixture`
   with the `claude` provider and its review profile, and the fake runner:
   - a refused permission with a findings answer: exit `1`, `F1` and `F2`,
     `permissionRefused: true` in stdout and in the persisted record, then
     `roundfix review dispose F1 --dismiss --evidence <text>` exits `0`;
   - a refused permission with `No findings.`: exit `0`, outcome `reviewed`;
   - a refused permission with an answer that has no verdict: exit `2` and the
     suffix in the reason;
   - an answer whose first line is the oraculum message
     `Prompt is too long · the request is ~1294791 tokens (limit 1000000) ...`,
     once with a `BatchFailureError` whose reason is `agent/protocol error`
     and once with a parsed result: exit `2` and the reason
     `review prompt too long: Prompt is too long · the request is ~1294791 tokens (limit 1000000)`
     followed by the rest of the line.
5. MUST keep `TestReviewCommandBlocksOnTransportAnomaly` unchanged and
   passing.
6. MUST NOT change the review prompt, the verdict grammar, finding validation
   or the lineage rules.

## Subtasks

- [ ] Carry the refusal into the record and classify by verdict.
- [ ] Suffix a blocked reason after a refusal.
- [ ] Name a prompt-too-long answer.
- [ ] Add the review tests, including dispose.

## Acceptance Criteria

- [ ] A refused permission no longer blocks a delivered verdict, and its
      findings can be disposed.
- [ ] A prompt-too-long answer is named in the reason after a runtime failure
      and after a parsed result.
- [ ] A transport anomaly from any other exit still blocks.

## Context

- interface: `internal/cli/review.go`
- creates: `internal/cli/review_permission_test.go`
- instruction: `internal/cli/review_test.go`
- instruction: `.agents/skills/roundfix/references/review.md`
- instruction: `docs/adr/0174-a-pre-pr-review-reads-the-final-message-and-keeps-a-record-per-checkout.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReviewCommandClassifiesFindingsAfterARefusedPermission|TestReviewCommandPassesNoFindingsAfterARefusedPermission|TestReviewCommandNamesTheRefusalWhenTheAnswerHasNoVerdict|TestReviewCommandNamesAPromptThatIsTooLong|TestReviewCommandBlocksOnTransportAnomaly)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestReviewCommandClassifiesFindingsAfterARefusedPermission TestReviewCommandPassesNoFindingsAfterARefusedPermission TestReviewCommandNamesTheRefusalWhenTheAnswerHasNoVerdict TestReviewCommandNamesAPromptThatIsTooLong TestReviewCommandBlocksOnTransportAnomaly; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the four new tests do not exist, so the command fails.

## References

- `_prd.md` → Core Feature 2; Core Feature 3; User Story 1; User Story 2; User Story 3; Success Metric 2; Success Metric 3; Goal 1; Goal 2
- `_techspec.md` → The refused permission; The prompt-too-long reason; API Contract 2; API Contract 3; Testing Approach 2; Build Order 3
- ADR-0227
- ADR-0174
