---
task: task_05
spec: 0002-acpx-migration
status: completed
type: backend
complexity: high
---

# Task 05: Cut over — delete the SDK layer and prove parity

## Overview

Complete ADR-0017: remove the hand-rolled ACP client (the SDK runner, its fs/terminal/permission handlers, and the `coder/acp-go-sdk` module dependency), leaving the acpx runner as the only agent layer, and prove parity with the full gate plus a gated integration test against a real pinned acpx. This is the point of no return the previous tasks staged for; verifiable by the module shrinking and the entire suite passing unchanged.

## Requirements

1. MUST delete the SDK-based runner implementation and every SDK-only support type (fs handlers, terminal handlers, permission auto-approval, teardown plumbing) that has no acpx-path consumer; the stream→Run Event conversion and prompt builders stay.
2. MUST remove `coder/acp-go-sdk` from the module using the module tools (never hand-editing module files), leaving no reference in the module graph.
3. MUST keep the public behavioral contract byte-identical: no stdout, exit-code, state, flag, or journal changes; the full suite passes with no assertion edits beyond fakes already extended in earlier tasks.
4. MUST add one env-guarded integration test, skipped by default, that drives the real pinned acpx with a trivial command-override ACP echo agent to cover: a prompt round-trip with raw NDJSON journaled, cooperative cancel, and crash-resume (kill the adapter mid-Run, next Work Item succeeds) — the PRD's induced-kill metric at the runner seam.
5. MUST verify no orphan acpx owner processes outlive the integration test's Runs (session closed, owner gone or idle-TTL-bound).

## Subtasks

- [x] Delete the SDK runner and SDK-only handlers
- [x] Module dependency removal via the module tools
- [x] Parity sweep: full suite green with unchanged assertions
- [x] Gated real-acpx integration test (round-trip, cancel, crash-resume)
- [x] Orphan-process check in the integration test teardown

## Acceptance Criteria

- [x] `sh -c "! grep -q acp-go-sdk go.mod"` — the module no longer lists the SDK; the agent layer's owned line count shrinks (state the before/after in the Result).
- [x] The full verification gate passes with zero behavioral test edits in this task.
- [x] The gated integration test passes on a machine with Node and the pinned acpx (documented skip otherwise), covering round-trip, cancel, and crash-resume.
- [x] No references to the SDK import path remain anywhere in the repository.

## Verification

- `rtk go test ./...` — expected: full suite passes.
- `sh -c "! grep -q acp-go-sdk go.mod"` — expected: exits 0 (dependency gone).
- `make verify` — expected: full gate passes.

## References

`_prd.md` → User Stories 2, 3; Goals (layer deleted); Success Metrics. `_techspec.md` → System Architecture (deleted), Testing Approach (gated integration test), Build Order 5, Risks. ADR-0017.

## Result

Implemented the hard cutover to acpx as the only Agent layer:

- Deleted the SDK runner and SDK-only support path (`internal/agent/acp_runner.go`), including the hand-rolled SDK client, fs/terminal handlers, permission auto-approval, and SDK teardown code. Kept prompt building and stream-to-Run Event conversion by replacing SDK-typed conversion with raw JSON-RPC session/update decoding.
- Removed the SDK module using `rtk env GOCACHE=/private/tmp/roundfix-gocache go mod tidy`; `rtk go list -m all` no longer includes it.
- Agent layer line count shrank from 4,914 to 3,853 total Go lines. Production Agent code shrank from 2,953 to 1,829 lines.
- Added `TestRealACPXCommandOverrideRoundTripCancelCrashResume`, gated behind `ROUNDFIX_REAL_ACPX=1`. The test builds a trivial command-override ACP echo agent and covers prompt round-trip with raw NDJSON journaling, cooperative cancel, crash-resume after an induced adapter exit, and teardown checks for helper/acpx owner process orphans.
- The local machine has Node (`v25.6.1`) but no `acpx` binary on PATH. With `ROUNDFIX_REAL_ACPX=1`, the integration test skips with `acpx is required for real acpx integration: exec: "acpx": executable file not found in $PATH`; without the env guard it is skipped by default in the full suite.

Verification evidence:

- `rtk go test ./...` passed: 445 tests across 16 packages.
- `rtk sh -c "! grep -q acp-go-sdk go.mod"` exited 0.
- An exact SDK import-path search found no references.
- `rtk env ROUNDFIX_REAL_ACPX=1 GOCACHE=/private/tmp/roundfix-gocache go test ./internal/agent/ -run TestRealACPXCommandOverrideRoundTripCancelCrashResume -count=1 -v` exited 0 with the documented `acpx`-missing skip.
- `rtk make verify` passed: `rtk go test ./...`, `rtk go run -buildvcs=false ./cmd/roundfix skills check`, and `rtk go build -buildvcs=false -o bin/roundfix ./cmd/roundfix`.
