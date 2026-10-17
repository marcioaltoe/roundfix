---
task: task_05
spec: 0276-the-history-cleanup-reaches-every-adopter-and-the-review-sees-a-large-spec
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes the history cleanup write a
successor's missing lineage line, refuse only an ADR without an accepted
successor and repair live links; makes the pre-PR review omit deleted history
and generated files, review a large candidate in Review Segments and exempt the
QA Task from C2; adds the Spec Consistency Check gap
`SC-SKILL-VERSION-UNTASKED`; adds `roundfix archive --plan --offline`; and
makes the Delivery Queue run a fresh copy of each item binary. Every row runs
on temporary repositories and homes. No row reads or writes the real
`~/.roundfix`, reaches the network, sends a request to a provider or the judge,
reads a credential, or reads another repository; the gate never runs
`roundfix spec judge` and never runs an archive plan without `--offline`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1 through the built `bin/roundfix` with `HOME`
   set to a temporary home, `NODE_OPTIONS` unset and every Judge Key Variable
   unset, on a fresh temporary Git repository that holds a copy of the audited
   tree's `docs/agents/setup-context.json` and `docs/agents/docs-layout.md`, so
   the retirement rule is on, one minimal active Spec, and: a superseded ADR
   whose accepted successor does not name it, a superseded ADR with no
   accepted successor, a rejected ADR, a `docs/references/` page and a
   `docs/user-guide/` page each linking a file an earlier commit held and the
   head no longer does, and the annotated tag `history-full`. Record that
   `roundfix history sanitize` lists the lineage line as a citation with the
   exact bytes of API Contract 1, the refusal by path of API Contract 2, the
   `adr` unit removing the other two ADRs and the `live-links` line of API
   Contract 3; that after `roundfix history sanitize --apply --batch <n>` the
   successor ends with that line once and both pages hold Git Locators; that
   the `SC-HISTORY-RETAINED` findings of `roundfix spec check --strict` name
   only the kept ADR; and that a second plan lists only that refusal. Record
   the `baseline update` JSON shape from the output of
   `TestBaselineUpdatePlansLineageAndLiveLinks` (Invariants 1 to 3).
3. MUST verify Success Metrics 2 and 3 by running the tests task_02 names and
   `TestReviewBlocksAboveTheBoundWithoutAProviderCall`,
   `TestReviewCommandKeepsTheCodexDiffByteBound` and
   `TestReviewCommandSendsAClaudePromptBeyondTheCodexByteBound`, and record
   from their output the omitted reasons, the segment count, the merged record
   fields and each blocking reason of API Contracts 4 and 5 and the C2 text of
   API Contract 6 (Invariants 4 to 6). The review is exercised only through the
   fake runner of those tests; no reviewer is started.
4. MUST verify Success Metric 4 through the built binary on a temporary Git
   repository holding `skills/testdata/owned-skill-versions.json` with one
   owned skill and an active Spec whose pending Task declares that skill's
   reference without its mirror, `SKILL.md` pair, record and step: record that
   `roundfix spec check <slug> --strict` reports `SC-SKILL-VERSION-UNTASKED`
   with the summary of API Contract 7 and exits 1; then add the declarations
   and the step, commit, and record that the code is absent.
5. MUST verify Success Metric 5 through Surface Transcript 1: on a temporary
   Git repository holding a committed active Spec whose only candidate is one
   `references/` file, run `bin/roundfix archive <slug> --plan --offline` and
   compare its standard output, standard error and exit code with the
   transcript's lines under its controlled variations; record that
   `TestArchivePlanOfflineSkipsTheJudge` passes, which proves no judge request
   with a key present, and that `--offline` without `--plan` is refused
   (Invariant 7).
6. MUST verify Success Metric 6 by running the delivery tests task_04 names
   and recording the copy path, the unchanged bytes and identity of the held
   copy and the absent temporary file (Invariant 8). The macOS kill itself is
   environment-dependent and is not a row.
7. MUST verify Invariant 9: record that every count in
   `internal/docscontract/testdata/corpus-golden.json` at the Delivery Base,
   the merge base of the audited head with `main`, is unchanged at the audited
   head, that `SC-SKILL-VERSION-UNTASKED` appears with the count the corpus
   test reports, and that `TestCheckCorpusGolden` and
   `TestCheckActiveCorpusHasNoErrors` pass.
8. MUST execute the tests task_01 to task_04 name in their Verification, and
   prove each new guard fails without its implementation, each on a temporary
   copy of the audited tree made with `git archive`: remove the
   `detectSkillVersion` call; make `PlanHistoryKind` return the successor
   error again; make `reviewScopedDiff` read derived declarations from the
   head; make `stepExecutable` return the build output path; record that the
   matching named tests fail each time.
9. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the PR-Agent review documentation, the
   GitHub Docs page on `linguist-generated`, the ADR Tools description and
   Apple's "Updating Mac Software", quoting each passage the PRD cites and
   stating whether it supports the design. Because this Spec's authorization
   forbids network access, each lookup is recorded as
   `blocked (environment: network prohibited by Spec authorization)`, which
   never decides a qualifying `partial`.
10. MUST perform the glossary check of `docs/agents/domain.md`: record that
    `CONTEXT.md` defines **Review Segment**, that the glossary contract tests
    pass, and whether any other term the Spec introduced needs an entry.
11. MUST verify from the repository history that each Task's changed files
    stay within its declarations or its `## Recorded paths`; that every
    Governed Path changed is bounded in `_authorization.md`; that the
    Makefile, the CI workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every
    Baseline module and `internal/speccheck/history_retained.go` did not
    change; that the Roundfix, write-tasks and qa-gate skill versions were
    raised and recorded; that no `### QA settlement` section changed; that no
    guide or skill copy still says a superseded ADR whose successor does not
    name it refuses the unit, still names `roundfix/delivery-conventions/v2`,
    or tells a QA Agent to run the archive plan without `--offline`. This row
    reads Task commits, so it declares a `commit_range` input.
12. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance, and support it with the pre-PR equivalent of each
    control: approval from the recorded delivery authority, checks from the
    Daemon's repository Verification at the audited head, no unresolved
    threads because no Pull Request exists, no Merge-Ready acceptance yet, and
    the audited head named as the claimed candidate.
13. MUST, for the Outcome for the Archive Record, run only
    `bin/roundfix archive 0276-the-history-cleanup-reaches-every-adopter-and-the-review-sees-a-large-spec --plan --offline`
    with every Judge Key Variable unset. It MUST NOT accept a row whose only
    evidence is that a file was read, read or write `~/.roundfix`, record any
    environment value, send any request to a provider or the judge, or start a
    workflow.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Replay the history, spec check and offline archive shapes through the built binary.
- [ ] Run the review and delivery tests and sabotage each guard on temporary copies.
- [ ] Record the outside evidence, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The adopter shapes apply through the built binary and leave only the kept ADR reported.
- [ ] The offline archive plan matches Surface Transcript 1 and contacts no judge.
- [ ] Each removed guard makes its named test fail, and no earlier corpus count moved.
- [ ] The scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0276-the-history-cleanup-reaches-every-adopter-and-the-review-sees-a-large-spec/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-7; Core Features 1-9; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; API Contract 7; API Contract 8; API Contract 9; API Contract 10; Surface Transcript 1; Invariants 1-9; Current behavior; Build Order 5
- ADR-0282; ADR-0278; ADR-0254; ADR-0281; ADR-0273; ADR-0104; ADR-0276
