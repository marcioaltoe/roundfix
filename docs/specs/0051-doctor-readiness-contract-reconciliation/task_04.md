---
task: task_04
spec: 0051-doctor-readiness-contract-reconciliation
status: pending
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

- [ ] Add the outside-Git profile-workdir regression and checker non-invocation
      assertion.
- [ ] Restore separate profile and repository root resolution.
- [ ] Correct missing-root detail and mixed remediation composition.
- [ ] Move Doctor-specific tests to the source-matched test file without
      weakening assertions.
- [ ] Update the public Doctor example and explanatory user guidance.
- [ ] Exercise inside-Git, outside-Git, typed ownership, and unclassified error
      paths through the public runner.

## Acceptance Criteria

- [ ] Outside Git, profile proof receives `os.Getwd()` and the skills checker is
      not invoked.
- [ ] Inside Git, profile proof and Repository Skill Set inspection both receive
      the resolved Git root.
- [ ] Missing-root output uses the canonical Repository Skill Set term.
- [ ] Mixed remediation prints one exact `&&`-joined chain and a symlinked lock
      prints only external remediation.
- [ ] Codex and other independent checks still run after profile or skills
      failure and keep their existing order.
- [ ] Doctor behavior tests live in `internal/cli/doctor_test.go`; CLI-wide
      dispatch/help tests remain in `internal/cli/cli_test.go`.
- [ ] User guidance matches the shipped text contract and no protected skill
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
