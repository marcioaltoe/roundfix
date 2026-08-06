---
task: task_02
spec: 0069-review-run-targets-its-pull-request
status: pending
type: backend
complexity: high
---

# Task 02: Keep the target while the Run writes

## Overview

task_01 stops a Run that starts on the wrong branch. This slice stops one that
*becomes* wrong: the checkout moves while the Run is Active, and every write
after that lands somewhere the Pull Request never named.

The PRD assumed a late check already existed to be moved. Authoring disproved
it — `checkout branch mismatch` is not in the tree, and the nearest live code
is Round artifact reuse in `rounds.go`. This slice builds the guard.

The distinction it must preserve is the expensive one: a Review Issue that
failed on its merits and a Run stopped by its environment cost different
things. One deserves attention, the other deserves a re-run unchanged. The
session this Spec came from failed two legitimate security findings for an
environmental reason, and they had to be redone from scratch.

## Requirements

1. MUST record the Run's target branch and its revision at Preflight, so the
   mid-Run comparison has an anchor that does not depend on re-querying the
   forge.
2. MUST re-read the checkout and compare it against that anchor before every
   write boundary: each Batch commit, the review artifact commit, and Final
   Push.
3. MUST reach a terminal outcome distinct from `Failed` when the checkout no
   longer matches, so an interruption never reads as a Review Issue defect.
4. MUST leave the affected Review Issues unsettled rather than failed, so a
   re-run starts from their real state.
5. MUST make every review artifact commit and push target the Pull Request's
   head branch, asserted from Git rather than from the log.
6. MUST settle through the normal outcome path under ADR-0052; an interruption
   that cannot reach a terminal state is worse than the defect it replaces.
7. MUST NOT change what a Run does when the checkout never moves, asserted over
   the existing review tests unchanged.
8. MUST NOT move, check out, or restore the working tree.

## Subtasks

- [ ] Record the target branch and revision on the Run.
- [ ] Add the write-boundary re-read and the comparison.
- [ ] Add the terminal interruption outcome and its report line.
- [ ] Assert issues stay unsettled and artifacts land on the head branch.

## Acceptance Criteria

- [ ] A checkout moved before a Batch commit reaches the interruption outcome
      and commits nothing.
- [ ] A checkout moved before the review artifact commit reaches it too.
- [ ] A checkout moved before Final Push reaches it and pushes nothing.
- [ ] The interruption is a distinct terminal outcome, not `Failed`.
- [ ] Review Issues affected by an interruption are left unsettled, not failed.
- [ ] Every review artifact commit lands on the Pull Request's head branch,
      asserted from Git.
- [ ] A Run whose checkout never moves behaves exactly as it does today.

## Context

- interface: `internal/watch/watch.go`
- interface: `internal/store`
- instruction: `docs/adr/0036-review-artifacts-are-committed-in-a-separate-docs-commit.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `output="$(go test ./internal/watch -count=1 -run 'Interrupt|Moved|Target|Boundary' -v 2>&1)"; status=$?; printf '%s\n' "$output"; [ "$status" -eq 0 ] || exit "$status"; printf '%s\n' "$output" | grep -q -- "--- PASS"`
  — expected: exit 0; the interruption tests ran and passed.
- `go test ./internal/watch ./internal/cli ./internal/store -count=1`
  — expected: exit 0.
- `go test -parallel 16 ./...` — expected: exit 0.
- `git diff --name-only HEAD | grep -E "^(\.agents|skills)/" | grep -q . && exit 1 || exit 0`
  — expected: exit 0; the Skill is task_03's bounded scope.

## References

- `_prd.md` → Core Features 3 and 4; Success Metrics 2, 3 and 4.
- `_techspec.md` → Disproven premise; Build Order 2.
- ADR-0052, ADR-0036.
