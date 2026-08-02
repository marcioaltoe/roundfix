---
task: task_10
spec: 0057-baseline-capability-evidence-and-retention
status: pending
type: backend
complexity: high
---

# Task 10: Warn only about carriers nobody manages

## Overview

A clean adoption warns about the thirteen files it just wrote, so a real nested
carrier conflict is buried in noise about managed artifacts warning about
themselves. This Task classifies carriers and narrows warnings to the ones that
actually need a human — but only on positive evidence, because a warning that
disappears wrongly is worse than one that fires too often.

## Requirements

1. MUST classify each carrier as a current managed artifact, a stale managed
   artifact, a recognized repository extension, or an unmanaged nested carrier.
2. MUST warn only about unmanaged nested carriers.
3. MUST treat a stale managed artifact as upgrade input requiring clause-level
   retention rather than as a warning.
4. MUST name, in every warning it does emit, the boundary between managed and
   preserved bytes.
5. MUST narrow only on positive evidence: a carrier that cannot be classified
   keeps the warning it produces today.
6. MUST leave every non-carrier diagnostic unchanged.

## Subtasks

- [ ] Classify carriers into the four kinds.
- [ ] Suppress warnings for current managed artifacts and recognized
      extensions.
- [ ] Route stale managed artifacts to retention rather than to a warning.
- [ ] Name the managed/preserved boundary in remaining warnings.
- [ ] Keep the warning for anything unclassifiable.

## Acceptance Criteria

- [ ] An idempotent re-plan after a verified apply reports zero warnings about
      the managed guides the apply wrote.
- [ ] An unmanaged nested carrier still warns.
- [ ] A stale managed artifact appears as retention input, not as a warning.
- [ ] A recognized repository extension does not warn.
- [ ] A carrier that cannot be classified still warns, proven by an
      unclassifiable fixture.
- [ ] Every emitted warning names the managed versus preserved byte boundary.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/` and
      this task file.

## Context

- interface: `internal/baseline/classification.go`
- interface: `internal/baseline/preservation.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline -run TestCarrierClassification -count=1` —
  expected: exit 0; all four kinds classify and only unmanaged nested carriers
  warn.
- `go test ./internal/baseline -run TestUnclassifiableCarrierStillWarns -count=1`
  — expected: exit 0; narrowing requires positive evidence.
- `go test ./internal/baseline -run TestBaselinePlanCharacterization -count=1` —
  expected: exit 0 with goldens re-recorded only for the intended warning
  suppression; no warning disappears without a positive classification.
- `go test ./internal/baseline -count=1` — expected: exit 0.

## References

- `_prd.md` → User Story 6; Core Features 2 and 12; Success Metrics
  (idempotent re-plan reports zero warnings).
- `_techspec.md` → Risks (carrier classification can under-warn);
  Build Order 8.
