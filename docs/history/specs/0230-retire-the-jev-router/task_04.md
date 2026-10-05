---
task: task_04
spec: 0230-retire-the-jev-router
status: failed
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec retires the Jev Router (its
selection, inline provider, relay, spend and credit gate, credit floor and
`router-prompt` lines) and refuses any `opencode` selection that reaches an
OpenAI or Anthropic model through OpenRouter, in configuration and again in
the runner, while the advisory judge stays unchanged. Every behavior row is
exercised against the built tree with disposable Homes, disposable
repositories and the fake ACP adapter. No command reaches OpenRouter or
TypeSafe.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify, by executing task_02's tests and by running the built
   `roundfix` against a disposable repository and Home, Success Metric 1 and
   API Contracts 1 to 3: each of the four refused models fails in Project
   Config, User Config, a one-Run override and `roundfix profiles configure`
   with the field, the model and the rule, the configure command writes
   nothing, and the allowed selections load; and MUST reproduce Surface
   Transcript 1 and Surface Transcript 2 through the built product.
3. MUST verify, by executing task_01's tests, Success Metric 2 and API
   Contract 4: a refused selection never starts the fake ACP adapter through
   `RunPrompt` or `PrepareSession`, and before Agent work the fallback takes
   the Task with `reason_code` `subscription_only`.
4. MUST verify Success Metric 3 and API Contract 6 from the tree: outside the
   judge no Go source names OpenRouter's host, sets
   `OPENCODE_CONFIG_CONTENT` or writes `router-prompt`; the `internal/jevrouter`
   package is gone; the judge package and the `spec judge` command source are
   byte-identical to the Spec's base; and the judge's existing tests pass
   without reaching the network.
5. MUST verify Success Metric 4 and API Contract 5 by executing task_02's
   deprecated-key test and Surface Transcript 2, and that the existing
   deprecated-key warnings are unchanged.
6. MUST verify Success Metric 5: the configuration guide and the Roundfix
   Skill's `runtime` reference state the rule and no longer describe the
   router's gate, floor, relay or refusal codes; each mirror equals its
   canonical file and the raised version is recorded; the model reference and
   the `.roundfixrc.yml` comment state the rule and no configuration value
   changed; ADR-0218 and ADR-0234 are under `docs/history/adr/` as superseded
   by ADR-0235 and ADR-0231 carries its note; the backlog entry on routed
   spend is declined under `docs/history/backlog/`; and the
   `### QA settlement` section of every skill is unchanged.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the OpenRouter activity export (its
   SHA-256 and per-model sums, read only if the operator makes the file
   available, otherwise carried from the PRD), the "Relay confirmation,
   2026-10-05" and "Maintainer decision, 2026-10-05" sections of Spec 0218's
   archived measurement record, and OpenRouter's Auto Router and model
   variants guides. For each source it reaches it MUST record what it read
   and whether it still supports the design; for each it cannot reach it MUST
   record the row as blocked with the reason. It MUST NOT send any request to
   OpenRouter or TypeSafe.
8. MUST verify that this Spec's own artifacts satisfy the promise rule, and
   perform the glossary check of `docs/agents/domain.md`: record whether
   **Judge Log**, **Agent Selection**, **Fallback Chain** and **User Config**
   still describe the behavior, and whether "subscription rule" needs a
   glossary term.
9. MUST verify from the repository history that each Task's changed files stay
   within its declarations or its `## Recorded paths`, that every Governed
   Path changed is bounded in `_authorization.md`, and that `Makefile`,
   `go.mod` and the CI workflows did not change. This row reads Task commits,
   so it declares a `commit_range` input.
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

- `newest="$(find 'docs/specs/0230-retire-the-jev-router/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals; User Stories 1-5; Core Features 1-6; Success Metrics 1-5;
Acceptance evidence; `_techspec.md` → Interfaces; API Contract 1; API
Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6;
Surface Transcript 1; Surface Transcript 2; Testing Approach; Build Order 4;
ADR-0235; ADR-0080; ADR-0091.
