---
task: task_03
spec: 0315-a-task-that-deletes-a-tracked-file-settles
status: pending
type: qa
complexity: high
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the three readers that
list tracked files and then read them from the working tree
(`CitationSources`, `HistoryCitations`, `StaleHistoryCitations`) skip a
tracked path the working tree no longer holds, so a gated Task whose declared
`deletes:` path was removed from disk only settles through the Spec
Consistency settlement check, while every other read failure, the deletes
check, link findings and the authorization check keep judging the deletion.

Every row runs on temporary repositories and Git configuration under the
Run's temporary directory, with synthetic files. No row reads or writes the
real `~/.roundfix`, reads another repository, reaches the network, sends a
request to a provider, the judge or a reviewer, or reads a credential. Every
`go test` it runs carries an explicit `-timeout`, and every row that runs
`bin/roundfix` sets `HOME` and `TMPDIR` to fresh temporary directories,
because the CLI loads `$HOME/.roundfix/config.yml`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 4 and Invariant 3: run
   `TestAGatedTaskWhoseDeclaredDeletionIsOnlyOnDiskSettles` (both subtests),
   `TestSettlementStillReportsALinkToADeletedFile` and
   `TestSettlementRefusesAnUnauthorizedGovernedDeletion` in
   `./internal/daemon`, and `TestATaskThatLeavesItsDeletesPathDoesNotSettle`
   and `TestATaskThatRemovesItsDeletesPathSettles`, and record each pass.
3. MUST verify Success Metrics 5 and 6 and Invariants 1 and 2: run
   `TestTrackedSourceReadersSkipAPathTheWorkingTreeLacks` and
   `TestTrackedSourceReadersStillRefuseAPathBehindAFile` in `./internal/spec`
   and record each of their six subtests, and run
   `TestStaleHistoryCitationsReportsFilesTheTagHolds` and
   `TestStaleHistoryCitationsIgnoresDirectoriesAndPresentFiles` unchanged.
4. MUST exercise the command end to end: in a temporary copy of the audited
   head made with `git archive` and committed in a fresh temporary
   repository with one file under `docs/history/findings/` committed and
   then removed with `git rm` in a second commit, build `bin/roundfix` there
   with `go build -buildvcs=false -o bin/roundfix ./cmd/roundfix`, and with
   `HOME` and `TMPDIR` at fresh temporary directories:
   - run `bin/roundfix spec check 0315-a-task-that-deletes-a-tracked-file-settles --strict`
     and record its exit code and findings;
   - remove one tracked Go source file from disk only (`rm`, no `git rm`),
     run the same command again, and record that its output holds no
     `inspect citing file` error and that the exit code and findings are
     those of the first run.
5. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, and record that the named subtest fails, with the case
   whose input makes the mutation observable:
   - in `CitationSources`, skip on any error from `historyRewriteRegularFile`
     instead of a missing entry only:
     `TestTrackedSourceReadersStillRefuseAPathBehindAFile/citation-sources`
     fails on the tracked `notes/a.md` behind the file `notes`;
   - in `HistoryCitations`, skip on any error instead of a missing entry
     only: `TestTrackedSourceReadersStillRefuseAPathBehindAFile/history-citations`
     fails on the same input;
   - in `StaleHistoryCitations`, skip on any read error instead of a missing
     entry only:
     `TestTrackedSourceReadersStillRefuseAPathBehindAFile/stale-history-citations`
     fails on the same input;
   - restore `CitationSources`' `docs/history`-only condition:
     `TestAGatedTaskWhoseDeclaredDeletionIsOnlyOnDiskSettles/disk-only` fails
     with `inspect citing file packages/backend/test/gone.test.ts`, and
     `TestTrackedSourceReadersSkipAPathTheWorkingTreeLacks/citation-sources`
     fails.
6. MUST execute the tests task_01 and task_02 name in their Verification.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - read the adopted Inbox Entry under `references/` in the audited
     checkout, written by oraculum before this Spec existed, and record that
     its refusal holds `settlement check: spec consistency` and `inspect
     citing file`, the shape the `disk-only` sabotage row reproduces;
   - run `git ls-files -h` and record the `--cached` and `--deleted` lines;
     in a temporary repository, remove a committed file from disk only and
     record that `git ls-files --cached` and `git ls-files --deleted` both
     list it, then `git rm` it and record that neither does.
8. MUST record Success Metric 7 as exactly one row, `blocked (declared
   unreachable)`, against the PRD's single Unreachable Acceptance entry, and
   count it once in `rows_blocked_declared`: no row may read oraculum's
   checkout or Roundfix Home.
9. MUST perform the glossary and documentation checks:
   - record that the Spec declares no glossary term and that `CONTEXT.md`'s
     **Citation Source** entry still reads "A tracked regular text file"
     (`docs/agents/domain.md`);
   - search the guides, the agent guides and the skill copies for wording
     that tells an Agent to stage a deletion with `git rm` for a `deletes:`
     Task, and record what it finds.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - no Governed Path changed, matching `paths: []` in `_authorization.md`;
    - these are byte-identical to the Delivery Base: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `CONTEXT.md`, every
      file under `internal/baseline/`, `docs/agents/`, `docs/user-guide/`,
      `docs/history/`, `.agents/skills/` and `skills/`, every production
      `.go` file other than `internal/spec/citation_sources.go`,
      `internal/spec/history_entries.go` and `internal/spec/stale_citations.go`;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
12. MUST, for the Outcome for the Archive Record, run only
    `HOME="$(mktemp -d)" bin/roundfix archive 0315-a-task-that-deletes-a-tracked-file-settles --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's `## Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test, the end-to-end strict check rows and each sabotage on temporary copies.
- [ ] Record the outside sources, the declared row, the glossary and documentation checks and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A gated Task whose declared deletion is only on disk settles, a staged one still settles, and a deletion is still judged by the link finding and the authorization check.
- [ ] Each sabotage makes its named subtest fail.
- [ ] The scope audit finds no Governed Path change and no change outside the declared files.

## Verification



The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0315-a-task-that-deletes-a-tracked-file-settles/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-3; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-4; Invariants 1-3; Testing Approach; Build Order 3
