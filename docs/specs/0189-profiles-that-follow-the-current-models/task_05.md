---
task: task_05
spec: 0189-profiles-that-follow-the-current-models
status: pending
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec changes static data and its
readers in `internal/agent`, `internal/config`, `internal/cli` and
`internal/baselineacp`, and the documents that state them. Every behavior row
is exercised through the built binary or by executing the named tests against
the built tree. Each binary command runs in a disposable repository with a
disposable Roundfix Home. No command reaches a real ACP adapter, a provider or
the network.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify through the built binary, in a disposable repository and Roundfix Home with no configured profile, that `roundfix profiles show --json` reports schema `roundfix/profiles/v2`, that each required category's effective profile equals its recommendation rows with source `built-in`, and that no tuple names `gpt-5.5`.
3. MUST verify, by executing task_01's and task_03's catalog tests against the built tree, that each catalog is exactly the Reference data table, and MUST verify the two effort lists and that Doctor refuses an adapter below each floor and accepts it at the floor.
4. MUST verify, by executing task_02's tests against the built tree, the ten Recommended Profiles, the runtime rule for `review` and the other categories, the catalog invariant, and the `profiles show` text and JSON.
5. MUST verify, by executing task_03's tests against the built tree, the built-in profiles, the generated config, the legacy Codex default, the review provider check and the Baseline analysis models, and MUST confirm that no non-test Go file under `internal` or `cmd` contains `gpt-5.5`.
6. MUST verify, by executing task_04's test and the documentation contract against the built tree, that the reference states the shipped snapshot and that `TestProfilesDocumentationContractMatchesPublicGuidance` passes.
7. MUST record, as evidence this Spec did not author, the adapter advertisement and the retirement dates quoted in `_prd.md` → Acceptance evidence, and compare the built catalog with that advertisement: every catalog model is advertised, and the only advertised Codex model the catalog omits is `gpt-5.5`. When the published page is reachable, it MUST record what it read there; when it is not, it MUST record the row as blocked with its reason.
8. MUST verify that the two guides, the Roundfix skill, its mirror and both manifests carry the floors `2.0.1` and `0.84.0`, `roundfix/profiles/v2` and the Recommended Profile, that the `### QA settlement` section of the skill is unchanged, and that `make skills-sync-check` exits `0`.
9. MUST verify that this Spec's own artifacts satisfy the promise rule, and check that `CONTEXT.md` carries **Recommended Profile** and lists the retired term under its `_Avoid_` line.
10. MUST verify from the repository history that each Task's changed files stay within its declarations or its `## Recorded paths`, that every Governed Path is bounded in `_authorization.md`, and that `.roundfixrc.yml` changed in comment lines only.
11. MUST record the non-waivable Pull Request row as `blocked (environment: no open Pull Request)`, naming the Pull Request row in its provenance, and support it with the pre-PR equivalent of each control:
    - approval: the maintainer's recorded delivery authority;
    - checks and status: the Daemon's repository Verification at the audited head;
    - unresolved threads: none, because no Pull Request exists;
    - Merge-Ready acceptance: none yet, because the pre-PR review follows archive;
    - review-artifact ancestry: the audited head named as the claimed candidate.
12. MUST NOT accept a row whose only evidence is that a file was read.
13. MUST NOT write to the live Run Database under `~/.roundfix`, and MUST NOT start a session on a real ACP adapter.

## Subtasks

- [ ] Build the matrix from the Requirements above and the sources no declaration waives.
- [ ] Execute each row against the built tree and record its evidence.
- [ ] Write the dated QA Report with its verdict.

## Acceptance Criteria

- [ ] The QA Report records a verdict and names, in each row's provenance, the sources that row covers.
- [ ] The repository Verification result is recorded as a fact, not re-derived.

## Context

- instruction: `.agents/skills/qa-gate/SKILL.md`

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0189-profiles-that-follow-the-current-models/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-4; Core Features 1-6; Success Metrics 1-5; Acceptance
evidence; `_techspec.md` → Testing Approach 1-7; API Contracts 1-5; ADR-0037;
ADR-0049; ADR-0050; ADR-0069; ADR-0079; ADR-0080; ADR-0088; ADR-0091;
ADR-0093; ADR-0094; ADR-0096; ADR-0104; ADR-0107; ADR-0117; ADR-0140;
ADR-0147; ADR-0151; ADR-0153; ADR-0155; ADR-0156; ADR-0166; ADR-0167;
ADR-0176; ADR-0180.

## Result
