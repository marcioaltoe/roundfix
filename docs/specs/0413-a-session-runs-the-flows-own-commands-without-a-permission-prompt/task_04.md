---
task: task_04
spec: 0413-a-session-runs-the-flows-own-commands-without-a-permission-prompt
status: completed
type: backend
complexity: high
---

# Task 04: The Doctor Command and the Update Check say where a project's list stands, and a harness's permission files are Governed Paths

## Overview

This Task makes a project whose applied list fell behind hear about it, and makes the files a harness reads its rules from tooling configuration. The Doctor Command gains the check `flow-permissions` and the Update Check the optional part `flowPermissions`; both read the plan of task_02 for each harness, print its state and the command to run, and write nothing. The governed set gains one pattern for the project permission files of the three harnesses. A test proves that no other part of Roundfix uses the list. It is verifiable on its own through the Doctor's entry point, the Update Check's tests, the governed-path tests and the skill and guide texts.

## Requirements

1. MUST add the check `flow-permissions` as `_techspec.md` → API Contract 5 states for the Doctor Command, in `internal/cli/doctor_flow_permissions.go`, with its name `HealthCheckFlowPermissions = "flow-permissions"` in `internal/cli/health.go` and its place in `internal/cli/doctor.go`: after the last Baseline check, `baseline-claude-link` or a `baseline-` check a Spec merged before this one added after it, such as the `baseline-export` and `baseline-docs` of Spec 0410, and before `residue`. Its detail joins one part per present harness with `; `; its statuses, codes, next action and the closing `; to see the list: roundfix permissions` are those of API Contract 5. Outside a Git repository the line is `flow-permissions: skipped (outside a Git repository)`. The check never returns `failed`, never changes the Doctor's exit code, reads no file outside the repository and writes nothing.
2. MUST update every test that pins the Doctor Command's ordered lines, its count of lines or a whole report on the default branch when the Task runs, among them `internal/cli/doctor_test.go`, `internal/cli/doctor_baseline_rules_test.go` and `internal/cli/doctor_claude_link_test.go` at the authoring base and the Doctor tests Spec 0410 adds, so that each expects the `flow-permissions:` line in its place; a test that pins `residue:` right after a Baseline check moves to `flow-permissions:` there and `residue:` right after it.
3. MUST add to the Doctor's help in `internal/cli/cli.go`, right after the sentence that ends "AGENTS.md and never fails Doctor.", the sentence "The flow-permissions: line reports whether the Flow Command List is applied and current for each Agent harness this repository holds, and never fails Doctor."
4. MUST add `internal/cli/doctor_flow_permissions_test.go` with `TestDoctorReportsFlowPermissions`, which drives the Doctor Command's entry point with the fakes the existing Doctor tests use, in temporary repositories, never through the built binary and never reaching the forge, with these subtests, each named as written: `no_harness_files`, `outside_a_git_repository`, `claude_not_applied`, `claude_current`, `claude_behind`, `claude_unreadable`, `codex_not_applied`, `codex_behind`, `opencode_not_needed`, `opencode_own_rules`, `behind_and_not_applied`, `behind_first_harness_named`, `unreadable_without_behind`, `never_fails`, `exit_unchanged`, `after_the_baseline_checks_before_residue`, `writes_nothing`. `claude_behind` and `behind_and_not_applied` assert the line of `_prd.md` → User Experience: `flow-permissions: warn (DR-FLOW-PERMISSIONS-BEHIND: claude: behind, <a> to add, <r> to remove; codex: not applied; next: roundfix permissions --harness claude)`, with the counts of the planted file.
5. MUST add to `internal/cli/update_check.go` the optional part `flowPermissions` of `_techspec.md` → API Contract 5, as a field of the report left out of the JSON when empty: one object per harness whose state is `current`, `behind` or `unreadable`, with `harness`, `state`, `add` and `remove` as counts and `next`, the preview command, empty when current. The text report prints, after the recommendations and before the `Next:` line, one line per object: `Flow permissions: <name> current`, `Flow permissions: <name> behind, <a> to add, <r> to remove — roundfix permissions --harness <name>` or `Flow permissions: <name> not readable — roundfix permissions --harness <name>`. The report's `next` field and its exit code do not change. The Update Check's help in `internal/cli/cli.go` gains, right after the sentence that ends "with the command that adopts it.", the sentence "It also reports each Agent harness where this repository applied the Flow Command List, and whether that list is current."
6. MUST add to `internal/cli/update_check_test.go` the test `TestUpdateCheckReportsFlowPermissions` with the subtests `left_out_when_not_applied`, `left_out_without_harness_files`, `current_text_and_json`, `behind_text_and_json`, `unreadable_text_and_json`, `codex_behind`, `next_and_exit_unchanged`, `writes_nothing`, and update every test of that file, or any other, that pins the Update Check's help, text or JSON on the default branch when the Task runs.
7. MUST add to `internal/speccheck/governed.go` one entry of kind `harness permission file`, with the tooling-authority clause and the pattern of `_techspec.md` → Data Models, the governed set: `(^|/)\.claude/settings(\.local)?\.json$|(^|/)\.codex/(rules/.*|config\.toml)$|(^|/)opencode\.jsonc?$`. No other entry changes and no path leaves the set.
8. MUST add to `TestGovernedPath` in `internal/speccheck/governed_test.go` the rows `Claude Code shared settings` (`.claude/settings.json`, governed), `Claude Code personal settings` (`.claude/settings.local.json`, governed), `Codex rules file` (`.codex/rules/roundfix-flow.rules`, governed), `Codex project configuration` (`.codex/config.toml`, governed), `OpenCode configuration` (`opencode.json`, governed), `OpenCode configuration with comments` (`opencode.jsonc`, governed), `nested Claude Code settings` (`packages/app/.claude/settings.json`, governed), `settings look-alike` (`docs/claude/settings.json.md`, not governed) and `Author Subagent file` (`.claude/agents/roundfix-tasks-author.md`, not governed); and to `TestGovernedSetOnlyGrows` in `internal/speccheck/governed_repocontract_test.go`, keeping its build tag and every path it already holds, a loop over the six governed paths above, each with a label, whose subtests are named `harness permission file <label>` with the labels `claude-shared-settings`, `claude-personal-settings`, `codex-rules-file`, `codex-config`, `opencode-json` and `opencode-jsonc`, and assert that the path is governed and carries the tooling-authority clause.
9. MUST add `internal/flowpermissions/boundary_test.go` with `TestOnlyTheCommandAndTheReportsUseTheList`, which parses the imports of every Go file under `internal/` and `cmd/` that is not a test and fails, naming the file, when one outside `internal/flowpermissions` imports the package other than `internal/cli/permissions.go`, `internal/cli/doctor_flow_permissions.go` and `internal/cli/update_check.go`; its subtests are `this_repository` and `planted_baseline_import_is_reported`. It proves `_prd.md` → Core Feature 11 for the Baseline, skill installation, setup and the Daemon.
10. MUST change the Roundfix skill, canonical copy first, then `make skills-sync` for the mirrors:
    - `references/modes.md`: in step 5 of `## --update`, right after the sentence that ends "and leave `.roundfixrc.yml` to the user as well.", the sentence of `_techspec.md` → API Contract 7: "When `flowPermissions` names a harness that is `behind` or `unreadable`, name its `next` command, and never apply it unasked.";
    - `references/setup.md`: at the end of the section `## Flow permissions`, the sentences "`roundfix doctor` reports where each harness stands on its `flow-permissions:` line, and `roundfix update check` names each harness whose applied list is behind or not readable. Neither of them applies the list.";
    - then raise the skill's version fields to the next version and record it with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`, on the tree the Task runs on. No `### QA settlement` section changes.
11. MUST add `skills/flow_permissions_reports_test.go` with `TestTheRoundfixSkillNamesTheFlowPermissionsReports`, whose subtests read both copies and pin one guarantee each: `update_mode_names_next`, `update_mode_never_applies_unasked`, `setup_names_both_reports`, `reports_never_apply`, `copies_identical`.
12. MUST change the user guide:
    - `docs/user-guide/commands/doctor.md`: in the list of checks, after the last Baseline check and before `residue:`, an item for `flow-permissions:` that names its states, `DR-FLOW-PERMISSIONS-BEHIND`, `DR-FLOW-PERMISSIONS-UNREADABLE` and the sentence "It warns when a harness is behind or not readable and never fails Doctor."; and, after the subsection of the Claude Instruction Link, a subsection `#### Flow permissions` with one example line of each status;
    - `docs/user-guide/commands/update.md`: the part `flowPermissions`, its three text lines and the sentence "The update check never applies the Flow Command List.";
    - `docs/user-guide/configuration.md`: a section `## Agent harness permission files`, right before `## Full example`, with the sentences "No Project Config key records the Flow Command List: the Harness Permission File itself is the record, and `roundfix permissions` shows or applies it." and "Each Harness Permission File is a Governed Path, so a Task that changes one needs a grant that names it."
13. MUST change `CONTEXT.md` through the `domain-modeling` skill, and no other entry: in **Governed Path**, right before the sentence that begins "The declared set is held to the record", the sentence "Each Harness Permission File is governed too, because it decides which commands an Agent session runs without asking."; and at the end of the definition of **Harness Permission File**, the sentence "Each one is a Governed Path."
14. MUST update `docs/references/skill-coverage.json` only where `TestTheSkillCoverageMapIsCurrent` requires it for the new check and the new report part, and re-record `docs/references/behavior-surfaces.json` with `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` as the Task's last step, after the last edit of a guide, a skill or a help text.
15. MUST NOT change `.roundfixrc.yml`, any file under `docs/agents/` or `internal/baseline/assets/`, `internal/cli/cli_test.go`, `docs/references/coverage-record.json`, or any Harness Permission File of this repository; the tooling-authority clause keeps its text.

## Subtasks

- [ ] Add the Doctor check in its place and move every test that pins the Doctor's lines.
- [ ] Add the Update Check part and its text lines.
- [ ] Add the governed pattern and its rows in both governed-path tests.
- [ ] Add the test that only the command and the two reports use the list.
- [ ] Add the skill, guide, help and glossary texts, raise and record the skill version, and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] For each state of `_techspec.md` → Data Models, States, the Doctor Command prints its line with the next command where one is due, warns only on behind and not readable, and never fails or changes its exit code.
- [ ] The Update Check reports each harness whose list is current, behind or not readable as text and JSON, and leaves the part out when none is applied.
- [ ] The six paths of the governed pattern are Governed Paths, every path governed at the Delivery Base stays governed, and a look-alike file and an Author Subagent file stay ordinary.
- [ ] No package but the command and the two reports imports the list, and neither report writes a file.
- [ ] Both copies of the Roundfix skill, the three guides, the two help texts and the glossary carry their texts.

## Context

- creates: `internal/cli/doctor_flow_permissions.go`
- creates: `internal/cli/doctor_flow_permissions_test.go`
- creates: `internal/flowpermissions/boundary_test.go`
- creates: `skills/flow_permissions_reports_test.go`
- interface: `internal/cli/health.go`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/doctor_test.go`
- interface: `internal/cli/doctor_baseline_rules_test.go`
- interface: `internal/cli/doctor_claude_link_test.go`
- interface: `internal/cli/update_check.go`
- interface: `internal/cli/update_check_test.go`
- interface: `internal/cli/cli.go`
- interface: `internal/speccheck/governed.go`
- interface: `internal/speccheck/governed_test.go`
- interface: `internal/speccheck/governed_repocontract_test.go`
- interface: `.agents/skills/roundfix/references/modes.md`
- interface: `skills/roundfix/references/modes.md`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/user-guide/commands/doctor.md`
- interface: `docs/user-guide/commands/update.md`
- interface: `docs/user-guide/configuration.md`
- interface: `docs/references/skill-coverage.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Consumers

- The Doctor Command's lines, their order and the check `flow-permissions` — consumers: `./internal/cli`
- The Update Check's report, text and JSON — consumers: `./internal/cli`
- `GovernedPath` and the governed set — consumers: `./internal/speccheck`
- The Roundfix skill's `modes` and `setup` references and its version — consumers: `./skills`
- `CONTEXT.md`, the guides, the help texts and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 40m -v -run 'Doctor|UpdateCheck' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in no_harness_files outside_a_git_repository claude_not_applied claude_current claude_behind claude_unreadable codex_not_applied codex_behind opencode_not_needed opencode_own_rules behind_and_not_applied behind_first_harness_named unreadable_without_behind never_fails exit_unchanged after_the_baseline_checks_before_residue writes_nothing; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestDoctorReportsFlowPermissions/$name " || { echo "missing pass: TestDoctorReportsFlowPermissions/$name" >&2; exit 1; }; done; for name in left_out_when_not_applied left_out_without_harness_files current_text_and_json behind_text_and_json unreadable_text_and_json codex_behind next_and_exit_unchanged writes_nothing; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestUpdateCheckReportsFlowPermissions/$name " || { echo "missing pass: TestUpdateCheckReportsFlowPermissions/$name" >&2; exit 1; }; done` — expected: exit 0, every Doctor and Update Check test of the package passes, the new subtests among them.
- `out="$(go test -count=1 -timeout 20m -v -run '^TestGovernedPath$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in Claude_Code_shared_settings Claude_Code_personal_settings Codex_rules_file Codex_project_configuration OpenCode_configuration OpenCode_configuration_with_comments nested_Claude_Code_settings settings_look-alike Author_Subagent_file; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestGovernedPath/$name " || { echo "missing pass: TestGovernedPath/$name" >&2; exit 1; }; done; rc="$(go test -count=1 -timeout 20m -tags repocontract -v -run '^TestGovernedSetOnlyGrows$' ./internal/speccheck 2>&1)" || { printf '%s\n' "$rc"; exit 1; }; for label in claude-shared-settings claude-personal-settings codex-rules-file codex-config opencode-json opencode-jsonc; do printf '%s\n' "$rc" | grep -qF -- "--- PASS: TestGovernedSetOnlyGrows/harness_permission_file_$label " || { echo "missing pass: harness permission file $label" >&2; exit 1; }; done` — expected: exit 0, the six paths are governed under the tooling-authority clause and the look-alikes stay ordinary.
- `out="$(go test -count=1 -timeout 20m -v -run '^TestOnlyTheCommandAndTheReportsUseTheList$' ./internal/flowpermissions 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in this_repository planted_baseline_import_is_reported; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestOnlyTheCommandAndTheReportsUseTheList/$name " || { echo "missing pass: $name" >&2; exit 1; }; done; sk="$(go test -count=1 -timeout 20m -v -run '^(TestTheRoundfixSkillNamesTheFlowPermissionsReports|TestTheRoundfixSkillDescribesFlowPermissions|TestEveryOwnedSkillVersionIsRecorded)$' ./skills 2>&1)" || { printf '%s\n' "$sk"; exit 1; }; for name in update_mode_names_next update_mode_never_applies_unasked setup_names_both_reports reports_never_apply copies_identical; do printf '%s\n' "$sk" | grep -q -- "--- PASS: TestTheRoundfixSkillNamesTheFlowPermissionsReports/$name " || { echo "missing pass: $name" >&2; exit 1; }; done; dc="$(go test -count=1 -timeout 20m -tags docscontract -run '^(TestTheSkillCoverageMapIsCurrent|TestGlossaryPhrasesAreDefinedTerms|TestOwnedSkillsNameOnlyListedOrExcludedCommands)$' ./internal/docscontract 2>&1)" || { printf '%s\n' "$dc"; exit 1; }; for pair in 'skills/roundfix/references/modes.md|command, and never apply it unasked.' '.agents/skills/roundfix/references/setup.md|Neither of them applies the list.' 'docs/user-guide/commands/doctor.md|It warns when a harness is behind or not readable and never fails Doctor.' 'docs/user-guide/commands/doctor.md|DR-FLOW-PERMISSIONS-UNREADABLE' 'docs/user-guide/commands/update.md|The update check never applies the Flow Command List.' 'docs/user-guide/configuration.md|the Harness Permission File itself is the record' 'CONTEXT.md|Each Harness Permission File is governed too, because it decides which commands an Agent session runs without asking.' 'internal/cli/cli.go|The flow-permissions: line reports whether the Flow Command List is applied and current'; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; grep -F -- '**Governed Path**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Governed Path\*\*(:| —( |$))'` — expected: exit 0, only the command and the reports import the list, and the skill, the guides, the help and the glossary carry their texts with the skill version and the Behavior Surface Record current.

## References

- `_prd.md` → Goals 3 and 5; User Stories 8, 10 and 11; Core Features 11, 12, 13, 15 and 16; Success Metric 6; Success Metric 7; Success Metric 8; Success Metric 9; Success Metric 10; User Experience; Glossary
- `_techspec.md` → The two reports; The governed set; Interfaces; Invariant 1; Data Models (States, the governed set); API Contract 5; API Contract 7; Testing Approach; Build Order 3
- ADR-0130, ADR-0444, ADR-0256


## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, domain-modeling, tech-writer, evidence-gate, systematic-debugging, no-workarounds.

Implemented the assigned report and governed-path slice. The initial tree held no Flow Command List consumer in Doctor or Update Check; the only pre-existing changed path was this Task file. Its Daemon-owned `in_progress` status is retained.

### Acceptance evidence

- Doctor reads each present harness's existing preview, prints the required state and counts, selects the first harness needing attention, and gives behind precedence for its warning code. It follows `baseline-export` and `baseline-docs` and precedes `residue`. All required `TestDoctorReportsFlowPermissions` subtests passed through CLI dispatch with temporary repositories and existing readiness fakes. Additional cases cover Codex current, mixed current/not-applied lists, missing Guard Entries, and a target linked outside the repository. Repository and isolated-home snapshots remain identical. An initial outside-Git fixture expected exit zero incorrectly: the existing `skills:` check fails outside Git. Correcting that new fixture's expectation preserved the existing product contract; the new line stays skipped.
- Update Check adds optional `flowPermissions` objects and the three text forms after recommendations and before `Next:`. All eight required `TestUpdateCheckReportsFlowPermissions` subtests passed, including explicit zero counts and empty current `next` in JSON, omitted empty parts, Codex additions/removals, unreadable files for all three harnesses, unchanged top-level `next` for both `none` and `apply`, unchanged exits, and unchanged repository/home snapshots. The release-available case uses a fresh isolated cache and never reaches the forge.
- The governed set gains exactly one pattern and retains every prior entry and the tooling-authority clause. Focused `TestGovernedPath` cases passed for all six root paths, nested Claude settings, the look-alike and the Author Subagent. The six new `TestGovernedSetOnlyGrows/harness_permission_file_*` cases passed with the existing `repocontract` tag and tooling-authority checks. The full historical loop remains for Daemon Verification.
- The import-boundary test parses production Go imports under `internal/` and `cmd/`, regardless of build tags. Both `this_repository` and `planted_baseline_import_is_reported` passed, naming the planted Baseline file while accepting the command and reports. Both reports' snapshot assertions passed without applying a plan.
- Canonical skill references were edited first and mirrored with `make skills-sync`. The required owned-version recorder advanced both Roundfix version fields from `0.1.282` to `0.1.283` and recorded the digest. All five new skill-report subtests passed for both copies. The three guide pages, two help texts, and only the two specified glossary entries carry the required wording. The Doctor coverage row names its new source; the existing Update Check source row already covers that report. The required Behavior Surface recorder ran after the final guide/help/skill edit and exited zero.

### Focused checks and required recorders

All commands below exited zero; `rtk proxy` preserved raw output for Go checks.

- `go test ./internal/cli -run '^Test(DoctorReportsFlowPermissions|UpdateCheckReportsFlowPermissions)$' -count=1 -v`
- `go test ./internal/cli -run '^Test(DoctorReportsFlowPermissions|UpdateCheckReportsFlowPermissions|RunDoctorProfileReadinessProvesEffectiveCategoriesAndReportsCounts|RunDoctorDerivesExternalSkillRequirementFromSetupManifest|DoctorReportsTheClaudeInstructionLink|UpdateCheckTranscripts)$' -count=1`
- `go test ./internal/cli -run '^Test(DoctorReportsMisplacedRepositoryRules|DoctorReportsTheDocumentationRules|RunDoctorAdapterReadinessReportsRequiredProfileRuntimes|RunDoctorRepositorySkillReadiness|RunDoctorMissingRepositoryRoot|RunDoctorRealRepositoryCheckDoesNotMutateState)$' -count=1`
- `go test ./internal/speccheck -run '^TestGovernedPath$/(Claude_Code|Codex_|OpenCode_|nested_Claude|settings_look|Author_Subagent)' -count=1`
- `go test -tags repocontract ./internal/speccheck -run '^TestGovernedSetOnlyGrows$/harness_permission_file_' -count=1`
- `go test ./internal/flowpermissions -run '^TestOnlyTheCommandAndTheReportsUseTheList$/(this_repository|planted_baseline_import_is_reported)$' -count=1 -v`
- `go test ./skills -run '^TestTheRoundfixSkillNamesTheFlowPermissionsReports$' -count=1 -v`
- `go test -tags docscontract ./internal/docscontract -run '^Test(GlossaryPhrasesAreDefinedTerms|OwnedSkillsNameOnlyListedOrExcludedCommands)$' -count=1`
- `make skills-sync`
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
- `git diff --check` and a changed-path audit: only this Task, its declared source/guide/skill/record paths, and `internal/cli/doctor_documentation_test.go`, the existing Spec 0410 report-order test explicitly covered by Requirement 2. No prohibited path changed.

The shell executor rejected the configured symlinked writable root even with the physical worktree as cwd; task commands used approved escalation. No change to executor configuration was made. Declared Verification, selected repository Verification and incremental Verification were not run; the Daemon owns them and Task settlement. No commit, push, Pull Request, other Task, Task Graph, repository Harness Permission File, Baseline asset or agent-guide mutation was made. No follow-up implementation was added.


### Verification Feedback repair — attempt 1

Inspected the Daemon diagnostic artifact `batch-004-attempt-1.log`. Its failures were the four help variants under `TestUpdateCheckUsage`: the pinned help transcript omitted the Flow Command List sentence required by Requirement 5 and already present in the command help. Updated only that expectation in `internal/cli/update_check_test.go`, as Requirement 6 requires. The full exact-output assertions, usage-error checks and no-side-effect assertions remain intact.

Focused evidence: `go test ./internal/cli -run '^TestUpdateCheckUsage$' -count=1 -v` exited zero, with all nine subtests passing, including `help`, `group_help`, `short_help` and `short_group_help`. `git diff --check` exited zero. This repair changes no production code, guide, help source, skill or generated record, so the existing recorder evidence remains applicable. The declared Verification was not rerun; Task status and settlement remain Daemon-owned.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/cli/doctor_documentation_test.go`

## Carry-forward provenance

- Source Run: `run_20261120T033653Z_96d2b5330174b3f3`
- Source commit: `df1529f9e2845738434bff1d4a99ccdb1b250426`
