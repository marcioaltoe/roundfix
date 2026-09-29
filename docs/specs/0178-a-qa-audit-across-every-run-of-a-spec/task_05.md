---
task: task_05
spec: 0178-a-qa-audit-across-every-run-of-a-spec
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers and
settles the Spec on evidence. This Spec allows these production changes:

- the QA step's Delivery Base, Task commit selection, staleness warning event
  and report section, settlement checks and self-audit prompt line in
  `internal/daemon/`;
- the mechanical stage's earlier-Runs skip and the authorization audit table in
  `internal/speccheck/mechanical.go` and `internal/speccheck/report.go`;
- the auditor evidence, the `user_flow_binary` key and the Precondition Refusal
  report's auditor identity in `internal/spec/`;
- the staleness reasons in `internal/app/version.go`.

Every row is judged against that scope.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify, by executing the built tree's Task 01 tests and by replaying
   `qaMechanicalRequest` in a disposable repository with an `origin/HEAD`
   symbolic ref, that a Task commit made by an earlier Run and an older commit
   of a Task each produce a `QA-AUTH-PATHS` finding naming that commit, and that
   a commit already on the default branch is not audited.
3. MUST verify in the same way that a grant widened on the Spec branch does not
   authorize the widened path, and the mirror operation: the same widening
   landed on the default branch and merged into the Spec branch does.
4. MUST verify that without a determinable default branch the stage keeps the
   Run start head and the seeded report's mechanical skips name
   `Task commits of earlier Runs`.
5. MUST verify that a seeded QA Report's authorization audit table carries the
   `Commit` column with one row per audited commit, and that the qa-gate skill
   and its mirror tell the gate to audit by command only a Task commit the table
   does not list.
6. MUST verify, with an injected Auditing Binary in a disposable repository,
   that a Daemon built from the Delivery Base seeds `auditor_staleness` as
   `current: commit ancestry: build commit does not predate the delivery base`
   while the audited head carries candidate commits and publishes no staleness
   warning, and that a Daemon built from an ancestor of the Delivery Base
   publishes exactly one `daemon.qa` event with phase `auditor_staleness`
   naming the staleness line and the rebuild instruction, records the
   `## Auditor staleness warning` section in the seeded report, writes no
   `precondition_check`, and still runs repository Verification and the QA
   Agent; and that an `unknown` auditor publishes no warning.
7. MUST verify that in a self-audit Daemon settlement accepts a `pass` whose
   `user_flow_binary` is the built binary's `roundfix --version` line from the
   audited head, and refuses a missing `user_flow_binary`, one from another
   commit, one without a build commit, and a rewritten `auditing_binary` or
   `auditor_staleness`, each with a reason naming its cause; that outside a
   self-audit the key is not checked; and that the self-audit QA prompt names
   `user_flow_binary`.
8. MUST verify through the built binary that `roundfix qa-report accept` and
   `roundfix archive` keep today's acceptance: every archived `pass` Spec in
   `docs/history/specs/` stays archive-eligible, and a `pass` report without
   `user_flow_binary` is accepted by `roundfix qa-report accept` in a disposable
   Spec.
9. MUST replay, in a disposable clone, Spec 0172's retained history
   `256ad156..ebac7cbe` with `origin/main` set to a revision of the default
   branch that contains `fc296df0`: record every Task commit the built tree's
   selection lists — every Task commit of that history that changes a
   governed path, which the authorization audit selects by design (task_01,
   task_02, task_04 and task_05, including those before `0532c0fe chore: merge
   main into 0172`; task_03 and task_07 change no governed path) and every `QA-AUTH-PATHS` finding it
   produces for the in-PR grant widening `c0818081`. MUST also recompute the
   staleness of the Daemon builds recorded by the QA Reports of Specs 0171,
   0172 and 0174 against each audited head's Delivery Base and record each
   result (0171 and 0174 `current`, 0172 `stale`). MUST record this history
   and these reports as evidence this Spec did not author, or the row as
   blocked with its reason when those objects are absent.
10. MUST verify that `CONTEXT.md` carries the Delivery Base entry and that its
    Auditing Binary and QA Report entries describe the delivery-base staleness,
    the stale warning, the Daemon-owned auditor fields and `user_flow_binary`.
11. MUST verify from Git evidence that the only changed governed paths are the
    two qa-gate skill files listed in `_authorization.md`, and that
    `internal/speccheck/mechanical_test.go` and
    `docs/references/coverage-record.json` are unchanged.
12. MUST verify that this Spec's own artifacts satisfy the promise rule.
13. MUST NOT accept a row whose only evidence is that a file was read.

## Subtasks

- [ ] Build the matrix from the Requirements above and the sources no
      declaration waives.
- [ ] Execute each row against the built tree and record its evidence.
- [ ] Write the dated QA Report with its verdict.

## Acceptance Criteria

- [ ] The QA Report records a verdict and names, in each row's provenance, the
      sources that row covers.
- [ ] The repository Verification result is recorded as a fact, not re-derived.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0178-a-qa-audit-across-every-run-of-a-spec/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-5; Core Features 1-4; Success Metrics 1-5; Acceptance
evidence; `_techspec.md` → Testing Approach 1-5; API Contracts 1-6; ADR-0014;
ADR-0015; ADR-0023; ADR-0057; ADR-0080; ADR-0089; ADR-0091; ADR-0093; ADR-0096;
ADR-0097; ADR-0104; ADR-0117; ADR-0130; ADR-0132; ADR-0138; ADR-0155; ADR-0156.
