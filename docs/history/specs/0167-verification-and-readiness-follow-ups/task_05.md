---
task: task_05
spec: 0167-verification-and-readiness-follow-ups
status: completed
type: backend
complexity: low
---

# Task 05: Only a verdict replaces a verdict; setup writes a project template

## Overview

Corrective Task from the pre-PR review of 2026-09-25. `retainCollectedVerificationFailures` treats a command the exclusive retry reached but ended in unknown as having a retry verdict, so its first-run deterministic failure is dropped while a first-run temporary failure is carried as the repair target. And `roundfix setup` still writes the User Config template, with `runs.max_active`, as a new Project Config.

## Requirements

1. MUST let a retry verdict replace a first-run verdict only when the retry passed or failed that command; a command that ended in unknown keeps its first-run failure.
2. MUST never carry a first-run temporary failure as a repair target.
3. MUST make `roundfix setup` write the project-scope template, without `runs.max_active`.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] Unknown on retry keeps the first-run deterministic failure as the repair target.
- [ ] Setup followed by Load warns nothing.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/cli/setup.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestUnknownOnRetryKeepsTheFirstRunFailure|TestFirstRunTemporaryFailureIsNotARepairTarget|TestSetupWritesNoRunCeilingIntoProjectConfig)$" ./internal/daemon ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestUnknownOnRetryKeepsTheFirstRunFailure TestFirstRunTemporaryFailureIsNotARepairTarget TestSetupWritesNoRunCeilingIntoProjectConfig; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — Build Order

## Result

Implementation:

- Retry merging now replaces a first-run command failure only when the retry observed a pass or failure for that command. An unobserved retry keeps the first-run deterministic failure, and the first-run temporary failure is excluded from repair feedback.
- Setup now generates a new Project Config from the project-scope template, while a new User Config continues to use the user-scope template.
- Added the three named regression tests and updated the existing retry-merge expectation so temporary failures remain retry control flow rather than repair targets.

Focused checks:

- `GOCACHE=/tmp/roundfix-task05-gocache go test -count=1 -run '^(TestUnknownOnRetryKeepsTheFirstRunFailure|TestFirstRunTemporaryFailureIsNotARepairTarget)$' ./internal/daemon` — exit 0.
- `GOCACHE=/tmp/roundfix-task05-gocache go test -count=1 -run '^TestSetupWritesNoRunCeilingIntoProjectConfig$' ./internal/cli` — exit 0.
- `GOCACHE=/tmp/roundfix-task05-gocache go test -count=1 -run '^(TestRetryVerdictReplacesTheFirstRun|TestAFailureRepeatedOnRetryReachesRepairOnce|TestRetryKeepsFirstRunFailuresForCommandsItDidNotReach|TestIndependentVerificationReplacesCollectedFailuresForRetriedCommands|TestIndependentVerificationKeepsTemporaryRetryHandling)$' ./internal/daemon` — exit 0.
- `GOCACHE=/tmp/roundfix-task05-gocache go test -count=1 -run '^(TestRunSetupFreshMachineAcceptsOffers|TestRunSetupHealthyMachineIsIdempotent|TestRunSetupNewerACPXIsReadyWithoutInstall|TestSetupCommandCompatibility|TestSetupWritesNoRunCeilingIntoProjectConfig)$' ./internal/cli` — exit 0.
- `GOCACHE=/tmp/roundfix-task05-gocache go test -count=1 -run '^(TestProjectConfigCannotRaiseTheRunCeiling|TestProjectInitWritesNoRunCeiling|TestProjectInitThenLoadWarnsNothing)$' ./internal/config` — exit 0.
- `git diff --check` — exit 0.

Acceptance evidence:

- Unknown on retry keeps the first-run deterministic failure as the repair target: `TestUnknownOnRetryKeepsTheFirstRunFailure` exercises the merge directly and exits 0.
- Setup followed by Load warns nothing: `TestSetupWritesNoRunCeilingIntoProjectConfig` runs setup, loads the emitted Project Config through the real config loader, asserts no `runs.max_active`, and exits 0 with no warning.
- A first-run temporary failure is never a repair target: `TestFirstRunTemporaryFailureIsNotARepairTarget` exits 0.

The Task's declared `## Verification` command was not run; Daemon Verification owns that command and the terminal Task status.
