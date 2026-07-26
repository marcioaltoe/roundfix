---
task: task_03
spec: 0051-doctor-readiness-contract-reconciliation
status: completed
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

- [x] Add pre-cancelled public API regressions at the owning package layer.
- [x] Add deterministic cancellation-during-walk or read coverage without a
      production-only test hook.
- [x] Thread context through repository and folder-hash call chains.
- [x] Migrate the Doctor dependency seam and every compile-time caller.
- [x] Correct symlinked lock ownership and path-parameter naming.
- [x] Re-run repository security, no-mutation, and cancellation cases together.

## Acceptance Criteria

- [x] A pre-cancelled context returns promptly from both public filesystem APIs
      and remains identifiable with `errors.Is`.
- [x] Cancellation observed during controlled traversal or reading stops before
      processing the remaining entries.
- [x] Doctor passes its command context to the Repository Skill Set seam.
- [x] A symlinked lock produces a typed error with external ownership and never
      reads its target.
- [x] Existing shared-ancestor, owned-root, and external-root symlink tests keep
      their current blocking classifications.
- [x] All callers compile with context-first signatures and no contextless
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

## Result

Implemented context-first Repository Skill Set inspection. `CheckRepository`
and `SkillFolderHash` now propagate the Doctor Command context through anchored
repository checks, owned comparisons, external walks, and file reads. Walks
check cancellation before each entry and immediately before and after file
reads, without goroutines or test-only production hooks. Cancellation remains
discoverable through wrapped errors.

`skills-lock.json` symlinks now return `RepositoryReadinessError` with external
ownership before the target can be decoded. Shared ancestors and owned and
external skill roots retain their prior blocking behavior. Repository helper
parameters no longer shadow the imported `path` package, and every in-repository
caller uses the context-first signatures.

Acceptance evidence:

- `TestCheckRepositoryHonorsPreCanceledContext` and
  `TestSkillFolderHashHonorsPreCanceledContext` cover canceled and expired
  contexts before a missing root can be inspected; `errors.Is` matches both
  `context.Canceled` and `context.DeadlineExceeded`, including a
  `RepositoryReadinessError` chain.
- `TestSkillFolderHashStopsAfterCancellationDuringRead` cancels through a real
  `fs.FS` read boundary and proves that only the first of two files opens.
- `TestRunDoctorPassesCommandContextToRepositorySkillReadiness` proves the
  injected seam receives the command context and resolved repository root.
- `TestCheckRepositoryRejectsSymlinkedAuthoritiesBeforeReadingTargets` proves
  lock symlinks are externally owned, shared ancestors stay unclassified, and
  the malformed target is never decoded. Existing owned and external symlink
  and no-mutation cases pass in the focused package run.
- Repository-wide compilation and the caller search confirm that only the
  context-first exported signatures remain.

Verification:

- `GOCACHE=/private/tmp/roundfix-task03-go-cache rtk go test ./skills -run 'Test(CheckRepository|SkillFolderHash)'`:
  passed.
- `GOCACHE=/private/tmp/roundfix-task03-go-cache rtk go test ./internal/cli -run 'TestRunDoctor'`:
  passed.
- `GOCACHE=/private/tmp/roundfix-task03-go-cache rtk go test -race ./skills ./internal/cli -run 'Test(CheckRepository|SkillFolderHash|RunDoctor)'`:
  passed.
- `GOCACHE=/private/tmp/roundfix-task03-go-cache rtk make verify`: passed;
  2,420 tests across 23 packages, skill checks, and build completed.
- `rtk git -c core.fsmonitor=false diff --check`: passed.

The first focused command using the host-default Go cache could not open its
sandbox-inaccessible cache path. The same command with a writable isolated
`GOCACHE` passed; the Daemon remains authoritative for the task's verbatim
Verification commands.
