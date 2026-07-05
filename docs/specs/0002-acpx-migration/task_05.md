---
task: task_05
spec: 0002-acpx-migration
status: pending
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

- [ ] Delete the SDK runner and SDK-only handlers
- [ ] Module dependency removal via the module tools
- [ ] Parity sweep: full suite green with unchanged assertions
- [ ] Gated real-acpx integration test (round-trip, cancel, crash-resume)
- [ ] Orphan-process check in the integration test teardown

## Acceptance Criteria

- [ ] `sh -c "! grep -q acp-go-sdk go.mod"` — the module no longer lists the SDK; the agent layer's owned line count shrinks (state the before/after in the Result).
- [ ] The full verification gate passes with zero behavioral test edits in this task.
- [ ] The gated integration test passes on a machine with Node and the pinned acpx (documented skip otherwise), covering round-trip, cancel, and crash-resume.
- [ ] No references to the SDK import path remain anywhere in the repository.

## Verification

- `rtk go test ./...` — expected: full suite passes.
- `sh -c "! grep -q acp-go-sdk go.mod"` — expected: exits 0 (dependency gone).
- `make verify` — expected: full gate passes.

## References

`_prd.md` → User Stories 2, 3; Goals (layer deleted); Success Metrics. `_techspec.md` → System Architecture (deleted), Testing Approach (gated integration test), Build Order 5, Risks. ADR-0017.
