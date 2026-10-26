---
task: task_05
spec: 0322-runs-are-not-lost-to-an-authorization-row-or-a-provider-at-capacity
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the shared QA partial
policy accept an outside-evidence row in the form `blocked (environment:
<subject> prohibited by Spec authorization)` with a bounded subject, counting
each row toward one exempt kind; lets the next Fallback Selection take any
Task or QA turn after a Provider System Error, reopening a QA verdict the
failed turn wrote; rewords the two Baseline clauses; and documents both rules
in the skills, guides, ADR-0335 and `CONTEXT.md`.

Every row runs on this repository's own files, temporary copies of the
audited head, temporary directories and package fakes. No row reads or writes
the real `~/.roundfix`, starts an Agent session, reaches the network, sends a
request to a provider, the judge or a reviewer, reads a credential, or reads
another repository. Every `go test` it runs carries an explicit `-timeout`,
and every row that runs a `roundfix` binary sets `HOME` and `TMPDIR` to fresh
temporary directories, because the CLI loads `$HOME/.roundfix/config.yml`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1-3 and Invariants 1-2 by running, each with
   `-v` and recording each named subtest's pass:
   `TestANetworkProhibitedRowNamesWhatTheAuthorizationForbids`,
   `TestTheArchivedAuthorizationRowsQualify` and the existing partial-policy
   tests of `./internal/spec`;
   `TestQAReportAcceptArchiveAndSettleAgreeOnEveryPartialShape` in
   `./internal/cli`; and `TestTheDaemonSettlesEveryPartialShapeAsTheCommandsDo`
   in `./internal/daemon`.
3. MUST verify Success Metrics 4-6 and Invariants 3-6 by running, each with
   `-v`, every test task_02's first Verification command names in
   `./internal/daemon`, `TestDeliveryRunResultReadsAnExhaustedProviderSystemError`
   in `./internal/cli` and `TestCoverageEquivalence` in `./internal/spec`,
   recording each named subtest's pass.
4. MUST exercise the built binary end to end on Spec 0309's shape. In a fresh
   temporary directory write `spec/_prd.md` holding one `## Unreachable
   Acceptance` entry (`criterion`, `reason`, `satisfied-by` lines) and
   `spec/qa/qa-report-2026-10-26.md` with frontmatter `verdict: partial`,
   `rows_blocked_environment: 4`, `rows_blocked_finding: 0`,
   `rows_blocked_declared: 1` and a `## Results` table (`#`, `Status`,
   `Provenance`) holding the six rows of the `Spec 0309` case of
   `internal/spec/qa_prohibited_subject_test.go`
   `TestTheArchivedAuthorizationRowsQualify` (R05 `pass`, R06 and R08 the
   network status, R07 `another repository`, R09 `blocked (declared
   unreachable)`, R12 the Pull Request row):
   - with the audited `bin/roundfix`, record that
     `bin/roundfix qa-report accept <dir>/spec/qa/qa-report-2026-10-26.md`
     exits `0` with no output;
   - with a binary built from the Delivery Base in a temporary `git archive`
     copy (`go build -buildvcs=false -o <tmp>/roundfix-base ./cmd/roundfix`),
     record that the same command exits `1` and prints
     `rows_blocked_environment is 4, 1 outside the pre-PR Pull Request row and
     the network-denied and network-prohibited outside-evidence rows; expected
     0 outside them`, the defect the adopted Backlog Entry reports;
   - change R07's subject to `another repository (quaestor)` and record that
     the audited binary exits `1` with the same message.
5. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, and record that the named case fails, naming the input
   that makes the mutation observable:
   - `qaReportResultsCounts` compares with `QANetworkProhibitedStatus` by
     `==` again: `TestANetworkProhibitedRowNamesWhatTheAuthorizationForbids/accepts_Spec_0309_R07_another_repository`
     and `TestTheArchivedAuthorizationRowsQualify/Spec_0309` fail;
   - the subject is counted with `len` instead of `utf8.RuneCountInString`:
     `.../accepts_a_subject_of_80_multibyte_characters` fails (80 × `é` is 160
     bytes);
   - the parenthesis rule is dropped:
     `.../refuses_a_subject_with_a_parenthesis` fails;
   - the Network-Denied Row check becomes a second `if`:
     `.../a_row_in_both_forms_counts_once` fails (the both-forms row hides the
     `credentials unavailable` row);
   - `recoverProviderSystemError` keeps `!owner.handedOff` for Tasks:
     `TestProviderSystemErrorAfterFirstHandoffActivatesTheFallback` fails;
   - it keeps `reportPending()` for QA:
     `TestProviderSystemErrorRecoveryEligibility/QA_with_a_finished_report_falls_back`
     and `TestQAProviderSystemErrorAfterAVerdictFallsBackAndReopensIt/partial`
     fail;
   - `resetUnfinishedQAVerdict` is not called:
     `TestQAProviderSystemErrorAfterAVerdictFallsBackAndReopensIt/pass` fails;
   - `Run` does not resend the failed repair prompt:
     `TestProviderSystemErrorAfterFirstHandoffActivatesTheFallback` fails.
6. MUST execute the tests task_03 and task_04 name in their Verification, run
   `go test -timeout 15m -count=1 ./skills` and
   `go test -timeout 20m -count=1 -tags docscontract ./internal/docscontract`,
   and record, in a temporary copy of the audited head with `HOME` at a fresh
   temporary directory, that
   `go run -buildvcs=false ./cmd/roundfix baseline update --repo . --no-skills --no-history --yes --format text`
   changes no file.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - read `docs/history/specs/0317-the-baseline-sends-handoffs-and-messages-through-ai-memory.md`
     and `docs/history/specs/0309-a-python-baseline-profile-for-uv-pytest-and-ruff.md`
     in the audited checkout and quote each `qa_override_reason` phrase
     "outside the accepted wording" with its row id;
   - because this Spec's authorization forbids network access, record the TAP
     14 specification row and the LangSmith fallback row each as
     `blocked (environment: network prohibited by Spec authorization)` with an
     `outside-evidence row` provenance item, quoting the passage the PRD
     cites. Such rows never decide a qualifying `partial`.
8. MUST record the PRD's single Unreachable Acceptance entry as exactly one
   row, `blocked (declared unreachable)`, counted once in
   `rows_blocked_declared`: no row may reach a real provider.
9. MUST perform the glossary and documentation checks:
   - record that `CONTEXT.md` revises **Network-Prohibited Row**,
     **Provider System Error**, **Fallback Chain** and **First Handoff** as
     task_04 requires (`docs/agents/domain.md`);
   - search `docs/user-guide`, `.agents/skills`, `skills`, `docs/agents`,
     `internal/baseline/assets/modules` and `CONTEXT.md` for
     "authorization forbids network access", "record the row exactly as" and
     "take the Task under the same condition", and record that nothing
     matches;
   - record that the `qa-gate`, `archive-spec` and `roundfix` skills' two
     version fields are equal in each canonical copy and mirror and recorded
     in `skills/testdata/owned-skill-versions.json`.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - these are byte-identical to the Delivery Base: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `internal/agent/`,
      every Baseline module other than `spec-workflow` and `autonomous-work`,
      every skill other than `qa-gate`, `archive-spec` and `roundfix`, and
      every ADR;
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
    `HOME="$(mktemp -d)" bin/roundfix archive 0322-runs-are-not-lost-to-an-authorization-row-or-a-provider-at-capacity --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's `## Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - start an Agent session;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test, the end-to-end binary rows and each sabotage on temporary copies.
- [ ] Record the outside sources, the declared row, the glossary and documentation checks and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Spec 0309's shape is accepted by the audited binary where the Delivery Base binary refuses it, and a subject with a parenthesis is still refused.
- [ ] A repair turn and a QA turn that wrote a verdict complete on the fallback after a Provider System Error, and an exhausted chain still settles as runtime infrastructure.
- [ ] Each sabotage makes its named case fail, and the scope audit finds no change outside the declared files.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0322-runs-are-not-lost-to-an-authorization-row-or-a-provider-at-capacity/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; User Experience; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; Interfaces; Data Models; API Contracts 1-8; Invariants 1-6; Testing Approach; Build Order 5
