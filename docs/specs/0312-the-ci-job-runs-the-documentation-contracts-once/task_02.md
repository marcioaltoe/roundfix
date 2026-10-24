---
task: task_02
spec: 0312-the-ci-job-runs-the-documentation-contracts-once
status: completed
type: qa
complexity: high
---

# Task 02: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes a pull request's
"Verify docs" step skip the `internal/docscontract` run that "Verify changed"
already made, keeps `make verify-docs` alone and pushes to main unchanged,
and adds a repository test that keeps every documentation contract file
`//verify:always`.

Every row runs `make` over stub tools in temporary directories, temporary
copies of the audited head and of the Delivery Base, and this repository's
own documents. No row reads or writes `~/.roundfix`, reaches the network,
sends a request to a provider, the judge or a reviewer, reads a credential,
or reads another repository. Every `go test` it runs carries an explicit
`-timeout`, and every row that runs `bin/roundfix` sets `HOME` to a fresh
temporary directory, because the CLI loads `$HOME/.roundfix/config.yml`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 and 2 and Invariants 1 and 2 for task_01:
   - run `TestVerifyDocsRunsTheDocumentationContractsOnceUnlessSkipped`,
     `TestTheCIPullRequestSkipsOnlyTheRepeatedDocumentationContracts`,
     `TestEveryDocumentationContractFileIsAlwaysSelected`,
     `TestMainAndReleaseRunEveryContract` and `TestVerifyContractsRunsEveryContract`
     with `-v` and record each pass;
   - in a temporary copy of the audited head made with `git archive`, run
     `make -n verify-docs` and `make -n verify-docs VERIFY_DOCS_TEST=docs-test-skipped`
     and record the documentation contract lines each prints (one, then
     only the `spec-budget` line);
   - record the "Verify docs" step of `ci-verify.yml` at the audited head
     and at the Delivery Base.
3. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, and record that the named test fails, with the case
   whose input makes the mutation observable:
   - make `verify-docs` depend on `docs-test` directly:
     `TestVerifyDocsRunsTheDocumentationContractsOnceUnlessSkipped/skipped`
     fails, counting one documentation contract run;
   - remove the `env` of the "Verify docs" step:
     `TestTheCIPullRequestSkipsOnlyTheRepeatedDocumentationContracts` fails;
   - replace `//verify:always` by `//verify:relevant` in
     `internal/docscontract/corpus_test.go`:
     `TestEveryDocumentationContractFileIsAlwaysSelected` fails naming that
     file.
4. MUST execute the tests task_01 names in its Verification.
5. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - read the adopted Backlog Entry under `references/` and record its 764 s
     and 541 s with five active Specs and its target;
   - because this Spec's authorization forbids network access, record the
     GitHub Actions expression documentation as `blocked (environment:
     network prohibited by Spec authorization)`, quoting the sentence the PRD
     cites. Such a row never decides a qualifying `partial`.
6. MUST record Success Metric 3 as exactly one row, `blocked (declared
   unreachable)`, against the PRD's single Unreachable Acceptance entry, and
   count it once in `rows_blocked_declared`: no row may read a CI run.
7. MUST perform the glossary check: record that the Spec declares no term,
   that `CONTEXT.md` is byte-identical to the Delivery Base, and whether any
   term the Spec introduced needs an entry (`docs/agents/domain.md`).
8. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`,
      and the Makefile and `.github/workflows/ci-verify.yml` changed only in
      the lines task_01 names;
    - these are byte-identical to the Delivery Base: `go.mod`, `go.sum`,
      `.roundfixrc.yml`, every workflow other than `ci-verify.yml`, every
      non-test `.go` file, every file under `internal/docscontract/`,
      `.agents/`, `skills/`, `docs/agents/` other than
      `docs/agents/setup-context.json` (which changed only through the
      sanctioned regeneration task_01 Requirement 7 names), and
      `docs/user-guide/`, and `CONTEXT.md`;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
9. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
10. MUST, for the Outcome for the Archive Record, run only
    `HOME="$(mktemp -d)" bin/roundfix archive 0312-the-ci-job-runs-the-documentation-contracts-once --plan --offline`
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
- [ ] Run each named test, the dry runs and each sabotage on temporary copies.
- [ ] Record the outside sources, the declared row, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A pull request's job runs the documentation contracts once, and `make verify-docs` alone still runs them.
- [ ] A documentation contract file that stops being `//verify:always` fails a test, and each sabotage makes its named test fail.
- [ ] The scope audit finds no change outside the declared lines and files.

## Verification



The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0312-the-ci-job-runs-the-documentation-contracts-once/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Feature 1; Success Metric 1; Success Metric 2; Success Metric 3; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contract 1; Invariants 1-2; Testing Approach; Build Order 2
