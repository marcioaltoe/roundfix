---
task: task_03
spec: 0311-the-suite-audit-and-no-test-hard-links-its-binary
status: completed
type: qa
complexity: high
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the `internal/agent`
package adapters copies of the test binary instead of hard links, keeps the
re-executed fixture children from making any, makes the fake acpx helper name
the `RunPrompt` error, and adds the Hard Link rule with its guide line and
the **Single-Name Test Binary** glossary term (ADR-0324).

Every row runs on temporary directories, temporary copies of the audited head
and of the Delivery Base, compiled test binaries and this repository's own
documents. No row reads or writes `~/.roundfix`, reaches the network, sends a
request to a provider, the judge or a reviewer, reads a credential, or reads
another repository. The gate never runs `roundfix spec judge`, never runs
`roundfix review` against a real provider, and never runs an archive plan
without `--offline`. Every `go test` it runs carries an explicit `-timeout`
(or `-test.timeout` for a compiled binary), and every row that runs
`bin/roundfix` sets `HOME` to a fresh temporary directory, because the CLI
loads `$HOME/.roundfix/config.yml`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 2 and Invariants 1, 2 and 4 for task_01: run
   `TestPackageAdaptersAreCopiesOfTheTestBinary`,
   `TestASpawnedFixtureChildLeavesNoPackageAdapter`,
   `TestSpawnedFixtureDeathFailsFast`, `TestFixtureBinarySurvivesConcurrentExec`,
   `TestCheckAdapterProvesOfficialClaudePackageAndVersion` and every test of
   `./internal/agent` with `-count=1 -timeout 15m`, and record that each
   passes; record that `grep -rn 'os[.]Link(' internal` prints only lines of
   `internal/testfixture/hard_link_rule_test.go`, the rule's own test, whose
   seeded source strings (the forbidden call and the shadowed receiver) must
   contain the pattern for the rule to be tested.
3. MUST verify Success Metric 1 with the load harness, on a macOS host
   only (record `blocked (environment: not macOS)` otherwise, never a pass):
   - in two temporary copies made outside the checkout, one of the
     Delivery Base and one of the audited head, compile the
     `internal/agent` test binary of each with `go test -c`;
   - from the copy's `internal/agent` directory, run eight rounds of twelve
     concurrent runs of the binary with `-test.count=3 -test.parallel 16
     -test.timeout 8m`, first all rounds of the audited head, then all rounds
     of the Delivery Base;
   - for each round record the start and end time in BRT, the number of
     runs that exited non-zero, the 1-minute load average, and the number of
     `syspolicyd` lines "Terminating process due to Malware rejection" that
     `/usr/bin/log show` reports in that window (the absolute path matters:
     `log` is a shell builtin in zsh);
   - the row passes when all 96 runs of the audited head exit 0. Record the
     Delivery Base's failing rounds as context; they are expected (four of
     eight on 2026-10-24) but not required.
4. MUST verify Success Metric 3 and Invariant 3 for task_02: run
   `TestNoTestCreatesAHardLink` with `-v` and record its repository and
   seeded subtests; record the line task_02 adds to
   `docs/agents/specific-repository.md` and the **Single-Name Test Binary**
   entry in `CONTEXT.md`.
5. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, and record that the named test fails, with the case
   whose input makes the mutation observable:
   - replace the copy in `provisionPackageAdapterLinks` by
     `os.Link(os.Args[0], target)`:
     `TestPackageAdaptersAreCopiesOfTheTestBinary` fails on `os.SameFile`
     for a package adapter target, and `TestNoTestCreatesAHardLink/repository` fails
     naming `internal/agent/acpx_runner_test.go` and the line;
   - remove the `ROUNDFIX_TEST_SPAWNED_FIXTURE_DEATH` check before the
     provisioning in `TestMain`:
     `TestASpawnedFixtureChildLeavesNoPackageAdapter` fails, finding one
     `roundfix-agent-package-adapters-` root left by the fixture child that
     ended through `os.Exit(23)`;
   - add `_ = syscall.Link("a", "b")` to any `internal` test file that
     imports `syscall`, or import it: `TestNoTestCreatesAHardLink/repository`
     fails naming that file and line;
   - make the rule ignore aliased imports:
     the seeded aliased-import subtest of `TestNoTestCreatesAHardLink` fails.
6. MUST execute the tests task_01 and task_02 name in their Verification.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - on a macOS host, run `/usr/bin/log show --start '2026-10-23 23:10:00'
     --end '2026-10-23 23:10:20' --style compact --predicate 'process ==
     "syspolicyd" AND eventMessage CONTAINS "Malware rejection"'` and record
     whether the line for pid 14242 is still retained; the unified log keeps
     a few days, so a missing line is recorded as `blocked (environment: log
     no longer retained)` and the row then rests on the copy in
     `docs/references/2026-10-24-test-suite-audit.md`, which the planning
     change committed before any Task ran;
   - because this Spec's authorization forbids network access, record the
     Apple Developer Forums thread as `blocked (environment: network
     prohibited by Spec authorization)`, quoting the two sentences the PRD
     cites. Such a row never decides a qualifying `partial`.
8. MUST record Success Metrics 4 and 5 as exactly two rows, each `blocked
   (declared unreachable)` against its own entry of the PRD's Unreachable
   Acceptance, and count them as two in `rows_blocked_declared`: no row may
   read a CI run or a Run's Verification log.
9. MUST perform the glossary check: record that `CONTEXT.md` carries
   **Single-Name Test Binary**, that its definition cites ADR-0324, and
   that no other term changed (`docs/agents/domain.md`).
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - these are byte-identical to the Delivery Base: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every non-test `.go`
      file, every file under `internal/testfixture/` except the new
      `hard_link_rule_test.go`, every `internal/agent` test file except
      `acpx_runner_test.go` and the new `package_adapter_test.go`, and every
      file under `.agents/`, `skills/` and `docs/user-guide/`;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
12. MUST, for the Outcome for the Archive Record, run only
    `HOME="$(mktemp -d)" bin/roundfix archive 0311-the-suite-audit-and-no-test-hard-links-its-binary --plan --offline`
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
- [ ] Run each named test, the load harness and each sabotage on temporary copies.
- [ ] Record the outside sources, the declared rows, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The package adapters are copies, a fixture child leaves nothing, and 96 loaded runs of the audited head pass.
- [ ] The Hard Link rule passes on the repository, and each sabotage makes its named test fail.
- [ ] The scope audit finds no change outside the declared files.

## Verification



The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0311-the-suite-audit-and-no-test-hard-links-its-binary/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-2; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-2; Invariants 1-4; Testing Approach; Build Order 3
