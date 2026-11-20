---
task: task_01
spec: 0413-a-session-runs-the-flows-own-commands-without-a-permission-prompt
status: completed
type: backend
complexity: high
---

# Task 01: The Flow Command List is closed, rendered for each harness, and equal to what the owned skills name

## Overview

This Task adds the package `internal/flowpermissions` with the Flow Command List as Go data: its entries, its Named Exclusions with their reasons, its Guard Entries, the words that are no command and the retired set, and the three renderings a harness reads. It adds the two contracts that keep the list honest: no rendering holds a rule the PRD forbids, and every command an owned skill names is on the list or is a Named Exclusion. No command uses the package yet; it is verifiable on its own through its tests and the documentation contract.

## Requirements

1. MUST add the package `internal/flowpermissions` with the types and functions of `_techspec.md` → Interfaces that belong to the list and its renderings: `Entry`, `Exclusion`, `Guard`, `Entries`, `Exclusions`, `Guards`, `NotCommands`, `Retired`, `ClaudeRule`, `RenderCodexRules` and `ReadCodexRules`, and `PlainEntries() []string`, which returns each entry in plain words, such as `roundfix spec` and `gh pr view`, for the OpenCode report. The data lives in one unexported value of an unexported type, for example `type list struct` holding the entries, the exclusions, the guards and the two retired sets, in a file `internal/flowpermissions/list.go`; every exported function reads that one value, so that a later Task can compute a plan from another value in its tests without changing this file. Nothing computes the list at run time.
2. MUST derive the list on the tree the Task runs on by the two rules of `_prd.md` → Core Features 2 and 3 and ADR-0445, with `_techspec.md` → Data Models as the result at the authoring base:
   - a top-level `roundfix` subcommand is an entry when the binary registers it and an owned skill names it inside a code span, unless it is `permissions`, whose reason is `the command that applies this list`, or the Roundfix skill holds a sentence that tells a session to run it, or one of its forms, only after it asked the user or when the user asks, whose reason is `the Roundfix skill tells a session to ask you first`;
   - a Git or GitHub command an owned skill names inside a code span is an entry when, in the form the skill names it, it changes no ref, index, working-tree file or configuration of the repository and nothing on GitHub; otherwise it is a Named Exclusion, with the reason `they change the repository or its configuration` for a Git command, `it changes your credentials` for `gh auth`, and `it changes GitHub` for any other GitHub command;
   - `failed` and `has`, which follow `gh` in quoted diagnostics, are the words that are no command; a new such word is added there only when it is no command.
   - `git -C`, which the `deliver` reference of the Roundfix skill names before `branch`, `reset --hard` and `cherry-pick`, is a Named Exclusion with the reason `they change the repository or its configuration`; an option such as `-C` is never a word that is no command, because that would cover every Git command written after it.
   Specs 0414, 0409, 0410 and 0412 are delivered before this one and may make a skill name another command. The Task MUST classify every command the contract of requirement 10 reports on its tree by these rules, and record in its Result each entry, exclusion or word that differs from the TechSpec's table, with the sentence or the form that decided it.
3. MUST order `Entries()` with the `roundfix` group first, then Git, then GitHub, each group sorted, with no entry twice and no word on both sides; an entry names a command down to its subcommand, `gh pr view` included, and never an argument.
4. MUST make `ClaudeRule(entry)` return `Bash(<the entry's words> *)`, for example `Bash(roundfix spec *)`, and hold one Guard Entry per option of `_techspec.md` → Data Models, Guard Entries table, for each guarded entry that is on the list: a long option renders as `Bash(<command> *<option>*)` and a short option as `Bash(<command>* <option>*)`, where `<command>` is the entry's words, for example `Bash(git show *--output*)` and `Bash(gh api* -X*)`. A Guard Entry belongs only to an entry that is on the list.
5. MUST make `RenderCodexRules()` return the whole text of `_techspec.md` → Data Models, the Codex file: the two header lines `# Roundfix Flow Command List for Codex. Written by roundfix permissions; do not edit.` and `# To remove it: roundfix permissions --harness codex --remove`, then one `prefix_rule(` block whose pattern is `["roundfix", [<every roundfix subcommand, sorted, quoted, comma-separated>]]` with the justification `Roundfix Flow Command List: a roundfix command of the flow`, one block for the Git reads in the same form with the justification `Roundfix Flow Command List: a Git read the skills name`, and one block per GitHub entry, such as `["gh", "api"]` and `["gh", "pr", "view"]`, with the justification `Roundfix Flow Command List: a GitHub read the skills name`. Each block holds `decision = "allow",` and every line is indented as in the TechSpec. The rendering holds no rule for a Named Exclusion and no Guard Entry.
6. MUST pin that text with the golden file `internal/flowpermissions/testdata/roundfix-flow.rules.golden`, written by the test flag `-record-codex-rules` of `TestCodexRulesFile`; when the rendering and the golden differ, the test fails with `the Codex rules file changed; record it with: go test ./internal/flowpermissions -run '^TestCodexRulesFile$' -record-codex-rules`.
7. MUST make `ReadCodexRules(data)` return the entries named by the `pattern` lines of exactly the rendered form, with `ok` true, and `ok` false for any other content: a line that is neither a header line, a line of a `prefix_rule` block of that form nor empty, a `decision` other than `allow`, or an unbalanced block.
8. MUST hold the retired set empty in this first release and record every entry, in plain words, and every Claude Code rule and Guard Entry rule this release ships, one per line and sorted, in `internal/flowpermissions/testdata/released.txt`. `TestRetiredSetKeepsEveryReleasedEntry` fails when a line of that record is neither current nor retired, when a current entry or rule is missing from the record, or when a retired entry or rule is also current; its message names the line and says that a release that drops an entry adds it to the retired set and that a release that adds one appends it to the record.
9. MUST add, in `internal/flowpermissions`, these tests with these subtests, each named as written:
   - `TestFlowCommandList`: `roundfix_group_first_and_sorted`, `git_and_gh_groups_sorted`, `no_entry_twice`, `no_word_on_both_sides`, `every_exclusion_has_a_reason`, `permissions_is_excluded`, `entries_name_no_argument`;
   - `TestRenderingsHoldNoForbiddenRule`, over the Claude Code rules, the Codex text and the plain entries, with a table of forbidden first words that holds at least `bash`, `sh`, `zsh`, `python`, `python3`, `node`, `go`, `make`, `npm`, `bun` and `npx`, and refusing a `*` in the place of a command or a subcommand and any rule for a Named Exclusion: `claude`, `codex`, `opencode`, `planted_shell_is_reported`, `planted_interpreter_is_reported`, `planted_wildcard_subcommand_is_reported`, `planted_exclusion_is_reported`; the planted subtests run the same checker over a rendering of a list value that holds the planted entry and expect a finding naming it;
   - `TestGuardEntriesMatchTheirOption`, which reads each Guard Entry's rule with `*` as any text, the documented reading the PRD's outside evidence quotes: `right_after_the_command`, `after_another_argument`, `with_an_equals_sign`, `short_option_form`, `long_option_form`, and `read_without_the_option_is_not_guarded`, which expects no Guard Entry to match `git show HEAD:README.md`, `git diff --name-only HEAD`, `git grep -l hi HEAD`, `git archive --format=tar HEAD README.md` and `gh api repos/o/r/releases`;
   - `TestCodexRulesFile`: `golden`, `every_entry_allowed`, `no_exclusion_rule`, `no_guard`, `reads_back_every_entry`, `other_content_is_not_ok`;
   - `TestRetiredSetKeepsEveryReleasedEntry`: `every_recorded_line_is_current_or_retired`, `every_current_entry_and_rule_is_recorded`, `retired_is_not_current`.
10. MUST add, in `internal/docscontract/command_documentation_test.go`, keeping its build tag `docscontract` and its `//verify:always` line, the first contract of `_techspec.md` → API Contract 6:
    - `TestOwnedSkillsNameOnlyListedOrExcludedCommands` reads every Markdown file under `.agents/skills/<name>/` for each owned skill, the keys of `skills/testdata/owned-skill-versions.json`, and collects, inside inline code spans and fenced blocks, each `roundfix <word>` whose word is a top-level command the root help registers, each `git <word>`, and each `gh <word>` with its second word when the first is `pr` or `auth`. A collected command is covered when the words of an entry or of a Named Exclusion are a prefix of it, so that `gh auth login` is covered by the exclusion `gh auth`, or when the word after `git` or `gh` is a word that is no command. For each command not covered it reports `<file>: <command> is neither on the Flow Command List nor a Named Exclusion`, with the file relative to the repository root.
    - `TestAnUnlistedSkillCommandIsReported` proves the same checker over planted text and planted lists, with the subtests `unlisted_roundfix_subcommand`, `unlisted_git_command`, `unlisted_gh_pr_command`, `entry_taken_off_the_list`, which takes `spec` off the entries and expects a report naming `roundfix spec`, `excluded_command_is_covered`, `exclusion_prefix_covers_its_forms`, `unregistered_roundfix_word_is_ignored` and `not_a_command_word_is_covered`.
11. MUST keep every rendering free of a rule for a shell, an interpreter, a package manager, a build tool, a commit, a push, a merge or a Pull Request command (`_prd.md` → Core Feature 4), and MUST NOT read or write any file under the user's home directory, start a harness or run `gh`.

## Subtasks

- [ ] Write the list as data, deriving it again on this tree.
- [ ] Render it for Claude Code, Codex and OpenCode, and read the Codex file back.
- [ ] Pin the Codex rendering with its golden and the shipped entries with the release record.
- [ ] Add the tests over the list, the renderings and the Guard Entries.
- [ ] Add the documentation contract over the owned skills and its planted cases.

## Acceptance Criteria

- [ ] Every `roundfix` subcommand, Git command and GitHub command an owned skill names is an entry or a Named Exclusion, and the contract reports a planted unlisted command and an entry taken off the list by file and command.
- [ ] No rendering holds a rule for a shell, an interpreter, a wildcard in the place of a command or subcommand, or a Named Exclusion, and each planted case is reported.
- [ ] Each Guard Entry matches its option right after the command, after another argument and with an equals sign, and no Guard Entry matches a read written without the option.
- [ ] The Codex rendering equals its golden, allows every entry and holds no rule for a Named Exclusion and no guard; reading it back returns every entry.
- [ ] Outside evidence: the Claude Code permissions documentation (https://code.claude.com/docs/en/permissions, quoted in `_prd.md` → Acceptance evidence and in `references/2026-11-19-harness-rule-measurements.md`) states "Rules are evaluated in order: deny, then ask, then allow."; each Guard Entry is therefore an ask rule of its own, and row A of the measurement record shows such a rule winning over the allow entry it guards.

## Context

- instruction: `internal/baseline/owned_skills_lock.go`
- instruction: `skills/testdata/owned-skill-versions.json`
- creates: `internal/flowpermissions/list.go`
- creates: `internal/flowpermissions/render.go`
- creates: `internal/flowpermissions/list_test.go`
- creates: `internal/flowpermissions/render_test.go`
- creates: `internal/flowpermissions/testdata/roundfix-flow.rules.golden`
- creates: `internal/flowpermissions/testdata/released.txt`
- interface: `internal/docscontract/command_documentation_test.go`

## Consumers

- The Flow Command List and its renderings (`Entries`, `Exclusions`, `Guards`, `NotCommands`, `Retired`, `ClaudeRule`, `RenderCodexRules`, `ReadCodexRules`, `PlainEntries`) — consumers: `./internal/flowpermissions`, `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestFlowCommandList|TestRenderingsHoldNoForbiddenRule|TestGuardEntriesMatchTheirOption|TestCodexRulesFile|TestRetiredSetKeepsEveryReleasedEntry)$' ./internal/flowpermissions 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestFlowCommandList/roundfix_group_first_and_sorted TestFlowCommandList/git_and_gh_groups_sorted TestFlowCommandList/no_entry_twice TestFlowCommandList/no_word_on_both_sides TestFlowCommandList/every_exclusion_has_a_reason TestFlowCommandList/permissions_is_excluded TestFlowCommandList/entries_name_no_argument TestRenderingsHoldNoForbiddenRule/claude TestRenderingsHoldNoForbiddenRule/codex TestRenderingsHoldNoForbiddenRule/opencode TestRenderingsHoldNoForbiddenRule/planted_shell_is_reported TestRenderingsHoldNoForbiddenRule/planted_interpreter_is_reported TestRenderingsHoldNoForbiddenRule/planted_wildcard_subcommand_is_reported TestRenderingsHoldNoForbiddenRule/planted_exclusion_is_reported TestGuardEntriesMatchTheirOption/right_after_the_command TestGuardEntriesMatchTheirOption/after_another_argument TestGuardEntriesMatchTheirOption/with_an_equals_sign TestGuardEntriesMatchTheirOption/short_option_form TestGuardEntriesMatchTheirOption/long_option_form TestGuardEntriesMatchTheirOption/read_without_the_option_is_not_guarded TestCodexRulesFile/golden TestCodexRulesFile/every_entry_allowed TestCodexRulesFile/no_exclusion_rule TestCodexRulesFile/no_guard TestCodexRulesFile/reads_back_every_entry TestCodexRulesFile/other_content_is_not_ok TestRetiredSetKeepsEveryReleasedEntry/every_recorded_line_is_current_or_retired TestRetiredSetKeepsEveryReleasedEntry/every_current_entry_and_rule_is_recorded TestRetiredSetKeepsEveryReleasedEntry/retired_is_not_current; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0, every named subtest passes.
- `out="$(go test -count=1 -timeout 20m -tags docscontract -v -run '^(TestOwnedSkillsNameOnlyListedOrExcludedCommands|TestAnUnlistedSkillCommandIsReported|TestEveryCommandIsNamedInTheRoundfixSkill)$' ./internal/docscontract 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestOwnedSkillsNameOnlyListedOrExcludedCommands TestAnUnlistedSkillCommandIsReported/unlisted_roundfix_subcommand TestAnUnlistedSkillCommandIsReported/unlisted_git_command TestAnUnlistedSkillCommandIsReported/unlisted_gh_pr_command TestAnUnlistedSkillCommandIsReported/entry_taken_off_the_list TestAnUnlistedSkillCommandIsReported/excluded_command_is_covered TestAnUnlistedSkillCommandIsReported/exclusion_prefix_covers_its_forms TestAnUnlistedSkillCommandIsReported/unregistered_roundfix_word_is_ignored TestAnUnlistedSkillCommandIsReported/not_a_command_word_is_covered TestEveryCommandIsNamedInTheRoundfixSkill; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0, the contract passes on this repository's owned skills and reports each planted case.
- `f=internal/flowpermissions/testdata/roundfix-flow.rules.golden; grep -qx '# To remove it: roundfix permissions --harness codex --remove' "$f" && grep -q 'pattern = \["roundfix", \[' "$f" && grep -q 'pattern = \["gh", "pr", "view"\],' "$f" && ! grep -q '"permissions"' "$f" && ! grep -q '"push"' "$f" && grep -qx 'roundfix spec' internal/flowpermissions/testdata/released.txt && grep -qx 'Bash(gh api\* -X\*)' internal/flowpermissions/testdata/released.txt && GOOS=windows go build -buildvcs=false ./internal/flowpermissions` — expected: exit 0, the golden and the release record hold the shipped list and the package builds for Windows.

## References

- `_prd.md` → Goals 3 and 5; User Stories 4 and 5; Core Features 1, 2, 3, 4, 5 and 14; Success Metric 1; Success Metric 2; Success Metric 5; Acceptance evidence
- `_techspec.md` → The list and its renderings; The contracts over the list; Interfaces; Invariants 6, 7 and 9; Data Models (the list at the authoring base, Guard Entries, the Codex file); API Contract 6; Testing Approach; Build Order 1
- [`references/2026-11-19-harness-rule-measurements.md`](references/2026-11-19-harness-rule-measurements.md)
- [`references/2026-11-19-the-flows-commands-and-the-harness-permission-rules.md`](references/2026-11-19-the-flows-commands-and-the-harness-permission-rules.md)
- ADR-0445

## Result

Skills loaded: implement-task, coding-guidelines, golang-cli, golang-safety, golang-structs-interfaces, golang-testing, testing-boss, tech-writer, evidence-gate, systematic-debugging, no-workarounds.

Implemented the list half of the package: one static unexported `list` value holds 45 entries (31 Roundfix, 12 Git, two GitHub), 16 Named Exclusions, 19 individual Claude Code Guard Entries, the diagnostic words `failed` and `has`, and empty retired sets. Accessors copy nested word slices. Claude allow rules, OpenCode plain entries and the complete Codex file render from that data. The Codex reader accepts the rendered block syntax, including empty lines and older lists, and rejects other lines, malformed patterns, other decisions, incorrect justifications and unbalanced blocks. The first-release inventory records 109 sorted lines: 45 plain entries, 45 allow rules and 19 ask rules.

Derivation on this tree: read the owned-skill registry and scanned Markdown under all twelve canonical skill folders; the documentation contract obtains registered top-level commands from root help. The allow entries, four ask-first Roundfix exclusions, other Git/GitHub exclusions and diagnostic words match the TechSpec's authoring table. The sole additional classification is `git -C`, a Named Exclusion with reason `they change the repository or its configuration`: `.agents/skills/roundfix/references/deliver.md` names `git -C '<worktree>' branch 'roundfix-amended-<run-id>' HEAD`, `git -C '<worktree>' reset --hard '<first-amendment>^'` and `git -C '<worktree>' cherry-pick '<amendment>' ...`. It is an option before writes, never a diagnostic word. The existing Roundfix skill sentences about asking before applying a Baseline plan, approving profile adoption, running a release when the user asks, and removing or refreshing a user-level skill only when the user asks retain the exclusions `baseline`, `profiles`, `release` and `skills`. `permissions` remains excluded as the command that applies the list.

Focused implementation evidence, all commands below exited 0:

- `go test ./internal/flowpermissions -run '^TestCodexRulesFile$' -record-codex-rules` generated the golden through the required recorder, which declares its exact command with `suiteguard.DeclareSanctionedRegeneration` before writing.
- `go test ./internal/flowpermissions -run '^(TestFlowCommandList|TestRenderingsHoldNoForbiddenRule|TestGuardEntriesMatchTheirOption|TestRetiredSetKeepsEveryReleasedEntry|TestListAccessorsReturnIndependentData)$' -count=1` exercised ordering, disjoint classifications, exclusion reasons, argument-free entries, defensive copies, all three renderings, planted forbidden entries, guard matching, and the release inventory.
- `go test -tags docscontract ./internal/docscontract -run '^(TestOwnedSkillsNameOnlyListedOrExcludedCommands|TestAnUnlistedSkillCommandIsReported|TestOwnedSkillCommandCollectionReadsOnlyCode)$' -count=1` checked every owned skill, each required planted command and the removed `roundfix spec` entry. Findings assert the exact repository-relative file and command. The collector checks inline code and backtick/tilde fenced blocks, including indented fences, and ignores prose. An initial collector test exposed ambiguity in its indented-fence fixture; the scanner now handles those fences explicitly and the fixture expects the command inside them to be collected.
- `go test ./internal/flowpermissions -run '^TestCodexRulesFile$' -count=1`, after adding the empty-line reader case, exercised the golden, every allowed entry, absence of exclusion rules and guards, complete and older-list round trips, permitted empty lines, and malformed-content rejection.
- `go test ./internal/flowpermissions -run '^TestRenderingsHoldNoForbiddenRule$' -count=1`, after strengthening the independent Git-read oracle, checked all forbidden first words and rejected commits, pushes, merges, `git clone`, `git fetch` and GitHub PR writes even when no skill names them. Every required planted case names the offending entry in each harness rendering.
- `git diff --check` reported no whitespace errors during review.

Acceptance evidence:

1. Command coverage and planted findings: the focused documentation contract passed over the current tree and asserted the exact file-and-command findings for all required negative cases; `gh auth` and `git -C` exclusions cover their longer forms by word prefix.
2. Forbidden rules: the three rendering checks and their planted shell, interpreter, wildcard-subcommand and exclusion cases passed; the oracle also rejects a wildcard command and unnamed Git/GitHub writes.
3. Guard behavior: the independent 19-option matrix passed every required option-position, equals-sign, short-form and long-form subtest. No guard matched any of the five specified ordinary reads. Each guard belongs to a current entry and has its own ask-rule text.
4. Codex file: the recorded golden and focused reader tests passed, with all 45 entries allowed, no Named Exclusion or Guard Entry rendered, and every entry read back. The release-record tests passed with both retired sets empty.
5. Outside evidence: inspected the committed PRD's quoted permissions-documentation passage, `Rules are evaluated in order: deny, then ask, then allow.`, and row A of `references/2026-11-19-harness-rule-measurements.md`. Row A records the ordinary `git show` read running while `git show --output=probe-out.txt HEAD` is refused with a separate ask rule; row B records the same output option running without that ask rule. This is carried recorded evidence, not a fresh harness measurement. No network lookup, harness session or real `gh` invocation was performed.

The initial working-tree change was this Task's Daemon-owned `status: in_progress`; it was preserved. All new implementation paths are declared by this Task. No other Task, Task Graph manifest, owned skill, harness permission file or user-home file was changed. The execution environment rejected sandbox startup because the configured writable root contains a symlink; approved escalated shell access was used for the Task's reads, edits and focused checks.

Declared Verification, the Windows build it contains, and repository Verification remain for the Daemon. No terminal Task verdict is asserted here.

## Carry-forward provenance

- Source Run: `run_20261120T033653Z_96d2b5330174b3f3`
- Source commit: `29ea0e371184228e3f9709f336e3db3925105459`
