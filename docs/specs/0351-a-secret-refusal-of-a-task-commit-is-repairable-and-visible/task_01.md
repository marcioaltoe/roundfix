---
task: task_01
spec: 0351-a-secret-refusal-of-a-task-commit-is-repairable-and-visible
status: pending
type: backend
complexity: high
---

# Task 01: The Task commit's Secret Scan runs as the last Settlement Check of every attempt

## Overview

Today the Daemon scans the lines a Task commit would add only after the
Task's last Verification attempt (`executeTask`, `internal/daemon/task_engine.go`),
so a Secret Match fails the Task with no repair turn. After this Task the scan
is the check `settlement check: secret scan`, the last check of every
Verification attempt of every non-QA Task, with or without a QA gate. A
first-attempt match returns to the Agent as Verification Feedback with path,
line, kind and a `Fix:` line, never the value; the second attempt scans again;
a final match keeps today's reason bytes; a scan error still halts the Run.
This answers the report reported by fiscus on 2026-11-04, "recusa do scanner
de segredos no commit da Task não aparece em runs causes nem no park"
([references](references/2026-11-04-fiscus-a-task-commit-secret-refusal-is-invisible.md)).

## Requirements

1. MUST, before changing code, characterize the starting tree: run
   `go test -count=1 -timeout 10m -v -run '^TestTaskCommitWithASecretMatchSettlesFailedWithoutACommit$' ./internal/daemon`
   and record in the Result that it passes with one `secret_scan` event, no
   repair prompt and the reason `Task commit carries a secret: fixture.txt:2:
   credential-assignment`.
2. MUST add `const SettlementCheckSecretScan = "settlement check: secret scan"`
   to `internal/runevent/event.go`, and make `ProjectStreamEvent`
   (`internal/runevent/stream.go`) disclose that command as it discloses the
   two existing labels (`_techspec.md` → API Contract 5); add the label to the
   table of `TestASettlementCheckLabelIsProjectedAndARepositoryCommandIsNot`
   and change its invariant comment from "the two exact Roundfix labels" to
   "the three exact Roundfix labels", keeping every existing case.
3. MUST create `internal/daemon/settlement_secret_scan.go` with
   `secretScanSettlementCheck` exactly as `_techspec.md` → API Contract 1
   states, including the event payload `phase: secret_scan`, `attempt`,
   `matches`, the progress line `Task <id> commit refused: <n> Secret
   Match(es)` and the two-line diagnostic whose second line is, verbatim:
   `Fix: replace each named value with one no Secret Pattern matches, such as letters only or fewer than 20 characters; copying a value the repository already tracks is refused the same way.`
4. MUST append that check in `verifyTask` after the two in-process checks and
   outside the `plan.settlementChecks` condition (API Contract 2), add
   `HaltOnError` to `verificationCheck` and honor it in
   `runVerificationAttempt` (API Contract 3), render the scan's final failure
   as its first diagnostic line in `verificationTerminalReason` (API Contract
   4), and remove the post-attempt `scanTaskCommitPaths` block from
   `executeTask`. `scanTaskCommitPaths`, `prepareTaskCommit` and the
   `secretscan` package keep their bytes.
5. MUST add to `internal/daemon/task_secret_scan_test.go`, each calling
   `t.Parallel()` first and building values at runtime as the file already
   does:
   - `TestTaskCommitSecretMatchReturnsAsVerificationFeedback`: the fake Agent
     writes `SERVICE_API_KEY = <strings.Repeat("K8", 15)>` on its first call
     and `SERVICE_API_KEY = <strings.Repeat("z", 30)>` on its second; assert
     the Task settles `completed` and its commit holds the letters-only line,
     exactly one repair prompt was sent and it contains
     `settlement check: secret scan`, `fixture.txt:2: credential-assignment`
     and the `Fix:` line of requirement 3, exactly one `secret_scan` event
     exists with `attempt` 1, and no eight consecutive bytes of the value
     occur in any prompt, any `*-settlement-secret-scan.log` file under the
     Run's artifacts, any event summary or payload, or the progress output.
   - `TestTaskCommitSecretScanRunsAfterTheSettlementChecks`: subtest
     `gated graph`, built with `settlementAuditFixture` and an empty
     `fakeSettlementChecker`, reads the `started` verification events of
     attempt 1 and asserts the check labels end with `settlement check: spec
     consistency`, `settlement check: authorization`, `settlement check:
     secret scan` in that order; subtest `graph without a QA gate`, built
     with `newTaskCycleFixture`, asserts the only check label after the
     commands is `settlement check: secret scan`.
   - `TestTaskSecretScanReadErrorHaltsTheRun`: the fake Agent writes an
     untracked file and sets its mode to `0o000` (restored in `t.Cleanup`);
     `TaskCycle` returns an error, no Verification Feedback prompt is sent,
     no Task commit exists and the Task is not settled `completed`. It skips
     with a reason on Windows and when `os.Geteuid() == 0`.
6. MUST change, and only as `_techspec.md` → Declared breaks states:
   - `TestTaskCommitWithASecretMatchSettlesFailedWithoutACommit`: two
     refusals with `attempt` 1 and 2, one repair prompt, the same reason
     bytes, and the value search extended to the prompts and the
     `*-settlement-secret-scan.log` files;
   - `TestTaskCycleVerificationSequenceStopsAtFirstFailureWithOneVerdictPerAttempt`:
     the expected phases become
     `waiting|started|command-passed|started|failed|started|command-passed|verdict|waiting|started|command-passed|started|failed|started|command-passed|verdict`;
   - `TestTaskCycleExecutesAgentVerifySettleCommitContract`: the snapshot
     script becomes `{nil, {"src/one.go"}, {"src/one.go"}, nil, {"src/two.go"}}`
     with its comment naming the scan's preparation;
   - `TestRunImplementTemporaryVerificationFlowRepeatedTemporaryPreservesTaskWorktree`
     (`internal/cli/implement_test.go`): two retry `started` events, the
     command's and the scan's, with the failure message naming both.
   Every other existing test, helper and fixture stays byte-identical.
7. MUST keep every test file compiling for Windows
   (`GOOS=windows go vet ./internal/daemon ./internal/runevent`) and run each
   sabotage `_techspec.md` → Testing Approach lists for this Task on a
   temporary copy, recording in the Result which named test fails on which
   case.
8. MUST NOT touch `internal/secretscan`, `internal/runcause`,
   `internal/delivery`, any other `internal/cli` file, the QA Report scan,
   any document or skill, the Makefile, CI workflows, `go.mod` or `go.sum`.

## Subtasks

- [ ] Characterize today's refusal.
- [ ] Add the label constant and its stream projection.
- [ ] Move the scan into the attempt as a halting check with today's final reason.
- [ ] Add the feedback, order and read-error tests; adjust the four declared tests.
- [ ] Run the sabotages on a temporary copy.

## Acceptance Criteria

- [ ] A first-attempt Secret Match reaches the Agent as Verification Feedback with location and `Fix:` line and no value, and a repaired Task completes.
- [ ] A final Secret Match settles `failed` with today's reason and no commit, after exactly one repair turn; the scan runs last in gated and ungated graphs; a scan error halts the Run.

## Context

- interface: `internal/runevent/event.go`
- interface: `internal/runevent/stream.go`
- interface: `internal/runevent/settlement_label_projection_test.go`
- creates: `internal/daemon/settlement_secret_scan.go`
- interface: `internal/daemon/engine.go`
- interface: `internal/daemon/task_engine.go`
- interface: `internal/daemon/task_secret_scan_test.go`
- interface: `internal/daemon/task_engine_test.go`
- interface: `internal/daemon/settlement_checks_test.go`
- interface: `internal/daemon/settlement_spec_consistency_test.go`
- interface: `internal/cli/implement_test.go`
- instruction: `internal/secretscan/scan.go`

## Consumers

- The attempt's check list, the `secret_scan` event and the final Task reason — consumers: `./internal/daemon`, `./internal/cli`, `./internal/runevent`

## Verification

- `out="$(go test -count=1 -timeout 20m -v -run '^(TestTaskCommitSecretMatchReturnsAsVerificationFeedback|TestTaskCommitWithASecretMatchSettlesFailedWithoutACommit|TestTaskCommitSecretScanRunsAfterTheSettlementChecks|TestTaskSecretScanReadErrorHaltsTheRun|TestTaskCommitScansOnlyTheLinesItAdds|TestTaskCommitAcceptsALettersOnlyFixture|TestTaskSecretScanUsesHeadAndPreparedPaths|TestTaskCycleVerificationSequenceStopsAtFirstFailureWithOneVerdictPerAttempt|TestTaskCycleExecutesAgentVerifySettleCommitContract|TestSettlement.*|TestSpecConsistency.*|TestGatelessGraphNeverCallsTheSettlementChecker)$' ./internal/daemon 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestTaskCommitSecretMatchReturnsAsVerificationFeedback TestTaskCommitSecretScanRunsAfterTheSettlementChecks TestTaskSecretScanReadErrorHaltsTheRun; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { echo "missing pass: $name" >&2; exit 1; }; done` — expected: exit 0. Before this Task none of the three new tests exists, so the first pass line is missing; afterwards the scan returns as feedback, runs last and halts on a read error, and the existing scan and Settlement Check tests pass.
- `grep -q 'SettlementCheckSecretScan' internal/runevent/event.go && go test -count=1 -timeout 20m -run '^TestASettlementCheckLabelIsProjectedAndARepositoryCommandIsNot$' ./internal/runevent && go test -count=1 -timeout 20m -run '^TestRunImplementTemporaryVerificationFlow' ./internal/cli` — expected: exit 0. Before this Task the constant does not exist, so the first check fails; afterwards the label is projected and the temporary-retry flows of the Implement Command pass with the scan's events.

## References

- `_prd.md` → Sources 1; Goals; Core Features 1-3; Success Metric 1; Success Metric 2; Success Metric 3
- `_techspec.md` → Current behavior 1-4; API Contract 1; API Contract 2; API Contract 3; API Contract 4; API Contract 5; Invariants 1-5; Declared breaks; Testing Approach; Build Order 1
