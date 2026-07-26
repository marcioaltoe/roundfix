---
task: task_03
spec: 0050-doctor-skill-readiness-hardening
status: pending
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

- [ ] Introduce the anchored repository read boundary.
- [ ] Validate shared ancestors and the lock authority without following links.
- [ ] Adapt owned and external skill inspection to rooted operations.
- [ ] Relocate repository tests and helpers to their owning file.
- [ ] Add ancestor and lock symlink regression fixtures.
- [ ] Prove existing missing, outdated, malformed, and no-mutation behavior.

## Acceptance Criteria

- [ ] A symlinked `.agents` fails readiness and its target is not read.
- [ ] A symlinked `.agents/skills` fails readiness and its target is not read.
- [ ] A symlinked `skills-lock.json` fails readiness and its target is not
      decoded.
- [ ] A complete repository still reports ready; missing and outdated
      classifications remain sorted and unchanged.
- [ ] Nested links and special entries remain rejected without escaping the
      anchored repository.
- [ ] Repository checks remain read-only.
- [ ] `skills/repository.go` has its canonical test file and
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

