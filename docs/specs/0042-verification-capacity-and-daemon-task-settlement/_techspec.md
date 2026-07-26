---
spec: 0042-verification-capacity-and-daemon-task-settlement
prd: _prd.md
created: 2026-07-18
---

# Verification Capacity and Daemon Task Settlement — Technical Spec

## Executive Summary

Add a Task-cycle-owned Verification gate whose capacity is configured
independently from the existing ready-Task scheduler. Normal Verification
attempts share that gate; a positively classified Temporary Verification
Failure releases its normal permit and may acquire the whole gate once for an
exclusive retry. The Daemon writes `in_progress` before Agent work and remains
the only Task-status writer through final settlement, so a reloaded
Agent-authored terminal value can no longer bypass Verification. Additive Run
Events make the queue, attempt, retry, repair, and effective capacities
durable; Attach and the Live Run View derive per-Task phases from that journal.
The accepted trade-off is per-Run rather than machine-wide coordination: it
solves contention caused by one Implement Run without introducing a daemon
service or claiming control over external processes.

## Project Constraints

- Identifier strategy: not applicable — the design extends existing Run Event
  and Verification-attempt records without minting a new project-owned
  identity. Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — all new behavior remains inside
  local configuration, Agent Session, Verification, Run Event, and TUI
  boundaries. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0014, ADR-0025, ADR-0038, ADR-0051,
  ADR-0056, and ADR-0057 bind Daemon Verification, Task readiness, bounded
  repair, per-Task Agent Selection, independent capacity, and exclusive Task
  status authorship. Source: `docs/agents/domain.md`.
- Tooling authority: applicable — on 2026-07-26, the maintainer expressly
  authorizes changes to exactly `.agents/skills/roundfix/SKILL.md`,
  `skills/roundfix/SKILL.md`, `.agents/skills/implement-task/SKILL.md`, and
  `skills/implement-task/SKILL.md`; no other protected tooling mutation is
  authorized. Source: `docs/agents/agent-instructions.md`.

## System Architecture

- **`internal/config`** gains a top-level `Verification` configuration section
  with positive integer `concurrency`, built-in default `1`, and the existing
  built-in → User Config → Project Config overlay. The existing
  `defaults.verification` review command remains a separate string contract.
- **`internal/cli`** passes both effective capacities into `daemon.TaskPlan`,
  includes them in initial and attached `LiveRunView` values, and preserves all
  Implement arguments, stdout, stderr, and exit codes.
- **`internal/daemon`** creates one cancellation-aware Verification gate per
  `TaskCycle`. Task Worktree copies of `TaskPlan` share the same gate. The
  Daemon owns Task status, shared/exclusive acquisition, exit-code
  classification, bounded retry state, diagnostic artifacts, Agent repair,
  and final settlement.
- **`internal/agent`** changes Task and Verification Feedback prompts from
  terminal-status authorship to implementation-ready handoff. Agents may run
  focused checks and edit the Result section but do not execute declared Task
  Verification or edit Task status.
- **`internal/runevent`** adds Verification phases and additive payload fields
  under the existing `daemon.verification` kind and `roundfix-events/v1`
  projection. No journal schema or store migration is required.
- **`internal/tui`** folds per-Task `daemon.task` and
  `daemon.verification` events into explicit Agent working, Waiting for
  Verification, Verifying, and terminal labels. It shows Task Capacity and
  Verification Capacity in the existing Run header without adding navigation
  or changing panel positions.
- **Docs and skills** align configuration, command guidance,
  `implement-task`, and the canonical/embedded Roundfix Skill with the new
  ownership boundary. The full project gate remains mandatory and is run by
  the Daemon for each Task.

The Task lifecycle becomes:

```text
Daemon status=in_progress -> Agent working -> Waiting for Verification
    -> shared Verification attempt N
       pass -> Daemon status=completed -> commit
       deterministic fail -> release -> Verification Feedback Agent
          -> Waiting for Verification -> shared attempt 2 -> settle
       exit 75 -> release -> Waiting for exclusive Verification
          -> exclusive retry of attempt N
             pass -> continue as passed
             deterministic fail -> optional existing Agent repair
             exit 75 again -> Daemon status=failed
```

## Implementation Design

### Interfaces

Configuration and Task-cycle input:

```go
type Verification struct {
    Concurrency int
}

type TaskPlan struct {
    // existing fields omitted
    Concurrency             int // Task Capacity
    VerificationConcurrency int
}
```

The cycle creates and internally threads one gate; it is not a public package
or a process-global singleton:

```go
type verificationMode uint8

const (
    verificationShared verificationMode = iota
    verificationExclusive
)

type verificationGate interface {
    Acquire(context.Context, verificationMode) (release func(), err error)
}
```

The concrete gate uses one mutex, active shared count, exclusive-active flag,
exclusive-waiter count, and a close-and-replace notification channel. Shared
acquisition succeeds only when capacity is available and no exclusive waiter
is queued; exclusive acquisition succeeds atomically only when no attempt is
active. Blocking new shared entrants behind an exclusive waiter prevents
starvation. Cancellation selects on `ctx.Done()` without spawning a goroutine.
Release updates state once and notifies all waiters; tests must prove no permit
leak rather than expose gate internals as production test hooks.

Temporary-failure classification wraps the existing typed command error so
diagnostic and exit information remain inspectable:

```go
const TemporaryVerificationExitCode = 75

type TemporaryVerificationFailureError struct {
    CommandFailure *VerificationCommandError
}

func (*TemporaryVerificationFailureError) Error() string
func (*TemporaryVerificationFailureError) Unwrap() error
```

`ExecVerifier` uses `errors.As` to inspect `*exec.ExitError` after retaining
the diagnostic artifact. Exit `75` returns
`TemporaryVerificationFailureError`; other non-zero exits return the existing
`VerificationCommandError`. Cancellation, process-start failure, and artifact
filesystem failure stay infrastructure errors and are never retryable.

The Verification request gains retry identity and acquisition mode without
renumbering the existing Agent-repair attempts:

```go
type verificationAttemptRequest struct {
    // existing fields omitted
    Attempt int // 1 or 2
    Retry   int // 0 normally, 1 for exclusive retry
    Mode    verificationMode
}
```

### Data Models

No SQLite migration is needed. The Task-cycle-start `daemon.status` payload
keeps existing `concurrency` as a compatibility alias for Task Capacity and
adds:

```json
{
  "task_capacity": 2,
  "verification_capacity": 1
}
```

Each `daemon.verification` payload continues to carry `task`, `batch`,
`attempt`, `phase`, command, diagnostic path, and verdict where applicable.
Additive fields are:

- `retry: 1` for the exclusive retry, omitted otherwise;
- `mode: "shared" | "exclusive"` on waiting and started events;
- `capacity` with the effective Verification Capacity on waiting events;
- `classification: "temporary"` when exit `75` is observed;
- `retry_available: true | false` on a temporary failure;
- `reason: "temporary_verification_failure"` on retry/exhaustion evidence.

Add `waiting` to `VerificationPhase`. The existing `started`,
`command-passed`, `failed`, and `verdict` phases keep their meanings. A Task
emits `waiting` before every acquisition, even when the gate is immediately
available, then emits `started` only after acquisition. Before Verification
Feedback invokes the Agent, `daemon.task` emits phase
`verification_feedback`; the Live Run View maps both initial `started` and
`verification_feedback` to Agent working. Settlement remains the terminal
`daemon.task` event.

Failed normal attempts keep the current artifact path:

```text
runs/<run-id>/verification/batch-<NNN>-attempt-<1|2>.log
```

An exclusive retry uses a distinct path so it cannot overwrite the signal
that caused it:

```text
runs/<run-id>/verification/batch-<NNN>-attempt-<1|2>-retry-1.log
```

The one-retry budget belongs to the Task execution, not to each numbered
attempt. Consequently a Task can never receive two exclusive retries when a
deterministic failure and Agent repair occur in the same lifecycle.

Task status remains in the Task file. `runTaskAgent` sets `in_progress` through
the existing status writer before publishing Task start. After each Agent
turn, reload preserves Result content but rewrites status to `in_progress`
when the Agent changed it. A successful Agent call means
implementation-ready regardless of the reloaded terminal value. Agent runtime
failure, unreadable task content, Stop Request, and Daemon infrastructure
failure retain their explicit branches; only the Daemon decides whether to
settle failed or leave the Task resumable.

### API Contracts

The new YAML contract is:

```yaml
verification:
  concurrency: 1
```

Omission uses `1`; User and Project Config values overlay per key; values less
than `1` fail strict validation naming `verification.concurrency`. Generated
config includes the safe default and explains its independence from
`worktree.concurrency`. There is no CLI flag because capacity is repository
policy rather than a one-Run convenience override.

Exit code `75` is a child Verification-command protocol, not a new Roundfix
CLI exit code. Project scripts return it only after positively identifying a
temporary infrastructure/capacity condition. Roundfix does not inspect the
output to validate or discover that decision.

`roundfix events` remains newline-delimited `roundfix-events/v1` JSON with the
same filters and forward-compatible unknown-phase rule. Attach replay reads
the capacity fields from the Task-cycle-start event; for older Runs it falls
back to the recorded legacy `concurrency` and current configured Verification
Capacity. `LiveRunView` retains its internal Task Capacity value, adds
Verification Capacity, and renders both labels only for spec Runs.

No stdout report, Implement flag, Run state, terminal outcome, notification,
or top-level exit code changes.

## Coverage Map

- Goal 1 and Stories 1–2 → config, TaskPlan wiring, shared Verification gate.
- Goal 2 and Story 4 → Daemon-owned status transition and prompt contracts.
- Goals 3–4 and Stories 3, 5 → authoritative handoff, Run Events, Attach, Live
  Run View phase projection.
- Goal 5 and Stories 6–7 → exit-75 typed classification, Task-scoped retry
  budget, exclusive acquisition, diagnostic artifacts.
- Goal 6 and Story 8 → cancellation-aware acquisition, unchanged settlement,
  Task Worktree, dependency, and recovery policies.
- Core Features 5 and 7 → ADR-0038-compatible attempt numbering and one Agent
  repair independent from the exclusive retry.
- Core Feature 10 → Agent prompts, authorial workflow skill, user docs, and
  canonical/embedded Roundfix Skill synchronization.

## Integration Points

- **Shell Verification boundary:** commands still execute verbatim through
  `sh -c` in the Task Worktree. Only exit `75` receives the new typed meaning.
- **Task files:** the existing frontmatter status writer is reused. Result
  prose remains Agent-authored; status is Daemon-authored during Implement.
- **ACP Runtime:** initial and repair work stay in the Task's selected Agent
  Session. Capacity is never held while ACP work runs.
- **Run Event Journal:** existing append/replay storage carries all new
  evidence. Per-Task events, not the aggregate Run state, are authoritative
  when concurrent Tasks occupy different phases.
- **Bubble Tea v2:** the current synchronous model/event refresh seam derives
  labels. No terminal emulation, new keybinding, panel reflow, or mouse change
  is required.
- **Roundfix skills:** one final tooling-only Task owns the authorized
  `.agents/skills/implement-task/SKILL.md`,
  `skills/implement-task/SKILL.md`, `.agents/skills/roundfix/SKILL.md`, and
  `skills/roundfix/SKILL.md` changes. It updates each canonical/generated pair
  directly and byte-identically without running the broad `make skills-sync`
  mutation target. Upstream-managed skills remain untouched.

## Testing Approach

Tests name the invariant and use the lowest existing owning suite. Config
table tests cover default `1`, User/Project precedence, generated YAML, strict
unknown fields, and zero/negative rejection. `ExecVerifier` tests use real
shell processes and temporary directories to prove exit `0`, exit `1`, exit
`75`, cancellation, typed wrapping, and distinct retained artifact paths.

Task-cycle tests drive observable worker and Verifier boundaries with channels
rather than sleeps. They prove: Task Capacity `2` overlaps Agent work while
Verification Capacity `1` never exceeds one active attempt; capacity `2`
permits two; waiting precedes started; an exclusive retry begins only after all
normal attempts exit; an exclusive waiter is not starved; Agent repair runs
without a held permit and attempt 2 reacquires; one Task gets at most one
exclusive retry; and cancellation while queued starts no command or leaked
worker. Each positive path has a deterministic failure or cancellation
counterpart. Run the daemon suite under `-race`.

Status tests replace the old bypass assertion: Agent-authored `failed` and
`completed` are both normalized to `in_progress`, Daemon Verification still
runs, only settlement writes terminal state, Result prose survives, and a real
Agent execution failure is still Daemon-settled correctly. Prompt tests assert
the implementation-ready handoff and prohibition on declared Verification,
without snapshotting unrelated prose.

Run Event tests assert additive projection fields and ordering. Attach tests
replay new capacities and cover legacy fallback. Bubble Tea tests drive
`model.Update` synchronously and assert exact text labels for interleaved Tasks,
styled/no-color parity, and non-wrapping rendering at supported widths; no
whole-screen snapshot is added solely to freeze layout.

A CLI-level disposable-repository flow exercises Task Capacity `2`,
Verification Capacity `1`, a project wrapper returning `75` once, exclusive
retry success, journal replay, final Task status, and unchanged stdout/exit
behavior. Run `go test -race ./...`, `make verify`, and the skill-sync check.
No dependency beyond the standard library is introduced.

## Build Order

1. Add Verification configuration, generated/default documentation, TaskPlan
   plumbing, and capacity fields in the Task-cycle-start event.
2. Move Implement Task status ownership and Agent prompt contracts fully to
   the Daemon (depends on: 1).
3. Add the cancellation-aware shared/exclusive Verification gate, waiting
   events, shared attempt acquisition, and concurrency tests (depends on: 1).
4. Add typed exit-75 classification, distinct retry artifacts, the
   Task-scoped retry budget, exclusive retry flow, and repair interaction
   tests (depends on: 3).
5. Extend public Run Event projection, Attach capacity replay, and plain Live
   Run View capacity output (depends on: 1, 3, 4).
6. Add Bubble Tea per-Task phase derivation and accessible labels from
   interleaved Task/Verification events (depends on: 2, 3, 5).
7. Update configuration/usage guidance, align `docs/agents/autonomous-work.md`
   and the `CONTEXT.md` Agent Session definition with ADR-0051, preserve
   ADR/finding traceability, and add end-to-end QA fixtures (depends on: 2–6).
8. In one dedicated tooling-only Task, update the exact four authorized
   canonical/generated Roundfix and `implement-task` Skill files, then run
   read-only sync and full repository verification (depends on: 7).

## Risks & Considerations

- A project can misuse exit `75` and hide a deterministic defect for one run.
  Bounded one-time retry, retained diagnostics, and explicit events make that
  misuse visible; Roundfix cannot validate project-owned classification
  without returning to prohibited heuristics.
- An exclusive retry can starve if new shared attempts bypass it. The gate
  blocks new shared acquisition once an exclusive waiter exists and grants the
  exclusive request atomically after active attempts drain.
- Capacity `1` can lengthen a Wave with heavy gates. That is the safe default;
  repositories with isolated Verification resources may explicitly raise it.
- The aggregate Run state can say Verifying while another Task's Agent is
  working. Consumers must use per-Task journal phases for task-level truth;
  no new combinatorial Run states are added.
- A Stop Request during capacity wait leaves `in_progress` in the Task file.
  This is intentional resumable evidence, not a terminal verdict. Resume and
  terminal Run behavior must remain consistent with existing Stop contracts.
- Longer phase labels can pressure narrow panes. Reuse existing truncation and
  fixed layout accounting; do not resize panels or rely on color to convey the
  distinction.
- Changing Agent authorship can conflict with repository instructions that
  demand a full gate before a completion claim. The prompt must state that the
  Agent is not claiming Task completion: focused checks support handoff and
  the Daemon performs the mandatory full Verification before settlement.
- Current autonomous-work guidance and the Agent Session glossary still
  describe one Agent per Run and a separate frontend Spec, contrary to
  ADR-0051. Task 07 corrects those documentation targets; it must not remove
  the frontend Task from this graph or weaken per-Task Task Type routing.

## Decisions

- Use one fair, cancellation-aware shared/exclusive gate per Task cycle; do not
  add a global daemon or filesystem lock.
- Keep numbered Verification attempts `1` and `2`; identify the one
  infrastructure retry with `retry: 1` and a distinct artifact path.
- Treat exit `75` as Temporary Verification Failure only at the child-process
  boundary and preserve typed error inspection through wrapping.
- Apply one exclusive retry budget per Task lifecycle, independent from the
  one Agent repair.
- Preserve legacy Task-cycle `concurrency` event data while adding explicit
  Task and Verification capacity fields.
- Derive task-level UI truth from the Run Event Journal rather than multiplying
  aggregate Run states.
- Apply ADR-0051 per-Task Agent Session selection to this mixed Task Graph;
  older one-Agent-per-Run and separate-frontend-Spec prose is a documentation
  defect, not a reason to split or retype Task 05.
- Keep protected Skill mutations in the dedicated final tooling Task; code,
  tests, operator docs, manifests, and other Skill files remain outside that
  Task's changed-file scope.
- See [ADR-0056](../../adr/0056-spec-runs-separate-task-and-verification-capacity.md).
- See [ADR-0057](../../adr/0057-daemon-exclusively-owns-implement-task-status.md).
