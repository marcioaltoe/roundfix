---
task: task_02
spec: 0061-repository-derived-skill-requirements
status: completed
type: backend
complexity: medium
---

# Task 02: Derive the requirement from the repository's Setup Manifest

## Overview

Resolve which external skills a repository requires from what it actually
selected: the modules recorded in its Setup Manifest and the `requiredSkills`
those modules declare in the embedded Baseline catalog. Wire the Doctor path to
that resolution so a TypeScript repository stops being told it needs Go skills.

## Requirements

1. MUST read the repository's Setup Manifest and take the modules it records
   as the selected set.
2. MUST union the `requiredSkills` those modules declare in the embedded
   Baseline catalog, subtract the Roundfix-owned skill names, and produce a
   deduplicated set in deterministic order.
3. MUST report zero external requirements when the Setup Manifest is absent or
   unreadable, and name Baseline adoption as the next action rather than
   falling back to the embedded recommendation list.
4. MUST fail legibly when the manifest names a module the catalog does not
   know, naming that module, instead of silently requiring nothing.
5. MUST keep the resolver outside the embedded-bundle package so it does not
   depend on the Baseline catalog.
6. MUST keep the Doctor Command diagnosis-only: no file, config, or repository
   state is written.

## Subtasks

- [ ] Add the resolver reading the Setup Manifest and the catalog modules.
- [ ] Subtract owned names and order the result deterministically.
- [ ] Wire the Doctor readiness path to the resolver.
- [ ] Cover the TypeScript, Go/TUI, absent-manifest, and unknown-module cases.

## Acceptance Criteria

- [ ] A repository whose manifest selects the TypeScript modules requires no
      `golang-*`, `bubbletea`, or `tui-design` skill.
- [ ] A repository whose manifest selects `go`, `cli-surface`, and
      `tui-surface` still requires those skills.
- [ ] A repository with no Setup Manifest reports zero external requirements
      and a next action naming Baseline adoption.
- [ ] A manifest naming an unknown module fails with that module named.
- [ ] The embedded-bundle package still does not import the Baseline catalog.
- [ ] Running Doctor writes nothing.

## Context

- interface: `internal/cli/doctor.go`
- interface: `internal/cli/doctor_test.go`
- interface: `internal/baseline/catalog.go`
- interface: `internal/baseline/plan.go`

## Verification

- `go build -buildvcs=false ./...` — expected: clean build.
- `go test -count=1 ./internal/cli/ -run 'TestRunDoctor'` — expected: pass, including the four derivation cases.
- `grep -rn 'roundfix/internal/baseline' skills/*.go | grep -v _test ; test $? -eq 1` — expected: no matches; the bundle package stays a leaf.

## References

`_prd.md` → User Stories 1, 3, 4, Core Features 1, 3, 5; `_techspec.md` →
Build Order 2, Interfaces: resolveExternalSkillRequirement.

## Result

### Implementation

- Doctor now resolves external skill requirements from
  `docs/agents/setup-context.json` and the embedded Baseline catalog, removes
  names owned by the Roundfix bundle, deduplicates them, and sorts the result
  before repository readiness checks.
- An absent, non-regular, unreadable, or malformed Setup Manifest produces an
  empty external requirement. Doctor reports `0 external`, fails readiness
  with the Setup Manifest diagnosis, and names `roundfix baseline` as the next
  action.
- An unknown catalog module stops requirement resolution and appears by name
  in the `skills:` failure. The existing external remediation behavior remains
  unchanged for Task 03.

### Focused checks

- Pre-change signal:
  `GOCACHE="$PWD/.gocache" rtk go test ./internal/cli -run "^TestResolveExternalSkillRequirement$"`
  reached compilation and failed because
  `resolveExternalSkillRequirement` did not exist.
- `GOCACHE="$PWD/.gocache" rtk go test ./internal/cli -run "^Test(ResolveExternalSkillRequirementUnreadableManifest|RunDoctorDerivesExternalSkillRequirementFromSetupManifest|RunDoctorRealRepositoryCheckDoesNotMutateState|RunDoctorRepositorySkillReadiness|RunDoctorPassesCommandContextToRepositorySkillReadiness|RunDoctorMissingRepositoryRoot)$"`:
  passed, 16 tests.
- `GOCACHE="$PWD/.gocache" rtk go test ./internal/cli -run "^TestRunDoctor(ProfileReadinessProvesEffectiveCategoriesAndReportsCounts|ProfileReadinessMatchesProfilesValidateFailureEvidence)$"`:
  passed, 5 tests.
- `rtk rg -n 'roundfix/internal/baseline' skills --glob '*.go'`: no matches
  (expected exit 1).
- `rtk git diff --check`: passed.
- The Task's declared Verification commands were not run; Daemon Verification
  owns them.

### Acceptance criteria evidence

- TypeScript selection: the
  `TypeScript_modules_exclude_Go_and_TUI_skills` Doctor subtest passed and
  also proved sorted, deduplicated output with the owned `evidence-gate` name
  removed.
- Go/TUI selection: the
  `Go_CLI_and_TUI_modules_retain_their_skills` Doctor subtest passed with the
  exact nine-skill external set.
- No Setup Manifest: the
  `absent_manifest_requires_Baseline_adoption_but_no_external_skills` Doctor
  subtest passed with `0 external` and `next: roundfix baseline`; the focused
  resolver test also covered malformed manifest input.
- Unknown module: the
  `unknown_module_fails_before_repository_readiness` Doctor subtest passed and
  asserted that `unknown-module` appears before the repository checker runs.
- Bundle leaf: the focused import search found no Baseline import under
  `skills/*.go`; the new dependency remains in `internal/cli`.
- Diagnosis-only: `TestRunDoctorRealRepositoryCheckDoesNotMutateState` passed
  after exercising the real resolver and
  `skills.CheckRepositoryWithExternal`.
