---
task: task_07
spec: 0043-context-driven-baseline-coverage-and-skill-restoration
status: completed
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

- [x] Update the canonical setup skill's audit, decision, plan, apply, and
      restoration recipes.
- [x] Update Context-Driven user guidance and maintainer asset documentation.
- [x] Add documentation contract checks for commands, schemas, exits, and
      ownership boundaries.
- [x] Regenerate the embedded setup skill from the canonical source.
- [x] Confirm no upstream-managed skill content changed.

## Acceptance Criteria

- [x] The canonical skill gives an agent a complete non-interactive sequence
      from audit through resolved plan, digest confirmation, final audit, and
      explicit drift restoration.
- [x] User guidance names `setup-context-driven/audit-v1` and
      `setup-context-driven/restore-v1`, all supported exit categories, and the
      exact confirmation behavior.
- [x] Maintainer guidance explains coverage/rule/dispatch/reference and
      snapshot-v2 validation plus the source synchronization boundary.
- [x] Documentation never suggests a generic skill refresh, automatic restore,
      project-specific generated architecture, or extra-skill removal.
- [x] Canonical and embedded setup skills are byte-identical after sanctioned
      synchronization.
- [x] The complete repository verification gate passes with no modified
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

## Result

Updated the canonical setup workflow to carry one non-interactive contract from
local audit and Decision Plan resolution through complete plan review, exact
digest confirmation, apply, final audit, and explicit Repository Skill Set
restoration. The user guide now documents both JSON schemas, plan and finding
fields, exits `0` through `3`, stale-plan behavior, offline acquisition, and the
boundary between network-free documentation setup and Git-backed restoration.

Added maintainer guidance for semantic coverage, rule carriers, skill dispatch,
typed references, snapshot-v2 provenance and complete-tree validation, the Spec
0036 lock compatibility gate, and canonical-to-embedded synchronization. Added
contract tests that fail when commands, schemas, exits, plan fields, prohibited
broad skill-management recipes, or ownership boundaries disappear. Refreshed
the setup skill's repo-owned content digest in all bundled snapshots and ran the
sanctioned synchronization workflow; no upstream-managed skill content changed.

Acceptance evidence:

- Complete agent sequence: `test_operator_docs_publish_schema_exit_and_confirmation_contract`
  passed and verifies audit → resolved plan → apply → final audit → explicit
  restoration ordering.
- User contract: the focused workflow test passed with both schema identifiers,
  every plan field, `--confirm-plan`, `restore-skills`, and exits `0`–`3` present
  in the canonical skill and user guide.
- Maintainer contract: `test_asset_maintenance_doc_publishes_catalog_and_source_boundaries`
  passed for coverage, `requiredRules`, `skillDispatch`, typed references,
  snapshot-v2, `treeDigest`, Spec 0036/Doctor ownership, and synchronization.
- Prohibited guidance: the contract test rejects generic install/update commands;
  inspection confirmed the docs require explicit restoration, forbid branch or
  default-revision fallback, preserve extra skills, and leave project-specific
  architecture repository-owned.
- Distribution: `rtk make skills-sync-check` and
  `rtk diff -r .agents/skills/setup-context-driven skills/setup-context-driven`
  both exited `0` after `rtk make skills-sync`.
- Scope and repository gate: `rtk git status --short` showed changes only in the
  repo-owned setup skill, its embedded copy, this Task, and the Context-Driven
  user guide. `rtk make verify` passed with 124 Python tests, 1,687 Go tests,
  canonical and embedded asset validation, the Roundfix skill check, and build.

Verification:

- `rtk env PYTHONDONTWRITEBYTECODE=1 python3 -B -m unittest discover -s .agents/skills/setup-context-driven/tests -p 'test_workflow.py'` — passed (7 tests).
- `rtk grep -F 'restore-skills' .agents/skills/setup-context-driven/SKILL.md docs/user-guide/context-driven-development.md` — passed; both preview and confirmed recipes are present.
- `rtk make skills-sync-check` — passed.
- `rtk make verify` — passed after rerunning with access to the Go build cache; 124 Python tests and 1,687 Go tests passed, asset/skill checks passed, and the build completed.
