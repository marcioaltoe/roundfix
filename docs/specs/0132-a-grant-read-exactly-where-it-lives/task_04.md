---
task: task_04
spec: 0132-a-grant-read-exactly-where-it-lives
status: pending
type: backend
complexity: high
---

# Task 04: Resolve a citation against the artifact that carries it

## Overview

For an external Spec Root the constraint reader produces a display path that
reaches outside the code repository, so joining a Spec-relative link to it
yields a path the containment check rejects and a valid record is refused. This
slice resolves the reference against the artifact's own location and anchors
containment to the Spec Root. It is verifiable alone: the external citation
resolves and the default one is unchanged.

This is an authorized tooling Task. It may change only
`internal/speccheck/constraints.go`,
`internal/speccheck/constraints_characterization_test.go`, and this Task file.
Stop before any other mutation. The bounded set comes from
[_authorization.md](_authorization.md).

## Requirements

1. MUST resolve a Spec-relative citation against the artifact that carries it,
   not against the code repository root.
2. MUST anchor the containment check to the resolved Spec Root, so a record
   beside its PRD in an external root is inside the boundary rather than
   outside it.
3. MUST keep rejecting a reference that escapes the Spec Root, including one
   that traverses upward past it or resolves through a symlink out of it.
4. MUST keep the default root's resolution identical, including which record a
   row citing several references resolves to.
5. MUST update only the characterization rows this Task intentionally moves.

## Subtasks

- [ ] Resolve the reference relative to the carrying artifact.
- [ ] Anchor containment to the Spec Root.
- [ ] Prove escape is still rejected.
- [ ] Prove the default root's answers are unchanged.

## Acceptance Criteria

- [ ] A record cited as `[_authorization.md](_authorization.md)` beside a PRD in
      an external Spec Root resolves; it was refused before this Task.
- [ ] A reference traversing upward past the Spec Root, and one resolving
      through a symlink out of it, are still rejected.
- [ ] The default root resolves the same record it resolves today, including for
      a row citing an approved grant beside a proposed record.
- [ ] The mechanical audit receives the resolved reference for the external case
      and produces a real audit rather than a skip.

## Context

- interface: `internal/speccheck/constraints.go`
- interface: `internal/speccheck/mechanical.go`

## Verification

- `grep -q 'func TestCitationResolvesInExternalSpecRoot' internal/speccheck/constraints_characterization_test.go && go test -count=1 ./internal/speccheck -run '^TestCitationResolvesInExternalSpecRoot$'` — an external Spec-relative citation resolves; this fails today.
- `grep -q 'func TestCitationRejectsEscapeFromSpecRoot' internal/speccheck/constraints_characterization_test.go && go test -count=1 ./internal/speccheck -run '^TestCitationRejectsEscapeFromSpecRoot$'` — upward traversal and symlink escape are still rejected.
- `grep -q 'func TestCitationResolvesInExternalSpecRoot' internal/speccheck/constraints_characterization_test.go || exit 1; go test -count=1 ./internal/speccheck` — the package suite passes with the default-root answers unmoved.

## References

- `_prd.md` → Goals 2; User Stories 2; Core Features 3.
- `_techspec.md` → Implementation Design: Paths derived from the resolved root; Build Order 4.
