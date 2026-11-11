---
task: task_02
spec: 0377-a-push-whose-tests-pass-runs-every-contract-under-a-measured-budget
status: pending
type: qa
complexity: low
---

# Task 02: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the
matrix covers and settles the Spec on evidence. The Spec makes the suite-time
budget of a push to `main` one measured number, 420 s, in the `Makefile`,
records the verdict in `make test-budget` and judges it in
`make suite-budget-check`, and wires `.github/workflows/ci-verify.yml` so a
push whose tests pass runs `Verify docs` and `Verify every contract` before
the budget is judged in its last step.

Every row runs in this repository or in temporary directories it creates
and removes. No row creates, removes or changes any file in `/tmp`,
`/private/tmp` or `$HOME` outside a directory made by `mktemp -d`, reaches
the network, runs `gh` or `git push`, starts or reruns a workflow, reads
another repository, starts an Agent session or an ACP adapter, reads or
writes the real `~/.roundfix`, sends a request to a provider, the judge or a
reviewer, records an environment value, or sends a signal to any process.
Every `go test` carries an explicit `-timeout` and keeps the session's own
`HOME` and `TMPDIR`. Every row that runs `bin/roundfix` sets `HOME` and
`TMPDIR` to fresh temporary directories, with `GOMODCACHE` and `GOCACHE` kept
at their `go env` values, and removes them afterwards. The QA never writes
into the checkout while a test suite runs. The gate keeps its evidence under
names the repository does not ignore, `.txt` for command output, and links
only files it kept.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST run, with `go test -count=1 -timeout 20m -v ./internal/verifyselect`,
   the tests `TestTheSuiteTimeBudgetIsRecordedThenJudged` and
   `TestTheCIPushJudgesTheSuiteTimeBudgetLast`, recording the pass line of
   each of the eight subtests; then the whole `./internal/verifyselect`
   package, `go test -count=1 -timeout 20m -tags docscontract ./internal/docscontract -run '^(TestTheSetupManifestIsCurrent|TestGlossaryPhrasesAreDefinedTerms)$'`
   and `GOOS=windows go vet ./internal/verifyselect`.
3. MUST record Success Metrics 1 and 4 against the starting behavior: in a
   `git archive` of the Delivery Base (the merge base of the audited head
   with `main`) unpacked in a temporary directory and made a Git repository
   there, copy the audited head's
   `internal/verifyselect/suite_budget_test.go` over and run both tests with
   `-timeout`; record that `TestTheCIPushJudgesTheSuiteTimeBudgetLast` and
   seven subtests fail there, that only
   `without a verdict file an overrun fails at once` passes there and on the
   audited head (the kept path), and that all pass on the audited head
   (row 2).
4. MUST record Success Metric 2 through the workflow file itself: run
   `actionlint .github/workflows/ci-verify.yml` when `command -v actionlint`
   finds it and record its empty output, or record the row as
   `blocked (environment: actionlint not installed)` with the YAML-parse test
   of row 2 as the equivalent evidence.
5. MUST record Success Metric 3: read the 33 suite times in
   `references/2026-11-10-ci-verify-push-suite-times.md` and record that the
   largest is 343 s, below the `SUITE_BUDGET_SECONDS` default of the audited
   `Makefile`.
6. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive` and made a Git repository, with `-timeout` on every
   `go test`, and record that the named test or subtest fails:
   - `test-budget` exits 1 instead of 0 after recording the verdict:
     `TestTheSuiteTimeBudgetIsRecordedThenJudged/an_overrun_is_recorded_without_failing`;
   - `suite-budget-check` exits 0 on a line starting
     `suite-time budget exceeded:`:
     `TestTheSuiteTimeBudgetIsRecordedThenJudged/the_check_fails_an_exceeded_verdict`;
   - `test-budget` no longer removes the verdict file before the suite runs:
     `TestTheSuiteTimeBudgetIsRecordedThenJudged/a_failing_suite_records_no_verdict`;
   - the `SUITE_BUDGET_SECONDS` default set back to 360:
     `TestTheSuiteTimeBudgetIsRecordedThenJudged/the_default_budget_is_420_seconds`;
   - `suite-budget-check` exits 0 when no verdict was recorded:
     `TestTheSuiteTimeBudgetIsRecordedThenJudged/the_check_fails_without_a_verdict`;
   - `suite-budget-check` exits 0 on any other content:
     `TestTheSuiteTimeBudgetIsRecordedThenJudged/the_check_fails_an_unreadable_verdict`;
   - `SUITE_BUDGET_SECONDS: "240"` put back into the `Verify` step's `env`:
     `TestTheCIPushJudgesTheSuiteTimeBudgetLast`;
   - the budget step's `if` reduced to `github.event_name == 'push'`:
     `TestTheCIPushJudgesTheSuiteTimeBudgetLast`.
7. MUST execute the tests the non-QA Task names in its Verification.
8. MUST record, as the outside-evidence row, sources this Spec did not
   write, with their dates and not blocked: the adopted Backlog Entry
   (`references/`, recorded by the operator from CI runs this Spec did not
   design: "Run 42747838420 on `5be18746` took 215 seconds and run
   42749378382 on `9fda0b65` took 311"), and the published passages quoted
   in `_prd.md` → Acceptance evidence (Akinshin: "the “average” performance
   levels still can be up to three times different across builds"; hotpath:
   "100 runs landed on 6 different CPU models"; GitHub Docs on
   `steps.<step_id>.outcome` and `runner.temp`), compared with the budget of
   row 5 and the condition and verdict path of row 2. The row reads the
   files in this repository, never the network.
9. MUST record the PRD's one Unreachable Acceptance entry, the first push
   run of `CI | Verify` on `main` after the merge, as exactly one row
   `blocked (environment: push to main and a GitHub-hosted workflow run required)`
   with `declared: ` followed by the entry's criterion text, counted once in
   `rows_blocked_declared`.
10. MUST perform the glossary check (`docs/agents/domain.md`): record that
    the PRD and TechSpec declare `None.` under `## Glossary` and that the Task
    commit range does not change `CONTEXT.md`.
11. MUST verify from the repository history:
    - the Task's changed files stay within its declarations or its
      `## Recorded paths`;
    - every Governed Path changed is listed in `_authorization.md` `paths`,
      and `docs/agents/setup-context.json` changes only its two `Makefile`
      `declarationDigest` values;
    - ADR-0259 gains only the `updated_at` change and the amendment
      paragraph;
    - the Task commit range changes none of these: `.github/workflows/release.yml`
      or any other workflow, `go.mod`, `go.sum`, `.roundfixrc.yml`,
      `skills-lock.json`, `CHANGELOG.md`, `internal/app/version.go`,
      `dist/npm/**`, `CONTEXT.md`, any guide, skill, Go production file, any
      Baseline module or other derived Baseline file, or the
      `### QA settlement` section of any skill.

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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0377-a-push-whose-tests-pass-runs-every-contract-under-a-measured-budget --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read another repository, or run `git fetch`, `git push` or `gh`;
    - read or write `~/.roundfix`;
    - create a file in `/tmp`, `/private/tmp` or `$HOME` outside a
      `mktemp -d` directory it removes;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - send a signal to any process;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification and the named suites.
- [ ] Characterize the Delivery Base with the new tests and run each sabotage on a temporary copy.
- [ ] Record the outside sources, the declared unreachable row, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] An overrun with a verdict file is recorded and judged by the check on the audited head, and fails at once on the Delivery Base.
- [ ] The workflow records the verdict in `Verify`, keeps the docs and contract steps, and judges the verdict last; each sabotage makes its named test fail.
- [ ] The one Unreachable Acceptance entry is one declared row, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0377-a-push-whose-tests-pass-runs-every-contract-under-a-measured-budget/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-3; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; Testing Approach; Build Order 2
