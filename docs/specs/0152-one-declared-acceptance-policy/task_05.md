---
task: task_05
spec: 0152-one-declared-acceptance-policy
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
2. MUST verify that one exported decision answers QA Report eligibility, and
   that both archive and the derived Verification reach it.
3. MUST verify that a newest report whose verdict is `partial` with
   declared-unreachable blocked rows settles the terminal `qa` Task.
4. MUST verify that `fail`, a missing report, an unparseable report and each
   disqualifying `partial` shape still refuse, and that a `pass` carrying an
   environment-blocked row is still accepted — this gate's own report carries
   one.
5. MUST verify that archive accepts and refuses exactly what it did before, for
   the same reasons.
6. MUST verify that the shipped skill and its mirror state the one policy, and
   that every Task commit stayed inside the approved authority record.
7. MUST verify that this Spec's own artifacts satisfy the promise rule.
8. MUST NOT accept a row whose only evidence is that a file was read, and MUST
   NOT treat this gate's own settlement as evidence that the rule is correct.

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
It selects the newest report by name and delegates the verdict judgement to the
one eligibility decision; the matrix above remains the QA Agent's
responsibility.

- `newest="$(find 'docs/specs/0152-one-declared-acceptance-policy/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; roundfix qa-report accept "$newest"`

## References

`_prd.md` → Goals 1-3; Core Features 1-4; Success Metrics 1-3; Acceptance
evidence; `_techspec.md` → Testing Approach 1-6; API Contracts 1-3; ADR-0080;
ADR-0091; ADR-0093; ADR-0096; ADR-0097; ADR-0104; ADR-0117; ADR-0130; ADR-0155;
ADR-0156.

## Invalidation

- Date: `2026-09-21`
- QA Report: `qa/qa-report-2026-09-21.md`
- Dependencies not completed: `task_06`
