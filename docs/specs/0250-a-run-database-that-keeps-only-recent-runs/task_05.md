---
task: task_05
spec: 0250-a-run-database-that-keeps-only-recent-runs
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec removes terminal Runs past a User
Config window from the Run Database, keeps the Runs that are running, that a
Delivery Queue references or whose Run Worktree exists, runs the Run
Retention Sweep at most once a day at Run and Delivery Queue start and through
`roundfix gc`, and compacts the database incrementally. Every behavior row
runs the built binary against fixture Run Databases in temporary homes. No
row reads or writes the real `~/.roundfix`, reaches a provider, starts a real
Agent Session or reads a credential.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify API Contract 1, API Contract 2, API Contract 3 and API
   Contract 5 by executing task_01's eight tests and task_02's ten tests and
   the updated `TestRunGCSkipsWhenJournalRetentionIsZero`.
3. MUST verify API Contract 4 and API Contract 6 by executing task_03's eight
   tests.
4. MUST run Surface Transcript 1, Surface Transcript 2, Surface Transcript 3,
   Surface Transcript 4 and Surface Transcript 5 through the built
   `bin/roundfix` with `HOME` set to a temporary home and `NODE_OPTIONS`
   unset, from a temporary Git repository. For Transcripts 1 and 2 it seeds
   the fixture the transcripts describe with the store API of task_01 or with
   SQL on the fixture file, and for Transcript 3 it writes
   `store.run_retention_days: 10` to the temporary User Config. It records
   each output and exit code and compares them with the transcript.
5. MUST replay Success Metric 1 and Success Metric 2 on that fixture: after
   the dry run the database file and every artifact directory are
   byte-identical, measured by checksum; after the live run the removable Run
   has no row in `runs`, `run_events`, `run_agent_selections`,
   `run_token_usage` or `active_run_locks` and no directory, while every kept
   and recent Run keeps the same row counts, and `run_windows`,
   `interactive_defaults` and the Delivery Queue tables have the same
   content.
6. MUST replay Success Metric 3 through task_03's budget test and by starting
   `roundfix gc` twice on a fixture: the second live run after a completed
   sweep removes nothing new; and record that a start within 24 hours does
   nothing, from task_03's once-a-day test.
7. MUST replay Success Metric 4: a fixture created by the built binary
   reports `PRAGMA auto_vacuum` 2 and its page count falls after the live
   `roundfix gc`, and a fixture created in mode 0 reports
   `Compaction: full (converted to incremental)` and mode 2 afterwards.
8. MUST measure Success Metric 5. It generates, in a temporary home, a
   fixture Run Database through the built binary's schema with at least
   200 terminal Runs and 100,000 events of about 1.3 KB each, half of them
   completed before the cutoff. It times one live `roundfix gc` and one
   automatic start sweep under the two-second budget, and records the removal
   and compaction times and file sizes beside the TechSpec's Current behavior
   measurements. The generator is kept under the QA evidence directory. The
   budgeted start MUST return within 3 s; the row records whether it paused.
9. MUST replay Success Metric 6 through
   `TestDeliverStatusIsUnchangedByRunRetention` and record its output.
10. MUST record, as evidence this Spec did not author, each source of
    `_prd.md` → Acceptance evidence:
    - SQLite's PRAGMA documentation and VACUUM documentation;
    - Jim Nelson's GNOME Blogs article of 2015-01-06;
    - the 2026-10-08 measurement of the live Run Database in the TechSpec.

    For each source it reaches, it records what it read and whether the
    source still supports incremental compaction in bounded slices and the
    conversion by `VACUUM`. For each source it cannot reach, it records the
    row as blocked with the reason, `blocked (environment: network denied: <host>)`
    when the sandbox denied the host.
11. MUST perform the glossary check of `docs/agents/domain.md`. It records
    that `CONTEXT.md` defines **Run Retention** and **Run Retention Sweep**
    and carries the revised **GC Command** and **Journal Retention**, and
    whether any other term the Spec introduced needs one.
12. MUST verify from the repository history:
    - that each Task's changed files stay within its declarations or its
      `## Recorded paths`;
    - that every Governed Path changed is bounded in `_authorization.md`;
    - that `go.mod`, `go.sum`, the `Makefile`, `.roundfixrc.yml`, every
      Baseline module and `internal/cli/cli_test.go` did not change;
    - that the Roundfix skill's version was raised and recorded, and that no
      `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
13. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance. It supports the row with the pre-PR equivalent of each
    control:
    - approval: the maintainer's recorded delivery authority;
    - checks and status: the Daemon's repository Verification at the audited
      head;
    - unresolved threads: none, because no Pull Request exists;
    - Merge-Ready acceptance: none yet, because the pre-PR review follows
      archive;
    - review-artifact ancestry: the audited head named as the claimed
      candidate.
14. MUST NOT accept a row whose only evidence is that a file was read. It
    MUST NOT read or write `~/.roundfix`, send any request to a provider, or
    start a workflow.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Execute the Task tests and the five Surface Transcripts on fixture homes.
- [ ] Replay the removal, kept, budget, compaction and queue rows.
- [ ] Measure the generated fixture and record the outside evidence and the glossary check.
- [ ] Write the QA report with the scope audit and the Pull Request row blocked.

## Acceptance Criteria

- [ ] Every behavior row passes on recorded evidence.
- [ ] The dry run changes no byte, and the live run removes exactly the removable Run.
- [ ] The budgeted start returns within 3 s on the generated fixture.
- [ ] The scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0250-a-run-database-that-keeps-only-recent-runs/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-6; User Stories 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; Surface Transcript 5; Invariants 1-10; Current behavior; Build Order 5
- ADR-0255; ADR-0033; ADR-0171; ADR-0172; ADR-0053
