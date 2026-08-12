---
task: task_04
spec: 0001-implement-command
status: completed
type: backend
complexity: low
---

# Task 04: Detect the repository default branch in preflight

## Overview

Add default-branch detection to the git preflight inspection so the Implement Command can veto Runs started on the repository default branch (ADR-0013). Verifiable on its own through preflight unit tests with the existing fake git runner.

## Requirements

1. MUST detect the default branch from `refs/remotes/origin/HEAD` via the existing git runner wrapper.
2. MUST fall back to a `main`/`master` branch-name match when `origin/HEAD` is unset (fresh clones, no remote), and report when no default could be determined.
3. MUST return the detected default branch name so callers can compose a veto message that names both the current branch and the detected default — a false positive must be diagnosable from the message alone.
4. MUST NOT change any existing preflight behavior or output for fetch, resolve, watch, or stop.

## Subtasks

- [x] Default-branch probe through the git runner
- [x] Name-match fallback and the undetermined case
- [x] Unit tests over the fake git runner covering all three outcomes

## Acceptance Criteria

- [x] With `origin/HEAD` pointing at a branch, detection returns that branch regardless of its name.
- [x] Without `origin/HEAD`, a repository whose current branch is `main` or `master` is detected as being on the default; any other branch is not vetoed.
- [x] The detection result carries the default branch name (or the undetermined marker) for use in the veto message.
- [x] The full existing test suite passes unchanged.

## Verification

- `rtk go test ./internal/preflight/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 4; Core Feature 2. `_techspec.md` → System Architecture (preflight), Build Order 4, Risks (default-branch detection). ADR-0013.

## Result

Default-branch detection now lives in `internal/preflight` as a standalone
helper on the existing `GitRunner` seam. `DetectDefaultBranch(ctx, workDir,
currentBranch, runner)` probes `git symbolic-ref refs/remotes/origin/HEAD`;
when that ref resolves, its branch name is authoritative regardless of what
it is called. When the probe errors or returns a malformed target, a current
branch named `main` or `master` becomes the default by name match, and any
other branch yields the undetermined result. The `DefaultBranch` result
carries `Name` plus a `Source` (`origin/HEAD`, `name-match`, `undetermined`),
so the task_06 veto message can name both the current branch and the detected
default and a false positive is diagnosable from the message alone.
`DefaultBranch.IsDefault(branch)` encapsulates the veto comparison and never
matches on an undetermined detection. No existing preflight code path was
touched: the helper is a new file and nothing calls it yet, so fetch,
resolve, watch, and stop behavior is byte-identical.

Commands run:

- `rtk go test ./internal/preflight/` — 25 tests passed.
- `rtk go test ./...` — 333 tests passed in 16 packages.
- `make verify` — passed (full suite, `roundfix skills check`, build).
- `rtk gofmt -l internal/preflight/` — no output; formatting clean.

Evidence per acceptance criterion:

1. `origin/HEAD` set → that branch wins: subtests "origin/HEAD names the
   default regardless of branch name" (`trunk`), "origin/HEAD wins over a
   main current branch", and "origin/HEAD keeps slashes in the default
   branch name" in `TestDetectDefaultBranch`.
2. `origin/HEAD` unset → name match only: subtests "unset origin/HEAD
   matches main by name", "unset origin/HEAD matches master by name", and
   "unset origin/HEAD leaves a feature branch undetermined", plus the
   `TestDefaultBranchIsDefault` cases showing non-default branches never
   match.
3. Result carries the name or the undetermined marker: every table case
   asserts the full `DefaultBranch{Name, Source}` value, including
   `Source: DefaultBranchUndetermined` with an empty `Name`.
4. Existing suite unchanged: `rtk go test ./...` green with zero edits to
   `preflight.go` or `preflight_test.go`.

Follow-ups:

- task_06 wires the veto using `DetectDefaultBranch` + `IsDefault` and
  composes the message from `Name`/`Source`.
- `CONTEXT.md` has no glossary entry for "default branch"; if the veto copy
  needs a canonical term, add one rather than inventing per-command wording.
