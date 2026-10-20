---
task: task_04
spec: 0283-a-git-locator-is-reachable-in-ci-and-the-storage-report-counts-every-page
status: completed
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec delivers two things:

- the Default Branch Head, read by the Git Locator, the QA Delivery Base, the
  Archive Command's `source_revision` warning and the pre-PR review's default
  base, so a check answers the same in a CI checkout and a developer clone;
- the Storage Snapshot, so `roundfix storage report` groups every page,
  reconciles exactly and works while a writer holds the write-ahead log.

Every row runs on temporary repositories, clones, homes and stores. No row
reads or writes the real `~/.roundfix`, reaches the network, starts Docker,
sends a request to a provider, the judge or a reviewer, reads a credential, or
reads another repository. Clones made by a row read only this repository's own
objects through a local path.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 and 2:
   - Run the tests task_01 names. Record that the CI-like clone and the
     stale-local clone each resolved the locator to its commit, and that
     `SC-HISTORY-RETAINED` reported no finding until the remote-tracking ref
     was deleted.
   - On a temporary copy of the audited head made with `git archive`, make
     `ResolveDefaultBranchHead` read `refs/heads/<name>` before
     `refs/remotes/origin/<name>`. Record that
     `TestResolveGitLocatorIgnoresAStaleLocalDefaultBranch` fails
     (Invariant 2).
3. MUST record, as evidence this Spec did not author, the CI checkout's ref
   layout rebuilt from this repository's history:
   - in a temporary directory, `git init`, then fetch from the audited
     checkout `+refs/heads/*:refs/remotes/origin/*` and
     `+refs/tags/*:refs/tags/*`;
   - delete every tag that contains the commit of Spec 0279's locator, named
     in `_prd.md` → Acceptance evidence;
   - check out the audited head detached, with no origin/HEAD and no local
     `main`.

   Run the audited `bin/roundfix spec check --format json` there and record
   that no `SC-HISTORY-RETAINED` finding names Spec 0279's Archive Record. Then
   build the binary from a `git archive` of the Delivery Base into a temporary
   directory, run the same command in the same clone, and record its
   `SC-HISTORY-RETAINED` finding for that locator. Write no file under the
   checkout.
4. MUST verify Success Metric 3: run the tests task_02 names, and record that
   without origin/HEAD the QA Delivery Base resolved, the archive warning named
   `origin/main`, and the review's `baseTipCommit` was the `origin/main` tip
   over a stale local `main`.
5. MUST verify Success Metrics 4 and 5:
   - Run the tests task_03 names. Record the report's grouped and snapshot
     bytes for the pointer-map store, and the positive WAL bytes for the live
     WAL store. Record that the database and WAL files were unchanged and that
     an idle temporary store gained no `-wal` or `-shm` file.
   - Reproduce Surface Transcript 1, "storage report without a Run Database",
     with `bin/roundfix storage report` and a temporary `HOME`, and record
     every line, empty standard error and exit 0.
   - On a temporary copy of the audited head, drop the pointer-map group from
     the grouped bytes, and record that
     `TestStorageReportGroupsPointerMapPages` fails (Invariant 3).
6. MUST record, as evidence this Spec did not author, the 2026-10-18
   measurement in the adopted Source 2,
   `references/2026-10-18-storage-report-fails-while-the-run-database-wal-is-live.md`:
   - divide its file bytes by 4096 to get the page count;
   - apply the page-structure helper's formula, through a throwaway test in a
     temporary copy of the audited head;
   - record that pointer-map and lock-byte pages times 4096 equal its
     reported difference.
7. MUST execute the tests task_01, task_02 and task_03 name in their
   Verification.
8. MUST record each published source of `_prd.md` → Acceptance evidence, with
   the passage the PRD quotes. Because this Spec's authorization forbids
   network access in the gate, record each web source as `blocked
   (environment: network prohibited by Spec authorization)`. Such a row never
   decides a qualifying `partial`.
9. MUST perform the glossary check of `docs/agents/domain.md`:
   - record that `CONTEXT.md` defines **Default Branch Head** and **Storage
     Snapshot** and carries the revised **Git Locator** and **Delivery Base**;
   - record that the glossary contract tests pass;
   - record whether any other term the Spec introduced needs an entry.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its
      `## Recorded paths`;
    - every Governed Path changed is bounded in `_authorization.md`;
    - the Roundfix skill version was raised and recorded;
    - no `### QA settlement` section changed;
    - no guide or skill copy still says the locator reads "the local default
      branch".

    This row reads Task commits, so it declares a `commit_range` input.
11. MUST verify that these files are byte-identical to the Delivery Base: the
    Makefile, every file under `.github/workflows/`, `go.mod`, `go.sum`,
    `.roundfixrc.yml`, every Baseline module, `docs/agents/setup-context.json`,
    `internal/cli/implement.go`, `internal/worktree/worktree.go`,
    `internal/cli/deliver.go`, `internal/cli/deliver_workflow.go` and
    `internal/cli/deliver_linux_lane.go`.
12. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance. It supports the row with the pre-PR equivalent of each
    control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
13. MUST, for the Outcome for the Archive Record, run only
    `bin/roundfix archive 0283-a-git-locator-is-reachable-in-ci-and-the-storage-report-counts-every-page --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset. It MUST NOT
    do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - start a workflow.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test and each sabotage on temporary copies of the audited head.
- [ ] Rebuild the CI ref layout from this repository and compare the audited and Delivery Base binaries.
- [ ] Recompute the 2026-10-18 storage difference, record the outside evidence, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A locator the fetched default reaches is valid in a CI-like clone and a stale-local clone, and the Delivery Base binary reproduces the 2026-10-18 failure on this repository's own history.
- [ ] The QA Delivery Base, the archive warning and the review base read the fetched default without origin/HEAD.
- [ ] The storage report reconciles exactly on a pointer-map store and a live WAL store, the 2026-10-18 difference is recomputed to the byte, and the scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0283-a-git-locator-is-reachable-in-ci-and-the-storage-report-counts-every-page/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## Unreachable Acceptance

- The GitHub Actions run of the delivery pull request is the only place the
  real CI checkout runs. Requirement 3 rebuilds its ref layout offline; the
  run itself needs an open Pull Request and the network.

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; API Contract 7; API Contract 8; API Contract 9; API Contract 10; Surface Transcript 1; Invariants 1-4; Current behavior; Build Order 4
- ADR-0295; ADR-0296; ADR-0104; ADR-0276
