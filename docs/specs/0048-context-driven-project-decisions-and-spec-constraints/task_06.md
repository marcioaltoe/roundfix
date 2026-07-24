---
task: task_06
spec: 0048-context-driven-project-decisions-and-spec-constraints
status: pending
type: docs
complexity: high
---

# Task 06: Enforce Project Constraints downstream

## Overview

Carry the approved Project Constraint snapshot into Task decomposition,
execution, and final QA. Tooling mutations remain blocked outside the exact
files authorized in the active Spec.

## Requirements

1. MUST make `write-tasks` refuse a non-archived Spec whose PRD or TechSpec
   lacks complete Project Constraints.
2. MUST make `write-tasks` refuse tooling Tasks when express authorization and
   bounded files are absent.
3. MUST make `implement-task` stop before any tooling mutation outside the
   approved bounded files.
4. MUST make `qa-gate` verify applicability, source paths, authorization, and
   actual changed-file scope.
5. MUST exempt existing completed or archived Specs from forced rewriting.
6. MUST update generated Spec workflow guidance with the same contract.
7. MUST preserve Task status and dependency ownership boundaries.

## Subtasks

- [ ] Add decomposition preconditions to `write-tasks`.
- [ ] Add bounded mutation enforcement to `implement-task`.
- [ ] Add Project Constraint checks to `qa-gate`.
- [ ] Update generated Spec workflow guidance.
- [ ] Add active, archived, authorized, and refusal tests.

## Acceptance Criteria

- [ ] An active new Spec without complete constraints cannot produce a Task
  Graph.
- [ ] A tooling Task without bounded authorization cannot start.
- [ ] An authorized tooling Task can change only the listed files.
- [ ] QA detects both missing authorization and out-of-scope tooling changes.
- [ ] Completed and archived legacy Specs remain byte-identical.

## Context

- instruction: `docs/adr/0077-new-specs-carry-a-readable-project-constraint-snapshot.md`
- interface: `.agents/skills/write-tasks/SKILL.md`
- interface: `.agents/skills/implement-task/SKILL.md`
- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `internal/baseline/assets/modules/spec-workflow.json`
- interface: `internal/baseline/assets/templates/guides/spec-routing.md`

## Verification

- `rtk go test -count=1 ./skills ./internal/baseline -run 'TestProjectConstraintTaskGate|TestProjectConstraintImplementationGate|TestProjectConstraintQAGate|TestLegacySpecConstraintExemption'` — expected: decomposition, execution, QA, legacy, and ownership boundaries pass.
- `rtk go run -buildvcs=false ./cmd/roundfix skills check` — expected: all changed repo-owned workflow skills pass.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goals 3–4; User Stories 4–6; Core Features 12–15.
- `_techspec.md` → Implementation Design: API Contracts; Build Order 5.
- ADR-0077 → downstream Project Constraint enforcement.
