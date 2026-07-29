---
task: task_03
spec: 0061-repository-derived-skill-requirements
status: pending
type: backend
complexity: low
---

# Task 03: Report every missing skill with a per-skill install command

## Overview

A missing external skill currently fails on the first gap and prints a
package-wide install that pulls the entire upstream catalog. Report the whole
gap at once and name the command that installs exactly the missing skills.

## Requirements

1. MUST report every missing external skill in one failure rather than
   stopping at the first.
2. MUST print a per-skill install command for each missing skill, using the
   upstream CLI's skill-scoped form.
3. MUST NOT print a package-wide install as the remediation for a skill-level
   gap.
4. MUST keep the existing failure classification, ownership label, and exit
   behavior unchanged.
5. SHOULD keep the owned-skill remediation exactly as it is today.

## Subtasks

- [ ] Collect every missing external skill before failing.
- [ ] Render the per-skill install commands in the next action.
- [ ] Cover a single gap, several gaps, and the unchanged owned remediation.

## Acceptance Criteria

- [ ] A repository missing three external skills names all three in one
      failure.
- [ ] Each named skill appears with a command that installs only it.
- [ ] No printed next action installs a whole package for a skill-level gap.
- [ ] The owned-skill failure and its remediation are unchanged.

## Context

- interface: `skills/repository.go`
- interface: `internal/cli/doctor.go`
- interface: `internal/cli/doctor_test.go`

## Verification

- `go build -buildvcs=false ./...` — expected: clean build.
- `go test -count=1 ./skills/ ./internal/cli/ -run 'TestRunDoctor|TestCheckRepository'` — expected: pass with the multi-gap coverage.

## References

`_prd.md` → User Story 2, Core Feature 4; `_techspec.md` → Build Order 3,
API Contracts.
