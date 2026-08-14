---
task: task_03
spec: 0058-npm-trusted-publishing-and-release-preflight
status: completed
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

## Result

Implemented a manual Publication Preflight rehearsal in the authorized release
workflow. `workflow_dispatch` now requires a `version` input, while the release
job resolves `TARGET_VERSION` from that input for dispatch runs and from
`github.ref_name` for tag runs. `Validate tag` and `Publication preflight`
consume the resolved target, and cross-compilation, npm publication, and GitHub
Release creation are each restricted to `push` events. Dispatch validation
reports `preflighting`; the existing tag path still reports `releasing` and the
mutating scripts retain `GITHUB_REF_NAME` as their version authority.

Focused checks:

- `rtk ruby -e '<YAML parse>' .github/workflows/release.yml` parsed the edited
  workflow successfully.
- `rtk ruby -ropen3 -rjson -ryaml -e '<validation path probe>'
  .github/workflows/release.yml` executed the actual YAML-embedded validation
  script. A tag event using the checked-in version exited 0 and reported
  `releasing`; a dispatch using the same version exited 0 and reported
  `preflighting`; an invalid dispatch version exited non-zero with the input
  named in the semver diagnostic.
- `rtk ruby -e '<dispatch structure audit>'
  .github/workflows/release.yml` observed the retained `v*` tag trigger, the
  required dispatch input, trigger-based target resolution, shared validation
  and preflight target use, and a `push`-only condition on all three mutating
  stages.
- `rtk git diff --check` exited 0.

Acceptance evidence:

1. The structural audit observes both `push` with the existing `v*` tag filter
   and `workflow_dispatch` in the workflow trigger.
2. The dispatch trigger exposes a `version` input with `required: true` and
   `type: string`.
3. The structural audit resolves each of `Cross-compile and stage`, `Publish to
   npm`, and `GitHub Release` and proves its condition permits only
   `github.event_name == 'push'`.
4. The target-resolution expression selects `inputs.version` only for a
   dispatch and otherwise selects `github.ref_name`. The validation-path probe
   exercises both sources through the same semver and checked-in-version
   checks, while the structural audit proves Publication Preflight consumes the
   same target. Tag-only mutating scripts remain based on `GITHUB_REF_NAME`.
5. `rtk git -c core.fsmonitor=false status --short` lists only
   `.github/workflows/release.yml` and this task file. The task file's
   pre-existing `status: in_progress` change remains Daemon-owned.

The commands under `## Verification` were not run; the Daemon owns that gate
and Task settlement.
