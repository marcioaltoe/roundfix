---
task: task_06
spec: 0180-a-prepared-queue-that-revalidates-before-each-spec
status: pending
type: qa
complexity: medium
---

# Task 06: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers and
settles the Spec on evidence. This Spec changes production code in
`internal/cli`, `internal/delivery` and `internal/store`, adds one Run Database
schema version, and rewrites the owned `implement-spec` skill. Every behavior
row below is exercised through the built binary or by executing the named tests
against the built tree. Every binary command runs in a disposable repository
with a disposable Roundfix Home. The only queue the gate starts has a deadline
already passed, so its owner parks the item without a worktree and never starts
an Agent.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify through the built binary, in a disposable repository holding one
   approved Spec and one Spec whose authorization lacks `merge`:
   - `roundfix deliver plan` with both slugs exits `1`, reports the first
     `approved` and the second `blocked` with `authorization lacks merge`, and
     `--json` carries the same verdicts;
   - the plan creates no Run Database and leaves the checkout's HEAD and
     status unchanged;
   - `roundfix deliver start` with both slugs exits `2`, and `roundfix deliver
     status` then reports that no queue exists;
   - `roundfix deliver --help` and `roundfix --help` name the plan command.
3. MUST verify, by executing the engine and workflow tests against the built
   tree, that:
   - an item whose starting main fails its strict check parks as
     `revalidation-failed`, and one whose declared production Go file an
     earlier item's merge changed parks as `premise-changed`, both before any
     Run;
   - a clean item runs;
   - a retry of an item with no Run refuses while findings remain and
     acknowledges a changed premise.
4. MUST verify, against this repository's history, the intersection of `git
   diff --name-only 6fac37ea^ 6fac37ea` with the production Go `interface:`
   paths that Spec 0175's Tasks declare at `b92aefda` (read through `git show
   b92aefda:<path>`). It must be exactly `internal/cli/carryforward.go`,
   `internal/cli/deliver_workflow.go` and `internal/delivery/engine.go`, and the
   gate MUST record it as evidence this Spec did not author, beside the premise
   the revalidation names for the same shape. When the history is unavailable,
   such as in a shallow clone, record the row as blocked with its reason.
5. MUST verify through the built binary, with the approved Spec:
   - `roundfix deliver start --max-duration 1ns --max-retries 1 <slug>` prints
     a `Limits:` line naming a deadline and `retries per item 1`;
   - once its owner has exited, `roundfix deliver status` shows the item parked
     as `queue-deadline` with no worktree, the same `Limits:` line, one
     `Pending question:` line and an `Answer:` line naming `roundfix deliver
     start`;
   - `roundfix deliver retry <slug>` then exits `2` and leaves the status
     output unchanged;
   - a start with `--max-retries 0` or `--max-duration 0s` exits `2` and records
     no queue.
6. MUST verify, by executing the store, engine, question and command tests
   against the built tree, the limits round trip, the migration of the previous
   schema version, a retry refused at its limit with the item unchanged, a
   started item advancing past the deadline, one Pending Question for two
   parked items, and the question unchanged after an owner pass with the clock
   advanced.
7. MUST verify that the `implement-spec` skill and its mirror name `roundfix
   deliver plan`, `roundfix implement --spec` and `roundfix deliver start`,
   carry no implement-task loop, and that `make skills-sync-check` exits `0`.
8. MUST verify that `docs/user-guide/commands.md`, the Roundfix skill and its
   mirror, and `CONTEXT.md` describe the Delivery Plan, the start refusal, the
   Delivery Revalidation and its two blockers, the limit flags and the Pending
   Question.
9. MUST verify that this Spec's own artifacts satisfy the promise rule.
10. MUST verify from the repository history that the changed files stay within
    the paths the Tasks declare and that every governed path is bounded in
    `_authorization.md`.
11. MUST NOT accept a row whose only evidence is that a file was read.
12. MUST NOT write to the live Run Database under `~/.roundfix`.

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

- `newest="$(find 'docs/specs/0180-a-prepared-queue-that-revalidates-before-each-spec/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-6; Core Features 1-6; Success Metrics 1-5; Acceptance
evidence; `_techspec.md` → Testing Approach 1-7; API Contracts 1-8; ADR-0014;
ADR-0044; ADR-0052; ADR-0057; ADR-0080; ADR-0091; ADR-0093; ADR-0094;
ADR-0096; ADR-0104; ADR-0117; ADR-0130; ADR-0137; ADR-0139; ADR-0153; ADR-0155; ADR-0156; ADR-0158; ADR-0160.

## Result
