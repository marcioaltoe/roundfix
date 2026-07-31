---
task: task_03
spec: 0058-npm-trusted-publishing-and-release-preflight
status: pending
type: infra
complexity: medium
---

# Task 03: Expose a publish-free preflight rehearsal

## Overview

The preflight from task 02 can otherwise only be exercised by pushing a real
tag, which is exactly the irreversible act it exists to protect. This task adds
a manual trigger that runs tag validation and the preflight against the live
registry and then stops — never cross-compiling, never publishing, never
creating a GitHub Release. Verifiable on its own: the trigger exists, takes a
target version, and every mutating stage is guarded against it.

## Requirements

1. MUST add a manual `workflow_dispatch` trigger alongside the existing `push`
   tag trigger, taking a required version input that supplies the target
   version the preflight evaluates.
2. MUST run semver validation and the Publication Preflight on a dispatch run,
   using the dispatch input as the target version rather than a tag ref.
3. MUST guard every mutating stage — cross-compilation, npm publication, and
   GitHub Release creation — so none of them executes on a dispatch run.
4. MUST leave tag-triggered runs behaving exactly as before, including their
   use of the tag as the version authority.
5. MUST NOT introduce any path by which a dispatch run can publish a package or
   create a release.
6. MUST confine every change to the bounded authorized path
   `.github/workflows/release.yml`.

## Subtasks

- [ ] Add the `workflow_dispatch` trigger with its required version input.
- [ ] Resolve the target version from either the tag ref or the dispatch input.
- [ ] Guard cross-compilation, publication, and release creation against
      dispatch runs.
- [ ] Confirm tag-triggered behavior is unchanged.

## Acceptance Criteria

- [ ] The workflow declares both `push` on `v*` tags and `workflow_dispatch`.
- [ ] The dispatch trigger declares a required input carrying the target
      version.
- [ ] The cross-compilation, publication, and GitHub Release stages each carry
      a condition that excludes dispatch runs.
- [ ] The target version resolves from the tag on a tag run and from the input
      on a dispatch run, with the existing semver check applied to both.
- [ ] `git status --porcelain` shows no path outside
      `.github/workflows/release.yml` and this task file.

## Context

- interface: `.github/workflows/release.yml`

## Verification

- `grep -q 'workflow_dispatch' .github/workflows/release.yml` — expected: exit
  0; the manual trigger exists.
- `grep -q 'required: true' .github/workflows/release.yml` — expected: exit 0;
  the dispatch input is mandatory.
- `grep -q 'tags:' .github/workflows/release.yml` — expected: exit 0; the tag
  trigger survived.
- `grep -A2 'name: Cross-compile and stage' .github/workflows/release.yml | grep -q 'if:'`
  — expected: exit 0; cross-compilation is conditional.
- `grep -A2 'name: Publish to npm' .github/workflows/release.yml | grep -q 'if:'`
  — expected: exit 0; publication is conditional.
- `grep -A2 'name: GitHub Release' .github/workflows/release.yml | grep -q 'if:'`
  — expected: exit 0; release creation is conditional.
- `grep -q 'github.event_name' .github/workflows/release.yml` — expected: exit
  0; the guards discriminate on the trigger.

## References

- `_prd.md` → User Story 2; Success Metrics (workflow stops before
  cross-compilation).
- `_techspec.md` → API Contracts: Workflow trigger; Testing Approach;
  Build Order 4.

