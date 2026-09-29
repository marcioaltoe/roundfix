---
task: task_05
spec: 0181-gates-that-refuse-only-what-someone-can-act-on
status: pending
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix covers and settles the Spec on evidence. The Spec changes production code in `internal/spec`, `internal/daemon`, `internal/cli`, `internal/agent` and `internal/speccheck`, and it aligns the `qa-gate`, `write-tasks` and Roundfix skills. It rests one row on history and fleet evidence it did not author.

Every behavior row below is exercised through the built binary, or by executing the named tests against the built tree. Every binary command that writes runs in a disposable clone or repository with a disposable Roundfix Home. The gate starts no Run and writes nothing to the Secondbrain.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify, by executing the spec, daemon and settle tests task_01 names against the built tree, that:
   - a completed Task declaring `interface: internal/store/delivery.go` that also changes `internal/store/delivery_test.go` commits `## Recorded paths` naming exactly that test file, with the same `recorded_paths` in its `daemon.commit` event;
   - a new package records its files;
   - a Task with every path declared, a failed Task, a QA Task and a Governed Path record nothing;
   - the Task file's Context and `CarryForwardInputs` are unchanged by the section;
   - `roundfix settle` records the same way.
3. MUST verify through the built binary that `roundfix qa-report accept` behaves as follows:
   - it exits `0` on Spec 0179's archived `qa-report-2026-09-29.md`;
   - in a disposable Spec directory holding one Unreachable Acceptance declaration, a declared `partial` with one declared row, the pre-PR Pull Request row and one other environment-blocked row exits `1` naming `1 outside the pre-PR Pull Request row`;
   - the same report without the other row exits `0`;
   - a report whose no-PR row does not name the Pull Request row in its provenance exits `1` with `rows_blocked_environment is 1; expected 0`.
4. MUST verify, on evidence this Spec did not author, the outside-evidence row. It MUST record where each piece came from.
   - Run the built `roundfix qa-report accept` on Oraculum Spec 0027's report as mirrored at `~/dev/secondbrain/projects/oraculum/mirror/docs/history/specs/0027-agregadores-de-vendas/qa/qa-report-2026-08-07-02.md`. It MUST still exit `1`, because that report has environment-blocked rows other than the Pull Request row.
   - Count how many QA reports under the Secondbrain mirrors of projects other than Roundfix contain `blocked (environment: no open Pull Request)`, and in how many projects.
   - Read the Secondbrain read-only. When it is unavailable, record the row as blocked with that reason.
5. MUST replay the ADR-0161 cascade in a disposable clone of this repository at `ebeb997f`. Commit `docs/adr/0161-a-merged-spec-releases-its-runs-on-the-merged-head.md` from `06afa835` in a later commit, then run `spec check 0180-a-prepared-queue-that-revalidates-before-each-spec --strict` twice:
   - a binary built from this Spec's Delivery Base MUST report `SC-ADR-RELATED` for ADR-0161;
   - the candidate binary MUST report no `SC-ADR-RELATED` and exit `0`.

   Then copy that Spec folder under a new, uncommitted slug. The candidate binary MUST report `SC-ADR-RELATED` for ADR-0161 on the copy. When the history is unavailable, such as in a shallow clone, record the row as blocked with its reason.
6. MUST verify, by executing the eligibility, prompt, horizon and citation projection tests task_02, task_03 and task_06 name against the built tree, the remaining negative cases:
   - a `pass` unchanged;
   - an ADR committed before or with the PRD still reported;
   - a directory without Git kept on the full check;
   - `SC-ADR-UNLISTED` still reported for a citation in authored text;
   - no `SC-ADR-UNLISTED` for an unlisted ADR cited only in a Task's `## Result`, in a Daemon-owned section or under `qa/`.
7. MUST verify that the `qa-gate`, `write-tasks`, `archive-spec` and Roundfix skills, their mirrors, `docs/user-guide/commands.md`, `docs/user-guide/context-driven-development.md` and `CONTEXT.md` describe the recorded paths, the pre-PR Pull Request row and the horizon, that both skills declare `0.0.3`, and that `make skills-sync-check` exits `0`.
8. MUST verify that this Spec's own artifacts satisfy the promise rule, and check whether the Spec introduced, changed or retired a glossary term that `CONTEXT.md` does not carry.
9. MUST verify from the repository history that each Task's changed files stay within its declarations, and that every Governed Path is bounded in `_authorization.md`. The following also count as declared, and the scope row MUST name each one it relies on:
   - a path in the Task file's `## Recorded paths`;
   - a path rewritten by a sanctioned regeneration command that `_authorization.md` names;
   - a `_test.go` file that the Task's `## Result` names as a test the change invalidated.

   The last applies because this Spec's Run executes on a Daemon that does not record paths yet. It MUST fail the row for any other undeclared path and for any unbounded Governed Path.
10. MUST NOT accept a row whose only evidence is that a file was read.
11. MUST NOT write to the live Run Database under `~/.roundfix` or to any path under `~/dev/secondbrain`.

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

- `newest="$(find 'docs/specs/0181-gates-that-refuse-only-what-someone-can-act-on/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-4; Core Features 1-4; Success Metrics 1-4; Acceptance
evidence; `_techspec.md` → Testing Approach 1-6; API Contracts 1-5; ADR-0014;
ADR-0057; ADR-0080; ADR-0088; ADR-0091; ADR-0093; ADR-0094; ADR-0096;
ADR-0097; ADR-0104; ADR-0116; ADR-0117; ADR-0130; ADR-0155; ADR-0156;
ADR-0166; ADR-0167; ADR-0168; ADR-0176.

## Result
