---
task: task_04
spec: 0321-runs-trim-the-shared-build-cache-and-doctor-warns-on-low-disk
status: completed
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes every Implement Run
and the Delivery Queue's repository gate trim each declared Shared Build
Cache by themselves, with ADR-0327's rules, under a thirty-second budget that
never changes their outcome, records each trim as a `daemon.status` Run Event
or a Queue Owner log line, adds the `disk-space:` warning to `roundfix
doctor`, and documents both in the guides, the Roundfix skill references and
`CONTEXT.md`.

Every row runs on this repository's own files, temporary copies of the
audited head, temporary repositories and fake caches. No row reads or writes
the real `~/.roundfix` or the real shared Go build cache, fills a disk,
depends on the host's free space, starts an Agent session, reaches the
network, sends a request to a provider, the judge or a reviewer, reads a
credential, or reads another repository. Every `go test` it runs carries an
explicit `-timeout`, and every row that runs `bin/roundfix` sets `HOME` and
`TMPDIR` to fresh temporary directories, because the CLI loads
`$HOME/.roundfix/config.yml`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1-5 and Invariants 1-5 by running, each with
   `-v` and recording each named test's and subtest's pass:
   `TestMeasureContextStopsWhenTheContextEnds`,
   `TestApplyContextStopsBetweenEntries` and the existing `TestMeasure…` and
   `TestApply…` tests in `./internal/sharedcache`;
   `TestImplementStartTrimsTheDeclaredSharedBuildCache`,
   `TestImplementStartRunsOnWhenTheTrimDoesNot`,
   `TestImplementStartWithoutASharedBuildCacheRecordsNoTrim`,
   `TestDeliveryGateTrimsTheSharedBuildCacheBeforeItsVerification`, the four
   GC Shared Build Cache tests, `TestImplementStartRunsRunRetention`,
   `TestRunRetentionAtStartWarnsAndNeverBlocks`,
   `TestGateReturnsTheLaneEnvironmentSkip`,
   `TestDoctorWarnsWhenALocationIsLowOnDisk` and
   `TestDoctorPrintsTheLinuxLaneLineOnlyWhenTheLaneIsOn` in `./internal/cli`;
   and `TestAvailableBytesReadsTheFilesystemOfAPath` in `./internal/store`.
3. MUST exercise the built binary end to end on GC's unchanged report (API
   Contract 5). In a fresh temporary Git repository whose Project Config
   declares `{home}/gocache` (`go-build`, `max_size_mb: 1`, `idle: 2h`), with
   `HOME` at a fresh temporary directory holding a fake cache (the go
   command's README line and three 1 MiB entries named by 64 hexadecimal
   digits and `-d` in shard `ab`, aged 10 h, 3 h and 0), run
   `bin/roundfix gc --dry-run` with the audited binary and with a binary
   built from the Delivery Base in a temporary `git archive` copy
   (`go build -buildvcs=false -o <tmp>/roundfix-base ./cmd/roundfix`), each on
   the same `HOME`, `TMPDIR` and repository paths with the fixture recreated
   before each run, and record that standard output, standard error and exit
   codes are byte-identical except the values after `Cutoff: ` and
   `Run Retention cutoff: `, which print the invocation's wall clock (the same
   binary run twice differs there); record those two lines raw, then compare
   every other byte exactly after replacing only those two values with a fixed
   placeholder; then run `bin/roundfix gc` the same way with each binary and
   record the same, and that the 10 h and 3 h entries are gone and the 0 h
   entry and README remain.
4. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, and record that the named test fails, naming the case
   whose input makes the mutation observable:
   - the Run-start trim call removed from `runImplementCommand`:
     `TestImplementStartTrimsTheDeclaredSharedBuildCache` fails (the 10 h
     entry remains and no `shared_cache_trim` event exists);
   - a trim outcome other than `trimmed` or `unchanged` turned into an
     implement preflight failure:
     `TestImplementStartRunsOnWhenTheTrimDoesNot/refused_cache` fails (exit
     `2` instead of `0`);
   - the Run-start trim run under the Run's context without the budget:
     `TestImplementStartRunsOnWhenTheTrimDoesNot/budget_ends` fails (outcome
     `trimmed`, not `stopped`);
   - the idle check dropped from both `MeasureContext` and `ApplyContext`:
     `TestImplementStartTrimsTheDeclaredSharedBuildCache` fails (the 0 h
     entry is removed);
   - the Run Event append skipped:
     `TestImplementStartTrimsTheDeclaredSharedBuildCache` fails;
   - the gate trim moved after the gate's Verification:
     `TestDeliveryGateTrimsTheSharedBuildCacheBeforeItsVerification/idle_entries_removed_first`
     fails (its Verification sees the 10 h entry);
   - the doctor threshold compared with `<=`:
     `TestDoctorWarnsWhenALocationIsLowOnDisk/at_the_threshold` fails;
   - the next action printed once per finding:
     `TestDoctorWarnsWhenALocationIsLowOnDisk/below_the_threshold` fails;
   - the free space read at the configured path instead of its nearest
     existing ancestor:
     `TestDoctorWarnsWhenALocationIsLowOnDisk/missing_location` fails.
5. MUST execute the tests task_01, task_02 and task_03 name in their
   Verification, run `go test -timeout 15m -count=1 ./skills` and
   `go test -timeout 20m -count=1 -tags docscontract ./internal/docscontract`,
   and build the non-test code with
   `GOOS=windows go build -buildvcs=false ./internal/store ./internal/sharedcache ./internal/cli`
   and the same with `GOOS=linux`.
6. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - read `docs/references/2026-10-24-test-suite-audit.md` in the audited
     checkout and quote its `## Disk` sentence on the shared build cache
     (31 GB on 2026-10-24, 43.7 GB on 2026-10-22, and that `roundfix gc` does
     not reclaim it);
   - because this Spec's authorization forbids network access, record the Go
     command documentation and the Bazel remote caching documentation as two
     rows `blocked (environment: network prohibited by Spec authorization)`,
     each quoting the passage the PRD cites. Such a row never decides a
     qualifying `partial`.
7. MUST record the PRD's single Unreachable Acceptance entry as exactly one
   row, `blocked (declared unreachable)`, counted once in
   `rows_blocked_declared`: no row may read the real Roundfix Home or disk.
8. MUST perform the glossary and documentation checks:
   - record that `CONTEXT.md`'s **Shared Build Cache** entry holds the
     definition of `_techspec.md` → Documentation with its `_Avoid_` line
     (`docs/agents/domain.md`);
   - search `docs/user-guide`, `.agents/skills`, `skills`, `docs/agents`,
     `docs/adr` and `CONTEXT.md` for wording that says only the GC Command
     trims a Shared Build Cache, and record what it finds;
   - record that the `roundfix` skill's two version fields are equal in the
     canonical copy and the mirror and recorded in
     `skills/testdata/owned-skill-versions.json`, and that each canonical
     reference this Spec changed equals its mirror.
9. MUST verify from the repository history:
   - each Task's changed files stay within its declarations or its `##
     Recorded paths`;
   - every changed Governed Path is listed in `_authorization.md` `paths`;
   - these are byte-identical to the Delivery Base: the Makefile, the CI
     workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every file under
     `internal/config/`, `internal/baseline/` and `docs/agents/`, every skill
     other than `roundfix`, and every ADR, ADR-0327 and ADR-0334 included;
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
    `HOME="$(mktemp -d)" bin/roundfix archive 0321-runs-trim-the-shared-build-cache-and-doctor-warns-on-low-disk --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's `## Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix` or the real shared Go build cache;
    - start an Agent session;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test, the end-to-end GC rows and each sabotage on temporary copies.
- [ ] Record the outside sources, the declared row, the glossary and documentation checks and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] An Implement Run and a repository gate trim the declared cache first and keep their outcome whatever the trim does.
- [ ] `roundfix gc` and `roundfix gc --dry-run` keep their bytes against the Delivery Base binary (apart from the two wall-clock cutoff values), and the doctor line appears only below 10 GiB free.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0321-runs-trim-the-shared-build-cache-and-doctor-warns-on-low-disk/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; User Experience; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; Interfaces; Documentation; API Contracts 1-6; Invariants 1-5; Testing Approach; Build Order 4
