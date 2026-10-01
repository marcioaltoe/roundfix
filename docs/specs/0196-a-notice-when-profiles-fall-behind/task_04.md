---
task: task_04
spec: 0196-a-notice-when-profiles-fall-behind
status: completed
type: backend
complexity: high
---

# Task 04: Doctor and `upgrade` carry the notice, and nothing else reaches the comparison

## Overview

The comparison of task_02 is seen only by someone who runs `roundfix profiles check`. This Task puts it where a maintainer already looks: one Doctor line, and a notice on standard error after every release outcome of `roundfix upgrade` (ADR-0181). Neither may change what the hosting command prints on standard output or its exit code, so the upgrade behavior is characterized before it is touched. After an install, the notice comes from the installed executable. A source test keeps every other command away from the comparison.

## Requirements

1. MUST first write `internal/cli/upgrade_characterization_test.go` against the unchanged code. It MUST fix the exact standard output and the exit code of: no release published, already current, `--check` with a newer release, `--check` when current, a release installed, a failed download, an unknown flag and `--help`. It MUST assert nothing about standard error on the four release outcomes. The Result MUST state that it passed before any production change.
2. MUST give `runUpgradeCommand` the `commandEnvironment` its dispatcher already holds, and add `installedProfilesCheck` to `upgradeDependencies`, as `_techspec.md` → Interfaces states. The default MUST run the installed path with the arguments `profiles check`, in the process working directory, under a ten-second timeout, and return what it printed.
3. MUST write the notice to standard error after the outcome line, under the rules of `_techspec.md` → Upgrade: through `installedProfilesCheck` after an install, and in process otherwise, with the renderer task_02 added in `internal/cli/profiles_check.go`. It MUST NOT call `installedProfilesCheck` when nothing was installed.
4. MUST keep standard output and the exit code of every case of Requirement 1 byte-identical. A comparison that cannot run, in process or in the child, MUST write exactly one line, `roundfix: recommendations not checked: <reason>`, and MUST NOT change the exit code. Help, a usage error and a failed upgrade MUST write no notice.
5. MUST append the Doctor result `recommendations` right after the `profiles` result, with the statuses and detail of `_techspec.md` → Doctor and → Surface Transcripts. It MUST be computed from the configuration Doctor already loaded, call no runner, and never be `failed`.
6. MUST update only the existing tests this change invalidates, and name each in the Result. The measurement recorded in `_prd.md` → Acceptance evidence found five in `internal/cli/doctor_test.go` and two in `internal/cli/upgrade_test.go`. It MUST rename or remove no top-level test and MUST NOT edit `internal/cli/cli_test.go`.
7. MUST add `internal/cli/recommendation_check_scope_test.go`, which parses the non-test Go files under `internal` and `cmd` and fails when `CheckRecommendations` is referenced outside the files `_techspec.md` → Testing Approach 6 lists. It MUST fail, not skip, when it finds no reference at all.
8. MUST put the other new tests in `internal/cli/upgrade_notice_test.go` and `internal/cli/doctor_recommendations_test.go`. The upgrade tests MUST use the existing release fixtures and a fake `installedProfilesCheck` that records its arguments; no test may start a real child process, download a release or reach the network.
9. MUST prove each new gate can fail. The Result MUST record one sabotage for the exit-code rule (for example returning the child's failure), one for the scope test (for example calling the comparison from `internal/cli/implement.go`), and one for the Doctor rule (for example marking a difference `failed`), each with the test that failed, and that the code was restored.
10. MUST mention the notice in the help text of `upgrade` and the new line in the help text of `doctor` in `internal/cli/cli.go`, keeping every string `TestRunCommandHelp` requires.
11. MUST document both in `docs/user-guide/commands.md`: the `recommendations:` line in the Doctor sample output and its two statuses, and, under `upgrade`, the recommendation notice, where it is written, that it never changes standard output or the exit code, and that an upgrade performed by an older executable prints none. It MUST add one sentence on the recommendation notice to `docs/user-guide/usage.md`.
12. MUST add the Doctor line and the recommendation notice to the text under the skill heading task_01 created. It MUST put the skill text under the heading `### Recommendation check` of `.agents/skills/roundfix/SKILL.md`, add no text inside the `### QA settlement` section, then run `make skills-sync` and `make baseline-digests`, and name in the Result every file either command rewrote.
13. MUST, because this Task changes `.agents/skills/roundfix/SKILL.md`, raise both of its version fields by one patch step, regenerate `skills/roundfix/SKILL.md` with `make skills-sync`, and re-record `skills/testdata/owned-skill-versions.json` with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`; `go test -count=1 ./skills` MUST pass.

## Subtasks

- [ ] Characterize `upgrade` before changing it.
- [ ] Add the notice to `upgrade`, in process and through the installed executable.
- [ ] Add the Doctor line.
- [ ] Add the scope test and the other tests, and update the invalidated ones.
- [ ] Document both, then sync the mirror and the digests.
- [ ] Record one sabotage per new gate in the Result.

## Acceptance Criteria

- [ ] The characterization passes before and after the change.
- [ ] Each of the four release outcomes writes the notice to standard error.
- [ ] After an install, `installedProfilesCheck` receives the installed path, and its output is the notice. It is not called otherwise.
- [ ] When the comparison fails, `upgrade` exits as it would have and standard error carries one `not checked` line.
- [ ] Help, a usage error and a failed upgrade write no notice.
- [ ] Doctor prints `recommendations:` after `profiles:`, `ok`, `found` or `skipped`, and exits `0` when only that line reports a difference.
- [ ] A reference to `CheckRecommendations` from any other file fails the scope test.

## Context

- interface: `internal/cli/upgrade.go`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/health.go`
- interface: `internal/cli/cli.go`
- interface: `internal/cli/upgrade_test.go`
- interface: `internal/cli/doctor_test.go`
- creates: `internal/cli/upgrade_characterization_test.go`
- creates: `internal/cli/upgrade_notice_test.go`
- creates: `internal/cli/doctor_recommendations_test.go`
- creates: `internal/cli/recommendation_check_scope_test.go`
- instruction: `internal/cli/cli_test.go`
- interface: `docs/user-guide/commands.md`
- interface: `docs/user-guide/usage.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`

## Verification

- `out="$(go test -count=1 -v -run "^(TestUpgradeStdoutAndExitCodesAreCharacterized|TestUpgradeWritesTheNoticeOnEveryReleaseOutcome|TestUpgradeAsksTheInstalledExecutableAfterAnInstall|TestUpgradeKeepsItsExitCodeWhenTheCheckFails|TestUpgradeWritesNoNoticeOnHelpUsageErrorOrFailure|TestDoctorReportsRecommendationsAfterProfiles|TestDoctorRecommendationsNeverFailDoctor|TestRecommendationCheckIsReachedOnlyFromItsCommands|TestRunUpgradeFixtureMatrix|TestRunUpgradeCheckReportsAvailableWithoutInstalling|TestRunDoctorProfileReadinessProvesEffectiveCategoriesAndReportsCounts|TestRunCommandHelp)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestUpgradeStdoutAndExitCodesAreCharacterized TestUpgradeWritesTheNoticeOnEveryReleaseOutcome TestUpgradeAsksTheInstalledExecutableAfterAnInstall TestUpgradeKeepsItsExitCodeWhenTheCheckFails TestUpgradeWritesNoNoticeOnHelpUsageErrorOrFailure TestDoctorReportsRecommendationsAfterProfiles TestDoctorRecommendationsNeverFailDoctor TestRecommendationCheckIsReachedOnlyFromItsCommands TestRunUpgradeFixtureMatrix TestRunUpgradeCheckReportsAvailableWithoutInstalling TestRunDoctorProfileReadinessProvesEffectiveCategoriesAndReportsCounts TestRunCommandHelp; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && go test -count=1 -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills` — expected: exit 0; before this Task none of the eight new named tests exists, so the command fails.
- `for pair in "docs/user-guide/commands.md|recommendations: ok" "docs/user-guide/commands.md|recommendation notice" "docs/user-guide/usage.md|recommendation notice" ".agents/skills/roundfix/SKILL.md|recommendation notice" "skills/roundfix/SKILL.md|recommendation notice" ".agents/skills/roundfix/SKILL.md|recommendations:"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && diff -r .agents/skills/roundfix skills/roundfix >/dev/null && make skills-sync-check` — expected: exit 0; before this Task no document names the recommendation notice, so the command fails.

## References

- [_prd.md](_prd.md) — Goals 1 and 4; User Stories 1 and 4; Core Features 5, 6 and 7; Success Metrics 4, 5 and 6; Declared breaks; Recorded limits
- [_techspec.md](_techspec.md) — Interfaces; Doctor; Upgrade; Surface Transcripts; API Contract 5; API Contract 6; Testing Approach 4–6; Build Order 4; Risks & Considerations
- ADR-0181

## Result

Implemented this Task's slice for Daemon Verification; status, Subtasks,
Acceptance Criteria, and authored Verification remain Daemon-owned. The only
pre-existing working-tree change was the Daemon's `task_04.md` status update.
No commit, push, Pull Request, Task Graph edit, or other Task edit was made.

### Implementation and acceptance evidence

- Characterization was written first. Before any production change,
  `GOCACHE=/private/tmp/roundfix-task04-cache rtk proxy go test ./internal/cli -run '^TestUpgradeStdoutAndExitCodesAreCharacterized$' -count=1`
  passed (exit 0). It fixes stdout and exit codes for no release, current,
  newer/current `--check`, installation, download failure, unknown flag and
  help, with no release-outcome stderr assertions. Help is compared byte for
  byte with `commandUsage("upgrade")`, whose text this Task expressly extends.
  The release-outcome stdout literals and all exit codes remain unchanged.
  The characterization also passed after implementation and sabotage restoration.
- `runUpgradeCommand` now receives the dispatcher's environment. Its outcome
  carries the installed path; the path is empty when nothing was installed.
  The outcome is printed before the notice. `TestUpgradeWritesTheNoticeOnEveryReleaseOutcome`
  covers all four release outcomes and compares in-process differences with
  the shared `profiles check` renderer, including pinned/expired deviations.
- `installedProfilesCheck` is a dependency, defaulting to the installed path
  with exactly `profiles check`, using the process working directory and a
  ten-second child context. `TestUpgradeAsksTheInstalledExecutableAfterAnInstall`
  records the path and directory, checks the deadline and that stdout already
  holds the outcome, requires the child's exact output on stderr, and forbids
  child calls on the other outcomes. Release lookup/download and child execution
  use the existing fixtures and injected fakes; no new test starts a child or
  reaches the network.
- `TestUpgradeKeepsItsExitCodeWhenTheCheckFails` covers child failure, timeout,
  unavailable config home, unavailable working directory, and malformed YAML.
  Each keeps successful upgrade stdout/exit behavior and writes exactly one
  `roundfix: recommendations not checked: <reason>` line, discarding failed
  child output and normalizing multiline reasons. `TestUpgradeWritesNoNoticeOnHelpUsageErrorOrFailure`
  requires no notice or child call for help, an unknown flag, or failed download.
- Doctor computes the comparison from its already loaded Config, directly
  after `profiles`. `TestDoctorReportsRecommendationsAfterProfiles` covers
  exact position/details for current, differing, pinned and uncheckable configs;
  a difference and a pin exit 0 with otherwise healthy checks. Uncheckable
  required profiles retain the independent adapter-readiness failure while
  recommendations reports `skipped`. `TestDoctorRecommendationsNeverFailDoctor`
  requires `found` for a difference and `skipped` for an uncheckable comparison.
  The recommendation result invokes no runner and cannot report `failed`.
- `TestRecommendationCheckIsReachedOnlyFromItsCommands` parses every non-test
  Go file under `internal` and `cmd`, restricts references to all six TechSpec
  paths, and fails if there are zero references (a declaration alone does not
  count as a caller). The restored tree passed this gate.

### Focused checks

- After restoring all sabotages,
  `GOCACHE=/private/tmp/roundfix-task04-cache rtk proxy go test ./internal/cli -run 'TestUpgrade|TestRunUpgrade|TestDoctorRecommendations|TestDoctorReportsRecommendations|TestRecommendationCheck|TestRunDoctor|TestRunCommandHelp' -count=1 -v`
  passed (exit 0). All eight new named tests ran, along with the invalidated
  existing Doctor/upgrade tests and unchanged help contracts.
- `GOCACHE=/private/tmp/roundfix-task04-cache rtk proxy go test -count=1 ./skills`
  passed (exit 0).
- `rtk proxy git -c core.fsmonitor=false diff --check` passed. Source/diff
  inspection confirms only the Task's declared slice changed, existing
  top-level test names are preserved, and `internal/cli/cli_test.go` is untouched.
- The first `GOCACHE=/private/tmp/roundfix-task04-cache rtk make verify-incremental`
  reached the test suite but exited 2: two existing force-stop integration
  tests could not enumerate the process table (`operation not permitted`).
  No test or production code was altered to mask those environment failures.
  The same command rerun with process-table access passed (exit 0), including
  formatting, vet, the full Go suite (CLI freshly rerun), skill sync/readiness
  checks, and the binary build. Other unchanged packages reused their safe
  test cache, as intended by the incremental tier.

### Existing tests updated

Only the five invalidated Doctor tests gained the recommendation line or its
line-count/index shift:

- `TestRunDoctorProfileReadinessProvesEffectiveCategoriesAndReportsCounts`
- `TestRunDoctorAdapterReadinessReportsRequiredProfileRuntimes`
- `TestRunDoctorRepositorySkillReadiness`
- `TestRunDoctorMissingRepositoryRoot`
- `TestRunDoctorRealRepositoryCheckDoesNotMutateState`

Only the two invalidated upgrade tests changed their empty-stderr expectations:
`TestRunUpgradeFixtureMatrix` and
`TestRunUpgradeCheckReportsAvailableWithoutInstalling`. They now use isolated
configuration workspaces, and the shared existing fixture dependency supplies
an inert installed-check fake. No top-level test was renamed or removed.

### Sabotages and restoration

Each probe ran alone with
`GOCACHE=/private/tmp/roundfix-task04-cache rtk proxy go test ./internal/cli -run '^<test>$' -count=1`.
All three exited 1 with the named test reporting `FAIL`; each original file
was restored immediately afterward, then the focused suite above passed.

- Exit-code sabotage: temporarily returned `exitRunFailed` after a successful
  outcome and notice. `TestUpgradeKeepsItsExitCodeWhenTheCheckFails` failed on
  all five failure cases because exit 1 replaced exit 0. Restored `upgrade.go`.
- Scope sabotage: temporarily appended a function-value reference to
  `roundconfig.CheckRecommendations` in `internal/cli/cli.go`, outside the
  six allowed paths. `TestRecommendationCheckIsReachedOnlyFromItsCommands`
  failed with that file and source position. Restored `cli.go`; no other
  Task's file was touched.
- Doctor sabotage: temporarily assigned `CheckStatusFailed` to a difference.
  `TestDoctorRecommendationsNeverFailDoctor` failed with
  `difference status="failed" want found`. Restored `doctor.go`.

### Documentation and generated artifacts

Help names the new Doctor line and upgrade notice while retaining the strings
required by `TestRunCommandHelp`. `commands.md` includes the Doctor sample line,
`ok`/`found` semantics and `skipped` fallback, stderr placement, preserved
upgrade stdout/exit codes, and the older-executable limitation. `usage.md`
adds the requested recommendation-notice sentence.

The Roundfix skill text is under `### Recommendation check`; both version
fields rose from `0.0.11` to `0.0.12`. The `### QA settlement` section and all
later bytes were compared against HEAD and are unchanged.

- `rtk make skills-sync` passed. Its only byte-changed output was
  `skills/roundfix/SKILL.md`; all other owned mirrors remained byte-identical.
  The canonical source change is `.agents/skills/roundfix/SKILL.md`.
- `GOCACHE=/private/tmp/roundfix-task04-cache rtk proxy go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`
  passed and rewrote only `skills/testdata/owned-skill-versions.json`, recording
  Roundfix `0.0.12` with its content digest.
- `GOCACHE=/private/tmp/roundfix-task04-cache rtk make baseline-digests`
  passed with `changed:false`; it changed no derived file. Canonical and
  mirrored Roundfix skill bytes match.

The two authored `## Verification` commands were not run. Their execution and
Task settlement remain with the Daemon. No follow-up implementation outside
this slice was added.
