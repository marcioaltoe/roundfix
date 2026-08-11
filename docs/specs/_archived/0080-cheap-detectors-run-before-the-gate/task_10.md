---
task: task_10
spec: 0080-cheap-detectors-run-before-the-gate
status: completed
type: test
complexity: low
---

# Task 10: Give the prompt fixture the Git precondition the stage now needs

## Overview

Task 03 put the mechanical stage in front of the QA Agent Session, and the stage
resolves the repository head before anything else runs.
`TestQAGatePromptUsesTaskContextBuilderAndPreviousReportIdentity` builds its
fixture from an ordinary temp directory with no Git repository, so
`git rev-parse HEAD` exits 128, the TaskCycle stops before the prompt is
assembled, and the test never observes what it exists to observe. Both
verification tiers exit 2.

The QA gate recorded this as F-003 with `Blocks-Completion` impact and named the
cause precisely: Task 03 added a real Git-dependent stage without adapting this
fixture, and Task 09 repaired seeded-report selection without touching it.

## Requirements

1. MUST give the fixture a real Git repository with a resolvable head, using the
   repository's existing test Git helpers rather than a hand-rolled one.
2. MUST NOT bypass, stub, or skip production head resolution. The point of the
   test is that the prompt is assembled after the stage the production path
   runs; a fixture that avoids the stage proves nothing about that order.
3. MUST keep every assertion the test already makes about the Task Context
   builder and the previous QA Report identity.
4. MUST leave production code unchanged. If the stage genuinely cannot run
   against a fixture that a real repository can satisfy, stop and report which
   production assumption is unreachable rather than weakening the test.

## Subtasks

- [ ] Initialise the fixture repository with a resolvable head.
- [ ] Keep the existing prompt assertions intact.

## Acceptance Criteria

- [ ] The focused test passes with the mechanical stage running.
- [ ] No assertion about builder context or previous-report identity was removed.
- [ ] `git diff --name-only` lists only test files and this Task file.

## Bounded scope

This Task may create or modify only:

- `internal/daemon/task_context_test.go`
- `docs/specs/0080-cheap-detectors-run-before-the-gate/task_10.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/daemon -run '^TestQAGatePromptUsesTaskContextBuilderAndPreviousReportIdentity$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestQAGatePromptUsesTaskContextBuilderAndPreviousReportIdentity'` — expected: exits 0. The test fails with `exit status 128` against the unchanged tree.

## References

- `_prd.md` → Goal 1.
- `task_03.md` → the mechanical stage this fixture must satisfy.
- `qa/qa-report-2026-08-11-01.md` → F-003.

## Result

The QA prompt fixture now initializes `fixture.gitRoot` through
`gittest.InitRepo`, stages its fixture files with `gittest.Run`, and creates a
deterministic initial commit. Production `HEAD` resolution and the real
mechanical stage remain active. The prior-report fixture is also a minimal
structurally valid closed fail report, so the report-shape detector can inspect
it without blocking the prompt path; its filename, verdict, and asserted
previous-report identity remain unchanged.

Focused checks run after the implementation edit:

- Pre-change signal: `rtk proxy env GOCACHE=/private/tmp/roundfix-task10-gocache go test ./internal/daemon -run 'TestQAGatePromptUsesTaskContextBuilderAndPreviousReportIdentity$' -count=1`
  exited 1 with `resolve mechanical-stage Git head: exit status 128`.
- The same focused command exited 0 after the fixture gained its committed
  head and valid prior-report structure: `ok roundfix/internal/daemon`.

Acceptance evidence:

1. The focused test passes while using the default production mechanical stage;
   no stage dependency was replaced or stubbed.
2. The final diff preserves the complete assertion loops for Task Context
   bundle content, previous QA Report path and head, and prior-change resolver
   calls. The test change only strengthens fixture setup.
3. `rtk git diff --name-only` exited 0 and listed only
   `internal/daemon/task_context_test.go` and this Task file. `rtk git diff
   --check` also exited 0 with no diagnostics.

The command under this Task's `## Verification` was not run; the Daemon owns
that exact command and Task settlement.
