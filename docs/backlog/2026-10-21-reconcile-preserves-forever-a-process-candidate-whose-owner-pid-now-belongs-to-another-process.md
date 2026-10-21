---
type: fix
status: open
created: 2026-10-21
spec: null
reason: null
---

# Reconcile preserves forever a process candidate whose owner PID now belongs to another process

## Symptom

`roundfix reconcile` (read-only report) lists, on every invocation, `kind=process` preserved
candidates for terminal Runs from 2026-10-01, 2026-10-12 and 2026-10-14 with the refusal reason
`owned process tree cannot be proven: owner process <pid> failed step "prove owner process identity":
owner process identity cannot be proven: live process start identity does not match the recorded owner
identity`. The recorded owner exited long ago and the operating system reused its PID for an unrelated
process. The identity mismatch proves the recorded owner is gone, yet the candidate is never released,
so the report repeats the same entries indefinitely.

## Where

`internal/store/process.go` (`OwnerProcessControl`'s owner proof returns
`ErrOwnerProcessIdentityUnproven` when the live start identity differs from the recorded one) and the
reconcile debris classification of process candidates in `internal/cli` (the `Debris summary`
`preserved` count).

## Expected

A proven identity mismatch never signals the process now using the PID (Force Stop's rule stays). For
reconcile, the mismatch is proof the recorded owner tree is gone: `--apply` releases the stale process
record without sending any signal, and the report classifies it as releasable instead of preserved. An
unreadable identity stays preserved, as today.

## Evidence

`roundfix reconcile` on the maintainer's machine, 2026-10-21 ~00:40 BRT (roundfix 0.101.0): three
such candidates (Runs `run_20261001T000219Z_0e911d72b0d72d2f`, `run_20261012T110224Z_380dd005763aa861`,
`run_20261014T174911Z_8d771e92d24f52c5`) among `preserved=6`. Kept in the backlog by the maintainer's
cycle AN decision ("V → W → X"; Y stays an open entry).
