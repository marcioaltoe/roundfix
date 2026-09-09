---
task: task_03
spec: 0119-spec-contained-authorization
status: pending
type: backend
complexity: medium
---

# Task 03: Refuse a cited record that is not an operative grant

## Overview

Make authoring-time validation depend on the record's role instead of a date in
its filename, and make a record cited from inside its own Spec resolve. The
slice is verifiable on its own: a Spec citing a proposal starts failing
`roundfix spec check` where it passed before, and a Spec citing an approved
record keeps passing.

This is an authorized tooling Task. It may change only
`internal/speccheck/constraints.go`,
`internal/speccheck/constraints_characterization_test.go`, and this Task file.
Stop before any other mutation. The bounded set comes from the approved grant in
[_authorization.md](_authorization.md).

## Requirements

1. MUST resolve an authorization record cited from inside its own Spec,
   including a Spec-relative citation, to the record's repository-relative path.
2. MUST validate a resolved record through the typed reader regardless of
   whether its filename carries a date, and MUST NOT reintroduce the filename
   date as the trigger for validation.
3. MUST emit `SC-TOOLING-UNAPPROVED` when a cited record resolves but is not an
   operative grant for the citing Spec, naming which field withheld the grant
   and pointing at both the citing row and the record.
4. MUST keep `SC-TOOLING-UNAUTHORIZED`, `SC-TOOLING-UNTYPED` and
   `SC-TOOLING-UNBOUNDED` emitting for the conditions they already own, with
   their existing tokens unchanged.
5. MUST keep a record that legitimately passes today passing: a dated legacy
   record naming the citing Spec with bounded paths must not acquire a new
   refusal from this change.
6. MUST update the recorded characterization rows this Task intentionally
   changes and leave every other recorded row untouched, so the diff states the
   behavior change instead of hiding it.

## Subtasks

- [ ] Resolve a Spec-relative citation to its repository-relative record path.
- [ ] Route resolved records through the typed reader by role, not by filename.
- [ ] Emit `SC-TOOLING-UNAPPROVED` with the withholding field named.
- [ ] Preserve the existing tooling refusal tokens and their conditions.
- [ ] Update only the characterization rows this change intends to move.

## Acceptance Criteria

- [ ] A Spec whose Tooling authority row cites its own `_authorization.md` by a
      Spec-relative link resolves that record, where the same input resolved to
      none before this Task.
- [ ] A Spec citing a record with `status` proposed or a null grant date reports
      `SC-TOOLING-UNAPPROVED` naming the withholding field.
- [ ] A Spec citing an approved record for itself reports no tooling finding.
- [ ] A Spec citing a record that names a different consuming Spec still reports
      the unauthorized condition it reports today.
- [ ] The characterization rows changed by this Task are exactly the two the
      requirements name; running the suite with any other recorded row altered
      fails.

## Context

- interface: `internal/speccheck/citations.go`
- interface: `internal/speccheck/report.go`

## Verification

- `grep -q 'func TestConstraintsResolveSpecContainedRecord' internal/speccheck/constraints_characterization_test.go && go test -count=1 ./internal/speccheck -run '^TestConstraintsResolveSpecContainedRecord$'` — a Spec-relative citation resolves to its record path.
- `grep -q 'SC-TOOLING-UNAPPROVED' internal/speccheck/constraints.go && grep -q 'func TestConstraintsRefuseNonOperativeGrant' internal/speccheck/constraints_characterization_test.go && go test -count=1 ./internal/speccheck -run '^TestConstraintsRefuseNonOperativeGrant$'` — a proposal and a null grant date each refuse with the new code, and an approved record does not.
- `grep -q 'func TestConstraintsResolveSpecContainedRecord' internal/speccheck/constraints_characterization_test.go && go test -count=1 ./internal/speccheck` — the package suite runs with the new cases present, so the untouched characterization rows still pass.

## References

- `_prd.md` → User Stories 2; Core Features 4; Goals 2; Decisions: Declared intentional breaks 1.
- `_techspec.md` → Vocabulary Contract; Implementation Design: Audit and compatibility; Build Order 2.
- ADR-0096, ADR-0117.
