---
task: task_05
spec: 0209-sources-that-share-a-context-share-a-spec
status: pending
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec adds three mandatory Baseline
clauses, sources that share a context may share a Spec, grouping stops at four
implementation Tasks plus the gate, and an open record is extended instead of
duplicated, with their Source Baseline rows; teaches them in the authoring
skills; and adds an advisory grouping question to `roundfix spec judge`.
Every behavior row is exercised through the built binary or by executing the
named tests against the built tree, in a disposable repository with a
disposable Roundfix Home. No command reaches OpenRouter, TypeSafe or any other
network service, except the reads of published pages the outside-evidence row
names.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify, by executing task_01's tests against the built tree, Exact clause texts, Source Baseline rows and API Contracts 6 and 7: the three clauses with `mandatory` force in their modules, guides and Standard TypeScript Monorepo goldens, each Source Baseline row, and a ready adopter plan that records each new clause `retained` and none `unaccounted`. It MUST verify that every existing clause of the two modules is byte-identical to the starting commit, and that a second Managed Refresh of this repository changes no file.
3. MUST verify that `internal/judge/questions.json` holds the block under `_techspec.md` → The grouping question byte for byte, and that its question text, threshold of 0.3, statuses, scrub pattern and cut of 1,500 equal the measurement recorded in `_techspec.md` → Measured outside evidence.
4. MUST verify, by executing task_02's tests against the built tree, every rule of The grouping reader, Preparing a source, Grouping pairs and Asking and outcomes, including that an Inbox Entry, a `declined` Backlog Entry, a `done` Finding, a symbolic link, a non-English entry and a Go file reach no request on either transport, `suggested` at 0.30 and `clear` at 0.29, and Spec 0205's `TestRequestsCarryOnlySpecArtifactText` still passing.
5. MUST verify, by executing task_03's tests against the built tree, Surface Transcripts 1, 2 and 3 and API Contracts 1 and 2, and that Spec 0205's transcripts still pass with `0 suggested`.
6. MUST verify through the built binary, in a disposable repository holding the fixture of `_techspec.md` → Surface Transcripts, with `HOME` set to a disposable directory, `HTTPS_PROXY` and `HTTP_PROXY` pointed at a closed local port, and `ROUNDFIX_OPENROUTER_API_KEY` and `ROUNDFIX_TYPESAFE_API_KEY` removed, that `roundfix spec judge 0300-example --stage prd` with only a dummy `OPENROUTER_API_KEY` set reproduces Surface Transcript 3 and that the disposable Judge Log gained no line, and that `roundfix spec judge --help` names the grouping question as API Contract 5 states.
7. MUST verify that `write-idea`, `write-prd` and `write-techspec` carry `## Sources that share a context` immediately before `## Process` with the phrases task_04 requires, that `.agents/skills/roundfix/references/spec.md` and `docs/user-guide/commands/spec.md` describe the grouping question, that each mirror equals its canonical copy, that the `### QA settlement` section of every skill is unchanged, and that `make skills-sync-check` exits `0`.
8. MUST record, as evidence this Spec did not author, the grouping measurement of 2026-10-01 in `_techspec.md` → Measured outside evidence, and the published sources `_prd.md` → Acceptance evidence names: Sun et al., ASE 2011, and Zhang et al., TOSEM 2023, and the fluxus record read through the Secondbrain. For each page it reaches it MUST record what it read and whether it still supports the suggestion-only design; for each it cannot reach it MUST record the row as blocked with the reason.
9. MUST verify that this Spec's own artifacts satisfy the promise rule.
10. MUST verify from the repository history that each Task's changed files stay within its declarations or its `## Recorded paths`, that every Governed Path is bounded in `_authorization.md`, and that `Makefile`, `.roundfixrc.yml`, `go.mod`, `internal/cli/cli_test.go` and the CI workflows did not change.
11. MUST perform the glossary check of `docs/agents/domain.md` for the terms this Spec introduces, grouping suggestion and grouping bound, and record its outcome.
12. MUST record the non-waivable Pull Request row as `blocked (environment: no open Pull Request)`, naming the Pull Request row in its provenance, and support it with the pre-PR equivalent of each control:
    - approval: the maintainer's recorded delivery authority;
    - checks and status: the Daemon's repository Verification at the audited head;
    - unresolved threads: none, because no Pull Request exists;
    - Merge-Ready acceptance: none yet, because the pre-PR review follows archive;
    - review-artifact ancestry: the audited head named as the claimed candidate.
13. MUST NOT accept a row whose only evidence is that a file was read.
14. MUST NOT send any request to OpenRouter or TypeSafe, read the real `ROUNDFIX_OPENROUTER_API_KEY`, `ROUNDFIX_TYPESAFE_API_KEY` or `OPENROUTER_API_KEY`, or write under the real `~/.roundfix`.

## Subtasks

- [ ] Build the matrix from the Requirements above and the sources no declaration waives.
- [ ] Execute each row against the built tree and record its evidence.
- [ ] Write the dated QA Report with its verdict.

## Acceptance Criteria

- [ ] The QA Report records a verdict and names, in each row's provenance, the sources that row covers.
- [ ] The repository Verification result is recorded as a fact, not re-derived.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0209-sources-that-share-a-context-share-a-spec/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-5; User Stories 1-6; Core Features 1-8; Success Metrics
1-6; Acceptance evidence; Declared breaks; `_techspec.md` → Surface
Transcripts 1-3; Testing Approach 1-5; API Contracts 1-7; Measured outside
evidence; ADR-0058; ADR-0060; ADR-0080; ADR-0081; ADR-0083; ADR-0088;
ADR-0089; ADR-0091; ADR-0092; ADR-0093; ADR-0094; ADR-0096; ADR-0104;
ADR-0117; ADR-0149; ADR-0155; ADR-0156; ADR-0166; ADR-0167; ADR-0178;
ADR-0182; ADR-0184; ADR-0186; ADR-0187; ADR-0189; ADR-0193; ADR-0200;
ADR-0201; ADR-0203; ADR-0208; ADR-0209.
