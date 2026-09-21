---
task: task_07
spec: 0152-one-declared-acceptance-policy
status: completed
type: backend
complexity: medium
---

# Task 07: A rendered command that needs nothing installed

## Overview

Task 03 removed the awk copy of the rule by having the rendered command invoke
`roundfix qa-report accept`. Pre-PR review raised that to a P1, and it is right:
the command then runs whichever `roundfix` is first on `PATH`.

Measured on this machine, a Homebrew 0.14.1 shadowed the current build and has
no `qa-report` command at all, so a valid report was rejected before the policy
ran. `.agents/skills/write-tasks/SKILL.md` forbids exactly this — Verification
must be satisfiable "in a fresh worktree using only repository state", and an
installed binary is ambient machine state.

Weakening the rendered command alone would open a hole: `settle` replays a
Task's Verification and does not special-case a `qa` Task, so a weaker command
would let it complete a gate whose report is ineligible.

So both settlement paths apply the decision in process, and the rendered command
stops deciding.

## Requirements

1. MUST render a command that invokes no installed binary and depends on no
   ambient machine state.
2. MUST have the rendered command still prove, from repository state alone, that
   a newest QA Report exists and that its verdict is readable, failing when
   either does not hold.
3. MUST apply the one eligibility decision in `settle` before it completes a
   `qa` Task, in process.
4. MUST keep the Daemon's gate settlement applying the decision exactly as Task
   06 left it.
5. MUST keep `settle` unchanged for every non-`qa` Task.
6. MUST NOT reintroduce a copy of the eligibility rule in awk or any other
   rendered form.

## Subtasks

- [ ] Render a hermetic command that proves presence and readability only.
- [ ] Apply the decision in `settle` for a `qa` Task.
- [ ] Add tests for the ineligible-report refusal through `settle`.

## Acceptance Criteria

- [ ] The rendered command contains no invocation of `roundfix`.
- [ ] It fails when no report exists and when the newest report's verdict cannot
      be read.
- [ ] `settle` refuses to complete a `qa` Task whose newest report is ineligible,
      and completes one whose report qualifies.
- [ ] A non-`qa` Task settles exactly as before.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/spec/task.go`
- interface: `internal/cli/settle.go`

## Verification

- `test "$(grep -cF 'roundfix qa-report accept' internal/spec/task.go)" = "0"` — expected: exit 0; the rendered command invokes no installed binary. Before this Task it does, so the command fails.
- `out="$(go test -count=1 -v -run "^TestSettleAppliesEligibilityToAQATask" ./internal/cli 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestSettleAppliesEligibilityToAQATask"` — expected: exit 0; before this Task the case does not exist, so the command fails.

## References

- [_techspec.md](_techspec.md) — How the derived command reaches it

## Result

Implementation:

- The rendered QA Verification still selects the newest dated-and-sequenced
  report, but now checks only that the report has closed frontmatter and one
  non-empty verdict. It does not invoke `roundfix` or encode report
  eligibility.
- `settle` now reads the newest QA Report and applies
  `spec.QAReportEligibility` in process after the Task's Verification passes
  and before any Task status or commit change. Non-`qa` Tasks bypass this
  check.
- The Daemon gate settlement path was left unchanged from Task 06.

Focused-check evidence:

- Red signal: `rtk env GOCACHE=/private/tmp/roundfix-task07-gocache go test
  -count=1 ./internal/cli -run
  'TestSettleAppliesEligibilityToAQATask/refuses_an_ineligible_report_after_Verification'`
  failed before implementation because the rendered command still ran
  `roundfix qa-report accept` and never reached an in-process eligibility
  decision.
- `rtk env GOCACHE=/private/tmp/roundfix-task07-gocache go test -count=1
  ./internal/spec -run
  'TestDerivedQAVerification(ProvesNewestVerdictReadable|InvokesNoRoundfixBinary)'`
  passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task07-gocache go test -count=1
  ./internal/cli -run 'TestSettleAppliesEligibilityToAQATask/'` passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task07-gocache go test -count=1
  ./internal/spec ./internal/cli` passed `internal/spec`; the sandboxed CLI
  package run reached two unrelated force-stop integration failures because
  process-table reads returned `operation not permitted`.
- The host-permitted rerun, `rtk env
  GOCACHE=/private/tmp/roundfix-task07-gocache go test -count=1
  ./internal/cli`, passed.
- `rtk env GOCACHE=/private/tmp/roundfix-task07-gocache make
  verify-incremental` passed with host process-table access, covering `go vet`,
  the complete Go suite, skill checks, and the build.

Acceptance evidence:

- `TestDerivedQAVerificationInvokesNoRoundfixBinary` asserts the rendered text
  contains no `roundfix` invocation and executes it with a `PATH` containing
  only the report-selection utilities.
- `TestDerivedQAVerificationProvesNewestVerdictReadable` covers missing,
  malformed, duplicate, empty, and body-only verdicts and proves a readable
  `pass`, `fail`, or `partial` verdict satisfies the rendered command.
- `TestSettleAppliesEligibilityToAQATask` proves an ineligible newest report
  leaves the `qa` Task, worktree, and `HEAD` unchanged after Verification, and
  proves a qualifying `partial` settles the Task `completed`.
- The same settlement regression proves a backend Task retains its authored
  Verification and settles through the existing non-`qa` path.

Not run:

- The Task's declared `## Verification` commands — reserved for the Daemon by
  the assigned execution contract.
