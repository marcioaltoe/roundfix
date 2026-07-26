---
task: task_03
spec: 0050-doctor-skill-readiness-hardening
status: completed
type: backend
complexity: high
---

# Task 03: Anchor Repository Skill Set filesystem reads

## Overview

Make the supplied Git root the enforced filesystem boundary for every
Repository Skill Set authority. Static symbolic links in the shared
skill-tree ancestors or lock authority must fail before their targets are
read, while existing missing and outdated classifications remain stable.

## Requirements

1. MUST open one `os.Root` for the supplied repository and perform later
   authority reads through it.
2. MUST inspect `.agents`, `.agents/skills`, and `skills-lock.json` with
   non-following metadata operations before opening or reading them.
3. MUST reject a symlinked shared ancestor or lock authority without reading
   its target.
4. MUST inspect each required skill root without following its final path and
   keep nested symlink and special-file rejection.
5. MUST preserve missing-skill classification when `.agents`,
   `.agents/skills`, or an individual skill does not exist.
6. MUST wrap filesystem causes so callers retain `errors.Is` and `errors.As`
   behavior and ownership-specific remediation where ownership is known.
7. MUST move `CheckRepository` tests and their repository-specific helpers to
   `skills/repository_test.go` without weakening existing assertions.
8. MUST keep `SkillFolderHash` tests beside `skills.go`.

## Subtasks

- [x] Introduce the anchored repository read boundary.
- [x] Validate shared ancestors and the lock authority without following links.
- [x] Adapt owned and external skill inspection to rooted operations.
- [x] Relocate repository tests and helpers to their owning file.
- [x] Add ancestor and lock symlink regression fixtures.
- [x] Prove existing missing, outdated, malformed, and no-mutation behavior.

## Acceptance Criteria

- [x] A symlinked `.agents` fails readiness and its target is not read.
- [x] A symlinked `.agents/skills` fails readiness and its target is not read.
- [x] A symlinked `skills-lock.json` fails readiness and its target is not
      decoded.
- [x] A complete repository still reports ready; missing and outdated
      classifications remain sorted and unchanged.
- [x] Nested links and special entries remain rejected without escaping the
      anchored repository.
- [x] Repository checks remain read-only.
- [x] `skills/repository.go` has its canonical test file and
      `skills/skills_test.go` retains only `skills.go`-owned coverage and
      genuinely shared helpers.

## Context

- instruction: `docs/agents/go.md`
- instruction: `.agents/skills/coding-guidelines/SKILL.md`
- instruction: `.agents/skills/golang-error-handling/SKILL.md`
- instruction: `.agents/skills/golang-testing/SKILL.md`
- instruction: `.agents/skills/testing-boss/SKILL.md`
- interface: `skills/repository.go`
- interface: `skills/skills.go`
- interface: `skills/skills_test.go`

## Verification

- `rtk go test ./skills -run 'TestCheckRepository' -count=1` — expected:
  ready, missing, outdated, malformed, no-mutation, shared-ancestor symlink,
  lock symlink, and nested-link cases pass.
- `rtk go test -race ./skills -run 'TestCheckRepository' -count=1` — expected:
  anchored Repository Skill Set reads are race-free.

## References

- `_prd.md` → Goal 1; User Story 1; Core Features 1 and 5; Success Metrics.
- `_techspec.md` → Root-anchored repository reads; Test ownership and
  no-mutation proof; Testing Approach; Build Order 3.

## Result

Implemented one `os.Root` boundary for the supplied Git root. The repository
checker now uses rooted non-following metadata checks for `skills-lock.json`,
`.agents`, `.agents/skills`, and every required skill root before reading or
walking them. Shared-authority links fail with an unclassified
`RepositoryReadinessError`; owned and external skill failures retain their
ownership when it is known, and wrapped filesystem causes remain available to
`errors.Is` and `errors.As`.

Repository checks reuse the root-backed filesystem for owned comparison and
external hashing. Missing shared directories classify every required skill as
missing, while ready, outdated, malformed, nested-link, special-entry, stable
sorting, and no-mutation behavior remain covered.

Moved every `TestCheckRepository*` case and repository fixture helper to
`skills/repository_test.go`. `skills/skills_test.go` retains
`TestSkillFolderHash*`, the hash compatibility fixture, and the shared file
writer used by both canonical suites.

Verification:

- `rtk go test ./skills -run 'TestCheckRepository' -count=1` — passed,
  26 tests.
- `rtk go test -race ./skills -run 'TestCheckRepository' -count=1` — passed,
  26 tests.
- `rtk go test ./skills -count=1` — passed, 106 tests.
- `rtk make verify` — passed: 2,406 tests across 23 packages, 4 focused skill
  contract tests, Repository Skill Set check, and the Roundfix build.
- `rtk git -c core.fsmonitor=false diff --check` — passed.

Acceptance evidence:

- `TestCheckRepositoryRejectsSymlinkedAuthoritiesBeforeReadingTargets` proves
  `.agents`, `.agents/skills`, and `skills-lock.json` links fail at the
  authority path; the malformed lock target is never decoded.
- `TestCheckRepositoryReportsReadyRequiredSetWithoutMutation`,
  `TestCheckRepositoryClassifiesMissingAndOutdatedSkills`, and
  `TestCheckRepositoryClassifiesMissingSharedSkillDirectories` prove ready,
  read-only, missing, outdated, and sorted classification behavior.
- `TestCheckRepositoryHandlesNestedLinksSpecialEntriesAndStableOrdering`
  proves nested links and special entries remain rejected within the anchored
  repository view.
- `TestCheckRepositoryWrapsFilesystemCauses` proves `errors.Is` and
  `errors.As` survive repository-root failures.
- Source inspection found one `os.OpenRoot` call and no direct repository
  authority `os.ReadFile`, `os.Lstat`, or `os.DirFS` calls in
  `skills/repository.go`; test-name inspection places all
  `TestCheckRepository*` coverage in `skills/repository_test.go` and all
  `TestSkillFolderHash*` coverage in `skills/skills_test.go`.

Follow-ups: none.
