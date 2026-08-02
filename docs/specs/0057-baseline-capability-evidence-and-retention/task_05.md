---
task: task_05
spec: 0057-baseline-capability-evidence-and-retention
status: completed
type: backend
complexity: medium
---

# Task 05: Satisfy a portable Verification role from the repository

## Overview

A portable Verification role that a repository already satisfies with its own
command keeps re-appearing as an unresolved workspace divergence, so a
maintainer sees a gate they have already selected reported as missing. The
projection already distinguishes a repository command from a profile
expectation; this Task lets that mapping satisfy the role.

## Requirements

1. MUST let a portable Verification role be satisfied by mapping it to a
   declared repository command, so the divergence disappears once mapped.
2. MUST require positive evidence: a role reports satisfied only when its
   declared repository command is present, and an unmapped role keeps the
   divergence it produces today.
3. MUST record which command satisfied the role, so the mapping is auditable
   rather than implicit.
4. MUST NOT execute the mapped command; presence is established from
   declaration, not by running it.
5. MUST leave every unmapped role's divergence, message, and blocking status
   unchanged.

## Subtasks

- [ ] Accept a mapping from a portable role to a declared repository command.
- [ ] Satisfy the role only when the declared command is present.
- [ ] Record the satisfying command on the result.
- [ ] Confirm unmapped roles are untouched.

## Acceptance Criteria

- [ ] A portable role mapped to a present declared repository command reports
      satisfied and produces no divergence.
- [ ] The same role with no mapping still produces its current divergence.
- [ ] A mapping naming a command that is not declared does not satisfy the
      role, and reports why.
- [ ] The satisfying command is recorded on the result.
- [ ] No command is executed to establish satisfaction.
- [ ] The characterization corpus shows no change for repositories with no
      mapping.
- [ ] `git status --porcelain` shows no path outside `internal/baseline/` and
      this task file.

## Context

- interface: `internal/baseline/profile_alignment.go`

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline -run TestPortableVerificationRoleMapping -count=1`
  — expected: exit 0; mapped-and-present satisfies, unmapped and
  mapped-but-absent do not.
- `go test ./internal/baseline -run TestBaselinePlanCharacterization -count=1` —
  expected: exit 0; unmapped repositories are unchanged.
- `go test ./internal/baseline -count=1` — expected: exit 0.

## References

- `_prd.md` → Core Features 9 and 12; Success Metrics (a mapped repository gate
  satisfies its role).
- `_techspec.md` → Coverage Map; Build Order 9.

## Result

Implementation:

- `ProfileAlignmentRequest.VerificationRoleMappings` accepts an explicit
  portable-role-to-repository-command mapping at the existing alignment seam.
- A mapped role uses the repository command only when its `package.json` script
  or Make target declaration is present. The projection records the successful
  command in `SatisfiedByCommand` together with its declaration path and digest.
- A mapped-but-undeclared command keeps the role unsatisfied and emits
  `verification.role-mapping.undeclared`. With no mapping, the existing command
  resolution, divergence code, message, next action, requirement, and blocking
  status are unchanged.
- Mapping validation uses declaration inspection only. The regression fixture's
  mapped Make recipe would create a marker if executed, and the marker remains
  absent after alignment.

Focused checks:

- Pre-change signal: the mapped-role subtest failed to compile against the old
  interface because `VerificationRoleMappings` and `SatisfiedByCommand` did not
  exist. The first attempt was environment-blocked by the inherited unwritable
  Go cache; rerunning with `GOCACHE=/private/tmp/roundfix-task05-gocache`
  produced the expected compile failure.
- `rtk gofmt -w internal/baseline/profile_alignment.go internal/baseline/profile_alignment_test.go`
  — exit 0.
- Three separately filtered `TestPortableVerificationRoleMapping` subtest runs
  (`mapped_role_is_satisfied`, `unmapped_role`, and `mapped_role_rejects`) with
  the task-specific writable Go cache — exit 0 for each.
- `Test(ExecutableVerificationCommandRequiresLocalDeclaration|ProfileAlignmentDiscoversDeclaredRepositoryFormatter)$`
  with the task-specific writable Go cache — exit 0; 2 tests passed.
- The five named `TestBaselinePlanCharacterization` subtests with the
  task-specific writable Go cache — exit 0; the existing goldens matched.
- `rtk git -c core.fsmonitor=false status --porcelain` and
  `rtk git -c core.fsmonitor=false diff --name-only` listed only this task file,
  `internal/baseline/profile_alignment.go`, and
  `internal/baseline/profile_alignment_test.go`.

Acceptance evidence:

- Mapped and present: the `workspace` role maps to declared `make verify`, has
  no `verification.workspace` divergence, and records the command, Makefile
  path, and declaration digest.
- Unmapped: the same role retains the exact existing
  `verification.profile-expectation.unresolved` advisory divergence, message,
  next action, and non-blocking status.
- Mapped but absent: `make missing` does not satisfy the role, records no
  satisfying command, and reports the missing local declaration.
- Non-execution: the marker-producing Make recipe is never run during
  alignment.
- Characterization: all five no-mapping corpus cases matched their existing
  golden results without regeneration.
- Scope: the changed-path postflight contains no path outside
  `internal/baseline/` and this task file.

Daemon Verification was not run in this Agent turn.
