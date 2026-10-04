---
task: task_03
spec: 0225-a-jev-ceiling-the-maintainer-sets
status: pending
type: qa
complexity: medium
---

# Task 03: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes the monthly Jev ceiling the
User Config value `jev.monthly_ceiling_usd`, with US$5 when unset, read by
`roundfix spec judge` and by the Jev Router gate and its key-limit check, and
describes it in the guides and the Roundfix Skill. Every behavior row is
exercised through the built binary or by executing the named tests against
the built tree, with disposable repositories and a disposable Roundfix Home.
No command reaches OpenRouter, TypeSafe or the network, except the read of
the published page the outside-evidence row names.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST reproduce Surface Transcripts 1, 2, 3 and 4 through the built binary
   with `HOME` set to a disposable directory, `NODE_OPTIONS` and
   `ROUNDFIX_OPENROUTER_API_KEY` unset, and `ROUNDFIX_TYPESAFE_API_KEY` set to
   a dummy value, in a disposable repository holding an active Spec
   `0300-example`, with the Judge Log month written so the ceiling is already
   reached. It MUST record stdout, stderr and the exit code of each, and that
   no request left the machine.
3. MUST verify, by executing task_02's tests against the built tree, API
   Contracts 1 to 4 and Success Metrics 1 to 3, and record the messages the
   tests assert, including the `jev_router_key_unbounded` message at a
   configured ceiling of 50 and at the built-in one.
4. MUST verify from the diff that the Implement Run engine and the resolve
   engine pass the loaded ceiling to the Daemon, that every other
   `daemon.NewEngine` call either passes it or starts no Agent prompt, and
   that `questions.json`, `judge.Run`, `jevrouter.MonthSpend` and
   `Spend.CheckKeyLimit` are unchanged.
5. MUST verify that the configuration and `spec` guides and the Roundfix
   Skill's `spec` and `runtime` references describe the behavior task_01
   names and no longer state US$5 as the fixed ceiling, that each mirror
   equals its canonical file, that the version was raised and recorded, that
   ADR-0201 and ADR-0218 carry the supersession note naming ADR-0231, and
   that the `### QA settlement` section of every skill is unchanged.
6. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the maintainer's decision of 2026-10-04
   and the key status read that day; OpenRouter's limits reference; and a
   reproduction with the binary the operator has installed, run with a
   disposable User Config holding `jev.monthly_ceiling_usd: 50`, showing the
   refusal the PRD quotes. For each source it reaches it MUST record what it
   read and whether it still supports the design; for each it cannot reach it
   MUST record the row as blocked with the reason. It MUST NOT query the
   OpenRouter key endpoint or any Jev endpoint.
7. MUST verify that this Spec's own artifacts satisfy the promise rule, and
   perform the glossary check of `docs/agents/domain.md`: record whether
   **User Config**, **Project Config**, **Roundfix Home** and **Judge Log**
   still describe the behavior, and whether "Jev ceiling" needs a glossary
   term.
8. MUST verify from the repository history that each Task's changed files stay
   within its declarations or its `## Recorded paths`, that every Governed
   Path changed is bounded in `_authorization.md`, and that `Makefile`,
   `.roundfixrc.yml`, `go.mod` and the CI workflows did not change. This row
   reads Task commits, so it declares a `commit_range` input.
9. MUST record the non-waivable Pull Request row as
   `blocked (environment: no open Pull Request)`, naming the Pull Request row
   in its provenance, and support it with the pre-PR equivalent of each
   control:
   - approval: the maintainer's recorded delivery authority;
   - checks and status: the Daemon's repository Verification at the audited
     head;
   - unresolved threads: none, because no Pull Request exists;
   - Merge-Ready acceptance: none yet, because the pre-PR review follows
     archive;
   - review-artifact ancestry: the audited head named as the claimed
     candidate.
10. MUST NOT accept a row whose only evidence is that a file was read.
11. MUST NOT start a Delivery Queue, merge, re-run or otherwise change
    anything on GitHub, and MUST NOT read or write the real `~/.roundfix`.

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

- `newest="$(find 'docs/specs/0225-a-jev-ceiling-the-maintainer-sets/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References 
`_prd.md` → Goals; User Stories 1-5; Core Features 1-6; Success Metrics 1-3;
Acceptance evidence; `_techspec.md` → Testing Approach 1-4; API Contract 1;
API Contract 2; API Contract 3; API Contract 4; Surface Transcript 1; Surface
Transcript 2; Surface Transcript 3; Surface Transcript 4; Integration Points;
ADR-0080; ADR-0088; ADR-0091; ADR-0104; ADR-0155; ADR-0167; ADR-0231.
