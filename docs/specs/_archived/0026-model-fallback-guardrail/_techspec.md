---
spec: 0026-model-fallback-guardrail
prd: _prd.md
created: 2026-07-11
---

# Model Fallback Guardrail — Technical Spec

## Executive Summary

On an Agent selection Preflight Validation failure, Roundfix probes the
runtime's Model Catalog with disposable Agent Sessions to find a proven
Fallback Selection, then gates it behind explicit human confirmation: an
interactive prompt, or a non-interactive failure that prints the exact
explicit-flags re-run. The primary trade-off is probe cost — each candidate
costs a disposable adapter session (seconds) on a path that is already
failing — accepted because ADR-0039 established that only a live probe proves
a selection works, and a wrong fallback offer would burn the user's one
confirmation on another broken model. The design reuses the existing
selection probe, error, catalog, and effective-selection machinery; the only
new seams are one fallback-probe function in the agent layer and one
confirmation prompt in the CLI preflight path.

## System Architecture

- `internal/agent` — new fallback probe alongside the existing `Probe`: walks
  caller-supplied candidates (Model Catalog order, failed model excluded),
  proves a model with a disposable Agent Session, then proves the highest
  functional reasoning effort by trying caller-supplied effort values
  highest-first; no settable effort with a working model means a
  model-managed Fallback Selection (ADR-0040 semantics from spec 0025).
- `internal/cli` — the commands' selection preflight (the existing
  `Probe`-failure path used by resolve, watch, and implement) gains the
  guardrail: on `SelectionPreflightError`, run the fallback probe, then
  either prompt (interactive) or fail with the fallback report
  (non-interactive). The candidate model list comes from the existing Model
  Catalog; the candidate effort order comes from the existing per-runtime
  reasoning choices, highest first.
- Confirmation is a plain stderr prompt with a stdin answer, following the
  Interactive Input style — not the Live Run View. `--no-input`, `--detach`,
  and non-TTY stderr all mean non-interactive.
- Doctor and Setup keep their read-only contract: they report the selection
  failure and may name the probed fallback, but never prompt.
- Effective selection persistence (spec 0023) is reused untouched: a
  confirmed fallback simply becomes the RuntimeSpec the Run is created with,
  so the Run row, headers, and Attach already tell the truth.

## Implementation Design

### Interfaces

```go
// internal/agent — shape, not implementation
type FallbackCandidateSet struct {
    Models  []string // Model Catalog order, failed model excluded
    Efforts []string // runtime effort vocabulary, highest first
}

type FallbackSelection struct {
    Model           string
    ReasoningEffort string // empty = model-managed (ADR-0040)
}

// ProbeFallback proves candidates with disposable Agent Sessions and
// returns the first functional selection, or ok=false when none proves.
func (r ACPXRunner) ProbeFallback(ctx context.Context,
    runtime RuntimeSpec, c FallbackCandidateSet,
) (FallbackSelection, bool, error)
```

```go
// internal/cli — the guardrail wraps the existing probe failure path
// pseudo-shape: on SelectionPreflightError from Probe:
//   fb, ok := ProbeFallback(...)
//   interactive: confirm(fb) -> proceed with fb, or fail as today
//   non-interactive: fail with fallback report + explicit re-run command
```

### Data Models

None. The confirmed fallback flows through the existing RuntimeSpec and Run
row (`runs.model`, `runs.reasoning_effort`).

### API Contracts

- Interactive selection-failure prompt (stderr), one confirmation question;
  an empty or negative answer declines. Declining and the no-candidate case
  end in today's Preflight Validation failure (exit 2).
- Non-interactive selection failure adds a deterministic report naming the
  failed selection, the proven Fallback Selection (with `model-managed` for
  an empty effort), and one copy-paste re-run line with explicit `--model`
  and `--reasoning-effort` flags. Exit code stays 2. stdout stays reserved
  for command output; the report goes to stderr.
- No new flags, no new config keys (ADR-0041). `--model`/`--reasoning-effort`
  already exist as the confirmation vehicle.
- When the fallback probe itself finds nothing, the original
  `SelectionPreflightError` text is preserved, extended with the probed
  candidates.

## Coverage Map

- Story 1 (interactive recovery) → CLI guardrail prompt + `ProbeFallback`.
- Story 2 (non-interactive contract) → fallback failure report + re-run line.
- Story 3 (orchestrator relay) → Roundfix Skill guidance (both copies).
- Story 4 (honest effective selection) → existing effective-selection
  persistence, exercised by tests with a confirmed fallback.
- Core Feature 6 (no functional candidate) → preserved selection error with
  probe summary.

## Integration Points

acpx/codex-acp/claude-code-acp through the existing disposable-session
mechanics only (`applyDisposableSelection` building blocks). Probes reuse the
unique preflight session naming and always close sessions on every path.

## Testing Approach

Existing seams only: the argument-recording fake acpx runner in
`internal/agent`, and buffer-captured CLI runs with scripted stdin in
`internal/cli`.

- Unit (`internal/agent`): `ProbeFallback` walks candidates newest-first,
  excludes the failed model, stops at the first functional model, probes
  efforts highest-first, classifies no-settable-effort as model-managed,
  returns ok=false when nothing proves, and closes every disposable session.
- CLI (`internal/cli`): interactive confirm creates the Run with the
  fallback RuntimeSpec (persisted effort/model asserted on the Run row);
  decline and no-candidate exit 2 with no Run; `--no-input`/`--detach` print
  the fallback report with the exact re-run line; doctor/setup never prompt.
- Docs: `roundfix skills check` and the sync gate in `make verify`.

## Build Order

1. Fallback probe in the agent layer: `ProbeFallback`, candidate/effort
   ordering, model-managed classification, session hygiene, unit tests.
2. Non-interactive guardrail (depends on: 1): fallback failure report with
   the explicit re-run line on resolve, watch, and implement; doctor/setup
   report-only behavior; CLI tests.
3. Interactive confirmation (depends on: 1, 2): the stderr/stdin confirm
   flow, Run creation with the confirmed RuntimeSpec, decline path, CLI
   tests with scripted stdin.
4. Guidance (depends on: 2, 3): Roundfix Skill (both copies) with the
   orchestrator relay rule and re-run recipes, README, and the "Fallback
   Selection" glossary entry in CONTEXT.md.

## Risks & Considerations

- Probe latency: worst case is candidates × (1 session + efforts) adapter
  spawns on an already-failing path; candidates stop at the first proven
  model, efforts at the first accepted value.
- A probe can pass while the model still fails later at prompt time; the
  existing Run failure contracts cover that, as they do today for configured
  selections (ADR-0039 accepts this residual risk).
- Prompting must never block automation: every non-TTY or no-input path is
  fail-with-report by construction, and tests assert no prompt is emitted.
- The Daemon excludes Project Config from Batch staging, so no task may edit
  `.roundfixrc.yml`; this spec needs no config change.

## Decisions

- Probe-discovered, confirmation-gated fallback; no autonomous path. See
  ADR-0041.
- Explicit-flags re-run is the non-interactive confirmation; no new flag.
  See ADR-0041.
- The failed runtime's catalog only — never cross-runtime.
- Confirmation applies to that Run only; config is never written.
