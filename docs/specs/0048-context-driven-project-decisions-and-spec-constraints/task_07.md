---
task: task_07
spec: 0048-context-driven-project-decisions-and-spec-constraints
status: pending
type: infra
complexity: high
---

# Task 07: Synchronize project-decision assets

## Overview

Propagate project decisions, rendered clauses, tooling authority, and Spec
constraints through every affected Baseline Profile, canonical snapshot,
formatter fixture, public guide, and thin setup skill. Distribution guards
make the compiled CLI and restored repository workflow agree.

## Requirements

1. MUST update affected Profiles, modules, templates, source accounting,
   retention transitions, formatter fixtures, and canonical setup snapshots.
2. MUST update Roundfix user documentation with exact human and automation
   decision examples and refusal behavior.
3. MUST update the thin setup skill without adding an independent decision or
   rendering engine.
4. MUST synchronize repo-owned authorial skill snapshots from immutable
   provenance.
5. MUST preserve every upstream-managed skill byte.
6. MUST make catalog, documentation, formatter, skill-sync, and project-
   agnostic branding checks fail on drift.

## Subtasks

- [ ] Align embedded and canonical Baseline assets.
- [ ] Refresh formatter and source-accounting fixtures.
- [ ] Synchronize repo-owned authorial skill snapshots.
- [ ] Update public documentation and the thin setup skill.
- [ ] Add distribution, provenance, and branding guards.

## Acceptance Criteria

- [ ] Every affected Profile selects and renders the correct project decisions.
- [ ] Canonical and distributed setup assets agree byte-for-byte.
- [ ] Documentation examples parse through the public CLI contracts.
- [ ] The setup skill contains no independent execution behavior.
- [ ] Upstream skill digests and project-agnostic asset checks pass.

## Context

- instruction: `docs/agents/skill-governance.md`
- interface: `internal/baseline/assets`
- interface: `internal/baseline/assets_sync.go`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `.agents/skills/setup-context-driven/SKILL.md`
- interface: `skills`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/cli ./skills -run 'TestProjectDecisionAssets|TestProjectConstraintDocumentation|TestThinSetupSkill|TestAuthorialSkillSync'` — expected: asset, formatter, docs, skill ownership, provenance, and branding cases pass.
- `rtk go run -buildvcs=false ./cmd/roundfix baseline assets sync --help` — expected: canonical synchronization remains Go-owned.
- `rtk go run -buildvcs=false ./cmd/roundfix skills check` — expected: distributed repo-owned skills match their contracts.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goals 1–5; User Stories 1–6; Core Features 1–17; Success Metrics.
- `_techspec.md` → Integration Points; Build Order 6.
- ADR-0066 → CLI authority and thin skill.
- ADR-0076 and ADR-0077 → distributed decision and Spec contracts.
