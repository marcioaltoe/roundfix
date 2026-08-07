---
task: task_02
spec: 0082-the-manifest-already-answered-that
status: pending
type: backend
complexity: high
---

# Task 02: Refresh managed regions without asking who owns the prose

## Overview

Adds the third instruction-preservation mode: a managed refresh that regenerates
only Baseline-owned regions and leaves every other byte identical. This is the
slice that removes the supervised analyzer from the update path — in the mode,
planning produces no source baseline and no decision skeleton, which are already
the two conditions every caller uses to skip classification. It is verifiable on
its own through planning alone, before any new command exists.

## Requirements

1. MUST add a managed-refresh preservation mode accepted by root-preservation
   planning alongside the existing greenfield and preservation modes.
2. MUST produce, in that mode, a ready preservation plan carrying no source
   baseline entries and no decision skeleton, so no classification input is ever
   required and no semantic analyzer is invoked.
3. MUST still detect and report blocking carriers; an unsafe root carrier blocks
   a managed refresh exactly as it blocks adoption.
4. MUST treat a hand-edited managed marker as blocking rather than as a warning,
   because the mode's guarantee depends on the markers being trustworthy.
5. MUST plan no root instruction backup in this mode, per ADR-0100, and instead
   let the plan's preimages carry the preservation proof.
6. MUST leave the greenfield and preservation modes behaviorally unchanged.
7. MUST prove, by test, that every byte outside a managed marker is identical
   before and after a managed-refresh plan is applied — including authored prose
   and repository-rule blocks — asserting on region digests rather than a golden
   file.
8. MUST leave retention accounting on this path untouched, so an unaccounted
   managed clause still blocks planning.

## Subtasks

- [ ] Add the mode and accept it in root-preservation planning.
- [ ] Return a ready plan with no source baseline and no decision skeleton.
- [ ] Make hand-edited managed markers blocking in this mode.
- [ ] Suppress root backup planning in this mode.
- [ ] Build the mixed fixture: managed markers, authored prose, repository rules.
- [ ] Assert non-managed region digests survive a plan and apply unchanged.
- [ ] Assert an unaccounted managed clause still blocks in this mode.

## Acceptance Criteria

- [ ] A managed-refresh plan on an adopted repository reports a ready state with
      zero source baseline entries and no decision skeleton.
- [ ] Applying a managed-refresh plan against a stale catalog leaves every
      non-managed region byte-identical, proven by digest comparison.
- [ ] A repository with a hand-edited managed marker blocks in managed-refresh
      mode and names the offending path.
- [ ] A managed-refresh plan emits no root backup postimage.
- [ ] A managed-refresh plan whose transition leaves a managed clause unaccounted
      does not become applicable and names the unaccounted clause.
- [ ] The task_01 corpus still passes, proving adoption is unchanged.

## Context

- interface: `internal/baseline/preservation.go`
- interface: `internal/baseline/plan.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go test ./internal/baseline/ -run 'ManagedRefresh' -v 2>&1 | grep -q '^--- PASS: .*ManagedRefresh'` — expected: exits 0, proving the new mode's cases exist and pass rather than being selected out.
- `go test ./internal/baseline/ -run 'ManagedRefresh.*Preserv|Preserv.*ManagedRefresh' -v 2>&1 | grep -q '^--- PASS: '` — expected: exits 0, proving the byte-identity assertion ran.
- `go test ./internal/baseline/ ./internal/cli/ -count=1` — expected: exits 0, with the task_01 corpus passing unchanged.

## References

- `_techspec.md` → Build Order 2 and 3; Interfaces: `PreservationModeManagedRefresh`; Testing Approach.
- `_prd.md` → Core Features 2 and 3; User Stories 3 and 6; Goals 2 and 5.
- ADR-0058, ADR-0069, ADR-0070, ADR-0099, ADR-0100.
