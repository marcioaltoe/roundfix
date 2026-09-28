---
task: task_07
spec: 0175-cleanup-after-a-squash-merge
status: completed
type: qa
complexity: medium
---

# Task 07: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers and
settles the Spec on evidence. This Spec allows production changes in
`internal/worktree`, `internal/store`, `internal/cli`, `internal/delivery`,
`internal/daemon/reconcile.go` and `internal/speccheck/checkout.go`, plus the
guides, the Roundfix skill, `CONTEXT.md` and ADR-0161.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify, in a disposable repository driven through the built binary with
   a seeded Run Database, that `roundfix reconcile --format json` classifies
   the Spec 0172 shape `superseded` with evidence naming the merged head. The
   shape is Task 01–04 commits contained in the merged head, a unique redone
   Task 05 commit and a unique failed QA Report commit. The check runs once
   with a Delivery Queue merge record and once, without it, against a default
   branch carrying the archived Spec. `roundfix reconcile --apply` must then
   remove its Run Worktree and Run Branch, and a dry-run must remove nothing.
3. MUST verify through the built binary that a Run holding one unrepresented
   commit, a Run whose Task is not completed at the merged head, and a Run
   covered only by another Spec's merge record are each kept, with a reason
   naming the cause.
4. MUST verify, by executing the named merged-head tests against the built
   tree, the Spec 0164 shape: an archived Spec, and a default branch that later
   changed other files and a file the Run touched. That Run must be released.
5. MUST verify, by executing the delivery engine and workflow tests against the
   built tree, the order and outcomes after a delivery merge. The Spec's
   provable Runs are released before the item worktree and item branch are
   removed. A kept Run appears in the item's cleanup warning. An Active Run and
   another Spec's Run are untouched.
6. MUST verify through the built binary that an unkeyed legacy Run whose
   checkout was removed is listed by `roundfix runs list` and `roundfix
   reconcile` from the main checkout after a write-mode open. A Run keyed to
   another repository is not listed.
7. MUST verify through the built binary that a staging worktree left `locked
   initializing`, and one whose owner process is dead, are reported in
   `stagingCandidates` by dry-run and released by `--apply`. A live owner's
   staging must be kept in `preservedCandidates`, and every existing report
   field must keep its name.
8. MUST verify, by executing the item-cleanup and guard tests against the built
   tree, two properties. A half-removed item holding a nested registered
   worktree is refused and the nested worktree keeps its files. No production
   file outside `internal/worktree` runs `git worktree` add, remove, prune or
   move.
9. MUST record the built binary's read-only `roundfix reconcile --format json`
   in the maintainer's checkout at `/Users/marcio/dev/roundfix`, naming how each
   Run of Specs 0155–0164 cited by the adopted Finding that still exists is now
   classified, as evidence this Spec did not author. The row is recorded as
   blocked with its reason when none of those Runs remains. The gate MUST NOT
   run `--apply`, `--discard-superseded` or `--carry-forward` in that checkout.
10. MUST verify that `docs/user-guide/commands.md`, the Roundfix skill and its
    mirror, `CONTEXT.md` and ADR-0161 describe the merged-head proof, the
    automatic release, the cleanup warning and the staging sweep.
11. MUST verify that this Spec's own artifacts satisfy the promise rule.
12. MUST verify Project Constraint applicability, the operative source paths,
    the tooling authorization and the actual changed-file scope from Git
    evidence against `_authorization.md`.
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

- `newest="$(find 'docs/specs/0175-cleanup-after-a-squash-merge/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-5; Core Features 1-6; Success Metrics 1-5; Acceptance
evidence; `_techspec.md` → Testing Approach 1-7; API Contracts 1-6; ADR-0004;
ADR-0014; ADR-0023; ADR-0044; ADR-0052; ADR-0053; ADR-0057; ADR-0080;
ADR-0091; ADR-0093; ADR-0096; ADR-0097; ADR-0104; ADR-0115; ADR-0117;
ADR-0130; ADR-0155; ADR-0156; ADR-0158; ADR-0161; ADR-0163.
