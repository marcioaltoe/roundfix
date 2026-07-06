---
task: task_02
spec: 0010-run-robustness
status: pending
type: backend
complexity: medium
---

# Task 02: Session-close reaping in stop --force and the preflight sweep

## Overview

Stop the adapter-orphan leak (forty accumulated in one day): `stop --force`
closes the Run's Agent Sessions after the cancel so acpx terminates its
owner and adapter processes, and the implement preflight sweep closes
discovered roundfix-named sessions belonging to terminal Runs. Close-based,
conservative, never raw kills. Verifiable through the fake rig plus the
gated real-acpx orphan assertion.

## Requirements

1. MUST add session discovery to the agent package: parse
   `acpx <agent> sessions list` output into names, filtered to the
   `roundfix-` prefix, exposing the embedded run id (Run-level
   `roundfix-<run-id>` and per-Task `roundfix-<run-id>-<task_id>` shapes).
2. MUST extend `stop --force`: after the cooperative cancel, close the
   target Run's Run-level session and every discovered per-Task session of
   that Run; each close reports one stderr line; failures are notes, never
   fatal, and never block completion or lock release.
3. MUST extend the implement preflight sweep: discover roundfix-named
   sessions, resolve each embedded run id against the Run Database, and
   close only those whose Runs are terminal — one report line per close;
   sessions with unknown run ids or foreign names are never touched.
4. MUST keep product code free of raw process kills; the reap is entirely
   `sessions close` (acpx owns and terminates its process tree).
5. MUST extend the gated real-acpx integration test
   (`ROUNDFIX_REAL_ACPX=1`): after a force-stop of a live session, assert
   no roundfix-spawned adapter or owner process remains.

## Subtasks

- [ ] Sessions-list parsing with roundfix-prefix filtering
- [ ] stop --force close ordering with per-session reports
- [ ] Preflight sweep close for terminal Runs
- [ ] Conservative-scope tests (foreign and unknown sessions untouched)
- [ ] Gated real-acpx orphan assertion

## Acceptance Criteria

- [ ] Rig tests: force-stop issues cancel then close for Run and Task
      sessions with exact invocations captured; close failure still
      completes the Run Stopped with the note asserted.
- [ ] Sweep tests: terminal-Run sessions closed with report lines;
      active-Run and foreign sessions untouched.
- [ ] The real-acpx gated test passes the no-orphan assertion.
- [ ] Full suite passes.

## Verification

- `rtk go test ./internal/agent/ ./internal/cli/` — expected: all tests
  pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 2, 3; Core Feature 2; Decisions. `_techspec.md` →
Session-close reaping, Build Order 2. Work-plan finding R3-6. ADR-0022
(force-stop shape).
