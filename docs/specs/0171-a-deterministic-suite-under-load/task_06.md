---
task: task_06
spec: 0171-a-deterministic-suite-under-load
status: pending
type: qa
complexity: high
---

# Task 06: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers and
settles the Spec on evidence.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify through `go test` that `internal/testwait` bounds a wait by the
   test deadline, fails at once naming the result when the watched work ends
   first, and prefers a ready value over a simultaneous end.
3. MUST verify that `internal/cli/implement_test.go`,
   `internal/daemon/task_engine_test.go` and `internal/worktree/worktree_test.go`
   hold no fixed-duration wait and no unexamined one- or five-second bootstrap
   timeout, and that the four tests named in the adopted Backlog Entry pass a
   `-count=10 -cpu 1,4` stress run executed by the gate while the other
   packages of the suite run beside it.
4. MUST verify that `fetch`, `resolve`, `watch` and `implement` run by the suite
   record the offline lookup's tag in their HOME, that a built binary leaves a
   seeded fresh cache byte-identical, and that a per-test fake still produces
   the upgrade warning.
5. MUST verify that `go test -c ./internal/agent` followed by running the
   compiled binary through a relative path passes the adapter tests.
6. MUST record, in a row whose provenance names its outside source, that the
   wait bound follows Go's documented `testing.T.Deadline` semantics
   (https://pkg.go.dev/testing#T.Deadline) and that the four failures come from
   the QA reports of the Runs named in the adopted Backlog Entry; when a source
   cannot be obtained, record the row blocked with that reason.
7. MUST verify from Git evidence against the delivery target that the only
   production files changed are task_07's bounded corrective paths
   (`internal/worktree/worktree.go` and `internal/worktree/adminlock*.go`),
   and that `internal/cli/cli_test.go`,
   `docs/references/coverage-record.json` and the Makefile did not change.
8. MUST verify that this Spec's own artifacts satisfy the promise rule.
9. MUST NOT accept a row whose only evidence is that a file was read.

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

- `newest="$(find 'docs/specs/0171-a-deterministic-suite-under-load/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-4; Core Features 1-5; Success Metrics 1-5; Acceptance
evidence; `_techspec.md` → Testing Approach 1-6; API Contracts 1-3; ADR-0034;
ADR-0056; ADR-0080; ADR-0089; ADR-0091; ADR-0093; ADR-0096; ADR-0097; ADR-0104;
ADR-0117; ADR-0125; ADR-0126; ADR-0130; ADR-0155; ADR-0156.
