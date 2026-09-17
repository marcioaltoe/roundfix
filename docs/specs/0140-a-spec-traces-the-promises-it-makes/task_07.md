---
task: task_07
spec: 0140-a-spec-traces-the-promises-it-makes
status: pending
type: test
complexity: medium
---

# Task 07: Let the Daemon fixtures and the archive characterization catch up

## Overview

Two more places embed what this Spec changed. The Daemon's QA-gate tests build
synthetic Specs whose PRDs declare no promise, so the gate's strict precondition
refuses them. The archive characterization pins the corpus golden's recorded
content, which the new codes moved. This slice makes both catch up, and it is
the second and final corrective Task the contract allows.

## Requirements

1. MUST make each synthetic Spec fixture in the Daemon's tests declare its
   Success Metrics, using the explicit none form with a reason that says what
   the fixture measures.
2. MUST re-record the archive characterization against the corpus golden as this
   Spec leaves it, preserving what the characterization is about.
3. MUST NOT change production behavior, a detector, a severity, a verdict rule,
   or any Spec artifact of this repository.
4. MUST NOT weaken or delete an assertion to make a test pass.
5. MUST leave the repository Verification green.

## Subtasks

- [ ] Add the explicit-none declaration to every Daemon fixture PRD that needs
      one.
- [ ] Re-record the archive characterization from the golden this Spec produced.
- [ ] Run the two affected packages and the repository Verification.

## Acceptance Criteria

- [ ] The Daemon's QA-gate and task-cycle flows settle again with the rule
      active.
- [ ] The archive characterization passes and still pins the golden it names.
- [ ] `make verify` exits 0.
- [ ] No production file changes.

## Context

- interface: `internal/daemon/task_engine_test.go`
- interface: `internal/spec/archive_layout_characterization_test.go`

## Verification

- `grep -q "Success Metrics" internal/daemon/task_engine_test.go` — expected: exit 0; the Daemon fixture PRD declares its promises. Before this Task it carries no such section.
- `go test -count=1 -run "^TestTaskCycleQAVerdictMatrixSettlesRunAndCommitsReport$" ./internal/daemon` — expected: exit 0; the QA verdict matrix settles again. Before this Task it fails at the blocked mechanical stage.
- `go test -count=1 -run "^TestArchiveLayoutCharacterizationPinsCorpusGoldenAfterSpec0095$" ./internal/spec` — expected: exit 0; the characterization matches the golden this Spec produced. Before this Task it fails on the recorded content.

## References

`_prd.md` → Core Feature 1; Declared intentional breaks;
`_techspec.md` → Implementation Design: Declaring a promise; Data Models;
Risks & Considerations: A gap blocks decomposition; ADR-0130; ADR-0156.
