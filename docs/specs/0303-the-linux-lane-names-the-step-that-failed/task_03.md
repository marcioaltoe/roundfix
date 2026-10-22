---
task: task_03
spec: 0303-the-linux-lane-names-the-step-that-failed
status: pending
type: qa
complexity: high
---

# Task 03: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes the Linux Lane script
mark a failed package-selection `go list` with its exit code and exit `4`,
and makes the lane name `go list exited <n>` as a Lane Change Failure below
`128` and as a Lane Container Failure with its signal from `128`, while every
other outcome keeps its bytes.

Every row runs on temporary repositories and Git configuration under the
Run's temporary directory, with the compiled Docker and Go stand-ins of the
lane tests and the host's `sh`. No row starts a container, runs `docker` or
`roundfix deliver lane` against a real Docker daemon, pulls an image, reads
or writes the real `~/.roundfix`, reaches the network, sends a request to a
provider, the judge or a reviewer, reads a credential, or reads another
repository. The gate never runs `roundfix spec judge`, never runs `roundfix
review` against a real provider, and never runs an archive plan without
`--offline`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1-5:
   - Run `TestLinuxLaneScriptNamesAFailedGoList` and record, per row, the
     script's exit code, the marker line and whether `go test` ran
     (Invariant 1).
   - Run `TestLinuxLaneGoListFailureIsNamed`,
     `TestDeliverLaneGoListFailureTranscript`,
     `TestDeliverLaneGoListSignalTranscript` and
     `TestGateParksWithTheBlockerItNames`, and record the gate reasons, the
     blockers, the lane log endings, the two transcripts with their exit
     codes, and the Park Classes `finding` and `environment` (Invariant 2).
   - Run `TestLinuxLaneOutcomeClassification`,
     `TestLinuxLaneRunsTheChangedPackagesInAContainer`,
     `TestLinuxLaneFailureFailsTheGate`,
     `TestLinuxLaneContainerFailureIsToldApartFromAChangeFailure`,
     `TestLinuxLaneReportsThePackagesTheContainerSelected`,
     `TestDeliverLaneChangeFailureTranscript` and
     `TestDeliverLaneContainerFailureTranscript`, and record that every
     outcome other than exit `4` keeps its text and that the container's
     arguments are unchanged apart from the script branch (Invariants 3-5).
   - On temporary copies of the audited head made with `git archive`, apply
     one sabotage per copy and record that the named test fails, with the
     case whose input makes the mutation observable:
     - make the script's `go list` branch exit `"$code"` instead of `4`:
       `TestLinuxLaneScriptNamesAFailedGoList` fails in its `go list` exits
       `1` row (exit `1`, not `4`);
     - drop the branch's `printf` of the marker:
       `TestLinuxLaneScriptNamesAFailedGoList` fails in its `go list` exits
       `1` row (no marker line);
     - restore `&& if` after `go list` (the branch removed):
       `TestLinuxLaneScriptNamesAFailedGoList` fails in its `go list` exits
       `137` row;
     - remove the classifier's exit `4` case:
       `TestLinuxLaneOutcomeClassification` fails in its `go list change`
       row, and `TestLinuxLaneGoListFailureIsNamed` fails;
     - classify a marker code of `128` or more as a change failure:
       `TestLinuxLaneOutcomeClassification` fails in its `go list killed`
       row, and `TestDeliverLaneGoListSignalTranscript` fails on its exit
       code;
     - compare the code with `> 128` instead of `>= 128`:
       `TestLinuxLaneOutcomeClassification` fails in its `go list signal
       zero` row;
     - keep the first marker line instead of the last:
       `TestLinuxLaneOutcomeClassification` fails in its `go list last
       marker` row;
     - accept a marker anywhere in a line (`strings.Contains` instead of a
       prefix): `TestLinuxLaneOutcomeClassification` fails in its `go list
       embedded marker ignored` row;
     - stop dropping the trailing `\r`: `TestLinuxLaneOutcomeClassification`
       fails in its `go list carriage return` row;
     - accept any text after the prefix as the code (trim spaces before
       `strconv.Atoi`): `TestLinuxLaneOutcomeClassification` fails in its
       `go list spaced code ignored` row.
3. MUST execute the tests task_01 and task_02 name in their Verification.
4. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - read the adopted Backlog Entry under `references/` in the audited
     checkout, which records Quaestor's report of 2026-10-21, and record that
     its quoted finding names `go list` running before `go test` and that the
     `go list exits 1` row of `TestLinuxLaneScriptNamesAFailedGoList` shows
     `go test` did not run;
   - run `go help list` from the Go toolchain the gate uses and record that
     its text contains, with whitespace runs collapsed to one space, "With
     the -e flag, the list command never prints errors to standard error and
     instead processes the erroneous packages with the usual printing.",
     which the PRD quotes from https://pkg.go.dev/cmd/go; because this Spec's
     authorization forbids network access, the page itself is recorded as
     `blocked (environment: network prohibited by Spec authorization)`. Such
     a row never decides a qualifying `partial`.
5. MUST record Success Metric 7 as exactly one row, `blocked (declared
   unreachable)`, against the PRD's single Unreachable Acceptance entry, and
   count it once in `rows_blocked_declared`: no row may start a container.
6. MUST perform the glossary and documentation checks:
   - record that `CONTEXT.md` revises **Lane Change Failure** with the `go
     list` run and links ADR-0316, that **Lane Container Failure** keeps its
     bytes, that the glossary contract tests pass, and whether any other term
     the Spec introduced needs an entry (`docs/agents/domain.md`);
   - record the bullet and sentences `docs/user-guide/commands/deliver.md`
     and both copies of the Roundfix `deliver` reference gained, the
     `<cause>` outcome lines, and that `TestTheSkillCoverageMapIsCurrent`
     passes;
   - search the guides, the agent guides and the skill copies for wording
     that still says every Lane Change Failure is `go test exited <n>`, and
     record what it finds.
7. MUST verify from the repository history:
   - each Task's changed files stay within its declarations or its `##
     Recorded paths`;
   - every Governed Path changed is bounded in `_authorization.md`;
   - these are byte-identical to the Delivery Base: the Makefile, the CI
     workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every Baseline module
     under `internal/baseline/assets/`, `docs/agents/setup-context.json`,
     every file under `docs/history/`, every ADR,
     `internal/cli/deliver_lane.go`, `internal/cli/deliver_workflow.go`,
     every non-test file under `internal/delivery`,
     `docs/references/skill-coverage.json`, and every `_test.go` file under
     `internal/cli` and `internal/delivery` that exists at the Delivery Base
     other than `internal/cli/deliver_linux_lane_test.go`,
     `internal/cli/deliver_lane_test.go` and
     `internal/delivery/gate_blocker_test.go`;
   - the Roundfix skill's version was raised and recorded;
   - no `### QA settlement` section changed.

   This row reads Task commits, so it declares a `commit_range` input.
8. MUST record the non-waivable Pull Request row as
   `blocked (environment: no open Pull Request)`, naming the Pull Request row
   in its provenance. It supports the row with the pre-PR equivalent of each
   control:
   - approval from the recorded delivery authority;
   - checks from the Daemon's repository Verification at the audited head;
   - no unresolved threads, because no Pull Request exists;
   - no Merge-Ready acceptance yet;
   - the audited head named as the claimed candidate.
9. MUST, for the Outcome for the Archive Record, run only
   `bin/roundfix archive 0303-the-linux-lane-names-the-step-that-failed --plan --offline`
   with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
   report's `## Outcome` as one paragraph that does not end with a colon. It
   MUST NOT do any of these:
   - accept a row whose only evidence is that a file was read;
   - start a container or call a real Docker daemon;
   - read or write `~/.roundfix`;
   - record any environment value;
   - send any request to a provider, the judge or a reviewer;
   - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test and each sabotage on temporary copies of the audited head.
- [ ] Compare the Backlog Entry and `go help list` with the asserted behavior, and record the blocked outside source.
- [ ] Record the declared row, the glossary and documentation checks and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A failed `go list` is named with its exit code and Park Class through the gate and `deliver lane`, and each removed guard makes its named test fail.
- [ ] Every other lane outcome and the container's arguments are unchanged.
- [ ] The scope audit finds no undeclared Governed Path change and no change to the files row 7 names.

## Verification



The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0303-the-linux-lane-names-the-step-that-failed/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-3; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → Current behavior; API Contracts 1-5; Surface Transcripts; Invariants 1-5; Testing Approach; Build Order 3
- ADR-0316
