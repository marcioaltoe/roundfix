---
status: accepted
created_at: 2026-10-13T00:00:00Z
updated_at: 2026-10-18T13:00:00Z
deprecated_at: null
superseded_by: null
---

# A review record survives a fix, and a correction inside the archived Spec's declared paths returns to review

On 2026-10-13 two Delivery Retries of `corrective-spec-required` items were
refused although the operator had fixed and disposed every finding, and both
Pull Requests were opened by hand. `roundfix review` names its record
directory from the checkout root Git reports, which is the resolved path,
while the retry's correction proof named it from the item worktree path as the
Delivery Queue recorded it, under a symbolic-linked `worktree.location`, so
the proof looked in a directory nobody wrote. The operator then ran review
round 2 at the fixed head, and that round replaced the only record of the
parked candidate's findings. Both fixes also changed code, which ADR-0233's
retry rule refuses whatever the records say: only paths under the archived
Specs the blocker names may change. Asked whether to keep that rule, the
maintainer answered "Relaxar, com limites" (2026-10-13).

**The record is keyed by the resolved checkout.** The directory of a checkout's
Pre-PR Review Record is named from the checkout path after
`config.CanonicalizeResolvedPath`, which resolves the longest existing prefix
(ADR-0268), so every spelling of one checkout reads and writes one record.
Every comparison of a record's or a disposition's repository with a checkout,
or with another record, compares the two resolved paths. Records and ledger
lines keep the repository they recorded; nothing is rewritten or migrated.
Because `review` always recorded the path Git resolves, an existing record
keeps its directory.

**Each reviewed head keeps its record.** Beside the current record, a review
that ends in a verdict for a head (`reviewed`, `findings`,
`findings-dismissed` or `ceiling-closed`), and is not a reuse, also writes a
Reviewed Head Record at `heads/<head commit>.json` in the same checkout
directory, with the same content. A later round at another head never touches
it; a reused, blocked or omitted record never writes one. The current record
keeps its role: the next round, reuse and `roundfix review dispose` read it as
before. The correction proof reads the Reviewed Head Record of the parked
candidate, or the current record when it names the candidate and no head
record exists, which covers items parked before this decision.

**A correction inside the archived Spec's declared paths returns to review.**
This refines ADR-0233's third rule. A Delivery Retry of a
`corrective-spec-required` item whose head moved returns it to `reviewing`
when Git proves the head descends from the parked candidate, every standing
finding of the review recorded at that candidate has one disposition (fixed by
a commit the head contains, or dismissed with evidence), and every path
changed from the candidate to the head lies either under an archived Spec the
blocker names, as before, or among that Spec's declared paths: the
`creates:` and `interface:` entries and the `## Recorded paths` of its Task
files. The Task files are read at the parked candidate for a Spec archived as
a folder, and at the Archive Record's `source_revision` for one archived as a
record, recovered as any reader recovers it (ADR-0261); when they cannot be
read, the retry refuses. Paths match exactly. Review round 2 of the same
Reviewer Lineage and the repository gate and required checks stay mandatory.
Anything else still requires a corrective Spec (ADR-0165), and the refusal
names the first changed path outside both sets.

**Rejected.** Reading the round-1 findings from the round-2 record's lineage
(`previousFindings`) loses them whenever the next review starts a new lineage
(a moved merge base or another provider) or closes the ceiling, which replaces
the lineage. Numbering records by round collides, because a new lineage starts
again at round 1, and the proof asks by head, not by round. Recording the
resolved item path in the Delivery Queue would need a Run Database migration
and leave every parked item wrong, the alternative ADR-0268 rejected for Run
paths. A `review dispose --head` flag to dispose an earlier round's finding
adds a public flag the observed flow does not need. Accepting any code
correction would let a retry carry a change outside what the archived Spec
was authorized and reviewed to touch; reading the declared paths from the
item's current tree instead of the recorded revision would let the correction
widen its own bound.

## Consequences

- Under a symbolic-linked `worktree.location`, a fixed and disposed finding
  lets the retry return the item to `reviewing` when the other proofs hold,
  also after review round 2.
- A code fix to a file the archived Spec's Tasks created, changed or recorded
  no longer needs a hand-opened Pull Request; it is reviewed again in round 2
  and gated as before. A fix to any other file still parks for a corrective
  Spec, and the retry names that file.
- ADR-0174's per-checkout record stands; two spellings of one checkout are one
  checkout, and a reader still never takes another checkout's record.
- ADR-0197's lineage is unchanged; the ceiling still keeps the round-2 record
  in place for `review dispose`.
- A checkout directory gains one small file per reviewed head. Nothing reclaims
  them yet; they live in the Artifact Directory with the rest of the review
  artifacts.

**Amended by [ADR-0286](0286-a-fixed-round-2-finding-gets-one-correction-review-of-its-fix.md) (2026-10-18).** A retry of an item parked by round-2 findings whose every standing finding is disposed as fixed between the parked candidate and the item head skips the declared-path bound, because the Correction Review reads the whole fix; every other retry keeps it.
