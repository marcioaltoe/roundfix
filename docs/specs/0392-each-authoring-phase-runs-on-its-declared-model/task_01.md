---
task: task_01
spec: 0392-each-authoring-phase-runs-on-its-declared-model
status: completed
type: backend
complexity: high
---

# Task 01: A repository declares its authoring models and installation applies the declaration

## Overview

Project Config gains the Authoring Model Declaration, `skills.authoring`,
read and validated by the configuration loader and ignored with a warning in
User Config. `roundfix skills install --target project` and the skills stage
of `roundfix baseline update` load it and hand it to the one installer they
share, which writes the Declared Skill Header into each declared authoring
skill and one generated subagent file per declared entry, refuses a file it
did not generate, and removes its own left-over files. The slice is whole on
its own: in a temporary repository that declares the section, one
installation leaves the header and the generated files, a second changes
nothing, and a repository that declares nothing prints and installs what it
did before. The configuration, skills and baseline guides describe it.

## Requirements

1. MUST read `skills.authoring` from Project Config into the typed values of
   `_techspec.md` → Interfaces (`config.Skills`), as `_techspec.md` → API
   Contract 1 says. The section's keys are yaml-tagged struct fields
   (`skills`, `authoring`, `claude`, `idea`, `prd`, `techspec`, `tasks`,
   `research`, `model`, `reasoning_effort`, `tools`), so the `config keys`
   Behavior Surface moves with them. The `skills` node is validated before
   the strict decode, and each refusal returns, after the loader's usual
   `parse config "<label>": ` prefix, one of the twelve messages of API
   Contract 1, word for word. A model and a tool match the two patterns that
   contract gives. A section whose `authoring` or `claude` mapping is empty
   declares nothing and is accepted. MUST also refuse, never accept, the
   shapes the twelve messages do not name: a `skills`, `skills.authoring` or
   `skills.authoring.claude` node that is not a mapping is refused with
   `<path> must be a mapping`, and a `tools` value that is not a list of
   strings, or a `model` or `reasoning_effort` that is not a string, is
   refused with the tool, model or effort message naming the value as
   written.
2. MUST remove the whole `skills` node from User Config before decoding and
   write, once, through the configuration warnings, the line
   `config: skills in User Config is ignored; set skills in Project Config`.
   The loaded declaration comes from Project Config only. A repository
   without the section loads exactly as before, with no warning.
3. MUST add to the `skills` package the phase table of API Contract 2, the
   header writer of API Contract 3 and the generator of API Contract 4, with
   the signatures of `_techspec.md` → Interfaces and the bytes of those
   contracts exactly: the two frontmatter lines and their quoting rule, the
   two generated bodies word for word, LF line endings and one final newline,
   and `model:` and `effort:` lines only for declared values. The Research
   Subagent's `tools` line is built as `Read, Grep, Glob` followed by the
   names it is given, never copied from a file. A file at a generated name is
   Roundfix's only when it is a regular file whose first non-blank line after
   the frontmatter is the marker of the key its name belongs to. The
   production code of `skills` imports neither `internal/config` nor
   `internal/cli`.
4. MUST add the installation step of API Contract 5 to the installer that
   both commands share, in the order (a) to (e) it gives. The hazards it
   bounds are stated here because the step overwrites and deletes files:
   - with a declaration, a file that is not Roundfix's at any of the five
     generated names makes the installation return an error before it writes
     anything, the skills included;
   - a generated file is overwritten only when it is Roundfix's, and a file
     that is Roundfix's and has no declared entry is deleted; no other path
     under `.claude/agents` is read, written or deleted;
   - an installed authoring `SKILL.md` that held a top-level `model` or
     `effort` key no entry covers is named once in the result before it is
     overwritten;
   - a nil declaration writes the bundle unchanged, refuses nothing and
     deletes only files that are Roundfix's;
   - only the project target without a directory override applies a
     declaration: a user-level target and a `--dir` target receive no header
     and no generated file, whatever the request carries;
   - a second installation with the same declaration changes no byte and no
     file mode.
5. MUST make `roundfix skills install`, for `--target project` without
   `--dir`, load the configuration through the command's configuration
   loader with the repository root as the working directory, writing the
   loader's warnings to stderr, and pass the declaration to the installer. A
   load error stops the command before anything is written, with
   `roundfix: skills install failed: <error>` and exit `2`. After the
   `Installed Roundfix skill for project:` line and before the lock lines it
   prints the three stdout lines of API Contract 5 in the phase table's
   order, headers first, then generated files, then removals; on stderr the
   notice line of that contract once per uncovered skill, and the refusal
   line with exit `2`. The lock entry of each installed owned skill is
   recorded after the header is written, so it equals the installed tree's
   hash.
6. MUST make the skills stage of `roundfix baseline update` load and pass the
   declaration the same way. A load or installation error fails the stage
   with `apply the Authoring Model Declaration: <error>`. Its text report
   prints `Authoring models: <h> header(s), <s> generated subagent(s), <r> removed`
   after the `Skills installed: <n>` line and its `- installed <name>` lines,
   only when a header, a generated file or a removal exists. Its JSON report
   carries `skills.authoring` with the members `headers`, `subagents` and
   `removed`, each omitted when empty and the whole member omitted when all
   three are. The stage writes the uncovered notice to its stderr. A caller
   of the stage that supplies no declaration loader applies no declaration
   and reads no configuration, so the stage's existing tests keep their
   temporary homes.
7. MUST add these tests, each top-level test calling `t.Parallel()` first,
   using temporary repositories and homes through the packages' existing
   seams, never the host's home, and never editing
   `internal/cli/cli_test.go`:
   - `internal/config`:
     `TestTheAuthoringModelDeclarationIsReadFromProjectConfig` (the example
     of `_prd.md` → User Experience, whole; an absent section and an empty
     `authoring` or `claude` mapping declare nothing without a warning);
     `TestTheAuthoringModelDeclarationRefusesWhatItDoesNotKnow`, one subtest
     per message of API Contract 1 asserting the whole error, named
     `skills_key`, `harness`, `phase_key`, `phase_field`, `research_field`,
     `empty_phase`, `empty_research`, `model_with_a_space`, `effort_hihg`
     (the effort `hihg`), `tool_Bash` (the tool `Bash`), `duplicate_tool` and
     `scalar_phase`; and
     `TestUserConfigIgnoresTheAuthoringModelDeclaration` (the warning line
     whole, no declared key, and a Project Config entry that survives beside
     it).
   - `skills`: `TestNoShippedSkillDeclaresAModelOrEffort` (every embedded and
     every canonical `SKILL.md`, failing when it read fewer than the owned
     skills twice); `TestTheDeclaredHeaderIsTheShippedSkillPlusTwoLines`
     (each authoring skill; a plain and a quoted value; model only and effort
     only; the bytes equal the shipped file with the lines removed; the two
     errors of API Contract 3);
     `TestGeneratedSubagentsMatchTheirContract` (the literal bytes of API
     Contract 4 for the four Author Subagents and for research, with and
     without each optional line);
     `TestTheResearchSubagentNeverHoldsAShellOrAWriteTool` (every one- and
     two-tool list drawn from `WebSearch`, `WebFetch`, `mcp__` names and the
     six names of `_techspec.md` → Invariant 3 goes through the
     configuration loader, which this test file may import; each list the
     loader accepts is generated, and its `tools` line holds `Read, Grep,
     Glob` and the accepted names only; the test fails when the loader
     accepted no list);
     `TestInstallAppliesTheDeclarationAndASecondRunChangesNothing`;
     `TestInstallRefusesAFileItDidNotGenerate`, with the subtests
     `hand_written_file`, `another_keys_marker`,
     `marker_below_the_first_line` and `undeclared_name`, each asserting the
     whole error and that nothing but the fixture file exists afterwards;
     `TestInstallRemovesOnlyItsOwnLeftOverSubagents` (an unmarked
     `my-reviewer.md`, a marked file at a name that is not one of the five
     and a marked file in a nested directory all survive a narrowed
     declaration and a nil one); and `TestUserLevelTargetsReceiveNoHeader`,
     with the subtests `claude`, `all` and `project_with_dir`.
   - `internal/cli`: `TestSkillsInstallAppliesTheAuthoringModelDeclaration`,
     with the subtests `surface_transcript_2` (Surface Transcript 2 whole:
     every stdout and stderr line in order and the exit code),
     `lock_records_the_installed_tree` (a `skills-lock.json` with stale
     entries for `roundfix` and `write-prd`: afterwards each entry's hash
     equals the installed tree's and the `write-prd` hash differs from the
     shipped skill's), `refuses_a_file_it_did_not_generate` (exit `2`, the
     whole stderr, nothing installed), `names_an_uncovered_header` and
     `an_invalid_declaration_stops_before_installing`;
     `TestBaselineUpdateAppliesTheAuthoringModelDeclaration` (the stage with
     the production declaration loader: the result, the lock hash, the text
     line, the JSON member, a narrowed declaration that reports one removal,
     the uncovered notice on the stage's stderr, and the two failures with
     their whole error); and
     `TestARepositoryWithoutTheDeclarationPrintsWhatItPrinted`, with the
     subtests `skills_install` and `baseline_update`: in a temporary
     repository without the section, once with no Project Config and once
     with a Project Config that loads without a warning, the whole stdout of
     `roundfix skills install --target project` is its one
     `Installed Roundfix skill for project:` line and stderr is empty, the
     installed tree equals the shipped bundle byte for byte, nothing exists
     under `.claude`, and the skills stage prints no `Authoring models:` line
     and no `authoring` member.
8. MUST document the slice in the guides: `docs/user-guide/configuration.md`
   gains the section `## Authoring models` that `_techspec.md` → API Contract
   15 describes, with the example of `_prd.md` → User Experience introduced
   by the sentence that contract gives, the validation rules of API Contract
   1, the User Config rule, what installation writes, and the contract's
   paragraph word for word; `docs/user-guide/commands/skills.md` documents
   API Contract 5 for `skills install`, showing in `text` blocks, as the
   contract writes them, the three stdout lines, the notice line and the
   refusal line, and says that a user-level target and `--dir` receive
   nothing; and `docs/user-guide/commands/baseline.md` documents the stage's
   failure, shows its text line in a `text` block as requirement 6 writes
   it, and describes its JSON member. MUST then re-record the Behavior
   Surface Record, as the Task's last step, after the last guide edit, with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
9. MUST put the new code and tests in the files `## Context` lists under
   `creates:`. MUST NOT change `.roundfixrc.yml`, any file under
   `.agents/skills/` or any shipped skill file under `skills/`, `CONTEXT.md`,
   any ADR, any Baseline file, the Makefile, the CI workflows, `go.mod`,
   `go.sum` or the lint configuration; the inspector, the three reports, the
   user-level notice and `roundfix skills remove` belong to tasks 02 and 03.
   Every test file compiles for Windows, and no file holds a real-looking
   secret.

## Subtasks

- [ ] Read and validate the section in Project Config, and ignore it in User Config.
- [ ] Add the phase table, the header writer and the subagent generator.
- [ ] Add the installation step with its refusal and removal rules.
- [ ] Load and pass the declaration in `skills install` and in the skills stage, and print their lines.
- [ ] Add the named tests.
- [ ] Write the three guides and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] A Project Config that holds the PRD's example loads into typed values, every refusal of API Contract 1 names its key or value with the contract's words, and User Config ignores the section with one warning.
- [ ] For every accepted declaration, one installation leaves each declared authoring skill equal to the shipped file plus its declared lines and each generated file equal to its contract; a second installation changes nothing; the lock records the installed tree; a file Roundfix did not generate is neither overwritten nor deleted; and no accepted configuration puts a shell, a file-writing tool or the tool that starts subagents in the Research Subagent's tools.
- [ ] A user-level target and a `--dir` target receive no header and no generated file, no shipped `SKILL.md` holds a model or an effort, and a repository without the section prints and installs what it did before.
- [ ] The three guides describe the section and what installation writes, and the Behavior Surface Record is current.

## Context

- instruction: `docs/adr/0401-each-authoring-phase-runs-on-the-model-its-repository-declares.md`
- instruction: `docs/adr/0402-the-research-subagent-reads-and-never-runs-a-command-or-writes-a-file.md`
- instruction: `internal/cli/owned_skills_lock_test.go`
- interface: `internal/config/config.go`
- creates: `internal/config/authoring.go`
- creates: `internal/config/authoring_test.go`
- interface: `skills/skills.go`
- creates: `skills/authoring.go`
- creates: `skills/authoring_test.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/baseline_update.go`
- creates: `internal/cli/authoring_models.go`
- creates: `internal/cli/authoring_models_test.go`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/user-guide/commands/skills.md`
- interface: `docs/user-guide/commands/baseline.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- Project Config's accepted keys and its load errors — consumers: `./internal/config`
- The installer's request, result and written bytes, and the output of `roundfix skills install` and of the skills stage — consumers: `./skills`, `./internal/cli`
- The configuration, skills and baseline guides and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTheAuthoringModelDeclarationIsReadFromProjectConfig|TestTheAuthoringModelDeclarationRefusesWhatItDoesNotKnow|TestUserConfigIgnoresTheAuthoringModelDeclaration)$' ./internal/config 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheAuthoringModelDeclarationIsReadFromProjectConfig TestTheAuthoringModelDeclarationRefusesWhatItDoesNotKnow/effort_hihg TestTheAuthoringModelDeclarationRefusesWhatItDoesNotKnow/tool_Bash TestTheAuthoringModelDeclarationRefusesWhatItDoesNotKnow/scalar_phase TestUserConfigIgnoresTheAuthoringModelDeclaration; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/config` — expected: exit 0. Before this Task the three tests do not exist, so their pass lines are missing; afterwards the PRD's example loads, every refusal names its key or value with the contract's words, User Config ignores the section with one warning, and the configuration suite passes.

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestNoShippedSkillDeclaresAModelOrEffort|TestTheDeclaredHeaderIsTheShippedSkillPlusTwoLines|TestGeneratedSubagentsMatchTheirContract|TestTheResearchSubagentNeverHoldsAShellOrAWriteTool|TestInstallAppliesTheDeclarationAndASecondRunChangesNothing|TestInstallRefusesAFileItDidNotGenerate|TestInstallRemovesOnlyItsOwnLeftOverSubagents|TestUserLevelTargetsReceiveNoHeader)$' ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestNoShippedSkillDeclaresAModelOrEffort TestTheDeclaredHeaderIsTheShippedSkillPlusTwoLines TestGeneratedSubagentsMatchTheirContract TestTheResearchSubagentNeverHoldsAShellOrAWriteTool TestInstallAppliesTheDeclarationAndASecondRunChangesNothing TestInstallRefusesAFileItDidNotGenerate/hand_written_file TestInstallRemovesOnlyItsOwnLeftOverSubagents TestUserLevelTargetsReceiveNoHeader/claude TestUserLevelTargetsReceiveNoHeader/project_with_dir; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./skills` — expected: exit 0. Before this Task the eight tests do not exist, so their pass lines are missing; afterwards the header and the generated files hold their contract's bytes, installation applies, refuses and removes as API Contract 5 says, no shipped skill declares a model or an effort, and the skills suite passes.

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestSkillsInstallAppliesTheAuthoringModelDeclaration|TestBaselineUpdateAppliesTheAuthoringModelDeclaration|TestARepositoryWithoutTheDeclarationPrintsWhatItPrinted)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSkillsInstallAppliesTheAuthoringModelDeclaration/surface_transcript_2 TestSkillsInstallAppliesTheAuthoringModelDeclaration/lock_records_the_installed_tree TestSkillsInstallAppliesTheAuthoringModelDeclaration/refuses_a_file_it_did_not_generate TestBaselineUpdateAppliesTheAuthoringModelDeclaration TestARepositoryWithoutTheDeclarationPrintsWhatItPrinted/skills_install TestARepositoryWithoutTheDeclarationPrintsWhatItPrinted/baseline_update; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 30m -run 'Skill|BaselineUpdate|Lock' ./internal/cli` — expected: exit 0. Before this Task the three tests do not exist, so their pass lines are missing; afterwards `roundfix skills install` prints Surface Transcript 2, the lock records the installed tree, the skills stage applies the declaration, a repository without the section prints what it printed, and the CLI's skills, Baseline update and lock tests pass.

- `grep -q '^## Authoring models$' docs/user-guide/configuration.md || { echo 'configuration guide lacks the section' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- 'The names below are an example: Roundfix ships no model name and accepts what the harness accepts.' || { echo 'configuration guide lacks the example sentence' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/configuration.md | grep -qF -- 'The section needs the Roundfix release that introduced it.' || { echo 'configuration guide lacks the older binary paragraph' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/skills.md | grep -qF -- 'Applied declared header: <skill> (model <model>, effort <effort>)' || { echo 'skills guide lacks the installation lines' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/baseline.md | grep -qF -- 'Authoring models: <h> header(s), <s> generated subagent(s), <r> removed' || { echo 'baseline guide lacks the stage line' >&2; exit 1; }; go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the configuration guide has no such section and the two command guides do not show the lines; afterwards they do, and the documentation contracts, the Behavior Surface Record included, pass.

## References

- `_prd.md` → User Stories 1, 2, 3, 4 and 7; Core Features 1, 2, 4 and 8; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 6
- `_techspec.md` → Interfaces; Data Models; API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 15; Surface Transcript 2; Invariants 1 to 4, 7 and 8; Declared breaks; Testing Approach; Build Order 1
- ADR-0401, ADR-0402


## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-context, golang-error-handling, golang-safety, golang-structs-interfaces, golang-testing, testing-boss, tech-writer, evidence-gate.

Implemented this Task's declaration and installation slice. The pre-existing
change was the Daemon's `pending` to `in_progress` status transition; it is
preserved. No Task Graph, other Task, shipped skill, Baseline asset, ADR,
Project Config or dependency/tooling configuration was edited.

Acceptance-criterion evidence:

1. Project Config decodes `skills.authoring` into the specified tagged types.
   Validation runs on the YAML nodes before strict decoding, with all twelve
   contract refusals and additional non-mapping/non-string shape cases.
   Empty authoring/harness mappings normalize to no declaration. User Config
   discards the entire node and warns once; a Project Config entry survives.
   The initial focused example test failed to compile because `Config.Skills`
   and the declaration types did not exist. After implementation,
   `go test ./internal/config -run 'Test(TheAuthoringModelDeclaration|UserConfigIgnoresTheAuthoring)' -count=1`
   exited 0, including the full PRD example, exact refusals and warning.
2. The shared installer preflights all five generated names before writing
   skills, detects uncovered headers before replacing the installed files,
   adds headers to the shipped bytes, generates the literal author/research
   bodies, and removes only its own files at those names. Generated ownership
   requires a regular file and its own marker as the first non-blank body
   line. CRLF checkouts retain ownership and are restored to LF. Identical
   files are not rewritten, preserving bytes, modes and modification times.
   `go test ./skills -run 'Test(NoShippedSkillDeclares|TheDeclaredHeader|GeneratedSubagents|TheResearchSubagent|InstallAppliesTheDeclaration|InstallRefusesAFile|InstallRemovesOnlyItsOwn|UserLevelTargets)' -count=1`
   exited 0. The tools test sends every one- and two-tool combination through
   the real config loader, generates accepted lists and rejects any generated
   shell/write/subagent-starting tool. The final CRLF refinement was checked
   with `go test ./skills -run '^TestInstallAppliesTheDeclarationAndASecondRunChangesNothing$' -count=1`
   (exit 0).
   Both CLI paths load with the resolved repository root and captured command
   home, forward loader warnings and uncovered notices to stderr, and record
   lock entries after installation. The Baseline stage has an optional loader;
   omitting it reads no configuration. Text/JSON reports omit empty authoring
   results. `go test ./internal/cli -run '^(TestSkillsInstall|TestBaselineUpdateSkillsStage|TestBaselineUpdateAppliesTheAuthoringModelDeclaration|TestARepositoryWithoutTheDeclaration)' -count=1`
   exited 0, covering Surface Transcript 2 in full, installed-tree lock hashes,
   exact refusal/load errors, Baseline text/JSON, removals, notices and existing
   skills-stage/lock behavior.
3. The skills checks above sweep every embedded and canonical owned
   `SKILL.md` for top-level model/effort keys, and verify `claude`, `all` and
   project-directory overrides receive the unchanged bundle and no generated
   file. The CLI compatibility cases cover both no Project Config and a
   warning-free Project Config without the declaration, asserting the full
   original stdout, empty stderr, shipped bytes and no `.claude` directory.
   Tests use temporary homes/repositories and the existing injected seams.
   `GOOS=windows go test -c` with explicit `/tmp/roundfix-0392-<package>.test.exe`
   outputs succeeded for `./internal/config`, `./skills` and `./internal/cli`;
   the latter two were refreshed after the final marker change. These were
   compilation checks only, not Windows runtime tests.
4. The configuration guide contains `## Authoring models`, the complete PRD
   example, exact introduction and compatibility paragraph, validation and
   User Config rules, and generated-file behavior. The skills and Baseline
   guides show the contracted streams, failures and authoring report fields.
   After the last guide edit, the sanctioned command
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
   exited 0 and regenerated only the config-key and three guide fingerprints
   in `docs/references/behavior-surfaces.json`. No record was hand-edited.

Formatting and diff whitespace checks returned no findings. Commands ran
through `rtk proxy`; shell execution required the approved escalation path
because the executor rejected the symlinked writable root, including when
its physical checkout path was supplied.

The Task's declared Verification commands, repository Verification and full
suites were not run; they remain Daemon-owned. Task status and checkboxes
remain unchanged. No commit, push, Pull Request, ai-memory message or Spec
Event Announcement was made. The inspector/reports, removal command and
shipped skill changes remain with Tasks 02–04.
