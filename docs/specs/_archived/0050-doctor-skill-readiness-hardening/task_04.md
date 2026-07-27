---
task: task_04
spec: 0050-doctor-skill-readiness-hardening
status: completed
type: backend
complexity: medium
---

# Task 04: Harden Doctor coordination and evidence

## Overview

Make Doctor use the resolved Git root without a process-directory fallback,
complete deterministic error remediation, and prove the full public command is
read-only with the real Repository Skill Set checker. Reconcile the canonical
Doctor Command definition without changing archived Specs or current output
delimiters.

## Requirements

1. MUST remove Doctor's `os.Getwd()` fallback for Repository Skill Set
   readiness.
2. MUST avoid calling `checkSkills` when `roundconfig.Loaded.GitRoot` is empty
   and print repository-specific remediation instead.
3. MUST preserve eager execution and ordering of Node, acpx, Adapter
   Readiness, Agent Selection Profile Readiness, Repository Skill Set, and
   codex checks.
4. MUST provide both existing ownership remediation commands, in their current
   order, when an unclassified checker error has no narrower safe action.
5. MUST preserve the `"; next: "` boundary, stdout/stderr placement, and
   Doctor exit codes.
6. MUST add a public `Run([]string{"doctor"}, ...)` no-mutation test that uses
   the real repository checker and snapshots repository, User Config,
   `.roundfix`, lock, and skill state.
7. MUST restore the canonical Doctor Command wording for detected acpx version
   reporting.
8. MUST leave the archived Spec 0036, upstream-managed skills, current lock,
   and branch history unchanged.

## Subtasks

- [x] Separate missing-Git-root handling from repository checking.
- [x] Complete unclassified error remediation.
- [x] Add exact-output missing-root and generic-error cases.
- [x] Add the real-checker public no-mutation fixture.
- [x] Restore the canonical Doctor Command wording.
- [x] Run focused, race, and repository-wide verification.

## Acceptance Criteria

- [x] An empty loaded Git root never invokes the repository checker or falls
      back to the process working directory.
- [x] Missing-root and generic checker failures print deterministic `next:`
      actions and exit `1`.
- [x] All independent checks still run and render in their established order.
- [x] The public Doctor no-mutation test proves all relevant snapshots are
      byte-identical after execution with the real checker.
- [x] `CONTEXT.md` again states that Doctor reports the detected acpx version
      against the minimum.
- [x] Archived Spec 0036 and every protected or upstream-managed path outside
      the approved Task 01 files remain unchanged.
- [x] The complete repository Verification passes.

## Context

- instruction: `CONTEXT.md`
- instruction: `docs/agents/cli.md`
- instruction: `.agents/skills/agentic-cli-design/SKILL.md`
- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/golang-cli/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/cli_test.go`
- interface: `skills/repository.go`

## Verification

- `rtk go test ./internal/cli -run 'TestRunDoctor' -count=1` — expected:
  exact output, missing-root handling, generic remediation, independent check
  ordering, and public no-mutation cases pass.
- `rtk go test -race ./internal/skillhash ./skills ./internal/baseline ./internal/cli -run 'Test(Sum|SkillFolderHash|CheckRepository|SkillsRestore|RunDoctor)' -count=1` — expected:
  the assembled correction is race-free across affected packages.
- `rtk make verify` — expected: formatting, tests, Repository Skill Set
  integrity, and build all pass.

## References

- `_prd.md` → Goals 4–6; User Stories 4–5; Core Features 4–6; Success Metrics.
- `_techspec.md` → Doctor coordination and remediation; Test ownership and
  no-mutation proof; Contract reconciliation; Testing Approach; Build Order 4.

## Result

Doctor now treats the loaded Git root as the sole Repository Skill Set
authority. An empty root produces a repository-specific failed result without
calling the checker, and an unclassified checker error falls back to both
ownership remediation commands in their established order. The public check
sequence executes eagerly in its rendered order.

The public `Run([]string{"doctor"}, ...)` regression fixture uses
`skills.CheckRepository` against a complete disposable Repository Skill Set.
It compares repository, User Config, user and repository `.roundfix`, lock,
and skill-tree snapshots before and after the command. `CONTEXT.md` also
restores the detected-acpx-version wording.

Acceptance evidence:

- AC 1–3: `TestRunDoctorMissingRepositoryRoot` and
  `TestRunDoctorRepositorySkillReadiness` assert checker call counts, exact
  stdout, exit `1`, the `"; next: "` boundary, and the ordered Node, acpx,
  Adapter Readiness, Agent Selection Profile Readiness, Repository Skill Set,
  and codex calls.
- AC 4: `TestRunDoctorRealRepositoryCheckDoesNotMutateState` passes with the
  real repository checker and byte-identical snapshots for every required
  state surface.
- AC 5: `CONTEXT.md` states that the Doctor Command reports the detected acpx
  version against the minimum.
- AC 6: the changed-file postflight contains only `CONTEXT.md`, this Task file,
  `internal/cli/doctor.go`, and `internal/cli/cli_test.go`; archived Spec 0036,
  upstream-managed skills, `skills-lock.json`, `skills/recommended.txt`,
  `go.mod`, and `go.sum` have no Task 04 diff. Branch history remains at
  `be91ebc`.
- AC 7: focused, affected-package race, and full repository gates passed.

Verification:

- Initial red signal:
  `rtk go test ./internal/cli -run 'TestRunDoctor' -count=1` exposed the
  process-directory fallback, precomputed check order, and missing generic
  remediation. The sandbox then required a task-specific `GOCACHE`.
- `rtk env GOCACHE=/private/tmp/roundfix-task04-go-cache go test ./internal/cli -run 'TestRunDoctor' -count=1`
  — passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task04-go-cache go test -race ./internal/skillhash ./skills ./internal/baseline ./internal/cli -run 'Test(Sum|SkillFolderHash|CheckRepository|SkillsRestore|RunDoctor)' -count=1`
  — passed across all four affected packages.
- `rtk env GOCACHE=/private/tmp/roundfix-task04-go-cache make verify` — passed:
  2,409 tests across 23 packages, four skill contract tests, Repository Skill
  Set integrity, and the build.
- `rtk git -c core.fsmonitor=false diff --check` — passed.

Follow-ups: none.
