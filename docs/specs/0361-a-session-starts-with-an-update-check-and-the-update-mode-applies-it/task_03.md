---
task: task_03
spec: 0361-a-session-starts-with-an-update-check-and-the-update-mode-applies-it
status: completed
type: docs
complexity: medium
---

# Task 03: The Roundfix skill's update mode applies the plan it showed

## Overview

The Roundfix skill gains a fourth Skill Mode, `update`, read from the request
text in every harness like the other three. It runs the Update Check with
`--refresh`, upgrades the binary when a release is newer, previews this
repository's Baseline plan with the binary it then has and applies exactly that
plan by its Plan Digest, asking before history units, Retired Skill deletions
or uncommitted changes, and never commits. The skill's Modes paragraph, its
`modes` reference, its contract tests and the glossary's **Skill Mode** entry
say so.

## Requirements

1. MUST change `.agents/skills/roundfix/SKILL.md` as `_techspec.md` → API
   Contract 8 says, word for word: the Modes paragraph, the `--update` bullet
   after the `--unattended` bullet, the frontmatter `description` naming "its
   `--follow-up`, `--status`, `--unattended` and `--update` Skill Modes", and
   the `modes` index row naming the `--update` Skill Mode beside the other
   three.
2. MUST change `.agents/skills/roundfix/references/modes.md` as API Contract 8
   says, word for word: "How a mode is chosen" names `follow-up`, `status`,
   `unattended` or `update`; the follow-up command sentence becomes the two
   sentences given; and a new last section `## --update` holds the opening
   line and the five numbered steps of the fenced block. MUST NOT change any
   other section, and MUST NOT add or change text inside any `### QA
   settlement` section.
3. MUST change `TestTheRoundfixSkillReadsItsModeFromTheRequest` in
   `skills/roundfix_modes_test.go` to pin the new Modes paragraph, include
   `--update` among the modes the Modes section names, and require the phrase
   "`follow-up`, `status`, `unattended` or `update`" in "How a mode is
   chosen"; and add `TestTheUpdateModeAppliesOnlyThePlanItShowed` in the same
   file, over both copies (`roundfixLayoutSources`), requiring the `--update`
   section to open with "never in a Daemon-assigned Run", to hold items 1 to
   5, and, by item: 1 `roundfix update check --refresh --format json`; 2
   `roundfix upgrade` and "without applying a Baseline plan"; 3 `roundfix
   baseline update --repo <root> --format json` and "never pass
   `--adopt-suggested` unasked"; 4 `--confirm-plan <planDigest>`,
   `history.units`, `skills.retired`, `git status --porcelain` and "apply no
   Baseline plan over uncommitted changes"; 5 "never commits, pushes or opens
   a pull request".
4. MUST change through `domain-modeling` the `CONTEXT.md` entry **Skill Mode**
   to the text of API Contract 10, word for word, adding no other capitalized
   phrase.
5. MUST run `make skills-sync` to copy the canonical files to
   `skills/roundfix/`, then raise the Roundfix skill version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
   as the Task's last step; never hand-edit a mirror, a version field or the
   record. Of the Governed Paths the Task may change only
   `.agents/skills/roundfix/SKILL.md`,
   `.agents/skills/roundfix/references/modes.md` and
   `skills/roundfix/SKILL.md`.

## Subtasks

- [ ] Change the Modes paragraph, bullet, description and index row.
- [ ] Change "How a mode is chosen", the command sentence, and add the `--update` section.
- [ ] Change and add the modes contract tests.
- [ ] Change the glossary entry, sync and record the skill version.

## Acceptance Criteria

- [ ] `update` and `--update` select the mode in every harness, as the three other modes do.
- [ ] The `--update` section checks, upgrades, previews and applies exactly the previewed digest, asks before deletions and uncommitted changes, and never commits.
- [ ] Both copies of the skill, the contract tests and `CONTEXT.md` agree.

## Context

- instruction: `docs/adr/0371-a-session-starts-with-an-update-check-and-the-skill-applies-the-update.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- instruction: `.agents/skills/roundfix/references/baseline.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/modes.md`
- interface: `skills/roundfix/references/modes.md`
- interface: `skills/roundfix_modes_test.go`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `CONTEXT.md`

## Consumers

- The Roundfix skill's Modes paragraph and `modes` reference — consumers: `./skills`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheUpdateModeAppliesOnlyThePlanItShowed|TestTheRoundfixSkillReadsItsModeFromTheRequest)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheUpdateModeAppliesOnlyThePlanItShowed TestTheRoundfixSkillReadsItsModeFromTheRequest; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./skills` — expected: exit 0. Before this Task the new test does not exist, so its pass line is missing; afterwards both copies hold the `update` mode as specified and the skills suite, the recorded version included, passes.

- `for f in .agents/skills/roundfix/references/modes.md skills/roundfix/references/modes.md; do grep -q '^## --update$' "$f" || { echo "missing update section: $f" >&2; exit 1; }; tr -s '[:space:]' ' ' < "$f" | grep -q -- '--confirm-plan <planDigest> --format json.. Ask first' || { echo "missing apply step: $f" >&2; exit 1; }; done; tr -s '[:space:]' ' ' < CONTEXT.md | grep -q -- 'follow-up., .status., .unattended. or .update., with or without a leading' || { echo 'Skill Mode entry unchanged' >&2; exit 1; }; grep -F -- '**Skill Mode**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Skill Mode\*\*(:| —( |$))' || { echo 'missing term: Skill Mode' >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task neither copy has an `--update` section and the glossary names three modes; afterwards both copies apply the previewed digest, the glossary names `update`, and the documentation and glossary contracts pass.

## References

- `_prd.md` → Goal 3; Core Feature 3; Core Feature 5; Success Metric 5; Success Metric 7; Glossary
- `_techspec.md` → Current behavior 6; API Contract 8; API Contract 10; Declared breaks; Testing Approach; Build Order 3
- ADR-0371; ADR-0320

## Result

Skills loaded: implement-task, domain-modeling, tech-writer, testing-boss, golang-testing, coding-guidelines, evidence-gate, systematic-debugging

Implemented the fourth `update` Skill Mode in the canonical Roundfix skill and its `modes` reference, including the exact check, upgrade, preview, Plan Digest confirmation, deletion/uncommitted-change questions, and no-commit contract. Updated the mode contract tests for both layouts and changed the `CONTEXT.md` **Skill Mode** definition to include `update` and its Plan Digest behavior. After Verification feedback identified the 20,151-byte skill entry budget failure, shortened only the free-form frontmatter description while preserving the required mode names and contract text.

Focused checks:

- `make skills-sync`: passed after the budget repair; the mirror was generated from the canonical files.
- `go test -count=1 ./skills -run '^TestRoundfixSkillEntryFileStaysWithinItsBudget$'`: passed, covering both layouts (3 passed in 1 package).
- `go test -count=1 ./skills -run '^TestTheRoundfixSkillReadsItsModeFromTheRequest$'`: passed, covering both layouts (3 passed in 1 package).
- `go test -count=1 ./skills -run '^TestTheUpdateModeAppliesOnlyThePlanItShowed$'`: passed, covering both layouts (3 passed in 1 package).
- `cmp .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` and the corresponding `modes.md` copies: passed.
- `git diff --check`: passed.

Acceptance evidence:

- `update` and `--update` are named by the canonical Modes paragraph, the mode bullet, the modes index row, the `How a mode is chosen` reference, and the focused contract test.
- `TestTheUpdateModeAppliesOnlyThePlanItShowed` asserts the opening guard, five ordered items, all required commands and safety phrases for both layouts.
- Canonical and mirror skill files match, the contract test pins the shared wording, and `CONTEXT.md` carries the API Contract 10 glossary text.

The Daemon-owned `## Verification` commands were not rerun; the failure was the skill byte-budget check in the suite. The owned-skill version recorder is rerun as the final repair step after this Result update; Daemon Verification and Task settlement remain pending.
