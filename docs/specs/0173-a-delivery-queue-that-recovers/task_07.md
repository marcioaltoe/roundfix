---
task: task_07
spec: 0173-a-delivery-queue-that-recovers
status: pending
type: qa
complexity: medium
---

# Task 07: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers and
settles the Spec on evidence. This Spec changes production code in
`internal/cli`, `internal/daemon`, `internal/delivery` and `internal/store`;
every behavior row below is exercised through the built binary or by executing
the named tests against the built tree.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify through the built binary, in a disposable repository, that
   `roundfix reconcile <run-id> --carry-forward` carries both Tasks of a Run
   whose Tasks settled out of Task Graph order, one declaring an input the other
   changed, and that a checkout change to that input still refuses the whole set
   with exit `2` and HEAD unchanged.
3. MUST verify, by executing the Daemon tests against the built tree, that a
   Task whose frozen authorization bounds `.roundfixrc.yml` commits it, that an
   unbounded Task settles `failed` with reason `Task commit lost output:
   .roundfixrc.yml (Project Config outside the Spec's authorization)` and a
   dropped-path event, and that Batch and QA Report commits omit it with their
   published reason.
4. MUST verify, by executing the engine and store tests against the built tree,
   every row of the re-entry table, the refusals that leave the item unchanged,
   the recorded Run ID of a `run-unresolved` park, the guarded item transition,
   the idle owner release, and that a retried archived item pushes, opens its
   pull request and merges exactly once.
5. MUST verify through the built binary, in a disposable repository with a
   disposable Roundfix Home, that `roundfix deliver retry <slug>` refuses a
   missing slug, an unknown slug and an item that is not parked with exit `2`
   and the queue unchanged as `roundfix deliver status` prints it, and that
   `roundfix deliver --help` and `roundfix --help` name the command.
6. MUST verify, by executing the delivery workflow and command tests against the
   built tree, that a retried `run-unresolved` item carries its Run's completed
   Tasks onto the item branch and runs only its unfinished Tasks, that the retry
   starts an owner, hands the item to a live owner, or reclaims a stale owner
   record, and that the owner makes a further pass for an item retried during
   its pass.
7. MUST verify through the built binary that `roundfix profiles validate --json`
   with a proof that times out twice reports classification `temporary` and
   rerun advice, or record the row as blocked with its reason when no adapter
   can be made to time out in the QA environment, backed by the executed
   `internal/cli` proof tests.
8. MUST verify that the Run Event Journal of Run
   `run_20260925T182023Z_7e3e8c4b8bc65d93` in the live Run Database, replayed
   through the built `roundfix events`, shows `task_04` settling completed
   before `task_01`, and that the live Delivery Queue, read through the built
   `roundfix deliver status` and a read-only query of `delivery_queue_items`,
   shows the three wave-one items parked `run-unresolved` with no Run ID,
   recording both as evidence this Spec did not author, or the row as blocked
   with its reason when the Run Database is absent or unreadable; the QA MUST
   NOT write to the live Run Database.
9. MUST verify that `docs/user-guide/commands.md`, the Roundfix skill and its
   mirror, and `CONTEXT.md` describe the integration order, the Project Config
   rule, the Delivery Retry and its hand-off, and the proof retry.
10. MUST verify that this Spec's own artifacts satisfy the promise rule.
11. MUST verify from the repository history that the changed files stay within
    the paths the Tasks declare and that every governed path is bounded in
    `_authorization.md`.
12. MUST NOT accept a row whose only evidence is that a file was read.

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

- `newest="$(find 'docs/specs/0173-a-delivery-queue-that-recovers/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-5; Core Features 1-6; Success Metrics 1-5; Acceptance
evidence; `_techspec.md` → Testing Approach 1-7; API Contracts 1-7; ADR-0014;
ADR-0026; ADR-0044; ADR-0052; ADR-0053; ADR-0057; ADR-0080; ADR-0091;
ADR-0093; ADR-0096; ADR-0104; ADR-0107; ADR-0117; ADR-0130; ADR-0138;
ADR-0153; ADR-0155; ADR-0156; ADR-0157; ADR-0158; ADR-0160.
