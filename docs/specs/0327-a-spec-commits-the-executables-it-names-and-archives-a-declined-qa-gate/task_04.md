---
task: task_04
spec: 0327-a-spec-commits-the-executables-it-names-and-archives-a-declined-qa-gate
status: pending
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec lets a Task commit keep the
untracked executables a Spec's `_authorization.md` lists under `executables`
and refuses every other one, also inside a new directory, in Task and Batch
commits (ADR-0340, amending ADR-0157); makes `roundfix archive` accept a
Declared QA Decline without an override (ADR-0154); and makes
`roundfix settle` name each worktree it removes.

Every row runs on temporary directories, temporary copies of the audited head
made with `git archive`, compiled test binaries, the built `bin/roundfix` and
this repository's own documents. No row reads or writes `~/.roundfix`,
reaches the network, sends a request to a provider, the judge or a reviewer,
reads a credential, or reads another repository. Every `go test` carries an
explicit `-timeout`, and every row that runs `bin/roundfix` sets `HOME` to a
fresh temporary directory, because the CLI loads `$HOME/.roundfix/config.yml`.
Evidence is written into the checkout only after the suites end.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST execute the commands task_01, task_02 and task_03 name in their
   Verification and record each exit status.
3. MUST verify Success Metrics 1 to 4 and Invariants 1 and 2: run, with
   `-count=3 -timeout 15m -v`, the tests task_01 names in
   `./internal/authorization` and `./internal/daemon`, then the whole
   `./internal/daemon` package once with `-count=1 -timeout 20m`, and record
   that each passes.
4. MUST verify Success Metric 5 and Invariant 3 through the built product,
   with a fresh `HOME`: in a temporary Git repository on branch `main`, commit
   a Spec `0001-demo` made of a `_prd.md` with `spec`, `status: active`,
   `created` and `surfaces` frontmatter and a title paragraph, a `_tasks.md`
   declaring `qa: declined` with `qa_reason: documentation-only demo` and one
   node, and a `task_01.md` with `status: completed`; run
   `bin/roundfix archive 0001-demo` and record exit `0`, its `archived 0001-demo -> `
   line and `disposition: no-qa` in the record. In a second copy of that
   repository, run it with `--qa-override --approval qa --reason qa` and
   record exit `2` with `QA archive override is not allowed because the Spec
   already qualifies for normal archive`. At authoring, the same two commands
   with the Delivery Base's binary exited `2` with `no QA Report found` and
   `0` with `archived 0001-demo with QA override`; record those as context
   only. Also run the three tests task_02 names with `-count=3`.
5. MUST verify Success Metric 6 and Invariant 4: run the six tests task_03
   names with `-count=3 -timeout 15m -v` and record each passes, and record
   from `bin/roundfix settle --help` (fresh `HOME`) that the usage names the
   removal.
6. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, and record that the named test fails, with the case
   whose input makes the mutation observable:
   - drop the authorized-set check in the stage filter:
     `TestTaskCommitKeepsAnAuthorizedExecutable` and
     `TestTaskCycleCommitsAnAuthorizedExecutableWithItsMode` fail (input: a
     listed `scripts/launch.sh` with mode `0755`);
   - pass the authorized set to the QA Report commit as well:
     `TestQAReportCommitRefusesAListedExecutable` fails (input: a listed path
     under `qa/evidence/`);
   - expand untracked directories after the filter in `prepareTaskCommit`
     again: `TestTaskCycleRefusesAnUnlistedExecutableInANewDirectory` fails
     (input: `newdir/helper.sh` in a directory Git has never seen);
   - remove the Batch commit's directory expansion:
     `TestResolveCycleRefusesAnExecutableInsideANewBatchDirectory` fails
     (input: the snapshot entry `newdir/`);
   - accept `*` in an `executables` entry:
     `TestAuthorizationRefusesAMalformedExecutablesEntry/glob` fails;
   - read the QA Report for a declined graph again:
     `TestArchiveAcceptsADeclaredQADecline` fails with `no QA Report found`;
   - drop the override refusal for a declined graph:
     `TestArchiveRefusesAnOverrideForAQualifyingDeclaredQADecline` fails
     (input: every Task completed, `--qa-override` given);
   - print `Run Worktree removed:` before checking for a pending
     integration: `TestSettleNamesNoRemovalWhenIntegrationIsPending` fails
     (input: a user branch that gained a commit after the Run started).
7. MUST verify Success Metric 7: record that each file task_01, task_02 and
   task_03 name in their phrase checks carries its rule, and that the
   Roundfix `implement` reference shows the Task commit's refusal line with
   ` unless the Spec's _authorization.md lists the path under executables`.
8. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - because this Spec's authorization forbids network access, record the
     Git documentation and the Git test as `blocked (environment: network
     prohibited by Spec authorization)`, naming the outside-evidence row in
     its provenance and quoting the passages the PRD cites; such a row never
     decides a qualifying `partial`;
   - run, in a temporary directory, `/bin/sh -c 'mkdir gone && cd gone && rmdir ../gone && git rev-parse HEAD; echo exit=$?'`
     and record the `fatal: Unable to read current working directory` line
     and `exit=128`, Git's own behaviour on the QA host;
   - quote Quaestor's report from the adopted Backlog Entry under
     `references/`.
9. MUST record the PRD's single Unreachable Acceptance entry as exactly one
   row, `blocked (declared unreachable)`, counted as one in
   `rows_blocked_declared`.
10. MUST perform the glossary check: record that `CONTEXT.md` defines
    **Authorized Executable** and **Declared QA Decline**, that the revised
    **Run Worktree**, **Settle Command** and **Archive Command** entries carry
    the sentences API Contracts 8 and 9 give, and that no bolded phrase of the
    Spec's artifacts lacks a glossary entry or a declaration
    (`docs/agents/domain.md`).
11. MUST verify from the Task commits of this Run (a `commit_range` input):
    - each Task's changed files stay within its declarations or its
      `## Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`;
    - no Task commit changes the Makefile, `.github/workflows/`, `go.mod`,
      `go.sum`, `.roundfixrc.yml`, `internal/baseline/`, `docs/agents/`,
      `docs/adr/` or the text of any `### QA settlement` section.
12. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request
    row in its provenance, supported by the pre-PR equivalent of each
    control: approval from the recorded delivery authority; checks from the
    Daemon's repository Verification at the audited head; no unresolved
    threads, because no Pull Request exists; no Merge-Ready acceptance yet;
    the audited head named as the claimed candidate.
13. MUST, for the Outcome for the Archive Record, run only
    `HOME="$(mktemp -d)" bin/roundfix archive 0327-a-spec-commits-the-executables-it-names-and-archives-a-declined-qa-gate --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's `## Outcome` as one paragraph that does not end with a colon.
    It MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge or a reviewer;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary, run the repository Verification and every Task Verification.
- [ ] Run the named tests, the built-product archive rows and each sabotage on temporary copies.
- [ ] Record the outside sources, the declared row, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] Every named test passes three times, and the built product archives a declined Spec and refuses the override for it.
- [ ] Each sabotage makes its named test fail.
- [ ] The scope audit finds no change outside the declared files and the granted Governed Paths.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0327-a-spec-commits-the-executables-it-names-and-archives-a-declined-qa-gate/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-9; Declared breaks; Invariants 1-4; Testing Approach; Build Order 4
