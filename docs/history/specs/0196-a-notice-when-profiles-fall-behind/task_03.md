---
task: task_03
spec: 0196-a-notice-when-profiles-fall-behind
status: completed
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

Implemented the adoption slice for Daemon Verification; status remains
Daemon-owned. No commit, push, Pull Request, other Task edit or Task Graph edit
was performed. At entry, only this Task file was modified (`in_progress`).
The new adoption file and its eight named tests were absent.

### Implementation and acceptance evidence

- Extracted `writeProfilesConfiguration` mechanically from
  `runProfilesConfigureCommand`, from preparation through final output and
  rollback. Both configure and adoption call it. Existing configure tests were
  unchanged. `rtk proxy go test -count=1 ./internal/cli -run
  '^TestProfilesConfigure'` exited 0 before extraction (0.834s) and after
  extraction (2.423s). An intermediate compile attempt while the dispatcher
  referenced the not-yet-created adoption function failed with an undefined
  function; creating that function resolved it.
- `TestProfilesCheckApplyWritesTheRecommendedProfileOfEachDifferingCategory`
  checks project adoption of backend, docs and qa; all complete written
  profiles equal the recommendation and carry no deviation. The fake runner's
  exact proof tuples equal all selected preferred/fallback tuples. User Config
  stays byte-identical, and the following read-only check reports `0 differ,
  1 pinned`.
- `TestProfilesCheckApplyLeavesAPinnedCategoryUntouched` compares the pinned
  review block's bytes and rejects selected pinned/current categories. The
  original assertion included a newly appended docs block after review; it
  was corrected to compare the pinned block itself, after inspecting output
  showing its bytes were preserved.
- `TestProfilesCheckApplyProvesBeforeItWrites` reads every fixture file from
  inside each proof callback and the confirmation callback, requiring the
  original bytes. It requires actual runner calls and a later successful
  write. Its failed-proof case requires exit 2, a proof diagnostic and every
  original file byte unchanged.
- `TestProfilesCheckApplyDryRunWritesNothing` requires proof calls, preview,
  dry-run output, no confirmation and unchanged bytes.
  `TestProfilesCheckApplyDeclinedConfirmationWritesNothing` exercises the real
  confirmation reader with `n`, requires proof first, exit 1, JSON refusal,
  a decline diagnostic and unchanged bytes.
- `TestProfilesCheckApplyUserScopeNamesProjectDefinedCategories` requires
  backend and qa on stderr with `--scope project` advice, selects only docs,
  verifies its Recommended Profile, and compares Project Config bytes.
- `TestProfilesCheckApplyWithNothingToAdoptChangesNothing` covers built-ins,
  a pinned-only configuration and a project-defined difference excluded by
  user scope. The forbidden runner and confirmation callback reject any proof
  or prompt. Text is the exact nothing-to-adopt outcome; JSON has schema
  `roundfix/profiles-configure/v1`, `changed: false` and empty arrays. Every
  file stays byte-identical and the command exits 0.
- `TestProfilesCheckApplyFlagRules` covers missing/invalid scope, write flags
  without apply, false apply, unknown flags, invalid boolean input and extra
  arguments. Each case requires exit 2, no stdout, a diagnostic, no runner
  call and unchanged bytes. The read-only implementation is unchanged.
- Added flags to top-level, profiles and profiles-check help, retaining the
  existing read-only usage line and help strings. Documented adoption flags,
  scope precedence, proof Sessions and no-write outcomes in the commands guide;
  added the requested usage sentence.

### Focused checks

- `rtk proxy go test -count=1 ./internal/cli -run
  '^TestProfiles(CheckApply|Configure|Check)'`: exit 0, after correcting the new
  test's runner field name and pinned-block boundary.
- After both sabotages were restored, `rtk proxy go test -count=1
  ./internal/cli -run
  '^TestProfiles(CheckApply|Configure|Check)|^TestRunCommandHelp$'`: exit 0
  (0.729s). This includes all new adoption tests, unchanged configure tests,
  read-only check tests and help characterization.
- `rtk proxy go test ./skills -run
  '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`: exit 0.
- `rtk proxy go test -count=1 ./skills`: exit 0 (1.347s).
- `rtk git diff --check`: exit 0.
- Initial `rtk make verify-incremental`: exit 2. Process-owner integration
  tests could not read the host process table in the sandbox. Baseline and CLI
  suite guards also detected my concurrent Result edit. The rerun uses the
  required process permission and keeps all repository files unchanged while
  the check runs.
- `rtk make verify-incremental` rerun with required process permission:
  exit 0. Formatting, vet, repository tests (CLI 159.809s), skill sync/check
  and build passed. No repository edits occurred while the rerun was active.
- Python inspection confirmed the skill mirror equals its canonical source,
  `### QA settlement` is byte-identical to HEAD, and status is `in_progress`.
- Authored `## Verification` commands were not run; they remain Daemon-owned.

### Sabotage evidence

1. Changed the selector to skip only current categories, deliberately selecting
   pinned categories. `rtk proxy go test -count=1 ./internal/cli -run
   '^TestProfilesCheckApplyLeavesAPinnedCategoryUntouched$'` exited 1:
   `selected pinned/current category` identified review. Restored the original
   selector from a task-scoped temporary copy.
2. Inserted `PersistProfilesConfig` immediately before proof in the shared
   flow. `rtk proxy go test -count=1 ./internal/cli -run
   '^TestProfilesCheckApplyProvesBeforeItWrites$'` exited 1 in both `success`
   and `failed proof`: the proof callback found changed file bytes. Restored
   the shared flow from its temporary copy. The subsequent focused check
   above passed with both sabotages removed.

### Skill regeneration and scope

Adoption text lives only under `### Recommendation check`; the skill's
`### QA settlement` section is untouched. Both version fields rose from
0.0.10 to 0.0.11. The protected source write used the Spec's explicit
`_authorization.md` grant through sandbox escalation.

- `rtk proxy make skills-sync`: exit 0; rewrote
  `skills/roundfix/SKILL.md` from `.agents/skills/roundfix/SKILL.md`.
- `rtk proxy make baseline-digests`: exit 0; reported no changes and rewrote
  no files (derived artifacts already matched their canonical sources).
- The version-recording command rewrote
  `skills/testdata/owned-skill-versions.json`, adding version 0.0.11 and its
  canonical digest.

No follow-up implementation was added to this slice.
