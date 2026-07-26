---
task: task_02
spec: 0036-doctor-skill-readiness
status: pending
type: docs
complexity: medium
---

# Task 02: Synchronize Doctor skill-readiness guidance

## Overview

Publish the completed Doctor behavior through the canonical glossary, user
documentation, and Roundfix-owned skill bundle while preserving externally
managed skill updates exactly as produced by the upstream installer. This slice
ensures a developer, Supervisor, or Agent receives the same ownership,
failure, and remediation contract as the shipped CLI.

## Requirements

1. MUST keep the canonical **Repository Skill Set** term in `CONTEXT.md` and
   update **Doctor Command** to include repository skill readiness.
2. MUST document the `skills: ok|failed` line, blocking exit behavior, owned and
   external authorities, offline/read-only guarantee, and exact update commands
   in README and the Doctor user guide.
3. MUST update the canonical `.agents/skills/roundfix` instructions and OpenAI
   manifest wherever Doctor checks or recovery behavior are enumerated.
4. MUST regenerate the embedded `skills/roundfix` copy with `make skills-sync`
   and extend the shipped skill contract checks so future CLI/skill drift fails
   verification.
5. MUST preserve the authorial ownership boundary: do not manually edit any
   externally managed skill, and retain only upstream changes produced by the
   approved `make skills-update` flow.
6. MUST keep `skills-lock.json`, the installed required external skill trees,
   and `skills/recommended.txt` synchronized after the update.
7. MUST describe that unrelated extra installed skills or lock entries are
   ignored and that Doctor never deletes or updates skills automatically.
8. MUST keep all repository content in English and use only canonical terms
   from `CONTEXT.md`.
9. MUST describe Repository Skill Set readiness after, and independently from,
   the Agent Selection Profile readiness delivered by Spec 0041.

## Subtasks

- [ ] Finalize Repository Skill Set glossary and Doctor definition.
- [ ] Update README and Doctor command/user guidance.
- [ ] Update the canonical Roundfix Skill and manifest.
- [ ] Regenerate embedded owned skills and recommended external names.
- [ ] Extend skill contract assertions for the new Doctor wording.
- [ ] Verify upstream-managed skill changes and lock hashes without authorial
      modification.
- [ ] Run focused skill synchronization and documentation checks.

## Acceptance Criteria

- [ ] A reader can identify all required skill authorities, Doctor outcomes,
      update commands, and the no-network/no-mutation guarantee from supported
      user documentation.
- [ ] The canonical Roundfix Skill tells an Agent to surface a failed skill
      check and its remediation before work continues, without claiming Doctor
      performs the update.
- [ ] `.agents/skills/roundfix` and `skills/roundfix` are byte-identical after
      synchronization, and `roundfix skills check` enforces the new wording.
- [ ] All required external installed trees hash to their current
      `skills-lock.json` entries, including the four upstream updates discovered
      during specification work.
- [ ] No unexpected external authorial edit is introduced, and unrelated extra
      skills are neither flagged nor removed.
- [ ] `CONTEXT.md`, docs, CLI help, and skill text use Repository Skill Set,
      Doctor Command, and Roundfix Skill consistently.

## Context

- instruction: `CONTEXT.md`
- instruction: `docs/agents/skill-dispatch.md`
- docs: `README.md`
- docs: `docs/user-guide/commands.md`
- docs: `docs/user-guide/usage.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/agents/openai.yaml`
- interface: `skills/skills.go`
- interface: `Makefile`
- data: `skills-lock.json`

## Verification

- `rtk make skills-sync` — expected: embedded Roundfix-owned skills and
  recommended external names regenerate from canonical local authorities.
- `rtk make skills-sync-check` — expected: no canonical/embedded owned-skill or
  recommendation drift.
- `rtk go run -buildvcs=false ./cmd/roundfix skills check` — expected: shipped
  skill contracts pass with the new Doctor wording.
- `rtk git diff --check` — expected: no whitespace errors in docs or manifests.
- `rtk make verify` — expected: formatting, Go tests, setup-context tests, skill
  synchronization, shipped skill validation, and build all pass.

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

Pending.
