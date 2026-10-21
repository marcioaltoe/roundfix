---
task: task_01
spec: 0291-the-queue-shows-a-skipped-linux-lane-and-waits-for-a-run-slot
status: pending
type: backend
complexity: high
---

# Task 01: A Lane Environment Skip reaches the item and the Pull Request body

## Overview

With `delivery.linux_lane: true`, the gate of three Specs skipped the Linux
Lane on a missing image and merged them with one console notice
(`_prd.md` → Sources 1). This Task makes the lane tell an expected skip from a
Lane Environment Skip, makes the gate return the environment skip's reason
with its passing result, and makes the Delivery Engine record it on the item,
where `roundfix deliver status` prints it, and in the Pull Request body. It is
verifiable on its own through the lane's compiled Docker stand-in and the
engine fakes.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-20,
   [`references/2026-10-20-the-queue-skips-the-linux-lane-without-a-trace.md`](references/2026-10-20-the-queue-skips-the-linux-lane-without-a-trace.md)
   (`_prd.md` → Sources 1), by implementing API Contracts 1 to 3 of
   `_techspec.md` in `internal/cli/deliver_linux_lane.go`,
   `internal/cli/deliver_workflow.go` and `internal/delivery/engine.go`, with
   `SkipExpected`, `GateResult.LaneSkip`, `WarningLinuxLaneSkipped` and
   `LaneSkipBody` as `_techspec.md` → Interfaces sketches them.
2. MUST set `SkipExpected` only for `no go.mod at the repository root` and `no
   changed Go package`, so every other skip reason, present or future, is a
   Lane Environment Skip. The console notice, the lane log, the reasons, their
   order and `roundfix deliver lane`'s exit codes do not change.
3. MUST state the hazard in the code's behavior, as API Contracts 2 and 3 do.
   The data source is the lane result and `linuxLaneOptIn`'s unreadable
   reason, both computed on the Queue Owner's machine. The destinations are
   the item warning in the Run Database, which `roundfix deliver status`
   prints, and the Pull Request body published to the forge. The reason holds
   only the lane's own text (a path, an image name, Docker's first output
   line) and never an environment value. A skip still passes the gate and
   never parks the item.
4. MUST make `gateCandidate` replace any earlier `linux-lane-skipped: `
   segment, so a later gate that ran the lane or skipped it for another reason
   leaves only its own outcome, and keep the segment last in the warning.
5. MUST keep every existing test's name and assertions, including
   `TestLinuxLaneSkipsWithANamedReason`,
   `TestGateSkipsTheLaneWhenTheFetchedDefaultCannotBeRead`,
   `TestGateWithoutTheLaneStartsNoDocker`, `TestLinuxLaneFailureFailsTheGate`
   and `TestDeliverStatusPrintsAParkLinePerParkedItem`; this Task declares no
   break.
6. MUST add these tests, each top-level test calling `t.Parallel()` first,
   compiling for Windows, ending every child it starts in `t.Cleanup`, using
   temporary homes and repositories, and declaring no real-time budget under
   ten seconds:
   - in `internal/cli/deliver_lane_skip_visible_test.go`,
     `TestGateReturnsTheLaneEnvironmentSkip` with the subtests
     `_techspec.md` → Testing Approach lists, reusing `linuxLaneTools` and
     `newLinuxLaneFixture`, and `TestDeliverStatusShowsTheLaneEnvironmentSkipTranscript`,
     asserting Surface Transcript 1 of `_techspec.md` whole: every standard
     output line in order, empty standard error and exit `0`, with the
     warning seeded as the engine records it (`WarningLinuxLaneSkipped`
     joined to the reason);
   - in `internal/delivery/lane_skip_warning_test.go`,
     `TestAGateLaneSkipIsRecordedOnTheItemAndInThePullRequestBody` and
     `TestAQuietGateLeavesNoLaneWarningAndReplacesAnEarlierOne`, with the
     cases `_techspec.md` → Testing Approach lists, asserting the exact body
     the fake pull request boundary receives.
7. MUST document the change in the Linux Lane section of
   `docs/user-guide/commands/deliver.md`: after the skip reasons, state that
   `no go.mod at the repository root` and `no changed Go package` are expected
   skips that leave only the notice, and that every other skip of an enabled
   lane, and a fetched default whose configuration cannot be read, is a Lane
   Environment Skip: the item records the warning `linux-lane-skipped:
   <reason>`, which `roundfix deliver status` prints on its `Warning:` line,
   and the Pull Request body adds the paragraph `Linux Lane: not run. The
   Delivery Queue gate skipped it: <reason>. This candidate's Go packages were
   not tested on Linux before this pull request.` MUST search the guides, the
   agent guides and the skill copies for wording that says a lane skip leaves
   only the notice, and change each one this Task owns; the Roundfix
   `deliver` reference is task_02's.
8. MUST re-record the Behavior Surface Record, because this Task edits a user
   guide page: `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`.
   The Agent runs it itself before handing back; no record is hand-edited.

## Subtasks

- [ ] Mark the two expected skips in the lane result.
- [ ] Return the environment skip's reason from the gate.
- [ ] Record and replace the item warning, and extend the Pull Request body.
- [ ] Add the gate, status and engine tests.
- [ ] Document the Lane Environment Skip in the `deliver` guide and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] An enabled lane skipped because its image is absent passes the gate and leaves `linux-lane-skipped: <reason>` on the item and the paragraph in the Pull Request body.
- [ ] `roundfix deliver status` reproduces Surface Transcript 1.
- [ ] An expected skip, a disabled lane and a lane that ran leave no warning, and a later gate replaces an earlier one.
- [ ] Every existing lane, gate and status test passes unchanged.

## Consumers

- `delivery.GateResult` and the item warning `linux-lane-skipped` — consumers: `./internal/delivery`, `./internal/cli`
- the Linux Lane result — consumers: `./internal/cli`

## Context

- instruction: `docs/adr/0303-a-linux-lane-skipped-for-its-environment-is-shown-where-the-delivery-is-read.md`
- instruction: `docs/adr/0288-the-delivery-queue-gate-runs-the-changed-go-packages-in-a-linux-lane.md`
- instruction: `docs/adr/0294-the-linux-lane-runs-like-ci.md`
- instruction: `internal/cli/deliver_linux_lane_test.go`
- instruction: `internal/cli/deliver_park_status_test.go`
- instruction: `internal/delivery/engine_test.go`
- instruction: `internal/cli/deliver.go`
- interface: `internal/cli/deliver_linux_lane.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/delivery/engine.go`
- creates: `internal/cli/deliver_lane_skip_visible_test.go`
- creates: `internal/delivery/lane_skip_warning_test.go`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestGateReturnsTheLaneEnvironmentSkip|TestDeliverStatusShowsTheLaneEnvironmentSkipTranscript|TestLinuxLaneSkipsWithANamedReason|TestGateSkipsTheLaneWhenTheFetchedDefaultCannotBeRead|TestGateWithoutTheLaneStartsNoDocker|TestLinuxLaneFailureFailsTheGate|TestDeliverStatusPrintsAParkLinePerParkedItem)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestGateReturnsTheLaneEnvironmentSkip TestDeliverStatusShowsTheLaneEnvironmentSkipTranscript TestLinuxLaneSkipsWithANamedReason TestGateSkipsTheLaneWhenTheFetchedDefaultCannotBeRead TestGateWithoutTheLaneStartsNoDocker TestDeliverStatusPrintsAParkLinePerParkedItem; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v ./internal/delivery 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestAGateLaneSkipIsRecordedOnTheItemAndInThePullRequestBody TestAQuietGateLeavesNoLaneWarningAndReplacesAnEarlierOne; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- 'linux-lane-skipped: <reason>' || { echo 'missing warning in the deliver guide' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/deliver.md | grep -qF -- 'Linux Lane: not run. The Delivery Queue gate skipped it: <reason>.' || { echo 'missing body paragraph in the deliver guide' >&2; exit 1; }; GOOS=windows go build -buildvcs=false ./internal/cli ./internal/delivery || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the new tests do not exist, so the command fails at its first pass check. Afterwards the gate returns each environment skip, status reproduces the transcript, the engine records and replaces the warning and extends the body, every existing lane, gate, status and delivery test passes, the guide documents the Lane Environment Skip, both packages build for Windows, and the Behavior Surface Record is current.

## References

- `_prd.md` → Sources; Goals; Core Feature 1; Core Feature 5; Success Metric 1; Success Metric 2
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 9; Surface Transcript 1; Invariants 1-2, 7; Interfaces; Data Models; Current behavior; Testing Approach; Build Order 1
- ADR-0303; ADR-0288; ADR-0294; ADR-0297; ADR-0184; ADR-0259; ADR-0271; ADR-0285; ADR-0280
