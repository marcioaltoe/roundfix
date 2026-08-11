---
task: task_07
spec: 0041-agent-selection-runtime-readiness
status: completed
type: backend
complexity: medium
---

# Task 07: Prove profile configuration before persistence

## Overview

Make `profiles configure` prove candidate Agent Selection Profiles before
confirmation or persistence. Structural fallback failures must remain
pre-proof, and every failure or declined confirmation must leave the selected
configuration scope byte-identical.

## Requirements

1. MUST validate the complete candidate profile and distinct Fallback Chain
   before creating a disposable Session.
2. MUST prove every distinct candidate Preferred Selection and fallback through
   the shared readiness coordinator before confirmation or write.
3. MUST run exact proof for interactive, `--file`, `--dry-run`, `--json`, and
   `--yes` inputs without allowing any mode to bypass readiness.
4. MUST report that one additional distinct authorized and proven selection is
   required when fallback configuration is missing or duplicates Preferred.
5. MUST leave User Config or Project Config byte-identical on validation or
   proof failure, cleanup failure, output failure, or declined confirmation.
6. MUST keep recommendations advisory and never insert a ranked model or
   reasoning value into the candidate.
7. MUST report `changed: false` in JSON for every non-persisted outcome.

## Subtasks

- [x] Separate candidate construction from persistence.
- [x] Run structural fallback validation before Agent Session work.
- [x] Prove candidate tuples through the shared readiness coordinator.
- [x] Align interactive, file, dry-run, JSON, and yes modes.
- [x] Add distinct-fallback remediation without automatic substitution.
- [x] Snapshot config bytes across every non-persisted outcome.

## Acceptance Criteria

- [x] A complete proven candidate reaches confirmation and writes only after
      authorization.
- [x] Missing or duplicate fallback fails before any disposable Session and
      names the required distinct proof or authorization.
- [x] Unsupported adapter, model, or reasoning capability fails before
      confirmation and reports the shared classification.
- [x] `--dry-run` performs exact proof but writes nothing.
- [x] `--yes` bypasses confirmation only; it cannot bypass structural or
      runtime proof.
- [x] Decline, proof failure, cleanup failure, and JSON output failure leave
      the target bytes unchanged and report `changed: false` when JSON applies.
- [x] No recommendation is inserted automatically into Preferred Selection or
      the Fallback Chain.

## Context

- instruction: `docs/agents/autonomous-work.md`
- interface: `internal/cli/profiles_configure.go`
- interface: `internal/cli/profiles_validate.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/config/profile_config.go`
- interface: `internal/config/profiles.go`

## Verification

- `rtk go test ./internal/cli -run 'TestProfilesConfigure.*(Proof|Fallback|DryRun|JSON|Yes|Decline|NoMutation)' -count=1` — expected: all input modes prove candidates and every non-persisted path preserves config bytes.
- `rtk go test ./internal/config -run 'Test(WriteProfilesConfig|NormalizeProfilesFragment).*Fallback' -count=1` — expected: missing and duplicate fallback contracts remain deterministic and pre-proof.
- `rtk go test -race ./internal/cli -run 'TestProfilesConfigure.*(Proof|NoMutation)' -count=1` — expected: candidate proof and persistence boundaries are race-free.

## References

- `_prd.md` → User Stories 2, 3, and 8; Core Features 4, 9, and 10; Success
  Metrics.
- `_techspec.md` → Shared Readiness Coordinator; Setup Transaction Boundary;
  Error Taxonomy and Diagnostics; Build Order 7.
- `../../adr/0055-agent-selection-encoding-follows-advertised-acp-capabilities.md`
  → exact advertised selection proof.

## Result

Implemented profile configuration as a proposal transaction: Roundfix builds
and validates the candidate config without mutation, proves every distinct
candidate Preferred Selection and fallback through the shared readiness
coordinator, asks for confirmation when required, and only then persists the
proposal. Output failures roll the proposal back to the exact pre-command
snapshot.

Acceptance evidence:

- `TestProfilesConfigureProofRunsBeforeConfirmationAndWrite` proves both
  candidate tuples before confirmation observes the unchanged target, then
  verifies the authorized write.
- `TestProfilesConfigureFallbackFailurePrecedesProofAndPreservesBytes` and
  `TestNormalizeProfilesFragmentRequiresDistinctFallback` prove missing and
  duplicate fallbacks fail before any proof call and name the additional
  distinct authorized and proven Agent Selection requirement.
- `TestProfilesConfigureProofFailureYesJSONPreservesBytes` and
  `TestProfilesConfigureProofCleanupFailureAndDeclinePreserveBytes` prove
  classified readiness, cleanup, `--yes`, and decline paths preserve bytes and
  report `changed: false` in JSON.
- `TestProfilesConfigureDryRunAndFailedConfigurationLeaveBytesUnchanged`
  proves `--dry-run` executes both exact candidate proofs without writing.
- `TestProfilesConfigureJSONOutputFailureDoesNotMutateConfig` proves a JSON
  writer failure restores the byte-identical pre-command config.
- `TestProfilesConfigureInteractiveProofKeepsRecommendationsAdvisory` proves
  Interactive Input sends only the entered preferred and fallback selections
  to readiness and inserts no recommendation.

Verification:

- `rtk go test ./internal/cli -run 'TestProfilesConfigure.*(Proof|Fallback|DryRun|JSON|Yes|Decline|NoMutation)' -count=1` — passed, 14 tests.
- `rtk go test ./internal/config -run 'Test(WriteProfilesConfig|NormalizeProfilesFragment).*Fallback' -count=1` — passed, 3 tests.
- `rtk go test -race ./internal/cli -run 'TestProfilesConfigure.*(Proof|NoMutation)' -count=1` — passed, 9 tests.
- `rtk go test ./internal/config ./internal/cli -count=1` — passed.
- `rtk make verify` — passed.

Follow-ups: none.
