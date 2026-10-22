---
task: task_01
spec: 0297-a-corrective-task-after-a-completed-gate-passes-the-wave-collision-check
status: completed
type: backend
complexity: medium
---

# Task 01: Every Wave collision message names the edge that resolves it

## Overview

The `SC-WAVE-COLLISION` finding, whose code and summary are what a QA Report
records when the gate's precondition refuses, names no remedy, while the
Daemon's pre-dispatch refusal names the `needs` edge. This Task builds both
texts from one helper pair on `spec.WaveCollision`, makes the finding name the
edge, and pins oraculum's corrective shape: the gate precondition, the strict
check, the settlement check and the pre-dispatch refusal agree before and
after the edge.

## Requirements

1. MUST answer the Secondbrain Inbox Entry of 2026-10-21,
   [`references/2026-10-21-sc-wave-collision-no-portao-com-needs-e-check-estrito-limpo.md`](references/2026-10-21-sc-wave-collision-no-portao-com-needs-e-check-estrito-limpo.md)
   (`_prd.md` → Sources 1), by implementing API Contracts 1-5 of
   `_techspec.md`.
2. MUST add `func (collision WaveCollision) SharedPathList() string` and
   `func (collision WaveCollision) ResolvingEdge() string` to
   `internal/spec/collision.go` exactly as API Contracts 1-2 state.
3. MUST make `detectWaveCollisions` in `internal/speccheck/coherence.go`
   write the summary `Task <First> and Task <Second> can run in the same Wave
   and share repository files: <SharedPathList>; <ResolvingEdge>` and the fix
   ``Give Task <Second> a `needs` edge on Task <First> in <manifest>.``,
   keeping the code, severity and locations (API Contract 3).
4. MUST make `refuseTaskPlanWaveCollisions` in
   `internal/daemon/task_engine.go` build each pair's text from
   `SharedPathList` and `ResolvingEdge`, byte for byte as today (API
   Contract 4).
5. MUST NOT change which pairs collide: `Collisions`,
   `priorRunTaskTouches`, `taskDependencyClosure`, `declaredTaskTouches`,
   `spec.Load`, `speccheck.Check`, `specGatePrecondition`,
   `SpecCheckSettlementChecker`, the Precondition Refusal report and the
   serial-capacity exemption keep their behavior (API Contract 5). No
   existing test, golden or fixture changes.
6. MUST state the hazard in the tests: the evidence comes from the
   repository's own files and from local Git history on every ref; the
   tests build both in temporary repositories with the `gittest` helpers and
   never read the real `~/.roundfix`, another repository or the network.
7. MUST add the tests `_techspec.md` → Testing Approach lists, each top-level
   test calling `t.Parallel()` first, compiling for Windows and waiting on no
   real time:
   - in the new file `internal/spec/collision_edge_test.go`,
     `TestWaveCollisionNamesTheEdgeThatOrdersThePair`, whose two cases list
     the manifest as `task_01, task_02, task_05` and as `task_01, task_05,
     task_02` and assert the two exact edges, `SharedPathList`, no collision
     after the named edge and a `topologicalOrder` with no cycle; and
     `TestWaveCollisionSharedPathListSortsPathsWithTheirSources`;
   - in the new file `internal/speccheck/wave_collision_edge_test.go`,
     `TestWaveCollisionFindingNamesTheResolvingEdge`, asserting the exact
     summary, fix and `RefusalReason`;
   - in the new file `internal/daemon/qa_gate_wave_collision_test.go`,
     `TestQAGateJudgesACorrectiveTaskOnTheCurrentGraph`, with the two cases
     and the exact reason of `_prd.md` → Success Metric 1, the gate run
     through `qaGatePrecondition` in a worktree added with `git worktree add
     --detach`, the strict check through `speccheck.Check`, `PromoteGaps` and
     `GatePrecondition` in the item checkout, the settlement check through
     `SpecCheckSettlementChecker.RefusingFindings`, and the pre-dispatch
     refusal through `refuseTaskPlanWaveCollisions` at `Concurrency: 2`.

## Subtasks

- [ ] Add `SharedPathList` and `ResolvingEdge`.
- [ ] Build the finding's summary and fix, and the Daemon's refusal, from them.
- [ ] Add the edge, finding and corrective-shape tests.

## Acceptance Criteria

- [ ] The `SC-WAVE-COLLISION` summary and fix name the edge, and the refusal reason a QA Report records carries it.
- [ ] The Daemon's pre-dispatch refusal keeps its bytes.
- [ ] In oraculum's shape the gate, the strict check, the settlement check and the pre-dispatch refusal all refuse before the edge and none refuses after it.
- [ ] The existing Wave collision tests pass unchanged.

## Consumers

- `spec.WaveCollision.SharedPathList` and `ResolvingEdge`, the `SC-WAVE-COLLISION` finding text and the Daemon's pre-dispatch refusal text — consumers: `./internal/spec`, `./internal/speccheck`, `./internal/daemon`

## Context

- instruction: `internal/spec/spec.go`
- instruction: `internal/spec/collision_test.go`
- instruction: `internal/speccheck/mechanical.go`
- instruction: `internal/speccheck/citations.go`
- instruction: `internal/speccheck/coherence_test.go`
- instruction: `internal/speccheck/citation_characterization_test.go`
- instruction: `internal/speccheck/citations_test.go`
- instruction: `internal/speccheck/wave_instruction_test.go`
- instruction: `internal/speccheck/wave_collision_command_files_test.go`
- instruction: `internal/daemon/settlement_checks.go`
- instruction: `internal/daemon/task_engine_test.go`
- instruction: `internal/gittest/gittest.go`
- interface: `internal/spec/collision.go`
- interface: `internal/speccheck/coherence.go`
- interface: `internal/daemon/task_engine.go`
- creates: `internal/spec/collision_edge_test.go`
- creates: `internal/speccheck/wave_collision_edge_test.go`
- creates: `internal/daemon/qa_gate_wave_collision_test.go`

## Verification

- `out="$(go test -count=1 -v -run '^(TestWaveCollisionNamesTheEdgeThatOrdersThePair|TestWaveCollisionSharedPathListSortsPathsWithTheirSources|TestCollisionsExcludesTransitivelyOrderedTasks|TestCollisionsReadsSettlementsOnRunBranchesNotOnlyHead)$' ./internal/spec 2>&1 && go test -count=1 -v -run '^(TestWaveCollisionFindingNamesTheResolvingEdge|TestWaveCollisionIsReportedAtAuthoring|TestWaveCollisionCheckStillReportsASharedInterfacePath|TestWaveCollisionRefusesTheSameCommandFile)$' ./internal/speccheck 2>&1 && go test -count=1 -v -run '^(TestQAGateJudgesACorrectiveTaskOnTheCurrentGraph|TestTaskCycleRefusesAWaveCollision|TestTaskCycleDispatchesSerializedCollisionGraph)$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestWaveCollisionNamesTheEdgeThatOrdersThePair TestWaveCollisionSharedPathListSortsPathsWithTheirSources TestCollisionsExcludesTransitivelyOrderedTasks TestCollisionsReadsSettlementsOnRunBranchesNotOnlyHead TestWaveCollisionFindingNamesTheResolvingEdge TestWaveCollisionIsReportedAtAuthoring TestWaveCollisionCheckStillReportsASharedInterfacePath TestWaveCollisionRefusesTheSameCommandFile TestQAGateJudgesACorrectiveTaskOnTheCurrentGraph TestTaskCycleRefusesAWaveCollision TestTaskCycleDispatchesSerializedCollisionGraph; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; GOOS=windows go vet ./internal/spec ./internal/speccheck ./internal/daemon` — expected: exit 0. Before this Task the four new tests do not exist and `SharedPathList` is undefined, so the command fails at its first test run. Afterwards the edge, finding and corrective-shape tests pass beside the unchanged Wave collision tests, and the three packages compile for Windows.

## References

- `_prd.md` → Sources; Goals 1-3; Core Features 1-3; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4
- `_techspec.md` → Current behavior; API Contracts 1-5; Invariants 1-5; Testing Approach; Build Order 1
- ADR-0117
- ADR-0182


## Result

Implemented API Contracts 1–5 within this Task's slice. `WaveCollision`
now renders the sorted path/source list and the resolving edge. The Spec
Consistency finding uses those methods, names both Tasks in its fix, and
carries the edge into the durable refusal reason. The Daemon uses the same
methods while preserving its refusal text byte for byte. Collision discovery,
loading, gate/settlement classification and the serial-capacity exemption
are unchanged.

Added the four required top-level tests in the three declared new test files.
Each starts with `t.Parallel()`. Both manifest orders are loaded from disk,
and the named edge is written into the manifest before reloading and checking
that no collision or cycle remains. The corrective-shape fixture builds a
plan, trailer-bearing settlements on `roundfix/run-run_first`, fast-forward
integration, an archive commit and revert, the corrective declaration and the
reopened gate. The ordered case also leaves a corrective settlement on
`roundfix/run-run_second`. The gate reads a detached linked worktree; the
strict and settlement checks read the item checkout. All repositories and
refs are disposable and use `gittest`; the tests document the evidence hazard
and access no real Roundfix Home, other repository or network, and perform no
real-time waits.

| Acceptance criterion | Implementation and focused evidence |
| --- | --- |
| Finding summary, fix and recorded reason name the edge | `TestWaveCollisionFindingNamesTheResolvingEdge` asserts exact summary, fix and `RefusalReason`, plus unchanged severity and locations. The corrective-shape test asserts the exact `PreconditionRefusal.Reason`. |
| Daemon refusal keeps its bytes | The corrective-shape test compares the entire pre-dispatch error to the existing sentence, including the merged sources and edge, at `Concurrency: 2`. |
| Four readers agree before and after the edge | Both cases of `TestQAGateJudgesACorrectiveTaskOnTheCurrentGraph` pass. Gate, strict check and settlement return only the exact expected collision before ordering and no findings afterward; dispatch refuses with the exact text before ordering and returns nil afterward. |
| Existing collision tests remain unchanged | The focused collision run below passes the existing `TestCollisions*`, `TestWaveCollision*`, `TestTaskCycleRefusesAWaveCollision` and `TestTaskCycleDispatchesSerializedCollisionGraph` tests alongside the new tests. No existing test, golden or fixture file was edited. |

Focused checks:

- Red starting point: `rtk proxy go test ./internal/spec -run '^TestWaveCollision' -count=1`
  exited 1 because `SharedPathList` and `ResolvingEdge` were undefined.
- Final focused regressions: `rtk proxy go test ./internal/spec ./internal/speccheck ./internal/daemon -run '(Collisions|WaveCollision|SerializedCollisionGraph|QAGateJudges)' -count=1`
  exited 0 for all three packages. During fixture development, the loader
  exposed the missing QA dependency/reopening, and the stronger sole-refusal
  assertion exposed the corrective Task's missing Consumer Declaration;
  those fixture declarations were corrected without changing production
  validation or filtering away diagnostics.
- `GOOS=windows rtk proxy go test -c -o /private/tmp/roundfix-0297-spec.test.exe ./internal/spec`
  exited 0.
- `GOOS=windows rtk proxy go test -c -o /private/tmp/roundfix-0297-speccheck.test.exe ./internal/speccheck`
  exited 0.
- `GOOS=windows rtk proxy go test -c -o /private/tmp/roundfix-0297-daemon.test.exe ./internal/daemon`
  exited 0 after the final fixture edit.
- `rtk proxy git diff --check` exited 0.

Environment: the configured writable root contains a symlink, so ordinary
sandbox execution refused before starting a process. Authorized elevated
execution reached the resolved Run checkout; this is an environment condition.
The pre-existing diff was the Daemon's `status: in_progress` change in this
Task file, and that field was preserved.

The declared Task Verification, Windows vet and repository Verification were
not run; they remain Daemon-owned. No Task Graph manifest, other Task file,
existing test, golden, fixture, documentation/skill surface from task_02 or
configuration was changed. No commit, push or pull request was made.
