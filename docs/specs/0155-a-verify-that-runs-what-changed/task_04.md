---
task: task_04
spec: 0155-a-verify-that-runs-what-changed
status: pending
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
