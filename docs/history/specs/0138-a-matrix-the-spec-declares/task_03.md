---
task: task_03
spec: 0138-a-matrix-the-spec-declares
status: completed
type: qa
complexity: medium
---

# Task 03: Run the final QA gate

## Overview

Execute the `qa-gate` skill against the assembled tree and write the canonical
dated report under `qa/`. This Task declares its own coverage under the rule
this Spec introduces:

- Requirements 1-6 below start with `MUST verify` or `MUST run`, so each is a
  declared verification Requirement.
- Requirement 1 verifies every Acceptance Criterion of Tasks 01 and 02, so each
  of those criteria is a coverage source of its own.
- The remaining Requirements constrain how the gate runs.

The gate plans rows whose provenance covers every declared source and every
non-waivable source, and nothing else. Rows may group sources that one
observation settles.

## Requirements

1. MUST verify every Acceptance Criterion of Tasks 01 and 02 against the
   assembled tree, recording the observed result for each. No row passes by
   inference from a sibling Task's own Result.
2. MUST verify PRD Goals 1 and 2 through what the gate Agent actually reads: the
   qa-gate skill in the assembled tree, and the QA contract as the prompt
   builder renders it. Both MUST state the declaration rule, the coverage rule
   and scoped blocking, and neither may still carry a removed rule.
3. MUST verify the outside-evidence row by replaying evidence this Spec did not
   write, under the rules as now written. Name every source path, and record the
   row as blocked with its reason when a source cannot be read. Replay:
   - Spec 0119's six reports. For each real finding, record the covered source or
     mechanical finding that still reaches it: the missing operations list, the
     dropped template clauses, the rewritten setup context, the Verification
     timeout and the out-of-grant changed paths. For each noise row, record the
     rule that removes it: the six finding-cascaded blocks and the aggregate
     row.
   - Spec 0136's failed and passing reports. Record the source that still covers
     the governed-rename push finding, and why the report-closure row covers no
     source.
   - Spec 0134's recorded analyzer evidence.
   - The QA Reports of Specs 0134, 0135 and 0137. For each, record which of its
     `qa` Task's declared sources no row named.
4. MUST run the repository's own verification gate, `make verify`, and record it
   clean with the resolved Go toolchain version.
5. MUST run `go vet ./internal/agent`, the Go package this Spec changes, and
   compare the result against the delivery target. Record diagnostics identical
   on the delivery target as observed and pre-existing, naming Spec 0123 as the
   owner of the runner mutex copying.
6. MUST verify whether the work introduced, changed or retired a term the
   project glossary should carry, and record the result.
7. MUST exercise the gate against a binary rebuilt from the assembled tree.
8. MUST record the Pull Request row on its equivalent-evidence path.
9. MUST NOT change implementation code, tests, skills or any sibling Task.

## Subtasks

- [ ] Rebuild the binary the gate exercises.
- [ ] Plan rows whose provenance covers every coverage source, then write the
      pending matrix.
- [ ] Execute the Acceptance Criteria of Tasks 01 and 02.
- [ ] Verify Goals 1 and 2 through the skill and the rendered QA contract.
- [ ] Replay the evidence from Specs 0119 and 0134 to 0137.
- [ ] Run `make verify` and `go vet ./internal/agent`, and classify any
      diagnostics.
- [ ] Write the dated report with exact counters and the terminal verdict.

## Acceptance Criteria

- [ ] Every row's provenance names the sources it covers, and every coverage
      source appears in at least one row:
      - Requirements 2-6;
      - each Acceptance Criterion of Tasks 01 and 02;
      - the outside-evidence row;
      - the Pull Request row;
      - the repository Verification.

      No row names a source outside that set.
- [ ] Every Acceptance Criterion of Tasks 01 and 02 has an observed result with
      recorded evidence.
- [ ] The replay names, with source paths, each real finding's covering source,
      each noise row's removing rule, and each earlier report's uncovered
      declared sources. Otherwise it is recorded as blocked with its reason.
- [ ] `make verify` is recorded clean with the toolchain version stated.
- [ ] `go vet ./internal/agent` is recorded, with diagnostics identical on the
      delivery target attributed to Spec 0123.
- [ ] The glossary check result is recorded.
- [ ] The report records the terminal verdict and exact counters.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`
- instruction: `docs/agents/domain.md`

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It checks the report marker only; the matrix above remains the QA Agent's
responsibility.

- `newest="$(find 'docs/specs/0138-a-matrix-the-spec-declares/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && verdict == \"pass\" ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals 1-4; Acceptance evidence; Success Metrics; Regression locks.
- `_techspec.md` → Coverage Map; Testing Approach 3-4; Build Order 3.
- ADR-0104; ADR-0155.
