---
task: task_03
spec: 0217-a-cursor-runtime-to-measure-grok-on
status: completed
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

## Result

Implemented the Roundfix Skill documentation slice for Daemon Verification.
Task status and acceptance checkboxes remain Daemon-owned; no declared
Verification command, commit, push, or pull request was performed.

### Implementation

- Added `### Cursor` to the canonical runtime reference. It documents the
  opt-in `cursor` runtime through `cursor-agent acp` via acpx, verbatim Cursor
  model values with `reasoning_effort: ""`, refusal of non-empty effort, the
  maintainer-owned login check through `cursor-agent status`, and the
  `cursor_login_required` classification when the login is missing.
- Raised both Roundfix Skill front-matter version fields from `0.1.15` to
  `0.1.16`.
- Regenerated the distributed `skills/roundfix` mirror with `make skills-sync`
  and recorded the new `0.1.16` digest in
  `skills/testdata/owned-skill-versions.json`.
- No other skill or `### QA settlement` section was changed.

### Acceptance evidence

1. `rg` found `### Cursor`, `cursor-agent acp`,
   `cursor_login_required`, and `reasoning_effort: ""` in the canonical
   runtime reference. The synchronized mirror contains the same section.
2. `make skills-version-check` exited 0. The required recording command first
   hit the host Go cache permission boundary; rerunning it with
   `GOCACHE=/private/tmp/roundfix-task03-gocache` exited 0 and added the
   `0.1.16` Roundfix record. `make skills-sync` completed successfully, and
   the canonical and mirror diffs show matching skill and runtime-reference
   changes.

### Focused checks

- `git -c core.fsmonitor=false diff --check` exited 0.
- The initial unscoped recording attempt was blocked only by the environment's
  permission error opening the host Go cache; no repository failure was
  observed.
- The authored Verification commands and repository-wide Verification remain
  for the Daemon, as required by the daemon-assigned execution contract.

## Carry-forward provenance

- Source Run: `run_20261002T190644Z_6dc4ccb1092ddf45`
- Source commit: `ed28e3562b57c1a66bbebeffe5ef7b7ea2501b48`
