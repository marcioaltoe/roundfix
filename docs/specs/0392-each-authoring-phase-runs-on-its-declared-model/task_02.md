---
task: task_02
spec: 0392-each-authoring-phase-runs-on-its-declared-model
status: completed
type: backend
complexity: medium
---

# Task 02: The checks report an unapplied Authoring Model Declaration

## Overview

One read-only inspector compares the installed authoring skills and the
generated subagent files with the Authoring Model Declaration, and `roundfix
skills check`, Doctor and the Update Check print what it finds as a warning
that names `roundfix skills install --target project` and never changes an
exit code. The slice is whole on its own: in a temporary repository whose
declaration changed after an installation, each of the three reports names
every problem, and after a new installation Doctor prints its
`authoring-models` line as `ok`. A repository that declares nothing prints
what it printed before. The skills, doctor and update guides describe the
lines.

## Requirements

1. MUST add to the `skills` package the inspector of `_techspec.md` → API
   Contract 6, with the signature of `_techspec.md` → Interfaces
   (`InspectAuthoring`) and the state of `_techspec.md` → Data Models
   (`AuthoringState`). It only reads: it writes, creates, links and deletes
   nothing. It returns the problems of that contract, word for word and in
   its order, each marked as fixable by the installation or not. Where the
   contract leaves a case open, the inspector does this:
   - a declared phase whose installed `SKILL.md` does not exist gets no
     header problem;
   - the header differs when either the top-level `model` or the top-level
     `effort` the installed file holds differs from the entry, compared as
     parsed values, so a quoted and a plain spelling of one value are equal;
   - `.claude/skills/<skill>` leads to `.agents/skills/<skill>` only when
     both resolve to the same directory, so a missing `.claude/skills`, a
     copied folder and a dangling link all get the link problem, which is
     never fixable by the installation;
   - a file at a generated name that is not Roundfix's is the conflict of
     that contract, never `differs from the declaration` and never fixable
     by the installation, which refuses it: at a declared entry's path, and,
     while at least one key is declared, at a generated name without an
     entry. The inspector names the file and says that installation will not
     overwrite it, with the contract's words;
   - with a nil declaration the harness is `claude` and no key is declared,
     and a file at a generated name that is not Roundfix's is not a problem.
2. MUST make `roundfix skills check` print the warning of `_techspec.md` →
   API Contract 7 on stderr, after its result line and before the user-level
   warning, on the paths where the check passes or reports unversioned
   skills, never when it fails; the `next:` line only when a problem is
   fixable by the installation; and, when the configuration cannot be
   loaded, only the line
   `roundfix: warning: Authoring Model Declaration not read: <error>`. The
   check loads the configuration with the repository root as the working
   directory and prints none of the loader's own warnings. Its exit code
   never depends on the declaration.
3. MUST make Doctor print the `authoring-models` line of API Contract 7,
   after `skills:` and after `user-skills:` when that line prints: `ok` with
   the detail `claude: <declared keys joined by ", "> applied`, or `warn`
   with the detail and next action that contract gives, the next action only
   when a problem is fixable by the installation. Doctor prints no such line
   when nothing is declared and no generated file exists, outside a Git
   repository, and when Project Config was not read. Doctor's exit code
   never depends on the line.
4. MUST add `authoringModels` to the Update Check's JSON, declared right
   after `userSkills`, with the members of `_techspec.md` → Data Models:
   `state` (`applied` or `unapplied`), `harness`, `declared`, `problems`
   (always an array, empty when applied) and `next` (omitted when empty). It
   is present whenever a key is declared or a problem exists and absent
   otherwise, including when Project Config was not read. The text report
   prints the line of API Contract 7 after `User skills:` only when the state
   is `unapplied`, with `; apply with ...` only when a problem is fixable by
   the installation. The report's own `next` and its exit code do not change,
   and the schema version stays as it is.
5. MUST add these tests, each top-level test calling `t.Parallel()` first,
   using temporary repositories and homes through the packages' existing
   seams and never editing `internal/cli/cli_test.go`:
   - `skills`: `TestInspectAuthoringNamesEachProblem`, with the subtests
     `clean`, `header`, `link`, `missing`, `differs`, `conflict`,
     `no_declared_entry`, `order` (fixture B of `_techspec.md` → Surface
     Transcripts: its four problems in the transcript's order) and
     `nothing_declared`. Each builds a real installation, changes one thing,
     asserts the whole list of problems with its fixable marks, and asserts
     that the inspected tree, links included, did not change. `conflict`
     puts a hand-written file without the marker at a declared entry's path
     and another at a generated name that has no entry: each is named with
     the conflict text and marked as not fixable by the installation, and
     the same subtest shows why, by installing again and asserting that the
     installation returns its refusal and leaves both files unchanged.
   - `internal/cli`:
     `TestSkillsCheckWarnsOfAnUnappliedDeclarationAndKeepsItsExit`, with the
     subtests `declaration` (fixture B with a temporary home that holds no
     skill copy: stdout is the `Roundfix skill check passed: <names>` line,
     stderr is exactly the first six stderr lines of Surface Transcript 3,
     exit `0`; then, after a new installation, only the link problem remains
     and no `next:` line prints; then, with the link made, stderr is empty)
     and `not_read` (a Project Config that does not load: the one `not read`
     line whole, exit `0`). Task 03 adds the subtest for the whole
     transcript.
   - `TestDoctorReportsTheAuthoringModels`, with the subtests `ok`, `warn`,
     `warn_without_next` and `after_user_skills`: the line word for word
     right after `skills:`, and, with an older user-level copy present, after
     the line that starts with `user-skills:`, matched by that prefix only
     because task 03 changes its text; exit `0` each time.
   - `TestUpdateCheckReportsTheAuthoringModels`, with the subtests `applied`,
     `unapplied` and `unapplied_without_next`: the JSON member whole, its
     place between `userSkills` and `next`, the text line after
     `User skills:`, and the report's `next` unchanged.
   - MUST extend `TestARepositoryWithoutTheDeclarationPrintsWhatItPrinted`
     with the subtests `skills_check`, `doctor` and `update_check`: with the
     skills installed and no section, the stderr of `skills check` is empty,
     Doctor prints no `authoring-models` line, and the Update Check prints no
     `Authoring models:` line and no `authoringModels` member.
6. MUST document the lines in the guides, each shown in a `text` block as
   the command prints it, with `<problem>`, `<error>` and `<declared keys>`
   as placeholders: `docs/user-guide/commands/skills.md` shows the three
   lines of the `skills check` warning and the line
   `roundfix: warning: Authoring Model Declaration not read: <error>`, and
   lists the problems of API Contract 6 in its order, saying of the conflict
   that installation will not overwrite a file it did not generate and that
   the user moves or renames it;
   `docs/user-guide/commands/doctor.md` shows
   `authoring-models: ok (claude: <declared keys> applied)` and the `warn`
   form with its next action, and says when the line is not printed; and
   `docs/user-guide/commands/update.md` shows the text line that begins
   `Authoring models: the installed files do not match the declaration:` and
   describes `authoringModels` with its five members. MUST then re-record the Behavior Surface Record, as the
   Task's last step, after the last guide edit, with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
7. MUST put the new code and tests in the files `## Context` lists under
   `creates:`; `internal/cli/authoring_models_test.go` is the file task 01
   created, extended here. MUST NOT change the text of the user-level notice, the usage,
   any file under `.agents/skills/` or any shipped skill file under
   `skills/`, `CONTEXT.md`, any ADR, `.roundfixrc.yml`, the Makefile, the CI
   workflows, `go.mod`, `go.sum` or the lint configuration. No recorded test
   is removed or renamed. Every test file compiles for Windows.

## Subtasks

- [ ] Add the inspector and its test.
- [ ] Print the warning in `skills check`, the line in Doctor and the member and line in the Update Check.
- [ ] Add the CLI tests and extend the test of a repository without the declaration.
- [ ] Write the three guides and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] For every state of the installed files, the inspector returns the problems of API Contract 6 in its order with the right fixable marks, never marks a file Roundfix did not generate as fixable by the installation, and changes no file.
- [ ] With a declaration whose files are missing or different, `roundfix skills check`, Doctor and the Update Check each print the warning with the installation command and keep their exit codes; with matching files Doctor prints its `authoring-models` line as `ok`.
- [ ] Without the section and without a generated file, none of the three reports prints a new line or member.
- [ ] The three guides describe the lines, and the Behavior Surface Record is current.

## Context

- instruction: `docs/adr/0401-each-authoring-phase-runs-on-the-model-its-repository-declares.md`
- creates: `skills/authoring_inspect.go`
- creates: `skills/authoring_inspect_test.go`
- creates: `internal/cli/authoring_reports.go`
- creates: `internal/cli/authoring_reports_test.go`
- creates: `internal/cli/authoring_models_test.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/health.go`
- interface: `internal/cli/update_check.go`
- interface: `docs/user-guide/commands/skills.md`
- interface: `docs/user-guide/commands/doctor.md`
- interface: `docs/user-guide/commands/update.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The inspector's state and the lines of the three reports — consumers: `./skills`, `./internal/cli`
- The skills, doctor and update guides and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^TestInspectAuthoringNamesEachProblem$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestInspectAuthoringNamesEachProblem TestInspectAuthoringNamesEachProblem/order TestInspectAuthoringNamesEachProblem/link TestInspectAuthoringNamesEachProblem/conflict TestInspectAuthoringNamesEachProblem/nothing_declared; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./skills` — expected: exit 0. Before this Task the inspector and its test do not exist, so the pass lines are missing; afterwards every problem of API Contract 6 is named in its order with its fixable mark, a file Roundfix did not generate is a conflict the installation does not fix, the inspected tree is unchanged, and the skills suite passes.

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestSkillsCheckWarnsOfAnUnappliedDeclarationAndKeepsItsExit|TestDoctorReportsTheAuthoringModels|TestUpdateCheckReportsTheAuthoringModels|TestARepositoryWithoutTheDeclarationPrintsWhatItPrinted)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSkillsCheckWarnsOfAnUnappliedDeclarationAndKeepsItsExit/declaration TestSkillsCheckWarnsOfAnUnappliedDeclarationAndKeepsItsExit/not_read TestDoctorReportsTheAuthoringModels/ok TestDoctorReportsTheAuthoringModels/warn TestDoctorReportsTheAuthoringModels/after_user_skills TestUpdateCheckReportsTheAuthoringModels/applied TestUpdateCheckReportsTheAuthoringModels/unapplied TestARepositoryWithoutTheDeclarationPrintsWhatItPrinted/skills_check TestARepositoryWithoutTheDeclarationPrintsWhatItPrinted/doctor TestARepositoryWithoutTheDeclarationPrintsWhatItPrinted/update_check; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m -run 'Skill|Doctor|UpdateCheck' ./internal/cli` — expected: exit 0. Before this Task the three report tests and the three new subtests do not exist, so their pass lines are missing; afterwards each report prints the contract's lines and keeps its exit code, a repository without the section prints nothing new, and the CLI's skills, Doctor and Update Check tests pass.

- `tr -s '[:space:]' ' ' < docs/user-guide/commands/skills.md | grep -qF -- 'roundfix: warning: Authoring Model Declaration not read: <error>' || { echo 'skills guide lacks the not read line' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/doctor.md | grep -qF -- 'authoring-models: ok (claude: <declared keys> applied)' || { echo 'doctor guide lacks the authoring-models line' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/update.md | grep -qF -- 'Authoring models: the installed files do not match the declaration:' || { echo 'update guide lacks the Authoring models line' >&2; exit 1; }; go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the three guides do not show the lines; afterwards they do, and the documentation contracts, the Behavior Surface Record included, pass.

## References

- `_prd.md` → User Story 6; Core Feature 5; Success Metric 4
- `_techspec.md` → Interfaces; Data Models; API Contract 6; API Contract 7; API Contract 15; Surface Transcript 3; Invariants 1, 5 and 8; Testing Approach; Build Order 2
- ADR-0401

## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, tech-writer, evidence-gate

Implemented this Task's inspector and reporting slice. `InspectAuthoring`
compares parsed top-level headers, directory identities and generated bytes,
using the installer's ownership-marker reader. It reports ordered problems
with installation fixability; links and handwritten-file conflicts never
recommend installation as their own fix. A nil declaration still reports
Roundfix-generated leftovers and ignores unmarked files.

The three reports consume that state without changing their exit contracts.
`skills check` loads at the resolved repository root, suppresses loader
warnings and prints the dedicated not-read warning on configuration failure.
Doctor reports after skills/user-skills and omits the line outside Git or
when Project Config was not read. Update Check includes the optional member
immediately after userSkills, emits arrays for declared/problems and keeps
its report-level next action and schema unchanged. The user-level notice,
usage and shipped skills are unchanged.

The three command guides show the exact warning forms, ordered problems,
fixability and omission rules. The sanctioned recorder was run after the
last guide edit; its diff changes only those three guide fingerprints.

Focused evidence (commands were executed through `rtk proxy`, except the
Windows cross-compilations):

- Before implementation, `go test ./skills -run '^TestInspectAuthoringNamesEachProblem$/order$' -count=1` exited 1 because `AuthoringProblem` did not exist, establishing the missing inspector API.
- `go test ./skills -run '^TestInspectAuthoringNamesEachProblem$/(order|conflict|clean|header|link|missing|differs|no_declared_entry|nothing_declared|quoted_header|effort_header|extra_header|missing_skill)$' -count=1` exited 0. Whole problem lists and fixability matched, including fixture B's four problems in order. Tree snapshots include directory metadata, file bytes/modes/times and link destinations; inspection preserved them. The conflict case also asserted the installer's exact refusal and that neither handwritten file nor the rest of the tree changed.
- `go test ./skills -run '^TestInspectAuthoringNamesEachProblem$/(leftovers_without_declaration|quoted_header|missing_skill|link)$' -count=1` exited 0 after adding nil-declaration leftover coverage; the link cases include absent/copied directories, a dangling link and a whole-directory link.
- `go test ./internal/cli -run '^(TestSkillsCheckWarnsOfAnUnappliedDeclarationAndKeepsItsExit|TestDoctorReportsTheAuthoringModels|TestUpdateCheckReportsTheAuthoringModels|TestARepositoryWithoutTheDeclarationPrintsWhatItPrinted)/(declaration|not_read|unversioned|ok|warn|warn_without_next|after_user_skills|project_config_not_read|outside_repository|applied|unapplied|unapplied_without_next|skills_check|doctor|update_check)$' -count=1 -timeout 5m` exited 0. Tests assert fixture B's six stderr lines, reinstallation leaving only the link warning without next, silence after linking, exact configuration-error output, unversioned and failed-check behavior, Doctor placement/status, complete Update Check JSON including member order and empty problems array, and unchanged exits/report-level next. The three no-declaration subtests add no authoring output.
- `go test ./internal/cli -run '^(TestSkillsCheckWarnsOfStaleUserSkillCopiesAndKeepsItsExit|TestDoctorReportsStaleUserSkillCopiesWithoutFailing|TestUpdateCheckNamesStaleUserSkillCopies|TestUpdateCheckTranscripts|TestSkillsInstallAppliesTheAuthoringModelDeclaration|TestBaselineUpdateAppliesTheAuthoringModelDeclaration)$' -count=1 -timeout 5m` exited 0, preserving nearby user-copy notices, release-report transcripts and declaration installation behavior.
- `go test -tags docscontract ./internal/docscontract -run '^(TestUserGuideLinksResolve|TestEveryCommandIsNamedInTheUserGuide|TestCommandReferenceLinksResolve)$' -count=1` exited 0.
- `GOOS=windows go test -c -o /private/tmp/roundfix-0392-skills.test.exe ./skills` and `GOOS=windows go test -c -o /private/tmp/roundfix-0392-cli.test.exe ./internal/cli` each exited 0. This proves compilation, not Windows execution.
- Final sanctioned generation: `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` exited 0 after the final guide edit.
- `git diff --check` exited 0 before recording this Result; final scope inspection is limited to the declared implementation/tests, three guides, Behavior Surface Record and this Task file.

Acceptance evidence:

| Criterion | Implementation and focused evidence |
| --- | --- |
| Inspector problems, order, fixability and no writes | Inspector cases assert complete lists, refusal and unchanged tree/link snapshots. |
| Three reports warn without changing exits; matching Doctor line is ok | Focused command cases assert exact warning/ok lines, placement, JSON and exits, with installation next omitted for link-only problems. |
| No declaration and no generated file adds nothing | `skills_check`, `doctor` and `update_check` subtests pass with a real installed bundle and temporary home/repository. |
| Guides and current Behavior Surface Record | Exact text blocks and JSON-field reference added; guide-link checks and final sanctioned recorder exited 0. |

Only this Task file was modified at entry (the Daemon-owned status change),
and that status is preserved. No Task Graph or other Task file was edited.
The declared Verification commands, package-wide suites and repository gates
were not run; Daemon Verification and settlement remain outstanding. Task 03
retains ownership of the new user-copy notice and the whole Surface
Transcript 3 test. No commit, push or pull request was made.
