---
task: task_05
spec: 0051-doctor-readiness-contract-reconciliation
status: completed
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

- [x] Capture the stale Roundfix Skill wording as the red signal.
- [x] Update the canonical skill with the settled Doctor contract.
- [x] Apply the exact canonical bytes to the shipped copy.
- [x] Prove byte equality, repository skill synchronization, and the
      changed-file allowlist.

## Acceptance Criteria

- [x] Both Roundfix Skill copies describe canonical Repository Skill Set
      missing-root behavior.
- [x] Both copies show mixed remediation as one owned-then-external chain
      joined by `&&`.
- [x] The canonical and shipped copies are byte-identical.
- [x] Repository-owned skill synchronization checks pass.
- [x] Newly changed paths are limited to the two authorized skill files and
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

## Result

The Roundfix Skill now states the exact missing-Git Repository Skill Set
failure and next action. Mixed ownership remediation is one fail-closed shell
chain in Roundfix-owned-then-external order. The canonical and shipped copies
contain the same bytes.

Acceptance criterion evidence:

1. Before the edit, exact-string greps for the canonical missing-Git result
   and full mixed-remediation chain both exited `1`; after the edit, each exact
   string matched the canonical skill.
2. `rtk cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` exited
   `0`, proving both copies contain the same guidance.
3. `rtk make skills-sync-check` passed its four repository-owned skill tests.
4. The changed-file postflight listed only the canonical Roundfix Skill, its
   shipped copy, and this Task file.

Verification:

- `rtk cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` —
  passed.
- `rtk make skills-sync-check` — passed; four tests passed.
- `rtk grep -n '&& bunx skills experimental_install' .agents/skills/roundfix/SKILL.md`
  — passed; matched the fail-closed mixed chain.
- The declared allowlist pipeline exited `0`; Git emitted an fsmonitor IPC
  diagnostic before the pipe. Repeating the same path predicate with
  `rtk git -c core.fsmonitor=false status --porcelain` passed cleanly with no
  out-of-allowlist path.
- `rtk git diff --check` — passed.

Follow-ups: none.
