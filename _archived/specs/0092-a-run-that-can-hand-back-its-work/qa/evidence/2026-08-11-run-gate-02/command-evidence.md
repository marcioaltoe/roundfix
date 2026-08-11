# Command evidence — QA rerun 02

Build: `32ccc445c9591a59687b2fe2d90e737402f49974`

## Authoring precondition

`rtk roundfix spec check 0092-a-run-that-can-hand-back-its-work --strict`
exited 0 with `No findings.`

## Cold static gate

`rtk env GOCACHE="$PWD/.gocache" go clean -testcache` exited 0 before each
attempt.

The first `rtk make verify` attempt reached the repository integration test and
was denied network access to `api.github.com` by the sandbox. The exact gate was
rerun after another cache clear with that integration boundary authorized.

The authorized cold `rtk make verify` exited 2. Every listed package was
uncached. The only failed case was:

```text
--- FAIL: TestRunCommandHelp (0.00s)
    --- FAIL: TestRunCommandHelp/reconcile (0.00s)
        cli_test.go:458: expected help output to contain
        "roundfix reconcile [run-id] [--apply] [--format <text|json>]",
        got "roundfix reconcile [run-id] [--apply | --discard-superseded | --carry-forward] [--format <text|json>]"
FAIL roundfix/internal/cli
make: *** [test] Error 1
```

This is a code-caused exact-help-contract regression introduced by the Task 10
documentation change, not an environment failure.

The focused uncached reproduction also exited 1:

`rtk env GOCACHE="$PWD/.gocache" go test ./internal/cli -run
'^TestRunCommandHelp$/^reconcile$' -count=1 -v`

## Built CLI and prior-finding retest

`rtk go build -buildvcs=false -o /tmp/roundfix-qa-0092 ./cmd/roundfix` exited
0.

- `rtk /tmp/roundfix-qa-0092 --help` exited 0 and included the reconcile
  synopsis with `--discard-superseded` and `--carry-forward`.
- `rtk /tmp/roundfix-qa-0092 reconcile --help` exited 0 and listed both flags
  with the descriptions `Discard Run Branches proven superseded` and `Hand a
  stopped Run's settled Tasks back to the checkout`.
- Focused source and glossary reads found `Selection Failure` at
  `CONTEXT.md:199-201`, `Branch Disposition` at `CONTEXT.md:239-241`, emitted
  `agent_selection_failed`, and both accepted reconcile flag names.
- `rtk /tmp/roundfix-qa-0092 reconcile --format json` exited 0 in read-only
  mode with `mode: "dry-run"`, `applied: 0`, no candidates, and
  `carryForwards: []`.

Prior F-02 and F-03 are fixed on this build. F-04 remains the exact help-test
regression described above.

## Selection, settlement, and Run outcome

The following current-build commands exited 0:

- `rtk env GOCACHE="$PWD/.gocache" go test ./internal/agent -run
  '^(TestWorkStartedBoundary.*|TestMarkBatchFailedKeepsAlreadySettledIssues)$'
  -count=1 -v` — three work-started/selection cases and six persisted-status
  Batch cases passed.
- `rtk env GOCACHE="$PWD/.gocache" go test ./internal/daemon -run
  '^(TestFallbackEligibility.*|TestRunOutcomeDerived.*|TestResolveCycleVerificationFailureFailsBatchAndContinues|TestResolveCycleContinuesToNextBatchAfterFailedBatch)$'
  -count=1 -v` — three fallback cases and four outcome/continuation cases
  passed.
- `rtk env GOCACHE="$PWD/.gocache" go test ./internal/cli -run
  '^(TestAgentSelectionProfilesMacro|TestRunResolveVerificationFailureDoesNotCommit|TestRunResolveAgentFailureMarksBatchFailed|TestRunResolveAgentFailureContinuesWithLaterBatches|TestRunResolveClosesAgentSessionForTerminalOutcomes|TestFailedVerificationJournalsFailureWithoutCommitEvents)$'
  -count=1 -v` — the public macro's no-output fallback, post-output guard, and
  mixed-profile persistence cases passed with the selected Run outcome cases.

The outside-evidence origin is the pre-Spec measurement in
`docs/handoffs/2026-08-08-two-quotas-out-and-a-gate-that-said-no.md:16-20`:
Codex was exhausted and each attempted Run ended in nineteen seconds with the
usage-limit text wrapped as a protocol error. Lines 167-171 record that the
Fallback Chain did not activate because work-started had already been emitted.

## Branch disposition and carry-forward

`rtk env GOCACHE="$PWD/.gocache" go test ./internal/cli ./internal/daemon
./internal/spec -run
'^(TestCarryForward.*|TestReconcileDiscard.*|TestReconcileWithoutTheFlag.*|TestSupersededBranch.*|TestWriteBranchDispositionRecord.*|TestRecordCarryForward.*|TestRunReconcileJSONMatchesTextFields)$'
-count=1 -v` exited 0.

The selected cases passed success, moved-input refusal, mixed-set atomicity,
read-only no-flag behavior, unreachable-commit refusal, record-write failure,
record persistence after discard, reachability classification, later-Run Task
coverage, Task provenance, and exact reconcile JSON fields.

## Commit scope and protected-tooling audit

For every commit from Task 01 through Task 10 and both earlier QA reports, the
gate ran `rtk git -c core.fsmonitor=false diff-tree --no-commit-id --name-only
-r <commit>`. Every command exited 0. The Task delivery commits resolved to:

| Task | Commit | Changed paths |
| --- | --- | --- |
| 01 | `3fc542bf` | `task_01.md`; characterization test |
| 02 | `94ff10b2` | `task_02.md`; bounded Agent/Daemon source and tests |
| 03 | `1152beb8` | `task_03.md`; bounded Agent source/tests and characterization test |
| 04 | `ecd8feb4` | `task_04.md`; bounded CLI/Daemon source and tests |
| 05 | `e1ecb02e` | `task_05.md`; bounded reconcile CLI/Daemon source and tests |
| 06 | `9ee875c3` | `task_06.md`; bounded reconcile/spec source and tests |
| 08 | `1a36a3d7` | `task_08.md`; bounded CLI/characterization tests |
| 09 | `b60deccf` | `task_09.md`; bounded exact JSON contract test |
| 10 | `32ccc445` | `task_10.md`; `internal/cli/cli.go`; `CONTEXT.md` |

The authoring and repair commits (`6cdf7298`, `15fda0d0`, `14f174da`,
`d1b5608b`, `7eff0fea`, `d55eab9d`, `79dd4d3a`, `7c5ae3a5`, and `59011d50`)
changed only their declared Spec artifacts or the derived coverage record. The
two QA commits changed only their report, evidence, and gate Task status. No
commit touched repository-tooling configuration, scripts, ignore files,
plugin declarations, or version pins, so no tooling-authorization chronology
row applies.

## Report contract

The report-specific checks exited 0 for the QA directory, evidence file,
terminal `verdict`, all three blocked-row keys, and both coined terms. Focused
row counting found 8 `pass`, 2 `fail`, and 7 `blocked (environment: ...)` rows.
A focused pending-row search exited 1 with no matches, as expected.

`rtk git -c core.fsmonitor=false diff --check` exited 0. The final status scan
listed only this report and this evidence file as untracked paths.
