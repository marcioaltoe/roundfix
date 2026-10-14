---
task: task_05
spec: 0270-acceptance-rules-that-match-the-authorization-and-the-queue
status: completed
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes an outside-evidence row that
the Spec's own authorization keeps offline count toward a qualifying partial,
reports a QA requirement that pins a byte invariant to a commit, and states
both rules in the Baseline guides, the owned skills, the user guide and the
glossary. Every behavior row runs against the built binary, temporary
repositories, temporary homes and the repository's own history. No row
reaches a provider or the network, starts an Agent Session, reads a
credential, or reads or writes the real `~/.roundfix`.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact,
   with `bin/roundfix` rebuilt from the audited head first.
2. MUST execute every test that task_01, task_02, task_03 and task_04 created
   and each Task's phrase checks, and record each result.
3. MUST measure Success Metric 1. It exports the folder
   `docs/specs/0258-authoring-that-proves-what-it-claims` from commit
   `078869bf`, the `source_revision` of that Spec's Archive Record, into a
   temporary directory. It runs `qa-report accept` on its
   `qa/qa-report-2026-10-13.md` with a binary built from the delivery base and
   with `bin/roundfix`, and records both exit codes and the refusal text: 1
   with the environment-row refusal, then 0.
4. MUST measure Success Metric 2 through `bin/roundfix` in a temporary
   repository, for the five network-prohibited shapes of `_techspec.md` →
   Testing Approach. It records the exit code and refusal text of
   `qa-report accept`, `archive` and `settle` for each shape, and the results
   of the Daemon shape test and the park test for the same shapes.
5. MUST measure Success Metric 3. In a temporary clone of the audited head it
   places the Spec 0264 folder as it stood at the parent of commit `4af8655e`,
   whose QA Task pinned a byte invariant to a commit, and runs
   `bin/roundfix spec check` on it with `--strict`. It records exactly one
   `SC-QA-BYTES-PINNED` finding and exit 1. It then replaces the folder with
   the one at `4af8655e` and records no finding of that code.
   When either commit is missing locally, it fetches
   `refs/pull/514/head` (0258) or `refs/pull/524/head` (0264) from `origin`
   into a temporary ref first and records that it did.
6. MUST sabotage each detector in a temporary copy of the audited head: make
   `PinnedByteInvariants` return nothing, and drop the network-prohibited
   count from `EnvironmentRowsNeedingOverride`. It records that the task_02
   and task_01 tests fail on the sabotaged copy and pass on the audited head.
7. MUST measure Success Metric 4: `bin/roundfix spec check --strict` over
   every active Spec at the audited head reports zero `SC-QA-BYTES-PINNED`
   findings, this Spec included, and `TestCheckCorpusGolden` and
   `TestCheckActiveCorpusHasNoErrors` pass.
8. MUST measure Success Metric 5. It records each text of `_techspec.md` →
   Exact texts in the files named there, and that in a temporary clone of the
   audited head a second `make baseline-digests` reports `"changed":false` and
   a Managed Refresh with a temporary `HOME` reports `File changes: 0`.
9. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the archived 0258 report and the 0264 QA
   Task from the repository history (requirements 3 and 5), and the Git and
   GitHub documentation through the passages `_prd.md` quotes. For each it
   records whether the source still supports the design, and it opens no
   network connection to do so.
10. MUST perform the glossary check of `docs/agents/domain.md`. It records that
    `CONTEXT.md` defines **Network-Prohibited Row** and **Pinned Byte
    Invariant**, that the **Delivery Base** entry carries the new sentence,
    and whether any other term the Spec introduced needs an entry.
11. MUST repeat the wording sweep. It searches `docs/`, `skills/`,
    `.agents/skills/`, `CONTEXT.md` and `AGENTS.md`, outside `docs/history/`
    and this Spec's folder, for every list of exempt outside-evidence rows
    that names only the network-denied row, and records every remaining copy
    with whether it is allowed.
12. MUST verify from the repository history:
    - that each Task's changed files stay within its declarations or its
      `## Recorded paths`;
    - that every Governed Path changed is bounded in `_authorization.md`;
    - that `go.mod`, `go.sum`, `.roundfixrc.yml`, the Makefile, the CI
      workflows, every Baseline module other than `spec-workflow` and
      `autonomous-work`, and every production Go file outside
      `internal/spec/qa.go` and `internal/speccheck` did not change;
    - that the `### QA settlement` section of `qa-gate`, `archive-spec` and
      `roundfix` is byte-identical with its bytes at the delivery base, the
      merge base of the audited head with `main`, which later Specs may have
      changed;
    - that nothing under `docs/history/` changed.

    This row reads Task commits, so it declares a `commit_range` input.
13. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval: the maintainer's recorded delivery authority;
    - checks and status: the Daemon's repository Verification at the audited
      head;
    - unresolved threads: none, because no Pull Request exists;
    - Merge-Ready acceptance: none yet, because the pre-PR review follows
      archive;
    - review-artifact ancestry: the audited head named as the claimed
      candidate.
14. MUST NOT accept a row whose only evidence is that a file was read. It MUST
    NOT send any request to a provider, record an environment value, or write
    outside temporary directories and the Spec's `qa/` folder.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Execute the Task tests and phrase checks, and the two sabotages.
- [ ] Replay the 0258 report and the 0264 QA Task, and measure Success Metrics 1 to 5.
- [ ] Record the outside evidence, the glossary check, the wording sweep and the scope audit.
- [ ] Write the QA report with the Pull Request row blocked.

## Acceptance Criteria

- [ ] The archived 0258 report is refused before the change and accepted after it.
- [ ] Every network-prohibited shape gets the same answer from every command and the Daemon.
- [ ] The pre-fix 0264 QA Task yields one `SC-QA-BYTES-PINNED` gap and the fixed one none.
- [ ] The active corpus reports zero findings of the new code.
- [ ] The scope audit finds no undeclared Governed Path and an unchanged QA settlement section.

## Consumers

None. The gate changes no Shared Contract; it writes only its own QA report.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0270-acceptance-rules-that-match-the-authorization-and-the-queue/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`


## References

- [_prd.md](_prd.md) — Goals; User Stories 1-4; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Adopter notes
- [_techspec.md](_techspec.md) — Current behavior; API Contract 1; API Contract 2; API Contract 3; API Contract 4; Exact texts; Invariants 1-7; Testing Approach; Risks & Considerations; Build Order 5
- ADR-0276; ADR-0240; ADR-0104; ADR-0167; ADR-0093; ADR-0094; ADR-0265; ADR-0148; ADR-0259; ADR-0271; ADR-0244; ADR-0256; ADR-0184
