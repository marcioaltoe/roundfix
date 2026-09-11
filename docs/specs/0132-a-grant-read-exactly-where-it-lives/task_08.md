---
task: task_08
spec: 0132-a-grant-read-exactly-where-it-lives
status: pending
type: backend
complexity: medium
---

# Task 08: Let the mechanical audit consume the resolved reference

## Overview

Task 04 resolves a Spec-relative citation against the artifact that carries it.
The changed-path audit is the consumer that must read that resolved reference
and judge bounded paths against the project repository, and it lives in a file
Task 04's bounded scope forbids. This slice is that consumer. It is verifiable
alone: an external Spec Root produces a real audit where it produces a skip
today.

`internal/speccheck/mechanical.go` is not a governed path, so this Task needs no
tooling grant for it. Its regression lives in
`internal/speccheck/constraints_characterization_test.go`, which the approved
grant covers, rather than in `internal/speccheck/mechanical_test.go`, which
belongs to Spec 0119's grant and not to this one.

## Requirements

1. MUST carry the resolved reference from the constraint reader into the
   mechanical authorization request, so the audit never re-derives it with a
   narrower rule.
2. MUST judge bounded paths and ancestor authority against the project
   repository root and delivery target, never against the Spec repository, per
   the two-root boundary the TechSpec names.
3. MUST produce a real changed-path audit for a Spec whose record lives in an
   external Spec Root, where it records a skip today.
4. MUST keep every audit outcome that works today: the out-of-grant refusal, the
   grant-edited-in-the-consuming-commit refusal, and the presence-aware skip for
   a Spec that genuinely declares no authorization.
5. MUST NOT let an unresolved reference read as a skip, because a skip says
   there was nothing to audit and that is the failure this pairing removes.

## Subtasks

- [ ] Carry the resolved reference into the mechanical request.
- [ ] Judge paths against the project root while reading the record from the Spec root.
- [ ] Prove the external case produces a real audit.
- [ ] Prove the preserved outcomes and the unresolved-is-not-a-skip rule.

## Acceptance Criteria

- [ ] A Spec whose record lives in an external Spec Root produces a changed-path
      audit with Task commits; it records a skip today.
- [ ] The audit judges bounded paths against the project repository, proven by a
      case where the two roots hold different trees and judging the Spec root
      would reach the wrong answer.
- [ ] A governed change outside the grant still refuses, and a Spec declaring no
      authorization still records the presence-aware skip.
- [ ] An unresolved reference refuses rather than skipping.

## Context

- interface: `internal/speccheck/mechanical.go`

## Verification

- `grep -q 'func TestMechanicalAuditConsumesResolvedReference' internal/speccheck/constraints_characterization_test.go && go test -count=1 ./internal/speccheck -run '^TestMechanicalAuditConsumesResolvedReference$'` — an external Spec Root produces a real audit; this fails today.
- `grep -q 'func TestMechanicalAuditJudgesTheProjectRoot' internal/speccheck/constraints_characterization_test.go && go test -count=1 ./internal/speccheck -run '^TestMechanicalAuditJudgesTheProjectRoot$'` — bounded paths are judged against the project repository, proven where judging the Spec root gives the wrong answer.
- `grep -q 'func TestMechanicalAuditConsumesResolvedReference' internal/speccheck/constraints_characterization_test.go || exit 1; go test -count=1 ./internal/speccheck` — the package suite passes with the preserved outcomes unmoved.

## References

- `_prd.md` → Goals 2; Core Features 3; User Stories 2.
- `_techspec.md` → Implementation Design: Two roots, named separately; The boundary both roots travel in; Testing Approach observation 2.
