---
status: accepted
created_at: 2026-10-22T21:00:00Z
updated_at: 2026-10-22T21:00:00Z
deprecated_at: null
superseded_by: null
---

# A Result Status Note of a completed Task is dismissed without the validator

The Daemon sets a Task `in_progress` before the Agent turn and writes
`completed` only after the Task's Verification passes
([ADR-0057](0057-daemon-exclusively-owns-implement-task-status.md)); the
Agent writes the Task's `## Result` in between, so a Result often reports the
status as `in_progress` or Daemon-owned. The pre-PR reviewer kept raising that
sentence against the `completed` front matter, and the sealed validator of
[ADR-0196](0196-a-pre-pr-review-finding-parks-only-after-validation.md) kept
it standing on Specs 0290 and 0297, reading convention C2's "agrees with a
`completed` status" as a requirement the Result failed.

## Decision

**C2 says who writes the Result and when.** Delivery Conventions version
`roundfix/delivery-conventions/v4` states that the Agent writes a non-QA
Task's `## Result` during its turn, while the Daemon holds the Task
`in_progress`, so a Result that reports the status as `in_progress`, unchanged
or Daemon-owned states what was true when it was written and does not
contradict a `completed` status. C2 no longer says the Daemon writes the
Result; its status, carry-forward and recorded-paths clauses keep their
wording.

**The pair is dismissed mechanically.** Before the sealed validator is asked,
Roundfix dismisses a finding by validation with the rule `convention:C2` and a
fixed reason, without a model, when all of these hold at the review head: the
whole anchor lies in the finding's Task file C2 region (its front matter and
settlement sections); that Task's front matter has one `status:` line whose
value is `completed` and no `type: qa`; its `## Result` section holds a Result
Status Note, a line containing `in_progress` or `Daemon-owned` in any letter
case; and the finding's text contains `completed` and one of `in_progress`,
`in progress` or `Daemon-owned` in any letter case. Every other finding keeps
today's path: the sealed validator, failing closed.

**A review whose only findings are dismissed this way passes.** When no
finding is left to ask, the record reads `findings-dismissed` if every finding
is dismissed by validation or by the operator's evidence, as after a validator
run, and the record's `validator` stays `not-needed`.

## Considered Options

- Rewording C2 alone was rejected: the verdict stays a model's reading, which
  is what failed twice.
- Telling the Agent not to report the status in its Result was rejected as
  the fix: Results already written keep the sentence, and an Agent may still
  write it.
- Dismissing on the anchor alone, without reading the finding's text, was
  rejected: a finding anchored on a Result Status Note can name another
  failure, such as a missing test.
- Dismissing every finding anchored in a completed Task's Result was
  rejected for the same reason.

## Consequences

- The dismissal reads only the head tree and the finding's text, so the same
  candidate and answer give the same record on every host.
- A finding that names the pair and also another failure is dismissed with
  it; the reviewer is asked for one finding per list item, and a second
  failure stated alone still reaches the validator.
- A Task whose status is not `completed`, a QA Task, a Result without a
  Result Status Note and an anchor outside the C2 region are never dismissed
  this way.
