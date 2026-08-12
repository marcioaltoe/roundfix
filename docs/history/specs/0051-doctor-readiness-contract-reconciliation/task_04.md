---
task: task_04
spec: 0051-doctor-readiness-contract-reconciliation
status: completed
type: backend
complexity: medium
---

# Task 04: Reconcile Doctor readiness contracts

## Overview

Restore the Agent Selection Profile working-directory fallback without
weakening the Git-root-only Repository Skill Set boundary, then make Doctor's
terminology and mixed remediation exact. Public command tests and user
guidance make the corrected output contract independently reviewable.

## Requirements

1. MUST pass the trimmed Git root to profile proof when present and the process
   working directory when it is absent.
2. MUST pass only a non-empty trimmed Git root to Repository Skill Set
   inspection; an empty root must not invoke the checker.
3. MUST report `Repository Skill Set readiness requires a Git repository` for
   the missing-root detail and keep the existing run-from-Git next action.
4. MUST join multiple ownership remediation actions with ` && ` in
   owned-then-external order while retaining the existing `"; next: "` result
   boundary.
5. MUST preserve eager check execution, line order, stdout/stderr discipline,
   and exit codes.
6. MUST move Doctor-specific behavior tests and helpers from the CLI-wide test
   file to `doctor_test.go`, while leaving dispatch and help registry tests in
   the CLI-wide suite.
7. MUST update user guidance for the canonical term and fail-closed command
   chain without editing either protected Roundfix Skill file in this Task.

## Subtasks

- [x] Add the outside-Git profile-workdir regression and checker non-invocation
      assertion.
- [x] Restore separate profile and repository root resolution.
- [x] Correct missing-root detail and mixed remediation composition.
- [x] Move Doctor-specific tests to the source-matched test file without
      weakening assertions.
- [x] Update the public Doctor example and explanatory user guidance.
- [x] Exercise inside-Git, outside-Git, typed ownership, and unclassified error
      paths through the public runner.

## Acceptance Criteria

- [x] Outside Git, profile proof receives `os.Getwd()` and the skills checker is
      not invoked.
- [x] Inside Git, profile proof and Repository Skill Set inspection both receive
      the resolved Git root.
- [x] Missing-root output uses the canonical Repository Skill Set term.
- [x] Mixed remediation prints one exact `&&`-joined chain and a symlinked lock
      prints only external remediation.
- [x] Codex and other independent checks still run after profile or skills
      failure and keep their existing order.
- [x] Doctor behavior tests live in `internal/cli/doctor_test.go`; CLI-wide
      dispatch/help tests remain in `internal/cli/cli_test.go`.
- [x] User guidance matches the shipped text contract and no protected skill
      file changes in this Task.

## Context

- instruction: `docs/agents/cli.md`
- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `.agents/skills/golang-cli/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- instruction: `.agents/skills/tech-writer/SKILL.md`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/cli_test.go`
- interface: `docs/user-guide/commands.md`
- interface: `CONTEXT.md`

## Verification

- `rtk go test ./internal/cli -run 'Test(RunDoctor|Doctor)'` — expected:
  working-directory, output, eager-order, and remediation cases pass.
- `rtk go test -race ./internal/cli -run 'Test(RunDoctor|Doctor)'` — expected:
  Doctor contracts pass under the race detector.
- `rtk grep -n '&& bunx skills experimental_install' docs/user-guide/commands.md` — expected: the public mixed-remediation example uses a fail-closed chain.
- `rtk git diff --exit-code -- .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md` — expected: this non-tooling Task leaves both protected skill files unchanged.

## References

- `_prd.md` → Core Features 3, 4, and 6; User Stories 3–4; User Experience;
  Success Metrics.
- `_techspec.md` → Doctor working directories; Ownership and remediation;
  Documentation and skill synchronization; Testing Approach; Build Order 4.

## Result

Doctor now resolves separate profile and repository working directories:
profile proof falls back to the process working directory outside Git, while
Repository Skill Set inspection runs only for a non-empty trimmed Git root.
Missing-root detail uses the canonical term, and multiple ownership actions
form one owned-then-external `&&` chain without changing the existing
`"; next: "` output boundary.

Doctor behavior tests and their helpers now live in
`internal/cli/doctor_test.go`. The CLI-wide suite retains command dispatch and
help registry coverage. The user guide documents the outside-Git behavior and
the fail-closed mixed-remediation example.

Acceptance criterion evidence:

1. `TestRunDoctorMissingRepositoryRoot` compares the profile work directory
   with the process directory and asserts zero Repository Skill Set checker
   calls.
2. `TestRunDoctorRepositorySkillReadiness` passes a whitespace-padded Git root
   and asserts both profile proof and Repository Skill Set inspection receive
   `/repo/project`.
3. `TestRunDoctorMissingRepositoryRoot` asserts the exact canonical
   `Repository Skill Set readiness requires a Git repository` detail and the
   existing run-from-Git action.
4. `TestRunDoctorRepositorySkillReadiness` covers exact mixed `&&` composition,
   external-only remediation for a symlinked lock error, and conservative
   mixed remediation for an unclassified error.
5. `TestRunDoctorContinuesChecksAfterProfileReadinessFailure` and
   `TestRunDoctorRepositorySkillReadiness` assert eager calls, line order,
   stdout/stderr discipline, and exit codes through the public runner.
6. Source inspection places all `TestRunDoctor*` behavior tests in
   `internal/cli/doctor_test.go`; `TestRunCommandHelp` and
   `TestProfilesDocumentationContractMatchesPublicGuidance` remain in
   `internal/cli/cli_test.go` and passed.
7. The user-guide grep found the shipped `&&` chain, and the protected-skill
   diff check was empty.

Verification:

- `rtk go test ./internal/cli -run 'Test(RunDoctor|Doctor)'` — passed, 18 tests.
- `rtk go test -race ./internal/cli -run 'Test(RunDoctor|Doctor)'` — passed,
  18 tests.
- `rtk go test ./internal/cli -run 'Test(RunCommandHelp|ProfilesDocumentationContractMatchesPublicGuidance)'`
  — passed, 11 tests.
- `rtk grep -n '&& bunx skills experimental_install' docs/user-guide/commands.md`
  — passed; matched the public mixed-remediation example.
- `rtk git diff --exit-code -- .agents/skills/roundfix/SKILL.md skills/roundfix/SKILL.md`
  — passed with no protected skill changes.
- `rtk git diff --check` — passed.

Follow-ups: none.
