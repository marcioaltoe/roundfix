---
task: task_05
spec: 0256-an-archive-record-that-survives-a-squash-merge
status: completed
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes the Archive Record name its
delivery pull request, recovers a squash-dropped `source_revision` from the
pull request head ref in every reader, makes the Archive Command warn, and
makes the archive-root readers read only real records. Every behavior row
runs against the built tree, temporary repositories and a local bare
repository standing in for the forge. No row reaches a provider or a real
forge, starts a real Agent Session, reads a credential, or reads or writes the
real `~/.roundfix`.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact,
   with `bin/roundfix` rebuilt from the audited head first.
2. MUST execute every test that task_02, task_03 and task_04 created, the
   existing tests they keep (`TestRunCausesReportAnArchiveRecordWithoutItsRevision`,
   `TestLeftoversKeepTheWorktreeWithoutTheSourceRevision`,
   `TestArchiveCommandLeavesTheArchiveRecord`,
   `TestDeliveryArchiveStageCommitsTheArchiveRecord`), and task_01's phrase
   checks, and record each result.
3. MUST reproduce each of `_techspec.md` → Surface Transcripts 1-5 through
   `bin/roundfix` in a temporary repository whose `origin` is a local bare
   repository with `origin/HEAD` naming `main`, and compare standard output,
   standard error and exit code under the transcript conventions.
4. MUST measure Success Metric 1 from task_04's end-to-end and retry tests:
   the merged record's `pull_request` and the count of Pull Request Record
   Commits after an interrupted run.
5. MUST measure Success Metric 2. In a temporary clone that lacks a record's
   `source_revision`, with a bare remote holding `refs/pull/7/head`, it runs
   `bin/roundfix runs causes --format json` and reports `source_available`
   for a record with `pull_request: "7"` and for one whose adding commit's
   subject ends with `(#7)`. It records `git for-each-ref` before and after,
   and the absence of `FETCH_HEAD`.
6. MUST measure Success Metric 3, with and without the warning, and Success
   Metric 4 from task_02's README test and its recorded characterization.
7. MUST repeat the wording sweep. For each replaced sentence of task_01's
   Requirement 4, it searches `docs/`, `skills/`, `.agents/skills/`,
   `CONTEXT.md` and `AGENTS.md`, outside `docs/history/` and this Spec's
   folder, and records every remaining hit with whether it is allowed.
8. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - the committed measurement
     `docs/references/2026-10-12-squash-merged-archive-source-revisions.md`,
     comparing its three records, pull request numbers and no-ref result with
     the recovery contract;
   - GitHub Docs and the Git `git fetch` documentation, through the passages
     `_prd.md` quotes;
   - the adopter evidence in the two adopted Backlog Entries.

   For each source it reaches, it records what it read and whether the source
   still supports the design. For a source it cannot reach, it records the row
   as blocked with the reason, using
   `blocked (environment: network denied: <host>)` when the Run sandbox denied
   the host. The committed measurement needs no network.
9. MUST perform the glossary check of `docs/agents/domain.md`. It records that
   `CONTEXT.md` carries **Source Revision Recovery**, **Pull Request Record
   Commit** and the revised **Archive Record** and **Archive Command**, and
   whether any other term the Spec introduced needs one.
10. MUST verify from the repository history:
    - that each Task's changed files stay within its declarations or its
      `## Recorded paths`;
    - that every Governed Path changed is bounded in `_authorization.md`;
    - that `go.mod`, `.roundfixrc.yml`, the Makefile, the CI workflows, every
      Baseline module and derived Baseline file, `internal/spec/archive.go`
      and every existing record under `docs/history/specs/` did not change;
    - that the `### QA settlement` section of every skill is byte-identical.

    This row reads Task commits, so it declares a `commit_range` input.
11. MUST record the Unreachable Acceptance criterion of `_prd.md` as not
    measured by this gate, with its `satisfied-by` evidence.
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
    NOT send any request to a provider or a real forge, and MUST NOT write
    outside temporary directories and the Spec's `qa/` folder.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Execute the Task tests, the kept tests and the phrase checks.
- [ ] Reproduce the five transcripts and measure Success Metrics 1 to 4.
- [ ] Repeat the wording sweep; record outside evidence, the glossary check and the scope audit.
- [ ] Write the QA report with the Pull Request row blocked.

## Acceptance Criteria

- [ ] Every behavior row passes on recorded evidence.
- [ ] Every transcript matches through the built product.
- [ ] Recovery writes no ref, and an unreachable forge keeps today's outcome.
- [ ] The scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0256-an-archive-record-that-survives-a-squash-merge/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`


## References

- [_prd.md](_prd.md) — Goals 1-5; User Stories 1-5; Core Features 1-6; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence; Unreachable Acceptance
- [_techspec.md](_techspec.md) — Reader inventory; API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; Surface Transcript 5; Skill and document texts; Testing Approach; Risks & Considerations; Build Order 5
- ADR-0261; ADR-0247; ADR-0104; ADR-0080; ADR-0091; ADR-0097; ADR-0240
