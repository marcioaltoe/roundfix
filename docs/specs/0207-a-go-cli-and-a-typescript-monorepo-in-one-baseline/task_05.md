---
task: task_05
spec: 0207-a-go-cli-and-a-typescript-monorepo-in-one-baseline
status: pending
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers and settles the Spec on evidence. This Spec adds an optional decision with clauses that apply under its values, a composed Setup Snapshot, a built-in composed profile and a check of the root gate's parts. Every behavior row is exercised by executing the named tests against the built tree, by running the built binary against disposable repositories, or by comparing assets through Git. No command reaches a provider or a live remote.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify, by executing task_01's tests against the built tree, that the frontend guide renders the suggestion sentence and the six current clauses with no recorded layout, the same clauses under `systems`, and only the repository-layout clause under `repository-defined`; that an unrecorded optional decision is not missing; that an invalid optional decision and a gate on a required decision are refused; and that a clause a recorded layout turns off is a reasoned rejection while a clause missing for no recorded reason is unaccounted.
3. MUST verify through the built binary, in a disposable copy of a Standard TypeScript Monorepo adoption, that `roundfix baseline plan --decision frontend.layout=repository-defined` with the other decisions answered renders the repository-layout clause, and that the same plan without that flag renders both systems clauses and reports no missing decision.
4. MUST verify, by executing task_02's tests, that the composed snapshot equals the union of its components, that a drifted, unknown-component and conflicting composition is each refused with its code, and that an asset sync rewrites the composed snapshot with its components and leaves it byte-identical with no source change.
5. MUST verify, by executing task_03's tests and Spec 0200's dispatch check, that `go-cli-typescript-monorepo` loads on `go-cli-typescript-bun`, renders every guide of its modules with no clause line repeated, carries the Go guide's scope sentence, and converges.
6. MUST verify, by executing task_04's tests and through the built binary in a disposable repository with a Go module, the Standard TypeScript Monorepo workspaces and a root Makefile, that a gate reaching both parts yields no part divergence and a gate reaching one yields exactly one non-blocking `verification.gate.part.missing`.
7. MUST verify through the built binary, in a disposable copy of this repository, that `roundfix baseline update --repo <copy> --no-skills --format text` reports `File changes: 0`, and that `make baseline-digests` run in that copy reports `"changed":false`.
8. MUST verify through Git that, between the Spec's base commit and the audited head, the Source Baseline corpus, the retention transitions and `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json` are byte-identical, that every clause of `internal/baseline/assets/modules/frontend.json` present at the base keeps its identifier, enforcement and guidance, and that no setup snapshot other than `go-cli-typescript-bun` changed.
9. MUST record, as evidence this Spec did not author: React's file-structure documentation (<https://legacy.reactjs.org/docs/faq-structure.html>) and its sentence that React has no opinion on folder layout; the renderer layout of `~/dev/argus`, read only; the frontend guides and Setup Manifests of the adopters the PRD's Acceptance evidence names under `~/dev/secondbrain/projects/*/mirror/docs/agents/`, read only; and the aggregate Make gates of the published repositories the PRD names. The row records what each source held when read and where it came from. When a source is unavailable or no longer carries the text, it MUST record that fact as the row's evidence or as a blocked row with its reason. The gate MUST NOT write to the Secondbrain or to another repository.
10. MUST verify that this Spec's own artifacts satisfy the promise rule, and check whether the Spec introduced, changed or retired a glossary term that `CONTEXT.md` does not carry; the candidates are the Frontend Layout Decision and the composed Setup Snapshot.
11. MUST verify from the repository history that each Task's changed files stay within its declarations or its `## Recorded paths`, that every Governed Path is bounded in `_authorization.md`, and that every regenerated file is an output of a sanctioned command.
12. MUST record the non-waivable Pull Request row as `blocked (environment: no open Pull Request)`, naming the Pull Request row in its provenance, and support it with the pre-PR equivalent of each control:
    - approval: the maintainer's recorded delivery authority;
    - checks and status: the Daemon's repository Verification at the audited head;
    - unresolved threads: none, because no Pull Request exists;
    - Merge-Ready acceptance: none yet, because the pre-PR review precedes archive and publication;
    - review-artifact ancestry: the audited head named as the claimed candidate.
13. MUST NOT accept a row whose only evidence is that a file was read, except the outside-evidence row, whose subject is what an outside source holds.
14. MUST NOT write to the live Run Database under `~/.roundfix`.

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

- `newest="$(find 'docs/specs/0207-a-go-cli-and-a-typescript-monorepo-in-one-baseline/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-5; Stories 1-6; Core Features 1-6; Success Metrics 1-7;
Declared breaks; Prerequisites; Acceptance evidence; `_techspec.md` → Testing
Approach 1-5; API Contracts 1-6; ADR-0058; ADR-0059; ADR-0060; ADR-0061;
ADR-0063; ADR-0067; ADR-0072; ADR-0073; ADR-0080; ADR-0081; ADR-0088;
ADR-0091; ADR-0093; ADR-0094; ADR-0096; ADR-0099; ADR-0103; ADR-0104;
ADR-0117; ADR-0130; ADR-0149; ADR-0155; ADR-0156; ADR-0166; ADR-0167;
ADR-0186; ADR-0190; ADR-0191; ADR-0204; ADR-0205.
