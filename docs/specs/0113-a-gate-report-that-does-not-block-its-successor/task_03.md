---
task: task_03
spec: 0113-a-gate-report-that-does-not-block-its-successor
status: completed # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: low
---

# Task 03: Name the literal the blocked cause requires

## Overview

A row typed `blocked (finding: …)` counts only when its text also carries the
exact string `" — waits on "`. A row missing it is refused with a message listing
the three type prefixes the row already satisfies, never the literal actually
wanted. An author reading that rewrites a row that was typed correctly.

## Requirements

1. MUST name, and quote, the literal a typed blocked status must carry, when a row
   carries a valid type without it.
2. MUST keep refusing a row whose blocked cause is not one of the three types,
   with the message that names the three.
3. MUST distinguish the two cases in the message, so a reader can tell a wrong
   type from a missing literal.
4. MUST NOT change which rows count toward which typed total.

## Subtasks

- [ ] Separate the missing-literal case from the wrong-type case.
- [ ] Quote the required literal in the diagnostic.
- [ ] Cover both cases and each of the three types.

## Acceptance Criteria

- [ ] A finding-typed row without the literal is refused with a message quoting
      the literal.
- [ ] A row with an unrecognised blocked cause is still refused with the message
      naming the three types.
- [ ] The two messages differ.
- [ ] Each of the three typed causes still counts as it does today.

## Verification

- `go test -count=1 ./internal/speccheck -run 'TestBlockedCauseDiagnosticNamesTheLiteral' -v > /tmp/0113-t03.log 2>&1; s=$?; grep -q '^--- PASS: TestBlockedCauseDiagnosticNamesTheLiteral' /tmp/0113-t03.log || { cat /tmp/0113-t03.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today, where no such test exists.
- `test -s /tmp/0113-t03.log || { echo 'the suite produced no output'; exit 1; }; grep -q '^--- PASS: TestBlockedCauseDiagnosticNamesTheLiteral' /tmp/0113-t03.log || { echo 'the diagnostic test did not run'; exit 1; }; grep -c '^    --- PASS\|^--- PASS' /tmp/0113-t03.log > /tmp/0113-t03-n.txt; test "$(cat /tmp/0113-t03-n.txt)" -ge 5 || { echo "expected the missing-literal case, the wrong-type case, and the three typed causes, got $(cat /tmp/0113-t03-n.txt)"; cat /tmp/0113-t03.log; exit 1; }` — expected: exits 0, refusing a vacuous run and proving each case runs on its own.
- `grep -rq 'func TestBlockedCauseDiagnosticNamesTheLiteral' internal/speccheck || { echo 'the diagnostic case does not exist'; exit 1; }; grep -qE '(Detail|Fix):.*waits on' internal/speccheck/mechanical.go || { echo 'the literal appears only in the parser condition, never in a diagnostic'; exit 1; }; grep -n 'blocked cause outside environment, finding, or declared' internal/speccheck/mechanical.go > /tmp/0113-t03-msg.txt; test -s /tmp/0113-t03-msg.txt || { echo 'the wrong-type message was removed rather than kept beside the new one'; exit 1; }; exit 0` — expected: exits 0, proving the literal reaches a diagnostic rather than only the parser's condition, and that the wrong-type message survives beside it. Fails today: the detector already contains the literal at its comparison, so the check is on where it appears, not whether.

## Context

- interface: `internal/speccheck/mechanical.go`

## References

`_techspec.md` → Build Order 3; Implementation Design, Interfaces. `_prd.md` →
Core Features 3 and 5; Goal 2; User Story 2. ADR-0133.

## Result

### Implementation

- The report-shape detector now separates a finding-typed blocked status that
  lacks the required literal from a blocked status with an unrecognised cause.
- The missing-literal diagnostic quotes `" — waits on "`; the existing wrong-type
  diagnostic and its three named cause types remain unchanged.
- The environment, finding and declared count increments remain on their existing
  accepted-status branches.

### Focused-check evidence

- Before the production change,
  `GOCACHE=/tmp/roundfix-task-03-gocache rtk proxy go test ./internal/speccheck -run '^TestBlockedCauseDiagnosticNamesTheLiteral$'`
  failed because the missing-literal row received the wrong-type diagnostic.
- After the production change, the same focused command passed.
- `GOCACHE=/tmp/roundfix-task-03-gocache rtk make verify-incremental` passed when
  run outside the managed sandbox. The sandboxed attempt could not let two
  unrelated `internal/cli` process-tree tests read the process table; all other
  reported packages, including `internal/speccheck`, passed in that attempt.
- `rtk git diff --check` passed.

### Acceptance-criterion evidence

- Finding-typed row without the literal: subtest
  `finding_type_without_required_literal_names_the_literal` requires the quoted
  literal in the diagnostic and passed.
- Unrecognised blocked cause: subtest
  `unrecognised_blocked_cause_names_the_three_types` requires the existing three-
  type diagnostic and passed.
- Distinct messages: both refusal subtests reject the other case's diagnostic;
  both passed.
- Typed totals: the environment, finding and declared subtests each provide the
  matching frontmatter count and require zero report-shape findings; all passed.

### Daemon verification

The commands under `## Verification` were not run; the Daemon owns those commands
and Task settlement.
