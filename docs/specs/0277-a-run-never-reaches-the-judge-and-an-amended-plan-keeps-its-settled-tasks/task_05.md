---
task: task_05
spec: 0277-a-run-never-reaches-the-judge-and-an-amended-plan-keeps-its-settled-tasks
status: pending
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

The authored terminal gate for this Spec. It declares what the matrix covers
and settles the Spec on evidence. This Spec gives every Agent Session and
Verification command of a Run the Run Marker and makes the judge send nothing
while it is set, withholds `TYPESAFE_API_KEY` from Runs, removes the race of
the bootstrap serialization test, makes every test compile for Windows with a
contract that proves it, lets Task Carry-Forward keep a settled Task whose
Binding Text an amendment left unchanged, and makes the history cleanup read
every ADR reference form adopters use. Every row runs on temporary
repositories and homes. No row reads or writes the real `~/.roundfix`, reads a
shell startup file or a runtime's state directory, reaches the network, sends
a request to a provider or the judge, reads a credential, or reads another
repository; the gate never runs `roundfix spec judge` and never runs an archive
plan without `--offline` except the Surface Transcript row, which runs with
`ROUNDFIX_RUN_ID` set and every Judge Key Variable and `TYPESAFE_API_KEY`
unset.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1 and Surface Transcript 1 through the built
   `bin/roundfix` with `HOME` set to a temporary home, `NODE_OPTIONS` unset and
   every Judge Key Variable and `TYPESAFE_API_KEY` unset: on a temporary Git
   repository holding a committed active Spec whose only candidate is one
   `references/` file, run `ROUNDFIX_RUN_ID=qa-run bin/roundfix archive <slug> --plan`
   and compare its standard output, standard error and exit code with the
   transcript's lines under its controlled variations. Record that the tests
   task_01 names pass, which prove with a counting transport and a local server
   that no request is made with a key present, that each prepare and prompt
   invocation and each Verification command carries `ROUNDFIX_RUN_ID` with its
   Run ID, and that none receives a Judge Key Variable or `TYPESAFE_API_KEY`
   (Invariants 1 to 3).
3. MUST verify Success Metric 2: run
   `go test -count=50 -run '^TestBootstrapSerializesAcrossSiblings$' ./internal/worktree`
   and record that it passes; then, on a temporary copy of the Delivery Base,
   the merge base of the audited head with `main`, made with `git archive`,
   change `tryWriteBootstrapRelease` to keep its writer open after writing
   until the released sibling's end event has arrived and a non-blocking
   write-only open of the same FIFO finds a reader, and record that
   `TestBootstrapSerializesAcrossSiblings` fails with `watched work ended`;
   apply the same hold, keyed to the released sibling's own FIFO, to a
   temporary copy of the audited head and record that it passes (Invariant 4).
4. MUST verify Success Metric 3: record that
   `GOOS=windows go vet -tags 'repocontract docscontract' ./...` exits 0 at the
   audited head and exits 1 at the Delivery Base; that
   `TestEveryTestCompilesForWindows` passes; that, on a temporary copy of the
   audited head, adding a use of `syscall.Kill` to an unconstrained test file
   makes it fail with a diagnostic naming that file; and that the coverage
   record lists each test task_02 moved for `darwin` and `linux` only while
   `TestCoverageEquivalence` passes (Invariant 5).
5. MUST verify Success Metric 4 by running the tests task_03 names and
   recording, from their output, the carried Task with `amendedInputs`, the
   refusal text for a changed named section, and the unchanged retry refusal;
   and, on a temporary copy of the audited head, make `PlanningBinding` always
   return `whole` and record that
   `TestCarryForwardCarriesATaskOverAnAmendmentOutsideItsBindingText` fails
   (Invariant 6).
6. MUST verify Success Metric 5 through the built binary with a temporary
   home: on a temporary Git repository holding a copy of the audited tree's
   `docs/agents/setup-context.json` and `docs/agents/docs-layout.md`, one
   minimal active Spec, the annotated tag `history-full`, a superseded ADR
   under `docs/history/adr/` whose `superseded_by` is written as number and
   slug, and an accepted successor under `docs/adr/` whose `supersedes` names it
   by number and slug, record that `roundfix history sanitize` removes the
   superseded ADR with no refusal and no lineage citation; record that the
   tests task_04 names pass (Invariant 7).
7. MUST execute the tests task_01 to task_04 name in their Verification, and
   prove the Run Marker guard fails without its implementation: on a temporary
   copy of the audited head made with `git archive`, remove the `RunID` check
   from `judge.Run` and record that `TestJudgeSendsNothingInsideARun` fails.
8. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence: the Codex shell snapshot source, Codex
   issue 30971, POSIX `read()` and Go's build constraint text, quoting each
   passage the PRD cites and stating whether it supports the design. Because
   this Spec's authorization forbids network access, each lookup that needs it
   is recorded as
   `blocked (environment: network prohibited by Spec authorization)`, which
   never decides a qualifying `partial`; the Go build constraint text is read
   from the local `go help buildconstraint`.
9. MUST perform the glossary check of `docs/agents/domain.md`: record that
   `CONTEXT.md` defines **Run Marker** and **Binding Text** and carries the
   revised **Judge Key Variables** and **Task Carry-Forward**, that the
   glossary contract tests pass, and whether any other term the Spec
   introduced needs an entry.
10. MUST verify from the repository history that each Task's changed files
    stay within its declarations or its `## Recorded paths`; that every
    Governed Path changed is bounded in `_authorization.md`; that the
    Makefile, the CI workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every
    Baseline module, `internal/worktree/worktree.go`, `internal/daemon/qa_format.go`
    and `internal/notify/notify.go` did not change; that the Roundfix skill
    version was raised and recorded; that no `### QA settlement` section
    changed; and that no guide or skill copy still says that a Task whose Spec
    moved is never replayed or that a superseded ADR names its successor only
    as `ADR-NNNN`. This row reads Task commits, so it declares a
    `commit_range` input.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance, and support it with the pre-PR equivalent of each
    control: approval from the recorded delivery authority, checks from the
    Daemon's repository Verification at the audited head, no unresolved
    threads because no Pull Request exists, no Merge-Ready acceptance yet, and
    the audited head named as the claimed candidate.
12. MUST, for the Outcome for the Archive Record, run only
    `bin/roundfix archive 0277-a-run-never-reaches-the-judge-and-an-amended-plan-keeps-its-settled-tasks --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset. It MUST NOT
    accept a row whose only evidence is that a file was read, read or write
    `~/.roundfix`, record any environment value, send any request to a
    provider or the judge, write into the checkout while a test suite runs, or
    start a workflow.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Replay the archive transcript and the history shape through the built binary.
- [ ] Run the bootstrap stress, the Windows compile and each sabotage on temporary copies.
- [ ] Record the outside evidence, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The archive plan inside a Run matches Surface Transcript 1 and no judge request is made.
- [ ] The old bootstrap design fails under the held writer and the new one passes 50 runs.
- [ ] Every test compiles for Windows, and the contract fails on a planted Unix-only call.
- [ ] Each removed guard makes its named test fail, and the scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0277-a-run-never-reaches-the-judge-and-an-amended-plan-keeps-its-settled-tasks/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-5; Core Features 1-9; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Acceptance evidence; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; API Contract 7; API Contract 8; API Contract 9; Surface Transcript 1; Invariants 1-7; Current behavior; Build Order 5
- ADR-0283; ADR-0284; ADR-0285; ADR-0273; ADR-0104; ADR-0276
