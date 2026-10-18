---
status: accepted
created_at: 2026-10-18T13:00:00Z
updated_at: 2026-10-18T13:00:00Z
deprecated_at: null
superseded_by: null
---

# A fixed round-2 finding gets one Correction Review of its fix

On 2026-10-17 the pre-PR review of Spec 0276's delivery candidate reached round
2 and found two real defects in the operator's CI fix: a retry with a sleep on
`ETXTBSY`, a timing allowance ADR-0271 rejects, and a cancelled wait that
returned a stale error. The Delivery Queue parked the item
`corrective-spec-required`. The operator replaced the fix structurally, wrote
the item binary copy under the fork lock, and found no sanctioned way forward.
A Delivery Retry would have refused the correction, because it changed a test
file outside the paths 0276's Tasks declared (ADR-0272). Had the retry
accepted it, the next review would have closed the lineage as `ceiling-closed`
without any reviewer reading the new code (ADR-0197). So the operator reviewed
the delta by hand in a scratch clone with `roundfix review --base`. That review
found one more defect, which was fixed and then reviewed clean. The operator
ran `make verify-changed` and merged pull request #568 by hand with the
maintainer's approval.

## Decision

**A fixed round-2 finding gets one Correction Review.** A Correction Review
runs when all of these hold:

- Round 2 of a Reviewer Lineage ended with outcome `findings`.
- Every standing round-2 finding has exactly one valid disposition, and it is
  `fixed` by a commit that descends from the round-2 head and that the current
  head contains.
- The current head descends from the round-2 head.

Under those conditions the next review of that lineage calls the reviewer once
more instead of closing at the ceiling. It reviews only the delta from the
round-2 head to the current head. Its prompt carries the round-2 findings and
the commit that fixed each one. Like round 2, it is never segmented, a delta
above the review bound blocks it, and it opens a fresh Agent Session.

The record keeps round 2 as the lineage's round. It adds `correction: true`,
names the round-2 head as `previousHead`, and lists the round-2 findings and
their fixing dispositions.

**The outcomes follow round 2's.** A clean Correction Review lets the candidate
proceed exactly as a clean round 2 does. Findings park it exactly as round-2
findings do: `corrective-spec-required` when the findings name an archived
Spec, `review-findings` otherwise. A Correction Review that blocks repeats as a
Correction Review, as a blocked round repeats its own round.

**There is never a second Correction Review.** The next review of a lineage
whose last record is a Correction Review reaches the ceiling as round 2 does
today. That review either closes `ceiling-closed` on valid dispositions or
blocks and names the findings to dispose.

**The retry admits the fix of a round-2 park whatever paths it touches.**
`roundfix deliver retry` on an item parked `corrective-spec-required` accepts a
moved head when the parked record is a round-2 review, not itself a Correction
Review, and every standing finding has exactly one valid `fixed` disposition
between the parked candidate and the item head. It then skips the
declared-path bound of ADR-0272, because the Correction Review reads every line
the fix changed. Every other retry keeps the bound.

**The ceiling's purpose stays.** The ceiling exists so that no candidate loops
through reviews without end. A Correction Review is bounded: it runs at most
once per lineage, it never starts another round, and only a disposition with a
fixing commit opens it.

## Considered Options

- Allowing a third round for any disposition was rejected: a dismissal is the
  operator's verdict, not a fix, so there is nothing new to review.
- Leaving the ceiling unchanged and asking the operator for a corrective Spec
  was rejected. The defect lay in delivery code written after archive, and a
  corrective Spec for a two-commit fix costs a full planning, Run and QA cycle
  that adds no review the Correction Review lacks.
- Keeping ADR-0272's declared-path bound for the fix of a round-2 park was
  rejected. A structural fix needs its own test, and the measured case changed
  `internal/cli/deliver_item_binary_fork_linux_test.go`, which no Task of Spec
  0276 declared. With the bound, the Correction Review would be unreachable in
  the case it exists for.
- Naming the review "round 3" was rejected. ADR-0197 and the glossary use
  "third round" for what never calls the reviewer, and the Correction Review is
  not a round of discovery.

## Consequences

- `ceiling-closed` now covers only lineages where the operator dismissed at
  least one round-2 finding, or where a Correction Review already ran. A
  round-2 fix whose dispositions mix `fixed` and `dismissed` still closes at
  the ceiling unreviewed, as before.
- The Delivery Queue's repository gate, required checks and merge order are
  unchanged. A Correction Review adds one reviewer call before them.
- This amends ADR-0197's ceiling and ADR-0272's declared-path bound for exactly
  this case. ADR-0165's corrective-Spec park and ADR-0196's validation of
  findings are unchanged.
