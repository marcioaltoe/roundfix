---
task: task_02
spec: 0351-a-secret-refusal-of-a-task-commit-is-repairable-and-visible
status: completed
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


## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-structs-interfaces, golang-testing, testing-boss, systematic-debugging, no-workarounds, tech-writer, evidence-gate.

The Causes builder now emits one `secret-scan` item per refusal event, including
legacy events without payload Task or attempt. The item carries only match
locations and kinds, uses `repository-convention`, and has null signature and
evidence. The failed scan check marks its attempt seen without emitting a
verification item, so its later failed verdict cannot duplicate the refusal.
Other items serialize `matches` as null; the schema and text formatter remain
unchanged.

Acceptance evidence:

- Each refusal appears once in text and JSON: the new report test seeds one
  terminal Run with a legacy refusal and a new refusal containing two matches,
  checks Work Item fallback and payload Task precedence, null/1 attempts,
  locations, classes, match arrays and summary counts, and rejects the duplicate
  verification item. The new CLI test checks the entire prescribed two-line
  transcript, exit 0, empty stderr and the JSON fields against an isolated Run
  Database with the pinned fixture ID and creation time.
- No fixture value reaches an item: the report test builds the value at runtime,
  includes it in unknown input match fields, and rejects every consecutive
  eight-byte slice in each serialized item. It also checks explicit null
  `matches` for verification and corrective items. Prefix comparisons against
  `HEAD` confirm every byte of both existing test files remains unchanged before
  the appended tests. Byte comparisons confirm the signature table and
  `internal/cli/runs_causes.go` remain unchanged. Signature SHA256:
  `53e5e5a764109b92806c6efe4d88603c13e5de2c6d21816755406deb11b6dd5f`.

Characterization and focused checks:

- Before changing code, an untracked temporary test supplied through a Go
  overlay created a terminal Run with a `task_01` start, one legacy
  `secret_scan` match and a failed settlement. Command:
  `go test -overlay "$probe/overlay.json" -count=1 -timeout 5m -v -run '^TestCharacterizeSecretRefusalMissing$' ./internal/runcause`.
  Exit 0: the characterization asserted one terminal Run and zero Causes items.
  The temporary test was never added to the worktree.
- Before the builder change,
  `go test -count=1 -timeout 5m -run '^TestBuildReportsASecretRefusalAsItsOwnItem$' ./internal/runcause`
  exited 1 with `secret scan duplicated as a verification item for task_02 attempt 1`
  and only one item instead of two.
- After the builder change,
  `go test -count=1 -timeout 5m -v -run '^TestBuildReportsASecretRefusalAsItsOwnItem$' ./internal/runcause`
  exited 0.
- After the last test edit,
  `go test -count=1 -timeout 5m ./internal/runcause` exited 0, covering the
  existing package tests and the new report test.
- `go test -count=1 -timeout 10m -run 'TestRunsCauses|TestBuildLeavesTheRunDatabaseUnchanged' ./internal/cli`
  exited 0, covering the existing Causes command tests and the new transcript
  and JSON test.
- Required sabotage: copied `report.go` to a temporary directory, removed only
  the `p.Command == runevent.SettlementCheckSecretScan` skip, and mapped that
  copy through a Go overlay. Command:
  `go test -overlay "$sabotage/overlay.json" -count=1 -timeout 5m -v -run '^TestBuildReportsASecretRefusalAsItsOwnItem$' ./internal/runcause`.
  Expected exit 1 observed. The named test reported
  `secret scan duplicated as a verification item for task_02 attempt 1` and
  three items rather than two. Production source was unchanged by the probe.
- `git diff --check` exited 0; `gofmt -l` on the three changed Go files printed
  nothing. Changed-path inspection found only those files and this Task file.

Execution note: the sandbox could not start with the configured symlinked
writable root, so commands used the same physical checkout with escalated
execution. The only pre-existing diff was the Daemon's `status: in_progress`;
that field was preserved. No other Task or manifest was edited. Declared Task
Verification and repository Verification are left to the Daemon. No commit,
push or pull request was created. No follow-up work was identified.

## Carry-forward provenance

- Source Run: `run_20261105T015847Z_709cb0723a1fbaae`
- Source commit: `94a172ae9c39a07c755024b7b334d8641c79df90`
