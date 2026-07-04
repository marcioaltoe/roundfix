---
task: task_07
spec: 0001-implement-command
status: pending
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

- [ ] `--qa` flag and plan wiring into the engine cycle
- [ ] All-Tasks-completed gating including the QA-only Run
- [ ] Agent invocation with the QA prompt and verdict reading
- [ ] Verdict settling of the Run outcome
- [ ] QA Report commit and `daemon.qa` journaling

## Acceptance Criteria

- [ ] Verdict matrix tests: `pass` → `Clean` exit 0; `partial`, `fail`, missing, unreadable → `Unresolved` exit 1 — with the QA Report committed in its own commit in every case a report exists.
- [ ] A Run with any failed or skipped Task never invokes the QA step even with `--qa`.
- [ ] A QA-only Run (all Tasks previously completed, `--qa` set) runs the step and settles the outcome from the verdict.
- [ ] The `daemon.qa` event appears in the Run Event Journal with the verdict; stdout carries the verdict line after the Task lines.
- [ ] Without `--qa`, output and outcomes are byte-identical to the previous task's behavior.

## Verification

- `rtk go test ./internal/daemon/ ./internal/cli/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Story 5; Core Feature 9; Success Metrics (QA opt-in). `_techspec.md` → API Contracts (commit contract), Build Order 7, Decisions (QA verdict). ADR-0015.
