---
task: task_01
spec: 0214-measure-before-changing
status: completed
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
reproduces Surface Transcripts 1 to 4 and leaves the home byte-identical, except for the `roundfix.db-wal` and `roundfix.db-shm` sidecars SQLite creates when `store.OpenReader` opens a WAL database read-only; the database file, its lock and every artifact log stay byte-identical.

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

## Result

Implemented the task_01 slice for Daemon Verification; status remains
Daemon-owned. No declared Verification command was run, and no commit, push
or pull request was made. The only pre-existing changed path was this Task
file, whose Daemon-written `status: in_progress` was preserved.

### Implementation

- Added the SQL-filtered, parameter-bound `Store.RunEventsOfKinds` query,
  preserving raw payloads and per-Run cursor order without writes.
- Embedded the exact signature table and its SHA-256;
  `Load` rejects unknown classes/sources, duplicate IDs and invalid patterns.
  Classification and triggers preserve table order and source separation.
- Built terminal, repository-scoped, half-open-window reports with one item
  per failed attempt, per-Task verdict/feedback/Run counts, historical graph
  membership, one corrective item per Task, missing-Spec reporting and summary.
  Diagnostic reads are bounded to the last 65,536 bytes of regular files in
  `<artifact_dir>/runs/<run-id>/`, rejecting leaf and parent symbolic links.
  Task evidence is limited to 4,096 bytes of title and Overview.
- Added `runs causes`, text/JSON output, UTC date validation, help and exit
  contracts. It opens only `store.OpenReader`; an absent database produces an
  empty report. Classification sees original commands; exported checks redact
  absolute paths and credential assignments/flags. No diagnostic or Task text
  is serialized. A missing Spec Root permits archive lookup and missing-Spec
  reporting instead of requiring execution eligibility.
- Added `internal/spec/cause_graph.go` as a necessary ordinary path in this
  slice: the existing execution loader requires an active PRD, so the new
  read-only projection uses the Spec package's manifest parser for archived
  graph membership and title/Overview evidence. No execution loader or
  settlement behavior changed.
- Added the record consistency harness, valid and sabotaged synthetic records,
  documentation sections and explicit docs-contract command anchors.

### Acceptance evidence

| Acceptance criterion | Focused evidence |
| --- | --- |
| Surface Transcripts 1–4 match byte for byte | `TestRunsCausesPrintsAClassifiedWindow`, `TestRunsCausesPrintsAnEmptyWindow` (absent and present databases), `TestRunsCausesRefusesAMalformedDate`, and `TestRunsCausesRequiresAGitRepository` compare full stdout, stderr and exit codes against the authored transcripts. |
| Every signature/trigger and first-match order | `TestClassifyAppliesEachSignatureInOrder` covers all 12 signatures, overlapping golden/lint/test failures, source isolation, absent evidence and `any` source order. `TestTriggerNamesWhatAddedTheCorrectiveTask` covers all three triggers, precedence and unknown. `TestLoadRefusesABrokenTable` covers all specified refusals and the embedded-byte digest. |
| Unsafe diagnostic files are absent evidence | `TestBuildIgnoresALogOutsideTheArtifactDirectory` covers leaf symlinks, parent symlinks, outside paths, another Run's directory, a directory and a missing file. `TestBuildBoundsEvidenceAndExportsNoPrivateCommandText` covers tail/text limits and command privacy. |
| Corrective/QA boundary, archive and missing graphs | `TestBuildReportsAttemptsTasksAndCorrectives` checks earlier Tasks, QA, later Tasks, active/archive equivalence, null corrective fields and `specs_not_found`. `TestBuildWindowAndMultipleRuns` checks distinct-Run counts, first-Run corrective attachment, nonmembers, other repositories, Active Runs and inclusive/exclusive window boundaries. |
| Home bytes preserved after Build and command | `TestBuildLeavesTheRunDatabaseUnchanged` exists in both `internal/runcause` and `internal/cli`. The CLI test hashes every file before Build, after Build, after JSON command and after reader close, allowing only new SQLite WAL/SHM sidecars. Existing database, lock, logs and all other file hashes must match. Build and command reports must also agree. |
| Saved record passes and sabotage fails | `TestCausesRecordIsConsistent`, `TestCausesRecordRejectsASabotagedFixture` and `TestCausesRecordDecisionBoundaries` validate recomputed summary, digest, closed window, Run count and the single decision-rule verdict. Explicit valid flags passed; a missing record and a missing document each exited 1. |
| Command reference, skill and version contracts | The two new Why Verification failed sections are present; both Roundfix Skill version fields rose from 0.1.14 to 0.1.15. `TestEveryCommandIsNamedInTheRoundfixSkill` and `TestEveryCommandIsNamedInTheUserGuide` explicitly require `runs causes`. Focused docs-contract checks passed; source/mirror bytes match and QA settlement is byte-identical to HEAD. |

The store-focused `TestRunEventsOfKindsReturnsOnlyTheNamedKindsInCursorOrder`
also proves SQL exclusion by corrupting an excluded event's timestamp: a
reader that scanned and filtered afterward would fail. It checks another Run,
duplicate requested kinds, empty selections, cursor order and exact payloads.
All fixture writes use temporary repositories and temporary Roundfix Homes.

### Commands and outcomes

All Go checks used `GOCACHE=/tmp/roundfix-task01-cache`.

- Starting-commit evidence: `git cat-file -e HEAD:internal/cli/runs_causes.go`
  exited 128; the new command does not exist in the starting revision.
- `go test ./internal/runcause ./internal/store -run
  'Test(Classify|Trigger|LoadRefuses|RunEventsOfKinds)' -count=1`: exit 0.
- `go test ./internal/runcause -run 'TestBuild' -count=1`: exit 0.
- `go test ./internal/runcause -run 'TestCausesRecord' -count=1`: exit 0.
- `go test ./internal/cli -run
  'TestRunsCauses|TestBuildLeaves|TestRunCommandHelp' -count=1`: exit 0.
- After restoring sabotages and adding privacy/bounds checks,
  `go test ./internal/runcause ./internal/cli ./internal/store -run
  'Test(Build|RunsCauses|Classify|Trigger|LoadRefuses|CausesRecord|RunEventsOfKinds)'
  -count=1`: exit 0, 51 checks including subtests.
- `go test -tags docscontract ./internal/docscontract -run
  'TestEveryCommand|TestCommandPaths' -count=1`: exit 0, three tests.
- `go test ./internal/runcause -run '^TestCausesRecordIsConsistent$'
  -causes-record=testdata/causes-record.json
  -causes-document=testdata/causes-measurement.md -count=1`: exit 0.
  Replacing the record path with an absent temporary path exited 1 with
  `read cause record`; replacing the document path exited 1 with
  `read cause document`.
- A Python byte comparison of the TechSpec's JSON block and embedded file
  passed; SHA-256 is
  `53e5e5a764109b92806c6efe4d88603c13e5de2c6d21816755406deb11b6dd5f`.
  Source/mirror and unchanged QA settlement comparisons also passed.
- `make verify-incremental`: the sandbox run exited 2 because two existing
  force-stop integration tests could not enumerate the process table.
  The permitted rerun exited 0, including formatting, Go vet, all Go tests,
  skill sync/check and binary build. No product change was made for the
  sandbox restriction.
- `git -c core.fsmonitor=false diff --check`: exit 0. Scope inspection found
  no Task Graph, other Task or archived-file changes.

### Sabotage evidence

Each mutation was restored from its original bytes in a `finally` block.
Subsequent focused checks and incremental verification used restored code.

| Gate | Deliberate mutation | Observed failure |
| --- | --- | --- |
| First-match order | Reversed signature iteration in `classify.go`, allowing implementation-defect to beat repository-convention. | `TestClassifyAppliesEachSignatureInOrder/record-or-golden` and `/lint-or-format` failed with implementation-defect/go-test-failure; the killed/timed-out overlap also failed. Test command exited 1. |
| Symbolic-link refusal | Replaced `os.Lstat(current)` with `os.Stat(current)` in `report.go`. | `TestBuildIgnoresALogOutsideTheArtifactDirectory/leaf_symlink` and `/parent_symlink` failed because the linked log classified as environment. Test command exited 1. |
| Record verdict | Removed the unclassified-share condition from the decision rule in `record_test.go`. | `TestCausesRecordIsConsistent` failed: recomputed verdict became reopen, while the document says inconclusive for 13 unclassified items out of 36. Test command exited 1. |

### Skill regeneration

- `make skills-sync`: exit 0. Byte changes were
  `skills/roundfix/SKILL.md` and `skills/roundfix/references/runs.md`.
  The target physically recopied every owned-bundle file listed below;
  every other copied file retained its original bytes.
- `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$'
  -record-skill-versions`: exit 0; rewrote
  `skills/testdata/owned-skill-versions.json` with Roundfix 0.1.15 and its digest.
- `make baseline-digests`: exit 0, `changed:false`; no derived digest file
  gained a byte change.

Files recopied by `make skills-sync`:

```text
skills/archive-spec/SKILL.md
skills/brainstorming/SKILL.md
skills/business-analyst/SKILL.md
skills/council/SKILL.md
skills/council/assets/synthesis-template.md
skills/council/references/archetypes.md
skills/council/references/debate-protocols.md
skills/evidence-gate/SKILL.md
skills/implement-spec/SKILL.md
skills/implement-task/SKILL.md
skills/qa-gate/SKILL.md
skills/roundfix/SKILL.md
skills/roundfix/agents/openai.yaml
skills/roundfix/references/archive.md
skills/roundfix/references/baseline.md
skills/roundfix/references/deliver.md
skills/roundfix/references/events.md
skills/roundfix/references/implement.md
skills/roundfix/references/profiles.md
skills/roundfix/references/reconcile.md
skills/roundfix/references/release.md
skills/roundfix/references/review-runs.md
skills/roundfix/references/review.md
skills/roundfix/references/runs.md
skills/roundfix/references/runtime.md
skills/roundfix/references/settle.md
skills/roundfix/references/setup.md
skills/roundfix/references/spec-delivery.md
skills/roundfix/references/spec.md
skills/roundfix/references/stop.md
skills/roundfix/references/storage.md
skills/setup-context-driven/SKILL.md
skills/write-idea/SKILL.md
skills/write-idea/references/idea-template.md
skills/write-idea/references/opportunity-scan.md
skills/write-prd/SKILL.md
skills/write-prd/references/prd-template.md
skills/write-tasks/SKILL.md
skills/write-tasks/references/task-template.md
skills/write-techspec/SKILL.md
skills/write-techspec/references/concrete-contracts.md
skills/write-techspec/references/techspec-template.md
```

No follow-up from another Task was implemented. Declared Verification and
Task settlement remain with the Daemon.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/docscontract/command_documentation_test.go`
- `internal/docscontract/user_guide_contract_test.go`
- `internal/spec/cause_graph.go`

## Carry-forward provenance

- Source Run: `run_20261002T171516Z_94169d908e48d811`
- Source commit: `424ad384a3fb762d5ea13402d371a84f950c7830`
