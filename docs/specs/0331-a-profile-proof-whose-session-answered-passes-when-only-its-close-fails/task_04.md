---
task: task_04
spec: 0331-a-profile-proof-whose-session-answered-passes-when-only-its-close-fails
status: pending
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes an Exact Agent
Selection Proof decide on its setup alone, closes its disposable Agent
Session under the attempt's Proof Setup Deadline, records an Unclosed Proof
Session that the next proof closes, names a timed-out setup step, stops the
CLI from retrying a close deadline as a setup timeout, reports the warning in
`profiles validate`, and keeps Baseline semantic analysis's stop.

Every row runs on this repository's own files, the compiled fake acpx of
`internal/agent`, fake provers and fake runtimes, and temporary copies of the
audited head. No row starts acpx or an ACP adapter, reads or writes the real
`~/.roundfix` or `~/.acpx`, reads another repository or `~/dev/skills`,
reaches the network, sends a request to a provider, the judge or a reviewer,
or reads a credential. Every `go test` it runs carries an explicit
`-timeout`, and every row that runs `bin/roundfix` sets `HOME` and `TMPDIR` to
fresh temporary directories, because the CLI loads `$HOME/.roundfix/config.yml`,
and removes them afterwards. The QA never writes into the checkout while a
test suite runs. No row compares one command's output with another's byte for
byte.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 and 2 and Invariants 1, 3 and 5 by running
   `go test -count=1 -timeout 20m -v` over `./internal/agent`,
   `./internal/baselineacp` and `./internal/testfixture` with the test names
   task_02 names in its Verification, and recording each pass.
3. MUST verify Success Metrics 3, 4 and 5 and Invariants 2, 4 and 6 by
   running `go test -count=1 -timeout 20m -v` over `./internal/agent` and
   `./internal/cli` with the test names task_03 names in its Verification,
   and recording each pass.
4. MUST verify through the built binary that the help and the JSON schema
   name stay put: `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t"
   bin/roundfix profiles validate --help; status=$?; rm -rf "$h" "$t"`
   exits `0` and prints the `--json` flag, and that
   `TestProfilesValidateReportsAnUnclosedProofSession` asserts the schema
   `roundfix/profiles-validate/v1` with the `warning` field.
5. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that
   the named test fails, with the case whose input makes the mutation
   observable:
   - let a close error after a passing setup fail `ProveExactSelection`
     again: `TestAProofWhoseSessionAnsweredPassesWhenItsCloseHangs` fails on
     its non-nil error;
   - close under `acpxPreflightCleanupTimeout` instead of the setup deadline:
     `TestTheProofCloseDeadlineFollowsTheSetupDeadline` fails because its
     warning names `5s`, not `2h0m0s`;
   - return the proof when the setup failed and the close failed:
     `TestASessionThatNeverAnsweredStillFailsWhenItsCloseFails` fails on the
     missing `*SelectionRejectedError`;
   - drop the step wrapping: `TestATimedOutProofStepIsNamed` fails on
     `profile proof step "read runtime catalogue"`;
   - skip writing the record: `TestAProofWhoseSessionAnsweredPassesWhenItsCloseHangs`
     fails on the missing record file;
   - continue the reaping pass after a failed close:
     `TestAFailedReapKeepsTheRecordAndStops` fails because the newer
     record's close is invoked;
   - delete a record whose close failed:
     `TestAFailedReapKeepsTheRecordAndStops` fails on the missing older
     record;
   - restore `profileProofSetupTimedOut` to a bare
     `errors.Is(err, context.DeadlineExceeded)`:
     `TestACloseDeadlineIsNotRetriedAsASetupTimeout` fails on three calls;
   - remove the `Unclosed` check from `Analyzer.proveSelection`:
     `TestAProofThatLeftItsSessionOpenStopsSemanticAnalysis` fails because a
     sealed prompt starts;
   - drop the `warning` line from `printProfilesValidateSuccess`:
     `TestProfilesValidateReportsAnUnclosedProofSession` fails on its text
     case.
6. MUST execute the tests task_01 to task_03 name in their Verification.
7. MUST search the repository, excluding this Spec's folder and
   `docs/history/`, for `closes the Session successfully` and for `The
   cleanup deadline, the Fallback Chain probe after Run creation and Run
   sessions keep their existing deadlines`, and record that neither appears
   in `CONTEXT.md`, `docs/user-guide/` or a skill copy.
8. MUST record, as the outside-evidence row, the published acpx sessions
   guide passages `_prd.md` → Acceptance evidence quotes (openclaw/acpx
   `docs/sessions.md`, `https://acpx.sh/sessions.html`, retrieved
   2026-10-27): "`sessions close` does not delete anything. It marks the
   record `closed: true` with `closedAt`, asks any active queue owner to send
   ACP `session/close`, and tears down adapter processes.", "`sessions prune`
   is the explicit way to delete closed records." and "After the queue
   drains, the owner stays alive for an idle TTL (default `300s`)", and the
   Vortex report quoted in the adopted Backlog Entry under `references/`
   (ai-memory message `01a17b18`, 2026-10-26). Compare offline the reaping
   close a passing `TestTheNextProofClosesARecordedSessionFirst` asserts
   (`sessions close <recorded session>` in the recorded working directory)
   with that contract. The row is never blocked; name openclaw/acpx and the
   message id with its date as the source.
9. MUST record Success Metric 6 as exactly one row, `blocked (declared
   unreachable)`, against the PRD's single Unreachable Acceptance entry, and
   count it once in `rows_blocked_declared`.
10. MUST perform the glossary check (`docs/agents/domain.md`): record that
    `CONTEXT.md` defines **Unclosed Proof Session** and that **Exact Agent
    Selection Proof** and **Proof Setup Deadline** cite ADR-0344, each
    citing no Spec.
11. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      anything under `internal/baseline/assets/` or `docs/agents/`, any
      skill other than the Roundfix files task_01 declares, and
      `internal/agent/sealed.go`;
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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" bin/roundfix archive 0331-a-profile-proof-whose-session-answered-passes-when-only-its-close-fails --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's ` Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix` or `~/.acpx`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests and the binary row.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside source, the declared row, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A proof whose setup answered passes with a warning and a record when its close hangs or fails; a rejected setup still fails.
- [ ] A recorded session is closed by the next runner's first proof; a failed close keeps the record.
- [ ] The CLI never retries a close deadline and names a timed-out step.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification



The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0331-a-profile-proof-whose-session-answered-passes-when-only-its-close-fails/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-5; User Experience; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-8; Invariants 1-6; Testing Approach; Build Order 4
