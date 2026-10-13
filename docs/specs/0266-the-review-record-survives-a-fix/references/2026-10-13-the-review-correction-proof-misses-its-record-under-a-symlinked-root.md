---
type: fix
status: promoted
created: 2026-10-13
spec: 0266-the-review-record-survives-a-fix
reason: null
---

# The review correction proof misses its record under a symlinked worktree root

## Symptom

After a pre-PR review finding is fixed and disposed with `--fixed-by`,
`roundfix deliver retry` refuses the item as `corrective-spec-required` because
no review correction record is found. It happened on 0259 and 0258, and both
needed a pull request opened by hand.

## Where

`ProveReviewCorrection` reads the record under
`reviewCheckoutDir(artifactDir, workDir)` with the item's worktree path as
recorded under `~/.roundfix/worktrees` (a symbolic link), while
`roundfix review` and `roundfix review dispose` key the record by the physical
path under `/Volumes/...`.

## Expected

Both sides key the review record by the same resolved checkout path, so a
disposed and fixed finding lets the retry proceed under a symbolic-linked
`worktree.location`.

## Evidence

Operator log entries 308 and 312 (2026-10-13): PRs #512 and #514 were opened by
hand after the retry was refused.
