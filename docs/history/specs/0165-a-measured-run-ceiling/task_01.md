---
task: task_01
spec: 0165-a-measured-run-ceiling
status: completed
type: backend
complexity: medium
---

# Task 01: A machine-wide Active Run ceiling

## Overview

Nothing bounds how many Implement Runs one machine starts; five at once on 2026-09-24 overloaded the machine enough to fail an unrelated QA gate.

## Requirements

1. MUST add `runs.max_active` to User Config: integer, default 3, `0` disables, a negative value refused by validation, documented in `docs/user-guide/configuration.md`.
2. MUST count Active Runs of kind implement across every repository in the Run Database.
3. MUST make `implement` refuse with exit 2 before creating a Run when that count has reached the ceiling, listing each holding Run's id, repository and Spec.
4. MUST leave the per-checkout Active Run check unchanged.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] At `max_active: 2` with Active Runs in two repositories, a third `implement` refuses naming both and creates nothing.
- [ ] `0` disables the bound; a negative value is refused.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/config/config.go`
- interface: `internal/store/store.go`
- interface: `internal/cli/implement.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestConfigRunsMaxActiveDefaultsToThree|TestConfigRefusesANegativeRunsMaxActive|TestActiveImplementRunsAreCountedAcrossRepositories|TestImplementRefusesAtTheActiveRunCeiling|TestImplementRunCeilingZeroDisables)$" ./internal/config ./internal/store ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestConfigRunsMaxActiveDefaultsToThree TestConfigRefusesANegativeRunsMaxActive TestActiveImplementRunsAreCountedAcrossRepositories TestImplementRefusesAtTheActiveRunCeiling TestImplementRunCeilingZeroDisables; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — The ceiling

## Result

Implemented `runs.max_active` with a built-in value of `3`, strict rejection
of negative values, and `0` as the disabled value. The Run Database now exposes
every non-terminal Implement Run across repositories, and `implement` checks
that machine-wide set after the unchanged per-checkout Active Run check and
before Run creation. A refusal exits `2` and lists every holder's Run id,
repository, and Spec. The configuration guide documents the default, refusal,
and disable behavior.

Acceptance evidence:

- With two Active Implement Runs in different repositories and
  `runs.max_active: 2`, `TestImplementRefusesAtTheActiveRunCeiling` observed
  exit `2`, found both holders' ids, repositories, and Specs in stderr, and
  confirmed the Run Database still contained only the two seeded Runs.
- With the same machine-wide pressure and `runs.max_active: 0`,
  `TestImplementRunCeilingZeroDisables` allowed the third Implement Run to
  reach its normal Clean fixture outcome. `TestConfigRefusesANegativeRunsMaxActive`
  confirmed that `-1` fails validation.
- `TestActiveImplementRunsAreCountedAcrossRepositories` confirmed that the
  store returns non-terminal Implement Runs from two repositories while
  excluding a terminal Implement Run and an Active Fetch Run.
- `TestRunImplementPreflightRejectsActiveRunInWorkingTree` passed with its
  existing refusal text and no added Run, preserving the per-checkout check.

Focused checks:

- Pre-change signal: `rtk env GOCACHE=/tmp/roundfix-0165-go-cache go test -count=1 -run '^TestConfig(RunsMaxActiveDefaultsToThree|RefusesANegativeRunsMaxActive)$' ./internal/config` failed to build because `Config.Runs` did not exist.
- `rtk env GOCACHE=/tmp/roundfix-0165-go-cache go test -count=1 -run '^TestConfig(RunsMaxActiveDefaultsToThree|RefusesANegativeRunsMaxActive)$' ./internal/config` — passed.
- `rtk env GOCACHE=/tmp/roundfix-0165-go-cache go test -count=1 -run '^TestActiveImplementRunsAreCountedAcrossRepositories$' ./internal/store` — passed.
- `rtk env GOCACHE=/tmp/roundfix-0165-go-cache go test -count=1 -run '^TestImplement(RefusesAtTheActiveRunCeiling|RunCeilingZeroDisables)$' ./internal/cli` — passed.
- Adjacent config generation, store listing, and the existing per-checkout
  refusal checks passed in focused package runs.
- `rtk env GOCACHE=/tmp/roundfix-0165-go-cache make verify-incremental` — the
  sandboxed attempt reached the existing force-stop integration tests but
  could not enumerate the process table (`operation not permitted`); the
  rerun with process-table access passed vet, the full Go suite, skill checks,
  and the build.

The Daemon-owned command under `## Verification` was not run in this Agent
turn.
