---
task: task_01
spec: 0072-qa-is-a-task-not-a-flag
status: pending
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
