---
task: task_03
spec: 0264-baseline-group-usage-errors-name-their-own-help
status: pending
type: qa
complexity: medium
---

# Task 03: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes the `baseline`
dispatcher's usage errors name the group once and point at the help of the
command path the user typed, guards every group's usage error in a
repository test, and describes the bare-group usage error in the `roundfix`
skill. Every row runs against the built binary, a temporary home and an empty
temporary working directory. No row reaches a provider or the network, starts
an Agent Session, reads a credential, or reads or writes the real
`~/.roundfix`.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact,
   with `bin/roundfix` rebuilt from the audited head first, and record that
   Spec 0260 is merged in the audited head's ancestry.
2. MUST reproduce Surface Transcripts 1, 2, 3, 4, 5 and 6 of `_techspec.md`
   through the built binary, from an empty temporary working directory with
   `HOME` set to a temporary directory, comparing standard output, standard
   error and exit code under the transcript conventions, and record each
   result. For each transcript it also runs the pointed command with
   `--help` and records that it exits 0 and lists every command of the group.
3. MUST measure Success Metric 1: run every group of the built binary's root
   usage, derived as `_techspec.md` → API Contract 3 states, with the word
   `no-such-command`, plus the three bare `baseline` groups, in the same
   temporary home and directory, and record the number of groups, the number
   of calls, and every call whose standard error repeats a word in a row or
   points at a help that does not exit 0 or does not list every command of
   the group. The metric holds at 0 such calls.
4. MUST measure Success Metric 2: `TestEveryGroupUsageErrorPointsAtItsCommands`
   passes at the audited head, and on a temporary copy of the tree where
   `printBaselineUsageError` again prints the `baseline profile --help`
   pointer it fails naming a `baseline` group.
5. MUST record that `baseline profile` with an unknown command keeps today's
   text and exit 2, that `docs/references/behavior-surfaces.json` is
   byte-identical to its bytes at `3820bdd4`, and that the existing
   `baseline_profile_test.go` cases pass (Invariant 5).
6. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence, through the passages `_prd.md` quotes: the
   Command Line Interface Guidelines Errors passage and the TS-16 error
   message passage. For each it records whether the passage still supports
   the behavior, and it opens no network connection to do so.
7. MUST perform the glossary check of `docs/agents/domain.md` and record that
   the Spec declares no term and whether any term it introduced needs one.
8. MUST verify from the repository history that task_01's and task_02's
   changed files stay within their declarations or their `## Recorded paths`,
   that every Governed Path changed is bounded in `_authorization.md`, and
   that the user guide, `docs/references/skill-coverage.json`,
   `internal/cli/cli_test.go`, the Makefile, the CI workflows, `go.mod`,
   `go.sum` and `.roundfixrc.yml` did not change, and that
   `internal/cli/command_help_test.go` changed only in the rows task_02
   Requirement 5 names. This row reads Task commits, so it declares a
   `commit_range` input.
9. MUST record the non-waivable Pull Request row as
   `blocked (environment: no open Pull Request)`, naming the Pull Request row
   in its provenance, supported by the pre-PR equivalent of each control: the
   maintainer's recorded delivery authority, the Daemon's repository
   Verification at the audited head, no unresolved threads because no Pull
   Request exists, no Merge-Ready acceptance yet, and the audited head named
   as the claimed candidate.
10. MUST NOT accept a row whose only evidence is that a file was read, and
    MUST NOT write outside temporary directories and the Spec's `qa/` folder.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Reproduce Surface Transcripts 1-6 and run the group sweep.
- [ ] Measure Success Metric 2 with the sabotaged copy.
- [ ] Record the outside evidence, the glossary check and the scope audit.
- [ ] Write the QA report with the Pull Request row blocked.

## Acceptance Criteria

- [ ] Every Surface Transcript reproduces through the built binary and every
      pointer resolves to a help that lists the group's commands.
- [ ] The group sweep records 0 failing calls.
- [ ] The class guard fails on the sabotaged copy and passes at the audited
      head.
- [ ] The scope audit finds no undeclared Governed Path and an unchanged
      Behavior Surface Record.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0264-baseline-group-usage-errors-name-their-own-help/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`


## References

- [_prd.md](_prd.md) — Goals; User Stories 1-3; Core Features 1-4; Success Metric 1; Success Metric 2; Acceptance evidence
- [_techspec.md](_techspec.md) — API Contract 1; API Contract 2; API Contract 3; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; Surface Transcript 5; Surface Transcript 6; Invariants 1-5; Testing Approach; Build Order 3
- ADR-0256; ADR-0184; ADR-0259
