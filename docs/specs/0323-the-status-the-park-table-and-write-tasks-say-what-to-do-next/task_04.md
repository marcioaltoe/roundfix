---
task: task_04
spec: 0323-the-status-the-park-table-and-write-tasks-say-what-to-do-next
status: pending
type: qa
complexity: high
---

# Task 04: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes
`roundfix follow-up --status` name, in the manual classification state, how
many commits need a manual release classification and the command that
classifies them, with no blank field in text or JSON; makes the Park Class
table of the `deliver` guide and both Roundfix reference copies list blockers
only and state the `secret-found` next action as the queue prints it, pinned
by a documentation contract; and adds the shared-test-file rule to the
`write-tasks` skill and the **Recorded Path** glossary entry.

Every row runs on this repository's own files, temporary copies of the
audited head and temporary repositories. No row reads or writes the real
`~/.roundfix`, starts an Agent session, reaches the network, sends a request
to a provider, the judge or a reviewer, reads a credential, or reads another
repository. Every `go test` it runs carries an explicit `-timeout`, and every
row that runs a built `roundfix` sets `HOME` and `TMPDIR` to fresh temporary
directories, because the CLI loads `$HOME/.roundfix/config.yml`, and puts
first on `PATH` a temporary directory holding a `gh` stub that prints
`gh stub` to standard error and exits `1`, so the status view's Pull Request
source never reaches GitHub.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1-4 and Invariants 1-4 by running
   `TestFollowUpStatusNamesTheManualClassification` and
   `TestFollowUpStatusReportsTheNextRelease` in `./internal/cli`, and
   `TestTheParkActionsAreDocumentedAsTheQueueClassifiesThem` and
   `TestTheParkActionContractReportsAnActionInTheClassTable` with
   `-tags docscontract` in `./internal/docscontract`, each with `-v`,
   recording each named test's and subtest's pass.
3. MUST exercise the built binary end to end on the reported shape. In a
   fresh temporary Git repository on branch `main` with one commit
   `feat: seed widget`, tag `v0.4.1`, then commit
   `feat: expose widget status`, `chore: change widget behavior` and
   `chore: change widget storage`, each changing its own `internal/widget<N>.go`
   file, and with no remote:
   - with the audited `bin/roundfix`, record that `follow-up --status` exits
     `0` and prints `Next release: not proposed (manual_classification_required): 2 commit(s) since v0.4.1 need a manual release classification`
     and `  release classification: classify them with roundfix release plan --from v0.4.1 --to HEAD --impact <none|patch|minor|major> --reason <text>`,
     and neither `Next release:  ` nor `  release : `; and that
     `--format json` carries the first sentence in `release.reason`, the
     second as the `detail` of the one `release_classification` item of
     `maintainer`, and `release.proposed_version` `""`;
   - with a binary built from the Delivery Base in a temporary `git archive`
     copy (`go build -buildvcs=false -o <tmp>/roundfix-base ./cmd/roundfix`),
     record that the same command prints
     `Next release:  from v0.4.1 (manual_classification_required)` and
     `  release : manual classification required`, the defect the first
     adopted Backlog Entry reports;
   - in a second temporary repository with the tag and one `fix:` commit,
     record that both binaries print the same `Next release: v0.4.2 from v0.4.1 (ready)` line.
4. MUST run each sabotage on its own temporary copy of the audited head made
   with `git archive`, and record that the named test fails, naming the case
   whose input makes the mutation observable:
   - `runFollowUpStatus` counts `len(plan.Changes)` instead of the blocking
     commits: `TestFollowUpStatusNamesTheManualClassification` fails (three
     changes, two blocking);
   - the manual-state case is removed from `renderFollowUpStatusText`:
     `TestFollowUpStatusNamesTheManualClassification` fails on
     `Next release:  from v0.4.1`;
   - `release.reason` is left nil in the manual state:
     `TestFollowUpStatusNamesTheManualClassification` and
     `TestFollowUpStatusReportsTheNextRelease/manual` fail;
   - the joined `finding` row is restored in
     `docs/user-guide/commands/deliver.md`:
     `TestTheParkActionsAreDocumentedAsTheQueueClassifiesThem` fails naming
     that file;
   - the `secret-found` sentence of `skills/roundfix/references/deliver.md`
     loses its final clause `then run roundfix deliver retry <slug>`:
     `TestTheParkActionsAreDocumentedAsTheQueueClassifiesThem` fails naming
     that file.
5. MUST execute the commands task_01, task_02 and task_03 name in their
   Verification, and run `go test -timeout 15m -count=1 ./skills` and
   `go test -timeout 20m -count=1 -tags docscontract ./internal/docscontract`.
6. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - read `docs/adr/0166-the-daemon-records-the-paths-a-task-changed-without-declaring-them.md`
     in the audited checkout and quote "Recording discloses a change after
     the fact, and it reserves nothing.", the premise of the `write-tasks`
     rule;
   - read the third adopted Backlog Entry under `references/` and record the
     Fluxus report it quotes (ai-memory message `01a16a54`, 2026-10-23),
     recorded session evidence the Spec did not produce;
   - because this Spec's authorization forbids network access, record the
     Command Line Interface Guidelines as `blocked (environment: network
     prohibited by Spec authorization)`, quoting the passage the PRD cites.
     Such a row never decides a qualifying `partial`.
7. MUST record the PRD's single Unreachable Acceptance entry as exactly one
   row, `blocked (declared unreachable)`, counted once in
   `rows_blocked_declared`: no row may read the operator's real repository
   with the released binary.
8. MUST perform the glossary and documentation checks:
   - record that `CONTEXT.md`'s **Recorded Path** entry holds the sentence of
     `_techspec.md` → Shared test files in the Task Graph after its unchanged earlier sentences and
     before its unchanged `_Avoid_` line (`docs/agents/domain.md`);
   - search `docs/user-guide`, `.agents/skills`, `skills` and `CONTEXT.md`
     for a Park Class table row that carries a next action after a colon,
     and for the follow-up text `manual classification required`, and record
     what it finds;
   - record that the `roundfix` and `write-tasks` skills' two version fields
     are equal in each canonical copy and mirror and recorded in
     `skills/testdata/owned-skill-versions.json`.
9. MUST verify from the repository history:
   - each Task's changed files stay within its declarations or its `##
     Recorded paths`;
   - every changed Governed Path is listed in `_authorization.md` `paths`;
   - these are byte-identical to the Delivery Base: the Makefile, the CI
     workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every file under
     `internal/releaseplan/`, `internal/delivery/`, `internal/spec/`,
     `internal/speccheck/`, `internal/daemon/`, `internal/baseline/`,
     `docs/agents/` and `docs/adr/`, and every skill other than `roundfix`
     and `write-tasks`;
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
    `HOME="$(mktemp -d)" bin/roundfix archive 0323-the-status-the-park-table-and-write-tasks-say-what-to-do-next --plan --offline`
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

- [ ] The audited binary names the manual classification count and command with no blank field where the Delivery Base binary prints a blank version.
- [ ] The Park Class `finding` row lists blockers only in every copy, and the contract ties the `secret-found` sentence to the queue's action.
- [ ] Each sabotage makes its named test fail, and the scope audit finds no change outside the declared files.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0323-the-status-the-park-table-and-write-tasks-say-what-to-do-next/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-3; User Experience; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; Interfaces; Shared test files in the Task Graph; API Contracts 1-3; Invariants 1-5; Testing Approach; Build Order 4
