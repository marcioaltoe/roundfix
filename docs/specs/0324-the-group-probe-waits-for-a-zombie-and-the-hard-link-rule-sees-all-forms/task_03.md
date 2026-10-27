---
task: task_03
spec: 0324-the-group-probe-waits-for-a-zombie-and-the-hard-link-rule-sees-all-forms
status: completed
type: qa
complexity: high
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the unreaped-group
test wait for a zombie member before its group probe, adds a test proving
the zombie reading rejects an exiting member, and widens the Hard Link rule
to every reference and every module test file (ADR-0324). Only test files
change.

Every row runs on temporary directories, temporary copies of the audited head
and of the Delivery Base made with `git archive`, compiled test binaries and
this repository's own documents. No row reads or writes `~/.roundfix`,
reaches the network, sends a request to a provider, the judge or a reviewer,
reads a credential, or reads another repository. Every `go test` carries an
explicit `-timeout` (or `-test.timeout` for a compiled binary), and every row
that runs `bin/roundfix` sets `HOME` to a fresh temporary directory, because
the CLI loads `$HOME/.roundfix/config.yml`. Evidence is written into the
checkout only after the suites and harness rounds end.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 2 and 3 and Invariants 1 to 3 for task_01:
   run `TestDetachFixtureGroupWithOnlyAnUnreapedMemberHasEnded`,
   `TestDetachFixtureGroupWithOnlyAnExitingMemberHasEnded` and
   `TestDetachFixtureZombieReadingRejectsAnExitingMember` with `-count=20
   -timeout 15m -v`, then the whole `./internal/cli` package once with
   `-count=1 -timeout 20m`, the shape that failed on 2026-10-25, and record
   that each passes; record from `git diff` against the Delivery Base that the
   unreaped-group test changed only its wait and the comment above it.
3. MUST verify Success Metric 1 with the load harness, on a macOS host only
   (record `blocked (environment: not macOS)` otherwise, never a pass):
   - in two temporary copies made outside the checkout, one of the audited
     head and one of the Delivery Base, compile each `internal/cli` test
     binary with `go test -c`;
   - from the copy's `internal/cli` directory, run three rounds of eight
     concurrent runs of the binary, each with
     `-test.run '^TestDetachFixtureGroupWithOnlyAnUnreapedMemberHasEnded$' -test.count=2000 -test.timeout 5m`,
     first every round of the audited head, then every round of the
     Delivery Base;
   - during each round keep eight CPU-bound loops (`sh -c 'while :; do :; done'`)
     running, record their PIDs when started, and end exactly those PIDs
     when the round ends, also on failure; never end a process by name;
   - for each round record the start and end time in BRT, the exit status of
     each of the eight runs, the number of lines containing `want EPERM`,
     and the 1-minute load average at the end;
   - the row passes when all 24 runs of the audited head exit 0 with no
     `want EPERM` line. Record the Delivery Base's failures as context; the
     audit addendum measured 87 in 48,000 at load averages 49 to 55, but they
     are not required.
4. MUST verify Success Metric 4 and Invariant 4 for task_02: run
   `TestNoTestCreatesAHardLink` with `-v` and record each subtest; record
   that `grep -rn --include='*_test.go' -e 'os[.]Link' -e 'unix[.]Link' -e 'syscall[.]Link' cmd internal skills`
   prints only lines of `internal/testfixture/hard_link_rule_test.go`, the
   rule's own test, whose seeded sources must contain those patterns for the
   rule to be tested.
5. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, and record that the named test fails, with the case
   whose input makes the mutation observable:
   - make the darwin `detachFixtureMemberIsZombie` also return true when
     `P_WEXIT` is set: `TestDetachFixtureZombieReadingRejectsAnExitingMember`
     fails naming an exiting state read as a zombie (its input is the
     observed exiting window);
   - change the Darwin expectation of the unreaped-group test from `EPERM`
     to success: that test fails, because a zombie-only group probe returns
     `EPERM`;
   - stop recording dot imports in `hardLinkViolations`: the subtests
     `dot import of os`, `dot import of unix Linkat` and
     `dot-imported syscall Link as a value` fail;
   - report a selector only when it is a call's function: the subtest
     `os Link as a function value` fails;
   - walk `internal/` only: the subtests `a test outside internal` and
     `a test under skills` fail;
   - stop excluding struct literal keys: the subtest
     `dot import beside Link fields methods keys and labels` fails.
6. MUST execute the commands task_01 and task_02 name in their Verification.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - because this Spec's authorization forbids network access, record the
     XNU `kern_sig.c` source as `blocked (environment: network prohibited by
     Spec authorization)`, quoting the two lines the PRD cites, and rest the
     row on their copy in the 2026-10-26 addendum of
     `docs/references/2026-10-24-test-suite-audit.md`; such a row never
     decides a qualifying `partial`;
   - quote the Spec 0317 QA failure line from the adopted Backlog Entry
     under `references/`.
8. MUST record the PRD's single Unreachable Acceptance entry as exactly one
   row, `blocked (declared unreachable)`, counted as one in
   `rows_blocked_declared`: no row may read a later Run's log.
9. MUST perform the glossary check: record that the Spec declares no term,
   that `CONTEXT.md` is byte-identical to the Delivery Base, and that no
   bolded phrase of the Spec's artifacts lacks a glossary entry
   (`docs/agents/domain.md`).
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its
      `## Recorded paths`;
    - no Governed Path changed (`_authorization.md` has `paths: []`);
    - these are byte-identical to the Delivery Base: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `CONTEXT.md`, every
      non-test `.go` file, every file under `docs/agents/`,
      `docs/references/`, `docs/user-guide/`, `.agents/` and `skills/`,
      every `internal/cli` file except the four task_01 declares, and every
      `internal/testfixture` file except `hard_link_rule_test.go`;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance, supported by the pre-PR equivalent of each
    control: approval from the recorded delivery authority; checks from the
    Daemon's repository Verification at the audited head; no unresolved
    threads, because no Pull Request exists; no Merge-Ready acceptance yet;
    the audited head named as the claimed candidate.
12. MUST, for the Outcome for the Archive Record, run only
    `HOME="$(mktemp -d)" bin/roundfix archive 0324-the-group-probe-waits-for-a-zombie-and-the-hard-link-rule-sees-all-forms --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's `## Outcome` as one paragraph that does not end with a colon.
    It MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite or harness round runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests, the load harness and each sabotage on temporary copies.
- [ ] Record the outside sources, the declared row, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The detach tests pass, and 24 loaded runs of the audited head, 48,000 runs of the unreaped-group test, pass.
- [ ] The Hard Link rule passes on the repository, and each sabotage makes its named test fail.
- [ ] The scope audit finds no change outside the declared files.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0324-the-group-probe-waits-for-a-zombie-and-the-hard-link-rule-sees-all-forms/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-3; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-3; Invariants 1-4; Testing Approach; Build Order 3
