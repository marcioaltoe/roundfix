---
task: task_02
spec: 0029-launch-and-recovery-fixes
status: pending
type: backend
complexity: medium
---

# Task 02: Classify Agent Model not-advertised rejections from acpx stderr

## Overview

Turn the opaque `agent/protocol error` into a typed, actionable failure: when acpx rejects an Agent Model with its "did not advertise that model" stderr, the agent layer recognizes it and produces a selection error carrying the runtime, the rejected model, the parsed advertised-model list, and the same recovery guidance the Preflight rejection already renders. Parsing is best-effort by design — anything unrecognized keeps today's infrastructure error, so this change can only improve messages.

## Requirements

1. MUST add a typed error recognized from the acpx stderr tail containing "did not advertise that model", carrying runtime, rejected model, and the advertised list parsed best-effort from the "Available models:" segment (possibly empty).
2. MUST render the same recovery guidance the existing preflight selection error uses: update the runtime or adapter, choose an advertised Agent Model, or pass a one-Run model override.
3. MUST fall back to the existing infrastructure error untouched when the stderr does not match — no rejection may become less informative than today.
4. MUST NOT change any invocation shape toward acpx; this task is classification only, with wiring into Batch settlement and Doctor handled by later tasks.

## Subtasks

- [ ] Typed error with runtime, model, advertised list, and recovery rendering
- [ ] Best-effort stderr-tail parser for the rejection message and the advertised list
- [ ] Classification hook where acpx failures are currently wrapped
- [ ] Unit tests: exact message, wrapped/multi-line tail, list parsing, unparseable garbage falling back to the existing error

## Acceptance Criteria

- [ ] A fixture stderr tail with the rejection yields the typed error with the correct model and advertised list, and its message names the recovery paths
- [ ] A garbage stderr tail yields exactly the current infrastructure error
- [ ] `errors.As` extracts the typed error through the existing wrap chain
- [ ] The full test suite passes

## Context

- interface: `internal/agent/acpx_runner.go`

## Verification

- `grep -q "ModelNotAdvertisedError" internal/agent/acpx_runner.go` — expected: exit 0 (typed error exists)
- `go test ./internal/agent/...` — expected: all tests pass, including the classification fixtures
- `go build -buildvcs=false ./...` — expected: clean build

## References

`_prd.md` → Goal 3, Core Feature 2, Problem 2; `_techspec.md` → Build Order 2, Interfaces (ModelNotAdvertisedError), Integration Points (acpx), Risks (stderr parsing); ADR-0037, ADR-0039, ADR-0041.
