---
task: task_05
spec: 0048-context-driven-project-decisions-and-spec-constraints
status: completed
type: docs
complexity: high
---

# Task 05: Add Project Constraints to authoring

## Overview

Make every new PRD and TechSpec carry a concise, readable snapshot of effective
project decisions and universal clauses. The repo-owned authorial skills refuse
completion until the snapshot is applicable, sourced, and complete.

## Requirements

1. MUST add `Project Constraints` to the PRD and TechSpec templates.
2. MUST snapshot identifier, authentication and HTTP, active ADR, and tooling
   constraints as applicable or not applicable with a reason.
3. MUST cite the operative source paths under `docs/agents/`.
4. MUST require express tooling authorization and bounded files in prose when a
   Spec proposes protected mutation.
5. MUST block `write-prd` and `write-techspec` completion on an absent or
   incomplete section.
6. MUST add no tooling-authorization frontmatter.
7. MUST modify only repository-owned authorial skills.

## Subtasks

- [x] Add the readable PRD constraint template.
- [x] Add the readable TechSpec constraint template.
- [x] Teach both skills to resolve effective constraint sources.
- [x] Add completion and authorization gates.
- [x] Add ownership and artifact-contract tests.

## Acceptance Criteria

- [x] A newly authored PRD and TechSpec contain all four constraint areas and
  source paths.
- [x] Each constraint records applicability and a reason.
- [x] Tooling work cannot conclude authoring without express bounded
  authorization.
- [x] No new authorization frontmatter field appears.
- [x] Upstream-managed skill bytes remain unchanged.

## Context

- instruction: `docs/agents/skill-governance.md`
- instruction: `docs/adr/0077-new-specs-carry-a-readable-project-constraint-snapshot.md`
- interface: `.agents/skills/write-prd/SKILL.md`
- interface: `.agents/skills/write-prd/references/prd-template.md`
- interface: `.agents/skills/write-techspec/SKILL.md`
- interface: `.agents/skills/write-techspec/references/techspec-template.md`

## Verification

- `rtk go test -count=1 ./skills -run 'TestWritePRDProjectConstraints|TestWriteTechSpecProjectConstraints|TestAuthoringConstraintOwnership'` — expected: templates, completion gates, bounded authorization, and skill ownership cases pass.
- `rtk go run -buildvcs=false ./cmd/roundfix skills check` — expected: changed repo-owned skills pass their shipped contract.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goal 3; User Stories 4 and 6; Core Features 12–14.
- `_techspec.md` → Implementation Design: API Contracts; Build Order 5.
- ADR-0077 → readable mandatory Project Constraint snapshots.

## Result

Added the readable Project Constraint contract to both authoring templates and
taught `write-prd` and `write-techspec` to resolve the effective identifier,
authentication and HTTP, active ADR, and tooling constraints from semantic
guides under `docs/agents/`. Both skills now refuse to report completion when a
snapshot is absent or incomplete. Protected tooling work additionally requires
express maintainer authorization and exact bounded files in the Tooling
authority row; authorization remains out of frontmatter.

Acceptance evidence:

- PRD and TechSpec artifact-contract tests require all four constraint rows,
  applicability markers, reasons, and `docs/agents/` source paths.
- Completion-gate assertions require both authoring skills to stop on an
  incomplete snapshot and on protected tooling work without express bounded
  authorization.
- Frontmatter assertions reject any authorization field in either template.
- Ownership assertions prove `write-prd` and `write-techspec` remain
  repository-owned, canonical and shipped bytes match, and the locked
  upstream-managed skill tree retains its prior digest.

Verification:

- `rtk env GOCACHE=/private/tmp/roundfix-task05-go-cache go test -count=1 ./skills -run 'TestWritePRDProjectConstraints|TestWriteTechSpecProjectConstraints|TestAuthoringConstraintOwnership'` — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task05-go-cache go run -buildvcs=false ./cmd/roundfix skills check` — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task05-go-cache make verify` — passed: 2,264 tests in 22 packages, skill synchronization/checks, and build.

Follow-ups: None.
