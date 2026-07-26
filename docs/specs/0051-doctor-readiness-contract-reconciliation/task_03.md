---
task: task_03
spec: 0051-doctor-readiness-contract-reconciliation
status: pending
type: backend
complexity: high
---

# Task 03: Make Repository Skill Set inspection cancellable

## Overview

Propagate the Doctor Command's context through every repository-controlled
blocking Repository Skill Set entry point and preserve cancellation in its
error chain. This slice also completes the typed external lock authority and
removes path-name shadowing inside the same repository inspection boundary.

## Requirements

1. MUST make `CheckRepository` and `SkillFolderHash` context-first APIs and
   migrate every repository caller.
2. MUST propagate the same context through anchored repository inspection,
   owned-skill comparison, external tree collection, and the injected Doctor
   seam.
3. MUST check cancellation before traversal, while walking entries, and before
   reading the next file without spawning an interrupt goroutine.
4. MUST preserve `context.Canceled` and `context.DeadlineExceeded` for
   `errors.Is`, including when wrapped in `RepositoryReadinessError`.
5. MUST classify every `skills-lock.json` authority failure, including a
   symbolic link, as external ownership.
6. MUST rename helper parameters that shadow the imported `path` package
   without changing path validation behavior.
7. MUST preserve anchored no-symlink checks, ownership-specific skill-root
   behavior, offline execution, and no-mutation guarantees.

## Subtasks

- [ ] Add pre-cancelled public API regressions at the owning package layer.
- [ ] Add deterministic cancellation-during-walk or read coverage without a
      production-only test hook.
- [ ] Thread context through repository and folder-hash call chains.
- [ ] Migrate the Doctor dependency seam and every compile-time caller.
- [ ] Correct symlinked lock ownership and path-parameter naming.
- [ ] Re-run repository security, no-mutation, and cancellation cases together.

## Acceptance Criteria

- [ ] A pre-cancelled context returns promptly from both public filesystem APIs
      and remains identifiable with `errors.Is`.
- [ ] Cancellation observed during controlled traversal or reading stops before
      processing the remaining entries.
- [ ] Doctor passes its command context to the Repository Skill Set seam.
- [ ] A symlinked lock produces a typed error with external ownership and never
      reads its target.
- [ ] Existing shared-ancestor, owned-root, and external-root symlink tests keep
      their current blocking classifications.
- [ ] All callers compile with context-first signatures and no contextless
      compatibility wrapper remains.

## Context

- instruction: `docs/agents/go.md`
- instruction: `.agents/skills/golang-context/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `skills/repository.go`
- interface: `skills/repository_test.go`
- interface: `skills/skills.go`
- interface: `skills/skills_test.go`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/cli_test.go`

## Verification

- `rtk go test ./skills -run 'Test(CheckRepository|SkillFolderHash)'` — expected:
  cancellation, ownership, confinement, hash, and no-mutation cases pass.
- `rtk go test ./internal/cli -run 'TestRunDoctor'` — expected: every migrated
  Doctor caller compiles and the existing coordination contract remains green.
- `rtk go test -race ./skills ./internal/cli -run 'Test(CheckRepository|SkillFolderHash|RunDoctor)'` — expected: affected package contracts pass under the race detector.

## References

- `_prd.md` → Core Features 1 and 4; User Story 1; Success Metrics.
- `_techspec.md` → Interfaces; Context propagation and errors; Ownership and
  remediation; Module and naming hygiene; Build Order 3.
