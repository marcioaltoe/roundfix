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

Implemented project-neutral disposable-repository journeys that compose
confirmed setup, the profile's declared formatter contract, the persisted
Verification, two clean audits, and byte-identical reapply for all three
maintained profiles. Added strict Standard TypeScript Monorepo Profile journeys
for REST and Post-only HTTP Contract Decisions with typed exceptions, selected
Inngest evidence, absent Docker evidence, and exact frontend/backend workspace
paths.

Added a complete Baseline Readoption journey whose inventory contains every
structural Source Baseline Entry kind and exercises all four dispositions:
current managed entry, typed repository document, Repository-Specific Normative
Rules, and individual rejection. The journey proves confirmed creation,
typed-document preservation, maintainer edits surviving reapply, stale-plan
rejection, required-capability no-write blocking, recommended-capability
warnings, and atomic rollback on write or postwrite failure.

Expanded the owned-version contract to a 66-field mutation matrix spanning Go,
Python, npm packages and dependency pins, both Roundfix-owned skill trees,
profile/setup/source schemas and snapshots, the Release Plan schema, and the
changelog. Every one-field mutation is rejected, while the existing protected
Run Database, external lock, protocol, and upstream-skill fixture remains
accepted unchanged.

Added a public Release Plan Command journey backed by a real disposable Git
repository and bare remote plus a fake two-page GitHub provider. It inventories
shared, local-only, and remote-only stable tags, exhausts all release pages,
returns approval-required exit `3`, preserves stdout/stderr discipline, leaves
the repository unchanged, and records zero Git or GitHub mutation calls.

Acceptance evidence:

- Every maintained profile: `ProfileMacroFlowTests.test_supported_profiles_apply_audit_clean_and_reapply_without_changes` executes formatter compatibility, the persisted Verification, two audits, and reapply with identical repository bytes.
- TypeScript contracts: `ReadoptionApplyTests.test_rest_and_post_only_contracts_persist_typed_exceptions_and_exact_workspace_evidence` preserves both modes, ordered typed exceptions, optional-module status, and exact workspace evidence in the Setup Snapshot and reapply.
- Readoption completeness: `ReadoptionApplyTests.test_every_structural_entry_kind_has_an_individual_destination_and_preserves_repository_bytes` covers `file`, `managed-block`, `manifest-record`, and `unmarked-span`; existing stale-confirmation and rollback cases remain green in the same suite.
- Capability safety: the Readoption required-capability case proves zero writes, the new recommended-capability case proves stable non-blocking warnings, and `RepositoryCapabilityTests.test_capability_evaluation_has_no_write_install_network_or_script_side_effects` rejects forbidden evaluation side effects.
- Version ownership: `VersionContractTests.test_authoritative_distribution_surfaces_report_0_0_1` validates and mutation-tests all 66 owned fields; `test_protected_versions_match_the_operational_and_upstream_fixture` proves protected versions remain unchanged.
- Reset planning: `TestReleasePlanResetInventoriesTemporaryGitRemoteAndPaginatedGitHubReadOnly` proves real Git inventory, fake GitHub pagination, exit `3`, repository preservation, and zero mutation calls.
- Distribution parity: `make skills-sync-check` and a recursive tree comparison reported no canonical/distributed difference.

Verification:

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests` — passed, 246 tests.
- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s skills/setup-context-driven/tests` — passed, the same 246 tests.
- `rtk go test ./internal/releaseplan ./internal/cli ./skills` — passed, 740 tests across 3 packages.
- `rtk make verify` — passed: both 246-test setup suites, 1,727 Go tests across 20 packages, asset loading, Repository Skill Set check, and build.

Follow-ups: none.
