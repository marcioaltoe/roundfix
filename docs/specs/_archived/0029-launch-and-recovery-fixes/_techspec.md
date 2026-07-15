---
spec: 0029-launch-and-recovery-fixes
prd: _prd.md
created: 2026-07-15
---

# Launch and Recovery Fixes — Technical Spec

## Executive Summary

Three bug fixes with root causes established by live reproduction (see the PRD's findings references; the detach cause was isolated this session with a minimal spawn repro). The primary trade-off: the Detached Run start must stay caller-bounded — a broken child must fail fast — while tolerating an arbitrarily slow but healthy Preflight Validation, whose model probe launches a real acpx disposable session. The design resolves that with a two-phase handshake: liveness is proven within the existing short deadline, run creation gets its own generous ceiling. The model-rejection fix deliberately does not try to make the preflight probe infallible — codex-acp's advertised-model list is assembled per session and observed to differ between sessions minutes apart — it instead makes the Batch-time rejection carry the same actionable content the preflight rejection already has. The settle fix is an ordering correction in surface resolution plus mandatory surface reporting.

## System Architecture

All changes land in existing modules; no new packages or seams.

- **`internal/cli` (detach)** — `runDetachedCommand` gains the two-phase wait; the child side writes a liveness marker on entry to child mode (before config load and Preflight Validation) and the existing run-id line unchanged afterward. All three failure branches (phase-1 timeout, phase-2 timeout, child exit before handshake) print explicit diagnostics.
- **`internal/agent`** — the acpx stderr classification recognizes the "did not advertise that model" rejection and produces a selection-failure error carrying the model, the advertised list parsed from the message, and the same recovery text `SelectionPreflightError` renders today.
- **`internal/daemon`** — a Batch that fails with that classification settles its issues/Task with a terminal reason naming the rejected model, and the Run report/console shows the actionable message instead of `agent/protocol error`.
- **`internal/cli` (doctor)** — a `model:` check line after the existing `agent:` line, reporting the effective Agent Model probe outcome for the configured runtime; failures include the advertised list.
- **`internal/cli` (settle)** — surface resolution loads the Task's status per candidate surface and picks the first with status `failed` (Task Worktree → Run Worktree → current repository); a `Settle surface: <path>` stderr line always precedes verification; when no candidate qualifies, the refusal names each candidate path and the status found there.
- **`skills/roundfix` (via `.agents/skills/roundfix`)** — Detached Run, Doctor, and Settle sections updated with the shipped behavior; skill-check anchors updated together.

## Implementation Design

### Interfaces

Two-phase detach handshake (internal protocol between the parent and its re-exec'd child; same binary, same release — no compatibility window needed):

```go
// Child writes one liveness byte ('\x06') on the handshake fd immediately
// on entering child mode, then the existing "<runID>\t<consoleLog>\n" line
// once the Run exists. Parent-side wait:
//   phase 1 (liveness): detachHandshakeTimeout (10s, unchanged)
//   phase 2 (run id):   detachStartCeiling (5m) — covers Preflight
//                       Validation including the real agent probe
// Any child exit ends the wait immediately via pipe EOF, as today.
```

Failure diagnostics (all on the caller's stderr, before the relayed console temp):

```text
roundfix: Detached Run child produced no liveness signal within 10s; killed (exit: signal: killed)
roundfix: Detached Run child did not create a Run within 5m0s; killed (exit: signal: killed)
roundfix: Detached Run child exited before the handshake (exit status 2); console output follows
roundfix: Detached Run child exited before the handshake (exit status 1) and produced no output
```

Model-rejection classification (`internal/agent`):

```go
// ModelNotAdvertisedError is recognized from the acpx stderr tail
// ("Cannot apply --model ...: the ACP agent did not advertise that model.
//  Available models: ...").
type ModelNotAdvertisedError struct {
    Runtime    string
    Model      string
    Advertised []string // parsed best-effort; may be empty
}
// Error() renders the same recovery guidance as the preflight
// SelectionPreflightError: update the runtime/adapter, choose an
// advertised model, or use a one-Run --model override.
```

Settle surface selection (`internal/cli/settle.go`):

```go
// resolveSettleSurface returns the first candidate whose task file has
// status failed, in order: Task Worktree, Run Worktree, current repository.
// It returns every candidate's (path, status) for refusal reporting.
```

### Data Models

None — no store, artifact, or config changes. (`detachStartCeiling` is a compile-time constant like the existing timeouts; making it configurable is deliberately out of scope until someone needs it.)

### API Contracts

- The four-line detach success report, all exit codes, and the settle stdout lines (`verify …`, `commit <path>`, `settled …`) are unchanged.
- New stderr diagnostics on detach failure paths (shapes above) — additive; today those paths print nothing.
- New stderr line `Settle surface: <path>` before settle verification — additive.
- Settle refusal for a non-failed Task now enumerates candidate surfaces and their statuses instead of judging only the first resolved surface.
- Doctor gains one `model: ok (<model>)` / `model: failed …` line — additive, mirroring the existing check-line format with `next:` on failure.
- Batch failure reason for a rejected model becomes `Agent Model "<model>" not advertised by runtime "<runtime>"; advertised: <list>` (in Run Events, task/issue terminal reasons, and the final report) instead of `agent/protocol error`.

## Coverage Map

- PRD Goal 1 and 2, Core Feature 1 → two-phase handshake + failure diagnostics (`internal/cli` detach).
- PRD Goal 3, Core Features 2–3 → `ModelNotAdvertisedError` classification (`internal/agent`), Batch settlement wiring (`internal/daemon`), Doctor model line (`internal/cli`).
- PRD Goal 4, Core Feature 4 → settle surface resolution + reporting (`internal/cli/settle.go`).
- Core Feature 5 → skill and docs sync (`.agents/skills/roundfix`, mirror regenerated).

## Integration Points

- **acpx** (existing boundary): the classification parses acpx's stderr tail — best-effort by design; when parsing fails the current infrastructure error remains, so the change can only improve messages, never mask them. The Doctor model probe reuses the existing disposable-session probe path; no new acpx invocation shapes.
- No GitHub, store schema, or config surface involvement.

## Testing Approach

Existing seams only.

- **Detach unit/integration**: the passing detach test gains siblings — a child stub that sleeps past phase 1 before liveness (parent reports the phase-1 diagnostic), a child that signals liveness then sleeps past a shrunk phase-2 ceiling (constants injectable for tests, as the current timeout already is), a child that exits 2 pre-handshake with output (relay includes it), and a child that exits silently (the "produced no output" line). The real-probe slow-start case is covered by liveness-then-slow-preflight, no live agent needed.
- **Classification unit tests**: stderr-tail fixtures (exact message, wrapped message, unparseable garbage) → error type, advertised-list parsing, and fallback to the existing infrastructure error.
- **Engine tests**: a fake runner failing with the not-advertised stderr settles the Task/issues with the new terminal reason and the report shows it.
- **Doctor test**: fake probe success and failure render the `model:` line with `next:` guidance.
- **Settle tests**: regression of the field case — two seeded kept Runs where the stale worktree has the Task `pending` and the checkout has it `failed` → settle picks the checkout, prints the surface line; a no-qualifying-surface case renders the per-candidate refusal.

## Build Order

1. **Detach two-phase handshake and diagnostics** — child liveness marker, parent phase split, all failure-branch messages, tests with injectable timings.
2. **Model-rejection classification** (`internal/agent`) — error type, stderr parsing, recovery text shared with the preflight error; unit tests.
3. **Batch settlement and report wiring** (depends on: 2) — daemon consumes the classification for terminal reasons, Run Events, and the final report.
4. **Doctor model line** (depends on: 2) — probe outcome rendering with the advertised list on failure.
5. **Settle surface preference and reporting** — resolution ordering by Task status, surface line, per-candidate refusal; regression tests.
6. **Skill and docs sync** (depends on: 1, 3, 4, 5) — Detached Run failure semantics, Doctor line, Settle surface line in the canonical skill; mirror regenerated; anchors updated.

## Risks & Considerations

- **Phase-2 ceiling is a new hang bound, not a behavior change**: a healthy preflight of any realistic length fits 5m; a genuinely hung child is still killed and now reported. If a repository's Worktree Bootstrap pushes Run creation past 5m, the diagnostic names the phase — revisit configurability only then.
- **stderr-tail parsing is inherently fragile**: treated as progressive enhancement; unparsed rejections keep today's error, so no regression is possible.
- **Settle ordering change could surprise a user who wanted the stale-worktree surface**: the mandatory `Settle surface:` line makes the choice visible, and the old behavior was the reproduced bug.
- The dogfood machine's probe latency (~11s) sits just past the old deadline — the fix must land before any further `--detach` use; until merged, the nohup workaround stands.

## Decisions

- Two-phase handshake over simply raising the single timeout: liveness stays fast-failing (broken spawn detected in 10s) while Run creation gets a ceiling that covers real probes — one number cannot serve both.
- Batch-time classification over probe hardening: the advertised list is per-session-nondeterministic upstream, so the probe can never guarantee the Batch's view; making the failure actionable is the durable fix (the preflight probe stays as-is).
- Settle prefers the surface where the Task is `failed`; surface choice is always reported. No new flags.
- No ADRs minted: every decision here restores or clarifies contracts owned by earlier ADRs (ADR-0028 detached runs, ADR-0037/0039/0041 model selection, the Settle Command contract) without reversing any.
