---
status: accepted
created_at: 2026-10-03T00:00:00Z
updated_at: 2026-10-03T00:00:00Z
deprecated_at: null
superseded_by: null
---

# A delivery keeps its work across archive, requeue and review

Between 2026-10-01 and 2026-10-03 the operator finished five deliveries by
hand whose work was already correct. A test read its own Spec's active
TechSpec and failed once the archive commit moved it. Another active Spec's
Task Context named that active path and failed the corpus check on the
delivery's Pull Request. A Delivery Retry refused once a correction committed
after the archive moved the item head. `deliver start` for a Spec whose item
branch held completed Tasks created a new item from the default branch, twice.
And the pre-PR review raised, for Specs 0215 and 0218, that an archive made
with an authorized QA Archive Override had a failed QA Task.

Five rules settle them:

- A file outside the Spec roots and the history root that names an active
  Spec's directory, other than Markdown, is a Spec Consistency error for that
  Spec. The Daemon's Settlement Check therefore returns it to the Task that
  wrote it, and the Archive Command refuses to archive a Spec while such a
  file exists. Roundfix never rewrites code to follow an archive.
- A Task Context path under an active Spec's directory resolves to the same
  path under the archive root once that Spec is archived. Another active
  Spec is never edited by an archive, so its declaration stays valid without
  a rewrite.
- A Delivery Retry of an archived item whose head descends from its last
  candidate head records that head as a new candidate and returns the item to
  review. A park for a corrective Spec still refuses a moved head, because
  an archived Spec never absorbs a review finding.
- `deliver start` continues the one item branch of the Spec that holds
  commits the default branch lacks, and refuses, naming each branch, when
  there is more than one.
- An archive whose `_prd.md` records `qa_override: true` with an approval
  source and a reason is a Delivery Convention. The review is told so, and
  the validator may dismiss a finding that only restates it. An archive
  without that record is still a finding.

Two alternatives were rejected. Having the archive rewrite other Specs' paths
would make the archive commit more than the Spec's move, which the queue
treats as a stale review. Rewriting a test's path would hide that the test
reads a document that may be deleted.

## Consequences

The archive commit stays an exact move. A post-archive correction is
reviewed again before publication, so no unreviewed change reaches a Pull
Request. A continued item branch keeps every Task it completed and may lag
the default branch; the existing conflict handling deals with the lag.
Markdown outside the Spec roots that names an active Spec's path, and the
relative links inside a Spec that the move breaks, are not covered.
