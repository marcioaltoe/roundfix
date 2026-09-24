---
spec: 0158-daemon-verification-and-access-readiness
status: archived
created: 2026-09-24
surfaces: [backend, cli]
archived: "2026-09-24"
source_slug: 0158-daemon-verification-and-access-readiness
---


# Daemon verification and access readiness

Three gaps in what the Daemon proves before and after a Task make a Run spend
turns it did not need to spend.

**The repair turn sees one failure.** A Task's Verification commands run in
order and the attempt stops at the first failure. The single repair turn that
ADR-0038 allows then fixes that failure, the next command fails for a reason
that was already true, and the Task fails with its repair spent. The Agent
never saw the second failure.

**A known-red repository blocks the Task that would fix it.** When a Task's
Verification carries the configured repository command, the Daemon runs it
before the Task and settles the Task failed with "repository not green on
entry". That is right for ordinary work and wrong for a Task written to repair
that very red gate: it cannot start, so the repair has to leave the Spec
workflow.

**An access policy the adapter refuses is found after the Run starts.** Profile
readiness proves the model and effort on a disposable session but never applies
the requested access mode. `.roundfixrc.yml` records the consequence: with
`agent_full_access: true`, codex-acp answered `session/set_mode "full-access"`
with ACP -32602 and every Task in the Run failed before any Agent work. A
runtime with no full-access mode at all accepts the request silently and runs
without it.

This Spec is delivery 3b of the restructured queue, first half: Spec 0122 Core
Features 6 and 7 and Spec 0123 Core Feature 3.

## Project Constraints

- Identifier strategy: not applicable — no new identifier; Task identifiers
  named by an authorization record are the existing Task Graph identifiers.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — local processes and ACP adapters
  only; no credential and no network call. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0011 keeps full access an explicit
  opt-in, ADR-0038 allows one Verification repair, ADR-0056 separates Task and
  Verification Capacity, which independent commands still share, ADR-0057 makes the Daemon the
  only writer of Task status, ADR-0093 checks Spec consistency by citation,
  ADR-0104 accepts on evidence a Spec did not author, ADR-0130 keeps a path
  governed once bounded, ADR-0155 makes the `qa` Task declare the matrix and
  ADR-0156 makes a declared promise name a consuming Task. This Spec's gate is
  bound by ADR-0080, ADR-0091, ADR-0096, ADR-0097 and ADR-0117. All hold.
  This Spec adds ADR-0159 and ADR-0160. Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — the intersection of this Spec's changed
  paths with `internal/speccheck/governed.go` is empty. Source:
  `docs/agents/agent-instructions.md`, `docs/agents/spec-routing.md`.

## Goals

- The single repair turn sees every failure it could have fixed.
- A red repository gate can be repaired inside the Spec workflow, by a Task
  authorized for it and by no other.
- An access policy the runtime cannot honour is refused before any Run starts.

## Core Features

1. **Independent Verification is declared and collected.** A Task that
   declares its Verification commands independent runs every command even after
   one fails, and the repair turn receives every failure with its diagnostics.
   A Task that declares nothing keeps today's order and stops at the first
   failure. Independence is never inferred from the command text. The retry
   ceiling is unchanged.
2. **A named Task may repair a known-red precondition.** The Spec's
   authorization record, as frozen when the Run starts, may name Tasks allowed
   to enter while the repository Verification is red. Such a Task enters, and
   settles completed only when the same repository command and its focused
   repair assertion both pass. Every other Task stays blocked. A Task file or
   Agent edit cannot grant the entry, and a Task whose Verification does not
   carry the configured command verbatim is refused before the Run starts.
3. **Access policy is proven in readiness.** Profile readiness applies the
   requested access mode on the disposable session it already proves, for the
   preferred selection and each fallback, and records the effective policy. A
   runtime without the mode, or an adapter that refuses it, is refused before
   the Run starts, with the failed predicate and its remedy named.

## Non-Goals / Out of Scope

- Settlement semantics, the QA Archive Override and the authoring checks of
  delivery 3b (Spec 0122 Core Features 5 and 8, Spec 0129 Core Features 1, 4
  and 5), which follow in their own Spec together with the authoring guidance
  for the two new declarations.
- Independence between subsets of commands; a Task declares all of its commands
  independent or none.
- Changing the codex sandbox-preset warning that follows a granted mode.

## Success Metrics

1. A Task with two failing independent commands hands both failures to its one
   repair turn; the same Task without the declaration stops after the first.
2. With the repository command red, an authorized repair Task enters and an
   unauthorized Task is blocked with "repository not green on entry".
3. A requested full access on a runtime or adapter that cannot honour it fails
   readiness, and no Run is created.

## Recorded limits

The corrective ceiling of two Tasks was spent on the defects the first pre-PR
review of 2026-09-24 found. The second review found three minor defects,
carried to Spec 0167:

- In independent mode, a command that fails in both the first run and the
  exclusive retry is reported twice to the repair turn, and a command the retry
  shows passing is still carried as failed. Reproduction: a Task declaring
  `verification: independent` with one always-failing command and one command
  that fails temporarily once.
- A named repair Task that completed after the Agent removed or padded the
  repository command in its Task file makes the next `implement` of the Spec
  refuse at planning, because the verbatim check also inspects completed Tasks.
  Reproduction: let the Agent delete the command, settle, then run `implement`
  again with a later Task failed.
- A degraded full-access policy reaches only `profiles validate --json`; the
  text output and Doctor still print `passed`.

## Decisions

- **Declared, never inferred.** Reading independence from shell text would guess
  at setup order; a declaration keeps the author responsible for it.
- **Authority lives in the frozen record.** Only the authorization record read
  at Run start can open the entry, so no Agent turn inside the Run can widen it.
- **Refuse rather than degrade.** A requested access policy that cannot be
  honoured is refused, not silently dropped, because work run under a weaker
  policy than requested fails later and less legibly.

## Acceptance evidence

Each Core Feature requires positive and negative public-contract evidence in the
Task Graph. The negative cases carry the weight: an ungranted Task that enters a
red repository, or an undeclared Task that keeps running after a failure, would
pass the happy path while widening what the Daemon allows.

## Research basis

The repair-turn gap follows from `runVerificationAttempt` in
`internal/daemon/engine.go`, which returns at the first failed command. The
precondition gate is `verifyRepositoryPrecondition` in
`internal/daemon/task_engine.go`. The access-policy failure is recorded in
`.roundfixrc.yml`, measured on codex-acp 1.1.9 on 2026-08-04, and
`applyFullAccess` in `internal/agent/acpx_runner.go` runs only when a working
session starts, never during the disposable proof.

## Technical candidate

The [_techspec.md](_techspec.md) records the implementation map, coverage and
build order.
