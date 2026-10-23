---
task: task_02
spec: 0303-the-linux-lane-names-the-step-that-failed
status: completed
type: backend
complexity: medium
---

# Task 02: The lane script marks a failed go list and the lane names it with its Park Class

## Overview

The Linux Lane script selects packages with `go list` before `go test`, but a
failed `go list` ends the script with its own exit code, so the gate, its
blocker, the lane log and `roundfix deliver lane` report `go test exited <n>`
for a test run that never started, and a `go list` ended by a signal as the
container's signal. This Task gives the failed `go list` its own exit `4` and
marker line, and makes the classifier name `go list` and its code with the
Park Class ADR-0316 sets, while every other outcome keeps its bytes.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-22,
   [`references/2026-10-22-the-linux-lane-reports-a-failed-go-list-as-go-test.md`](references/2026-10-22-the-linux-lane-reports-a-failed-go-list-as-go-test.md)
   (`_prd.md` → Sources 1), by implementing API Contracts 1-4 of
   `_techspec.md`.
2. MUST change `linuxLaneScript` in `internal/cli/deliver_linux_lane.go`
   exactly as API Contract 1 states: the `&& if` after the package-selection
   `go list` becomes `|| { code=$?; printf 'roundfix-linux-lane: go list:
   exited %s\n' "$code" >&2; exit 4; }; if`, and every other byte of the
   script stays.
3. MUST make `classifyLinuxLaneOutcome` classify exit `4` exactly as API
   Contract 2 states, after the cap and before the other codes: the last
   whole output line, with one trailing `\r` dropped, that starts with
   `roundfix-linux-lane: go list: exited ` and ends in ASCII digits that
   `strconv.Atoi` accepts names `<n>`; below `128` it is `Linux Lane change
   failure: go list exited <n>` with `delivery.BlockerGateFailed`; `128` or
   more is `Linux Lane container failure: go list exited <n>: the process
   ended by signal <n-128>` with `delivery.BlockerLinuxLaneContainer`; no such
   line is `Linux Lane change failure: the lane script exited 4 without
   reporting the go list exit code` with `delivery.BlockerGateFailed`.
4. MUST keep API Contract 3: every other exit, the cap and exit `0` keep their
   summaries and blockers byte for byte, and `runLinuxLane`,
   `internal/cli/deliver_lane.go`, `internal/cli/deliver_workflow.go` and
   every file under `internal/delivery` that is not a test keep their bytes.
5. MUST state the hazard in the code's behavior: the outcome decides whether
   the operator repairs Docker or the change, so only the script's own exit
   `4` with a whole marker line names `go list`, and a code of `128` or more
   stays a container outcome.
6. MUST add and change the tests `_techspec.md` → Testing Approach lists, each
   new top-level test calling `t.Parallel()` first, compiling for Windows,
   using temporary repositories and passing no real-time duration:
   - `TestLinuxLaneScriptNamesAFailedGoList` and
     `TestLinuxLaneGoListFailureIsNamed` in the new file
     `internal/cli/deliver_linux_lane_go_list_test.go`; the first cuts
     `linuxLaneScript` at its first `pkgs="$(go list`, fails when that text is
     absent, and runs the tail with `sh -ceu` and a `PATH` holding only a
     stand-in `go` built with `testfixture.FixtureBinary`, with the four rows
     listed there, asserting the exit code, that the marker is a whole line,
     and whether `go test` ran;
   - `TestDeliverLaneGoListFailureTranscript` and
     `TestDeliverLaneGoListSignalTranscript` in
     `internal/cli/deliver_lane_test.go`, asserting with
     `assertDeliverLaneTranscript` the two whole outputs that `_techspec.md`
     → Surface Transcripts describes (standard output, empty standard error,
     exit `1` and `3`);
   - in `internal/cli/deliver_linux_lane_test.go`, the row `other change
     exit` of `TestLinuxLaneOutcomeClassification` moves from `4` to `5`
     (`go test exited 5`), the table gains the ten `go list` rows listed
     there, and the independent script literal of
     `TestLinuxLaneRunsTheChangedPackagesInAContainer` gains the branch of
     API Contract 1;
   - in `internal/delivery/gate_blocker_test.go`,
     `TestGateParksWithTheBlockerItNames` gains the rows `go list change`
     (`finding`) and `go list signal` (`environment`).
   No other existing test, golden or fixture changes.

## Subtasks

- [ ] Add the `go list` branch to the lane script.
- [ ] Classify exit `4` through the marker line with its Park Class.
- [ ] Add the script, gate, transcript and Park Class tests and update the pinned script and classification rows.

## Acceptance Criteria

- [ ] A failed `go list` exits `4` with its marker and never starts `go test`.
- [ ] Exit `4` names `go list exited <n>` as a change failure below `128` and as a container failure with its signal from `128`.
- [ ] Every other outcome, `go test exited <n>` included, keeps its bytes, and `roundfix deliver lane` exits `1` or `3` by the blocker as before.

## Consumers

- `linuxLaneScript`'s `go list` marker and exit `4` — consumers: `./internal/cli`
- `classifyLinuxLaneOutcome`'s `go list` summaries and blockers — consumers: `./internal/cli`, `./internal/delivery`

## Context

- instruction: `docs/adr/0316-the-linux-lane-names-a-failed-go-list-apart-from-go-test.md`
- instruction: `docs/adr/0294-the-linux-lane-runs-like-ci.md`
- instruction: `docs/user-guide/commands/deliver.md`
- instruction: `internal/cli/deliver_lane.go`
- instruction: `internal/cli/deliver_workflow.go`
- instruction: `internal/delivery/park_class.go`
- instruction: `internal/testfixture/external_tool.go`
- creates: `internal/cli/deliver_linux_lane_go_list_test.go`
- interface: `internal/cli/deliver_linux_lane.go`
- interface: `internal/cli/deliver_linux_lane_test.go`
- interface: `internal/cli/deliver_lane_test.go`
- interface: `internal/delivery/gate_blocker_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestLinuxLaneScriptNamesAFailedGoList|TestLinuxLaneGoListFailureIsNamed|TestDeliverLaneGoListFailureTranscript|TestDeliverLaneGoListSignalTranscript|TestLinuxLaneOutcomeClassification|TestLinuxLaneRunsTheChangedPackagesInAContainer|TestLinuxLaneFailureFailsTheGate|TestLinuxLaneContainerFailureIsToldApartFromAChangeFailure|TestLinuxLaneReportsThePackagesTheContainerSelected|TestDeliverLaneChangeFailureTranscript|TestDeliverLaneContainerFailureTranscript|TestGateParksWithTheBlockerItNames)$' ./internal/cli ./internal/delivery 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestLinuxLaneScriptNamesAFailedGoList TestLinuxLaneGoListFailureIsNamed TestDeliverLaneGoListFailureTranscript TestDeliverLaneGoListSignalTranscript TestLinuxLaneOutcomeClassification TestLinuxLaneRunsTheChangedPackagesInAContainer TestLinuxLaneFailureFailsTheGate TestLinuxLaneContainerFailureIsToldApartFromAChangeFailure TestLinuxLaneReportsThePackagesTheContainerSelected TestDeliverLaneChangeFailureTranscript TestDeliverLaneContainerFailureTranscript TestGateParksWithTheBlockerItNames; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for row in 'go list killed' 'go list signal zero' 'go list embedded marker ignored'; do grep -qF -- "$row" internal/cli/deliver_linux_lane_test.go || { printf 'missing classification row: %s\n' "$row" >&2; exit 1; }; done; GOOS=windows go vet ./internal/cli ./internal/delivery` — expected: exit 0. Before this Task the four new tests do not exist, so the command fails at its first pass check. Afterwards the script, gate, transcript, classification and Park Class tests pass with the unchanged lane failure, container failure, package count and transcript tests, the classification table holds the `go list` rows, and both packages' tests compile for Windows.

## References

- `_prd.md` → Sources; Goals; Core Features 1-2; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5
- `_techspec.md` → Current behavior; API Contracts 1-4; Surface Transcripts; Invariants 1-5; Testing Approach; Build Order 2
- ADR-0316


## Result

Implemented the API Contracts 1–4 slice. A failed package-selection `go list`
prints its whole marker line, exits the script with `4`, and never starts
`go test`. The classifier handles that reserved exit after the cap, keeps the
last valid ASCII-decimal marker (dropping one trailing carriage return), and
uses the change or container blocker at the `128` boundary. Its comment names
the operator hazard: choosing whether to repair Docker or the change.

Acceptance evidence from focused checks:

- Failed `go list` exits `4` without testing: the four rows of
  `TestLinuxLaneScriptNamesAFailedGoList` pass against the real `sh -ceu`
  script tail with a compiled Go stand-in. Before the production edit, the
  same test failed for list exits `1` and `137`: the script returned those
  codes with empty output instead of exit `4` and the marker. Successful
  selection retains the failing test's exit `1`; empty selection exits `0`.
- Exit `4` names and classifies the selection failure: the ten new rows in
  `TestLinuxLaneOutcomeClassification` pass, including `127`, `128`, `137`,
  last-marker selection, missing/embedded/signed/spaced markers and CRLF.
  `TestLinuxLaneGoListFailureIsNamed` passes for change and signal results,
  asserting gate reason, blocker and the lane log's final summary.
  The two new rows in `TestGateParksWithTheBlockerItNames` pass with
  `finding` and `environment` respectively.
- Other outcomes and command behavior: the unchanged classification rows,
  independent container-script assertion, existing gate failures, package
  counts and change/container transcripts pass. Both new whole command
  transcripts pass, asserting standard output, empty standard error and
  exits `1` and `3`. A byte comparison against `HEAD` confirms the exact
  script substitution with every other script byte preserved, an unchanged
  `runLinuxLane`, unchanged `deliver_lane.go` and `deliver_workflow.go`, and
  no production change under `internal/delivery`.

Focused commands (each exited `0` after the implementation):

- `go test ./internal/cli -count=1 -run 'TestLinuxLane(ScriptNamesAFailedGoList|GoListFailureIsNamed|OutcomeClassification|RunsTheChangedPackagesInAContainer)$|TestDeliverLaneGoList'`
- `go test ./internal/delivery -count=1 -run '^TestGateParksWithTheBlockerItNames$/go_list'`
- `go test ./internal/cli -count=1 -run 'TestLinuxLane(FailureFailsTheGate|ContainerFailureIsToldApartFromAChangeFailure|ReportsThePackagesTheContainerSelected)$|TestDeliverLane(ChangeFailure|ContainerFailure)Transcript$'`
- `build_dir=$(mktemp -d /tmp/roundfix-0303-windows.XXXXXX)` followed by
  `GOOS=windows go test -c -o "$build_dir/" ./internal/cli ./internal/delivery`
  compiled both Windows test binaries without running them.
- `git diff --check` and the focused byte-comparison probe above.

All new top-level tests begin with `t.Parallel()`, use disposable repositories
and compiled stand-ins, and pass no real-time duration. No real container or
network service was exercised. The checkout's symlinked writable root was
rejected by the sandbox; checks and edits used its resolved Run path with
approved escalation. The pre-existing Task status remains Daemon-owned.
Declared Task Verification and repository Verification were not run; the
Daemon owns both and settlement. No follow-up outside this slice was found.
