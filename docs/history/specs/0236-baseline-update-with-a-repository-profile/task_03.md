---
task: task_03
spec: 0236-baseline-update-with-a-repository-profile
status: completed
type: qa
complexity: high
---

# Task 03: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes the skill snapshot
comparison, skill restore and lock reconciliation resolve a repository-owned
Baseline Profile as `baseline update` does, so `roundfix baseline update`
and Doctor work in repositories adopted with their own profile. Every
behavior row runs against the built tree in disposable repositories with a
temporary home; no row reaches the network, reads or writes the real
`~/.roundfix`, or touches an adopter's repository.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify Success Metric 1 and Surface Transcript 1 through the built
   `roundfix` on a disposable repository adopted with `baseline profile init
   --id <id> --from go-cli-tui`, `baseline plan --profile <id>` and
   `baseline apply`: `baseline update --format text` matches the transcript's
   stdout, stderr and exit, and the JSON form reports `state: current`.
3. MUST verify Success Metric 2 and Success Metric 3 by executing task_02's
   preview and skills-stage tests, and confirm on the disposable repository
   that `baseline update --yes` installs the Roundfix-owned skills without
   failing on the snapshot comparison (an external restore that needs the
   network may degrade to a drifted warning; the row records it).
4. MUST verify Success Metric 4 and API Contract 4 by executing task_02's
   Doctor test, and by running the built `roundfix doctor` in the disposable
   repository: its `skills:` line carries no "Unknown built-in Baseline
   Profile".
5. MUST verify Success Metric 5, API Contract 2 and API Contract 3 through
   Surface Transcript 2 and Surface Transcript 3 with the built binary:
   stdout, stderr and exit match.
6. MUST verify Success Metric 6: every existing baseline, restore,
   reconcile, Doctor and update test passes, and the only changed
   expectations are those task_01 and task_02 declare.
7. MUST verify that `docs/user-guide/commands/baseline.md`,
   `docs/user-guide/commands/doctor.md` and the Roundfix Skill's baseline
   reference describe repository profiles and the new findings, that the
   Vocabulary Contract's pattern is documented where it says, that each
   skill mirror equals its canonical file, and that the raised owned-skill
   version is recorded.
8. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the Pantheon and Oraculum reports of
   2026-10-05 in the Secondbrain's `inbox/roundfix/_triaged/`. For each it
   reaches it MUST record the failure it describes and whether the fixed
   binary's behavior on an equivalent disposable repository answers it; for
   each it cannot reach it MUST record the row as blocked with the reason.
9. MUST verify that this Spec's own artifacts satisfy the promise rule, and
   perform the glossary check of `docs/agents/domain.md`: record whether
   **Baseline Profile** and the Setup Snapshot terms still describe the
   behavior now that a repository profile takes contracts from several
   snapshots, and whether a new term is needed.
10. MUST verify from the repository history that each Task's changed files
    stay within its declarations or its `## Recorded paths`, that every
    Governed Path changed is bounded in `_authorization.md`, and that
    `Makefile`, `go.mod`, the CI workflows and the embedded Baseline catalog
    did not change. This row reads Task commits, so it declares a
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
12. MUST NOT accept a row whose only evidence is that a file was read, and
    MUST NOT run any command in an adopter's repository.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Adopt a disposable repository with a repository profile and exercise update, Doctor and restore.
- [ ] Reproduce the three Surface Transcripts through the built binary.
- [ ] Record the outside evidence, glossary check and scope audit.
- [ ] Write the QA report with the Pull Request row blocked.

## Acceptance Criteria

- [ ] Every behavior row passes on recorded evidence.
- [ ] The three Surface Transcripts match through the built binary.
- [ ] The outside-evidence row records what each adopter report describes
      or why it was unreachable.
- [ ] The scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0236-baseline-update-with-a-repository-profile/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Vocabulary Contract; Build Order 3
- ADR-0241; ADR-0080; ADR-0091; ADR-0104
