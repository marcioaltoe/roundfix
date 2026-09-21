---
task: task_06
spec: 0152-one-declared-acceptance-policy
status: completed
type: backend
complexity: medium
---

# Task 06: The decider that actually settles the gate

## Overview

Pre-PR review found that this Spec unified two deciders and missed the one that
matters. `runQAGate` settles the terminal `qa` Task from the raw verdict:

```go
// internal/daemon/task_engine.go:2417
if verdict == spec.VerdictPass {
	qaStatus = spec.StatusCompleted
}
```

Every other verdict settles the Task `failed`. It consults neither
`archiveUnprovenActions` nor the rendered command. So a qualifying `partial`
still fails the gate in a real Implement Run, and Core Feature 3 is not
delivered by Tasks 01 through 05.

The comment above that line paraphrases ADR-0015 as "Every value except
`spec.VerdictPass` ends the Run Unresolved". ADR-0015 says the Daemon "ends the
Run as Unresolved on a **failing** verdict". A qualifying `partial` is not a
failing verdict, so this Task settles it under the accepted decision rather than
against it.

Settling in process also answers the second review finding: the rendered command
invokes whatever `roundfix` is first on `PATH`, and an installed 0.14.1 rejects
`qa-report` outright. A settlement path that calls the decision directly cannot
inherit a stranger's binary.

## Requirements

1. MUST apply the one eligibility decision when settling the terminal `qa`
   Task, in process, so a qualifying `partial` settles `completed`.
2. MUST keep settling `completed` for a `pass`, including one carrying an
   environment-blocked row.
3. MUST keep settling `failed` for `fail`, for a missing or unreadable report,
   and for a `partial` that does not qualify, with the reason each gives today.
4. MUST NOT depend on a `roundfix` resolved from `PATH` to settle.
5. MUST cover the qualifying-partial path with a TaskCycle regression, not only
   a unit test of the decision.
6. MUST NOT change what `archiveUnprovenActions` or the rendered command decide.

## Subtasks

- [ ] Apply the decision in the gate settlement path.
- [ ] Add a TaskCycle regression for the qualifying partial.
- [ ] Pin the unchanged settlements for pass, fail, missing and unreadable.

## Acceptance Criteria

- [ ] A Run whose newest report is a qualifying `partial` settles its `qa` Task
      `completed` and does not end Unresolved.
- [ ] A `pass` with an environment-blocked row still settles `completed`.
- [ ] `fail`, missing, unreadable and a non-qualifying `partial` still settle
      `failed` with today's reason.
- [ ] Settlement performs no `PATH` lookup of `roundfix`.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/spec/qa.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestTaskCycleSettlesAQualifyingPartial" ./internal/daemon 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestTaskCycleSettlesAQualifyingPartial"` — expected: exit 0; before this Task the case does not exist, so the command fails.
- `test "$(grep -cF 'if verdict == spec.VerdictPass {' internal/daemon/task_engine.go)" = "0"` — expected: exit 0; the gate no longer settles on a raw verdict comparison. Before this Task line 2417 reads exactly that, so the command fails.

## References

- [_techspec.md](_techspec.md) — The one decision

## Result

Implementation:

- The Daemon now reads the newest QA Report and applies
  `spec.QAReportEligibility` in process before settling the terminal `qa` Task.
  The raw report verdict remains available for reporting and commit messages,
  while `TaskCycleResult.QAAccepted` carries the shared eligibility decision to
  Implement Run settlement.
- Implement Run disposition now uses that accepted/not-accepted result, so a
  qualifying `partial` can end `Clean` without changing what the report says.
- Archive eligibility and the rendered QA Verification command were not
  changed.

Focused-check evidence:

- Red signal: `rtk env GOCACHE=/private/tmp/roundfix-task06-go-cache go test
  -run '^TestTaskCycle(QAVerdictMatrixSettlesRunAndCommitsReport|SettlesAQualifyingPartial)$'
  ./internal/daemon` failed before implementation because `TaskCycleResult`
  had no accepted eligibility result.
- The same focused daemon command passed after implementation. It covers the
  qualifying partial and the pass, fail, missing, unreadable, and
  non-qualifying-partial settlement matrix.
- `rtk env GOCACHE=/private/tmp/roundfix-task06-go-cache go test -run
  '^TestRunImplementQAVerdictMatrix$' ./internal/cli` passed. The matrix proves
  a qualifying partial ends `Clean` and a non-qualifying partial remains
  `Unresolved`.
- `rtk env GOCACHE=/private/tmp/roundfix-task06-go-cache go test
  ./internal/daemon ./internal/cli` passed the daemon package. The CLI package
  reached two unrelated force-stop integration failures because sandboxed
  process-table enumeration returned `operation not permitted`; its focused QA
  disposition matrix had already passed.
- `rtk git diff --check` passed.

Acceptance evidence:

- `TestTaskCycleSettlesAQualifyingPartial` gives the Spec one unreachable
  declaration and a newest `partial` report with one declared-blocked row. The
  TaskCycle records the report as accepted and settles the `qa` Task
  `completed`.
- The same regression replaces `PATH` with an isolated directory containing
  `git` but no `roundfix`; settlement still succeeds through the in-process
  decision.
- `TestTaskCycleQAVerdictMatrixSettlesRunAndCommitsReport` keeps a `pass` with
  one environment-blocked row `completed`. It keeps `fail`, missing,
  unreadable, and a partial with no declared-blocked row `failed`, and pins the
  existing reasons `QA verdict fail`, `QA verdict missing`, `QA verdict
  unreadable`, and `QA verdict partial`.
- `TestRunImplementQAVerdictMatrix` proves the accepted qualifying partial does
  not end the Run `Unresolved`; the negative cases still do.

Not run:

- The Task's declared `## Verification` commands — reserved for the Daemon by
  the assigned execution contract.
