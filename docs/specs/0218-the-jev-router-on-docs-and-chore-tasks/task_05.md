---
task: task_05
spec: 0218-the-jev-router-on-docs-and-chore-tasks
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes the Jev Router an opt-in
OpenCode selection in `internal/agent` and `internal/config`, adds the
spend reader and router line in `internal/jevrouter`, gates routed prompts in
`internal/daemon`, documents it, and measures it live. Behavior rows run the
built binary and the named tests against fakes and an `httptest` key
endpoint; no row reaches OpenRouter or the network, and the only live evidence
is the measurement record Task 04 committed.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify through the built binary, in a temporary repository and Home:
   a User Config naming the router reproduces Surface Transcript 1; a Project
   Config naming it with an empty effort loads in `roundfix profiles show`;
   one with a non-empty effort is refused; and a one-Run override naming it is
   refused with exit `2`.
3. MUST verify, by executing the Task 01 tests with a sentinel key, that the
   sentinel reaches no acpx argument, configuration value, Run Event or
   Judge Log line, that a routed session's environment holds the fixed inline
   provider of the TechSpec, and that no production file names
   `"OPENROUTER_API_KEY"`.
4. MUST verify, by executing the Task 02 and Task 03 tests, the spend
   formula, the fallback before work, the failure after work, the unreadable
   spend refusal, one `router-prompt` line per routed prompt read back by
   the judge package, and no gate call for a non-routed prompt.
5. MUST read the measurement record Task 04 committed and confirm one row per
   replay that ran, a routed cost per routed row, a reading, no key in the
   record, and that the scratch clones were of this repository. When Task 04
   ended with `jev_router_key_missing`, the row MUST be recorded as
   `blocked (environment: ROUNDFIX_OPENROUTER_API_KEY not set)`.
6. MUST record, as evidence this Spec did not author:
   - OpenRouter's model listing and page for `typesafe/jev-router`
     (<https://openrouter.ai/typesafe/jev-router>), for its parameters and
     variable price;
   - OpenRouter's key reference
     (<https://openrouter.ai/docs/api/api-reference/api-keys/get-current-api-key>)
     and limits page (<https://openrouter.ai/docs/api_reference/limits>), for
     `usage_monthly`;
   - OpenCode's configuration and providers references
     (<https://opencode.ai/docs/config/>, <https://opencode.ai/docs/providers/>),
     for `OPENCODE_CONFIG_CONTENT`, `{env:...}` and the compatible provider;
   - the live measurement record, as a service this Spec did not build.

   When a source is unavailable, it MUST record that row as blocked with its
   reason.
7. MUST verify that this Spec's own artifacts satisfy the promise rule, and
   check whether it introduced, changed or retired a glossary term; record
   whether `CONTEXT.md` needs `Jev Router`, routed selection or
   `router-prompt`.
8. MUST verify from the repository history that each Task's changed files
   stay within its declarations or its `## Recorded paths`, that the only
   Governed Paths changed are the three `_authorization.md` bounds, and that
   the Recommended Profile, the built-in profiles and `.roundfixrc.yml` are
   unchanged. This row reads Task commits, so it declares a `commit_range`
   input.
9. MUST record the non-waivable Pull Request row as
   `blocked (environment: no open Pull Request)`, naming the Pull Request row
   in its provenance. It MUST support that row with the pre-PR equivalent of
   each control:
   - approval: the maintainer's recorded delivery authority;
   - checks and status: the Daemon's repository Verification at the audited
     head;
   - unresolved threads: none, because no Pull Request exists;
   - Merge-Ready acceptance: none yet, because the pre-PR review follows
     archive;
   - review-artifact ancestry: the audited head named as the claimed
     candidate.
10. MUST NOT accept a row whose only evidence is that a file was read.
11. MUST NOT write to the live Run Database or Judge Log under
    `~/.roundfix`, read a real `ROUNDFIX_OPENROUTER_API_KEY`, or send a
    prompt or a request anywhere.

## Subtasks

- [ ] Build the matrix from the Requirements above and the sources no
      declaration waives.
- [ ] Execute each row against the built tree and record its evidence.
- [ ] Write the dated QA Report with its verdict.

## Acceptance Criteria

- [ ] The QA Report records a verdict and names, in each row's provenance, the
      sources that row covers.
- [ ] The repository Verification result is recorded as a fact, not re-derived.
- [ ] The measurement row is either satisfied by the committed record or
      blocked by the named key blocker.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0218-the-jev-router-on-docs-and-chore-tasks/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-4; User Stories 1-4; Core Features 1-7; Success Metrics
1-5; Acceptance evidence; `_techspec.md` → Surface Transcript 1; API
Contracts 1-4; Testing Approach 1-5; The measurement protocol; Build Order 5;
ADR-0080; ADR-0088; ADR-0091; ADR-0104; ADR-0155; ADR-0156; ADR-0167;
ADR-0218.
