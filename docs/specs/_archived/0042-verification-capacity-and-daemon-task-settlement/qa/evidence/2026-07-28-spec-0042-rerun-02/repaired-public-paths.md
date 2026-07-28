# Repaired public paths

## Queued Stop Request

The current regression composes the same public boundary that failed on build
`8593002`:

- `TestRunStopByRunIDRecordsStopRequest` proves the public Stop Command writes
  the durable Run-store Stop Request.
- `TestTaskCycleStopRequestWhileQueuedForVerificationStartsNoCommandAndStaysResumable`
  records that durable request while one Task owns Verification Capacity and
  another is waiting.

Observed on `ffd6852`: the waiting Task publishes the Stop event, starts zero
Verification commands, remains `in_progress` in its Task Worktree, and
receives no terminal settlement. The already-running Task completes and
integrates, the Run returns `ErrStopRequested`, and worker/capacity state
drains. The adjacent direct-cancellation and later-acquisition cases also pass.

## Attach flag order

Current built help prints:

```text
roundfix attach [<run-id>] [--no-input]
```

`TestAttachAcceptsDocumentedFlagOrders` replays a real stored terminal Run
through the public CLI runner with all three supported forms:

```text
roundfix attach <run-id> --no-input
roundfix attach --no-input <run-id>
roundfix attach <run-id>
```

All exit 0 and replay the same read-only result. The current built
`bin/roundfix` also accepted the exact trailing and leading forms against
stored Run `run_20260728T134451Z_ec12a53008910524`; both exited 0 and rendered
the same Unresolved Run with Task Capacity 1 and Verification Capacity 1.
