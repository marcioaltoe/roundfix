---
task: task_05
spec: 0262-a-run-worktree-under-a-symlinked-root-is-cleaned
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec resolves a recorded worktree path
within its Worktree Root, compares worktree paths through their longest
existing prefix, proves every consumer under a symbolic-linked worktree
location, documents the reclaim with `roundfix reconcile --apply`, and makes
the disposal lock tests wait on a synchronization point. No row reads or
writes the real `~/.roundfix`, follows the maintainer's real worktree link,
reaches a provider or reads a credential.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify API Contract 1 and API Contract 2 by executing task_01's and
   task_02's tests by name and the pinned refusal tests
   `TestInspectTerminalRunUnsafePath` and
   `TestCountRetainedTerminalRunsStillReportsASymlinkedRecordedRoot`.
3. MUST verify API Contract 3 and Success Metric 1 through the built
   `bin/roundfix` with `HOME` set to a temporary home, `NODE_OPTIONS` unset, a
   worktree location that is a symbolic link to a temporary directory and a
   temporary Git repository: start no Agent, seed
   one terminal implement Run whose Run Worktree under the linked location is
   clean and whose Run Branch is contained in `main`, using the repository's
   test helpers or a recorded fixture, run `roundfix reconcile --format json`
   and record the classification, then `roundfix reconcile --apply` and
   record that the Run Worktree directory, its `git worktree list` entry and
   the Run Branch are gone. When seeding a Run without an Agent is not
   possible through the binary, the row records the CLI tests of task_03 as
   its evidence and says why.
4. MUST verify Success Metric 2 by recording, in a disposable copy of the
   Run Worktree layout, that a recorded path whose repository directory is a
   link to the real one stays `unknown` with
   `recorded Run Worktree path is unsafe` and that the real worktree and Run
   Branch survive, through task_01's subtests.
5. MUST verify Success Metric 3 by executing task_02's item tests and
   task_03's release and daemon tests by name.
6. MUST verify API Contract 4 and Success Metric 4: run the two disposal tests
   with `-count=50 -race -parallel 16` and record the pass counts; then, in a
   disposable copy of the repository only, make `runWorktreeCommand` run the
   command without `withWorktreeAdminLock`, run the two tests once, record
   that both fail, and discard the copy. It MUST record that
   `internal/worktree/disposal_test.go` contains no `context.WithTimeout` or
   `context.WithDeadline`.
7. MUST repeat the measurement of `_techspec.md` → Current behavior with the
   installed Git: add a worktree through a symbolic-linked directory in a
   temporary repository and record that `git worktree list --porcelain` names
   the resolved directory.
8. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: Git's `builtin/worktree.c` passage and the
   Go blog passages, quoted in the PRD. For each source it reaches it records
   whether the quoted passage is still there and still supports the design; for
   a source it cannot reach it records
   `blocked (environment: network denied: <host>)` and relies on the quoted
   passage and the measurement of Requirement 7.
9. MUST perform the glossary check of `docs/agents/domain.md`. It records that
   `CONTEXT.md` defines **Worktree Root**, and whether any other term the Spec
   introduced needs one.
10. MUST search the user guides, `CONTEXT.md`, the agent guides and the
    Roundfix skill references for wording that says a symbolic link in a
    recorded Run Worktree path is always unsafe, and record the result.
11. MUST verify from the repository history that each Task's changed files
    stay within its declarations or its `## Recorded paths`; that no Governed
    Path changed (`paths: []`); that no production file under `internal/cli`,
    `internal/daemon` or `internal/delivery` changed; that no skill, Baseline
    module, Makefile, CI workflow, `go.mod` or `go.sum` changed; and that no
    `### QA settlement` section changed. This row reads Task commits, so it
    declares a `commit_range` input.
12. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance, and support it with the pre-PR equivalent of each
    control: approval from the recorded delivery authority, checks from the
    Daemon's repository Verification at the audited head, no unresolved
    threads because no Pull Request exists, no Merge-Ready acceptance yet, and
    the audited head named as the claimed candidate.
13. MUST NOT accept a row whose only evidence is that a file was read. It MUST
    NOT read or write `~/.roundfix`, use a broad `pkill`, send any request to a
    provider, or start a workflow.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Execute the Task tests and the pinned refusal tests.
- [ ] Run reconcile and apply through the binary on a linked location.
- [ ] Run the lock tests 50 times and the sabotaged copy once.
- [ ] Record the Git measurement, the outside evidence, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A Run under a linked location is classified and released by the binary.
- [ ] A link at the repository or worktree directory stays unsafe.
- [ ] The lock tests pass 50 times and fail with the lock removed.
- [ ] The scope audit finds no undeclared change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0262-a-run-worktree-under-a-symlinked-root-is-cleaned/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-7; User Stories 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; Invariants 1-5; Current behavior; Build Order 5
- ADR-0268; ADR-0053; ADR-0259
