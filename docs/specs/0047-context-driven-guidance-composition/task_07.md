---
task: task_07
spec: 0047-context-driven-guidance-composition
status: pending
type: infra
complexity: high
---

# Task 07: Synchronize composed Baseline assets

## Overview

Carry the hierarchy, document contracts, semantic destinations, and Profile
adaptation behavior through every maintained embedded and canonical Baseline
asset. Formatter goldens and retention accounting become the distribution
contract for the composed result.

## Requirements

1. MUST update every affected Profile, module, decision effect, template,
   source corpus, coverage declaration, and retention transition coherently.
2. MUST regenerate canonical setup snapshots only from approved immutable
   sources through the Go-owned synchronization command.
3. MUST preserve all upstream-managed skill content and provenance.
4. MUST update formatter fixtures for every affected maintained Profile.
5. MUST make catalog, source accounting, formatter, audit, and empty reapply
   validation fail on any missing asset or clause.
6. MUST keep portable assets free of Fluxus and Oraculum names or policy.

## Subtasks

- [ ] Align embedded modules, templates, profiles, and coverage.
- [ ] Update source corpora and retention transitions.
- [ ] Refresh formatter-compatible golden fixtures.
- [ ] Synchronize canonical setup snapshots.
- [ ] Add asset completeness and branding guards.

## Acceptance Criteria

- [ ] The embedded catalog loads with one deterministic digest.
- [ ] Every required clause and semantic destination has complete source
  accounting.
- [ ] All maintained Profiles produce formatter-stable generated output.
- [ ] Canonical and distributed setup snapshots agree byte-for-byte.
- [ ] Upstream skill digests and project-agnostic branding guards pass.

## Context

- instruction: `docs/adr/0059-generated-output-is-formatter-stable-in-the-target-repository.md`
- instruction: `docs/adr/0060-source-baselines-are-exhaustive-and-project-agnostic.md`
- interface: `internal/baseline/assets`
- interface: `internal/baseline/assets_sync.go`
- interface: `internal/baseline/testdata/parity-corpus/v1`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestGuidanceCompositionAssets|TestFormatterComposition|TestSourceBaseline|TestBaselineAssetsSync'` — expected: catalog, accounting, formatter, provenance, sync, and branding cases pass.
- `rtk go run -buildvcs=false ./cmd/roundfix baseline assets sync --help` — expected: the Go-owned synchronization contract remains available.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goals 1–5; User Stories 1–6; Core Features 1–17; Success Metrics.
- `_techspec.md` → Integration Points; Testing Approach; Build Order 6.
- ADR-0059 → Formatter-Stable Output.
- ADR-0060 → exhaustive project-agnostic assets.
