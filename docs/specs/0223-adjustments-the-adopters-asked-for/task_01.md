---
task: task_01
spec: 0223-adjustments-the-adopters-asked-for
status: pending
type: docs
complexity: low
---

# Task 01: The qa-gate skill, the Roundfix Skill and the guides describe the adjustments

## Overview

This Task writes every piece of text the Spec changes outside Go code: the
`qa-gate` skill's rule for a report without row inputs, the Roundfix Skill's
release reference and the release runbook for the release plan's two checks,
and the Context-Driven Development guide for the optional branch prefix, as
the TechSpec states them. The repository's skill-sync rule requires a change
to CLI behavior to ship the Roundfix Skill update, and each owned skill whose
text changes raises its version.

## Requirements

1. MUST add to `.agents/skills/qa-gate/SKILL.md`, section "Row input
   declaration", directly after the paragraph that ends "Do not add inputs
   after execution to make carry-forward eligible.", the paragraph of the
   TechSpec's "The QA gate paragraph" verbatim. This answers the Backlog Entry
   of 2026-10-03, "A reopened QA gate has no path for a report written before
   row inputs".
2. MUST describe in `.agents/skills/roundfix/references/release.md`, after the
   paragraph that says the command mutates no repository or release state, that
   a range plan also reports the skills and baseline checks read-only after its
   `Next action:` line: one `skills:` line with the comparison `roundfix
   doctor` makes and one `baseline:` line with what `roundfix baseline update
   --no-skills` would report; that the JSON carries both under `checks`; that a
   line that is not `ok` or `current` ends with the next action "complete the
   skills and guides check in the release runbook before the release Pull
   Request"; and that the checks never change the decision state, the proposed
   version or the exit code. This answers the Backlog Entry of 2026-09-30, "The
   release plan does not report the skill and guide checks".
3. MUST say in `docs/user-guide/release-runbook.md`, in step 1 of "Cutting a
   release", that the plan reports the skills and baseline checks read-only in
   its `skills:` and `baseline:` lines and that they never change its
   decision.
4. MUST add to `docs/user-guide/context-driven-development.md` a "Branch
   prefix" row to the suggested-values table, `<type>/` marked optional, and a
   paragraph before "The frontend layout is optional." stating that the branch
   prefix is optional, that while none is recorded ("No branch prefix is
   recorded") the agent-instructions guide names new work branches
   `<type>/<description>` from the work's Conventional Commit type, and that a
   recorded prefix keeps its sentence and an update never asks for it. This
   answers the Backlog Entry of 2026-10-03, "A repository cannot leave the
   branch prefix to the commit-type rule".
5. MUST raise the `qa-gate` skill's version from 0.0.7 to 0.0.8 and the
   Roundfix Skill's version by one patch level from the value on the tree the
   Task starts from, each in both front-matter fields, run `make skills-sync`
   so the mirrors equal their canonical files, and re-record the versions with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`.
6. MUST NOT edit the `### QA settlement` section of any skill, any other
   section of the `qa-gate` skill, any other command reference or guide,
   `CONTEXT.md` or `.roundfixrc.yml`.

## Subtasks

- [ ] Add the re-execution paragraph to the `qa-gate` skill.
- [ ] Describe the two checks in the release reference and the runbook.
- [ ] Describe the optional branch prefix in the Context-Driven Development guide.
- [ ] Raise both versions, sync the mirrors and record the versions.

## Acceptance Criteria

- [ ] The `qa-gate` skill and its mirror carry the re-execution paragraph.
- [ ] The release reference, its mirror and the runbook carry "reports the
      skills and baseline checks read-only".
- [ ] The Context-Driven Development guide carries "No branch prefix is
      recorded" and "Conventional Commit type".
- [ ] Each mirror equals its canonical file, and both raised versions are
      recorded.

## Context

- interface: `.agents/skills/qa-gate/SKILL.md`
- interface: `skills/qa-gate/SKILL.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/release.md`
- interface: `skills/roundfix/references/release.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/release-runbook.md`
- interface: `docs/user-guide/context-driven-development.md`
- instruction: `docs/adr/0228-an-unrecorded-branch-prefix-follows-the-commit-type-rule.md`

## Verification

- `tr -s '[:space:]' ' ' < .agents/skills/qa-gate/SKILL.md | grep -qF -- "A gate reopened over a report written before rows declared inputs" || { printf 'missing phrase in %s: %s\n' .agents/skills/qa-gate/SKILL.md "A gate reopened over a report written before rows declared inputs" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/qa-gate/SKILL.md | grep -qF -- "declares its inputs for that pass" || { printf 'missing phrase in %s: %s\n' .agents/skills/qa-gate/SKILL.md "declares its inputs for that pass" >&2; exit 1; }; tr -s '[:space:]' ' ' < skills/qa-gate/SKILL.md | grep -qF -- "A gate reopened over a report written before rows declared inputs" || { printf 'missing phrase in %s: %s\n' skills/qa-gate/SKILL.md "A gate reopened over a report written before rows declared inputs" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/release.md | grep -qF -- "reports the skills and baseline checks read-only" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/release.md "reports the skills and baseline checks read-only" >&2; exit 1; }; tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/release.md | grep -qF -- "complete the skills and guides check in the release runbook before the release Pull Request" || { printf 'missing phrase in %s: %s\n' .agents/skills/roundfix/references/release.md "complete the skills and guides check in the release runbook before the release Pull Request" >&2; exit 1; }; tr -s '[:space:]' ' ' < skills/roundfix/references/release.md | grep -qF -- "reports the skills and baseline checks read-only" || { printf 'missing phrase in %s: %s\n' skills/roundfix/references/release.md "reports the skills and baseline checks read-only" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/release-runbook.md | grep -qF -- "reports the skills and baseline checks read-only" || { printf 'missing phrase in %s: %s\n' docs/user-guide/release-runbook.md "reports the skills and baseline checks read-only" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/context-driven-development.md | grep -qF -- "No branch prefix is recorded" || { printf 'missing phrase in %s: %s\n' docs/user-guide/context-driven-development.md "No branch prefix is recorded" >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/context-driven-development.md | grep -qF -- "Conventional Commit type" || { printf 'missing phrase in %s: %s\n' docs/user-guide/context-driven-development.md "Conventional Commit type" >&2; exit 1; }` — expected: exit 0; before this Task none of these phrases is in these files, so the command fails.
- `cmp .agents/skills/qa-gate/SKILL.md skills/qa-gate/SKILL.md && cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md && cmp .agents/skills/roundfix/references/release.md skills/roundfix/references/release.md && tr -s '[:space:]' ' ' < skills/qa-gate/SKILL.md | grep -qF -- "declares its inputs for that pass" && out="$(go test -count=1 -v -run "^(TestEveryOwnedSkillVersionIsRecorded)$" ./skills 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestEveryOwnedSkillVersionIsRecorded; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the `qa-gate` mirror does not carry the new phrase, so the command fails; after it the mirrors equal their canonical files and both raised versions and their content digests are recorded.

## References

- `_prd.md` → Core Features 4 and 6; User Stories 3 and 4; Success Metric 4
- `_techspec.md` → The QA gate paragraph; Vocabulary Contract; API Contract 1; API Contract 2; API Contract 5; Build Order 1
- ADR-0187; ADR-0189; ADR-0228
