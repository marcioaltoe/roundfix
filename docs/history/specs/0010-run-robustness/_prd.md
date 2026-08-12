---
spec: 0010-run-robustness
status: active
created: 2026-07-06
surfaces: [cli, infra]
---

# Run Robustness

The 0009 cycle exposed three process-lifecycle weaknesses in one day.
Removing a config key hard-failed every Run on a machine whose config still
carried it — the QA "ran overnight" as a corpse. Externally killed Runs
leaked forty orphaned adapter processes, because cancelling an Agent Session
never reaps its OS process tree. And a Run's lifetime silently depends on
its caller's: three Runs died mid-flight when the invoking session reaped
its background tasks, leaving `nohup` gymnastics as the only defense. This
Spec closes all three: config migrations that never break users, Runs that
clean up their processes, and Runs that survive their caller.

## Goals

- An existing configuration keeps working across roundfix upgrades:
  recognized removed keys warn and are ignored, never rejected. See
  ADR-0027.
- No Run leaves adapter processes behind: force-stop and the preflight
  sweep close the Run's Agent Sessions so their process trees terminate.
- A Run can outlive its caller: `--detach` starts it as an independent
  session leader with the run id reported and Attach/Stop as the follow and
  control surfaces. See ADR-0028.

## User Stories

1. As a developer upgrading roundfix, I want configs carrying removed keys
   to keep working with one clear deprecation warning, so that an upgrade
   never turns my own configuration into a Preflight failure.
2. As a developer force-stopping a dead Run, I want its Agent Sessions
   closed and their processes gone, so that kills stop accumulating
   orphaned adapters.
3. As a developer whose machine collected orphans anyway, I want the
   preflight sweep to close sessions belonging to terminal Runs, so that
   debris self-heals on the next Run.
4. As a developer running roundfix from a script, CI step, or supervised
   session, I want `--detach` to hand me the run id and exit, so that the
   Run's survival never depends on my process staying alive.
5. As a developer with a Detached Run, I want its console stream in a
   per-Run log and the attach/stop hints printed at start, so that
   following and controlling it is obvious.

## Core Features

1. **Config deprecation path.** A deprecated-keys table in the config
   package: keys listed there (starting with `resolve.concurrent` →
   `worktree.concurrency`) parse without error, are ignored, and emit one
   stderr warning naming the replacement; truly unknown keys keep failing
   strict validation. See ADR-0027.
2. **Session-close reaping.** `stop --force` closes the Run's Agent
   Sessions (Run-level and per-Task) after the cooperative cancel, so acpx
   terminates the owner and adapter processes; the implement preflight
   sweep additionally closes roundfix-named sessions whose Runs are
   terminal. Best-effort, reported one line per closed session, never
   fatal.
3. **Detached Runs.** `--detach` on resolve, watch, and implement: roundfix
   re-executes itself as a session leader decoupled from the terminal,
   receives the run id through a startup handshake, prints it plus the
   attach and stop hints on stdout, and exits 0. Detach implies
   non-interactive mode; the detached process writes its console stream to
   a per-Run log under the Artifact Directory, named in the handshake
   output. Failures before the handshake surface on the caller's stderr
   with the usual exit codes. See ADR-0028.

## User Experience

One new flag on the three operational commands; one new warning line shape
for deprecated keys; force-stop and the sweep gain per-session close report
lines. `--detach` stdout is deterministic: the run id line, the log path
line, and the attach/stop hint lines — nothing else. Everything not named
here is byte-stable.

## Non-Goals / Out of Scope

- The long-lived `roundfix serve` daemon and work queue (work-plan item 3)
  — detach is the stepping stone, not the destination.
- Auto-migration that rewrites the user's config file — warn-and-ignore
  only.
- A general process-manager or `doctor` command — session-close reaping
  covers the observed leak; a doctor waits for real demand.
- Windows process-group semantics.

## Success Metrics

- The exact 0009 failure replays green: a config carrying
  `resolve.concurrent` runs with a warning instead of dying at Preflight.
- After force-stopping a live Run, no roundfix-spawned adapter process
  remains (asserted in the gated real-acpx test); the sweep closes sessions
  of terminal Runs.
- A Detached Run survives its caller's death (caller killed immediately
  after the handshake) and completes to a journaled terminal outcome,
  attachable throughout.
- This Spec's own implement Run executes its independent Wave at
  concurrency 2 — the live test of 0009's scheduler.

## Decisions

- Deprecated keys warn and are ignored via a single registry table; unknown
  keys still fail. See ADR-0027.
- Reaping is session-close based (acpx owns its process tree); no raw
  pkill in product code. Explicit-PID kills remain an operator action.
- Detach re-execs self with a handshake for the run id; detached console
  goes to a per-Run log under the Artifact Directory unconditionally (it is
  the Run's only record). See ADR-0028.
- The glossary gains **Detached Run**.

## Open Questions

None.
