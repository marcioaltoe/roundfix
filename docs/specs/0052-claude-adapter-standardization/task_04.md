---
task: task_04
spec: 0052-claude-adapter-standardization
status: pending
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
