---
task: task_02
spec: 0026-model-fallback-guardrail
status: pending
type: backend
complexity: medium
---

# Task 02: Report the fallback on non-interactive selection failures

## Overview

Wire the fallback probe into the commands' selection preflight for the
non-interactive contract: when Agent selection fails and no prompt is
possible, the command fails Preflight Validation with a deterministic stderr
report naming the failed selection, the proven Fallback Selection, and one
copy-paste re-run line with explicit model and reasoning-effort flags. The
slice is verifiable through buffer-captured CLI runs.

## Requirements

1. MUST run the fallback probe when Agent selection fails Preflight
   Validation on the resolve, watch, and implement commands, using the Model
   Catalog for candidates and the runtime's reasoning vocabulary
   highest-first for efforts.
2. MUST fail Preflight Validation in non-interactive contexts (no-input,
   detach, or non-interactive stderr) with a report naming the failed
   selection, the proven Fallback Selection (rendering an empty effort as
   model-managed), and one exact re-run command line carrying explicit model
   and reasoning-effort flags for the same command.
3. MUST preserve today's Preflight Validation failure, extended with the
   probed candidates, when no fallback proves functional.
4. MUST keep exit codes unchanged, keep stdout free of the report, and start
   no Run and no Agent work from this path.
5. MUST keep the Doctor and Setup Commands prompt-free: their selection
   failure output may name the proven fallback but their report contracts
   and exit codes stay as they are.
6. MUST NOT add any flag or configuration key that pre-authorizes the
   fallback.

## Subtasks

- [ ] Hook the fallback probe into the shared selection preflight failure
      path with catalog and effort candidates.
- [ ] Compose the deterministic fallback report and the explicit-flags
      re-run line per command.
- [ ] Handle the no-candidate outcome by extending the existing failure
      text.
- [ ] Cover no-input and detach failures, report shape, re-run line
      correctness, and unchanged exit codes with CLI tests.

## Acceptance Criteria

- [ ] A no-input or detached command with a failing configured selection
      exits with the existing Preflight Validation code, prints the fallback
      report on stderr, and creates no Run.
- [ ] The printed re-run line names the same command with explicit model and
      reasoning-effort flags, and executing that selection passes preflight
      in the test fake.
- [ ] A model-managed fallback renders its effort as model-managed in the
      report and as an explicit empty reasoning-effort flag in the re-run
      line.
- [ ] With no functional candidate, the failure preserves the existing
      actionable selection error extended with the probed candidates.

## Verification

- `rtk go test ./internal/cli` - expected: non-interactive fallback report,
  re-run line, no-candidate, and exit-code tests pass.
- `rtk make verify` - expected: formatting, all tests, Skill sync checks,
  Skill checks, and build pass.

## References

- `_prd.md` → User Story 2; Core Features 1, 3, 6.
- `_techspec.md` → API Contracts; Build Order 2.
- ADR-0041.
