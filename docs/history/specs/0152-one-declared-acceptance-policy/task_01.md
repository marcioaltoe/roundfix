---
task: task_01
spec: 0152-one-declared-acceptance-policy
status: completed
type: backend
complexity: medium
---

# Task 01: The one eligibility decision

## Overview

Whether the newest QA Report is acceptable is decided in two places that
disagree. This Task creates the one decision both will use. It states exactly
the rule `internal/spec/archive.go` states today, so nothing that is acceptable
now stops being acceptable.

## Requirements

1. MUST provide one exported decision that takes a parsed QA Report and answers
   whether it is acceptable.
2. MUST return the reason when it is not, so a caller can report what archive
   reports today.
3. MUST accept `pass` outright, **including one carrying environment-blocked
   rows**. `readQAReport` has already refused a `pass` with finding-, declared-
   or precondition-blocked rows, and ADR-0080 keeps an environment-blocked row
   distinct from a failure. Every QA Report in this repository carries one, so a
   decision that refused it would archive nothing ever again.
4. MUST accept `partial` only when it carries no finding-blocked rows, no
   environment-blocked rows, at least one declared blocked row, and a declared
   count no greater than the Spec's unreachable acceptance declarations.
5. MUST refuse every other verdict, and refuse a `partial` failing any clause
   above, with the reason `archiveUnprovenActions` gives today.
6. MUST take the Spec directory alongside the report, because counting
   unreachable declarations needs it.
7. MUST NOT change which report is newest, or what makes a blocked row
   unreachable.

## Subtasks

- [x] Add the decision and its reasons.
- [x] Add unit tests for every accepted and refused shape.

## Acceptance Criteria

- [x] A `pass` is accepted, including one carrying an environment-blocked row.
- [x] A qualifying `partial` is accepted; each disqualifying clause refuses.
- [x] Each refusal names the reason `archiveUnprovenActions` gives today.
- [x] Newest-report selection is untouched.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/spec/qa.go`
- interface: `internal/spec/archive.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestQAReportEligibility" ./internal/spec 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestQAReportEligibility"` — expected: exit 0; before this Task the case does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — The one decision

## Result

Implementation:

- Added exported `QAReportEligibility`, which returns `nil` for an acceptable
  parsed report and preserves archive's current refusal reasons as errors.
- Added a table-driven unit suite for plain and environment-blocked `pass`,
  qualifying `partial` reports, every disqualifying `partial` clause, `fail`,
  and an unsupported verdict.

Focused-check evidence:

- Red signal: `rtk env GOCACHE=/tmp/roundfix-go-cache go test
  ./internal/spec -run '^TestQAReportEligibility$'` failed to compile before
  implementation because `QAReportEligibility` was undefined.
- `rtk env GOCACHE=/tmp/roundfix-go-cache go test ./internal/spec -run
  '^(TestQAReportEligibility|TestNewestQAReport)'` passed after implementation.
- `rtk env GOCACHE=/tmp/roundfix-go-cache make verify-incremental` passed with
  host process-table access. The first sandboxed run reached and passed
  `internal/spec`, then failed two unrelated `internal/cli` force-stop tests
  because process-table access was denied; the permitted rerun exited zero.

Acceptance evidence:

- `TestQAReportEligibility/pass` and `pass_with_environment-blocked_rows`
  exercise both required `pass` shapes.
- The qualifying `partial` cases exercise equal and lower declared-row counts;
  the refusal cases exercise finding-blocked, environment-blocked, zero
  declared, and over-declared reports.
- Every refusal case compares the complete error string with the reason
  `archiveUnprovenActions` currently returns.
- The production diff adds only `QAReportEligibility` below `ReadQAReport`;
  `NewestQAReport` and `NewestQAReportFromPaths` are unchanged, and their
  focused tests pass.

The Daemon-owned command under `## Verification` was not run.
