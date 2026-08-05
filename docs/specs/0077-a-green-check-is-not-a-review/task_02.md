---
task: task_02
spec: 0077-a-green-check-is-not-a-review
status: pending
type: backend
complexity: medium
---

# Task 02: Name the refusal so the stall is legible

## Overview

With the default inverted, a refused head already stalls — but it stalls as
`pending`, which reads as "no signal yet" rather than "the reviewer declined".
An operator cannot tell a refusal from an absence.

This slice recognises the documented refusal shapes and resolves them `skipped`
with the reason preserved. It explains the stall; it does not create it.

## Requirements

1. MUST recognise the documented shapes in which the Review Source declines to
   review the expected head, including the rate-limit refusal and the
   path-filter skip, and resolve them as `skipped`.
2. MUST recognise by class rather than by one literal: title casing and the
   documented variants all resolve the same way.
3. MUST read the authoritative signal the vendor documents — the comment as well
   as the check — because the check conclusion is green by design.
4. MUST preserve the refusal reason verbatim in the evidence detail, bounded by
   the existing helper.
5. MUST NOT let recognition widen what reaches `verified`. If refusal
   recognition and the closed default ever disagree, the default wins.
6. MUST NOT retry, re-request, or wait for capacity. That is deferred to the
   follow-on Spec and is out of scope here.

## Subtasks

- [ ] Recognise the documented refusal shapes and their reasons.
- [ ] Resolve them `skipped` ahead of every other classification.
- [ ] Add the class table and the #107 replay.
- [ ] Assert stale-head isolation for refusals.

## Acceptance Criteria

- [ ] The Pull Request #107 payload resolves `skipped` with a reason naming the
      rate limit.
- [ ] A path-filter skip resolves `skipped`, unchanged from today.
- [ ] Title-case variants of each documented refusal resolve identically,
      asserted by a table rather than one case.
- [ ] The refusal reason appears in the evidence detail.
- [ ] A refusal recorded against an earlier commit does not settle the current
      head.
- [ ] No payload reaches `verified` that did not reach it after task_01.
- [ ] No retry, re-request, or capacity wait exists in this Task's changes.

## Context

- interface: `internal/reviewsource/coderabbit/coderabbit.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/reviewsource/... -count=1 -run 'Refusal|Skip|RateLimit|Class' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the refusal tests ran and passed.
- `go test ./internal/reviewsource/... ./internal/watch -count=1` — expected:
  exit 0.
- `go test -parallel 16 ./... 2>&1 | grep -q "^FAIL" && exit 1 || exit 0`
  — expected: exit 0.

## References

- `_prd.md` → Core Features 1, 3 and 5; Success Metrics 1, 2 and 4.
- `_techspec.md` → Interfaces; Build Order 2.
- ADR-0054.
