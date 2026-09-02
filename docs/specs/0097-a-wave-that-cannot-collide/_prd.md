---
spec: 0097-a-wave-that-cannot-collide
status: active
created: 2026-08-12
surfaces: [backend, infra]
---

# A wave that cannot collide

Raising Task Worktree concurrency made three latent failures visible at once, and
all three cost whole Runs of finished Agent work. Sibling worktrees bootstrap
against one shared Git directory and collide on its lock, reporting failure after
having done every byte of the work. Tasks the graph declares independent edit the
same file and die at integration, which never appeared while Tasks ran one at a
time. And creating a second Task Worktree under suite load fails with a raw
filesystem errno that names neither the Run, the Task, nor the concurrency that
produced it. Concurrency is configured as a number today and verified as nothing;
this Spec makes the safety a property the graph can prove before it dispatches.

## Project Constraints

- Identifier strategy: applicable — Run, Task Worktree, Run Worktree, and Wave are
  glossary terms this Spec reports on, and a collision message that invents a
  synonym for one of them is a defect. The closing node checks whether the work
  introduced or changed a term the glossary should carry. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no network boundary, credential, or
  request is created or read; the work is process orchestration, filesystem
  isolation, and error reporting inside a local daemon. Source:
  `docs/agents/specific-repository.md`.
- Active ADR obligations: applicable — ADR-0014 gives the Daemon ownership of task
  verification and status settlement, so a bootstrap or worktree failure this Spec
  reclassifies must settle through that owner rather than beside it. The decisions
  that extend that ownership are accounted and none applies here: ADR-0020 ranks a
  parsed prompt result above the runtime's exit code, ADR-0038 allows one
  Verification repair, ADR-0057 gives the Daemon exclusive ownership of Implement
  Task status, and ADR-0096 adds the gate's mechanical stage — this Spec changes
  none of those four, and its collision report is a pre-dispatch fact rather than a
  settlement. ADR-0056 separates Task Capacity from Verification Capacity and does
  apply: this Spec constrains which Tasks may share a wave and changes neither
  capacity. ADR-0117 places a check with the stage that can produce its defect, and
  applies directly — the collision this Spec detects is produced at dispatch, which
  is where the check is placed rather than at integration. No accepted ADR governs
  wave composition, which is why the rule this Spec adds is new rather than a
  restatement. ADR-0127 places process residue in the readiness diagnostic; this Spec schedules waves before dispatch and inspects no process table, so it does not apply. ADR-0148 establishes that a rule enforced at authoring and at Run time lives in one extracted prober both callers use, so a checker cannot approve what the Run later refuses; the collision rule follows that shape. ADR-0135 makes an absent diagnostic a reported state rather than an empty message, which is the principle the worktree failure message applies to a raw errno. ADR-0093 bounds the Spec Consistency Check to what artifacts say rather than to inference, which is the line the collision rule works inside: it reads paths that resolve to repository files and never infers a Task's intent. Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — no protected tooling mutation is proposed
  or authorized. The work changes daemon, worktree, and graph-validation behavior
  in production Go and its tests, creating or editing no linter, formatter,
  test-runner, build, or skill configuration. Source:
  `docs/agents/agent-instructions.md`.

## Goals

1. Two Tasks that would edit the same file never run in the same wave unnoticed.
2. Bootstrapping many Task Worktrees is safe at any configured concurrency.
3. A worktree that cannot be created explains itself in the loop's own terms.
4. Raising concurrency stops being a wager the Supervisor places without evidence.

## User Stories

1. As a Supervisor raising Task Worktree concurrency, I want the tool to tell me
   which Tasks in a wave touch the same files, so that I fix the graph before a
   Run rather than after losing one.
2. As a Supervisor running a Spec at concurrency above one, I want worktree
   bootstrap to succeed regardless of how many siblings start together, so that a
   Run does not die on a lock while its work is already done.
3. As a Supervisor reading a failed Run, I want a worktree failure to name the
   Run, the Task, and the concurrency, so that I can act without translating a
   filesystem error.

## Core Features

1. **Wave collision detection.** Before dispatching a wave, the loop reports when
   two Tasks in it declare or are known to touch the same file, naming the Tasks
   and the paths. It either serializes them or refuses the wave; the choice is a
   decision this Spec's design settles, not a per-run option.
2. **Serialized worktree bootstrap.** Bootstrap is serialized across sibling Task
   Worktrees even when the Tasks themselves run in parallel, so concurrent writes
   to a shared Git directory and a shared package cache cannot collide.
3. **A bootstrap failure that does not lie.** A bootstrap that fails after
   completing its work reports that state distinctly from one that failed before
   starting, so a maintainer is not told nothing happened when everything did.
4. **A worktree error in Roundfix's own words.** A Task Worktree that cannot be
   created reports the Run, the Task, and the concurrency level, carrying the
   underlying error as evidence rather than as the message.

## Non-Goals / Out of Scope

- Changing the default concurrency, or recommending one.
- Deciding wave composition rules in the authoring skills; this Spec detects a
  collision, and teaching `write-tasks` to avoid authoring one is a separate Spec.
- Integration conflict resolution. This Spec prevents the collision; merging two
  conflicting changes remains out of scope.
- Any change to how Tasks are settled or verified.

## Success Metrics

- A graph whose same-wave Tasks share an edit target is reported before dispatch,
  proven against a graph shape measured in a repository this Spec did not build.
- Bootstrapping at concurrency above one produces no lock collision across
  repeated runs.
- A Task Worktree creation failure names the Run, the Task, and the concurrency.
- A Run that previously died at integration on a same-file collision instead
  reports the collision before any Agent Session opens.

## Decisions

- The collision check runs before dispatch rather than at integration, because the
  cost being removed is the Agent work already done when integration fails.

## Open Questions

- Whether a same-wave collision serializes the Tasks automatically or refuses the
  wave and returns to the Supervisor. Refusing is the safer default until the
  TechSpec settles it, because serializing silently changes the execution plan the
  Supervisor authored.
- Whether the collision is computed from declared Task context, from the files
  each Task's Verification names, or from a prior Run's changed-file record. The
  default until answered is the declared context, which is the only source
  available before a Run exists.
