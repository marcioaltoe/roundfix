---
task: task_05
spec: 0214-measure-before-changing
status: completed
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec adds the read-only
`roundfix runs causes` command, a flag-gated harness that re-scores the Task
acceptance question through the judge, and three measured reference
documents with verdicts from rules fixed in advance; it removes no archived
file and changes no gate. Every behavior row is exercised through the built
binary or by executing the named tests against the built tree, in a
disposable repository with a disposable Roundfix Home, except the
outside-evidence rows, which read the maintainer's Run Database read-only and
published pages. No row sends a request to OpenRouter or TypeSafe.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify that `internal/runcause/signatures.json` is byte-identical to the block under `_techspec.md` → The signature table, and `internal/judge/testdata/task-acceptance-question.json` to the block under The Task acceptance question.
3. MUST verify, by executing task_01's tests against the built tree, the store query, every refusal of `Load`, every signature and trigger, first-match order, the symbolic-link and outside-directory log refusals, the corrective-Task rule, `specs_not_found`, and that the temporary home is byte-identical after `Build` and after the command.
4. MUST verify through the built binary, in a disposable repository holding Spec `0300-example` and a disposable `HOME` whose Run Database the fixture writes, that `roundfix runs causes` reproduces Surface Transcripts 1 and 2, and that Surface Transcripts 3 and 4 are reproduced as written; that `--format json` prints every field of API Contract 2; and that the help names the synopsis of API Contract 3. After each run it MUST record that every file under the disposable home kept its SHA-256.
5. MUST verify, by executing task_02's tests against the built tree, that no request body carries a Task's `## Result` text or source-file text, each key reaches only its own endpoint, only the generic keys give a `blocked` record and no request, an unpinned model is skipped, the Judge Log gains one line per request, the ceiling stops requests, and `TestMeasureTaskAcceptance` is skipped without its flag.
6. MUST run `TestCausesRecordIsConsistent` and `TestTaskAcceptanceRecordIsConsistent` with the committed records and documents as task_04's Verification does, and record each recomputed verdict.
7. MUST run, as the outside-evidence row on a record this Spec did not design, `bin/roundfix runs causes` over the committed cause record's window against the maintainer's Roundfix Home (read-only; `NODE_OPTIONS` unset), and compare its summary with the committed record. A difference MUST be explained by Runs that Journal Retention or the GC Command removed since, named by their identifiers; when the Run Database is unavailable, the row is blocked with that reason.
8. MUST verify, against Git, that the archive document's seven header lines equal `git ls-tree -r -l` at the commit it names, and that `git diff --name-only --diff-filter=DR` from the merge base of this Spec's branch with `main` to the audited head lists no path under `docs/history/`.
9. MUST verify that each Backlog Entry a verdict or a candidate set calls for exists, is `open`, and cites its reference document, and that none exists for a rule that did not fire.
10. MUST record, as evidence this Spec did not author, the sources of `_prd.md` → Acceptance evidence: the two failure-taxonomy papers, the two AUC interval papers and the Git book page. For each page it reaches, it MUST record what it read and whether it still says what `_prd.md` quotes; for each it cannot reach, it MUST record the row as blocked with the reason. It MUST record the 2026-09-30 Task-lint figures of the adopted Backlog Entry and the Run Database figures of `_techspec.md` → Measured outside evidence as evidence this Spec did not produce.
11. MUST verify that `docs/user-guide/commands/runs.md` carries `#### Why Verification failed`, that `.agents/skills/roundfix/references/runs.md` carries `### Why Verification failed`, that each mirror equals its canonical copy, that the `### QA settlement` section of every skill is unchanged, and that `make skills-sync-check` exits `0`.
12. MUST verify that no non-test Go file outside `internal/judge` gained an import of `net/http`, and that `internal/runcause` imports no package that sends network requests.
13. MUST verify that this Spec's own artifacts satisfy the promise rule.
14. MUST verify from the repository history that each Task's changed files stay within its declarations or its `## Recorded paths`, that every Governed Path is bounded in `_authorization.md`, and that `internal/cli/cli_test.go`, `.roundfixrc.yml`, `Makefile` and `go.mod` did not change.
15. MUST perform the glossary check of `docs/agents/domain.md` for the terms this Spec introduces or leans on — cause class, signature, corrective Task (which `CONTEXT.md` lists under Corrective Spec as a term to avoid, while `docs/agents/autonomous-work.md` uses it for a Task added after the gate) — and record its outcome.
16. MUST record the non-waivable Pull Request row as `blocked (environment: no open Pull Request)`, naming the Pull Request row in its provenance, and support it with the pre-PR equivalent of each control:
    - approval: the maintainer's recorded delivery authority;
    - checks and status: the Daemon's repository Verification at the audited head;
    - unresolved threads: none, because no Pull Request exists;
    - Merge-Ready acceptance: none yet, because the pre-PR review follows archive;
    - review-artifact ancestry: the audited head named as the claimed candidate.
17. MUST NOT accept a row whose only evidence is that a file was read.
18. MUST NOT send any request to OpenRouter or TypeSafe, read a real `ROUNDFIX_OPENROUTER_API_KEY`, `ROUNDFIX_TYPESAFE_API_KEY`, `OPENROUTER_API_KEY` or `TYPESAFE_API_KEY`, run the harness with its flag, write under the real `~/.roundfix`, or delete, move or rewrite any file under `docs/history/`.

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

- `newest="$(find 'docs/specs/0214-measure-before-changing/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-5; User Stories 1-5; Core Features 1-11; Success Metrics
1-5; Acceptance evidence; Recorded limits; `_techspec.md` → Surface
Transcripts 1-4; API Contracts 1-4; Testing Approach 1-6; Decision rules;
Measured outside evidence; ADR-0004; ADR-0008; ADR-0014; ADR-0033; ADR-0035;
ADR-0038; ADR-0080; ADR-0081; ADR-0088; ADR-0089; ADR-0091; ADR-0093;
ADR-0094; ADR-0096; ADR-0104; ADR-0117; ADR-0120; ADR-0121; ADR-0130;
ADR-0149; ADR-0155; ADR-0156; ADR-0166; ADR-0167; ADR-0178; ADR-0182;
ADR-0184; ADR-0187; ADR-0189; ADR-0193; ADR-0200; ADR-0201; ADR-0214;
ADR-0215.
