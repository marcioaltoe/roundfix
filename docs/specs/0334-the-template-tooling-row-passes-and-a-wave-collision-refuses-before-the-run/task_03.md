---
task: task_03
spec: 0334-the-template-tooling-row-passes-and-a-wave-collision-refuses-before-the-run
status: completed
type: qa
complexity: high
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the Spec Check accept
the templates' no-mutation Tooling authority row when its record grants no
Governed Path, writes the row's ending into both templates under a contract
test, keeps a completed Task out of every Wave Collision, and makes
`roundfix implement` refuse a real collision before it records a Run.

Every row runs on this repository's own files, temporary repositories, the
implement test fixtures with the fake Agent runner, and temporary copies of
the audited head. No row starts an Agent session or an ACP adapter, reads or
writes the real `~/.roundfix`, reads another repository, reaches the network,
sends a request to a provider, the judge or a reviewer, or reads a
credential. Every `go test` it runs carries an explicit `-timeout`. Every row
that runs `bin/roundfix` sets `HOME` and `TMPDIR` to fresh temporary
directories, because the CLI loads `$HOME/.roundfix/config.yml`, and removes
them afterwards. The QA never writes into the checkout while a test suite
runs. No row compares one command's output with another's byte for byte.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1, 2 and 3 and Invariants 1 and 2 by running
   `go test -count=1 -timeout 20m -v` over `./internal/speccheck` with
   `TestANoMutationRowCitingAnEmptyGrantIsNotUnbounded` and over
   `./internal/docscontract` (with `-tags docscontract`) with
   `TestTheTemplateNoMutationToolingRowPassesTheSpecCheck`, recording each
   pass.
3. MUST verify Success Metrics 4 to 6 and Invariants 3 to 6 by running
   `go test -count=1 -timeout 20m -v` with `TestACompletedTaskNeverCollides`
   over `./internal/spec`, `TestQAGateJudgesACorrectiveTaskOnTheCurrentGraph`
   over `./internal/daemon` (recording both its `corrective Task beside a
   completed Task` and `two unordered corrective Tasks` subtests), and the
   three tests of `internal/cli/implement_wave_collision_test.go` over
   `./internal/cli`, recording each pass.
4. MUST reproduce Surface Transcript 1 through the built binary in a
   temporary Git repository the row creates and removes: a work branch, a
   committed `shared/notes.txt`, and a Spec `0001-wave` with `qa: declined`
   and two pending `type: backend` Tasks whose only Verification is
   `grep -q settled shared/notes.txt`, with
   `h="$(mktemp -d)"; t="$(mktemp -d)"` and `worktree:` `concurrency: 2` in
   `$h/.roundfix/config.yml`; run `HOME="$h" TMPDIR="$t" bin/roundfix
   implement --spec 0001-wave --no-input`, compare stdout, stderr and the
   exit code with the transcript (`<slug>` is `0001-wave`), record that no
   `.roundfix/worktrees` directory was created under `$h`, then
   `rm -rf "$h" "$t"` and the repository. With concurrency 2 the refusal
   comes before any Agent runtime is probed, so no Agent starts.
5. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that
   the named test fails, with the case whose input makes the mutation
   observable:
   - make every row whose reason says `no protected tooling mutation`
     declare no mutation, ignoring the cited record:
     `TestANoMutationRowCitingAnEmptyGrantIsNotUnbounded` fails on its case
     whose record lists `Makefile`;
   - restore `declaresProtectedToolingMutation` to return `true` for any
     cited record: `TestTheTemplateNoMutationToolingRowPassesTheSpecCheck`
     fails on its `paths: []` case;
   - delete the `bounded files: none` ending from the canonical
     `techspec-template.md`: `TestTheTemplateNoMutationToolingRowPassesTheSpecCheck`
     fails naming that template;
   - remove the completed-status skip from `Collisions`:
     `TestACompletedTaskNeverCollides` fails on its `completed` case and
     `TestImplementRunsATaskThatSharesAFileWithACompletedTask` fails on the
     collision text;
   - skip every Task whose status is not `pending`:
     `TestACompletedTaskNeverCollides` fails on its `failed` and
     `in_progress` cases;
   - remove the `RefuseWaveCollisions` call from `runImplementCommand`:
     `TestImplementRefusesAWaveCollisionBeforeTheRun` fails on exit `1` and
     one recorded Run;
   - refuse at a Task Capacity of 1 (`taskCapacity < 1` in place of
     `taskCapacity <= 1`):
     `TestImplementAtTaskCapacityOneLeavesCollidingTasksToTheRun` fails on
     the collision text.
6. MUST execute the tests task_01 and task_02 name in their Verification.
7. MUST search `CONTEXT.md`, `docs/user-guide/` and every skill copy for
   `SC-WAVE-COLLISION` and record that each passage that describes which
   Tasks collide either states that a completed Task never collides or
   describes only Tasks that are not completed; and search the two template
   copies of each authoring skill for `With no protected tooling mutation`
   and record that each names the `bounded files: none` ending.
8. MUST record, as outside-evidence rows, two sources this Spec did not
   write. First, Quaestor's report quoted verbatim in both adopted Backlog
   Entries under `references/` (ai-memory message `01a17e01`,
   2026-10-27): compare offline its item 1 row (`applicable — no protected
   tooling mutation proposed or authorized`, `paths: []`, `bounded files:
   none`) with the row `TestTheTemplateNoMutationToolingRowPassesTheSpecCheck`
   composes, and its item 2 shape (a corrective Task beside a completed Task
   sharing one file) with `TestImplementRunsATaskThatSharesAFileWithACompletedTask`.
   Second, the POSIX `make` passage `_prd.md` → Acceptance evidence quotes
   (IEEE Std 1003.1-2024, `make`): a target is brought up to date in
   parallel only when "The target's prerequisites (if any) are up-to-date";
   compare it with API Contract 3's rule that a completed Task never joins
   a Wave. Neither row is blocked; name each source with its date.
9. MUST record the Unreachable Acceptance criterion as exactly one row,
   `blocked (declared unreachable)`, against the PRD's single entry, and
   count it once in `rows_blocked_declared`.
10. MUST perform the glossary check (`docs/agents/domain.md`): record that
    `CONTEXT.md` defines **Wave Collision**, citing ADR-0347 and no Spec.
11. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      anything under `internal/baseline/assets/` or `docs/agents/`, any skill
      other than the `write-prd`, `write-techspec` and Roundfix files the
      Tasks declare;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
12. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance. It supports the row with the pre-PR equivalent of
    each control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
13. MUST, for the Outcome for the Archive Record, run only
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" bin/roundfix archive 0334-the-template-tooling-row-passes-and-a-wave-collision-refuses-before-the-run --plan --offline; rm -rf "$h" "$t"`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's ` Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests and reproduce the transcript through the binary.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the declared row, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] The templates' no-mutation row with an empty grant passes, and a row whose record lists a Governed Path is still refused.
- [ ] A completed Task never collides in any reader, and two pending Tasks still do.
- [ ] `roundfix implement` reproduces Surface Transcript 1 and records no Run.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification



The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0334-the-template-tooling-row-passes-and-a-wave-collision-refuses-before-the-run/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-4; Surface Transcript 1; Invariants 1-6; Declared breaks; Testing Approach; Build Order 3
