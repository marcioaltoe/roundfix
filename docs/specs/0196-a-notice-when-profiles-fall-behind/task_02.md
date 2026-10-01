---
task: task_02
spec: 0196-a-notice-when-profiles-fall-behind
status: pending
type: backend
complexity: high
---

# Task 02: `profiles check` compares the configuration with the recommendation

## Overview

Nothing compares a configured profile with the Recommended Profile. This Task adds the comparison as one pure function in `internal/config`, the read-only command `roundfix profiles check` that prints it as text or JSON, and the status of each category in `roundfix profiles show` (ADR-0181). It is verifiable on its own: with a configuration that differs, the command names the category, both profiles and the counts, and exits `0`.

## Requirements

1. MUST add `CheckRecommendations` and its types in the new file `internal/config/recommendation_check.go`, as `_techspec.md` → Interfaces and → The comparison state. It MUST read only its argument, `RecommendedProfile` and `ModelRecommendationSnapshotVersion`.
2. MUST report a category as `pinned` only when it differs and its deviation's `from` equals the shipped snapshot. A deviation with another date MUST leave the category `differs` and stay on the row. A category that equals the recommendation MUST be `current` whether or not it has a deviation.
3. MUST add `roundfix profiles check [--json]` in the new file `internal/cli/profiles_check.go`, dispatched from `runProfilesCommand`, with the text of `_techspec.md` → Surface Transcripts and the JSON and exit codes of API Contract 1. The text renderer MUST be one function that takes a `RecommendationCheck`, because task_04 reuses it.
4. MUST keep the command read-only and offline: it loads the configuration as `profiles show` does, calls no runner, opens no Agent Session, reaches no network and writes no file. It MUST NOT accept `--apply`, `--scope`, `--dry-run` or `--yes` in this Task; task_03 adds them.
5. MUST make `roundfix profiles show` print the lines and JSON fields of `_techspec.md` → Show, with status `inherited` for an optional category the configuration does not define. Its schema, flags and exit codes MUST not change.
6. MUST add `profiles check` to the usage text of `profiles` and add its own help text in `internal/cli/cli.go`. It MUST keep every string `TestRunCommandHelp` requires, so `internal/cli/cli_test.go` does not change.
7. MUST put the new tests in `internal/config/recommendation_check_test.go`, `internal/cli/profiles_check_test.go` and `internal/cli/profiles_show_status_test.go`. The command tests MUST use a temporary home and repository. One MUST prove, with a runner that fails the test when called and a byte comparison of the configuration files before and after, that the command opens no session and writes nothing.
8. MUST prove each new gate can fail. The Result MUST record one sabotage for the status rule (for example treating any deviation as a pin) and one for the read-only test (for example calling the runner), each with the test that failed, and that the code was restored.
9. MUST document the command in the `profiles` section of `docs/user-guide/commands.md`, with its statuses, its exit codes and the schema `roundfix/profiles-check/v1`, and the `Recommendation status` line of `profiles show`. It MUST add one paragraph on the command to `docs/user-guide/usage.md`.
10. MUST describe the command and its three statuses under the skill heading task_01 created. It MUST put the skill text under the heading `### Recommendation check` of `.agents/skills/roundfix/SKILL.md`, add no text inside the `### QA settlement` section, then run `make skills-sync` and `make baseline-digests`, and name in the Result every file either command rewrote.
11. MUST, because this Task changes `.agents/skills/roundfix/SKILL.md`, raise both of its version fields by one patch step, regenerate `skills/roundfix/SKILL.md` with `make skills-sync`, and re-record `skills/testdata/owned-skill-versions.json` with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`; `go test -count=1 ./skills` MUST pass.

## Subtasks

- [ ] Add the comparison and its tests.
- [ ] Add `profiles check`, its renderer, its JSON and its help text.
- [ ] Add the status and deviation lines to `profiles show`.
- [ ] Document the command, then sync the mirror and the digests.
- [ ] Record one sabotage per new gate in the Result.

## Acceptance Criteria

- [ ] A configuration on built-in profiles reports every required category `current`, and the text is the single summary line.
- [ ] A differing category prints its configured profile with its source and the recommended profile, and the summary counts it.
- [ ] A deviation naming the shipped snapshot makes the category `pinned`; one naming another date leaves it `differs` and says so.
- [ ] `roundfix profiles check --json` reports schema `roundfix/profiles-check/v1` and every configured category.
- [ ] The command exits `0` with differences and `2` for an unknown flag, an extra argument or a configuration that does not load.
- [ ] The command calls no runner and changes no file.
- [ ] `roundfix profiles show` states the recommendation status of every category.

## Context

- creates: `internal/config/recommendation_check.go`
- creates: `internal/config/recommendation_check_test.go`
- creates: `internal/cli/profiles_check.go`
- creates: `internal/cli/profiles_check_test.go`
- creates: `internal/cli/profiles_show_status_test.go`
- interface: `internal/cli/profiles.go`
- interface: `internal/cli/cli.go`
- instruction: `internal/config/recommendations.go`
- instruction: `internal/config/profiles.go`
- instruction: `internal/cli/cli_test.go`
- interface: `docs/user-guide/commands.md`
- interface: `docs/user-guide/usage.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`

## Verification

- `out="$(go test -count=1 -v -run "^(TestCheckRecommendationsReportsCurrentDiffersAndPinned|TestCheckRecommendationsEndsADeviationOfAnotherSnapshot|TestCheckRecommendationsReportsBuiltinsAsCurrent|TestCheckRecommendationsLeavesOutUndefinedOptionalCategories|TestProfilesCheckPrintsEachDifferenceAndTheSummary|TestProfilesCheckExitsZeroWithDifferences|TestProfilesCheckJSONIsSchemaV1|TestProfilesCheckRefusesUnknownFlagsAndArguments|TestProfilesCheckOpensNoAgentSessionAndWritesNothing|TestProfilesShowStatesTheRecommendationStatus|TestRunCommandHelp)$" ./internal/config ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestCheckRecommendationsReportsCurrentDiffersAndPinned TestCheckRecommendationsEndsADeviationOfAnotherSnapshot TestCheckRecommendationsReportsBuiltinsAsCurrent TestCheckRecommendationsLeavesOutUndefinedOptionalCategories TestProfilesCheckPrintsEachDifferenceAndTheSummary TestProfilesCheckExitsZeroWithDifferences TestProfilesCheckJSONIsSchemaV1 TestProfilesCheckRefusesUnknownFlagsAndArguments TestProfilesCheckOpensNoAgentSessionAndWritesNothing TestProfilesShowStatesTheRecommendationStatus TestRunCommandHelp; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && go test -count=1 -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills` — expected: exit 0; before this Task none of the ten new named tests exists, so the command fails.
- `for pair in "docs/user-guide/commands.md|roundfix profiles check" "docs/user-guide/commands.md|roundfix/profiles-check/v1" "docs/user-guide/commands.md|Recommendation status" "docs/user-guide/usage.md|roundfix profiles check" ".agents/skills/roundfix/SKILL.md|roundfix profiles check" "skills/roundfix/SKILL.md|roundfix profiles check" ".agents/skills/roundfix/SKILL.md|roundfix/profiles-check/v1"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && diff -r .agents/skills/roundfix skills/roundfix >/dev/null && make skills-sync-check` — expected: exit 0; before this Task no document names `roundfix profiles check`, so the command fails.

## References

- [_prd.md](_prd.md) — Goals 2 and 3; User Stories 2 and 3; Core Features 2 and 3; Success Metrics 1 and 2; Declared breaks
- [_techspec.md](_techspec.md) — Interfaces; The comparison; The renderer; Show; Surface Transcripts; API Contract 1; API Contract 4; Testing Approach 2; Build Order 2
- ADR-0181; ADR-0180; ADR-0107

## Result
