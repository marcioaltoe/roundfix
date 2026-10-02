---
task: task_05
spec: 0206-stack-rules-with-force-and-the-rules-adopters-repeat
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers and settles the Spec on evidence. This Spec changes eight Baseline modules, three profiles (the composed profile Spec 0207 shipped only loses the removed Bun rule), one retention transition, the Rust guide template and the Standard TypeScript Monorepo Source Baseline, and adds four test files. Every behavior row is exercised by executing the named tests against the built tree, by running the built binary's Managed Refresh in disposable repositories, or by comparing assets through Git. No command reaches a provider or a live remote.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify, by executing task_01's tests against the built tree, that the instruction and skill guides of the Standard TypeScript Monorepo and Go CLI/TUI profiles state the eight promoted core clauses with their force, that the TypeScript and Bun guide no longer states the conditional warnings clause, and that a Source Baseline adopter's plan records that clause `replaced` by `clause.core.lint-warnings-block` while the same plan without `replaces` is refused as `unaccounted`. It MUST also verify, by executing `TestTheComposedProfileRendersEveryGuideWithNoRepeatedClause` and `TestTheComposedProfilePlanConverges`, that the composed `go-cli-typescript-monorepo` profile renders its guides with no repeated clause and converges with the promoted clauses of every Task in place.
3. MUST verify, by executing task_02's tests, that no Go, CLI or TUI rule carries rule-level guidance, that the rendered guides label every rule line with its force and state the recorded-reason exception, the cross-platform build and the skill-ships-with-behavior clause, and that the Go guide names no library it forbids.
4. MUST verify, by executing task_03's tests, that the Rust CLI plan renders the scope sentence and the typed-error, entry-point and no-panic clauses with their force, and that the Rust CLI profile requires the new rule.
5. MUST verify, by executing task_04's tests, that the Spec-routing, docs-layout and TypeScript guides state the four clauses, and that every one of the twelve promoted Standard TypeScript Monorepo clauses has its Source Baseline row with its force.
6. MUST verify through the built binary, in a disposable Standard TypeScript Monorepo adoption created from the Spec's base commit and then refreshed with the audited binary, that `roundfix baseline update --yes --no-skills --format text` applies, lists the Bun warnings clause under retention as `replaced`, and renders `A lint warning fails Verification.` in the instruction guide.
7. MUST verify through the built binary, in a disposable copy of this repository, that `roundfix baseline update --repo <copy> --no-skills --format text` reports `File changes: 0`, and that `make baseline-digests` run in that copy reports `"changed":false`.
8. MUST verify through Git, between the Spec's base commit and the audited head, that the only clause identifier removed from any module is `clause.bun.block-warnings-when-profile-treats-them-as-errors`, that every other clause present at the base keeps its identifier and enforcement, that every backend and frontend clause keeps its guidance, that the Source Baseline keeps every row it had at the base, that the only change to `internal/baseline/assets/profiles/go-cli-typescript-monorepo.json` is the removal of `rule.bun.warning-free-verification` from its `requiredRules`, and that no skill file and no setup snapshot changed.
9. MUST record, as evidence this Spec did not author, the adopter files the TechSpec's candidate table names for rows 1, 8, 9, 13 and 20, read only from `~/dev/secondbrain/projects/*/mirror/` and `~/dev/onioncry/`, and what oraculum's and vortex's lint scripts and Makefile Verification hold when read, from their mirrors or from `~/dev/oraculum` and `~/dev/vortex` read only. The row MUST also record onioncry's `src/` counts of `.unwrap()`, `.expect(`, `panic!(` and `anyhow` outside tests, and the published sources the PRD's Acceptance evidence names. When a source is unavailable or no longer carries the text, it MUST record that fact as the row's evidence or as a blocked row with its reason. The gate MUST NOT write to the Secondbrain, onioncry or any other repository.
10. MUST verify that this Spec's own artifacts satisfy the promise rule, and check whether the Spec introduced, changed or retired a glossary term that `CONTEXT.md` does not carry; the candidate is the promotion of a Repository-Specific Normative Rule into a Normative Clause.
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

- `newest="$(find 'docs/specs/0206-stack-rules-with-force-and-the-rules-adopters-repeat/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-5; Stories 1-6; Core Features 1-5; Success Metrics 1-8;
Declared breaks; Adopter notice; Prerequisites; Acceptance evidence;
`_techspec.md` → Candidate rules; Testing Approach 1-5; API Contracts 1-5;
ADR-0058; ADR-0059; ADR-0060; ADR-0073; ADR-0080; ADR-0081; ADR-0088;
ADR-0091; ADR-0093; ADR-0094; ADR-0096; ADR-0099; ADR-0103; ADR-0104;
ADR-0117; ADR-0130; ADR-0149; ADR-0155; ADR-0156; ADR-0166; ADR-0167;
ADR-0182; ADR-0186; ADR-0190; ADR-0194; ADR-0195; ADR-0202; ADR-0203.
