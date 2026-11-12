---
task: task_01
spec: 0382-an-older-roundfix-refuses-a-repository-a-newer-roundfix-wrote
status: completed
type: backend
complexity: high
---

# Task 01: Configuration loading refuses a newer Recorded Roundfix Version and the recovery commands warn

## Overview

`config.Load` reads the Recorded Roundfix Version from the repository's Setup
Manifest before Project Config and refuses a release newer than the running
binary, naming both releases and `roundfix upgrade`; `config.Init` refuses the
same way for Project Config. Doctor, the Update Check, the Run readers and
`roundfix stop` load with a new option, report the refusal and continue without
Project Config. This answers the message reported by conexus on 2026-11-11
(`references/2026-11-11-conexus-an-older-roundfix-ran-against-a-newer-repository.md`):
a 0.96.0 binary ran against a repository configured for 0.137.7.

## Requirements

1. MUST characterize first, in the Task's `## Result`: on the starting tree no
   file the repository holds records a Roundfix release, and `config.Load`
   with a Project Config holding an unknown top-level key fails with `field
   <key> not found in type config.configOverlay`.
2. MUST create `internal/config/recorded_version.go` with `SetupManifestPath`,
   `IsReleaseVersion`, `NewerRepositoryError`, `RecordedRoundfixVersion` and
   `CheckRecordedRoundfixVersion` as `_techspec.md` → Interfaces and API
   Contracts 1 and 2 say, with both messages word for word:
   `Setup Manifest "<path>" records Roundfix <recorded>, newer than this binary's <running>; a newer roundfix wrote this repository's Baseline, so use that binary or run 'roundfix upgrade'`
   and
   `Setup Manifest "<path>" records roundfixVersion <raw JSON value>, which is not a Roundfix release version such as 0.140.0; restore the value the Roundfix that wrote it recorded`.
   The reader MUST treat an absent or unreadable file, content that is not a
   JSON object and an object without the field as no value, MUST accept only
   a JSON string matching `^[0-9]+\.[0-9]+\.[0-9]+$`, and MUST write nothing.
   The check MUST skip reading for `app.IsDevelopmentVersion(running)` and
   compare with `app.CompareVersions`, refusing only a strictly newer value.
3. MUST change `internal/config/config.go` as API Contract 3 says: add
   `AllowNewerRepository bool` to `LoadOptions` and `RepositoryVersionErr
   error` to `Loaded`; in `Load`, inside a Git repository, call
   `CheckRecordedRoundfixVersion(gitRoot, app.Version)` after User Config and
   before Project Config, return its error unwrapped and read no Project
   Config, or, with `AllowNewerRepository`, store it in
   `RepositoryVersionErr`, read no Project Config and succeed; set
   `ProjectConfigPath` either way. In `Init` with scope `project`, return the
   check's error before writing.
4. MUST change `internal/cli/cli.go`, `internal/cli/doctor.go`,
   `internal/cli/health.go` and `internal/cli/update_check.go` as API Contract
   5 says: add `loadRecoveryCommandConfig` (with `AllowNewerRepository`, no
   warning of its own) and `warnRepositoryVersion`, whose one stderr line is
   word for word `roundfix: warning: <refusal>; Project Config not read`;
   `loadRunReaderCommandConfig` sets the option and warns before its
   worktree-location warning; `runStopCommand` loads with
   `loadRecoveryCommandConfig(environment, stderr)` and warns to stderr;
   `runUpdateCommand` loads with `loadRecoveryCommandConfig(environment,
   io.Discard)` and warns to stderr; `runDoctorCommand` sets the option and,
   only when `RepositoryVersionErr` is set, puts first the result named
   `HealthCheckBaselineVersion = "baseline-version"`, status failed, detail
   the refusal followed by `; Project Config not read`, next action `roundfix
   upgrade`, then runs every other check. No other command's loader changes.
5. MUST create `internal/config/recorded_version_test.go` with the six tests
   `_techspec.md` → Testing Approach names for it, each calling `t.Parallel()`
   first, with the inputs that section lists; the newer release is a test
   constant above any shipped version, and the running release reads
   `app.Version`, never a literal of a shipped version.
6. MUST create `internal/cli/recorded_version_test.go` with
   `TestACommandRefusesARepositoryRecordedByANewerRoundfix` and
   `TestReadOnlyCommandsWarnOfARepositoryRecordedByANewerRoundfix`, each
   calling `t.Parallel()` first, with the Surface Transcripts fixture, the
   subtests and the assertions Testing Approach lists: the `runs list` subtest
   asserts Surface Transcript 1 whole; the `update check` subtest asserts exit
   `0`, the whole stderr warning line and that standard output begins with the
   development build's release line; the `doctor` subtest asserts the first
   line word for word, a `node: ok` line after it and exit `1`, through
   `withDoctorFakeLoadedAndReadiness` with an empty profile proof and a
   `loadConfig` that fails the test unless `AllowNewerRepository` is set and
   then calls `config.Load`. Name the fixture helper
   `newRecordedVersionRepository(t) (repo, refusal string)` so task_02 reuses
   it.
7. MUST add to `CONTEXT.md`, through `domain-modeling`, as `_techspec.md` →
   API Contract 12 says: the **Recorded Roundfix Version** entry word for word
   right after the **Setup Manifest** entry, and the sentence "It also carries
   the Recorded Roundfix Version (ADR-0392)." at the end of the first
   paragraph of the **Setup Manifest** definition, each `ADR-0392` a Markdown
   link to
   `docs/adr/0392-an-older-roundfix-refuses-a-repository-whose-baseline-a-newer-roundfix-wrote.md`;
   MUST add no other capitalized phrase.
8. MUST change `docs/user-guide/commands/doctor.md`: right before the heading
   `## Worktree Build Caches`, the sentence "When a newer Roundfix wrote this
   repository's Baseline, Doctor prints this line first, runs every other
   check without Project Config, and exits `1`:" and this `text` block, word
   for word:

   ```text
   baseline-version: failed (Setup Manifest "<path>" records Roundfix <recorded>, newer than this binary's <running>; a newer roundfix wrote this repository's Baseline, so use that binary or run 'roundfix upgrade'; Project Config not read; next: roundfix upgrade)
   ```

9. MUST change `docs/user-guide/commands/upgrade.md`: at its end, a paragraph
   that begins, word for word, "A binary older than the Recorded Roundfix
   Version of a repository" and says that every command that loads that
   repository's configuration and `roundfix init --scope project` refuse
   there and name `roundfix upgrade`; that `roundfix upgrade` reads no
   configuration before it installs; and that `roundfix doctor`, `roundfix
   update check`, `roundfix runs`, `attach`, `events` and `stop` warn and
   continue without Project Config.
10. MUST, as its last step and after the last guide edit, re-record the
    Behavior Surface Record with
    `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
11. MUST NOT change any Baseline file, the Update Check's report format, any
   other loader, any exit code without the record, any existing test or
   golden; no test reads the host's home or real configuration.

## Subtasks

- [ ] Record the characterization.
- [ ] Add the reader, the check and its six tests.
- [ ] Call it from `Load` and `Init` with the option.
- [ ] Wire the warning commands and Doctor's line, with the two command tests.
- [ ] Add the glossary entries, the Doctor and upgrade guide text, then re-record the Behavior Surface Record.

## Acceptance Criteria

- [ ] A newer recorded release refuses every configuration-loading command before Project Config is read, and `init --scope project` leaves Project Config unchanged.
- [ ] The Run readers, `stop` and the Update Check warn and continue; Doctor fails its first line and runs every other check.
- [ ] An equal, older or absent record and a development build behave as before.
- [ ] `CONTEXT.md` defines **Recorded Roundfix Version**, and the Doctor and upgrade guides describe the refusal.

## Context

- instruction: `docs/adr/0392-an-older-roundfix-refuses-a-repository-whose-baseline-a-newer-roundfix-wrote.md`
- instruction: `internal/app/version.go`
- instruction: `internal/app/release.go`
- instruction: `internal/store/store.go`
- instruction: `internal/cli/doctor_test.go`
- instruction: `internal/cli/update_check_test.go`
- instruction: `.agents/skills/domain-modeling/SKILL.md`
- creates: `internal/config/recorded_version.go`
- creates: `internal/config/recorded_version_test.go`
- creates: `internal/cli/recorded_version_test.go`
- interface: `internal/config/config.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/health.go`
- interface: `internal/cli/update_check.go`
- interface: `CONTEXT.md`
- interface: `docs/user-guide/commands/doctor.md`
- interface: `docs/user-guide/commands/upgrade.md`
- interface: `docs/references/behavior-surfaces.json`

## Consumers

- `config.CheckRecordedRoundfixVersion`, `NewerRepositoryError` and the `LoadOptions`/`Loaded` fields — consumers: `./internal/config`, `./internal/cli`
- The refusal, the warning line and the `baseline-version:` Doctor line — consumers: `./internal/cli`
- The glossary entries, the two guides and the Behavior Surface Record — consumers: `./internal/docscontract`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestARepositoryRecordedByANewerRoundfixRefusesLoad|TestAnEqualOrOlderRecordedRoundfixVersionLoadsProjectConfig|TestARepositoryWithoutARecordedRoundfixVersionLoads|TestAMalformedRecordedRoundfixVersionRefusesLoad|TestAllowNewerRepositoryLoadsWithoutProjectConfig|TestADevelopmentBuildIgnoresTheRecordedRoundfixVersion)$' ./internal/config 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestARepositoryRecordedByANewerRoundfixRefusesLoad TestAnEqualOrOlderRecordedRoundfixVersionLoadsProjectConfig TestARepositoryWithoutARecordedRoundfixVersionLoads TestAMalformedRecordedRoundfixVersionRefusesLoad TestAllowNewerRepositoryLoadsWithoutProjectConfig TestADevelopmentBuildIgnoresTheRecordedRoundfixVersion; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done; go test -count=1 -timeout 20m ./internal/config` — expected: exit 0. Before this Task the six tests do not exist; afterwards they pass with the rest of `./internal/config`.

- `out="$(go test -count=1 -timeout 40m -v -run 'RecordedByANewerRoundfix|UpdateCheck|Doctor|Init|RunsList|Stop' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestACommandRefusesARepositoryRecordedByANewerRoundfix TestReadOnlyCommandsWarnOfARepositoryRecordedByANewerRoundfix; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the two command tests do not exist; afterwards they pass, and every existing Update Check, Doctor, Init, Run reader and Stop test of `./internal/cli` still passes.

- `grep -F -- '**Recorded Roundfix Version**' < CONTEXT.md | grep -qE '^([-*+] )?\*\*Recorded Roundfix Version\*\*(:| —( |$))' || { echo 'missing term: Recorded Roundfix Version' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'and a plan that changes none keeps the value' || { echo 'Recorded Roundfix Version definition differs' >&2; exit 1; }; tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- 'It also carries the Recorded Roundfix Version' || { echo 'Setup Manifest not revised' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/doctor.md | grep -qF -- 'Project Config not read; next: roundfix upgrade)' || { echo 'doctor guide unchanged' >&2; exit 1; }; tr -s '[:space:]' ' ' < docs/user-guide/commands/upgrade.md | grep -qF -- 'A binary older than the Recorded Roundfix Version of a repository' || { echo 'upgrade guide unchanged' >&2; exit 1; }; go test -count=1 -timeout 30m -tags docscontract ./internal/docscontract` — expected: exit 0. Before this Task the term and the guide text are absent; afterwards the glossary defines the term, the **Setup Manifest** entry names it, the two guides carry their text, and the documentation contracts, the glossary phrase contract and the Behavior Surface Record included, pass.

## References

- `_prd.md` → Goals; User Story 1; User Story 3; Core Features 2-5; Success Metric 1; Success Metric 3; Success Metric 6; Success Metric 7; Glossary
- `_techspec.md` → Current behavior 1, 4, 5, 7; Interfaces; API Contracts 1-5; API Contract 10; API Contract 12; Surface Transcript 1; Invariants 1, 2, 5; Declared breaks; Testing Approach; Build Order 1
- ADR-0392
- [The adopted message](references/2026-11-11-conexus-an-older-roundfix-ran-against-a-newer-repository.md)

## Result

Skills loaded: implement-task, coding-guidelines, agentic-cli-design, golang-cli, golang-error-handling, golang-structs-interfaces, golang-testing, testing-boss, domain-modeling, tech-writer, evidence-gate.

### Starting-tree characterization

- Before implementation, inspection found no repository compatibility release record: `docs/agents/setup-context.json`, `.roundfixrc.yml` and `skills-lock.json` contain no `roundfixVersion`; the Setup Manifest carries schema/generator versions rather than a Roundfix release. No reader or guard exists in `internal/config` or `internal/cli`.
- A temporary `go test -overlay` probe, `TestRecordedVersionStartingTreeCharacterization`, used an empty temporary home and a temporary Git repository with `future_section: {enabled: true}` in Project Config. `config.Load` returned `parse config "<repo>/.roundfixrc.yml": yaml: unmarshal errors: line 1: field future_section not found in type config.configOverlay`. The probe passed (exit 0); it wrote no repository test file.
- The initial worktree change was solely the Daemon's `status: in_progress` in this Task file; it is preserved.

### Implementation and acceptance evidence

1. Configuration-loading commands refuse before Project Config is read: added the read-only Setup Manifest reader, release predicate and typed refusal, with the specified messages and absolute path; `Load` checks after User Config, returns the refusal unwrapped, and retains `ProjectConfigPath`. `Init` checks before any project write. The new config refusal test proves the version error precedes an unknown YAML key; the command test covers `spec check x`, `follow-up` and forced project `init`, each exiting 2 with unchanged Project Config bytes. Other command loaders are unchanged.
2. Recovery commands report and continue without Project Config: the Run reader loader opts in and warns before any worktree-location warning; `stop` and the Update Check use the new recovery loader and stderr warning. The CLI regression asserts Surface Transcript 1 in full, the Update Check's exact warning and development release line with exit 0, a successful `stop` that records a Stop Request, and Doctor's exact first failure line followed by `node: ok`, exit 1 and execution of the remaining health checks. Doctor's test uses the existing fake-readiness seam, an empty profile proof and an option-checking loader that calls real `config.Load`. The Update Check report format is unchanged; its Baseline preview guard belongs to task_02.
3. Equal, older and absent records retain prior behavior: the six new config tests cover the running release via `app.Version`, an older release, absent/unreadable/non-object/non-JSON manifests, malformed fields, recovery loads with valid and unknown Project Config keys, retained User Config, and development builds that accept both newer and malformed records. Future releases use the test constant `999999.0.0`; tests use temporary homes and repositories. Existing precedence, init, stop, Run reader, Doctor and Update Check transcript checks also pass.
4. Documentation: through domain-modeling, added the exact Recorded Roundfix Version definition immediately after Setup Manifest and its required sentence with accepted ADR links. Added the prescribed Doctor sentence/block and the upgrade paragraph. Focused glossary phrase, ADR-link, Avoid-line and guide-link checks pass. No Baseline source, manifest, existing test or golden was changed.

### Focused checks

All commands below exited 0. They are implementation checks, not Daemon Verification.

- `go test -count=1 -timeout 10m ./internal/config ./internal/cli -run 'RecordedRoundfixVersion|RecordedByANewerRoundfix|AllowNewerRepository'` — new config/CLI regressions before adding the successful-stop subtest.
- `go test -count=1 -timeout 5m ./internal/config -run '^(TestLoadAppliesConfigPrecedence|TestProjectInitThenLoadWarnsNothing|TestARepositoryRecordedByANewerRoundfixRefusesLoad|TestAnEqualOrOlderRecordedRoundfixVersionLoadsProjectConfig|TestARepositoryWithoutARecordedRoundfixVersionLoads|TestAMalformedRecordedRoundfixVersionRefusesLoad|TestAllowNewerRepositoryLoadsWithoutProjectConfig|TestADevelopmentBuildIgnoresTheRecordedRoundfixVersion)$'` — all six new config tests and selected existing behavior.
- `go test -count=1 -timeout 10m ./internal/cli -run '^(TestACommandRefusesARepositoryRecordedByANewerRoundfix|TestReadOnlyCommandsWarnOfARepositoryRecordedByANewerRoundfix|TestRunStopByRunIDRecordsStopRequest|TestRunInitForceOverwritesExistingConfig|TestRunRunsListEmptyResultExitsZero|TestRunDoctorContinuesChecksAfterProfileReadinessFailure|TestUpdateCheckTranscripts)$'` — both new command tests, including successful stop, and selected existing public contracts.
- `go test -count=1 -timeout 5m -tags docscontract ./internal/docscontract -run '^(TestGlossaryPhrasesAreDefinedTerms|TestGlossaryLinksNameAcceptedADRs|TestGlossaryAvoidLinesArePlainLists|TestUserGuideLinksResolve)$'` — focused glossary and guide contracts.
- `git diff --check` — no whitespace errors.

The required Behavior Surface regeneration is the final implementation step after the last guide edit; its outcome is recorded below. Task status and declared Verification remain Daemon-owned. No commit, push or Pull Request was made.

- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage` — exited 0 after the final guide edit. The generated diff changes only the Doctor and upgrade guide fingerprints.
- Final scope inspection includes all tracked and untracked paths: exactly the Task's declared interfaces/creates and this Task file; status remains `in_progress`. `git diff --check` exited 0.
