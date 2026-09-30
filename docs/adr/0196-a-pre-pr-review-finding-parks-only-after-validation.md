---
status: accepted
created_at: 2026-09-30T00:00:00Z
updated_at: 2026-09-30T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A pre-PR review finding parks only after validation

The Pre-PR Review Command parked a delivery on every finding the reviewer
wrote. Between 2026-09-29 and 2026-09-30 the operator dismissed 11 of the 34
findings recorded for this repository with evidence. Five of the eleven
described the delivery's own designed order: a QA Report that cannot name the
commit that records it, a Task status the Daemon wrote, and a planning
candidate whose Tasks are pending. The reviewer was never told that order.

Now the review prompt carries the Delivery Conventions, a closed and numbered
list of what a delivery writes by design, each with the paths it owns. Each
finding must open with a `path:line` anchor and state the failure it causes.
Roundfix then validates every finding before the verdict decides anything:

- A finding whose anchor names no line the candidate diff shows is dismissed
  as `unanchored`, as a pull request review comment must sit on the diff.
- A sealed validator prompt, with no tools, may dismiss an anchored finding
  as restating one named convention, and only when the anchor lies in that
  convention's paths. It may dismiss a finding as stating no failure only
  when the finding carries no failure clause.
- Every other finding stands, and only a standing finding parks.

A dismissed finding stays in the record, with its rule and its reason. It is
never dropped.

## Consequences

Validation fails closed. When the validator cannot run, times out, uses a tool
or answers outside its grammar, every anchored finding stands, and the record
says why. An answer in which no finding carries a parseable anchor is blocked,
because the reviewer did not follow the grammar. A reviewer that ignores it
therefore cannot pass by having everything dismissed.

A record whose findings are all dismissed, by validation or by the operator's
evidence, reads `findings-dismissed`, as before. The operator cannot dispose a
finding that validation already dismissed.

This decision refines
[ADR-0153](0153-pre-pr-review-is-an-explicit-provider-policy.md), whose
enabled review still requires a complete result for the current candidate. It
also refines
[ADR-0174](0174-a-pre-pr-review-reads-the-final-message-and-keeps-a-record-per-checkout.md),
whose final message is the answer the grammar reads. And it refines
[ADR-0169](0169-the-pre-pr-review-diffs-the-candidate-from-its-merge-base.md),
whose merge-base diff is the diff an anchor must fall in. It supersedes none of
them.

The TypeSafe judgment model was considered as the validator and rejected. The
maintainer's data grant covers Spec artifacts, not source code or diffs, and
that model never decides a gate.
