---
task: task_07
spec: 0001-implement-command
status: completed
type: backend
complexity: medium
---

# Task 07: Add the opt-in QA gate

## Overview

Deliver the QA gate as one vertical slice: the `--qa` flag, the engine's qa-gate step gated on a fully completed Task Graph, verdict settling from the QA Report frontmatter, and the QA Report commit. Verifiable through engine and CLI tests over a fake Agent that writes QA Reports.

## Requirements

1. MUST add the `--qa` flag to the Implement Command; without it, behavior is exactly as shipped in the previous task.
2. MUST run the qa-gate step only when every Task in the Task Graph is `completed` at the end of the cycle (including Tasks completed by earlier Runs); otherwise the step is skipped and the outcome comes from the Task results alone.
3. MUST invoke the Agent with the QA prompt, then read the verdict from the newest QA Report frontmatter: only `pass` lets the Run end `Clean` — `partial`, `fail`, missing report, and unreadable verdict all end the Run `Unresolved` (ADR-0015 plus the recorded `partial` decision).
4. MUST commit the QA Report in its own commit either way, with the message `docs(qa): qa report for <slug> (<verdict>)` and the `Roundfix-Spec` trailer; a missing report commits nothing.
5. MUST journal the step as a `daemon.qa` Run Event carrying the report path and verdict, and print the verdict line on stdout after the per-Task lines.
6. MUST support the QA-only Run: every Task already completed plus `--qa` creates a Run consisting of only the qa-gate step.

## Subtasks

- [x] `--qa` flag and plan wiring into the engine cycle
- [x] All-Tasks-completed gating including the QA-only Run
- [x] Agent invocation with the QA prompt and verdict reading
- [x] Verdict settling of the Run outcome
- [x] QA Report commit and `daemon.qa` journaling

## Acceptance Criteria

- [x] Verdict matrix tests: `pass` → `Clean` exit 0; `partial`, `fail`, missing, unreadable → `Unresolved` exit 1 — with the QA Report committed in its own commit in every case a report exists.
- [x] A Run with any failed or skipped Task never invokes the QA step even with `--qa`.
- [x] A QA-only Run (all Tasks previously completed, `--qa` set) runs the step and settles the outcome from the verdict.
- [x] The `daemon.qa` event appears in the Run Event Journal with the verdict; stdout carries the verdict line after the Task lines.
- [x] Without `--qa`, output and outcomes are byte-identical to the previous task's behavior.

## Verification

- `rtk go test ./internal/daemon/ ./internal/cli/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 5; Core Feature 9; Success Metrics (QA opt-in). `_techspec.md` → API Contracts (commit contract), Build Order 7, Decisions (QA verdict). ADR-0015.

## Result

The Implement Command gained the opt-in QA gate as one vertical slice.

**Behavior shipped.** `roundfix implement --qa` ends the Run with the
qa-gate step, which runs only when every Task in the Task Graph —
including Tasks completed by earlier Runs — is `completed` at cycle end;
any failed, skipped, or pending Task skips the step and the outcome comes
from the Task results alone. The step takes a before-snapshot, runs the
Agent with `BuildQAPrompt` as the next Batch ordinal under
`ResolvingWithAgent`, and settles the verdict from the newest QA Report
frontmatter: the report values `pass`/`fail`/`partial`, plus the Daemon
settlements `missing` (no report) and `unreadable` (report exists, verdict
unreadable). Only `pass` ends the Run `Clean` (exit 0); everything else
ends it `Unresolved` (exit 1). Whenever a report exists the QA Report gets
its own commit — `docs(qa): qa report for <slug> (<verdict>)` with the
`Roundfix-Spec` trailer — staged from the QA step's snapshot diff with the
report ensured; a missing report commits nothing. The step journals a
`daemon.qa` Run Event (payload `verdict` + `report`, the report path
relative to the working tree) on the QA Batch ordinal, and stdout gains
the deterministic verdict line `qa <verdict> — <report path>` (missing:
`qa missing — no QA Report found`) after the per-Task lines and before
the outcome line. All Tasks already completed plus `--qa` creates a
QA-only Run; without `--qa` that stays the no-Run report and every other
behavior is byte-identical to the previous slice. A non-stop Agent failure
during the QA step halts the cycle as an infrastructure failure (Run
`Failed`), since the step has no per-item settlement to fall back on.
`spec.NewestQAReport` was added (additive) so the engine can carry the
report path; `spec.QAVerdict` behavior is unchanged.

**Verification (all fresh, all passing).**

- `rtk go test ./internal/daemon/ ./internal/cli/` — 161 passed in 2 packages.
- `rtk go test ./...` — 394 passed in 16 packages.
- `make verify` — fmt-check, full tests, `roundfix skills check`, and build all passed.
- `rtk go run ./cmd/roundfix implement --help` — lists `--qa` with truthful copy.

**Evidence per acceptance criterion.**

- Verdict matrix: `TestTaskCycleQAVerdictMatrixSettlesRunAndCommitsReport`
  (engine: verdict, report path, own-commit message and paths, Batch
  ordinal, daemon.qa payload) and `TestRunImplementQAVerdictMatrix` (CLI:
  exit codes 0/1, Run states Clean/Unresolved, exact stdout, QA commit in
  every case a report exists, journaled daemon.qa event).
- QA never invoked with a failed/skipped Task:
  `TestTaskCycleQAStepSkippedUnlessEveryTaskCompleted` and
  `TestRunImplementQAStepSkippedWhenAnyTaskFails` (zero QA prompts, no
  daemon.qa event, no verdict line).
- QA-only Run: `TestTaskCycleQAOnlyRunWhenEveryTaskAlreadyCompleted` and
  `TestRunImplementQAOnlyRunSettlesOutcomeFromVerdict` (Run created, one
  Agent call as Batch 1, outcome settled from the verdict).
- daemon.qa in the Run Event Journal + stdout ordering: asserted in both
  matrix tests via the journal reader and exact-stdout comparison.
- Without `--qa` byte-identical: every pre-existing implement and daemon
  test passes unchanged; the only two test edits were the guards that
  pinned `--qa` as not-yet-shipped (help vocabulary and the
  flag-not-defined preflight case), which this task supersedes.

**Follow-ups.**

- task_09: Interactive Input gained no QA field; `--qa` is flag-only until
  a picker decision is made.
- task_10: the cockpit can render the `daemon.qa` event (payload `verdict`,
  `report`) on the QA Batch ordinal.
- task_11: the roundfix skill must document `--qa`, the verdict line shape,
  and the QA Report commit contract in the same PR that ships this CLI
  behavior (repo hard rule).
