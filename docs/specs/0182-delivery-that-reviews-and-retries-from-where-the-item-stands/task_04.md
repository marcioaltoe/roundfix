---
task: task_04
spec: 0182-delivery-that-reviews-and-retries-from-where-the-item-stands
status: pending
type: qa
complexity: medium
---

# Task 04: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers and
settles the Spec on evidence. This Spec changes production code in
`internal/cli` and `internal/delivery` and the guides that describe the review,
Task Carry-Forward and the Delivery Retry. Every behavior row is exercised
through the built binary or by executing the named tests against the built
tree. Every binary command runs in a disposable repository with a disposable
Roundfix Home, and no command reaches a reviewer, a provider or the network.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify through the built binary, in a disposable repository whose
   default branch holds a commit the candidate branch does not share (an orphan
   branch passed as `--base`), that `roundfix review --base <orphan>` exits
   `2`, names both commits, says the head shares no history with the base, and
   writes no reviewer answer file.
3. MUST verify, by executing the review tests against the built tree, that with
   the default branch advanced after the candidate was cut the reviewer's diff
   carries only the candidate's change, the record's `baseCommit` is the merge
   base and `baseTipCommit` the tip, a main-side archived Spec edit leaves
   `archivedSpecs` empty while the candidate's own is listed, and a findings
   record from before main moved is reused with no reviewer call.
4. MUST verify, against this repository's history, that
   `git diff --name-status 159b41db 74a51f23 -- docs/specs docs/history/specs`
   lists Specs 0177–0180 as deleted and 0126/0127 as reactivated, while
   `git diff --name-status $(git merge-base 159b41db 74a51f23) 74a51f23` lists
   none of them. The gate MUST record this as evidence this Spec did not author:
   Spec 0175's reviewed head and the main tip it was reviewed against, written
   by other sessions on 2026-09-28. When either commit is unavailable, such as
   after the unreachable head was pruned, the gate MUST record the row as
   blocked with its reason.
5. MUST verify, by executing the carry-forward tests against the built tree,
   that a Task completed on the checkout is reported
   `already completed; nothing to carry` while the rest carries, a repeated
   carry exits `0` with `HEAD` unchanged, a moved input still refuses the
   remaining set, `TestCarryForwardRefusesRatherThanCarryingASubset` passes,
   and the implement Preflight names only ready Tasks and stays silent for a Run
   with nothing to carry.
6. MUST verify, by executing the delivery recovery and engine tests against the
   built tree, that a retry recording an older Run carries the newer Run's
   Tasks and records the newer Run, disjoint Runs carry newest first, a
   completed Run with a gone worktree is skipped, a refusal names the Runs
   already carried, and the retry output prints one line per carried Run.
7. MUST read, through the built `roundfix runs list --all --state all --limit 0`
   without writing to the live Run Database, the two Spec 0175 Runs
   `run_20260928T154312Z_63679c6540b34603` and
   `run_20260928T174529Z_c0cad0da5734a85f`. It MUST record their outcome and
   shared local branch as evidence this Spec did not author for the shape the
   retry now carries. When either Run is no longer listed, it MUST record the
   row as blocked with its reason.
8. MUST verify that `docs/user-guide/commands.md`, the Roundfix skill and its
   mirror, and `CONTEXT.md` describe the merge-base review with
   `baseTipCommit`, the `already completed; nothing to carry` action and the
   retry over every Run newest first, and that `make skills-sync-check` exits
   `0`.
9. MUST verify that this Spec's own artifacts satisfy the promise rule.
10. MUST verify from the repository history that the changed files stay within
    the paths the Tasks declare and that every governed path is bounded in
    `_authorization.md`.
11. MUST record the non-waivable Pull Request row as `blocked (environment: no open Pull Request)` and support it with the pre-PR equivalent of each control. This gate is a Task of the Spec's own graph and precedes the Pull Request by construction, so the controls a Pull Request would carry are observed as follows. Approval: the maintainer's recorded delivery authority, which is implementation authority and not review approval. Checks and status: the Daemon's repository Verification at the audited head. Unresolved review threads: none, because no Pull Request exists. Merge-Ready acceptance: none yet, because the configured pre-PR review follows archive. Review-artifact ancestry: the audited head named as the claimed candidate. A row with all five recorded this way satisfies the environment policy for `pass`; the delivery queue observes the real controls after the Pull Request opens.
12. MUST NOT accept a row whose only evidence is that a file was read.
13. MUST NOT write to the live Run Database under `~/.roundfix`.

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

- `newest="$(find 'docs/specs/0182-delivery-that-reviews-and-retries-from-where-the-item-stands/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-4; Core Features 1-4; Success Metrics 1-4; Acceptance
evidence; `_techspec.md` → Testing Approach 1-4; API Contracts 1-4; ADR-0014;
ADR-0026; ADR-0044; ADR-0052; ADR-0053; ADR-0057; ADR-0080; ADR-0090; ADR-0091;
ADR-0093; ADR-0094; ADR-0096; ADR-0104; ADR-0117; ADR-0139; ADR-0151; ADR-0153;
ADR-0155; ADR-0156; ADR-0158; ADR-0165; ADR-0169; ADR-0170.
