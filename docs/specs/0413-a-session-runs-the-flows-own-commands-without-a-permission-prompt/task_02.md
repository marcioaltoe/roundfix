---
task: task_02
spec: 0413-a-session-runs-the-flows-own-commands-without-a-permission-prompt
status: pending
type: backend
complexity: high
---

# Task 02: A harness's file is written or cleaned only as a preview with the same digest showed it

## Overview

This Task gives `internal/flowpermissions` its file side: the editor that changes the two lists of a Claude Code settings file and keeps every other byte, the reader of an OpenCode configuration, and the plan that previews, applies and removes the list for a harness. A write recomputes the plan and goes ahead only when its digest equals the one a preview printed, so nothing reaches a harness's file that a preview of the same list and the same file did not show. It is verifiable on its own through table tests over settings files and plans in temporary repositories; no command uses it yet.

## Requirements

1. MUST add `Present`, `Preview` and `Apply`, the types `Harness`, `Action`, `State` and `Plan`, and the errors `ErrPreviewChanged`, `ErrNotStrictJSON` and `ErrUnexpectedShape` of `_techspec.md` → Interfaces, plus `ErrNotWritable` for a write to OpenCode or with an empty digest, `ErrUnreadable` for a target the harness's documented form cannot read, which wraps `ErrNotStrictJSON` and `ErrUnexpectedShape` for Claude Code, and `ErrOutsideRepository` for a target that invariant 8 refuses, in new files of the package; MUST NOT change `list.go` or `render.go`, which task_01 wrote. The plan is computed by one unexported function that takes the list value task_01 defined, so that tests can plan with a list that holds a retired rule, and the exported functions pass the shipped list.
2. MUST make `Present(repoRoot)` return, in the order `claude`, `codex`, `opencode`, each harness whose files the repository root holds: `.claude` for Claude Code, `.codex` for Codex, and one of `opencode.json`, `opencode.jsonc` and `.opencode` for OpenCode.
3. MUST add `ReadClaudeSettings(data)` and `EditClaudeSettings(data, addAllow, addAsk, remove)` as `_techspec.md` → Data Models, the Claude Code file, and invariants 3, 4 and 5 state, reading byte positions with the `encoding/json` decoder as `internal/baseline/owned_skills_lock.go` does and writing spans:
   - `nil` data is an absent file; a created file, `permissions` object or list uses two-space indentation, and a created file ends with a newline;
   - a new element is appended after the elements the list holds, in list order, and copies the separator and indentation of the element before it; a rule the list already holds is not added again;
   - a removal deletes each element with its separator, writes a list it empties as `[]`, and removes no key and no file;
   - data that is not strict JSON, that holds a comment or a trailing comma, whose `permissions`, `allow` or `ask` has another type, or whose top-level object or `permissions` object holds a key twice returns `ErrNotStrictJSON` or `ErrUnexpectedShape` and no bytes;
   - before it returns, the editor decodes its own result and compares the two lists with the ones it meant to produce; a result that does not decode or whose lists differ is an error and no bytes are returned.
4. MUST add `ReadOpenCode(repoRoot)` as `_techspec.md` → Data Models, States, says for OpenCode: it reads `opencode.json` and `opencode.jsonc` at the repository root after dropping `//` and `/* */` comments and trailing commas outside strings, returns `not_needed` when neither exists or neither holds a rule that asks about or denies a shell command, `own_rules` with the file's name when one does, and `unreadable` with the reason when a file cannot be read after that cleaning.
5. MUST compute, for a harness and an action, the state and the lists of `_techspec.md` → Data Models, States, with these targets: `.claude/settings.json` for Claude Code and `.codex/rules/roundfix-flow.rules` for Codex.
   - For Claude Code, the file is `not_applied` when it is absent or its allow list holds no `roundfix` rule of the list, current or retired; `current` when it holds every current rule and Guard Entry and no retired one; and `behind` otherwise. An apply adds each current rule and Guard Entry the file lacks and removes each retired one it holds; a removal removes every element of the allow list that equals a current or retired rule and every element of the ask list that equals a current or retired Guard Entry rule.
   - For Codex, the file is `not_applied` when it is absent, `current` when its bytes equal the rendering, `unreadable` when `ReadCodexRules` reports `ok` false, and `behind` otherwise. An apply writes the rendering whole; a removal deletes that one file.
   - An `unreadable` target carries its reason in `Plan.Reason` and has nothing to write.
6. MUST compute the digest as `_techspec.md` → Data Models, the digest, defines: the first twelve hexadecimal characters of the SHA-256 of the lines holding the harness, the action, the target path, the SHA-256 of the target's bytes or `absent`, and then each rule to add, each Guard Entry to add and each rule to remove, sorted within its group and prefixed `+`, `?` and `-`. A plan with nothing to write has an empty digest.
7. MUST make `Apply(repoRoot, harness, action, digest)` compute the plan again and write only when the digest is not empty and equals the plan's digest; otherwise it returns `ErrPreviewChanged`, or `ErrNotWritable` for an empty digest or for OpenCode, and writes nothing. It writes at most one file, through a temporary file in the target's directory and one rename; it creates `.claude` or `.codex/rules` only on a write; and it refuses with `ErrOutsideRepository`, writing nothing, a target that is a link or that has a directory on its path inside the repository that is a link resolving outside it. An `unreadable` target is never written: `Apply` returns `ErrUnreadable` with the reason before it compares the digest.
8. MUST NOT open any file under the user's home directory (invariant 10): no function of the package reads `HOME`, `os.UserHomeDir` or a path outside the repository root it is given.
9. MUST keep test fixtures as Go string literals or as files whose paths are not a Harness Permission File's path, such as `internal/flowpermissions/testdata/settings/user-rules.json`; every file named `.claude/settings.json`, `.codex/rules/...`, `opencode.json` or `opencode.jsonc` exists only inside a test's temporary directory.
10. MUST add these tests with these subtests, each named as written:
    - `TestEditClaudeSettings`, asserting bytes, never only decoded values: `absent_file_is_created`, `empty_object_gains_permissions`, `missing_ask_list_is_created`, `user_rules_kept`, `other_keys_kept`, `tab_indentation_kept`, `one_line_lists`, `no_final_newline_kept`, `windows_line_ends_kept`, `second_apply_changes_nothing`, `removal_restores_first_bytes`, `removal_keeps_user_rules`, `emptied_list_is_brackets`, `comment_is_not_strict_json`, `trailing_comma_is_not_strict_json`, `allow_of_another_type`, `duplicate_key`;
    - `TestReadOpenCode`: `no_file_is_not_needed`, `no_shell_rule_is_not_needed`, `permissions_list_asks_about_shell`, `permission_bash_ask`, `permission_value_deny`, `jsonc_comments_and_trailing_commas`, `comment_markers_inside_strings_kept`, `unreadable_file`;
    - `TestPlan`, in temporary repositories: `present_harnesses_in_order`, `claude_not_applied`, `claude_current`, `claude_behind_missing_rule`, `claude_behind_retired_rule`, `claude_git_rule_only_is_not_applied`, `claude_unreadable`, `codex_not_applied`, `codex_current`, `codex_behind`, `codex_unreadable`, `digest_covers_the_file_bytes`, `digest_is_stable`, `nothing_to_write_has_no_digest`, `home_settings_ignored`;
    - `TestApply`, in temporary repositories: `writes_exactly_the_preview`, `second_apply_writes_nothing`, `stale_digest_writes_nothing`, `empty_digest_is_refused`, `opencode_is_refused`, `unreadable_target_is_not_written`, `remove_claude_entries_only`, `remove_codex_file_only`, `link_outside_the_repository_is_refused`, `directory_link_outside_the_repository_is_refused`, `directories_created_only_on_write`. `stale_digest_writes_nothing` previews, changes one byte of the target and applies with the old digest; `home_settings_ignored` points `HOME` at a temporary directory holding a `.claude/settings.json` with a `roundfix` rule and expects the repository's plan to stay `not_applied`. The two link subtests create the link with `os.Symlink` and skip only when the platform refuses to create one.

## Subtasks

- [ ] Add the settings editor with its own check of what it wrote.
- [ ] Add the OpenCode reader.
- [ ] Add the plan, its states and its digest for both actions.
- [ ] Add the write that takes a digest, stays inside the repository and replaces one file.
- [ ] Add the table tests over settings files, plans and writes.

## Acceptance Criteria

- [ ] Over settings files that hold a user's rules, other keys, tabs, one-line lists, no final newline and Windows line ends, an apply changes no byte outside the elements it adds, a second apply changes nothing, and a removal after an apply returns the first bytes when both lists existed and held no rule of the list.
- [ ] A file that is not strict JSON, or whose lists have another type, is reported and never written.
- [ ] A write happens only with the digest of the current plan; a digest taken before the file changed writes nothing.
- [ ] A target reached through a link that leaves the repository is refused, and a settings file under the home directory is never read.
- [ ] Each state of the States table is reached by one subtest, and the OpenCode reader tells a configuration that asks about shell commands from one that does not.

## Context

- instruction: `internal/baseline/owned_skills_lock.go`
- creates: `internal/flowpermissions/claude_settings.go`
- creates: `internal/flowpermissions/claude_settings_test.go`
- creates: `internal/flowpermissions/opencode.go`
- creates: `internal/flowpermissions/opencode_test.go`
- creates: `internal/flowpermissions/plan.go`
- creates: `internal/flowpermissions/plan_test.go`

## Consumers

- `Present`, `Preview`, `Apply`, `Plan`, `ReadClaudeSettings`, `EditClaudeSettings` and `ReadOpenCode` — consumers: `./internal/flowpermissions`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestEditClaudeSettings|TestReadOpenCode)$' ./internal/flowpermissions 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in absent_file_is_created empty_object_gains_permissions missing_ask_list_is_created user_rules_kept other_keys_kept tab_indentation_kept one_line_lists no_final_newline_kept windows_line_ends_kept second_apply_changes_nothing removal_restores_first_bytes removal_keeps_user_rules emptied_list_is_brackets comment_is_not_strict_json trailing_comma_is_not_strict_json allow_of_another_type duplicate_key; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestEditClaudeSettings/$name " || { echo "missing pass: TestEditClaudeSettings/$name" >&2; exit 1; }; done; for name in no_file_is_not_needed no_shell_rule_is_not_needed permissions_list_asks_about_shell permission_bash_ask permission_value_deny jsonc_comments_and_trailing_commas comment_markers_inside_strings_kept unreadable_file; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestReadOpenCode/$name " || { echo "missing pass: TestReadOpenCode/$name" >&2; exit 1; }; done` — expected: exit 0, every named subtest passes.
- `out="$(go test -count=1 -timeout 20m -v -run '^(TestPlan|TestApply)$' ./internal/flowpermissions 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in present_harnesses_in_order claude_not_applied claude_current claude_behind_missing_rule claude_behind_retired_rule claude_git_rule_only_is_not_applied claude_unreadable codex_not_applied codex_current codex_behind codex_unreadable digest_covers_the_file_bytes digest_is_stable nothing_to_write_has_no_digest home_settings_ignored; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestPlan/$name " || { echo "missing pass: TestPlan/$name" >&2; exit 1; }; done; for name in writes_exactly_the_preview second_apply_writes_nothing stale_digest_writes_nothing empty_digest_is_refused opencode_is_refused unreadable_target_is_not_written remove_claude_entries_only remove_codex_file_only directories_created_only_on_write; do printf '%s\n' "$out" | grep -q -- "--- PASS: TestApply/$name " || { echo "missing pass: TestApply/$name" >&2; exit 1; }; done; for name in link_outside_the_repository_is_refused directory_link_outside_the_repository_is_refused; do printf '%s\n' "$out" | grep -qE -- "--- (PASS|SKIP): TestApply/$name " || { echo "missing run: TestApply/$name" >&2; exit 1; }; done` — expected: exit 0, every named subtest passes; the two link subtests pass, or skip only where the platform refuses a link.
- `test -f internal/flowpermissions/plan.go && ! grep -rqE 'UserHomeDir|Getenv\("HOME"\)' --include='*.go' --exclude='*_test.go' internal/flowpermissions && GOOS=windows go build -buildvcs=false ./internal/flowpermissions` — expected: exit 0, no source file of the package reads the home directory, and the package builds for Windows.

## References

- `_prd.md` → Goals 2 and 4; User Stories 2, 3, 7 and 9; Core Features 6, 7, 8, 9 and 10; Success Metric 3; Success Metric 4; Success Metric 7
- `_techspec.md` → The settings editor; The plan and its digest; Interfaces; Invariants 1, 2, 3, 4, 5, 8 and 10; Data Models (the Claude Code file, the Codex file, States, the digest); Testing Approach; Risks & Considerations; Build Order 1
- ADR-0444, ADR-0446
