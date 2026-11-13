---
task: task_05
spec: 0385-a-queue-item-stops-only-for-its-own-work
status: completed
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec gives the production
Delivery Queue workflow the owner's log, so a repository gate that waits for
the Machine Verification Capacity logs the wait instead of panicking; keeps
the owner record of a Queue Owner that stops with an item left, so `roundfix
deliver status` reports it; reads a review anchor followed by `.` or `;`; lets
the review validator's sealed session carry a copy of the reviewed checkout's
`.ai-memory.toml`; documents all of it in the `deliver` and `review` guides,
the Roundfix skill's references and the **Queue Owner** entry; and records
ADR-0395.

Every row runs on this repository's own files, temporary repositories,
temporary homes and temporary copies of the audited head. No row starts an
Agent session or an ACP adapter, calls ai-memory, reads or writes the real
`~/.roundfix` or the real home directory, runs the real `gh`, reaches the
network, reads another repository, sends a request to a provider, the judge or
a reviewer, or reads a credential. Every `go test` it runs carries an explicit
`-timeout`. Every row that runs `bin/roundfix` sets `HOME` and `TMPDIR` to
fresh temporary directories, with `GOMODCACHE` and `GOCACHE` kept at their `go
env` values, unsets `GH_TOKEN`, `GITHUB_TOKEN`, `TYPESAFE_API_KEY` and every
Judge Key Variable for that command, puts a fake `gh` script that prints `[]`
first on `PATH`, and removes the temporary directories afterwards. The QA never
writes into the checkout while a test suite runs. The gate keeps its evidence
under names the repository does not ignore, `.txt` for command output, and
links only files it kept.

The `go test` package and test runs of this gate do not run `bin/roundfix` and
keep the session's own `HOME` and `TMPDIR`: their tests make their own
temporary homes and repositories, and some, such as
`TestGuidanceCompositionJourney/rust-cli`, call an installed formatter that
needs the session's toolchain configuration under `HOME`. Only rows that run
`bin/roundfix` use fresh temporary homes.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 4 and Invariants 1 to 5 by running
   `go test -count=1 -timeout 40m -v` over `./internal/cli` with
   `TestTheQueueOwnerLogsAMachineCapacityWaitAndPassesItsGate`,
   `TestADeliveryWorkflowWithoutALogWritesToStandardError`,
   `TestAQueueOwnerThatStopsWithAnItemLeftKeepsItsOwnerRecord`,
   `TestRepositoryGateNamesCapacityHoldersBeforeACancelledWait`,
   `TestTheRepositoryGateHoldsTheMachineCapacityWhileItRuns`,
   `TestDeliverStatusReproducesSurfaceTranscriptOwnerNotRunning`,
   `TestReviewFindingAnchorParsesEachForm`,
   `TestReviewReadsAnAnchorFollowedBySentencePunctuation`,
   `TestReviewBlocksFindingsThatNameNoFileAndLine` and
   `TestTheReviewValidatorSessionCarriesTheRepositoryMarker`, and over
   `./internal/agent` with
   `TestASealedWorkingDirectoryHoldsNothingButTheIdentityMarker` and
   `TestCopySealedIdentityMarkerCopiesOnlyARegularMarker`, recording every
   test and subtest by name; then the whole `./internal/cli`,
   `./internal/agent`, `./internal/baselineacp` and `./skills` packages and
   `-tags docscontract ./internal/docscontract`, recording each pass.
3. MUST record, through the built binary, that `bin/roundfix deliver --help`
   and `bin/roundfix review --help` exit `0`. The Spec changes no command
   surface, flag or output line of its own: the gate's wait line and the
   `Owner:` status line keep their texts, and the states that print them need
   a held capacity unit and a dead owner process, which the named tests of
   requirement 2 build in-process.
4. MUST run each sabotage of `_techspec.md` → Testing Approach on its own
   temporary copy of the audited head made with `git archive` and committed
   in a fresh Git repository (the suite guard refuses a tree outside Git),
   with `-timeout` on every `go test`, and record that the named test fails on
   the named case: sabotages 1 to 10 as listed there.
5. MUST execute the tests tasks 01, 02, 03 and 04 name in their Verification.
6. MUST search and record:
   - the `deliver` and `review` guides carry the four texts of API Contracts 6
     and 7, and both copies of the Roundfix `deliver` and `review` references
     carry them too (API Contract 8);
   - `.agents/skills/roundfix` and `skills/roundfix` hold the same `deliver`
     and `review` bytes;
   - `CONTEXT.md` carries the sentence of API Contract 9, and ADR-0395 is
     `accepted`;
   - no non-test Go file under `internal/cli` reads `workflow.log` except
     `console()` (`grep -n 'workflow\.log' internal/cli/*.go`, excluding
     `_test.go` files, prints only the lines of `console()`);
   - the refusal text of API Contract 4 appears in `internal/agent` exactly
     once outside tests, and `internal/baselineacp` has no reference to
     `CopySealedIdentityMarker`.
7. MUST record, as outside-evidence rows, three published passages this Spec
   did not write, each with its date (read 2026-11-12) and none blocked,
   through the passages `_prd.md` → Acceptance evidence quotes:
   - the Go specification, "Selectors", rule 6: compare "calling or evaluating
     the method `x.f` causes a run-time panic" with the frame `Gate.func1` in
     the adopted Backlog Entry and the wiring assertion of
     `TestTheQueueOwnerLogsAMachineCapacityWaitAndPassesItsGate`;
   - the Go specification, "Handling panics": compare "Any functions deferred
     by `F` are then executed as usual." with the `panic` subtest of
     `TestAQueueOwnerThatStopsWithAnItemLeftKeepsItsOwnerRecord`;
   - ai-memory, "Marker file: `.ai-memory.toml`": compare "A non-git directory
     outside `$HOME` needs a marker in its exact cwd." with the `marker`
     subtest of `TestTheReviewValidatorSessionCarriesTheRepositoryMarker`.
8. MUST record the PRD's one Unreachable Acceptance entry as one row,
   `blocked (declared: a review validator session that runs after the release is recorded by ai-memory under this repository's project, and no new roundfix-review-validator-<n> project appears)`,
   so `rows_blocked_declared` is `1`.
9. MUST perform the glossary check (`docs/agents/domain.md`): record that the
   PRD and TechSpec declare `changes:` **Queue Owner**, that `CONTEXT.md`
   revises it with the text of API Contract 9, and that
   `TestGlossaryPhrasesAreDefinedTerms` passes.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`,
      `dist/npm/roundfix/package.json`, `docs/agents/setup-context.json`, any
      file under `docs/agents/`, `internal/baseline/`, `internal/baselineacp/`
      or `internal/delivery/`, `docs/references/coverage-record.json`, any
      Vendored Skill under `.agents/skills`, any skill other than the Roundfix
      skill;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance. It supports the row with the pre-PR equivalent of each
    control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
12. MUST, for the Outcome for the Archive Record, run only
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0385-a-queue-item-stops-only-for-its-own-work --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix` or the real home directory;
    - run the real `gh`, call ai-memory or reach the network;
    - start an Agent session or an ACP adapter;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests, the smoke rows and the whole packages.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the declared row, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A gate that waits for the capacity logs the wait and passes; an owner that stops with an item left keeps its record.
- [ ] An anchor may end a sentence, and the validator's session carries the repository's marker and nothing else.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0385-a-queue-item-stops-only-for-its-own-work/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-4; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contracts 1-9; Invariants 1-5; Declared breaks; Testing Approach; Build Order 5
