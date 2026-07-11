---
task: task_02
spec: 0026-model-fallback-guardrail
status: completed
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

- [x] Hook the fallback probe into the shared selection preflight failure
      path with catalog and effort candidates.
- [x] Compose the deterministic fallback report and the explicit-flags
      re-run line per command.
- [x] Handle the no-candidate outcome by extending the existing failure
      text.
- [x] Cover no-input and detach failures, report shape, re-run line
      correctness, and unchanged exit codes with CLI tests.

## Acceptance Criteria

- [x] A no-input or detached command with a failing configured selection
      exits with the existing Preflight Validation code, prints the fallback
      report on stderr, and creates no Run.
- [x] The printed re-run line names the same command with explicit model and
      reasoning-effort flags, and executing that selection passes preflight
      in the test fake.
- [x] A model-managed fallback renders its effort as model-managed in the
      report and as an explicit empty reasoning-effort flag in the re-run
      line.
- [x] With no functional candidate, the failure preserves the existing
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

## Result

The resolve, watch, and implement selection preflights now probe the failed
ACP Runtime's Model Catalog and highest-first reasoning vocabulary. A
non-interactive failure reports the failed selection and proven Fallback
Selection on stderr, then prints one copy-paste re-run with explicit
`--model` and `--reasoning-effort` flags. The path preserves exit code 2,
creates no Run, starts no Agent work, and leaves stdout empty. Doctor and
Setup Command behavior remains unchanged.

Verification:

- `rtk go test ./internal/cli`: passed — 377 tests.
- `rtk make verify`: passed — 1,062 tests across 19 packages, Roundfix skill
  checks, and the CLI build.

Acceptance evidence:

1. Buffer-captured resolve, watch, implement, simulated Detached Run child,
   and non-TTY stderr tests returned exit 2, kept stdout empty, recorded zero
   Agent calls, and found no created Run.
2. The resolve test executed the printed explicit fallback selection through
   the fake preflight and reached a clean Run; the recorded probe used the
   reported model and reasoning effort.
3. The watch test rendered an empty fallback effort as `model-managed` and
   printed `--reasoning-effort ""` in the re-run command.
4. The no-candidate test retained the original actionable
   `SelectionPreflightError` and appended the attempted Model Catalog entries
   plus the highest-first reasoning vocabulary.

Follow-ups: none.
