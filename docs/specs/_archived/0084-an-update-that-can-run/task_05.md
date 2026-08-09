---
task: task_05
spec: 0084-an-update-that-can-run
status: completed
type: backend
complexity: low
---

# Task 05: Make an unresolved Baseline Profile diagnose itself

## Overview

A repository whose Setup Manifest names a Baseline Profile the checkout can no
longer resolve currently fails with a raw filesystem error naming an internal
path. It is the sixth of the eight measured refusals and the only one whose
repair the maintainer could perform immediately if the message said what was
missing. This slice replaces that error with a diagnosis: the profile identity,
where the command looked, and the action that restores it.

## Requirements

1. MUST report, when a recorded Baseline Profile does not resolve, the profile
   identity taken from the Setup Manifest.
2. MUST report every location the command searched for that profile, as
   repository-relative paths.
3. MUST report the action that restores the repository, distinguishing a
   repository-owned profile the checkout is missing from a profile identity the
   catalog does not know.
4. MUST NOT leak a raw filesystem error string as the whole message.
5. MUST keep the failing exit code and the manifest-incompatible classification
   unchanged, because the condition still requires adoption rather than a refresh.
6. MUST report the same diagnosis in both the text and JSON surfaces.

## Subtasks

- [x] Carry the searched locations out of profile resolution.
- [x] Distinguish a missing repository-owned profile from an unknown identity.
- [x] Compose the diagnosis message with identity, locations, and action.
- [x] Report it on both output surfaces.
- [x] Cover the missing repository-owned profile case.
- [x] Cover the unknown catalog identity case.

## Acceptance Criteria

- [x] A repository whose Setup Manifest names a repository-owned profile absent
      from the checkout reports the profile identity, the searched paths, and the
      restoring action.
- [x] A repository whose Setup Manifest names an identity the catalog does not
      know reports a distinct action naming adoption.
- [x] Neither message contains a raw `lstat` or `open` error string as its whole
      text.
- [x] Both cases keep the exit code and the classification the command emitted
      before this task.
- [x] The JSON result carries the same identity, locations, and action as the text
      output.

## Context

- interface: `internal/baseline/update.go`
- interface: `internal/cli/baseline_update.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exits 0.
- `go test ./internal/baseline/ ./internal/cli/ -run 'UnresolvedProfile' -v > /tmp/0084-task-05-a.log 2>&1 && grep -q '^--- PASS: .*UnresolvedProfile' /tmp/0084-task-05-a.log` — expected: exits 0, proving the diagnosis cases exist and pass rather than being selected out.
- `go test ./internal/baseline/ ./internal/cli/ -count=1` — expected: exits 0.

## References

- `_techspec.md` → Build Order 5; API Contracts.
- `_prd.md` → Core Feature 6; User Story 5; Goal 1.
- `references/2026-08-08-the-update-refuses-six-of-the-eight-copies-it-exists-to-update.md`
  → the measured repository that produces this failure.

## Result

### Implementation

- Profile resolution now carries every repository-relative location consulted
  after an identity misses the embedded catalog.
- Setup Manifest resolution projects a typed unresolved-profile diagnosis with
  the recorded identity, searched locations, and repair action while preserving
  the underlying error chain for programmatic inspection.
- The update command emits the diagnosis as structured JSON under
  `unresolvedProfile` and renders the same identity, locations, and action in
  text output. Its user-facing message no longer contains the raw filesystem
  error.

### Focused checks

- Red signal: with the worktree-local Go cache, the new focused tests initially
  failed to compile because `UnresolvedProfileKind`, the diagnosis constants,
  and the result fields did not exist.
- `rtk go test ./internal/baseline -run '^TestResolveManifestInputUnresolvedProfileDiagnosis$' -count=1`
  with the worktree-local `GOCACHE`: passed, 3 tests.
- `rtk go test ./internal/cli -run '^TestBaselineUpdateUnresolvedProfileDiagnosis$' -count=1`
  with the worktree-local `GOCACHE`: passed, 3 tests.
- `rtk go test ./internal/baseline -run '^(TestResolveManifestInput|TestCustomProfileRejectsUnsafePathsAndNonRepositorySources)' -count=1`
  with the worktree-local `GOCACHE`: passed, 18 tests.
- `rtk go test ./internal/cli -run '^TestBaselineUpdate(NoManifest|UnresolvedProfile|NewDecision)' -count=1`
  with the worktree-local `GOCACHE`: passed, 5 tests.
- `rtk git -c core.fsmonitor=false diff --check`: passed.
- `rtk gofmt -d internal/baseline/custom_profile.go internal/baseline/update.go internal/baseline/update_test.go internal/cli/baseline_update.go internal/cli/baseline_update_test.go`:
  passed with no output.

### Acceptance evidence

1. `TestResolveManifestInputUnresolvedProfileDiagnosis/missing_repository-owned_Profile`
   and the matching CLI subtest assert the Setup Manifest identity,
   `.roundfix/baseline/profiles/repository-backend.json`, and a restore action.
2. The `unknown_catalog_identity` subtests assert the distinct
   `catalog_identity_unknown` diagnosis and an action that names adoption.
3. Both owning suites reject `lstat` and `open` leakage from the diagnosis;
   the CLI suite checks the JSON message plus text stdout and stderr.
4. The package test preserves `ErrManifestIncompatible`,
   `ManifestInputIncompatible`, and `profile_unresolved`; the CLI test preserves
   exit `2`, state `failed`, and category `manifest` for both cases.
5. The CLI test decodes the JSON identity, searched locations, and action, then
   requires those same values in the text surface.

The Daemon-owned commands under `## Verification` were not run in this Agent
turn.
