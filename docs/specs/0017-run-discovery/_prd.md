---
spec: 0017-run-discovery
status: active
created: 2026-07-06
surfaces: [cli, docs]
---

# Run discovery

Detached Runs and Attach solved Run ownership, but discovering what to attach to
still depends on the caller having captured the run id from the detach report.
A user who lost the id — or an agent session that started fresh — has no command
that answers "which Runs exist for this repository, and which are Active?".
Run discovery closes that gap with a deterministic listing command and an
Attach picker over the same data, so any session can find and follow a Run
using only the CLI.

## Goals

- A user or agent lists this repository's Runs — id, state, kind, and target —
  from the Run Database with one command and no prior knowledge.
- `roundfix attach` with no argument becomes usable: it offers the Runs that
  can be attached instead of failing for a missing run id.
- Both surfaces read the same Run data, so the picker never shows a Run the
  listing would not.

## User Stories

1. As a user driving Detached Runs, I want to list this repository's Runs with
   their id, state, kind, and target, so that I can attach to or stop a Run
   whose id I did not capture.
2. As an agent in a fresh session, I want a deterministic, stable-format Run
   listing, so that I can discover the Active Run for a repository and report
   or follow it without scraping logs.
3. As a user at an interactive terminal, I want `roundfix attach` without
   arguments to offer a picker of Runs, so that I can open the Live Run View
   without copying a run id from another terminal.

## Core Features

1. A listing command prints this repository's Runs, newest first, one Run per
   line: run id, state, kind, and the Run's target (the Open Pull Request
   number or the Spec slug). Active Runs are visibly distinguishable from
   terminal Runs.
2. The listing is scoped to the current repository by default and can be
   widened to every repository in the Run Database with a flag.
3. The listing supports filtering to Active Runs only, so scripts can answer
   "is anything running here?" with one call.
4. With no matching Runs, the listing prints a clear empty result and exits
   zero — an empty repository is not an error.
5. `roundfix attach` with no run id in an interactive terminal opens
   Interactive Input listing the same Runs the listing command shows for the
   repository; selecting one opens the Live Run View for that Run. Attach
   semantics are unchanged: read-only replay and follow, never owning or
   stopping the Run.
6. `roundfix attach` with no run id in a non-interactive context fails with an
   actionable error that names the listing command.

## User Experience

- The listing is plain text on stdout, deterministic order, no interactivity.
- The Attach picker follows the existing Interactive Input pattern: a numbered
  list accepting a number or a run id, newest Runs first with Active Runs at
  the top, and the Run's state, kind, and target visible per entry.
- Cancelling the picker exits without attaching and without side effects.

## Non-Goals / Out of Scope

- JSON or other machine-readable output modes — the stable text columns are
  the contract for now.
- Mutating actions from the listing or picker (stop, gc, settle) — discovery
  is read-only.
- Cross-repository picking in the Attach picker — the picker serves the
  current repository; the wide listing flag covers the rest.
- Run detail views beyond what Attach already renders.
- Pagination or a follow/watch mode for the listing.

## Success Metrics

- Attaching to a Run whose id was not captured requires only Roundfix commands
  (list, then attach — or attach's picker alone), with no Run Database or log
  spelunking.
- An agent can discover the Active Run for a repository with one deterministic
  command call.

## Decisions

- Ship both surfaces in one Spec: the listing command is the deterministic
  contract, the Attach picker is a thin interactive layer over the same query.
- Scope defaults to the current repository; a flag widens to all repositories.
- The picker lists attachable Runs regardless of state — Attach replays
  terminal Runs from the Run Event Journal, so terminal Runs stay selectable.

## Open Questions

None.
