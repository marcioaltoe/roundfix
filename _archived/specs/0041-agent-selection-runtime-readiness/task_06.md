---
task: task_06
spec: 0041-agent-selection-runtime-readiness
status: completed
type: backend
complexity: high
---

# Task 06: Generate and persist only proven Setup profiles

## Overview

Make Setup build adapter and profile changes in memory, prove the resulting
effective selections, and persist only after proof and authorization. Generated
Codex profiles must use Sol/high with GPT-5.5/xhigh fallback while retaining
the established frontend policy and official model catalog.

## Requirements

1. MUST generate Sol/high as Preferred and GPT-5.5/xhigh as fallback for
   `general`, `backend`, `qa`, and `review`; frontend policy remains unchanged.
2. MUST keep Sol, Terra, and Luna as accepted official identifiers without
   emitting Terra or Luna as generated operational defaults.
3. MUST construct proposed ACPX, User Config, and Project Config bytes in
   memory and resolve their effective profiles before writing any scope.
4. MUST prove adapter identity and every distinct generated Preferred Selection
   and fallback before offering or applying persistence.
5. MUST propose a deterministic supported official adapter command instead of
   creating a bare PATH override.
6. MUST require confirmation before migrating an existing stale adapter
   override and leave all bytes unchanged when declined.
7. MUST make `--no-input`, proof failure, cleanup failure, and any write error
   preserve every not-yet-committed target without silently choosing another
   model or reasoning effort.

## Subtasks

- [x] Update built-in and rendered generated profile defaults.
- [x] Build effective Setup proposals before persistence.
- [x] Prove adapter identity and each distinct proposed selection.
- [x] Replace bare adapter override generation with supported provenance.
- [x] Add explicit stale-override migration authorization.
- [x] Enforce no-input, decline, failure, and cleanup no-mutation behavior.
- [x] Cover successful User and Project Config persistence.

## Acceptance Criteria

- [x] Built-ins and rendered default YAML use Sol/high plus GPT-5.5/xhigh for
      all required Codex profiles and never emit Terra/max.
- [x] Terra and Luna remain valid Model Catalog and recommendation identifiers.
- [x] Setup proves every distinct effective tuple once before writing ACPX,
      User Config, or Project Config.
- [x] A legacy override produces one migration offer; declining it leaves the
      ACPX config and both Roundfix config scopes byte-identical.
- [x] `--no-input`, unsupported tuple, invalid evidence, and cleanup failure
      create or change no target file.
- [x] `--yes` cannot bypass adapter or exact profile proof.
- [x] A successful proven proposal writes the authorized scopes and subsequent
      profile validation observes the same effective tuples.

## Context

- instruction: `docs/agents/autonomous-work.md`
- interface: `internal/cli/setup.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/config/profiles.go`
- interface: `internal/config/config.go`
- interface: `internal/config/config_test.go`

## Verification

- `rtk go test ./internal/config -run 'Test(BuiltinProfiles|DefaultConfigYAML|ModelCatalog)' -count=1` — expected: Sol/high and GPT-5.5/xhigh defaults pass while Sol, Terra, and Luna remain valid catalog entries.
- `rtk go test ./internal/cli -run 'TestRunSetup.*(Profile|Adapter|NoInput|Decline|Cleanup|NoMutation)' -count=1` — expected: proof precedes persistence and every failure/decline path preserves target bytes.
- `rtk go test -race ./internal/cli ./internal/config -run 'Test(RunSetup|BuiltinProfiles|DefaultConfigYAML)' -count=1` — expected: Setup proposal, proof, and persistence boundaries are race-free.

## References

- `_prd.md` → User Stories 1, 2, and 9; Core Features 1, 5, 6, and 10; Success
  Metrics.
- `_techspec.md` → Generated Defaults and Official Identifiers; Adapter
  Provisioning and Identity; Setup Transaction Boundary; Build Order 6.
- `references/validation.md` → proven official adapter and model evidence.

## Result

Setup now resolves ACPX, User Config, and Project Config proposals entirely in
memory, proves the effective adapter and three distinct generated Agent
Selections, collects authorization, and then persists each authorized target
atomically. Legacy Codex overrides migrate only after one explicit offer to the
pinned official command; decline, `--no-input`, unsupported tuples, malformed
capability evidence, Agent Session cleanup failure, and write failure preserve
every not-yet-committed target.

Acceptance evidence:

- Built-in and rendered `general`, `backend`, `qa`, and `review` profiles use
  Sol/high with GPT-5.5/xhigh fallback. Generated YAML contains no Terra, Luna,
  or `max` operational default, while the Model Catalog still contains Sol,
  Terra, and Luna.
- `TestRunSetupProfileProofsEveryDistinctTupleOnceBeforePersistence` observes
  one proof each for Sol/high, GPT-5.5/xhigh, and Claude Fable/medium before
  User Config or Project Config persistence.
- `TestRunSetupAdapterMigrationDeclinePreservesAllTargets` observes one
  migration prompt and byte-identical ACPX, User Config, and Project Config
  state after decline. `TestRunSetupAdapterMigrationPersistsSupportedCommand`
  proves and writes `npx -y @agentclientprotocol/codex-acp@1.1.4` rather than a
  bare PATH override.
- Setup failure tests cover `--yes`, `--no-input`, unsupported selection,
  invalid capability evidence, cleanup failure, and atomic write failure with
  no unauthorized target mutation.
- `TestRunSetupProfilePersistenceMatchesSubsequentValidation` writes User
  Config and Project Config proposals, resolves the persisted precedence, and
  observes the same three distinct tuples through the shared validation path.

Verification:

- `rtk go test ./internal/config -run 'Test(BuiltinProfiles|DefaultConfigYAML|ModelCatalog)' -count=1` — passed, 3 tests.
- `rtk go test ./internal/cli -run 'TestRunSetup.*(Profile|Adapter|NoInput|Decline|Cleanup|NoMutation)' -count=1` — passed, 16 tests.
- `rtk go test -race ./internal/cli ./internal/config -run 'Test(RunSetup|BuiltinProfiles|DefaultConfigYAML)' -count=1` — passed, 27 tests.
- `rtk make verify` — passed.

Follow-ups: none.
