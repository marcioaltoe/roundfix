---
task: task_04
spec: 0205-an-advisory-judge-for-spec-authoring
status: pending
type: docs
complexity: low
---

# Task 04: The authoring skills run the judge and answer what it raises

## Overview

task_03 ships `roundfix spec judge`, and no authoring workflow runs it. This
Task adds one heading to each of the `write-prd` and `write-techspec` skills
that runs the command at its own stage after the skill's checker step and
tells the authoring model how to answer each raised judgment (ADR-0200). It
is verifiable on its own: both skills and their mirrors carry the heading and
its phrases, and each skill's raised version is recorded.

## Requirements

1. MUST add to `.agents/skills/write-prd/SKILL.md` a new second-level section `## Advisory judgment`, placed after the section that holds the Report step and before `## Anti-patterns`, that tells the authoring model to run `roundfix spec judge <slug> --stage prd` once `roundfix spec check <slug> --stage prd --run-verification` is clean and before reporting.
2. MUST add to `.agents/skills/write-techspec/SKILL.md` the same section, `## Advisory judgment`, at the same place, running `roundfix spec judge <slug> --stage techspec` once `roundfix spec check <slug> --stage techspec --run-verification` is clean.
3. Each section MUST state, in these words or longer sentences that contain them: `answer each raised judgment`, by correcting the artifact and re-running the checker, or by keeping the text and stating in the report why it stands; `The judge is advisory and never a gate`; `a skipped result is neither a failure nor a clean result`, so the model reports the skip reason and continues; and that the model never prints, stores or asks for either Jev key (`ROUNDFIX_OPENROUTER_API_KEY` or `ROUNDFIX_TYPESAFE_API_KEY`).
4. MUST change no other text of either skill, of its references, or of any other skill, and MUST add no text inside the `### QA settlement` section of any skill.
5. MUST raise both version fields of each of the two skills by one patch step from the value on this Task's starting commit, following the version rule in force on that commit, then run `make skills-sync`, record the versions with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`, run `make baseline-digests`, and name in the Result every file those commands rewrote.

## Subtasks

- [ ] Add the heading to `write-prd`.
- [ ] Add the heading to `write-techspec`.
- [ ] Raise both versions, sync the mirrors and record the versions.

## Acceptance Criteria

- [ ] Both skills carry `## Advisory judgment` with the command of their own stage and the three required phrases.
- [ ] `skills/write-prd` and `skills/write-techspec` equal their canonical copies.
- [ ] `TestEveryOwnedSkillVersionIsRecorded` passes with both new versions recorded.

## Context

- interface: `.agents/skills/write-prd/SKILL.md`
- interface: `skills/write-prd/SKILL.md`
- interface: `.agents/skills/write-techspec/SKILL.md`
- interface: `skills/write-techspec/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`

## Verification

- `for pair in ".agents/skills/write-prd/SKILL.md|## Advisory judgment" ".agents/skills/write-prd/SKILL.md|roundfix spec judge <slug> --stage prd" ".agents/skills/write-prd/SKILL.md|answer each raised judgment" ".agents/skills/write-prd/SKILL.md|The judge is advisory and never a gate" ".agents/skills/write-prd/SKILL.md|a skipped result is neither a failure nor a clean result" ".agents/skills/write-techspec/SKILL.md|## Advisory judgment" ".agents/skills/write-techspec/SKILL.md|roundfix spec judge <slug> --stage techspec" ".agents/skills/write-techspec/SKILL.md|answer each raised judgment" ".agents/skills/write-techspec/SKILL.md|The judge is advisory and never a gate" ".agents/skills/write-techspec/SKILL.md|a skipped result is neither a failure nor a clean result"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; diff -r .agents/skills/write-prd skills/write-prd >/dev/null || { printf 'mirror differs: %s\n' skills/write-prd >&2; exit 1; }; diff -r .agents/skills/write-techspec skills/write-techspec >/dev/null || { printf 'mirror differs: %s\n' skills/write-techspec >&2; exit 1; }; out="$(go test -count=1 -v -run "^(TestEveryOwnedSkillVersionIsRecorded|TestAuthorialSkillSync|TestWritePRDProjectConstraints|TestWriteTechSpecProjectConstraints)$" ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestEveryOwnedSkillVersionIsRecorded TestAuthorialSkillSync TestWritePRDProjectConstraints TestWriteTechSpecProjectConstraints; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; make skills-sync-check` — expected: exit 0; before this Task neither skill carries `## Advisory judgment`, so the command fails.

## References

- [_prd.md](_prd.md) — User Stories 1 and 2; Core Feature 11
- [_techspec.md](_techspec.md) — Testing Approach 6; Build Order 4
- ADR-0200; ADR-0189
