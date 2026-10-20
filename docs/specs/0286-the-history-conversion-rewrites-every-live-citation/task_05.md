---
task: task_05
spec: 0286-the-history-conversion-rewrites-every-live-citation
status: completed
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec delivers four corrections to
the history conversion that `roundfix baseline update` and `roundfix history
sanitize` share: every listed citation in every Citation Source is rewritten,
a lineage line names and dates a retired ADR the way its repository does, a
pending Task's `creates:` target survives, and a broken link is resolved by
its file name. It also extends `SC-HISTORY-RETAINED` to report a removed
history citation that is not a Git Locator.

Every row runs on temporary repositories and homes with the built binary or
the named tests. No row reads or writes the real `~/.roundfix`, reaches the
network, sends a request to a provider, the judge or a reviewer, reads a
credential, or reads another repository. The gate never runs `roundfix spec
judge`, never runs `roundfix review` against a real provider, and never runs
an archive plan without `--offline`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1:
   - Run the tests task_01 names in its Verification and record each result
     (Invariants 1-2).
   - On a temporary copy of the audited head made with `git archive`, restore
     the old `scope` condition in `PlanHistoryRewrite` and record that
     `TestHistoryRewriteReachesEveryCitationSource` fails; on another, stop
     resolving `./` and `../` prose paths and record that
     `TestHistoryRewriteResolvesARelativeProsePath` fails; on a third, let
     `CitationSources` keep `testdata` files and record that
     `TestCitationSourcesSkipSkillTreesTestdataAndBinaryFiles` fails.
   - Through the built product, reproduce Surface Transcript 1 of
     `_techspec.md`: build its fixture in a temporary repository with `HOME`
     set to a temporary directory, run `bin/roundfix history sanitize`, and
     record that standard output, standard error and the exit code match the
     transcript. Then run `bin/roundfix history sanitize --apply --batch 10`
     and record the output of `git grep -n -E
     'history/(findings|backlog|adr)/'` over the repository without
     `docs/history/`, skill trees and `testdata`: every remaining line carries
     `@` followed by the commit. Record that the skill file and the
     `testdata` file are byte-identical to their fixture bytes.
3. MUST verify Success Metric 2:
   - Run the tests task_02 names in its Verification.
   - On a temporary copy of the audited head, make `historyADRIdentifier`
     return `ADR-` with the leading digits for every name and record that
     `TestHistoryADRLineageNamesADateNamedADRByItsFileName` fails; on another,
     date the line with `HEAD`'s committer date and record that
     `TestHistoryADRLineageDatesTheSupersession` fails.
   - Record the lineage line of Surface Transcript 1 from the run above.
4. MUST verify Success Metric 3:
   - Run the tests task_03 names in its Verification.
   - On a temporary copy of the audited head, drop the `creates:` exemption
     from `PlanLiveLinks` and record that
     `TestLiveLinksLeaveAPendingTaskCreatesTarget` fails; on another, drop the
     file-name resolution and record that
     `TestLiveLinksResolveABrokenLinkByTheRemovedFileName` fails.
   - From the Surface Transcript 1 run, record that `_prd.md` of the fixture
     Spec still links `references/medicao.md` after apply and that
     `references/note.md` carries the Git Locator of
     `docs/history/findings/2026-01-01-x.md`.
5. MUST verify Success Metric 4:
   - Run the test task_04 names in its Verification.
   - Through the built product: in a temporary repository holding Surface
     Transcript 1's fixture, delete the three history files with `git rm`
     and commit, without any conversion. Add the retirement rule to the
     guide as `docs/agents/setup-context.json` and
     `docs/agents/docs-layout.md` declare it, and commit again. Then run
     `bin/roundfix spec check --strict`. Record its exit code `1` and every
     `SC-HISTORY-RETAINED` line. There is one error per citation for the
     Makefile, the ADR, the agent guide, the reference page and the README,
     each naming the file, the line and `Replace the citation with` its
     locator. None names the skill file, the `testdata` file or the fixture
     Spec's `references/medicao.md` link. Other codes the minimal fixture
     Spec raises are recorded and ignored.
   - On a temporary copy of the audited head, disable the removed-citation
     finding and record that
     `TestHistoryRetainedReportsARemovedHistoryCitationWithoutALocator`
     fails.
6. MUST execute the tests every implementation Task names in its
   Verification, and record the wall time of `bin/roundfix spec check
   --strict` at the audited head as a gate fact for the Risks row of the
   TechSpec.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - from the adopted Inbox Entry under `references/`, the adopter's numbers:
     45 citations listed, 30 rewritten, 18 left, a lineage line naming `ADR-` and the year
     `2026`, dated 2026-10-16 against an `updated_at` of 2026-10-01, and `spec check
     --strict` exiting 0;
   - because this Spec's authorization forbids network access, record each
     web source as `blocked (environment: network prohibited by Spec
     authorization)`, quoting the passage the PRD cites. Such a row never
     decides a qualifying `partial`.
8. MUST perform the glossary check of `docs/agents/domain.md`:
   - record that `CONTEXT.md` carries **Citation Source** and **Removed
     History Path** and the revised **History Citation** and **Git Locator**,
     each linking ADR-0299;
   - record that the glossary contract tests pass;
   - record whether any other term the Spec introduced needs an entry.
9. MUST verify from the repository history:
   - each implementation Task's changed files stay within its declarations
     or its `## Recorded paths`;
   - every Governed Path changed is bounded in `_authorization.md`;
   - these are byte-identical to the Delivery Base: the Makefile, the CI
     workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every Baseline module
     under `internal/baseline/assets/`, `docs/agents/setup-context.json` and
     `internal/docscontract/testdata/corpus-golden.json`;
   - the Roundfix skill version was raised and recorded;
   - no `### QA settlement` section changed;
   - no guide, agent guide or skill copy still says that other citations are
     listed and left unchanged, that the lineage date is HEAD's committer
     date, or that the update never rewrites citing files.

   This row reads Task commits, so it declares a `commit_range` input.
10. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
11. MUST, for the Outcome for the Archive Record, run only
    `bin/roundfix archive 0286-the-history-conversion-rewrites-every-live-citation --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset. It MUST NOT
    do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - start a workflow or push a tag.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test and each sabotage on temporary copies of the audited head.
- [ ] Reproduce Surface Transcript 1, apply it, and run the check on the shape an earlier conversion left, through the built binary.
- [ ] Record the outside evidence, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Every listed citation is rewritten in every Citation Source, and each removed guard makes its named test fail.
- [ ] The lineage line, the pending `creates:` target and the broken link behave as Surface Transcript 1 shows.
- [ ] `spec check --strict` fails on the shape an earlier conversion left with one error per citation, and the scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0286-the-history-conversion-rewrites-every-live-citation/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`


## References

- `_prd.md` → Goals; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; API Contract 7; Surface Transcript 1; Invariants 1-6; Current behavior; Build Order 5
- ADR-0299; ADR-0278; ADR-0282; ADR-0104; ADR-0240; ADR-0276
