# A gate that could have failed — Technical Spec

## Executive Summary

A Task's `## Verification` commands decide whether an Agent turn produced the
work it was asked for. Today they are run once, after the Agent, and a zero exit
is read as proof. Spec 0089 showed that a command can exit zero against a tree
nobody touched, and that nothing in the record distinguishes a command the Daemon
ran from a command an Agent described.

This Spec adds one execution of those same commands **before** the Agent turn,
classifies the result, and records what actually ran. The mechanism is small
because the seam already exists: `verifyRepositoryPrecondition` is a Verification
that runs before the Agent, and the pre-work probe is its sibling.

It also gives Task Verification a third outcome. A command that timed out, ran
partially, or never reached the surface it names is not a pass and not a
refutation of the work — it is `unknown`, and it stops the Task instead of
settling it either way.

## Project Constraints

- Identifier strategy: not applicable — no new persisted entity or resource
  identity. The probe keys off the existing Task ID and Run ID.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — the probe runs through the Daemon's
  existing command runner and reaches no network surface.
  Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0083 makes `make verify` the only
  authoritative gate, so its determinism is in scope here. ADR-0091 keeps the
  authored QA gate before any Pull Request; this Spec moves one defect class
  earlier without displacing it. ADR-0096 already establishes that the QA gate
  proves machine facts before spending an Agent turn, and this Spec applies that
  same principle to a Task's own Verification. ADR-0104 requires an acceptance
  row on evidence the Spec did not author; the PRD's `## External evidence`
  section names it. Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — every component below lives in
  `internal/daemon`, `internal/spec`, and `internal/runevent`, which are ordinary
  source. The `write-tasks` skill would benefit from stating the rule this Spec
  enforces, and that edit is protected tooling; it is deliberately **out of
  scope** and needs its own express authorization with bounded files.
  Source: `docs/agents/agent-instructions.md`.

## Vocabulary Contract

- emits: `internal/runevent/event.go`
  pattern: `verification_(vacuous|unknown)`
  documented-in: `CONTEXT.md`

Spec 0089 shipped a coined encoding, `runtime_deferred`, that its own TechSpec
never declared here — so the automated glossary check was skipped and the term
reached the built CLI undefined, which its QA gate then had to catch by hand.
This declaration exists so the same omission fails mechanically. `Negative
Control` is the other coined term and belongs in the same glossary entry.

## System Architecture

No new package. Three existing seams carry the whole design.

**Dispatch (`internal/daemon/task_engine.go`).** `executeTaskWorker` already runs
`verifyRepositoryPrecondition` before it creates an Agent Session owner. The
pre-work probe is inserted immediately after it, on the same shape: run
commands, classify, and settle the Task `failed` without spending a turn.

**Outcome classification (`internal/daemon/engine.go`).**
`verificationAttemptOutcome` carries `Failure`, `CommandFailure`, and
`TemporaryFailure`. It gains an explicit `Unknown` cause so a partial or
unobserved execution stops being folded into either pass or fail.

**Evidence (`internal/runevent`).** The Daemon already publishes
`KindDaemonVerification` events with a classification and reason. The probe
publishes on the same channel with its own classification, so the Run Event
stream — not a Task's prose — becomes the record of which commands ran.

The three controls the PRD names are delivered asymmetrically, on purpose:

| Control | Delivered how | Why here |
| --- | --- | --- |
| observability | the pre-work probe, mechanically, for every Task | The Daemon can run it with no author cooperation and no extra Agent turn. |
| positive | the existing post-Agent Verification | Already exists; it is the run that passes on correct work. |
| negative | authored per Task, recorded, not synthesised | Roundfix cannot mutate a repository to manufacture a defect safely. What it can do is require the author to name one and record whether it was exercised. |

## Implementation Design

### Interfaces

```go
// VerificationProbe classifies one Task's Verification commands against the
// tree as it stands before the Agent runs.
type VerificationProbe struct {
    TaskID   string
    Commands []VerificationProbeResult
}

type VerificationProbeResult struct {
    Command string
    // Vacuous is true when the command exited zero before any work was done.
    Vacuous bool
    // Unknown is true when the command could not be observed: a timeout, a
    // partial execution, or a runner error that is not the command's verdict.
    Unknown bool
}

func (probe VerificationProbe) Vacuous() []string
func (probe VerificationProbe) Unknown() []string
```

### Data Models

`verificationAttemptOutcome` gains one field:

```go
type verificationAttemptOutcome struct {
    Failure          string
    CommandFailure   *VerificationCommandError
    TemporaryFailure *TemporaryVerificationFailureError
    UnknownCause     *VerificationUnknownError // new
}
```

`UnknownCause` is set when the runner could not observe the command's verdict.
It is never set together with `CommandFailure`: a command that ran and exited
non-zero has a verdict, and that verdict is failure.

`spec.Task` gains one parsed section, mirroring `Verification`:

```go
type Task struct {
    // ...
    Verification    []string
    NegativeControl []string // parsed from ## Negative Control
}
```

### API Contracts

Two Run Event classifications join the existing set:

- `verification_vacuous` — the named command exited zero before the Agent ran.
  Payload carries the command and the Task ID.
- `verification_unknown` — the command's verdict could not be observed. Payload
  carries the command, the reason, and the diagnostic path.

Both are `KindDaemonVerification`, so every existing consumer of the event stream
keeps working.

## Coverage Map

| PRD goal | Component |
| --- | --- |
| 1 — a command that cannot fail is rejected before the turn | The pre-work probe in `executeTaskWorker` |
| 2 — a recorded Result is distinguishable from a claimed one | `verification_vacuous` and the probe's command record on the Run Event stream |
| 3 — the gate returns the same verdict for the same tree | The shared wait budget sweep, and a gate-determinism test |
| 4 — an unobserved Verification returns `unknown` | `UnknownCause` on `verificationAttemptOutcome` |

| PRD core feature | Component |
| --- | --- |
| Three controls, not one check | The probe (observability), existing post-Agent Verification (positive), `## Negative Control` (negative) |
| Gate health is recorded, not assumed | The probe's per-command record, published per Task |
| A rubric that predates the implementation | The probe reads `## Verification` before the Agent, which is the mechanical proof it predates the work |
| Recorded rather than narrated evidence | The Run Event payload replaces prose as the record |
| A gate that is the same twice | The wait-budget sweep and its regression |

## Integration Points

- `internal/daemon/task_engine.go` — probe insertion, Task settlement on a
  vacuous command.
- `internal/daemon/engine.go` — `UnknownCause`, and the reason strings that
  reach a Task's terminal reason.
- `internal/daemon/daemon.go` — the command runner distinguishes "ran and
  failed" from "could not be observed".
- `internal/spec/spec.go`, `internal/spec/task.go` — parse `## Negative
  Control`.
- `internal/runevent` — two classifications.
- `internal/agent/acpx_runner_test.go` and `internal/store/process_unix_test.go`
  — the two literal wait budgets that made the gate non-deterministic.

## Testing Approach

Every component ships its own gate, and each gate is proven able to fail by
running it against a deliberately reverted implementation before the Task
settles. That practice is the Spec's own subject, so it is not optional here.

The characterization corpus is captured first, on the current behaviour: a Task
whose Verification passes before the Agent runs today settles `completed`, and a
Verification that times out today settles `failed`. Both are recorded as
declared breaks before the change lands.

The acceptance case that no fixture substitutes for: a Task graph carrying Spec
0089's exact `task_05` gate — `grep -q 'reasoning_effort: xhigh' .roundfixrc.yml`
against a file that already contains it — must be refused at dispatch, with no
Agent turn spent.

## Build Order

1. **Characterization corpus.** Record today's behaviour for a vacuous
   Verification and for an unobservable one. Declares both breaks. Depends on
   nothing.
2. **`UnknownCause` on the outcome.** The command runner separates "ran and
   returned non-zero" from "could not be observed", and the outcome carries it.
   Depends on step 1.
3. **The pre-work probe.** Run the Task's `## Verification` in
   `executeTaskWorker` before the Agent Session owner exists; classify each
   command; settle the Task `failed` when any command is vacuous, naming it.
   Depends on step 2, because a probe command that cannot be observed must be
   `unknown` rather than "not vacuous".
4. **Run Event classifications.** Publish `verification_vacuous` and
   `verification_unknown` with their payloads, so the record lives on the event
   stream. Depends on step 3.
5. **`## Negative Control` parsing.** Parse the section, carry it on
   `spec.Task`, and record whether it was exercised. Depends on step 1 only;
   runs parallel to 2–4.
6. **Gate determinism sweep.** Source both literal wait budgets from their
   shared constants and add the regression that fails when a budget is restated
   at a call site. Depends on nothing; runs parallel to everything.
7. **QA gate.** The authored terminal gate. Depends on 4, 5, and 6.

## Risks & Considerations

**A legitimately invariant Verification.** Some Tasks assert an invariant that
holds before and after — a Task whose whole job is "this must never regress".
Under this design that Task is refused. The PRD calls this correct: the
`write-tasks` contract already forbids a Verification that does not prove its
own effect. The risk is that the refusal arrives at dispatch, where it is
expensive to fix, rather than at authoring, where it is cheap. That is the
argument for the out-of-scope skill edit, and it is why this Spec names it.

**Probe cost.** Every Task pays one extra execution of its own Verification
commands. Those are Task-scoped and cheap by contract — the authoritative
repository gate is a separate, already-existing precondition. Worth measuring in
the QA gate rather than assuming.

**`unknown` as a new terminal state.** A Task that ends `unknown` is neither
completed nor failed. This Spec deliberately does not add a fourth `spec.Status`:
an unobservable Verification settles the Task `failed` with an `unknown` cause on
its terminal reason and its event. Introducing a fourth Task status would change
the Run outcome contract, which is Spec 0092's subject, not this one.

## Decisions

- The probe runs commands rather than parsing them. A static rule about `grep`
  shapes would have missed Task 05's command, which is well-formed and would be
  correct against a file that did not already contain its needle. See ADR-0109.
- The negative control is authored, not synthesised. Roundfix will not mutate a
  repository to manufacture a defect; it requires the author to name one and
  records whether it was exercised. See ADR-0110.
- `unknown` is a cause on an existing outcome, not a fourth Task status. See
  ADR-0111.
