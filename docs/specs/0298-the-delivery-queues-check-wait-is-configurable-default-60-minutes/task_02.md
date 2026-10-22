---
task: task_02
spec: 0298-the-delivery-queues-check-wait-is-configurable-default-60-minutes
status: completed
type: backend
complexity: medium
---

# Task 02: The Check Wait reads delivery.check_timeout, 60 minutes by default

## Overview

The Delivery Queue waits at most 30 minutes, a constant the `deliver` command
cannot change, for a Pull Request's current-head checks, and parks the item
with a bare `checks-timeout`. This Task adds the User Config key
`delivery.check_timeout` (a Go duration, default `60m`, greater than 0 and at
most `24h`, ignored with a warning in Project Config), makes the engine's
default the same constant, passes the loaded value to the engine the command
builds, and makes the park name the key and the value it waited, while each
wait keeps starting its own deadline.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-22,
   [`references/2026-10-22-the-delivery-queue-parks-checks-timeout-while-ci-is-still-running.md`](references/2026-10-22-the-delivery-queue-parks-checks-timeout-while-ci-is-still-running.md)
   (`_prd.md` → Sources 1), by implementing API Contracts 1-5 of
   `_techspec.md` in `internal/config/config.go`,
   `internal/delivery/engine.go` and `internal/cli/deliver_workflow.go`, with
   the names `_techspec.md` → Interfaces gives.
2. MUST add to `internal/config` the exported constants
   `DeliveryCheckTimeoutKey` (`"delivery.check_timeout"`),
   `DefaultDeliveryCheckTimeout` (60 minutes) and `MaxDeliveryCheckTimeout`
   (24 hours), the field `Delivery.CheckTimeout`, the overlay field
   `check_timeout` decoded with the existing `durationValue`, the built-in
   value `DefaultDeliveryCheckTimeout`, and the validation that refuses a
   value `<= 0` or above `MaxDeliveryCheckTimeout` with exactly
   `delivery.check_timeout must be greater than 0 and at most 24h`. A Project
   Config source has the key removed before decoding with
   `warnIgnoredProjectSetting(DeliveryCheckTimeoutKey)`, and the overlay
   applies the value only for a source other than Project Config. The init
   templates stay unchanged.
3. MUST make `defaultCheckTimeout` in `internal/delivery/engine.go` read
   `config.DefaultDeliveryCheckTimeout`, add `checksTimeoutBlocker(timeout
   time.Duration) string` returning `checks-timeout: current-head checks not
   settled within delivery.check_timeout <timeout.String()>` built from
   `BlockerChecksTimeout` and `config.DeliveryCheckTimeoutKey`, and make
   `Engine.checkCandidate` park with it at the deadline. The three deadline
   call sites in `Engine.createPullRequest` and `Engine.checkCandidate` keep
   computing their deadline from their own start, and the single restart
   after a re-run (`timeoutRestarted`) is unchanged. `BlockerChecksTimeout`,
   `ClassifyPark` and its `Next` text keep their bytes.
4. MUST extract `commandDeliveryEngineDependencies(runStore *store.Store,
   loaded roundconfig.Loaded) delivery.EngineDependencies` from
   `newCommandDeliveryEngine` in `internal/cli/deliver_workflow.go`, building
   the workflow and the GitHub CLI exactly as today and adding
   `CheckTimeout: loaded.Config.Delivery.CheckTimeout`; `newCommandDeliveryEngine`
   returns `delivery.NewEngine(runStore, commandDeliveryEngineDependencies(runStore, loaded))`.
5. MUST state the hazard in the code's behavior. The value comes from the
   operator's User Config only, never from a repository: a Project Config
   value, including one on the fetched default branch, is removed and warned
   about. The wait abandons nothing and cancels nothing: at its deadline the
   item parks and its Pull Request, CI runs and branch are left as they are.
   The destination of the park text is the item's blocker in the Run
   Database, which `roundfix deliver status` prints; it carries the key and a
   duration and never a secret.
6. MUST declare this break, the only change to existing tests: the timeout
   assertions of `TestPersistentCheckReadErrorsParkAtTheDeadline` and
   `TestChecksTimeoutParksTheItem` (`internal/delivery/engine_test.go`),
   `TestABlockedMergeStateParksAtTheChecksTimeout`
   (`internal/delivery/merge_state_test.go`),
   `TestAnOldCheckAttemptCannotKeepResettingTheTimeout`
   (`internal/delivery/check_rerun_test.go`) and
   `TestStaleRerunsAreBoundedPerTipAndRestartTheTimeoutOnlyOnce`
   (`internal/delivery/stale_check_test.go`) compare the blocker with the
   literal new park text for their configured timeout (`... 2s`, `... 3s`,
   `... 5s`) instead of `BlockerChecksTimeout`; their inputs, polls and other
   assertions do not change. No other existing test, golden or fixture
   changes.
7. MUST add the tests `_techspec.md` → Testing Approach lists, each top-level
   test calling `t.Parallel()` first, compiling for Windows, using temporary
   homes and the engine's fake clock and sleeper, and waiting on no real time:
   - in `internal/config/delivery_check_timeout_test.go`,
     `TestDeliveryCheckTimeoutIsReadFromUserConfig` (absent → 60 minutes and
     `DefaultDeliveryCheckTimeout` equals 60 minutes; `45m`; `24h`),
     `TestProjectConfigCannotSetTheDeliveryCheckTimeout` (User Config `45m`,
     Project Config `0`, `5m`, `-1h`, `invalid`; the value stays 45 minutes
     and standard error equals the warning line) and
     `TestDeliveryCheckTimeoutRefusesAValueOutsideItsBounds` (`0`, `0s`,
     `-1m`, `24h1s` give the range message; `soon` gives `invalid
     duration`);
   - in `internal/delivery/check_wait_test.go`,
     `TestTheCheckWaitDefaultsToSixtyMinutes` (no `CheckTimeout`, checks
     pending forever: parked after exactly 60 minutes and 241 polls with the
     blocker `checks-timeout: current-head checks not settled within
     delivery.check_timeout 1h0m0s`, Park Class `environment`),
     `TestTheCheckWaitUsesTheConfiguredDuration` (45 minutes: parked after
     exactly 45 minutes with `... 45m0s`),
     `TestEachCheckWaitStartsItsOwnDeadline` (45 minutes: a first call that
     sees `pass` at the 161st poll moves to `merging` after 40 minutes; the
     item set back to `checking` and pending forever parks exactly 45 minutes
     after the second call starts) and
     `TestThePullRequestHeadWaitUsesTheCheckWait` (an existing Pull Request
     that keeps reporting the earlier candidate commit makes
     `createPullRequest` fail with `PR Head Branch is at "earlier", expected
     "candidate"` after exactly 45 minutes);
   - in `internal/cli/deliver_check_timeout_test.go`,
     `TestTheDeliveryEngineReceivesTheConfiguredCheckWait`
     (`roundconfig.Builtin()` gives `DefaultDeliveryCheckTimeout`; a loaded
     config with `Delivery.CheckTimeout` 45 minutes and `Runs.SlotWait` 90
     minutes gives both).

## Subtasks

- [ ] Add the key, its constants, its validation and its Project Config removal.
- [ ] Make the engine default the config constant and the park text name the key and value.
- [ ] Extract the dependency builder and pass the loaded value.
- [ ] Add the configuration, engine and wiring tests and move the five timeout assertions to the new text.

## Acceptance Criteria

- [ ] With nothing set, a pending Check Wait parks after exactly 60 minutes with the new park text in the `environment` Park Class.
- [ ] A User Config value reaches the engine the `deliver` command builds and bounds both the check wait and the Pull Request head wait.
- [ ] A Project Config value is ignored with the warning, and `0`, negative and over-`24h` values are refused.
- [ ] A second wait starts its own deadline, and the re-run tests still restart it once.

## Consumers

- `config.Delivery.CheckTimeout`, `config.DeliveryCheckTimeoutKey`, `config.DefaultDeliveryCheckTimeout` and `config.MaxDeliveryCheckTimeout` — consumers: `./internal/config`, `./internal/delivery`, `./internal/cli`
- The `checks-timeout` park text and the engine default — consumers: `./internal/delivery`, `./internal/cli`
- `commandDeliveryEngineDependencies` — consumers: `./internal/cli`

## Context

- instruction: `docs/adr/0311-the-delivery-queues-check-wait-is-a-user-config-duration-sixty-minutes-by-default.md`
- instruction: `docs/adr/0304-a-queue-item-waits-for-a-run-slot-up-to-runs-slot-wait.md`
- instruction: `docs/adr/0236-a-failed-check-is-judged-on-a-merge-with-the-current-default-branch.md`
- instruction: `docs/user-guide/configuration.md`
- instruction: `docs/user-guide/commands/deliver.md`
- instruction: `internal/delivery/park_class.go`
- instruction: `internal/delivery/github_test.go`
- instruction: `internal/delivery/conflict_test.go`
- instruction: `internal/config/run_slot_wait_test.go`
- instruction: `internal/config/profile_proof_setup_timeout_test.go`
- instruction: `internal/config/config_test.go`
- instruction: `internal/cli/deliver.go`
- instruction: `internal/cli/deliver_park_status_test.go`
- interface: `internal/config/config.go`
- interface: `internal/delivery/engine.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/delivery/engine_test.go`
- interface: `internal/delivery/merge_state_test.go`
- interface: `internal/delivery/check_rerun_test.go`
- interface: `internal/delivery/stale_check_test.go`
- creates: `internal/config/delivery_check_timeout_test.go`
- creates: `internal/delivery/check_wait_test.go`
- creates: `internal/cli/deliver_check_timeout_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestDeliveryCheckTimeoutIsReadFromUserConfig|TestProjectConfigCannotSetTheDeliveryCheckTimeout|TestDeliveryCheckTimeoutRefusesAValueOutsideItsBounds|TestRunSlotWaitIsReadFromUserConfig|TestProjectConfigCannotSetTheRunSlotWait)$' ./internal/config 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDeliveryCheckTimeoutIsReadFromUserConfig TestProjectConfigCannotSetTheDeliveryCheckTimeout TestDeliveryCheckTimeoutRefusesAValueOutsideItsBounds TestRunSlotWaitIsReadFromUserConfig TestProjectConfigCannotSetTheRunSlotWait; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v ./internal/delivery 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheCheckWaitDefaultsToSixtyMinutes TestTheCheckWaitUsesTheConfiguredDuration TestEachCheckWaitStartsItsOwnDeadline TestThePullRequestHeadWaitUsesTheCheckWait TestPersistentCheckReadErrorsParkAtTheDeadline TestChecksTimeoutParksTheItem TestABlockedMergeStateParksAtTheChecksTimeout TestAnOldCheckAttemptCannotKeepResettingTheTimeout TestStaleRerunsAreBoundedPerTipAndRestartTheTimeoutOnlyOnce TestCheckRerunRestartsTheTimeout TestAConflictingPullRequestStopsTheCheckWaitAtOnce; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v -run '^(TestTheDeliveryEngineReceivesTheConfiguredCheckWait|TestDeliverStatusPrintsAParkLinePerParkedItem)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheDeliveryEngineReceivesTheConfiguredCheckWait TestDeliverStatusPrintsAParkLinePerParkedItem; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go vet ./internal/config ./internal/delivery ./internal/cli` — expected: exit 0. Before this Task the seven new tests do not exist, so the command fails at its first pass check. Afterwards the key loads, refuses and ignores values as specified, the whole `internal/delivery` package passes with the default, configured, per-wait and Pull Request head waits and the five moved timeout assertions, the `deliver` command's engine receives the value, the unchanged re-run, conflict, slot-wait and status tests pass, and the three packages' tests compile for Windows.

## References

- `_prd.md` → Sources; Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5
- `_techspec.md` → Current behavior; API Contracts 1-5; Invariants 1-6; Testing Approach; Build Order 2
- ADR-0236
- ADR-0304
- ADR-0311


## Result

Implemented this Task's configuration, engine and CLI slice; Daemon Verification
and settlement remain pending. The initial worktree change was this Task's
Daemon-owned `in_progress` status, which this turn preserved.

- Added the exported key/default/maximum constants, `Delivery.CheckTimeout`,
  duration overlay, 60-minute built-in value and exact range validation.
  Project Config removes the key before decoding and warns; the overlay also
  refuses to apply it from that source. `DeliveryConfigAtCommit` uses
  `ResolveConfigProposal`, so fetched default-branch Project Config reaches
  the same removal. Init templates are unchanged.
- The engine default uses the config constant. Deadline parking stores
  `checks-timeout: current-head checks not settled within delivery.check_timeout <duration>`.
  Only the key and duration reach the Run Database blocker and status output;
  no secret is included. Parking leaves the Pull Request, CI runs and branch
  in place and cancels or abandons none of them.
- Extracted `commandDeliveryEngineDependencies` without changing its workflow
  or GitHub CLI construction; it passes the loaded check timeout beside the
  existing Run Slot Wait.
- Declared existing-test break applied only to the five named timeout
  assertions: the two engine tests now expect the literal `2s` blocker, the
  blocked merge-state test `3s`, and the old-attempt and stale-rerun tests `5s`.
  Their inputs, polling and other assertions are unchanged. No other existing
  tests, goldens or fixtures changed.

Acceptance evidence from focused checks:

| Criterion | Evidence |
| --- | --- |
| Default wait and environment park | `TestTheCheckWaitDefaultsToSixtyMinutes` passed: exactly 60 fake minutes, 241 polls, the literal `1h0m0s` blocker persisted, and environment Park Class. |
| User value reaches both waits | `TestDeliveryCheckTimeoutIsReadFromUserConfig`, `TestTheDeliveryEngineReceivesTheConfiguredCheckWait`, `TestTheCheckWaitUsesTheConfiguredDuration` and `TestThePullRequestHeadWaitUsesTheCheckWait` passed: absent/45m/24h loading, default/45m dependency wiring with independent 90m slot wait, and both waits bounded at exactly 45 fake minutes. |
| Project source and bounds | `TestProjectConfigCannotSetTheDeliveryCheckTimeout` passed for 0, 5m, -1h and invalid with exact warning and User Config 45m preserved. `TestDeliveryCheckTimeoutRefusesAValueOutsideItsBounds` passed for 0, 0s, -1m and 24h1s with exact range error and soon with invalid-duration decoding. |
| Independent deadlines and bounded restart | `TestEachCheckWaitStartsItsOwnDeadline` passed: first wait merges at poll 161 after 40 fake minutes, second parks 45 minutes after its own start. Existing `TestCheckRerunRestartsTheTimeout`, `TestAnOldCheckAttemptCannotKeepResettingTheTimeout` and `TestStaleRerunsAreBoundedPerTipAndRestartTheTimeoutOnlyOnce` passed. All three deadline assignments and `timeoutRestarted` logic are unchanged. |

Commands and outcomes:

- Red starting signal: `rtk proxy go test ./internal/config -run '^TestDeliveryCheckTimeoutIsReadFromUserConfig$' -count=1`
  failed to compile because the new constants and field were absent.
- An initial focused run reported Go vet's suspect-or diagnostic for a
  redundant comparison against equal default constants; removed the redundant
  comparison, retaining the separate config constant-equality assertion.
- Final focused regression command:
  `rtk proxy go test -count=1 -run 'CheckWait|DeliveryCheckTimeout|PersistentCheckReadErrorsParkAtTheDeadline|ChecksTimeoutParksTheItem|ABlockedMergeStateParksAtTheChecksTimeout|AnOldCheckAttemptCannotKeepResettingTheTimeout|StaleRerunsAreBoundedPerTipAndRestartTheTimeoutOnlyOnce|CheckRerunRestartsTheTimeout|AConflictingPullRequestStopsTheCheckWaitAtOnce' ./internal/config ./internal/delivery ./internal/cli`
  exited 0 for all three packages. This includes all eight new top-level
  tests, the five moved assertions and the unchanged re-run/conflict checks.
- Windows test compilation: created a disposable output directory with
  `task_windows_dir=$(mktemp -d /tmp/roundfix-check-wait-windows.XXXXXX)`, then
  `GOOS=windows rtk proxy go test -c -o "$task_windows_dir/" ./internal/config ./internal/delivery ./internal/cli`
  exited 0. No Windows binary was executed.
- `rtk proxy git diff --check` exited 0. Diff inspection confirmed only this
  Task's declared implementation/test paths and this Task file changed.

Environment: normal execution refused the symlinked writable root; commands
ran with authorized escalation in the resolved Run checkout. No host config,
forge, real-time wait, commit, push or Pull Request was used. No follow-up
outside this slice was found. The declared Verification, full delivery suite,
Windows vet, status regression and repository gates were not run by this
Agent; the Daemon owns them and Task settlement.


### Verification Feedback — attempt 1

Inspected the Daemon diagnostic artifact
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261022T163107Z_e127f4e4a764d3ab/verification/batch-002-attempt-1.log`.
The reported failure was the Behavior Surface Record's stale config-key
fingerprint: this Task adds the `check_timeout` YAML tag after task_01's
record was generated. The owning contract computes its config-key surface
from production YAML struct tags, so this is derived fallout of this Task's
new key.

Ran the exact generator already sanctioned by this Spec's `_authorization.md`:
`rtk proxy go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`
— exit 0. This regenerated `docs/references/behavior-surfaces.json`; reviewed
the diff and confirmed that only its `config keys` fingerprint changed. No
fingerprint was edited by hand, no contract was weakened, and no other Task
file or source was changed in this feedback turn.

Focused confirmation without record mode:
`rtk proxy go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -count=1`
— exit 0. `rtk proxy git diff --check` also exited 0.

This adds the generated Behavior Surface Record to this Task's implementation
fallout. The earlier implementation evidence remains applicable because no
Go source or test changed in this feedback turn. Task status remains
Daemon-owned and unchanged; the full configured Verification sequence and
settlement remain pending with the Daemon. No declared Verification command,
`make verify-changed`, commit, push or Pull Request was run by this Agent.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `docs/references/behavior-surfaces.json`
