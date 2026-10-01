---
task: task_01
spec: 0214-measure-before-changing
status: pending
type: backend
complexity: high
---

# Task 01: `roundfix runs causes` says why Verification failed and why corrective Tasks were added

## Overview

Nothing in Roundfix classifies why a Verification attempt failed or why a
corrective Task was added. This Task adds a read-only store query for events
of chosen kinds, the `internal/runcause` package with the embedded signature
table, and the `roundfix runs causes` command that prints the items, the
per-Task attempt counts as JSON, and a summary by class. It documents the
command in the command reference and the Roundfix Skill. It is verifiable on
its own: against a fixture Run Database in a temporary home, the command
reproduces Surface Transcripts 1 to 4 and leaves the home byte-identical.

## Requirements

1. MUST add `Store.RunEventsOfKinds` as `_techspec.md` → Interfaces states: one Run's events whose kind is listed, in cursor order, with payloads, filtered in SQL, writing nothing.
2. MUST embed `internal/runcause/signatures.json` byte-identical to the block under `_techspec.md` → The signature table, and MUST implement `Load` with every refusal that section names and `Table.SHA256` as the SHA-256 of the embedded bytes.
3. MUST implement `Table.Classify` and `Table.Trigger` as `_techspec.md` → The signature table and Interfaces state: signatures in table order, each read only against its `source` text, the first match deciding, and `unclassified` with no signature when none matches; triggers in table order, else `unknown`.
4. MUST implement `Build` exactly as `_techspec.md` → Building the report states, steps 1 to 6: terminal Spec Runs of the repository root in the half-open window, failed attempts and their checks, the diagnostic tail read only from a regular file inside the Run's own artifact directory and never through a symbolic link, per-Task counts, corrective Tasks found by their number after the QA Task in the active or archived Task Graph, `specs_not_found`, item order and the summary.
5. MUST add `causes` to `runRunsCommand` and implement API Contracts 1 and 2: flags `--since`, `--until` and `--format`, the repository root resolved as `runs list` resolves it, the Run Database opened only through `store.OpenReader`, an absent database treated as an empty report, and the exit codes that contract names. The text form MUST reproduce Surface Transcripts 1 to 4 byte for byte, and the JSON form MUST print the document of `_techspec.md` → Data Models with no absolute path, key or diagnostic text.
6. MUST add the help of API Contract 3 to the `runs` usage and the top-level usage in `internal/cli/cli.go`, keeping every string `TestRunCommandHelp` requires.
7. MUST add `TestCausesRecordIsConsistent` in `internal/runcause/record_test.go`. With `-causes-record=<path>` and `-causes-document=<path>` it MUST fail when either file is missing, and otherwise recompute the summary from the record's items, require `signatures_sha256` to equal `Table.SHA256`, require a closed window with both ends and at least one Run, and require the document's single `Verdict: ` line to equal what `_techspec.md` → Decision rules, rule 1, gives for the recomputed counts. Without the flags it MUST run on a fixture record and document under `internal/runcause/testdata/`, and `TestCausesRecordRejectsASabotagedFixture` MUST show that a fixture whose summary or verdict disagrees fails the same check.
8. MUST describe the command under a new section `#### Why Verification failed` in `docs/user-guide/commands/runs.md` (synopsis, window, the five classes and `unclassified`, which classes count as repository knowledge, the corrective-Task rule, the triggers, that it writes nothing and makes no network request, the exit codes), and under a new heading `### Why Verification failed` in `.agents/skills/roundfix/references/runs.md` (synopsis, that it is read-only, and that an `unclassified` item is not a class). It MUST add no text inside the `### QA settlement` section of any skill.
9. MUST raise the Roundfix Skill's version in both front-matter fields of `.agents/skills/roundfix/SKILL.md` by one patch step from the value on this Task's starting commit, then run `make skills-sync`, record the version with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`, run `make baseline-digests`, and name in the Result every file those commands rewrote.
10. MUST put the tests in the files under Context, each with a temporary repository and a temporary Roundfix Home, writing fixture Runs through `store.Open` in that home only. `TestBuildLeavesTheRunDatabaseUnchanged` MUST compare the SHA-256 of every file under the temporary home before and after `Build` and after the command.
11. MUST prove each new gate can fail. The Result MUST record one sabotage of the first-match order (for example matching `implementation-defect` before `repository-convention`), one of the symbolic-link refusal, and one of the record check (for example ignoring `unclassified`), each with the test that failed, and that the code was restored.

## Subtasks

- [ ] Add the store query and its test.
- [ ] Embed the table; load, classify and trigger.
- [ ] Build the report from Runs, events and Task Graphs.
- [ ] Add the command, its text and JSON forms, and its help.
- [ ] Add the record consistency test with its fixtures.
- [ ] Document the command; describe it in the skill and record its version.
- [ ] Record one sabotage per gate in the Result.

## Acceptance Criteria

- [ ] Surface Transcripts 1 to 4 are reproduced byte for byte by `internal/cli/runs_causes_test.go`.
- [ ] Every signature and every trigger of the table has a case that matches it, and an item matching two signatures takes the earlier one.
- [ ] A failed attempt whose log is a symbolic link, or lies outside its Run's artifact directory, is classified without a diagnostic.
- [ ] A Task numbered after the QA Task is `corrective`, the QA Task and earlier Tasks are not, and a Spec found in neither root yields `corrective: null` and a `specs_not_found` entry.
- [ ] After `Build` and after the command, every file under the temporary home has the SHA-256 it had before.
- [ ] `TestCausesRecordIsConsistent` passes on its fixture and fails on the sabotaged fixture.
- [ ] The command reference and the skill reference carry their new sections, the skill's version changed and is recorded, and the docs contract tests name the new command.

## Context

- creates: `internal/store/run_events_of_kinds.go`
- creates: `internal/store/run_events_of_kinds_test.go`
- creates: `internal/runcause/signatures.json`
- creates: `internal/runcause/signatures.go`
- creates: `internal/runcause/classify.go`
- creates: `internal/runcause/report.go`
- creates: `internal/runcause/classify_test.go`
- creates: `internal/runcause/report_test.go`
- creates: `internal/runcause/record_test.go`
- creates: `internal/runcause/testdata/causes-record.json`
- creates: `internal/runcause/testdata/causes-measurement.md`
- creates: `internal/runcause/testdata/causes-record-sabotaged.json`
- creates: `internal/cli/runs_causes.go`
- creates: `internal/cli/runs_causes_test.go`
- interface: `internal/cli/runs.go`
- interface: `internal/cli/cli.go`
- interface: `docs/user-guide/commands/runs.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `.agents/skills/roundfix/references/runs.md`
- interface: `skills/roundfix/references/runs.md`
- interface: `skills/testdata/owned-skill-versions.json`
- instruction: `internal/store/store.go`
- instruction: `internal/store/journal.go`
- instruction: `internal/runevent/event.go`
- instruction: `internal/spec/spec.go`
- instruction: `internal/spec/archive.go`
- instruction: `internal/cli/cli_test.go`
- instruction: `internal/docscontract/command_documentation_test.go`
- instruction: `internal/docscontract/user_guide_contract_test.go`

## Verification

- `out="$(go test -count=1 -v -run "^(TestRunEventsOfKindsReturnsOnlyTheNamedKindsInCursorOrder|TestLoadRefusesABrokenTable|TestClassifyAppliesEachSignatureInOrder|TestTriggerNamesWhatAddedTheCorrectiveTask|TestBuildReportsAttemptsTasksAndCorrectives|TestBuildIgnoresALogOutsideTheArtifactDirectory|TestBuildLeavesTheRunDatabaseUnchanged|TestCausesRecordIsConsistent|TestCausesRecordRejectsASabotagedFixture|TestRunsCausesPrintsAClassifiedWindow|TestRunsCausesPrintsAnEmptyWindow|TestRunsCausesRefusesAMalformedDate|TestRunsCausesRequiresAGitRepository|TestRunsCausesPrintsJSON|TestRunsCausesHelp|TestRunCommandHelp|TestEveryOwnedSkillVersionIsRecorded)$" ./internal/store ./internal/runcause ./internal/cli ./skills 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestRunEventsOfKindsReturnsOnlyTheNamedKindsInCursorOrder TestLoadRefusesABrokenTable TestClassifyAppliesEachSignatureInOrder TestTriggerNamesWhatAddedTheCorrectiveTask TestBuildReportsAttemptsTasksAndCorrectives TestBuildIgnoresALogOutsideTheArtifactDirectory TestBuildLeavesTheRunDatabaseUnchanged TestCausesRecordIsConsistent TestCausesRecordRejectsASabotagedFixture TestRunsCausesPrintsAClassifiedWindow TestRunsCausesPrintsAnEmptyWindow TestRunsCausesRefusesAMalformedDate TestRunsCausesRequiresAGitRepository TestRunsCausesPrintsJSON TestRunsCausesHelp TestRunCommandHelp TestEveryOwnedSkillVersionIsRecorded; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task fifteen of the seventeen named tests do not exist, so the command fails.
- `tmp="$(mktemp -d)" || exit 1; awk 'BEGIN { fence = sprintf("%c%c%c", 96, 96, 96) } $0 == "### The signature table" { h = 1 } h && $0 == fence "json" { f = 1; next } f && $0 == fence { exit } f { print }' docs/specs/0214-measure-before-changing/_techspec.md > "$tmp/want.json" || exit 1; test -s "$tmp/want.json" || { printf 'no signature block in the TechSpec\n' >&2; exit 1; }; cmp "$tmp/want.json" internal/runcause/signatures.json` — expected: exit 0; before this Task `internal/runcause/signatures.json` does not exist, so `cmp` fails.
- `out="$(go test -count=1 -v -tags docscontract -run "^(TestEveryCommandIsNamedInTheRoundfixSkill|TestEveryCommandIsNamedInTheUserGuide)$" ./internal/docscontract 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestEveryCommandIsNamedInTheRoundfixSkill TestEveryCommandIsNamedInTheUserGuide; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; for pair in "docs/user-guide/commands/runs.md|#### Why Verification failed" "docs/user-guide/commands/runs.md|roundfix runs causes [--since <YYYY-MM-DD>]" "docs/user-guide/commands/runs.md|repository knowledge" ".agents/skills/roundfix/references/runs.md|### Why Verification failed" ".agents/skills/roundfix/references/runs.md|roundfix runs causes"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done; add="$(git log -1 --format=%H --diff-filter=A -- internal/cli/runs_causes_test.go)" || exit 1; base="${add:+$add^}"; base="${base:-HEAD}"; tmp="$(mktemp -d)" || exit 1; git show "$base:.agents/skills/roundfix/SKILL.md" > "$tmp/before.md" || exit 1; awk 'substr($0,1,8)=="version:" {print; exit}' "$tmp/before.md" > "$tmp/version-old"; awk 'substr($0,1,8)=="version:" {print; exit}' .agents/skills/roundfix/SKILL.md > "$tmp/version-new"; if cmp -s "$tmp/version-old" "$tmp/version-new"; then printf 'the skill version did not change\n' >&2; exit 1; fi; diff -r .agents/skills/roundfix skills/roundfix >/dev/null && make skills-sync-check` — expected: exit 0; before this Task no guide or skill file names `roundfix runs causes`, so the docs contract tests and the phrase checks fail.

## References

- [_prd.md](_prd.md) — Goals 1 and 5; User Stories 1 and 2; Core Features 1, 2, 3, 4 and 5; Success Metrics 1 and 2; Declared breaks
- [_techspec.md](_techspec.md) — The signature table; Interfaces; Building the report; Data Models; Decision rules; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; API Contract 1; API Contract 2; API Contract 3; API Contract 4; Testing Approach 1; Testing Approach 2; Testing Approach 3; Testing Approach 5; Build Order 1
- [references/2026-09-30-measure-why-corrective-tasks-happen.md](references/2026-09-30-measure-why-corrective-tasks-happen.md)
- ADR-0214; ADR-0215; ADR-0004; ADR-0008; ADR-0033; ADR-0035; ADR-0038; ADR-0089; ADR-0184; ADR-0187; ADR-0189
