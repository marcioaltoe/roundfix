---
task: task_02
spec: 0351-a-secret-refusal-of-a-task-commit-is-repairable-and-visible
status: pending
type: backend
complexity: medium
---

# Task 02: The Causes report lists each secret refusal as a secret-scan item

## Overview

`roundfix runs causes` (`runcause.Build`, `internal/runcause/report.go`)
reads only failed Verification attempts and corrective Tasks, so Fiscus's Run
refused by the Task commit's Secret Scan produced no item. After this Task
each `secret_scan` event of a terminal Spec Run becomes one item of kind
`secret-scan` with Spec, Task, Run, attempt, path, line and kind, in text and
JSON, never the value, for events written before and after task_01; the
attempt's failed `settlement check: secret scan` is not listed again as a
`verification` item. This answers the first problem of the report reported by
fiscus on 2026-11-04
([references](references/2026-11-04-fiscus-a-task-commit-secret-refusal-is-invisible.md)).

## Requirements

1. MUST, before changing code, characterize the starting tree with a
   temporary test that is not committed: a terminal Run holding a `task_01`
   start, a `secret_scan` event with one match and a failed settlement yields
   no item for that Run; record it in the Result.
2. MUST change `runcause.Build` exactly as `_techspec.md` → API Contract 6
   states: add `Match` and `Item.Matches` (`json:"matches"`, `null` on every
   other item), read `matches` from the event payload, build the
   `secret-scan` item with class `repository-convention`, `signature` and
   `evidence` `null`, `attempt` from the payload or `null` when absent, the
   Task from the payload `task` or else the event's Work Item, and the check
   `<path>:<line>: <kind>` joined by `; ` through `publicCheck`; and make a
   failed verification event whose command is
   `runevent.SettlementCheckSecretScan` mark its attempt seen without adding
   a `verification` item. `signatures.json`, the schema name, the text line
   shape of `printRunsCauses` and the help text keep their bytes.
3. MUST add `TestBuildReportsASecretRefusalAsItsOwnItem` to
   `internal/runcause/report_test.go`, calling `t.Parallel()` first, over one
   terminal Run with:
   - a 0.130.1-shaped `secret_scan` event without `task` or `attempt` in its
     payload and with Work Item `task_01`;
   - a task_01-shaped attempt 1 of `task_02`: a failed verification event for
     `settlement check: secret scan`, its `secret_scan` event with `attempt` 1
     and two matches, and its failed verdict.
   It asserts exactly two `secret-scan` items with their attempts (`null`,
   then 1), checks, classes and `matches`, no `verification` item for
   `task_02` attempt 1, the summary counting both under
   `repository-convention`, and that no item field contains eight consecutive
   bytes of a value the fixture builds at runtime.
4. MUST add `TestRunsCausesPrintsASecretRefusal` to
   `internal/cli/runs_causes_test.go`, calling `t.Parallel()` first, seeding
   one terminal Run of `0300-example` whose ID and creation time are pinned
   as `seedCausesCLI` pins them, holding only a `task_01` start, a failed
   verification event for `settlement check: secret scan`, the `secret_scan`
   event with `attempt` 1 and one match at `fixture.txt:2:
   credential-assignment`, and the failed verdict. It asserts the whole text
   output of `roundfix runs causes`, two lines:
   `secret-scan 0300-example/task_01 run_20261001T120000Z_0000000000000001 attempt-1 repository-convention -: fixture.txt:2: credential-assignment`
   and
   `Causes: 1 item(s) from 1 terminal Spec Run(s) in window *..*; scope-or-authorization 0, shared-section-contract 0, repository-convention 1, implementation-defect 0, environment 0, unclassified 0; repository knowledge 1 of 1 classified; signatures <digest>`,
   where `<digest>` is `causesDigest(t)`; and, for `--format json`, the one
   item's `kind`, `attempt`, `check`, `class`, null `signature` and
   `evidence`, and `matches` `[{"path":"fixture.txt","line":2,"kind":"credential-assignment"}]`.
5. MUST keep every existing test of `internal/runcause` and
   `internal/cli/runs_causes_test.go` byte-identical and passing, and run
   the `runcause.Build` sabotage of `_techspec.md` → Testing Approach on a
   temporary copy, recording which named test fails.
6. MUST NOT touch `internal/runcause/signatures.json`,
   `internal/runcause/testdata/`, `internal/daemon`, `internal/delivery`,
   `internal/cli/runs_causes.go`, any document or skill, the Makefile, CI
   workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Characterize the missing item.
- [ ] Build the secret-scan item and skip the duplicate verification item.
- [ ] Add the report and CLI tests.
- [ ] Run the sabotage on a temporary copy.

## Acceptance Criteria

- [ ] Each secret refusal of a terminal Spec Run appears once as a `secret-scan` item with location and kind, for old and new events, in text and JSON.
- [ ] No value reaches the report, the signature digest is unchanged, and every existing Causes test passes.

## Context

- interface: `internal/runcause/report.go`
- interface: `internal/runcause/report_test.go`
- interface: `internal/cli/runs_causes_test.go`
- instruction: `internal/cli/runs_causes.go`
- instruction: `internal/runcause/signatures.json`

## Consumers

- The Causes report's items and JSON document — consumers: `./internal/runcause`, `./internal/cli`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestBuildReportsASecretRefusalAsItsOwnItem|TestBuild.*|TestCauses.*|TestClassify.*)$' ./internal/runcause 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestBuildReportsASecretRefusalAsItsOwnItem ' || { echo 'missing pass: TestBuildReportsASecretRefusalAsItsOwnItem' >&2; exit 1; }; out="$(go test -count=1 -timeout 20m -v -run '^(TestRunsCauses.*|TestBuildLeavesTheRunDatabaseUnchanged)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestRunsCausesPrintsASecretRefusal ' || { echo 'missing pass: TestRunsCausesPrintsASecretRefusal' >&2; exit 1; }` — expected: exit 0. Before this Task neither new test exists, so the first pass line is missing; afterwards a secret refusal is its own item in the report and the command, and the existing Causes tests pass.

## References

- `_prd.md` → Goals; Core Feature 4; Success Metric 4
- `_techspec.md` → Current behavior 5; API Contract 6; Invariant 1; Declared breaks; Testing Approach; Build Order 2
