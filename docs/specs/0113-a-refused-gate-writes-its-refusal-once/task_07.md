---
status: pending
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
