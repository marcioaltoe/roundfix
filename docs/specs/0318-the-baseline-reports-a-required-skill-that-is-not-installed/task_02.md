---
task: task_02
spec: 0318-the-baseline-reports-a-required-skill-that-is-not-installed
status: completed
type: backend
complexity: medium
---

# Task 02: skills check names a missing required skill of the repository it runs in

## Overview

`roundfix skills check` validates only the binary's embedded bundle, so it
passes in a repository whose Baseline dispatches a skill that is not
installed. This Task makes it inspect the adopted repository that contains its
work directory with task_01's helper, add one diagnostic per Missing Required
Skill to its existing failure list and exit `1`, keep its bundle-only result
everywhere else, make the two existing tests that run it independent of this
repository, and document it.

## Requirements

1. MUST answer the Backlog Entry
   [`references/2026-10-25-the-baseline-update-preview-reports-current-while-a-required-skill-is-missing.md`](references/2026-10-25-the-baseline-update-preview-reports-current-while-a-required-skill-is-missing.md)
   of 2026-10-25 (`_prd.md` → Sources 1), the `skills check` half of the
   maintainer's decision, by implementing API Contract 5 of `_techspec.md`.
2. MUST add `internal/cli/skills_check_repository.go` with the repository
   section of API Contract 5, using `environment.resolveWorkDir`,
   `environment.dependencies.resolveSkillsProjectRoot`,
   `environment.dependencies.doctor.resolveExternal`,
   `environment.dependencies.doctor.checkSkills` and task_01's
   `missingRequiredSkills`; it adds nothing when any of them fails, when the
   Setup Manifest is absent or unreadable, or when it names no external skill.
3. MUST, in `runSkillsCommand` (`internal/cli/cli.go`, `case "check"`), append
   those diagnostics after the bundle diagnostics and print them through the
   existing failure block, exit `1`; and change the `skills` usage entry for
   `check` to the text of `_techspec.md` → Interfaces. No flag is added, and
   the pass and `unversioned` lines keep their bytes.
4. MUST add `internal/cli/skills_check_repository_test.go` with the
   `Suite:`/`Invariant:`/`Boundary` header and
   `TestSkillsCheckNamesAMissingRequiredSkill` exactly as `_techspec.md` →
   Testing Approach describes, with those subtest names, calling
   `t.Parallel()` first in the test and each subtest; the missing-skill
   subtest asserts the exact stderr, an empty stdout and exit `1`.
5. MUST add `setCommandWorkDirForTest(t, t.TempDir())` as the first setup
   line of `TestRunSkillsCheck` in `internal/cli/cli_test.go` and of each
   subtest of `TestDoctorAndSkillsCheckReportSharedOwnedSkillReadiness` in
   `internal/cli/doctor_test.go`, so neither reads this repository; nothing
   else in those files changes.
6. MUST append the paragraph of `_techspec.md` → Interfaces to
   `docs/user-guide/commands/skills.md`, and the `setup` reference paragraph
   to the canonical `.agents/skills/roundfix/references/setup.md` after the
   paragraph that ends "nothing is restored.". Never touch a `### QA
   settlement` section.
7. MUST run `make skills-sync`, then
   `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
   then `make skills-sync` again.
8. MUST, as the Task's last step, after every guide, skill and usage edit,
   re-record the Behavior Surface fingerprint ("re-record after the last
   guide edit") with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
9. MUST NOT change the `baseline update` output, Doctor's output, any Baseline
   module or derived Baseline file, `docs/agents/`, the Makefile or the CI
   workflows.

## Subtasks

- [ ] Add the repository section and wire it into skills check with the new usage text.
- [ ] Add the skills check test and make the two existing tests use a temporary work directory.
- [ ] Add the guide and reference paragraphs, sync and record the Roundfix skill, then re-record the Behavior Surface fingerprint last.

## Acceptance Criteria

- [ ] In an adopted repository with a missing required skill, skills check exits `1` with the skill's diagnostic on stderr.
- [ ] In a complete adopted repository, a repository without a Setup Manifest and outside a Git repository, skills check keeps its pass line and exit `0`.
- [ ] The existing skills check tests no longer read this repository and still pass.

## Consumers

None. The `skills check` output and exit code change only inside an adopted repository with a Missing Required Skill; this repository's `make verify` runs it at the root, where the Repository Skill Set is complete, so it keeps passing (verified in the prototype), and no Go package parses its output.

## Context

- instruction: `internal/cli/doctor.go`
- instruction: `internal/cli/baseline_plan_test.go`
- instruction: `Makefile`
- creates: `internal/cli/skills_check_repository.go`
- interface: `internal/cli/cli.go`
- creates: `internal/cli/skills_check_repository_test.go`
- interface: `internal/cli/cli_test.go`
- interface: `internal/cli/doctor_test.go`
- interface: `docs/user-guide/commands/skills.md`
- interface: `.agents/skills/roundfix/references/setup.md`
- interface: `skills/roundfix/references/setup.md`
- interface: `.agents/skills/roundfix/SKILL.md`
- interface: `skills/roundfix/SKILL.md`
- interface: `skills/testdata/owned-skill-versions.json`
- interface: `docs/references/behavior-surfaces.json`

## Verification

- `out="$(go test -timeout 15m -count=1 -v -run '^(TestSkillsCheckNamesAMissingRequiredSkill|TestRunSkillsCheck|TestDoctorAndSkillsCheckReportSharedOwnedSkillReadiness)$' ./internal/cli 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSkillsCheckNamesAMissingRequiredSkill/adopted_repository_with_a_missing_required_skill TestSkillsCheckNamesAMissingRequiredSkill/adopted_repository_with_every_required_skill TestSkillsCheckNamesAMissingRequiredSkill/repository_without_a_Setup_Manifest TestSkillsCheckNamesAMissingRequiredSkill/outside_a_Git_repository TestRunSkillsCheck TestDoctorAndSkillsCheckReportSharedOwnedSkillReadiness; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done; go test -timeout 10m -count=1 ./skills && go test -timeout 15m -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract && tr -s '[:space:]' ' ' < docs/user-guide/commands/skills.md | grep -qF -- 'Outside a Git repository, or without a Setup Manifest, it checks only the shipped bundle.' && tr -s '[:space:]' ' ' < .agents/skills/roundfix/references/setup.md | grep -qF -- 'validates the shipped bundle and, inside an adopted Git repository, fails with exit'` — expected: exit 0. Before this Task `TestSkillsCheckNamesAMissingRequiredSkill` does not exist, so the command fails at its first pass check, and neither paragraph exists. Afterwards every subtest passes, the two existing tests pass from a temporary work directory, the skill version, mirror and Behavior Surface record are current, and both paragraphs exist.

## References

- `_prd.md` → Sources; Goals; Core Features 2 and 4; Success Metrics 3 and 6
- `_techspec.md` → Current behavior; Interfaces; API Contract 5; Invariants 2 and 3; Testing Approach; Build Order 2

## Result

Implemented the repository section of `skills check` using the existing work-directory,
repository-root and Doctor dependency seams and task_01's missing-skill helper.
Named missing or lock-mismatched external skills append diagnostics after bundle
findings. Resolution failures, unreadable or absent manifests, and empty external
requirements preserve the bundle-only result. Pass and unversioned output bytes
are unchanged. Added the exact authored usage and documentation paragraphs.

Acceptance evidence from focused implementation checks:

- Missing required skill: the new regression first reproduced exit `0` and the
  bundle pass line. After implementation, its missing-skill subtest passed with
  exact stderr, empty stdout and exit `1`.
- Complete adopted repository, repository without a Setup Manifest, and outside
  Git: all three subtests passed, asserting the exact bundle pass line, empty
  stderr and exit `0`. Every new test and subtest calls `t.Parallel()` first.
- Existing checks: `TestRunSkillsCheck` and the three shared-readiness subtests
  now set a fresh temporary command work directory; no other edits were made
  in those two test files. Their focused run passed.

Focused command:
`go test ./internal/cli -count=1 -run '^(TestSkillsCheckNamesAMissingRequiredSkill|TestRunSkillsCheck|TestDoctorAndSkillsCheckReportSharedOwnedSkillReadiness)$'`
returned exit `0` (10 tests passed).

Required regeneration sequence ran successfully: `make skills-sync`,
`go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions`,
then `make skills-sync`. The final implementation step is the required
`go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`
recorder, after every usage, guide and skill edit.

`git diff --check` returned exit `0`. The initial worktree contained only the
Daemon-owned task status change; it was preserved. Commands used the resolved
`/Volumes/Brain-SSD/cache/.roundfix/worktrees/` checkout because the sandbox
rejected the configured symlinked writable root.

Declared Verification and repository Verification remain Daemon-owned and were
not run. No Task status, Task Graph, other Task, Doctor implementation, Baseline
module, derived Baseline guide, Makefile, CI workflow or QA settlement section
was changed. No commit, push or PR was made. No follow-up work identified.
