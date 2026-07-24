---
task: task_05
spec: 0048-context-driven-project-decisions-and-spec-constraints
status: pending
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

- [ ] Add the readable PRD constraint template.
- [ ] Add the readable TechSpec constraint template.
- [ ] Teach both skills to resolve effective constraint sources.
- [ ] Add completion and authorization gates.
- [ ] Add ownership and artifact-contract tests.

## Acceptance Criteria

- [ ] A newly authored PRD and TechSpec contain all four constraint areas and
  source paths.
- [ ] Each constraint records applicability and a reason.
- [ ] Tooling work cannot conclude authoring without express bounded
  authorization.
- [ ] No new authorization frontmatter field appears.
- [ ] Upstream-managed skill bytes remain unchanged.

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
