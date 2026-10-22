---
spec: 0303-the-linux-lane-names-the-step-that-failed
prd: _prd.md
created: 2026-10-22
---

# The Linux Lane names the step that failed — Technical Spec

## Executive Summary

`linuxLaneScript` in `internal/cli/deliver_linux_lane.go` runs
`pkgs="$(go list -e ...)" && ... && exec go test -count=1 "$@"` under
`sh -ceu`, so a failed `go list` ends the script with `go list`'s own exit
code, and `classifyLinuxLaneOutcome` reads it as `go test exited <n>` or, for
a signal, as the container's signal. The script gains an explicit branch: a
failed `go list` prints `roundfix-linux-lane: go list: exited <n>` and exits
`4`, a code no other step and no `go test` outcome produces. The classifier
reads exit `4` through the last well-formed marker line: a code below `128` is
a Lane Change Failure with cause `go list exited <n>`, `128` or more a Lane
Container Failure with cause `go list exited <n>: the process ended by signal
<n-128>`. The gate, `roundfix deliver lane` and the Park Class mapping already
follow the blocker prefix, so they need no change.

The trade-off accepted: exit `4`, which the classifier read as `go test exited
4` before, now always means the script's `go list` branch; `go test` never
exits `4` (ADR-0294's measurement), so no `go test` outcome changes.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; every test uses temporary repositories,
  compiled stand-ins and `sh`. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — Build Orders 1 and 2 follow
  ADR-0316, ADR-0294, ADR-0288, ADR-0297 and ADR-0303, named with their titles
  in `_prd.md` → Project Constraints. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AP approval of 2026-10-22 ("same as cycle AO"), carrying the cycle AJ grant
  "Mexer em Governed Paths declarados". Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0303-the-linux-lane-names-the-step-that-failed/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`,
  `.agents/skills/roundfix/references/deliver.md`,
  `skills/roundfix/SKILL.md`.

## Current behavior

Read and characterized on 2026-10-22 at `42f523c2`. Line numbers move when
other Specs merge; the function names do not.

1. `linuxLaneScript` (`internal/cli/deliver_linux_lane.go:27`) checks its
   tools, the history mount, unpacks and records the candidate, each failure
   printing `roundfix-linux-lane: container: <cause>` and exiting `3` through
   `fail`. It then runs `pkgs="$(go list -e -f '{{if or .GoFiles .TestGoFiles
   .XTestGoFiles}}{{.ImportPath}}{{end}}' "$@")" && if [ -z "$pkgs" ]; then
   ...; exit 0; fi && set -- $pkgs && echo "roundfix-linux-lane: packages: $#"
   && exec go test -count=1 "$@"`. Under `sh -e` a failed command inside an
   `&&` list does not stop the shell early, but the list is the script's last
   command, so the script exits with `go list`'s status.
2. Characterization of that tail, cut from the constant and run with
   `/bin/sh -ceu` (bash in POSIX mode) and `/bin/dash -ceu` with a stand-in
   `go` on `PATH`: `go list` exiting `1`, `2` and `137` made the script exit
   `1`, `2` and `137`, with `go list`'s standard error and no line naming
   `go list`. The proposed tail of API Contract 1 exited `4` with
   `roundfix-linux-lane: go list: exited 1`, `... 137` in both shells; with
   `go list` succeeding it printed `roundfix-linux-lane: packages: 1` and
   exited with the stand-in `go test`'s `1`, and with empty output it printed
   `no package builds on linux` and exited `0`.
3. `classifyLinuxLaneOutcome` (`:280`) returns, for a started run: the cap;
   exit `3` through the last `roundfix-linux-lane: container: ` line; `126`;
   `127`; `128` or more as `docker run exited <n>: the container process ended
   by signal <n-128>`; and every other non-zero exit as `Linux Lane change
   failure: go test exited <n>` with `delivery.BlockerGateFailed` (`:300-301`).
   So the characterized `go list` exits `1` and `2` read as `go test exited
   1`/`2`, and `137` as the container's signal 9.
4. `runLinuxLane` (`:61`) passes the code to the classifier through its local
   `failure` closure, writes the summary as the lane log's last line, and sets
   `Reason` to `<summary>; see <log>` and `Blocker` to `<blocker>: <detail>;
   see <log>`, where a container blocker drops the `Linux Lane container
   failure: ` prefix. `commandDeliveryWorkflow.Gate`
   (`internal/cli/deliver_workflow.go`, the `runLinuxLane` call) returns that
   reason and blocker; `delivery.ClassifyPark`
   (`internal/delivery/park_class.go`) maps `linux-lane-container` to
   `environment` and `gate-failed` to `finding`. `runDeliverLane`
   (`internal/cli/deliver_lane.go`) prints `roundfix: <reason>` and exits `3`
   for a `linux-lane-container` blocker and `1` otherwise.
5. `go list -e` with the lane's environment (`HOME`, `GOMODCACHE` an empty
   directory, `GOFLAGS=-mod=readonly`, `GOPROXY=off`, `GOTOOLCHAIN=local`,
   `GOSUMDB=off`; go 1.26.5 on the host, no container) and the lane's `-f`
   template exited `1` for a `go.mod` with an unterminated `require (` block
   (`go: errors parsing go.mod`), for `go 1.99` (`go: go.mod requires go >=
   1.99 (running go 1.26.5; GOTOOLCHAIN=local)`) and for a `go.work` that uses
   a missing directory; it exited `0` for a requirement absent from the module
   cache, an import of an unknown module, a Go syntax error, a `toolchain`
   line above the running one and a pattern naming a missing directory. `go
   test` exited `1` in every module case.
6. Tests that pin the current bytes:
   `TestLinuxLaneRunsTheChangedPackagesInAContainer`
   (`internal/cli/deliver_linux_lane_test.go`) holds the whole script as an
   independent literal; `TestLinuxLaneOutcomeClassification` holds the row
   `{"other change exit", 4, ..., "Linux Lane change failure: go test exited
   4", ...}`; `TestLinuxLaneFailureFailsTheGate` (codes `1`, `2`),
   `TestDeliverLaneChangeFailureTranscript`
   (`internal/cli/deliver_lane_test.go`) and `TestGateParksWithTheBlockerItNames`
   (`internal/delivery/gate_blocker_test.go`) hold `go test exited 1`.

Sweep of every producer and reader of the lane's outcome text (`rg -n 'go test
exited|exited 3 before|pkgs=' internal docs .agents skills`): the classifier
and script above; the tests in item 6; `docs/user-guide/commands/deliver.md`
and both copies of the Roundfix `deliver` reference (the Lane Container
Failure cause list, the Lane Change Failure paragraph, the gate reason and
blocker block, the lane log's outcome lines and the `deliver lane` output
lines); ADR-0294 and ADR-0288's amendment note. No other code reads the text.

## System Architecture

- `internal/cli/deliver_linux_lane.go`: `linuxLaneScript`'s `go list` branch
  and `classifyLinuxLaneOutcome`'s exit `4`.
- `docs/user-guide/commands/deliver.md`, the Roundfix `deliver` reference and
  its mirror, the Roundfix skill's versions, the Behavior Surface Record,
  `CONTEXT.md`.

No package, import edge, exit code of a command, or Park Class mapping
changes.

## Implementation Design

### Interfaces

No signature changes. `classifyLinuxLaneOutcome(code int, capped bool,
output []byte) (summary, blocker string)` gains a case for `code == 4`.

### Data Models

None.

### API Contracts

1. API Contract: the script branch. In `linuxLaneScript`, the text
   `pkgs="$(go list -e -f '{{if or .GoFiles .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' "$@")" && if`
   becomes
   `pkgs="$(go list -e -f '{{if or .GoFiles .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' "$@")" || { code=$?; printf 'roundfix-linux-lane: go list: exited %s\n' "$code" >&2; exit 4; }; if`;
   every other byte of the script stays. `go test` starts only after a
   successful `go list`.
2. API Contract: the classification of exit `4`, placed after the cap and
   before the `126` case. The classifier reads every output line, drops one
   trailing `\r`, and keeps the last line that starts with
   `roundfix-linux-lane: go list: exited ` and whose remainder is one or more
   ASCII digits that `strconv.Atoi` accepts:
   - a code `n` below `128`: summary `Linux Lane change failure: go list
     exited <n>`, blocker `delivery.BlockerGateFailed`;
   - a code of `128` or more: summary `Linux Lane container failure: go list
     exited <n>: the process ended by signal <n-128>`, blocker
     `delivery.BlockerLinuxLaneContainer`;
   - no such line: summary `Linux Lane change failure: the lane script exited
     4 without reporting the go list exit code`, blocker
     `delivery.BlockerGateFailed`.
   A line holding the prefix after other text, a signed, spaced or
   non-decimal code, and a code `Atoi` rejects are ignored.
3. API Contract: unchanged outcomes. The cap, exit `0`, `3`, `126`, `127`,
   `128` or more and every other non-zero exit keep today's summaries and
   blockers byte for byte; `go test exited <n>` stays the cause of every exit
   other than `3`, `4`, `126`, `127` and `128` or more.
4. API Contract: the surfaces. `runLinuxLane`, `Gate`, `runDeliverLane` and
   `delivery.ClassifyPark` are unchanged and carry the new summaries: the
   gate's reason `<summary>; see <lane log>`, the blocker `gate-failed:
   <summary>; see <lane log>` or `linux-lane-container: go list exited <n>:
   the process ended by signal <n-128>; see <lane log>`, the lane log's last
   line `<summary>`, and `roundfix deliver lane`'s line `roundfix: <summary>;
   see <lane log>` with exit `1` or `3`.
5. API Contract: documentation. The `deliver` guide and both reference copies
   name the `go list` branch and its three causes, the gate reason and blocker
   for `go list`, the lane log line and the `deliver lane` line; the glossary
   revises **Lane Change Failure**.

### Surface Transcripts

None. No command, flag, stream placement or exit-code mapping changes, and
running `roundfix deliver lane` through the built product needs a Docker
daemon, which no QA row may call. `roundfix deliver lane` keeps its one
outcome line on standard output, empty standard error and its exits `1` and
`3`; only the cause text of a `go list` failure is new. The tests of API
Contract 4 assert the whole output of the command with the existing Docker
stand-in: with the stand-in exiting `4` after `roundfix-linux-lane: go list:
exited 1`, standard output is exactly `roundfix: Linux Lane change failure:
go list exited 1; see <lane log>`, standard error is empty and the exit is
`1`; after `roundfix-linux-lane: go list: exited 137`, standard output is
exactly `roundfix: Linux Lane container failure: go list exited 137: the
process ended by signal 9; see <lane log>`, standard error is empty and the
exit is `3`.

### Invariants

```text
1. go test runs only after go list succeeded; a failed go list never reaches go test.
2. Exit 4 is produced only by the script's go list branch, and its class is decided by the marker's code: below 128 finding, 128 or more environment.
3. Every outcome other than exit 4 keeps its summary and blocker bytes.
4. Only a whole line that starts with the marker prefix and ends in decimal digits names the code; the last such line wins.
5. The container's arguments, mounts, environment and skips are unchanged.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Feature 1, Success Metric 1, Invariant 1 → the script branch (API Contract 1)
- Goal 2, Core Feature 2, Success Metrics 2-5, Invariants 2 and 4 → the classification and the surfaces (API Contracts 2 and 4)
- Goal 3, Success Metric 5, Invariants 3 and 5 → unchanged outcomes (API Contract 3)
- Goal 4, Core Feature 3, Success Metric 6 → documentation (API Contract 5)
- Success Metric 7 → the PRD's Unreachable Acceptance

## Absence Evidence

1. The `go list` marker and its cause — command: `grep -rn -e 'go list: exited' -e 'go list exited' internal docs .agents skills`; result: no match at `42f523c2`.
2. A classification of exit `4` — command: `grep -n 'code == 4' internal/cli/deliver_linux_lane.go`; result: no match at `42f523c2`; exit `4` falls to the `go test exited <n>` default.

## ADR Alternatives

- ADR-0316: no rejected alternative is taken; the script does not exit with `go list`'s own code, exit `3` is not reused, every `go list` failure is not parked in `environment`, and no error text is parsed.
- ADR-0294: its classification is kept for every exit other than `4`; the rejected "every non-zero exit is the change's failure" is not reintroduced, since a `go list` signal stays a Lane Container Failure.
- ADR-0288: the lane still runs only the changed packages in a never-pull, network-free container; no rejected alternative (a contract test shelling out to Docker, pulling an image, mounting the checkout) is taken.
- ADR-0297: the history mount and the alternates line stay byte for byte; nothing is fetched or bundled.
- ADR-0303: every skip and its visibility stay; no skip becomes a failure.

## Integration Points

The lane script inside the container and the outcome text the gate, the
blocker, the lane log and `roundfix deliver lane` carry. No test,
Verification command or QA row starts a container, pulls an image or reaches
the network.

## Testing Approach

Every new top-level test calls `t.Parallel()` first, compiles for Windows,
uses temporary repositories, and passes no real-time duration.

- `internal/cli/deliver_linux_lane_go_list_test.go` (new):
  `TestLinuxLaneScriptNamesAFailedGoList` cuts `linuxLaneScript` at its first
  `pkgs="$(go list` (failing when it is absent) and runs that tail with
  `sh -ceu <tail> linux-lane ./...`, with `PATH` holding only a stand-in `go`
  built by `testfixture.FixtureBinary`. The stand-in, on `list`, writes its
  configured standard output and exits its configured code; on `test`, it
  records that it ran and exits its configured code. Rows: `go list` exits
  `1` (script exit `4`, the line `roundfix-linux-lane: go list: exited 1` on
  its own line, `go test` never ran); exits `137` (exit `4`, `exited 137`,
  `go test` never ran); succeeds with one package and `go test` exits `1`
  (exit `1`, `roundfix-linux-lane: packages: 1`, no `go list:` line, `go
  test` ran); succeeds with no package (exit `0`, `no package builds on
  linux`, `go test` never ran).
- `TestLinuxLaneGoListFailureIsNamed` (new, in
  `internal/cli/deliver_linux_lane_go_list_test.go`), through
  `commandDeliveryWorkflow.Gate` with the existing stand-in Docker
  (`LANE_EXIT=4`, `LANE_OUTPUT` the marker): code `1` gives the reason, the
  `gate-failed` blocker and the lane log ending of Success Metric 2; code
  `137` gives the reason and `linux-lane-container` blocker of Success
  Metric 3.
- `TestDeliverLaneGoListFailureTranscript` and
  `TestDeliverLaneGoListSignalTranscript` (new, in
  `internal/cli/deliver_lane_test.go`) assert, with
  `assertDeliverLaneTranscript`, the two whole outputs that Surface
  Transcripts describes.
- `TestLinuxLaneOutcomeClassification` (changed): the row `other change
  exit` moves from `4` to `5` with summary `Linux Lane change failure: go test
  exited 5`, and the table gains the rows `go list change` (`exited 1`),
  `go list below signals` (`exited 127`, a change failure), `go list signal
  zero` (`exited 128`, signal `0`), `go list killed` (`exited 137`, signal
  `9`), `go list last marker` (`exited 7` then `exited 1`, cause `go list
  exited 1`), `go list no marker`, `go list embedded marker ignored` (`x
  roundfix-linux-lane: go list: exited 1`), `go list signed code ignored`
  (`exited +1`), `go list spaced code ignored` (`exited 1 `) and `go list
  carriage return` (`exited 1\r\n`, cause `go list exited 1`); the embedded,
  signed and spaced forms give the no-marker cause.
- `TestLinuxLaneRunsTheChangedPackagesInAContainer` (changed): its
  independent script literal gains the branch of API Contract 1.
- `TestGateParksWithTheBlockerItNames`
  (`internal/delivery/gate_blocker_test.go`, changed): rows `go list change`
  (`gate-failed: Linux Lane change failure: go list exited 1; see lane.log`,
  `finding`) and `go list signal` (`linux-lane-container: go list exited 137:
  the process ended by signal 9; see lane.log`, `environment`).
- Unchanged and run by name: `TestLinuxLaneFailureFailsTheGate`,
  `TestLinuxLaneContainerFailureIsToldApartFromAChangeFailure`,
  `TestLinuxLaneReportsThePackagesTheContainerSelected`,
  `TestDeliverLaneChangeFailureTranscript`,
  `TestDeliverLaneContainerFailureTranscript`.

## Build Order

1. The `deliver` guide, the Roundfix `deliver` reference, the Roundfix
   skill's versions, the Behavior Surface Record and **Lane Change Failure**
   (no dependency: the text states the contracts of `_techspec.md`, and its
   Verification reads only documents).
2. The script branch, the classification and their tests (depends on: 1,
   which writes the documents for the surface this step changes, so the code
   and its text land in one chain).
3. QA gate (depends on: 1, 2).

No Spec is listed in `requires`. Specs 0299 and 0300 and the earlier cycle AP
Specs deliver first and raise the Roundfix skill's version; this Spec rebases
on them and touches none of their files apart from the skill's version
fields.

## Risks & Considerations

- A module the read-only cache lacks still fails as `go test exited 1`,
  because `go list -e` reports it as an erroneous package and exits `0`; this
  Spec does not reclassify `go test` failures.
- `/bin/sh` in `golang:1.26` is dash; the tail was characterized under both
  dash and bash in POSIX mode, and the test runs it under the host's `sh`.
  The real container run is the record's unproven step (Success Metric 7).
- A `go list` that prints a line starting with the marker prefix before it
  fails could only be overridden by the script's own line, which comes last.

## Glossary

- changes: **Lane Change Failure**

## Decisions

- A failed `go list` exits `4` with a marker line carrying its code; see
  ADR-0316.
- Its Park Class is `finding` below `128` and `environment` for a signal; see
  ADR-0316.
- The gate, `deliver lane` and the Park Class mapping stay unchanged, because
  they follow the classifier's blocker.
