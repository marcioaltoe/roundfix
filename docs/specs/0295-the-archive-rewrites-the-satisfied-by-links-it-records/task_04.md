---
task: task_04
spec: 0295-the-archive-rewrites-the-satisfied-by-links-it-records
status: completed
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes `roundfix archive`
rebase each `unproven` entry of the Archive Record it writes exactly as it
rebases the outcome, and makes the Record Link Repair of `roundfix history
sanitize` and `roundfix baseline update`, and the `SC-HISTORY-RETAINED` fix
line, point a promoted target at its `docs/references/` copy.

Every row runs on temporary repositories, homes and Git configuration under
the Run's temporary directory. No row reads or writes the real `~/.roundfix`,
reaches the network, sends a request to a provider, the judge or a reviewer,
reads a credential, or reads another repository. The gate never runs
`roundfix spec judge`, never runs `roundfix review` against a real provider,
and never runs an archive plan without `--offline`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1 through the built product:
   - In a temporary directory resolved with `pwd -P`, with `HOME` and
     `GIT_CONFIG_GLOBAL` set to files under it, build the Characterization
     Repository `_techspec.md` → Current behavior describes, in its two
     shapes, one fresh copy per case: the delivery-branch shape (a bare
     `origin` with `origin/HEAD`, `references/promoted.md` added on the
     branch `docs/0001-demo`) archived with `bin/roundfix archive 0001-demo
     --promote references/promoted.md --pull-request 91`, and the main shape
     (every file committed on `main`) archived with `bin/roundfix archive
     0001-demo --promote references/promoted.md`.
   - For each shape, record the archive's exit code and output, the
     record's whole `unproven` list, and that the promoted entry is
     `` the operator records it in [`references/promoted.md`](../../references/promoted.md) ``
     and the removed entry is the label followed by the Git Locator of
     `docs/specs/0001-demo/references/removed.md` (an `@<commit>` locator in
     both shapes, because `main` holds the file).
   - After committing each archive, run `bin/roundfix spec check --strict
     0002-other` and record that its output holds no `SC-HISTORY-RETAINED`
     line.
3. MUST verify Success Metric 2 and Surface Transcript 1:
   - Build a second binary from a `git archive` copy of the Delivery Base,
     outside the checkout, and archive a fresh delivery-branch copy with it;
     record that its record holds both links verbatim.
   - On that copy, after committing the archive, run `roundfix history
     sanitize` with both binaries and record that the Delivery Base binary
     names `docs/history/specs/references/promoted.md` (pull request #91) for
     the promoted link and the audited binary names
     `docs/references/promoted.md`, while both name the removed file's Git
     Locator; run `spec check --strict 0002-other` with the audited binary
     and record the promoted link's fix line naming
     `../../references/promoted.md` and `docs/references/promoted.md`.
   - On a fresh main-shape copy archived with the Delivery Base binary and
     committed, run `bin/roundfix history sanitize` and record that standard
     output, standard error and the exit code match Surface Transcript 1.
   - Create the annotated tag `history-full` on that copy, run `bin/roundfix
     history sanitize --apply --batch 1`, and record that the record's
     entries become `` [`references/promoted.md`](../../references/promoted.md) ``
     and the removed file's Git Locator at the record's `source_revision`,
     the forms the audited archive writes in Requirement 2's main shape, and
     that a second plan reports nothing pending.
4. MUST verify Success Metric 3: on the main-shape copy of Requirement 2, add
   a third declaration whose `satisfied-by` holds no link, and record that its
   entry is byte-identical to the declaration's text; and record that `roundfix
   history sanitize` without `--apply` leaves every file of the Requirement 3
   copy byte-identical.
5. MUST verify the sabotage proofs. On temporary copies of the audited head
   made with `git archive`, one sabotage per copy, run the named tests and
   record that each fails:
   - remove the `unproven` rebase from `BuildArchiveRecord`:
     `TestArchiveRebasesTheUnprovenLinksItRecords` and
     `TestArchiveOfAQualifyingPartialPassesTheStrictHistoryCheck` fail;
   - make `RecordLinkPromotion` never answer a copy:
     `TestPlanRecordLinksPointsAPromotedTargetAtItsCopy`,
     `TestHistorySanitizePointsAPromotedRecordLinkAtItsCopyTranscript` and
     `TestHistoryRetainedNamesThePromotedCopyForARecordLink` fail;
   - drop the condition that the tree holds the copy:
     `TestRecordLinkPromotionRequiresTheListedCopy` fails.
6. MUST execute the tests task_01, task_02 and task_03 name in their
   Verification and record each outcome, including
   `TestBaselineUpdatePlansAPromotedRecordLinkAtItsCopy` and
   `TestLegacyConversionRebasesTheUnprovenLinksItRecords` (Invariants 1-5).
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - run `go doc net/url.URL.ResolveReference` and record that the passage
     the PRD quotes occurs in its whitespace-normalized output;
   - because this Spec's authorization forbids network access, record the
     GitHub Docs page as `blocked (environment: network prohibited by Spec
     authorization)`, quoting the passage the PRD cites. Such a row never
     decides a qualifying `partial`;
   - record the Delivery Base binary's record and strict-check output of
     Requirement 3 as the reproduction of the symptom the adopted Inbox Entry
     reports.
8. MUST record that the PRD declares no Unreachable Acceptance, so
   `rows_blocked_declared` is 0 and no row is recorded `blocked (declared
   unreachable)`.
9. MUST perform the glossary and documentation checks:
   - record the revised **Legacy Link Rebase** and **Lenient Legacy Reading**
     entries and the new **Record Link Repair** entry in `CONTEXT.md`, that
     the glossary contract tests pass, and whether any other term the Spec
     introduced needs an entry (`docs/agents/domain.md`);
   - record the sentences the `archive`, `history`, `baseline` and `spec`
     guides gained, the sentences both copies of the Roundfix `archive`,
     `baseline` and `spec` references gained, the raised and recorded skill
     version, and that `TestTheSkillCoverageMapIsCurrent` passes;
   - search the guides, the agent guides and the skill copies for wording
     that says the archive copies `satisfied-by` text verbatim or that the
     `record-links` repair writes only Git Locators, and record what it
     finds.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every Governed Path changed is bounded in `_authorization.md`;
    - these are byte-identical to the Delivery Base: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every Baseline module
      under `internal/baseline/assets/`, `docs/agents/setup-context.json`,
      every file under `docs/history/`, every file under `docs/adr/`,
      `docs/references/skill-coverage.json`, `internal/spec/archive.go`,
      `internal/spec/repository_path.go`, `internal/spec/cancel.go`,
      `internal/spec/legacy_link_rebase.go`, `internal/spec/git_locator.go`,
      `internal/cli/archive.go`, `internal/cli/cancel.go`,
      `internal/cli/history.go`, `internal/cli/baseline_update_history.go`,
      `internal/cli/deliver_workflow.go`, every file under `skills/` and
      `.agents/skills/` other than the Roundfix `SKILL.md` pair, the
      Roundfix `references/archive.md`, `references/baseline.md` and
      `references/spec.md` pairs and `skills/testdata/owned-skill-versions.json`,
      and every file under `docs/user-guide/` other than
      `commands/archive.md`, `commands/history.md`, `commands/baseline.md`
      and `commands/spec.md`;
    - every `_test.go` file under `internal/spec`, `internal/speccheck` and
      `internal/cli` that exists at the Delivery Base is byte-identical to it;
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
    `bin/roundfix archive 0295-the-archive-rewrites-the-satisfied-by-links-it-records --plan --offline`
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
- [ ] Archive both Characterization Repository shapes, compare with the Delivery Base binary, and reproduce Surface Transcript 1.
- [ ] Run each named test and each sabotage on temporary copies of the audited head.
- [ ] Record the outside evidence, the glossary and documentation checks and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Every archived record holds rebased `unproven` entries, the strict check reports no `SC-HISTORY-RETAINED`, and each removed guard makes its named test fail.
- [ ] The repair names `docs/references/promoted.md` where the Delivery Base names a missing path, and Surface Transcript 1 matches.
- [ ] No record already written changed without an applied repair, and the scope audit finds no undeclared Governed Path change.

## Verification


The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0295-the-archive-rewrites-the-satisfied-by-links-it-records/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Behavior changes; Success Metric 1; Success Metric 2; Success Metric 3; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-5; Surface Transcript 1; Invariants 1-5; Testing Approach; Build Order 4
- ADR-0308; ADR-0247; ADR-0277; ADR-0260; ADR-0261; ADR-0248; ADR-0299; ADR-0104; ADR-0240; ADR-0276
