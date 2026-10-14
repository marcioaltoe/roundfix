---
task: task_05
spec: 0269-a-killed-child-names-its-signal-and-a-cancelled-run-starts-nothing
status: completed
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec reports a child that a signal
ended as a Signal Exit in the Delivery Queue and in every test helper that runs
a child Roundfix binary, makes a cancelled Run start no Verification command
after an immediate capacity grant, and documents both. Every row runs against
the audited head, a built binary, temporary homes and disposable copies of the
tree. No row reaches a provider or the network, starts an Agent Session, reads
a credential, sends a signal to a process the row did not start, uses `sudo`,
or reads or writes the real `~/.roundfix`.

## Requirements

1. MUST run the repository Verification with `bin/roundfix` rebuilt from the
   audited head first, and record its result as a gate fact. When any child in
   that run reports a Signal Exit, the report records the test, the named
   signal and the core state as an observation for the Backlog Entry
   `2026-10-13-trace-the-sender-of-the-f6-signal.md`; F6 recurring is not a
   failure of this Spec.
2. MUST measure Success Metric 1:
   `TestRoundfixBinaryHelperNamesTheSignalThatKilledItsChild` passes at the
   audited head, and on a disposable copy of the tree where
   `runRoundfixBinaryMacro` again returns `exitCodeFromWait(err)` it fails,
   reporting code `1` and no Signal Exit line. Record both outputs.
3. MUST measure Success Metric 2: count the `exitCodeFromWait(` calls in
   `_test.go` files of the audited head (0 expected) and record that each of
   the nine sites `_techspec.md` → Current behavior lists calls
   `testfixture.ChildExitCode`.
4. MUST measure Success Metric 3 on a disposable copy of the audited tree:
   insert, as the first statement of `acquireVerificationCapacity`, a 50 ms
   `time.Sleep` that runs only when the `*fairVerificationGate` behind
   `plan.verificationGate` has `activeShared > 0` (read under its mutex); run
   `TestTaskCycleVerificationCapacityCancellationWhileQueuedStartsNoCommandOrSettlement`
   with `-count=20` and record 0 failures; then remove only the post-grant
   cancellation check task_04 added and record 20 failures of 20, each with
   `total starts=2`. Also record that the two new tests of task_04 fail on a
   copy without that check and pass at the audited head.
5. MUST verify API Contracts 1-2 and Invariant 4: the two signaled delivery
   tests of task_02 pass and their asserted lines match the contracts; on a
   disposable copy where `stepExecutable` prints `probe.exitCode` again, the
   probe test fails naming `exited -1`; and
   `internal/cli/deliver_item_binary_test.go` and
   `internal/cli/deliver_budget_park_test.go` are byte-identical to the
   delivery base, the merge base of the audited head with `main`, and pass.
6. MUST verify Invariant 1 by running `TestExitDescribesHowAChildEnded` and
   recording `GOOS=windows go build -buildvcs=false` and
   `GOOS=linux go build -buildvcs=false` of `./internal/childexit`,
   `./internal/testfixture` and `./internal/cli` as exit 0.
7. MUST verify the documentation: the `deliver` reference, its mirror and the
   `deliver` user guide carry the `<outcome>` placeholder and both Signal Exit
   forms; `make skills-sync-check`, the owned-skill version test and
   `TestTheSkillCoverageMapIsCurrent` pass; the new Backlog Entry is open,
   follows the Backlog Operational Contract and links no Spec folder.
8. MUST record, as evidence this Spec did not author, the passages
   `_prd.md` → Acceptance evidence quotes, read without the network: the
   `os.ProcessState.ExitCode` passage from `go doc os.ProcessState.ExitCode`,
   the Select passage from the toolchain's `doc/go_spec.html` under
   `go env GOROOT`, and the two `x/sync` v0.22.0 passages from
   `semaphore/semaphore.go` under `go env GOMODCACHE` when that module is in
   the cache, otherwise as quoted in `_prd.md`. For each it records whether the
   passage still supports the behavior.
9. MUST perform the glossary check of `docs/agents/domain.md`: `CONTEXT.md`
   defines **Signal Exit**, links no Spec, and no other term the Spec
   introduced is missing.
10. MUST verify from the repository history that each Task's changed files
    stay within its declarations or its `## Recorded paths`, that every
    Governed Path changed is bounded in `_authorization.md`, that
    `internal/cli/baseline_release_gate_test.go` changed only inside
    `runBaselineReleaseCLI`, that `internal/daemon/task_engine_test.go` and
    the production `exitCodeFromWait` are byte-identical to the delivery base,
    and that the Makefile, the CI workflows, `go.mod`, `go.sum` and
    `.roundfixrc.yml` did not change. This row reads Task commits, so it
    declares a `commit_range` input.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance, supported by the pre-PR equivalent of each control: the
    maintainer's recorded delivery authority, the Daemon's repository
    Verification at the audited head, no unresolved threads because no Pull
    Request exists, no Merge-Ready acceptance yet, and the audited head named
    as the claimed candidate.
12. MUST NOT accept a row whose only evidence is that a file was read, MUST
    NOT persist an environment value or an environment dump, and MUST NOT
    write outside temporary directories and the Spec's `qa/` folder.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Measure Success Metrics 1-3 with their disposable copies.
- [ ] Verify the delivery contracts, the package builds and the documentation.
- [ ] Record the outside evidence, the glossary check and the scope audit.
- [ ] Write the QA report with the Pull Request row blocked.

## Acceptance Criteria

- [ ] Each Success Metric holds at the audited head and its sabotaged copy fails as stated.
- [ ] The delivery contracts hold and the exited forms are unchanged.
- [ ] The scope audit finds no undeclared Governed Path.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0269-a-killed-child-names-its-signal-and-a-cancelled-run-starts-nothing/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- [_prd.md](_prd.md) — Goals; User Stories 1-3; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Acceptance evidence
- [_techspec.md](_techspec.md) — API Contract 1; API Contract 2; API Contract 3; Invariants 1-5; Testing Approach; Build Order 5
- ADR-0056; ADR-0266; ADR-0271; ADR-0259; ADR-0256; ADR-0189; ADR-0233; ADR-0244
