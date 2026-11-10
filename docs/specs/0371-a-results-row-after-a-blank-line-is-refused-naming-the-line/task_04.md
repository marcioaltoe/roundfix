---
task: task_04
spec: 0371-a-results-row-after-a-blank-line-is-refused-naming-the-line
status: pending
type: qa
complexity: medium
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes one reader decide
which lines of a QA Report are Results rows and refuses a table line that no
header leads, such as a row after a blank line, on every QA settlement path,
in `roundfix archive` and in the mechanical stage, naming the report path,
the row's line and the repair. The Daemon tells the gate to continue the
seeded table and returns a broken report once; the qa-gate skill, the guides,
the references and the glossary say the same.

Every row runs on this repository's own files, temporary directories,
temporary Git repositories and temporary copies of the audited head. No row
starts an Agent session or an ACP adapter, reads or writes the real
`~/.roundfix`, reads another repository, reaches the network, sends a request
to a provider, the judge or a reviewer, or reads a credential. Every `go test`
it runs carries an explicit `-timeout` and keeps the session's own `HOME` and
`TMPDIR`. Every row that runs `bin/roundfix` sets `HOME` and `TMPDIR` to fresh
temporary directories, with `GOMODCACHE` and `GOCACHE` kept at their `go env`
values, and removes them afterwards. The QA never writes into the checkout
while a test suite runs. No row compares one command's output with another's
byte for byte. This report itself continues the seeded Results table directly
below its last row, with no blank line before a row. The gate keeps its own
evidence under `qa/evidence/<report name>/`, under names the repository does
not ignore, `.txt` for command output, and links only files it kept.

## Requirements

1. MUST build the binary with `make build` and record the Daemon's repository
   Verification result as a gate fact.
2. MUST verify Success Metrics 1 to 7 and Invariants 1 to 5 by running
   `go test -count=1 -timeout 20m -v` with every test task_01 and task_02 name
   in their Verification over `./internal/spec`, `./internal/speccheck`,
   `./internal/cli` and `./internal/daemon`, then the whole `./internal/spec`,
   `./internal/speccheck` and `./skills` packages and
   `-tags docscontract ./internal/docscontract`, recording each pass.
3. MUST reproduce the defect's shape through the built binary in a temporary
   directory the row creates and removes, outside any Git repository: write
   `docs/specs/9001-fixture/_prd.md` with the `## Unreachable Acceptance`
   declaration and `docs/specs/9001-fixture/qa/qa-report-2026-11-10.md` with
   the reduced Pantheon shape, both exactly as `_techspec.md` → Testing
   Approach gives them. Run
   `bin/roundfix qa-report accept <absolute report path>` with fresh `HOME`
   and `TMPDIR` and record exit `1`, empty stdout, and the stderr line
   `roundfix: QA Report eligibility refused: QA Report <absolute report path>: Results row at line 16 is outside the Results table that starts at line 12; remove the blank line at <absolute report path>:15 so the row continues the table`.
   Then delete line 15 of the report, rerun, and record exit `0` with empty
   stdout and stderr; then remove the temporary homes and directory.
4. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that the
   named test fails on the case whose input makes the mutation observable:
   - `ReadQAResults` reads a table line after a blank line as a row of the
     table above: `TestQAResultsRefusesARowAfterABlankLine` fails on the
     reduced Pantheon shape;
   - `QAResultsBreak.Repair` names only the first of several blank lines:
     `TestQAResultsNamesEveryBlankLineOrTheRowsPlace` fails on the two-blank
     -line case;
   - `StatusTables` returns only the first table:
     `TestQAResultsReadsEveryStatusTable` fails on the second table;
   - the break check removed from `QAReportEligibility`:
     `TestQAResultsBreakStopsEveryCountingReader` fails in its verdict
     subtests, which call `QAReportEligibility` directly. Record also that
     `TestQAReportAcceptRefusesARowAfterABlankLine` passes under this
     mutation and why: `RequireOwnRunEvidence` returns the same break after
     the eligibility check, so the command still refuses the report;
   - the break guard removed from `EnvironmentRowsNeedingOverride`:
     `TestQAResultsBreakStopsEveryCountingReader` fails on the reduced shape,
     whose frontmatter declares one environment row;
   - `detectMechanicalReportShape` still reconciles the counts after a break:
     `TestMechanicalReportShapeNamesARowAfterABlankLine` fails on its count
     finding;
   - the results-table repair removed from `runQAGate`:
     `TestQAResultsRowAfterABlankLineGetsOneRepairTurn` fails;
   - the seeded-table line removed from the QA prompt:
     `TestQAPromptTellsTheGateToContinueTheSeededTable` fails.
5. MUST execute the tests task_01, task_02 and task_03 name in their
   Verification.
6. MUST search `docs/user-guide/`, every copy of the Roundfix skill and every
   copy of the `qa-gate` skill and record: the qa-report guide and both copies
   of the Roundfix `settle` reference begin the rule with "Every Results row
   must belong to a table under" and carry the refusal line; the settle guide
   carries the settle refusal line; the implement guide and both copies of the
   Roundfix `implement` reference carry "When a Results row of the newest QA
   Report lies outside every table, the Daemon returns the QA turn once"; both
   copies of `qa-gate` carry `### Results rows stay inside the table`; and no
   guide says that only one table under `## Results` is read.
7. MUST record, as the outside-evidence row, the two sources `_prd.md` →
   Acceptance evidence names, each with its date and not blocked: the Pantheon
   report the adopted reference records
   (`references/2026-11-10-a-qa-report-results-row-after-a-blank-line-is-read-by-no-reader.md`,
   measured 2026-11-09), comparing offline its reduced shape with the fixture
   of requirement 3; and the GitHub Flavored Markdown passage the PRD quotes
   ("The table is broken at the first empty line, or beginning of another
   block-level structure"), comparing it with the refusal of a row after a
   blank line.
8. MUST record that the PRD declares no Unreachable Acceptance and count no
   row in `rows_blocked_declared`.
9. MUST perform the glossary check (`docs/agents/domain.md`): record that the
   PRD and TechSpec declare `changes: **QA Report**`, that `CONTEXT.md`
   carries the new sentence in that entry, that the Task commit range changes
   no other glossary entry, and that
   `go test -timeout 20m -count=1 -tags docscontract ./internal/docscontract -run '^TestGlossaryPhrasesAreDefinedTerms$'`
   passes.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`,
      `dist/npm/roundfix/package.json`, any non-test file under
      `internal/cli/`, anything under `internal/agent/`, `internal/baseline/`
      or `docs/agents/`, `internal/speccheck/report.go`, any file under
      `docs/adr/`, or any skill other than the Roundfix and `qa-gate` files
      the Tasks declare;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance. It supports the row with the pre-PR equivalent of each
    control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
12. MUST, for the Outcome for the Archive Record, run only
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0371-a-results-row-after-a-blank-line-is-refused-naming-the-line --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, cite evidence an earlier Run left, or keep
      evidence under a name the repository ignores.

## Subtasks

- [ ] Build the binary and record the repository Verification.
- [ ] Run the named tests and reproduce the defect's shape through the binary.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A Results row after a blank line is refused with its line and repair by `qa-report accept`, the Settle Command, the eligibility decision and the mechanical stage, in the tests and through the binary, and the repaired report is accepted.
- [ ] The Daemon tells the gate where rows go and runs exactly one repair turn for a broken report.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0371-a-results-row-after-a-blank-line-is-refused-naming-the-line/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-11; Invariants 1-5; Declared breaks; Testing Approach; Build Order 4
