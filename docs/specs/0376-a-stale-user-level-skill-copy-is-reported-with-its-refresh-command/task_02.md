---
task: task_02
spec: 0376-a-stale-user-level-skill-copy-is-reported-with-its-refresh-command
status: completed
type: docs
complexity: medium
---

# Task 02: The guides, the Roundfix references, the --update mode and the glossary describe the report

## Overview

The `skills`, `doctor` and `update` command guides and both copies of the
Roundfix `setup` reference describe the report task_01 prints; step 1 of the
`--update` Skill Mode asks before running a refresh the Update Check names;
the glossary adds **Stale User Skill Copy**. This answers the Backlog Entry
recorded on 2026-11-10, "A user-level skill copy older than the project copy
is loaded unnoticed".

## Requirements

1. MUST change `docs/user-guide/commands/skills.md` as `_techspec.md` → API
   Contract 8 says: after the paragraph that ends "it checks only the
   shipped bundle.", a paragraph that names the five directories in their
   order (`$CLAUDE_HOME/skills` or `~/.claude/skills`, `$CODEX_HOME/skills`
   or `~/.codex/skills`, `$OPENCODE_HOME/skills` or `~/.opencode/skills`,
   `~/.agents/skills` and `~/.config/opencode/skills`), says that symbolic
   links are followed, that a folder several links reach is named once under
   the first directory, that a copy at or above the binary's version, one
   without a comparable version and an unreadable one are not named, that
   nothing is written, and that the exit code does not change; then this
   `text` block, word for word:

   ```text
   roundfix: warning: user-level Roundfix skill copies older than this binary; an Agent harness may load them in place of the project copy:
     <skill> <found> (binary <version>) at <path>[ -> <link target>]
     refresh: <command>
   ```

   and a paragraph naming the refresh forms of API Contract 1 and that a
   refresh writes through a link to the folder ` -> <link target>` names.
2. MUST change `docs/user-guide/commands/doctor.md`: after the paragraph that
   ends "The release plan's `skills` check reads this same line.", a sentence
   that Doctor prints this line right after `skills:` when a Stale User Skill
   Copy exists, never changes its exit code, prints no such line otherwise,
   and that the release plan does not read it, then this `text` block, word
   for word:

   ```text
   user-skills: warn (user-level Roundfix skill copies older than this binary: <skill> <found> (binary <version>) at <path>[ -> <link target>], ...; next: <refresh> && <refresh>)
   ```

3. MUST change `docs/user-guide/commands/update.md`: after the paragraph that
   ends "The final outside-repository line skips the preview.", a sentence
   that this line follows the Baseline line when a Stale User Skill Copy
   exists, inside or outside a Git repository and for a development build
   too, then this `text` block, word for word:

   ```text
   User skills: copies older than this binary: <skill> <found> (binary <version>) at <path>[ -> <link target>], ...; refresh with <refresh> && <refresh>.
   ```

   then "It does not change `next`; the `--update` Skill Mode asks before
   running a refresh."; and, right before the `next` row of the JSON table,
   the row whose first cell is `userSkills`, saying: each Stale User Skill Copy
   with `skill`, `path`, `resolved` (the link target, omitted for a plain
   folder), `found`, `embedded` and `refresh`; omitted when there is none.
4. MUST change the canonical `.agents/skills/roundfix/references/setup.md`:
   right after the paragraph that ends "with `roundfix baseline skills
   restore --profile <id> --skill <name>`.", the paragraph of API Contract 9,
   whose first sentence begins, word for word, "`roundfix skills check`,
   `roundfix doctor` and `roundfix update check` also name each Stale User
   Skill Copy:" and which ends, word for word, "None of them changes a copy;
   run a refresh only when the user asks, because it writes outside the
   repository and through the link the report names."
5. MUST change step 1 of `## --update` in the canonical
   `.agents/skills/roundfix/references/modes.md`: right after "and the
   Baseline line.", the sentence of API Contract 9 word for word, and nothing
   else in the file; MUST keep five numbered steps and every phrase
   `skills/roundfix_modes_test.go` pins; MUST NOT add or change text inside
   any `### QA settlement` section. MUST add the phrases "`userSkills`" and
   "ask once" to item 1 of the `checks` map of
   `TestTheUpdateModeAppliesOnlyThePlanItShowed` in
   `skills/roundfix_modes_test.go`, and change nothing else in that file.
6. MUST add to `CONTEXT.md`, through `domain-modeling`, the
   **Stale User Skill Copy** entry of API Contract 10 word for word, right
   before **Missing Required Skill**, with `ADR-0386` as a Markdown link to
   `docs/adr/0386-a-user-level-copy-of-a-roundfix-owned-skill-older-than-the-binary-is-reported-never-changed.md`;
   MUST add no other capitalized phrase.
7. MUST run `make skills-sync` to copy the two references to
   `skills/roundfix/references/setup.md` and
   `skills/roundfix/references/modes.md`, then raise the Roundfix skill
   version with
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`;
   never hand-edit a mirror, a version field or the record. Of the Governed
   Paths the Task may change only `.agents/skills/roundfix/SKILL.md`,
   `.agents/skills/roundfix/references/modes.md`,
   `.agents/skills/roundfix/references/setup.md` and
   `skills/roundfix/SKILL.md`.
8. MUST, as its last step and after the last guide edit, re-record the
   Behavior Surface Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.

## Subtasks

- [ ] Document the report in the three command guides.
- [ ] Add the `setup` reference paragraph and the `--update` sentence with its contract phrases.
- [ ] Add the glossary term.
- [ ] Sync, record the skill version, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] The three guides and both `setup` copies carry the prescribed blocks and sentences.
- [ ] Both `modes` copies ask before running a refresh in step 1 of `--update`, and the contract test pins it.
- [ ] `CONTEXT.md` defines **Stale User Skill Copy**, and the skills, documentation and glossary contracts pass.

## Context

- instruction: `docs/adr/0386-a-user-level-copy-of-a-roundfix-owned-skill-older-than-the-binary-is-reported-never-changed.md`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- interface: `docs/user-guide/commands/skills.md`
- interface: `docs/user-guide/commands/doctor.md`
- interface: `docs/user-guide/commands/update.md`
- interface: `docs/references/behavior-surfaces.json`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/references/modes.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/modes.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `skills/roundfix_modes_test.go`
- interface: `CONTEXT.md`

## Consumers

- The `setup` and `modes` references and the `--update` step contract — consumers: `./skills`
- The three command guides, the Behavior Surface Record and the glossary entry — consumers: `./internal/docscontract`

## Verification

- `tr -s '[:space:]' ' ' < docs/user-guide/commands/skills.md | grep -qF -- 'roundfix: warning: user-level Roundfix skill copies older than this binary; an Agent harness may load them in place of the project copy:' || { echo 'skills guide unchanged' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/doctor.md | grep -qF -- 'user-skills: warn (user-level Roundfix skill copies older than this binary:' || { echo 'doctor guide unchanged' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/update.md | grep -qF -- 'User skills: copies older than this binary:' || { echo 'update guide unchanged' >&2; exit 1; }; grep -qE -- '^[|] .userSkills. [|]' docs/user-guide/commands/update.md || { echo 'update JSON row missing' >&2; exit 1; }; for f in .agents/skills/roundfix/references/setup.md skills/roundfix/references/setup.md; do tr -s '[:space:]' ' ' < "$f" | grep -qE -- 'roundfix update check. also name each Stale User Skill Copy:' || { echo "setup reference unchanged: $f" >&2; exit 1; }; done; for f in .agents/skills/roundfix/references/modes.md skills/roundfix/references/modes.md; do grep -qE -- 'When the report carries .userSkills., show each copy with its .refresh. command and ask once' "$f" || { echo "modes reference unchanged: $f" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./skills && go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task no guide or reference describes the report; afterwards they do, and the skills suite with its recorded version and the documentation contracts with the Behavior Surface Record pass.

- `grep -F -- '**Stale User Skill Copy**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Stale User Skill Copy\*\*(:| —( |$))' || { echo 'missing term: Stale User Skill Copy' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'none of them changes it' || { echo 'Stale User Skill Copy definition differs' >&2; exit 1; }; go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^TestGlossaryPhrasesAreDefinedTerms$'` — expected: exit 0. Before this Task the term is absent; afterwards the glossary defines it and the glossary phrase contract passes.

## References

- `_prd.md` → User Story 2; User Story 3; Core Feature 6; Success Metric 7; Glossary
- `_techspec.md` → Current behavior 6; API Contracts 8-10; Declared breaks; Build Order 2
- ADR-0386

## Result

Skills loaded: implement-task, domain-modeling, tech-writer, evidence-gate

Implementation and focused-check evidence:

- Added the prescribed stale-copy report paragraphs and text blocks to the `skills`, `doctor` and `update` command guides, including the `userSkills` JSON row and the refresh/exit-code rules.
- Added the API Contract 9 paragraph to the canonical Roundfix `setup` reference and the API Contract 9 prompt to `--update` step 1; `make skills-sync` copied both references to `skills/`.
- Added `userSkills` and `ask once` to `TestTheUpdateModeAppliesOnlyThePlanItShowed`, and added the **Stale User Skill Copy** glossary entry linked to ADR-0386.
- `rtk make skills-sync`: passed.
- `rtk go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`: passed; Roundfix skill version recorded as `0.1.217`.
- `rtk go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`: passed; Behavior Surface Record re-recorded after the final guide edit.
- `rtk go test -count=1 ./skills -run '^TestTheUpdateModeAppliesOnlyThePlanItShowed$'`: passed, 3 subtests.
- `rtk go test -count=1 -tags docscontract ./internal/docscontract -run '^(TestGlossaryPhrasesAreDefinedTerms|TestGlossaryLinksNameAcceptedADRs|TestTheSkillCoverageMapIsCurrent|TestCommandPathsReadSubcommandsAndAlternatives)$'`: passed, 4 tests.
- `cmp .agents/skills/roundfix/references/setup.md skills/roundfix/references/setup.md` and the equivalent `modes.md` comparison: passed; `git diff --check`: passed.

Acceptance evidence:

- Criterion 1: all three guides and both synced `setup` references contain the prescribed report and refresh descriptions.
- Criterion 2: both `modes` copies are synchronized with the ask-before-refresh sentence, and the update-mode contract test pins `userSkills` and `ask once`.
- Criterion 3: `CONTEXT.md` defines **Stale User Skill Copy** with the ADR-0386 link, and the focused glossary/coverage contracts pass.

The Task's declared Verification and repository Verification remain for the Daemon; they were not run in this child Agent turn.
