---
task: task_04
spec: 0050-doctor-skill-readiness-hardening
status: pending
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

- [ ] Separate missing-Git-root handling from repository checking.
- [ ] Complete unclassified error remediation.
- [ ] Add exact-output missing-root and generic-error cases.
- [ ] Add the real-checker public no-mutation fixture.
- [ ] Restore the canonical Doctor Command wording.
- [ ] Run focused, race, and repository-wide verification.

## Acceptance Criteria

- [ ] An empty loaded Git root never invokes the repository checker or falls
      back to the process working directory.
- [ ] Missing-root and generic checker failures print deterministic `next:`
      actions and exit `1`.
- [ ] All independent checks still run and render in their established order.
- [ ] The public Doctor no-mutation test proves all relevant snapshots are
      byte-identical after execution with the real checker.
- [ ] `CONTEXT.md` again states that Doctor reports the detected acpx version
      against the minimum.
- [ ] Archived Spec 0036 and every protected or upstream-managed path outside
      the approved Task 01 files remain unchanged.
- [ ] The complete repository Verification passes.

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

