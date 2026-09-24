---
task: task_03
spec: 0155-a-verify-that-runs-what-changed
status: pending
type: infra
complexity: low
---

# Task 03: The selective target and the Daemon's gate

## Overview

Wire the selector into `make verify-changed` and make it the repository
Verification the Daemon runs. `.roundfixrc.yml` sets `defaults.verification`,
which `internal/cli/implement.go` passes to the Daemon as the command every Run's
QA gate executes.

## Requirements

1. MUST add `make verify-changed`, running `fmt-check`, `vet` and `build` over
   the whole tree and then the test sets the selector chooses.
2. MUST run the Baseline set's skill checks only when the Baseline set is
   selected.
3. MUST leave `make verify` running the complete suite, unchanged.
4. MUST set `defaults.verification` in `.roundfixrc.yml` to
   `make verify-changed`, and update its comment to name `make verify` as the
   complete gate.
5. MUST honor a `VERIFY_BASE` variable, defaulting to the repository's main
   branch.

## Subtasks

- [ ] Add the Makefile target.
- [ ] Point the Daemon's repository Verification at it.

## Acceptance Criteria

- [ ] `make verify-changed` runs to completion on a core-only change without
      running a Baseline-set test.
- [ ] `make verify` still runs every package.
- [ ] `.roundfixrc.yml` names `make verify-changed`.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `Makefile`
- interface: `.roundfixrc.yml`

## Verification

- `grep -q '^verify-changed:' Makefile && grep -q 'verify-select' Makefile && grep -q 'verification: make verify-changed' .roundfixrc.yml` — expected: exit 0; before this Task none of the three holds, so the command fails.

## References

- [_techspec.md](_techspec.md) — The targets
