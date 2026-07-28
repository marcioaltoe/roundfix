---
task: task_01
spec: 0042-verification-capacity-and-daemon-task-settlement
status: completed
type: backend
complexity: high
---

# Task 01: Configure independent Verification Capacity

## Overview

Give Spec Runs an effective Verification Capacity that is independent from
Task Capacity. The complete slice starts at strict User/Project Config,
travels through the Implement Command into the Task cycle, and is durably
published when the cycle starts without changing the command's public flags,
output ownership, or exit codes.

## Requirements

1. MUST add top-level `verification.concurrency` with built-in default `1`,
   positive-integer validation, and existing per-key configuration precedence.
2. MUST keep `defaults.verification` as the unchanged review Verification
   command and reject unknown fields in the new section.
3. MUST include the safe default and explanatory distinction from
   `worktree.concurrency` in generated User and Project Config.
4. MUST pass the effective value through the Implement Command into the Task
   cycle without adding a CLI flag or late validation.
5. MUST publish Task Capacity and Verification Capacity in the
   Task-cycle-start event while retaining the existing `concurrency` payload
   value as the Task Capacity compatibility alias.
6. MUST reject missing or non-positive TaskPlan capacity before starting an
   Agent, Verification command, Run mutation beyond existing preflight, or
   worker goroutine.
7. MUST preserve review Run configuration, stdout/stderr discipline, public
   command help, and exit-code behavior.

## Subtasks

- [x] Add the independent configuration model, overlay, default, and validation.
- [x] Extend generated config and focused configuration tests.
- [x] Thread effective Verification Capacity through Implement planning.
- [x] Publish both capacities from Task-cycle start with the compatibility alias.
- [x] Add negative no-side-effect tests for invalid capacity.
- [x] Protect review Verification and public CLI contracts with regression tests.

## Acceptance Criteria

- [x] Omitted `verification.concurrency` resolves to `1` regardless of the
      effective `worktree.concurrency` value.
- [x] User Config is overridden by Project Config for the new key, while an
      omitted Project value inherits the User value.
- [x] Values `0` and `-1`, non-integer values, and unknown Verification keys
      fail strict validation naming `verification.concurrency` before a Run is
      created.
- [x] The Task-cycle-start event carries equal legacy `concurrency` and
      `task_capacity` values plus the independently resolved
      `verification_capacity`.
- [x] No new Implement flag appears and existing review Verification tests are
      byte-stable.
- [x] Focused config, CLI, daemon, race, and build gates pass with no new
      dependency.

## Context

- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `.agents/skills/golang-cli/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `internal/config/config.go`
- interface: `internal/config/config_test.go`
- interface: `internal/cli/implement.go`
- interface: `internal/cli/implement_test.go`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/daemon/task_engine_test.go`
- interface: `internal/runevent/event.go`

## Verification

- `rtk go test ./internal/config -run 'Test(Load.*Verification|DefaultConfigYAML.*Verification|Validate.*Verification)' -count=1` — expected: default, precedence, generated YAML, strict decoding, and invalid-value cases pass.
- `rtk go test ./internal/daemon ./internal/cli -run 'Test(TaskCycle.*Capacit|ExecuteImplementCycle.*VerificationCapacit|RunImplement.*VerificationCapacit)' -count=1` — expected: effective capacities reach the Task cycle and its start event without a new CLI flag.
- `rtk go test -race ./internal/config ./internal/daemon ./internal/cli -run 'Test(Load.*Verification|TaskCycle.*Capacit|ExecuteImplementCycle.*VerificationCapacit)' -count=1` — expected: the configuration and planning path is race-free.
- `rtk go build -buildvcs=false ./...` — expected: all packages compile with the extended internal configuration and TaskPlan contracts.

## References

- `_prd.md` → Goals 1 and 4; User Stories 1–3; Core Features 1 and 8; User Experience; Success Metrics.
- `_techspec.md` → System Architecture; Implementation Design: Interfaces, Data Models, and API Contracts; Build Order 1.
- `../../adr/0056-spec-runs-separate-task-and-verification-capacity.md` → independent capacity decision and compatibility boundary.

## Result

Implemented independent Verification Capacity from strict User/Project Config
through the Implement Command and `TaskPlan`. Task-cycle start now publishes
`concurrency` and `task_capacity` as equal Task Capacity values plus the
independently resolved `verification_capacity`. Invalid TaskPlan capacities
return before Run Events, Run state mutation, Agent or Verification calls,
progress output, or scheduler work.

Verification used a Task-local `GOCACHE` under `/private/tmp` because the
sandbox cannot write the default macOS Go build cache:

- `rtk env GOCACHE=/private/tmp/roundfix-task01-go-cache go test ./internal/config -run 'Test(Load.*Verification|DefaultConfigYAML.*Verification|Validate.*Verification)' -count=1` — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task01-go-cache go test ./internal/daemon ./internal/cli -run 'Test(TaskCycle.*Capacit|ExecuteImplementCycle.*VerificationCapacit|RunImplement.*VerificationCapacit)' -count=1` — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task01-go-cache go test -race ./internal/config ./internal/daemon ./internal/cli -run 'Test(Load.*Verification|TaskCycle.*Capacit|ExecuteImplementCycle.*VerificationCapacit)' -count=1` — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task01-go-cache go build -buildvcs=false ./...` — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task01-go-cache go test ./internal/config ./internal/daemon ./internal/cli -count=1` — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task01-go-cache make verify` — passed.
- `rtk git -c core.fsmonitor=false diff --check` — passed.

Acceptance evidence:

1. `TestLoadVerificationCapacityHierarchy` proves the built-in value stays `1`
   when Task Capacity differs, User Config inheritance, and Project Config
   override.
2. `TestValidateVerificationCapacity` and
   `TestRunImplementRejectsInvalidVerificationCapacityBeforeRunCreation` prove
   zero, negative, non-integer, and unknown-key failures name the new key and
   create no Run Database.
3. `TestTaskCycleRejectsInvalidCapacitiesBeforeSideEffects` proves missing and
   non-positive TaskPlan capacities cause no Daemon side effect.
4. `TestTaskCyclePublishesCapacities` and
   `TestRunImplementPassesVerificationCapacityIntoTaskCycle` prove the
   compatibility alias and both explicit capacity fields in the durable start
   event.
5. `TestDefaultConfigYAMLVerificationCapacity` proves generated User and
   Project Config include the independent safe default while retaining the
   exact review `defaults.verification: make verify` contract.
6. `TestRunImplementVerificationCapacityDoesNotAddFlag`, the race check, full
   build, and `make verify` prove the public CLI and repository gates remain
   intact. `go.mod` and `go.sum` are unchanged.

Follow-ups: none.
