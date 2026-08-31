---
status: completed
type: qa
---

# Task: QA gate

Verify every deliverable of this Spec against the running commands.

## Work

- A `qa` Task's Verification is supplied by Roundfix, rendered into the Task
  file, and an authored one is refused by name
- The derived command cannot accept a verdict outside the domain, exercised
  against the measured case that did
- A finding blocks the rows it names; unnamed rows in the same matrix are
  measured; a finding naming every row still blocks every row
- Withholding is unchanged: a blocking machine fact before a matrix exists still
  withholds the Agent Session
- The two measured citation forms parse, a bare number outside an obligations
  line is still not a citation, and an unrecognised citation names the form that
  is
- The authoring skill carries the characterization obligation and its ordering
- The gate applies the equivalent-evidence path to the Pull Request row and
  records the evidence; a row with no evidence stays blocked
- The skill changes stay inside the recorded authorization's bounded files,
  checked from Git evidence, with the authorization commit preceding the skill
  commit and the generated copies matching `make skills-sync`
- The glossary check: whether this Spec introduced, changed, or retired a term
  the domain context should carry

## Outside evidence

One acceptance row rests on evidence this Spec did not author. Of 201 failed
Tasks measured across five repositories, 123 were the QA gate returning a
verdict rather than code breaking, and one Spec paid six of its eight gate
executions for the Pull Request row alone — repositories this Spec did not
build, measured before this Spec existed. The row records that this measurement
is what establishes the requirement, rather than a rehearsal of the Spec's own
premise.

Read against the counter-number in the same evidence set: eleven non-passing
verdicts in one session all failed on contract rather than business logic, and
the same gate found four real defects no suite would catch. A change that made
the gate cheaper by making it weaker would satisfy the first number and betray
the second.

## References

- All user stories and core features

## Verification
- `newest="$(find docs/specs/0105-the-gates-own-economics/qa -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && verdict == \"pass\" ? 0 : 1) }" "$newest"`
