---
task: task_01
spec: 0259-proven-signals-and-a-shared-build-cache
status: pending
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
