---
task: task_02
spec: 0080-cheap-detectors-run-before-the-gate
status: pending
type: backend
complexity: medium
---

# Task 02: Tell the gate prompt what changed

## Overview

The QA prompt is about fifteen lines and carries no changed-file context at
all, while every Task prompt already assembles a Spec Context Bundle. That is
why scoping a gate round has been impossible: the gate cannot narrow what it
re-observes because it is never told what moved.

The same prompt's contract text names only two of the three typed
blocked-cause counts, while ADR-0080, the qa-gate skill, and the report reader
all require three. An agent that gets the third right today gets it from the
skill, not from the contract it was handed.

## Requirements

1. MUST assemble and pass the Spec Context Bundle to the QA prompt, using the
   same builder the Task prompt uses rather than a parallel implementation.
2. MUST carry the previous QA Report's identity — its path and the head it was
   taken at — when one exists for this Spec, and say plainly when none does.
3. MUST name all three blocked-cause counts in the prompt's contract text,
   matching ADR-0080, the skill, and the report reader exactly.
4. MUST keep the existing prompt facts unchanged: slug, spec directory, PRD
   path, Run Worktree branch, Spec target branch, user checkout, and the
   pull-request fact with its resolved/unresolved distinction.
5. MUST respect the Spec Context Bundle's existing path limit rather than
   introducing a second one.
6. MUST NOT change verdict semantics, report naming, or anything the gate does
   with the information.
7. MUST name the two new cases `TestBuildQAPromptCarriesTheSpecContextBundle`
   and `TestBuildQAPromptCarriesThePreviousReportIdentity`. A bare `QAPrompt`
   pattern matches the five cases that already pass, so it would approve this
   Task before any work happened; the declared Verification names cases that do
   not exist yet.

## Subtasks

- [ ] Assemble the bundle for the QA prompt through the shared builder.
- [ ] Carry the previous report's identity, including its absence.
- [ ] Repair the third blocked-cause count in the contract text.

## Acceptance Criteria

- [ ] The assembled QA prompt contains the changed-path context for a Spec
      whose branch has commits, and says so plainly when it has none.
- [ ] The prompt names `rows_blocked_environment`, `rows_blocked_finding`, and
      `rows_blocked_declared`.
- [ ] Every pre-existing prompt fact still appears.
- [ ] No verdict, report-naming, or gate behaviour changed.

## Context

- interface: internal/agent/spec_prompt.go
- interface: internal/daemon/task_context.go

## Verification

- `output="$(go test -count=1 ./internal/agent -run '^TestBuildQAPromptCarriesTheSpecContextBundle$' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
- `output="$(go test -count=1 ./internal/agent -run '^TestBuildQAPromptCarriesThePreviousReportIdentity$' -v 2>&1)"; st=$?; printf '%s\n' "$output" | grep -q -- '--- PASS' && [ "$st" -eq 0 ]`
  — expected: exit 0; the QA prompt tests are selected and pass.
- `grep -q 'rows_blocked_declared' internal/agent/spec_prompt.go && grep -q 'rows_blocked_environment' internal/agent/spec_prompt.go && grep -q 'rows_blocked_finding' internal/agent/spec_prompt.go`
  — expected: exit 0; all three counts are named in the prompt contract, which
  is false on the current surface.
  — expected: exit 0; the bundle reaches the QA prompt and both packages stay
  green.

These commands are deliberately absent: `go build -buildvcs=false ./...` and a
whole-package `go test` sweep both pass against a tree where no work has
happened, so each approves the Task before it starts. Compilation and
regression are the Run-level gate's job; the commands above name cases that
do not exist yet.

## References

- `_prd.md` → Core Features 4, 8; User Story 4.
- `_techspec.md` → System Architecture (the QA prompt); Build Order 2.
- ADR-0080, ADR-0096.
