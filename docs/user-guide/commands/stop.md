## stop

```bash
roundfix stop <run-id>
roundfix stop --force <run-id>
roundfix stop --force --owner-identity-unreadable <run-id>
```

Selectors: positional `<run-id>`, `--run-id`, `--pr`, `--spec`, or
`--head-repo` plus `--head-branch`. Graceful stop records a Stop Request and
reports `Stop Request recorded; the Run stops after the current Work Item
settles.`

During a watch Run's Review Source status, retry, quiet-period, or
merge-readiness wait, the owner checks for the Stop Request before the next
status access and after each interruptible sleep. It reaches Stopped by the
next configured poll boundary. After detecting the request, it does not run
another fetch, check, commit, push, or Review Source mutation. A Work Item
already in flight still settles before the graceful stop completes.

Force Stop is for dead, stuck, or runaway Runs. It cancels only registered
Agent Sessions whose current lifecycle is active, then terminates the recorded
owner process and proves that process exited. Only then does Roundfix report
Stopped, release the Active Run lock, and reap kept terminal Worktrees whose
branch has no commits beyond its base. An already-absent registered Agent
Session is an idempotent cleanup result.

Roundfix captures and compares owner identity through a direct kernel read:
procfs on Linux and sysctl on macOS. It spawns no subprocess, so Force Stop can
prove ownership on a host that cannot fork. If capture fails when a Run starts,
Roundfix prints one startup warning that the Run has PID-only reuse protection.
`roundfix runs list` keeps that condition visible by appending
`owner_identity_unproven=true` to the Run.

Owner identity proof can refuse Force Stop for two different reasons:

- A proven mismatch means the live process has a different comparable start
  identity from the recorded owner. Investigate PID reuse and do not signal
  that process. The Run remains Active, its lock stays retained, and
  `--owner-identity-unreadable` never applies.
- An unreadable identity means the host could not read or compare the identity.
  A kernel-read diagnostic includes the host error and directs you to resolve
  the host resource failure, then retry the normal Force Stop.

Use `--owner-identity-unreadable` only as an operator action of last resort,
after the normal Force Stop specifically reports an unreadable identity. It
authorizes PID-only termination for that condition. If identity is readable or
proves a mismatch, Roundfix exits `2` without signaling the process. The flag
cannot be enabled by configuration, environment, a default, or a timeout.

If owner exit cannot be proven, Force Stop fails with no stdout success report.
The diagnostic names the Run ID, owner PID, and failed process-control step;
the Run remains Active and its Active Run lock stays retained. Inspect the Run
with `roundfix runs list --state active`, resolve the reported owner-process
failure, then retry `roundfix stop --force <run-id>`. Agent Session cleanup
failures remain visible as secondary warnings after the primary failure. They
do not replace that failure or authorize terminal completion while the owner
is still alive.

Exit codes: `0` for a recorded Stop Request, a completed Force Stop, and the
idempotent already-Stopped report; `1` when Force Stop fails operationally
because owner exit cannot be proven; `2` for Preflight Validation failures
such as an invalid selector, no matching Active Run, an invalid
`--owner-identity-unreadable` precondition, or stopping a Run that already
holds a different terminal outcome.

Terminal results are stable. Repeating Force Stop for an already Stopped Run
reports the existing outcome without repeating process or Agent Session
actions. Force Stop against a different terminal outcome is rejected and
leaves that outcome unchanged.

Orphaned locks rarely need `--force` anymore: Runs record their owner process
id, and any command blocked by a lock whose owner is provably dead reclaims it
automatically — the Run completes Failed with the reason journaled and one
stderr warning names the reclaimed run id. A live owner, a PID-less legacy
Run, or any liveness result short of proof still blocks; a warning alone never
authorizes owner reclamation.

The terminology and behavior trace to the
[Roundfix glossary](../../../CONTEXT.md#language),
[ADR-0052](../../adr/0052-run-completion-is-compare-and-set.md),
[Spec 0037](../../history/specs/0037-terminal-outcome-integrity/_prd.md), and the
[detached-watch finding](../../history/findings/2026-07-16-vortex-pr87-detached-watch-notification.md#4-cleanup-noise-appeared-before-the-actionable-failure).

