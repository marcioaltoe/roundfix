---
task: task_03
spec: 0155-a-verify-that-runs-what-changed
status: completed
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

## Result

### Implementation

- Added `verify-changed` with `fmt-check`, `vet`, and `build` as whole-tree
  prerequisites. It passes `VERIFY_BASE` to `verify-select`, defaulting to
  `main`, and dispatches only the sets the selector prints.
- Added separate core and Baseline recipes. The core recipe excludes the
  Baseline CLI test pattern; the Baseline recipe runs the Baseline packages,
  the matching CLI tests, `skills-sync-check`, and `skills-check`.
- Left the `verify` recipe unchanged and pointed the repository Verification
  setting at `make verify-changed`; its comment names `make verify` as the
  complete gate for `main` and releases.

### Focused checks

- `rtk env GOCACHE=/tmp/roundfix-task03-selector-gocache go test
  ./internal/verifyselect ./cmd/verify-select` — passed; the selector package
  completed in 11.743 seconds and the entry point compiled.
- `rtk env GOCACHE=/tmp/roundfix-task03-prereq-gocache make fmt-check vet
  build` — passed.
- `rtk env GOCACHE=/tmp/roundfix-task03-core-gocache make
  verify-changed-core` — the sandboxed attempt was blocked when the existing
  suite reached `api.github.com`; the permitted rerun passed, including
  `internal/cli` in 55.694 seconds.
- `rtk make -n verify` — exited 0 and retained `go test -parallel 16 ./...`,
  followed by both skill checks and the build.
- `rtk rg -n
  "^VERIFY_BASE|^VERIFY_SELECT|^verify:|^verify-changed:|skills-sync-check skills-check|verification: make verify-changed"
  Makefile .roundfixrc.yml` — exited 0 and found the expected wiring.
- The Task's declared `## Verification` command was not run; the Daemon owns
  it.

### Acceptance evidence

1. The selector package's core-only cases passed, `verify-changed` dispatches
   one helper for each selector output line, and `verify-changed-core` completed
   without Baseline packages or skill checks. The exact top-level core-only
   journey needs a committed Makefile baseline and remains part of Task 05's
   terminal QA matrix.
2. The `verify` recipe is unchanged. Its focused dry run still expands the
   complete `./...` package test plus both skill checks and the build.
3. `.roundfixrc.yml` now contains `verification: make verify-changed`, as shown
   by the focused wiring inspection.

## Carry-forward provenance

- Source Run: `run_20260924T124305Z_70f339f47cee7834`
- Source commit: `cc0d2d4439907c2d7103f8ae42fbecacd62ebd8b`
