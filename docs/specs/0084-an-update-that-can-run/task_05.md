---
task: task_05
spec: 0084-an-update-that-can-run
status: pending
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

- [ ] Carry the searched locations out of profile resolution.
- [ ] Distinguish a missing repository-owned profile from an unknown identity.
- [ ] Compose the diagnosis message with identity, locations, and action.
- [ ] Report it on both output surfaces.
- [ ] Cover the missing repository-owned profile case.
- [ ] Cover the unknown catalog identity case.

## Acceptance Criteria

- [ ] A repository whose Setup Manifest names a repository-owned profile absent
      from the checkout reports the profile identity, the searched paths, and the
      restoring action.
- [ ] A repository whose Setup Manifest names an identity the catalog does not
      know reports a distinct action naming adoption.
- [ ] Neither message contains a raw `lstat` or `open` error string as its whole
      text.
- [ ] Both cases keep the exit code and the classification the command emitted
      before this task.
- [ ] The JSON result carries the same identity, locations, and action as the text
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
