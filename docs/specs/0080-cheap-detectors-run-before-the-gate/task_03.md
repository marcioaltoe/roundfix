---
task: task_03
spec: 0080-cheap-detectors-run-before-the-gate
status: completed
type: backend
complexity: high
---

# Task 03: Run the mechanical stage before the Agent Session

## Overview

The wiring: the Daemon's QA step runs the mechanical stage first, materializes
its result, and withholds the Agent Session when a blocking fact is present.
This is where the round stops costing thirty minutes to discover something a
Git read already knew.

The stage runs under the Verification ownership ADR-0014 already gives the
Daemon, and reports without ever settling — ADR-0057 keeps Task status
exclusively the Daemon's, and the stage is not the Daemon's settlement path.

## Requirements

1. MUST run the mechanical stage inside the existing QA Task step, before the
   Agent Session is created, and never as a command outside the Task Graph —
   ADR-0091 keeps the gate one terminal node and ADR-0088 removed the per-run
   choice.
2. MUST withhold the Agent Session when the result is blocking, and close a
   mechanical-stage-only report carrying `verdict: fail` with zero pending
   rows.
3. MUST hand the non-blocking result's materialized rows and skips to the
   Agent Session as the report it completes.
4. MUST NOT consume the single Verification repair ADR-0038 allots, nor borrow
   the Verification Capacity ADR-0056 keeps separate from Task Capacity.
5. MUST leave write authority exactly where it is: Task status stays the
   Daemon's, the QA Report stays the only artifact this step writes, and the
   report's naming and recency contract is untouched. The stage itself writes
   nothing and pushes nothing.
6. MUST keep every existing gate behaviour intact: the Daemon still records the
   QA Report either way, any non-pass verdict ends the Run Unresolved, and a
   missing or unreadable verdict still counts as fail.
7. MUST emit the stage's outcome on the existing Run Event surface so its cost
   and its refusals are visible without an agent turn.

## Subtasks

- [ ] Run the stage in the QA step and materialize its result.
- [ ] Withhold the session on a blocking result; seed it otherwise.
- [ ] Assert the repair, capacity, and status boundaries.

## Acceptance Criteria

- [ ] A blocking result produces a closed report and no Agent Session.
- [ ] A non-blocking result produces a seeded report and the Agent Session runs.
- [ ] The Verification repair is not consumed and Verification Capacity is not
      borrowed.
- [ ] Task status is never set by the stage.
- [ ] The recorded report, unresolved outcome, and naming contract are
      unchanged.

## Context

- interface: internal/daemon/task_engine.go
- interface: internal/daemon/engine.go

## Verification

  — expected: exit 0; the stage and gate tests are selected and pass.
- `output="$(go test -count=1 ./internal/daemon -run 'MechanicalStageWithholdsAgentSession' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; the withholding behaviour is proven by a named test, not
  inferred from the absence of a session.
  — expected: exit 0; the packages carrying gate behaviour and verdict
  semantics stay green.

These commands are deliberately absent: `go build -buildvcs=false ./...` and a
whole-package `go test` sweep both pass against a tree where no work has
happened, so each approves the Task before it starts. Compilation and
regression are the Run-level gate's job; the commands above name cases that
do not exist yet.

The broad `MechanicalStage|QAGate` pattern is deliberately absent: by the time
this Task is dispatched, Task 01's cases already match it and pass, so it would
approve this Task before any work happened. The command above names the one case
this Task must create.

## References

- `_prd.md` → Core Features 1, 2, 7; User Stories 1, 3.
- `_techspec.md` → System Architecture; Build Order 3.
- ADR-0096, ADR-0091, ADR-0088, ADR-0015, ADR-0014, ADR-0057, ADR-0038,
  ADR-0056.

## Result

Implemented the Daemon-owned mechanical stage inside `runQAGate`. The stage
now reads the PRD's tooling authorization, selects the Spec's applicable Task
commits from their Roundfix trailers, evaluates the real mechanical detectors,
and atomically creates the next collision-safe QA Report. A blocking result
materializes `verdict: fail`, publishes its mechanical Run Event, and reaches
the existing verdict settlement and report commit without creating an Agent
Session. A non-blocking result seeds the same report path that the QA prompt
requires the Agent to complete in place.

Focused checks run after the final implementation edit:

- Pre-change signal: `rtk rg -n 'TestMechanicalStageWithholdsAgentSession' internal/daemon/task_engine_test.go`
  exited 1 because the required withholding case did not exist. The first
  focused compile then failed on the deliberately introduced test reference to
  the absent `Dependencies.MechanicalStage` seam.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task03-gocache go test ./internal/daemon -run '^(TestMechanicalStage|TestQAMechanicalRequest|TestWriteMechanicalQAReport|TestTaskCycleQA|TestPerWorkAgentSessionMixedTaskTypesAndQA|TestQAPullRequest)' -count=1`
  passed.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task03-gocache go test ./internal/speccheck -run '^(TestMechanicalAuthorizationReadsThePRDBoundedDeclaration|TestMaterializeMechanicalResult)$' -count=1`
  passed.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task03-gocache go test ./internal/spec -run '^Test(QAVerdict|NewestQAReport)' -count=1`
  passed.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task03-gocache go vet ./internal/daemon ./internal/speccheck`
  passed with no diagnostics.
- `rtk git diff --check` passed.

Acceptance evidence:

1. `TestMechanicalStageWithholdsAgentSession` observes the QA Task still
   pending while the stage runs, then proves the blocking result writes a
   closed fail report with its finding-blocked row, no pending text, no Agent
   request, and the existing QA Report commit message.
2. `TestMechanicalStageSeedsReportBeforeAgentSession` proves the Agent sees
   the materialized skips in the exact seeded path, runs once, and returns the
   existing passing verdict. The prompt forbids creating a second report.
3. Both stage-branch cases assert zero Verifier calls and zero
   `daemon.verification` events. The mechanical path never calls the
   Verification-capacity or Verification-repair functions.
4. The blocking fake reads Task status during the stage and observes the
   original pending value; only the existing Daemon settlement path later
   writes failed.
5. The existing QA verdict matrix remains green for pass, partial, fail,
   missing, and unreadable reports. The report writer's same-day collision
   case produces `qa-report-2026-01-01-01.md` without overwriting the prior
   report, and the existing `internal/spec` recency and verdict cases pass.
   Both branches publish a `daemon.qa` mechanical event containing blocking,
   finding, skip, duration, and report-path facts.

The commands under this Task's `## Verification` were not run; the Daemon owns
that complete selection and settlement evidence.
