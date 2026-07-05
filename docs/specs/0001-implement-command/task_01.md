---
task: task_01
spec: 0001-implement-command
status: completed
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

- [x] Active-Spec discovery from PRD frontmatter and directory layout
- [x] Task Graph manifest parsing and schema validation
- [x] Graph validation (missing files, unknown needs, cycles) with typed errors
- [x] Deterministic topological ordering
- [x] Task file parsing: frontmatter, title, Verification commands
- [x] Status rewrite preserving file content
- [x] QA Report verdict reader

## Acceptance Criteria

- [x] Table tests cover: valid graph, cycle, unknown `needs` id, missing task file, wrong `schema`, inactive Spec, missing Spec, and a task file without Verification commands — each returning its typed error with the offending Task or check in the message.
- [x] The same manifest always yields the same topological order (asserted by running the ordering repeatedly).
- [x] A status rewrite changes only the `status` line: the rewritten file matches the original byte-for-byte outside that value.
- [x] Verdict reader table tests cover `pass`, `fail`, `partial`, no report present, and unreadable frontmatter, plus newest-report selection when several reports exist.
- [x] Task parsing extracts multi-command Verification sections in section order, verbatim.

## Verification

- `rtk go test ./internal/spec/` — expected: all tests pass.
- `rtk go build ./...` — expected: builds cleanly.

## References

`_prd.md` → User Stories 1, 4, 7; Core Features 2, 12. `_techspec.md` → System Architecture (internal/spec), Interfaces (internal/spec), Build Order 1. `docs/agents/issue-tracker.md` → tracker conventions.

## Result

### What changed

New package `internal/spec` now owns the on-disk Spec contract end to end; nothing outside it needs to touch spec markdown:

- `ListActive(gitRoot)` discovers Specs under `docs/specs/` (skipping `_archived/`) whose `_prd.md` frontmatter is `status: active`, sorted by slug. Directories without a readable active PRD are skipped so the Interactive Input picker stays robust; `Load` names the exact problem when a slug is requested explicitly.
- `Load(gitRoot, slug)` parses the `_tasks.md` frontmatter (`schema: spec-tasks/v1`, `graph.nodes[]`) as the authoritative graph source, validates it (existing task files, known `needs`, acyclic, plus duplicate/empty node guards), parses every task file, and returns Tasks in deterministic topological order — Kahn's algorithm picking the first ready node in manifest order.
- Every validation failure is a typed error naming the offending Task or check: `SpecNotFoundError`, `InactiveSpecError`, `ManifestError`, `ManifestSchemaError`, `UnknownNeedError`, `CycleError`, `MissingTaskFileError`, `TaskFileError`, `MissingVerificationError`.
- Task parsing extracts the frontmatter (`task`, `spec`, `status`, `type`, `complexity`), the title from the `# Task NN: <title>` heading, and Verification commands verbatim from backticked bullet entries in section order; zero parseable commands is a validation error.
- `SetStatus(taskPath, status)` rewrites only the status frontmatter value (trailing comments and spacing preserved) and keeps the rest of the file byte-for-byte; `ReloadTask(gitRoot, task)` re-reads a task file after an Agent edit, refreshing Status/Title/Type/Verification while leaving manifest-owned fields alone. `Task.File` is repository-root-relative so it doubles as the commit pathspec for the Daemon.
- `QAVerdict(specDir)` reads the newest `qa/qa-report-*.md` frontmatter and reports `pass`/`fail`/`partial` as values, missing as `ErrNoQAReport` (`errors.Is`-able), and unreadable as `QAReportError`.
- No new module dependencies; the frontmatter split mirrors the Round artifact approach over the existing `gopkg.in/yaml.v3`.

### Commands run

- `rtk go test ./internal/spec/` — pass (42 tests).
- `rtk go build ./...` — builds cleanly.
- `make verify` — pass (fmt-check, 289 tests in 16 packages, `roundfix skills check`, build).

### Evidence per acceptance criterion

- Typed-error table: `TestLoadReturnsTypedValidationErrors` covers missing Spec, inactive Spec, missing manifest, wrong schema, empty graph, duplicate id, unknown `needs`, cycle, missing task file, no Verification commands, unparseable frontmatter, and unsupported status — each asserting the typed error via `errors.As` and the offending Task/check in the message; `TestLoadAcceptsValidGraph` covers the valid case.
- Deterministic order: `TestLoadReturnsTasksInDeterministicTopologicalOrder` runs `Load` 20 times over a diamond manifest and asserts the identical manifest-order-tiebroken sequence every time.
- Byte-preserving rewrite: `TestSetStatusRewritesOnlyTheStatusValue` compares the whole rewritten file against the original with only the status value substituted, including a trailing-comment case; `TestSetStatusRejectsInvalidInput` proves failed rewrites leave the file untouched.
- Verdict reader: `TestQAVerdictReadsSupportedVerdicts` (`pass`/`fail`/`partial`), `TestQAVerdictReportsMissingReports` (no qa dir, no matching report), `TestQAVerdictReportsUnreadableReports` (no frontmatter, no verdict field, unsupported value), `TestQAVerdictSelectsTheNewestReport` (three dated reports, newest wins).
- Verbatim multi-command extraction: `TestLoadParsesTaskFiles` asserts two commands extracted in section order, bullets without backticks skipped, and backticked text outside `## Verification` ignored.

### Follow-ups

- `ListActive` intentionally tolerates unreadable `_prd.md` frontmatter (skip, not error) to keep the picker usable; revisit if task_08 wants a diagnostic listing.
- Consumers in task_05/task_06 should join `Task.File` to the git root for `ReloadTask`/`SetStatus` paths and reuse it directly as the per-Task commit pathspec.
