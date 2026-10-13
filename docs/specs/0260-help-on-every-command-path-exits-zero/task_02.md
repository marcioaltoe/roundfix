---
task: task_02
spec: 0260-help-on-every-command-path-exits-zero
status: completed
type: qa
complexity: medium
---

# Task 02: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes `--help` and `-h` exit 0
with a usage block on `qa-report accept` and the `baseline capabilities`,
`baseline skills` and `baseline assets` groups, sweeps every command path for
help in a repository test, fingerprints each Behavior Surface from its own
help, and corrects the `roundfix` skill. Every row runs against the built
binary, a temporary home and an empty temporary working directory. No row
reaches a provider or the network, starts an Agent Session, reads a
credential, or reads or writes the real `~/.roundfix`.

## Requirements

1. MUST run the repository Verification and record its result as a gate fact,
   with `bin/roundfix` rebuilt from the audited head first.
2. MUST reproduce Surface Transcripts 1, 2, 3, 4, 5 and 6 of `_techspec.md`
   through the built binary, from an empty temporary working directory with
   `HOME` set to a temporary directory, comparing standard output, standard
   error and exit code under the transcript conventions, and record each
   result. For transcripts 4-6 it also records that standard output equals
   `roundfix baseline --help`'s standard output.
3. MUST measure Success Metric 1: run `--help` and `-h` over every command
   path and group prefix of the built binary's root usage, derived as
   `_techspec.md` → API Contract 3 states, in the same temporary home and
   directory, and record the number of paths, the number of calls, and every
   call that does not exit 0 with standard output starting with `Usage:` and
   an empty standard error. The metric holds at 0 failing calls.
4. MUST measure Success Metric 2: `TestEveryCommandPathPrintsItsHelp` passes
   at the audited head, and on a temporary copy of the tree where the help
   check added to `qa-report accept` is removed it fails naming
   `qa-report accept --help`.
5. MUST measure Success Metric 3: the two Behavior Surface Record contract
   tests pass under `-tags docscontract` at the audited head,
   `docs/references/behavior-surfaces.json` is byte-identical to its bytes at
   `7a8b7f70`, and on the same temporary copy the contract fails naming
   `qa-report accept`.
6. MUST record that a bare `baseline capabilities` and `qa-report accept`
   with no path still exit 2 with today's usage error, and that
   `qa-report accept` on a report path keeps its eligibility result
   (Invariant 3), running the existing `qa_report_test.go` and
   `baseline_profile_test.go` cases.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence, through the passages `_prd.md` quotes: the
   GNU Coding Standards help-option passage and the Command Line Interface
   Guidelines help passages. For each it records whether the passage still
   supports the behavior, and it opens no network connection to do so.
8. MUST perform the glossary check of `docs/agents/domain.md` and record that
   the Spec declares no term and whether any term it introduced needs one.
9. MUST verify from the repository history that task_01's changed files stay
   within its declarations or its `## Recorded paths`, that every Governed
   Path changed is bounded in `_authorization.md`, and that the user guide,
   `docs/references/skill-coverage.json`, the Makefile, the CI workflows,
   `go.mod`, `go.sum` and `.roundfixrc.yml` did not change. This row reads
   Task commits, so it declares a `commit_range` input.
10. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance, supported by the pre-PR equivalent of each control: the
    maintainer's recorded delivery authority, the Daemon's repository
    Verification at the audited head, no unresolved threads because no Pull
    Request exists, no Merge-Ready acceptance yet, and the audited head named
    as the claimed candidate.
11. MUST NOT accept a row whose only evidence is that a file was read, and
    MUST NOT write outside temporary directories and the Spec's `qa/` folder.
12. Addendum (2026-10-13, after QA fail F-01 of Run
    `run_20261013T031317Z_0f4d6fa287f6d6a7`): the comparison revision of
    Requirement 5 is the delivery base, the commit on the default branch from
    which task_01 started (the parent of task_01's commit; at the time of this
    addendum `3820bdd4`), not `7a8b7f70`. Unrelated Specs merged between
    authoring and delivery updated `docs/references/behavior-surfaces.json`;
    Success Metric 3 asks that this Spec leave the record unchanged, so the row
    MUST compare the audited record with its bytes at the delivery base and
    MUST still report its SHA-256 at `7a8b7f70` for reference only, without
    failing on it. Requirements 1 and 4 and Success Metric 1 keep `7a8b7f70`
    as the characterization revision of the baseline count.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Reproduce Surface Transcripts 1-6 and run the help sweep.
- [ ] Measure Success Metrics 2 and 3 with the sabotaged copy.
- [ ] Record the outside evidence, the glossary check and the scope audit.
- [ ] Write the QA report with the Pull Request row blocked.

## Acceptance Criteria

- [ ] Every Surface Transcript reproduces through the built binary.
- [ ] The help sweep records 0 failing calls.
- [ ] Both new guards fail on the sabotaged copy and pass at the audited head.
- [ ] The scope audit finds no undeclared Governed Path and an unchanged
      Behavior Surface Record.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0260-help-on-every-command-path-exits-zero/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`


## References

- [_prd.md](_prd.md) — Goals; User Stories 1-3; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Acceptance evidence
- [_techspec.md](_techspec.md) — API Contract 1; API Contract 2; API Contract 3; API Contract 4; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; Surface Transcript 5; Surface Transcript 6; Invariants 1-4; Testing Approach; Build Order 2
- ADR-0256; ADR-0184; ADR-0259
