---
task: task_07
spec: 0043-context-driven-baseline-coverage-and-skill-restoration
status: pending
type: docs
complexity: medium
---

# Task 07: Synchronize setup workflow guidance

## Overview

Publish one consistent operator contract for coverage, reference audit,
digest-bound apply, and explicit Repository Skill Set restoration. The slice
updates the canonical skill and maintainer/user guidance, then regenerates the
embedded copy without modifying upstream-managed skill content.

## Requirements

1. MUST document the semantic coverage and ownership boundaries, universal
   Verification decision, typed-reference findings, complete plan fields, and
   confirmation flow.
2. MUST document external snapshot provenance, complete-tree audit,
   `restore-skills` preview/apply recipes, structured remediation, failure
   exits, offline source option, and absence of branch fallback.
3. MUST explain that audit and documentation apply remain local and
   network-free while explicit restoration owns Git acquisition.
4. MUST document the Spec 0036 lock compatibility gate without moving Doctor
   behavior into setup guidance.
5. MUST update the canonical repo-owned skill first and regenerate its embedded
   copy through the repository synchronization workflow.

## Subtasks

- [ ] Update the canonical setup skill's audit, decision, plan, apply, and
      restoration recipes.
- [ ] Update Context-Driven user guidance and maintainer asset documentation.
- [ ] Add documentation contract checks for commands, schemas, exits, and
      ownership boundaries.
- [ ] Regenerate the embedded setup skill from the canonical source.
- [ ] Confirm no upstream-managed skill content changed.

## Acceptance Criteria

- [ ] The canonical skill gives an agent a complete non-interactive sequence
      from audit through resolved plan, digest confirmation, final audit, and
      explicit drift restoration.
- [ ] User guidance names `setup-context-driven/audit-v1` and
      `setup-context-driven/restore-v1`, all supported exit categories, and the
      exact confirmation behavior.
- [ ] Maintainer guidance explains coverage/rule/dispatch/reference and
      snapshot-v2 validation plus the source synchronization boundary.
- [ ] Documentation never suggests a generic skill refresh, automatic restore,
      project-specific generated architecture, or extra-skill removal.
- [ ] Canonical and embedded setup skills are byte-identical after sanctioned
      synchronization.
- [ ] The complete repository verification gate passes with no modified
      upstream-managed skill file.

## Context

- instruction: `docs/agents/skill-governance.md`
- interface: `.agents/skills/setup-context-driven/SKILL.md`
- interface: `docs/user-guide/context-driven-development.md`
- interface: `Makefile`

## Verification

- `rtk grep -F 'restore-skills' .agents/skills/setup-context-driven/SKILL.md docs/user-guide/context-driven-development.md`
  — expected: canonical agent and user workflows both document the explicit
  restoration command.
- `rtk make skills-sync-check` — expected: every repo-owned canonical skill
  matches its embedded copy and recommended external names remain synchronized.
- `rtk make verify` — expected: documentation contracts, setup macro flows,
  skill checks, tests, formatting, and build all pass.

## References

- `_prd.md` → Core Feature 10; User Experience; Non-Goals; Success Metrics.
- `_techspec.md` → System Architecture: canonical distribution; Integration
  Points; Risks & Considerations; Build Order 7.
- `docs/agents/skill-governance.md` → authorial ownership and synchronization
  contract.
