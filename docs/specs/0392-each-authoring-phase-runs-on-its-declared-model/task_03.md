---
task: task_03
spec: 0392-each-authoring-phase-runs-on-its-declared-model
status: pending
type: backend
complexity: high
---

# Task 03: A reported user-level skill copy carries the command that removes it

## Overview

`roundfix skills remove` lists, and with `--apply` removes, the
Roundfix-owned skills in a user-level skill directory. `roundfix skills
check`, Doctor and the Update Check name each Stale User Skill Copy with that
removal command before its refresh command, and name too, whatever its
version, a copy in Claude Code's user-level directory that hides a declared
phase's Declared Skill Header. The slice is whole on its own: with a
temporary home, the three reports print the new notice, the command lists,
removes, unlinks a link while its target survives and refuses the project.
The usage, the skills, doctor and update guides, the Skill Coverage Map and
the Roundfix skill's `setup` reference name the command.

## Requirements

1. MUST change the user-level reader of the `skills` package as
   `_techspec.md` → API Contract 8 says, with the fields, JSON tags and
   signatures of `_techspec.md` → Interfaces: each of the five directories
   carries its target, its removal command and whether a copy there hides a
   header (the first directory only); each reported copy carries `remove`,
   `stale` and `hidesDeclaredHeader`, added beside the existing fields and
   renaming none; `ReportedUserSkillCopies` returns each older copy and, in a
   directory that hides a header, each readable copy of a declared skill
   whatever its version, a folder still named once; `Describe` appends one of
   the three reasons; and the removal commands and the refresh commands are
   each returned once per directory, the refresh commands for older copies
   only. A copy whose frontmatter cannot be parsed is skipped, as before. The
   reader writes nothing.
2. MUST change the three reports to the texts of API Contract 8, word for
   word: the `skills check` block with its new first line, the reason on each
   copy, `remove:` lines before `refresh:` lines; Doctor's `user-skills`
   detail, its `; or refresh with ...` part only when a copy is older, and
   its next action, the removals joined by ` && `; and the Update Check's
   `User skills:` line with `, or refresh with ...` only when a copy is
   older. The declared skills are the skills of the declared phases of the
   loaded configuration and are empty when it is not loaded; `skills check`
   takes them from the configuration task 02 already loads there and does not
   load it twice. No exit code changes and no report writes a file.
3. MUST add `roundfix skills remove` as API Contract 9 says, every stdout
   line and every refusal word for word. The hazards it bounds are stated
   here because the command deletes files in the user's home:
   - without `--apply` it only lists, and it changes nothing;
   - it acts only on the first directory of each target, or on `--dir` for a
     single target, and there only on the twelve Roundfix-owned names, in
     sorted order, in the order `claude`, `codex`, `opencode`;
   - a symbolic link is unlinked and its target is never opened or removed;
     `<link target>` in the output is the link's own text, read without
     following it; a folder is removed with its content; a file at an owned
     name, a missing name and every other name are left alone;
   - it resolves every target's directory before it touches anything; when
     none is known for a target it prints
     `roundfix: skills remove failed: no user-level skill directory is known for <target>`
     and exits `1`;
   - `--target project` is refused, and a missing or unknown `--target`
     value gets the `requires --target` refusal, each with exit `2`;
   - the closing `Nothing was removed; run the command again with --apply.`
     line prints once, at the end, only when the listing named an entry and
     `--apply` was not given.
4. MUST make `roundfix skills --help` and `roundfix skills remove --help`
   print Surface Transcript 1 of `_techspec.md`, and add to the top-level
   usage the line API Contract 10 gives, right after the `skills install`
   line.
5. MUST add these tests, each top-level test calling `t.Parallel()` first,
   with temporary homes through the packages' existing seams, never the
   host's user-level directories, and never editing
   `internal/cli/cli_test.go`:
   - `skills`:
     `TestReportedUserSkillCopiesNameACurrentCopyThatHidesAHeader`, in one
     temporary home: a copy of a declared skill at the binary's version in
     the Claude directory is named with `HidesDeclaredHeader` true and
     `Stale` false, and the same copy in the Codex directory is not named;
     an older copy of a declared skill in the Claude directory carries both
     reasons; a copy of a declared skill there without a comparable version
     is named with `Found` as read; a current copy of a skill that is not
     declared is not named; a folder reached by a link from two directories
     is named once; and without declared skills only the older copy is
     named.
   - `internal/cli`: `TestSkillsRemoveListsThenRemovesOnlyOwnedCopies`, with
     the subtests `surface_transcript_4` and `surface_transcript_5` (each
     transcript whole: every stdout and stderr line and the exit code, and
     the copy still there afterwards), `apply` (a folder with a nested file,
     a link to a folder outside the directory, a folder Roundfix does not
     own, a file at an owned name and a copy under a repository's
     `.agents/skills`: the listing whole, then with `--apply` the output
     whole, the folder and the link gone, and the link's target with its
     files, the other folder, the file and the repository's copy still
     there), `target_all`, `dir`, `refusals` (the three refusals word for
     word, and a removal that fails with exit `1`) and `no_known_directory`;
     and `TestSkillsUsageNamesTheRemoval` (Surface Transcript 1 whole for
     both help forms, and the top-level line right after the `skills
     install` line).
   - MUST add to
     `TestSkillsCheckWarnsOfAnUnappliedDeclarationAndKeepsItsExit` the
     subtest `surface_transcript_3`: fixture B with
     `<home>/.claude/skills/write-prd/SKILL.md` declaring the binary's
     version, the whole stdout and stderr of Surface Transcript 3 in order,
     and exit `0`.
   - MUST update every test that pins the old user-level text on the default
     branch when the Task runs, keeping each test's name. At the authoring
     base these are the three tests of `internal/cli/user_skills_test.go`
     and three tests of `skills/user_copies_test.go`.
     `TestSkillsCheckWarnsOfStaleUserSkillCopiesAndKeepsItsExit` pins, with
     its two older copies, the `remove:` lines before the `refresh:` lines.
   - MUST add, in a new test file of the `skills` package,
     `TestTheRoundfixSkillOffersTheRemovalOfAUserLevelCopy`, reading both
     copies of the `setup` reference through `roundfixLayoutSources` with one
     subtest per copy (`canonical`, `embedded`) and, under it, one named
     subtest per phrase, each found exactly once in whitespace-normalized
     text:
     - `reads_five_directories`: "they read the five user-level directories in the `skills check` order, follow symbolic links and name a folder once under the first directory that reaches it"
     - `names_hidden_headers`: "a copy of a declared phase's skill in Claude Code's user-level directory, because it hides the Declared Skill Header"
     - `keeps_exit_codes`: "`skills check` warns on stderr and keeps its exit code, Doctor adds `user-skills: warn` after `skills:`, and the Update Check adds `User skills:` and `userSkills`"
     - `names_removal_first`: "each names the removal command, `roundfix skills remove --target <claude|codex|opencode> --apply`, with `--dir <path>` for the other two directories, before the refresh command"
     - `changes_nothing_unasked`: "None of them changes a copy; run a removal or a refresh only when the user asks"
     - `writes_outside_the_repository`: "because both write outside the repository, a refresh through the link the report names"
     - `removal_lists_first`: "Without `--apply`, `roundfix skills remove` only lists what it would remove, and it never touches the project's copies"
6. MUST replace, in the canonical
   `.agents/skills/roundfix/references/setup.md`, the paragraph that begins
   "`roundfix skills check`, `roundfix doctor` and `roundfix update check`
   also name each Stale User Skill Copy:" with the fifth block of
   `_techspec.md` → API Contract 12, word for word, and change nothing else
   in that file: the documentation contract requires the Roundfix skill to
   name every command, so the command and its reference text arrive
   together. The paragraph keeps what it guaranteed before, as the phrases of
   requirement 5 pin. MUST NOT add or change text inside any `### QA
   settlement` section. MUST then run, in this order and never by hand,
   `make skills-sync` and the version record
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   which raises the version of `roundfix` from the tree the Task runs on
   (task 04 raises it again; never write a number).
7. MUST document the slice in the guides:
   `docs/user-guide/commands/skills.md` gains the synopsis line
   `roundfix skills remove --target claude|codex|opencode|all [--dir <path>] [--apply]`,
   shows the new `skills check` block in a `text` block with `<reason>` and
   `remove: <command>`, names the removal and the refresh command of each of
   the five directories, and describes `skills remove` with its six stdout
   lines and its four stderr lines in `text` blocks;
   `docs/user-guide/commands/doctor.md` shows the new `user-skills` line; and
   `docs/user-guide/commands/update.md` shows the new `User skills:` line and
   describes `remove`, `stale` and `hidesDeclaredHeader`. MUST add to
   `docs/references/skill-coverage.json`, by hand and in the file's existing
   order and shape, the entry of the new Behavior Surface: id
   `command skills remove`, covering skill
   `.agents/skills/roundfix/references/setup.md`, and sources
   `docs/user-guide/commands/skills.md` and the Go file that holds the
   command, never `internal/cli/cli.go`. MUST then re-record the Behavior
   Surface Record, as the Task's last step, after the last guide and usage
   edit, with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
8. MUST put the new code and tests in the files `## Context` lists under
   `creates:`; `internal/cli/authoring_reports.go` and
   `internal/cli/authoring_reports_test.go` are the files task 02 created,
   changed here. Of the Governed Paths the Task may change only
   `.agents/skills/roundfix/references/setup.md`,
   `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`. MUST
   NOT change any other skill or reference, `CONTEXT.md`, any ADR,
   `.roundfixrc.yml`, the Makefile, the CI workflows, `go.mod`, `go.sum` or
   the lint configuration. No recorded test is removed or renamed. Every
   test file compiles for Windows, and a case that needs a file mode Windows
   lacks is skipped there.

## Subtasks

- [ ] Extend the user-level reader with the removal command, the stale mark and the copy that hides a header.
- [ ] Print the new notice in the three reports.
- [ ] Add `roundfix skills remove` and the usage.
- [ ] Add and update the tests.
- [ ] Rewrite the `setup` paragraph, sync the mirror and record the Roundfix skill's version.
- [ ] Write the three guides and the Skill Coverage Map entry, and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] With a temporary home, the three reports name the removal command before the refresh command, a copy at the binary's version of a declared authoring skill in Claude Code's user-level directory is named and the same copy elsewhere is not, and no exit code changes.
- [ ] `roundfix skills remove` lists without `--apply`, removes with it, unlinks a link while its target survives, leaves every name it does not own and every repository copy alone, and refuses the project target.
- [ ] The usage, the three guides, the Skill Coverage Map and both copies of the Roundfix `setup` reference name the command, the Roundfix skill's version is recorded, and the Behavior Surface Record is current.

## Context

- instruction: `docs/adr/0405-a-user-level-skill-copy-is-reported-with-the-command-that-removes-it.md`
- instruction: `docs/adr/0386-a-user-level-copy-of-a-roundfix-owned-skill-older-than-the-binary-is-reported-never-changed.md`
- instruction: `skills/roundfix_layout_test.go`
- interface: `skills/user_copies.go`
- interface: `skills/user_copies_test.go`
- creates: `skills/user_copy_removal_guidance_test.go`
- interface: `internal/cli/user_skills.go`
- interface: `internal/cli/user_skills_test.go`
- creates: `internal/cli/skills_remove.go`
- creates: `internal/cli/skills_remove_test.go`
- creates: `internal/cli/authoring_reports.go`
- creates: `internal/cli/authoring_reports_test.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/update_check.go`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/skills.md`
- interface: `docs/user-guide/commands/doctor.md`
- interface: `docs/user-guide/commands/update.md`
- interface: `docs/references/skill-coverage.json`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The user-level reader's fields and the Roundfix `setup` reference — consumers: `./skills`
- The three reports' user-level notice, the Update Check's `userSkills` members and the usage — consumers: `./internal/cli`
- The skills, doctor and update guides, the Skill Coverage Map and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestReportedUserSkillCopiesNameACurrentCopyThatHidesAHeader|TestTheRoundfixSkillOffersTheRemovalOfAUserLevelCopy|TestEveryOwnedSkillVersionIsRecorded)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestReportedUserSkillCopiesNameACurrentCopyThatHidesAHeader TestTheRoundfixSkillOffersTheRemovalOfAUserLevelCopy/canonical/names_removal_first TestTheRoundfixSkillOffersTheRemovalOfAUserLevelCopy/embedded/names_removal_first TestEveryOwnedSkillVersionIsRecorded; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; diff -r .agents/skills/roundfix skills/roundfix >/dev/null || { echo 'mirror differs: roundfix' >&2; exit 1; }; go test -count=1 -timeout 20m ./skills` — expected: exit 0. Before this Task the two tests do not exist, so their pass lines are missing; afterwards the reader names the copy that hides a header only in Claude Code's directory, both copies of the `setup` reference carry each phrase once, the Roundfix skill's version is recorded, and the skills suite passes.

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestSkillsRemoveListsThenRemovesOnlyOwnedCopies|TestSkillsUsageNamesTheRemoval|TestSkillsCheckWarnsOfAnUnappliedDeclarationAndKeepsItsExit|TestSkillsCheckWarnsOfStaleUserSkillCopiesAndKeepsItsExit)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSkillsRemoveListsThenRemovesOnlyOwnedCopies/surface_transcript_4 TestSkillsRemoveListsThenRemovesOnlyOwnedCopies/surface_transcript_5 TestSkillsRemoveListsThenRemovesOnlyOwnedCopies/apply TestSkillsRemoveListsThenRemovesOnlyOwnedCopies/refusals TestSkillsUsageNamesTheRemoval TestSkillsCheckWarnsOfAnUnappliedDeclarationAndKeepsItsExit/surface_transcript_3 TestSkillsCheckWarnsOfStaleUserSkillCopiesAndKeepsItsExit; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m -run 'Skill|Doctor|UpdateCheck|Usage|Help' ./internal/cli` — expected: exit 0. Before this Task the removal and usage tests and the transcript subtest do not exist, so their pass lines are missing; afterwards the command lists, removes and refuses as API Contract 9 says, the usage is Surface Transcript 1, `roundfix skills check` prints Surface Transcript 3, and the CLI's skills, Doctor, Update Check and usage tests pass.

- `tr -s '[:space:]' ' ' < docs/user-guide/commands/skills.md | grep -qF -- 'roundfix skills remove --target claude|codex|opencode|all [--dir <path>] [--apply]' || { echo 'skills guide lacks the command' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/skills.md | grep -qF -- 'Nothing was removed; run the command again with --apply.' || { echo 'skills guide lacks the listing lines' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/doctor.md | grep -qF -- 'user-skills: warn (user-level Roundfix skill copies that an Agent harness may load in place of the project copy:' || { echo 'doctor guide lacks the new user-skills line' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/update.md | grep -qF -- 'User skills: copies that an Agent harness may load in place of the project copy:' || { echo 'update guide lacks the new User skills line' >&2; exit 1; }; grep -qF -- '"id": "command skills remove"' docs/references/skill-coverage.json || { echo 'coverage map lacks the command' >&2; exit 1; }; go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the guides and the Skill Coverage Map do not name the command; afterwards they do, and the documentation contracts pass: every command is named in the Roundfix skill, and the Skill Coverage Map and the Behavior Surface Record are current.

## References

- `_prd.md` → User Story 5; Core Features 6 and 7; Success Metric 5
- `_techspec.md` → Interfaces; API Contract 8; API Contract 9; API Contract 10; API Contract 12; API Contract 15; Surface Transcript 1; Surface Transcript 3; Surface Transcript 4; Surface Transcript 5; Invariants 5, 6 and 8; Declared breaks; Testing Approach; Build Order 3
- ADR-0405, ADR-0386
