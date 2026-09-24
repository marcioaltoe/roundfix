---
task: task_04
spec: 0165-a-measured-run-ceiling
status: completed
type: backend
complexity: medium
---

# Task 04: A ceiling that frees itself and cannot be raced

## Overview

Corrective Task from the pre-PR review of 2026-09-24. The ceiling counts Active Runs whose owner process is dead, forever, and its refusal names no way out; Project Config can raise or disable the machine-wide ceiling; and the count is read outside the transaction that creates the Run, so two concurrent `implement`s in different repositories can both pass it.

## Requirements

1. MUST reclaim each holder whose owner process is dead, the way the per-checkout path already does, before comparing the count.
2. MUST make the refusal's next action name `roundfix stop <id>` for each holder and the `runs.max_active` key.
3. MUST ignore, with a warning, `runs.max_active` from Project Config; only User Config sets it; state this in `docs/user-guide/configuration.md`.
4. MUST count Active Implement Runs inside the same write transaction that creates the Run and refuse there with a typed ceiling error, keeping the preflight check for the early message.

## Subtasks

- [ ] Implement the requirements above.
- [ ] Add a test for each acceptance criterion.

## Acceptance Criteria

- [ ] A holder with a dead owner is reclaimed and no longer blocks.
- [ ] The refusal names `roundfix stop` and `runs.max_active`.
- [ ] A Project Config value does not raise the User Config ceiling.
- [ ] Two creations racing at the ceiling cannot both succeed.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/cli/implement.go`
- interface: `internal/store/store.go`
- interface: `internal/config/config.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestRunCeilingReclaimsADeadHolder|TestRunCeilingRefusalNamesTheWayOut|TestProjectConfigCannotRaiseTheRunCeiling|TestRunCeilingIsCountedInsideTheCreateTransaction)$" ./internal/config ./internal/store ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestRunCeilingReclaimsADeadHolder TestRunCeilingRefusalNamesTheWayOut TestProjectConfigCannotRaiseTheRunCeiling TestRunCeilingIsCountedInsideTheCreateTransaction; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || exit 1; done` — expected: exit 0; before this Task none of the named cases exists, so the command fails.

## References

- [_techspec.md](_techspec.md) — The ceiling

## Result

### Implementation

- The Implement Command now runs the existing orphan-reclamation path for every
  Active Implement Run holder before its early ceiling comparison.
- One typed Store error now reports every holder with its
  `roundfix stop <id>` command and names the User Config key
  `runs.max_active`; both preflight and transactional refusals use it.
- Project Config `runs.max_active` values are removed before overlay decoding,
  emit one stderr warning, and cannot raise or disable the User Config ceiling.
  The user guide records that scope.
- `CreateRun` now counts Active Implement Runs and refuses at the configured
  ceiling inside the same machine-wide write transaction that inserts a Run.

### Acceptance evidence

- Dead holders no longer block: `TestRunCeilingReclaimsADeadHolder` passed
  after exercising and recording reclamation for two dead-owner holders before
  a new Implement Run proceeded.
- Refusal guidance is actionable: `TestRunCeilingRefusalNamesTheWayOut` passed
  with both holders' `roundfix stop <id>` commands and `runs.max_active` in the
  CLI diagnostic.
- Project Config cannot change the ceiling: `TestProjectConfigCannotRaiseTheRunCeiling`
  passed for both a disabling value (`0`) and a raised value (`9`), preserving
  the User Config value and emitting the warning.
- Concurrent creation is bounded: `TestRunCeilingIsCountedInsideTheCreateTransaction`
  passed with one creation and one typed refusal across two Store instances;
  the same case passed under the race detector and across 20 repeated runs.

### Focused checks

- `GOCACHE=/tmp/roundfix-task04-gocache go test -count=1 ./internal/config` — passed.
- `GOCACHE=/tmp/roundfix-task04-gocache go test -count=1 ./internal/store` — passed.
- `GOCACHE=/tmp/roundfix-task04-gocache go test -count=1 -run '^(TestRunCeiling|TestImplementRunCeilingZeroDisables|TestRunImplementReclaimsDeadOwnerActiveRun|TestRunImplementPreflightRejectsActiveRunInWorkingTree)' ./internal/cli` — passed.
- `GOCACHE=/tmp/roundfix-task04-gocache go test -race -count=1 -run '^TestRunCeilingIsCountedInsideTheCreateTransaction$' ./internal/store` — passed.
- `GOCACHE=/tmp/roundfix-task04-gocache go test -count=20 -run '^TestRunCeilingIsCountedInsideTheCreateTransaction$' ./internal/store` — passed.
- The Task's declared Verification command was not run; it remains Daemon-owned.

### Context lookup

- Read `secondbrain/wiki/index.md`; the required qmd query could not initialize
  its local Metal context and returned no results, so no Secondbrain content
  informed the implementation.

### Follow-ups

- None found within this Task's slice.
