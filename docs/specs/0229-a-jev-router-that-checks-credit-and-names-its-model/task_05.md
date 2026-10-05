---
task: task_05
spec: 0229-a-jev-router-that-checks-credit-and-names-its-model
status: failed
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes the Jev Router gate read the
OpenRouter account credit and refuse below the User Config floor
`jev.router_min_credit_usd` (US$15 when unset) with `openrouter_credit_low`,
names an OpenRouter HTTP 402 on a routed request `openrouter_credit_refused`,
and carries routed requests through a loopback relay so each `router-prompt`
Judge Log line records the routed model, provider and last response id. Every
behavior row is exercised by executing the named tests against the built tree
with local stand-ins for OpenRouter, a disposable Roundfix Home and
disposable repositories. No command sends a prompt to OpenRouter or reaches
TypeSafe.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify, by executing task_02's tests against the built tree, API
   Contracts 1 to 4 and Success Metrics 1 and 4, and record the refusal text
   the tests assert for an account at US$5.11 with US$40.31 left on the key,
   and that a refusal before Agent work falls back with a receipt naming
   `openrouter_credit_low` while one after it fails the Work Item.
3. MUST verify, by executing task_03's tests, API Contract 6 and Success
   Metric 3: the relay forwards streamed and whole answers byte-identical,
   refuses a path without an open token before any upstream request, and a
   routed prompt reporting `a/one` then `b/two` writes one `router-prompt`
   line whose `model` is `a/one, b/two`; and MUST search the test output, the
   disposable Judge Log and every file the tests wrote for the sentinel key,
   recording that it appears nowhere.
4. MUST verify, by executing task_04's tests, API Contract 5 and Success
   Metric 2: a 402 before Agent output is a failed selection the Fallback
   Chain takes with a receipt naming `openrouter_credit_refused`, and after
   Agent output the Task fails with that reason and never with
   `agent/protocol error`.
5. MUST verify from the diff that the relay listens on `127.0.0.1` only,
   neither logs nor writes a header or a body, is closed when its last session
   token is released, and that its serving goroutine is awaited by `Close`;
   that both the Implement Run and resolve engines pass the configured floor;
   that every other `daemon.NewEngine` call either passes it or starts no Agent
   prompt; and that `MonthSpend`, `Spend.CheckKeyLimit` and the Judge Log
   schema are unchanged.
6. MUST verify that the configuration guide and the Roundfix Skill's
   `runtime` reference describe the behavior task_01 names, that each mirror
   equals its canonical file, that the version was raised and recorded, that
   ADR-0218 carries the note naming ADR-0234, that the `docs` and `chore`
   default profiles and every Recommended Profile are unchanged, and that the
   `### QA settlement` section of every skill is unchanged.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the 2026-10-04 measurement addendum and the
   run records it cites; OpenRouter's limits page, its credits reference and
   its router guide. For each source it reaches it MUST record what it read
   and whether it still supports the design; for each it cannot reach it MUST
   record the row as blocked with the reason. It MAY read
   `GET https://openrouter.ai/api/v1/key` and `GET https://openrouter.ai/api/v1/credits`
   once each with `ROUNDFIX_OPENROUTER_API_KEY`, recording only the HTTP status
   and the field names, never a value or the key; it MUST NOT send any other
   request to OpenRouter.
8. MUST verify that this Spec's own artifacts satisfy the promise rule, and
   perform the glossary check of `docs/agents/domain.md`: record whether
   **Jev Router**, **Judge Log**, **Fallback Chain** and **User Config** still
   describe the behavior, and whether "credit floor" or "relay" needs a
   glossary term.
9. MUST verify from the repository history that each Task's changed files stay
   within its declarations or its `## Recorded paths`, that every Governed
   Path changed is bounded in `_authorization.md`, and that `Makefile`,
   `.roundfixrc.yml`, `go.mod` and the CI workflows did not change. This row
   reads Task commits, so it declares a `commit_range` input.
10. MUST record the non-waivable Pull Request row as
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
11. MUST NOT accept a row whose only evidence is that a file was read.
12. MUST NOT start a Delivery Queue, merge, re-run or otherwise change
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

- `newest="$(find 'docs/specs/0229-a-jev-router-that-checks-credit-and-names-its-model/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals; User Stories 1-5; Core Features 1-6; Success Metrics 1-4;
Acceptance evidence; `_techspec.md` → Testing Approach 1-6; API Contract 1;
API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract
6; Integration Points; ADR-0080; ADR-0088; ADR-0091; ADR-0104; ADR-0155;
ADR-0167; ADR-0234.
