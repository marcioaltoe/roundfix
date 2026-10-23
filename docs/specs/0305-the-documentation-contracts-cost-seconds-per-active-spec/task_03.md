---
task: task_03
spec: 0305-the-documentation-contracts-cost-seconds-per-active-spec
status: pending
type: qa
complexity: high
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the removed-history
citation pass of `spec.HistoryCitationsInSource` a set lookup, which makes one
`speccheck.Check` of this repository about ten times cheaper with the same
findings, and makes the four corpus tests of `internal/docscontract` share
one Check per active Spec per test process.

Every row runs on in-memory content, temporary fixture corpora, temporary
copies of the audited head and of the Delivery Base, and this repository's
own documents. No row reads or writes `~/.roundfix`, reaches the network,
sends a request to a provider, the judge or a reviewer, reads a credential,
or reads another repository. The gate never runs `roundfix spec judge`,
never runs `roundfix review` against a real provider, and never runs an
archive plan without `--offline`. Every `go test` it runs carries an explicit
`-timeout`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 and 5 for task_01:
   - Run `TestHistoryCitationsInSourceWorkDoesNotGrowWithTheRemovedPaths`
     with `-v` in `./internal/spec` and record both allocation counts it logs
     (Invariant 2).
   - Run `TestHistoryCitationsInSourceNamesEachRemovedPathAndItsDescendants`,
     `TestHistoryRewriteReachesEveryCitationSource`,
     `TestHistoryRewriteResolvesARelativeProsePath`,
     `TestHistoryRewriteWritesAPlainLocatorOutsideMarkdown`,
     `TestHistoryCitationsNameMarkdownThatCitesARemovedPath`, every
     `TestHistoryRetained` test in `./internal/speccheck` and every
     `TestBaselineUpdate` test whose name contains `Citation` in
     `./internal/cli`, and record that each passes (Invariant 1).
3. MUST verify Success Metric 2: in two temporary clones of the checkout
   made with `git clone --local` outside the checkout, one checked out at the
   Delivery Base and one at the audited head, so that both keep the
   repository's history, compile `go test -c -tags docscontract` of
   `./internal/docscontract` and run
   `TestCheckLoopOrderRepositoryAgrees` from each binary three times
   alternately, with `-test.timeout 10m`, recording each run's user CPU time
   and the 1-minute load average. Record the ratio of the median user CPU
   times; the row passes when the audited head's median is at most a fifth
   of the Delivery Base's.
4. MUST verify Success Metrics 3 and 4 for task_02:
   - Run `TestTheCorpusTestsCheckEachSpecOncePerProcess` and record its
     pointer and operation results (Invariant 3).
   - Run `go test -count=1 -timeout 20m -tags docscontract -v
     ./internal/docscontract` once and record the number of active Specs, the
     package's wall time, each corpus test's elapsed time and the
     `TestCheckCorpusBudget` log line; record that only the first corpus test
     to run pays the sweep and that `TestStageScopeDefaultSweepIsUnchanged`
     pays one `CheckStage` per Spec.
5. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, and record that the named test fails, with the case
   whose input makes the mutation observable:
   - restore the loop over every target in the `Locate` closure of
     `HistoryCitationsInSource`, with `strings.HasPrefix(candidate,
     target+"/")`: `TestHistoryCitationsInSourceWorkDoesNotGrowWithTheRemovedPaths`
     fails, because its 2,000 removed paths are longer than 32 bytes and each
     comparison allocates;
   - look up only the candidate itself, without its parent directories:
     `TestHistoryCitationsInSourceNamesEachRemovedPathAndItsDescendants`
     fails in its descendant and nested removed paths cases;
   - stop at the first target found: the same test fails in its nested
     removed paths case;
   - make `sharedCorpusSweep` compute a new sweep on every request:
     `TestTheCorpusTestsCheckEachSpecOncePerProcess` fails on its pointer
     assertion;
   - make the sweep call `checkCorpusSpec` twice per Spec:
     the same test fails on its operation count of 2;
   - add one to any count of `testdata/corpus-golden.json`:
     `TestCheckCorpusGolden` fails, reading the shared sweep;
   - make `speccheck.CheckStage` append one skipped detector for
     `StageAll`: `TestStageScopeDefaultSweepIsUnchanged` fails in each Spec's
     subtest, comparing the shared `Check` result with the fresh
     `CheckStage` call;
   - rename the `## Project Constraints` heading of one active Spec's
     `_prd.md` with a TechSpec to `## Project Constraintz` in the copy:
     `TestCheckActiveCorpusHasNoErrors` fails with `SC-CONSTRAINT-MISSING`
     for that Spec, reading the shared sweep.
6. MUST execute the tests task_01 and task_02 name in their Verification.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - read the adopted Backlog Entry under `references/` in the audited
     checkout and record its measured 764 s and 541 s with five active Specs;
     record the per-Check ratio of row 3 next to them as a ratio only, never
     as a predicted CI time;
   - because this Spec's authorization forbids network access, record the Go
     `testing.AllocsPerRun` documentation as `blocked (environment: network
     prohibited by Spec authorization)`, quoting the passages the PRD cites,
     and record that `go doc testing.AllocsPerRun` in the audited checkout
     prints the same two sentences. Such a row never decides a qualifying
     `partial`.
8. MUST record Success Metric 6 as exactly one row, `blocked (declared
   unreachable)`, against the PRD's single Unreachable Acceptance entry, and
   count it once in `rows_blocked_declared`: no row may read a CI run.
9. MUST perform the glossary check: record that the Spec declares no term,
   that `CONTEXT.md` is byte-identical to the Delivery Base, and whether any
   term the Spec introduced needs an entry (`docs/agents/domain.md`).
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - no Governed Path changed, as `_authorization.md` records `paths: []`;
    - these are byte-identical to the Delivery Base: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`,
      `internal/testfixture/parallel_tests_test.go`, every file under
      `internal/speccheck/`, `internal/docscontract/testdata/corpus-golden.json`,
      `internal/spec/git_locator.go`, `internal/spec/live_citations.go`, every
      `_test.go` file under `internal/spec` and `internal/cli` that exists at
      the Delivery Base, every file under `docs/user-guide/`, `.agents/` and
      `skills/`, and `CONTEXT.md`;
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
    `bin/roundfix archive 0305-the-documentation-contracts-cost-seconds-per-active-spec --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset and with `HOME`
    set to a fresh empty temporary directory created for this row (the CLI
    loads User Config from `$HOME/.roundfix/config.yml`, so the default home
    would read `~/.roundfix`), and write the
    report's `## Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test, the per-Check measurement and each sabotage on temporary copies.
- [ ] Record the outside sources, the declared row, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The removed-history pass keeps every citation, its allocations do not grow with the Removed History Paths, and one Check costs at most a fifth of the Delivery Base's user CPU.
- [ ] The corpus tests share one Check per active Spec, and each sabotage makes its named test fail.
- [ ] The scope audit finds no Governed Path change and no change to the files row 10 names.

## Verification



The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0305-the-documentation-contracts-cost-seconds-per-active-spec/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-2; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; Measurements; API Contracts 1-3; Invariants 1-4; Testing Approach; Build Order 3
