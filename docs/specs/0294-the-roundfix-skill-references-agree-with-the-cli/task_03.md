---
task: task_03
spec: 0294-the-roundfix-skill-references-agree-with-the-cli
status: completed
type: qa
complexity: high
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the Roundfix skill's
`review`, `events` and `deliver` references, the `review` and `deliver` user
guides and the glossary's Run Event Stream entry state what the CLI does, and
adds documentation contract tests that read the reviewer's conventions and the
stream's categories from the code.

Every row runs in the audited checkout, on temporary copies of it made with
`git archive`, or under a temporary `HOME` in the Run's temporary directory.
No row reads or writes the real `~/.roundfix`, reaches the network, sends a
request to a provider, the judge or a reviewer, reads a credential, or reads
another repository. The gate never runs `roundfix spec judge`, never runs
`roundfix review` against a real provider, and never runs an archive plan
without `--offline`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 and 2 through the guards:
   - Run `TestTheDeliveryConventionsAreDocumentedAsTheReviewerPromptsThem`,
     `TestTheDeliveryConventionContractReportsAStaleConvention`,
     `TestTheEventCategoriesAreDocumentedAsTheStreamAcceptsThem` and
     `TestTheEventCategoryContractReportsAMissingCategory` with `-tags
     docscontract`, and record each result (Invariants 1 and 2).
   - Run `TestReviewPromptCarriesTheDeliveryConventionsAndTheFindingGrammar`
     in `./internal/cli` and `TestDefaultFilterIncludesUsage` and
     `TestUsageFilterSelectsOnlyUsage` in `./internal/runevent`, and record
     that the prompt and the stream the documents describe are unchanged.
   - On temporary copies of the audited head made with `git archive`, apply
     one sabotage per copy and record that the named test fails:
     - put the old C3 sentence ("The archive commit moves a completed Spec's
       directory to the archive root and stamps its archive front matter.")
       back into the canonical `review` reference:
       `TestTheDeliveryConventionsAreDocumentedAsTheReviewerPromptsThem`
       fails naming that file and `C3`;
     - append one word to the C4 prompt in
       `internal/cli/review_conventions.go`, leaving every document:
       `TestTheDeliveryConventionsAreDocumentedAsTheReviewerPromptsThem`
       fails naming `C4`, proving the guard reads the source;
     - remove `usage` from the canonical `events` reference's category
       sentence, and in a second copy remove its `usage` table row:
       `TestTheEventCategoriesAreDocumentedAsTheStreamAcceptsThem` fails in
       each, naming `usage`;
     - add a seventh category to `AllStreamCategories` in
       `internal/runevent/stream.go`, leaving every document:
       `TestTheEventCategoriesAreDocumentedAsTheStreamAcceptsThem` fails
       naming it, proving the guard reads the package.
3. MUST verify Success Metric 3:
   - run `TestGlossaryClosedListsEqualTheirConstants` and
     `TestGlossaryContractReportsAPlantedViolationOfEachRule` and record the
     result;
   - on a temporary copy, remove the `usage` code span from the glossary's
     Run Event Stream entry and record that
     `TestGlossaryClosedListsEqualTheirConstants` fails naming Run Event
     Stream.
4. MUST exercise the documented surfaces through the built binary, under a
   temporary `HOME` holding a temporary Git repository:
   - record `bin/roundfix events --help` and that its `--filter` line names
     the six categories the `events` reference lists;
   - run `bin/roundfix events run_missing --filter usage` and record that it
     exits 2 for the missing Run Database or Run, not for the filter, and run
     `bin/roundfix events run_missing --filter bogus` and record that it
     exits 2 naming the unknown event category.
5. MUST verify Success Metric 4: read `conventionRegions`,
   `reviewSpecQAOverride`, the Delivery Queue's prerequisite check (the
   function in `internal/cli/deliver_workflow.go` that reads `requires`) and
   `provenDeliveryEvidence` in `internal/worktree/merged_head.go`, and record
   for each corrected statement of task_01 the code it now matches. Then
   repeat task_01's sweep over the Roundfix skill, the `review`,
   `review-report-shape`, `events` and `deliver` user guides and the glossary,
   and record every remaining match with its reason.
6. MUST execute the tests task_01 and task_02 name in their Verification.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - quote from the adopted Inbox Entry under `references/` the two
     contradictions and the lines it cites, and record that the audited head
     no longer holds either statement;
   - because this Spec's authorization forbids network access, record the
     Nassif et al. paper and the Metzner post as `blocked (environment:
     network prohibited by Spec authorization)`, quoting the passages the PRD
     cites. Such a row never decides a qualifying `partial`.
8. MUST record Success Metric 5 as exactly one row, `blocked (declared
   unreachable)`, against the PRD's single Unreachable Acceptance entry, and
   count it once in `rows_blocked_declared`: no row may read another
   repository, the real `~/.roundfix` or a provider.
9. MUST perform the glossary and documentation checks:
   - record that `CONTEXT.md` revises **Run Event Stream** with the six
     categories, that the glossary contract tests pass, and whether any other
     term the Spec introduced needs an entry (`docs/agents/domain.md`);
   - record that `TestTheSkillCoverageMapIsCurrent` passes and that the skill
     mirror is in sync (`make skills-sync-check`).
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every Governed Path changed is bounded in `_authorization.md`;
    - these are byte-identical to the Delivery Base: every non-test Go file,
      every `_test.go` file other than
      `internal/docscontract/command_documentation_test.go` and
      `internal/docscontract/user_guide_contract_test.go`, the Makefile, the
      CI workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every Baseline
      module under `internal/baseline/assets/`,
      `docs/agents/setup-context.json`, every file under `docs/history/` and
      `docs/adr/`, `docs/user-guide/commands/events.md` and
      `docs/references/skill-coverage.json`;
    - the Roundfix skill version was raised and recorded;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance. It supports the row with the pre-PR equivalent of each
    control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
12. MUST, for the Outcome for the Archive Record, run only
    `bin/roundfix archive 0294-the-roundfix-skill-references-agree-with-the-cli --plan --offline`
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
- [ ] Exercise `events` through the binary and read the code each corrected statement matches.
- [ ] Record the outside evidence, the declared row, the glossary and documentation checks and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Each guard passes on the audited head and fails on its sabotage, including a change to the code with the documents left alone.
- [ ] Every corrected statement matches the code it describes, and the sweep finds no remaining contradiction.
- [ ] The scope audit finds no production Go change and no undeclared Governed Path change.

## Verification


The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0294-the-roundfix-skill-references-agree-with-the-cli/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-7; Invariants 1-4; Testing Approach; Build Order 3
