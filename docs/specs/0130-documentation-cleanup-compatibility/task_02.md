---
task: task_02
spec: 0130-documentation-cleanup-compatibility
status: pending
type: qa
complexity: medium
---

# Task 02: Verify the cleanup compatibility boundary

## Overview

Execute the qa-gate skill against the bounded repair using independent evidence.
Complete the canonical report at `qa/qa-report-2026-09-09.md`; preserve materialized
rows, counters and outcomes rather than replacing the report.

## Requirements

1. MUST verify the PRD's four Core Features and OE-1 against real filesystem,
   Git and regeneration boundaries in isolated temporary repositories.
2. MUST validate the five-file authorization scope and immutable main ancestry;
   separate skipped original historical objects from exercised controlled cases.
3. MUST leave implementation/tests untouched; own only this Task and Spec-local
   QA artifacts. Confirm no glossary term changed under the Vocabulary Contract.

## Subtasks

- [ ] Materialize and execute the canonical QA matrix, including refused cases.
- [ ] Exercise regenerated outputs and strict/frozen boundaries independently.
- [ ] Audit source scope, references, outside evidence and glossary currency.
- [ ] Record the report and this Task's criterion evidence.

## Acceptance Criteria

- [ ] CF-1–CF-4 and OE-1 have auditable observed results, with no inferred pass.
- [ ] Missing historical objects never count as executed historical Task audits.
- [ ] No path outside the approved five-file implementation scope was changed.
- [ ] The report honestly records the terminal verdict and exact counters.
- [ ] The glossary check confirms the Vocabulary Contract.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`
- instruction: `docs/agents/domain.md`
- instruction: `docs/agents/specific-repository.md`

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It checks the report marker only; the required matrix, source-scope audit and
independent evidence above remain the QA Agent's responsibility.

- `newest="$(find 'docs/specs/0130-documentation-cleanup-compatibility/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && verdict == \"pass\" ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals 1–3, Core Features 1–4, Outside Evidence OE-1.
- `_techspec.md` → Design, Vocabulary Contract, Research and decision evidence.
- `docs/agents/autonomous-work.md` — Daemon verification and terminal QA ownership.
