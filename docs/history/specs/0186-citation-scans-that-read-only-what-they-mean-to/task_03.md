---
task: task_03
spec: 0186-citation-scans-that-read-only-what-they-mean-to
status: completed
type: qa
complexity: medium
---

# Task 03: Run the final QA gate

## Overview

This Task is the authored terminal gate for this Spec. It declares what the matrix covers and settles the Spec on evidence. The Spec changes the Spec citation walk in `internal/speccheck` and the Relocation Citation scan in `internal/baseline`. No command output, flag, exit code or schema changes.

- **Where commands run.** Every binary command runs in a disposable repository with the built `roundfix`. No command reaches a provider or the network, except fetching the published sources named in Requirement 4.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify through the built binary, in a disposable repository holding an active Spec, that `roundfix spec check <slug> --strict` reports `SC-ADR-UNLISTED` for an unlisted ADR cited only in `references/task_example.md`'s `## Result`, and reports nothing for the same citation in the top-level `task_01.md` `## Result`.
3. MUST verify, by executing the projection and root-containment tests task_01 and task_02 name against the built tree, the negative cases:
   - a nested `task_*.md` file is read in full;
   - a directory swapped for an outside symbolic link yields no read;
   - an unraced repository still reports its citations;
   - the FIFO and symbolic-link cases are unchanged.
4. MUST verify the outside-evidence row against published sources this Spec did not produce:
   - the Go blog "Traversal-resistant file APIs" (https://go.dev/blog/osroot);
   - the `os.Root` documentation (https://pkg.go.dev/os#Root);
   - MITRE CWE-367 (https://cwe.mitre.org/data/definitions/367.html).

   The gate reproduces the blog's parent-directory symlink attack against the built scan in a disposable repository, where a tracked directory is replaced with a link to an outside directory. It records that nothing outside the repository was read, and where each source came from. When a source cannot be fetched, the row is recorded as blocked with its reason.
5. MUST verify that `GOOS=windows GOARCH=amd64 go build -buildvcs=false -o /dev/null ./cmd/roundfix` exits `0`.
6. MUST verify that this Spec's own artifacts satisfy the promise rule, and check whether the Spec introduced, changed or retired a glossary term that `CONTEXT.md` does not carry.
7. MUST verify from the repository history that each Task's changed files stay within its declarations or its `## Recorded paths`, and that no Governed Path changed.
8. MUST NOT accept a row whose only evidence is that a file was read.
9. MUST NOT write to the live Run Database under `~/.roundfix` or to any path under `~/dev/secondbrain`.

## Subtasks

- [ ] Build the matrix from the Requirements above and the sources no declaration waives.
- [ ] Execute each row against the built tree and record its evidence.
- [ ] Write the dated QA Report with its verdict.

## Acceptance Criteria

- [ ] The QA Report records a verdict and names, in each row's provenance, the sources that row covers.
- [ ] The repository Verification result is recorded as a fact, not re-derived.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0186-citation-scans-that-read-only-what-they-mean-to/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals 1-3; Core Features 1-2; Success Metric 1; Success Metric 2; Success Metric 3; Acceptance evidence
- `_techspec.md` → Testing Approach 1-3
- ADR-0080; ADR-0091; ADR-0093; ADR-0094; ADR-0096; ADR-0097; ADR-0104;
  ADR-0117; ADR-0155; ADR-0156; ADR-0173; ADR-0176; ADR-0177

## Result
