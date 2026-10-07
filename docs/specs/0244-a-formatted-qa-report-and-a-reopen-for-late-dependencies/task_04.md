---
task: task_04
spec: 0244-a-formatted-qa-report-and-a-reopen-for-late-dependencies
status: pending
type: backend
complexity: medium
---

# Task 04: Reopen returns a completed QA gate to pending over a Late Dependency

## Overview

Gives `roundfix reopen` a second trigger. Today it reopens only a gate above a
dependency that is no longer completed. It now also reopens over a Late
Dependency proven from Git. This Task answers the second report of the Backlog
Entry "An unformatted QA report breaks the next Run's precondition" of
2026-10-06: a corrective Task, already `completed`, became a dependency of a
completed QA gate, reopen refused, and the operator set `status: pending` by
hand. It is verifiable on its own through the spec and CLI tests below.

## Requirements

1. MUST add `internal/spec/gate_closure.go` with
   `QAGateClosure(manifestPath string, content []byte) (string, []string, error)`,
   which parses manifest bytes with the package's manifest parser and returns
   the QA Task id and its sorted transitive dependency closure.
2. MUST add `ReopenGateForLateDependencies(taskPath, reportPath string, taskIDs []string, date time.Time) error`
   to `internal/spec/task.go`. It sets the QA Task `pending` in one atomic
   replacement and appends `## Invalidation` with `- Date:`, `- QA Report:` and
   `- Dependencies added after the QA Report: ` followed by each id in
   backticks (`_techspec.md` → Invariant 11). `ReopenGate` and its record stay
   byte-identical.
3. MUST change reopen so that a completed gate whose dependencies are all
   completed is checked for a Late Dependency per `_techspec.md` → Invariant
   10: resolve the newest QA Report, find the oldest commit in `HEAD`'s history
   that added it with `git log --diff-filter=A`, read the Spec's `_tasks.md` at
   that commit, and compare its QA closure with the working tree's. Git runs in
   the checkout's Git root with the command's context.
4. MUST keep every existing outcome: a stale gate reopens as today, and an
   absent commit, an absent or unparsable recorded manifest, a report outside
   the Git root, or no Late Dependency keeps today's refusal, message and exit
   code. The pre-write recheck MUST re-derive the Late Dependency ids and
   refuse with exit 2 when they changed (`_techspec.md` → Invariant 9 and
   Invariant 11).
5. MUST update the reopen usage text so it names both triggers.
6. MUST add `internal/spec/gate_closure_test.go` with
   `TestQAGateClosureFromManifestBytes` and
   `TestReopenGateForLateDependenciesRecordsTheLine`, and
   `internal/cli/reopen_late_dependency_test.go` with
   `TestReopenReopensALateDependencyThroughTheBuiltBinary` (a temporary
   repository whose newest QA Report was committed before `task_03` existed:
   exit 0, QA Task `pending`, the record line, the report unchanged),
   `TestReopenRefusesWhenTheNewestReportIsUncommitted` and
   `TestReopenRefusesWhenTheRecordedClosureMatches` (both exit 2 with today's
   message and no mutation).
7. MUST NOT change the Task Graph loader, `implement`, `settle` or `deliver`,
   and MUST NOT change any existing reopen or spec test expectation.

## Subtasks

- [ ] Add the closure helper and the Late Dependency record.
- [ ] Prove a Late Dependency from Git in reopen and recheck it.
- [ ] Update the usage text.
- [ ] Write the spec and CLI tests.

## Acceptance Criteria

- [ ] A completed gate with a Late Dependency reopens with the new record line.
- [ ] Without a recording commit or a Late Dependency, reopen refuses exactly
      as before.
- [ ] Every existing reopen and spec test passes unchanged.

## Context

- interface: `internal/cli/reopen.go`
- interface: `internal/spec/task.go`
- creates: `internal/spec/gate_closure.go`
- creates: `internal/spec/gate_closure_test.go`
- creates: `internal/cli/reopen_late_dependency_test.go`
- instruction: `internal/cli/reopen_test.go`
- instruction: `internal/spec/spec.go`
- instruction: `docs/adr/0249-a-qa-step-formats-its-qa-directory-and-reopen-sees-a-late-dependency.md`

## Verification

- `out="$(go test -count=1 -v -run '^(TestQAGateClosureFromManifestBytes|TestReopenGateForLateDependenciesRecordsTheLine)$' ./internal/spec 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestQAGateClosureFromManifestBytes TestReopenGateForLateDependenciesRecordsTheLine; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { printf 'missing passing test %s\n' "$t" >&2; exit 1; }; done` — expected: exit 0; before this Task neither test exists, and after it both pass.
- `out="$(go test -count=1 -v -run '^(TestReopenReopensALateDependencyThroughTheBuiltBinary|TestReopenRefusesWhenTheNewestReportIsUncommitted|TestReopenRefusesWhenTheRecordedClosureMatches)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for t in TestReopenReopensALateDependencyThroughTheBuiltBinary TestReopenRefusesWhenTheNewestReportIsUncommitted TestReopenRefusesWhenTheRecordedClosureMatches; do printf '%s\n' "$out" | grep -q -- "--- PASS: $t " || { printf 'missing passing test %s\n' "$t" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the three tests exists, and after it the built binary reopens over a Late Dependency and refuses the two other cases.

## References

- `_prd.md` → Core Feature 3; Goals; Success Metric 4; Success Metric 5
- `_techspec.md` → API Contract 3; Invariant 9; Invariant 10; Invariant 11; Build Order 4
- ADR-0249
