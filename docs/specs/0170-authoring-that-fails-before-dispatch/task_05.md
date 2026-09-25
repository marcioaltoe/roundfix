---
task: task_05
spec: 0170-authoring-that-fails-before-dispatch
status: pending
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers and
settles the Spec on evidence.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST run `make verify-docs` and record its result as a gate fact.
3. MUST verify through the built binary, on a scratch Spec whose pending Task
   declares `interface: Makefile`, that `roundfix spec check <slug>` exits `1`
   with exactly one `SC-TOOLING-UNDECLARED` naming `Makefile` while the
   authorization record omits it, and reports none once the record and both
   `bounded files:` rows name it.
4. MUST verify with outside evidence, in a scratch copy of the repository, that
   archived Spec 0155 restored as active with its Tasks pending and its
   authorization `paths:` emptied makes the built binary report
   `SC-TOOLING-UNDECLARED` for each of `Makefile`, `.roundfixrc.yml` and
   `.github/workflows/ci-verify.yml`, the three paths whose undeclared edits
   ended 0155's Runs Unresolved; record that the evidence comes from 0155's
   archived Tasks and grant, which this Spec did not author, and record the row
   as blocked with its reason if the copy cannot be made.
5. MUST verify through the built binary that a scratch Spec whose pending Task
   names `internal/cli/cli_test.go` without a guide reports `SC-CLI-UNDOCUMENTED`
   as a gap, exits `1` under `--strict`, and reports nothing once the Task or a
   Task it depends on names `.agents/skills/roundfix/SKILL.md`.
6. MUST verify through the built binary that the captured and inline
   pipe-to-`grep` emptiness forms are refused with `SC-VERIFY-INVERTED-EXIT`,
   that a wrapped `None.` Success Metrics paragraph produces no
   `SC-METRIC-UNDECLARED`, that an adopted Backlog Entry left at its original
   path is refused with `SC-REF-UNRESOLVED`, and that a reason-only archived
   Finding is refused with `SC-ARCHIVE-LICENSE`.
7. MUST verify that neither template contains `bounded paths:` and that a PRD
   whose Tooling authority row is filled from the PRD template passes the built
   binary without `SC-TOOLING-UNBOUNDED`.
8. MUST verify that `roundfix spec check 0170-authoring-that-fails-before-dispatch --strict`
   and a check of every active Spec report no finding, and that the corpus golden
   records `0` for `SC-TOOLING-UNDECLARED` and `SC-CLI-UNDOCUMENTED`.
9. MUST verify from Git that every Governed Path this Spec's Task commits changed
   is in `_authorization.md` `paths:` or is produced by `make skills-sync`.
10. MUST verify that this Spec's own artifacts satisfy the promise rule.
11. MUST NOT accept a row whose only evidence is that a file was read.

## Subtasks

- [ ] Build the matrix from the Requirements above and the sources no
      declaration waives.
- [ ] Execute each row against the built tree and record its evidence.
- [ ] Write the dated QA Report with its verdict.

## Acceptance Criteria

- [ ] The QA Report records a verdict and names, in each row's provenance, the
      sources that row covers.
- [ ] The repository Verification result is recorded as a fact, not re-derived.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0170-authoring-that-fails-before-dispatch/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-4; Core Features 1-4; Success Metrics 1-6; Acceptance
evidence; `_techspec.md` → Testing Approach 1-5; API Contracts 1-3;
ADR-0080; ADR-0083; ADR-0091; ADR-0092; ADR-0093; ADR-0094; ADR-0096; ADR-0097;
ADR-0104; ADR-0116; ADR-0117; ADR-0130; ADR-0131; ADR-0149; ADR-0155; ADR-0156.
