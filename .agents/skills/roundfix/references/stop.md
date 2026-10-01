## Stopping Runs

Use `roundfix stop` for a graceful stop. Every selector keeps its existing
shape: `<run-id>`, `--run-id`, `--pr`, `--spec`, or `--head-repo` plus
`--head-branch`. For an Active Run, the default records a Stop Request in the
Run Database and returns this report line:

```text
Stop Request recorded; the Run stops after the current Work Item settles.
```

The owning Run finishes the in-flight Work Item's verification, settlement,
and commit boundary first, then ends Stopped through the normal outcome path.
During a watch Run's Review Source status, retry, quiet-period, or
merge-readiness wait, the owner checks for the Stop Request before the next
status access and after each interruptible sleep. It reaches Stopped by the
next configured poll boundary. After observing the request, it does not run
another fetch, check, commit, push, or Review Source mutation.

Use `roundfix stop --force` only for a dead, stuck, or runaway Run. It first
validates the recorded owner PID, terminates the recorded owner process, and
proves that process exited. Until owner exit is proven, registered Agent
Sessions and their Agent Selection lifecycles remain active.

Owner identity comes from a direct kernel read: procfs on Linux and sysctl on
macOS. Roundfix spawns no subprocess for capture or comparison, so the proof
still works when the host cannot fork. A Run whose identity capture fails at
creation still starts, prints this warning once, and carries a durable marker:

```text
roundfix: warning: Run <run-id> started with PID-only reuse protection because owner identity capture failed at creation.
```

`roundfix runs list` appends `owner_identity_unproven=true` to that Run.

Owner identity proof has two distinct refusal conditions:

- A proven mismatch means the live process has a different comparable start
  identity from the recorded owner. Leave the Run Active and its lock retained,
  investigate PID reuse, and do not signal that process. The
  `--owner-identity-unreadable` flag never overrides this refusal.
- An unreadable identity means the host could not read or compare the identity.
  For a kernel-read failure, the diagnostic includes the host error and tells
  the operator to resolve the host resource failure, then retry the normal
  Force Stop.

Only when the normal Force Stop specifically reports an unreadable owner
identity may an operator use this last-resort command:

```bash
roundfix stop --force --owner-identity-unreadable <run-id>
```

The flag authorizes PID-only termination for that one condition. It exits `2`
and signals nothing when identity is readable or proves a mismatch; no
configuration, environment variable, default, or timeout can activate it.

After owner exit proof, Force Stop cancels and closes only registered Agent
Sessions whose latest Agent Selection lifecycle is active. No active lifecycle
record means no session action, and an already-absent registered session is an
idempotent cleanup result. Other cleanup failures remain visible as secondary
warnings.

Only after owner exit proof does Roundfix complete the Run as Stopped, release
its Active Run lock, and reap eligible kept terminal Worktrees. If owner exit
cannot be proven, Force Stop prints no stdout success report; its diagnostic
names the Run ID, owner PID, and failed process-control step. The Run remains
Active with its Agent Sessions unchanged and its Active Run lock retained.
Inspect it with `roundfix runs list --state active`, resolve the reported
owner-process failure, and retry `roundfix stop --force <run-id>`.

After owner exit proof and successful Stopped completion, the force-stop report
title includes:

```text
Roundfix Run force-stopped
```

Terminal completion is compare-and-set. The winning transition alone publishes
the terminal outcome event and notification. Repeating Force Stop for an
already Stopped Run reports the stored outcome without repeating process,
session, event, or notification actions. A different terminal outcome is
rejected and preserved; a losing owner observes that stored outcome and exits
without publishing another terminal event or notification.

When an Active-Run lock records an owner PID and Roundfix can prove that owner
process no longer exists, preflight reclaims the orphan automatically: the Run
settles Failed, the Run Event Journal records the reclamation, one stderr
warning names the Run id and PID, and the blocked command proceeds. A live
owner, a PID-less legacy Run, or any liveness result short of proof still
blocks with the existing `roundfix stop <id>` guidance; `stop --force` remains
the manual path for those cases.

Force stop also reaps kept Run or Task Worktrees and branches for terminal Runs
whose branch has no commits beyond its base. Each removed pair is reported on
stderr with this shape:

```text
roundfix: reaped terminal Worktree path=<path> branch=<branch>
```

The Implement Command preflight sweep uses the same worktree reaping report and
the same `roundfix: closed session <session>` /
`roundfix: could not close session <session>: <reason>` session-close reports
for roundfix-named Agent Sessions whose Runs are terminal. Active, unknown, and
non-roundfix sessions are ignored.

