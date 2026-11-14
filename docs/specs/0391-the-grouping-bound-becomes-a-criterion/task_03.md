---
task: task_03
spec: 0391-the-grouping-bound-becomes-a-criterion
status: pending
type: qa
complexity: medium
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec rewords the Baseline's
grouping clause so a Spec is sized by one context and by Tasks that each earn
their verification cycle, with both reasons stated and four implementation
Tasks plus the QA gate kept as a reference, not a cap, while every guarantee
of the old clause stays; it carries the same criterion into the `write-idea`,
`write-prd`, `write-techspec` and `write-tasks` skills, the Roundfix `spec`
reference and the `spec` command guide; and it records the decision in
ADR-0400, its notice in ADR-0208, and the glossary.

Every row runs on this repository's own files, temporary homes, temporary
repositories and temporary copies of the audited head. No row starts an Agent
session or an ACP adapter, calls ai-memory, reads or writes the real
`~/.roundfix` or the real home directory, runs the real `gh`, reaches the
network, reads another repository, sends a request to a provider, the judge
or a reviewer, or reads a credential. Every `go test` it runs carries an
explicit `-timeout`. Every row that runs `bin/roundfix` sets `HOME` and
`TMPDIR` to fresh temporary directories, with `GOMODCACHE` and `GOCACHE` kept
at their `go env` values, puts a fake `gh` script that prints `[]` first on
`PATH`, unsets `GH_TOKEN`, `GITHUB_TOKEN`, `TYPESAFE_API_KEY` and every Judge
Key Variable for that command, and removes the temporary directories
afterwards. The QA never writes into the checkout while a test suite runs. The
gate keeps its evidence under names the repository does not ignore, `.txt` for
command output, and links only files it kept.

The `go test` package and test runs of this gate do not run `bin/roundfix`
and keep the session's own `HOME` and `TMPDIR`: their tests make their own
temporary homes, and some, such as `TestGuidanceCompositionJourney/rust-cli`,
call an installed formatter (`rustfmt` through rustup) that needs the
session's toolchain configuration under `HOME`. Only rows that run
`bin/roundfix` use fresh temporary homes.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 3 and Invariants 1, 2 and 4 by running
   `go test -count=1 -timeout 20m -v` with
   `TestTheGroupingBoundIsACriterionThatKeepsEveryGuarantee`,
   `TestTheGroupingClausesCarryTheirForceAndText`,
   `TestTheGroupingClausesRenderInTheGuides`,
   `TestTheGroupingClausesHaveSourceBaselineRows` and
   `TestAnAdopterRetainsTheGroupingClauses` over `./internal/baseline`, and
   `TestTheAuthoringSkillsSizeASpecByContextAndVerificationCycle` and
   `TestEveryOwnedSkillVersionIsRecorded` over `./skills`, recording every row
   and subtest by name; then the whole `./internal/baseline` and `./skills`
   packages, `-run 'Baseline|Guidance|SetupManifest' ./internal/cli` and
   `-tags docscontract ./internal/docscontract`, recording each pass.
3. MUST reproduce the managed refresh in a temporary copy of the Delivery Base
   (the merge base of the audited head with the default branch) made with
   `git archive` and committed: record that
   `bin/roundfix baseline update --repo . --no-skills --no-history --yes`
   writes into `docs/agents/spec-routing.md`, exactly once each after
   whitespace normalization, "Size a Spec by two criteria, never by a count of
   its Tasks.", "Tasks that each earn their verification cycle", "joins a
   neighbouring Task of the same context instead of becoming its own Task",
   "a Spec that holds several contexts grows giant", "Four implementation
   Tasks plus the QA gate is a reference, not a cap" and "give each source
   exactly one owning Spec"; that the file no longer holds "only while the
   Spec fits four implementation Tasks plus its QA gate"; and that a second run
   reports the repository current.
4. MUST run each sabotage of `_techspec.md` → Testing Approach on its own
   temporary copy of the audited head made with `git archive`, with `-timeout`
   on every `go test`, and record that the named test fails in the named
   subtest: sabotages 1 to 6 as listed there. Sabotages 1 to 4 also run the
   module record command and `make baseline-digests` in that copy before the
   test.
5. MUST execute the tests tasks 01 and 02 name in their Verification.
6. MUST search and record:
   - both copies of `write-idea`, `write-prd` and `write-techspec` carry API
     Contract 2, both copies of `write-tasks` carry API Contract 3, both
     copies of the Roundfix `spec` reference and
     `docs/user-guide/commands/spec.md` carry API Contract 4;
   - this repository's `docs/agents/spec-routing.md`, the formatter golden of
     the Standard TypeScript profile and
     `internal/baseline/assets/modules/spec-workflow.json` carry API Contract
     1, the clause id, its `mandatory` force and its absence of `replaces` are
     unchanged, and `clause.spec.sources-01-group-by-shared-context` is
     byte-identical to the Delivery Base;
   - every file under
     `internal/baseline/assets/source-baselines/` is byte-identical to the
     Delivery Base;
   - `CONTEXT.md` carries API Contract 5 and ADR-0208 carries API Contract 6;
     ADR-0400 is `accepted`;
   - `.agents/skills` and `skills` hold the same bytes for each changed skill
     file;
   - the phrases "fits four implementation Tasks plus its QA gate" and "The
     grouping bound is four implementation Tasks" appear in none of
     `.agents/skills`, `skills`, `docs/agents`, `docs/user-guide`,
     `internal/baseline/assets/modules` and
     `internal/baseline/assets/formatter-fixtures`, searching only those paths
     (this Spec's folder, `docs/adr/`, `docs/history/`, `CHANGELOG.md`, the
     Source Baseline corpus and the test that names the absent phrases quote
     the old words and are excluded);
   - no Go file outside tests, no strict-check rule and no judge question in
     `internal/judge/questions.json` counts a Spec's Tasks (Invariant 3).
7. MUST record, as outside-evidence rows, the published sources this Spec did
   not write, each with its date and none blocked, through the passages
   `_prd.md` → Acceptance evidence quotes, looking nothing up:
   - Reinertsen, *The Principles of Product Development Flow*, Chapter 1 and
     principle B11 (read 2026-11-14): compare "batch sizes have their own
     economic tradeoffs, which lead to U-curve optimization problems." and
     "Economic batch size is a U-curve optimization." with requirement 3's
     guide, which sizes a Spec by two opposing costs instead of a count;
   - Reinertsen, "The science of batch size", FlowCon France 2019 (read
     2026-11-14): compare "the correct answer is never on one extreme or the
     other" and "when you're near the bottom of the U-curve, you don't need a
     very exact answer." with "Four implementation Tasks plus the QA gate is a
     reference, not a cap".
8. MUST record that the PRD declares no Unreachable Acceptance entry.
9. MUST perform the glossary check (`docs/agents/domain.md`): record that the
   PRD and TechSpec declare `changes:` for **Grouping Bound**, that
   `CONTEXT.md` defines it with the text of API Contract 5, and that
   `TestGlossaryPhrasesAreDefinedTerms` passes.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `CHANGELOG.md`, `internal/app/version.go`,
      `dist/npm/roundfix/package.json`,
      `docs/references/coverage-record.json`, any file under `internal/cli/`,
      any file under `internal/judge/`, `internal/speccheck/` or
      `internal/spec/`, any Vendored Skill under `.agents/skills`, any file
      under `internal/baseline/assets/modules/` other than
      `spec-workflow.json`, any file under
      `internal/baseline/assets/source-baselines/`, any file under
      `docs/user-guide/` other than `commands/spec.md`, `.claude/`, `.codex/`
      or `.opencode/`, any skill other than `write-idea`, `write-prd`,
      `write-techspec`, `write-tasks` and `roundfix`, and any Roundfix
      reference other than `spec`;
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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$(go env GOMODCACHE)" GOCACHE="$(go env GOCACHE)" bin/roundfix archive 0391-the-grouping-bound-becomes-a-criterion --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's Outcome as one paragraph that does not end with a colon. It MUST
    NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix` or the real home directory;
    - run the real `gh`, call ai-memory or reach the network;
    - send any message;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the
      repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests and reproduce the managed refresh.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] An adopter's managed refresh sizes a Spec by one context and by Tasks that earn their verification cycle, states both reasons, keeps four Tasks as a reference, keeps every guarantee of the old clause, and records the clause retained.
- [ ] The authoring skills, write-tasks, the Roundfix spec reference and the spec command guide state the same criterion in both copies, and the glossary and ADR-0208 agree with ADR-0400.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0391-the-grouping-bound-becomes-a-criterion/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-3; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contracts 1-6; Invariants 1-4; Declared breaks; Testing Approach; Build Order 3
