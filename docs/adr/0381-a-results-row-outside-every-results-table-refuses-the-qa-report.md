---
status: accepted
created_at: 2026-11-10T19:00:00Z
updated_at: 2026-11-10T19:00:00Z
deprecated_at: null
superseded_by: null
---

# A Results row outside every Results table refuses the QA Report

On 2026-11-09 Pantheon's QA gate closed a report whose `## Results` table had a
blank line in the middle: the Daemon had seeded the table with carried rows, and
the QA Agent wrote its own rows after a blank line, without a header. Roundfix
0.137.7 read nothing after that blank line. Its eligibility check answered
`rows_blocked_environment is 1; expected 0`, because the pre-PR Pull Request row
and the declared row were among the rows it never read, and the mechanical stage
asked to set both blocked counts to the number of rows it had read. Both named
the wrong cause, and Pantheon archived with the QA Archive Override. The
maintainer chose on 2026-11-10: "O leitor passa a ler a tabela de resultados
inteira ou a recusar apontando a linha em branco exata, e a skill/prompt de QA
passa a mandar continuar a tabela sem linha em branco."

## Decision

**One reader decides which lines are Results rows.** In the section under the
first `## Results` heading, up to the next heading of depth one or two and
outside fenced blocks, a table starts at a table line followed by a separator
line, and its data rows are the table lines below it up to the first line that
is not one. The data rows of every such table whose header has a `Status`
column are the report's Results rows, under a deeper heading too. This is what
the eligibility check already counted, so every archived report reads as before.
`spec.ReadQAResults` is that reader, and every reader of Results rows takes them
from it: the eligibility decision that `roundfix qa-report accept`, the Settle
Command, `roundfix archive` and the Daemon's QA settlement share, the count of
environment rows needing the QA Archive Override, the Carried Row exemption of
the Borrowed Evidence check, the mechanical stage that validates and carries a
previous report's rows, and the Daemon's count of the rows it seeds.

**A table line that no header leads is refused, never skipped.** A table line in
that section that belongs to no table, such as a row after a blank line, makes
the report unacceptable. The refusal names the report path, the line of the
first such row, the table above it, and the repair: the blank line or lines to
remove so the row continues the table, the line below which to move a row that
prose separates from the table, or the separator line a header lacks. The
mechanical stage reports the same row as one report-shape finding and
reconciles no blocked count against a table it could not read whole.

**The cause is closed where rows are written.** The Daemon's prompt for the
seeded report and the qa-gate skill tell the gate to add each row directly
below the last row of the seeded table, with no blank line and with the same
three columns, `#`, `Status` and `Provenance`.

**The Daemon gives the gate one repair turn.** When the QA turn ends with such a
row in the newest report, the Daemon returns the turn once, naming the line and
the repair, in the same turn as an unresolved evidence link (ADR-0351) or
Borrowed Evidence (ADR-0376). A row still outside after it refuses settlement
with the same text.

**Archived reports keep their bytes.** Nothing rewrites a QA Report; archive
never re-reads an archived Spec's report.

## Considered Options

- Read rows through the blank line. The GitHub Flavored Markdown specification
  says "The table is broken at the first empty line, or beginning of another
  block-level structure", so GitHub, formatters and reviewers show those lines
  as a paragraph; Roundfix would count rows that no reader sees as rows, and a
  row with a different column count could be attributed to a header only by a
  guess. A refusal with the exact line keeps the record what it renders.
- Read only the first table under `## Results`, as the mechanical stage did.
  The eligibility check read every `Status` table, deeper headings included,
  and an archived pass report carries a second table with the same header;
  narrowing would miscount archived and adopter reports.
- Refuse a second `Status` table. Rejected for the same archived reports.
- No repair turn. As in ADR-0351, the gate that wrote the row is the only
  actor that knows where it belongs, and a refusal without a turn costs a whole
  QA Run for one deleted line.

## Consequences

A report with a row outside every table, which an earlier version accepted
while ignoring the row, is now refused by every QA settlement path, by
`roundfix archive` and by the Delivery Queue's QA archive classification, which
no longer calls it an environment partial. The mechanical stage now validates
the rows of every `Status` table under `## Results`, not only the first.
