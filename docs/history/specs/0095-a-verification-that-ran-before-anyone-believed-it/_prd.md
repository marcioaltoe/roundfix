---
spec: 0095-a-verification-that-ran-before-anyone-believed-it
status: archived
created: 2026-08-12
surfaces: [backend, cli, docs]
archived: "2026-08-14"
source_slug: 0095-a-verification-that-ran-before-anyone-believed-it
---


# A Verification that ran before anyone believed it

A Task's Verification is checked for form and never for execution, and the gap is
the most expensive one measured across the fleet. Six defects in one night in a
repository this Spec did not build were pure shell semantics — a count that prints
a filename, a pipeline that exits with the status of its last stage, a platform
that pads a number with spaces — and four of them failed work that was correct.
A twelve-Task graph elsewhere came within one runner's exit code of being entirely
vacuous, saved by the tool's luck rather than by design. Meanwhile the rule that a
Verification passes only by exiting zero is enforced by the Daemon and written
nowhere, so a natural-looking assertion whose success is an empty result fails its
Task for doing the right thing.

## Project Constraints

- Identifier strategy: applicable — Verification, Task, Task Graph, and the
  consistency-check code vocabulary are glossary terms this Spec extends, and each
  new check code is coined vocabulary the glossary must own. The closing node
  checks whether the work introduced or changed a term. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — no network boundary, credential, or
  request is created or read. The work is static analysis of authored commands and
  their execution in a disposable local checkout. Source:
  `docs/agents/specific-repository.md`.
- Active ADR obligations: applicable — ADR-0093 checks Spec consistency by
  citation rather than inference, which bounds what a new static detector may
  conclude; ADR-0117 places a check with the stage that can produce its defect,
  which is this Spec's whole thesis moved one stage earlier — from Run time to
  authoring. ADR-0094 makes the consistency check artifact-presence aware, so a
  detector must skip rather than fail where its artifact is absent. ADR-0111 makes
  an unobserved Verification unknown rather than a verdict, which is what an
  execution mode must report when it could not run a command at all. ADR-0124
  makes authoring and Run time share one prober, so the earlier check cannot
  approve what the later one refuses; ADR-0014 keeps the Daemon the owner of task
  verification and settlement, which sharing the prober moves code around rather
  than changes; ADR-0081 keeps a grant drawn around the authorized cause with its
  sanctioned regeneration following as fallout, which is how the skill edit and
  its mirror are accounted; ADR-0091 makes the QA gate a Task node of its own type
  and ADR-0104 makes a Spec accept on evidence it did not author, which together
  shape the closing node. The decisions built on the Daemon's ownership are
  accounted and unchanged: ADR-0020 ranks a parsed prompt result above the acpx
  exit code, ADR-0038 allows one Verification repair per Task, and ADR-0057 makes
  the Daemon the sole writer of Task status — all three govern a Run, and the
  authoring caller this Spec adds starts no Run, so none of them applies to it;
  ADR-0096 proves machine facts before the gate spends an Agent turn, which is the
  same economy one stage earlier and is followed rather than altered. ADR-0056
  separates Task Capacity from Verification Capacity, which is Run bookkeeping the
  Daemon keeps around the shared prober rather than inside it, so the authoring
  caller acquires neither. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — the authoring contract in the task-authoring
  skill gains the exit-zero rule. Express maintainer authorization:
  `docs/workflow/authorizations/2026-08-12-the-authoring-and-baseline-corrections.md`,
  granted 2026-08-12. Bounded files: `.agents/skills/write-tasks/SKILL.md`.
  Source: `docs/agents/agent-instructions.md`.

## Goals

1. A Verification command is proven to run before a Daemon spends a turn on it.
2. A command whose success condition is an empty result is authored correctly the
   first time, because the rule is written where the author reads.
3. A command that cannot fail, and a command that fails when the work is right,
   are both refused at authoring.
4. A Task can declare the file it creates without the reference check refusing it.

## User Stories

1. As a Supervisor authoring a Task Graph, I want each Verification line executed
   before a Run, so that a shell-semantics slip costs a second instead of a Run.
2. As a Supervisor writing a Verification whose success is an empty result, I want
   the authoring contract to state the working form, so that I do not author a
   gate that fails the Task for succeeding.
3. As a Supervisor authoring a Task that creates a file, I want to declare that
   file as its output, so that the reference check does not refuse a path that
   cannot exist yet.
4. As a Supervisor, I want a Verification depending on undeclared environment
   state refused at authoring, so that a false red does not consume a repair turn.

## Core Features

1. **Verification execution at authoring.** A mode executes each authored
   Verification line in a disposable checkout and reports its exit status, so a
   command that cannot run, cannot fail, or fails on correct work is visible
   before any Run.
2. **A refusal for a reversed exit condition.** Known forms whose exit status
   inverts the author's intent are refused statically, with the working form
   named.
3. **A refusal for a non-hermetic Verification.** A command referencing an
   environment variable the repository does not declare, guarding itself on one,
   or depending on state outside the repository is refused, since all three
   produced false reds in measured sessions.
4. **The vacuity refusal returns.** The detector that refuses a command already
   satisfied before any work runs is re-enabled, its ten measured findings
   accounted, and its cost paid at authoring rather than at Run time.
5. **A Task may declare what it creates.** The context contract accepts a
   declared output whose existence is not required, so a Task that creates a file
   can name it without the reference check refusing the graph.
6. **The exit-zero rule is written where it is read.** The authoring contract
   states that a Verification command passes only by exiting zero, beside the
   vacuity rule, with the working forms as the worked answer.

## User Experience

The execution mode reports one line per Verification command: the command, its
exit status, and whether that status matches what a passing Task would produce.
A static refusal names the command, the form it matched, and the working
replacement.

## Non-Goals / Out of Scope

- Changing how the Daemon executes Verification at Run time, or how it settles a
  Task from the result.
- Deciding a Task's Verification content on the author's behalf.
- Any change to the QA gate's own verification or its verdict semantics.
- A general shell linter. The refusals cover measured forms, not shell correctness
  in general.
- Changing the authoring skill beyond the one bounded file the grant names.

## Success Metrics

- The six shell-semantics defects measured in one night in a repository this Spec
  did not build are each caught by the execution mode or a static refusal. Source:
  a session on 2026-08-07/08 whose defects were recorded outside this Spec.
- A graph whose Verification commands all pass on the pre-work tree is refused as
  vacuous rather than dispatched.
- A Task declaring a created file passes the reference check.
- The three measured non-hermetic forms are each refused.
- An authored graph that would have failed a correct Task on a reversed exit
  condition is refused before dispatch.

## Decisions

- The exit-zero rule is written into the authoring contract rather than inferred
  from the Daemon's behavior, because the Daemon's behavior is correct and its
  silence is what cost the Run.

## Open Questions

- Whether the execution mode runs against the pre-work tree, proving a command can
  fail, or the post-work tree, proving it can pass, or offers both. Both as
  distinct modes is the default until answered, because the vacuity rule and the
  exit-zero rule ask opposite questions of the same command.
- Whether re-enabling the vacuity refusal is gated on the gate-economics Spec
  landing first. The default is no: its ten measured findings each cost a daemon
  cycle today, and the refusal is one line in a staged registry.
