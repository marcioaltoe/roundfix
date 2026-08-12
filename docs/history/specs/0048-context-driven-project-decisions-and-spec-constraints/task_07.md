---
task: task_07
spec: 0048-context-driven-project-decisions-and-spec-constraints
status: completed
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

- [x] Align embedded and canonical Baseline assets.
- [x] Refresh formatter and source-accounting fixtures.
- [x] Synchronize repo-owned authorial skill snapshots.
- [x] Update public documentation and the thin setup skill.
- [x] Add distribution, provenance, and branding guards.

## Acceptance Criteria

- [x] Every affected Profile selects and renders the correct project decisions.
- [x] Canonical and distributed setup assets agree byte-for-byte.
- [x] Documentation examples parse through the public CLI contracts.
- [x] The setup skill contains no independent execution behavior.
- [x] Upstream skill digests and project-agnostic asset checks pass.

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

## Result

- Added catalog validation and drift guards that require the maintained
  TypeScript Profile to select `identifier.strategy` and every applicable
  capability-gated project decision. Golden clauses, Source Baseline
  accounting, retention transitions, formatter fixtures, and project-agnostic
  branding are covered by `TestProjectDecisionAssets`.
- Refreshed normalized catalog and parity-corpus fixtures, then synchronized
  the three canonical setup snapshots from the immutable skills source at
  commit `236847f6956134bf468abb641bac0493a899bca5`. Repo-owned authorial skill
  digests now match canonical and distributed trees byte-for-byte; the pinned
  upstream aggregate digest remains unchanged.
- Expanded the public guide and the thin setup skill with exact human and
  automation decision examples, conditional Better Auth behavior, complete
  refusal semantics, and bounded Project Constraint examples. Documentation
  contract tests parse the published Decision Document and plan commands
  through public CLI parsers.
- Kept execution authority in the compiled CLI. `TestThinSetupSkill` and
  `TestProjectConstraintDocumentation` assert that the setup skill does not
  collect, derive, validate, or render decisions independently.
- Verification evidence:
  - `rtk go test -count=1 ./internal/baseline ./internal/cli ./skills -run
    'TestProjectDecisionAssets|TestProjectConstraintDocumentation|TestThinSetupSkill|TestAuthorialSkillSync'`
    passed.
  - `rtk go run -buildvcs=false ./cmd/roundfix baseline assets sync --help`
    passed and documents Go-owned synchronization.
  - `rtk go run -buildvcs=false ./cmd/roundfix baseline assets sync
    --source-dir /Users/marcio/dev/skills/setups --check --format json`
    passed with no drift findings.
  - `rtk go run -buildvcs=false ./cmd/roundfix skills check` passed for all
    14 repo-owned skills.
  - `rtk make verify` passed: 2,293 Go tests, four focused skill contract
    tests, the complete skill check, and the CLI build.
