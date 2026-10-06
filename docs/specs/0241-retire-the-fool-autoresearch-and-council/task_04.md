---
task: task_04
spec: 0241-retire-the-fool-autoresearch-and-council
status: pending
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec retires `council` and `the-fool`
from the Baseline (modules, triggers, Setup Snapshots), removes `council` from
the Roundfix-owned bundle, drops `the-fool` and `autoresearch` from this
repository, and makes `roundfix baseline update` list retired copies an
adopter still holds. Behavior rows run against the built tree in disposable
repositories with a temporary home; the asset sync row uses a local clone of
`~/dev/skills`. No row reaches the network, reads or writes the real
`~/.roundfix`, or touches an adopter's repository.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify Success Metric 1 and Success Metric 2 by executing task_01's
   catalog and sync tests, and by running the TechSpec's refresh procedure's
   `baseline assets sync --check` with the built binary against a disposable
   clone of `~/dev/skills` at `b3c45a45f1bccd3b33aaecaaa22947d942f2fc02`:
   exit `0`, and no snapshot names `council` or `the-fool`. When the local
   checkout lacks the commit, the row is blocked with that reason.
3. MUST verify Success Metric 3 and API Contract 3: `roundfix skills list`
   from the built binary names no `council`, and `roundfix skills install
   --target project` into a disposable directory writes no `council` tree.
4. MUST verify Success Metric 4: in this repository the built `roundfix
   doctor` prints `skills: ok (41 required: 13 Roundfix-owned, 28 external)`
   (record the line; a machine-dependent failure of another Doctor line does
   not fail this row), no lock entry, recommended line or `.agents/skills/`
   directory names `the-fool`, `autoresearch` or `council`, the kept skills
   `grilling`, `grill-with-docs`, `write-idea`, `business-analyst` and
   `handoff` are installed, and `baseline update --repo . --format json`
   reports `current`.
5. MUST verify Success Metric 5, API Contract 2 and Surface Transcript 1
   through the built binary on a disposable adopted repository that holds
   `.agents/skills/council` and a locked `.agents/skills/the-fool`: the text
   output's `Skills retired` lines, its state and exit match the transcript;
   `--no-skills` prints no `Skills retired`; after the Release note's removal
   commands, the next update prints no `Skills retired` and the JSON has no
   `skills.retired`.
6. MUST verify Success Metric 6: `docs/user-guide/commands/baseline.md`
   carries the removal commands of `_prd.md` → Release note, and `CONTEXT.md`
   defines **Retired Skill**; perform the glossary check of
   `docs/agents/domain.md` and record whether **Retired Skill**, **Setup
   Snapshot** terms and **Repository Skill Set** still describe the behavior.
7. MUST verify that each changed owned skill's mirror equals its canonical
   copy, that the raised versions of `write-idea`, `write-prd` and `roundfix`
   are recorded, and that no `### QA settlement` section changed.
8. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the upstream setup lists at `b3c45a4`
   (that they still name both skills) and vercel-labs/skills issue #977. For
   each it reaches it MUST record what it shows; for each it cannot reach it
   MUST record the row as blocked with the reason.
9. MUST verify that this Spec's own artifacts satisfy the promise rule.
10. MUST verify from the repository history that each Task's changed files
    stay within its declarations or its `## Recorded paths`, that every
    Governed Path changed is bounded in `_authorization.md`, that the
    `Makefile` change is only the `OWNED_SKILLS` line and its comment, that
    no vendored skill's content was edited (only deleted), and that
    `go.mod`, the CI workflows, the Source Baselines and every parity fixture
    other than `asset-sync.json` did not change. This row reads Task commits,
    so it declares a `commit_range` input.
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
- [ ] Run the asset sync check against the local upstream clone.
- [ ] Exercise skills, Doctor and update in this repository and in a disposable adopted repository.
- [ ] Reproduce Surface Transcript 1 and the removal commands.
- [ ] Record the outside evidence, glossary check and scope audit.
- [ ] Write the QA report with the Pull Request row blocked.

## Acceptance Criteria

- [ ] Every behavior row passes on recorded evidence.
- [ ] Surface Transcript 1 matches through the built binary.
- [ ] The outside-evidence row records what each source shows or why it was unreachable.
- [ ] The scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0241-retire-the-fool-autoresearch-and-council/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Release note; Acceptance evidence
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; Surface Transcript 1; The refresh procedure; Vocabulary Contract; Build Order 4
- ADR-0246; ADR-0080; ADR-0091; ADR-0104
