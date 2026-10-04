---
task: task_05
spec: 0221-a-delivery-queue-that-runs-the-binary-its-item-builds
status: failed
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec adds `roundfix migrate --check`,
the `delivery.item_binary` Project Config declaration and the queue owner's
choice of executable for an item's `implement`, `archive` and `review` steps,
and describes them in the Roundfix Skill and three guides. Every behavior row
is exercised through the built binary or by executing the named tests against
the built tree, with disposable repositories, a disposable Roundfix Home and a
compiled fake item binary. No command reaches GitHub, a provider or the
network, except the reads of published pages the outside-evidence row names.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST reproduce Surface Transcripts 1, 2 and 3 through the built binary
   with `HOME` set to a disposable directory and `NODE_OPTIONS` unset: no
   database; a database whose `PRAGMA user_version` is this binary's schema
   version; and one whose version is 21, each seeded with the `sqlite3` CLI
   in the disposable home. It MUST record stdout, stderr and the exit code of
   each, the database file's checksum and `PRAGMA user_version` before and
   after, and the Roundfix Home's file list after; and that a newer version
   exits `2` naming `roundfix upgrade`.
3. MUST verify, by executing task_03's tests against the built tree, API
   Contract 4, and by loading a disposable repository whose Project Config
   declares `delivery.item_binary` with `roundfix spec check` or another
   command that loads configuration, that the built binary accepts a valid
   declaration and refuses `delivery.item_binary.extra` with its message.
4. MUST verify, by executing task_04's tests against the built tree, API
   Contracts 2, 3 and 5 and Success Metrics 1 to 3, and record the console-log
   lines the tests assert.
5. MUST verify that the Roundfix Skill's `deliver` and `setup` references and
   the `deliver`, `migrate` and configuration guides describe the behavior
   task_01 names, that each mirror equals its canonical file, that the
   version was raised and recorded, and that the `### QA settlement` section
   of every skill is unchanged.
6. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the operator's intervention log entries
   105 to 107; a reproduction with the released binary the operator has
   installed, run in a disposable repository whose Project Config declares
   `delivery.item_binary`, showing the refusal the PRD quotes; the Rust
   bootstrap guide; and SQLite's file-format reference. For each source it
   reaches it MUST record what it read and whether it still supports the
   design; for each it cannot reach it MUST record the row as blocked with the
   reason.
7. MUST verify that this Spec's own artifacts satisfy the promise rule, and
   perform the glossary check of `docs/agents/domain.md`: record whether
   **Delivery Queue**, **Project Config** and **Run Database** still describe
   the behavior, and whether "item binary" needs a glossary term.
8. MUST verify from the repository history that each Task's changed files stay
   within its declarations or its `## Recorded paths`, that every Governed
   Path changed is bounded in `_authorization.md`, and that `Makefile`,
   `.roundfixrc.yml`, `go.mod` and the CI workflows did not change. This row
   reads Task commits, so it declares a `commit_range` input.
9. MUST record the non-waivable Pull Request row as
   `blocked (environment: no open Pull Request)`, naming the Pull Request row
   in its provenance, and support it with the pre-PR equivalent of each
   control:
   - approval: the maintainer's recorded delivery authority;
   - checks and status: the Daemon's repository Verification at the audited
     head;
   - unresolved threads: none, because no Pull Request exists;
   - Merge-Ready acceptance: none yet, because the pre-PR review follows
     archive;
   - review-artifact ancestry: the audited head named as the claimed
     candidate.
10. MUST NOT accept a row whose only evidence is that a file was read.
11. MUST NOT start a Delivery Queue, merge, re-run or otherwise change
    anything on GitHub, and MUST NOT write to the live Run Database under
    `~/.roundfix`.

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

- `newest="$(find 'docs/specs/0221-a-delivery-queue-that-runs-the-binary-its-item-builds/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals; User Stories 1-5; Core Features 1-8; Success Metrics 1-4;
Acceptance evidence; `_techspec.md` → Testing Approach 1-3; API Contract 1;
API Contract 2; API Contract 3; API Contract 4; API Contract 5; Surface
Transcript 1; Surface Transcript 2; Surface Transcript 3; Integration Points;
ADR-0080; ADR-0088; ADR-0091; ADR-0104; ADR-0155; ADR-0167; ADR-0225.
