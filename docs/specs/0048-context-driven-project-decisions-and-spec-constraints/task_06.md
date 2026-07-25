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
8. MUST keep the release-journey fixture hermetic when the changed generated
   guidance causes its repository formatter to run: provision the exact
   fixture-owned formatter dependency before invoking `bunx --no-install`.

## Subtasks

- [x] Add decomposition preconditions to `write-tasks`.
- [x] Add bounded mutation enforcement to `implement-task`.
- [x] Add Project Constraint checks to `qa-gate`.
- [x] Update generated Spec workflow guidance.
- [x] Add active, archived, authorized, and refusal tests.
- [ ] Provision the release-journey fixture's declared formatter dependency.

## Acceptance Criteria

- [x] An active new Spec without complete constraints cannot produce a Task
  Graph.
- [x] A tooling Task without bounded authorization cannot start.
- [x] An authorized tooling Task can change only the listed files.
- [x] QA detects both missing authorization and out-of-scope tooling changes.
- [x] Completed and archived legacy Specs remain byte-identical.
- [ ] The release-journey test passes without relying on a globally installed
  formatter or network access during `bunx --no-install`.

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
- `rtk go test -count=1 ./internal/cli -run '^TestGuidanceCompositionJourney$'`
  — expected: the disposable repository provisions and runs its exact
  formatter dependency hermetically.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goals 3–4; User Stories 4–6; Core Features 12–15.
- `_techspec.md` → Implementation Design: API Contracts; Build Order 5.
- ADR-0077 → downstream Project Constraint enforcement.

## Result

The downstream workflow now fails closed at decomposition, execution, and QA.
Active Specs must carry complete constraint applicability, reasons, and
`docs/agents/` sources. Tooling Tasks remain pending without express bounded
authorization; authorized Tasks preflight every target path and postflight the
actual worktree delta. QA resolves committed changed paths from Git instead of
trusting reported file lists. Generated Spec guidance carries the same rules,
while legacy completed or archived Specs remain exempt from rewriting and Task
Graph dependency and Task-status ownership stay unchanged.

Acceptance evidence:

- Active incomplete Spec refusal:
  `TestProjectConstraintTaskGate` passed and rejects removal of any required
  constraint row, source, refusal, authorization, or ownership clause.
- Tooling start and bounded mutation:
  `TestProjectConstraintImplementationGate` passed and requires authorization
  before `in_progress`, an exact allowlist, per-mutation checks, and a
  changed-file postflight.
- QA authorization and scope:
  `TestProjectConstraintQAGate` passed and requires applicability, source,
  authorization, `git diff-tree` evidence, and failure on missing or
  out-of-scope paths.
- Legacy preservation:
  `TestLegacySpecConstraintExemption` passed across `write-tasks`,
  `implement-task`, `qa-gate`, the generated module, and its formatter fixture.
- Generated identity:
  `TestCatalogCompatibility`, `TestFormatterComposition`,
  `TestPlanDeterminismMatchesMaintainedManagedEntryFixture`, and
  `TestBaselineCompatibilityCorpus` passed after refreshing the maintained
  guide, catalog, and parity identities.

Verification:

- `rtk env GOCACHE=/private/tmp/roundfix-go-cache go test -count=1 ./skills ./internal/baseline -run 'TestProjectConstraintTaskGate|TestProjectConstraintImplementationGate|TestProjectConstraintQAGate|TestLegacySpecConstraintExemption'`
  — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-go-cache go run -buildvcs=false ./cmd/roundfix skills check`
  — passed for every repository-owned workflow skill.
- `rtk env TMPDIR=/private/tmp GOCACHE=/private/tmp/roundfix-go-cache make verify`
  — failed in the pre-existing
  `TestGuidanceCompositionJourney/standard-typescript-monorepo` release
  journey. Its disposable repository invokes
  `bunx --no-install oxfmt@0.59.0` without a local dependency bootstrap, and
  Bun 1.3.14 reports that no existing `oxfmt` binary can be run. The same
  failure reproduced after installing the exact version globally and running
  the focused test outside the sandbox.

Reopened because the generated guidance now makes the release journey exercise
the repository formatter. The fixture must provision its own exact formatter
dependency before `bunx --no-install`; a user-global installation is not valid
acceptance evidence.
