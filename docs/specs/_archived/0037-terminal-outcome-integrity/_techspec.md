---
spec: 0037-terminal-outcome-integrity
prd: _prd.md
created: 2026-07-17
---

# Terminal outcome integrity — Technical Spec

## Executive Summary

Run completion becomes a guarded store operation: the first non-terminal-to-terminal transition wins, an identical replay is idempotent, and a conflicting terminal result returns a typed conflict without mutating the row or lock. Force Stop reverses its current unsafe order by proving owner-process exit before calling completion. Review Source waits gain the Run Database as a Stop Request source, and Agent Session cleanup queries the durable Agent Selection lifecycle instead of deriving an unconditional Agent Session name. The primary trade-off is availability for safety: a Force Stop that cannot prove owner exit leaves the Run Active and locked, requiring the user to resolve the owner process instead of risking two live writers.

## Project Constraints

- Identifier strategy: not applicable — all persistence and process operations
  reuse existing Run, scope, and Agent Session identities. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the implementation changes local
  process, database, and polling boundaries without changing authentication or
  HTTP behavior. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0022, ADR-0044, ADR-0051, and
  ADR-0052 govern Stop Request transport, owner-death proof, Work Item Agent
  Sessions, and terminal compare-and-set behavior. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — on 2026-07-26, the maintainer expressly
  authorizes changes to exactly `.agents/skills/roundfix/SKILL.md` and
  `skills/roundfix/SKILL.md`. On 2026-07-27, the maintainer additionally
  expressly authorizes the deterministic Skill-digest fallout of that edit in
  exactly `internal/baseline/assets/setups/typescript-bun.json`,
  `internal/baseline/testdata/catalog.digest`,
  `internal/baseline/testdata/catalog.normalized.json`,
  `internal/baseline/testdata/parity-corpus/v1/fixtures/asset-sync.json`, and
  `internal/baseline/testdata/parity-corpus/v1/manifest.json`. No other
  protected tooling mutation is authorized. Source:
  `docs/agents/agent-instructions.md`.

## System Architecture

- **`internal/store`** extends Run completion with compare-and-set semantics and a typed terminal conflict. It also exposes a guarded Integration Pending reconciliation operation and a query that returns Agent Selection scopes whose latest persisted status is active.
- **`internal/cli`** coordinates Force Stop in the order Agent Session cancellation → owner termination → terminal completion. Small platform files own process termination and exit proof, following the existing detached-process platform split.
- **`internal/watch`** accepts a Stop Request source and checks it at every wait boundary. It returns the existing stop classification so the CLI completes the Run as Stopped.
- **`internal/daemon`** keeps persisting Agent Selection `active` only after Agent Session preparation succeeds and `closed` after idempotent close. Those records become the cleanup source of truth.
- **Run Event Journal** receives one terminal outcome from the completion winner and secondary cleanup diagnostics only after the primary failure event.

No new package is required. The design reuses the Run Database, Agent Selection lifecycle, and existing OS-specific CLI files.

## Implementation Design

### Interfaces

```go
type TerminalOutcomeConflictError struct {
    RunID     string
    Stored    string
    Requested string
}

type CompleteRunResult struct {
    Run          Run
    Transitioned bool
}

func (s *Store) CompleteRun(ctx context.Context, id, outcome string) (CompleteRunResult, error)
func (s *Store) ReconcileIntegration(ctx context.Context, req IntegrationReconciliation) (Run, error)
```

`CompleteRun` is internal to Roundfix, so its result becomes explicit instead of making callers infer whether they won. Its behavior changes as follows:

1. Update the Run only when its stored state is non-terminal.
2. On zero affected rows, read the Run in the same transaction.
3. Return the Run with `Transitioned: false` when the stored and requested outcomes match.
4. Return `TerminalOutcomeConflictError` when the stored outcome differs.
5. Delete the Active Run lock and return `Transitioned: true` only when this call won the transition.

```go
type StopRequestSource interface {
    StopRequested(context.Context, string) (bool, error)
}

type OwnerProcessController interface {
    TerminateAndWait(context.Context, int) error
}
```

The watch dependency remains optional only in isolated unit tests. Operational watch Runs always supply the Store-backed source. `TerminateAndWait` sends the platform's graceful termination signal, waits within the existing bounded stop window, escalates through the platform's force-kill path, and returns success only after process absence is proven. Permission errors, unsupported termination, and an owner still alive at the deadline are failures.

### Data Models

No schema migration is required for the Agent Session registry: `run_agent_selections` already records Run ID, scope, runtime, selection role, fallback index, and lifecycle status. The latest record per `(run_id, scope_kind, scope_id)` determines whether cleanup is eligible. Agent Session names remain deterministic from scope and fallback index; work directories come from the Run and Task Worktree rules already used by the Daemon.

`IntegrationReconciliation` carries the Run ID, proven Run Branch head, target branch, target head, and timestamp. The store accepts it only when the current state is Integration Pending, changes the state to Clean, and appends a reconciliation Run Event containing both states and the evidence in the same transaction. Spec 0038 owns collection of that Git evidence.

### API Contracts

`roundfix stop --force <run-id>` keeps its command and success report. Its failure contract becomes:

- stdout remains empty;
- stderr names the Run ID, recorded owner PID, failed step, and states that the Run remains Active with its lock retained;
- exit code is the existing Run failure code;
- no Stopped outcome event or notification is emitted.

The command refuses to target an already terminal Run unless the requested Stopped outcome is already stored, in which case it reports the existing result idempotently and performs no process or Agent Session action.

Graceful Stop Requests preserve their current command response. During a Review Source wait, the owner checks the flag before status access and after each interruptible sleep. A detected request returns the existing stop sentinel; it does not run another fetch, check, commit, push, or Review Source mutation.

Agent Session cleanup rules:

- no active lifecycle record means no cancel or close call;
- each active scope is canceled and closed once in deterministic scope order;
- an adapter response stating that the named Agent Session was not found is treated as an idempotent close for a registered scope;
- other failures are secondary warnings and do not authorize terminal completion while the owner remains alive.

## Coverage Map

- Goal 1, Stories 1 and 3 → Store compare-and-set completion, conflict type, winner-only outcome publication.
- Goal 2, Story 1 → owner process controller and reordered Force Stop flow.
- Goal 3, Story 2 → Store-backed Stop Request source in every watch wait.
- Goal 4, Story 4 → active Agent Selection scope query and ordered secondary diagnostics.
- Story 5 → guarded Integration Pending reconciliation operation.

## Integration Points

- **Operating-system process control**: existing Unix process-group and Windows process termination patterns. The recorded PID is the only permitted owner target. A reused PID produces a conservative refusal because Roundfix cannot prove ownership.
- **acpx Agent Sessions**: existing cancel and close commands, selected from durable active lifecycle records.
- **Run Database**: compare-and-set state, lock release, selection lifecycle reads, and reconciliation event transaction.

## Testing Approach

- Store table tests cover non-terminal completion, identical replay, conflicting replay, lock retention on conflict, and the sole Integration Pending reconciliation transition.
- A deterministic concurrency test races Stopped against another terminal outcome and asserts one winner, one stable row, and one released lock only after the winning completion.
- CLI tests use a fake owner controller for graceful exit, force-kill exit, permission failure, and deadline failure. Failure cases assert Active state and retained lock.
- Unix integration tests start a helper owner process, force stop it, and verify process exit before Stopped is stored. Platform-specific tests preserve Windows build coverage without assuming Unix signals.
- Watch tests issue Stop Requests during status wait, quiet period, retry sleep, and merge-readiness wait; each asserts no later fetch or Review Source call.
- Agent Session cleanup tests seed active, failed, and closed Agent Selection scopes. Only active scopes are targeted, absence of the named Agent Session is silent, and secondary errors follow the primary failure.

## Build Order

1. Store compare-and-set completion, typed conflict, idempotent replay, and guarded integration reconciliation with store tests.
2. Active Agent Selection scope query and idempotent cleanup of registered Agent Sessions (depends on: 1).
3. Platform owner-process controller and reordered force-stop CLI flow (depends on: 1, 2).
4. Stop Request source across all Review Source wait boundaries (depends on: 1).
5. Winner-only terminal event/notification wiring and primary-before-secondary diagnostics (depends on: 1, 2, 3, 4).
6. User guide, command help, and spec finding traceability updates (depends on: 3, 4, 5).
7. Dedicated tooling-only update of
   `.agents/skills/roundfix/SKILL.md` and `skills/roundfix/SKILL.md`, with
   direct byte-identical edits and read-only sync verification (depends on: 6).

## Risks & Considerations

- PID identity is intentionally conservative. If the PID has been reused or cannot be signaled, force stop refuses rather than risking termination of an unrelated process.
- Same-outcome replay is idempotent but does not republish a lost notification. Notification retry history remains outside this spec.
- The Agent Selection lifecycle is append-only. Cleanup must select the latest status per scope, not every historical active record.
- A Stop Request is observed by the next configured poll boundary, not by a new high-frequency database loop. This avoids additional polling configuration and journal noise.
- Spec 0038 depends on the guarded Integration Pending reconciliation operation. Spec 0039 depends on winner-only outcome publication and stop-aware retry waits.

Rollout is an atomic binary change because no schema or public command syntax changes. Before replacing a running binary, operators must let Active Runs finish or stop them; mixing old and new owners against the same Run Database would preserve the original race. Rollback to the prior binary is data-compatible because existing terminal state values and tables remain unchanged, although it restores the unsafe completion behavior and therefore must not occur while a force stop is in progress.

## Decisions

- Force Stop fails closed until owner exit is proven.
- `CompleteRun` remains the ordinary completion boundary and becomes compare-and-set; a separate operation owns the only terminal reconciliation transition.
- Existing Agent Selection lifecycle records serve as the Agent Session registry.
- Protected Roundfix Skill publication is isolated in one Task whose changed
  files are the two authorized `SKILL.md` paths and its own Task file; it does
  not run the broad `make skills-sync` mutation target.
- See [ADR-0052](../../adr/0052-run-completion-is-compare-and-set.md).
