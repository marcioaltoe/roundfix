---
task: task_03
spec: 0301-the-light-tier-yields-to-the-category-profile
status: completed
type: backend
complexity: medium
---

# Task 03: Doctor warns when the Light Tier cannot serve full access

## Overview

No Doctor check reads `openrouter.light_models` or
`defaults.agent_full_access`, so an operator learns that the light runtime
cannot serve full access only when a Run skips the Light Tier. This Task adds a
`light-tier:` line, printed after `recommendations:` only when the light model
list is not empty and full access is requested, that warns with
`DR-LIGHT-ACCESS-MODE` when no light runtime offers a full-access mode.

## Requirements

1. MUST implement API Contract 4 of `_techspec.md`: add the constant
   `HealthCheckLightTier = "light-tier"` to `internal/cli/health.go`, and
   `doctorLightTierResult` to `internal/cli/doctor.go` with the signature of
   `_techspec.md` → Interfaces, calling `lighttier.ServingSelections`.
2. MUST print, through `readinessResult`, exactly `light-tier: warn
   (DR-LIGHT-ACCESS-MODE: light runtime opencode offers no full-access mode;
   light Tasks run on their category profile; next: set
   openrouter.light_models: [] in User Config, or defaults.agent_full_access:
   false)` when no light selection is kept, and `light-tier: ok (light
   runtime <runtime> serves full access)` when one is, naming the first kept
   selection's runtime.
3. MUST make `runDoctorCommand` append the result right after the
   `recommendations:` result and before the `pre-pr-review:` result, only when
   `doctorLightTierResult` reports it, passing a factory over
   `runtimeForProfileSelectionWithOptions` with `EnableFullAccess` set from
   the loaded configuration's `defaults.agent_full_access`.
4. MUST state the hazard in the code's behavior: the check reads only the
   loaded configuration and the static runtime capabilities; it opens no
   Agent Session, starts no process, reads no key, never reports `failed` and
   never changes Doctor's exit code.
5. MUST keep every existing Doctor line byte-identical, and leave
   `internal/cli/profiles_validate.go`, `internal/lighttier` and
   `internal/daemon` unchanged. No existing test, golden or fixture changes.
6. MUST add, in the new file `internal/cli/doctor_light_tier_test.go`, each
   top-level test calling `t.Parallel()` first, compiling for Windows and
   passing no real-time duration:
   - `TestDoctorLightTierResultNamesTheLightRuntimeThatCannotServeFullAccess`,
     with the cases of `_techspec.md` → Testing Approach: not printed with
     full access off or an empty light model list; the exact `warn` line with
     the static factory; the exact `ok` line with a fake factory that gives
     `opencode` a full-access candidate;
   - `TestDoctorPrintsTheLightTierLineOnlyWithFullAccessAndLightModels`,
     through `roundfix doctor` with the fake checker and profile readiness the
     `linux-lane` line test uses: no `light-tier:` line with the built-in
     configuration; with `defaults.agent_full_access` `true` the `warn` line
     right after the `recommendations:` line and before `pre-pr-review:`,
     the output otherwise byte-identical to the built-in run, exit 0; and no
     line with an empty light model list.

## Subtasks

- [ ] Add the `light-tier` check name and `doctorLightTierResult`.
- [ ] Print the line after `recommendations:` only for the combination.
- [ ] Add the result and command tests.

## Acceptance Criteria

- [ ] With full access and the default light model, Doctor prints the exact `warn` line after `recommendations:` and exits 0.
- [ ] With full access off or no light model, Doctor prints no `light-tier:` line and every other line is unchanged.
- [ ] The existing Doctor tests named in the Verification pass unchanged.

## Consumers

- `doctorLightTierResult` and `HealthCheckLightTier` — consumers: `./internal/cli`

## Context

- instruction: `docs/adr/0314-a-light-tier-that-cannot-serve-the-run-yields-to-the-category-profile.md`
- instruction: `internal/cli/readiness.go`
- instruction: `internal/cli/profiles_validate.go`
- instruction: `internal/cli/doctor_linux_lane_test.go`
- instruction: `internal/cli/doctor_test.go`
- instruction: `docs/user-guide/commands/doctor.md`
- creates: `internal/cli/doctor_light_tier_test.go`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/health.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestDoctorLightTierResultNamesTheLightRuntimeThatCannotServeFullAccess|TestDoctorPrintsTheLightTierLineOnlyWithFullAccessAndLightModels|TestDoctorPrintsTheLinuxLaneLineOnlyWhenTheLaneIsOn|TestRunDoctorProfileReadinessProvesEffectiveCategoriesAndReportsCounts|TestRunDoctorMissingRepositoryRoot|TestDoctorReportsPrePRReviewPolicy)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDoctorLightTierResultNamesTheLightRuntimeThatCannotServeFullAccess TestDoctorPrintsTheLightTierLineOnlyWithFullAccessAndLightModels TestDoctorPrintsTheLinuxLaneLineOnlyWhenTheLaneIsOn TestRunDoctorProfileReadinessProvesEffectiveCategoriesAndReportsCounts TestRunDoctorMissingRepositoryRoot TestDoctorReportsPrePRReviewPolicy; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go vet ./internal/cli` — expected: exit 0. Before this Task the two new tests and `doctorLightTierResult` do not exist, so the command fails at its first pass check. Afterwards the light-tier result and command tests pass with the unchanged Doctor tests, and the package's tests compile for Windows.

## References

- `_prd.md` → Goals; Core Feature 3; Success Metric 6
- `_techspec.md` → API Contracts 4-5; Invariant 6; Testing Approach; Build Order 3
- ADR-0314

## Result

Added the conditional `light-tier:` readiness result immediately after
`recommendations:` and before `pre-pr-review:`. The command builds static
runtime capabilities with `EnableFullAccess` from the loaded configuration,
and the result calls `lighttier.ServingSelections`. The check opens no Agent
Session, starts no process, reads no key, and returns only `ok` or `warn`.

Acceptance evidence from focused checks:

- Full access with the built-in light model: the new result test asserts the
  exact `DR-LIGHT-ACCESS-MODE` warning; the command test asserts that line's
  position, exactly one occurrence, empty stderr and exit 0. The fake serving
  runtime case asserts the exact `ok` line naming `opencode`.
- Full access off or an empty light model list: both new tests assert the
  absence of the line. The result test also rejects any runtime-factory call
  for these cases. Removing the warning from the full-access command output
  leaves it byte-identical to the built-in run; the empty-model run also
  matches that baseline byte for byte.
- Existing Doctor behavior: the four existing tests named below pass without
  edits. `profiles_validate.go`, `internal/lighttier`, `internal/daemon`, all
  existing tests, goldens and fixtures are unchanged by this Task.

Commands and outcomes:

- Before implementation, `go test ./internal/cli -run
  '^TestDoctorLightTierResultNamesTheLightRuntimeThatCannotServeFullAccess$'
  -count=1` failed to compile with `undefined: doctorLightTierResult`, proving
  the missing seam with the new tests in place.
- `go test ./internal/cli -run
  '^(TestDoctorLightTierResultNamesTheLightRuntimeThatCannotServeFullAccess|TestDoctorPrintsTheLightTierLineOnlyWithFullAccessAndLightModels)$'
  -count=1` — exit 0.
- `go test ./internal/cli -run
  '^(TestRunDoctorProfileReadinessProvesEffectiveCategoriesAndReportsCounts|TestRunDoctorMissingRepositoryRoot|TestDoctorReportsPrePRReviewPolicy|TestDoctorPrintsTheLinuxLaneLineOnlyWhenTheLaneIsOn)$'
  -count=1` — exit 0.
- `GOOS=windows go test -c -o /tmp/roundfix-task03-cli.test.exe
  ./internal/cli` — exit 0; package tests compile for Windows.
- `git diff --check` — exit 0.

Environment: the executor refused the symlinked writable workspace root;
commands ran in its resolved `/Volumes/Brain-SSD/cache/.roundfix/worktrees/`
checkout with approved escalation. No repository configuration changed.

The declared Verification (including Windows vet) and repository Verification
were not run; the Daemon owns those checks and Task settlement. Task status
and authored checkboxes remain unchanged. No follow-up work was identified.
