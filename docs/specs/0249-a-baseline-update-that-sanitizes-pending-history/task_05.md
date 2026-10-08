---
task: task_05
spec: 0249-a-baseline-update-that-sanitizes-pending-history
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes `roundfix baseline update`
plan and apply the Pending History in the change it plans, makes
`roundfix upgrade` name it, accepts the legacy list-of-maps `unproven` and
prints every Refused Unit reason on one line. Every behavior row runs against
the built tree in temporary repositories with synthetic history and a
temporary home. No row reaches a provider, starts a real Agent Session, reads
a credential, pushes a tag, reads or writes the real `~/.roundfix`, or touches
an adopter repository.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact,
   with `bin/roundfix` rebuilt from the audited head first.
2. MUST verify API Contract 8, Invariant 11 and Invariant 12 by executing
   task_02's five tests and the existing History Sanitize Command tests.
3. MUST verify API Contracts 1 to 6 and Invariants 1 to 9 by executing
   task_03's eight tests and every existing `baseline update` test.
4. MUST verify API Contract 7 and Invariant 10 by executing task_04's three
   tests and every existing upgrade test.
5. MUST run Surface Transcript 1, Surface Transcript 2, Surface Transcript 3
   and Surface Transcript 4 through `go run -buildvcs=false ./cmd/roundfix` in
   that order, against one temporary adopted repository built as the
   TechSpec's Surface Transcripts section describes, under a temporary home.
   It records each output and exit code. Around Surface Transcript 2 it
   records the commit count, `git for-each-ref` and `git cat-file -t
   history-full`, which prove that no commit was created, that no ref other
   than the new tag changed and that the tag is annotated at the preview's
   `HEAD`. That covers Success Metric 1, Success Metric 2 and Success Metric 4.
6. MUST measure Success Metric 3: in a temporary adopted repository without
   history, `roundfix baseline update --format json` reports a `planDigest`
   equal to the `planDigest` of `roundfix baseline plan` for the same Profile
   and decisions, or, when the two commands cannot share inputs, cites
   task_03's digest test. It also records the history status `current`.
7. MUST replay Success Metric 5 through task_04's tests, because a real
   `roundfix upgrade` resolves releases through the GitHub CLI and the network.
   It records the exact notice lines they assert inside and outside a
   repository, and that stdout and the exit code equal the outcome's existing
   values.
8. MUST replay Success Metric 6 through the built binary: in a temporary
   repository with a committed Legacy Archive Folder whose `unproven` uses the
   two Fluxus key sets, `roundfix history sanitize` plans it, and
   `--apply --batch 1` after an annotated `history-full` tag writes a record
   whose `unproven` lines it records. It also records one plan whose refusal
   reason is one line.
9. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - the Fluxus report triaged in the Secondbrain;
   - the adopter release notes sent through the Secondbrain inbox;
   - Django's unapplied-migrations notice;
   - the Git documentation for `git push`.

   For each source it reaches, it records what it read and whether the source
   still supports the design: the key sets and the multi-line refusal,
   notify-then-run, and a local tag that must be pushed. For each source it
   cannot reach, it records the row as blocked with the reason, using
   `blocked (environment: network denied: <host>)` when the Run sandbox denied
   the host.
10. MUST perform the glossary check of `docs/agents/domain.md`. It records that
    `CONTEXT.md` defines **Pending History** and carries the revised
    **Managed Refresh**, **History Sanitize Command**, **Refused Unit**,
    **Lenient Legacy Reading**, **Sanitize Batch** and **History Full Tag**,
    and whether any other term the Spec introduced needs one.
11. MUST verify from the repository history:
    - that each Task's changed files stay within its declarations or its
      `## Recorded paths`;
    - that every Governed Path changed is bounded in `_authorization.md`;
    - that `go.mod`, `.roundfixrc.yml`, the Makefile, the CI workflows, the
      Baseline modules and the Baseline guides did not change;
    - that the `### QA settlement` section of every skill is byte-identical.

    This row reads Task commits, so it declares a `commit_range` input.
12. MUST record the non-waivable Pull Request row as
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
13. MUST NOT accept a row whose only evidence is that a file was read. It MUST
    NOT send any request to a provider, push a tag or write outside temporary
    directories and the Spec's `qa/` folder.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Execute the Task tests and the four Surface Transcripts.
- [ ] Replay the digest, notice and list-of-maps rows.
- [ ] Record the outside evidence, the glossary check and the scope audit.
- [ ] Write the QA report with the Pull Request row blocked.

## Acceptance Criteria

- [ ] Every behavior row passes on recorded evidence.
- [ ] The confirmed update converts, tags and stays uncommitted, and a second update is current.
- [ ] A repository without Pending History keeps today's digest and output.
- [ ] The scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0249-a-baseline-update-that-sanitizes-pending-history/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-7; User Stories 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; API Contract 7; API Contract 8; Invariants 1-12; Surface Transcripts; Current behavior; Risks & Considerations; Build Order 5
- ADR-0254; ADR-0248; ADR-0251; ADR-0100; ADR-0184
