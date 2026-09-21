---
task: task_03
spec: 0152-one-declared-acceptance-policy
status: pending
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
4. MUST still refuse `fail`, an undeclared `partial`, a missing report, an
   unparseable report, and a `pass` carrying blocked rows.
5. MUST fail closed: a judgement that cannot be reached is a refusal, never an
   acceptance.

## Subtasks

- [ ] Replace the awk verdict tail with a delegated judgement.
- [ ] Add tests for the qualifying partial and every refusal.
- [ ] Cover the missing and unparseable report cases explicitly.

## Acceptance Criteria

- [ ] A qualifying `partial` is accepted by the derived command.
- [ ] `fail`, an undeclared `partial`, a missing report, an unparseable report
      and a `pass` with blocked rows are each refused.
- [ ] The rendered command still names the report directory it searches.

## Context

- instruction: `.agents/skills/implement-task/SKILL.md`
- interface: `internal/spec/task.go`

## Verification

- `out="$(go test -count=1 -v -run "^TestDerivedQAVerification" ./internal/spec 2>&1)" || { printf "%s\n" "$out"; exit 1; }; printf "%s\n" "$out" | grep -q -- "--- PASS: TestDerivedQAVerificationSettlesAQualifyingPartial"` — expected: exit 0; before this Task the case does not exist, so the command fails.
- `test "$(grep -cF 'verdicts == 1 && verdict ==' internal/spec/task.go)" = "0"` — expected: exit 0; the rendered command no longer carries its own copy of the verdict rule. Before this Task line 117 reads ``END { exit(closed && verdicts == 1 && verdict == \"pass\" ? 0 : 1) }``, so the command fails. Matched with `-F` because the literal contains backslash-escaped quotes that a regex would have to spell exactly.

## References

- [_techspec.md](_techspec.md) — How the derived command reaches it
