---
task: task_05
spec: 0183-run-storage-that-says-what-it-holds
status: pending
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers and
settles the Spec on evidence. This Spec changes production code in
`internal/store`, `internal/cli`, `internal/config` and `internal/worktree`,
and the guides that describe `gc`, the operational sweep and Doctor. Every
behavior row is exercised by executing the named tests against the built tree
or through the built binary. Tests run in disposable repositories with
disposable Roundfix Homes. The only reads of the live `~/.roundfix` are the
read-only outside-evidence rows below, and no command reaches a provider or the
network.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify, by executing the store and `gc` tests against the built tree,
   that a second `gc` reports `Runs pruned: 0` and a dry run then
   `Runs eligible: 0`, that a Run whose events are gone but whose artifact
   directory survived is removed and reported, that a store holding only
   emptied candidates is pruned while another holder keeps the write lock, and
   that the operational sweep prints no line over emptied Runs.
3. MUST verify, by executing the Doctor storage tests against the built tree,
   that a reclaimable Run prints `storage: found` naming `roundfix gc`, free
   bytes at the threshold name `roundfix gc compact`, nothing left prints
   `storage: ok`, a missing database prints `ok (no Run Database)` and is not
   created, an unreadable database prints `partial`, Doctor's exit code stays
   unaffected and the database bytes are unchanged.
4. MUST verify, by executing the sanitation tests against the built tree, that
   in a bare clone with a linked worktree a root derived from the checkout path
   is classified `orphaned` and `--apply` removes its eligible directory,
   while a root equal to neither default stays `overridden` with both defaults
   named.
5. MUST verify, by executing the retained-count tests against the built tree,
   that a removed checkout yields no failure and counts a Run Branch through
   the repository key, a bare key is listed, an absent repository counts only
   an existing Run Worktree, and a symlinked recorded root still fails.
6. MUST run the built `roundfix runs list` from this repository's checkout,
   which only reads the live Run Database, and record whether stderr carries
   any `inspect retained terminal Runs` warning. The gate MUST record this as
   evidence this Spec did not author: the Runs of the removed checkouts
   `~/dev/roundfix-wt-0173` to `-0180` were recorded by other sessions on
   2026-09-28, and v0.20.0 printed 21 such warnings on 2026-09-29. When the
   live Run Database cannot be read by the built binary, it MUST record the row
   as blocked with its reason.
7. MUST count, read-only, the files under `~/.roundfix/artifacts` whose text
   carries `pruned Run storage runs=` followed by `journal_rows=0 artifact_bytes=0`,
   and record the count and the command as evidence this Spec did not author:
   console logs written by earlier Runs, 133 files on 2026-09-29. When the
   directory is unavailable, it MUST record the row as blocked with its reason.
8. MUST record, as published evidence, that SQLite documents
   `PRAGMA freelist_count` as the number of unused pages in the database file
   (https://sqlite.org/pragma.html) and DBSTAT as one row per btree page
   (https://sqlite.org/dbstat.html), and verify from the Doctor storage code
   that it reads the pragma and never `dbstat`.
9. MUST verify that `docs/user-guide/commands.md`, the Roundfix skill and its
   mirror, and `CONTEXT.md` describe the reclaimable count, the silent sweep
   and the `storage` check, and that `make skills-sync-check` exits `0`.
10. MUST verify that this Spec's own artifacts satisfy the promise rule.
11. MUST verify from the repository history that the changed files stay within
    the paths the Tasks declare or the Daemon recorded, and that every governed
    path is bounded in `_authorization.md`.
12. MUST NOT accept a row whose only evidence is that a file was read.
13. MUST NOT write to the live Run Database or any Artifact Root under
    `~/.roundfix`, and MUST NOT run `roundfix gc` or `gc sanitize --apply`
    against the live Roundfix Home.

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

- `newest="$(find 'docs/specs/0183-run-storage-that-says-what-it-holds/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-5; Core Features 1-4; Success Metrics 1-5; Acceptance
evidence; `_techspec.md` → Testing Approach 1-5; API Contracts 1-4; ADR-0023;
ADR-0032; ADR-0033; ADR-0053; ADR-0080; ADR-0090; ADR-0091; ADR-0093; ADR-0094;
ADR-0096; ADR-0104; ADR-0107; ADR-0117; ADR-0155; ADR-0156; ADR-0171; ADR-0172.
