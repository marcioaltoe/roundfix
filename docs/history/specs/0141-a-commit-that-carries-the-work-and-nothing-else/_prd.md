---
spec: 0141-a-commit-that-carries-the-work-and-nothing-else
status: archived
created: 2026-09-17
surfaces: [backend, cli, docs]
archived: "2026-09-17"
source_slug: 0141-a-commit-that-carries-the-work-and-nothing-else
---


# A commit that carries the work and nothing else

A Task's settlement commit is supposed to be the work. Today it disagrees with
the work in three ways, and each disagreement is silent.

The Daemon refuses every executable file it was asked to stage, reading the mode
bit as proof of a build artifact. This repository tracks three executable source
files — two Git hooks and a debugging script — so a Task that edits one of them
loses the change at the commit boundary.

A refused path is then only a console line and an event. The Task still settles
completed, and the Run can still reach Clean, with the delivered tree missing
what the work produced. Nothing in the outcome says a path was dropped.

Finally, the QA Report commit takes its baseline when the QA step begins, before
the Daemon runs the repository Verification. Anything that Verification writes
into the worktree afterwards is a candidate for that commit. This repository is
not affected by its own command, whose output is ignored by Git, which is why
the exposure has gone unmeasured.

This Spec makes the commit boundary carry the work and nothing else. It is the
first slice carved from Spec 0122, which keeps the rest.

## Project Constraints

- Identifier strategy: applicable — the drop reasons, the Run and Task outcome vocabulary and the Run Event kinds keep their spelling and meaning; this Spec adds no identifier and renames none. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the commit boundary reads and writes the local worktree and index only, and opens no credential store or network transport. Source: `docs/agents/agent-instructions.md` and `docs/agents/cli.md`.
- Active ADR obligations: applicable — the settlement path is governed by accepted decisions this Spec must preserve. Source: `docs/agents/domain.md` and `docs/agents/spec-routing.md`.
  ADR-0057 applies: the Daemon owns Task status, so a dropped path changes the status the Daemon writes and never asks an Agent to record it.
  ADR-0096 applies: the QA gate proves machine facts before it spends an Agent turn, so the commit baseline moves without changing what the mechanical stage decides.
  ADR-0117 applies: a defect is checked by the stage that can produce it, so the drop is refused at the commit boundary that creates it.
  ADR-0035 applies: a Spec Root outside the repository keeps its artifacts uncommitted, so a path external to the repository is refused by design and is not a dropped work output.
  ADR-0036 applies: review artifacts are committed in their own docs commit, which this Spec neither creates nor changes.
  ADR-0029 applies: review artifacts live with the Spec, and this Spec moves no artifact between locations.
  ADR-0142 applies: head-bound Review Source Evidence decides the watch outcome, and this Spec changes no watch, review or evidence binding.
  ADR-0104 applies: a Spec accepts on evidence it did not author, which here is this repository's three tracked executable files.
  ADR-0156 applies: this Spec declares its Success Metrics and API Contracts as numbered units and names each in a Task.
  ADR-0080 applies: a QA row that no environment can run is recorded as environment-blocked and never as a product failure, which is how this Spec's gate records a Pull Request that does not exist yet.
  ADR-0091 applies: the QA gate is a Task node of its own type, which is the terminal Task this Spec authors.
  ADR-0093 applies: Spec consistency is checked by citation and never by inference, and this Spec changes no consistency rule.
  ADR-0097 applies: a QA row carries forward only on declared unmoved evidence, and this Spec moves no row's evidence.
  ADR-0155 applies: the `qa` Task declares the gate's matrix, so this Spec's terminal gate covers what its Requirements name.
  ADR-0157 applies: tracked source is work output whatever its mode, and a drop stops the Task from settling completed.
- Tooling authority: applicable — no protected tooling mutation proposed or authorized. The repair is ordinary source in `internal/daemon`, which is not a Governed Path, and its test files are ordinary too. Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`, `docs/agents/specific-repository.md`.

## Goals

- A Task that changes a tracked file keeps that change in its settlement commit,
  whatever the file's mode.
- A path the work produced that cannot be staged stops the Task from settling
  completed, and the outcome names the paths and the cause.
- The QA Report commit carries the report, its evidence and the Agent's own
  writes, and not what the repository Verification wrote.
- Every other reason a path is refused today keeps its meaning, and an untracked
  executable stays refused.

## User Stories

1. As a Supervisor whose Task edits a Git hook, I want the change to reach the
   settlement commit, so that the delivered tree contains the work the Task
   reported.
2. As a Supervisor reading a Run outcome, I want a dropped path named with its
   cause, so that I learn about it when it happens rather than when something
   later breaks.
3. As a Supervisor rerunning a Spec, I want a Run that dropped required work to
   refuse Clean, so that a rerun is what recovers it.
4. As a Supervisor whose repository Verification writes into the worktree, I
   want those files kept out of the QA Report commit, so that the report commit
   stays the report.

## Core Features

1. **A tracked path is stageable whatever its mode.**
   - A path Git already tracks is staged, including an executable regular file.
   - An executable file that Git does not track is still refused, with the
     reason and mode it reports today.
2. **A drop is an outcome, not a warning.**
   - When any path the work produced is refused, the Task does not settle
     completed.
   - The Task and Run outcome name each dropped path and its cause.
   - The existing console line and Run Event stay; they stop being the only
     record.
3. **The QA Report commit excludes what the Verification wrote.**
   - The paths that appear while the Daemon's repository Verification runs are
     not carried by the QA Report commit.
   - The report the mechanical stage seeded, its evidence, and what the QA Agent
     writes afterwards are carried exactly as today.
4. **Every other refusal keeps its meaning.** A path external to the repository,
   a path that crosses a symbolic link, and a path absent from both worktree and
   index keep their reasons, their event payloads and their console lines.

## User Experience

A Task edits a tracked hook and the hook's change appears in the settlement
commit. A Task produces a file the Daemon cannot stage, and the Run stops with
an outcome that names the path and why: the Supervisor reruns instead of
discovering the gap in a later Spec. A QA step whose repository Verification
writes a tracked file still commits only the report and its evidence.

## Non-Goals / Out of Scope

- Running the repository Verification against the integrated candidate before
  declaring Clean. That is Spec 0122's own Core Feature 3 and needs its commit
  boundary decided first.
- One declared-acceptance eligibility policy across settlement, the derived QA
  command and archive; the QA Archive Override; the known-red precondition
  repair; and independent Verification groups. Those remain with Spec 0122.
- Changing which artifacts an external Spec Root keeps uncommitted.
- Changing Run Event payload shapes or the projection that reads them. The
  `daemon.verification` payload disagreement is captured separately.
- Guessing whether an untracked file is a build artifact by reading its content.
- Changing the QA Report's shape, its verdict rules or its naming.

## Declared intentional breaks

- A Run that today reaches Clean while dropping a path the work produced no
  longer reaches Clean. That is the point of the change, and it is visible the
  first time a Task edits an untracked executable.
- A Task that deliberately produces an untracked executable must add it to the
  repository before the Daemon will carry it.

## Regression locks

- The three existing drop reasons keep their text, their payloads and their
  console lines.
- A file Git ignores is still never staged.
- The QA Report commit still carries the report and its evidence, including a
  report that was already dirty when the step began.
- Task status stays Daemon-written, and the mechanical stage keeps deciding what
  it decides today.

## Acceptance evidence

At least one acceptance row rests on evidence this Spec did not author: the
three executable files this repository already tracks — two Git hooks and a
debugging script — and the recorded finding that the QA Report commit can carry
what a repository Verification wrote. The replay edits one tracked hook in a
fixture Run and shows the change reaching the settlement commit, where today it
is dropped. Where the outside artifact cannot be read, the row records that
reason and does not block.

## Success Metrics

1. A fixture Task that edits a tracked executable settles with that change in
   its commit, and the same Task on the unchanged tree loses it.
2. A fixture Task whose produced path cannot be staged does not settle
   completed, and its outcome names the path and the cause.
3. A fixture whose repository Verification writes a tracked file produces a QA
   Report commit that does not carry that file.

## Decisions

- **Tracked beats mode.** A path Git tracks is work output whatever its mode; an
  untracked executable stays refused. See ADR-0157.
- **A drop fails the Task.** The console line and the Run Event stay, but they
  stop being the whole record. See ADR-0157.
- **The baseline moves, the report does not.** The QA Report commit excludes the
  paths that appear while the repository Verification runs, and keeps everything
  the mechanical stage and the Agent write.
- **Slice, not portfolio.** Spec 0122 carries nine Core Features across five
  families; this Spec takes the commit-content family and leaves the rest, the
  way Spec 0140 was carved from Spec 0129.

## Research basis

**Secondbrain.** Consulted before authoring, index first, then two queries: one
on committing what a task produced rather than build output, and one on failing
loudly instead of warning in automation. The results were dominated by mirrors
of this repository, which are references rather than independent knowledge. The
one non-mirror source, a process-automation chapter in the wiki, argues that an
automated step which continues after a partial failure moves the cost to whoever
finds it later; it supports making a drop terminal rather than advisory, and
changed no other decision.

The repository's pending Inbox Entries were read before authoring. None is a
source for this Spec; the nearest concerns a Run Event payload the projection
cannot read, which this Spec deliberately leaves alone.

**Exa MCP.** Consultation was attempted, and no Exa MCP tool was available in
this session. No external source was read, so no external validation is claimed.

**Local measurement.** The three tracked executable files, the filter that
refuses them, and the QA baseline taken before the repository Verification were
each read in this repository's source before this PRD was written.

## Open Questions

- Should a drop fail the Task immediately, or let the Task finish and settle the
  Run as Unresolved with the drop named? Default until the maintainer decides
  otherwise: fail the Task, because a Task whose output is incomplete has not
  done its work, and the Run's other Tasks then skip on the dependency they
  never received.
