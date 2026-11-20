---
task: task_06
spec: 0413-a-session-runs-the-flows-own-commands-without-a-permission-prompt
status: completed
type: backend
complexity: low
---

# Task 06: Every text preview names every Named Exclusion, `git -C` included

## Overview

Corrective Task 1 of 2. The QA gate of 2026-11-20 (`qa/qa-report-2026-11-20.md` on the Run Branch of `run_20261120T033653Z_96d2b5330174b3f3`, audited head `df1529f9`, finding F-01, row Q15) found that the text preview of `roundfix permissions` for Claude Code and for Codex names 15 Named Exclusions while the JSON preview of the same plan names 16: the text omits `git -C`, whose reason is `they change the repository or its configuration`. `printPermissionsExclusions` in `internal/cli/permissions.go` skips every catalog entry whose second word begins with a hyphen, so the maintainer reviewing the default text preview cannot see the whole catalog of exclusions.

task_01 Requirement 2 makes `git -C` a Named Exclusion with that reason, task_03 Requirement 2 generates the block `Left out, and still asked about:` from `Exclusions()`, and task_05 Requirement 15 requires every preview to name the Named Exclusions with their reasons. The example line that task_03 Requirement 2 and `_techspec.md` → API Contract 2 print for Git predates the classification of `git -C`; with it, the line the generator writes is the one below, and this Task makes it the line every text preview, the transcript oracle and the guide page carry:

`  git -C, config, init, merge, mv, push, remote, repack, switch, worktree: they change the repository or its configuration`

## Requirements

1. MUST answer finding F-01 of the QA gate of 2026-11-20: the block `Left out, and still asked about:` of the text preview, for Claude Code and for Codex, names every entry `Exclusions()` returns, each on the line of its reason, `git -C` included. Remove the skip of entries whose second word begins with a hyphen from `printPermissionsExclusions` in `internal/cli/permissions.go`, with its comment, and keep the rest of the generator: one line per reason, the reason of `roundfix permissions` first, later commands of a line written without the first word they share with the line's first command, and the fixed line about a commit, a push, a merge and a Pull Request last.
2. MUST print the Git line exactly as the Overview shows it, with `git -C` first because the catalog orders it first; every other line of the block keeps its bytes.
3. MUST change the transcript oracle `permissionsClaudeTranscript` in `internal/cli/permissions_test.go` so its Git line is the line of the Overview, and change no other line of the oracle. `TestPermissionsSurfaceTranscripts/claude_preview` and `TestPermissionsCommand/exclusions_printed` then pass against the corrected preview.
4. MUST add `TestPermissionsTextPreviewNamesEveryExclusion` to `internal/cli/permissions_test.go`, with the subtests `claude` and `codex`. Each runs, in a temporary Git repository with one commit and no harness file, the text preview `permissions --harness <name>` and the JSON preview `permissions --harness <name> --format json`, both through the command's entry point. It reads the lines of the text block between `Left out, and still asked about:` and the fixed line about a commit, splits each line at its last `: ` into commands and a reason, splits the commands at `, ` and writes each later command with the first word of the line's first command in front of it when it does not already start with that word. It then asserts that the set of commands the text names equals the set of `command` values of `harnesses[0].exclusions` in the JSON, that each command's reason in the text equals its `reason` in the JSON, and, as a separate assertion, that the result is not empty. A failure names each command present on one side only. The test reads no constant of the production renderer: what it compares are the two outputs of the command.
5. MUST change `docs/user-guide/commands/permissions.md` so that every place it prints or lists the Git Named Exclusions names `git -C` with the others, the Git line of its preview example being the line of the Overview.
6. MUST re-record the Behavior Surface Record after the last edit of `internal/cli/permissions.go` and the guide page with `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`, the regeneration `_authorization.md` sanctions; no fingerprint is written by hand.
7. MUST keep every existing subtest of `TestPermissionsSurfaceTranscripts` and `TestPermissionsCommand` passing with no other change, and change nothing else: no other command, flag, exit code, JSON field or Project Config key; not the Flow Command List, its Named Exclusions or Guard Entries in `internal/flowpermissions`; no skill, glossary entry, ADR, other guide page, Governed Path, Harness Permission File of this repository, the Makefile, a workflow, `go.mod`, `go.sum` or a lint configuration.

## Subtasks

- [ ] Remove the skip of hyphen-led entries from `printPermissionsExclusions`.
- [ ] Correct the Git line of the transcript oracle and add the text and JSON parity test for both harnesses.
- [ ] Correct the guide page and re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] The text previews for Claude Code and Codex name the same 16 Named Exclusions, with the same reasons, as their JSON previews, and the Git line reads as the Overview shows.
- [ ] Removing `git -C` from the text block again, or any other silent omission, makes `TestPermissionsTextPreviewNamesEveryExclusion` fail for both harnesses.
- [ ] Every existing subtest of the two permissions tests passes and the Behavior Surface Record is current.

## Context

- interface: `internal/cli/permissions.go`
- interface: `internal/cli/permissions_test.go`
- interface: `docs/user-guide/commands/permissions.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- The text preview of `roundfix permissions` and its tests — consumers: `./internal/cli`
- The guide page of `roundfix permissions` and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 40m -v -run '^(TestPermissionsSurfaceTranscripts|TestPermissionsCommand|TestPermissionsTextPreviewNamesEveryExclusion)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestPermissionsTextPreviewNamesEveryExclusion/claude TestPermissionsTextPreviewNamesEveryExclusion/codex TestPermissionsSurfaceTranscripts/help TestPermissionsSurfaceTranscripts/claude_preview TestPermissionsSurfaceTranscripts/stale_digest TestPermissionsSurfaceTranscripts/opencode_report TestPermissionsCommand/exclusions_printed TestPermissionsCommand/codex_preview_says_no_guard TestPermissionsCommand/json_preview_schema TestPermissionsCommand/home_files_untouched; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; test -z "$(gofmt -l internal/cli/permissions.go internal/cli/permissions_test.go)" || { echo 'gofmt reports the permissions files' >&2; exit 1; }` — expected: exit 0. Before this Task `TestPermissionsTextPreviewNamesEveryExclusion` does not exist, so its two subtests are missing and the command exits 1; with only the test added, both subtests fail naming `git -C`; afterwards every named subtest passes and both files are formatted.
- `tr -s '[:space:]' ' ' < docs/user-guide/commands/permissions.md | grep -qF -- 'git -C, config, init, merge, mv, push, remote, repack, switch, worktree: they change the repository or its configuration' || { echo 'the permissions guide page does not name git -C among the Named Exclusions' >&2; exit 1; }; out="$(go test -count=1 -timeout 20m -tags docscontract -v -run '^(TestTheSkillCoverageMapIsCurrent|TestOwnedSkillsNameOnlyListedOrExcludedCommands)$' ./internal/docscontract 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTheSkillCoverageMapIsCurrent TestOwnedSkillsNameOnlyListedOrExcludedCommands; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the guide page prints the Git line without `git -C`, so the command exits 1; afterwards the page names it and the Behavior Surface Record and the skill command contract are current.

## References

- QA gate of 2026-11-20 (`qa/qa-report-2026-11-20.md`, audited head `df1529f9`) → F-01, row Q15
- `_techspec.md` → API Contract 2
- task_01 → Requirement 2
- task_03 → Requirement 2
- task_05 → Requirement 15

## Result

Skills loaded: implement-task, coding-guidelines, golang-cli, golang-testing, testing-boss, tech-writer, evidence-gate, no-workarounds.

The text preview now names every entry `Exclusions()` returns, `git -C` included. Removing the hyphen-led skip in `printPermissionsExclusions` makes the Git line read `git -C, config, init, merge, mv, push, remote, repack, switch, worktree: they change the repository or its configuration`. The transcript oracle and the guide page carry the same line. A new `TestPermissionsTextPreviewNamesEveryExclusion` compares the text block with the JSON `exclusions` for `claude` and `codex` and names any command present on one side only.

Acceptance evidence:

| Criterion | Implementation and focused evidence |
| --- | --- |
| Text previews name all 16 Named Exclusions with correct reasons, Git line matches Overview | `TestPermissionsTextPreviewNamesEveryExclusion` passed for both `claude` and `codex`; the text block now matches the JSON preview's 16 entries, reasons included. `TestPermissionsSurfaceTranscripts/claude_preview` passed against the corrected transcript oracle. The guide line grep found the exact Overview Git line. |
| Omitting `git -C` fails the new test | Before the production fix, the focused new test failed for both harnesses with `JSON names "git -C", text preview does not`. After the fix, the same test passed. |
| Existing subtests pass and the Behavior Surface Record is current | `TestPermissionsSurfaceTranscripts` and `TestPermissionsCommand` passed. `TestTheSkillCoverageMapIsCurrent` and `TestOwnedSkillsNameOnlyListedOrExcludedCommands` passed after the sanctioned re-record. `gofmt -l` printed nothing for the two Go files. |

Focused checks (not the Task's declared Verification):

- `go test -count=1 -timeout 10m -run '^TestPermissionsTextPreviewNamesEveryExclusion$' ./internal/cli` — failed before the fix for both subtests, passed after.
- `go test -count=1 -run '^TestPermissionsSurfaceTranscripts$' ./internal/cli` — passed.
- `go test -count=1 -run '^TestPermissionsCommand$' ./internal/cli` — passed.
- `gofmt -l internal/cli/permissions.go internal/cli/permissions_test.go` — printed nothing.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$'` — passed after the re-record.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestOwnedSkillsNameOnlyListedOrExcludedCommands$'` — passed.
- `tr -s '[:space:]' ' ' < docs/user-guide/commands/permissions.md | grep -qF -- 'git -C, config, init, merge, mv, push, remote, repack, switch, worktree: they change the repository or its configuration'` — passed.

The Behavior Surface Record was re-recorded with `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` after the last edit; the only changed fingerprint is the guide page's. The pre-existing Task status delta is Daemon-owned and was preserved. No Task Graph or other Task file changed, and no commit, push or Pull Request was made. Declared Verification and terminal settlement remain with the Daemon.
