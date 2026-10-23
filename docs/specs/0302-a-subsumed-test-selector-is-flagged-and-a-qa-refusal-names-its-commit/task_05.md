---
task: task_05
spec: 0302-a-subsumed-test-selector-is-flagged-and-a-qa-refusal-names-its-commit
status: completed
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the Spec Consistency
Check report `SC-VERIFY-SUBSUMED-SELECTOR` for a pytest node id that another
collection argument of the same invocation contains, states the rule in the
`write-tasks` skill and the `spec` documents, and makes every QA Report the
Daemon closes without an Agent turn record `build`, the commit it audited,
with the `qa-gate` refusal template, the Context-driven development guide and
the glossary saying so.

Every row runs on temporary repositories and Git configuration under the
Run's temporary directory. No row runs pytest, reads or writes the real
`~/.roundfix`, reaches the network, sends a request to a provider, the judge
or a reviewer, reads a credential, or reads another repository. The gate
never runs `roundfix spec judge`, never runs `roundfix review` against a real
provider, and never runs an archive plan without `--offline`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 and 2:
   - Run `TestSubsumedSelectorRaisesANodeIdBesideItsContainingPath`,
     `TestSubsumedSelectorAcceptsSelectorsNoArgumentContains` and
     `TestSubsumedSelectorJudgesOnlyPendingNonQATasks` in
     `./internal/speccheck`, and record, per subtest, the command and the
     finding count, and for one raising case the exact summary, fix and line
     (Invariants 1-3).
   - Through the public CLI, on a temporary copy of the audited head made
     with `git archive` and committed in a fresh temporary repository: add to
     this Spec's copy a pending `backend` Task `task_06.md` with an
     `## Overview`, a `## Consumers` section holding `None.` and a reason,
     and one Verification, `uv run pytest tests/a.py::test_novo tests/a.py`,
     the shape of the adopted Quaestor report; in the copy's `_tasks.md` add
     the node `task_06` needing `task_04`, add `task_06` to `task_05`'s
     `needs`, and add its row to the table; commit, and run `bin/roundfix
     spec check <slug> --strict --format json` from that repository with the
     binary of row 1; record the one `SC-VERIFY-SUBSUMED-SELECTOR` finding
     with its summary and fix.
     Then change that Verification to `uv run pytest tests/a.py::test_novo`,
     commit, rerun, and record that the code is gone. Record any other
     finding code the added Task raises without treating it as a failure of
     this row.
   - Run `TestVerificationPathUndeclaredRaisesTheMissingFile` and
     `TestVerificationPathUndeclaredIgnoresWhatThePlanProvides`, and record
     that the neighbouring detector is unchanged (Invariant 4 for the
     detectors; the prober's files are covered by row 8).
3. MUST verify Success Metrics 3 and 4:
   - Run `TestPreconditionRefusalReportRecordsItsAuditedCommit` in
     `./internal/spec` and record the opening frontmatter lines, the single
     `build` key and the blank-commit case.
   - Run `TestMechanicalQAReportRecordsTheAuditedCommitWhenItIsFinal` and
     `TestWriteMechanicalQAReportWritesThePreconditionRefusal` in
     `./internal/daemon` and record the `build` line of the refusal and of
     the blocking report, the fixture `HEAD` it names, and that the seed
     carries none (Invariants 5-6).
   - Run `TestWritePreconditionRefusalReportWritesOneTerminalRow`,
     `TestWritePreconditionRefusalReportIsReadableAsARefusal`,
     `TestPreconditionRefusalReportNamesItsAuditor`,
     `TestMechanicalQAReportSeedStartsPending`,
     `TestWriteMechanicalQAReportRecordsTheRefusal` and
     `TestMechanicalQAReportStartsUnsuffixedOnAFreshDate`, and record that
     the report shapes they pin are unchanged.
4. MUST verify Success Metric 5: on temporary copies of the audited head made
   with `git archive`, apply one sabotage per copy and record that the named
   test fails in the named case, whose input makes the mutation observable:
   - make the selector check judge an invocation that passes
     `--keep-duplicates`: `TestSubsumedSelectorAcceptsSelectorsNoArgumentContains`
     fails in its `--keep-duplicates` case;
   - make it read the value of a value-taking option as a collection
     argument: the same test fails in its `--deselect` and `-k` cases;
   - make a path contain any argument it is a string prefix of, without the
     `/` boundary: the same test fails in its `tests/test_ab.py::test_new
     tests/test_a` case;
   - make it compare a node id with an argument of identical text: the same
     test fails in its case that passes the same node id twice;
   - make it stop splitting a command at `;`, `&&`, `||`, `|`, `&`,
     parentheses and braces: the same test fails in its two-invocation case,
     and `TestSubsumedSelectorRaisesANodeIdBesideItsContainingPath` fails in
     its `cd tools &&` case;
   - make it skip the bodies of `$(...)` substitutions:
     `TestSubsumedSelectorRaisesANodeIdBesideItsContainingPath` fails in its
     `out="$(…)"` case;
   - make it ignore the `run` wrappers: the same test fails in its `uv run`
     cases;
   - make it judge Tasks of every status:
     `TestSubsumedSelectorJudgesOnlyPendingNonQATasks` fails in its completed
     and in-progress cases;
   - drop the `build` line from `WritePreconditionRefusalReport`:
     `TestPreconditionRefusalReportRecordsItsAuditedCommit` fails in its
     resolved-commit case;
   - write the commit unquoted: the same test fails in its resolved-commit
     case, whose commit reads as a YAML number;
   - write the key for a blank commit: the same test fails in its
     unresolved-commit case;
   - write `build` in the non-blocking seed:
     `TestMechanicalQAReportRecordsTheAuditedCommitWhenItIsFinal` fails in its
     seed case;
   - stop setting `AuditedCommit` in `resolveQAAuditorEvidence`: the same
     test fails in its written-refusal case, and
     `TestWriteMechanicalQAReportWritesThePreconditionRefusal` fails.
5. MUST execute the tests task_01 to task_04 name in their Verification.
6. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - read the adopted Quaestor Backlog Entry under `references/` in the
     audited checkout and record that the command shape it reports, a node id
     with its own file, is the shape row 2 ran through the CLI and that the
     check reported it;
   - read the adopted oraculum Backlog Entry and record that its refusal
     report `qa-report-2026-10-21-01.md` had no `build:` while its later
     reports did, beside the `build` line row 3 observed in the Daemon's
     refusal;
   - because this Spec's authorization forbids network access, record pytest
     9.0.0's changelog (#12083) as `blocked (environment: network prohibited
     by Spec authorization)`, quoting the passages the PRD cites. Such a row
     never decides a qualifying `partial`.
7. MUST record Success Metric 6 as exactly one row, `blocked (declared
   unreachable)`, against the PRD's single Unreachable Acceptance entry, and
   count it once in `rows_blocked_declared`: no row may read another
   repository.
8. MUST perform the glossary, documentation and scope checks:
   - record that `CONTEXT.md` defines **Subsumed Test Selector** linking
     ADR-0315 and revises **Precondition Refusal** and **QA Report**, that
     the glossary contract tests pass, and whether any other term the Spec
     introduced needs an entry (`docs/agents/domain.md`);
   - record the row `docs/user-guide/commands/spec.md` gained, the bullet
     both copies of the Roundfix `spec` reference gained, the section both
     copies of the `write-tasks` skill gained, the template line and
     paragraph both copies of the `qa-gate` skill gained, the guide's refusal
     paragraph, and that `TestTheSkillCoverageMapIsCurrent` passes;
   - verify from the repository history that each Task's changed files stay
     within its declarations or its `## Recorded paths`, that every Governed
     Path changed is bounded in `_authorization.md`, that the Roundfix,
     `write-tasks` and `qa-gate` skill versions were raised and recorded, and
     that no `### QA settlement` section changed;
   - verify that these are byte-identical to the Delivery Base: the
     Makefile, the CI workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every
     Baseline module under `internal/baseline/assets/`,
     `docs/agents/setup-context.json`, every file under `docs/history/`,
     `internal/daemon/verification_probe.go`, `internal/cli/spec_check.go`,
     `internal/speccheck/verification.go`,
     `internal/speccheck/verification_path.go`,
     `internal/speccheck/citations.go`, `internal/speccheck/mechanical.go`,
     `internal/spec/qa_test.go`, `internal/daemon/qa_hollow_report_test.go`,
     `internal/daemon/qa_report_sequence_test.go`,
     `internal/docscontract/testdata/corpus-golden.json` and
     `docs/references/skill-coverage.json`.

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
    `bin/roundfix archive 0302-a-subsumed-test-selector-is-flagged-and-a-qa-refusal-names-its-commit --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's `## Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - run pytest or any other External Test Tool;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test, the CLI flow on a temporary copy, and each sabotage on temporary copies of the audited head.
- [ ] Record the outside sources, the declared row, the glossary and documentation checks and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] An overlapping selector is reported through the CLI and the tests, every non-overlapping command and non-pending Task is not, and each removed guard makes its named test fail.
- [ ] The Daemon's final reports name the audited commit, the seed does not, and the pinned report shapes are otherwise unchanged.
- [ ] The scope audit finds no undeclared Governed Path change and no change to the files row 8 names.

## Verification



The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0302-a-subsumed-test-selector-is-flagged-and-a-qa-refusal-names-its-commit/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-8; Invariants 1-6; Testing Approach; Build Order 5
- ADR-0315
- ADR-0148
- ADR-0132
