---
task: task_01
spec: 0072-qa-is-a-task-not-a-flag
status: completed
type: backend
complexity: high
---

# Task 01: Teach the graph the gate and its invalidation

## Overview

The Task Graph gains the QA gate as a first-class node. A new Task Type `qa`
joins the canonical set; the manifest frontmatter declares the authoring
decision — `qa: task_NN` naming the gate node, or `qa: declined` with
`qa_reason` — and validation makes the gate structurally honest: it must be
unique, terminal, and depend on every leaf, and a settled gate above a
dependency that settled after it (or never settled) is a stale-gate error
naming the inserted Tasks. A graph with neither declaration is a legacy
graph and keeps working byte-for-byte unchanged.

## Requirements

1. MUST add `TaskTypeQA` (`"qa"`) to the canonical Task Type set in
   `internal/spec`, accepted in task frontmatter and in the manifest
   projection table.
2. MUST parse the manifest frontmatter declaration: `qa: <task-id>` or
   `qa: declined` plus a non-empty `qa_reason`; expose the result on the
   loaded `Graph`.
3. MUST validate, for a graph declaring a gate: the named node exists, has
   `type: qa`, is the only `qa` node, has no dependents, and depends
   (directly or transitively) on every other leaf node.
4. MUST reject: a `qa`-typed node in a graph with no `qa:` declaration; a
   declaration naming a node of another type; `qa: declined` with a gate
   node present or with an empty reason.
5. MUST fail loading with a stale-gate error when the gate node's status is
   settled but any of its dependencies is not `completed`; the error names
   the offending Task ids and says the gate result is invalidated.
6. MUST leave graphs without any declaration loading exactly as today,
   proven against the archived Specs' real manifests.
7. MUST keep dependencies owned only by the manifest and status owned only
   by task files; the declaration adds no second copy of either.

## Subtasks

- [ ] Add the type, the frontmatter parse, and the `Graph` fields.
- [ ] Implement terminal-coverage validation and its error vocabulary.
- [ ] Implement the stale-gate check with the inserted-Task naming.
- [ ] Add fixtures: authored gate, declined, legacy, non-terminal gate,
      appended-after-report.
- [ ] Characterize the archived Specs' manifests as legacy pass-through.

## Acceptance Criteria

- [ ] A fixture graph whose terminal `qa` node depends on every leaf loads,
      and `Graph` reports the gate's id.
- [ ] A fixture with `qa: declined` and a reason loads gateless; the same
      fixture without the reason fails naming `qa_reason`.
- [ ] Appending a Task to a fixture whose gate is `completed` makes the next
      load fail with the inserted Task named (PRD Success Metric 4).
- [ ] A gate node that is not terminal, or not covering a leaf, fails
      validation naming the uncovered node.
- [ ] Every manifest under `docs/specs/_archived/` still loads with
      unchanged results.
- [ ] `git status --porcelain` shows no path outside `internal/spec/` and
      this task file.

## Verification

- `go test ./internal/spec -count=1 -run 'QA|Gate' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the new contract tests run and pass.
- `go test ./internal/spec -count=1` — expected: exit 0; nothing regressed.
- `go build -buildvcs=false ./...` — expected: exit 0.

## References

- `_prd.md` → Core Features 1, 2, 4, 7; Success Metrics 1, 4, 5.
- `_techspec.md` → Implementation Design (Interfaces); Decisions; ADR-0091.

## Result

### Implementation

- Added `TaskTypeQA` to the closed Task Type set used by task frontmatter and
  manifest projection rows.
- Added `Graph.QATaskID`, `Graph.QADeclined`, and `Graph.QAReason`, with strict
  parsing for `qa: <task-id>` and `qa: declined` plus a non-empty
  `qa_reason`.
- Added typed QA gate validation for declaration/type consistency, uniqueness,
  terminal placement, leaf coverage, and declined-gate conflicts.
- Added `StaleGateError`, which invalidates a settled gate above any
  non-completed transitive dependency and names those Task ids in manifest
  order.
- Added graph fixtures for authored, declined, invalid-shape, and appended
  gate cases, plus a characterization that rehosts all 62 real archived
  manifests and their task files under an active temporary PRD without
  modifying the archived sources.

### Focused checks

- Red signal: `rtk env GOCACHE=/private/tmp/roundfix-task01-gocache go test ./internal/spec -run '^TestLoadQAGateContract$' -count=1` exited 1 before implementation because `Graph` had no QA declaration fields.
- `rtk env GOCACHE=/private/tmp/roundfix-task01-gocache go test ./internal/spec -run '^(TestTaskTypeCanonicalValuesLoadThroughTaskGraph|TestLoadQAGateContract|TestLoadQADeclinedContract|TestLoadRejectsInvalidQAGateShape|TestLoadInvalidatesSettledQAGateAfterTaskAppend|TestQAGateLegacyArchivedManifestsLoadUnchanged)$' -count=1 -v` exited 0; every named test and negative subtest reported `PASS`.
- `rtk env GOCACHE=/private/tmp/roundfix-task01-gocache go test ./internal/spec -run '^(TestTaskTypeRejectsInvalidFrontmatterValues|TestTaskTypeProjectionMustMatchTaskFile|TestTaskTypeProjectionTablePresenceValidatesRows|TestLoadReturnsTypedValidationErrors|TestLoadReturnsTasksInDeterministicTopologicalOrder)$' -count=1` exited 0 for the adjacent legacy parser/error contracts.
- `rtk git diff --check` exited 0.
- The Task's declared `## Verification` commands were not run; they remain
  Daemon-owned.

### Acceptance criteria evidence

1. `TestLoadQAGateContract` loads a terminal `qa` Task through its projection
   row and asserts `Graph.QATaskID == "task_03"`.
2. `TestLoadQADeclinedContract` loads a reasoned decline and rejects both a
   missing and an empty reason with an error naming `qa_reason`.
3. `TestLoadInvalidatesSettledQAGateAfterTaskAppend` first loads a completed
   gate, appends pending `task_04`, and then asserts `StaleGateError` names
   `task_04` and says the gate result is invalidated.
4. `TestLoadRejectsInvalidQAGateShape` rejects an uncovered `task_02`, a
   non-terminal gate with dependent `task_03`, an undeclared QA Task, a
   mismatched declared type, multiple QA Tasks, and a declined graph with a
   gate node.
5. `TestQAGateLegacyArchivedManifestsLoadUnchanged` enumerates all 62
   `_tasks.md` files under `docs/specs/_archived/`, loads each real manifest
   and its real task files through the legacy path, asserts zero QA declaration
   state, and compares the source manifest bytes before and after.
6. `rtk git -c core.fsmonitor=false status --porcelain` exited 0 and listed
   only `internal/spec/` files and this Task file.
