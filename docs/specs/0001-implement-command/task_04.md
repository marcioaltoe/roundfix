---
task: task_04
spec: 0001-implement-command
status: pending
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

- [ ] Default-branch probe through the git runner
- [ ] Name-match fallback and the undetermined case
- [ ] Unit tests over the fake git runner covering all three outcomes

## Acceptance Criteria

- [ ] With `origin/HEAD` pointing at a branch, detection returns that branch regardless of its name.
- [ ] Without `origin/HEAD`, a repository whose current branch is `main` or `master` is detected as being on the default; any other branch is not vetoed.
- [ ] The detection result carries the default branch name (or the undetermined marker) for use in the veto message.
- [ ] The full existing test suite passes unchanged.

## Verification

- `rtk go test ./internal/preflight/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 4; Core Feature 2. `_techspec.md` → System Architecture (preflight), Build Order 4, Risks (default-branch detection). ADR-0013.
