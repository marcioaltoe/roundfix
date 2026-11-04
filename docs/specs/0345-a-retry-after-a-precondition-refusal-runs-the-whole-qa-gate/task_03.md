---
task: task_03
spec: 0345-a-retry-after-a-precondition-refusal-runs-the-whole-qa-gate
status: completed
type: qa
complexity: medium
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the QA Report Import
refuse a candidate commit that adds or modifies a report recording a
Precondition Refusal, so a retry after a precondition refusal runs the whole
QA gate; pins that the refusal report goes through the configured Format
Command before its commit; and documents the import in the Roundfix skill's
implement reference, the configuration guide, `CONTEXT.md` and ADR-0355.

Every row runs on this repository's own files, temporary repositories and
temporary copies of the audited head. No row starts an Agent session or an
ACP adapter, reads or writes the real `~/.roundfix`, reads another
repository, reaches the network, sends a request to a provider, the judge or
a reviewer, or reads a credential. Every `go test` it runs carries an
explicit `-timeout`. Every row that runs `bin/roundfix` sets `HOME` and
`TMPDIR` to fresh temporary directories, because the CLI loads
`$HOME/.roundfix/config.yml`, runs with every Judge Key Variable and
`TYPESAFE_API_KEY` unset, and removes the directories afterwards. The QA
never writes into the checkout while a test suite runs.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 4 and Invariants 1 to 4 by running
   `go test -count=1 -timeout 20m -v` over `./internal/daemon` with
   `TestPriorQAPassRefusesAPreconditionRefusal`,
   `TestARetryAfterAPreconditionRefusalRunsTheWholeGate`,
   `TestQAFormatFormatsAPreconditionRefusalReport`, every test whose name
   starts with `TestPriorQAPass`, `TestTwoGatePasses` or `TestQAFormat`,
   `TestAPriorPassIsImportedOnlyFromARecordedRunCommit`,
   `TestQAGateCarriesARowFromAnUnintegratedFailedPass`,
   `TestQAPromptNamesTheImportedPassHead` and
   `TestQAGateRefusesOnFailedRepositoryVerification`, recording each pass and
   each case of the three new tests.
3. MUST verify Invariant 1 from the repository history: the diff of
   `internal/daemon/qa_prior_pass_test.go`,
   `internal/daemon/qa_prior_pass_skip_test.go`,
   `internal/daemon/qa_two_pass_carry_test.go` and
   `internal/daemon/qa_format_test.go` between the Delivery Base (the merge
   base of the audited head with `main`) and the audited head removes no
   line, and the diff of `internal/daemon/qa_prior_pass.go` only adds lines
   inside `copyPriorQAPass`.
4. MUST verify Invariant 5 by reading `carryForwardTaskCommits`
   (`internal/cli/carryforward.go`) at the audited head and recording that it
   keeps only commits whose `Roundfix-Task` trailer names a Task, and that the
   Task commit range does not change that file.
5. MUST run each sabotage `_techspec.md` → Testing Approach lists on its own
   temporary copy of the audited head made with `git archive`, with
   `-timeout` on every `go test`, and record that the named test fails, with
   the case whose input makes the mutation observable.
6. MUST execute the tests task_01 and task_02 name in their Verification.
7. MUST search `docs/user-guide/`, `docs/agents/`, `CONTEXT.md`, `docs/adr/`
   and every copy of the Roundfix skill, and record that: both copies of the
   implement reference carry the four phrases task_02 requires and the list
   of what the import never brings in; the configuration guide carries "the
   Precondition Refusal report included"; ADR-0355 exists with status
   `accepted`; and no live document outside `docs/history/` says that a
   Precondition Refusal is imported or becomes a later pass's previous report.
8. MUST record, as outside-evidence rows, two sources this Spec did not
   write, each with its date and none blocked:
   - the Pantheon report quoted in
     `references/2026-11-03-pantheon-retry-reimports-a-precondition-refusal.md`
     (2026-11-03): compare offline its Run 2 shape (the refusal recommitted,
     no Agent session, a precondition failing on the imported report) with
     the starting-tree characterization task_01 records and with
     `TestARetryAfterAPreconditionRefusalRunsTheWholeGate`;
   - the Bazel user manual passage `_prd.md` → Acceptance evidence quotes
     ("the test failed"): compare it with a refusal never being reused as the
     next pass's previous report.
9. MUST record that the PRD declares no Unreachable Acceptance and count no
   row in `rows_blocked_declared`.
10. MUST perform the glossary check (`docs/agents/domain.md`): record that the
    PRD declares `adds: **QA Report Import**` and `changes: **Precondition
    Refusal**`, that `CONTEXT.md` defines the new term and carries the
    revision, and that the Task commit range changes no other glossary entry.
11. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`, `dist/`,
      `docs/references/skill-coverage.json`, any existing Archive Record, any
      Baseline asset, `docs/agents/`, `internal/daemon/task_engine.go`,
      `internal/daemon/qa_format.go`, `internal/spec`, `internal/speccheck`,
      `internal/cli`, any Roundfix reference other than `implement`, any
      skill other than the Roundfix files the Tasks declare;
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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" bin/roundfix archive 0345-a-retry-after-a-precondition-refusal-runs-the-whole-qa-gate --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's ` Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests and check the history diffs of the import tests and carry-forward.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A recorded refusal is refused and leaves no file, and a retry after it reaches the Agent and commits none of its files, in the tests.
- [ ] A failed pass that built its matrix is still imported byte for byte, the refusal report is formatted before its commit, each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0345-a-retry-after-a-precondition-refusal-runs-the-whole-qa-gate/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-6; Invariants 1-5; Declared breaks; Testing Approach; Build Order 3
