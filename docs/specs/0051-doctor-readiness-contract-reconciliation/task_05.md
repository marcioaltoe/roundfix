---
task: task_05
spec: 0051-doctor-readiness-contract-reconciliation
status: pending
type: docs
complexity: low
---

# Task 05: Synchronize Roundfix Doctor guidance

## Overview

Align the canonical Roundfix Skill and its shipped copy with the settled
Doctor contract after code and user guidance are verified. This dedicated
tooling slice has an exact three-path mutation allowlist and changes no runtime
code, upstream-managed skill, lock authority, or generated Baseline artifact.

## Requirements

1. MUST document the canonical missing-Git Repository Skill Set result and the
   `&&`-joined mixed remediation chain in the Roundfix Skill.
2. MUST keep `.agents/skills/roundfix/SKILL.md` and
   `skills/roundfix/SKILL.md` byte-identical.
3. MUST limit every Task mutation to
   `.agents/skills/roundfix/SKILL.md`,
   `skills/roundfix/SKILL.md`, and
   `docs/specs/0051-doctor-readiness-contract-reconciliation/task_05.md`.
4. MUST NOT run a broad mutation that rewrites unrelated owned skills.
5. MUST NOT edit an upstream-managed skill, `skills-lock.json`,
   `skills/recommended.txt`, a Baseline artifact, `.coderabbit.yaml`, or
   `.roundfixrc.yml`.

## Subtasks

- [ ] Capture the stale Roundfix Skill wording as the red signal.
- [ ] Update the canonical skill with the settled Doctor contract.
- [ ] Apply the exact canonical bytes to the shipped copy.
- [ ] Prove byte equality, repository skill synchronization, and the
      changed-file allowlist.

## Acceptance Criteria

- [ ] Both Roundfix Skill copies describe canonical Repository Skill Set
      missing-root behavior.
- [ ] Both copies show mixed remediation as one owned-then-external chain
      joined by `&&`.
- [ ] The canonical and shipped copies are byte-identical.
- [ ] Repository-owned skill synchronization checks pass.
- [ ] Newly changed paths are limited to the two authorized skill files and
      this Task file.

## Context

- instruction: `docs/agents/agent-instructions.md`
- instruction: `docs/agents/skill-dispatch.md`
- instruction: `.agents/skills/tech-writer/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `docs/user-guide/commands.md`

## Verification

- `rtk cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` — expected:
  canonical and shipped files are byte-identical.
- `rtk make skills-sync-check` — expected: repository-owned skill pairs are
  synchronized.
- `rtk grep -n '&& bunx skills experimental_install' .agents/skills/roundfix/SKILL.md` — expected: canonical guidance uses the fail-closed chain.
- `rtk git status --porcelain | rtk awk '{path=substr($0,4); if (path != ".agents/skills/roundfix/SKILL.md" && path != "skills/roundfix/SKILL.md" && path != "docs/specs/0051-doctor-readiness-contract-reconciliation/task_05.md") {print; bad=1}} END {exit bad}'` — expected: no out-of-allowlist path.

## References

- `_prd.md` → Core Features 6; Non-Goals; Success Metrics.
- `_techspec.md` → Documentation and skill synchronization; Build Order 5.
