# A Run that can hand back its work — Technical Spec

## Executive Summary

Four defects, one shape: a surface decides terminal state for a resource another
surface created. The design keeps each decision with its owner.

The Fallback Chain becomes eligible again by moving one signal. `ACPXRunner.Run`
publishes `AgentWorkStartedStatus` after `PrepareSession` and before the first
prompt, and `agentSessionOwner.Run` does the same after `activate`. A first
prompt that returns an adapter-level refusal with no Agent output therefore
arrives with the flag already set. The signal moves to the first Agent output.

Batch settlement stops deciding the Run's outcome. A Batch records what its Agent
achieved; the Run derives `Unresolved` from Review Issues still open. Both facts
survive.

A stopped Run gains two dispositions it lacks: settled Tasks can be handed to the
checkout on evidence, and a superseded Run Branch can be discarded through
Roundfix with its contents recorded first.

## Project Constraints

- Identifier strategy: not applicable — no new persisted entity.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — Run lifecycle, git plumbing and
  local storage. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0010 lets a Batch settle and the Run
  continue. ADR-0050 has a configured fallback activate only after its notification; this Spec narrows when a chain is *ineligible*
  and widens nothing. ADR-0091 makes the QA gate its own Task node. ADR-0104
  requires an outside-evidence acceptance row. This Spec adds ADR-0113,
  ADR-0114 and ADR-0115. Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — ordinary source under `internal/daemon`,
  `internal/agent`, `internal/rounds`, `internal/spec` and `internal/cli`.
  Source: `docs/agents/agent-instructions.md`.

## Vocabulary Contract

- emits: `internal/agent/agent.go`
  pattern: `agent_(work_started|selection_failed)`
  documented-in: `CONTEXT.md`

## System Architecture

**Selection boundary (`internal/agent/acpx_runner.go`,
`internal/daemon/agent_session_owner.go`).** `AgentWorkStartedStatus` moves from
"the session is prepared" to "the Agent produced its first output". A turn that
ends with no Agent output and an adapter-level refusal publishes a new selection
failure instead, which leaves the Fallback Chain eligible.

**Settlement and outcome (`internal/agent/agent.go`,
`internal/daemon/engine.go`).** `MarkBatchFailed` re-reads each Review Issue and
preserves Terminal outcomes, matching its sibling `SettleAssignedIssues`. The Run
outcome, previously implied by the Batch's failure, is computed from the Review
Issues that remain unresolved.

**Branch disposition (`internal/cli`, `internal/daemon`).** A named act
classifies a Run Branch as superseded — its commits already reachable from the
target, or its Run replaced by a later one covering the same Tasks — writes what
it held, and then removes branch and worktree. `reconcile` keeps reporting and
never disposes.

**Task carry-forward (`internal/spec`, `internal/cli`).** A stopped Run's Tasks
whose Verification passed and whose declared inputs are unchanged since
settlement can be handed to the checkout, with the Run Worktree commits as the
record. It is an explicit maintainer act.

## Implementation Design

### Interfaces

```go
// SelectionFailure marks a turn that ended before the Agent produced output,
// so a Fallback Selection stays eligible.
type SelectionFailure struct {
    Runtime string
    Reason  string // quota, authentication, adapter startup
    Err     error
}

// BranchDisposition is why a Run Branch may be discarded.
type BranchDisposition struct {
    Branch     string
    Superseded bool
    Reachable  bool     // every commit already reachable from the target
    Commits    []string // recorded before removal
}

// CarryForward is one settled Task a stopped Run can hand back.
type CarryForward struct {
    TaskID       string
    RunID        string
    Commit       string
    InputsMoved  bool // true blocks the carry-forward
}
```

### Data Models

Review Issue statuses are unchanged. What changes is who reads them: the Run
outcome derives from a count of unresolved issues rather than from a Batch
failure flag.

### API Contracts

- `roundfix reconcile` — unchanged, read-only.
- `roundfix reconcile --discard-superseded` — the named disposition. Refuses any
  branch it cannot prove superseded, and writes the branch record before
  removing anything.
- `roundfix reconcile --carry-forward` — hands settled Tasks back, refusing any
  Task whose declared inputs moved.

## Coverage Map

| PRD goal | Component |
| --- | --- |
| 1 — the chain activates on a runtime that could not serve | The selection boundary and `SelectionFailure` |
| 2 — a failed Batch preserves outcomes, Run still Unresolved | `MarkBatchFailed` preservation plus the derived Run outcome |
| 3 — a stopped Run hands settled Tasks back | `CarryForward` and `reconcile --carry-forward` |
| 4 — a superseded branch has a disposition | `BranchDisposition` and `reconcile --discard-superseded` |
| 5 — no Run reports Clean on unfinished work | The derived Run outcome, proven by the six tests that previously encoded the opposite |

## Integration Points

- `internal/agent/acpx_runner.go`, `internal/agent/agent.go`
- `internal/daemon/agent_session_owner.go`, `internal/daemon/engine.go`
- `internal/rounds/rounds.go` — unresolved-issue count
- `internal/cli` — the two reconcile flags
- `internal/cli/cli_test.go`, `internal/daemon/engine_test.go` — the six tests
  that encode the current outcome contract and must be rewritten deliberately

## Testing Approach

The corpus is captured first and is unusually load-bearing here, because one of
the changes was already attempted and reverted. The six tests that encode "a
failed Batch means an Unresolved Run" are recorded as declared breaks with the
contract they assert, so the Task that changes them changes them on purpose.

The acceptance no fixture substitutes for: a real Run against a runtime whose
selection cannot be served, observed activating its configured fallback. The
2026-08-08 instance was a quota exhaustion, which cannot be summoned on demand;
an unreachable adapter command is the reproducible equivalent and the gate says
which one it used.

## Build Order

1. **Characterization corpus.** Record all four behaviours: work-started before
   the first prompt, Batch overwrite, a stopped Run's pending Tasks, and
   Preflight refusing on a superseded branch. Name the six outcome-contract
   tests as declared breaks. Depends on nothing.
2. **The selection boundary.** Move `AgentWorkStartedStatus` to the first Agent
   output; publish a selection failure otherwise; leave the chain eligible.
   Depends on step 1.
3. **Settlement preserves outcomes.** `MarkBatchFailed` re-reads and preserves
   Terminal statuses. Depends on step 1.
4. **The Run outcome is derived.** Compute `Unresolved` from remaining
   unresolved issues, and rewrite the six tests to the new contract. Depends on
   step 3, and only step 3 makes it safe.
5. **Superseded branch disposition.** Classify, record, then remove. Depends on
   step 1.
6. **Task carry-forward.** Hand settled Tasks back on unchanged inputs. Depends
   on step 1.
7. **QA gate.** Depends on steps 2, 4, 5 and 6.

## Risks & Considerations

**Rewriting six tests is the risk, not the work.** They encode a real contract,
and the Spec that changes them must show the new contract is better rather than
merely different. Step 4 exists as its own Task for that reason, and its
rehearsal cases name each rewritten test.

**A carry-forward that is wrong is silent.** A Task reported complete against
inputs it never saw is worse than re-execution. Hence the explicit act and the
inputs-unchanged condition, and hence carry-forward refusing rather than warning.

**Discarding a branch is irreversible.** The record is written first, and the
disposition refuses any branch it cannot prove superseded.

## Decisions

- Batch settlement and Run outcome are separate, per ADR-0113.
- Opening a session is selection, not work, per ADR-0114.
- Disposal is a named act, never a widening of `reconcile`, per ADR-0115.
