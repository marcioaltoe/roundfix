---
task: task_05
spec: 0139-a-suite-that-passes-where-it-runs
status: completed
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

Execute the `qa-gate` skill against the assembled tree and write the canonical
dated report under `qa/`. The matrix is declared by the Requirements below: each
Requirement that starts with `MUST verify` or `MUST run` is a row source, and the
remaining Requirements constrain how the gate runs.

## Requirements

1. MUST verify every Acceptance Criterion of Tasks 01, 02, 03 and 04 against the
   assembled tree, recording the observed result for each.
2. MUST verify the outside-evidence row against sources this Spec did not write,
   naming each:
   - Read Spec 0138's QA Report of 2026-09-15 at commit `990bb65b` on the kept Run
     Branch `roundfix/run-run_20260915T105843Z_41c059d11301702f`, and record its
     finding F-001.
   - Run `make verify` in a detached worktree at delivery target revision
     `7a15e970`, then at the assembled tree, and record both exit statuses. The
     base is expected to fail in `internal/agent`, `internal/daemon` and
     `internal/speccheck`, and the assembled tree to pass.

   When the commit, branch or revision no longer resolves, record the row as
   blocked, with that reason.
3. MUST verify that `TestAuditRefusesAGrantWidenedAfterItsConsumingCommit` runs
   rather than skips in a fresh `git clone --no-local` of the assembled branch.
4. MUST run `go test -count=5 -parallel 16 ./internal/agent`, and
   `go test -count=3 -parallel 16 ./internal/daemon -run 'TestTaskCycle'` while
   `go test ./internal/baseline` runs concurrently, and record both clean.
5. MUST run the repository's own verification gate, `make verify`, and record it
   clean with the resolved Go toolchain version. When the QA prompt states that
   the Daemon already ran it outside the Agent sandbox, record that observation
   and its evidence instead of running it again.
6. MUST run `go vet ./internal/agent ./internal/daemon ./internal/speccheck`, the
   packages this Spec changes, and compare the result against the delivery
   target. A diagnostic identical on the delivery target is recorded as observed
   and pre-existing, naming Spec 0123 for the runner mutex copying.
7. MUST verify whether the work introduced, changed or retired a term the project
   glossary should carry, and record the result.
8. MUST exercise the gate against a binary rebuilt from the assembled tree.
9. MUST NOT change implementation code, tests or any sibling Task.

## Subtasks

- [ ] Rebuild the binary the gate exercises.
- [ ] Execute the Acceptance Criteria of Tasks 01-04.
- [ ] Measure `make verify` at the base and at the assembled tree.
- [ ] Prove the audit fixture runs in a fresh clone.
- [ ] Run the stressed package runs and scoped `go vet`.
- [ ] Write the dated report with exact counters and the terminal verdict.

## Acceptance Criteria

- [ ] Every Acceptance Criterion of Tasks 01-04 has an observed result with
      recorded evidence.
- [ ] The before and after `make verify` measurement is recorded with both exit
      statuses, or recorded as blocked with its reason.
- [ ] The audit fixture test is shown running in a fresh clone.
- [ ] The stressed ACPX and Task-cycle runs are recorded clean.
- [ ] The repository Verification is recorded clean. When the Daemon ran it, its
      evidence is cited.
- [ ] Scoped `go vet` is recorded, with pre-existing diagnostics attributed.
- [ ] The glossary check result and the terminal verdict with exact counters are
      recorded.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`
- instruction: `docs/agents/domain.md`

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It checks the report marker only; the matrix above remains the QA Agent's
responsibility.

- `newest="$(find 'docs/specs/0139-a-suite-that-passes-where-it-runs/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && verdict == \"pass\" ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals 1-4; Acceptance evidence; Success Metrics; Regression locks.
- `_techspec.md` → Coverage Map; Testing Approach 5; Risks & Considerations:
  Bootstrap; Build Order 5.
