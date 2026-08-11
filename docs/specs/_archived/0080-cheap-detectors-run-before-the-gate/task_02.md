---
task: task_02
spec: 0080-cheap-detectors-run-before-the-gate
status: completed
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

## Result

Implemented the QA prompt context without changing Task status, verdict
evaluation, or report naming. The Daemon now builds the QA gate's
`SpecContextBundle` through the same bounded builder used for Task prompts,
adds the newest prior QA Report path with the Run-start head when one exists,
and records an explicit absence otherwise. The prompt contract now names all
three blocked-cause counts.

Focused checks:

- Red signal: `rtk proxy env GOCACHE=/private/tmp/roundfix-task02-gocache go test ./internal/agent -run '^TestBuildQAPromptCarries(TheSpecContextBundle|ThePreviousReportIdentity)$'` exited 1 before implementation because `QAPromptRequest` had no context or previous-report fields.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task02-gocache go test ./internal/agent -run '^TestBuildQAPromptCarries(TheSpecContextBundle|ThePreviousReportIdentity)$'` exited 0 after implementation.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task02-gocache go test -count=1 ./internal/daemon -run '^(TestAssembleTaskContextBundleReservesExplicitPathsAndCountsOmittedPriorFiles|TestQAGatePromptUsesTaskContextBuilderAndPreviousReportIdentity)$'` exited 0.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task02-gocache go test ./internal/agent -run '^TestBuildQAPrompt'` exited 0.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task02-gocache go test -count=1 ./internal/agent -run '^TestBuild(TaskPrompt|QAPrompt)'` exited 0.
- `rtk git diff --check` exited 0.

Acceptance evidence:

1. `TestBuildQAPromptCarriesTheSpecContextBundle` covers both a populated
   changed-path list and the explicit `Prior changed files: none` form. The
   Daemon integration case proves the QA gate receives paths from the shared
   builder and the existing 200-path limit remains its only limit.
2. `TestBuildQAPromptStatesQAGateContract` passed with
   `rows_blocked_environment`, `rows_blocked_finding`, and
   `rows_blocked_declared`, including the existing pass semantics for each
   cause.
3. The focused `^TestBuildQAPrompt` regression selection passed all existing
   slug, Spec directory, PRD, branch, checkout, and resolved/unresolved Pull
   Request fact cases.
4. Diff inspection found no change to verdict parsing, report selection or
   naming, QA settlement, or Agent execution. The Daemon change only assembles
   and forwards prompt context before the existing QA Agent Session.

The Task's declared `## Verification` commands were not run; the Daemon owns
them after this handoff.
