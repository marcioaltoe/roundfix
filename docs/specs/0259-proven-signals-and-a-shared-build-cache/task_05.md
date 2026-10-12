---
task: task_05
spec: 0259-proven-signals-and-a-shared-build-cache
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec makes every signal go only to a
proven process, gives every waiting test helper a lifeline and has the suite
guard kill Helper Residue, has `make` share one Go build cache, and lets
`doctor` and `gc` report and reclaim Worktree Build Caches. Its central row is
ten consecutive full suites without F6. No row reads or writes the real
`~/.roundfix`, signals a process the row did not start, reaches a provider or
reads a credential.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify API Contract 1, API Contract 2, API Contract 3 and API Contract
   4 by executing task_01's Signal rule and its six store tests, and by seeding,
   in a disposable copy of the repository only, one direct `syscall.Kill` with
   `SIGKILL` in a test file and recording that the rule fails naming its file
   and line.
3. MUST verify API Contract 5 and API Contract 6 by executing task_02's two
   death tests and three suite guard tests, and by recording the suite guard
   lines of every package in the ten suites of Requirement 6.
4. MUST verify API Contract 7 by executing `TestMakeSharesOneGoBuildCache` and
   by running `make -n test` with `GOCACHE` unset, recording the exported
   path, and once inside `codex sandbox -P :workspace` when the `codex` CLI is
   installed, recording whether `make build` succeeds there with no variable
   set; when it is not installed, the sandbox part of the row is recorded as
   `blocked (environment: codex CLI absent)`.
5. MUST run Surface Transcript 1, Surface Transcript 2 and Surface Transcript 3
   through the built `bin/roundfix` with `HOME` set to a temporary home,
   `NODE_OPTIONS` unset and a temporary Git repository whose `.roundfixrc.yml`
   declares `.gocache`, seeding one terminal Run whose Run Worktree under the
   worktree location holds a `.gocache` of known size and one Active Run whose
   worktree holds another. It records each output and exit code, compares them
   with the transcripts, and records that the dry run changed no byte (by
   checksum of the home and both worktrees) and that the live run removed only
   the terminal Run's `.gocache`.
6. MUST measure Success Metric 1: at the audited head, run
   `go test -count=1 -parallel 16 ./...` ten consecutive times on one machine,
   with the shared cache, recording for each run its start time, load average,
   wall time, exit status and the count of lines matching `exited -1`,
   `signal: killed` or `suiteguard: helper residue`. The row passes only when
   all ten exit 0 with every count zero. On the first failing run it stops,
   keeps the full output, records the failing test, the child's exit and any
   signal, and fails the row; it never reruns a failed suite to obtain a pass.
   This row is the acceptance of F6 that ADR-0266 decides.
7. MUST replay Success Metric 2, Success Metric 3, Success Metric 4 and Success
   Metric 5 from Requirements 2 to 5 and record each against its metric.
8. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the LWN article, the Unix Stack Exchange
   answer, psutil issue 2400, the Codex permissions guide and `go help cache`
   from the installed toolchain. For each source it reaches it records the
   quoted passage and whether it still supports the design; `go help cache`
   is read offline. For a source it cannot reach it records
   `blocked (environment: network denied: <host>)`.
9. MUST perform the glossary check of `docs/agents/domain.md`. It records that
   `CONTEXT.md` defines **Proven Signal**, **Helper Lifeline**, **Helper
   Residue** and **Worktree Build Cache** and carries the revised **GC
   Command**, and whether any other term the Spec introduced needs one.
10. MUST verify from the repository history that each Task's changed files
    stay within its declarations or its `## Recorded paths`; that every
    Governed Path changed is bounded in `_authorization.md`; that the
    Makefile changed only in task_03 and only its `GOCACHE` default; that the
    CI workflows, `go.mod`, `go.sum` and every Baseline module did not change;
    that the Roundfix skill's version was raised and recorded; and that no
    `### QA settlement` section changed. This row reads Task commits, so it
    declares a `commit_range` input.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance, and support it with the pre-PR equivalent of each
    control: approval from the recorded delivery authority, checks from the
    Daemon's repository Verification at the audited head, no unresolved
    threads because no Pull Request exists, no Merge-Ready acceptance yet, and
    the audited head named as the claimed candidate.
12. MUST NOT accept a row whose only evidence is that a file was read. It MUST
    NOT read or write `~/.roundfix`, signal a process it did not start, use a
    broad `pkill`, send any request to a provider, or start a workflow.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Execute the Task tests, the seeded rule violation and the Makefile cases.
- [ ] Run the three Surface Transcripts on a fixture home.
- [ ] Run the ten full suites and record each.
- [ ] Record the outside evidence, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Ten consecutive full suites exit 0 with no child killed by a signal and no Helper Residue line.
- [ ] The seeded signal violation fails the rule with its file and line.
- [ ] The three transcripts match and the dry run changes no byte.
- [ ] The scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0259-proven-signals-and-a-shared-build-cache/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-8; User Stories 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; API Contract 7; API Contract 8; API Contract 9; API Contract 10; API Contract 11; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Invariants 1-7; Current behavior; Build Order 5
- ADR-0266; ADR-0213; ADR-0172; ADR-0171; ADR-0259
