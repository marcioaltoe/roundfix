---
task: task_05
spec: 0203-a-reviewer-that-validates-its-findings-and-remembers-its-rounds
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. The Spec changes `roundfix review` in
`internal/cli`, one field of the acpx runner in `internal/agent`, one mapping
of the Delivery Queue adapter, and the review reference, command guide and
version of the Roundfix skill.

Every behavior row is exercised by running the named tests against the built
tree. The tests use temporary repositories, temporary Artifact Directories,
fake runners and fake ACP streams. The live machine is read only for the
outside-evidence rows below, and nothing on it is written. No command reaches
a reviewer, an adapter, a provider or the network.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify, by executing task_01's tests against the built tree:
   - the prompt carries the grammar sentence and the four Delivery
     Conventions, with today's grammar sentence unchanged;
   - the anchor parser and diff index accept and refuse at their edges;
   - Surface Transcript 2 (an unanchored finding beside an anchored one) and
     Surface Transcript 3 (findings with no anchor) are reproduced;
   - reuse and `dispose` follow the validation status.
3. MUST verify, by executing task_02's tests against the built tree:
   - Surface Transcript 1 (a `C2` dismissal) is reproduced;
   - a dismissal outside a convention region is refused;
   - `no-failure` is never asked for a finding with a failure clause;
   - every validator failure cause leaves the asked findings standing;
   - each region admits and refuses at its edges;
   - Surface Transcript 7 (the `dispose` refusal) is reproduced.
4. MUST verify, by executing task_03's tests against the built tree:
   - every row of the lineage table;
   - that the round-2 prompt holds the delta and the round-1 findings and not
     the round-1 diff;
   - that round-2 anchors read the full diff;
   - Surface Transcript 5 (the blocked ceiling, with no runner call and the
     round-2 record unchanged) and Surface Transcript 6 (`ceiling-closed`);
   - the Delivery Queue mapping.
5. MUST verify, by executing task_04's tests against the built tree, and with
   every Surface Transcript now complete, including its `lineage` field:
   - the ACP session id capture;
   - that round 1 with findings leaves its session open;
   - Surface Transcript 4 (round 2 on a continued session);
   - the fallback to a fresh session;
   - that a lineage change or a dismissed reuse ends the open session.
6. MUST record, as evidence this Spec did not author, a read-only replay of
   `~/.roundfix/artifacts/339f8dac2b687a04/pre-pr-review-dispositions.jsonl`.
   Other sessions wrote that ledger on 2026-09-29 and 2026-09-30. For each
   disposition, the replay MUST:
   - resolve the named head in this repository;
   - take its merge base with `origin/main`;
   - check the finding's leading anchor against the new-side hunks of
     `git diff --no-ext-diff --no-textconv --no-color <merge base> <head> --`;
   - apply the implemented region rules to the anchor at that head.

   It MUST record how many of the 34 findings are anchored, and which fall in a
   convention region. It MUST record whether the result matches `_prd.md` →
   Acceptance evidence: the five convention dismissals in `C1`, `C2` or `C4`,
   and none of the 23 `fixed` findings. When the ledger or a head is absent,
   it MUST record the row as blocked with its reason.
7. MUST record, as published evidence, the Agent Client Protocol session
   setup page (<https://agentclientprotocol.com/protocol/v1/session-setup>)
   on `session/resume` and `sessionCapabilities.resume`, and GitHub's review
   comment reference (<https://docs.github.com/en/rest/pulls/comments>) on
   comments made on the diff. When a page is unreachable, it MUST record the
   row as blocked with its reason.
8. MUST read, without running them, the installed codex-acp and
   claude-agent-acp adapters and acpx under the global Node modules
   directory. It MUST record each version, whether each adapter advertises
   `loadSession` and `sessionCapabilities.resume`, and whether acpx tries
   resume before load before a new session. When one is absent, it MUST
   record that row as blocked with its reason.
9. MUST verify that the review reference of the Roundfix skill, its mirror and
   `docs/user-guide/commands/review.md` describe the conventions, validation,
   the lineage, the continued session and the ceiling. The skill's version
   MUST be recorded in `skills/testdata/owned-skill-versions.json`,
   `make skills-sync-check` MUST exit `0`, and the `### QA settlement`
   section of every skill MUST be unchanged.
10. MUST verify that this Spec's own artifacts satisfy the promise rule, and
    that Specs 0191 and 0194, which this Spec lists as prerequisites, are
    merged into the audited head's ancestry.
11. MUST verify from the repository history that each Task's changed files
    stay within its declarations or its `## Recorded paths`, that every
    Governed Path is bounded in `_authorization.md`, and that
    `internal/cli/cli_test.go` and `.roundfixrc.yml` did not change.
12. MUST record the non-waivable Pull Request row as `blocked (environment: no
    open Pull Request)`, naming the Pull Request row in its provenance, and
    support it with the pre-PR equivalent of each control:
    - approval: the maintainer's recorded delivery authority;
    - checks and status: the Daemon's repository Verification at the audited
      head;
    - unresolved threads: none, because no Pull Request exists;
    - Merge-Ready acceptance: none yet, because the pre-PR review follows
      archive;
    - review-artifact ancestry: the audited head named as the claimed
      candidate.
13. MUST NOT accept a row whose only evidence is that a file was read.
14. MUST NOT run `roundfix review` or `roundfix deliver` against a real
    reviewer, MUST NOT start a session on a real ACP adapter, and MUST NOT
    write to `~/.roundfix` or `~/.acpx`.

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

- `newest="$(find 'docs/specs/0203-a-reviewer-that-validates-its-findings-and-remembers-its-rounds/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-5; User Stories 1-5; Core Features 1-9; Success Metrics
1-6; Success Metric 7; Declared breaks; Acceptance evidence; `_techspec.md` → Testing Approach
1-7; API Contracts 1-5; Surface Transcripts 1-7; ADR-0017; ADR-0018;
ADR-0051; ADR-0080; ADR-0091; ADR-0093; ADR-0094; ADR-0096; ADR-0104;
ADR-0117; ADR-0151; ADR-0153; ADR-0155; ADR-0156; ADR-0165; ADR-0166;
ADR-0167; ADR-0168; ADR-0169; ADR-0174; ADR-0176; ADR-0183; ADR-0184;
ADR-0189; ADR-0196; ADR-0197.
