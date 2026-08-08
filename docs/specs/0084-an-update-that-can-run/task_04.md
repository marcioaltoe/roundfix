---
task: task_04
spec: 0084-an-update-that-can-run
status: pending
type: test
complexity: medium
---

# Task 04: Prove the update converges on a second run

## Overview

Adds the property that tells a finished sweep from an unfinished one: after an
applied refresh, running the update again against an unchanged catalog reports
the repository current and proposes nothing. The property has no test today, and
it is the only check that fails if Setup Manifest republication ever regresses —
which is precisely the regression that made the cold-start block permanent in the
first place.

## Requirements

1. MUST build an adopted repository copy on the filesystem, age its Setup Manifest
   so the recorded digests no longer describe its managed regions, and leave the
   managed regions themselves untouched.
2. MUST apply a managed refresh to that copy and assert the apply verifies.
3. MUST assert the second plan over the same copy and the same catalog reports the
   repository current and proposes zero file changes.
4. MUST assert the republished Setup Manifest's recorded digests describe the
   bytes on disk after the apply, so convergence is proven by the record and not
   only by the reported state.
5. MUST assert that every byte outside a managed marker is identical before and
   after the apply, so convergence is not bought by rewriting authored prose.
6. MUST fail if manifest republication is removed: the test must distinguish a
   converged repository from one whose second run proposes the same work again.

## Subtasks

- [ ] Build the aged-manifest fixture from an adopted copy.
- [ ] Apply the refresh and assert verification.
- [ ] Assert the second run reports current with zero file changes.
- [ ] Assert the republished manifest describes the on-disk bytes.
- [ ] Assert non-managed regions are byte-identical across the apply.
- [ ] Demonstrate the negative: a copy whose manifest is not republished does not
      report current.

## Rehearsal Cases

- Case: an adopted copy whose Setup Manifest records digests that no longer
  describe its untouched managed regions; Observation: the first run reaches a
  ready plan rather than an action-required state.
- Case: the same copy after an approved apply; Observation: the second run
  reports the repository current and proposes zero file changes.
- Case: the same copy after an approved apply, reading the Setup Manifest
  directly; Observation: every recorded managed-artifact digest equals the digest
  of the corresponding on-disk region.
- Case: a copy whose Setup Manifest is deliberately left un-republished after the
  refresh; Observation: the second run proposes the same changes again, proving
  the assertion can fail.
- Case: the same copy, comparing every region outside a managed marker before and
  after the apply; Observation: every region digest is unchanged.

## Acceptance Criteria

- [ ] The aged-manifest copy reaches a ready plan on the first run.
- [ ] Applying that plan verifies.
- [ ] The second run over the applied copy reports the repository current with
      zero proposed file changes.
- [ ] Every managed-artifact digest recorded in the applied copy's Setup Manifest
      equals the digest of the corresponding on-disk region.
- [ ] Every region outside a managed marker has an identical digest before and
      after the apply.
- [ ] The negative case, with republication suppressed, does not report current.

## Context

- interface: `internal/baseline/preservation.go`
- interface: `internal/baseline/plan.go`
- interface: `internal/baseline/apply.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go test ./internal/baseline/ -run 'Converge' -v > /tmp/0084-task-04-a.log 2>&1 && grep -q '^--- PASS: .*Converge' /tmp/0084-task-04-a.log` — expected: exits 0, proving the convergence cases exist and pass rather than being selected out.
- `go test ./internal/baseline/ -run 'Converge' -count=2 > /tmp/0084-task-04-b.log 2>&1` — expected: exits 0, proving the fixture does not depend on state left by a previous run.
- `go test ./internal/baseline/ ./internal/cli/ -count=1` — expected: exits 0.

## References

- `_techspec.md` → Build Order 4; Testing Approach.
- `_prd.md` → Core Feature 4; User Story 4; Goal 2; Success Metrics.
- `references/2026-08-08-the-update-refuses-six-of-the-eight-copies-it-exists-to-update.md`
  → the measured second-run behavior this task turns into a test.
- ADR-0103, ADR-0100.
