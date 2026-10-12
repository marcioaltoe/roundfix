---
task: task_01
spec: 0259-proven-signals-and-a-shared-build-cache
status: completed
type: backend
complexity: high
---

# Task 01: Every signal goes to a proven process, and a rule keeps it so

## Overview

Seven places signal a recorded PID or process group without proving which
process holds it now (`_techspec.md` → Current behavior). This Task adds the
two Proven Signal operations to the store, makes Force Stop prove the owner
again before its forced kill, routes the production detach and QA-format group
kills and the three test senders through the operations, and adds the
repository rule that forbids every other signal path. It is verifiable on its
own: the rule passes on the repository and fails on seeded violations, and the
store's real-process tests prove that a reused PID receives nothing.

## Requirements

1. MUST answer the Backlog Entry of 2026-10-09,
   [a test's child process is killed by a signal nobody identified](references/2026-10-09-a-test-child-killed-by-an-unknown-signal.md),
   by fixing its class: every signal to a recorded PID verifies the process's
   identity first.
2. MUST add `SignalProvenProcess` and `SignalProvenProcessGroup` to the store
   for Unix, following `_techspec.md` → API Contract 1 and API Contract 2,
   reading start identity with the reader `OwnerProcessIdentity` already uses.
   A mismatched or absent target MUST receive nothing and MUST NOT be reported
   as an error.
3. MUST make `TerminateAndWait` prove the owner again with the same recorded
   identity before the forced signal, following API Contract 3. A mismatch
   after the graceful signal MUST return nil without the forced signal. The
   PID-only path for an empty recorded identity and the
   `--owner-identity-unreadable` flag MUST stay as they are.
4. MUST route through `SignalProvenProcessGroup`, with the leader's identity
   recorded right after the process starts: the detach group kill in
   production, the QA-format cancellation group kill in production, the
   store tree helper's cleanup, `killDetachFixtureGroup`, and the detach
   caller group kill in the implement tests. MUST route the detach survivor
   teardown, which today looks the PID up with `os.FindProcess`, through
   `SignalProvenProcess` with the identity recorded when the PID was read, or
   through the `*os.Process` of the child the test started.
5. MUST add `TestEverySignalGoesToAProvenProcess` following API Contract 4,
   with the seeded subtests `direct-kill`, `bare-pid-lookup`, `group-kill`,
   `shell-kill`, `zero-probe` and `allowlisted`. Each failing seed MUST report
   `<file>:<line>:` for the seeded line. The allowlist names functions, not
   files, and lists only the Proven Signal operations, the store's per-platform
   `signalOwnerProcess` and the self-interrupt in the run UI.
6. MUST add the store tests named in Verification, using real helper processes
   started by the test with an identity recorded at start, and the controller's
   injected readers for the re-proof test. A mismatch test MUST prove the
   helper is still alive after the call.
7. MUST keep Windows and Linux builds of the non-test code compiling; the
   Windows detach and QA-format files keep their behavior.
8. MUST add **Proven Signal** to `CONTEXT.md` through `domain-modeling`,
   linked to ADR-0266, with `Avoid` synonyms such as safe kill.
9. MUST add one sentence to the detached Runs user guide, where it shows the
   `killed` diagnostic, saying that Roundfix signals a detached child's
   process group only after proving that its leader is the child it started.
   MUST record a Coverage Review for the `command resolve`, `command watch`
   and `guide docs/user-guide/commands/detached-runs.md` entries of the Skill
   Coverage Map, whose sources this Task changes, each naming `Spec 0259 task_01`
   and stating that only the
   detached child's group kill now goes through a Proven Signal and that no
   flag, output or exit code changes, and MUST re-record the Behavior Surface
   Record with
   `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`.
10. MUST NOT change any test's assertion other than the teardown it converts,
    the Makefile, the CI workflows or the module files.

## Subtasks

- [ ] Add the two Proven Signal operations and their tests.
- [ ] Make Force Stop prove the owner again before its forced kill.
- [ ] Route the production detach and QA-format group kills through them.
- [ ] Convert the three test senders and the survivor teardown.
- [ ] Add the Signal rule with its seeded cases.
- [ ] Add the glossary term.

## Acceptance Criteria

- [ ] The rule passes on the repository and fails on each seeded violation with its file and line.
- [ ] A Proven Signal with a recorded identity that no longer matches leaves the process running.
- [ ] Force Stop sends no forced kill when the owner's identity changed during the grace.
- [ ] Windows and Linux builds of the non-test code pass.

## Context

- instruction: `docs/adr/0266-a-signal-goes-only-to-a-proven-process-and-worktrees-share-one-build-cache.md`
- instruction: `docs/adr/0213-a-test-fixture-process-ends-with-the-test-binary-that-started-it.md`
- instruction: `docs/adr/0044-orphaned-run-locks-are-reclaimed-on-proven-owner-death.md`
- instruction: `internal/testfixture/parallel_tests_test.go`
- creates: `internal/store/process_signal_unix.go`
- creates: `internal/store/process_signal_unix_test.go`
- creates: `internal/testfixture/signal_rule_test.go`
- interface: `internal/store/process.go`
- interface: `internal/store/process_unix_test.go`
- interface: `internal/cli/detach.go`
- interface: `internal/cli/detach_unix.go`
- interface: `internal/cli/detach_windows.go`
- interface: `internal/cli/detach_test.go`
- interface: `internal/cli/implement_detach_teardown_test.go`
- interface: `internal/cli/implement_test.go`
- interface: `internal/daemon/qa_format.go`
- interface: `internal/daemon/qa_format_unix.go`
- interface: `internal/daemon/qa_format_windows.go`
- instruction: `.agents/skills/roundfix/references/review-runs.md`
- interface: `docs/user-guide/commands/detached-runs.md`
- instruction: `docs/user-guide/commands/resolve.md`
- instruction: `docs/user-guide/commands/watch.md`
- interface: `docs/references/skill-coverage.json`
- interface: `docs/references/behavior-surfaces.json`
- interface: `CONTEXT.md`

## Verification

- `out="$(go test -count=1 -v -run '^TestEverySignalGoesToAProvenProcess$' ./internal/testfixture 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestEverySignalGoesToAProvenProcess TestEverySignalGoesToAProvenProcess/direct-kill TestEverySignalGoesToAProvenProcess/bare-pid-lookup TestEverySignalGoesToAProvenProcess/group-kill TestEverySignalGoesToAProvenProcess/shell-kill TestEverySignalGoesToAProvenProcess/zero-probe TestEverySignalGoesToAProvenProcess/allowlisted; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done` — expected: exit 0. Before this Task the rule does not exist, so no PASS line is printed and the command fails.
- `out="$(go test -count=1 -v -run '^(TestSignalProvenProcessSendsToTheRecordedProcess|TestSignalProvenProcessSendsNothingToAReusedPID|TestSignalProvenProcessRefusesAnUnprovableTarget|TestSignalProvenProcessGroupProvesItsLeader|TestSignalProvenProcessGroupSignalsMembersOfAGoneLeader|TestOwnerProcessControllerReprovesBeforeForceKill)$' ./internal/store 2>&1)" || { printf '%s\n' "$out"; exit 1; }; for name in TestSignalProvenProcessSendsToTheRecordedProcess TestSignalProvenProcessSendsNothingToAReusedPID TestSignalProvenProcessRefusesAnUnprovableTarget TestSignalProvenProcessGroupProvesItsLeader TestSignalProvenProcessGroupSignalsMembersOfAGoneLeader TestOwnerProcessControllerReprovesBeforeForceKill; do printf '%s\n' "$out" | grep -q -- "--- PASS: $name " || { printf 'missing passing test %s\n' "$name" >&2; exit 1; }; done; go test -count=1 -race -short ./internal/store || exit 1; GOOS=windows go build -buildvcs=false ./... || exit 1; GOOS=linux go build -buildvcs=false ./... || exit 1` — expected: exit 0. Before this Task the six tests do not exist, so the PASS check fails.
- `tr -s '[:space:]' ' ' < CONTEXT.md | grep -qF -- '**Proven Signal**:' || { printf 'CONTEXT.md does not define **Proven Signal**\n' >&2; exit 1; }; go test -count=1 -run 'Detach|ForceStop|Orphan' ./internal/cli || exit 1; go test -count=1 -run '^TestQAFormat' ./internal/daemon || exit 1; grep -q 'Spec 0259 task_01' docs/references/skill-coverage.json || { printf 'the Coverage Reviews of resolve, watch and the detached Runs guide are missing\n' >&2; exit 1; }; go test -count=1 -tags docscontract -run '^TestTheSkillCoverageMapIsCurrent$' ./internal/docscontract || exit 1` — expected: exit 0. Before this Task the glossary lacks the term, so the command fails; afterwards the converted detach, Force Stop and QA-format tests pass and the Skill Coverage Map is current.

## References

- `_prd.md` → Goals; User Story 2; Core Feature 1; Core Feature 2; Core Feature 8; Success Metric 2; Acceptance evidence; Skills
- `_techspec.md` → Current behavior; API Contract 1; API Contract 2; API Contract 3; API Contract 4; Invariants 1-2; Testing Approach; Build Order 1
- ADR-0266; ADR-0213; ADR-0044


## Result

Implemented the Task 01 slice for Daemon Verification. Task status, authored
Verification, Subtasks and Acceptance Criteria checkboxes remain Daemon-owned.

### Implementation

- Added `SignalProvenProcess` and `SignalProvenProcessGroup`. Invalid PID,
  self PID, blank identity and zero signal refuse with
  `ErrOwnerProcessIdentityUnproven`. Each operation reads the existing start
  identity immediately before sending; absence or mismatch returns no-send
  without an error. Live identity read failures send nothing and propagate.
- Group proof reads live membership after leader absence. Darwin uses
  `kern.proc.pgrp`; Linux reads the group field in `/proc/<pid>/stat`.
  Exiting/zombie-only groups handle `EPERM` as no-send. Other Unix platforms
  return `ErrOwnerProcessUnsupported`. The platform membership readers live
  in `internal/store/process_darwin.go`, `process_linux.go` and
  `process_unix_other.go` so Darwin APIs do not enter Linux builds; these
  changes implement only this Task's group operation.
- `TerminateAndWait` repeats `proveOwner` with the recorded identity after
  the grace and before escalation. A changed identity proves the original
  owner ended and returns nil without the forced kill. Initial refusal,
  unreadable identity errors, empty-identity PID-only behavior and the CLI's
  explicit unreadable-identity override retain their existing contracts.
- Detach records the child identity immediately after Start and carries it
  through handshake teardown. QA-format splits Run into Start/Wait and
  synchronizes cancellation with the post-Start identity publication.
  Windows retains child-only detach teardown and default formatter
  cancellation.
- Converted the store tree cleanup, fixture group sender and implement
  caller group sender to group proof with recorded identities. Survivor
  teardown uses its owned child handle when available; a PID learned from
  the record uses a Proven Signal followed by `Wait4` to reap the child the
  test binary started. Existing behavior assertions were retained.
- Added an AST rule over every Go file, including files for other platforms,
  with the required artifact-directory exclusions. Its exceptions name
  functions in the store and the run UI, never entire files. Import aliases
  are resolved. Shell command strings and split command/argument calls are
  checked, while zero probes and diagnostic prose remain permitted.
  The direct-kill seed also has an allowlisted neighbor in the same file.
- Added the six required store regression tests with real test-owned helpers
  and identities captured after Start. A token mismatch models PID reuse;
  an input/output exchange proves the helper still runs after the refusal.
  The controller test injects identity/liveness readers and uses a real
  child kill in its forced branch, checking both the absence of that branch
  and continued helper liveness after changed/unreadable proof.
- Added Proven Signal to CONTEXT with its ADR-0266 link and Avoid terms.
  Added the detached Runs guide sentence and dated Coverage Reviews for
  resolve, watch and the guide, each naming Spec 0259 task_01. The Behavior
  Surface Record was generated with the required record command.

### Focused checks and acceptance evidence

All commands below were invoked through `rtk proxy`; final runs exited 0.
These are focused implementation checks, not the authored Verification.

| Acceptance criterion | Evidence |
| --- | --- |
| Rule accepts the repository and rejects seeded violations with file and line | `go test -race -count=1 -v ./internal/testfixture -run 'TestEverySignalGoesToAProvenProcess/\|TestEveryTestInAParallelTestPackageRunsInParallel/repository$'`: explicit PASS lines for the repository and all six required seeds. Each failing seed asserts its diagnostic at the seeded file and line; the shell seed covers kill, pkill and killall. The Parallel Test Package repository check also passed. |
| Mismatched identity leaves the process running | `go test -race -count=1 -timeout=2m ./internal/store -run 'SignalProven\|ReprovesBeforeForceKill\|TestOwnerProcessController(Terminate\|Force\|ProveOwner)'`: passed. Process and group mismatch tests exchange data with the unchanged live helpers; matching signals stop the helpers; group proof also stops a member after its leader is reaped and declines an empty group. |
| Force Stop sends no forced kill after identity changes during grace | The same store selection passed `TestOwnerProcessControllerReprovesBeforeForceKill`, covering changed, absent, unreadable, matching and legacy identities. Changed proof returns nil, the forced branch is not entered, and the real helper remains responsive. |
| Windows and Linux non-test code compiles | `GOOS=windows go build -buildvcs=false ./internal/store ./internal/cli ./internal/daemon` and `GOOS=linux go build -buildvcs=false ./internal/store ./internal/cli ./internal/daemon`: passed for all affected non-test packages. Whole-module builds remain for Daemon Verification. |

Additional final focused evidence:

- `go test -race -count=1 ./internal/cli -run 'TestRunDetachedCommand|TestDetachedChildIsTerminatedAtTeardown|TestRunImplementDetachSurvivesCallerProcessGroupKill|TestImplementDetachChildEndsWhenItsTestBinaryDies|TestDetachSurvivorEndsWhenItsTestBinaryDies|TestDetachFixtureGroupWithOnlyAnUnreapedMemberHasEnded'`: passed. This covers both handshake timeout phases, existing diagnostics/exit codes, survivor reaping and caller/group teardown.
- `go test -race -count=1 ./internal/daemon -run 'TestQAFormat(CancellationRestoresFiles|StartFailureRestoresFiles|CommitsTheFormattedReport|RevertsAFailingFormatter)$'`: passed, including cancellation restoration and Start failure.
- `go test -count=1 -tags docscontract ./internal/docscontract -run '^TestTheSkillCoverageMapIsCurrent$' -record-skill-coverage`: passed after the final guide edit. This is the source record operation expressly required by Requirement 9, not the Daemon's non-recording Verification command. Only the detached Runs guide fingerprint changed.
- `git diff --check`: passed before this Result entry; final scope and authored-Task projection are checked again after the entry.

### Corrections during implementation

The starting inspection found the unproven signal paths and no Proven Signal
operations or repository rule. Early focused runs exposed scanner false
positives on its own seed/comparison strings and loss of released-child
reaping in survivor teardown; both were corrected without file exemptions or
weaker behavior assertions. The controller harness initially used zero grace,
which this controller treats as waiting until its stop deadline. Its later
real-helper version omitted the actual mocked force kill, exposing duplicate
Wait calls in failure cleanup. It now uses a bounded grace, actually kills
its owned child in the force branch, and gives Wait one owner, including on
failure. An early long-running check overlapped source edits and reported
suiteguard boundary changes; final focused runs used stable sources and
exited 0, including race checks.

### Execution boundary and follow-up

The sandbox could not start a command because the configured writable root
contains a symlink. Approved escalated local execution was used in the same
Task Worktree. No Makefile, CI workflow, module file, Task Graph manifest or
other Task file was changed. No commit, push or pull request was made.
Authored Verification and repository Verification were not run.

Task 02's future `killResidue` sender must extend the function-scoped Signal
rule only together with its environment-mark re-proof, as API Contract 4
requires. Helper Lifeline/Residue and build-cache work stay in their assigned
Tasks. This Result does not settle Task status or claim a full-suite/F6 verdict.


### Verification Feedback — attempt 1

The Daemon's configured `make verify-changed` exited 2. Its diagnostic is
`/Users/marcio/.roundfix/artifacts/339f8dac2b687a04/runs/run_20261012T221735Z_f00baa2097cf9b98/verification/batch-001-attempt-1.log`.
The only test failure recorded there is
`TestBudgetExceededRunRecordsItsReason`, whose expectation of one completed
Task and one commit was unmet. No diagnostic log body is embedded here.
The artifact records successful store and testfixture package checks; those
are Daemon-provided facts, not additional Agent reruns.

Investigation of `internal/daemon/task_engine_test.go` found a fixed 500 ms
wall-clock Run Budget surrounding real Git/setup/settlement work for the
first fixture Task. Its assertion assumes that work finishes before expiry
and that the second Task then blocks. `fixture.plan()` configures no
formatter and uses a fake Agent; `formatQADirectory` returns before any
process setup when FormatCommand is empty. The budget watchdog is unchanged
by this Task and expires from the real wall-clock deadline. None of the
changed Proven Signal operations is invoked by this fixture.

Focused diagnostic evidence:

- `rtk proxy go test -race -count=1 ./internal/daemon -run '^TestBudgetExceededRunRecordsItsReason$'`: exited 0 on one isolated run. This is not a fix or evidence that the configured suite failure is resolved.
- Prepared `/private/tmp/roundfix-task01-budget-baseline-hop9rk1x/overlay.json` without repository edits. The overlay restored every changed production Go file to HEAD, removed the new signal operations from that diagnostic compilation, and used the committed budget test with its assertions unchanged. Its only test perturbation delayed the first fake Task by 501 ms through the existing `afterTask` seam, modeling work that exceeds the fixture's budget.
- `rtk proxy go test -race -count=1 -overlay=/private/tmp/roundfix-task01-budget-baseline-hop9rk1x/overlay.json ./internal/daemon -run '^TestBudgetExceededRunRecordsItsReason$'`: exited 1 with the same unmet completion/commit expectation. This deliberately failing baseline diagnostic demonstrates the existing timing dependency; it does not establish the exact latency source in the Daemon's earlier run.
- Diff inspection confirmed the diagnostic overlay introduced no repository source/test changes. The feedback turn changes only this Result entry; Task status and the authored Task projection remain unchanged.

Follow-up outside Task 01: make the budget-reason fixture arrange expiration
after the first settlement through a deterministic test seam, preserving all
its current outcome, reason, commit and Task-state assertions. Increasing its
timeout or accepting zero completed Tasks would not repair that contract.
The assigned slice and Requirement 10 do not authorize a budget fixture or
scheduler repair, so `task_engine_test.go` and budget production code remain
unchanged. The configured gate failure remains unresolved and is not waived
by the isolated pass. No in-slice signal defect was identified in the feedback.

No declared Verification command, `make verify-changed`, or full suite was
rerun by the Agent. No Task status, other Task, manifest, commit, push or pull
request was changed. The Daemon retains Verification and settlement ownership.

## Recorded paths

The Daemon recorded these paths, which this Task changed without declaring them in `## Context`.

- `internal/store/process_darwin.go`
- `internal/store/process_linux.go`
- `internal/store/process_unix_other.go`

## Carry-forward provenance

- Source Run: `run_20261012T221735Z_f00baa2097cf9b98`
- Source commit: `3bdb6f58e36a403e605f6c63b62d6da68c373f99`
