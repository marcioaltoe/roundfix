---
task: task_03
spec: 0131-a-failed-gate-accepts-its-repair
status: pending
type: qa
complexity: medium
---

# Task 03: Run the final QA gate

## Overview

Execute the `qa-gate` skill against the assembled change with evidence
independent of Task 02's own tests, and write the canonical dated report under
`qa/`.

## Requirements

1. MUST verify both PRD Goals, both User Stories and all four Core Features
   against the real command surface, not against the implementation Task's
   fixtures.
2. MUST verify the declared intentional break occurs and that nothing outside
   it changed observably, using the characterization as the reference.
3. MUST verify the outside-evidence row: the repository's own active and
   archived Spec corpus still loads, recording where that corpus came from, or
   record the row blocked with its reason.
4. MUST confirm no Task status was written by loading, and that the stale-gate
   error identity is unchanged.
5. MUST exercise the gate against a binary rebuilt from the assembled tree.
6. MUST NOT change implementation code, tests or any sibling Task.

## Subtasks

- [ ] Rebuild the binary the gate exercises.
- [ ] Materialize and execute the QA matrix, including the refusal cases.
- [ ] Verify the declared break and the absence of undeclared ones.
- [ ] Verify the corpus row and the unchanged error identity.
- [ ] Write the dated report with honest counters and the terminal verdict.

## Acceptance Criteria

- [ ] Every Goal, User Story and Core Feature has an observed result with
      recorded evidence; no row passes by inference.
- [ ] The declared break is observed and no undeclared behavior change appears.
- [ ] The corpus row records its provenance, or is recorded blocked with the
      reason it could not be read.
- [ ] The report records the terminal verdict and exact counters.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`
- instruction: `docs/agents/domain.md`

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It checks the report marker only; the matrix and independent evidence above
remain the QA Agent's responsibility.

- `newest="$(find 'docs/specs/0131-a-failed-gate-accepts-its-repair/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && verdict == \"pass\" ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals 1-2, User Stories 1-2, Core Features 1-4, Success Metrics.
- `_techspec.md` → Coverage Map, Testing Approach, Build Order 3.
- `docs/agents/autonomous-work.md` — Daemon verification and terminal QA ownership.
