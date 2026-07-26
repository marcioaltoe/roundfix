---
spec: 0042-verification-capacity-and-daemon-task-settlement
status: active
created: 2026-07-18
surfaces: [cli, backend, frontend, docs]
---

# Verification Capacity and Daemon Task Settlement

Spec Runs currently use `worktree.concurrency` to limit an entire Task
lifecycle, so increasing Task Capacity also lets multiple repository-wide
Verification suites run at once. In Vortex this coupled two useful concurrent
Agent implementations to two simultaneous integration suites, exhausting
local listeners and setup capacity. The same Run also let Agents mark Tasks
`failed` from their own gate runs, which bypassed Daemon Verification,
Verification Feedback, and repair. This Spec absorbed the incident evidence;
the original report remains available through Git history.

This feature separates Verification Capacity from Task Capacity, makes the
Daemon the sole Task-status writer during an Implement Run, exposes each
Task's Verification phase, and gives explicitly declared temporary capacity
failures one bounded exclusive retry. Projects keep their full Verification
commands unchanged while concurrent Agents can continue preparing independent
Tasks.

## Project Constraints

- Identifier strategy: not applicable — capacity, retry, and Task settlement
  reuse existing Run, Task, Batch, and Verification attempt identities and
  create no project-owned Internal Identifier. Source:
  `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the feature is confined to local
  configuration, Agent Sessions, Verification processes, Run Events, and TUI
  projection. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0014 preserves Daemon-run
  Verification, ADR-0025 preserves Task Graph readiness, ADR-0038 preserves one
  Agent repair, ADR-0051 routes each Task through its Task Type-selected Agent
  Session, and ADR-0056 plus ADR-0057 define capacity and Task-status ownership.
  Source: `docs/agents/domain.md`.
- Tooling authority: applicable — on 2026-07-26, the maintainer expressly
  authorizes changes to exactly `.agents/skills/roundfix/SKILL.md`,
  `skills/roundfix/SKILL.md`, `.agents/skills/implement-task/SKILL.md`, and
  `skills/implement-task/SKILL.md`; no other protected tooling mutation is
  authorized. Source: `docs/agents/agent-instructions.md`.

## Goals

- Let multiple independent Tasks perform Agent work concurrently without
  overlapping repository Verification by default.
- Make the Daemon authoritative for every Task status transition during an
  Implement Run, including `in_progress`, `completed`, and `failed`.
- Ensure every implementation-ready Task reaches Daemon Verification even if
  its Agent wrote a premature terminal status.
- Distinguish Agent work, Waiting for Verification, active Verification,
  Verification Feedback, and settlement in durable Run evidence and the Live
  Run View.
- Retry only a positively identified Temporary Verification Failure once,
  under exclusive Verification Capacity, without consuming the existing Agent
  repair.
- Preserve the full project gate, failed Task Worktrees, dependency blocking,
  independent Task continuation, and existing terminal Run outcomes.

## User Stories

1. As a maintainer of a repository with a heavy integration suite, I want
   Verification Capacity independent from Task Capacity, so concurrent Agents
   do not exhaust the suite's shared local resources.
2. As a developer optimizing a Spec Run, I want Task Capacity `2` and
   Verification Capacity `1`, so implementation overlaps while full gates run
   sequentially.
3. As a Supervisor following an Implement Run, I want to see whether each Task
   is Agent working, Waiting for Verification, or Verifying, so a healthy
   queue is not mistaken for a stalled Agent.
4. As a maintainer relying on the Daemon-owned Verification contract, I want
   Agent-authored `completed` or `failed` statuses ignored as verdicts, so no
   Task skips the authoritative gate.
5. As a Task Agent, I want to hand back implementation-ready work and receive
   only Daemon Verification Feedback when the gate fails, so I do not duplicate
   the full suite or decide the terminal outcome.
6. As a project author, I want an explicit Temporary Verification Failure
   signal, so known capacity failures can be retried without fragile log-text
   matching.
7. As a developer diagnosing an unstable gate, I want the temporary failure,
   exclusive retry, diagnostic artifacts, and final classification journaled,
   so a retry cannot silently hide a deterministic failure.
8. As a user stopping a Run, I want a Task waiting for Verification Capacity
   to honor cancellation without leaking work or inventing a terminal Task
   verdict.

## Core Features

1. **Independent capacities.** `worktree.concurrency` remains the Task
   Capacity for concurrently executing Task Worktree lifecycles. A new
   `verification.concurrency` Project/User Config key controls the number of
   Task Verification attempts that may run concurrently within one Implement
   Run. It defaults to `1`, accepts only positive integers, and follows the
   existing built-in, User Config, and Project Config precedence.
2. **Per-Run Verification queue.** Every Task Verification attempt enters one
   cancellation-aware queue owned by its Task cycle. Capacity is held for the
   complete sequence of commands in one attempt and released before Agent
   Verification Feedback. A later attempt queues and acquires capacity again.
3. **Daemon-owned Task status.** At Task start the Daemon writes
   `in_progress`. It alone writes `completed` or `failed` after the applicable
   Agent, Verification, repair, infrastructure, or settlement outcome.
   Agent-authored status changes are normalized back to the Daemon's current
   status and never act as a verdict.
4. **Authoritative Verification handoff.** Initial and repair Agents must not
   run a Task's declared `## Verification` commands or claim a terminal Task
   outcome. Successful Agent handoff always proceeds to the Daemon's verbatim
   Verification sequence. Focused implementation checks remain allowed and
   belong in the Task Result evidence.
5. **Existing Agent repair preserved.** A deterministic command failure still
   produces bounded diagnostics, one Verification Feedback turn in the same
   Task Agent Session, and one final Daemon Verification attempt. Capacity is
   not held while the Agent repairs.
6. **Explicit temporary-failure protocol.** Exit code `75` from a project
   Verification command is the only Temporary Verification Failure signal.
   Roundfix never infers this classification from command output, timing,
   package names, ports, or error text. All other non-zero command exits retain
   the deterministic failure and Agent-repair behavior.
7. **One exclusive retry.** A Task receives at most one exclusive retry across
   its complete Verification lifecycle. The retry waits until no other Task
   Verification is running in the Run, executes alone, retains separate
   diagnostics, and does not consume or add an Agent repair. A repeated
   Temporary Verification Failure settles the Task failed; a deterministic
   failure from the retry may still use the one existing Agent repair when it
   has not been used.
8. **Task-phase observability.** The Run Event Stream records effective Task
   and Verification capacities, Waiting for Verification before each
   acquisition, Verification start, temporary classification, exclusive
   retry, Verification Feedback, verdict, and settlement. Attach replays the
   same evidence and the Live Run View shows per-Task phase labels rather than
   mapping every `in_progress` Task to `Executing`.
9. **Cancellation and recovery.** A Stop Request or context cancellation while
   waiting exits the queue promptly, starts no Verification command, and
   leaves the Task resumable rather than assigning a false terminal verdict.
   Failed Task Worktrees and existing Settle Command recovery remain
   unchanged.
10. **Contract-aligned guidance.** Agent prompts, the repository-owned
    `implement-task` workflow, configuration and command documentation, and
    the canonical/embedded Roundfix Skill describe the same Daemon-owned
    status and Verification behavior.

## User Experience

A repository that wants two concurrent Task lifecycles but only one full gate
configures:

```yaml
worktree:
  concurrency: 2

verification:
  concurrency: 1
```

Omitting `verification.concurrency` uses the safe built-in value `1`.
Configuration values below `1` fail strict validation before a Run is created.
No new Implement flag or Roundfix command exit code is introduced.

The Run Event Stream emits a Waiting for Verification record before a Task
attempt acquires capacity, even when capacity is immediately available. This
gives every attempt the same auditable ordering: waiting, started, command
results, and verdict. The Live Run View reports both effective capacities and
uses the text labels `Agent working`, `Waiting for Verification`, and
`Verifying`; meaning never depends on color alone.

A project wrapper that positively recognizes its own temporary infrastructure
condition returns exit code `75`. Roundfix retains that failure's diagnostics,
records that an exclusive retry is waiting, and runs the retry alone. A
successful retry continues normally. A second exit `75` records exhaustion
and fails the Task. A non-`75` failure follows the existing Verification
Feedback contract.

## Non-Goals / Out of Scope

- Weakening, removing, splitting, or increasing timeouts for project
  Verification commands.
- Parsing logs or maintaining a Roundfix allowlist of framework-specific
  timeout, port, container, database, or listener errors.
- Automatically converting arbitrary non-zero exits into retryable failures.
- Coordinating Verification across separate Roundfix processes, repositories,
  Runs, CI jobs, or commands started outside Roundfix.
- Protecting shared databases or other resources mutated during Agent work or
  Worktree Bootstrap; projects still configure safe Task Capacity and isolated
  bootstrap behavior.
- Adding configurable retry counts, a second Agent repair, unbounded retries,
  backoff policy, or a new terminal Run outcome.
- Changing review Batch Verification, the Settle Command contract, Task Graph
  dependency semantics, Task Worktree preservation, or Run integration.
- Treating exit code `75` as proof that a project classified its failure
  correctly; the project owns the wrapper that emits the signal.
- Replacing Task files as the durable owner of Task status; this feature
  changes who may write that field during an Implement Run.

## Success Metrics

- With Task Capacity `2` and Verification Capacity `1`, two Task Agents can
  overlap while the observed maximum number of active Task Verification
  attempts is exactly `1`.
- With Verification Capacity `2`, two ready Verification attempts can overlap;
  values `0` and negative values fail before Run creation.
- Every Task Verification attempt records Waiting for Verification before its
  first started event, and Attach replays the same ordering.
- Zero Agent-authored terminal statuses can bypass Daemon Verification or
  directly settle an Implement Task.
- A deterministic first failure produces exactly one Agent repair and at most
  two numbered Verification attempts, preserving the existing contract.
- Exit code `75` produces at most one separately identified exclusive retry
  per Task; no log content can trigger that retry.
- Cancellation while queued starts zero Verification commands and leaves no
  blocked goroutine under race-enabled tests.
- The Live Run View reports both capacities and exact per-Task working,
  waiting, verifying, and terminal labels in styled and no-color modes.
- `go test -race ./...`, the full `make verify` gate, and Roundfix Skill
  synchronization all pass.

## Decisions

- Deliver the full remediation: independent capacity, Daemon-only Task status,
  observable waiting, and one positively classified exclusive retry.
- Default Verification Capacity to `1`; do not inherit Task Capacity.
- Use fixed exit code `75` as the project-authored Temporary Verification
  Failure signal instead of configuration allowlists or marker artifacts.
- Give each Task at most one exclusive retry across both numbered Verification
  attempts; the retry does not consume the existing Agent repair.
- Keep Verification Capacity scoped to one Implement Run rather than adding a
  machine-wide coordination service.
- Preserve the full Verification gate and diagnose resource ownership at its
  source rather than weakening tests or extending timeouts.
- Keep backend, frontend, and documentation Tasks in this Task Graph. ADR-0051
  governs per-Task Agent Session selection by Task Type, superseding older
  one-Agent-per-Run and separate-frontend-Spec guidance; Task 07 aligns that
  guidance and the Agent Session glossary.
- Isolate all protected Skill changes in one tooling-only Task bounded to the
  four expressly authorized canonical/generated `roundfix` and
  `implement-task` `SKILL.md` files.
- See [ADR-0056](../../adr/0056-spec-runs-separate-task-and-verification-capacity.md).
- See [ADR-0057](../../adr/0057-daemon-exclusively-owns-implement-task-status.md).

## Open Questions

None.
