---
task: task_03
spec: 0094-one-history-root-under-docs
status: pending # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: medium
---

# Task 03: Decide a review's liveness from local Git alone

## Overview

An orphan Review Artifact is finished when its Pull Request is, and that fact
lives outside the repository. This slice answers it from local Git instead — the
head a round recorded, against the default branch and the ref namespace — and
answers three ways rather than two, because a head the repository cannot decide
must stay live. Verifiable on its own against fixture repositories, with no
network and no caller yet.

## Requirements

1. MUST report a review finished when the head its newest round recorded is an
   ancestor of the default branch.
2. MUST report a review finished when that head is unreachable and the branch it
   was pushed to no longer exists.
3. MUST report a review live when that head is reachable and is not an ancestor of
   the default branch.
4. MUST report a review undecidable when no round records a head, when the
   metadata cannot be read, or when Git cannot answer, and MUST treat undecidable
   as live for every downstream purpose.
5. MUST decide from the newest round's recorded head when a review holds several
   rounds with different heads.
6. MUST NOT contact the hosting provider, read a credential, or make any network
   request.
7. MUST name the reason alongside the answer, so a caller can report why a review
   was left live.

## Subtasks

- [ ] Read the recorded head from a review's newest round metadata.
- [ ] Answer from local ancestry, reachability, and branch existence.
- [ ] Resolve every undecidable case to live with its reason.
- [ ] Cover each answer with a fixture repository.

## Acceptance Criteria

- [ ] A head merged into the default branch reports finished.
- [ ] A head whose branch was deleted and which is unreachable reports finished.
- [ ] A head on a live branch, not an ancestor, reports live.
- [ ] A review whose metadata records no head reports undecidable, and a test
      fails if that case ever reports finished.
- [ ] A review with several rounds decides from the newest round's head.
- [ ] The reason is populated for every answer.
- [ ] No test in this slice requires a network, and the classification makes no
      outbound call.

## Verification

- `go test -count=1 ./internal/spec -run 'ReviewLiveness|ClassifyReview' -v > /tmp/0094-task-03.log 2>&1; s=$?; grep -q '^--- PASS: .*ReviewLiveness\|^--- PASS: .*ClassifyReview' /tmp/0094-task-03.log || { cat /tmp/0094-task-03.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing liveness tests; fails when the named tests do not exist.
- `! grep -qi 'no tests to run' /tmp/0094-task-03.log` — expected: exits 0, refusing a vacuous run.
- `! grep -rn 'runGH\|gh api\|net/http' internal/spec` — expected: exits 0, proving the classification takes no provider or network dependency.
- `go build -buildvcs=false ./...` — expected: exits 0.

## Context

- interface: `internal/gittest/gittest.go`
- interface: `internal/rounds/rounds.go`

## References

`_techspec.md` → Build Order 3; Interfaces: `ClassifyReview`, `ReviewLiveness`;
Testing Approach: review liveness; Integration Points. `_prd.md` → Core Feature 4;
Goal 4; User Story 4. ADR-0123.
