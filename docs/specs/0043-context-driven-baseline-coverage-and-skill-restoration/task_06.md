---
task: task_06
spec: 0043-context-driven-baseline-coverage-and-skill-restoration
status: pending
type: backend
complexity: high
---

# Task 06: Restore external Repository Skill Set members

## Overview

Add the explicit restoration surface for missing or drifted external skills in
the selected Repository Skill Set. The command acquires only the declared
immutable source, previews every directory and lock mutation, requires the
same plan digest for apply, and rolls back the repository on any failed proof.

## Requirements

1. MUST add the non-interactive `restore-skills` command contract, structured
   schema, categorized exits, repeatable skill filter, offline source option,
   and plan confirmation defined by the TechSpec.
2. MUST acquire supported GitHub sources by exact commit with argv-only Git,
   prompting disabled, commit identity verified, and one acquisition per
   unique source/ref pair.
3. MUST verify each acquired source subtree against snapshot authority before
   planning any repository write and MUST never fall back to a branch or
   default revision.
4. MUST preview every created, refreshed, and removed skill path plus the exact
   targeted lock edit under one digest-bound Change Plan.
5. MUST atomically swap staged skill directories and the lock file, preserve
   unrelated skills and lock entries, reject stale confirmation, verify final
   tree digests, and roll back every target on failure.
6. MUST keep persisted lock provenance portable and use the isolated lock
   adapter only on already verified temporary bytes; writes remain disabled if
   its compatibility fixture disagrees with Spec 0036 Task 01.
7. MUST bound acquired file count and bytes, reject traversal, links, devices,
   and unsupported providers, and never execute downloaded skill content.

## Subtasks

- [ ] Add restore selection, JSON/text output, exit, and help contracts.
- [ ] Add exact-commit Git and offline-object-store acquisition adapters.
- [ ] Build complete directory and lock operations under the shared plan
      digest.
- [ ] Add isolated lock compatibility normalization with portable provenance.
- [ ] Add staged directory/lock swap, postwrite proof, and rollback.
- [ ] Add security limits and unsafe-tree rejection before mutation.
- [ ] Add disposable-Git integration flows and injected-boundary failure cases.

## Acceptance Criteria

- [ ] Preview for a drifted nested skill names every file that will be created,
      refreshed, or removed and the exact lock entry that will change.
- [ ] Confirmation restores bytes from the declared full commit, produces the
      expected complete-tree digest, and writes no absolute or machine-local
      lock value.
- [ ] Multiple selected skills from one source/ref use one acquisition while
      retaining separate path operations and digest proofs.
- [ ] Missing confirmation or stale state performs no repository mutation and
      returns the current structured plan; malformed input exits `2`.
- [ ] Missing Git, unreachable commit, commit mismatch, unsupported provider,
      digest mismatch, unsafe tree, size breach, or incompatible lock adapter
      exits non-zero before mutation with a specific next action.
- [ ] An injected directory-swap, lock-write, or postwrite-verification failure
      restores every targeted directory and original lock byte.
- [ ] Unrelated skills and lock entries remain byte-identical after success.
- [ ] A second restoration against the resulting repository has an empty plan,
      final audit exits `0`, and no confirmation is required.
- [ ] Canonical and embedded setup skill trees are synchronized after the
      slice.

## Context

- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `docs/agents/skill-governance.md`
- instruction: `docs/specs/0036-doctor-skill-readiness/task_01.md`
- interface: `.agents/skills/setup-context-driven/scripts/context_setup.py`
- interface: `skills-lock.json`

## Verification

- `rtk python3 -B .agents/skills/setup-context-driven/tests/test_restore_skills.py`
  — expected: preview, exact acquisition, confirmation, portable lock, security,
  rollback, final audit, and idempotency flows pass against disposable local
  Git sources.
- `rtk make verify` — expected: the full repository gate passes with the
  canonical and embedded restoration surface synchronized.

## References

- `_prd.md` → User Story 5; Core Features 8–9; User Experience; Non-Goals;
  Success Metrics.
- `_techspec.md` → Interfaces: SkillSource and SkillLockAdapter; API Contracts:
  restore-skills; Integration Points; Risks & Considerations; Build Order 6.
- Spec 0036 Task 01 → external lock compatibility fixture and Doctor ownership
  boundary.
- `docs/agents/skill-governance.md` → prohibition on authorial changes to
  upstream-managed skill content.
