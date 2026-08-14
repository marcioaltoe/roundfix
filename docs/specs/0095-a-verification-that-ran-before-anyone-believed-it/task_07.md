---
task: task_07
spec: 0095-a-verification-that-ran-before-anyone-believed-it
status: pending # pending | in_progress | completed | failed — only implement-task changes this
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

- [ ] Accept the declared-output entry kind.
- [ ] Keep existence required for input kinds.
- [ ] Cover a declared output that does not exist and a mistyped input.

## Acceptance Criteria

- [ ] A Task declaring an output that does not exist passes the reference check.
- [ ] A Task declaring an input that does not exist still fails it.
- [ ] An unclean, absolute, or duplicate declared path is refused.
- [ ] The per-Task entry ceiling counts declared outputs.

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
