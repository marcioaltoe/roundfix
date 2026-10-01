---
task: task_05
spec: 0194-a-skill-and-a-command-guide-read-one-command-at-a-time
status: completed
type: qa
complexity: medium
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec moves text between files and
changes how contracts read it, so the gate proves the move against Git history,
proves that each guard fails when sabotaged, and exercises the built binary.
Every sabotage runs in a disposable copy of the repository, and every binary
command uses a disposable target directory and Roundfix Home. No command
reaches GitHub, a provider or a live remote.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact.
2. MUST verify through the built binary that `roundfix skills check` exits `0`,
   and that `roundfix skills install` into a disposable target directory writes
   `roundfix/SKILL.md` and every `roundfix/references/<name>.md` of the
   eighteen reference files.
3. MUST repeat the move proof of the TechSpec's "The move proof" section for
   both documents, comparing each splitting commit with its parent, and record
   the two commits and the compared line counts. The skill comparison leaves
   out the front matter and the reference index. The command reference
   comparison leaves out the command index and ignores link targets.
4. MUST prove in a disposable copy that each guard fails when sabotaged, and
   record each observed failure:
   - deleting one non-blank line from a reference file fails the move proof;
   - copying one section into a second reference file fails the move proof;
   - removing one row from the reference index fails the skill layout tests;
   - moving a required phrase out of every file of the skill makes
     `roundfix skills check` report it;
   - adding one word to the `### QA settlement` section of the entry file
     fails `TestSettlementGuidanceIsOneTable`.
5. MUST verify, by executing the tests task_01, task_02 and task_03 name
   against the built tree, that every contract which pins text of either
   document passes. It MUST run the sweep command of the TechSpec's "Pins read
   the tree" section and record that no reader of either document's content
   remains outside the reader, apart from the three the TechSpec keeps on the
   entry file.
6. MUST verify that the skill's `SKILL.md` is at most 20,000 bytes, that its
   index names every file under `references/` and no other, and that the
   command index names every file under `docs/user-guide/commands/`.
7. MUST verify, by executing task_04's tests against the built tree, that two
   same-Wave Tasks declaring different command files raise no
   `SC-WAVE-COLLISION` and that two declaring the same one do. It MUST verify
   that the `write-tasks` skill, its mirror and its template carry
   `declares the one command file it changes` and the two example paths, and
   that `make skills-sync-check` exits `0`.
8. MUST record, as evidence this Spec did not author:
   - the published guidance at
     <https://learn.microsoft.com/en-us/agent-framework/agents/skills>, which
     asks for a `SKILL.md` under 500 lines with reference material in separate
     files loaded on demand, compared with the line count of the entry file
     after the split;
   - this repository's history: the serial four-Task chain of Spec 0187's
     Task Graph, each Task declaring the Roundfix skill, and the refusals on
     `TestSettlementGuidanceIsOneTable` recorded in the QA reports of Specs
     0181 and 0187.

   When a source is unavailable, it MUST record the row as blocked with its
   reason.
9. MUST verify that this Spec's own artifacts satisfy the promise rule, and
   check whether it introduced, changed or retired a glossary term that
   `CONTEXT.md` does not carry.
10. MUST verify from the repository history that each Task's changed files stay
    within its declarations or its `## Recorded paths`, that every Governed
    Path is bounded in `_authorization.md`, and that every other change to a
    governed file is an output of a sanctioned regeneration command. It MUST
    record that the `Makefile`, `go.mod` and the CI workflows are unchanged.
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
13. MUST NOT write to the live Run Database under `~/.roundfix`, and MUST NOT
    leave a sabotage in the audited worktree.

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

- `newest="$(find 'docs/specs/0194-a-skill-and-a-command-guide-read-one-command-at-a-time/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

`_prd.md` → Goals 1-4; Core Features 1-4; Success Metrics 1-4; Acceptance
evidence; `_techspec.md` → Testing Approach 1-5; API Contracts 1-2; The move
proof; ADR-0080; ADR-0081; ADR-0088; ADR-0091; ADR-0093; ADR-0094; ADR-0096;
ADR-0104; ADR-0117; ADR-0130; ADR-0155; ADR-0156; ADR-0166; ADR-0167;
ADR-0176; ADR-0187.

## Result
