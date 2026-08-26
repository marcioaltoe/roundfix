---
status: completed
type: backend
---

# Task: Parse results from the Results table only

The shape detector parses rows from any markdown table in the QA Report, so a
table a gate writes in its prose to justify a row is read as a result. Measured
on Spec 0098 on 2026-08-26: an evidence table comparing the three hook-refusal
cases produced findings against rows named `Case`, `82-line function vs 80`,
`2462-line generated file vs 500`, and `` `sort()` vs `toSorted()` `` — cells of
a comparison, not results. Each became its own blocker on the run that carried
the fix, and no edit to the Results table could clear them.

## Work

- Bound row parsing to the table under the report's `## Results` heading. A
  table under any other heading, or in prose with no heading, is not a Results
  matrix and yields no rows.
- Keep the existing column contract for the Results table itself: a row still
  needs its id, a terminal status, and its provenance.
- A report with no `## Results` heading keeps today's behavior — the shape
  detector still reports that the matrix is missing, rather than silently
  finding zero rows and passing.
- Cover the measured shape directly: a report whose Results table is valid and
  whose prose carries a second table with non-terminal cells yields findings for
  the Results rows alone.

## References

- User Story 1: Gate writes valid report
- User Story 2: Refusal recorded and auditable
- Core Feature 4: Only the Results table is read as results

## Verification
- `grep -q "## Results" internal/speccheck/mechanical.go && go test -count=1 ./internal/speccheck 2>&1 | grep -q "^ok"`

## Result
A markdown table a gate writes as evidence in its prose cannot become a result
row, so it cannot block a later run.

### What the detector actually did

`parseMechanicalReport` already bound the scan to the `## Results` section — it
started at that heading and stopped at the next `## `. The measured blockers
came from inside that window. Spec 0098's report
(`docs/specs/0098-a-hook-that-cannot-outrank-the-gate/qa/qa-report-2026-08-25.md`)
opens `## Results` at line 63 and does not reach `## Findings` until line 231,
so the four `### Row detail` subsections between them were still scanned, and
the comparison table at line 118 was read as rows.

The fix bounds collection to the *table*, not the section:
`internal/speccheck/mechanical.go:897` now delegates to
`parseMechanicalResultsRows`, which finds the `## Results` heading, collects the
first table under it, and stops at the next heading of any depth or at the first
line after the table that is not a table row. `markdownHeading`
(`internal/speccheck/citations.go:1472`) supplies the depth beside the existing
`markdownCells` / `markdownSeparator` primitives.

### Acceptance criteria

**Row parsing is bound to the table under `## Results`.** Measured on the real
artifact by running `RunMechanicalStage` against 0098's report through a scratch
test, before and after the change (scratch file removed after measuring):

| | Findings on `qa-report-2026-08-25.md` |
| --- | --- |
| before | 7 — rows `2`, `14`, `22` (real Results rows), plus `Case` at line 118, `82-line function vs 80` at 120, `2462-line generated file vs 500` at 121, `` `sort()` vs `toSorted()` `` at 122 |
| after | 3 — rows `2`, `14`, `22` only |

The four blockers named in this Task are gone; the three that remain are genuine
Results-table rows whose status cell reads `**fail** — F-2` / `**fail** — F-1`,
which is a different contract question and not this slice.

**The Results table keeps its existing column contract.** `id`, `status`,
`provenance`, and `evidence` are read from the same header-named columns as
before; a row still needs its id, a terminal status, and its provenance. The
whole pre-existing shape suite passes unchanged —
`TestMechanicalReportShape`, `TestBlockedCauseDiagnosticNamesTheLiteral`,
`TestCountDisagreementReportsItsCause`,
`TestMechanicalStageAcceptsThePreconditionRefusalRow`,
`TestMechanicalEvidencePath`, `TestMechanicalCorpusNonRegression`.

**A report with no `## Results` heading keeps today's behavior.**
`TestResultsTableIsTheOnlyRowSource/a_report_with_no_Results_heading_still_reports_the_missing_matrix`
asserts exactly one finding, `Results table has no report rows` — not silent
zero-row success. `TestMechanicalFindingsWithoutRowHintsBlockTheirRefusalCode`
(a report whose only table sits under `## Mechanical rows`) still passes
unchanged.

**The measured shape is covered directly.** `TestResultsTableIsTheOnlyRowSource`
(`internal/speccheck/mechanical_test.go:1082`) carries 0098's comparison table
verbatim in shape and runs it through five placements: under a `### Row detail`
subsection, in prose under no heading, under a later `## Findings` section, and
beside both a valid and a defective Results row. The defective-row case asserts
that the single finding names `row R01`, not a comparison cell.

The canonical passing fixture `testdata/mechanical/report-green.md` now carries
the same comparison under a `### Row detail` subsection, so any regression fails
`TestMechanicalReportShape/green_fixture_has_terminal_typed_rows_and_exact_counts`
as well.

### Focused checks

- `go test -count=1 -run 'TestResultsTableIsTheOnlyRowSource|TestMechanicalReportShape|TestMechanicalEvidencePath|TestMechanicalFindingsWithoutRowHints|TestMechanicalStageAcceptsThePreconditionRefusalRow' ./internal/speccheck -v` — all pass, 5 new subtests plus every pre-existing shape subtest.
- The new tests were run against `HEAD:internal/speccheck/mechanical.go` to prove
  they are not vacuous: `a_comparison_under_a_row-detail_subsection_is_not_a_result`,
  `a_comparison_in_prose_under_no_heading_is_not_a_result`,
  `a_defective_Results_row_is_the_only_finding_beside_a_comparison`, and
  `green_fixture_has_terminal_typed_rows_and_exact_counts` all fail on the old
  parser and pass on the new one. The other two subtests guard behavior the
  section binding already had.
- `gofmt -l internal/speccheck/` — clean. `go vet ./internal/speccheck/...` — clean.
- `go test -count=1 ./...` — every package `ok`, no failures.

### Notes for later Tasks

No new emitted token, finding code, or user-visible string was introduced; every
diagnostic sentence is unchanged. `golangci-lint` is not installed in this
environment, so lint coverage comes from the Daemon's Verification run.
