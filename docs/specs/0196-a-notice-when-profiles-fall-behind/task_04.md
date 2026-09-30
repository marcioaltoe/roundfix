---
task: task_04
spec: 0196-a-notice-when-profiles-fall-behind
status: pending
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
- [ ] Doctor prints `recommendations:` after `profiles:`, `ok` or `found`, and exits `0` when only that line reports a difference.
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

## Verification

- `out="$(go test -count=1 -v -run "^(TestUpgradeStdoutAndExitCodesAreCharacterized|TestUpgradeWritesTheNoticeOnEveryReleaseOutcome|TestUpgradeAsksTheInstalledExecutableAfterAnInstall|TestUpgradeKeepsItsExitCodeWhenTheCheckFails|TestUpgradeWritesNoNoticeOnHelpUsageErrorOrFailure|TestDoctorReportsRecommendationsAfterProfiles|TestDoctorRecommendationsNeverFailDoctor|TestRecommendationCheckIsReachedOnlyFromItsCommands|TestRunUpgradeFixtureMatrix|TestRunUpgradeCheckReportsAvailableWithoutInstalling|TestRunDoctorProfileReadinessProvesEffectiveCategoriesAndReportsCounts|TestRunCommandHelp)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestUpgradeStdoutAndExitCodesAreCharacterized TestUpgradeWritesTheNoticeOnEveryReleaseOutcome TestUpgradeAsksTheInstalledExecutableAfterAnInstall TestUpgradeKeepsItsExitCodeWhenTheCheckFails TestUpgradeWritesNoNoticeOnHelpUsageErrorOrFailure TestDoctorReportsRecommendationsAfterProfiles TestDoctorRecommendationsNeverFailDoctor TestRecommendationCheckIsReachedOnlyFromItsCommands TestRunUpgradeFixtureMatrix TestRunUpgradeCheckReportsAvailableWithoutInstalling TestRunDoctorProfileReadinessProvesEffectiveCategoriesAndReportsCounts TestRunCommandHelp; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task none of the eight new named tests exists, so the command fails.
- `for pair in "docs/user-guide/commands.md|recommendations: ok" "docs/user-guide/commands.md|recommendation notice" "docs/user-guide/usage.md|recommendation notice" ".agents/skills/roundfix/SKILL.md|recommendation notice" "skills/roundfix/SKILL.md|recommendation notice" ".agents/skills/roundfix/SKILL.md|recommendations:"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && diff -r .agents/skills/roundfix skills/roundfix >/dev/null && make skills-sync-check` — expected: exit 0; before this Task no document names the recommendation notice, so the command fails.

## References

- [_prd.md](_prd.md) — Goals 1 and 4; User Stories 1 and 4; Core Features 5, 6 and 7; Success Metrics 4, 5 and 6; Declared breaks; Recorded limits
- [_techspec.md](_techspec.md) — Interfaces; Doctor; Upgrade; Surface Transcripts; API Contract 5; API Contract 6; Testing Approach 4–6; Build Order 4; Risks & Considerations
- ADR-0181

## Result
