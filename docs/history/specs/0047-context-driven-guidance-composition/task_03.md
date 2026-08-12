---
task: task_03
spec: 0047-context-driven-guidance-composition
status: completed
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

- [x] Define strict segmentation snapshot and proposal schemas.
- [x] Validate range coverage and source identities.
- [x] Materialize canonical segment entries.
- [x] Connect preferred, fallback, and lossless manual paths.
- [x] Add adversarial range and determinism tests.

## Acceptance Criteria

- [x] Gap, overlap, duplicate, reordered, empty, stale, and out-of-bounds
  proposals fail closed.
- [x] Concatenating admitted segment bytes reproduces the exact source bytes.
- [x] Equivalent proposals produce byte-identical Source Baseline identities.
- [x] ACP failure retains the complete original entry without mutation.
- [x] No segmentation payload exposes a checkout path or write capability.

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

## Result

Implemented a sealed, digest-bound Segmentation Snapshot and strict Proposal
contract. The deterministic validator admits only ordered, non-empty,
byte-exhaustive ranges, derives bytes, digests, offsets, and Source Baseline
Entry identities locally, and keeps empty structural evidence unchanged.
Materialized entries pass the existing classification snapshot contract.

Connected segmentation to the existing exact preferred and fallback Agent
Selections. Both attempts reuse the same canonical snapshot in fresh private
directories; invalid, unavailable, tool-using, or incomplete analysis falls
back without checkout access. When neither attempt is admitted, the manual
proposal preserves every non-empty original entry as one lossless range.

Acceptance evidence:

- Gap, overlap, duplicate, reordered, empty, stale, unknown, rewritten-digest,
  and out-of-bounds proposals fail closed in
  `TestSegmentationProposalRejectsInvalidRangesAndStaleIdentity`; strict JSON
  and output bounds are covered by
  `TestSegmentationProposalParserRejectsUntrustedJSON`.
- Exact byte reconstruction and classification compatibility pass in
  `TestRuleSegmentationMaterializesExactBytesAndStableIdentities`.
- Equivalent parsed and direct proposals produce deeply identical Source
  Baselines and Source Baseline Entry identities in
  `TestRuleSegmentationMaterializesExactBytesAndStableIdentities`.
- ACP unavailability preserves the complete original entry without mutation
  in `TestRuleSegmentationACPFailureRetainsCompleteOriginalEntry` and
  `TestRuleSegmentationManualFallbackRetainsOriginalEntry`.
- Checkout-free canonical payloads and fresh private, empty ACP directories are
  verified by `TestSegmentationSnapshotIsCanonicalAndCheckoutFree` and
  `TestRuleSegmentationUsesPrivateCheckoutFreeDirectories`.

Verification:

- `rtk go test -count=1 ./internal/baseline ./internal/baselineacp -run 'TestRuleSegmentation|TestSegmentationSnapshot|TestSegmentationProposal'`
  — passed, 22 tests in 2 packages.
- `rtk make verify` — passed.

Follow-ups: none.
