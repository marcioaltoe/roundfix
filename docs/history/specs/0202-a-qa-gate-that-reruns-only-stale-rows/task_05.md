---
task: task_05
spec: 0202-a-qa-gate-that-reruns-only-stale-rows
status: completed
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec changes production code in
`internal/speccheck`, `internal/daemon` and `internal/agent`, the qa-gate
skill and the Context-Driven Development guide. Every behavior row is
exercised by executing the named tests against the built tree or through the
built binary in a disposable repository with a disposable Roundfix Home. No
command reaches GitHub, a provider or a live remote.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify, by executing task_01's tests against the built tree:
   - a passing row with only repository inputs is recorded at the audited
     head, and every row that cannot carry is not;
   - an Agent-written `evidence_snapshots` value is replaced or removed, and
     no other byte of the report changes;
   - the recorder's output is carried by the existing mechanical-stage reader;
   - the QA Report commit carries the key at the audited head.
3. MUST verify, by executing task_02's tests against the built tree:
   - a head recorded by a QA Report commit on another branch carries, and one
     recorded by a Task commit or by nothing does not;
   - always-observed rows never carry;
   - every prior row has one disposition with a reason from the closed list;
   - a carried row keeps its provenance;
   - a result without dispositions renders today's bytes, and the governed
     carry tests pass unedited.
4. MUST verify, by executing task_03's tests against the built tree, that a
   failed pass committed only on a side branch is imported byte for byte, that
   each refusal leaves nothing behind, and that the second pass carries the
   unmoved row end to end.
5. MUST verify, by executing task_06's two-pass tests, task_03's end-to-end test and the event tests of
   task_01 and task_02 against the built tree, that two consecutive gate
   passes publish the `prior_report`, `mechanical` and `evidence_snapshots`
   phases with their payloads. It MUST also verify that the second pass's
   committed QA Report carries a row the first pass recorded and whose inputs
   no Task changed, and read that report back from the test's repository as
   the independent confirmation. No public command runs the QA stage without
   a live ACP Runtime, so this row runs through the Daemon's own test seam.
6. MUST record, as evidence this Spec did not author:
   - Spec 0179's archived QA Reports, read from the repository. They show
     that every pass declared repository inputs and none recorded
     `evidence_snapshots`. Recompute, from the reports and
     `git diff --name-only` between consecutive `build` heads, how many
     re-executed passing rows had no declared input among the changed files.
     The PRD states 24 of 41;
   - Spec 0192's QA Report commits `20db691c` and `3327baa1`, read without
     writing, when this repository still holds them. Record whether their
     audited heads are ancestors and which files their trees differ in;
   - the two published sources the PRD's Acceptance evidence names.

   When a source is unavailable, it MUST record that row as blocked with its
   reason.
7. MUST verify that the qa-gate skill and its mirror teach the carry, that
   their `### QA settlement` sections are unchanged, that
   `make skills-sync-check` exits `0`, and that the guide names every emitted
   word of the TechSpec's Vocabulary Contract.
8. MUST verify that this Spec's own artifacts satisfy the promise rule, and
   check whether it introduced, changed or retired a glossary term that
   `CONTEXT.md` does not carry. **Evidence Snapshot**, **Carried Row** and
   **Carry Disposition** are the terms it coined.
9. MUST verify from the repository history that each Task's changed files stay
   within its declarations or its `## Recorded paths`, and that every Governed
   Path is bounded in `_authorization.md`. This row reads Task commits, so it
   declares a `commit_range` input.
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
12. MUST NOT write to the live Run Database under `~/.roundfix`.

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

- `newest="$(find 'docs/specs/0202-a-qa-gate-that-reruns-only-stale-rows/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References
`_prd.md` → Goals 1-4; User Stories 1-4; Core Features 1-6; Success Metrics
1-5; Acceptance evidence; `_techspec.md` → Testing Approach 1-5; API Contracts
1-5; Vocabulary Contract; ADR-0053; ADR-0057; ADR-0080; ADR-0088; ADR-0091;
ADR-0096; ADR-0097; ADR-0104; ADR-0155; ADR-0167; ADR-0170; ADR-0194;
ADR-0195.
