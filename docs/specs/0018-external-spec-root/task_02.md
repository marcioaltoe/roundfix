---
task: task_02
spec: 0018-external-spec-root
status: pending
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

- [ ] Spec loading signature shift and task-path resolution against the root
- [ ] Implement Command threading: preflight, execution reload, status
      writes, QA
- [ ] Settle, Archive, Attach, and Interactive Input threading
- [ ] Review artifact locations follow the root
- [ ] Startup report line for a non-default root
- [ ] Tests: spec package against an external temp root; CLI implement,
      settle, and archive end to end with a configured external root

## Acceptance Criteria

- [ ] The Implement Command completes a Run in a test repository whose Spec
      Root is a directory outside the repository, reading the Task Graph and
      writing task statuses and QA Reports there — including from the Run
      Worktree.
- [ ] Settle and Archive operate on the same external root.
- [ ] The active-Spec interactive listing shows Specs from the external root.
- [ ] With the default layout, the full existing test suite passes unchanged.
- [ ] Run startup names the resolved root on stderr when it is not the
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
