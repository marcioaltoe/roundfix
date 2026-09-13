---
task: task_04
spec: 0134-a-governed-deletion-the-gate-can-see
status: pending
type: qa
complexity: medium
---

# Task 04: Run the final QA gate

## Overview

Execute the `qa-gate` skill against the assembled tree and write the canonical
dated report under `qa/`. The matrix is narrow by decision: the Acceptance
Criteria of every non-QA Task, both PRD Goals, and the repository suite and
linter running clean. The Agent review before the Pull Request carries the
design and defect reading, and the mechanical stage keeps the changed-path
audit against the grant, so this gate does not repeat either.

## Requirements

1. MUST verify every Acceptance Criterion of Tasks 01, 02 and 03 against the
   assembled tree, recording the observed result for each; no row passes by
   inference from a sibling Task's own report.
2. MUST verify both PRD Goals through a call the described reader actually
   makes, not through the classifier in isolation.
3. MUST run the repository suite and the formatting and lint gates and record
   them clean, using the repository's pinned Go toolchain rather than whichever
   Go the ambient PATH resolves.
4. MUST verify the outside-evidence row: that a command-only regeneration
   declaration resolves the outputs the repository's own Baseline ownership
   declarations derive for it, recording where that list came from, or record
   the row blocked with the reason it could not be obtained.
5. MUST verify both declared intentional breaks occur, and that no refusal
   token changed spelling or condition.
6. MUST exercise the gate against a binary rebuilt from the assembled tree.
7. MUST NOT change implementation code, tests or any sibling Task.

## Subtasks

- [ ] Rebuild the binary the gate exercises.
- [ ] Execute the Acceptance Criteria matrix for Tasks 01-03.
- [ ] Verify both Goals through the reader's own call path.
- [ ] Run the suite, the formatting check and the linter on the pinned toolchain.
- [ ] Record the outside-evidence row with its provenance.
- [ ] Write the dated report with honest counters and the terminal verdict.

## Acceptance Criteria

- [ ] Every Acceptance Criterion of Tasks 01, 02 and 03 has an observed result
      with recorded evidence.
- [ ] Both PRD Goals have an observed result obtained through the reader's own
      call path.
- [ ] The suite, the formatting check and the linter are recorded clean, with
      the resolved toolchain version stated.
- [ ] The outside-evidence row records the ownership list's provenance, or is
      recorded blocked with its reason.
- [ ] Both declared breaks are observed and every existing refusal token is
      recorded unchanged.
- [ ] The report records the terminal verdict and exact counters.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`
- instruction: `docs/agents/domain.md`

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It checks the report marker only; the matrix above remains the QA Agent's
responsibility.

- `newest="$(find 'docs/specs/0134-a-governed-deletion-the-gate-can-see/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && verdict == \"pass\" ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals 1-2, Core Features 1-4, Declared intentional breaks 1-2.
- `_techspec.md` → Coverage Map, Testing Approach, Build Order 4.
