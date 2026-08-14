---
task: task_03
spec: 0094-one-history-root-under-docs
status: completed # pending | in_progress | completed | failed — only implement-task changes this
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
- `grep -q 'func ClassifyReview' internal/spec/*.go && ! grep -rn 'runGH\|gh api\|net/http' internal/spec` — expected: exits 0, proving the classification exists and takes no provider or network dependency. The two clauses are one command on purpose: the negative guard alone passes on a tree where nothing was written, so it proves nothing until it is anchored to the code it guards.
- `grep -q 'ReviewUndecidable' internal/spec/*.go` — expected: exits 0, proving the third answer ADR-0123 requires exists rather than a boolean that can only say finished or live.

## Context

- interface: `internal/gittest/gittest.go`
- interface: `internal/rounds/rounds.go`

## References

`_techspec.md` → Build Order 3; Interfaces: `ClassifyReview`, `ReviewLiveness`;
Testing Approach: review liveness; Integration Points. `_prd.md` → Core Feature 4;
Goal 4; User Story 4. ADR-0123.

## Result

### Implementation

- `ClassifyReview` returns `finished`, `live`, or `undecidable` with a non-empty
  reason. Expected metadata and local-Git uncertainty returns `undecidable`
  without an operational error, so downstream discovery leaves the Review
  Artifact live.
- Classification reads the highest-numbered `round-*` directory and its
  `round.md` frontmatter. A missing head, unreadable or malformed metadata,
  ambiguous Round numbering, an invalid head, or a Git answer that cannot be
  established remains undecidable.
- Local Git resolves the default through `origin/HEAD`, preferring the matching
  up-to-date local branch, then falls back to local `main` or `master`. It checks
  ancestry, refs containing the recorded head, and local or remote-tracking refs
  for the recorded PR Head Branch.
- An unmerged head stays live when any local ref reaches it. An unreachable head
  retires only when its recorded branch is also absent; a rewritten branch that
  still exists remains undecidable.
- The implementation invokes only read-only local `git` commands with optional
  locks and filesystem-monitor refresh disabled. It has no hosting-provider,
  credential, HTTP, or network dependency.

### Focused-check evidence

- The pre-change focused run,
  `GOCACHE=/tmp/roundfix-0094-task-03-gocache rtk go test ./internal/spec -run
  '^TestClassifyReviewLocalGit$'`, failed to compile because `ClassifyReview` and
  the three `ReviewLiveness` constants did not exist.
- The same focused command passed 11 tests after implementation. Named fixtures
  cover a head merged into `main`, a nonstandard default named by `origin/HEAD`,
  an unreachable head after branch deletion, a reachable live branch, a deleted
  branch whose head remains reachable through a tag, and a rewritten branch that
  still exists.
- The focused suite also passed the several-Round fixture, where Round 001's
  default-branch head is finished but Round 002's head is live, proving the
  newest Round decides.
- The no-head fixture asserts both `ReviewUndecidable` and explicitly refuses
  `ReviewFinished`. Malformed newest metadata and a repository with no locally
  identifiable default branch also return `ReviewUndecidable`.
- Every fixture routes through an assertion that rejects an empty reason; the
  no-head negative case performs the same check directly.
- `GOCACHE=/tmp/roundfix-0094-task-03-gocache rtk go test ./internal/spec`
  passed all 298 package tests.
- `GOCACHE=/tmp/roundfix-0094-task-03-gocache rtk go test -race
  ./internal/spec -run '^TestClassifyReviewLocalGit$'` passed all 11 focused
  tests with the race detector.
- Source inspection with `rtk rg -n
  'exec\.Command|runGH|gh api|net/http|https?://'` over the new source and test
  files found only `exec.Command("git", gitArgs...)`; no provider or transport
  call is present.
- `rtk git diff --check` exited 0, `rtk gofmt -d` reported no diff for the two
  new Go files, and `GOCACHE=/tmp/roundfix-0094-task-03-gocache rtk make
  verify-incremental` exited 0 after the implementation and Result edits,
  covering formatting, the Go suite, skill checks, and the build.

The Daemon-owned commands in `## Verification` were not run in this Agent turn.
