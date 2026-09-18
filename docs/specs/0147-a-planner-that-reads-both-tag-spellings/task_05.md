---
task: task_05
spec: 0147-a-planner-that-reads-both-tag-spellings
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers, in
the form ADR-0155 defines, and settles the Spec on evidence.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify that both tag spellings parse and that malformed, prerelease and
   build-metadata inputs keep their current rejections.
3. MUST verify that selection spans spellings and that an ambiguous highest
   version refuses in preflight, naming both refs and the selector, with no
   proposed version.
4. MUST verify that a bare-tagged repository receives a plan whose base and
   proposal are bare, and that a prefixed repository receives what it receives
   today.
5. MUST verify that the reset inventory and digest keep each entry's spelling,
   ref and commit.
6. MUST verify that the shipped skill and its mirror state the spellings and the
   refusal, and that the copies are identical.
7. MUST verify, against this repository's own tag history, that its plan is
   unchanged.
8. MUST verify that this Spec's own artifacts satisfy the promise rule and that
   every Task commit stayed inside the approved authority record.
9. MUST NOT accept a row whose only evidence is that a file was read.

## Subtasks

- [ ] Build the matrix from the Requirements above and the sources no
      declaration waives.
- [ ] Execute each row against the built tree and record its evidence.
- [ ] Write the dated QA Report with its verdict.

## Acceptance Criteria

- [ ] The QA Report records a verdict and names, in each row's provenance, the
      sources that row covers.
- [ ] Every declared source above appears in some row, and no row covers
      anything else.
- [ ] The repository Verification result is recorded as a fact, not re-derived.
- [ ] The outside-evidence row names where its evidence came from.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It checks the report marker only; the matrix above remains the QA Agent's
responsibility.

- `newest="$(find 'docs/specs/0147-a-planner-that-reads-both-tag-spellings/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && verdict == \"pass\" ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-4; User Stories 1-4; Core Features 1-6; Success Metrics 1-3;
Acceptance evidence; Declared intentional breaks; Regression locks;
`_techspec.md` → Testing Approach 1-6; API Contracts 1-3; ADR-0155; ADR-0104.
