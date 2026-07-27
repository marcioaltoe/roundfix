---
status: pending
created_at: 2026-07-27
updated_at: 2026-07-27
---

# Force Stop — owner identity forks `/usr/bin/ps`, so the escape hatch fails exactly when the host is loaded (2026-07-27)

Spec 0037 gave Force Stop a real ownership proof: Runs record an opaque owner
start-time identity, and Force Stop refuses to signal a PID whose live identity
does not match. The proof is correct, but it is obtained by forking
`/usr/bin/ps` on every read. When the host cannot fork, the proof fails, and
Force Stop fails closed — refusing to stop the Run.

That inverts the availability property the command exists for.
`roundfix stop --force` is the escape hatch for a dead, stuck, or **runaway**
Run, and a runaway Run is precisely what exhausts a host's process slots. The
harder the machine is struggling, the less able Roundfix is to stop the Run
causing it.

This is not hypothetical. During the Spec 0038 Implement Run on 2026-07-27,
five store/CLI tests failed inside the Daemon's own `make verify` because the
host could not fork `/usr/bin/ps`. The same code path serves production Force
Stop.

## Mechanism

`OwnerProcessControl.proveOwner` (`internal/store/process.go:80`):

```go
liveIdentity, identityErr := processStartIdentity(ctx, pid)
if identityErr != nil {
    // exited between checks? absence is its own proof
    if absent, absentErr := processAbsent(pid); absentErr == nil && absent {
        return true, nil
    }
    return false, ownerProcessControlError(pid, "prove owner process identity",
        fmt.Errorf("%w: read live owner identity: %v", ErrOwnerProcessIdentityUnproven, identityErr))
}
```

`processStartIdentity` runs `ps -p <pid> -o lstart=`
(`internal/store/process_unix.go`). A transient
`fork/exec /usr/bin/ps: resource temporarily unavailable` is therefore
indistinguishable, in outcome, from a genuine PID-reuse refusal: the Run stays
Active, its lock is retained, and the operator is told the owner cannot be
proven. Retrying is the documented advice, but the retry forks again under the
same pressure.

## Second, quieter consequence

The same helper captures the identity when a Run is created. A capture failure
records NULL, which by design degrades to the legacy PID-only proof. Under host
pressure a Run can therefore be created **without** reuse protection, silently
and permanently, with no warning at creation and no indication later that the
Run is running unprotected. The degradation path exists for legacy rows
predating the column; it should not be reachable by transient load on a current
binary.

## Why forking is avoidable

The value is only ever compared for equality — it is an opaque token, never
parsed. Both supported platforms expose process start time without spawning
anything:

- Linux: field 22 (`starttime`) of `/proc/<pid>/stat`.
- macOS: `sysctl` `KERN_PROC_PID` → `kinfo_proc.kp_proc.p_starttime`.

Either yields a stabler token than `ps` output, needs no subprocess, and
removes the `TZ`/`LC_ALL` pinning that CodeRabbit already had to request on
PR #38 to keep the rendered text stable across locale changes.

## Suggested resolution

1. Replace the `ps` fork with a syscall/procfs read behind the existing
   build-tag structure, keeping the token opaque and equality-compared and
   leaving the Windows and non-Unix stubs as they are.
2. Separate "identity mismatch" from "identity unreadable" in the failure
   classification. A proven mismatch must keep failing closed. An unreadable
   identity is a different condition and deserves its own diagnostic naming the
   host resource failure, so an operator is not told to suspect PID reuse when
   the machine simply could not fork.
3. Make a failed identity capture at Run creation visible — one warning at
   startup and a recorded marker — so a Run running without reuse protection is
   never silent. Reserve the NULL degradation for genuinely legacy rows.
4. Decide, explicitly, whether an unreadable identity on a machine under
   pressure should keep failing closed. Failing closed is right when ownership
   is genuinely unknown, but the operator needs a documented way to stop a
   runaway Run when the host cannot answer at all.

## Suggested acceptance checks

- Owner identity is captured and compared with no subprocess spawned, proven by
  a test that fails if the implementation execs anything.
- Under simulated fork exhaustion, Force Stop distinguishes an unreadable
  identity from a mismatch and each carries its own diagnostic.
- A Run whose identity capture fails at creation emits a warning and is
  observably marked as unprotected.
- The existing reuse-refusal, matching-identity, absent-owner, and legacy
  no-token cases keep their current behavior.
- The full suite passes under parallel load without `ps`-related failures.

## What worked — keep

- Failing closed on a proven identity mismatch is correct and must survive any
  fix here; the problem is only that unreadable and mismatched share one
  outcome.
- Recording absence as its own proof (the owner exiting between the liveness
  check and the identity read) is right and already handled.
