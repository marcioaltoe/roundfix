---
task: task_04
spec: 0052-claude-adapter-standardization
status: completed
type: backend
complexity: medium
---

# Task 04: Report multi-runtime Adapter Readiness in Doctor

## Overview

Make the Doctor Command's `adapter:` check inspect every distinct ACP Runtime
referenced by the effective required Agent Selection Profiles instead of only
the `general` category's preferred runtime. With the frontend built-in now on
claude, a default machine gets Claude lineage evidence — or a failed check
with the official install action — from `roundfix doctor` alone.

## Requirements

1. MUST collect the distinct runtimes across every required category's
   Preferred Selection and Fallback Chain, deduplicated and ordered by
   runtime ID, and run the adapter check for each.
2. MUST keep a single `adapter:` report line whose detail carries
   per-runtime evidence (command, package, version where proven), in
   deterministic runtime order.
3. MUST fail the line when any runtime's check fails, naming the failing
   runtime and its install next action; the command's exit behavior is
   unchanged.
4. MUST keep the profile-proof evidence surface (`profiles:` line) working
   with the runtime-parameterized lineage and version errors.
5. SHOULD keep the existing single-runtime detail wording recognizable so
   operators and the documented examples migrate cleanly.

## Subtasks

- [ ] Resolve the required-profile runtime set for the adapter check.
- [ ] Render per-runtime evidence in the `adapter:` detail with
      deterministic order.
- [ ] Surface the failing runtime's classification and install action as the
      line's next action.
- [ ] Extend the Doctor tests: both-runtimes evidence, claude legacy lineage
      failing with the official action, deterministic check line-up
      preserved.

## Acceptance Criteria

- [ ] On built-in profiles, `roundfix doctor` reports adapter evidence for
      both `claude` and `codex` in one `adapter:` line, ordered by runtime
      ID.
- [ ] A machine whose claude command resolves to a legacy lineage reports
      `adapter: failed` with classification `adapter_lineage_unknown` and
      `next:` naming the official Claude install command, while codex
      evidence still appears.
- [ ] A fully official machine reports `adapter: ok` with package and
      version for both runtimes.
- [ ] The ordered Doctor check sequence (`node`, `acpx`, `adapter`,
      `profiles`, `skills`, `codex`) is unchanged.

## Context

- interface: `internal/cli/doctor.go`
- interface: `internal/cli/health.go`
- interface: `internal/cli/doctor_test.go`
- interface: `internal/cli/profiles_validate.go`

## Verification

- `go build -buildvcs=false ./...` — expected: clean build.
- `go test -count=1 ./internal/cli/ -run 'TestRunDoctor'` — expected: pass, including the new multi-runtime adapter cases.

## References

`_prd.md` → User Story 1, Core Feature 2, User Experience; `_techspec.md` →
Build Order 4, API Contracts; ADR-0055.

## Result

### Implementation

- Doctor now derives its Adapter Readiness runtime set from the effective
  required Agent Selection Profiles, including every Preferred Selection and
  Fallback Chain entry. It deduplicates by runtime ID and checks the sorted
  runtime list.
- The command keeps one `adapter:` line and joins each runtime's health
  evidence in runtime-ID order. It continues checking later runtimes after a
  failure, carries typed adapter classifications into the failing runtime's
  detail, and prints the first failing runtime's install action through the
  existing `next:` field.
- Adapter health failures now retain their typed error in `CheckResult`, which
  lets Doctor classify generalized Claude or Codex lineage/version errors
  without changing the existing profile-proof reporting path.
- Doctor tests cover official Claude and Codex evidence, legacy Claude
  lineage failure with retained Codex evidence, a runtime referenced only by
  a Fallback Chain, deduplication, sorted runtime order, and the unchanged six
  report-line sequence.

### Focused checks

- The first test attempt used Go's default cache and could not start because
  the sandbox denied writes under `~/Library/Caches/go-build`.
- Pre-change signal:
  `rtk proxy env GOCACHE=/private/tmp/roundfix-task04-gocache go test ./internal/cli -run '^TestRunDoctorAdapterReadinessReportsRequiredProfileRuntimes$' -count=1`
  exited 1. Doctor checked only Codex, omitted Claude evidence, and returned
  exit 0 for the legacy-Claude case.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task04-gocache go test ./internal/cli -run '^(TestRunDoctorAdapterReadinessReportsRequiredProfileRuntimes|TestRunDoctorAdapterReadinessIncludesFallbackOnlyRuntime)$' -count=1`
  — passed.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task04-gocache go test ./internal/cli -run '^(TestRunDoctorProfileReadinessProvesEffectiveCategoriesAndReportsCounts|TestRunDoctorProfileReadinessReportsLegacyAdapterThroughEffectiveProfile|TestRunDoctorProfileReadinessMatchesProfilesValidateFailureEvidence|TestRunDoctorContinuesChecksAfterProfileReadinessFailure|TestRunDoctorRepositorySkillReadiness|TestRunDoctorMissingRepositoryRoot|TestRunDoctorRealRepositoryCheckDoesNotMutateState)$' -count=1`
  — passed.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task04-gocache go test ./internal/cli -run '^TestHealthCheckerReportsFailedPrerequisitesWithNextActions$' -count=1`
  — passed.
- `rtk git diff --check` — passed.
- The Task's declared `## Verification` commands were not run; the Daemon
  owns that gate.

### Acceptance-criterion evidence

1. `TestRunDoctorAdapterReadinessReportsRequiredProfileRuntimes` observes one
   `adapter:` line containing `claude` then `codex`, and records exactly those
   runtime calls in sorted order. The fallback-only companion adds OpenCode
   only through `general.fallbacks` and observes `claude`, `codex`,
   `opencode` once each.
2. The legacy-Claude case observes `adapter: failed`, classification
   `adapter_lineage_unknown`, the official
   `ClaudeAdapterInstallCommand()`, and the successful Codex package/version
   evidence on the same line.
3. The official-adapters case observes `adapter: ok` with command, official
   package, and pinned version evidence for both Claude and Codex.
4. The same table asserts the six output line names remain `node`, `acpx`,
   `adapter`, `profiles`, `skills`, `codex`; the adjacent Doctor checks also
   pass with two internal adapter probes and the existing profile-proof
   failure evidence.

### Follow-ups

None discovered within this Task's slice.
