---
status: done
absorbed_by: 2026-08-06-rollup-review-and-delivery-convergence.md
created_at: 2026-08-14
updated_at: 2026-08-26
kind: finding
---

# A Review retires on whatever the object store happens to hold (2026-08-14)

Running the migration Spec 0094 shipped, against the repository that shipped it,
produced three different answers for the same fifty orphan Review Artifacts
within one session. Nothing about the pull requests changed between them. What
changed was the local Git object store.

ADR-0123 settles that an orphan Review Artifact retires on local Git
reachability: an ancestor of the default branch retires as merged, an
unreachable head whose branch is gone retires as abandoned, and a reachable
non-ancestor stays live. The decision was deliberate and its trade-off was
recorded — reachability can lag the provider, so an undecidable head stays live
rather than retiring. What the record did not anticipate is that the answer is
not a property of the pull request at all.

## 1. The same head classifies three ways in one session

- Symptom / evidence: measured 2026-08-14 on `roundfix` at `217d78a`, against
  head `4a03df27` recorded by `docs/specs/_reviews/pr-124`.

  | Local state | `rev-parse` | Classification | Outcome |
  | --- | --- | --- | --- |
  | Clone as found, no pull request refs | fails | `undecidable` | retained |
  | After `git fetch origin '+refs/pull/*/head:…'` | resolves | `live` | retained |
  | Refs deleted, object still in the store | resolves | finished, abandoned | migrated |

  All fifty folders moved in the third state and none in the first two. The
  third state is the most accidental of the three: deleting a ref does not delete
  its object, so a head with no ref pointing at it reads as "reachable, branch
  gone", which the rule calls abandoned.

- Root cause: the rule asks the object store a question the object store cannot
  answer. Whether a commit is present locally depends on the machine's fetch
  history, its garbage collection, and whether anyone ever fetched pull request
  refs — none of which says anything about whether the pull request finished.
  Two clones of one repository, at one commit, answer differently.

- Action / suggestion: make the answer a property of the artifact rather than of
  the machine. The candidates worth weighing are recording the provider's answer
  once, at the moment the Round settles, so later runs read a recorded fact
  rather than re-deriving it; or treating an orphan review whose pull request
  number is closed as finished on one authenticated read, accepting that the
  migration then needs a network for that family alone; or dropping liveness
  from the migration entirely and moving the orphan root as a legacy layout,
  which is what the underscore rename below already argues for.

## 2. A squash-merging repository can never retire a review by ancestry

- Symptom / evidence: with pull request refs fetched, `4a03df27` resolves and is
  **not** an ancestor of `origin/main`, because pull request 124 was
  squash-merged. This repository squash-merges every pull request. So the
  "ancestor of the default branch" path — the rule's primary way to prove a
  review finished — can never fire here.
- Root cause: a squash merge creates a new commit and discards the branch's own
  history, so ancestry cannot prove integration for any repository using it.
  ADR-0123 names this case and treats it as lag; measured, it is not lag but
  total absence of that path.
- Action / suggestion: whatever replaces the rule must not rest on ancestry
  alone. Content equivalence, a recorded settlement fact, or the provider's
  state each survive a squash merge; ancestry does not.

## 3. The underscored orphan root is never renamed

- Symptom / evidence: Spec 0094 moved the live orphan Review Artifact root from
  `docs/specs/_reviews/` to `docs/specs/reviews/` for new writes. The fifty
  existing folders stayed at the underscored path, and nothing migrates them
  there: discovery evaluates them for retirement, not for the rename. In the
  first two states above they would have sat at a path the resolver no longer
  writes to, indefinitely.
- Root cause: the underscore change is a resolver change, and the layout
  discovery that migrates legacy shapes does not treat the underscored orphan
  root as one of them.
- Action / suggestion: treat `docs/specs/_reviews/` as a legacy layout to
  migrate, independent of any liveness question. A rename needs no judgement
  about whether a review finished.

## What worked — keep

The relocation machinery itself is sound. All 501 files moved with each
destination verified against the content identity recorded before the move, the
emptied source directories were removed, and the second run reported nothing
left to do. The defect is in deciding *which* files to move, not in moving them.

---

Triage 2026-08-26: deferred out of the active queue. See docs/workflow/2026-08-26-triage-the-queue-earns-its-tokens.md.
