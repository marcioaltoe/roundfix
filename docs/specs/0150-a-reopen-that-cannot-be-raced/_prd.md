---
spec: 0150-a-reopen-that-cannot-be-raced
status: active
created: 2026-09-19
surfaces: [backend, cli, docs]
---

# A reopen that cannot be raced

Spec 0149 shipped `roundfix reopen` with two defects its third review round found
and its corrective ceiling would not let it fix. Both are carried here rather
than dropped.

**It decides, then writes, with nothing in between.** `preflightReopen` reads the
Task Graph, confirms the gate is stale, and returns a plan holding the QA Task's
path and the stale dependency ids. `ReopenGate` then rewrites the file from that
plan. If a dependency settles `completed` in the gap, the command resets a gate
that is no longer stale and records an invalidation naming dependencies that are
no longer incomplete — a false record, written by the command whose purpose is
to keep records honest.

**Its atomic write replaces a symlink instead of writing through it.** Spec 0149
added `replaceTaskFile` to make the status rewrite and the invalidation record
land together. It renames a temporary file onto the Task path. Where the Task
path is a symlink to a file inside the Spec Root — which the confinement check
allows — the rename replaces the link with a regular file. The real Task keeps
`status: completed` and the manifest now points at a copy. The predecessor,
`SetStatus`, used `os.WriteFile` and wrote through the link; the atomicity fix
silently changed that.

Neither defect is reachable in this repository today: no Run can hold the
checkout while `reopen` runs, and no Task file here is a symlink. Both are
reachable in the product.

## Project Constraints

- Identifier strategy: applicable — the command keeps addressing a Spec by slug
  and a Task by its Task Graph id, and mints no identity. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local file reads and writes only.
  Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0155 makes the `qa` Task declare the
  gate's matrix and ADR-0104 makes a Spec accept on evidence it did not author;
  ADR-0156 makes each declared promise name a consuming Task. All three hold.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the new refusal is an exit-2 condition, and
  exit codes are public API, so the shipped skill must name it. Express
  maintainer authorization: granted 2026-09-18 as a standing grant, consumed
  here and recorded in [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`. Sanctioned
  regeneration: `make skills-sync`. `internal/cli` and `internal/spec` are
  ordinary source that no authorization has bounded, measured through the
  governance probe. Source: `docs/agents/agent-instructions.md`,
  `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- A reopen writes only if the gate is still stale at the moment it writes.
- An atomic write preserves whatever the Task path is.

## Core Features

1. **The decision is rechecked against the tree it writes to.** Immediately
   before the replacement, the command re-derives the gate's staleness and the
   stale dependency set. If either changed since preflight, it refuses and
   writes nothing.
2. **A refusal from that recheck is a named condition.** It exits 2 like the
   other preflight refusals and says the gate changed underneath, rather than
   failing with a generic error.
3. **The atomic write follows a symlink to its target.** Where the Task path is
   a symlink, the replacement lands on the file the link points at, leaving the
   link intact — the behavior `SetStatus` had before Spec 0149 made the write
   atomic, now with atomicity kept.

## Non-Goals / Out of Scope

- A lock. The recheck closes the window that matters for a command that already
  refuses while a Run is active; a mutation lock across the Spec directory is a
  larger contract than this defect justifies.
- The shared `commandWantsHelp` behavior, where a bare `help` token anywhere in
  the arguments returns usage and exit 0 before flag parsing. It is identical
  across more than a dozen commands, so fixing it in one makes the surface
  inconsistent. It belongs to a Spec that owns the shared helper.
- Any change to when the loader raises its stale-gate refusal, or to the
  refusals Spec 0149 already ships.

## Success Metrics

1. A fixture whose stale dependency is completed between preflight and the write
   is refused, and its QA Task file is byte-identical afterwards.
2. A fixture whose Task path is a symlink to a file inside the Spec Root has the
   target rewritten and the link preserved as a link.
3. Every reopen behavior Spec 0149 shipped still holds.

## Decisions

- **Recheck rather than lock.** The command already refuses while a Run holds
  the tree, so the remaining window is narrow and a recheck closes it at a
  fraction of a lock's cost and blast radius.
- **Follow the link, do not reject it.** Rejecting a symlinked Task path would
  be simpler and would break a layout that worked before Spec 0149. The
  regression was introduced here, so it is repaired here.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in the
Task Graph. The negative cases carry the weight: a recheck that never refuses,
or a symlink test that asserts only that the write succeeded, would pass while
leaving both defects in place.

## Research basis

Both defects were reported by independent pre-PR review of Spec 0149 and
verified in source before this Spec was authored. `replaceTaskFile` was
confirmed new to Spec 0149 with `git log -S`, and `SetStatus` was confirmed to
have used `os.WriteFile`, which writes through a symlink — establishing the
change in behavior rather than assuming it. The governance probe classified
every path this Spec touches.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
