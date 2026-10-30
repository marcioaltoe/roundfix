---
spec: 0334-the-template-tooling-row-passes-and-a-wave-collision-refuses-before-the-run
prd: _prd.md
created: 2026-10-30
---

# The template's tooling row passes the Spec Check, and a Wave collision is refused before the Run — Technical Spec

## Executive Summary

Two defects, each one condition in an existing function plus one call site.
In `internal/speccheck/constraints.go`, `declaresProtectedToolingMutation`
returns `true` for every row that cites an authorization record before it
reads the row's own words, so the templates' no-mutation row, which cites
its record, is a mutation, and the record's `paths: []` leaves no bounded
file. The fix reads the words first and accepts them when the cited record is
an operative grant with an empty `paths` list. In `internal/spec/collision.go`,
`Collisions` pairs every Task of the graph, completed ones included; the fix
skips a pair with a completed Task. `roundfix implement` gains one Preflight
Validation call to the Daemon's collision refusal, now exported, before the
Run is recorded.

The trade-off accepted: the collision rule runs twice for a graph that
passes, once in the checkout before the Run and once in the Run Worktree
before dispatch, rather than moving it, because the Daemon's refusal also
guards plans that reach the Task cycle without the Implement Command's
preflight (ADR-0347).

## Project Constraints

- Identifier strategy: not applicable — no identifier, slug or record schema
  is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no request path, endpoint or
  credential is added. Source: `docs/agents/go.md`,
  `docs/agents/agent-instructions.md`.
- Active ADR obligations: applicable — ADR-0347 records both decisions;
  ADR-0131 and ADR-0179 govern the tooling row, as `_prd.md` → Project Constraints states and ADR
  Alternatives below relates. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — express maintainer authorization: the cycle
  AS answers of 2026-10-30 quoted in the authorization record; the Governed
  Paths are the Spec Check's constraint reader and the owned skills
  `write-prd`, `write-techspec` and `roundfix`; no Makefile, workflow, module
  or lint configuration changes. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`;
  Spec-contained authorization record:
  `docs/specs/0334-the-template-tooling-row-passes-and-a-wave-collision-refuses-before-the-run/_authorization.md`;
  bounded files: `.agents/skills/roundfix/SKILL.md`, `.agents/skills/roundfix/references/implement.md`, `.agents/skills/roundfix/references/spec.md`, `.agents/skills/write-prd/SKILL.md`, `.agents/skills/write-prd/references/prd-template.md`, `.agents/skills/write-techspec/SKILL.md`, `.agents/skills/write-techspec/references/techspec-template.md`, `internal/speccheck/constraints.go`, `skills/roundfix/SKILL.md`, `skills/write-prd/SKILL.md`, `skills/write-prd/references/prd-template.md`, `skills/write-techspec/SKILL.md`, `skills/write-techspec/references/techspec-template.md`.

## Current behavior

Read on 2026-10-30 at `35de7f74`; `_prd.md` → Validation records the runs.

1. `detectToolingRow` runs only for an `applicable` row. It calls
   `declaresProtectedToolingMutation(row, len(selection.Candidates) != 0)`,
   which returns `true` whenever the row cites any authorization record and
   otherwise returns `false` when the reason says `no protected tooling
   mutation`. A mutation without `row.BoundedFiles` is `SC-TOOLING-UNBOUNDED`.
2. `recordsBoundedFiles` treats a bounded-files value of exactly `none` as
   absent and any other non-empty value as present.
3. Both templates prescribe the no-mutation reason, the record citation, the
   bounded files copied from `paths`, and `paths: []` for a Spec with no
   Governed Path. They name no value for an empty list.
4. `Collisions` builds a touch set for every Task and reports every pair the
   dependency closure leaves unordered that shares a path, whatever the
   Tasks' status. Its readers are `detectWaveCollisions` (Spec Check, QA gate
   precondition, settlement checks) and `refuseTaskPlanWaveCollisions`
   (Daemon, at a Task Capacity above 1).
5. `runImplementCommand` records the Run with `createRunReclaimingOrphan`
   before `executeImplementCycle` reaches the Daemon's refusal, so a
   colliding graph exits `1` with `implement failed after Run start: task
   cycle: refuse Wave before dispatch because Tasks collide: …`.

## System Architecture

No package, command, flag or file format is added. Changes by component:

- `internal/speccheck/constraints.go` (task_01): the row reading.
- `write-prd` and `write-techspec` templates (task_01): the row's ending.
- `internal/docscontract/command_documentation_test.go` (task_01): the
  contract that composes the row from each template.
- `internal/spec/collision.go` (task_02): completed Tasks leave the pairs.
- `internal/daemon/task_engine.go` (task_02): the exported preflight form of
  the refusal, sharing the detail text with the pre-dispatch refusal.
- `internal/cli/implement.go` (task_02): the preflight call.
- Guides, Roundfix references and `CONTEXT.md` (task_02).

## Implementation Design

### Interfaces

```go
// internal/speccheck/constraints.go
func declaresProtectedToolingMutation(row constraintRow, citesAuthorization, recordGrantsNoPath bool) bool
// recordGrantsNoPath: the row's selected record resolves as
// spec.AuthorizationGranted for the asking Spec and its Record.Paths is empty.

// internal/daemon/task_engine.go
func RefuseWaveCollisions(repoRoot string, graph *spec.Graph, taskCapacity int) error
func refuseTaskPlanWaveCollisions(plan TaskPlan) error // unchanged message
```

### Data Models

None. No record or file format changes.

### API Contracts

1. API Contract: the no-mutation row. A Tooling authority row whose reason,
   whitespace-normalized and lower-cased, contains `no protected tooling
   mutation` declares no protected tooling mutation when it cites no
   authorization record, or when the record selected for it resolves as an
   operative grant for the asking Spec with an empty `paths` list. Otherwise
   a row citing a record declares one, and a row citing none declares one
   when its reason contains `protected tooling mutation` or `express
   maintainer authorization`, as today. `SC-TOOLING-UNBOUNDED` keeps its code,
   summary and fix.
2. API Contract: the templates. The no-mutation comment of
   `prd-template.md` and `techspec-template.md` keeps the reason
   `applicable — no protected tooling mutation proposed or authorized` and
   adds that the row still cites the Spec-contained authorization record and,
   when its `paths` is `[]`, ends with `bounded files: none`.
3. API Contract: completed Tasks. `spec.Collisions` reports no pair in which
   either Task's status is `completed`. A Task of any other status, `pending`
   after a reopen, `in_progress` or `failed`, still pairs. The touch sources
   and the `ResolvingEdge` text are unchanged.
4. API Contract: the preflight refusal. `RefuseWaveCollisions` returns nil
   when `taskCapacity` is 1 or less, otherwise the error `Tasks would collide
   in one Wave at Task Capacity <n>: <details>`, where `<details>` is the
   text the pre-dispatch refusal prints after its colon today: for each
   pair, `Task <first> and Task <second> share repository files: <paths>;`
   followed by the pair's `ResolvingEdge` text, the pairs joined by `; `.
   `runImplementCommand` calls it with the checkout root, the committed graph
   and `worktree.concurrency` after the graph loads and before the default
   branch check, prints it with `printPreflightFailure` and exits `2`. The
   pre-dispatch refusal keeps its message bytes.

### Surface Transcripts

1. Surface Transcript: `roundfix implement` refuses two pending Tasks that
   share a Verification operand, at a Task Capacity of 2, before the Run.

   ```transcript
   $ roundfix implement --spec <slug> --no-input
   stdout:
   stderr:
   Preflight failed

   Reason:
     Tasks would collide in one Wave at Task Capacity 2: Task task_01 and Task task_02 share repository files: shared/notes.txt (verification command); add a `needs` edge from Task task_02 to Task task_01

   No side effects:
     Roundfix did not create a Run, fetch Review Source issues, start an Agent, commit, or push.

   Usage:
     Run 'roundfix implement --help' for usage.
   exit: 2
   ```

### Invariants

```text
1. A row that cites a record listing any Governed Path, or that declares a mutation, still needs bounded files.
2. A row that cites no record keeps today's verdict.
3. A completed Task never pairs; every other status pairs as today.
4. A colliding graph at a Task Capacity above 1 records no Run and starts no Agent session.
5. At a Task Capacity of 1 Preflight Validation refuses nothing.
6. The Daemon's pre-dispatch refusal message keeps its bytes.
```

### Declared breaks

- `internal/daemon/qa_gate_wave_collision_test.go`
  `TestQAGateJudgesACorrectiveTaskOnTheCurrentGraph`, case `unordered
  corrective Task`: `task_05 needs [task_01]` shares
  `pkg/versions.test.ts` with the completed `task_02`, and the case expects
  all four readers to refuse. Under API Contract 3 none refuses. The case is
  renamed `corrective Task beside a completed Task` and expects no finding
  from any reader, and a new case `two unordered corrective Tasks` adds a
  pending `task_06` (`needs [task_01]`) declaring the same file, so all four
  readers agree on the refusal of `task_05` and `task_06`. The top-level test
  keeps its name, so the coverage record does not move.
- No other test changes: the prototype ran `internal/spec`,
  `internal/speccheck`, `internal/daemon` and `internal/cli` with only this
  case failing.

## Coverage Map

- Goal 1, Core Feature 1, Success Metric 1 →
  `declaresProtectedToolingMutation` (API Contract 1).
- Goal 2, Success Metric 2 → the same, Invariants 1 and 2.
- Goal 3, Core Feature 2, Success Metric 3 → the templates and
  `TestTheTemplateNoMutationToolingRowPassesTheSpecCheck` (API Contract 2).
- Goal 4, Core Feature 3, Success Metrics 4 and 6 → `Collisions`
  (API Contract 3), Invariant 3.
- Goal 5, Core Feature 4, Success Metric 5 → `RefuseWaveCollisions` and the
  preflight call (API Contract 4, Surface Transcript 1), Invariants 4–6.
- Goal 6, Core Feature 5, Success Metric 7 → the guides, references and
  `CONTEXT.md` (task_02).

Deviations from the PRD: None.

## Absence Evidence

1. Reading a cited empty grant as no mutation — command: `grep -n 'citesAuthorization' internal/speccheck/constraints.go`; result: the only use returns `true` before the reason is read; no code reads the cited record's `paths` for this decision.
2. Completed Tasks in `Collisions` — command: `grep -n 'Status' internal/spec/collision.go`; result: no match; status is never read.
3. A collision check in the Implement Command — command: `grep -n 'Collisions\|WaveCollision' internal/cli/implement.go`; result: no match; only the Daemon's `refuseTaskPlanWaveCollisions` calls `spec.Collisions` at Run time.
4. The templates' value for an empty grant — command: `grep -n 'bounded files: none' .agents/skills/write-prd/references/prd-template.md .agents/skills/write-techspec/references/techspec-template.md`; result: no match.

## ADR Alternatives

- ADR-0131: followed; the checker changes rather than the template, as it decided, and the row's words decide again when the cited record grants nothing.
- ADR-0179: followed; an explicitly empty `paths` list bounds nothing, so a row that says "no mutation" agrees with it; an absent or null list is still refused by the reader before this rule runs.
- ADR-0347: implemented as recorded, including its rejection of changing the templates to `not applicable`, of ignoring the cited record, and of refusing only in the Spec Check or only in the Daemon.

## Integration Points

None external. The preflight reads the checkout's files and local Git
history the way `spec.Collisions` already does.

## Testing Approach

All tests are hermetic: temporary repositories, the implement test fixtures
with the fake Agent runner, and no network.

- task_01, `internal/speccheck/report_test.go`:
  `TestANoMutationRowCitingAnEmptyGrantIsNotUnbounded`, three cases through
  `writeToolingRowFixture` with a Spec-contained record: the no-mutation row
  on `paths: []` (no finding), the same row on `paths` listing `Makefile`
  (`SC-TOOLING-UNBOUNDED`; this case catches a fix that ignores the record),
  and `applicable — protected tooling mutation proposed and authorized` on
  `paths: []` (`SC-TOOLING-UNBOUNDED`).
  `internal/docscontract/command_documentation_test.go` (a `//verify:always`
  file): `TestTheTemplateNoMutationToolingRowPassesTheSpecCheck` extracts
  the reason and the `bounded files: none` ending from each template,
  composes the row in a temporary repository and requires
  `speccheck.CheckStage` at the PRD and TechSpec stages to give no
  `SC-TOOLING-UNBOUNDED` on `paths: []` and to give it on `paths` listing
  `Makefile`.
- task_02, `internal/spec/collision_test.go`:
  `TestACompletedTaskNeverCollides`, one graph per status of `task_01`
  (`completed`, `pending`, `failed`, `in_progress`) beside a pending
  `task_02` with the same Verification operand: zero pairs for `completed`,
  one otherwise (the `failed` and `in_progress` cases catch a rule that
  skips every Task that is not `pending`).
  `internal/daemon/qa_gate_wave_collision_test.go`: the declared break.
  `internal/cli/implement_wave_collision_test.go` (new):
  `TestImplementRefusesAWaveCollisionBeforeTheRun` (Surface Transcript 1:
  empty stdout, the whole stderr, exit `2`, no Run in the store, no Agent
  call), `TestImplementRunsATaskThatSharesAFileWithACompletedTask` (the
  Quaestor shape: the Run is created and stderr has no collision text) and
  `TestImplementAtTaskCapacityOneLeavesCollidingTasksToTheRun` (the Run is
  created).
- The suites that read the changed contracts run in the same Verification
  commands: `internal/speccheck` tooling and constraint tests, the spec
  collision tests, the daemon collision and gate tests, and the CLI
  implement preflight tests.

## Build Order

1. The row reading, the template ending and the contract test (task_01).
2. Completed Tasks leave the pairs, the preflight refusal, the guides,
   references and glossary (task_02) (depends on: 1 — both raise owned-skill
   versions in `skills/testdata/owned-skill-versions.json`).
3. QA gate (task_03) (depends on: 1, 2).

Spec 0333 delivers first and may raise the `write-prd`, `write-techspec` or
Roundfix skill versions; each Task here edits after that merge and raises
the version through the record command, never by hand.

## Risks & Considerations

- The preflight reads the checkout's working tree for declared paths, while
  the Daemon reads the Run Worktree created from `HEAD`. An untracked file
  in the checkout at a path a Task declares under `creates:` could make the
  preflight see a touch the Run would not; the refusal then names the path
  and the edge, which is safe, and the Daemon's check is unchanged.
- `recordGrantsNoPath` reads the cited record a second time; records are
  small and the reader is already called for the same row.
- A row whose bounded-files value is `none` followed by a parenthesized note
  keeps counting as bounded (Non-Goals); a row that declares a mutation with
  such a value is still caught by `SC-TOOLING-UNDECLARED` when a Task
  declares a Governed Path.

## Glossary

- adds: **Wave Collision**

## Decisions

- The checker reads the row's words first and accepts "no mutation" with an
  empty cited grant. See ADR-0347.
- The templates add only the row's ending, so the contract test can compose
  it from them.
- A completed Task never pairs; a reopened one does. See ADR-0347.
- The Implement Command refuses in Preflight Validation at a Task Capacity
  above 1, through an exported Daemon function that shares the detail text;
  the Daemon keeps its refusal. See ADR-0347.
