---
task: task_03
spec: 0194-a-skill-and-a-command-guide-read-one-command-at-a-time
status: completed
type: docs
complexity: medium
---

# Task 03: The command reference is an index plus one file per command

## Overview

`docs/user-guide/commands.md` documents every command in one file, and almost every CLI Task edits it. This Task splits it as the TechSpec's "The command reference layout" section states: `commands.md` keeps its introduction, the global contract, its group headings and an index, and each command's section moves, byte for byte, into `docs/user-guide/commands/<name>.md`. It splits whatever content the file holds on this Task's starting commit. Only relative link targets change, so each link still reaches the same document. The Verification proves, against the starting commit, that no line was lost, duplicated or changed apart from link targets.

## Requirements

1. MUST cut `docs/user-guide/commands.md` at every `### ` heading outside a fenced code block and move each such section whole to `docs/user-guide/commands/<name>.md`, named by the TechSpec's rule. Two sections with the same name share a file, in their original order.
2. MUST move a second-level section that has no `### ` child the same way, except `Global contract` and `Agent boundaries`, which stay.
3. MUST keep in `commands.md` its title and introduction, `Global contract`, each group heading with the text before its first `### `, `Agent boundaries`, and a command index between the lines `<!-- roundfix:command-index:begin -->` and `<!-- roundfix:command-index:end -->` with one row per command file. The index is the only new text this Task writes.
4. MUST rewrite every relative link inside a moved section so it reaches the same target from the new directory: a relative target gains one leading `../`, and an in-page anchor to a section that moved becomes that section's file. A link with a URL scheme is unchanged. It MUST rewrite each link elsewhere in the repository that points into a moved section, which includes the three in `docs/user-guide/usage.md`, to the command's file.
5. MUST NOT change, add, remove or re-wrap anything else on any line. It MUST perform the move with a throwaway script or tool that copies lines, never by retyping. That script is not part of the change and stays out of the repository.
6. MUST put the new tests in `internal/docscontract/commands_index_test.go`, under the `docscontract` build tag. They cover the index and the links whose target lies under `docs/user-guide/`.
7. MUST edit no governed file, rename or remove no top-level test, and change no command's documented behavior.
8. MUST record in its Result the name of every command file it created that this Task does not declare, which can happen when an earlier Spec added a command section.
9. MUST NOT name any decision by its `ADR-` identifier in this Task's `## Result`.

## Subtasks

- [ ] Move every command section to its file with a script, rewriting only link targets.
- [ ] Write the command index and update the links that pointed into moved sections.
- [ ] Add the index and link tests.

## Acceptance Criteria

- [ ] Every non-blank line of the command reference on the starting commit is present, with the same count, in `commands.md` and the command files together, outside the index block, when link targets are ignored.
- [ ] The index names every file under `docs/user-guide/commands/` and names no file that does not exist.
- [ ] Every relative link in `commands.md`, the command files and `usage.md` whose target lies under `docs/user-guide/` names a file that exists.
- [ ] Every contract that pins the command reference's text passes while reading the index with its command files.

## Context

- interface: `docs/user-guide/commands.md`
- creates: `docs/user-guide/commands/setup.md`
- creates: `docs/user-guide/commands/doctor.md`
- creates: `docs/user-guide/commands/migrate.md`
- creates: `docs/user-guide/commands/review.md`
- creates: `docs/user-guide/commands/deliver.md`
- creates: `docs/user-guide/commands/upgrade.md`
- creates: `docs/user-guide/commands/init.md`
- creates: `docs/user-guide/commands/gc.md`
- creates: `docs/user-guide/commands/skills.md`
- creates: `docs/user-guide/commands/profiles.md`
- creates: `docs/user-guide/commands/baseline.md`
- creates: `docs/user-guide/commands/fetch.md`
- creates: `docs/user-guide/commands/resolve.md`
- creates: `docs/user-guide/commands/watch.md`
- creates: `docs/user-guide/commands/review-report-shape.md`
- creates: `docs/user-guide/commands/window.md`
- creates: `docs/user-guide/commands/implement.md`
- creates: `docs/user-guide/commands/reopen.md`
- creates: `docs/user-guide/commands/settle.md`
- creates: `docs/user-guide/commands/archive.md`
- creates: `docs/user-guide/commands/qa-report.md`
- creates: `docs/user-guide/commands/supersede.md`
- creates: `docs/user-guide/commands/reconcile.md`
- creates: `docs/user-guide/commands/runs.md`
- creates: `docs/user-guide/commands/attach.md`
- creates: `docs/user-guide/commands/events.md`
- creates: `docs/user-guide/commands/stop.md`
- creates: `docs/user-guide/commands/detached-runs.md`
- interface: `docs/user-guide/usage.md`
- creates: `internal/docscontract/commands_index_test.go`
- instruction: `docs/adr/0187-the-roundfix-skill-and-the-command-reference-are-read-one-command-at-a-time.md`

## Verification

- `out="$(go test -count=1 -tags docscontract -v -run "^(TestCommandIndexNamesEveryCommandFile|TestCommandReferenceLinksResolve|TestBaselineDocumentationContract|TestProfilesDocumentationContractMatchesPublicGuidance)$" ./internal/docscontract 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestCommandIndexNamesEveryCommandFile TestCommandReferenceLinksResolve TestBaselineDocumentationContract TestProfilesDocumentationContractMatchesPublicGuidance; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; add="$(git log -1 --format=%H --diff-filter=A -- docs/user-guide/commands/deliver.md)" || exit 1; base="${add:+$add^}"; base="${base:-HEAD}"; tmp="$(mktemp -d)" || exit 1; git show "$base:docs/user-guide/commands.md" > "$tmp/before.md" || exit 1; sed 's/](\([^)]*\))/]()/g' "$tmp/before.md" | awk 'NF {print}' | LC_ALL=C sort > "$tmp/old"; cat docs/user-guide/commands.md docs/user-guide/commands/*.md | awk '$0=="<!-- roundfix:command-index:begin -->" {skip=1; next} $0=="<!-- roundfix:command-index:end -->" {skip=0; next} !skip && NF {print}' | sed 's/](\([^)]*\))/]()/g' | LC_ALL=C sort > "$tmp/new"; cmp "$tmp/old" "$tmp/new" || { printf 'the split lost, duplicated or changed a line\n' >&2; exit 1; }` — expected: exit 0; before this Task the two new named tests do not exist, so the command fails.
- `test -f docs/user-guide/commands/deliver.md || { printf 'the command files are missing\n' >&2; exit 1; }; pins="$(go test -count=1 -v -run "^(TestEventsHelpDocumentsAgentSelectionFilter|TestNoPythonBaselineRuntime)$" ./internal/cli ./skills 2>&1)" || { printf "%s\\n" "$pins"; exit 1; }; for name in TestEventsHelpDocumentsAgentSelectionFilter TestNoPythonBaselineRuntime; do printf "%s\\n" "$pins" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task no command file exists, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 2; Goal 3; Core Feature 2; Success Metric 2, Success Metric 3
- [_techspec.md](_techspec.md) — The command reference layout; The move proof; Testing Approach 4; Build Order 3
- ADR-0187

## Result

Implemented the command-reference split from the task-start `HEAD` with a
throwaway line-copying script. `commands.md` now retains its introduction,
global contract, group headings, agent boundaries, and a marked index; moved
sections are stored under `docs/user-guide/commands/`. Relative links in moved
sections gain the required directory prefix, and the three command-section
links in `usage.md` now point to command files. The two command files created
from earlier-wave sections that were not declared in this Task are
`spec-audit.md` and `storage-report.md`.

Added `internal/docscontract/commands_index_test.go` under the `docscontract`
build tag. It checks that the index and command directory have the same file
set and that relative links targeting the user guide resolve from
`commands.md`, every command file, and `usage.md`.

Focused evidence after the final edit:

- `GOCACHE=/private/tmp/roundfix-task03-gocache go test -count=1 -tags docscontract ./internal/docscontract -run '^$'` passed; the docscontract package compiles with the new tests.
- `GOCACHE=/private/tmp/roundfix-task03-gocache go test -count=1 -tags docscontract ./internal/docscontract -run '^(TestUserGuideLinksResolve|TestEveryCommandIsNamedInTheUserGuide)$'` passed.
- `python3 /private/tmp/command_move_check.py` passed; the non-blank line multiset matches `HEAD` after link-target normalization and index removal.
- `git diff --check` passed.

The daemon must run the declared Verification commands and settle the Task.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `docs/user-guide/commands/spec-audit.md`
- `docs/user-guide/commands/storage-report.md`

## Carry-forward provenance

- Source Run: `run_20261001T010622Z_8bfe4bebad66f971`
- Source commit: `ecc91734e754867e8c37e847845d3e6e05f636c8`
