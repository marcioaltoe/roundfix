---
task: task_03
spec: 0223-adjustments-the-adopters-asked-for
status: pending
type: backend
complexity: medium
---

# Task 03: roundfix release plan reports the skills and baseline checks

## Overview

Let a successful range release plan print one `skills:` line and one
`baseline:` line, and carry both under a JSON `checks` object, as the
TechSpec's "The release plan checks" states. The checks reuse the Doctor's
skills comparison and the Baseline update's managed-refresh planning, read
local files only, write nothing, and never change the decision. This answers
the Backlog Entry of 2026-09-30, "The release plan does not report the skill
and guide checks".

## Requirements

1. MUST extract the Doctor's skills check into `repositorySkillsCheck` in
   `internal/cli/doctor.go`, called by the Doctor, with the Doctor's output and
   tests unchanged.
2. MUST add `collectReleasePlanChecks` and the two printers' check output in
   the new `internal/cli/releaseplan_checks.go`, and call it from
   `internal/cli/releaseplan_command.go` only after `releaseplan.Build`
   succeeds; the reset plan, a refused plan and a failed plan print no check.
3. MUST produce the statuses, details and next action of "The release plan
   checks" steps 1-4 and the text and JSON shapes of API Contracts 1 and 2,
   keeping `schemaVersion` unchanged.
4. MUST NOT let a check change `state`, `proposedVersion` or the exit code
   (API Contract 3), write any file, run `baseline update`'s skills stage, or
   contact any service.
5. MUST add the help sentence of API Contract 5 to the `release plan` usage in
   `internal/cli/cli.go`.
6. MUST add the tests of the TechSpec's Testing Approach 3 in the new file
   `internal/cli/releaseplan_checks_test.go`, and MUST NOT edit any existing
   test.

## Subtasks

- [ ] Extract the Doctor's skills check.
- [ ] Compute the baseline check from the Setup Manifest and a managed-refresh plan.
- [ ] Print both checks in text and JSON and add the help sentence.
- [ ] Add the four tests.

## Acceptance Criteria

- [ ] A range plan prints exactly one `skills:` and one `baseline:` line after
      `Next action:`, and its JSON carries `checks`.
- [ ] A failing check leaves the state, proposed version and exit code as they
      were.
- [ ] The plan leaves the repository's bytes and Git status unchanged.
- [ ] The Doctor's output is unchanged.

## Context

- interface: `internal/cli/doctor.go`
- interface: `internal/cli/releaseplan_command.go`
- interface: `internal/cli/cli.go`
- creates: `internal/cli/releaseplan_checks.go`
- creates: `internal/cli/releaseplan_checks_test.go`
- instruction: `docs/user-guide/release-runbook.md`

## Verification

- `out="$(go test -count=1 -v -run "^(TestReleasePlanReportsTheSkillsAndBaselineChecks|TestReleasePlanChecksNeverChangeTheDecision|TestReleasePlanChecksOnAnAdoptedRepository|TestReleasePlanChecksWriteNothing)$" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; for name in TestReleasePlanReportsTheSkillsAndBaselineChecks TestReleasePlanChecksNeverChangeTheDecision TestReleasePlanChecksOnAnAdoptedRepository TestReleasePlanChecksWriteNothing; do printf "%s\n" "$out" | grep -q -- "--- PASS: $name" || { printf 'missing pass: %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0; before this Task the four new tests do not exist, so the command fails; after it they pass.

## References

- `_prd.md` → Core Feature 5; User Story 4; Success Metric 3
- `_techspec.md` → The release plan checks; API Contract 1; API Contract 2; API Contract 3; API Contract 5; Surface Transcript 1; Surface Transcript 2; Testing Approach 3; Build Order 3
- ADR-0184
