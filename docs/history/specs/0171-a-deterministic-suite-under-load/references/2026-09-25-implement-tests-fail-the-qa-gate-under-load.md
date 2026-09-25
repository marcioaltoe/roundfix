---
type: fix
status: promoted
created: 2026-09-25
spec: 0171-a-deterministic-suite-under-load
reason: null
---

# Timing-dependent `internal/cli` implement tests fail QA gates under load

## Symptom

QA gates fail their repository-Verification precondition on `internal/cli`
tests that pass in isolation, so Specs that never touched `implement` need extra
Runs. On 2026-09-25 alone, Spec 0163's QA failed on
`TestRunImplementQueuedCancellationStartsNoChildAndKeepsResumableTasks`
("timed out waiting for 2 Agent starts; got 1", 91 s) and then on
`TestRunImplementBootstrapsEachConcurrentTaskWorktreeBeforeAgentWork`
("expected Clean exit, got 1"); Spec 0164's QA failed on the daemon test
`TestTaskCycleVerificationCapacityCancellationWhileQueuedStartsNoCommandOrSettlement`.
All three pass 3/3 in isolation on the Spec branch and on `main`. The same day,
`make verify` for a documentation-only branch failed on
`internal/worktree` `TestBootstrapSerializesAcrossSiblings`, which also passes
3/3 in isolation.

## Where

`internal/cli/implement_test.go` and `internal/daemon/task_engine_test.go`:
fixed wall-clock waits around concurrent Task worktrees, Agent starts and
cancellation.

## Expected

Waits bound to the test deadline and to observed events rather than fixed
durations, so the suite is deterministic on a loaded machine.

## Evidence

QA reports of Runs `run_20260925T153433Z_8603eb3e79157622`,
`run_20260925T155346Z_b801d0ee9b2a247a` and
`run_20260924T223414Z_54e6f85ae98b6550`; isolation reruns of 2026-09-25.
