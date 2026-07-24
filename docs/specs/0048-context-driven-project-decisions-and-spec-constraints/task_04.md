---
task: task_04
spec: 0048-context-driven-project-decisions-and-spec-constraints
status: pending
type: backend
complexity: medium
---

# Task 04: Protect tooling authority

## Overview

Ship tooling authority as one universal Normative Clause in the core Baseline.
No Profile, setup preference, narrower guide, or generic implementation request
can authorize or weaken the protected mutation surface.

## Requirements

1. MUST cover creation, editing, renaming, moving, and deletion across every
   confirmed tooling category.
2. MUST cover configurations, scripts, ignore files, plugin declarations, and
   version pins.
3. MUST include the clause unconditionally in every maintained Profile.
4. MUST expose no decision, default, prompt, or exclusion effect for the
   clause.
5. MUST reject any catalog declaration that allows a narrower artifact to
   weaken the clause.
6. MUST include the clause in source accounting and retention contracts.

## Subtasks

- [ ] Add the universal tooling-authority rule.
- [ ] Render it in the core agent guide.
- [ ] Add catalog non-excludability validation.
- [ ] Extend source and retention accounting.
- [ ] Add mutation and no-prompt tests.

## Acceptance Criteria

- [ ] Every maintained Profile renders the complete tooling-authority clause.
- [ ] Baseline human setup adds no tooling preference prompt.
- [ ] A Profile effect cannot exclude or weaken the clause.
- [ ] Every protected tooling artifact category appears in the operative text.
- [ ] Removing the clause or its accounting makes catalog validation fail.

## Context

- instruction: `docs/adr/0077-new-specs-carry-a-readable-project-constraint-snapshot.md`
- interface: `internal/baseline/assets/modules/core.json`
- interface: `internal/baseline/assets/templates/guides/agent-instructions.md`
- interface: `internal/baseline/assets/coverage.json`

## Verification

- `rtk go test -count=1 ./internal/baseline ./internal/cli -run 'TestToolingAuthorityClause|TestToolingAuthorityCannotBeDisabled|TestToolingAuthorityNoPrompt|TestToolingAuthorityAccounting'` — expected: completeness, universality, no-prompt, mutation, and retention cases pass.
- `rtk make verify` — expected: the full repository gate passes.

## References

- `_prd.md` → Goals 3–4; User Story 5; Core Features 10–11 and 13.
- `_techspec.md` → Implementation Design: Data Models; Build Order 4.
- ADR-0077 → readable authorization contract.
