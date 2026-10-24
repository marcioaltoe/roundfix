---
task: task_04
spec: 0314-gc-trims-the-shared-build-cache
status: pending
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec adds the
`internal/sharedcache` trim engine, the `worktree.shared_caches` key and a GC
section that trims each declared **Shared Build Cache** above its size by
removing only idle Go build cache entries (ADR-0327), declares this
repository's Makefile cache, and documents all of it.

Every row runs on temporary homes, temporary `TMPDIR` roots holding fixture
caches, temporary copies of the audited head and of the Delivery Base, and
this repository's own documents. No row reads or writes `~/.roundfix`,
measures or trims the host's real Go build caches, reaches the network,
sends a request to a provider, the judge or a reviewer, reads a credential,
or reads another repository. Every `go test` it runs carries an explicit
`-timeout`, and every row that runs `bin/roundfix` sets `HOME` and `TMPDIR`
to fresh temporary directories, because the CLI loads
`$HOME/.roundfix/config.yml` and the declared cache lives under `TMPDIR`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 3 and Invariants 1 to 3 for task_01: run
   the four `internal/sharedcache` tests with `-v` and record each pass.
3. MUST verify Success Metrics 4 and 5 and Invariants 4 and 5 for task_02:
   run `TestWorktreeSharedCachesLoadsAndResolvesItsTemplate`,
   `TestWorktreeSharedCachesRejectsAnUnsafeEntry` with its six subtests, the
   three `internal/cli` Shared Build Cache GC tests and every existing
   `TestGC` and `TestWorktreeCaches` test, and record each pass.
4. MUST exercise the command end to end with `bin/roundfix`, with `HOME` and
   `TMPDIR` at fresh temporary directories and the working directory a
   temporary Git repository outside the checkout:
   - build under `TMPDIR` a cache `roundfix-go-build-<uid>` holding the go
     command's README line, a 30 MB entry `aa/<64 hex>-d` dated two days
     ago and a 1 MB entry `aa/<other 64 hex>-d` dated now, and give the
     repository a `.roundfixrc.yml` that declares
     `{tmpdir}/roundfix-go-build-{uid}` with `kind: go-build`,
     `max_size_mb: 20` and `idle: 2h`;
   - run `bin/roundfix gc --dry-run`, record the Shared Build Cache block
     and that both entries remain; run `bin/roundfix gc`, record the block
     and that only the old entry is gone; run it again and record
     `Entries removed: 0`;
   - in a repository with no declaration, record that `gc --dry-run`
     prints no Shared Build Cache line.
5. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, and record that the named test fails, with the case
   whose input makes the mutation observable:
   - make `Measure` accept a directory whatever its README:
     `TestMeasureRefusesWhatIsNotAGoBuildCache` fails on its README-less
     directory;
   - compare entry times with `cutoff.Add(idle)` instead of `cutoff`:
     `TestMeasureKeepsEveryEntryWithinTheLimitOrInUse` fails, planning the
     entry used 30 minutes before `now`;
   - replace the entry-name match by `name != ""`:
     `TestMeasurePlansIdleEntriesOldestFirstUntilWithinTheLimit` fails,
     planning the 48-hour-old `aa/notes.txt`;
   - drop the modification-time reread in `Apply`:
     `TestApplySkipsAnEntryUsedSinceItWasPlanned` fails;
   - drop the minimum idle check in configuration validation:
     `TestWorktreeSharedCachesRejectsAnUnsafeEntry/short_idle` fails;
   - make the GC section apply the plan during a dry run:
     `TestGCDryRunReportsAnOversizedSharedBuildCacheAndChangesNothing`
     fails.
6. MUST execute the tests task_01, task_02 and task_03 name in their
   Verification.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - read `$(go env GOROOT)/src/cmd/go/internal/cache/cache.go` in the
     audited checkout and record the `mtimeInterval` comment and the
     `mtimeInterval` and `trimLimit` values;
   - run `go help cache` and record the sentence the PRD quotes.
8. MUST perform the glossary check: record that `CONTEXT.md` carries
   **Shared Build Cache**, that its definition cites ADR-0327, and that no
   other term changed (`docs/agents/domain.md`).
9. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - these are byte-identical to the Delivery Base: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, every production `.go` file outside
      `internal/sharedcache/`, `internal/config/config.go`,
      `internal/cli/gc.go` and `internal/cli/gc_shared_cache.go`, and every
      skill other than the Roundfix skill's storage reference and version
      fields;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
10. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
11. MUST, for the Outcome for the Archive Record, run only
    `HOME="$(mktemp -d)" bin/roundfix archive 0314-gc-trims-the-shared-build-cache --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's `## Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - measure or trim the host's real Go build caches;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test, the end-to-end GC rows and each sabotage on temporary copies.
- [ ] Record the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] GC trims only idle entries of a declared cache over its limit, reports the same in a dry run, and stays silent without a declaration.
- [ ] Each sabotage makes its named test fail.
- [ ] The scope audit finds no change outside the declared files.

## Verification



The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0314-gc-trims-the-shared-build-cache/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; User Experience; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-4; Invariants 1-5; Testing Approach; Build Order 4
