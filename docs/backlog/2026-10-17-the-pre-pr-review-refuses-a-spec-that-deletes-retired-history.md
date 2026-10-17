---
type: fix
status: open
created: 2026-10-17
spec: null
reason: null
---

# The pre-PR review refuses a Spec that deletes retired history

## Symptom

The Delivery Queue parked Spec 0272 as `review-blocked` on 2026-10-17 with
`review diff too large: 1180615 bytes after omitting 174 path(s) exceeds the
review bound of 917504 bytes`. The review never ran. Its diff carried the
deletion of 139 retired history files (158 KB), four Baseline plan
characterization goldens (270 KB), the coverage record (60 KB) and two QA
reports (97 KB).

## Where

The pre-PR review diff in `internal/cli/review.go` omits QA evidence,
upstream-managed skills and the deleted sources of a Spec archived in the
same candidate. It does not omit the deletion of a retired history file,
which ADR-0278 now makes routine, or a generated golden that a declared
regeneration command writes.

## Expected

A deletion of a retired history file and a file a sanctioned regeneration
command writes are listed as omitted, with their reason, the way QA evidence
is, so a Spec that sanitizes the history stays under the review bound. A
candidate that is still too large is reviewed in bounded segments rather than
refused.

## Evidence

Delivery of Spec 0272 on 2026-10-17. The operator reviewed the candidate in
three segments under the bound, using `roundfix review --base` in a scratch
clone, and opened the Pull Request by hand.
