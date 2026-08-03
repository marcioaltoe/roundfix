---
task: task_07
spec: 0064-spec-artifact-consistency-gate
status: pending
type: docs
complexity: medium
---

# Task 07: Bring this repository's own Specs to a clean report

## Overview

Run the check across this repository's active Specs, correct every reported
`error`, and add the test that holds them at zero from here on. This is the
step that makes a fail-closed gate possible: wiring the gate before the Specs
are clean would turn `make verify` red for every contributor.

The output is a declared break list, not a discovered one. Every correction is
a real artifact contradiction the check located, with both sides named — never
a rewrite for style.

## Requirements

1. MUST run the check across every active Spec and correct each reported
   `error` in the Spec artifacts, so the active corpus reports zero errors.
2. MUST leave archived Specs byte-identical. An archived Spec's findings are
   recorded by the corpus golden, never fixed.
3. MUST add a test that runs the check across every active Spec and fails,
   naming each Spec and code, when any `error` remains. The test lands with the
   corrections it holds in place.
4. MUST record each correction in the Task's Result as a declared break: the
   Spec, the code, and what changed on each side.
5. MUST NOT suppress, downgrade, or narrow a detector to make a Spec pass. A
   detector that over-reaches is a defect in the detector, and correcting it
   there is in scope; silencing it is not.
6. MUST leave every reported `gap` visible, dismissing one only by writing its
   reason into the Spec that carries it.
7. SHOULD keep each correction minimal — the smallest edit that removes the
   contradiction the check located.

## Subtasks

- [ ] Run the check across the active corpus and capture the report.
- [ ] Correct each reported error in the owning Spec artifact.
- [ ] Resolve or reason-dismiss each reported gap.
- [ ] Add the active-corpus zero-error test.
- [ ] Update the corpus golden's active-Spec counts to zero.
- [ ] Record every correction as a declared break in the Result.

## Acceptance Criteria

- [ ] The check reports zero `error` findings across every active Spec.
- [ ] A test asserts that and fails naming the Spec and code when it regresses.
- [ ] Every archived Spec file is byte-identical to its pre-Task content.
- [ ] No detector was disabled, narrowed, or made advisory to reach zero,
      asserted by the unchanged detector tests from tasks 01 through 03.
- [ ] Each reported gap is either resolved or carries a written reason in its
      Spec.
- [ ] The Result lists every correction with its Spec, code, and both sides.

## Context

- instruction: `docs/agents/spec-routing.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go run -buildvcs=false ./cmd/roundfix spec check` — expected: exit 0; the
  active corpus reports no error.
- `go test ./internal/speccheck -count=1 -run 'ActiveCorpus' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the active-corpus test ran and passed.
- `go test ./internal/speccheck -count=1` — expected: exit 0; the detector
  tests from tasks 01 through 03 still pass unchanged.
- `git diff --name-only HEAD -- docs/specs/_archived | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no archived Spec file changed.

## References

- `_prd.md` → Success Metrics; Decisions (non-regression).
- `_techspec.md` → Build Order 7; Risks & Considerations.
- ADR-0094.
