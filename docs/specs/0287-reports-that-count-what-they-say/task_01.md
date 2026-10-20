---
task: task_01
spec: 0287-reports-that-count-what-they-say
status: completed
type: backend
complexity: medium
---

# Task 01: The Linux Lane's pass line counts the packages that ran

## Overview

`roundfix deliver lane --all` printed `Linux Lane passed 1 package(s)` while
43 packages ran, because the pass line counts the package patterns the host
handed to the container (`_prd.md` → Sources 1). This Task makes the lane
script report how many packages `go list` selected and the host print that
number, in the Delivery Queue gate and in `roundfix deliver lane` alike. It is
verifiable on its own through the lane's compiled Docker stand-in.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-20,
   [`references/2026-10-20-the-linux-lane-counts-package-patterns-as-packages.md`](references/2026-10-20-the-linux-lane-counts-package-patterns-as-packages.md)
   (`_prd.md` → Sources 1), by implementing API Contracts 1 and 2 of
   `_techspec.md` in `internal/cli/deliver_linux_lane.go`, with
   `linuxLanePackageCount` as `_techspec.md` → Interfaces sketches it.
2. MUST change the lane script to API Contract 1's exact text, so it prints
   `roundfix-linux-lane: packages: <n>` to standard output before `exec go
   test`, and `roundfix-linux-lane: packages: 0` after `no package builds on
   linux`. Nothing else in the script changes.
3. MUST keep one builder: the count is read in `runLinuxLane`, so the gate and
   `roundfix deliver lane`, in changed-package and `--all` modes, print it
   alike; neither `deliver_workflow.go` nor `deliver_lane.go` changes.
4. MUST state the hazard in the code's behavior, as API Contract 2 does: the
   data source is the container's combined output, which the candidate's
   tests also write to, so only the first line that starts with
   `roundfix-linux-lane: packages: ` and holds only decimal digits after it is
   read; its destination is the console line and the lane log only. A pass
   without such a line prints `Linux Lane passed on <image>; the lane script
   reported no package count` and stays a pass. Failure lines, the outcome
   classification, blockers, exit codes, mounts, image and package selection
   do not change, and `result.Packages` keeps the patterns.
5. MUST declare these breaks and keep each test's name:
   - `linuxLaneToolSource` in `internal/cli/deliver_linux_lane_test.go`
     writes `roundfix-linux-lane: packages: ` followed by `LANE_PACKAGES` and
     a newline to standard output, before its fixture output, on `run` when
     `LANE_PACKAGES` is not empty; `newLinuxLaneFixture` adds
     `LANE_PACKAGES=1` to its environment, so the existing `passed 1
     package(s)` assertions keep their text;
   - `TestLinuxLaneRunsTheChangedPackagesInAContainer`'s independent script
     copy becomes API Contract 1's text.
6. MUST add these tests, each top-level test calling `t.Parallel()` first,
   compiling for Windows and ending every child it starts in `t.Cleanup`:
   - in `internal/cli/deliver_linux_lane_test.go`,
     `TestLinuxLaneReportsThePackagesTheContainerSelected`, through the gate
     on a candidate that changes `go.mod` (one pattern, `./...`), with
     subtests: `LANE_PACKAGES=43` gives the console line `roundfix: Delivery
     Queue item widget: Linux Lane passed 43 package(s) on golang:1.26.1` and
     a lane log ending with `Linux Lane passed 43 package(s) on
     golang:1.26.1`; `LANE_PACKAGES=0` gives `passed 0 package(s)`; an empty
     `LANE_PACKAGES` gives both no-count lines; `LANE_PACKAGES=43` with
     `LANE_OUTPUT` holding `roundfix-linux-lane: packages: 7` still gives 43;
     and the script constant holds `echo "roundfix-linux-lane: packages: $#"`
     before `exec go test`;
   - in `internal/cli/deliver_lane_test.go`,
     `TestDeliverLanePackageCountTranscript` and
     `TestDeliverLaneNoPackageCountTranscript`, asserting Surface Transcripts
     1 and 2 of `_techspec.md` whole: the one standard output line, the empty
     standard error and exit `0`.
7. MUST document the change in `docs/user-guide/commands/deliver.md` and
   `.agents/skills/roundfix/references/deliver.md`:
   - replace "`<n>` counts selected package arguments" with the statement that
     `<n>` is the number of packages `go list` selected in the container,
     which the lane script reports with `roundfix-linux-lane: packages: <n>`
     before `go test`;
   - add the no-count console and log lines, `Linux Lane passed on <image>;
     the lane script reported no package count`, beside each pass line the
     guide and the reference list.

   MUST search the guides, the agent guides and the skill copies for wording
   that says the pass line counts package arguments or patterns, and change
   each. No text inside any `### QA settlement` section may change.
8. MUST run these commands, because this Task edits a user guide page and the
   Roundfix reference. No record or mirror is hand-edited.
   Before them, raise both version fields of the Roundfix skill's canonical
   `SKILL.md` (`metadata.version` and the top-level `version`) by one patch
   over the value the tree holds when this Task starts, because the record
   command refuses changed content under an unchanged version. The Agent runs
   these commands itself before handing back: the Daemon only runs the
   Verification and never regenerates a mirror or a record.
   - `make skills-sync`;
   - `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
     to record the raised Roundfix skill version;
   - `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`,
     to re-record the Behavior Surface Record.

## Subtasks

- [ ] Add the count line to the lane script.
- [ ] Read the first count line and print the counted or the no-count pass lines.
- [ ] Teach the Docker stand-in the count line and update the pinned script copy.
- [ ] Add the count and transcript tests.
- [ ] Document the count in the guide and the reference; sync and record.

## Acceptance Criteria

- [ ] With the stand-in reporting 43, the gate and `roundfix deliver lane --all` print `Linux Lane passed 43 package(s)` although they handed the container one pattern.
- [ ] A count line a test prints later never replaces the script's.
- [ ] A pass without a count prints the no-count line, and `roundfix deliver lane` reproduces Surface Transcripts 1 and 2.
- [ ] The guide and the reference describe the count and the no-count line and no longer say the number counts package arguments.

## Consumers

- the Linux Lane's script and pass lines — consumers: `./internal/cli`, `./internal/docscontract`

## Context

- instruction: `docs/adr/0294-the-linux-lane-runs-like-ci.md`
- instruction: `docs/adr/0297-the-linux-lane-reads-the-repository-history-read-only.md`
- instruction: `internal/cli/deliver_lane.go`
- instruction: `internal/cli/deliver_workflow.go`
- interface: `internal/cli/deliver_linux_lane.go`
- interface: `internal/cli/deliver_linux_lane_test.go`
- interface: `internal/cli/deliver_lane_test.go`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -count=1 -v -run '^(TestLinuxLaneReportsThePackagesTheContainerSelected|TestDeliverLanePackageCountTranscript|TestDeliverLaneNoPackageCountTranscript|TestLinuxLaneRunsTheChangedPackagesInAContainer|TestDeliverLanePassedTranscript|TestDeliverLaneSelectsPackagesAndImage|TestLinuxLaneRunsTheConfiguredImage)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestLinuxLaneReportsThePackagesTheContainerSelected TestDeliverLanePackageCountTranscript TestDeliverLaneNoPackageCountTranscript TestLinuxLaneRunsTheChangedPackagesInAContainer TestDeliverLanePassedTranscript; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; grep -qF -- 'echo "roundfix-linux-lane: packages: $#"' internal/cli/deliver_linux_lane.go || { printf 'missing count line in the lane script\n' >&2; exit 1; }; for f in docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'roundfix-linux-lane: packages: <n>' || { printf 'missing count line in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'the lane script reported no package count' || { printf 'missing no-count line in %s\n' "$f" >&2; exit 1; }; if tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'counts selected package arguments'; then printf 'old count wording left in %s\n' "$f" >&2; exit 1; fi; done; GOOS=windows go build -buildvcs=false ./internal/cli || exit 1; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the three new tests do not exist, so the command fails at its first pass check. Afterwards the count and transcript tests pass with the existing lane and `deliver lane` tests, the script prints the count line, the guide and the reference describe the count and the no-count line without the old wording, the CLI builds for Windows, and the skill version and the Behavior Surface Record are current.

## References

- `_prd.md` → Sources; Goals; Core Feature 1; Core Feature 4; Success Metric 1
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 5; Surface Transcript 1; Surface Transcript 2; Invariants 1-2; Interfaces; Current behavior; Testing Approach; Build Order 1
- ADR-0294; ADR-0297; ADR-0288; ADR-0125; ADR-0259; ADR-0271; ADR-0285; ADR-0189; ADR-0256; ADR-0280


## Result

Implemented this Task's Linux Lane reporting slice; status and declared
Verification remain Daemon-owned.

### Implementation and acceptance evidence

- Container-selected count: the lane script matches TechSpec API Contract 1
  exactly, reporting zero on the empty selection and `$#` before `exec go
  test`. `runLinuxLane` builds both pass lines from `linuxLanePackageCount`;
  `result.Packages` still holds the host patterns. The gate regression changes
  `go.mod`, confirms the container receives only `./...`, and observes 43 in
  the exact console line and final lane-log line. The on-demand transcript
  regression observes 43 through `--all`, empty stderr and exit 0.
- First count wins: the gate regression observes 43 even when later candidate
  output reports 7. Additional cases cover zero, signed values, trailing
  spaces, an embedded marker, the first valid decimal line after malformed
  output, and an overflowing first decimal count. The code names the combined
  output's candidate-test hazard and limits the count's destination to the
  console and lane log.
- Missing count stays a pass: both gate lines use the specified no-count
  wording, and the on-demand test asserts the complete stdout line, empty
  stderr and exit 0. Existing lane failure, skip, selection, image, history
  mount and cancellation tests also pass in the focused run.
- Guide and reference: both explain the container's `go list` count and marker
  and include the no-count alternative beside all three pass-line lists.
  Searching `docs/user-guide`, `docs/agents`, `.agents/skills` and `skills`
  found no remaining wording that says the pass line counts package arguments
  or patterns. Both canonical Roundfix skill version fields rose from
  `0.1.113` to `0.1.114`; mirror and records were generated, not hand-edited.
  Inspection confirmed both `### QA settlement` sections are byte-identical
  to HEAD, as are `deliver_lane.go`, `deliver_workflow.go` and `_tasks.md`.

### Focused checks and generation

- Red signal before production edits:
  `go test ./internal/cli -run '^(TestLinuxLaneReportsThePackagesTheContainerSelected|TestDeliverLanePackageCountTranscript|TestDeliverLaneNoPackageCountTranscript)$' -count=1`
  exited 1: the old implementation printed 1 for the 43, zero and missing-count
  fixtures, and the script lacked its marker.
- Final focused run:
  `go test ./internal/cli -run '^(TestLinuxLane.*|TestDeliverLane.*)$' -count=1`
  exited 0 (`ok roundfix/internal/cli`, 8.595s). An earlier attempt overlapped
  skill generation and suiteguard refused those concurrent file writes after
  the assertions passed; the final run occurred after generation ended.
- `GOOS=windows go test -c -o <temporary-directory>/cli.test.exe ./internal/cli`
  exited 0; the temporary binary was removed.
- `make skills-sync` exited 0.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
  exited 0 and recorded `0.1.114`.
- `go test -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`
  exited 0 and regenerated the Behavior Surface Record.
- `make baseline-digests`, required by the repository after an owned-skill
  edit, exited 0 with `changed:false`; no Baseline artifact remains changed.
- Exact-script, mirror/version, QA-settlement and unchanged-caller inspections
  passed; `git diff --check` exited 0.

The Docker-refusal subtest retains its name and expected message; its fixture
now clears `LANE_PACKAGES`, since a container that never started cannot report
a script count. No production failure behavior changed. No other Task slice
was implemented; no declared Verification, repository Verification, commit,
push or Pull Request was run. Full settlement evidence is pending the Daemon.

## Carry-forward provenance

- Source Run: `run_20261020T144117Z_258265979932c5f1`
- Source commit: `ea87340bd96c307d7edf9f009d434361e476c053`
