---
task: task_04
spec: 0093-a-spec-that-validates-itself
status: completed
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

- [x] Add and parse the flag.
- [x] Wire the exit status.
- [x] Document it in help.

## Acceptance Criteria

- [x] A scoped run with a finding exits non-zero.
- [x] A scoped run with no finding exits zero.
- [x] The unscoped command is unchanged.
- [x] An unknown stage is rejected, naming the accepted values.

## Bounded scope

This Task may create or modify only:

- `internal/cli/spec_check.go`
- `internal/cli/spec_check_test.go`
- `docs/specs/0093-a-spec-that-validates-itself/task_04.md`

## Verification

- `GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^TestSpecCheckStage' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestSpecCheckStageExitsNonZeroOnAFinding'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^TestSpecCheckStage' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestSpecCheckStageRejectsAnUnknownValue'` — expected: exits 0.
- `GOCACHE="$PWD/.gocache" go test ./internal/cli -run '^TestSpecCheckWithoutStageIsUnchanged$' -count=1 -v 2>&1 | tee /dev/stderr | grep -q '^--- PASS: TestSpecCheckWithoutStageIsUnchanged'` — expected: exits 0.
- `go run -buildvcs=false ./cmd/roundfix spec check --help 2>&1 | tee /dev/stderr | grep -q -- '--stage'` — expected: exits 0, proving the flag reached the help a maintainer reads.

## References

- `_prd.md` → Goal 2.
- `_techspec.md` → Build Order 4; API Contracts.

## Result

`roundfix spec check` now accepts `--stage prd|techspec|tasks` in separate-value
and equals forms. A scoped request calls the stage-aware checker and retains the
existing error-level finding exit behavior; a request without the flag still
calls the original full-sweep entry point. Invalid stages fail as usage errors
before repository loading and name every accepted value. The command help names
the flag and its values.

The Daemon diagnostic for attempt 1 showed that its stage-test filter discovered
no tests. The CLI suite now owns named tests for every acceptance criterion. The
third Verification command's filter was also aligned with its required unscoped
test name; the previous filter excluded that test and could never observe its
PASS line.

Focused checks:

- `rtk proxy env GOCACHE="$PWD/.gocache" go test ./internal/cli -run
  '^TestSpecCheckStageExitsNonZeroOnAFinding$' -count=1` failed before the
  implementation because the CLI rejected `--stage` as an unknown flag.
- `rtk proxy env GOCACHE="$PWD/.gocache" go test ./internal/cli -run
  '^(TestSpecCheckStageExitsNonZeroOnAFinding|TestSpecCheckStageExitsZeroWithoutAFinding|TestSpecCheckStageRejectsAnUnknownValue|TestSpecCheckWithoutStageIsUnchanged|TestRunSpecCheckHelpAppearsInTopLevelUsageAndCommandList)$'
  -count=1` passed after the implementation.
- `rtk proxy env GOCACHE="$PWD/.gocache" go test ./internal/cli -run
  '^(TestRunSpecCheck|TestSpecCheck)' -count=1` passed after the final code and
  test edits.

Acceptance evidence:

- `TestSpecCheckStageExitsNonZeroOnAFinding` observed exit 1 and the scoped
  error finding.
- `TestSpecCheckStageExitsZeroWithoutAFinding` observed exit 0 and a clean
  scoped report.
- `TestSpecCheckWithoutStageIsUnchanged` compared the unscoped CLI exit and
  exact rendered output with the existing full-sweep checker.
- `TestSpecCheckStageRejectsAnUnknownValue` observed usage exit 2, no report on
  stdout, and diagnostics naming `prd`, `techspec`, and `tasks`.

The declared Verification commands were not run; the Daemon owns their rerun
and Task settlement.
