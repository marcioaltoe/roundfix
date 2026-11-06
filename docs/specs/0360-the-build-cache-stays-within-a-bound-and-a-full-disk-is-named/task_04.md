---
task: task_04
spec: 0360-the-build-cache-stays-within-a-bound-and-a-full-disk-is-named
status: pending
type: qa
complexity: medium
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec holds this repository's
Makefile Go build cache to the Build Cache Bound by removing only idle
entries, adds the temporary directory to Doctor's Free-Space Floor, warns in
the Delivery Queue below the floor, and parks a full filesystem during the
Item Binary build or the repository gate as a Disk-Full Park with an
uncounted retry (ADR-0370), and states the rule in the guides, the README,
the repository rules guide, the Roundfix skill, the glossary and the
amendment notices of ADR-0327 and ADR-0334.

Every row runs on this repository's own files, temporary homes, temporary
caches, temporary repositories and temporary copies of the audited head. No
row starts an Agent session or an ACP adapter, reads or writes the real
`~/.roundfix` or the real home directory, reads, trims or fills a real build
cache or the real disk, runs `go clean -cache`, reads another repository,
reaches the network, sends a request to a provider, the judge or a reviewer,
or reads a credential. Every `go test` it runs carries an explicit
`-timeout`. Every row that runs `bin/roundfix` sets `HOME` and `TMPDIR` to
fresh temporary directories, with `GOMODCACHE` and `GOCACHE` kept at their
`go env` values, because the CLI loads `$HOME/.roundfix/config.yml`, and
removes them afterwards. Every row that runs `make` or `gocache-bound` over a
cache points it at a fresh temporary cache with an explicit `GOCACHE` and
`GOCACHE_BOUND_MB`, never at the session's own cache. The QA never writes
into the checkout while a test suite runs. The gate keeps its evidence under
names the repository does not ignore, `.txt` for command output, and links
only files it kept.

The `go test` package and test runs of this gate do not run `bin/roundfix`
and keep the session's own `HOME` and `TMPDIR`: their tests make their own
temporary homes, and some, such as `TestGuidanceCompositionJourney/rust-cli`,
call an installed formatter (`rustfmt` through rustup) that needs the
session's toolchain configuration under `HOME`. Only rows that run
`bin/roundfix` use fresh temporary homes.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 5 and Invariants 1 to 6 by running
   `go test -count=1 -timeout 20m -v` with
   `TestBoundMainTrimsOnlyIdleEntriesAboveTheBound` and
   `TestBoundMainRefusesAnInvalidRequest` over `./internal/sharedcache`,
   `TestMakeBoundsItsOwnGoBuildCache`,
   `TestVerifyDocsRunsTheDocumentationContractsOnceUnlessSkipped` and
   `TestMakeSharesOneGoBuildCache` over `./internal/verifyselect`,
   `TestDoctorWarnsWhenALocationIsLowOnDisk`,
   `TestCommandFreeSpaceNamesEachLocationBelowTheFloor`,
   `TestTheItemBuildOnAFullDiskIsADiskFullFailure`,
   `TestTheRepositoryGateOnAFullDiskParksDiskFull` and
   `TestDeliveryStepParksWhenTheItemBuildFails` over `./internal/cli`,
   `TestTheQueueWarnsOfLowDiskBeforeTheRunAndTheGate`,
   `TestAFullDiskParksTheItemAsDiskFull`, `TestADiskFullRetryIsNotCounted`,
   `TestEveryBlockerHasAParkClass` and `TestParkClassesAndNextActions` over
   `./internal/delivery`, and `TestTheStoreDoesNotCountADiskFullRetry` over
   `./internal/store`, recording every row by name; then the whole
   `./internal/sharedcache`, `./internal/verifyselect`, `./internal/cli`,
   `./internal/delivery` and `./internal/store` packages, and
   `-tags docscontract ./internal/docscontract`, recording each pass.
3. MUST reproduce the bound through the real Makefile in a temporary copy of
   the audited head made with `git archive`, with a fresh temporary `GOCACHE`
   directory:
   - run `make build GOCACHE=<temporary cache> GOCACHE_BOUND_MB=1` and record
     the `gocache-bound: warning: <cache> holds <bytes> bytes, above its
     1048576-byte bound, and no other entry is idle for 2h0m0s; the build
     continues` line, and that `make` exits `0` and the cache still holds
     every entry;
   - set the modification time of some `-d` entries of that cache to more
     than two hours ago with `touch -t`, run `make build` again with the same
     variables, and record the `gocache-bound: <cache>: removed <n> idle
     entries (<bytes> bytes)` line and that only those entries went;
   - run `make vet GOCACHE=<temporary cache>` without `GOCACHE_BOUND_MB` and
     record that no `gocache-bound` line is printed.
   Remove the temporary copy and the temporary cache afterwards.
4. MUST reproduce the Free-Space Floor through the built binary with a
   fresh temporary `HOME` and `TMPDIR` in a temporary repository: record the
   `disk-space:` line of `bin/roundfix doctor`, which is absent when every
   location has 10,240 MiB or more, and record from requirement 2's
   `TestDoctorWarnsWhenALocationIsLowOnDisk` rows that the temporary
   directory is named below the floor; no row fills a disk to reach it.
5. MUST run each sabotage of `_techspec.md` → Testing Approach on its own
   temporary copy of the audited head made with `git archive`, with
   `-timeout` on every `go test`, and record that the named test fails on
   the named case:
   - the `$(MAKELEVEL)` condition dropped: `TestMakeBoundsItsOwnGoBuildCache`
     fails on `sub-make`;
   - `GOCACHE_BOUND_MB ?= 10240` moved outside the undefined-`GOCACHE`
     block: it fails on `explicit cache`;
   - the warning fallback removed from the recipe: it fails on `failing
     bound`;
   - `BoundMain` passing a zero limit:
     `TestBoundMainTrimsOnlyIdleEntriesAboveTheBound` fails on `within its
     bound`;
   - the temporary directory dropped from `inspectFreeSpace`:
     `TestDoctorWarnsWhenALocationIsLowOnDisk` fails;
   - the `low-disk` segment appended after the lane-skip segment:
     `TestTheQueueWarnsOfLowDiskBeforeTheRunAndTheGate` fails on its
     lane-skip row;
   - the log tail scan removed from `diskFullFailure`:
     `TestTheItemBuildOnAFullDiskIsADiskFullFailure` fails on `path in the
     log`;
   - `disk-full` dropped from the store's exemption:
     `TestTheStoreDoesNotCountADiskFullRetry` fails.
6. MUST execute the tests tasks 01 to 03 name in their Verification.
7. MUST search the `deliver` and `doctor` guides, both copies of the
   Roundfix `deliver` and `setup` references, the README, the repository
   rules guide, `CONTEXT.md` and the two amended ADRs, and record that each
   carries its text of `_techspec.md` → API Contract 8, that each `deliver`
   document states the `disk-full` next action exactly once, and that
   ADR-0327 and ADR-0334 carry the amendment notice naming ADR-0370 and an
   `updated_at` of 2026-11-06.
8. MUST record, as outside-evidence rows, three sources this Spec did not
   write, each with its date and none blocked:
   - the oraculum report quoted in
     `references/2026-11-05-the-roundfix-build-cache-fills-the-mac-mini-internal-disk.md`
     (2026-11-05): compare offline its 25 GiB cache and its fall from 7 GiB to
     3.3 GiB free with the bound of 10,240 MiB and the floor of 10,240 MiB;
   - golang/go issue 31948 (opened 2019-05-09), through the passage
     `_prd.md` → Acceptance evidence quotes ("the build can fail, because
     intermediate artifacts get deleted"): record that requirement 3 and
     Invariant 1 never empty the cache and never remove an entry younger
     than two hours;
   - `docs/references/2026-10-24-test-suite-audit.md` (2026-10-24): record
     that the 31 GB it measured is above the new bound.
9. MUST record the PRD's one Unreachable Acceptance entry, the cache on the
   operator's delivery machine across a day of queued Runs after the
   release, as one row counted in `rows_blocked_declared`, with its
   `satisfied-by`.
10. MUST perform the glossary check (`docs/agents/domain.md`): record that
    the PRD and TechSpec declare `adds:` for **Build Cache Bound**,
    **Free-Space Floor** and **Disk-Full Park**, that `CONTEXT.md` defines
    all three, and that `TestGlossaryPhrasesAreDefinedTerms` passes.
11. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the CI workflows,
      `go.mod`, `go.sum`, `skills-lock.json`, `CHANGELOG.md`,
      `internal/app/version.go`, `dist/npm/roundfix/package.json`,
      `docs/references/skill-coverage.json`,
      `docs/references/coverage-record.json`, any Vendored Skill under
      `.agents/skills`, anything under `internal/baseline/assets/`, any file
      under `docs/agents/` other than `setup-context.json` and
      `specific-repository.md`, any skill other than the Roundfix skill files
      tasks 02 and 03 declare;
    - the Makefile diff adds only the bound's variables, comment, target and
      prerequisites, and no line runs `go clean`;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
12. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
13. MUST, for the Outcome for the Archive Record, run only
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0360-the-build-cache-stays-within-a-bound-and-a-full-disk-is-named --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix` or the real home directory;
    - read, trim or fill a real build cache, or run `go clean -cache`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests and reproduce the bound and the floor.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the unreachable row, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The bound removes only idle entries above its bound, once per top-level make, and never fails a build.
- [ ] Doctor and the queue warn below the floor without failing or parking, and a full filesystem parks `disk-full` with an uncounted retry.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0360-the-build-cache-stays-within-a-bound-and-a-full-disk-is-named/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-8; Invariants 1-6; Declared breaks; Testing Approach; Build Order 4
