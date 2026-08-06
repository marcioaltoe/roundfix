---
task: task_07
spec: 0073-skill-versions-decoupled-from-the-binary
status: pending
type: backend
complexity: medium
---

# Task 07: Assert the property directly, and record the external provenance

## Overview

Corrective Task from the QA gate's F-001 (`Blocks-Completion`) and F-002.

**F-001 is a harness defect, not a product defect.**
`TestOwnedSkillEditLeavesMakeVerifyGreen` proves the Spec's headline property
by running the complete `make verify` **nested inside a test**. That makes the
Spec's own acceptance hostage to every other test in the repository: the gate
observed it fail on
`TestTaskCycleVerificationCapacityCancellationWhileQueuedStartsNoCommandOrSettlement`,
which passes in isolation and failed under concurrent load. The property being
proven is narrow — an owned-skill edit needs no regeneration step — and it
should be asserted narrowly.

**F-002 is a Success Metric with no Task.** The PRD requires that
`golang-dependency-management`, `golang-safety`, and `golang-structs-interfaces`
carry recorded provenance without becoming setup-snapshot entries. No Task
requirement ever asked for it, so nothing delivered it — the same shape as Spec
0059's F-001 in the same session.

## Requirements

1. MUST replace the nested-`make verify` proof with a direct assertion of the
   property: after editing an owned skill's content, the derived Baseline
   artifacts and both characterization corpora are unchanged, so no
   regeneration step is required.
2. MUST assert that property from the artifacts themselves — compare their
   bytes before and after the edit — rather than by invoking a repository-wide
   gate from inside a test.
3. MUST NOT weaken what is proven. The claim is unchanged; only the instrument
   changes, from "the whole suite passes" to "these artifacts did not move".
4. MUST record provenance for the three external Go Skills without creating a
   setup-snapshot entry for any of them, since the PRD keeps third-party skills
   outside the version contract.
5. MUST make the recorded provenance visible where an operator looks for it,
   and MUST NOT hold those skills to an owned minimum.
6. MUST leave every owned-skill readiness behaviour from task_02 and task_04
   unchanged, asserted over the existing tests.

## Subtasks

- [ ] Replace the nested gate with a direct artifact-stability assertion.
- [ ] Record provenance for the three external Go Skills.

## Acceptance Criteria

- [ ] Editing an owned skill leaves every derived Baseline artifact and both
      characterization corpora byte-identical, asserted directly.
- [ ] No test in this Spec invokes `make verify` from inside a test.
- [ ] The three external Go Skills carry recorded provenance.
- [ ] None of the three appears in any setup snapshot.
- [ ] None of the three is compared against an owned minimum.
- [ ] Owned-skill readiness behaviour is unchanged, asserted over the existing
      tests.

## Context

- interface: `skills/baseline_skill_contract_integration_test.go`
- interface: `skills/repository.go`
- interface: `skills-lock.json`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `grep -rn "make verify" skills/*_test.go | grep -q . && exit 1 || exit 0`
  — expected: exit 0; no test in this package invokes the repository gate.
- `output="$(go test ./skills -count=1 -run 'OwnedSkillEdit|Provenance' -v 2>&1)"; status=$?; printf '%s\n' "$output"; [ "$status" -eq 0 ] || exit "$status"; printf '%s\n' "$output" | grep -q -- "--- PASS"`
  — expected: exit 0; the property and provenance tests ran and passed.
- `go run -buildvcs=false ./cmd/roundfix skills list | grep -q "golang-safety"`
  — expected: exit 0; the recorded provenance is visible to an operator.
- `go test ./skills ./internal/baseline -count=1` — expected: exit 0.
- `make verify` — expected: exit 0.

## References

- `qa/qa-report-2026-08-06.md` → F-001, F-002.
- `_prd.md` → Success Metrics 1 and 7; Which skills the contract covers.
- `_techspec.md` → Testing Approach.
