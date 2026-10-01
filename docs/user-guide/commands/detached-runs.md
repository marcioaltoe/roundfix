## Detached Runs

`--detach` is available on `resolve`, `watch`, and `implement`. The foreground
command prints exactly five stdout lines and exits `0`:

```text
Run ID: <run-id>
Console Log: <path>
Attach: roundfix attach <run-id>
Supervisor monitor: roundfix events <run-id> --follow --filter outcome
Stop: roundfix stop <run-id>
```

The handshake is two-phase: the child writes a liveness marker immediately on
entering child mode, before configuration load and Preflight Validation, and
the run-id line once the Run exists. The parent waits 10 seconds for liveness
and up to 5 minutes for Run creation, so a slow but healthy preflight (a real
Agent probe takes seconds) never fails a detach start. Every failure branch
prints an explicit stderr diagnostic — for example:

```text
roundfix: Detached Run child produced no liveness signal within 10s; killed (exit: <exit or signal>)
```

— followed by the child's console output when any exists. The detached child
owns the terminal outcome. Supervisors use the printed outcome command;
humans use `roundfix attach <run-id>` for the read-only Live Run View. Detach
implies non-interactive mode: `--interactive` is rejected and `--no-input` is
implied.

Run Outcome Notification delivery is best-effort and never changes the Run
outcome or exit code. Each attempt appends a separate durable Run Event with
route, completion time, and receipt status `sent`, `skipped`, or `failed`;
receipt success means the local route accepted the request, not that a person
saw it. The original `notify.command` variables remain available:
`ROUNDFIX_RUN_ID`, `ROUNDFIX_OUTCOME`, `ROUNDFIX_KIND`, and
`ROUNDFIX_TARGET`. Terminal context adds `ROUNDFIX_REASON`,
`ROUNDFIX_CONSOLE_LOG`, `ROUNDFIX_ATTACH_COMMAND`,
`ROUNDFIX_REVIEW_ISSUES_KNOWN`, and `ROUNDFIX_NEXT_ACTION`.

The review Evidence, artifact inheritance, Detached outcome, and notification
contracts trace to the [Roundfix glossary](../../../CONTEXT.md#language),
[ADR-0054](../../history/adr/0054-review-source-evidence-determines-review-outcomes.md),
[Spec 0039](../../history/specs/0039-review-source-evidence-and-detached-outcomes/_prd.md),
and the
[detached-watch finding](../../history/findings/2026-07-16-vortex-pr87-detached-watch-notification.md).

