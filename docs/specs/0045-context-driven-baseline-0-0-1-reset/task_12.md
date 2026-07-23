---
task: task_12
spec: 0045-context-driven-baseline-0-0-1-reset
status: completed
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

- [x] Build reusable disposable repositories for every maintained profile.
- [x] Add the full apply/format/verify/audit/reapply journey.
- [x] Add TypeScript HTTP, optional-module, workspace, and capability variants.
- [x] Add complete Readoption, preservation, stale-plan, and rollback journeys.
- [x] Add owned/protected version-matrix coverage.
- [x] Add temporary Git and paginated GitHub reset-plan journeys.
- [x] Run the same embedded suite from both skill-tree locations.

## Acceptance Criteria

- [x] Every maintained profile completes the macro journey and finishes with a
      clean audit and byte-identical reapply.
- [x] REST and Post-only TypeScript fixtures preserve their typed contracts and
      explicit exceptions across snapshot and reapply.
- [x] Every Source Baseline Entry kind is dispositioned explicitly, and
      repository-owned files remain unchanged on later setup runs.
- [x] Required-capability absence blocks with no writes; recommended absence
      emits a stable warning and allows the journey to continue.
- [x] One-field mutations across every owned version surface fail, while
      protected operational and upstream versions remain accepted unchanged.
- [x] Reset-plan fixtures exhaust pagination, return approval-required exit 3,
      and record zero mutation calls.
- [x] Canonical and distributed suites are independently runnable and pass the
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

## QA Reopen

Final QA reproduced a clean `standard-typescript-monorepo` apply that accepted
no package, workspace, PostgreSQL, or LogTape evidence and wrote legacy
portable manifest and guide contracts. See
[F-01 in the 2026-07-23 QA report](qa/qa-report-2026-07-23.md#f-01--clean-standard-typescript-monorepo-adoption-bypasses-the-001-contract).

Add a public disposable-repository regression that starts without required
capability evidence and proves no-write blocking. After supplying the exact
requirements, the same journey MUST prove the frontend/backend architecture
carriers, strict manifest and snapshot identity, clean audit, and byte-empty
reapply through the public setup entry point. Run the case from both canonical
and distributed skill trees.

## Result

The QA-reopened clean-adoption journey is covered by the public disposable
repository regression
`ReadoptionApplyTests.test_clean_profile_adoption_uses_the_strict_change_plan`.
It starts with a complete `0.0.1` Decision File but no package, workspace,
PostgreSQL, or LogTape evidence; the public Setup Command exits `1` with
`capability.required.missing` and leaves every repository byte unchanged. The
same repository then receives the exact required evidence and completes
preview, digest-confirmed apply, clean audit, and byte-empty reapply.

The regression also proves the strict output contract: profile snapshot schema
`setup-context-driven/profile-snapshot/0.0.1`, manifest schema
`setup-context-driven/manifest/0.0.1`, version and generator version `0.0.1`,
baseline `baseline.standard-typescript-monorepo-0.0.1`, exact embedded Setup
Snapshot identity, and the frontend/backend architecture carriers required by
the Standard TypeScript Monorepo Profile. It already lived at the public script
boundary after Task 08 repaired the shared strict planning path, so Task 12 did
not add a duplicate macro assertion.

Acceptance evidence:

- Every maintained profile: the canonical and distributed suites passed
  `ProfileMacroFlowTests.test_supported_profiles_apply_audit_clean_and_reapply_without_changes`,
  covering apply, declared formatter, persisted Verification, two clean audits,
  and byte-identical reapply.
- TypeScript decisions and exact workspaces: both suites passed
  `ReadoptionApplyTests.test_rest_and_post_only_contracts_persist_typed_exceptions_and_exact_workspace_evidence`.
- Source Baseline Entry dispositions and repository-byte preservation: both
  suites passed the structural inventory, stale confirmation, preservation,
  and rollback cases in `ReadoptionApplyTests`.
- Capability safety: the QA-reopened public regression passed independently in
  each skill tree and again in both full suites; the required-capability,
  recommended-capability, and side-effect rejection cases also passed.
- Owned/protected versions: both suites passed the complete mutation matrix and
  protected operational/upstream version fixture.
- Reset planning: `rtk go test ./internal/releaseplan ./internal/cli ./skills`
  passed 740 tests, including temporary Git remote, paginated GitHub inventory,
  approval-required exit `3`, and zero mutation calls.
- Distribution parity: `rtk make skills-sync-check` passed, and the canonical
  and distributed suites independently passed the same 252 tests.

Verification:

- Focused QA-reopen regression in `.agents/skills/setup-context-driven/tests` —
  passed.
- Focused QA-reopen regression in `skills/setup-context-driven/tests` — passed.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests`
  — passed, 252 tests.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s skills/setup-context-driven/tests`
  — passed, 252 tests.
- `rtk go test ./internal/releaseplan ./internal/cli ./skills` — passed, 740
  tests across 3 packages.
- `rtk make skills-sync-check` — passed.
- `rtk make verify` — passed: both 252-test setup suites, 1,727 Go tests across
  20 packages, asset loading, Repository Skill Set check, and build.

Follow-ups: none.
