---
task: task_01
spec: 0420-roundfix-stops-verifying-on-linux
status: completed
type: backend
complexity: high
---

# Task 01: The Delivery Queue, its commands and doctor lose the Linux Lane, and its two keys become deprecated

## Overview

This Task delivers `_techspec.md` → Build Order 1, on the maintainer's decision of 2026-11-24, "Tirar de vez". The Delivery Queue's repository gate ends when the repository Verification passes; `roundfix deliver lane` becomes an unknown `deliver` command; doctor loses its `linux-lane` check; `delivery.linux_lane` and `delivery.linux_lane_image` join the registry of deprecated keys; what the lane stored on an item stays readable; this repository's Project Config loses the key. The lane's two source files and its 36 tests leave, the coverage record, the Behavior Surface Record and the Skill Coverage Map are written again, and the guide's section on the removed command leaves with it. It is verifiable on its own through ten new tests over the production gate, the `deliver`, `runs` and `doctor` commands, the configuration loader and the Delivery Engine, with fake runners, a recording stand-in for `docker`, temporary repositories, temporary homes and temporary Run Databases.

## Requirements

1. MUST make `(*commandDeliveryWorkflow).Gate` return a passed result as soon as the repository Verification passes, as `_techspec.md` → API Contract 1 states: it reads no Project Config at the default branch head for the lane, resolves and starts no `docker`, writes no `linux-lane.log` and prints no line of the lane. The Shared Build Cache trim, the unit of the Machine Verification Capacity under the step name `repository gate` and the Verification Core Budget stay as they are; a failed Verification parks the item as today.
2. MUST delete `internal/cli/deliver_linux_lane.go` and `internal/cli/deliver_lane.go`, remove the `lane` case from the dispatch of `runDeliverCommand`, the three lane lines from the usage of `internal/cli/deliver.go` and the word `lane` from the root help of `internal/cli/cli.go`, so that `roundfix deliver lane` reproduces Surface Transcript 1, `roundfix deliver --help` reproduces Surface Transcript 2 and the root help holds the line `  roundfix deliver <plan|start|status|resume|retry|stop> [<slug> ...]` (API Contract 2).
3. MUST remove from `internal/delivery/engine.go` the `LaneSkip` field of `GateResult`, the warning the gate step wrote, the paragraph the publication added to the Pull Request body, `LaneSkipBody`, `laneSkipWarning`, `BlockerLinuxLaneContainer` and `WarningLinuxLaneSkipped`, leaving `GateResult` as `_techspec.md` → Interfaces gives it. When the engine rewrites an item's low-disk warning, every other segment keeps its text and its order and the low-disk segment comes last (API Contract 5).
4. MUST keep in `internal/delivery/park_class.go` one case for the constant `retiredLaneContainerBlocker` of `_techspec.md` → Interfaces: a blocker that begins `linux-lane-container:` is classed `environment` with the next step `run roundfix deliver retry <slug>; the Linux Lane was retired, so the gate no longer runs it`; a blocker `gate-failed: Linux Lane change failure: …` keeps the `finding` class and the next step of every `gate-failed` item. This is the Delegated Decision of 2026-11-24 of the `techspec` phase, "Which Park Class and next step does status show for a stored `linux-lane-container` blocker once no code produces it?".
5. MUST remove from `internal/config/config.go` the fields `LinuxLane` and `LinuxLaneImage`, their overlay fields, their validation, their merge and the image reference pattern, and append to `deprecatedConfigKeys` the two entries of `_techspec.md` → Interfaces, in that order, so that API Contract 3 and Surface Transcript 3 hold: a Project Config or User Config that names either key loads whatever the value, prints `config: <key> is deprecated and ignored; the Linux Lane was retired, so remove it` once per key and load, `delivery.linux_lane` first; a configuration read at a commit removes them and prints nothing; any other unknown key under `delivery` still fails with today's message.
6. MUST remove from `internal/cli/doctor.go` the `linux-lane` check, its dependency field and its four codes `DR-LANE-IMAGE-UNKNOWN`, `DR-LANE-DOCKER-MISSING`, `DR-LANE-DOCKER-UNREACHABLE` and `DR-LANE-IMAGE-MISSING`, and from `internal/cli/health.go` the name `HealthCheckLinuxLane`, so that the `verification-cores` line is followed by the `codex` line and every other check keeps its name, place and text (API Contract 4).
7. MUST remove from `.roundfixrc.yml` the line `linux_lane: false` and the six comment lines above it, and nothing else, so that the `delivery` mapping begins with the comment of `require_defect_class` (API Contract 10, `_authorization.md` → What the grant covers).
8. MUST delete these eight test files, whose 36 tests are the characterization of the lane on the unchanged tree (`references/2026-11-24-techspec-measurements.md`, section 6): `internal/cli/deliver_linux_lane_test.go`, `internal/cli/deliver_lane_test.go`, `internal/cli/deliver_linux_lane_go_list_test.go`, `internal/cli/deliver_lane_skip_visible_test.go`, `internal/cli/doctor_linux_lane_test.go`, `internal/config/linux_lane_test.go`, `internal/config/linux_lane_image_test.go` and `internal/delivery/lane_skip_warning_test.go`; and remove from `internal/cli/deliver_core_budget_test.go` only `TestDeliverLanePassesItsVerificationCoreBudget`, keeping `TestTheRepositoryGateCommandReceivesTheCoreBudget`. The 36 names are, exactly:
   - `TestLinuxLaneRunsTheChangedPackagesInAContainer`, `TestLinuxLaneSkipsWithANamedReason`, `TestLinuxLaneFailureFailsTheGate`, `TestLinuxLaneCancellationRemovesTheContainer`, `TestGateWithoutTheLaneStartsNoDocker`, `TestGateSkipsTheLaneWhenTheFetchedDefaultCannotBeRead`, `TestLinuxLaneOptInPropagatesCancellationAndAnUnreadableConfig`, `TestLinuxLaneRunsTheConfiguredImage`, `TestLinuxLaneContainerFailureIsToldApartFromAChangeFailure`, `TestLinuxLaneOutcomeClassification`, `TestLinuxLaneMountsTheRepositoryHistoryReadOnly`, `TestLinuxLaneSkipsAnObjectStoreItCannotMount`, `TestLinuxLaneReportsThePackagesTheContainerSelected`;
   - `TestDeliverLanePassedTranscript`, `TestDeliverLaneChangeFailureTranscript`, `TestDeliverLaneContainerFailureTranscript`, `TestDeliverLaneSkipTranscript`, `TestDeliverLaneSelectsPackagesAndImage`, `TestDeliverLaneCancellationRemovesContainer`, `TestDeliverLaneObjectStoreSkipTranscript`, `TestDeliverLanePackageCountTranscript`, `TestDeliverLaneNoPackageCountTranscript`, `TestDeliverLaneGoListFailureTranscript`, `TestDeliverLaneGoListSignalTranscript`;
   - `TestLinuxLaneScriptNamesAFailedGoList`, `TestLinuxLaneGoListFailureIsNamed`, `TestGateReturnsTheLaneEnvironmentSkip`, `TestDeliverStatusShowsTheLaneEnvironmentSkipTranscript`, `TestDoctorLinuxLaneCheckReportsDockerAndTheImage`, `TestDoctorPrintsTheLinuxLaneLineOnlyWhenTheLaneIsOn`;
   - `TestDeliveryLinuxLaneDefaultsOffAndProjectConfigEnablesIt`, `TestDeliveryLinuxLaneImageDefaultsEmptyAndProjectConfigSetsIt`, `TestDeliveryLinuxLaneImageRefusesAnInvalidReference`, `TestAGateLaneSkipIsRecordedOnTheItemAndInThePullRequestBody`, `TestAQuietGateLeavesNoLaneWarningAndReplacesAnEarlierOne`, `TestDeliverLanePassesItsVerificationCoreBudget`.
9. MUST update the tests of another behavior that name the lane, each keeping its name and losing only the lane from a fixture, a stub or a pinned line, as `references/2026-11-24-techspec-measurements.md`, section 7, lists them: `TestTopLevelUsageNamesDeliverRetry` pins the root help's new line; `TestTheQueueOwnerLogsAMachineCapacityWaitAndPassesItsGate`, `TestTheRepositoryGateHoldsTheMachineCapacityWhileItRuns`, `TestDeliveryGateTrimsTheSharedBuildCacheBeforeItsVerification` and the helper `sharedCacheImplementWorkspace` stop writing `linux_lane: false` and reading the lane's field; `TestDoctorWarnsWhenALocationIsLowOnDisk`, `TestDoctorReportsTheVerificationCoreBudget` and `TestDoctorReportsTheMachineVerificationCapacity` lose the `linux-lane:` line from their pinned output; the external skill tests stop clearing the doctor's lane dependency; the comment of `internal/cli/releaseplan_release_preparation_test.go` stops naming the lane; `TestGateParksWithTheBlockerItNames` loses its four lane rows; `TestTheQueueWarnsOfLowDiskBeforeTheRunAndTheGate` asserts that a stored lane segment keeps its text before the low-disk segment; `TestParkClassesAndNextActions` asserts the retired blocker's class and next step of requirement 4. No test of another behavior is removed to make the change pass.
10. MUST add these tests, each top-level test calling `t.Parallel()` first, with a fake runner, a recording stand-in for `docker` on the PATH the command environment carries where the test names it, temporary repositories, temporary homes and temporary Run Databases (`_techspec.md` → Testing Approach):
    - `internal/cli/deliver_gate_without_lane_test.go`: `TestTheRepositoryGateEndsAtTheRepositoryVerification`, through the fixture the surviving gate tests use, over a repository whose fetched default branch says `linux_lane: true` and names an image: a passed result, the stand-in never started, no `linux-lane.log`, no line of the lane on the console.
    - `internal/cli/deliver_lane_removed_test.go`: `TestDeliverLaneReproducesSurfaceTranscriptUnknownCommand`, asserting Surface Transcript 1 whole (every line of standard output and standard error, in order, and the exit code) in its subtests `no_arguments` (`deliver lane`) and `with_arguments` (`deliver lane --all HEAD`), and in its subtest `help_flag` that `deliver lane --help` prints the standard output of Surface Transcript 2 and exits `0`; and `TestDeliverHelpReproducesSurfaceTranscriptWithoutTheLane`, asserting Surface Transcript 2 whole.
    - `internal/cli/retired_lane_keys_test.go`: `TestRunsListReproducesSurfaceTranscriptRetiredLaneKeys`, asserting Surface Transcript 3 whole, in a new Git repository whose `.roundfixrc.yml` holds a `delivery` mapping with `linux_lane: true` and `linux_lane_image: golang:1.26`, with a home directory that holds no Roundfix file.
    - `internal/cli/doctor_without_lane_test.go`: `TestDoctorRunsNoLinuxLaneCheck`, with `linux_lane: true` in the checkout's configuration and the recording stand-in: no line begins `linux-lane:`, no `DR-LANE-` code is printed, the stand-in was never started, and the `verification-cores` line is followed by the `codex` line.
    - `internal/cli/deliver_stored_lane_records_test.go`, over a temporary Run Database with one item parked with a `linux-lane-container:` blocker and one merged item whose warning holds a `linux-lane-skipped:` segment: `TestStatusPrintsWhatTheLaneRecorded` (status prints the stored warning as stored, the `Park:` line with the class and next step of requirement 4, and exits `0`) and `TestADeliveryRetryReturnsAnItemTheLaneParked` (with the owner fake of `internal/cli/deliver_retry_owner_proof_test.go`, the retry exits `0` and the item leaves `parked`, as for an item any other gate failure parked).
    - `internal/config/retired_lane_keys_test.go`: `TestTheRetiredLaneKeysLoadWithOneWarningEach`, with the subtests `project_config`, `user_config`, `values_of_every_kind` (a boolean, a text such as `"yes"`, a mapping, a list and nothing), `key_in_both_files_warns_once`, `only_retired_keys_leave_an_empty_delivery`, `unknown_key_still_fails` (today's message, `field bogus_key not found in type config.deliveryOverlay`) and `config_at_a_commit_prints_nothing`; each asserts the loaded configuration, the exact warnings and their count.
    - `internal/delivery/stored_lane_records_test.go`: `TestAPassedGateLeavesNoLaneWarningAndNoPullRequestParagraph` (an item reaches `publishing` with an empty warning and the Pull Request body is the publication's own) and `TestAStoredLaneBlockerKeepsItsParkClass`, with the subtests `container_blocker`, `change_failure_blocker` and `skip_segment_survives_the_low_disk_rewrite`.
11. MUST remove from `docs/user-guide/commands/deliver.md` the section `### Run the Linux Lane on demand` whole, the guide's documentation of the command this Task deletes, and change no other part of that page: the `#### Linux Lane` section, the Park Class table and the capacity and core-budget sentences are task_03's.
12. MUST remove the entry `command deliver lane` from `docs/references/skill-coverage.json` by hand, and change no other entry of the map. After the last edit of this Task, run, in this order, the two commands `_authorization.md` → Sanctioned regeneration declares for these records, and hand-edit none of their outputs:
    1. `go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1`, which writes `docs/references/coverage-record.json`: it loses the 36 names of requirement 8 and no other, and gains the ten names of requirement 10;
    2. `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`, the `-record-skill-coverage` step, which writes `docs/references/behavior-surfaces.json` without the surface `command deliver lane` and with the new fingerprints of the `deliver` commands, of `config keys` and of the guide's `deliver` page.
13. MUST keep `GOOS=windows go build -buildvcs=false ./...` and `GOOS=linux go build -buildvcs=false ./...` succeeding, since no file this Task touches carries a build constraint that would hide a broken reference on another platform.
14. MUST NOT change `internal/cli/cli_test.go`, `internal/docscontract/testdata/corpus-golden.json`, any skill, any other user-guide page, `CONTEXT.md`, any decision record, `.secondbrain-export`, the Makefile, any workflow, `go.mod`, `go.sum`, the Run Database's schema, or any file `_authorization.md` → Limits names. No test reaches the network, runs the real `docker` or `gh`, starts an Agent session, reads the real home directory or waits on elapsed time.

## Subtasks

- [ ] End the gate at the repository Verification and remove the lane's warning, Pull Request paragraph and constants from the engine, keeping the retired blocker's case.
- [ ] Delete the lane's command and its help lines, and doctor's check and codes.
- [ ] Register the two keys as deprecated and remove them from the loader and from `.roundfixrc.yml`.
- [ ] Delete the eight test files and the one test, and update the tests that only name the lane.
- [ ] Add the ten tests in seven files.
- [ ] Remove the guide's section on the removed command and the map's entry, then write the coverage record and the Behavior Surface Record with their commands.

## Acceptance Criteria

- [ ] Over a default branch that says `linux_lane: true`, the gate passes on a passed Verification without starting `docker`, writing a lane log or printing a lane line, and a passed item carries no lane warning and its Pull Request body no lane paragraph.
- [ ] `roundfix deliver lane` reproduces Surface Transcript 1, `roundfix deliver --help` reproduces Surface Transcript 2, and the root help names `<plan|start|status|resume|retry|stop>`.
- [ ] A configuration naming either key, with any value, in Project Config or User Config, loads with exactly one warning per key in Surface Transcript 3's form; any other unknown key still fails.
- [ ] Doctor prints no `linux-lane:` line and no `DR-LANE-` code, and every other line keeps its place.
- [ ] A stored `linux-lane-container` blocker is classed `environment` with the retired next step, a stored skip segment is printed as stored and survives a low-disk rewrite, and a Delivery Retry returns the parked item.
- [ ] The coverage record lost exactly the 36 names and gained the ten; the map and the Behavior Surface Record no longer hold `command deliver lane`; `./internal/cli`, `./internal/delivery`, `./internal/config`, `./internal/spec`'s coverage test and `./internal/docscontract` pass.

## Context

- instruction: `internal/cli/deliver_retry_owner_proof_test.go`
- instruction: `internal/delivery/engine_test.go`
- interface: `internal/cli/deliver.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/health.go`
- interface: `internal/delivery/engine.go`
- interface: `internal/delivery/park_class.go`
- interface: `internal/config/config.go`
- interface: `internal/cli/deliver_core_budget_test.go`
- interface: `internal/cli/deliver_retry_test.go`
- interface: `internal/cli/deliver_owner_log_test.go`
- interface: `internal/cli/doctor_disk_space_test.go`
- interface: `internal/cli/doctor_verification_cores_test.go`
- interface: `internal/cli/machine_verification_status_test.go`
- interface: `internal/cli/shared_cache_trim_test.go`
- interface: `internal/cli/external_skill_lock_test.go`
- interface: `internal/cli/external_skill_presence_test.go`
- interface: `internal/cli/releaseplan_release_preparation_test.go`
- interface: `internal/delivery/gate_blocker_test.go`
- interface: `internal/delivery/low_disk_warning_test.go`
- interface: `internal/delivery/park_class_test.go`
- interface: `.roundfixrc.yml`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `docs/references/coverage-record.json`
- interface: `docs/references/skill-coverage.json`
- interface: `docs/references/behavior-surfaces.json`
- creates: `internal/cli/deliver_gate_without_lane_test.go`
- creates: `internal/cli/deliver_lane_removed_test.go`
- creates: `internal/cli/retired_lane_keys_test.go`
- creates: `internal/cli/doctor_without_lane_test.go`
- creates: `internal/cli/deliver_stored_lane_records_test.go`
- creates: `internal/config/retired_lane_keys_test.go`
- creates: `internal/delivery/stored_lane_records_test.go`
- deletes: `internal/cli/deliver_linux_lane.go`
- deletes: `internal/cli/deliver_lane.go`
- deletes: `internal/cli/deliver_linux_lane_test.go`
- deletes: `internal/cli/deliver_lane_test.go`
- deletes: `internal/cli/deliver_linux_lane_go_list_test.go`
- deletes: `internal/cli/deliver_lane_skip_visible_test.go`
- deletes: `internal/cli/doctor_linux_lane_test.go`
- deletes: `internal/config/linux_lane_test.go`
- deletes: `internal/config/linux_lane_image_test.go`
- deletes: `internal/delivery/lane_skip_warning_test.go`

## Consumers

- the repository gate and `delivery.GateResult` — consumers: `./internal/delivery`, `./internal/cli`
- the loaded configuration and its registry of deprecated keys — consumers: `./internal/config`, `./internal/cli`
- the coverage record — consumers: `./internal/spec`
- the guide's `deliver` page, the Skill Coverage Map and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheRepositoryGateEndsAtTheRepositoryVerification|TestDeliverLaneReproducesSurfaceTranscriptUnknownCommand|TestDeliverHelpReproducesSurfaceTranscriptWithoutTheLane|TestTopLevelUsageNamesDeliverRetry|TestRunsListReproducesSurfaceTranscriptRetiredLaneKeys|TestDoctorRunsNoLinuxLaneCheck|TestStatusPrintsWhatTheLaneRecorded|TestADeliveryRetryReturnsAnItemTheLaneParked)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheRepositoryGateEndsAtTheRepositoryVerification TestDeliverLaneReproducesSurfaceTranscriptUnknownCommand/no_arguments TestDeliverLaneReproducesSurfaceTranscriptUnknownCommand/with_arguments TestDeliverLaneReproducesSurfaceTranscriptUnknownCommand/help_flag TestDeliverHelpReproducesSurfaceTranscriptWithoutTheLane TestTopLevelUsageNamesDeliverRetry TestRunsListReproducesSurfaceTranscriptRetiredLaneKeys TestDoctorRunsNoLinuxLaneCheck TestStatusPrintsWhatTheLaneRecorded TestADeliveryRetryReturnsAnItemTheLaneParked; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0, the gate ends at the Verification, Surface Transcripts 1, 2 and 3 are reproduced whole, doctor runs no lane check, and the lane's stored records stay readable and retryable.
- `out="$(go test -count=1 -timeout 20m -v -run '^TestTheRetiredLaneKeysLoadWithOneWarningEach$' ./internal/config 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in project_config user_config values_of_every_kind key_in_both_files_warns_once only_retired_keys_leave_an_empty_delivery unknown_key_still_fails config_at_a_commit_prints_nothing; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestTheRetiredLaneKeysLoadWithOneWarningEach/$name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0, both keys load with one warning each in every case and an unknown key still fails.
- `out="$(go test -count=1 -timeout 20m -v -run '^(TestAPassedGateLeavesNoLaneWarningAndNoPullRequestParagraph|TestAStoredLaneBlockerKeepsItsParkClass)$' ./internal/delivery 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestAPassedGateLeavesNoLaneWarningAndNoPullRequestParagraph TestAStoredLaneBlockerKeepsItsParkClass/container_blocker TestAStoredLaneBlockerKeepsItsParkClass/change_failure_blocker TestAStoredLaneBlockerKeepsItsParkClass/skip_segment_survives_the_low_disk_rewrite; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0, a passed gate leaves no lane trace and the stored blocker and skip keep their meaning.
- `for name in TestLinuxLaneRunsTheChangedPackagesInAContainer TestLinuxLaneSkipsWithANamedReason TestLinuxLaneFailureFailsTheGate TestLinuxLaneCancellationRemovesTheContainer TestGateWithoutTheLaneStartsNoDocker TestGateSkipsTheLaneWhenTheFetchedDefaultCannotBeRead TestLinuxLaneOptInPropagatesCancellationAndAnUnreadableConfig TestLinuxLaneRunsTheConfiguredImage TestLinuxLaneContainerFailureIsToldApartFromAChangeFailure TestLinuxLaneOutcomeClassification TestLinuxLaneMountsTheRepositoryHistoryReadOnly TestLinuxLaneSkipsAnObjectStoreItCannotMount TestLinuxLaneReportsThePackagesTheContainerSelected TestDeliverLanePassedTranscript TestDeliverLaneChangeFailureTranscript TestDeliverLaneContainerFailureTranscript TestDeliverLaneSkipTranscript TestDeliverLaneSelectsPackagesAndImage TestDeliverLaneCancellationRemovesContainer TestDeliverLaneObjectStoreSkipTranscript TestDeliverLanePackageCountTranscript TestDeliverLaneNoPackageCountTranscript TestDeliverLaneGoListFailureTranscript TestDeliverLaneGoListSignalTranscript TestLinuxLaneScriptNamesAFailedGoList TestLinuxLaneGoListFailureIsNamed TestGateReturnsTheLaneEnvironmentSkip TestDeliverStatusShowsTheLaneEnvironmentSkipTranscript TestDoctorLinuxLaneCheckReportsDockerAndTheImage TestDoctorPrintsTheLinuxLaneLineOnlyWhenTheLaneIsOn TestDeliveryLinuxLaneDefaultsOffAndProjectConfigEnablesIt TestDeliveryLinuxLaneImageDefaultsEmptyAndProjectConfigSetsIt TestDeliveryLinuxLaneImageRefusesAnInvalidReference TestAGateLaneSkipIsRecordedOnTheItemAndInThePullRequestBody TestAQuietGateLeavesNoLaneWarningAndReplacesAnEarlierOne TestDeliverLanePassesItsVerificationCoreBudget; do if grep -qF "\"$name\"" docs/references/coverage-record.json; then echo "coverage record still lists $name" >&2; exit 1; fi; done; for name in TestTheRepositoryGateEndsAtTheRepositoryVerification TestDeliverLaneReproducesSurfaceTranscriptUnknownCommand TestDeliverHelpReproducesSurfaceTranscriptWithoutTheLane TestRunsListReproducesSurfaceTranscriptRetiredLaneKeys TestDoctorRunsNoLinuxLaneCheck TestStatusPrintsWhatTheLaneRecorded TestADeliveryRetryReturnsAnItemTheLaneParked TestTheRetiredLaneKeysLoadWithOneWarningEach TestAPassedGateLeavesNoLaneWarningAndNoPullRequestParagraph TestAStoredLaneBlockerKeepsItsParkClass; do grep -qF "\"$name\"" docs/references/coverage-record.json || { echo "coverage record lacks $name" >&2; exit 1; }; done; if grep -qF '"command deliver lane"' docs/references/skill-coverage.json docs/references/behavior-surfaces.json; then echo 'command deliver lane is still a Behavior Surface' >&2; exit 1; fi; if tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- 'Run the Linux Lane on demand'; then echo 'the guide still documents the removed command' >&2; exit 1; fi; go test -count=1 -timeout 20m -run '^TestCoverageEquivalence$' ./internal/spec && go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0, the coverage record lost exactly the 36 lane tests and gained the ten new ones, the removed command is neither a Behavior Surface nor documented, and the coverage and documentation contracts pass.
- `test -d internal/cli && ! grep -rqE 'LinuxLane|linuxLane|LaneSkip|laneSkip|DR-LANE-|linux-lane\.log' internal/cli internal/delivery internal/config --include='*.go' --exclude='*_test.go' && ! grep -qE 'linux_lane|Linux Lane' .roundfixrc.yml && grep -qx '  require_defect_class: true' .roundfixrc.yml && GOOS=windows go build -buildvcs=false ./... && GOOS=linux go build -buildvcs=false ./... && go test -count=1 -timeout 40m ./internal/cli ./internal/delivery ./internal/config` — expected: exit 0, no production code names the lane, the Project Config holds no lane key, both other platforms build, and the consumers of the gate result and of the configuration pass whole.

## References

- `_prd.md` → Goals 1 and 2; User Stories 1, 2, 4 and 6; Core Features 1, 2, 3, 4, 5, 10 and 13; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Decisions
- `_techspec.md` → The repository gate; The command and the help; The two retired keys; Doctor; Interfaces; API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 10; API Contract 13; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Invariants; Testing Approach; Build Order 1; Decisions (the Delegated Decision of 2026-11-24 on the stored `linux-lane-container` blocker)
- `_authorization.md` → What the grant covers; Sanctioned regeneration; Limits
- [`references/2026-11-24-techspec-measurements.md`](references/2026-11-24-techspec-measurements.md), sections 2, 3, 4, 5, 6 and 7
- ADR-0288, ADR-0294, ADR-0297, ADR-0303, ADR-0316, ADR-0365, ADR-0384, ADR-0385, ADR-0398, ADR-0450

## Result

Skills loaded: implement-task, coding-guidelines, golang-cli, golang-testing, testing-boss, agentic-cli-design, golang-structs-interfaces, tech-writer, evidence-gate, systematic-debugging.

Implemented this Task's product slice. Status and the Task Graph remain Daemon-owned. Initial inspection found only this Task file already modified; its existing status was preserved. The executor rejected the configured symlinked writable root, including with the physical working directory; approved escalated execution provided access to the same checkout.

Acceptance evidence:

- **Repository gate and publication:** the gate returns immediately after repository Verification. The new gate regression uses the surviving item fixture, committed legacy configuration on the fetched default branch, a Git boundary fake that refuses configuration reads, and a recording Docker stand-in in the command environment. It asserts success, no Docker invocation, no lane log and no lane console output. The engine regression observes `publishing`, an empty warning and the publication's original body. The surviving low-disk test also proves an existing skip warning is preserved through gating and publication without a PR paragraph.
- **Command and help:** removed both lane source files and dispatch/help entries. The new command tests pin all stdout, stderr and exit codes for `deliver lane`, `deliver lane --all HEAD`, `deliver lane --help` and `deliver --help`. The surviving root-help test pins the updated root line.
- **Configuration:** removed the fields, overlay, validation, image pattern and merges; appended the two deprecated-key entries in the specified order. Tests exercise each key alone, both together, project/user scopes, duplicate scopes, boolean/text/mapping/list/null values, an otherwise empty delivery mapping, the unchanged unknown-key diagnostic and silent resolution at a commit. Project Config lost exactly the key and its six comment lines.
- **Doctor:** removed its dependency, check, health name and diagnostics. Its new regression supplies legacy checkout configuration and a recording Docker stand-in, checks no lane line/code/invocation, and checks `verification-cores` immediately precedes `codex`. Other doctor fixtures retain their checks and output positions.
- **Stored records:** the retired container blocker remains an environment park with the prescribed retry action; a change-failure blocker remains a finding. Status prints the stored warning and retired park action. Retry uses the existing owner-proof fixture, with a parked legacy item and a merged item holding a stored skip warning; it leaves parked with an empty blocker and preserves the merged item. Low-disk rewriting preserves other segments in order and appends low disk last.
- **Records and documentation:** deleted the eight obsolete test files and the one obsolete core-budget test, keeping the other tests by name. Removed only the guide's on-demand command section and only `command deliver lane` from the map. Both generated records were written by their sanctioned commands, in the required order, after the final source/test edit. Both records now omit the removed command.

Focused checks (all exited 0 on their final runs):

```text
go test ./internal/cli ./internal/config ./internal/delivery -run 'Test(TheRetiredLane|DeliverLaneReproduces|DeliverHelpReproduces|RunsListReproducesSurfaceTranscriptRetired|DoctorRunsNoLinux|TheRepositoryGateEnds|APassedGateLeaves|AStoredLaneBlocker|TheRepositoryGateCommandReceives|DoctorReportsTheVerificationCoreBudget|TheQueueWarnsOfLowDisk)' -count=1

rtk proxy go test ./internal/cli ./internal/delivery ./internal/config -run 'Test(TheRetiredLane|DeliverLaneReproduces|DeliverHelpReproduces|RunsListReproducesSurfaceTranscriptRetired|DoctorRunsNoLinux|TheRepositoryGateEnds|APassedGateLeaves|AStoredLaneBlocker|StatusPrintsWhatTheLaneRecorded|ADeliveryRetryReturnsAnItemTheLaneParked)' -count=1

rtk proxy go test ./internal/cli -run '^TestADeliveryRetryReturnsAnItemTheLaneParked$' -count=1

go test ./internal/spec -run '^TestCoverageEquivalence$' -update-coverage-record -count=1
go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage

rtk proxy git diff --check
```

A separate focused selection also exercised the surviving owner-capacity, gate-capacity, shared-cache trim, disk-space doctor, machine-capacity doctor, Park Class, gate-blocker and root-help tests; those passed. Its initial new status assertion used the wrong punctuation; the preserved command transcript established the correct `Park: <slug> environment: <next>` form, and the corrected new test passed in the replacement-test selection above. `systematic-debugging` was also used to inspect the removal's compile diagnostics and correct the orphaned syntax/imports before the passing focused runs.

The generated coverage delta is **36 removed and 18 added**: all ten authored replacement tests plus eight tests already present in HEAD but absent from the starting coverage record (four review-model CLI tests and four recovery/review documentation tests from earlier work). No unrelated test was removed, no existing test source was added by this Task, and the generated output was not hand-edited. This is baseline record drift exposed by the required generator, rather than an exact ten-entry addition.

A changed-path audit compared tracked and untracked changes with this Task's declared Context paths: none fell outside the declarations. A production-source scan found no `LinuxLane`, `linuxLane`, `LaneSkip`, `laneSkip`, `DR-LANE-` or `linux-lane.log` references in the three affected packages.

Declared Verification, whole consumer suites, the full documentation contracts and cross-platform builds remain for the Daemon; they were not run here. The remaining guide/skill text and decision/glossary changes belong to task_03 and task_04. No commit, push, PR or ai-memory message was made.

## Carry-forward provenance

- Source Run: `run_20261124T213033Z_283908b7147be213`
- Source commit: `acf74f278ba9bf265668975a9a115a993f2a8411`
