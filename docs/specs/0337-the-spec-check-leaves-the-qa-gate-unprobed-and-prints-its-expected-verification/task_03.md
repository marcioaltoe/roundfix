---
task: task_03
spec: 0337-the-spec-check-leaves-the-qa-gate-unprobed-and-prints-its-expected-verification
status: pending
type: qa
complexity: medium
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes
`roundfix spec check --run-verification` leave the QA gate's derived
Verification unprobed, through the `daemon.TaskIsQAGate` predicate the
Daemon's routing also reads, and makes `SC-QA-VERIFICATION-AUTHORED` print the
exact `## Verification` block to copy, in text and in JSON, which the
`write-tasks` skill, the Roundfix spec reference and the spec command guide
teach.

Every row runs on this repository's own files, temporary repositories and
temporary copies of the audited head. No row starts an Agent session or an
ACP adapter, reads or writes the real `~/.roundfix`, reads another
repository, reaches the network, sends a request to a provider, the judge or
a reviewer, or reads a credential. Every `go test` it runs carries an
explicit `-timeout`. Every row that runs `bin/roundfix` sets `HOME` and
`TMPDIR` to fresh temporary directories, because the CLI loads
`$HOME/.roundfix/config.yml`, and removes them afterwards. The QA never
writes into the checkout while a test suite runs. No row compares one
command's output with another's byte for byte.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 to 5 and Invariants 1 to 5 by running
   `go test -count=1 -timeout 20m -v` with
   `TestSpecCheckNeverProbesTheQAGate` and
   `TestSpecCheckLeavesCompletedTasksAlone` over `./internal/cli`,
   `TestAuthoredQAVerificationPrintsTheExpectedBlock`,
   `TestAFindingRendersItsExpectedBlock`,
   `TestAuthoredQAVerificationIsRefusedByTaskName` and
   `TestDerivedQAVerificationIsAcceptedByTaskStageChecker` over
   `./internal/speccheck`, and the whole `./internal/daemon` package, recording
   each pass.
3. MUST reproduce the Fluxus case through the built binary in a temporary
   clone of the audited head that the row creates and removes: commit a
   fixture Spec `9001-fixture` with a minimal `_prd.md` (`status: active`), a
   `_tasks.md` naming `qa: task_02`, a pending `type: backend` `task_01`
   whose only Verification is `test -f corrective-output`, a pending
   `type: qa` `task_02` carrying the derived block, and
   `qa/qa-report-2026-10-31.md` whose frontmatter is `verdict: fail`. With
   `h="$(mktemp -d)"; t="$(mktemp -d)"`, run
   `HOME="$h" TMPDIR="$t" bin/roundfix spec check 9001-fixture --run-verification`
   in text and in `--format json`, and record that the Verification lines
   report `task_01` `honest` and `task_02` `not probed` with `(QA gate)`, that
   JSON lists `task_02` under `verification.notProbed` with reason `qa_gate`,
   and that no Verification line is `vacuous`; then `rm -rf "$h" "$t"` and the
   clone. The fixture's PRD findings are expected and are not part of the row.
4. MUST reproduce the Quaestor case in the same kind of clone: replace the
   fixture QA Task's Verification with one authored bullet whose code span is
   `test -f docs/specs/9001-fixture/qa/report.md`, run
   `bin/roundfix spec check 9001-fixture --stage tasks` in text and in
   `--format json` with fresh `HOME` and `TMPDIR`, record that the text prints
   `expected block:` and a fenced `markdown` block after the finding's `fix:`
   line and that the JSON finding carries `expected`; write that `expected`
   text in place of the Task file's `## Verification` section, rerun, and
   record that no `SC-QA-VERIFICATION-AUTHORED` remains.
5. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, with `-timeout` on every `go test`, and record that
   the named test fails, with the case whose input makes the mutation
   observable:
   - `TaskIsQAGate` always returns false:
     `TestSpecCheckNeverProbesTheQAGate` fails on the subtest "a preserved
     failed QA report beside an honest corrective Task" (the gate is
     `vacuous`), and `./internal/daemon` fails, `TestAQATaskRecordsNoPaths`
     among its failures;
   - the QA branch removed from `probeSpecVerifications`:
     `TestSpecCheckNeverProbesTheQAGate` fails on its three subtests;
   - reason `completed` written for the QA gate:
     `TestSpecCheckNeverProbesTheQAGate` fails on the verdict line;
   - `Expected` left empty in `AuthoredQAVerification`:
     `TestAuthoredQAVerificationPrintsTheExpectedBlock` fails;
   - the bullet of `ExpectedQAVerificationBlock` written without its code
     span: `TestAuthoredQAVerificationPrintsTheExpectedBlock` fails;
   - `omitempty` dropped from the `expected` tag:
     `TestAFindingRendersItsExpectedBlock` fails on the second finding.
6. MUST execute the tests task_01 and task_02 name in their Verification.
7. MUST search `docs/user-guide/`, every copy of the Roundfix skill and every
   copy of the `write-tasks` skill and record that the guide and both copies
   of the `spec` reference state that the QA gate is never probed, with
   reason `qa_gate`, and document `SC-QA-VERIFICATION-AUTHORED` with its
   expected block and `expected` field; that both copies of the `write-tasks`
   template and `SKILL.md` teach copying the block; and that no passage still
   says a pending QA gate's vacuous command makes the check exit 1.
8. MUST record, as outside-evidence rows, three sources this Spec did not
   write, each with its date and none blocked:
   - the Fluxus report quoted in
     `references/2026-11-01-the-spec-check-probe-judges-the-qa-gates-derived-verification-vacuous.md`
     (2026-11-01): compare offline its shape (honest corrective Task, kept
     failed QA Report, gate `vacuous`) with the row of requirement 3;
   - the Quaestor report quoted in
     `references/2026-11-01-no-authoring-path-shows-the-qa-gates-derived-verification.md`
     (2026-11-01): compare offline its shape (QA Task from the template,
     `SC-QA-VERIFICATION-AUTHORED`, no way to learn the text) with the row of
     requirement 4;
   - the rustc book passage `_prd.md` → Acceptance evidence quotes ("An
     optional string of a suggested replacement for this span to solve the
     issue."): compare it with API Contract 5.
9. MUST record that the PRD declares no Unreachable Acceptance and count no
   row in `rows_blocked_declared`.
10. MUST perform the glossary check (`docs/agents/domain.md`): record that the
    PRD and TechSpec declare `changes: **Vacuous Verification**`, that
    `CONTEXT.md` defines it with the QA gate sentence, and that the Task
    commit range changes no other glossary entry.
11. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - the Task commit range changes none of these: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`,
      `docs/references/skill-coverage.json`, `internal/spec/task.go`,
      `internal/spec/qa.go`, anything under `internal/baseline/assets/` or
      `docs/agents/`, any skill other than the Roundfix and `write-tasks`
      files the Tasks declare;
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
    `h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" bin/roundfix archive 0337-the-spec-check-leaves-the-qa-gate-unprobed-and-prints-its-expected-verification --plan --offline; rm -rf "$h" "$t"`
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
- [ ] Run the named tests and reproduce both reports through the binary.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the wording search, the outside sources, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] An honest corrective Task beside a kept failed QA Report passes the probe, with the gate `not probed` and reason `qa_gate`, in the tests and through the binary.
- [ ] `SC-QA-VERIFICATION-AUTHORED` prints the expected block in text and JSON, and copying it clears the finding, in the tests and through the binary.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

- `newest="$(find 'docs/specs/0337-the-spec-check-leaves-the-qa-gate-unprobed-and-prints-its-expected-verification/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-4; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-6; Invariants 1-5; Declared breaks; Testing Approach; Build Order 3
