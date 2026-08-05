---
task: task_04
spec: 0067-derived-artifact-regeneration-boundary
status: pending
type: backend
complexity: low
---

# Task 04: Say what a human must do when the command cannot

## Overview

`make baseline-digests` prints "Read the failing test output above, fix the
canonical source it validates, then rerun make baseline-digests" — advice that
cannot be followed when the canonical source is already correct and only a
non-sanctioned corpus is stale. This slice makes the diagnostic read the
ownership record and name the actual human action.

## Requirements

1. MUST name, when a failure's remediation lies outside the sanctioned command,
   the exact invocation from the owning record rather than printing a command
   that will not help.
2. MUST name the owning record's path, so the reader can verify the claim.
3. MUST keep the existing diagnostic for failures the sanctioned command *can*
   fix, unchanged.
4. MUST state plainly, for a `frozen` artifact, that nothing regenerates it and
   why, instead of suggesting a regeneration.

## Subtasks

- [ ] Read the ownership record when a failure falls outside the command.
- [ ] Emit the declared invocation and the record path.
- [ ] Keep the in-scope diagnostic unchanged.

## Acceptance Criteria

- [ ] A stale `dedicated` corpus produces a diagnostic naming its exact
      invocation and its record path.
- [ ] A failure the sanctioned command can fix keeps today's diagnostic
      verbatim.
- [ ] A `frozen` artifact's diagnostic states nothing regenerates it and gives
      the recorded reason.
- [ ] No diagnostic suggests a command that cannot resolve the failure it
      reports, asserted across all three owner classes.

## Verification

- `go build -buildvcs=false ./...` — expected: exit 0.
- `go test ./internal/baseline -count=1 -run 'Diagnostic|Remediation|Ownership' -v | grep -q -- "--- PASS"`
  — expected: exit 0; the diagnostic tests ran and passed.
- `go test ./internal/baseline -count=1` — expected: exit 0.
- `make verify` — expected: exit 0.

## References

- `_prd.md` → Core Feature 4; Goals.
- `_techspec.md` → API Contracts; Build Order 4.
