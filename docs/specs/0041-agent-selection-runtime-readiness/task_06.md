---
task: task_06
spec: 0041-agent-selection-runtime-readiness
status: pending
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

- [ ] Update built-in and rendered generated profile defaults.
- [ ] Build effective Setup proposals before persistence.
- [ ] Prove adapter identity and each distinct proposed selection.
- [ ] Replace bare adapter override generation with supported provenance.
- [ ] Add explicit stale-override migration authorization.
- [ ] Enforce no-input, decline, failure, and cleanup no-mutation behavior.
- [ ] Cover successful User and Project Config persistence.

## Acceptance Criteria

- [ ] Built-ins and rendered default YAML use Sol/high plus GPT-5.5/xhigh for
      all required Codex profiles and never emit Terra/max.
- [ ] Terra and Luna remain valid Model Catalog and recommendation identifiers.
- [ ] Setup proves every distinct effective tuple once before writing ACPX,
      User Config, or Project Config.
- [ ] A legacy override produces one migration offer; declining it leaves the
      ACPX config and both Roundfix config scopes byte-identical.
- [ ] `--no-input`, unsupported tuple, invalid evidence, and cleanup failure
      create or change no target file.
- [ ] `--yes` cannot bypass adapter or exact profile proof.
- [ ] A successful proven proposal writes the authorized scopes and subsequent
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

