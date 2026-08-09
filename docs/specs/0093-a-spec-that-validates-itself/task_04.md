---
task: task_04
spec: 0093-a-spec-that-validates-itself
status: pending
type: backend
complexity: low
---

# Task 04: Surface the stage through the command

## Overview

The scope exists in the package; an authoring skill reaches it through the
command. This Task adds the flag, its help text, and the exit behaviour a caller
depends on to block a stage while a finding stands.

## Requirements

1. MUST accept `--stage prd|techspec|tasks` on `roundfix spec check`.
2. MUST exit non-zero when a scoped run reports any error-level finding, so an
   authoring step can block on it.
3. MUST leave the command without the flag behaving exactly as today, including
   its exit status and output shape.
4. MUST reject an unknown stage value with a message naming the accepted ones.
5. MUST document the flag in the command's help.

## Subtasks

- [ ] Add and parse the flag.
- [ ] Wire the exit status.
- [ ] Document it in help.

## Acceptance Criteria

- [ ] A scoped run with a finding exits non-zero.
- [ ] A scoped run with no finding exits zero.
- [ ] The unscoped command is unchanged.
- [ ] An unknown stage is rejected, naming the accepted values.

## Bounded scope

This Task may create or modify only:

- `internal/cli/spec_check.go`
- `internal/cli/spec_check_test.go`
- `docs/specs/0093-a-spec-that-validates-itself/task_04.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^TestSpecCheckStage' -count=1 -v 2>&1 | grep -q '^--- PASS: TestSpecCheckStageExitsNonZeroOnAFinding'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^TestSpecCheckStage' -count=1 -v 2>&1 | grep -q '^--- PASS: TestSpecCheckStageRejectsAnUnknownValue'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^TestSpecCheckStage' -count=1 -v 2>&1 | grep -q '^--- PASS: TestSpecCheckWithoutStageIsUnchanged'` — expected: exits 0.
- `go run -buildvcs=false ./cmd/roundfix spec check --help 2>&1 | grep -q -- '--stage'` — expected: exits 0, proving the flag reached the help a maintainer reads.

## References

- `_prd.md` → Goal 2.
- `_techspec.md` → Build Order 4; API Contracts.
