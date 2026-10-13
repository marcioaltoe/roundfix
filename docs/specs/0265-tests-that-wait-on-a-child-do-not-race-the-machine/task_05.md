---
task: task_05
spec: 0265-tests-that-wait-on-a-child-do-not-race-the-machine
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes every test that waits on a
child process, a goroutine or Git bound by real time only the budget it
asserts, adds two dependencies with production defaults, and adds the
Real-Time Budget rule. Its central rows run the fixed tests and three full
suites while synthetic CPU load runs. No row reads or writes the real
`~/.roundfix`, signals a process the row did not start, uses `pkill` or
`killall`, reaches a provider or reads a credential.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST run synthetic load only this way: `2 × hw.ncpu` CPU burners started
   under `nice -n 19`, each PID recorded when it starts, stopped by those PIDs
   through a shell `trap` on exit, interrupt and termination, so no burner
   outlives the window it covers; a window covers one measured command and
   nothing else. It records the load average at each window's start and end.
3. MUST measure Success Metric 1: inside one load window per package, run each
   test named in rows 1 to 26 of `_techspec.md` → Current behavior (for row 18,
   its five named tests) with `go test -race -count=50 -parallel 16 -run` on its
   package, and record for each package the exit status and the count of
   `--- FAIL` lines. The row passes only when every run exits 0 with no failure.
4. MUST measure Success Metric 2: run `go test -count=1 -parallel 16 ./...`
   three consecutive times, each inside its own load window, recording start
   time, load averages, wall time, exit status and the count of lines matching
   `signal: killed` or `exited -1`. While the audited head's
   `internal/daemon/task_engine_test.go` still contains the text
   `Sequential: shares process-wide wall-clock scheduling`, Spec 0259 has not
   merged, and the three suites pass
   `-skip '^(TestBudgetExceededRunRecordsItsReason|TestRemoveRegisteredWorktreeWaitsForTheAdminLock|TestAddDetachedWorktreeWaitsForTheAdminLock)$'`
   and record that skip; otherwise they pass no skip. The row passes only when
   all three exit 0 with every count zero. On the first failing suite it stops,
   keeps the full output, records the failing test and its diagnostic, and fails
   the row; it never reruns a failed suite to obtain a pass.
5. MUST verify API Contract 7 and Success Metric 3 by executing task_04's two
   rule tests, and by seeding, in a disposable copy of the repository only, one
   `context.WithTimeout(ctx, 500*time.Millisecond)` in a new test file of
   `internal/cli` and recording that the rule fails naming that file and line.
6. MUST verify Success Metric 4 and Invariant 1 from the diff of every non-test
   Go file between the delivery base and the audited head: only
   `internal/agent/acpx_runner.go`, `fallback.go`, `selection_assignment.go`,
   `sealed.go`, `internal/daemon/engine.go` and `task_engine.go` changed, the
   5 s cleanup constant and every production duration literal are unchanged, and
   the two new dependencies default to the values used before. It also executes
   `TestDisposableCleanupTimeoutDefaultsToTheProductionValue` and
   `TestTaskBudgetWatchdogUsesTheRealTimerByDefault`.
7. MUST verify API Contract 1 and Invariant 2 by comparing every assertion of
   rows 1 to 8 between the delivery base and the audited head: each expected
   exit status, stdout value and diagnostic substring is byte-identical.
8. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the RAFT preprint and the async-wait paper.
   For each source it reaches it records the quoted passage and whether it still
   supports the design. For a source it cannot reach it records
   `blocked (environment: network denied: <host>)`.
9. MUST perform the glossary check of `docs/agents/domain.md`. It records that
   `CONTEXT.md` defines **Real-Time Budget** linked to ADR-0271, and whether any
   other term the Spec introduced needs one.
10. MUST verify from the repository history that each Task's changed files stay
    within its declarations or its `## Recorded paths`; that the only Governed
    Path changed is `docs/agents/specific-repository.md`, bounded in
    `_authorization.md`; that the Makefile, the CI workflows, `go.mod`, `go.sum`,
    every skill, every Baseline module and
    `internal/testfixture/parallel_tests_test.go` did not change; and that no
    region the TechSpec's 0259 overlap table assigns to Spec 0259 changed. This
    row reads Task commits, so it declares a `commit_range` input.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance, and support it with the pre-PR equivalent of each
    control: approval from the recorded delivery authority, checks from the
    Daemon's repository Verification at the audited head, no unresolved
    threads because no Pull Request exists, no Merge-Ready acceptance yet, and
    the audited head named as the claimed candidate.
12. MUST NOT accept a row whose only evidence is that a file was read. It MUST
    NOT read or write `~/.roundfix`, signal a process it did not start, use a
    broad `pkill` or `killall`, leave a burner running after its window, send
    any request to a provider, or start a workflow.
13. Addendum (2026-10-13, after QA finding F-SUITE of Run
    `run_20261013T045342Z_a78df1f6c39a67e2`): while Spec 0259 has not merged
    (Requirement 4's condition), the skip of Requirement 4 also names
    `TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate`, the test that
    saw F6, which Spec 0259 owns (its task_02 Requirement 8) and whose
    `exited -1` is F6 itself. The skip pattern becomes
    `-skip '^(TestBudgetExceededRunRecordsItsReason|TestRemoveRegisteredWorktreeWaitsForTheAdminLock|TestAddDetachedWorktreeWaitsForTheAdminLock|TestDeliveryStepFallsBackWhenTheItemBinaryWouldMigrate)$'`;
    once 0259 has merged, no skip is passed.
14. Addendum (same QA): when any measured command runs with a temporary
    `HOME`, it MUST set `RUSTUP_HOME` and `CARGO_HOME` to the real rustup and
    cargo homes, so the Rust toolchain the Guidance composition journey
    invokes stays configured; it records the two values it set.
15. Addendum (same QA): Requirement 3's row for `internal/agent` runs row 19
    with `ROUNDFIX_REAL_ACPX=1` after task_06; a skip of that test fails the
    row.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each fixed test 50 times under load, package by package.
- [ ] Run the three full suites under load and record each.
- [ ] Execute the rule tests, the seeded violation and the production diff audit.
- [ ] Record the outside evidence, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Every fixed test passes 50 consecutive runs with `-race -parallel 16` under load.
- [ ] Three consecutive full suites exit 0 under load with no `signal: killed` or `exited -1` line.
- [ ] The seeded short budget fails the rule with its file and line.
- [ ] No production budget value changed, and the scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0265-tests-that-wait-on-a-child-do-not-race-the-machine/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; API Contract 7; Invariants 1-5; Current behavior; 0259 overlap; Build Order 5
- ADR-0271; ADR-0259; ADR-0213; ADR-0266
