---
task: task_05
spec: 0233-a-light-tier-on-open-models
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec runs every non-QA
`complexity: low` Task that declares no Governed Path on an open OpenRouter
model through OpenCode by default, under a User Config allow-list and monthly
ceiling, records light spend in a Light Spend Log, skips the tier with a
warning without a key, past the ceiling or with an unreadable log, escalates a
failed light Task once to its category's profile, and lets the advisory judge
suggest a tier at the `tasks` stage. Every behavior row is exercised against
the built tree with disposable Homes, disposable repositories and the fake ACP
adapter; only the live row of Requirement 8 reaches OpenRouter.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify, by executing task_02's tests and by running the built
   `roundfix` against a disposable repository and Home, Success Metric 5 and
   API Contracts 1 to 3: the defaults, each refused light model and ceiling,
   an empty list, and the Project Config warning for both keys.
3. MUST verify, by executing task_03's tests, Success Metrics 1 to 4 and
   API Contracts 4 to 6: the dispatch order and the unchanged dispatch of
   `medium`, `qa` and Governed-Path Tasks; no `Session setup.` prompt; each
   skip reason with its warning and Run Event; a light start failure taken
   before Agent work; spend lines that add up to the last reading; the inline
   option naming the variable with neither the key value nor
   `OPENROUTER_API_KEY` in the session's environment; and one escalation in a
   new session.
4. MUST reproduce Surface Transcript 3 with the built `roundfix implement`,
   the fake ACP adapter and a disposable Home with no OpenRouter key, and
   show that an `--agent` override keeps a light Task on the override.
5. MUST verify, by executing task_04's tests and running the built
   `roundfix spec judge` against a disposable Spec with no key, Success
   Metric 6, API Contract 7 and Surface Transcripts 1 and 2, and that the
   judge's output without `--stage` is unchanged.
6. MUST verify Core Feature 8: the configuration guide, the `spec` command
   reference, the model reference, the Roundfix Skill's `runtime` and `spec`
   references and the write-tasks skill describe the tier, the keys, the
   warning, the Light Spend Log, the escalation and the `tasks` stage; each
   mirror equals its canonical file and both raised versions are recorded;
   each Vocabulary Contract pattern appears in its emitting and documenting
   file; and the `### QA settlement` section of every skill is unchanged.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the measurement section of the adopted
   Backlog Entry, the installed acpx's handling of allowed tools, OpenCode's
   configuration and providers documentation, and OpenCode issue 13219. For
   each source it reaches it MUST record what it read and whether it still
   supports the design; for each it cannot reach it MUST record the row as
   blocked with the reason.
8. MUST verify Success Metric 7 live only when `ROUNDFIX_OPENROUTER_API_KEY`
   is present in the gate's environment: run at most two real light Tasks,
   each a `low` `docs` Task in a disposable repository, through the built
   `roundfix implement` with a disposable Home whose User Config keeps the
   default light model, read the key's usage with OpenRouter's key endpoint
   before and after, stop before total spend could pass US$5, and record the
   Run, the model, the Light Spend Log lines, OpenCode's cost and the key's
   movement. It MUST use only open models and MUST NOT print or store the key.
   Without the key the row is `blocked (environment: no OpenRouter key)`.
9. MUST verify that this Spec's own artifacts satisfy the promise rule, and
   perform the glossary check of `docs/agents/domain.md`: record whether
   **Agent Selection Profile**, **Preferred Selection**, **Fallback Chain**,
   **Verification Feedback**, **Judge Log** and **User Config** still describe
   the behavior, and whether **Light Tier** and **Light Spend Log** need
   glossary terms.
10. MUST verify from the repository history that each Task's changed files
    stay within its declarations or its `## Recorded paths`, that every
    Governed Path changed is bounded in `_authorization.md`, and that
    `Makefile`, `go.mod`, `.roundfixrc.yml`, the Run Database schema and the
    CI workflows did not change. This row reads Task commits, so it declares a
    `commit_range` input.
11. MUST record the non-waivable Pull Request row as
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
12. MUST NOT accept a row whose only evidence is that a file was read.
13. MUST NOT start a Delivery Queue, merge, re-run or otherwise change
    anything on GitHub, MUST NOT read or write the real `~/.roundfix`, and
    MUST NOT send any request to OpenRouter or TypeSafe outside Requirement 8.

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

- `newest="$(find 'docs/specs/0233-a-light-tier-on-open-models/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals; User Stories 1-6; Core Features 1-8; Success Metrics 1-7;
Acceptance evidence; `_techspec.md` → Interfaces; Invariants; API Contract 1;
API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract
6; API Contract 7; Surface Transcript 1; Surface Transcript 2; Surface
Transcript 3; Vocabulary Contract; Testing Approach; Build Order 5; ADR-0238;
ADR-0080; ADR-0091.
