---
task: task_10
spec: 0119-spec-contained-authorization
status: pending
type: qa
complexity: high
---

# Task 10: Run the final QA gate

## Overview

Execute the `qa-gate` skill against the assembled Spec with evidence independent
of the implementation Tasks' own tests. Derive the matrix from the PRD, walk the
declared user journeys through the real command surface, and write the canonical
dated report under `qa/`, preserving materialized rows, counters and outcomes
rather than replacing them.

## Requirements

1. MUST verify every PRD Goal, User Story and Core Feature against the real
   command surface and real Git boundaries, not against the implementation
   Tasks' own fixtures.
2. MUST verify each of the three declared intentional breaks actually occurs and
   that nothing outside that list changed observably, using the characterization
   corpus as the regression reference.
3. MUST verify the outside-evidence row: the preserved historical authorization
   records resolve, with their provenance recorded, and MUST record the row as
   blocked with its reason if the revision cannot be reached.
4. MUST audit the actual changed files against the approved grant's bounded
   paths and confirm the grant's ancestry precedes the consuming work.
5. MUST confirm the Vocabulary Contract runs rather than skips, and check
   whether the Spec introduced, changed or retired a glossary term.
6. MUST exercise the gate against a binary rebuilt from the assembled tree, so
   the gate reports defects the running artifact actually contains.
7. MUST classify every finding by user impact and record blocked or skipped rows
   with their reason and equivalent evidence, never as an inferred pass.
8. MUST NOT change implementation code, tests, or any sibling Task; this Task
   owns only itself and the Spec-local QA artifacts.

## Subtasks

- [ ] Rebuild the binary the gate exercises from the assembled tree.
- [ ] Materialize and execute the QA matrix, including the refusal journeys.
- [ ] Verify the three declared breaks and the absence of undeclared ones.
- [ ] Audit changed-file scope and grant ancestry against the approved record.
- [ ] Run the glossary and Vocabulary Contract checks.
- [ ] Write the dated report with honest counters and the terminal verdict.

## Acceptance Criteria

- [ ] Every PRD Goal, User Story and Core Feature has an observed result with
      recorded evidence; no row passes by inference.
- [ ] Each declared intentional break is observed, and the characterization
      corpus shows no undeclared behavior change.
- [ ] The outside-evidence row records the preserved historical records'
      provenance, or is recorded blocked with the reason the revision could not
      be reached.
- [ ] The changed-file audit shows no path outside the approved bounded set, and
      the grant's ancestry precedes the consuming work.
- [ ] The Vocabulary Contract check runs and the glossary currency check is
      recorded.
- [ ] The report records the terminal verdict and exact counters, including any
      blocked or skipped rows with their reasons.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`
- instruction: `docs/agents/domain.md`
- instruction: `docs/agents/specific-repository.md`

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It checks the report marker only; the required matrix, source-scope audit and
independent evidence above remain the QA Agent's responsibility.

- `newest="$(find 'docs/specs/0119-spec-contained-authorization/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && verdict == \"pass\" ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals 1-4, User Stories 1-3, Core Features 1-6, Success Metrics, Decisions: Declared intentional breaks.
- `_techspec.md` → Vocabulary Contract, Coverage Map, Testing Approach, Build Order 6.
- `_authorization.md` → approved bounded paths and sanctioned regeneration.
- `docs/agents/autonomous-work.md` — Daemon verification and terminal QA ownership.
