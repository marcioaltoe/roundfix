---
task: task_03
spec: 0010-run-robustness
status: completed
type: backend
complexity: high
---

# Task 03: Detached Runs

## Overview

Implement ADR-0028: `--detach` on resolve, watch, and implement re-executes
roundfix as a session leader independent of the caller, hands the run id
back through a startup handshake, prints the deterministic detach report,
and exits 0 — so a Run's survival never again depends on the invoking
process. Verifiable through handshake-protocol tests and a caller-killed
integration case.

## Requirements

1. MUST add `--detach` to resolve, watch, and implement: rejected with
   `--interactive` (existing conflicting-flags shape); implies
   non-interactive semantics exactly like `--no-input`.
2. MUST re-exec the binary (`os.Executable`) with the same arguments minus
   `--detach`, as a new session (`SysProcAttr{Setsid: true}`), stdio
   detached from the caller's terminal, with an internal handshake channel
   (pipe fd or temp file passed by env — pick the cleaner, test-driven
   option) never exposed as a public flag.
3. MUST have the child write `run_id` and console-log path to the handshake
   immediately after Run creation, then continue as a normal non-TTY Run;
   the child's console stream lands in
   `<artifact_dir>/runs/<run-id>/console.log` (temp-then-rename covering
   the pre-run-id window) regardless of any future log gating.
4. MUST have the caller print, on handshake success, exactly four stdout
   lines — `Run detached: <run-id>`, `Console log: <path>`,
   `Follow: roundfix attach <run-id>`, `Stop: roundfix stop <run-id>` — and
   exit 0.
5. MUST relay child failures that precede the handshake verbatim: the
   caller waits bounded, then surfaces the child's captured stderr and exit
   code (a detached launch never hides a Preflight message and never
   hangs).
6. MUST prove caller-independence: an integration-style test detaches,
   kills the caller process group, and asserts the child reaches a
   journaled terminal outcome and stays attachable throughout.

## Subtasks

- [x] Flag wiring, conflicts, and non-interactive implication
- [x] Re-exec with Setsid and the internal handshake channel
- [x] Console-log binding with the temp-then-rename window
- [x] Caller report shape and pre-handshake relay semantics
- [x] Caller-killed survival integration test

## Acceptance Criteria

- [x] Detach success prints the four lines byte-exactly and exits 0; the
      child Run completes with journal, worktree, and integration behavior
      identical to a foreground Run.
- [x] A Preflight failure under `--detach` (e.g. invalid spec) reaches the
      caller's stderr verbatim with exit 2 and no orphan child.
- [x] The caller-killed test proves the Run survives and terminal state
      lands in the Run Database.
- [x] `--detach --interactive` fails with the existing conflict shape;
      help text documents the flag on all three commands.
- [x] Full suite passes.

## Verification

- `rtk go test ./internal/cli/` — expected: all tests pass.
- `rtk go run ./cmd/roundfix implement --help` — expected: `--detach`
  documented, exit 0.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 4, 5; Core Feature 3; Decisions. `_techspec.md` →
Detach, Risks, Build Order 3. ADR-0028. Work-plan finding R3-7.

## Result

Implemented `--detach` for `resolve`, `watch`, and `implement`. The caller
now launches `os.Executable()` as a new session leader with `--detach`
removed, passes an internal handshake fd/env pair, and exits after printing
the four-line detach report. The detached child writes its console stream to
a temp log under the Artifact Directory, renames it to
`runs/<run-id>/console.log` immediately after Run creation, writes the
handshake, and continues through the normal non-interactive Run path.

Evidence:

- `rtk go test ./internal/cli/`: passed, 255 tests.
- `rtk go run ./cmd/roundfix implement --help`: passed, exit 0; output
  documents `--detach`.
- `rtk go test ./...`: passed, 728 tests across 17 packages.
- `rtk make verify`: passed; it ran the full test suite, `roundfix skills
  check`, and the build.

Acceptance evidence:

- Byte-exact detach report and Run completion: covered by
  `TestRunImplementDetachPrintsReportAndCompletesRun`, which asserts the
  four stdout lines, exit 0, Clean Run state, removed clean Run Worktree,
  console log content, and a journaled `daemon.outcome` Clean event.
- Preflight relay: covered by
  `TestRunImplementDetachRelaysPreflightFailureVerbatim`, which compares
  detached invalid-Spec stderr byte-for-byte with the foreground Preflight
  failure and verifies exit 2 with no Run Database.
- Caller independence: covered by
  `TestRunImplementDetachSurvivesCallerProcessGroupKill`, which kills the
  caller process group after the handshake, attaches to the still-active Run,
  releases the blocked fake acpx prompt, and verifies a Clean terminal state
  plus journaled outcome.
- Flag conflict and help text: covered by
  `TestRunDetachRejectsInteractiveWithExistingConflictShape`,
  `TestRunCommandHelp`, and
  `TestRunImplementHelpListsExactlyImplementedFlags`.

Follow-up:

- Task 04 owns README and canonical Roundfix skill documentation for the new
  Detached Run usage pattern.
