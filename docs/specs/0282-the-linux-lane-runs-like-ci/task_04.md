---
task: task_04
spec: 0282-the-linux-lane-runs-like-ci
status: pending
type: backend
complexity: high
---

# Task 04: `roundfix deliver lane` runs the Linux Lane on demand and exits with its outcome

## Overview

The adopted Backlog Entry ends with "After that, this repository turns the
lane back on" (`_prd.md` → Sources 1), and the operator must prove a lane run
passes before that change. Today the lane runs only inside a Delivery Queue
gate, only when `delivery.linux_lane` is on at the fetched default. This Task
adds `roundfix deliver lane [--all] [<revision>]`, which runs the same lane
against the checkout and exits with its outcome. It is verifiable on its own
through four Surface Transcript tests and an argument test with the compiled
Docker and Go stand-ins.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-18,
   `references/2026-10-18-the-linux-lane-container-lacks-ci-tools-and-an-init.md`
   (`_prd.md` → Sources 1), for turning the lane back on, by implementing API
   Contract 6 of `_techspec.md` (ADR-0294).
2. MUST complete `linuxLaneRequest` with `Label`, `LogPath`, `Name` and `All`
   (`_techspec.md` → Interfaces) and make the gate build it with today's
   label `Delivery Queue item <slug>: `, log path and container name, so every
   gate line, log and name stays byte-identical. `All` selects `./...` and
   reads no merge base. The `sorted nearest unique packages` subtest of
   `TestLinuxLaneRunsTheChangedPackagesInAContainer` follows the new call and
   keeps its assertions.
3. MUST implement `deliver lane` in the new `internal/cli/deliver_lane.go`,
   dispatched from `internal/cli/deliver.go`, with API Contract 6's
   arguments, configuration source, log path, container name, output lines
   and exit codes.
   - The data source is the checkout's own Git objects and its Project
     Config over User Config; it never reads `delivery.linux_lane`, the
     Delivery Queue, the Run Database or a Spec, and writes only the lane log.
   - Its readers are the operator and scripts: outcome lines go to standard
     output and preflight failures to standard error, with no other line.
   - An interrupt removes the container, as the gate's cancellation does,
     before exiting `130`.
4. MUST change the help per API Contract 6: `deliverUsage` gains the lane
   usage, command and flag lines, and the root usage line in
   `internal/cli/cli.go` gains `lane`. This is a declared break of
   `internal/cli/deliver_retry_test.go` (line 367), which moves to the new
   root line and keeps its name.
5. MUST add, in the new `internal/cli/deliver_lane_test.go`, each test calling
   `t.Parallel()` first, compiling for Windows and ending every child it
   starts in `t.Cleanup`:
   - `TestDeliverLanePassedTranscript`, `TestDeliverLaneChangeFailureTranscript`,
     `TestDeliverLaneContainerFailureTranscript` and
     `TestDeliverLaneSkipTranscript`, each asserting the whole Surface
     Transcript of the same number in `_techspec.md`: every standard output
     line, the empty standard error and the exit code, with `<lane log>` the
     only controlled variation;
   - `TestDeliverLaneSelectsPackagesAndImage`: with `delivery.linux_lane:
     false` in Project Config, `--all` hands `./...` to `docker run`, a
     revision argument archives that commit, the checkout's
     `delivery.linux_lane_image` is the image, `--init` is present, the log
     is `<artifact dir>/delivery/linux-lane.log`, and an extra argument or a
     revision that is not a commit exits `2` with the `deliver lane`
     Preflight failed layout and starts no Docker.
6. MUST document the command in `docs/user-guide/commands/deliver.md` and
   `.agents/skills/roundfix/references/deliver.md`: its synopsis, `--all`, the
   revision, the configuration source, the output lines, the exit codes, and
   that the operator runs it before turning `delivery.linux_lane` on. Both
   MUST name `roundfix deliver lane`. No text inside any `### QA settlement`
   section may change.
7. MUST add to `docs/references/skill-coverage.json` a `command deliver lane`
   entry shaped like the other `deliver` entries, covered by
   `.agents/skills/roundfix/references/deliver.md`, whose sources add
   `internal/cli/deliver_lane.go` and `internal/cli/deliver_linux_lane.go`.
8. MUST revise **Linux Lane** in `CONTEXT.md` through domain-modeling: it also
   runs on demand with `roundfix deliver lane`.
9. MUST run these commands, because this Task edits a user guide page, the
   Roundfix reference and the command help. No record or mirror is
   hand-edited.
   Before them, raise both version fields of each edited owned skill's
   canonical `SKILL.md` (`metadata.version` and the top-level `version`)
   by one patch over the value the tree holds when this Task starts,
   because the record command refuses changed content under an unchanged
   version. The Agent runs these commands itself before handing back: the
   Daemon only runs the Verification and never regenerates a mirror or a
   record.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
     to record the raised Roundfix skill version;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     to re-record the Behavior Surface Record, including the new surface and
     the moved `deliver` fingerprints.

## Subtasks

- [ ] Complete the lane request and keep the gate's lines byte-identical.
- [ ] Implement and dispatch `deliver lane`.
- [ ] Change the help and the pinned root line.
- [ ] Add the transcript and argument tests.
- [ ] Document the command, add its coverage entry, revise the glossary, sync and record.

## Acceptance Criteria

- [ ] `roundfix deliver lane` reproduces Surface Transcripts 1-4 with exit `0`, `1`, `3` and `3`, whatever `delivery.linux_lane` says.
- [ ] The gate's lines, log and container name are unchanged by the refactor.
- [ ] The help, the guide, the reference, the Skill Coverage Map and the glossary name `roundfix deliver lane`.

## Consumers

- `deliver lane` help, output and exit codes, and the root usage line — consumers: `./internal/cli`, `./internal/docscontract`

## Context

- instruction: `docs/adr/0294-the-linux-lane-runs-like-ci.md`
- instruction: `internal/cli/deliver_plan.go`
- interface: `internal/cli/deliver.go`
- interface: `internal/cli/cli.go`
- creates: `internal/cli/deliver_lane.go`
- creates: `internal/cli/deliver_lane_test.go`
- interface: `internal/cli/deliver_linux_lane.go`
- interface: `internal/cli/deliver_linux_lane_test.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver_retry_test.go`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `docs/references/skill-coverage.json`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestDeliverLanePassedTranscript|TestDeliverLaneChangeFailureTranscript|TestDeliverLaneContainerFailureTranscript|TestDeliverLaneSkipTranscript|TestDeliverLaneSelectsPackagesAndImage|TestLinuxLaneRunsTheChangedPackagesInAContainer|TestLinuxLaneRunsTheConfiguredImage|TestLinuxLaneContainerFailureIsToldApartFromAChangeFailure|TestLinuxLaneCancellationRemovesTheContainer)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestDeliverLanePassedTranscript TestDeliverLaneChangeFailureTranscript TestDeliverLaneContainerFailureTranscript TestDeliverLaneSkipTranscript TestDeliverLaneSelectsPackagesAndImage TestLinuxLaneRunsTheChangedPackagesInAContainer; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; help="$(go run -buildvcs=false ./cmd/roundfix deliver --help 2>&1)" || { printf '%s\n' "$help"; exit 1; }; printf '%s\n' "$help" | grep -qF -- 'roundfix deliver lane [--all] [<revision>]' || { printf 'deliver help lacks the lane\n' >&2; exit 1; }; help="$(go run -buildvcs=false ./cmd/roundfix --help 2>&1)" || { printf '%s\n' "$help"; exit 1; }; printf '%s\n' "$help" | grep -qF -- 'roundfix deliver <plan|start|status|resume|retry|stop|lane>' || { printf 'root help lacks the lane\n' >&2; exit 1; }; for f in docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'roundfix deliver lane' || { printf 'missing roundfix deliver lane in %s\n' "$f" >&2; exit 1; }; done; grep -qF -- '"command deliver lane"' docs/references/skill-coverage.json || { printf 'missing coverage entry\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Linux Lane**:' || exit 1; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'roundfix deliver lane' || { printf 'Linux Lane glossary entry not revised\n' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -run 'Deliver|Retry|Help|Gate|LinuxLane' ./internal/cli || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the five `deliver lane` tests do not exist, so the command fails at its first pass check. Afterwards the transcripts and the argument test pass, the gate's lane tests still pass, both helps name the lane, the guide, the reference, the coverage map and the glossary name the command, and the skill version and the Behavior Surface Record are current.

## References

- `_prd.md` → Sources; Goals; User Story 5; Core Feature 4; Core Feature 5; Success Metric 4; Glossary
- `_techspec.md` → API Contract 6; API Contract 7; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; Interfaces; Current behavior; Testing Approach; Build Order 4
- ADR-0294; ADR-0288; ADR-0125; ADR-0259; ADR-0285; ADR-0189; ADR-0256; ADR-0280
