---
spec: 0006-acpx-run-robustness
prd: _prd.md
created: 2026-07-05
---

# ACPX Run Robustness — Technical Spec

## Executive Summary

Three contained changes. The consequential one is classification: the acpx
runner currently maps any nonzero exit to failure, which twice discarded
finished work; per ADR-0020 the runner now returns success-with-result when
the NDJSON stream already carried the `session/prompt` response, journaling
the exit as an anomaly. The trade-off accepted: a Batch can now "succeed" on
a wounded transport — safe because the Daemon's verbatim verification (ADR-
0014) was always the real gate, and the anomaly is journaled loudly rather
than swallowed. The Settle Command reuses the Daemon's existing collaborators
(Verifier, Committer, spec status writer) with no engine involvement, and the
buffer work is investigation plus configuration, not new machinery.

## System Architecture

- `internal/agent` — the acpx runner's exit handling: track whether the
  prompt result line was parsed; on nonzero exit with a parsed result,
  return the result plus an anomaly (bounded stderr tail) instead of an
  error. No interface changes.
- `internal/daemon` — both engines forward the anomaly into the journal via
  existing daemon status events (kind reuse; no vocabulary change), then
  proceed to verification exactly as on a clean exit.
- `internal/cli` — new `settle` command file following the stop/implement
  conventions: flag parsing, Preflight Validation, verification loop through
  the existing `Verifier`, settlement through the spec status writer, commit
  through the existing `Committer` with `TaskCommitMessage`.
- `internal/spec` — nothing new (Load, ReloadTask, SetStatus suffice).
- acpx configuration — buffer mitigation applied in `~/.acpx/config.json`
  guidance or invocation flags if the pinned version exposes one (task
  verifies against acpx 0.12.0's source/help; nothing hardcoded on hope).

## Implementation Design

### Interfaces

```go
// internal/agent — ExecuteResult gains the anomaly; Runner unchanged
type ExecuteResult struct {
    // existing fields...
    TransportAnomaly string // "" normally; bounded exit+stderr summary when
                            // a nonzero exit followed a parsed prompt result
}
```

```
roundfix settle --spec <slug> --task <task_id> [--no-input]
```

- Preflight (exit 2, one actionable message each): repository resolved; Spec
  loads valid; task id exists in the Task Graph; task status is `failed`
  (anything else names the status and the right path: pending/in_progress →
  implement, completed → nothing to do); no Active Run for the spec target
  or the working tree; `--spec` and `--task` both required (no interactive
  picker — recovery is deliberate).
- Execution: run the Task's Verification commands verbatim through the
  existing Verifier in the repository root, streaming one stdout line per
  command: `verify <command> — ok|failed`.
- On all-pass: `SetStatus(completed)`, stage all current worktree changes
  plus the task file, commit with `TaskCommitMessage` (identical message and
  trailers to a Daemon settlement). stdout ends `settled task_NN completed —
  <commit short sha>`; exit 0.
- On any failure: nothing written, stdout ends `task_NN stays failed —
  verification failed`; exit 1.
- No Run row, no journal events, no push — a local recovery command in the
  stop family.

### Classification change (ADR-0020)

The runner already parses each stdout line; it records when the
`session/prompt` response arrives. On process exit: exit 0 → as today;
nonzero exit with recorded result → return the result with
`TransportAnomaly` set (exit code + bounded stderr tail; same bounding as
the infrastructure error); nonzero exit without a result → the existing
failure classification, untouched, including exits 2/4 as infrastructure
errors. Engines journal a non-empty anomaly through the existing daemon
status event with the anomaly text in the payload, then continue into
verification.

### Data Models

None. No store, journal schema, or config keys.

### API Contracts

- New command `settle` as specified above; house exit codes; deterministic
  stdout; diagnostics on stderr; help text truthful.
- No changes to existing commands' contracts. Batch behavior changes only in
  the ADR-0020 case, observable as: Run proceeds to verification, journal
  carries one anomaly status event.

## Coverage Map

- Story 1 → runner classification + engine anomaly journaling (ADR-0020)
- Story 2 → Settle Command (findings-2 item 3)
- Story 3 → buffer mitigation task (findings-2 item 1)
- Story 4 → anomaly journaling assertions + settle stdout report

## Integration Points

acpx only (existing boundary): classification reads the same stream and exit
the runner already owns; the mitigation task touches acpx configuration
surfaces verified against the pinned 0.12.0.

## Testing Approach

Existing rigs throughout: the helper-process fake acpx gains scripts that
emit a full result then exit nonzero (classification matrix: result+0,
result+1, no-result+1, no-result+2/4, result+130); engine tests assert the
anomaly event lands in the journal and verification still gates the commit.
Settle gets buffer-captured CLI tests over a real temp repo + store: every
preflight refusal, verification-pass settlement (commit message/trailers
byte-compared to the Daemon's), verification-fail no-op, and the
staging-includes-worktree contract. The gated real-acpx integration test
extends with the post-result kill case if reproducible cheaply; otherwise the
fake rig covers it. Full suite is the regression net.

## Build Order

1. Result-over-exit classification in the runner plus engine anomaly
   journaling (no deps)
2. Settle Command (no deps)
3. Buffer mitigation investigation against pinned acpx: apply config if it
   exists, else document + draft the upstream report; record the outcome in
   the round-2 findings log (no deps)
4. Docs and skill sync: settle command, ADR-0020 semantics, any buffer
   guidance (depends on: 1, 2, 3)

## Risks & Considerations

- Classification must never soften real failures: only a fully parsed prompt
  response flips the outcome; partial streams stay failures. The test matrix
  is the guard.
- Settle's stage-everything contract is deliberately blunt; the docs task
  must state it plainly so nobody expects path scoping (worktree-per-task
  later makes this precise).
- The buffer investigation may conclude "no mitigation exists" — that is an
  acceptable outcome with the upstream report as the deliverable; the task
  must not invent flags.

## Decisions

- Parsed result outranks exit code; verification stays the gate. See
  ADR-0020.
- `TransportAnomaly` rides the existing result type; no Runner interface
  change; anomaly journaled via existing event kinds.
- Settle: no Run, no journal, `failed`-only targets, stage-everything
  contract, `--spec`+`--task` required.
- Buffer mitigation is evidence-driven — config applied only if the pinned
  acpx actually exposes one.
