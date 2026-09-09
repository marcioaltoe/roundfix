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
3. MUST emit `SC-TOOLING-UNAPPROVED` when a Spec claims authorization its cited
   record does not sustain, naming which field withheld the grant and pointing
   at both the citing row and the record. A row asserting express maintainer
   authorization over a record that is proposed, undated, withdrawn, or
   consumed by another Spec is such a claim.
4. MUST NOT refuse a row that honestly declares its mutations as proposed and
   claims no operative grant. A Spec in authoring that says its grant is
   proposed is accurate, not defective, and seven active Specs are in exactly
   that state; refusing them would fail `make verify-docs` for every Spec in the
   portfolio and block the pull request boundary on correct declarations.
5. MUST resolve the operative record when a row cites more than one, preferring
   the record the row names as its express authorization over any record it
   mentions as proposed, and MUST refuse rather than guess when the row's claim
   cannot be matched to exactly one record.
6. MUST keep `SC-TOOLING-UNAUTHORIZED`, `SC-TOOLING-UNTYPED` and
   `SC-TOOLING-UNBOUNDED` emitting for the conditions they already own, with
   their existing tokens unchanged.
7. MUST keep a record that legitimately passes today passing: a dated legacy
   record naming the citing Spec with bounded paths must not acquire a new
   refusal from this change.
8. MUST leave every active Spec in this repository passing `roundfix spec check`
   after the change, proven by running the checker over the whole Spec Root
   rather than over this Spec alone.
9. MUST update the recorded characterization rows this Task intentionally
   changes and leave every other recorded row untouched, so the diff states the
   behavior change instead of hiding it.

## Subtasks

- [ ] Resolve a Spec-relative citation to its repository-relative record path.
- [ ] Route resolved records through the typed reader by role, not by filename.
- [ ] Emit `SC-TOOLING-UNAPPROVED` for a claim the record does not sustain.
- [ ] Leave an honest proposal declaration passing, and resolve a row citing two records.
- [ ] Preserve the existing tooling refusal tokens and their conditions.
- [ ] Update only the characterization rows this change intends to move.

## Acceptance Criteria

- [ ] A Spec whose Tooling authority row cites its own `_authorization.md` by a
      Spec-relative link resolves that record, where the same input resolved to
      none before this Task.
- [ ] A Spec whose row asserts express maintainer authorization over a record
      with `status` proposed or a null grant date reports
      `SC-TOOLING-UNAPPROVED` naming the withholding field.
- [ ] A Spec whose row honestly declares its mutations proposed and claims no
      operative grant reports no tooling finding.
- [ ] A row citing both an approved narrow grant and a proposed record resolves
      to the approved one; a row whose claim matches no single record refuses
      instead of guessing.
- [ ] Running the checker over the whole Spec Root reports no error for any
      active Spec, so the pull request boundary stays passable.
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
- `grep -q 'SC-TOOLING-UNAPPROVED' internal/speccheck/constraints.go && grep -q 'func TestConstraintsRefuseNonOperativeGrant' internal/speccheck/constraints_characterization_test.go && go test -count=1 ./internal/speccheck -run '^TestConstraintsRefuseNonOperativeGrant$'` — a row claiming express authorization over a proposal or a null grant date refuses with the new code, and an approved record does not.
- `grep -q 'func TestConstraintsAcceptHonestProposalDeclaration' internal/speccheck/constraints_characterization_test.go || exit 1; go test -count=1 ./internal/speccheck -run '^TestConstraintsAcceptHonestProposalDeclaration$' || exit 1; go run -buildvcs=false ./cmd/roundfix spec check` — an honest proposal declaration passes, and every active Spec in the repository still checks clean, so the pull request boundary stays passable.
- `grep -q 'func TestConstraintsResolveSpecContainedRecord' internal/speccheck/constraints_characterization_test.go && go test -count=1 ./internal/speccheck` — the package suite runs with the new cases present, so the untouched characterization rows still pass.

## References

- `_prd.md` → User Stories 2; Core Features 4; Goals 2; Decisions: Declared intentional breaks 1.
- `_techspec.md` → Vocabulary Contract; Implementation Design: Audit and compatibility; Build Order 2.
- ADR-0096, ADR-0117.
