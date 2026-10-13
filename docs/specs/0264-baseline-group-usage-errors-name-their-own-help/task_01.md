---
task: task_01
spec: 0264-baseline-group-usage-errors-name-their-own-help
status: completed
type: docs
complexity: low
---

# Task 01: The roundfix skill describes the bare-group usage error

## Overview

The `roundfix` skill's baseline reference covers the `baseline` Behavior
Surfaces whose source `internal/cli/baseline_profile.go` task_02 changes
(ADR-0256), and `docs/agents/cli.md` requires the skill to change with the
command. This Task adds one sentence on what a bare `baseline` group answers
and records the owned skill's version raise. It is implemented on the tree
after Spec 0260 merges and is verifiable on its own through the phrase, the
mirror and the version record.

## Requirements

1. MUST add to `.agents/skills/roundfix/references/baseline.md`, directly
   after the group-help sentence Spec 0260 adds (the one ending "which lists
   their commands, and exit `0`."), the sentence: "Run bare, or with a
   command they do not have, they exit `2`: the usage error names the group
   once on standard error and points at `roundfix baseline <group> --help`,
   and an unknown `baseline` command points at `roundfix baseline --help`."
2. MUST run `make skills-sync`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   to raise and record the `roundfix` skill version. MUST NOT edit text
   inside any `### QA settlement` section.
3. MUST NOT change the user guide, `docs/references/skill-coverage.json`, any
   other skill, or any sentence Spec 0260 added.

## Subtasks

- [ ] Add the sentence to the canonical baseline reference.
- [ ] Sync the mirror and record the version raise.

## Acceptance Criteria

- [ ] Both copies of the baseline reference carry the sentence.
- [ ] The `SKILL.md` pair matches and the owned-skill version record passes.

## Consumers

- The `roundfix` skill's baseline reference and its owned-skill version record — consumers: `./skills`

## Context

- instruction: `docs/agents/cli.md`
- instruction: `docs/adr/0256-a-release-waits-for-the-skills-that-describe-a-changed-surface.md`
- interface: `.agents/skills/roundfix/references/baseline.md`
- interface: `skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`

## Verification

- `for file in .agents/skills/roundfix/references/baseline.md skills/roundfix/references/baseline.md; do tr -s '[:space:]' ' ' < "$file" | grep -qF -- 'names the group once on standard error' || { printf 'missing phrase in %s\n' "$file" >&2; exit 1; }; done; cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md || exit 1; go test -count=1 -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills || exit 1` — expected: exit 0; before this Task the phrase is absent from both copies.

## References

- [_prd.md](_prd.md) — Goals; Core Feature 4
- [_techspec.md](_techspec.md) — Coverage Map; Build Order 1
- ADR-0256; ADR-0189; ADR-0233

## Result

- Added the required bare or unknown-group usage-error sentence immediately
  after the Spec 0260 group-help sentence in the canonical baseline reference.
- Ran `rtk make skills-sync`, which synchronized the mirror reference.
- Ran `rtk go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`; it passed (`1 passed in 1 packages`) and recorded `roundfix` skill version `0.1.71` in both `SKILL.md` fields and `skills/testdata/owned-skill-versions.json`.
- Acceptance criterion 1 evidence: the required phrase is present in both
  `.agents/skills/roundfix/references/baseline.md` and
  `skills/roundfix/references/baseline.md`; `cmp` confirms the reference
  copies match.
- Acceptance criterion 2 evidence: `cmp` confirms the canonical and mirrored
  `SKILL.md` files match, and the owned-skill version recorder test passed.
  `git diff --check` passed, and the protected `### QA settlement` sections
  remain identical.
- The Task's declared Verification commands remain for Daemon execution.

## Carry-forward provenance

- Source Run: `run_20261013T182909Z_0621af1dca36d211`
- Source commit: `9e97b7c15de62fadfa1080e77629f3755f9a4c47`
