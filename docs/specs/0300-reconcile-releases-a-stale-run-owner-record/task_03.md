---
task: task_03
spec: 0300-reconcile-releases-a-stale-run-owner-record
status: completed
type: qa
complexity: high
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the owner proof tell a
proven start identity mismatch from an unreadable identity, makes the
Reconcile Command report a terminal Run whose owner PID now belongs to another
process as a Stale Run Owner Record, and makes `--apply` release that record
without sending any signal, while every identity it cannot read stays
preserved and Force Stop is unchanged.

Every row runs on temporary homes, repositories and Git configuration under
the Run's temporary directory, with fake process tables or processes the test
itself started. No row reads or writes the real `~/.roundfix`, runs
`roundfix reconcile` outside a temporary Roundfix Home, signals a process the
row did not start, reaches the network, sends a request to a provider, the
judge or a reviewer, reads a credential, or reads another repository. The
gate never runs `roundfix spec judge`, never runs `roundfix review` against a
real provider, and never runs an archive plan without `--offline`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1, 2 and 3:
   - Run `TestRunReconcileReleasesAStaleRunOwnerRecordWithoutASignal` and
     record the text heading, the proof, the dry-run action, the summary
     substring, the JSON `kind`, `processIds` and `runOwnerRecordsReleased`
     values, the apply action, that no termination was requested, that the
     owner PID and identity were cleared, and that the second `--apply`
     reported nothing (Invariants 1, 4 and 5).
   - Run `TestRunReconcilePreservesARunOwnerRecordWhoseIdentityIsNotProvenStale`
     and record, per subtest, the refusal reason it kept and that the record
     was kept (Invariant 2).
   - Run `TestRunReconcileOffersAndAppliesOwnedProcessTreesAndRunBranches`,
     `TestRunReconcileApplyPreservesProcessWhenTerminationIsUnproven` and
     `TestReconcileListsAndReleasesAnItemBranchOfAMergedSpec`, and record that
     process trees, Run Branches, item branches and the pinned summary
     substrings are unchanged.
3. MUST verify Success Metrics 4, 5 and 6:
   - Run `TestOwnerProcessControllerTellsAProvenIdentityMismatchFromAnUnreadableIdentity`,
     `TestInspectTreeReportsAProvenMismatchOfALiveProcessWithoutSignallingIt`
     and `TestReleaseStaleRunOwnerRecordClearsOnlyTheRecordedTerminalOwner` in
     `./internal/store`, and record each row's sentinel matches, the exact
     mismatch message, that the live helper was still alive, and the release
     results and journal events (Invariants 3, 4 and 6).
   - Run `TestOwnerProcessControllerProveOwnerClassifiesIdentityEvidence`,
     `TestOwnerProcessControllerRefusesMismatchedOwnerIdentity`,
     `TestOwnerProcessControllerProveOwnerRefusesMismatchedIdentity`,
     `TestRunForceStopOwnerIdentityUnreadableFlagRequiresUnreadableProof` and
     `TestDoctorReportsProcessResidueUnreadable`, and record that Force Stop
     and `doctor` are unchanged.
   - On temporary copies of the audited head made with `git archive`, apply
     one sabotage per copy and record that the named test fails, with the
     case whose input makes the mutation observable:
     - make `inspectReconcileProcesses` report every error matching
       `ErrOwnerProcessIdentityUnproven` as a `runOwnerRecord`:
       `TestRunReconcilePreservesARunOwnerRecordWhoseIdentityIsNotProvenStale`
       fails in its bare unproven subtest;
     - make it also report an error matching `ErrOwnerIdentityUnreadable`:
       the same test fails in its read failure and uncomparable token
       subtests;
     - make `applyReconcileProcesses` pass `runOwnerRecord` entries to
       `TerminateTreeAndWait`:
       `TestRunReconcileReleasesAStaleRunOwnerRecordWithoutASignal` fails on
       its termination record;
     - make `proveOwner` return `ErrOwnerProcessIdentityMismatch` for the own
       PID and a non-positive PID:
       `TestOwnerProcessControllerTellsAProvenIdentityMismatchFromAnUnreadableIdentity`
       fails in its own-PID and PID 0 rows;
     - make it return the mismatch for a token without the platform's prefix:
       the same test fails in its foreign platform token row;
     - make `TerminateAndWait` continue past a proven mismatch: the same test
       fails in its comparable mismatch row, and
       `TestOwnerProcessControllerRefusesMismatchedOwnerIdentity` fails;
     - make `InspectTree` enumerate the tree after a proven mismatch: the same
       test fails in its comparable mismatch row, and
       `TestInspectTreeReportsAProvenMismatchOfALiveProcessWithoutSignallingIt`
       fails;
     - drop `AND owner_identity = ?` from the release's update:
       `TestReleaseStaleRunOwnerRecordClearsOnlyTheRecordedTerminalOwner`
       fails on its different-identity call;
     - drop the release's terminal-state refusal: the same test fails on its
       Active Run call;
     - drop the release's journal event: the same test fails on its event
       count.
4. MUST execute the tests task_01 and task_02 name in their Verification.
5. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - read the adopted Backlog Entry under `references/` in the audited
     checkout, which records the maintainer's `roundfix reconcile` output of
     2026-10-21, and record that its refusal text, with whitespace runs collapsed to one space, `owner process identity
     cannot be proven: live process start identity does not match the
     recorded owner identity` is a substring of the exact message the
     comparable mismatch row of
     `TestOwnerProcessControllerTellsAProvenIdentityMismatchFromAnUnreadableIdentity`
     asserts;
   - because this Spec's authorization forbids network access, record
     POSIX.1-2024 Base Definitions 4.17 and proc_pid_stat(5) as `blocked
     (environment: network prohibited by Spec authorization)`, quoting the
     passages the PRD cites. Such a row never decides a qualifying `partial`.
6. MUST record Success Metric 7 as exactly one row, `blocked (declared
   unreachable)`, against the PRD's single Unreachable Acceptance entry, and
   count it once in `rows_blocked_declared`: no row may read the real
   `~/.roundfix`.
7. MUST perform the glossary and documentation checks:
   - record that `CONTEXT.md` defines **Stale Run Owner Record**, that
     **Reconcile Command** gained the release sentence, that both link
     ADR-0313, that the glossary contract tests pass, and whether any other
     term the Spec introduced needs an entry (`docs/agents/domain.md`);
   - record the paragraph `docs/user-guide/commands/reconcile.md` gained, the
     bullet and sentence both copies of the Roundfix `reconcile` reference
     gained, and that `TestTheSkillCoverageMapIsCurrent` passes;
   - search the guides, the agent guides and the skill copies for wording
     that says every `processCandidates` entry names live processes or that
     a mismatched owner identity is only ever preserved, and record what it
     finds.
8. MUST verify from the repository history:
   - each Task's changed files stay within its declarations or its `##
     Recorded paths`;
   - every Governed Path changed is bounded in `_authorization.md`;
   - these are byte-identical to the Delivery Base: the Makefile, the CI
     workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every Baseline module
     under `internal/baseline/assets/`, `docs/agents/setup-context.json`,
     every file under `docs/history/`, `internal/store/store.go`,
     `internal/store/delivery.go`, `internal/store/process_signal_unix.go`,
     `internal/cli/cli.go`, `internal/cli/doctor.go`,
     `internal/cli/deliver.go`, every `_test.go` file under `internal/store`
     and `internal/cli` that exists at the Delivery Base, and
     `docs/references/skill-coverage.json`;
   - the Roundfix skill's version was raised and recorded;
   - no `### QA settlement` section changed.

   This row reads Task commits, so it declares a `commit_range` input.
9. MUST record the non-waivable Pull Request row as
   `blocked (environment: no open Pull Request)`, naming the Pull Request row
   in its provenance. It supports the row with the pre-PR equivalent of each
   control:
   - approval from the recorded delivery authority;
   - checks from the Daemon's repository Verification at the audited head;
   - no unresolved threads, because no Pull Request exists;
   - no Merge-Ready acceptance yet;
   - the audited head named as the claimed candidate.
10. MUST, for the Outcome for the Archive Record, run only
    `bin/roundfix archive 0300-reconcile-releases-a-stale-run-owner-record --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's `## Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - signal a process the row did not start;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test and each sabotage on temporary copies of the audited head.
- [ ] Compare the Backlog Entry's refusal text with the asserted mismatch message, and record the blocked outside sources.
- [ ] Record the declared row, the glossary and documentation checks and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A proven mismatch is reported once and released without a termination request, every other refusal stays preserved, and each removed guard makes its named test fail.
- [ ] Force Stop, `doctor` and the existing reconcile reports are unchanged.
- [ ] The scope audit finds no undeclared Governed Path change and no change to the files row 8 names.

## Verification



The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0300-reconcile-releases-a-stale-run-owner-record/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-7; Invariants 1-6; Testing Approach; Build Order 3
- ADR-0313
