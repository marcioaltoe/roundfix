---
task: task_04
spec: 0155-a-verify-that-runs-what-changed
status: completed
type: infra
complexity: low
---

# Task 04: Pull requests run the selective gate

## Overview

Pull requests should pay only for what they change; `main` and releases keep the
complete gate.

## Requirements

1. MUST make the pull request job in `.github/workflows/ci-verify.yml` run
   `make verify-changed` against the pull request's base branch.
2. MUST keep pushes to `main` running `make verify`.
3. MUST leave `.github/workflows/release.yml` running `make verify`.
4. MUST keep the full-history checkout the selector needs.

## Subtasks

- [ ] Split the pull request and push paths of the verify job.

## Acceptance Criteria

- [ ] The pull request path names `make verify-changed` with the base branch.
- [ ] The push-to-main path and the release workflow name `make verify`.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `.github/workflows/ci-verify.yml`

## Verification

- `grep -q 'make verify-changed' .github/workflows/ci-verify.yml && grep -q 'make verify$' .github/workflows/ci-verify.yml && grep -q 'make verify$' .github/workflows/release.yml` — expected: exit 0; before this Task the pull request job does not name `make verify-changed`, so the command fails.

## References

- [_techspec.md](_techspec.md) — Where each is used

## Result

### Implementation

- Split the shared CI setup into event-specific verification steps. Pull
  requests run `make verify-changed` with `VERIFY_BASE` set to the pull
  request base branch's head SHA; pushes to `main` keep the budgeted
  `make verify` gate.
- Kept `fetch-depth: 0` so the selector can resolve the merge base, and left
  `.github/workflows/release.yml` unchanged.

### Focused checks

- `rtk git diff --check` — passed.
- A focused Ruby structural check over both workflow files exited 0 and
  reported `pr_selective: ok`, `push_complete: ok`, `full_history: ok`, and
  `release_complete: ok`.
- The changed-file postflight found `.github/workflows/ci-verify.yml` as the
  only newly changed path beyond the pre-existing `.roundfixrc.yml` and
  daemon-owned `task_04.md` status change. The workflow is inside this Spec's
  approved tooling boundary.
- The Task's declared `## Verification` command was not run; the Daemon owns
  it.

### Acceptance evidence

1. The pull-request-only `Verify changed` step names `make verify-changed` and
   passes `${{ github.event.pull_request.base.sha }}` through `VERIFY_BASE`.
   The checkout still uses `fetch-depth: 0`.
2. The push-only `Verify` step still names `make verify`, and the unchanged
   release workflow still names `make verify` in its Verify gate.

## Carry-forward provenance

- Source Run: `run_20260924T124305Z_70f339f47cee7834`
- Source commit: `9d01f5a56fec48cf4011d70733997168f2c60408`
