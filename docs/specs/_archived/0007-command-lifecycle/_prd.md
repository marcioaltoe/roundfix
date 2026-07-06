---
spec: 0007-command-lifecycle
status: active
created: 2026-07-05
surfaces: [cli, infra]
---

# Command Lifecycle

Roundfix now implements Specs and cleans reviews end to end, but living with
it day to day still involves ceremony the tool should own: installing the
pinned acpx by hand (and discovering the local-adapter override the hard
way), no way to update an installed roundfix or learn a newer one exists,
stopping a live Run only from its own terminal with Ctrl-C, and spec Runs
that never push even when a repository wants its Clean branches published.
This Spec ships the four lifecycle pieces the first dogfood round asked for:
setup, upgrade, graceful stop, and configurable spec-Run push.

## Goals

- A fresh machine reaches a working Roundfix environment with one command —
  pinned acpx installed, runtimes probed, configs offered — idempotently.
- An installed roundfix can update itself to the latest release, and quietly
  tells the developer when it is behind.
- A live Run stops safely from any terminal: graceful by default (the
  current Work Item settles first), immediate with `--force`. See ADR-0022.
- A repository can opt its spec Runs into pushing at Clean; pull request
  creation stays out of scope. See ADR-0021.

## User Stories

1. As a developer on a fresh machine, I want `roundfix setup` to install the
   pinned acpx, verify Node and the runtime adapters, offer the acpx
   local-binary override when local adapters exist, and offer creating the
   missing configs, so that the first Run works without archaeology.
2. As a developer with roundfix installed, I want `roundfix upgrade` to
   fetch and install the latest released version, so that staying current is
   one command.
3. As a developer running any operational command, I want a quiet stderr
   note when my version is behind the latest release — never blocking, never
   failing offline — so that I learn about updates without asking.
4. As a developer watching a Run from another terminal, I want
   `roundfix stop` to end it gracefully — the current Task or Batch settles,
   then the Run stops — so that stopping never wastes verified work.
5. As a developer with a dead or runaway Run, I want `roundfix stop --force`
   to cancel the Agent cooperatively and release everything immediately, so
   that recovery is never blocked on a corpse.
6. As a developer whose repository wants Clean branches published, I want a
   Project Config key that makes spec Runs push at Clean, so that the
   publish step is config, not ritual.

## Core Features

1. **Setup Command.** Bootstraps and verifies the environment in one
   idempotent pass: Node present, pinned acpx installed (installing via the
   documented npm command on confirmation), runtime adapters probed, the
   acpx local-binary agents override offered when local adapter binaries are
   on PATH, and User/Project Config creation offered where missing (through
   the existing Init Command flows). A non-interactive accept-all flag
   exists; every action taken or skipped is reported deterministically.
2. **Upgrade Command and freshness check.** `roundfix upgrade` resolves the
   latest release for the current platform, downloads it, verifies it, and
   replaces the running binary atomically; "already current" and "no
   releases published" are clean, distinct outcomes. Operational commands
   check freshness at most once per day (cached in Roundfix Home), emit one
   stderr line when behind, and stay silent offline or on any check failure.
3. **Graceful Stop Requests.** `roundfix stop` on a live Active Run records
   a Stop Request in the Run Database; engines honor it at the next
   settlement boundary — the in-flight Work Item verifies, settles, and
   commits first — ending the Run Stopped with the standard exit semantics.
   The stop command reports that the Request was recorded and what happens
   next. See ADR-0022.
4. **Force stop.** `--force` sends the cooperative cancel through the Run's
   Agent Session, then completes the Run in the database and releases its
   locks immediately — the recovery path for dead processes, documented as
   such. See ADR-0022.
5. **Spec-Run push at Clean.** A Project Config key (default off) makes a
   spec Run push its branch after — and only after — a Clean outcome,
   through the existing push machinery and upstream detection; Unresolved,
   Failed, Stopped, and failing-QA Runs never push; no pull request is ever
   created. See ADR-0021.

## User Experience

Two new commands (`setup`, `upgrade`), one new stop flag, one new config
key, one daily stderr line. All reports are deterministic; exit codes follow
the house contract (0 success/no-op, 1 operation failed, 2 validation).
Setup in interactive mode asks per action; `--yes` accepts everything.

## Non-Goals / Out of Scope

- Pull request creation (permanently out of scope; ADR-0021).
- The long-lived daemon, queue, or webhook triggers (work-plan item 3) —
  the Stop Request channel is deliberately the primitive item 3 reuses.
- Retry budgets, escalation (item 7), permission policy (item 8).
- Package-manager distributions (homebrew, etc.) — upgrade targets the
  released binary directly.
- Pause/resume semantics — stop is terminal; resume stays `implement`.

## Success Metrics

- On this machine, `roundfix setup` reports every check green (or performs
  the missing installs) and is a no-op on the second run.
- A live spec Run stopped from another terminal settles its in-flight Task
  (commit present) and ends Stopped; `--force` on the same shape ends it
  immediately with the Agent cancel journaled.
- With the config key on, a Clean spec Run pushes its branch; every
  non-Clean outcome leaves the remote untouched — proven in tests.
- Upgrade against a fixture release replaces the binary and reports the
  version change; the freshness note appears exactly once per day when
  behind.

## Decisions

- Setup wraps the existing Init Command flows for config creation rather
  than duplicating them; acpx install uses the exact documented pin command.
- Upgrade reads GitHub Releases for the repository the binary was built
  from; until the first release exists it reports "no releases published"
  — building the command ahead of the channel is deliberate.
- Freshness checks are best-effort, cached one day in Roundfix Home, and
  never block or fail a command.
- Graceful is the stop default; `--force` is the escape hatch for dead Runs
  (the flag also covers today's force-complete behavior). See ADR-0022.
- The push key lives in Project Config only (per-repository decision), off
  by default. See ADR-0021.

## Open Questions

None.
