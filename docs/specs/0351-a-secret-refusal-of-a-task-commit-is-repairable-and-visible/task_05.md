---
task: task_05
spec: 0351-a-secret-refusal-of-a-task-commit-is-repairable-and-visible
status: completed
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec runs the Task commit's
Secret Scan as the last Settlement Check of every Verification attempt, so a
first-attempt Secret Match returns as Verification Feedback; lists each
secret refusal as a `secret-scan` item of `roundfix runs causes`; names the
failed and skipped Tasks of a `run-unresolved` park with the secret repair as
its next action; and documents the behavior in the Roundfix skill, three user
guides, `CONTEXT.md` and ADR-0361.

Every row runs on this repository's own files, temporary repositories,
temporary Run Databases and temporary copies of the audited head. No row
starts an Agent session or an ACP adapter, reads or writes the real
`~/.roundfix`, reads another repository, reaches the network, sends a request
to a provider, the judge or a reviewer, or reads a credential. Every `go test`
it runs carries an explicit `-timeout`. Every row that runs `bin/roundfix`
sets `HOME` and `TMPDIR` to fresh temporary directories, because the CLI loads
`$HOME/.roundfix/config.yml`, runs with every Judge Key Variable and
`TYPESAFE_API_KEY` unset, and removes the directories afterwards. The QA never
writes into the checkout while a test suite runs, and never writes a
real-looking secret literal into evidence: it records Secret Matches by path,
line and kind only.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 3 and Invariants 1 to 5 by running
   `go test -count=1 -timeout 20m -v` over `./internal/daemon` with
   `TestTaskCommitSecretMatchReturnsAsVerificationFeedback`,
   `TestTaskCommitWithASecretMatchSettlesFailedWithoutACommit`,
   `TestTaskCommitSecretScanRunsAfterTheSettlementChecks`,
   `TestTaskSecretScanReadErrorHaltsTheRun`,
   `TestTaskCommitScansOnlyTheLinesItAdds`,
   `TestTaskCommitAcceptsALettersOnlyFixture`,
   `TestTaskSecretScanUsesHeadAndPreparedPaths`, every test whose name starts
   with `TestSettlement`, and the two changed `TestTaskCycle` tests task_01
   names, and over `./internal/runevent` and `./internal/cli` with the tests
   task_01 names, recording each pass and each subtest.
3. MUST verify Success Metric 4 by running `go test -count=1 -timeout 20m -v`
   over `./internal/runcause` with `TestBuildReportsASecretRefusalAsItsOwnItem`
   and every `TestBuild` and `TestCauses` test, and over `./internal/cli` with
   every `TestRunsCauses` test, recording each pass, and by recording that
   `internal/runcause/signatures.json` and `internal/runcause/testdata/` are
   unchanged between the Delivery Base (the merge base of the audited head
   with `main`) and the audited head.
4. MUST verify Success Metric 5 and Invariant 6 by running
   `go test -count=1 -timeout 20m -v` over `./internal/delivery` with
   `TestARunUnresolvedParkNamesItsFailedAndSkippedTasks` and every
   `TestParkClass`, `TestRetry` and `TestExistingGenericParkAnswersAreUnchanged`
   test, and over `./internal/cli` with
   `TestDeliveryRunResultReadsFailedAndSkippedTasks`,
   `TestDeliverStatusNamesTheSecretRefusalOfARunUnresolvedPark` and every
   `TestDeliverStatus` test, recording each pass and each case.
5. MUST verify the declared breaks from the repository history: between the
   Delivery Base and the audited head, the diffs of
   `internal/daemon/task_secret_scan_test.go`,
   `internal/daemon/task_engine_test.go`, `internal/cli/implement_test.go`
   and `internal/runevent/settlement_label_projection_test.go` change only
   the tests `_techspec.md` → Declared breaks names and add the new tests,
   and the diff of `internal/secretscan` is empty.
6. MUST run each sabotage `_techspec.md` → Testing Approach lists on its own
   temporary copy of the audited head made with `git archive`, with
   `-timeout` on every `go test`, and record that the named test fails, with
   the case whose input makes the mutation observable.
7. MUST execute the tests task_01, task_02, task_03 and task_04 name in their
   Verification.
8. MUST search `docs/user-guide/`, `docs/agents/`, `CONTEXT.md`, `docs/adr/`
   and every copy of the Roundfix skill, and record that: both copies of the
   skill carry the `### Settlement Checks` text task_04 prescribes; each copy
   of the implement, runs and deliver references and its user guide carries
   the phrases task_04's Verification checks; ADR-0361 exists with status
   `accepted`; and no live document outside `docs/history/` still says that
   the Task commit's scan runs only after Verification or without a repair
   turn.
9. MUST record, as outside-evidence rows, two sources this Spec did not
   write, each with its date and none blocked:
   - the Fiscus report quoted in
     `references/2026-11-04-fiscus-a-task-commit-secret-refusal-is-invisible.md`
     (2026-11-04): compare offline its shape (a `credential-assignment` match
     in a new test after every Verification command and both Settlement
     Checks passed, skipped dependents, a bare `run-unresolved` park, no
     `runs causes` line) with
     `TestTaskCommitSecretMatchReturnsAsVerificationFeedback`,
     `TestRunsCausesPrintsASecretRefusal` and
     `TestDeliverStatusNamesTheSecretRefusalOfARunUnresolvedPark`;
   - the GitHub push protection passage `_prd.md` → Acceptance evidence
     quotes ("path: README.md:4" under the secret's type, and "provide a
     detailed message explaining the reason for the block"): compare it with
     the feedback, the `secret-scan` item and the park naming path, line and
     kind without the value.
10. MUST record that the PRD declares no Unreachable Acceptance and count no
    row in `rows_blocked_declared`.
11. MUST perform the glossary check (`docs/agents/domain.md`): record that the
    PRD declares `changes: **Settlement Check**`, `changes: **Secret Scan**`
    and `changes: **Cause Class**`, that `CONTEXT.md` carries each revision
    task_04 prescribes, and that the Task commit range changes no other
    glossary entry.
12. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`, `dist/`,
      `docs/references/skill-coverage.json`, any existing Archive Record, any
      Baseline asset, `docs/agents/`, `internal/secretscan`, `internal/store`,
      `internal/runcause/signatures.json`, any Roundfix reference other than
      `deliver`, `implement` and `runs`, any skill other than the Roundfix
      files the Tasks declare;
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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" bin/roundfix archive 0351-a-secret-refusal-of-a-task-commit-is-repairable-and-visible --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value or any matched value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named scan, Causes and park tests and check the history diffs.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A first-attempt Secret Match is repaired through Verification Feedback, a final one fails the Task with today's reason, a secret refusal is its own Causes item, and a run-unresolved park names its failed and skipped Tasks, in the tests.
- [ ] No value reaches any output, each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0351-a-secret-refusal-of-a-task-commit-is-repairable-and-visible/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-6; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-10; Invariants 1-6; Declared breaks; Testing Approach; Build Order 5
