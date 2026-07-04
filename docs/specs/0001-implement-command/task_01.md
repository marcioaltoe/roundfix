---
task: task_01
spec: 0001-implement-command
status: pending
type: backend
complexity: high
---

# Task 01: Build the Spec contract parser

## Overview

Create the package that owns the on-disk Spec contract end to end: discovering active Specs, loading and validating a Task Graph into deterministic topological order, parsing task files, rewriting task status, and reading the QA Report verdict. Nothing outside this package will touch spec markdown. Verifiable on its own through unit tests over literal markdown fixtures.

## Requirements

1. MUST discover active Specs: directories under the repository's `docs/specs/` (excluding `_archived/`) whose PRD frontmatter carries `status: active`.
2. MUST load a Spec's Task Graph from the `_tasks.md` frontmatter (`schema: spec-tasks/v1`, `graph.nodes[]` with `id`, `file`, `needs`); the frontmatter is the authoritative source, never the projection table.
3. MUST validate the graph: every referenced task file exists, every `needs` entry names a known node id, and the graph is acyclic — each failure surfaces as a typed error naming the offending Task or check so the CLI can render one actionable Preflight Validation message.
4. MUST produce a deterministic topological order, breaking ties by manifest node order.
5. MUST parse each task file: frontmatter (`task`, `spec`, `status`, `type`, `complexity`), the title from the `# Task NN: <title>` heading, and the Verification commands extracted verbatim from the backticked bullet entries; a task with no parseable Verification command is a validation error.
6. MUST rewrite only the `status` frontmatter value when settling a Task, preserving the rest of the file exactly, and support re-reading a task file after an Agent has modified it.
7. MUST read the verdict from the newest `qa/qa-report-*.md` frontmatter, reporting `pass`, `fail`, `partial`, missing, and unreadable as distinct results.
8. MUST introduce no new module dependencies; reuse the YAML and frontmatter-splitting approach already used for Round artifacts.

## Subtasks

- [ ] Active-Spec discovery from PRD frontmatter and directory layout
- [ ] Task Graph manifest parsing and schema validation
- [ ] Graph validation (missing files, unknown needs, cycles) with typed errors
- [ ] Deterministic topological ordering
- [ ] Task file parsing: frontmatter, title, Verification commands
- [ ] Status rewrite preserving file content
- [ ] QA Report verdict reader

## Acceptance Criteria

- [ ] Table tests cover: valid graph, cycle, unknown `needs` id, missing task file, wrong `schema`, inactive Spec, missing Spec, and a task file without Verification commands — each returning its typed error with the offending Task or check in the message.
- [ ] The same manifest always yields the same topological order (asserted by running the ordering repeatedly).
- [ ] A status rewrite changes only the `status` line: the rewritten file matches the original byte-for-byte outside that value.
- [ ] Verdict reader table tests cover `pass`, `fail`, `partial`, no report present, and unreadable frontmatter, plus newest-report selection when several reports exist.
- [ ] Task parsing extracts multi-command Verification sections in section order, verbatim.

## Verification

- `rtk go test ./internal/spec/` — expected: all tests pass.
- `rtk go build ./...` — expected: builds cleanly.

## References

`_prd.md` → User Stories 1, 4, 7; Core Features 2, 12. `_techspec.md` → System Architecture (internal/spec), Interfaces (internal/spec), Build Order 1. `docs/agents/issue-tracker.md` → tracker conventions.
