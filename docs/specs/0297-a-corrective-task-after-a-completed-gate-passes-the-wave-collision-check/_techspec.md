---
spec: 0297-a-corrective-task-after-a-completed-gate-passes-the-wave-collision-check
prd: _prd.md
created: 2026-10-22
---

# A corrective Task after a completed gate passes the Wave collision check — Technical Spec

## Executive Summary

The QA gate's precondition and `roundfix spec check --strict` call the same
`speccheck.Check` on the same working-tree Task Graph, and both read prior-Run
settlement evidence with `git log --all`; characterization on oraculum's own
history shows the gate refused its corrective Task once, before the `needs`
edge existed, and passed that check on every later head. The defect is the
message: the `SC-WAVE-COLLISION` finding, whose code and summary are what a QA
Report records as its `precondition_reason`, names no remedy, while the
Daemon's pre-dispatch refusal names the edge. Two methods on
`spec.WaveCollision`, `SharedPathList` and `ResolvingEdge`, now build both
texts, the finding's summary ends with the edge, and its fix line names the
same two Tasks. A Daemon-package test pins oraculum's corrective shape across
the gate precondition, the strict check, the settlement check and the
pre-dispatch refusal, before and after the edge.

The trade-off accepted: the summary carries a remedy, which diagnostic style
guides keep in a separate help line, because the QA Report records only the
code and the summary and the Daemon's refusal already reads that way.

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug, record or schema
  is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint,
  credential or recipient is added; every test uses temporary repositories
  and Git configuration. Source: `docs/agents/cli.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — Build Orders 1 and 2 follow ADR-0117
  and ADR-0182, and leave ADR-0096's mechanical stage unchanged, as
  `_prd.md` → Project Constraints names them with their titles. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AO approval of 2026-10-22 ("the same rules as cycle AN"), carrying the
  cycle AJ grant "Mexer em Governed Paths declarados". Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0297-a-corrective-task-after-a-completed-gate-passes-the-wave-collision-check/_authorization.md`;
  bounded files: `internal/speccheck/coherence.go`,
  `.agents/skills/roundfix/references/spec.md`,
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`.

## Current behavior

Read and characterized on 2026-10-22 at `34b4e246` (v0.107.0). No commit
between oraculum's binary `956066a7` and `34b4e246` changes
`internal/spec/collision.go`, `internal/speccheck/coherence.go`,
`internal/speccheck/citations.go`, `internal/speccheck/mechanical.go`,
`internal/spec/spec.go`, `internal/daemon/task_engine.go`,
`internal/daemon/settlement_checks.go` or `internal/daemon/qa_prior_pass.go`.
Line numbers are at `34b4e246`; the function names do not move.

1. The gate path. `Engine.runQAGate` builds its mechanical request through
   `qaMechanicalRequest`, which calls `qaGatePrecondition` and
   `specGatePrecondition` (`internal/daemon/task_engine.go:3262`, `:3266`):
   `speccheck.Check(plan.SpecsRoot, plan.WorkDir, slug)`, `PromoteGaps` and
   `GatePrecondition`. `plan.SpecsRoot` and `plan.WorkDir` are the Run
   Worktree's. A blocking precondition is written as the Precondition Refusal
   report, whose `precondition_reason` is `speccheck.RefusalReason` of each
   finding: the code, `: ` and the summary. The fix line is not recorded.
2. The strict check path. `runSpecCheckCommand` calls the same
   `speccheck.Check(resolvedSpecsRoot.Path, loaded.GitRoot, slug)`, then
   `PromoteGaps` under `--strict`, and `classifySpecCheckBoundary`
   (`internal/cli/spec_check.go:578`) applies the same `GatePrecondition`.
   The settlement checks call `SpecCheckSettlementChecker.RefusingFindings`
   (`internal/daemon/settlement_checks.go:22`), which is
   `specGatePrecondition` again.
3. The judgement. `speccheck.Check` loads the graph with
   `loadOptionalTaskGraph` → `spec.Load` from the working tree's `_tasks.md`
   (`internal/speccheck/citations.go:1641`) and calls `detectWaveCollisions`
   (`internal/speccheck/coherence.go:304`), which calls `spec.Collisions`
   (`internal/spec/collision.go:45`). `Collisions` skips every pair the
   dependency closure orders, so a pair with a `needs` path between them is
   never reported, whatever the evidence. The evidence is the Verification
   operands, the declared Context paths the tree holds, and
   `priorRunTaskTouches` (`:367`): the paths of each Task's newest commit
   carrying `Roundfix-Spec`/`Roundfix-Task` trailers on any ref
   (`git log --all`). No reader on this path opens an Archive Record, a
   `source_revision` or another commit's `_tasks.md`.
4. The two texts. `detectWaveCollisions` writes the summary `Task <First> and
   Task <Second> can run in the same Wave and share repository files: <path>
   (<source>), …` and the fix `Give one Task a `needs` edge on the other in
   <manifest>.`. `refuseTaskPlanWaveCollisions` (`task_engine.go:591`), which
   runs only at a Task Capacity above 1, writes `Task <First> and Task
   <Second> share repository files: <path> (<source>), …; add a `needs` edge
   from Task <Second> to Task <First>`. Each builds the path list with its own
   copy of the same loop.
5. Oraculum's Spec 0087, read without writing to oraculum through a mirror
   clone in the operator's scratch space. The item branch
   `roundfix/deliver-0087-o-rateio-repartido-pelo-percentual-f0d19f08d4a6668e`
   holds, in order: the three Task settlements of the first Run, the first QA
   Report (`partial`), the archive commit `04b985a` and its revert `677d816`,
   `57a7c55` adding `task_05` with `needs: [task_01]`, `0b48d9e` changing it
   to `needs: [task_02]`, `task_05`'s settlement and `75a9dd0` declaring the
   shared test under `task_05`. The QA commit `6552c7c` carries four reports.
   Only `qa-report-2026-10-21-01.md` is a Precondition Refusal, with the exact
   reason the Inbox Entry quotes and no `build:` key. Reports `-02`, `-03` and
   `-04` record `build:` `8736693`, `7472ad34` and `75a9dd0`, each with
   `task_05 needs: [task_02]`; each reached its matrix, wrote "Strict
   precondition … exited 0", and failed on row R13 for two documentation
   findings. Run times match: `run_20261021T055620Z` started at `57a7c55`'s
   graph, and the three `T06*` Runs at the later heads. `findRecordedPriorQAPass`
   (`internal/daemon/qa_prior_pass.go:75`) copies the newest QA Report commit
   the Run's history lacks into the next pass, so `-01` sat in every later Run
   Branch beside the newer report. `bin/roundfix spec check
   0087-o-rateio-repartido-pelo-percentual --strict` built at `34b4e246` on
   worktrees of the mirror reports the collision on `57a7c55` and nothing on
   `8736693`, `0b48d9e` and `75a9dd0`.
6. Reproduction. A temporary repository with a plan commit, a Run Branch
   `roundfix/run-run_A` settling `task_01` to `task_03` (`task_02` writing
   `pkg/versions.test.ts` without declaring it), a fast-forward, an archive
   commit and its revert, and a corrective `task_05` declaring that file under
   `interface:` with `needs: [task_01]`: `roundfix spec check 0087-shape
   --format json` reports `Task task_02 and Task task_05 can run in the same
   Wave and share repository files: pkg/versions.test.ts (declared context;
   prior Run settlement commit)`, oraculum's text. An overlay probe in
   `internal/daemon` that wrote nothing to the tree ran `specGatePrecondition`
   (the same finding) and `refuseTaskPlanWaveCollisions` at capacity 2 (`…;
   add a `needs` edge from Task task_05 to Task task_02`). With `needs:
   [task_02]` committed and a second Run Branch settling `task_05`, all three
   report no collision.
7. Contrast. Spec 0296's corrective `task_04` (`needs: [task_02]`, sharing
   `internal/cli/deliver_queue_subject_test.go`) passed its gate on
   2026-10-21; at its `source_revision` `3dc0b8dd` the strict check reports
   no finding. Both repositories took the same path; oraculum's later gates
   passed it too.

Sweep of every place that writes or reads the collision text, at `34b4e246`
(`grep -rn -e 'can run in the same Wave' -e 'share repository files' -e
'needs. edge' internal docs .agents`): the two writers above;
`TestTaskCycleRefusesAWaveCollision` (`internal/daemon/task_engine_test.go`),
which checks substrings of the Daemon text including `` `needs` edge ``;
`TestWaveCollisionIsReportedAtAuthoring` (`internal/speccheck/coherence_test.go`),
which checks substrings of the summary and that the fix contains `needs`. No
guide, skill or golden quotes either text; the `spec` user guide and the
Roundfix `spec` reference do not mention `SC-WAVE-COLLISION`.

## System Architecture

- `internal/spec/collision.go`: `WaveCollision.SharedPathList` and
  `WaveCollision.ResolvingEdge`.
- `internal/speccheck/coherence.go`: `detectWaveCollisions` builds its
  summary and fix from them.
- `internal/daemon/task_engine.go`: `refuseTaskPlanWaveCollisions` builds its
  text from them, byte for byte as today.
- `docs/user-guide/commands/spec.md` and the Roundfix `spec` reference: the
  finding.

No package, import edge or Context path changes; `internal/daemon` and
`internal/speccheck` already import `internal/spec`.

## Implementation Design

### Interfaces

```go
// SharedPathList renders the shared paths as "<path> (<source>)", sorted and
// joined with ", ".
func (collision WaveCollision) SharedPathList() string

// ResolvingEdge names the one needs edge that orders the pair: Second needs
// First.
func (collision WaveCollision) ResolvingEdge() string
```

### Data Models

None. `WaveCollision`, `speccheck.Finding`, the QA Report and its
Precondition Refusal keep their shapes.

### API Contracts

1. API Contract: `SharedPathList`. Each path of `Paths` as `<path>
   (<source>)`, sorted by that string, joined with `, `; the merged source
   form `declared context; prior Run settlement commit` is kept as
   `mergedTouchSource` builds it.
2. API Contract: `ResolvingEdge`. Exactly ``add a `needs` edge from Task
   <Second> to Task <First>``. `First` precedes `Second` in the graph's Task
   order and neither reaches the other in the dependency closure, so adding
   the edge orders the pair and cannot form a cycle.
3. API Contract: the finding. `detectWaveCollisions` writes the summary `Task
   <First> and Task <Second> can run in the same Wave and share repository
   files: <SharedPathList>; <ResolvingEdge>` and the fix ``Give Task <Second>
   a `needs` edge on Task <First> in <manifest>.``; code, severity and
   locations are unchanged. `RefusalReason` therefore records the edge in
   `precondition_reason`, in the settlement diagnostics and in the
   settlement baseline.
4. API Contract: the Daemon refusal. `refuseTaskPlanWaveCollisions` writes
   `Task <First> and Task <Second> share repository files: <SharedPathList>;
   <ResolvingEdge>` per pair, joined with `; ` after `task cycle: refuse Wave
   before dispatch because Tasks collide: `, byte for byte as today.
5. API Contract: unchanged judgement. `Collisions`, `priorRunTaskTouches`,
   `taskDependencyClosure`, `spec.Load`, `speccheck.Check`,
   `specGatePrecondition`, `SpecCheckSettlementChecker` and the
   serial-capacity exemption keep their bytes; a test proves the four readers
   agree in oraculum's shape.
6. API Contract: documentation. The `spec` user guide gains an
   `SC-WAVE-COLLISION` row and the Roundfix `spec` reference a bullet, each
   naming the evidence sources, the resolving edge and that the gate, the
   settlement checks and the pre-dispatch refusal read the same working-tree
   graph and settlement history.

### Surface Transcripts

None. The finding reaches an operator through `roundfix spec check`, whose
flags, exit codes and line layout do not change; only one finding's sentence
gains a clause, and the tests of API Contracts 2-4 assert the exact strings.

### Invariants

```text
1. Every SC-WAVE-COLLISION summary and fix, and every Daemon pre-dispatch refusal, names the edge ResolvingEdge returns for its pair.
2. Adding that edge removes the pair's collision and forms no cycle.
3. The gate precondition, the strict check, the settlement check and the pre-dispatch refusal report the same Wave collisions for the same tree and refs.
4. A pair the current graph lets run in one Wave that shares a file is still reported by all four.
5. The Daemon's refusal text keeps its bytes.
```

## Coverage Map

Deviations from the PRD: None.

- Goal 1, Core Features 1-2, Success Metric 3, Invariants 1-2 and 5 → `SharedPathList`, `ResolvingEdge`, `detectWaveCollisions`, `refuseTaskPlanWaveCollisions` (API Contracts 1-4)
- Goals 2-3, Core Feature 3, Success Metrics 1-2 and 4, Invariants 3-4 → the corrective-shape test and the unchanged judgement (API Contract 5)
- Goal 4, Core Feature 4 → the `spec` guide and the Roundfix `spec` reference (API Contract 6)
- Success Metric 5 → the PRD's Unreachable Acceptance

## Absence Evidence

1. `SharedPathList`, `ResolvingEdge` — command: `grep -rn -e SharedPathList -e ResolvingEdge internal`; result: no match at `34b4e246`; `detectWaveCollisions` and `refuseTaskPlanWaveCollisions` each render the path list in their own loop.
2. A remedy naming the edge in the Spec Consistency Check finding — command: `grep -n 'needs. edge' internal/speccheck/coherence.go`; result: only the fix line "Give one Task a `needs` edge on the other", which names neither Task.
3. A reader that loads the Task Graph or settlement evidence from an archived revision on the gate path — command: `grep -n -e source_revision -e ReadArchivedSpec -e 'git show' internal/spec/collision.go internal/speccheck/coherence.go internal/daemon/settlement_checks.go`; result: no match; the graph comes from `spec.Load` on the working tree and the evidence from `git log --all`.

## ADR Alternatives

- ADR-0117: no rejected alternative is taken; the check stays at authoring and at the gate as one judgement, and the gate is not made to read a different graph or evidence than the authoring check.
- ADR-0182: no rejected alternative is taken; the settlement checks keep reading the facts the gate reads, now with the same edge in the reason.
- ADR-0096: no rejected alternative is taken; the mechanical stage and its refusal record are untouched, and only the text of one finding it records gains its remedy.

## Integration Points

Git, read by `priorRunTaskTouches` with `git log --all` and `git diff-tree`,
unchanged. No test, Verification command or QA row reads the real
`~/.roundfix`, another repository, the network or a provider.

## Testing Approach

Every new top-level test calls `t.Parallel()` first, compiles for Windows,
uses temporary repositories and the shared `gittest` helpers, and waits on no
real time.

- `internal/spec/collision_edge_test.go` (new):
  - `TestWaveCollisionNamesTheEdgeThatOrdersThePair`: `task_02` and `task_05`
    both need `task_01` and declare `pkg/versions_test.go` under
    `interface:`; with the manifest order `task_01, task_02, task_05` the
    edge is ``add a `needs` edge from Task task_05 to Task task_02``, with
    `task_01, task_05, task_02` it is ``… from Task task_02 to Task
    task_05``; `SharedPathList` is `pkg/versions_test.go (declared
    context)`; after adding the named edge `Collisions` returns none and
    `topologicalOrder` accepts the nodes.
  - `TestWaveCollisionSharedPathListSortsPathsWithTheirSources`: two paths,
    one with the merged source, in sorted order.
- `internal/speccheck/wave_collision_edge_test.go` (new):
  `TestWaveCollisionFindingNamesTheResolvingEdge`: a three-Task fixture read
  through `speccheck.Check`; the summary, the fix and `RefusalReason` equal
  the exact strings of API Contract 3.
- `internal/daemon/qa_gate_wave_collision_test.go` (new):
  `TestQAGateJudgesACorrectiveTaskOnTheCurrentGraph`, two cases over
  oraculum's shape built with real Git (plan commit, a Run Branch
  `roundfix/run-run_first` with trailer-carrying settlements of `task_01` to
  `task_03`, `task_02` writing `pkg/versions.test.ts` undeclared, a
  fast-forward, an archive commit and `git revert`, and `task_05` declaring
  that file under `interface:`), the gate run in a linked worktree created
  with `git worktree add --detach`:
  - `needs: [task_01]`: `qaGatePrecondition` on the linked worktree,
    `speccheck.Check` + `PromoteGaps` + `GatePrecondition` on the item
    checkout and `SpecCheckSettlementChecker.RefusingFindings` each give
    exactly one Wave collision reason, equal to Success Metric 1's text; the
    Precondition Refusal reason carries it; `refuseTaskPlanWaveCollisions`
    at capacity 2 carries `pkg/versions.test.ts (declared context; prior Run
    settlement commit); add a `needs` edge from Task task_05 to Task
    task_02`.
  - `needs: [task_02]` with a second Run Branch `roundfix/run-run_second`
    settling `task_05` on the file: none reports a Wave collision and the
    pre-dispatch refusal returns nil.
- Unchanged and run by name: `TestTaskCycleRefusesAWaveCollision`,
  `TestTaskCycleDispatchesSerializedCollisionGraph`,
  `TestWaveCollisionIsReportedAtAuthoring`,
  `TestWaveCollisionCheckStillReportsASharedInterfacePath`,
  `TestWaveCollisionRefusesTheSameCommandFile`,
  `TestCollisionsExcludesTransitivelyOrderedTasks` and
  `TestCollisionsReadsSettlementsOnRunBranchesNotOnlyHead`.
- No existing test, golden or fixture changes.

Each sabotage the QA gate applies was run during authoring on a scratch copy
with the three new tests, and failed the named test: reversing
`ResolvingEdge` fails all three; dropping the edge from the finding's summary
fails the `speccheck` and Daemon tests; ignoring prior-Run settlement
evidence in `Collisions` fails the Daemon test, whose first case needs the
`prior Run settlement commit` source; returning no collisions fails all
three.

## Build Order

1. `SharedPathList`, `ResolvingEdge`, their use in `detectWaveCollisions` and
   `refuseTaskPlanWaveCollisions`, and the three new tests (no dependency).
2. The `spec` user guide, the Roundfix `spec` reference, the skill version
   and the Behavior Surface Record (depends on: 1, whose finding text the
   documents state).
3. QA gate (depends on: 1, 2).

No Spec is listed in `requires`. Spec 0297 delivers first of cycle AO; Specs
0298 and 0299 change `internal/delivery` and the Delivery Queue, and Spec
0300 `reconcile`, none of the files above.

## Risks & Considerations

- A QA Report written before this Spec keeps its old reason text; a Repeated
  Failure signature across the binary change may differ once.
- The precondition refusal report records no audited head, so a refusal
  imported beside a newer report can still be read as current; this Spec
  leaves that report shape unchanged and the `spec` guide's existing rule
  stands: only the newest report counts.
- A pair the closure leaves unordered for a reason the author intends (two
  Tasks meant to run together) still has only the remedy the rule allows:
  order them or stop sharing the file.

## Glossary

None.

## Decisions

- The gate and the strict check already judge one graph and one evidence;
  tests pin it in oraculum's shape and the collision path keeps its bytes.
- One helper pair builds both collision texts; the summary carries the edge
  because the QA Report records only code and summary.
- The edge makes the later Task in graph order need the earlier one, as the
  Daemon text already does.
- No new ADR.
