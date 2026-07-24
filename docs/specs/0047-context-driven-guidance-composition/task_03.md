---
task: task_03
spec: 0047-context-driven-guidance-composition
status: pending
type: backend
complexity: high
---

# Task 03: Segment preserved rules without byte loss

## Overview

Add the sealed segmentation stage that turns one structural Source Baseline
Entry into byte-exact clause entries before classification. Deterministic
validation prevents gaps, overlap, reordering, or rewritten source content.

## Requirements

1. MUST bind each segmentation proposal to one immutable snapshot and current
   Source Baseline identity.
2. MUST accept only ordered, non-empty ranges whose union covers every source
   byte exactly once.
3. MUST derive segment bytes and digests locally instead of trusting ACP output
   content.
4. MUST materialize admitted ranges as normal Source Baseline Entries for the
   existing classification contract.
5. MUST keep the original entry as one lossless segment when semantic
   segmentation is unavailable.
6. MUST remain bounded, checkout-free, proposal-only, and compatible with the
   preferred and fallback ACP sequence.

## Subtasks

- [ ] Define strict segmentation snapshot and proposal schemas.
- [ ] Validate range coverage and source identities.
- [ ] Materialize canonical segment entries.
- [ ] Connect preferred, fallback, and lossless manual paths.
- [ ] Add adversarial range and determinism tests.

## Acceptance Criteria

- [ ] Gap, overlap, duplicate, reordered, empty, stale, and out-of-bounds
  proposals fail closed.
- [ ] Concatenating admitted segment bytes reproduces the exact source bytes.
- [ ] Equivalent proposals produce byte-identical Source Baseline identities.
- [ ] ACP failure retains the complete original entry without mutation.
- [ ] No segmentation payload exposes a checkout path or write capability.

## Context

- instruction: `docs/adr/0069-baseline-semantic-analysis-is-read-only-and-supervised.md`
- instruction: `docs/adr/0074-repository-rules-use-hybrid-semantic-ownership.md`
- interface: `internal/baseline/classification.go`
- interface: `internal/baseline/preservation.go`
- interface: `internal/baselineacp`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/baselineacp -run 'TestRuleSegmentation|TestSegmentationSnapshot|TestSegmentationProposal'` — expected: byte coverage, bounds, sealed analysis, fallback, and determinism cases pass.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goals 2–3; User Story 2; Core Features 6–7.
- `_techspec.md` → Implementation Design: Interfaces and Data Models; Build Order 2.
- ADR-0069 → supervised proposal boundary.
- ADR-0074 → byte-exact clause segmentation.
