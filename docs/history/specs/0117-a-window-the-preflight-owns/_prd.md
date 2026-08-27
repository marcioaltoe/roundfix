---
spec: 0117-a-window-the-preflight-owns
status: archived
archived: 2026-08-26
created: 2026-08-26
surfaces: [cli, backend]
---

# A window the Preflight owns

A Supervisor running an unattended session bounds it by wall-clock time: work
until the cutoff, then stop opening new work and let what is running finish.
Today that bound lives in the Supervisor's own discipline — a script it must
remember to consult before each Run. A bound the caller must remember to check is
not a bound: skip the check once and Runs keep opening for the rest of the
night, with the guard installed and working. The Preflight already refuses Run
creation for a dozen reasons and is the one place a Run cannot get past. The
cutoff belongs there.

## Project Constraints

- Identifier strategy: applicable — this Spec coins one term for the bound it
  adds, and the obvious name collides: `Session` already denotes an Agent
  Session in the glossary, so a "session cutoff" would be ambiguous at the
  first reading. The closing node checks whether the work introduced or changed
  a term the glossary should carry. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no network boundary, credential, or
  request is created or read. The work is Preflight logic and Run Database
  state. Source: `docs/agents/specific-repository.md`.
- Active ADR obligations: applicable — ADR-0004 puts Run state in one central
  Run Database rather than per-repository files, which decides where the bound
  is stored and rules out the `.git/` state file the measured prior art used.
  ADR-0139 keys that database's locks by work target, supplying the scoping
  precedent this bound follows rather than inventing one.
  ADR-0137 requires the Run budget to be explained where it is configured, and
  this Spec adds a second time bound beside `budget.max_run_duration`, so both
  must be explained where a reader meets them. ADR-0133 requires a diagnostic
  to name the literal it requires, which governs the refusal message. ADR-0022
  routes Stop Requests through the same database and honours them at the next
  settlement boundary — the shape this Spec deliberately does not take: a
  cutoff that ended a running Run would be a Stop Request on a timer, and
  stopping a Run already under way is this Spec's first Non-Goal. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — no protected tooling mutation is proposed or
  authorized. The work is production Go in the CLI, preflight, and store
  packages plus their tests, and the user-guide page that documents the
  command. Source: `docs/agents/agent-instructions.md`.

## Goals

1. A Supervisor can bound an unattended session by wall-clock time without the
   bound depending on the Supervisor remembering to consult it.
2. A Run already under way is never interrupted by the bound.
3. A reader who meets either time bound — this one or the per-Run duration —
   learns how the two relate where they meet it.

## User Stories

1. As a Supervisor starting an unattended session, I want to declare when it
   should stop taking on new work, so that the machine stops opening Runs at
   that hour whether or not I am watching.
2. As a Supervisor whose cutoff has passed, I want the refusal to name the
   cutoff, the current time, and how to clear or move it, so that I can act
   without reading source or guessing the command.
3. As a Supervisor whose Run is running when the cutoff passes, I want that Run
   to finish and merge, so that the bound never leaves implemented work with no
   Pull Request.

## Core Features

1. **A stored bound on Run creation.** A Supervisor sets a cutoff as an
   instant; the Preflight refuses to create a Run once that instant has passed,
   with no Run created and no side effects, joining the existing refusals. The
   bound is stored in the Run Database, scoped to the repository the session
   works, and it is durable: a restarted session, a compaction, or a new
   process finds the same cutoff.
2. **The cutoff is the next occurrence of a wall-clock time, never a
   comparison against today.** Setting a cutoff of `07:00` at 23:00 means
   07:00 tomorrow. A naive same-day comparison reports the window already
   closed at the moment a night session starts, which is the failure the prior
   art documents in its own header.
3. **Setting a cutoff does not silently move one.** Re-declaring a cutoff that
   already exists reports the standing one and changes nothing unless the
   caller asks for the change explicitly. Re-opening a window by accident is
   the failure mode that costs a whole night of unattended Runs.
4. **The bound governs starting, never finishing.** A Run created before the
   cutoff runs to its own terminal outcome. Where a Run is created with less
   than `budget.max_run_duration` remaining before the cutoff — so it may cross
   it — the Run is still created and the crossing is reported, not refused.
5. **The two time bounds explain each other where they are met.** The command
   that sets the cutoff, the refusal it produces, and the configuration surface
   carrying `budget.max_run_duration` each state which bound governs starting
   and which governs duration.

## User Experience

Not applicable as a browser surface. The affected surfaces are a CLI command
that sets, shows, and clears the cutoff; the Preflight refusal on `implement`;
and the user-guide page documenting both.

## Non-Goals / Out of Scope

- Stopping, pausing, or shortening a Run that is already under way. The bound
  governs Run creation only; `roundfix stop` remains the way to end a Run.
- Scheduling. This Spec adds no timer, daemon, or wake-up: nothing starts work
  at an hour, and the cutoff is read only when a caller asks to create a Run.
- Changing what `budget.max_run_duration` bounds or where it is evaluated.
  This Spec explains the relationship and reports a crossing; ADR-0137 already
  settled that the evaluation point does not move without reproduced evidence.
- Extending the bound to `fetch`, `resolve`, or `watch`. Those Runs answer a
  Pull Request that is already open, and refusing them by clock would strand a
  review Round rather than bound an authoring session.
- A supervisor-loop skill. The discipline that consumes this bound is queued
  separately, and depends on a decision this Spec does not make.

## Success Metrics

- The measured failure mode becomes unreachable: with a cutoff set and passed,
  `roundfix implement` refuses, and the refusal is produced by the Preflight
  rather than by a caller that chose to check. Proven by invoking the command
  directly, not through a supervisor loop.
- The night-start case works: a cutoff of `07:00` set at 23:00 leaves the
  window open, which the prior art's own header records as the bug a naive
  comparison produces.
- A Run created before the cutoff and still running when it passes reaches its
  own terminal outcome.

## Decisions

- The bound is stored in the Run Database rather than in a file under `.git/`.
  Measured: the repository has no per-session state file today — only the
  database, its advisory lock, and a version-check cache — and ADR-0004 already
  decided that Run state is central rather than per-repository. A second state
  store would be a second thing to keep honest.
- A crossing is reported rather than refused, because refusing it would deny
  the last Run of a window on a prediction, and the measured intent of the
  prior art is that the cutoff traps the start and never the finish.
- The term is settled during authoring rather than at the gate, because
  `Session` is taken: the glossary's `Agent Session` is a different thing, and
  a coined term that reads as an existing one is a defect this Spec's own
  Identifier constraint would raise against it.

## Open Questions

- Whether a cutoff should be scoped to the repository or to the machine. The
  default until answered is the repository, because a Supervisor session works
  one checkout and a machine-wide cutoff would bound an unrelated repository's
  Runs — but the Run Database is machine-wide, so the narrower scope is a
  choice the TechSpec must make explicit rather than inherit.
