---
task: task_02
spec: 0308-recorded-paths-are-read-before-carry-forward-provenance
status: completed
type: qa
complexity: high
---

# Task 02: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the reader of a Task
file's `## Recorded paths` section find it when only Daemon-written settlement
sections follow it, such as `## Carry-forward provenance`, and makes the
writer replace it in place, so a carried Task's recorded paths count in the
corrective retry and in reconcile while the corrective-scope rule is
unchanged.

Every row runs on temporary repositories and Git configuration under the
Run's temporary directory, with synthetic Task files. No row reads or writes
the real `~/.roundfix`, reads another repository, reaches the network, sends a
request to a provider, the judge or a reviewer, or reads a credential. The
gate never runs `roundfix spec judge`, never runs `roundfix review` against a
real provider, and never runs an archive plan without `--offline`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 and 2: run
   `TestARecordedPathBeforeCarryForwardProvenanceCountsAsDeclared` in
   `./internal/cli` and record, per subtest, that the correction of
   `review.txt` was accepted and that the correction also changing
   `other.txt` was refused with exactly `changed path other.txt is outside
   the named archived Specs and the paths their Tasks declared` (Invariant
   4). Run `TestARecordedPathCountsAsDeclared`,
   `TestACodeCorrectionOutsideTheDeclaredPathsNamesTheFirstPath` and
   `TestDeclaredPathsAreReadAtTheRecordedRevision` and record that they pass
   unchanged.
3. MUST verify Success Metrics 3, 4 and 5:
   - Run `TestRecordedPathsAreReadBeforeTrailingSettlementSections` and
     `TestRecordTaskPathsReplacesTheSectionBeforeCarryForwardProvenance` in
     `./internal/spec` and record each row's result (Invariants 1-3).
   - Run `TestRecordedPathsKeepAFencedHeadingInAuthoredText`,
     `TestRecordedPathsKeepAnAuthoredHeadingFollowedBySections`,
     `TestRecordedPathsReplaceTheTrailingDaemonSection`,
     `TestRecordTaskPathsReplacesAnExistingSectionAndKeepsEveryOtherByte`
     and `TestRecordTaskPathsWithNoPathsRemovesTheSection`, and record that a
     file ending with its recorded section reads and rewrites as before
     (Invariant 2).
   - Run `TestArchivedSpecLeftoverRecordedBeforeCarryForwardProvenance` and
     `TestArchivedSpecLeftoverScopeKinds` in `./internal/worktree` and record
     the reconciliation states.
   - Run `TestTaskCommitRecordsAnUndeclaredTestFile`,
     `TestSettleRecordsTheUndeclaredPathsOfItsTaskCommit` and
     `TestVerificationPathUndeclaredIgnoresWhatThePlanProvides`, and record
     that the writers and the Verification path provision are unchanged.
   - On temporary copies of the audited head made with `git archive`, apply
     one sabotage per copy to `internal/spec/recorded_paths.go` and record
     that the named test fails, with the row whose input makes the mutation
     observable:
     - stop the walk at `## Result` and `## Carry-forward provenance` as at
       any other heading:
       `TestRecordedPathsAreReadBeforeTrailingSettlementSections` fails in
       its `carry-forward-provenance` row, and
       `TestARecordedPathBeforeCarryForwardProvenanceCountsAsDeclared` fails;
     - never stop the walk at an authored heading: the same test fails in its
       `authored-section-between` row, and
       `TestRecordedPathsKeepAnAuthoredHeadingFollowedBySections` fails;
     - treat lines inside a fence as headings: the same test fails in its
       `fenced-heading-before-provenance` row, and
       `TestRecordedPathsKeepAFencedHeadingInAuthoredText` fails;
     - make the writer drop every byte after the replaced section:
       `TestRecordTaskPathsReplacesTheSectionBeforeCarryForwardProvenance`
       fails in its `new-paths` row;
     - make the writer always append a new section:
       the same test fails in its `same-paths` row, and
       `TestRecordedPathsReplaceTheTrailingDaemonSection` fails.
4. MUST execute the tests task_01 names in its Verification.
5. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - read the adopted Inbox Entry under `references/` in the audited
     checkout, written by oraculum before this Spec existed, and record that
     its refusal text, with the path replaced by `other.txt` and whitespace
     runs collapsed to one space, equals the refusal
     `TestARecordedPathBeforeCarryForwardProvenanceCountsAsDeclared` asserts
     up to the semicolon;
   - because this Spec's authorization forbids network access, record the
     CommonMark Spec, 4.5 Fenced code blocks, as `blocked (environment:
     network prohibited by Spec authorization)`, quoting the passage the PRD
     cites. Such a row never decides a qualifying `partial`.
6. MUST record Success Metric 6 as exactly one row, `blocked (declared
   unreachable)`, against the PRD's single Unreachable Acceptance entry, and
   count it once in `rows_blocked_declared`: no row may read oraculum's
   Delivery Queue or Roundfix Home.
7. MUST perform the glossary and documentation checks:
   - record that the Spec declares no glossary term, and whether any term it
     introduced needs an entry (`docs/agents/domain.md`);
   - search the guides, the agent guides and the skill copies for wording
     that says the recorded section is read only when it ends the Task file,
     and record what it finds.
8. MUST verify from the repository history:
   - the Task's changed files stay within its declarations or its `##
     Recorded paths`;
   - no Governed Path changed, matching `paths: []` in `_authorization.md`;
   - these are byte-identical to the Delivery Base: the Makefile, the CI
     workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every file under
     `internal/baseline/assets/`, `docs/agents/setup-context.json`, every
     file under `docs/history/`, `docs/user-guide/`, `.agents/skills/` and
     `skills/`, `CONTEXT.md`, `internal/spec/spec.go`,
     `internal/spec/declared_paths.go`,
     `internal/cli/deliver_review_correction.go`, `internal/cli/settle.go`,
     `internal/cli/review_conventions.go`, `internal/worktree/merged_head.go`,
     `internal/speccheck/verification_path.go` and
     `internal/daemon/task_engine.go`;
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
    `bin/roundfix archive 0308-recorded-paths-are-read-before-carry-forward-provenance --plan --offline`
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
- [ ] Compare the Inbox Entry's refusal text with the asserted refusal, and record the blocked outside source.
- [ ] Record the declared row, the glossary and documentation checks and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A carried Task's recorded path is accepted by the corrective retry and reconcile, a path in no Task is still refused, and each sabotage makes its named row fail.
- [ ] A file ending with its recorded section reads and rewrites exactly as before.
- [ ] The scope audit finds no Governed Path change and no change to the files row 8 names.

## Verification



The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0308-recorded-paths-are-read-before-carry-forward-provenance/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-3; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-4; Invariants 1-4; Testing Approach; Build Order 2
- ADR-0166
