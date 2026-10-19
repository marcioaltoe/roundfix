---
task: task_03
spec: 0282-the-linux-lane-runs-like-ci
status: completed
type: backend
complexity: high
---

# Task 03: A Lane Container Failure is told apart from a Lane Change Failure

## Overview

On 2026-10-18 the operator saw only `gate-failed` for an item whose changed
packages all passed: the gate's reason never reached the blocker, and Docker's
own exit codes read as the change's failure (`_prd.md` → Sources 1). This Task
makes the lane script report its own setup failures, classifies every run
exit as a **Lane Container Failure** or a **Lane Change Failure**, carries the
outcome into the queue's blocker and Park Class, and ends the lane log with
it. It is verifiable on its own through the lane's compiled Docker stand-in
and engine tests with a fake gate.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-18,
   `references/2026-10-18-the-linux-lane-container-lacks-ci-tools-and-an-init.md`
   (`_prd.md` → Sources 1), for "A lane failure caused by the container is
   then told apart from a failure of the change", by implementing API
   Contract 2's script, API Contract 3 and API Contract 4 of `_techspec.md`
   (ADR-0294).
2. MUST replace `linuxLaneScript` with API Contract 2's script, byte for
   byte. This is a declared break of
   `TestLinuxLaneRunsTheChangedPackagesInAContainer`, whose independent
   script copy becomes the new text and keeps its name.
3. MUST classify a run that started per API Contract 3, through one function
   the tests can call with an exit code, whether the cap was exceeded and the
   run's output. The cause comes from the exit code or the last
   `roundfix-linux-lane: container: ` line of the container's own output;
   the readers are the operator, through the blocker and the log, so the
   cause names a tool, a step, an exit code or a signal and never echoes
   other output.
   - The reasons, blockers and the lane log's last line MUST be API Contract
     3's exact texts.
   - The cap's cancellation and the container's removal MUST keep their
     behavior; the cap's reason moves to the container failure text.
   - This is a declared break of `TestLinuxLaneFailureFailsTheGate`: exits
     `1` and `2` are Lane Change Failures with the new text, and `126` and
     `127` leave it for the new test. It keeps its name.
4. MUST add `Blocker` to `delivery.GateResult` and
   `BlockerLinuxLaneContainer = "linux-lane-container"` to
   `internal/delivery/engine.go`; make `gateCandidate` park `result.Blocker`
   when it is not empty and the bare `gate-failed` otherwise; make `Gate` in
   `internal/cli/deliver_workflow.go` set the blocker for both lane failures
   and leave it empty for a repository Verification failure.
5. MUST make `ClassifyPark` put `linux-lane-container` in the `environment`
   class with API Contract 4's next action. `gate-failed` MUST keep its class
   and its generic next action, and every other blocker its class and text.
6. MUST add these tests, each calling `t.Parallel()` first and compiling for
   Windows:
   - in `internal/cli/deliver_linux_lane_test.go`,
     `TestLinuxLaneContainerFailureIsToldApartFromAChangeFailure`: stand-in
     exits `3` with the marker line, `3` without it, `126`, `127` and `137`
     each fail the gate with the container reason and the
     `linux-lane-container: ` blocker, and the lane log ends with `Linux Lane
     container failure: <cause>`; and `TestLinuxLaneOutcomeClassification`,
     a table over the classification function that includes the cap and exit
     `130`;
   - in the new `internal/delivery/gate_blocker_test.go`,
     `TestGateParksWithTheBlockerItNames`: a fake gate's `Blocker` is the
     parked blocker, an empty one parks `gate-failed`, and `ClassifyPark` gives
     the parked `linux-lane-container` item the `environment` class;
   - in `internal/delivery/park_class_test.go`, a `linux-lane-container` row
     with Park Class `environment` and its exact next action in
     `TestParkClassesAndNextActions`.

   Every existing lane, gate, park and retry test MUST keep its name and
   pass.
7. MUST document the change in `docs/user-guide/commands/deliver.md` and
   `.agents/skills/roundfix/references/deliver.md`:
   - the "Linux Lane" section names both outcomes with API Contract 3's
     causes, reasons and blockers, the setup exit `3` and the log's last
     line;
   - the Park Class table's `environment` row gains `linux-lane-container`.

   MUST search the guides, the agent guides, the skill copies and
   `internal/` for the old texts `Linux Lane failed: docker run exited` and
   `Linux Lane exceeded 20m0s` and replace each. `CHANGELOG.md` and the ADRs
   stay as they are. No text inside any `### QA settlement` section may
   change.
8. MUST add **Lane Container Failure** and **Lane Change Failure** to
   `CONTEXT.md` through domain-modeling, with ADR-0294's link.
9. MUST run these commands, because this Task edits a user guide page and the
   Roundfix reference. No record or mirror is hand-edited.
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
     to re-record the Behavior Surface Record.

## Subtasks

- [ ] Replace the lane script and its pinned copy.
- [ ] Classify each exit and write the reasons, blockers and log line.
- [ ] Carry the blocker through the gate and classify the new blocker.
- [ ] Add the lane, classification, gate and park tests.
- [ ] Document both outcomes, add the glossary terms, sync and record.

## Acceptance Criteria

- [ ] A run that fails before `go test`, or whose exit Docker or a signal owns, parks `linux-lane-container` in `environment`.
- [ ] A run whose `go test` fails parks `gate-failed` with its reason in `finding`, and a repository Verification failure still parks the bare `gate-failed`.
- [ ] The lane log, the guides, the reference and the glossary name both outcomes, and no copy keeps the old failure texts.

## Consumers

- `delivery.GateResult.Blocker`, `delivery.BlockerLinuxLaneContainer` and its Park Class — consumers: `./internal/delivery`, `./internal/cli`

## Context

- instruction: `docs/adr/0294-the-linux-lane-runs-like-ci.md`
- instruction: `docs/adr/0288-the-delivery-queue-gate-runs-the-changed-go-packages-in-a-linux-lane.md`
- interface: `internal/cli/deliver_linux_lane.go`
- interface: `internal/cli/deliver_linux_lane_test.go`
- interface: `internal/cli/deliver_workflow.go`
- interface: `internal/delivery/engine.go`
- interface: `internal/delivery/park_class.go`
- interface: `internal/delivery/park_class_test.go`
- creates: `internal/delivery/gate_blocker_test.go`
- interface: `docs/user-guide/commands/deliver.md`
- interface: `.agents/skills/roundfix/references/deliver.md`
- interface: `skills/roundfix/references/deliver.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestLinuxLaneContainerFailureIsToldApartFromAChangeFailure|TestLinuxLaneOutcomeClassification|TestLinuxLaneFailureFailsTheGate|TestLinuxLaneRunsTheChangedPackagesInAContainer|TestLinuxLaneCancellationRemovesTheContainer|TestLinuxLaneSkipsWithANamedReason)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestLinuxLaneContainerFailureIsToldApartFromAChangeFailure TestLinuxLaneOutcomeClassification TestLinuxLaneFailureFailsTheGate TestLinuxLaneRunsTheChangedPackagesInAContainer; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; out="$(go test -count=1 -v ./internal/delivery 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestGateParksWithTheBlockerItNames TestParkClassesAndNextActions TestEveryBlockerHasAParkClass; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; grep -qF -- 'linux-lane-container' internal/delivery/park_class_test.go || { printf 'park class row missing\n' >&2; exit 1; }; for f in docs/user-guide/commands/deliver.md .agents/skills/roundfix/references/deliver.md; do tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'linux-lane-container' || { printf 'missing linux-lane-container in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Lane Container Failure' || { printf 'missing Lane Container Failure in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Lane Change Failure' || { printf 'missing Lane Change Failure in %s\n' "$f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -qF -- 'Linux Lane change failure: go test exited' || { printf 'missing change failure reason in %s\n' "$f" >&2; exit 1; }; done; if grep -rqF -e 'Linux Lane failed: docker run exited' -e 'Linux Lane exceeded 20m0s' docs/user-guide docs/agents .agents/skills skills internal; then printf 'old Linux Lane failure text remains\n' >&2; exit 1; fi; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Lane Container Failure**:' || { printf 'missing glossary term Lane Container Failure\n' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Lane Change Failure**:' || { printf 'missing glossary term Lane Change Failure\n' >&2; exit 1; }; make skills-sync-check || exit 1; go test -count=1 ./skills || exit 1; go test -count=1 -run 'Deliver|Gate|Park|Retry|LinuxLane' ./internal/cli || exit 1; go test -count=1 -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the new tests do not exist, so the command fails at its first pass check. Afterwards the lane, classification, gate and park tests pass, the guides, the reference and the glossary name both outcomes, the old texts are gone, and the skill version and the Behavior Surface Record are current.

## References

- `_prd.md` → Sources; Goals; User Story 4; Core Feature 3; Core Feature 5; Success Metric 3; Glossary
- `_techspec.md` → API Contract 2; API Contract 3; API Contract 4; API Contract 7; Invariants 4-5; Interfaces; Current behavior; Testing Approach; Build Order 3
- ADR-0294; ADR-0288; ADR-0125; ADR-0259; ADR-0271; ADR-0285; ADR-0189; ADR-0256; ADR-0280

## Result

Implemented Task 03's slice; status and declared Verification remain Daemon-owned.
The only pre-existing worktree change was this Task file's `in_progress` status;
that field was preserved. No other Task file or Task Graph manifest was edited.

- Acceptance criterion 1: the production script and its independent pinned test
  copy now match API Contract 2 byte for byte. `classifyLinuxLaneOutcome` handles
  the cap, setup exit 3 (last marker or the named fallback), exits 126/127 and
  signal exits, including 130. The lane supplies the exact container reason and
  `linux-lane-container` blocker; the engine persists that blocker, and Park
  Class is `environment` with API Contract 4's exact next action. Compiled
  stand-in tests cover marker selection, no marker, output without a trailing
  newline, 126, 127 and 137; the cap is covered without a real-time wait.
- Acceptance criterion 2: exits 1 and 2 retain the existing test name and now
  assert the exact Lane Change Failure reason, detailed `gate-failed` blocker
  and final log line. The fake-gate engine test proves the detailed blocker
  survives parking in `finding`, while an empty blocker falls back to bare
  `gate-failed`. The repository-failure test explicitly asserts its blocker
  stays empty. Existing generic Park Class answers are exercised unchanged.
- Acceptance criterion 3: both user/reference guides describe the causes,
  reasons, blockers, setup exit 3 and terminal log lines; their `environment`
  rows include `linux-lane-container`. `CONTEXT.md` defines both failure terms
  with ADR-0294 links. Both Roundfix skill version fields rose from 0.1.104
  to 0.1.105 before generation. Mirrors and records were generated, not
  hand-edited. A source-tree scan found neither old failure text in the
  required guides, agent guides, skill trees or `internal/`. Inspection also
  confirmed the successful lane log ends in its exact pass line and all
  `### QA settlement` sections remain byte-identical.

Focused-check evidence:

- Red signal: `go test ./internal/cli ./internal/delivery -run
  'Test(LinuxLaneOutcomeClassification|GateParksWithTheBlockerItNames)$'
  -count=1` failed to compile the newly added tests before implementation,
  naming the absent blocker constant, `GateResult.Blocker` and classifier.
- `go test ./internal/cli ./internal/delivery -run
  'Test(LinuxLane|GateWithoutTheLane|GateSkipsTheLane|GateParksWithTheBlockerItNames|ParkClassesAndNextActions|EveryBlockerHasAParkClass|ExistingGenericParkAnswersAreUnchanged)'
  -count=1`: exit 0 after the final test edits; CLI 2.735s, delivery 0.465s.
  This includes the existing skip, configured-image and cancellation/removal
  tests as well as the new classification and blocker tests.
- `GOOS=windows go test -c -o /tmp/task03-cli-windows.test.exe
  ./internal/cli` and `GOOS=windows go test -c -o
  /tmp/task03-delivery-windows.test.exe ./internal/delivery`: both exit 0
  after the final edits; compile-only evidence, not Windows execution.
- `make skills-sync`: exit 0.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
  -record-skill-versions`: exit 0; recorded Roundfix 0.1.105.
- `go test -tags docscontract ./internal/docscontract -run
  '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage -count=1`:
  exit 0; regenerated the deliver guide's Behavior Surface Record.
- `make baseline-digests`: exit 0, `changed:false`; derived Baseline
  artifacts already matched their canonical sources.
- A local Python inspection compared both scripts to API Contract 2,
  scanned the five required trees for old failure texts, compared canonical
  and mirrored Roundfix files and compared QA settlement sections to HEAD:
  all assertions passed. `git diff --check`: exit 0.

Environment: ordinary shell execution could not start because the configured
writable root contains a symlink. Focused work used authorized escalation in
this same checkout's resolved `/Volumes/Brain-SSD/cache/.roundfix/worktrees/`
location. No Docker daemon, network, provider or real user state was used.

The Task's declared Verification and repository Verification were not run;
they remain pending for the Daemon. No commit, push or pull request was made.
No follow-up outside this Task's slice was introduced.

## Carry-forward provenance

- Source Run: `run_20261019T224006Z_3d81ced188078879`
- Source commit: `ef08a2ca4ea862f9e2944a98c593583ed2f1bcb4`
