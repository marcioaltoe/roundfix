---
task: task_05
spec: 0193-baseline-guides-that-describe-the-product-as-it-is
status: pending
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers and settles the Spec on evidence. This Spec changes Baseline modules, the guides rendered from them, and two small readers of Backlog status. Every behavior row is exercised by executing the named tests against the built tree, by reading the rendered guides, or by running the built binary's managed refresh in a disposable copy of this repository. No command reaches a provider or a live remote.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify, by executing task_01's tests against the built tree, that no two clauses share text, that a repeated text is reported, that every clause keeps its enforcement level, and that the backend golden holds the boundary paragraph once.
3. MUST verify, by executing task_02's tests and reading `docs/agents/autonomous-work.md`, that the loop clause declares review before archive, names `roundfix deliver`, `roundfix reopen --spec <slug>` and Task Carry-Forward, and that an archive-before-review order is reported as a mismatch with the Delivery Queue.
4. MUST verify, by executing task_03's tests and reading `docs/agents/agent-instructions.md` and `docs/agents/spec-routing.md`, that each of the five clauses carries its new sentence and that none of the four removed phrases remains in any module, golden or rendered guide.
5. MUST verify, by executing task_04's tests, that no shipped guidance cites a Spec number or an ADR number, that a citation is reported, and that a `deferred` Backlog Entry is terminal.
6. MUST verify through the built binary, in a disposable copy of this repository, that `roundfix baseline update --repo <copy> --no-skills --format text` reports `File changes: 0`, and that `make baseline-digests` reports `"changed":false`.
7. MUST record, as evidence this Spec did not author, the adopter guides in the Secondbrain mirrors, read only: `~/dev/secondbrain/projects/conexus/mirror/docs/agents/backend.md` (the boundary paragraph twice), and `~/dev/secondbrain/projects/oraculum/mirror/docs/agents/autonomous-work.md` (archive before review, and `Spec 0078`). The row records what those files held when read and where they came from. It MUST also record Google's documentation best practices, "Update Docs with Code" (<https://github.com/google/styleguide/blob/gh-pages/docguide/best_practices.md>). When a mirror or the page is unavailable, or a mirror already carries the corrected text, it MUST record that fact as the row's evidence or as a blocked row with its reason. The gate MUST NOT write to the Secondbrain.
8. MUST verify that every clause this Spec rewrote keeps its enforcement level and identifier, and that the clauses the TechSpec lists as deliberately unchanged are byte-identical to `9e439dbb`.
9. MUST verify that this Spec's own artifacts satisfy the promise rule, and check whether it introduced, changed or retired a glossary term that `CONTEXT.md` does not carry.
10. MUST verify from the repository history that each Task's changed files stay within its declarations or its `## Recorded paths`, that every Governed Path is bounded in `_authorization.md`, and that every regenerated file is an output of a sanctioned command.
11. MUST record the non-waivable Pull Request row as `blocked (environment: no open Pull Request)`, naming the Pull Request row in its provenance, and support it with the pre-PR equivalent of each control:
    - approval: the maintainer's recorded delivery authority;
    - checks and status: the Daemon's repository Verification at the audited head;
    - unresolved threads: none, because no Pull Request exists;
    - Merge-Ready acceptance: none yet, because the pre-PR review precedes archive and publication;
    - review-artifact ancestry: the audited head named as the claimed candidate.
12. MUST NOT accept a row whose only evidence is that a file was read, except the outside-evidence row, whose subject is what an adopter's file holds.
13. MUST NOT write to the live Run Database under `~/.roundfix`.

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

- `newest="$(find 'docs/specs/0193-baseline-guides-that-describe-the-product-as-it-is/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-4; Core Features 1-4; Success Metrics 1-6; Acceptance
evidence; `_techspec.md` → Testing Approach 1-5; API Contracts 1-4; ADR-0014;
ADR-0073; ADR-0080; ADR-0081; ADR-0088; ADR-0091; ADR-0092; ADR-0093; ADR-0094;
ADR-0096; ADR-0103; ADR-0104; ADR-0117; ADR-0130; ADR-0149; ADR-0155; ADR-0156;
ADR-0163; ADR-0166; ADR-0167; ADR-0179; ADR-0186.

## Result
