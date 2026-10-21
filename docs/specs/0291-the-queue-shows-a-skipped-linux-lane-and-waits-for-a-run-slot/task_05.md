---
task: task_05
spec: 0291-the-queue-shows-a-skipped-linux-lane-and-waits-for-a-run-slot
status: completed
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix
covers and settles the Spec on evidence. The Spec makes a Linux Lane skipped
for an environment reason visible on the Delivery Queue item, in the Pull
Request body and in `roundfix doctor`, and makes an item refused only by the
Active Implement Run ceiling wait for a Run slot up to `runs.slot_wait`
before it parks as it does today.

Every row runs on temporary repositories and homes. No row runs real Docker,
reads or writes the real `~/.roundfix`, reaches the network, sends a request
to a provider, the judge, a reviewer or GitHub, reads a credential, or reads
another repository. The gate never runs `roundfix doctor` outside a test,
never runs `roundfix spec judge`, never runs `roundfix review` against a real
provider, and never runs an archive plan without `--offline`.

## Requirements

1. MUST build the binary with `make build` and run the repository
   Verification, recording its result as a gate fact.
2. MUST verify Success Metrics 1 and 2:
   - Run `TestGateReturnsTheLaneEnvironmentSkip`,
     `TestDeliverStatusShowsTheLaneEnvironmentSkipTranscript`,
     `TestAGateLaneSkipIsRecordedOnTheItemAndInThePullRequestBody` and
     `TestAQuietGateLeavesNoLaneWarningAndReplacesAnEarlierOne`, and record
     each reason and whether it reached the item (Invariants 1 and 2).
   - Through the built product, reproduce Surface Transcript 1 of
     `_techspec.md`. In a temporary copy of the audited head made with `git
     archive`, seed a temporary home outside every checkout with a throwaway
     test that opens the Run Database there, creates a Delivery Queue for a
     temporary Git repository with one commit holding one item
     `0285-example`, and records it `merged` with the warning
     `delivery.WarningLinuxLaneSkipped` + `": image golang:1.26 is not present
     locally"`. Then, with `HOME` set to that home and the temporary
     repository as the working directory, run `bin/roundfix deliver status`
     and record that standard output, standard error and the exit code match
     the transcript.
3. MUST verify Success Metric 3:
   - Run `TestDoctorLinuxLaneCheckReportsDockerAndTheImage` and
     `TestDoctorPrintsTheLinuxLaneLineOnlyWhenTheLaneIsOn`, and record each
     rendered line and the Docker stand-in's recorded arguments (Invariant 3).
   - Record that `internal/cli/doctor_test.go` and
     `internal/cli/doctor_characterization_test.go` are byte-identical to the
     Delivery Base and pass at the audited head.
4. MUST verify Success Metrics 4 and 5:
   - Run the seven tests of `internal/delivery/run_slot_wait_test.go`, the
     three of `internal/config/run_slot_wait_test.go`,
     `TestDeliveryRunResultRecognizesOnlyTheCeilingRefusal`,
     `TestRunSpecReturnsARunSlotRefusalWithItsHolders`,
     `TestRunSlotsExcludeAHolderWhoseOwnerIsGone` and
     `TestDeliverStatusShowsARunSlotWaitTranscript`, and record the park
     reason of each ended wait (Invariants 4 to 6).
   - Through the built product, reproduce Surface Transcript 2 of
     `_techspec.md` the way Requirement 2 reproduces Surface Transcript 1,
     seeding one item `0291-example` in stage `running` with worktree path
     `<worktree>` replaced by a temporary directory and the blocker
     `delivery.RunSlotWaitText` of the transcript's deadline and holder.
   - Record that `TestRunCeilingIsCountedInsideTheCreateTransaction` passes
     at the audited head and that `internal/store/store_test.go` is
     byte-identical to the Delivery Base.
5. MUST prove each new check fails when its guard is removed. On temporary
   copies of the audited head made with `git archive`, apply one sabotage per
   copy and record that the named test fails:
   - set `SkipExpected` for every skip reason:
     `TestGateReturnsTheLaneEnvironmentSkip` fails;
   - stop `gateCandidate` from recording the warning:
     `TestAGateLaneSkipIsRecordedOnTheItemAndInThePullRequestBody` fails;
   - print the doctor line whatever `delivery.linux_lane` says:
     `TestDoctorPrintsTheLinuxLaneLineOnlyWhenTheLaneIsOn` fails;
   - ignore the queue deadline in the wait:
     `TestTheQueueDeadlineEndsARunSlotWait` fails;
   - call `RunSpec` again without waiting for free slots:
     `TestARunCeilingRefusalWaitsAndStartsWhenASlotFrees` fails;
   - remove the ceiling count inside `CreateRun`'s write transaction:
     `TestTwoQueuesWaitingForOneSlotStartOneRun` fails;
   - recognize every exit `2` as the ceiling refusal:
     `TestDeliveryRunResultRecognizesOnlyTheCeilingRefusal` fails;
   - count a holder whose owner process is gone:
     `TestRunSlotsExcludeAHolderWhoseOwnerIsGone` fails.
6. MUST execute the tests task_01, task_02, task_03 and task_04 name in their
   Verification.
7. MUST record, as evidence this Spec did not author, each source of
   `_prd.md` → Acceptance evidence:
   - quote the two adopted Backlog Entries' Symptom sections for the
     maintainer's observations of 2026-10-20;
   - record the store ceiling race test row from Requirement 4;
   - because this Spec's authorization forbids network access, record each
     GitHub documentation source as `blocked (environment: network prohibited
     by Spec authorization)`, quoting the passage the PRD cites. Such a row
     never decides a qualifying `partial`.
8. MUST record Success Metric 6 as exactly one row, `blocked (environment:
   real Docker is out of the Spec's reach)`, naming the PRD's one Unreachable
   Acceptance entry: the maintainer's `roundfix doctor` against the real
   Docker daemon after the release. Count that single row once in
   `rows_blocked_declared`. Never run Docker.
9. MUST perform the glossary check of `docs/agents/domain.md`:
   - record that `CONTEXT.md` defines **Lane Environment Skip** linking
     ADR-0303 and **Run Slot Wait** linking ADR-0304;
   - record that the glossary contract tests pass;
   - record whether any other term the Spec introduced needs an entry.
10. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `##
      Recorded paths`;
    - every Governed Path changed is bounded in `_authorization.md`;
    - these are byte-identical to the Delivery Base: the Makefile, the CI
      workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, every Baseline module
      under `internal/baseline/assets/`, `docs/agents/setup-context.json`,
      every file under `docs/history/`, `internal/cli/deliver_linux_lane_test.go`,
      `internal/cli/deliver_park_status_test.go` and
      `internal/cli/implement_test.go`;
    - the Roundfix skill version was raised and recorded;
    - no `### QA settlement` section changed;
    - no guide, agent guide or skill copy still says that a lane skip leaves
      only the notice or that the queue parks at the ceiling without waiting.

    This row reads Task commits, so it declares a `commit_range` input.
11. MUST record the non-waivable Pull Request row as
    `blocked (environment: no open Pull Request)`, naming the Pull Request row
    in its provenance. It supports the row with the pre-PR equivalent of each
    control:
    - approval from the recorded delivery authority;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
12. MUST, for the Outcome for the Archive Record, run only
    `bin/roundfix archive 0291-the-queue-shows-a-skipped-linux-lane-and-waits-for-a-run-slot --plan --offline`
    with every Judge Key Variable and `TYPESAFE_API_KEY` unset, and write the
    report's `## Outcome` as one paragraph that does not end with a colon. It
    MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - run Docker, or read or write `~/.roundfix`;
    - record any environment value;
    - send any request to a provider, the judge, a reviewer or GitHub;
    - write into the checkout while a test suite runs.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run each named test and each sabotage on temporary copies of the audited head.
- [ ] Reproduce Surface Transcripts 1 and 2 through the built binary on temporary homes.
- [ ] Record the outside evidence, the single unreachable row, the glossary check and the scope audit, and write the QA report.

## Acceptance Criteria

- [ ] A Lane Environment Skip reaches the item, the status line and the Pull Request body, and doctor reports the lane's Docker and image only while the lane is on.
- [ ] A ceiling-only refusal waits, starts when a slot frees, and parks as today at the earlier deadline, with Surface Transcript 2 reproduced.
- [ ] Each removed guard makes its named test fail, and the store stays the ceiling's only authority.
- [ ] The real-Docker doctor run is one declared unreachable row, and the scope audit finds no undeclared Governed Path change.

## Verification

The following command is rendered unchanged from Roundfix's derived QA contract.
It proves a newest report exists and its verdict is readable; eligibility is
applied in process by whoever settles the Task.

- `newest="$(find 'docs/specs/0291-the-queue-shows-a-skipped-linux-lane-and-waits-for-a-run-slot/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; Core Features 1-5; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Acceptance evidence; Unreachable Acceptance; Glossary
- `_techspec.md` → API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; API Contract 7; API Contract 8; API Contract 9; Surface Transcript 1; Surface Transcript 2; Invariants 1-7; Current behavior; Build Order 5
- ADR-0303; ADR-0304; ADR-0288; ADR-0294; ADR-0199; ADR-0104; ADR-0240; ADR-0276
