---
type: fix
status: promoted
created: 2026-10-22
spec: 0304-the-review-drops-a-result-status-note-and-a-stale-fingerprint-names-its-command
reason: null
---

# The pre-PR review parks on a Result status note that convention C2 allows

## Symptom

The pre-PR review raised, and its sealed validator kept as `stands`, a finding
that a Task's `## Result` reports the status `in_progress` or Daemon-owned
while the Task's front matter says `completed`. The operator had to dismiss it
with evidence citing Delivery Convention C2 before delivery could continue.

- Spec 0290, review head `61dec34e`: finding F1 anchored at
  `docs/specs/0290-the-pull-request-record-commit-fits-commitlint/task_01.md:115`,
  "Failure: the Result says the Daemon-owned status remains `in_progress`,
  while the same candidate changes the task frontmatter to `completed`,
  leaving contradictory settlement evidence." Validator verdict `stands`:
  "The Result says the Daemon-owned status remains `in_progress`, but the
  task frontmatter is `completed`, contradicting C2's required agreement with
  a completed status."
- Spec 0297, review head `1aa6ae3a`, 2026-10-22: finding F1 anchored at
  `docs/specs/0297-a-corrective-task-after-a-completed-gate-passes-the-wave-collision-check/task_02.md:127`,
  "Failure: the Result says Task status remains `in_progress` while the
  candidate frontmatter sets `status: completed`, leaving the settlement
  record internally contradictory and misreporting daemon state." Validator
  verdict `stands`: "The Result reports status as in_progress while the
  candidate frontmatter records completed, contradicting the required
  completed-status agreement."

Both were dismissed by the operator with evidence citing C2. The same pair was
raised earlier on Specs 0199, 0202 and 0216 (2026-09-30 to 2026-10-02).

## Where

`roundfix review` (the Pre-PR Review Command): `internal/cli/review_conventions.go`
`deliveryConventions` (C2, `roundfix/delivery-conventions/v3`) and
`internal/cli/review_convention_validator.go` `validateReviewConventions`.
The Agent writes the sentence: the Daemon sets the Task `in_progress` before
the Agent turn and the Agent prompt says "The Daemon owns Task status during
Implement; do not edit the task file's status field.", so the
Result reports the status the Task had when it was written; the Daemon writes
`completed` after the Task's Verification passes.

## Expected

A Result that reports the status the Task had while the Agent wrote it, in a
Task the Daemon settled `completed`, is not reported, or is dismissed by
validation with a recorded reason, without depending on how a model reads the
convention. A real contradiction, such as a Task whose status is not
`completed`, a QA Task, or a finding about another failure, still stands.

## Evidence

The two review records under the operator's Roundfix artifact root (finding
texts and validator reasons quoted above), and the Task files at the Archive
Records' `source_revision` of Specs 0290 and 0297.
