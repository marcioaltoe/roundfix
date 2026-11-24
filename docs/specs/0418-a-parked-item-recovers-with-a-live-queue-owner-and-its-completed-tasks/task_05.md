---
task: task_05
spec: 0418-a-parked-item-recovers-with-a-live-queue-owner-and-its-completed-tasks
status: completed
type: qa
complexity: high
---

# Task 05: Run the final QA gate

## Overview

This is the authored terminal gate for this Spec. It declares what the matrix covers and settles the Spec on evidence. The Spec makes a Delivery Queue item record the Run that just ended, a Failed Run included; lets a Delivery Retry carry the completed Tasks of a QA-Only Failed Run; ends a Delivery Retry with a Queue Owner proven alive or parks the item again and fails with exit `2`; and makes `roundfix deliver status` print an `Owner:` line for a queue nobody holds and one `Run:` line per running item. It records the rules in ADR-0450, amends ADR-0279 and ADR-0395, revises the glossary, and rewrites the `deliver` guide and the Roundfix skill's `deliver` reference. It answers the findings reported by quaestor on 2026-11-23 (ai-memory messages 01a20714-3a6e and 01a20714-ed79) and the item reported by the Praetor on 2026-11-23. No command, flag, exit code, help text, Run Database table, column or stage is added.

Every row runs on this repository's own files, temporary homes, temporary Git repositories, temporary Run Databases and temporary copies of the audited head or of the Delivery Base. No row:

- starts an Agent harness, an Agent session, an Author Subagent or an ACP adapter, or asks a person a question;
- starts a real Run, a Daemon or a detached Delivery Owner outside a test's own temporary home;
- calls ai-memory, the Indexus or any MCP server, runs the real `gh` or reaches the network;
- reads another repository, quaestor's included;
- reads or writes the real `~/.roundfix`, or reads the real home directory except for the toolchain locations below;
- signals a process it did not start, or ends a process by name;
- sends a message or a request to a provider, the judge or a reviewer;
- reads a credential or asserts a wall-clock duration.

Every `go test` it runs carries an explicit `-timeout`.

The gate follows the `qa-gate` section "A fresh home keeps the host toolchains". Every `go test` package or test run, every `make` target and the repository Verification keep the session's own `HOME` and `TMPDIR`. Only a row that runs the built product uses a fresh `HOME`. Every such row runs `bin/roundfix`, the binary `make build` made from the audited head, or the Delivery Base binary named in requirement 8, and never a `roundfix` found on `PATH`, and prepares its environment the same way:

- before `HOME` changes it takes `m="$(go env GOMODCACHE)"`, `c="$(go env GOCACHE)"`, `r="${RUSTUP_HOME:-$HOME/.rustup}"` and `k="${CARGO_HOME:-$HOME/.cargo}"`;
- it sets `HOME` and `TMPDIR` to fresh temporary directories outside the checkout, with `GOMODCACHE="$m"`, `GOCACHE="$c"`, `RUSTUP_HOME="$r"` and `CARGO_HOME="$k"`, which the command only reads, and `PATH` otherwise unchanged;
- it puts first on `PATH` a temporary directory holding an executable `gh` that writes nothing and exits 1;
- it unsets `GH_TOKEN`, `GITHUB_TOKEN`, `TYPESAFE_API_KEY`, `CLAUDE_CODE_SESSION_ID` and every Judge Key Variable for that command;
- it reads standard input from `/dev/null` and removes the temporary directories afterwards, after `chmod -R u+w` on them.

A command that fails only because a fresh `HOME` hid a toolchain is a gate error, never an environment block: the gate runs it again with the toolchain locations and records that result. A temporary Git repository for a product row is made outside the checkout with `git init` and one commit; a temporary copy of the audited head or of the Delivery Base is made with `git archive` of that commit and committed in a fresh Git repository outside the checkout, because the suite guard refuses a tree outside Git. The QA never writes into the checkout while a test suite runs. It keeps its evidence under names the repository does not ignore, `.txt` for command output, and links only files it kept. Evidence records environment variable names, never their values.

## Requirements

1. MUST build the binary with `make build` and run the repository Verification, recording its result as a gate fact.
2. MUST verify Success Metric 1 and Core Feature 1 by running, with `-v`, `TestTheItemRecordsTheRunOfAFailedImplementExecutor` over `./internal/delivery` and `TestAFailedImplementChildNamesTheRunItCreated` over `./internal/cli`, recording each subtest by name, and record from the test that the park keeps the blocker `delivery-error: run Implement executor: roundfix implement failed with exit code 1` quaestor reported.
3. MUST verify Success Metric 2 and Core Features 2 and 3 by running, with `-v` over `./internal/cli`, `TestADeliveryRetryCarriesTheTasksOfAQAOnlyFailedRun`, recording its six subtests by name, the two controls of Success Metric 2 (`another_task_not_completed_carries_nothing` and `unresolved_run_carries_as_before`) among them; then `TestCarryForwardAcceptsBudgetExceededRun` with `-v` over `./internal/worktree`, recording that the accepted-outcome list keeps its three outcomes (invariant 3), and the existing reconcile and Preflight carry-forward tests of `./internal/cli` by name.
4. MUST verify Success Metric 3 and Core Feature 4 by running, with `-v` over `./internal/cli`, `TestDeliverRetryReproducesSurfaceTranscriptAfterAQAOnlyFailedRun`, `TestDeliverRetryReproducesSurfaceTranscriptOwnerDidNotStart`, `TestDeliverRetryReproducesSurfaceTranscriptOwnerEndedBeforeTheProof` and `TestDeliverRetryProvesTheQueueOwnerItStarted`, recording each test and subtest by name with the Surface Transcript it reproduces: Surface Transcript 4, Surface Transcript 5 and Surface Transcript 6, and `TestRevertDeliveryQueueItemRetryRestoresTheParkedItem` with `-v` over `./internal/store`; then the existing retry tests of `./internal/cli`, `TestDeliverRetryWithAnOpenPullRequestAndNoMergeEvidenceRetriesAsBefore` among them, recording that the successful retry transcripts keep their bytes.
5. MUST verify Success Metrics 4 and 5 and Core Features 5 and 6 by running, with `-v` over `./internal/cli`, `TestDeliverStatusReproducesSurfaceTranscriptNoRecordedOwner`, `TestDeliverStatusReproducesSurfaceTranscriptTasksInProgress`, `TestDeliverStatusReproducesSurfaceTranscriptNoTaskInProgress` and `TestDeliverStatusPrintsTheOwnerLineOfAQueueWithNoRecordedOwner`, recording each test and subtest by name with the Surface Transcript it reproduces: Surface Transcript 1, Surface Transcript 2 and Surface Transcript 3; then the ten existing status tests task_03 names and the other owner status tests of `internal/cli/deliver_owner_status_test.go`, recording that the line of a recorded owner proven gone, the uncheckable-owner diagnostic and the silence for a proven live owner keep their bytes.
6. MUST verify Success Metric 7, the outside-evidence row, with its date and not blocked: run `TestTheBuiltBinaryReportsTheQueueQuaestorMeasured` with `-v` over `./internal/cli`, recording both subtests by name; and compare, looking nothing up, the three lines quaestor captured at 03:12 BRT of 2026-11-23 in `references/sources.md`, section 2, with Surface Transcript 1: the item row and the `Limits:` line keep their shape, and the `Owner:` line naming `roundfix deliver resume` and the line `Run: <slug> no Active Run`, which quaestor's capture lacks, now appear between them. Quote the passage of the source the comparison rests on.
7. MUST verify Success Metric 6 and Core Features 8 and 9 by running, with `-tags docscontract` and `-v` over `./internal/docscontract`, `TestTheDeliveryRecoveryDecisionAmendsTheOwnerStatusRules`, `TestTheDeliverGuideAndReferenceDescribeTheRecoveryLines` and `TestTheSkillCoverageMapIsCurrent`, recording each test and subtest by name, then `TestEveryOwnedSkillVersionIsRecorded`, `TestAuthorialSkillSync` and `TestSettlementGuidanceIsOneTable` with `-v` over `./skills`; and record from the files: ADR-0450's header with `status: accepted` and `knowledge: true`; ADR-0279's and ADR-0395's `amended_by` holding ADR-0450 and each dated notice after its sentence; the path of ADR-0450 once in the Standard Export List, in ascending order; that ADR-0450 names no client and cites no Spec folder path; and that both copies of the Roundfix skill are byte-identical with `cmp`.
8. MUST verify invariants 1, 2 and 8 through the built binaries, in a fresh home prepared as above: build a Delivery Base binary from a `git archive` of the Delivery Base, the merge base of the audited head with the default branch, in a temporary copy outside the checkout; for each of `deliver`, `deliver plan`, `deliver start`, `deliver status`, `deliver resume`, `deliver retry`, `deliver stop` and `deliver lane`, record that `bin/roundfix <command> --help` and the Delivery Base binary's print the same bytes and exit code; and in a temporary Git repository with one commit and no Delivery Queue, record that `bin/roundfix deliver status` and `bin/roundfix deliver retry <slug>` print the same bytes and exit codes as the Delivery Base binary.
9. MUST verify that the regenerated files are the commands' output by running, on a temporary copy of the audited head, `make skills-sync`, `go test ./skills -run '^TestEveryOwnedSkillVersionIsRecorded$' -record-skill-versions` and `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`, and recording that they change no file.
10. MUST run each sabotage below on its own temporary copy of the audited head. Each sabotage keeps the packages compiling and regenerates nothing, so the named test fails on its assertion; record that it fails on the named test or subtest.
    1. `runCandidate` no longer calls `recordRun` before it reads the executor's error: `TestTheItemRecordsTheRunOfAFailedImplementExecutor/failed_executor_parks_with_its_run` fails.
    2. `qaOnlyFailedRun` always reports that the Run is not one: `TestADeliveryRetryCarriesTheTasksOfAQAOnlyFailedRun/carries_the_settled_tasks` fails.
    3. The retry treats a queue with no live owner after the start as proven: `TestDeliverRetryReproducesSurfaceTranscriptOwnerEndedBeforeTheProof` fails.
    4. `RevertDeliveryQueueItemRetry` keeps the retried count: `TestRevertDeliveryQueueItemRetryRestoresTheParkedItem/reverts_a_retried_item` fails.
    5. `deliveryOwnerLine` prints nothing when no owner is recorded: `TestDeliverStatusReproducesSurfaceTranscriptNoRecordedOwner` fails.
    6. `deliveryRunLines` ignores `settled` events: `TestDeliverStatusReproducesSurfaceTranscriptNoTaskInProgress` fails.
    7. The notice paragraph is removed from ADR-0395: `TestTheDeliveryRecoveryDecisionAmendsTheOwnerStatusRules` fails.
    8. The phrase `no Active Run` is removed from `skills/roundfix/references/deliver.md` only: `TestTheDeliverGuideAndReferenceDescribeTheRecoveryLines/mirror_reference` fails.
11. MUST execute the tests task_01, task_02, task_03 and task_04 name in their Verification.
12. MUST perform the glossary check (`docs/agents/domain.md`): record that `CONTEXT.md` defines **QA-Only Failed Run** as a line-anchored Glossary Definition, that **Task Carry-Forward**, **Delivery Retry** and **Queue Owner** are revised as `_techspec.md` → API Contract 8 states and each of the four cites ADR-0450, that no other entry of `CONTEXT.md` differs from the Delivery Base, and that no term the Spec uses needs a definition the glossary lacks.
13. MUST verify from the repository history:
    - each Task's changed files stay within its declarations or its `## Recorded paths`;
    - every changed Governed Path is listed in `_authorization.md` `paths`, and the Governed Paths changed are only those four (the three skill files and `docs/references/coverage-record.json`, granted by the third decision, ai-memory message `01a20b31-20fa-7650-8642-d4ee3b5d39a8`);
    - the Task commit range changes none of these: the Makefile, the CI workflows, `go.mod`, `go.sum`, `.roundfixrc.yml`, `skills-lock.json`, `AGENTS.md`, `CLAUDE.md`, `CHANGELOG.md`, `internal/app/version.go`, `dist/npm/roundfix/package.json`, any file under `docs/agents/`, any Vendored Skill, any owned skill other than `roundfix`, and any reference of the Roundfix skill other than `deliver`;
    - `docs/references/coverage-record.json` is changed only by task_03's regeneration command, and the only test name it removes is `TestDeliverStatusPrintsNoOwnerLineWithoutARecordedOwner`; every other difference is an addition (correction of 2026-11-24: the QA gate of Run run_20261124T043440Z_39a347ec4e2091ba failed only on finding F-01, this criterion contradicting the amended grant);
    - in `docs/references/skill-coverage.json`, only the `review` value of the `command reconcile` entry differs from the Delivery Base, and it begins with the dated sentence of task_01;
    - the Run Database schema version constant in `internal/store/store.go` equals the Delivery Base's and no table, column or stage was added;
    - no accepted ADR other than ADR-0279 and ADR-0395 changed, both stay `accepted`, and ADR-0450 is `accepted`;
    - no `### QA settlement` section changed.

    This row reads Task commits, so it declares a `commit_range` input.
14. MUST record the non-waivable Pull Request row as `blocked (environment: no open Pull Request)`, naming the Pull Request row in its provenance. It supports the row with the pre-PR equivalent of each control:
    - approval from the recorded delivery authority, the maintainer's two decisions of 2026-11-23 quoted in `_authorization.md`;
    - checks from the Daemon's repository Verification at the audited head;
    - no unresolved threads, because no Pull Request exists;
    - no Merge-Ready acceptance yet;
    - the audited head named as the claimed candidate.
15. MUST, for the Outcome for the Archive Record, run only `m="$(go env GOMODCACHE)"; c="$(go env GOCACHE)"; h="$(mktemp -d)"; t="$(mktemp -d)"; HOME="$h" TMPDIR="$t" GOMODCACHE="$m" GOCACHE="$c" bin/roundfix archive 0418-a-parked-item-recovers-with-a-live-queue-owner-and-its-completed-tasks --plan --offline </dev/null; rm -rf "$h" "$t"` with every Judge Key Variable and `TYPESAFE_API_KEY` unset. It writes the report's Outcome as one paragraph that does not end with a colon. It MUST NOT do any of these:
    - accept a row whose only evidence is that a file was read;
    - run a `roundfix` found on `PATH`;
    - read or write `~/.roundfix`, or read the real home directory except for the toolchain locations above;
    - run the real `gh`, call ai-memory, the Indexus or any MCP server, or reach the network;
    - start an Agent harness, an Agent session, an Author Subagent, an ACP adapter, a real Run, a Daemon or a detached Delivery Owner outside a test's own temporary home;
    - send any message, or any request to a provider, the judge or a reviewer;
    - assert a wall-clock duration;
    - write into the checkout while a test suite runs;
    - link a file it did not keep, or keep evidence under a name the repository ignores.

## Subtasks

- [ ] Build the binary and run the repository Verification.
- [ ] Run the named tests of task_01 to task_04 and the existing tests they keep.
- [ ] Compare the built binaries' help and no-queue output with the Delivery Base binary's, and record the outside-evidence row.
- [ ] Prove the regenerated records are the commands' output on a temporary copy.
- [ ] Run each sabotage on a temporary copy.
- [ ] Record the glossary check, the scope audit and the Pull Request row, and write the QA report.

## Acceptance Criteria

- [ ] A Failed Run is recorded on its item, a QA-Only Failed Run's completed Tasks are carried by a Delivery Retry and never the gate's Task, and the Reconcile Command and the Preflight keep their three outcomes.
- [ ] A Delivery Retry exits `0` only with a proven live owner or nothing left to advance, and otherwise exits `2`, parks the item again with its earlier blocker and retry count, and prints the failure report.
- [ ] Status prints the new `Owner:` and `Run:` lines in the six Surface Transcripts' shapes, every existing line and the help of every `deliver` command keep their bytes, and the built binary reproduces the state quaestor measured with the new lines.
- [ ] ADR-0450, its two amendments, the glossary, the guide and both copies of the skill reference agree with the behavior; every sabotage makes its named test fail; the scope audit finds no change outside the declared files and the three granted Governed Paths.

## Verification

- `newest="$(find 'docs/specs/0418-a-parked-item-recovers-with-a-live-queue-owner-and-its-completed-tasks/qa' -type f -name "qa-report-*.md" -print 2>/dev/null | awk "{ report=\$0; name=\$0; parts=split(name, path, \"/\"); name=path[parts]; name=substr(name, 11, length(name)-13); date=substr(name, 1, 10); suffix=substr(name, 11); shape=(length(date) == 10 && substr(date, 5, 1) == \"-\" && substr(date, 8, 1) == \"-\"); for (i=1; i <= 10 && shape; i++) { if (i != 5 && i != 8 && index(\"0123456789\", substr(date, i, 1)) == 0) shape=0 } year=substr(date, 1, 4)+0; month=substr(date, 6, 2)+0; day=substr(date, 9, 2)+0; leap=(year%400 == 0 || (year%4 == 0 && year%100 != 0)); days=(month == 2 ? 28+leap : ((month == 4 || month == 6 || month == 9 || month == 11) ? 30 : 31)); dated=(shape && month >= 1 && month <= 12 && day >= 1 && day <= days); if (!dated) date=\"\"; if (suffix == \"\") { sequenced=1; sequence=-1 } else { sequenced=(length(suffix) > 1 && substr(suffix, 1, 1) == \"-\"); for (i=2; i <= length(suffix) && sequenced; i++) { if (index(\"0123456789\", substr(suffix, i, 1)) == 0) sequenced=0 } sequence=(sequenced ? substr(suffix, 2)+0 : 0) } if (dated && sequenced) printf \"%s\\t%s\\t%s\\t%s\\t%s\\n\", dated, date, sequenced, sequence, report }" | sort -k1,1n -k2,2 -k3,3n -k4,4n -k5,5 | tail -1 | cut -f5-)"; test -n "$newest" || exit 1; awk "BEGIN { whitespace=\" \t\r\n\f\v\" } NR == 1 && \$0 == \"---\" { frontmatter=1; next } frontmatter && \$0 == \"---\" { closed=1; exit } frontmatter && index(\$0, \"verdict:\") == 1 { verdict=substr(\$0, 9); while (length(verdict) > 0 && index(whitespace, substr(verdict, 1, 1)) > 0) verdict=substr(verdict, 2); while (length(verdict) > 0 && index(whitespace, substr(verdict, length(verdict), 1)) > 0) verdict=substr(verdict, 1, length(verdict)-1); verdicts++ } END { exit(closed && verdicts == 1 && length(verdict) > 0 ? 0 : 1) }" "$newest"`

## References

- `_prd.md` → Goals; User Stories 1-5; Core Features 1-9; Success Metric 1; Success Metric 2; Success Metric 3; Success Metric 4; Success Metric 5; Success Metric 6; Success Metric 7; Acceptance evidence; Unreachable Acceptance; Skills; Glossary
- `_techspec.md` → Interfaces; API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; API Contract 6; API Contract 7; API Contract 8; API Contract 9; Surface Transcript 1; Surface Transcript 2; Surface Transcript 3; Surface Transcript 4; Surface Transcript 5; Surface Transcript 6; Invariants; Coverage Map; Testing Approach; Risks & Considerations; Build Order 5
- `_authorization.md` → What the grant covers; Sanctioned regeneration; Limits
- [`references/2026-11-23-characterization.md`](references/2026-11-23-characterization.md), sections 1 to 8
- [`references/sources.md`](references/sources.md), sections 1, 2 and 3 (reported by quaestor and by the Praetor, 2026-11-23)
- ADR-0279, ADR-0395, ADR-0170, ADR-0373, ADR-0450
