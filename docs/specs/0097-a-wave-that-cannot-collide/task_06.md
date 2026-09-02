---
status: pending
type: qa
---

# Task: QA gate

Verify every deliverable of this Spec against the running commands.

## Work

- A graph whose same-wave Tasks share a path is reported at authoring, naming
  both Tasks, each path, and its source
- The same graph is refused by the Run before any Agent Session opens or any
  Task Worktree is created
- A graph the `needs` chain already serializes is neither reported nor refused,
  including when the ordering is transitive rather than a direct edge
- A package selector, a flag, and a test name never become paths
- Each of the three sources produces a collision on its own, and a Task
  declaring no Context is still covered — the measured case
- The Run never reorders the plan; the refusal names the `needs` edge instead
- Bootstrap at capacity above one produces no lock collision across repeated
  attempts, and a bootstrap that failed after completing its work is classified
  apart from one that failed before starting
- A Task Worktree creation failure names the Run, the Task, and the concurrency,
  and carries the underlying error as evidence
- The glossary check: this Spec reports on a set of paths a Task is known to
  touch, which the glossary does not name — decide whether a term is owed

## Outside evidence

One acceptance row rests on evidence this Spec did not author. Raising Task
Worktree concurrency in repositories this Spec did not build produced the three
failures in its problem statement, including a `prepare` step that failed with
`could not lock config file` after doing its work. Those measurements were taken
before this Spec existed and are recorded in its PRD.

The collision shape is corroborated inside this repository and the corroboration
is the sharper evidence: Spec 0113's `task_05` and `task_07` ran as siblings on
2026-08-26 and died at integration on `internal/speccheck/mechanical_test.go`.
Neither declares a `## Context` section, and both Verifications name
`internal/speccheck/mechanical.go` — which is why this Spec reads Verification
paths and not declared Context alone. Verify that reading against those two
archived Task files rather than trusting this sentence.

## References

- All user stories and core features

## Verification
- `newest="$(find 'docs/specs/0097-a-wave-that-cannot-collide/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && verdict == \"pass\" ? 0 : 1) }" "$newest"`
