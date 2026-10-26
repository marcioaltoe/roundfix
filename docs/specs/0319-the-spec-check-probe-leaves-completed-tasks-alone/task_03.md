---
task: task_03
spec: 0319-the-spec-check-probe-leaves-completed-tasks-alone
status: completed
type: qa
complexity: high
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes
`roundfix spec check --run-verification` probe only the Tasks a Run would
still execute, chosen by the Daemon's exported `TaskAwaitsWork`, and report
each completed Task's command as `not probed` with its reason in text, JSON
and the verdict line, while a vacuous, unknown or malformed command of a Task
still to run keeps failing; and documents the rule in the `spec` guide, the
Roundfix and `qa-gate` skills, ADR-0148 and `CONTEXT.md`.

Every row runs on this repository's own files, temporary copies of the
audited head and temporary repositories. No row reads or writes the real
`~/.roundfix`, starts an Agent session, reaches the network, sends a request
to a provider, the judge or a reviewer, reads a credential, or reads another
repository. Every `go test` it runs carries an explicit `-timeout`, and every
row that runs `bin/roundfix` sets `HOME` and `TMPDIR` to fresh temporary
directories, because the CLI loads `$HOME/.roundfix/config.yml`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1-4 and Invariants 1-3 and 5 by running
   `TestSpecCheckLeavesCompletedTasksAlone`, `TestSpecCheckRunVerification`,
   `TestSpecCheckReportsAMalformedCommand`,
   `TestSpecCheckNamesAnUncommittedVerificationSource` and
   `TestSpecCheckNamesNoSourceWhenCommitted` in `./internal/cli`,
   `TestTaskAwaitsWorkSelectsEveryStatusButCompleted`,
   `TestTaskCycleRerunsStaleTasksAndSkipsCompletedTasks`,
   `TestTaskCycleQAStepSkippedUnlessEveryTaskCompleted`,
   `TestTaskCycleQAOnlyRunWhenEveryTaskAlreadyCompleted` and
   `TestTaskCycleQAStepRequiresEveryGraphTaskCompleted` in `./internal/daemon`,
   and `TestVerdictLineStatesProbeCoverage` in `./internal/speccheck`, each
   with `-v`, recording each named subtest's pass.
3. MUST exercise the built binary end to end on the Vortex shape. In a fresh
   temporary Git repository, copy `docs/agents`, `docs/adr`, `docs/workflow`
   and `docs/specs/clean` from the audited head's
   `internal/speccheck/testdata/repo`, then write the `clean` Spec's
   `_tasks.md` and Task files with: Task 01 `completed` whose Verification is
   `touch <absolute marker path>` and
   `roundfix-verification-binary-that-does-not-exist --run`, and Task 02
   `pending` with `test -f task-output-two.txt`, the shape
   `internal/cli/spec_check_completed_tasks_test.go` writes, and commit:
   - with the audited `bin/roundfix`, record that
     `bin/roundfix spec check clean --run-verification` exits `0`, prints
     both `- task_01: not probed — … (Task completed)` lines, `- task_02:
     honest` and a verdict line ending `not probed (completed Tasks): 2.`,
     that `--format json` carries both commands under
     `verification.notProbed` with reason `completed` and only Task 02 under
     `verification.commands`, and that the marker file does not exist;
   - with a binary built from the Delivery Base in a temporary `git archive`
     copy (`go build -buildvcs=false -o <tmp>/roundfix-base ./cmd/roundfix`),
     record that the same command exits `1` with Task 01's commands reported
     `vacuous` and `unknown`, the defect the adopted Backlog Entry reports;
   - change Task 02's command to `true`, commit, and record that the audited
     binary exits `1` with Task 02 `vacuous` while Task 01 stays `not
     probed`.
4. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, and record that the named subtest fails, naming the
   case whose input makes the mutation observable:
   - `TaskAwaitsWork` returns `true`:
     `TestSpecCheckLeavesCompletedTasksAlone/a_resumed_Spec_probes_only_the_Tasks_still_to_run`
     fails (Task 01's marker command runs and is `vacuous`), and
     `TestTaskCycleRerunsStaleTasksAndSkipsCompletedTasks` fails (three Agent
     runs instead of two), which proves both callers read the one function;
   - `TaskAwaitsWork` returns true only for `pending`:
     `.../a_resumed_Spec_probes_only_the_Tasks_still_to_run` fails (Tasks 03
     and 04 are `in_progress` and `failed`), and
     `TestTaskAwaitsWorkSelectsEveryStatusButCompleted` fails;
   - `probeSpecVerifications` replaces `daemon.TaskAwaitsWork` with its own
     rule that also skips `qa` Tasks:
     `TestSpecCheckLeavesCompletedTasksAlone/a_pending_QA_gate_is_still_probed`
     fails;
   - completed Tasks are skipped without a `notProbed` entry or text line:
     `.../every_Task_completed` fails;
   - `NotProbed` is not passed to the verdict line:
     `.../every_Task_completed` fails on its verdict line;
   - completed Task files are dropped from the uncommitted sources:
     `.../an_uncommitted_completed_Task_source_is_still_named` fails.
5. MUST execute the tests task_01 and task_02 name in their Verification,
   and run `go test -timeout 15m -count=1 ./skills` and
   `go test -timeout 20m -count=1 -tags docscontract ./internal/docscontract`.
6. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - read the adopted Backlog Entry under `references/` in the audited
     checkout and record that its verbatim message reports completed Task 01
     probed at `gesttione-solutions/vortex@a961a29` with exit 127, verdict
     `unknown` and final exit `1`, the shape row 3 reproduces with the
     Delivery Base binary;
   - read `docs/references/2026-08-14-a-clean-checker-left-eight-vacuous-commands-standing.md`
     in the audited checkout and quote "a check that cannot be complete should
     say what it did not cover";
   - because this Spec's authorization forbids network access, record the
     TAP 14 specification as `blocked (environment: network prohibited by
     Spec authorization)`, quoting the passage the PRD cites. Such a row never
     decides a qualifying `partial`.
7. MUST record the PRD's single Unreachable Acceptance entry as exactly one
   row, `blocked (declared unreachable)`, counted once in
   `rows_blocked_declared`: no row may read Vortex or another repository.
8. MUST perform the glossary and documentation checks:
   - record that `CONTEXT.md`'s **Vacuous Verification** entry holds the
     sentence of `_techspec.md` → Interfaces after its unchanged first two
     sentences (`docs/agents/domain.md`);
   - search `docs/user-guide`, `.agents/skills`, `skills`, `docs/agents`,
     `docs/adr` and `CONTEXT.md` for wording that says `--run-verification`
     probes completed Tasks or reports them `vacuous`, and record what it
     finds;
   - record that the `roundfix` and `qa-gate` skills' two version fields are
     equal in each canonical copy and mirror and recorded in
     `skills/testdata/owned-skill-versions.json`.
9. MUST verify from the repository history:
   - each Task's changed files stay within its declarations or its `##
     Recorded paths`;
   - every changed Governed Path is listed in `_authorization.md` `paths`;
   - these are byte-identical to the Delivery Base: the Makefile, the CI
     workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every file under
     `internal/baseline/` and `docs/agents/`, every skill other than
     `roundfix` and `qa-gate`, and every ADR other than ADR-0148;
   - no `### QA settlement` section changed.

   This row reads Task commits, so it declares a `commit_range` input.
10. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
11. MUST, for the Outcome for the Archive Record, run only
    `HOME="$(mktemp -d)" bin/roundfix archive 0319-the-spec-check-probe-leaves-completed-tasks-alone --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's `## Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - start an Agent session;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test, the end-to-end binary rows and each sabotage on temporary copies.
- [ ] Record the outside sources, the declared row, the glossary and documentation checks and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A completed Task's commands never run and are named `not probed`, and the check exits `0` on the Vortex shape where the Delivery Base binary exits `1`.
- [ ] A vacuous, unknown or malformed command of a Task still to run still fails the check.
- [ ] Each sabotage makes its named subtest fail, and the scope audit finds no change outside the declared files.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0319-the-spec-check-probe-leaves-completed-tasks-alone/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-3; User Experience; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; Interfaces; API Contracts 1-5; Invariants 1-5; Testing Approach; Build Order 3
