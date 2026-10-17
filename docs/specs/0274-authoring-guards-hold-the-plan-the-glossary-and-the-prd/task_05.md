---
task: task_05
spec: 0274-authoring-guards-hold-the-plan-the-glossary-and-the-prd
status: completed
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec adds two Spec Consistency Check
gaps, `SC-SURFACE-RECORD-UNTASKED` and `SC-DECISIONS-UNRECORDED`, and two
Repository Contract Tests, one keeping the documented QA secret refusal in
step with the Daemon and one holding `CONTEXT.md` to five structural rules,
with the skills, guides and glossary that describe them. No row reads or
writes the real `~/.roundfix`, reaches a provider, reads a credential or reads
another repository.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1 through the built `bin/roundfix` with `HOME`
   set to a temporary home and `NODE_OPTIONS` unset: export the tree of commit
   `e1e376dd` with `git archive` into a temporary directory, commit it into a
   fresh temporary repository, run
   `roundfix spec check 0268-a-turn-that-did-not-finish-is-not-accepted --strict`
   and record that it reports `SC-SURFACE-RECORD-UNTASKED` for task_01 at a
   Context line and at the Limits line and exits 1; then apply the operator's
   fix of 2026-10-14 to that copy (declare
   `docs/references/behavior-surfaces.json` under task_01's Context, add a
   requirement naming the `-record-skill-coverage` step, and drop the record
   from the Limits of `_authorization.md`), commit it, and record that the code
   is absent (API Contract 1). The fix commit is no longer reachable from any
   ref, so the row does not read it.
3. MUST verify Success Metric 2 on a temporary copy of the audited tree made
   with `git archive`: replace the colon form with
   `<path>:<line> (<kind>)` in `.agents/skills/roundfix/references/implement.md`
   and run `TestTheQASecretRefusalIsDocumentedAsTheDaemonPrintsIt`, then do the
   same in a fresh copy for `docs/user-guide/commands/implement.md`, recording
   that each run fails naming the file and that the unmodified audited tree
   passes (API Contract 2, Invariant 4).
4. MUST verify Success Metric 3: run the six glossary tests task_03 names on
   the audited tree, and on a temporary copy plant, one copy per rule, an
   undefined compound in a definition, a link to a file under `docs/specs/`, a
   Task Type list missing `qa`, a sentence in an `_Avoid_` line and one
   sentence of more than sixty characters repeated in a second entry,
   recording that exactly the matching rule's test fails each time (API
   Contract 4, Invariant 5). It MUST record the allowlist's size and each
   entry's reason.
5. MUST verify Success Metric 4 through the built binary on a temporary Git
   repository holding `.agents/skills/write-prd/references/clarification-record.md`
   added in a commit after one PRD and before three others: the earlier PRD
   without the record is skipped with the horizon text; a later PRD without
   `## Decisions` and one whose section holds neither shape each yield
   `SC-DECISIONS-UNRECORDED` at `--stage prd` and exit 1 with `--strict`; a
   later PRD with a dated answer and one with the autonomous-route line yield
   none (API Contract 3, Invariant 3).
6. MUST verify Success Metric 5 and Invariant 6: record that every count in
   `internal/docscontract/testdata/corpus-golden.json` at the delivery base,
   the merge base of the audited head with `main`, is unchanged at the audited
   head, that the two new codes appear with the counts the corpus test
   reports, and that `TestCheckCorpusGolden` and
   `TestCheckActiveCorpusHasNoErrors` pass.
7. MUST execute the tests task_01 to task_04 name in their Verification, and
   prove each new gate fails without its implementation: on a temporary copy,
   remove the `detectSurfaceRecord` call and then, in a fresh copy, the
   `detectClarificationRecord` calls, and record that the named detector
   tests fail.
8. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the Go Blog passage on testable examples,
   the Vale vocabularies passage on `reject.txt` and the MADR template's
   `decision-makers` field, quoting each and stating whether it still supports
   the design; for a source it cannot reach it records
   `blocked (environment: network denied: <host>)`.
9. MUST perform the glossary check of `docs/agents/domain.md`: record that
   `CONTEXT.md` defines **Clarification Record**, that it passes the glossary
   contract, and whether any other term the Spec introduced needs one.
10. MUST verify from the repository history that each Task's changed files stay
    within its declarations or its `## Recorded paths`; that every Governed
    Path changed is bounded in `_authorization.md`; that the Makefile, the CI
    workflows, `go.mod`, `go.sum`, every Baseline module,
    `internal/daemon/task_engine.go` and `internal/secretscan/scan.go` did not
    change; that the Roundfix, write-tasks, write-prd and write-techspec skill
    versions were raised and recorded; that no `### QA settlement` section
    changed; that the user guide changes and skill references name both
    codes; and that the glossary defines **Clarification Record**, the one
    term the PRD's Glossary declares (no Task adds a check code to
    `CONTEXT.md`; amended 2026-10-17 after the first QA gate read this row as
    asking the glossary for both codes). It also searches the repository's guides and skill copies
    for `<line> (<kind>)` and records that none remains. This row reads Task
    commits, so it declares a `commit_range` input.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance, and support it with the pre-PR equivalent of each
    control: approval from the recorded delivery authority, checks from the
    Daemon's repository Verification at the audited head, no unresolved
    threads because no Pull Request exists, no Merge-Ready acceptance yet, and
    the audited head named as the claimed candidate.
12. MUST NOT accept a row whose only evidence is that a file was read. It MUST
    NOT read or write `~/.roundfix`, record any environment value, send any
    request to a provider, or start a workflow.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Replay the 0268 plan and the Clarification horizon through the built binary.
- [ ] Sabotage each guard on temporary copies and record the failures.
- [ ] Record the outside evidence, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The 0268 plan reports the surface record gap and its fix does not.
- [ ] Each reverted document, planted glossary violation and removed detector makes its guard fail.
- [ ] No earlier corpus count moved.
- [ ] The scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0274-authoring-guards-hold-the-plan-the-glossary-and-the-prd/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-4; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; Invariants 1-6; Current behavior; Build Order 5
- ADR-0280; ADR-0256; ADR-0244; ADR-0273; ADR-0104; ADR-0167; ADR-0240; ADR-0276
