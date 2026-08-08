---
task: task_04
spec: 0088-a-third-runtime-that-can-run
status: pending
type: backend
complexity: medium
---

# Task 04: Prove every configured Agent Work Category

## Overview

`CONTEXT.md` requires Exact Agent Selection Proof for every distinct tuple, and
the Doctor Command resolves only the five required Agent Work Categories. A
profile configured for an optional category is parsed, registered, and never
proven, so `roundfix doctor` prints `profiles: ok` while that profile is failing
and never names its ACP Runtime. This Task widens readiness to what the effective
configuration actually defines.

## Requirements

1. MUST add a resolver returning every required Agent Work Category plus each
   optional category the effective configuration defines, and MUST exclude a
   category that only inherits `general`, because it contributes no distinct
   tuple.
2. MUST use that resolver for the Doctor Command's Agent Selection Profile
   Readiness, so a failing configured optional-category profile fails the check.
3. MUST use that resolver for the Doctor Command's ACP Runtime enumeration, so
   the adapter line names every runtime the configured tuples reference.
4. MUST use that resolver as the default scope of profile validation when no
   category is named, so both readiness surfaces answer the same question.
5. MUST use that resolver for the Setup Command's profile readiness over its
   proposed configuration, keeping the contract uniform; a generated proposal
   defines only the required five, so its behavior is unchanged.
6. MUST keep the reported counts honest: the tuple count stays a count of
   distinct tuples and the reference count stays a count of category references.
7. MUST NOT change the Run preflight's category derivation, which already follows
   the Task Graph's Task Types.
8. MUST re-record the coverage record in this Task's own commit if any test is
   renamed or removed.

## Subtasks

- [ ] Add the configured-category resolver beside the profile resolution it
      belongs to.
- [ ] Point the Doctor Command's profile readiness and runtime enumeration at it.
- [ ] Point profile validation's default scope at it.
- [ ] Point the Setup Command's profile readiness at it.
- [ ] Edit the break-half characterization test that pinned required-only scope,
      and declare the break in this Task's Result.

## Acceptance Criteria

- [ ] With a configured optional-category profile whose preferred selection
      fails, the Doctor Command reports the profiles check as failed and names
      that category among the affected categories.
- [ ] With a configured optional-category profile selecting a runtime no required
      category uses, the adapter line names that runtime.
- [ ] With no optional category configured, the Doctor Command's reported tuple
      and reference counts are unchanged from before this Task.
- [ ] A category present in configuration only by inheritance from `general` adds
      no tuple and no reference.
- [ ] Profile validation with no `--category` covers the same categories the
      Doctor Command covers.
- [ ] The Setup Command's readiness over a generated proposal behaves as before.

## Context

- interface: `internal/config/profiles.go`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/profiles_validate.go`
- interface: `internal/cli/setup.go`

## Bounded scope

This Task may create or modify only:

- `internal/config/profiles.go`
- `internal/config/config_test.go`
- `internal/cli/doctor.go`
- `internal/cli/doctor_test.go`
- `internal/cli/doctor_characterization_test.go`
- `internal/cli/profiles_validate.go`
- `internal/cli/setup.go`
- `internal/cli/cli_test.go`
- `docs/references/coverage-record.json`
- `docs/specs/0088-a-third-runtime-that-can-run/task_04.md`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go test ./internal/cli ./internal/config -count=1` — expected: exits 0.
- `go test ./internal/config -run 'ConfiguredWorkCategories' -count=1 -v` — expected: exits 0 and names at least one test; `no tests to run` fails this Task.
- `go test ./internal/cli -run 'DoctorProfile' -count=1 -v` — expected: exits 0 and names at least one test.
- `go test ./internal/spec -run '^TestCoverageEquivalence$' -count=1` — expected: exits 0.
- `grep -q 'ConfiguredWorkCategories' internal/cli/doctor.go` — expected: exits 0, proving the Doctor Command consumes the resolver rather than the required-only list.

## References

- `_prd.md` → Goal 3; Core Features: readiness that covers what is configured.
- `_techspec.md` → Implementation Design: Interfaces; Build Order 4.
- `references/2026-08-08-what-the-opencode-adapter-answers-before-its-first-prompt.md`
  → the measured `profiles: ok (5 distinct tuples; 10 category references)` over
  a failing configured profile.
- ADR-0107.
