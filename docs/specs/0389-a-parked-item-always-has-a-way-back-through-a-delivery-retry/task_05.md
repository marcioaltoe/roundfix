---
task: task_05
spec: 0389-a-parked-item-always-has-a-way-back-through-a-delivery-retry
status: pending
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes `roundfix deliver
retry` requeue an item parked before it has an Item Branch; reads the delivery
authorization of an archived Spec without a Spec Root directory; returns an
item whose Pull Request holds an earlier candidate to `gating`, so the
candidate is gated again and pushed; names the `review-stale` next action;
documents all of it in the `deliver` guide, the Roundfix skill's `deliver`
reference and the **Delivery Retry** entry; and records ADR-0398.

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
   `go test -count=1 -timeout 40m -v` over `./internal/delivery` with
   `TestRetryRequeuesAnItemParkedBeforeItHasABranch`,
   `TestCheckingReturnsACandidateThePullRequestLacksToTheGate`,
   `TestAReviewStaleParkNamesItsNextAction`,
   `TestExistingGenericParkAnswersAreUnchanged`,
   `TestPendingQuestionAnswersEachBlockerClass`,
   `TestRetryAfterArchiveReentersCheckingWithAPullRequest` and
   `TestDeliveryEngineParksAStaleReview`; over `./internal/cli` with
   `TestADeliveryRetryRequeuesAnItemParkedBeforeItHasABranch`,
   `TestTheDeliveryAuthorizationOfAnArchivedSpecNeedsNoSpecsRoot`,
   `TestTheDeliveryAuthorizationIsReadBeforeTheArchiveCommit`,
   `TestDeliverStatusPrintsTheItemWorktree` and
   `TestDeliverStatusPrintsAParkLinePerParkedItem`; and over `./internal/spec`
   with
   `TestASpecAuthorizationIsReadAtARevisionWhenTheSpecsRootLeftTheWorkingTree`
   and `TestOperationAuthorityDefaultRootUnchanged`, recording every test and
   subtest by name; then the whole `./internal/delivery`, `./internal/store`,
   `./internal/spec`, `./internal/cli` and `./skills` packages and
   `-tags docscontract ./internal/docscontract`, recording each pass.
3. MUST record, through the built binary, that `bin/roundfix deliver --help`
   exits `0`. The Spec adds no command, flag or exit code: the requeued
   retry prints the existing `Retried <slug>: <blocker> -> queued` line, and
   the `review-stale` `Park:` and `Answer:` lines need a parked item in a Run
   Database, which `TestDeliverStatusPrintsTheItemWorktree` and
   `TestDeliverStatusPrintsAParkLinePerParkedItem` of requirement 2 build
   in-process and assert whole, stdout, stderr and exit code.
4. MUST run each sabotage of `_techspec.md` → Testing Approach on its own
   temporary copy of the audited head made with `git archive` and committed
   in a fresh Git repository (the suite guard refuses a tree outside Git),
   with `-timeout` on every `go test`, and record that the named test fails on
   the named case: sabotages 1 to 8 as listed there.
5. MUST execute the tests tasks 01, 02, 03 and 04 name in their Verification.
6. MUST search and record:
   - the `deliver` guide carries the table row and Texts A, B, C and D of API
     Contract 6, and both copies of the Roundfix `deliver` reference carry
     Texts A to D (API Contract 7);
   - `.agents/skills/roundfix` and `skills/roundfix` hold the same `deliver`
     bytes;
   - `CONTEXT.md` carries the sentence of API Contract 8, and ADR-0398 is
     `accepted`;
   - the log line of API Contract 4 appears exactly once in
     `internal/delivery` outside tests;
   - no non-test file in the repository still produces the refusal `use item
     branch "": use item branch: branch is required` for a retry, which the
     whole-package runs of requirement 2 confirm by the new tests passing.
7. MUST record, as outside-evidence rows, sources this Spec did not write,
   each with its date (read 2026-11-13) and none blocked, through the
   passages `_prd.md` → Acceptance evidence quotes:
   - the git mailing list message of Junio C Hamano: compare "it notices that
     the directory is no longer needed to hold anything it cares about, and
     removes it" with the removed `docs/specs` of
     `TestTheDeliveryAuthorizationOfAnArchivedSpecNeedsNoSpecsRoot`;
   - the installed `git`, a measurement this Spec did not design beyond its
     command: in a temporary repository whose `docs/specs/example/` holds one
     committed file, `git rm -r -q docs/specs/example` leaves no `docs/specs`
     directory; record `git --version` and the absence;
   - actions/checkout `src/retry-helper.ts`: compare `defaultMaxAttempts = 3`
     and its use for every network command with ADR-0398's rejected option
     and `_prd.md` → Non-Goals.
8. MUST record that `_prd.md` → Unreachable Acceptance declares none, so
   `rows_blocked_declared` is `0`.
9. MUST perform the glossary check (`docs/agents/domain.md`): record that the
   PRD and TechSpec declare `changes:` **Delivery Retry**, that `CONTEXT.md`
   revises it with the text of API Contract 8, and that
   `TestGlossaryPhrasesAreDefinedTerms` passes.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`,
      `dist/npm/roundfix/package.json`, `docs/agents/setup-context.json`, any
      file under `docs/agents/` or `internal/baseline/`,
      `docs/references/coverage-record.json`, any Vendored Skill under
      `.agents/skills`, any skill other than the Roundfix skill, and any
      Roundfix reference other than `deliver`;
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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0389-a-parked-item-always-has-a-way-back-through-a-delivery-retry --plan --offline; rm -rf "$h" "$t"`
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
- [ ] Record the wording search, the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A branchless item retries to `queued`, the archived Spec's authorization reads without a Spec Root, and an unpushed candidate is gated and pushed.
- [ ] A `review-stale` park names its next action in status and the Pending Question.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0389-a-parked-item-always-has-a-way-back-through-a-delivery-retry/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-4; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contracts 1-8; Invariants 1-5; Declared breaks; Testing Approach; Build Order 5
