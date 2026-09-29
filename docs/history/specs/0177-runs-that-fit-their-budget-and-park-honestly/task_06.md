---
task: task_06
spec: 0177-runs-that-fit-their-budget-and-park-honestly
status: completed
type: qa
complexity: medium
---

# Task 06: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers and
settles the Spec on evidence. The Spec changes production code in
`internal/spec`, `internal/daemon`, `internal/cli`, `internal/config`,
`internal/delivery` and `internal/speccheck`, and the guides and skills named in
its Tasks; the gate verifies each through the built binary or by executing its
tests against the built tree.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify through the built binary, in a disposable Spec Root, that
   `roundfix spec check --strict` reports no `SC-WAVE-COLLISION` for two
   unordered Tasks that share only an `instruction:` path, and still reports
   one when they share an `interface:` path.
3. MUST verify by executing the Daemon and CLI budget tests against the built
   tree that a serial graph longer than one allowance ends `Clean`, that a
   settlement after the deadline does not renew and the Run ends
   `BudgetExceeded` with earlier work kept, that a stalled Task is cancelled
   one allowance after the last settlement, and that integration runs under the
   renewed deadline.
4. MUST verify through the built binary that `roundfix init --scope user` in a
   disposable home writes a User Config stating that an Implement Run's
   allowance renews at each Task settlement.
5. MUST verify by executing the delivery tests against the built tree that a
   `BudgetExceeded` Run parks its item `run-budget-exceeded` with its Run ID,
   that a retry carries forward from that Run, that an executor error still
   parks `delivery-error`, and that a Run existing before the invocation never
   decides the outcome.
6. MUST verify by executing the carry-forward tests against the built tree that
   `roundfix reconcile <run-id> --carry-forward` carries a proved Task through
   refusing `pre-commit` and `commit-msg` hooks with no repository hook
   running during staging and the checkout's `core.hooksPath` unchanged.
7. MUST verify with the installed Git, in a disposable repository, that
   `cherry-pick` runs `prepare-commit-msg` and `post-commit`, that `commit
   --amend` runs `pre-commit`, `prepare-commit-msg`, `commit-msg` and
   `post-commit`, and that `-c core.hooksPath=<an empty directory>` runs none
   of them. The row records Git's version as evidence this Spec did not author.
8. MUST verify through the built binary, in a disposable Spec Root, that
   `roundfix spec check` reports `SC-VERIFY-WRAP-FRAGILE` naming the phrase,
   the file and the wrap-tolerant fix for a line-bound multi-word grep against
   a Markdown file, and reports nothing for the wrap-tolerant form, an anchored
   pattern, a heading, a single word or a `.json` operand. The same row MUST
   verify that the Verifications of Spec 0173's `task_01` and `task_02` as
   authored in commit `f0faa780`, before their rewrite to the wrap-tolerant
   form, are reported when copied into a pending Task of that Spec Root. It
   records them as evidence this Spec did not author.
9. MUST verify by executing the suggested presence form through `sh` that it
   passes on a Markdown file whose phrase is wrapped and exits `1` with stderr
   naming the phrase and the file when the phrase is absent.
10. MUST verify, read only, that the live Run Database records Runs
    `run_20260928T111930Z_13d79b0227e06da5`,
    `run_20260928T125114Z_30db2f43ddd8a501` and
    `run_20260928T154312Z_63679c6540b34603` as `BudgetExceeded`, recording
    them as evidence this Spec did not author, or the row as blocked with its
    reason when the database or a Run is absent.
11. MUST verify that `docs/user-guide/commands.md`,
    `docs/user-guide/configuration.md`, the Roundfix and write-tasks skills and
    their mirrors, the Task template, `CONTEXT.md` and ADR-0164 describe the
    renewing budget, the `run-budget-exceeded` blocker, hook-free staging,
    instruction paths that never collide and the wrap-tolerant phrase form.
12. MUST verify that `.roundfixrc.yml` is unchanged from the Spec's base.
13. MUST verify that this Spec's own artifacts satisfy the promise rule.
14. MUST NOT accept a row whose only evidence is that a file was read.

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

- `newest="$(find 'docs/specs/0177-runs-that-fit-their-budget-and-park-honestly/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-5; Core Features 1-5; Success Metrics 1-5; Acceptance
evidence; `_techspec.md` → Testing Approach 1-6; API Contracts 1-7; ADR-0014;
ADR-0020; ADR-0025; ADR-0053; ADR-0056; ADR-0057; ADR-0080; ADR-0091;
ADR-0093; ADR-0096; ADR-0097; ADR-0104; ADR-0113; ADR-0117; ADR-0125;
ADR-0130; ADR-0133; ADR-0137; ADR-0148; ADR-0155; ADR-0156; ADR-0158.
