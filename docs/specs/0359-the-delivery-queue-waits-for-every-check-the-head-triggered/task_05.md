---
task: task_05
spec: 0359-the-delivery-queue-waits-for-every-check-the-head-triggered
status: completed
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the Delivery Queue's
checking stage read the GitHub Actions workflow runs of the pull request head
and the head's Expected Workflows beside its listed checks, wait until every
member of that Head Check Set is completed, park `checks-failed` naming each
failed member and its link, name the pending members in the console log, the
`checks-timeout` blocker and the `Turn:` line of `deliver status`, and records
the rule in ADR-0369, the notices in ADR-0311 and ADR-0365, `CONTEXT.md`, the
Roundfix skill and the user guides.

Every row runs on this repository's own files, temporary Run Databases,
scripted command runners and temporary copies of the audited head. No row
starts a real Queue Owner, an Agent session or an ACP adapter, reads or writes
the real `~/.roundfix`, reads another repository, reaches the network or
GitHub, sends a request to a provider, the judge or a reviewer, or reads a
credential. Every `go test` it runs carries an explicit `-timeout`, and `go
test` package and test runs keep the session's own `HOME` and `TMPDIR`; only a
row that runs `bin/roundfix` sets `HOME` and `TMPDIR` to fresh temporary
directories, because the CLI loads `$HOME/.roundfix/config.yml`, passes
`GOMODCACHE` and `GOCACHE` from `go env`, runs with every Judge Key Variable
and `TYPESAFE_API_KEY` unset, and removes the directories afterwards. The QA
never writes into the checkout while a test suite runs, and records
environment variable names only, never values.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 2 and 3 for the reader by running
   `go test -count=1 -timeout 20m -v -run '^(TestExpectedWorkflowsFollowTheTriggerRule|TestCurrentHeadChecksReadsTheWorkflowRunsOfTheHead)$' ./internal/delivery`,
   recording each pass and each subtest, including the superseded run 9, the
   run of another head and the `README.md` that the report leaves out.
3. MUST verify Success Metrics 1 to 3 and Invariants 1 to 3 by running
   `go test -count=1 -timeout 20m -v ./internal/delivery`, recording the
   passes of `TestTheQueueWaitsForWorkflowRunsWhoseChecksHaveNotRegistered`
   (with its subtest), `TestAnExpectedWorkflowWithoutARunKeepsTheItemWaiting`
   (with its subtest), `TestAFailedWorkflowRunWithoutAFailingCheckParks` (all
   eight cases), `TestTheChecksTimeoutNamesWhatItWaitedFor` (three subtests)
   and `TestAReadErrorIsLoggedOnce`, and that every other test of the package
   passes, including the re-run tests of `check_rerun_test.go` and
   `stale_check_test.go`.
4. MUST verify Success Metric 4 by recording the passes of
   `TestTurnLinesNameTheChecksAnItemWaitsFor` and
   `TestTheOwnerRecordsTheChecksAnItemWaitsFor` in the run of requirement 3,
   and of `TestDeliverStatusPrintsTheChecksAnItemWaitsFor` in a
   `go test -count=1 -timeout 30m -v -run 'Deliver|Delivery|Queue|Turn|FollowUp' ./internal/cli`
   run in which every selected test passes; and compare the expected text in
   `TestDeliverStatusPrintsTheChecksAnItemWaitsFor` line by line with Surface
   Transcript 1 of `_techspec.md`.
5. MUST verify the declared breaks from the repository history: between the
   Delivery Base (the merge base of the audited head with `main`) and the
   audited head, each existing test `_techspec.md` → Declared breaks names
   changes only the scripted `gh api` and `git ls-tree` steps or the expected
   blocker and log text it states, with no change to a poll count, an elapsed
   time, a stage or a re-run count; and no other existing top-level test
   function in `internal/delivery/*_test.go` or `internal/cli/*_test.go`
   changes its body.
6. MUST run each sabotage `_techspec.md` → Testing Approach lists on its own
   temporary copy of the audited head made with `git archive`, with
   `-timeout 3m` on every `go test`, and record that the named test fails,
   with the case whose input makes the mutation observable.
7. MUST execute the tests task_01, task_02, task_03 and task_04 name in their
   Verification, and record that `GOOS=windows go build -buildvcs=false ./internal/delivery ./internal/cli`
   and `GOOS=linux go build -buildvcs=false ./internal/delivery ./internal/cli`
   succeed.
8. MUST search `docs/user-guide/`, `docs/adr/`, `CONTEXT.md` and both copies
   of the Roundfix skill, excluding the Spec folder and `docs/history/`, and
   record that: both copies of the `deliver` reference are byte-identical and
   carry the paragraph task_01 prescribes; the `deliver` and configuration
   guides carry the text task_01 prescribes; no file outside `docs/history/`
   still says "even when the listed checks pass" or "When a required GitHub
   Actions check fails"; ADR-0369 exists with status `accepted` and the
   content `_techspec.md` → ADR-0369 gives; and ADR-0311 and ADR-0365 carry
   the notices `_techspec.md` → Amendment notices gives, with `updated_at`
   `2026-11-06T12:00:00Z`.
9. MUST record, as outside-evidence rows, three sources this Spec did not
   write, each with its date and none blocked:
   - the Fluxus report adopted in
     `references/2026-11-05-fluxus-delivery-queue-merged-before-checks.md`
     (2026-11-05): compare offline its "só "Validate PR title" tinha rodado"
     and "CI | Backend e CI | Frontend ainda estavam pending" with the first
     poll of `TestTheQueueWaitsForWorkflowRunsWhoseChecksHaveNotRegistered`,
     and its request for "o nome do check e o link" with the blocker that
     test asserts;
   - the GitHub REST passage `_prd.md` → Acceptance evidence quotes for
     `head_sha`: compare it with the runs read of API Contract 1;
   - the GitHub passage on the default `pull_request` activity types the
     same section quotes: compare it with the `types` rule of API Contract 2.
10. MUST record exactly one row for the PRD's one Unreachable Acceptance
    entry (a released Fluxus queue merging only after every workflow of the
    head completed) as
    `blocked (environment: another repository prohibited by Spec authorization)`,
    counting it once in `rows_blocked_declared`, and back it with the rows of
    requirements 2 and 3.
11. MUST perform the glossary check (`docs/agents/domain.md`): record that
    `CONTEXT.md` defines **Head Check Set** and **Expected Workflow** and
    carries the revised **Check Wait** and **Queue Turn** entries, each
    linking ADR-0369, exactly as task_01 prescribes, and that
    `go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^TestGlossary'`
    passes.
12. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`, `dist/`, the Run Database
      schema in `internal/store`, `internal/delivery/park_class.go`,
      `docs/references/skill-coverage.json`, any existing Archive Record, any
      Baseline asset, `docs/agents/`, any ADR other than ADR-0369, ADR-0311
      and ADR-0365, any Roundfix reference other than `deliver`, any skill
      other than the Roundfix files task_01 declares;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
13. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
14. MUST, for the Outcome for the Archive Record, run only
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0359-the-delivery-queue-waits-for-every-check-the-head-triggered --plan --offline; chmod -R u+w "$h" "$t"; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - start a real Queue Owner or `roundfix deliver start`;
    - send any request to a provider, the judge, a reviewer or GitHub;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the reader, completion rule and Turn tests and compare the transcript.
- [ ] Check the declared breaks in the history and run each sabotage on a temporary copy.
- [ ] Record the documents, the outside sources, the unreachable row, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The checking stage waits for every member of the Head Check Set and parks naming each failed member, and the `Turn:` line and the log read as the TechSpec states, in the tests and the transcript.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0359-the-delivery-queue-waits-for-every-check-the-head-triggered/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-6; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-6; Surface Transcript 1; Invariants 1-5; Declared breaks; Testing Approach; Build Order 5; ADR-0369; Amendment notices
