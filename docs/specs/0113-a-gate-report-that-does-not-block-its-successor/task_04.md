---
task: task_04
spec: 0113-a-gate-report-that-does-not-block-its-successor
status: completed # pending | in_progress | completed | failed — only implement-task changes this
type: backend
complexity: medium
---

# Task 04: Report one cause once

## Overview

A row that fails to parse is dropped from the typed totals, so the declared count
disagrees with the table and a second finding reports a counting defect. There is
one cause and two findings, and the second sends a reader hunting arithmetic that
is fine.

## Requirements

1. MUST report a count disagreement caused by unparsed rows as that parse failure,
   not as a separate counting finding.
2. MUST still report a count disagreement when every row parsed and the declared
   totals are genuinely wrong.
3. MUST name the rows that failed to parse when it reports the parse failure.
4. MUST NOT suppress either finding when both causes are present at once.

## Subtasks

- [ ] Attribute a count disagreement to unparsed rows when they explain it.
- [ ] Keep the genuine count finding.
- [ ] Cover both alone and both together.

## Acceptance Criteria

- [ ] A report whose count disagrees only because a row failed to parse yields one
      finding, naming that row.
- [ ] A report whose rows all parse and whose totals are wrong still yields the
      count finding.
- [ ] A report with both an unparsed row and a genuinely wrong total yields both.

## Verification

- `go test -count=1 ./internal/speccheck -run 'TestCountDisagreementReportsItsCause' -v > /tmp/0113-t04.log 2>&1; s=$?; grep -q '^--- PASS: TestCountDisagreementReportsItsCause' /tmp/0113-t04.log || { cat /tmp/0113-t04.log; exit 1; }; exit $s` — expected: exits 0 and the log names the passing test; fails today, where no such test exists.
- `test -s /tmp/0113-t04.log || { echo 'the suite produced no output'; exit 1; }; grep -q '^--- PASS: TestCountDisagreementReportsItsCause' /tmp/0113-t04.log || { echo 'the attribution test did not run'; exit 1; }; grep -c '^    --- PASS\|^--- PASS' /tmp/0113-t04.log > /tmp/0113-t04-n.txt; test "$(cat /tmp/0113-t04-n.txt)" -ge 3 || { echo "expected the parse-caused, genuine, and both-at-once cases, got $(cat /tmp/0113-t04-n.txt)"; cat /tmp/0113-t04.log; exit 1; }` — expected: exits 0, refusing a vacuous run and proving all three combinations run.
- `go test -count=1 ./internal/speccheck ./internal/daemon > /tmp/0113-t04-regress.log 2>&1; s=$?; grep -q 'FAIL' /tmp/0113-t04-regress.log && { echo 'the detector or its caller regressed:'; grep -B 3 -A 8 'FAIL' /tmp/0113-t04-regress.log | head -30; exit 1; }; grep -rq 'TestCountDisagreementReportsItsCause' internal/speccheck || { echo 'both packages pass, but the attribution case does not exist'; exit 1; }; exit $s` — expected: exits 0, proving the detector and the Daemon that calls it agree, anchored to the case this Task adds so a green pair cannot pass on an untouched tree.

## Context

- interface: `internal/speccheck/mechanical.go`

## References

`_techspec.md` → Build Order 4; Testing Approach, one cause one finding.
`_prd.md` → Core Feature 4; User Story 3. ADR-0133.

## Result

### Implementation

- The report-shape detector now tracks finding-typed rows that fail the required
  literal check separately from rows that contribute to parsed totals.
- A declared finding count equal to the parsed total plus those unparsed rows is
  attributed to the row-specific parse finding, so it produces no second count
  finding.
- A declared finding count that still differs after accounting for unparsed rows
  keeps the count finding alongside every row-specific parse finding.
- The change introduces no domain term and requires no glossary update.

### Focused-check evidence

- Before the production change,
  `GOCACHE=/tmp/roundfix-task-04-gocache rtk proxy go test ./internal/speccheck -run '^TestCountDisagreementReportsItsCause$'`
  failed because `unparsed_row_accounts_for_count_disagreement` returned both
  the `R-PARSE` literal finding and a `rows_blocked_finding` count finding.
- After the production change, the same focused command passed.
- `GOCACHE=/tmp/roundfix-task-04-gocache rtk proxy go test ./internal/speccheck -run '^(TestMechanicalReportShape|TestBlockedCauseDiagnosticNamesTheLiteral|TestCountDisagreementReportsItsCause|TestMechanicalFindingsWithoutRowHintsBlockTheirRefusalCode)$'`
  passed.
- `GOCACHE=/tmp/roundfix-task-04-gocache rtk make verify-incremental` passed
  outside the managed sandbox. The sandboxed attempt could not let two unrelated
  `internal/cli` process-tree tests read the process table; `internal/speccheck`
  and `internal/daemon` passed in that attempt.
- `rtk git diff --check` passed after the implementation edits.

### Acceptance-criterion evidence

- Parse-caused disagreement: subtest
  `unparsed_row_accounts_for_count_disagreement` requires exactly one report-shape
  finding, requires that finding to name `R-PARSE`, and rejects a count finding;
  it passed.
- Genuine disagreement: subtest
  `parsed_rows_expose_genuine_count_disagreement` uses only a parsed row and
  requires the `rows_blocked_finding` count finding; it passed.
- Both causes: subtest `unparsed_row_and_wrong_total_expose_both_causes` requires
  two findings, the parse failure naming `R-PARSE`, and the count finding; it
  passed.

### Daemon verification

The commands under `## Verification` were not run; the Daemon owns those commands
and Task settlement.
