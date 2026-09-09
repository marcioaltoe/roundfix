---
task: task_01
spec: 0119-spec-contained-authorization
status: pending
type: test
complexity: high
---

# Task 01: Characterize today's grant resolution and governed set

## Overview

Record what the constraint reader and the governed-path predicate do today,
before anything changes them, so the regression gate is the measured present
behavior rather than a test written afterwards to agree with new code. The
slice is verifiable on its own: the characterization runs against the unchanged
readers and states their current answers, including the answers this Spec
intends to change.

This is an authorized tooling Task. It may change only
`internal/speccheck/constraints_characterization_test.go`,
`internal/speccheck/governed_repocontract_test.go`, and this Task file. Stop
before any other mutation. The bounded set comes from the approved grant in
[_authorization.md](_authorization.md).

## Requirements

1. MUST record how the constraint reader resolves a Tooling authority row's
   authorization record today: which citation forms produce a resolved record
   path and which produce none, and that a record whose filename carries no
   date skips typed-grant validation entirely.
2. MUST record that a record cited from inside its own Spec by a
   Spec-relative link resolves to no record path today, so the row passes with
   its grant unread.
3. MUST record today's membership answers of the governed-path predicate for
   the owned shipped authoring templates, showing which of them the current
   predicate does not report as governed.
4. MUST assert the recorded answers as the present contract, so a later change
   to either reader fails this characterization rather than passing silently.
5. MUST NOT change the behavior of any reader, and MUST NOT weaken or delete an
   existing assertion in either file.
6. SHOULD keep the recorded expectations table-driven so a later Task edits one
   row per intended change and leaves the rest untouched.

## Subtasks

- [ ] Enumerate the citation forms the constraint reader accepts today and the
      answer each one yields.
- [ ] Record the date-keyed skip of typed-grant validation as present behavior.
- [ ] Record the Spec-relative citation resolving to no record path.
- [ ] Record governed-path membership for the owned shipped authoring templates.
- [ ] Assert every recorded answer so a behavior change fails here first.

## Acceptance Criteria

- [ ] A characterization case names each accepted citation form and its current
      resolved record path, and a case names a Spec-relative citation resolving
      to none.
- [ ] A characterization case shows a record whose filename carries no date
      skipping typed-grant validation while a dated record does not.
- [ ] A characterization case lists the owned shipped authoring templates with
      the predicate's current governed answer for each.
- [ ] Reverting any one recorded expectation to a different answer fails the
      suite, so the file cannot pass while disagreeing with the readers.
- [ ] Every assertion that existed in both files before this Task still runs
      and still passes.

## Context

- interface: `internal/speccheck/constraints.go`
- interface: `internal/speccheck/governed.go`

## Verification

- `grep -q 'func TestConstraintReaderCharacterizesGrantCitation' internal/speccheck/constraints_characterization_test.go && go test -count=1 ./internal/speccheck -run '^TestConstraintReaderCharacterizesGrantCitation$'` — the citation-form and date-keyed-skip characterization executes; its absence cannot pass.
- `grep -q 'func TestGovernedSetCharacterizesOwnedShippedTemplates' internal/speccheck/governed_repocontract_test.go && go test -count=1 -tags repocontract ./internal/speccheck -run '^TestGovernedSetCharacterizesOwnedShippedTemplates$'` — today's governed answer for each owned shipped template is recorded and asserted.
- `grep -q 'func TestConstraintReaderCharacterizesGrantCitation' internal/speccheck/constraints_characterization_test.go && go test -count=1 ./internal/speccheck` — the package suite runs with the new characterization present, so every assertion that existed before still passes beside it.

## References

- `_prd.md` → Decisions: Regression locks; Success Metrics.
- `_techspec.md` → Testing Approach observation 5; Build Order 1.
- `_authorization.md` → approved bounded paths.
- ADR-0130.
