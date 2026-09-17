---
task: task_06
spec: 0140-a-spec-traces-the-promises-it-makes
status: pending
type: test
complexity: low
---

# Task 06: Let the operational fixtures declare their promises

## Overview

The rule this Spec ships reaches the synthetic Specs that the operational
command tests build. Their PRD fixture carries Project Constraints but no
promise declaration, so the QA gate's strict precondition now refuses them and
the implement, attach and agent-selection flows fail with a blocked mechanical
stage. This slice makes those fixtures declare, exactly as an author would.

## Requirements

1. MUST make each synthetic Spec fixture declare its Success Metrics, using the
   explicit none form with a reason that says what the fixture measures.
2. MUST leave the fixtures' Project Constraints, authorization records, Task
   Graphs and seeds otherwise unchanged.
3. MUST NOT change production behavior, a detector, a severity, or any Spec
   artifact of this repository.
4. MUST NOT weaken a test's assertion to make it pass.

## Subtasks

- [ ] Find every fixture PRD the operational command tests write.
- [ ] Add the explicit-none Success Metrics declaration with its reason.
- [ ] Run the affected packages and confirm the flows settle again.

## Acceptance Criteria

- [ ] The implement, attach and agent-selection flows pass with the new rule
      active.
- [ ] The fixture declares its promise the way the rule prescribes, rather than
      suppressing the check.
- [ ] No production file changes.

## Context

- interface: `internal/cli/implement_test.go`

## Verification

- `grep -q "Success Metrics" internal/cli/implement_test.go` — expected: exit 0; the fixture PRD declares its promises. Before this Task the fixture carries no such section.
- `go test -count=1 -run "^TestRunImplementQAVerdictMatrix$" ./internal/cli` — expected: exit 0; the QA verdict matrix settles again. Before this Task it fails at the blocked mechanical stage.
- `go test -count=1 -run "^TestAttachReplaysCompletedSpecRunReadOnly$" ./internal/cli` — expected: exit 0; the attach replay seeds its Run again.

## References

`_prd.md` → Core Feature 1; Declared intentional breaks;
`_techspec.md` → Implementation Design: Declaring a promise; Risks &
Considerations: A gap blocks decomposition; ADR-0156.
