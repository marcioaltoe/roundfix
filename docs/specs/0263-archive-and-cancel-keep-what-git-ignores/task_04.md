---
task: task_04
spec: 0263-archive-and-cancel-keep-what-git-ignores
status: pending
type: qa
complexity: medium
---

# Task 04: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes `roundfix archive` and
`roundfix cancel` refuse while Git ignores a file under the Spec folder,
listing each one and writing nothing, and makes a failed folder removal keep
every write, exit 1 and print the command that finishes it. No row reads or
writes the real `~/.roundfix`, reaches a provider or reads a credential.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST run Surface Transcript 1, Surface Transcript 2, Surface Transcript 3
   and Surface Transcript 4 through the built `bin/roundfix` with `HOME` set to
   a temporary home, `NODE_OPTIONS` unset and temporary Git repositories built
   as the TechSpec's Surface Transcripts section describes, recording each
   output and exit code and comparing them with the transcripts.
3. MUST verify API Contract 1 and Success Metric 1: for both refusals, record a
   checksum of every file under the repository, ignored files included, and the
   output of `git status --porcelain --untracked-files=all --ignored` before
   and after, and that they are equal and no Archive Record exists. It MUST
   also record that an uncommitted change plus an ignored file yields today's
   uncommitted reason.
4. MUST verify API Contract 2, API Contract 3 and Success Metric 2: after each
   recovery transcript, record that the record and every other write exist,
   run the printed finishing command exactly as printed, and record that the
   folder is gone and `git status --porcelain --untracked-files=all` equals the
   status a successful run on an identical fixture leaves. For archive, repeat
   with `--pull-request 7` and record that the record carries it.
5. MUST verify API Contract 4 by recording `roundfix archive --help` and
   `roundfix cancel --help` and comparing the refusal and exit-code lines.
6. MUST execute the tests task_01 and task_02 name in their Verification, and
   replay Success Metric 3 from Requirement 2.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the Git manual for `git status`, compared
   offline with `git status -h` and the installed manual when present, and
   `go doc os.RemoveAll` from the installed toolchain. For each it records the
   quoted passage and whether it still supports the design; for a source it
   cannot reach it records `blocked (environment: network denied: <host>)`.
8. MUST perform the glossary check of `docs/agents/domain.md`: record that
   `CONTEXT.md` defines **Ignored Spec File** and **Folder Removal Recovery**,
   and whether any other term the Spec introduced needs one.
9. MUST verify from the repository history that each Task's changed files stay
   within its declarations or its `## Recorded paths`; that every Governed Path
   changed is bounded in `_authorization.md`; that the symbolic-link
   destination parent checks in `internal/spec` did not change; that the
   Makefile, the CI workflows, `go.mod`, `go.sum` and every Baseline module did
   not change; that both skill versions were raised and recorded; and that no
   `### QA settlement` section changed. This row reads Task commits, so it
   declares a `commit_range` input.
10. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance, and support it with the pre-PR equivalent of each
    control: approval from the recorded delivery authority, checks from the
    Daemon's repository Verification at the audited head, no unresolved
    threads because no Pull Request exists, no Merge-Ready acceptance yet, and
    the audited head named as the claimed candidate.
11. MUST NOT accept a row whose only evidence is that a file was read. It MUST
    NOT read or write `~/.roundfix`, send any request to a provider, or start a
    workflow.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the four Surface Transcripts on temporary repositories.
- [ ] Prove the refusals change no byte and the finishing command completes the removal.
- [ ] Record the outside evidence, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The four transcripts match the built binary.
- [ ] Both refusals leave every byte and the Git status unchanged.
- [ ] The printed finishing command leaves the tree a successful run leaves.
- [ ] The scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0263-archive-and-cancel-keep-what-git-ignores/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-2; Core Features 1-3; Success Metric 1; Success Metric 2; Success Metric 3; Acceptance evidence
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; Invariants 1-6; Current behavior; Build Order 4
- ADR-0269; ADR-0247; ADR-0262; ADR-0261
