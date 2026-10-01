---
task: task_03
spec: 0196-a-notice-when-profiles-fall-behind
status: pending
type: backend
complexity: high
---

# Task 03: `profiles check --apply` adopts the recommendation

## Overview

After task_02 a maintainer can see which categories differ, and still has to write a profile fragment by hand to adopt the recommendation. This Task adds `--apply` to `roundfix profiles check`: it builds the fragment from the differing categories and hands it to the write flow `roundfix profiles configure` already owns, with the same proof, preview, confirmation and exit codes (ADR-0181). A pinned category is skipped. It is verifiable on its own: after `--apply --scope project --yes`, a following `profiles check` counts no difference.

## Requirements

1. MUST extract the write flow of `runProfilesConfigureCommand`, from `PrepareProfilesConfig` to its final output, into one function in `internal/cli/profiles_configure.go` that both commands call. `roundfix profiles configure` MUST keep its outputs, prompts and exit codes byte-identical; its existing tests are the characterization and MUST pass unchanged. The Result MUST state that they were run before and after the extraction.
2. MUST add `--apply`, `--scope`, `--dry-run` and `--yes` to `roundfix profiles check` in the new file `internal/cli/profiles_check_apply.go`, routed from `runProfilesCommand` in `internal/cli/profiles.go` when `--apply` is among the arguments, with the rules and outputs of `_techspec.md` → Adoption and API Contract 2. `--apply` without `--scope`, and `--scope`, `--dry-run` or `--yes` without `--apply`, MUST be a usage error with exit `2` that writes nothing. Without `--apply`, the command of task_02 MUST behave exactly as before.
3. MUST select only the categories whose status is `differs`. A `pinned` category and a `current` category MUST NOT be written, and a pinned category's bytes in the configuration file MUST stay identical.
4. MUST, with `--scope user`, leave out every category whose effective profile comes from the Project Config, and name each on standard error with the advice to use `--scope project`.
5. MUST write the Recommended Profile of each selected category without a deviation, only after the exact proof of every tuple it writes, and only after confirmation unless `--yes` is given. `--dry-run`, a failed proof and a declined confirmation MUST leave every configuration byte unchanged.
6. MUST, when no category is selected, prepare, prove and write nothing, print the nothing-to-adopt outcome of API Contract 2, and exit `0`.
7. MUST add the flags to the help text of `profiles` and `profiles check` in `internal/cli/cli.go`, keeping every string `TestRunCommandHelp` requires.
8. MUST put the new tests in `internal/cli/profiles_check_apply_test.go`, with the existing fake runner, a temporary home and a temporary repository. The proof-order test MUST prove the runner was called before the file changed. Each refusal MUST compare the configuration bytes before and after.
9. MUST prove each new gate can fail. The Result MUST record one sabotage for the pinned rule (for example selecting pinned categories) and one for the proof order (for example writing before the proof), each with the test that failed, and that the code was restored.
10. MUST document `--apply` in the `profiles` section of `docs/user-guide/commands.md`, with its flags, the user-scope rule and the statement that it opens disposable Agent Sessions for proof, and add one sentence on it to `docs/user-guide/usage.md`.
11. MUST add `--apply` to the text under the skill heading task_01 created. It MUST put the skill text under the heading `### Recommendation check` of `.agents/skills/roundfix/SKILL.md`, add no text inside the `### QA settlement` section, then run `make skills-sync` and `make baseline-digests`, and name in the Result every file either command rewrote.
12. MUST, because this Task changes `.agents/skills/roundfix/SKILL.md`, raise both of its version fields by one patch step, regenerate `skills/roundfix/SKILL.md` with `make skills-sync`, and re-record `skills/testdata/owned-skill-versions.json` with `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`; `go test -count=1 ./skills` MUST pass.

## Subtasks

- [ ] Extract the shared write flow and run the existing configure tests before and after.
- [ ] Add `--apply` and its flag rules.
- [ ] Add the tests, each refusal with a byte comparison.
- [ ] Document `--apply`, then sync the mirror and the digests.
- [ ] Record one sabotage per new gate in the Result.

## Acceptance Criteria

- [ ] `roundfix profiles check --apply --scope project --yes` writes the Recommended Profile of each differing category, and a following `roundfix profiles check` counts no difference.
- [ ] A pinned category is not written, and its bytes are unchanged.
- [ ] Every tuple written was proved first, and a failed proof writes nothing.
- [ ] `--dry-run` and a declined confirmation write nothing.
- [ ] `--scope user` names a category the Project Config defines and leaves it.
- [ ] With nothing to adopt, the command writes nothing and exits `0`.
- [ ] The existing `profiles configure` tests pass without a change.

## Context

- creates: `internal/cli/profiles_check_apply.go`
- interface: `internal/cli/profiles.go`
- interface: `internal/cli/profiles_configure.go`
- interface: `internal/cli/cli.go`
- creates: `internal/cli/profiles_check_apply_test.go`
- instruction: `internal/cli/profiles_configure_test.go`
- instruction: `internal/cli/cli_test.go`
- instruction: `internal/config/profile_config.go`
- interface: `docs/user-guide/commands.md`
- interface: `docs/user-guide/usage.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`

## Verification

- `out="$(go test -count=1 -v -run "^(TestProfilesCheckApplyWritesTheRecommendedProfileOfEachDifferingCategory|TestProfilesCheckApplyLeavesAPinnedCategoryUntouched|TestProfilesCheckApplyProvesBeforeItWrites|TestProfilesCheckApplyDryRunWritesNothing|TestProfilesCheckApplyDeclinedConfirmationWritesNothing|TestProfilesCheckApplyUserScopeNamesProjectDefinedCategories|TestProfilesCheckApplyWithNothingToAdoptChangesNothing|TestProfilesCheckApplyFlagRules|TestProfilesConfigureExitCodes|TestProfilesConfigureChangeSummary|TestProfilesConfigureProofScope|TestProfilesConfigureProofRunsBeforeConfirmationAndWrite|TestProfilesConfigureDryRunAndFailedConfigurationLeaveBytesUnchanged|TestRunCommandHelp)$" ./internal/cli 2>&1)" || { printf "%s\\n" "$out"; exit 1; }; for name in TestProfilesCheckApplyWritesTheRecommendedProfileOfEachDifferingCategory TestProfilesCheckApplyLeavesAPinnedCategoryUntouched TestProfilesCheckApplyProvesBeforeItWrites TestProfilesCheckApplyDryRunWritesNothing TestProfilesCheckApplyDeclinedConfirmationWritesNothing TestProfilesCheckApplyUserScopeNamesProjectDefinedCategories TestProfilesCheckApplyWithNothingToAdoptChangesNothing TestProfilesCheckApplyFlagRules TestProfilesConfigureExitCodes TestProfilesConfigureChangeSummary TestProfilesConfigureProofScope TestProfilesConfigureProofRunsBeforeConfirmationAndWrite TestProfilesConfigureDryRunAndFailedConfigurationLeaveBytesUnchanged TestRunCommandHelp; do printf "%s\\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done && go test -count=1 -run '^TestEveryOwnedSkillVersionIsRecorded$' ./skills` — expected: exit 0; before this Task none of the eight new named tests exists, so the command fails.
- `for pair in "docs/user-guide/commands.md|roundfix profiles check --apply" "docs/user-guide/usage.md|roundfix profiles check --apply" ".agents/skills/roundfix/SKILL.md|roundfix profiles check --apply" "skills/roundfix/SKILL.md|roundfix profiles check --apply"; do file="${pair%%|*}"; phrase="${pair#*|}"; tr -s '[:space:]' ' ' < "$file" | grep -qF -- "$phrase" || { printf 'missing phrase in %s: %s\n' "$file" "$phrase" >&2; exit 1; }; done && diff -r .agents/skills/roundfix skills/roundfix >/dev/null && make skills-sync-check` — expected: exit 0; before this Task no document names `roundfix profiles check --apply`, so the command fails.

## References

- [_prd.md](_prd.md) — Goal 2; User Story 2; Core Feature 4; Success Metric 3; Recorded limits
- [_techspec.md](_techspec.md) — Adoption; API Contract 2; Testing Approach 3; Build Order 3; Risks & Considerations
- ADR-0181; ADR-0049; ADR-0140

## Result
