---
task: task_09
spec: 0047-context-driven-guidance-composition
status: completed
type: test
complexity: high
---

# Task 09: Prove composed Baseline journeys

## Overview

Create hermetic macro journeys that prove the composed guidance contract across
greenfield, update, Readoption, Profile adaptation, formatting, apply, audit,
and empty reapply. Live Fluxus and Oraculum checks remain final QA evidence,
not Task prerequisites.

## Requirements

1. MUST cover every affected maintained Profile through real CLI planning and
   apply boundaries.
2. MUST cover legacy generic carriers, semantic redistribution, residual
   retention, and zero-residual cleanup.
3. MUST reproduce the backend-only TypeScript divergence and reviewed Profile
   adaptation without project-specific fixture branding.
4. MUST prove formatter, repository Verification recommendation, audit, and
   empty reapply composition.
5. MUST assert complete managed-entry and Upgrade Retention ledgers plus exact
   rollback after injected failure.
6. MUST leave Fluxus and Oraculum journeys to separately authorized
   `qa-gate` execution and name their required evidence.

## Subtasks

- [x] Add all-profile greenfield and update macro fixtures.
- [x] Add semantic redistribution and residual fixtures.
- [x] Add Profile adaptation and universal-remediation fixtures.
- [x] Add formatter, audit, reapply, and rollback assertions.
- [x] Define the final live QA evidence matrix.

## Acceptance Criteria

- [x] Every maintained Profile completes its hermetic journey with zero
  second-pass delta.
- [x] Zero residual rules produce no generic or repository-specific carrier.
- [x] The backend-only fixture reaches a valid repository-owned Profile without
  waivers.
- [x] Injected failure restores all affected files, including a planned Profile
  file.
- [x] The Spec-local QA plan names fresh Fluxus greenfield/update and Oraculum
  divergence journeys.

## Context

- instruction: `docs/adr/0073-baseline-apply-uses-a-recoverable-multi-file-transaction.md`
- interface: `internal/baseline/compatibility_corpus_test.go`
- interface: `internal/baseline/release_gate_test.go`
- interface: `internal/cli/baseline_release_gate_test.go`
- interface: `internal/baseline/testdata/parity-corpus/v1`

## Verification

- `rtk go test -count=100 ./internal/baseline -run 'TestAssetsSyncProvenanceAndPreMutationRefusals/dirty_or_untracked_checkout'` — expected: temporary Git fixture cleanup passes deterministically without process-level Git configuration overrides.
- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestManagedRootFreshPlan|TestGuidanceCompositionJourney|TestProfileAdaptationJourney'` — expected: a verified greenfield, update, and repository-owned Profile apply each produce a valid zero-change fresh Plan while roots that still contain user-owned bytes continue to require an immutable backup.
- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestGuidanceCompositionJourney|TestSemanticRedistributionJourney|TestProfileAdaptationJourney|TestBaselineReleaseGate'` — expected: greenfield, update, adaptation, formatting, rollback, audit, and reapply journeys pass.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goals 1–6; User Stories 1–7; Core Features 1–20; Success Metrics.
- `_techspec.md` → Testing Approach; Build Order 8; Risks & Considerations.
- ADR-0073 → recoverable multi-file apply.

## QA repair cycle — 2026-07-24

The first full QA gate found three Blocks-Completion defects and two
documentation frictions. Repair the product behavior at its owning boundary,
extend the hermetic regression coverage, and rerun the complete live matrix.

- [x] F-01: an update MUST NOT propose setup-managed root bytes as
  repository-owned rules and then reject its own accepted proposal as stale,
  empty, or managed.
- [x] F-02: a reviewed, catalog-valid repository-owned Profile adaptation MUST
  reach the final portable Change Plan even when the existing Setup Manifest's
  Source Baseline has no unique maintained transition.
- [x] F-03: the first fresh Plan after a confirmed greenfield apply MUST have
  zero file changes and MUST NOT create a second digest-addressed root backup.
- [x] F-04: the documented complete greenfield Decision Document MUST include
  every required runtime decision.
- [x] F-05: formatter discovery and the reported post-apply recommendation MUST
  use the repository's real supported command.
- [x] The temporary Git fixture used by `TestFormatterComposition` MUST be
  isolated from host-level `commit.gpgsign=true` without requiring a process
  environment override.
- [x] Fresh regression tests MUST fail before each behavior repair and pass
  afterward.
- [ ] The complete Fluxus greenfield/update and Oraculum adaptation live matrix
  MUST pass from one fresh Roundfix build.

Evidence and exact reproductions:
`qa/qa-report-2026-07-24.md` and
`qa/evidence/2026-07-24-guidance-composition/`.

## QA repair cycle — 2026-07-25

The second full QA gate found a reproducible Git-fixture cleanup race before
the live matrix could start.

- [x] F-06: temporary Git fixtures MUST own all spawned Git work and prevent
  automatic background maintenance from outliving the command and racing
  `t.TempDir` cleanup.
- [x] The repair MUST live at the fixture or Git-command boundary; it MUST NOT
  suppress cleanup errors or depend on a process-level environment override.
- [x] The exact failing subtest MUST pass 100 consecutive repetitions without
  any override.
- [x] Unmodified `rtk make verify` MUST pass.
- [ ] The complete Fluxus greenfield/update and Oraculum adaptation live matrix
  MUST pass from one fresh Roundfix build.

Evidence and exact reproduction:
`qa/qa-report-2026-07-25.md` and
`qa/evidence/2026-07-25-guidance-composition/command-evidence.md`.

## QA persistence repair cycle — 2026-07-25

The third full QA gate passed the static gate and all initial public
Plan/apply/Verification paths, but every first fresh Plan failed at the same
root-backup validation boundary.

- [x] F-07: a root carrier composed entirely of verified setup-managed bytes
  MUST NOT require an immutable source backup on the first fresh Plan.
- [x] A root carrier that still contains any user-owned source bytes MUST
  continue to require and validate its immutable digest-addressed backup.
- [x] Greenfield, update, and repository-owned Profile adaptation MUST each
  complete Plan → apply → formatter → repository Verification → zero-change
  fresh Plan through public boundaries.
- [x] The repair MUST reconcile Plan assembly and Plan/apply validation at the
  same ownership boundary; it MUST NOT suppress the backup invariant globally.
- [x] Unmodified `rtk make verify` MUST pass.
- [ ] The complete Fluxus greenfield/update and Oraculum adaptation live matrix
  MUST pass from one fresh Roundfix build.

Evidence and exact reproductions:
`qa/qa-report-2026-07-25.md`,
`qa/evidence/2026-07-25-guidance-composition/command-evidence.md`, and
`qa/evidence/2026-07-25-guidance-composition/rerun/`.

## Result

Repaired the five QA findings at their owning Baseline boundaries and extended
the hermetic release-gate journeys.

- F-01: setup-managed markers now accept the shipped semantic version, so
  managed root bytes are excluded from repository-owned Readoption proposals.
- F-02: reviewed repository-owned Profile adaptation accepts the validated
  Source Profile Baseline as retention-compatible without inventing a
  maintained transition.
- F-03: an unchanged root containing only setup-managed guidance produces no
  redundant digest-addressed backup; Plan assembly and apply-contract
  validation enforce the same rule. The first fresh Plan after every
  maintained Profile's greenfield apply has zero file changes.
- F-04: the published complete greenfield Decision Document includes both
  required runtime decisions and is resolved against the embedded Profile
  catalog in its documentation contract test.
- F-05: formatter projection discovers the repository-declared `fmt` alias,
  persists `bun run fmt` into the Setup Manifest, and reports only executable
  repository recommendations after Plan and Apply.
- `TestFormatterComposition` sets repository-local `commit.gpgsign=false`, so
  the hermetic Git fixture does not inherit host signing policy or require a
  process environment override.
- The existing semantic redistribution, residual retention, managed-entry,
  Upgrade Retention, audit, empty-reapply, and exact rollback assertions remain
  green. Planned-Profile rollback still compares the complete visible
  repository tree with its preimage.
- `qa/live-journeys.md` remains the separately authorized evidence matrix for
  fresh Fluxus greenfield/update and Oraculum divergence/adaptation journeys.
  Task 09 did not read or mutate either live repository.

Verification:

- The focused red regressions reproduced managed-marker redistribution,
  incompatible source retention, redundant backup, incomplete documentation,
  formatter mismatch, and host signing failures before their repairs; the same
  focused suite passed afterward.
- `rtk env GOCACHE=/private/tmp/roundfix-0047-task09-go-cache go test -count=1
  ./internal/baseline ./internal/cli -run
  'TestGuidanceCompositionJourney|TestSemanticRedistributionJourney|TestProfileAdaptationJourney|TestBaselineReleaseGate'`
  passed across both packages.
- `rtk env GOCACHE=/private/tmp/roundfix-0047-task09-go-cache make verify` passed:
  2,201 repository tests across 22 packages, 4 skill contract tests,
  `roundfix skills check`, and the Roundfix build.
- `rtk git diff --check` passed.

Live QA remains pending by design and requires a separately authorized
`qa-gate` run using the Spec-local evidence matrix.

### QA repair cycle — 2026-07-25

Repaired F-06 at the asset-sync Git-command boundary. Every Git invocation
owned by the temporary fixture now disables automatic maintenance alongside
filesystem monitoring, preventing detached maintenance work from recreating
`.git/objects/info/packs` after the foreground command exits. The repair uses
no sleep, retry, cleanup suppression, or process-level Git environment
override.

Fresh acceptance evidence:

- Every maintained Profile completed the real CLI Plan/apply, formatter,
  repository Verification, audit, empty-reapply, and zero-second-pass-delta
  journey in `TestGuidanceCompositionJourney`; the focused journey command
  passed 34 tests across `internal/baseline` and `internal/cli`.
- `TestSemanticRedistributionJourney` proved exact semantic redistribution,
  canonical residual retention, and removal of generic and repository-specific
  carriers plus root pointers when the residual count is zero.
- `TestProfileAdaptationJourney` proved the unbranded backend-only TypeScript
  fixture reaches a catalog-valid repository-owned Profile, keeps universal
  capabilities required, applies the Profile postimage, and finishes with an
  empty reapply.
- `TestBaselineReleaseGate` passed the complete managed-entry and Upgrade
  Retention ledger checks plus transaction, semantic-guide, and planned-Profile
  rollback journeys. The planned-Profile rollback restores the complete
  visible repository preimage.
- `qa/live-journeys.md` names the separately authorized fresh Fluxus
  greenfield, Fluxus update, and Oraculum divergence/adaptation evidence
  requirements. This Task did not read or mutate those live repositories.

Verification:

- Before repair, `rtk go test -count=100 ./internal/baseline -run
  'TestAssetsSyncProvenanceAndPreMutationRefusals/dirty_or_untracked_checkout'`
  reproduced the cleanup race: 198 passing parent/subtest results and the
  failing subtest plus parent, with `TempDir RemoveAll cleanup` reporting
  `.git/objects/info` as non-empty.
- After repair, the same exact command passed 200 parent/subtest results across
  100 consecutive repetitions without a Git environment override.
- `rtk go test -count=1 ./internal/baseline ./internal/cli -run
  'TestGuidanceCompositionJourney|TestSemanticRedistributionJourney|TestProfileAdaptationJourney|TestBaselineReleaseGate'`
  passed 34 tests across both packages.
- Unmodified `rtk make verify` passed: 2,201 repository tests across 22
  packages, 4 skill contract tests, `roundfix skills check`, and the Roundfix
  build.
- `rtk git diff --check` passed.

The complete live Fluxus and Oraculum matrix remains pending by design for the
separately authorized `qa-gate` run.

### QA persistence repair cycle — 2026-07-25

Repaired F-07 at the shared root-source ownership boundary. A canonical
`AGENTS.md` root and a `CLAUDE.md` alias can select the same source bytes.
Plan assembly already omitted a redundant backup when that source's rendered
postimage was unchanged and entirely setup-managed, but Plan/apply validation
exempted only the direct carrier and then demanded a `CLAUDE.<digest>.md`
backup for the same source. Both boundaries now use the source identity and
the same unchanged setup-managed predicate. The existing source-byte equality
guard remains in Plan assembly.

Fresh acceptance evidence:

- `TestManagedRootFreshPlan` reproduced the live `AGENTS.md` plus `CLAUDE.md`
  alias failure before the repair. The setup-managed subtest failed with
  `root carrier source "AGENTS.md" has no immutable backup`, while its
  user-owned negative companion passed.
- After the repair, the setup-managed journey produced no backup, zero file
  changes, and a valid empty apply. The mixed user-owned journey still
  produced `backup:AGENTS.md` and passed strict Plan/apply validation.
- `TestGuidanceCompositionJourney` passed every maintained Profile through
  public Plan/apply, formatter, repository Verification, audit, and empty
  reapply boundaries with zero second-pass delta.
- `TestSemanticRedistributionJourney` passed exact redistribution, residual
  retention, and zero-residual carrier cleanup.
- `TestProfileAdaptationJourney` passed the unbranded backend-only TypeScript
  divergence and repository-owned Profile adaptation journey without a
  waiver.
- `TestBaselineReleaseGate` passed the complete managed-entry and Upgrade
  Retention ledgers plus exact transaction and planned-Profile rollback.
- `qa/live-journeys.md` still names the separately authorized fresh Fluxus
  greenfield/update and Oraculum divergence/adaptation evidence. This Task did
  not read or mutate those live repositories.

Verification:

- `rtk env GOCACHE=/private/tmp/roundfix-0047-task09-f07-go-cache go test
  -count=100 ./internal/baseline -run
  'TestAssetsSyncProvenanceAndPreMutationRefusals/dirty_or_untracked_checkout'`
  passed 100 consecutive repetitions.
- `rtk env GOCACHE=/private/tmp/roundfix-0047-task09-f07-go-cache go test
  -count=1 ./internal/baseline ./internal/cli -run
  'TestManagedRootFreshPlan|TestGuidanceCompositionJourney|TestSemanticRedistributionJourney|TestProfileAdaptationJourney|TestBaselineReleaseGate'`
  passed across both packages.
- The first full-gate attempt encountered one unrelated transient
  `.git/objects/maintenance.lock` cleanup failure in
  `TestBaselineMacroJourneysPublicCLI/consolidated_preservation_review`. The
  exact subtest passed immediately, and the final
  `rtk env GOCACHE=/private/tmp/roundfix-0047-task09-f07-go-cache make verify`
  rerun passed 2,204 repository tests across 22 packages, 4 skill contract
  tests, `roundfix skills check`, and the Roundfix build.
- `rtk git diff --check` passed.

The complete live Fluxus and Oraculum matrix remains pending by design for the
separately authorized `qa-gate` run.

### Verification Feedback — attempt 1

The Daemon's full gate exposed a temporary-repository lifecycle race in
`TestRepositoryInspectionNoMutation`: Git-internal state changed between the
test's before and after snapshots under full-suite load. The shared inspection
fixture now disables filesystem monitoring and automatic maintenance inline
for every Git setup, clone, checkout, add, and commit command. This keeps all
spawned Git work owned by the foreground fixture command without sleeps,
retries, cleanup suppression, or process-level Git configuration overrides.

Focused verification after the repair:

- `rtk env GOCACHE=/private/tmp/roundfix-0047-task09-f07-go-cache go test
  -count=100 ./internal/baseline -run
  '^TestRepositoryInspectionNoMutation$'` passed 100 consecutive repetitions.
- `rtk env GOCACHE=/private/tmp/roundfix-0047-task09-f07-go-cache go test
  -count=1 ./internal/baseline ./internal/cli -run
  'TestManagedRootFreshPlan|TestGuidanceCompositionJourney|TestProfileAdaptationJourney'`
  passed across both packages.
- `rtk git diff --check` passed.

The Daemon owns the single configured Verification rerun after this feedback
turn.
