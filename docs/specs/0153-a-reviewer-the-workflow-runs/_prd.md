---
spec: 0153-a-reviewer-the-workflow-runs
status: active
created: 2026-09-21
surfaces: [backend, cli, docs]
---

# A reviewer the workflow runs

Roundfix can say which reviewer a repository wants. It cannot ask one anything.

`internal/config/config.go` resolves a Pre-PR Review Policy, validates it against
`codex, claude, coderabbit, none`, defaults to `codex`, and records where the
value came from. `doctor` reports that policy's health. Then the trail ends:
measured by searching the tree, the only consumer of `Config.PrePRReview` is
`doctor.go`. Nothing calls a reviewer.

So the pre-PR review this project relies on is performed by whoever is driving —
in practice, a person or an agent typing `codex review --base main` in a
terminal and reading the answer. That review has caught real defects, including
several this queue shipped because of it. None of it is recorded anywhere the
workflow can read, none of it is bound to a candidate commit, and nothing stops
a changed candidate from inheriting an earlier answer.

This Spec is carved from Spec 0126 Core Feature 2, and takes only the Codex
provider. Core Feature 1 is already delivered.

## Project Constraints

- Identifier strategy: applicable — evidence names the repository, the base and
  the head commit exactly, so a review is bound to the candidate it examined.
  No identity is minted. Source: `docs/agents/domain.md`.
- Authentication and HTTP: applicable — the reviewer runs through the existing
  Agent Runtime, which owns its own transport and credentials. This Spec adds no
  credential handling and opens no connection of its own. Source:
  `docs/agents/cli.md` and `docs/agents/autonomous-work.md`.
- Active ADR obligations: applicable — ADR-0093 checks Spec consistency by
  citation, ADR-0104 accepts on evidence a Spec did not author, ADR-0130 keeps a
  path governed once bounded, ADR-0155 makes the `qa` Task declare the matrix
  and ADR-0156 makes a declared promise name a consuming Task. This Spec's own
  gate is bound by ADR-0080, which keeps an environment-blocked row distinct
  from a failure, ADR-0091, which makes the gate a Task node of its own type,
  ADR-0096, which proves machine facts before spending an agent turn, ADR-0097,
  which carries a row forward only on declared unmoved evidence, and ADR-0117,
  which checks a defect at the stage that can produce it. All hold. Source:
  `docs/agents/domain.md`.
- Tooling authority: applicable — a new command changes the public CLI surface
  the shipped skill documents. Express maintainer authorization: the standing
  grant of 2026-09-18 for the skill files and the standing grant of 2026-09-21
  for governed paths a slice needs, both recorded in
  [_authorization.md](_authorization.md); bounded files:
  `.agents/skills/roundfix/SKILL.md`, `skills/roundfix/SKILL.md`,
  `internal/cli/cli_test.go`. Sanctioned regeneration: `make skills-sync`. The
  bounded set is the intersection of this Spec's changed paths with the literal
  set in `internal/speccheck/governed.go`, computed rather than predicted.
  Source: `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`,
  `docs/agents/specific-repository.md`.

## Goals

- The configured reviewer is actually asked, by the workflow, not by a person.
- What it examined is recorded against the exact candidate it examined.
- An unselected provider costs nothing: no call, no readiness probe, no account.

## Core Features

1. **The workflow asks Codex.** When the policy resolves to `codex`, a command
   runs a read-only reviewer session over the current candidate through the
   existing Agent Runtime, and reports its outcome through the process exit
   status.
2. **Evidence names the candidate.** The record carries the repository, the base
   and head commits, the effective provider with the source that selected it,
   the execution outcome and the findings returned. A record that names one head
   says nothing about another.
3. **Explicit `none` costs nothing.** It performs no reviewer call and no
   readiness probe, records a configured omission for the candidate, and exits
   zero. It is a decision, not a failure.
4. **A blocked review is never a pass and never `none`.** A runtime failure, a
   timeout, a transport anomaly, a non-zero adapter exit, empty output, or
   output the command cannot read each block the selected mode with the reason
   named. None of them falls back to another provider and none becomes a
   configured omission. A fallback activates only for a selection that failed to
   start before the prompt was sent.
5. **An unimplemented provider refuses out loud.** `claude` and `coderabbit` are
   valid policy values this slice does not execute. Selecting one refuses with
   that reason rather than silently reviewing nothing.

## Non-Goals / Out of Scope

- The Claude and CodeRabbit providers. Core Feature 2 of Spec 0126 covers all
  three; this slice delivers one and refuses the others explicitly.
- Routing findings into bounded corrective work, which is Spec 0126 Core Feature
  5.
- Gating publication or merge on the evidence, which is Spec 0126 Core Feature
  6. This slice produces the record; nothing yet consumes it.
- Changing the policy resolution, its defaults, its precedence or its `doctor`
  check, all of which already work.

## Success Metrics

1. With the policy at `codex` and a stubbed runtime, the command runs a reviewer
   session and writes a record naming the repository, base and head.
2. With the policy at `none`, the command writes a configured-omission record,
   exits zero, and the stubbed runtime records no call and no probe.
3. A runtime failure, a timeout and unreadable output each block with a named
   reason, and none of the three produces a passing or omitted record.

## Decisions

- **Reuse the Agent Runtime rather than shell out.** The Daemon already owns
  running an agent session through `agent.Runner`. A second way to start a
  reviewer would be a second thing that can disagree about runtimes, models and
  credentials.
- **Record, do not yet gate.** Producing evidence and consuming it are separate
  changes with separate failure modes. This slice can be wrong without blocking
  anyone's merge, which is the right order to find out.
- **Refuse the providers this slice does not run.** A policy value that
  silently reviews nothing is worse than an error, because the repository
  believes it is covered.
- **Hand the reviewer its evidence; do not send it looking.** A first attempt
  gave the session commit identifiers and read-only file tools and asked it to
  inspect the range. Three review rounds showed why that fails: with write
  tools it could edit the candidate, without them it could read nothing, and
  with file tools but no diff it could answer `No findings` having never seen
  what changed. The command computes the diff and puts it in the prompt, so the
  reviewer judges content it was given rather than content it must go find.
- **Only an explicit clean answer is clean.** A timeout, a transport anomaly, a
  non-zero adapter exit, empty output and unparseable output are each a reason
  the review did not happen. Treating any of them as anything but blocked lets
  a broken reviewer look like a strict one.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in the
Task Graph. The negative cases carry the weight: a blocked review that recorded
a pass, or an unimplemented provider that recorded an omission, would satisfy
the happy path while making the record a lie.

## Research basis

Measured on this repository before authoring. `Config.PrePRReview` is consumed
only by `internal/cli/doctor.go`; a search for any execution of a review found
none. The policy's four values are validated in `internal/config/config.go:1676`
and default to `codex`. The reviewer will run through `agent.Runner`, whose
`Probe`, `Run` and `EndSession` methods the Daemon already uses, so the tests
stub that interface and the gate spends no paid API use.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
