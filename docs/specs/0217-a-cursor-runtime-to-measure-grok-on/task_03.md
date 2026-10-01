---
task: task_03
spec: 0217-a-cursor-runtime-to-measure-grok-on
status: pending
type: docs
complexity: low
---

# Task 03: The Roundfix Skill describes the Cursor runtime

## Overview

Tasks 01 and 02 change CLI behavior: a fourth runtime, a model-value rule and
a login check, and they document both in the guides. The repository's skill
sync rule requires the Roundfix Skill to match the shipped behavior. This
Task adds a `### Cursor` section to the skill's runtime reference and raises
the skill's version (ADR-0189).

## Requirements

1. MUST add a `### Cursor` section to
   `.agents/skills/roundfix/references/runtime.md` that states: `cursor` is
   opt-in and reached as `cursor-agent acp` through acpx; a Cursor selection
   names the advertised model value verbatim with `reasoning_effort: ""`, and
   a non-empty effort is refused; the login is the maintainer's own, checked
   with `cursor-agent status` and never performed, and a missing one is
   `cursor_login_required`.
2. MUST raise both version fields of `.agents/skills/roundfix/SKILL.md`, run
   `make skills-sync`, and re-record the version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
3. MUST NOT change the `### QA settlement` section of any skill, or any
   other skill.

## Subtasks

- [ ] Write the skill's `### Cursor` section and raise its version.
- [ ] Sync the mirror and re-record the version.

## Acceptance Criteria

- [ ] The skill reference carries the phrases the Verification names.
- [ ] The mirror equals the canonical skill, and the recorded version
      matches the skill's content.

## Context

- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/runtime.md`
- interface: `skills/roundfix/references/runtime.md`
- interface: `skills/testdata/owned-skill-versions.json`
- instruction: `docs/adr/0189-an-owned-skills-version-names-its-content-and-its-minimum-is-the-bundle.md`

## Verification

- `for phrase in '### Cursor' 'cursor-agent acp' 'cursor_login_required' 'reasoning_effort: ""'; do tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/runtime.md | grep -qF -- "$phrase" || { printf 'missing phrase in the runtime reference: %s\n' "$phrase" >&2; exit 1; }; done; diff -r .agents/skills/roundfix skills/roundfix >/dev/null || { printf 'mirror differs: skills/roundfix\n' >&2; exit 1; }; out="$(go test -count=1 -v -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; printf '%s\n' "$out" | grep -q -- '--- PASS: TestEveryOwnedSkillVersionIsRecorded' || { printf '%s\n' "$out"; exit 1; }; make skills-sync-check` — expected: exit 0; before this Task the runtime reference has no `### Cursor` section, so the command fails.

## References

- [_prd.md](_prd.md) — Core Feature 4; User Experience
- [_techspec.md](_techspec.md) — The Cursor selection rule; The login check; Testing Approach 4; Build Order 3
- ADR-0217; ADR-0187; ADR-0189
