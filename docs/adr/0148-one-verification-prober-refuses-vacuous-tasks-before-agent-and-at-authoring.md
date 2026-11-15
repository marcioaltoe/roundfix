---
status: accepted # proposed | accepted | rejected | deprecated | superseded
created_at: 2026-08-26T00:00:00Z
updated_at: 2026-11-15T11:56:04Z
deprecated_at: null # null or YYYY-MM-DDTHH:MM:SSZ
superseded_by: null # null or ADR-NNNN
---

# One Verification prober refuses vacuous Tasks before the Agent and at authoring

The Daemon runs a Task's `## Verification` commands once against the unchanged tree before creating the Agent Session and refuses the Task when any command already passes — a command that exits zero before the work exists proves nothing about the work (measured on Spec 0089, whose gate matched a string already in the file and settled a Task whose work was never done). The authoring-time check asks the identical question earlier, so the vacuous/failing/unknown classification lives in one extracted prober that both callers use: a second implementation would be free to disagree, and a checker that approves what the probe later refuses is the same defect one layer up. The Daemon keeps its Run bookkeeping around the shared loop; the authoring caller supplies a working directory and nothing else.

Consolidates ADR-0109 and ADR-0124 (2026-08-26); both are retired, and Git keeps their text.

Amendment (2026-10-26): both callers probe the same Tasks. The Daemon's
`TaskAwaitsWork`, true for every status but `completed`, chooses the Tasks the
Daemon schedules and the Tasks `roundfix spec check --run-verification`
probes. A completed Task's command passes because its work exists, so the
authoring check reports it `not probed` with its reason and never fails on it;
the vacuous, honest, unknown and malformed classification is unchanged.

Amendment (2026-11-01): neither caller probes the QA gate. The Daemon's
`daemon.TaskIsQAGate` identifies the Task of type `qa`, which the Daemon
routes to its QA step and never probes. Roundfix derives its Verification,
which would pass on a kept QA Report even when its verdict is fail, so the
check reports an unfinished gate `not probed` with reason `qa_gate`;
a completed gate keeps reason `completed`.

**Amended by [ADR-0406](0406-the-authoring-verification-probe-prepares-its-checkout-as-the-daemon-prepares-a-task-worktree.md) (2026-11-14).** The authoring caller no longer supplies a bare working directory: before its first probed command it prepares its disposable checkout with the `worktree.copy` files and the Worktree Bootstrap a Task Worktree receives, through the worktree package's provisioning function, so both callers probe a prepared tree. A failed bootstrap ends the check before any command is probed. The shared prober, its classification and the Daemon's Run bookkeeping still hold.
