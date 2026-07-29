---
task: task_03
spec: 0061-repository-derived-skill-requirements
status: completed
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

## Result

### Implementation

- Doctor now turns each named missing or outdated external skill into
  `bunx skills add marcioaltoe/skills@<skill>`, sorted and deduplicated before
  joining the next actions.
- External checker failures that do not identify a skill keep the existing
  generic external recovery action. The existing Roundfix-owned recovery
  action, failed status, ownership classification, and Doctor exit behavior
  remain unchanged.
- The repository readiness coverage now removes three explicitly required
  external skills and proves the checker returns all three in one readiness
  result; an additional removed recommendation remains outside that explicit
  requirement.

### Focused checks

- The first focused test attempt reached neither package because the default
  macOS Go cache was sandbox-denied. The same commands were rerun unchanged
  with `GOCACHE=/private/tmp/roundfix-task03-gocache`.
- Pre-change signal:
  `rtk proxy env GOCACHE=/private/tmp/roundfix-task03-gocache rtk go test -count=1 ./internal/cli -run '^TestRunDoctorRepositorySkillReadiness$'`
  failed in the single-gap, three-gap, and mixed-gap cases because Doctor
  still printed the package-wide external action. The companion repository
  checker case passed before the Doctor change, confirming that collection
  already occurred below the rendering layer.
- `rtk gofmt -w internal/cli/doctor.go internal/cli/doctor_test.go skills/repository_test.go`
  — passed.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task03-gocache rtk go test -count=1 ./skills -run '^TestCheckRepositoryWithExternalUsesExplicitRequirement$'`
  — passed, 1 focused test.
- `rtk proxy env GOCACHE=/private/tmp/roundfix-task03-gocache rtk go test -count=1 ./internal/cli -run '^TestRunDoctorRepositorySkillReadiness$'`
  — passed, all 7 table cases.
- `rtk git diff --check` — passed.

### Acceptance criteria evidence

- Three missing external skills:
  `TestCheckRepositoryWithExternalUsesExplicitRequirement` returns
  `agentic-cli-design`, `autoresearch`, and `bubbletea` together, while the
  `several_external_gaps_each_use_skill-scoped_remediation` Doctor subtest
  prints `agentic-cli-design`, `golang-testing`, and `testing-boss` in one
  failed `skills:` line.
- Per-skill commands: that Doctor subtest asserts one exact
  `bunx skills add marcioaltoe/skills@<skill>` action for each of its three
  names; the single-gap subtest asserts the same form for `testing-boss`.
- No package-wide skill-gap action: the exact expected output for the single,
  several, and mixed named-gap cases excludes
  `bunx skills experimental_install && bunx skills update -p -y`.
- Owned remediation unchanged: the owned-only subtest still requires exactly
  `roundfix skills install --target project`; the mixed case keeps that action
  first, and every failure case still expects `exitRunFailed`.

Daemon Verification was not run; the Daemon owns the commands in this Task's
`## Verification` section and terminal settlement.
