---
task: task_02
spec: 0018-external-spec-root
status: completed
type: backend
complexity: high
---

# Task 02: Spec Root threading through every Spec consumer

## Overview

Replace the derived `<execution-root>/docs/specs` assumption with the resolved
Spec Root, threaded explicitly from command start into every Spec consumer.
After this task, an external root works end to end for reads and writes —
Task Graph loading, task status updates, QA Reports, review artifacts, and the
interactive listings — from the user checkout and from Run and Task Worktrees
alike. Demoable by running the Implement Command with a configured external
root in a test repository.

## Requirements

1. MUST shift Spec loading to take the resolved Spec Root and resolve task
   file paths against it, never against the execution work directory.
2. MUST resolve the Spec Root once per command against the user's checkout
   and carry that absolute path into the Run, so Run and Task Worktrees read
   and write the same directory the checkout does.
3. MUST thread the root through every consumer: the Implement Command
   (preflight load, execution reload, task status writes, QA Reports), the
   Settle Command, the Archive Command, Attach's Task detail, Interactive
   Input's active-Spec listing, and the Spec-associated review artifact
   locations.
4. MUST report a non-default resolved Spec Root on stderr at Run startup.
5. MUST keep default-layout behavior byte-for-byte unchanged, including
   reports, exit codes, and commit contents.

## Subtasks

- [x] Spec loading signature shift and task-path resolution against the root
- [x] Implement Command threading: preflight, execution reload, status
      writes, QA
- [x] Settle, Archive, Attach, and Interactive Input threading
- [x] Review artifact locations follow the root
- [x] Startup report line for a non-default root
- [x] Tests: spec package against an external temp root; CLI implement,
      settle, and archive end to end with a configured external root

## Acceptance Criteria

- [x] The Implement Command completes a Run in a test repository whose Spec
      Root is a directory outside the repository, reading the Task Graph and
      writing task statuses and QA Reports there — including from the Run
      Worktree.
- [x] Settle and Archive operate on the same external root.
- [x] The active-Spec interactive listing shows Specs from the external root.
- [x] With the default layout, the full existing test suite passes unchanged.
- [x] Run startup names the resolved root on stderr when it is not the
      default.

## Verification

- `rtk go test ./internal/spec/ ./internal/cli/` — expected: all tests pass,
  including the new external-root tests.
- `rtk make verify` — expected: fmt-check, tests, skills check, and build all
  pass.

## References

`_prd.md` → User Story 2; Core Features 2-3. `_techspec.md` → System
Architecture; Interfaces: Load; Build Order 2; Risks (signature shift).
ADR-0035.

## Result

- Spec loading now accepts an explicit Spec Root, stores Task file paths
  relative to that root, and reloads/writes task files by joining against the
  threaded root. Evidence: `TestLoadUsesExplicitExternalSpecRoot`.
- Implement resolves `specs.root` once from the user checkout, maps default
  internal roots into Run/Task Worktrees, keeps external roots absolute, and
  passes that root through preflight load, execution reload, status writes,
  QA, Attach, and TUI task detail. Evidence:
  `TestRunImplementUsesConfiguredExternalSpecRootEndToEnd`.
- Settle and Archive use the configured external root for Spec loading and
  artifact writes. Evidence: `TestRunSettleUsesConfiguredExternalSpecRoot`
  and `TestRunArchiveUsesConfiguredExternalSpecRoot`.
- Interactive Input lists active Specs from the configured external root.
  Evidence: `TestRunImplementInteractiveInputListsConfiguredExternalSpecRoot`.
- Spec-associated review artifacts resolve under the configured Spec Root.
  Evidence: `TestResolveReviewRoot/existing_spec_under_external_root_stores_rounds_under_external_spec_reviews`.
- Non-default Spec Root startup reporting is on stderr for Implement Runs.
  Evidence: `TestRunImplementUsesConfiguredExternalSpecRootEndToEnd` asserts
  `Spec Root: <external-root>`.
- Verification passed:
  - `rtk go test ./internal/spec/ ./internal/cli/` — 357 tests passed in 2
    packages.
  - `rtk make verify` — `rtk go test ./...` passed 856 tests in 18 packages;
    `roundfix skills check` passed; build completed.

Follow-up note for Task 03: external task and QA artifact paths are threaded
to the commit boundary, but filtering external paths out of staging and the
settle-without-commit behavior remain Task 03's slice.
