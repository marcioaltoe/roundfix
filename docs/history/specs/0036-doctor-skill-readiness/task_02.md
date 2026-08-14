---
task: task_02
spec: 0036-doctor-skill-readiness
status: completed
type: docs
complexity: medium
---

# Task 02: Synchronize Doctor skill-readiness guidance

## Overview

Publish the completed Doctor behavior through the canonical glossary and user
documentation while preserving every protected tooling and externally managed
skill path. This slice gives developers the same ownership, failure, and
remediation contract as the shipped CLI; Task 03 owns the isolated Roundfix
Skill update.

## Requirements

1. MUST keep the canonical **Repository Skill Set** term in `CONTEXT.md` and
   update **Doctor Command** to include repository skill readiness.
2. MUST document the `skills: ok|failed` line, blocking exit behavior, owned and
   external authorities, offline/read-only guarantee, and exact update commands
   in README and the Doctor user guide.
3. MUST preserve the authorial ownership boundary: do not edit any protected
   tooling path, externally managed skill, `skills-lock.json`, or
   `skills/recommended.txt`; their content is read-only in this slice.
4. MUST describe that unrelated extra installed skills or lock entries are
   ignored and that Doctor never deletes or updates skills automatically.
5. MUST keep all repository content in English and use only canonical terms
   from `CONTEXT.md`.
6. MUST describe Repository Skill Set readiness after, and independently from,
   the Agent Selection Profile readiness delivered by Spec 0041.

## Subtasks

- [x] Finalize Repository Skill Set glossary and Doctor definition.
- [x] Update README and Doctor command/user guidance.
- [x] Verify required external skill hashes without authorial modification.

## Acceptance Criteria

- [x] A reader can identify all required skill authorities, Doctor outcomes,
      update commands, and the no-network/no-mutation guarantee from supported
      user documentation.
- [x] All required external installed trees hash to their current
      `skills-lock.json` entries.
- [x] No unexpected external authorial edit is introduced, and unrelated extra
      skills are neither flagged nor removed.
- [x] No protected tooling path changes in this Task.
- [x] `CONTEXT.md`, docs, and CLI help use Repository Skill Set, Doctor
      Command, and Roundfix Skill consistently.

## Context

- instruction: `CONTEXT.md`
- instruction: `docs/agents/skill-dispatch.md`
- instruction: `.agents/skills/tech-writer/SKILL.md`
- interface: `README.md`
- interface: `docs/user-guide/commands.md`
- interface: `docs/user-guide/usage.md`
- interface: `skills/skills.go`

## Verification

- `rtk git diff --check` — expected: no whitespace errors in docs or manifests.
- `rtk go test ./internal/cli -run 'Test(CommandUsage|DocumentationContract)' -count=1`
  — expected: supported command and documentation contracts describe the
  shipped Doctor behavior.

## References

- `_prd.md` → Goals 4–6; User Stories 4–5; Core Feature 6; User Experience;
  Non-Goals; Decisions.
- `_techspec.md` → Documentation and skill synchronization; Coverage Map; Build
  Order 3; Risks & Considerations.
- `task_01.md` → the shipped Doctor behavior this guidance must describe.
- `docs/specs/_archived/0041-agent-selection-runtime-readiness/_prd.md` → profile-aware
  Doctor prerequisite and separation of readiness concerns.
- `docs/agents/skill-dispatch.md` → canonical/embedded sync and
  upstream-managed ownership rules.

## Result

Completed after the explicitly authorized external skill refresh repaired the
repository prerequisite outside this Task's authorial slice. This Task's
changes remain limited to the canonical glossary and user documentation.

### Changes

- Defined the Repository Skill Set as required Roundfix-owned and externally
  managed Agent Skills while excluding unrelated extras.
- Defined Repository Skill Set readiness as an independent Doctor Command
  result after Agent Selection Profile Readiness.
- Documented exact `skills: ok` and `skills: failed` outcomes, exit `1`,
  ownership authorities, offline and read-only behavior, ignored extras, and
  both update commands in the README and user guides.

### Verification

- `rtk git diff --check` — passed with no whitespace errors.
- `rtk go test ./internal/cli -run 'Test(CommandUsage|DocumentationContract)' -count=1`
  — passed: 4 tests.
- `rtk go test ./skills -run 'Test(CheckRepositoryIgnoresUnrelatedSkillsAndLockEntries|CheckRepositoryReportsReadyRequiredSetWithoutMutation)' -count=1`
  — passed: 2 tests proving no mutation and ignored unrelated entries.
- CLI-compatible deterministic folder-hash audit — passed: all 25 required
  external installed trees match their `skills-lock.json` `computedHash`
  entries.
- `rtk git diff --name-only -- .agents/skills skills skills-lock.json` —
  passed with empty output; no protected tooling, external skill, recommendation,
  or lock path changed.

### Acceptance evidence

- Supported docs identify the running binary as authority for Roundfix-owned
  skills and `skills-lock.json` `computedHash` values as authority for
  externally managed skills. They include both exact remediation commands and
  the no-network/no-mutation guarantee.
- The fresh compatibility audit proves every required external tree matches its
  authoritative lock entry after the separately committed refresh.
- Fresh skills tests prove unrelated installed skills and lock entries do not
  affect readiness and that the checker does not mutate its repository.
- Git scope inspection proves this Task changed no protected tooling path.
- The glossary, docs, and existing Doctor help use Repository Skill Set,
  Doctor Command, and Roundfix Skill consistently.
