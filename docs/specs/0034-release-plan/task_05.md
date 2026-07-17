---
task: task_05
spec: 0034-release-plan
status: pending
type: docs
complexity: medium
---

# Task 05: Require Release Plan in release guidance

## Overview

Make the implemented Release Plan Command the mandatory first step for maintainers and Agents without changing the tag-triggered publication workflow. This slice synchronizes the release runbook, root instructions, user-facing command index, and owned Roundfix skill copies under an executable documentation contract.

## Requirements

1. MUST require `roundfix release plan` before changelog, version-file, tag, push, package, asset, or GitHub Release mutation in maintainer and Agent guidance.
2. MUST state that a generic release request authorizes only a conclusive patch plan and that minor, major, or version-zero breaking plans require explicit human approval.
3. MUST explain manual classification with `--impact` plus `--reason` and state that classification does not itself approve the resulting version.
4. MUST preserve the existing tag validation, artifact agreement, npm publication, asset upload, Upgrade Command, and GitHub Release workflow after approval.
5. MUST add the command to the user-facing command index and keep terminology aligned with `CONTEXT.md` and ADR-0048.
6. MUST update the canonical Roundfix skill first and regenerate its embedded copy through the repository skill-sync workflow.
7. MUST add or update documentation contract tests so command help, runbook, Agent pointer, and both skill copies cannot drift silently.

## Subtasks

- [ ] Add the mandatory planning step and approval boundaries to the release runbook.
- [ ] Add the root Agent pointer for release work.
- [ ] Add `release plan` to the user-facing command documentation.
- [ ] Update the canonical Roundfix skill command recipe.
- [ ] Regenerate the embedded Roundfix skill copy.
- [ ] Pin the documentation and skill-sync contract with tests.

## Acceptance Criteria

- [ ] Maintainer and Agent instructions start release work with the read-only Release Plan Command.
- [ ] Patch, minor, major, version-zero breaking, and manual-classification approval boundaries are explicit and consistent.
- [ ] Existing publication steps remain intact and occur only after the plan's required decision is satisfied.
- [ ] Root help, command help, runbook, command index, root Agent pointer, and both Roundfix skill copies use the same canonical terms and command shape.
- [ ] Canonical and embedded Roundfix skills are byte-identical after regeneration.

## Context

- instruction: `docs/agents/skill-governance.md`
- instruction: `docs/adr/0048-release-planning-is-read-only-and-confirmation-gated.md`
- interface: `docs/user-guide/release-runbook.md`
- interface: `docs/user-guide/usage.md`
- interface: `AGENTS.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`

## Verification

- `grep -F 'roundfix release plan' docs/user-guide/release-runbook.md && grep -F 'roundfix release plan' docs/user-guide/usage.md && grep -F 'roundfix release plan' AGENTS.md && grep -F 'roundfix release plan' .agents/skills/roundfix/SKILL.md` — expected: every required durable guidance surface names the command.
- `cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` — expected: canonical and embedded Roundfix skills are byte-identical.
- `make skills-sync-check` — expected: the owned skill bundle is synchronized.
- `go test ./internal/cli -run 'TestReleasePlanDocumentationContract' -count=1` — expected: root help, command help, runbook, Agent pointer, and skill contract pass.

## References

- `_prd.md` → Goal 5; User Story 2; Core Feature 9; Decisions.
- `_techspec.md` → System Architecture: documentation and skill surfaces; Integration Points; Testing Approach; Build Order 5-6.
- ADR-0048 → Release planning is read-only and confirmation-gated.
