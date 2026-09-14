---
task: task_02
spec: 0137-a-fixture-that-carries-its-committer
status: completed
type: qa
complexity: low
---

# Task 02: Run the final QA gate

## Overview

Execute the `qa-gate` skill against the assembled tree and write the canonical
dated report under `qa/`. The matrix is narrow: Task 01's Acceptance Criteria,
the PRD Goal, and the repository's own gate plus `go vet` over this Spec's
package.

## Requirements

1. MUST verify every Acceptance Criterion of Task 01 against the assembled
   tree, recording the observed result for each.
2. MUST verify the PRD Goal under the condition that reveals it: a Daemon
   commit in a fixture repository must succeed with host, global and system Git
   configuration excluded, so a discoverable host identity cannot mask the
   result.
3. MUST run the repository's own verification gate, `make verify`, and record it
   clean on the repository's pinned Go toolchain.
4. MUST run `go vet` over the package this Spec changes and record it clean. A
   diagnostic identical on the delivery target is pre-existing and is recorded
   as observed with its owning Spec named; Spec 0123 owns the runner mutex
   copying.
5. MUST verify the outside-evidence row: that Git answers `Author identity
   unknown` when it can discover no identity, recording the Pull Request gate
   failure as the measurement, or record the row blocked with its reason.
6. MUST exercise the gate against a binary rebuilt from the assembled tree.
7. MUST NOT change implementation code, tests or any sibling Task.

## Subtasks

- [ ] Rebuild the binary the gate exercises.
- [ ] Execute Task 01's Acceptance Criteria matrix.
- [ ] Prove the Goal with host configuration excluded.
- [ ] Run `make verify` and `go vet`.
- [ ] Write the dated report with honest counters and the terminal verdict.

## Acceptance Criteria

- [ ] Every Acceptance Criterion of Task 01 has an observed result with
      recorded evidence.
- [ ] The Goal is observed with host, global and system Git configuration
      excluded.
- [ ] `make verify` is recorded clean, with the resolved toolchain version
      stated.
- [ ] `go vet` over this Spec's package is recorded clean, with any
      pre-existing diagnostic attributed to its owning Spec.
- [ ] The outside-evidence row records Git's identity-resolution behavior, or is
      recorded blocked with its reason.
- [ ] The report records the terminal verdict and exact counters.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`
- instruction: `docs/agents/domain.md`

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It checks the report marker only; the matrix above remains the QA Agent's
responsibility.

- `newest="$(find 'docs/specs/0137-a-fixture-that-carries-its-committer/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && verdict == \"pass\" ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goal 1, Core Features 1-3.
- `_techspec.md` → Coverage Map, Testing Approach, Build Order 2.
