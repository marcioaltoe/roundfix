---
task: task_12
spec: 0045-context-driven-baseline-0-0-1-reset
status: pending
type: test
complexity: high
---

# Task 12: Prove complete profile journeys

## Overview

Exercise the assembled 0.0.1 behavior through real disposable repositories,
not isolated snapshots alone. Each maintained profile must survive setup,
formatter execution, persisted Verification, audit, and reapply while the
high-risk blocking and preservation paths remain observable.

## Requirements

1. MUST add project-neutral macro fixtures for every maintained profile and
   run both canonical and distributed setup skill suites in place.
2. MUST exercise apply, configured formatter, persisted Verification, audit,
   second audit after formatting, and idempotent reapply for each profile.
3. MUST cover the full Standard TypeScript Monorepo Profile with REST and
   Post-only decisions, typed exceptions, optional-module selection, and exact
   workspace evidence.
4. MUST cover incompatible Baseline Readoption with every structural entry
   kind, individual dispositions, repository-rules creation, typed-document
   moves, preservation, stale confirmation, and rollback.
5. MUST prove missing required capabilities block, recommended capabilities
   warn without blocking, and evaluation performs no forbidden side effects.
6. MUST verify the complete owned/protected version matrix across Go, Python,
   npm, skills, schemas, snapshots, changelog, and operational exclusions.
7. MUST exercise reset Release Plan behavior through temporary Git remotes and
   a fake paginated GitHub provider without mutating real release state.
8. MUST keep tests hermetic, deterministic, network-free, and based on
   observable behavior.

## Subtasks

- [ ] Build reusable disposable repositories for every maintained profile.
- [ ] Add the full apply/format/verify/audit/reapply journey.
- [ ] Add TypeScript HTTP, optional-module, workspace, and capability variants.
- [ ] Add complete Readoption, preservation, stale-plan, and rollback journeys.
- [ ] Add owned/protected version-matrix coverage.
- [ ] Add temporary Git and paginated GitHub reset-plan journeys.
- [ ] Run the same embedded suite from both skill-tree locations.

## Acceptance Criteria

- [ ] Every maintained profile completes the macro journey and finishes with a
      clean audit and byte-identical reapply.
- [ ] REST and Post-only TypeScript fixtures preserve their typed contracts and
      explicit exceptions across snapshot and reapply.
- [ ] Every Source Baseline Entry kind is dispositioned explicitly, and
      repository-owned files remain unchanged on later setup runs.
- [ ] Required-capability absence blocks with no writes; recommended absence
      emits a stable warning and allows the journey to continue.
- [ ] One-field mutations across every owned version surface fail, while
      protected operational and upstream versions remain accepted unchanged.
- [ ] Reset-plan fixtures exhaust pagination, return approval-required exit 3,
      and record zero mutation calls.
- [ ] Canonical and distributed suites are independently runnable and pass the
      same cases.

## Context

- instruction: `docs/adr/0059-setup-output-is-stable-under-declared-formatters.md`
- interface: `.agents/skills/setup-context-driven/tests/test_workflow.py`
- interface: `.agents/skills/setup-context-driven/tests/test_formatter_compatibility.py`
- interface: `.agents/skills/setup-context-driven/tests/test_macro_profiles.py`
- interface: `.agents/skills/setup-context-driven/tests/test_upgrade_retention.py`
- interface: `internal/cli/releaseplan_command_test.go`
- interface: `internal/releaseplan/build_test.go`

## Verification

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests` — expected: every canonical setup unit and macro journey passes without network access.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s skills/setup-context-driven/tests` — expected: the distributed copy runs in place and passes the same suite.
- `rtk go test ./internal/releaseplan ./internal/cli ./skills` — expected: reset planning, distribution contracts, and Repository Skill Set behavior pass.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → User Stories 1–6; Success Metrics; Non-Goals / Out of Scope.
- `_techspec.md` → Testing Approach; Risks & Considerations; Build Order 8.
- ADR-0059 → declared formatter stability.
- ADR-0060 through ADR-0065 → complete 0.0.1 behavior and safety boundaries.
