---
task: task_03
spec: 0152-one-declared-acceptance-policy
status: completed
type: backend
complexity: medium
---

# Task 03: The derived command delegates

## Overview

`DerivedQAVerification` renders an awk program whose tail reads
`exit(closed && verdicts == 1 && verdict == "pass" ? 0 : 1)`. It cannot see a
blocked-row declaration, so a qualifying `partial` fails it, and the Daemon runs
this command to settle the `qa` Task. The stricter copy wins by running first.

This Task makes the rendered command delegate the judgement instead of carrying
a second copy of the rule.

## Requirements

1. MUST keep finding the newest report by name in the rendered command, so a
   reader can still see which report is judged.
2. MUST delegate the verdict judgement to the one decision rather than
   evaluating it in awk.
3. MUST settle a newest report whose verdict is `partial` with
   declared-unreachable blocked rows.
4. MUST still refuse `fail`, a missing report, an unparseable report, and a
   `partial` that carries finding- or environment-blocked rows, declares none,
   or declares more than the Spec does.
5. MUST keep accepting a `pass` carrying an environment-blocked row, which every
   Spec in this repository depends on.
6. MUST fail closed: a judgement that cannot be reached is a refusal, never an
   acceptance.

## Subtasks

- [ ] Replace the awk verdict tail with a delegated judgement.
- [ ] Add tests for the qualifying partial and every refusal.
- [ ] Cover the missing and unparseable report cases explicitly.

## Acceptance Criteria

- [ ] A qualifying `partial` is accepted by the derived command.
- [ ] `fail`, a missing report, an unparseable report and each disqualifying
      `partial` shape are refused.
- [ ] A `pass` carrying an environment-blocked row is accepted.
- [ ] The rendered command still names the report directory it searches.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/spec/task.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestDerivedQAVerification" ./internal/spec 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestDerivedQAVerificationSettlesAQualifyingPartial"` — expected: exit 0; before this Task the case does not exist, so the command fails.
- `test "$(grep -cF 'verdicts == 1 && verdict ==' internal/spec/task.go)" = "0"` — expected: exit 0; the rendered command no longer carries its own copy of the verdict rule. Before this Task line 117 reads ``END { exit(closed && verdicts == 1 && verdict == \"pass\" ? 0 : 1) }``, so the command fails. Matched with `-F` because the literal contains backslash-escaped quotes that a regex would have to spell exactly.

## References

- [_techspec.md](_techspec.md) — How the derived command reaches it

## Result

Implementation:

- The rendered QA Verification still selects the newest dated-and-sequenced
  `qa-report-*.md` path in shell, then passes that exact path to
  `roundfix qa-report accept` instead of parsing or judging its verdict in awk.
- The read-only `qa-report accept` command parses the selected report and calls
  `QAReportEligibility` with its Spec directory. It emits no stdout and exits
  nonzero when the path, report, or shared eligibility decision refuses.
- `ReadQAReportFile` exposes the existing single-report parser so the command
  can judge the report selected by the rendered command without selecting a
  second report or copying parsing rules.

Focused-check evidence:

- Red signal: `rtk env GOCACHE=/tmp/roundfix-task03-go-cache go test -count=1
  ./internal/spec -run '^TestDerivedQAVerificationSettlesAQualifyingPartial$'`
  failed before implementation because the rendered awk still rejected
  `partial`.
- `rtk env GOCACHE=/tmp/roundfix-task03-go-cache go test -count=1
  ./internal/spec -run
  '^(TestDerivedQAVerificationSettlesAQualifyingPartial|TestDerivedQAVerificationDelegatesEligibility|TestDerivedQAVerificationFailsClosedWhenDelegateCannotRun|TestDerivedQAVerificationQuotesTheSpecPath|TestDerivedQAVerificationPassesTheChecker|TestReloadTaskDerivesOnlyQAVerification)$'`
  passed against a built Roundfix binary.
- `rtk env GOCACHE=/tmp/roundfix-task03-go-cache go test -count=1
  ./internal/cli -run '^TestRunQAReportAcceptCommand'` passed for the command's
  accepted and fail-closed paths.
- `rtk env GOCACHE=/tmp/roundfix-task03-go-cache make verify-incremental`
  passed with host process-table access. The first sandboxed run passed the
  changed packages and failed only the two force-stop integration tests whose
  process-table read was denied; the permitted rerun exited zero.
- `rtk git diff --check` passed.

Acceptance evidence:

- `TestDerivedQAVerificationSettlesAQualifyingPartial` accepts a newest
  `partial` report with one declared blocked row and one matching unreachable
  acceptance declaration.
- `TestDerivedQAVerificationDelegatesEligibility` refuses `fail`, a missing
  report, an unparseable report, finding- and environment-blocked partials, a
  partial with no declared rows, and a partial declaring more rows than the
  Spec. The same suite accepts a `pass` carrying an environment-blocked row.
- The derived-command fixture asserts that the rendered command names
  `docs/specs/<slug>/qa`; the existing path-quoting suite and checker contract
  also pass.
- Missing or malformed selection exits before delegation, while command-level
  missing, unparseable, and ineligible inputs all exit nonzero. A missing
  `roundfix` delegate also leaves the shell command nonzero, so an unreachable
  judgement cannot become acceptance.

Not run:

- The Task's declared `## Verification` commands — reserved for the Daemon by
  the assigned execution contract.
