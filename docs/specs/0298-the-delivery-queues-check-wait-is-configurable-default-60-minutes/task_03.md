---
task: task_03
spec: 0298-the-delivery-queues-check-wait-is-configurable-default-60-minutes
status: pending
type: qa
complexity: high
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the Delivery Queue's
Check Wait the User Config key `delivery.check_timeout`, a Go duration greater
than 0 and at most 24 hours, 60 minutes by default instead of a 30-minute
constant, ignored with a warning in Project Config, and makes the
`checks-timeout` park name the key and the value it waited, while each wait
keeps starting its own deadline.

Every row runs on temporary homes, the engine's fake clock and fake forges
under the Run's temporary directory. No row reads or writes the real
`~/.roundfix`, reaches the network, sends a request to a provider, the judge
or a reviewer, reads a credential, or reads another repository. The gate never
runs `roundfix spec judge`, never runs `roundfix review` against a real
provider, never runs `roundfix deliver` against a real forge, and never runs
an archive plan without `--offline`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1, 2 and 5:
   - Run `TestTheCheckWaitDefaultsToSixtyMinutes`,
     `TestTheCheckWaitUsesTheConfiguredDuration`,
     `TestEachCheckWaitStartsItsOwnDeadline` and
     `TestThePullRequestHeadWaitUsesTheCheckWait` in `./internal/delivery`
     and `TestTheDeliveryEngineReceivesTheConfiguredCheckWait` in
     `./internal/cli`, and record for each case the configured value, the
     fake-clock time at which it parked or failed, the poll count where the
     test asserts one, and the blocker text (Invariants 1, 4, 5 and 6).
   - Run `TestCheckRerunRestartsTheTimeout`,
     `TestAnOldCheckAttemptCannotKeepResettingTheTimeout` and
     `TestStaleRerunsAreBoundedPerTipAndRestartTheTimeoutOnlyOnce`, and
     record that a re-run still restarts the wait once and only once.
3. MUST verify Success Metrics 3 and 4:
   - Run `TestDeliveryCheckTimeoutIsReadFromUserConfig`,
     `TestProjectConfigCannotSetTheDeliveryCheckTimeout` and
     `TestDeliveryCheckTimeoutRefusesAValueOutsideItsBounds` in
     `./internal/config`, and record each value, the loaded duration or the
     refusal, and the warning text (Invariants 2 and 3).
   - With the built `bin/roundfix`, under a temporary `HOME` holding only a
     User Config, in a temporary Git repository with no remote, run
     `bin/roundfix deliver status` three times: with a User Config
     `delivery:` / `  check_timeout: 0`, with `check_timeout: 24h1s`, and
     with `check_timeout: 45m` plus a Project Config `.roundfixrc.yml`
     holding `check_timeout: 5m`. Record each exit code and standard error:
     the first two name `delivery.check_timeout must be greater than 0 and at
     most 24h`, and the third carries the line `config: delivery.check_timeout
     in Project Config is ignored; set delivery.check_timeout in User Config`
     whatever status then reports about the absent queue.
4. MUST, on temporary copies of the audited head made with `git archive`,
   apply one sabotage per copy and record that the named test fails on the
   named case:
   - make `defaultCheckTimeout` 30 minutes again:
     `TestTheCheckWaitDefaultsToSixtyMinutes` fails, because its engine sets
     no `CheckTimeout` and its checks stay pending;
   - make `checksTimeoutBlocker` omit `delivery.check_timeout`:
     `TestTheCheckWaitDefaultsToSixtyMinutes` fails on the blocker text;
   - make the Pull Request head wait in `Engine.createPullRequest` use a fixed
     30 minutes: `TestThePullRequestHeadWaitUsesTheCheckWait` fails, because
     it configures 45 minutes;
   - remove `CheckTimeout` from `commandDeliveryEngineDependencies`:
     `TestTheDeliveryEngineReceivesTheConfiguredCheckWait` fails on its
     45-minute case;
   - remove both the Project Config removal of `delivery.check_timeout` and
     the overlay's Project Config guard:
     `TestProjectConfigCannotSetTheDeliveryCheckTimeout` fails on its `5m`
     case;
   - let the validation accept `0`:
     `TestDeliveryCheckTimeoutRefusesAValueOutsideItsBounds` fails on its `0`
     and `0s` cases;
   - drop the `24h` upper bound:
     `TestDeliveryCheckTimeoutRefusesAValueOutsideItsBounds` fails on its
     `24h1s` case;
   - compute the `checking` deadline once per engine instead of once per
     call: `TestEachCheckWaitStartsItsOwnDeadline` fails on its second wait.
5. MUST execute the tests task_01 and task_02 name in their Verification.
6. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - run `git log -1 --format=%s 1f06abf1` and `git log -1 --format=%s
     c7bc381d` in the audited checkout and record both subjects, and run `git
     show c7bc381d:.github/workflows/ci-verify.yml` and record its
     `timeout-minutes: 60` line; then record that the default the tests of
     row 2 assert equals that job timeout;
   - because this Spec's authorization forbids network access, record
     GitHub's Actions limits documentation as `blocked (environment: network
     prohibited by Spec authorization)`, quoting the two passages the PRD
     cites. Such a row never decides a qualifying `partial`.
7. MUST record Success Metric 7 as exactly one row, `blocked (declared
   unreachable)`, against the PRD's single Unreachable Acceptance entry, and
   count it once in `rows_blocked_declared`: no row may reach a forge.
8. MUST perform the glossary and documentation checks:
   - record that `CONTEXT.md` defines **Check Wait** and links ADR-0311, that
     the glossary contract tests pass, and whether any other term the Spec
     introduced needs an entry (`docs/agents/domain.md`);
   - record the section and the row `docs/user-guide/configuration.md`
     gained, the sentences `docs/user-guide/commands/deliver.md` and both
     copies of the Roundfix `deliver` reference changed or gained, and that
     `TestTheSkillCoverageMapIsCurrent` passes;
   - search the guides, the agent guides and the skill copies for wording
     that states a 30-minute check wait or calls it "the existing checks
     timeout", and record what it finds.
9. MUST verify from the repository history:
   - each Task's changed files stay within its declarations or its `##
     Recorded paths`;
   - every Governed Path changed is bounded in `_authorization.md`;
   - these are byte-identical to the Delivery Base: the Makefile, the CI
     workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every Baseline module
     under `internal/baseline/assets/`, `docs/agents/setup-context.json`,
     every file under `docs/history/`, `internal/delivery/park_class.go`,
     `internal/store/`, `docs/references/skill-coverage.json`, and every
     existing `_test.go` file except the four task_02 declares;
   - the Roundfix skill version was raised and recorded;
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
    `bin/roundfix archive 0298-the-delivery-queues-check-wait-is-configurable-default-60-minutes --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's `## Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test and each sabotage on temporary copies of the audited head.
- [ ] Load the three User Config shapes through the built binary under a temporary home.
- [ ] Record the outside evidence, the declared row, the glossary and documentation checks and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] With nothing set every Check Wait is 60 minutes, a User Config value replaces it, and each removed guard makes its named test fail on its named case.
- [ ] A Project Config value is ignored with the warning, and out-of-bounds values are refused, through the built binary.
- [ ] The scope audit finds no undeclared Governed Path change and no change outside the declared files.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0298-the-delivery-queues-check-wait-is-configurable-default-60-minutes/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-6; Invariants 1-6; Testing Approach; Build Order 3
- ADR-0236
- ADR-0304
- ADR-0311
