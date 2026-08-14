---
task: task_07
spec: 0095-a-verification-that-ran-before-anyone-believed-it
status: completed # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: low
---

# Task 07: Let a Task declare the file it creates

## Overview

A Task's context accepts `interface:` and `instruction:` entries, and both
require an existing path. A Task that creates a file therefore cannot name its own
output without the reference check refusing the graph — it hit four task files in
one fleet session. This slice adds a declared output whose existence is not
required.

## Requirements

1. MUST accept a context entry declaring a path the Task creates, without
   requiring that path to exist.
2. MUST keep requiring existence for the entry kinds that name inputs, so a
   mistyped interface path is still refused.
3. MUST hold the declared output to the same path hygiene as the other kinds:
   clean, repository-relative, unique, and within the per-Task entry ceiling.
4. MUST NOT change how the Daemon fills the context bundle beyond including the
   declared path.

## Subtasks

- [x] Accept the declared-output entry kind.
- [x] Keep existence required for input kinds.
- [x] Cover a declared output that does not exist and a mistyped input.

## Acceptance Criteria

- [x] A Task declaring an output that does not exist passes the reference check.
- [x] A Task declaring an input that does not exist still fails it.
- [x] An unclean, absolute, or duplicate declared path is refused.
- [x] The per-Task entry ceiling counts declared outputs.

## Verification

- `go test -count=1 ./internal/spec ./internal/speccheck -run 'TestContextDeclaredOutput' -v > /tmp/0095-t07.log 2>&1; s=$?; grep -q '^--- PASS: TestContextDeclaredOutput' /tmp/0095-t07.log || { cat /tmp/0095-t07.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today, where no such test exists.
- `! grep -qi 'no tests to run' /tmp/0095-t07.log` — expected: exits 0, refusing a vacuous run.
- `grep -rq 'creates' internal/spec/spec.go && grep -c 'PASS' /tmp/0095-t07.log | { read n; test "$n" -ge 2 || { echo "expected the declared-output and the mistyped-input cases, got $n"; exit 1; }; }` — expected: exits 0, proving the entry kind reached the parser and that both directions are exercised.

## Context

- interface: `internal/spec/spec.go`

## References

`_techspec.md` → Build Order 7. `_prd.md` → Core Feature 5; User Story 3.
Evidence: `docs/findings/2026-08-12-three-consecutive-specs-measure-the-loop.md`
finding 3, which measured it against four task files.

## Result

### Implementation

- Added `creates:` as a Task Context kind. It uses the existing path cleaner and
  per-Task entry ceiling, and any duplicate involving a declared output is
  refused while the existing duplicate-input behavior remains unchanged.
- The reference check exempts only `creates:` entries from repository existence
  checks. `instruction:` and `interface:` entries still require an existing path.
- The Daemon includes a declared output through the existing interface-path
  bundle flow. No bundle field, limit, ordering, or prior-file behavior changed.

### Acceptance evidence

| Criterion | Focused evidence |
| --- | --- |
| A missing declared output passes the reference check | `internal/speccheck/TestContextDeclaredOutput` loads a fixture whose `generated/client.go` output does not exist and asserts that no `SC-REF-UNRESOLVED` finding names it. |
| A missing input still fails | The same test asserts the fixture produces exactly one `SC-REF-UNRESOLVED` finding and that it names the mistyped `interface:` path `missing/guide.md`. |
| Invalid or duplicate output paths are refused | `internal/spec/TestContextDeclaredOutput` passed its absolute, unclean, and duplicate-output refusal subtests; each returned a typed `TaskContextError` with the expected reason. |
| Declared outputs count toward the per-Task ceiling | `internal/spec/TestContextDeclaredOutput/counts_the_output_in_the_entry_ceiling` passed with 50 input entries plus one `creates:` entry and asserted the 50-entry refusal. |
| The Daemon includes the declared path without other bundle changes | `internal/daemon/TestAssembleTaskContextBundleReservesExplicitPathsAndCountsOmittedPriorFiles` passed with a nonexistent `creates:` path reserved under the unchanged 200-path limit. |

### Focused checks

- Red signal: with a task-scoped Go cache, the new parser test failed to compile
  because `ContextKindCreates` was undefined, and the checker test failed because
  `creates:` was an unsupported label.
- `GOCACHE=/tmp/roundfix-task07-gocache rtk go test ./internal/spec -run '^TestContextDeclaredOutput$'` passed 6 tests.
- `GOCACHE=/tmp/roundfix-task07-gocache rtk go test ./internal/speccheck -run '^TestContextDeclaredOutput$'` passed 1 test.
- `GOCACHE=/tmp/roundfix-task07-gocache rtk go test ./internal/daemon -run '^TestAssembleTaskContextBundleReservesExplicitPathsAndCountsOmittedPriorFiles$'` passed 1 test.
- The affected package suites passed: `internal/spec` (306 tests),
  `internal/speccheck` (177 tests), and `internal/daemon` (210 tests).
- `GOCACHE=/tmp/roundfix-task07-gocache rtk make verify-incremental` exited 0;
  the full Go suite, skill checks, and build passed.
- The Daemon-owned commands under `## Verification` were not run.
