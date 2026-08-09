# A proof that can refuse — Technical Spec

## Executive Summary

`ProveExactSelection` establishes what a runtime advertises by ensuring a
disposable session **with** the requested model and reading what comes back. On
`claude` that read is contaminated: the adapter reports the requested model as
current and adds it to its own list of options, so membership is decided against
an answer the question wrote.

The change is one extra read. Before asking about a selection, ensure the
disposable session without an override and keep that catalogue; then apply the
selection and decide membership against the earlier read. Everything downstream —
encodings, reasoning proof, receipts — is untouched.

A second, unrelated line of the same command is repaired here because it is in
the same failure path: a failed proof appends a session-close error for a session
that was never opened.

## Project Constraints

- Identifier strategy: not applicable — no new persisted entity; the proof keys
  off the existing runtime, model and effort tuple.
  Source: `docs/agents/domain.md`.
- Authentication and HTTP: not applicable — proof runs over ACP against a local
  adapter process. Source: `docs/agents/cli.md`.
- Active ADR obligations: applicable — ADR-0050 requires preflight to prove every
  configured tuple and substitute none, which a proof that cannot refuse
  defeats. ADR-0104 requires an acceptance row on outside evidence. This Spec
  adds ADR-0112. Source: `docs/agents/domain.md`.
- Tooling authority: not applicable — every file below is ordinary source under
  `internal/agent`. Source: `docs/agents/agent-instructions.md`.

## Vocabulary Contract

- emits: `internal/agent/selection_assignment.go`
  pattern: `catalogue_(contaminated|unavailable)`
  documented-in: `CONTEXT.md`

## System Architecture

No new package, no new command. Two functions change.

**`ProveExactSelection` (`internal/agent/selection_assignment.go`)** gains a
catalogue read as its first step, before any selection is applied. The result is
carried alongside the existing capabilities so membership is decided against it.

**`applyDisposableSelection` (`internal/agent/acpx_runner.go`)** keeps applying
the exact selection and observing effective state. What changes is that its
observation is no longer the only source of truth about what exists.

The existing stderr-matched refusal in `modelNotAdvertisedFromStderr` stays. It
is a fast path: when an adapter does refuse, Roundfix reports its message rather
than a membership verdict of its own, because the adapter's wording names the
advertised set more precisely than a catalogue read can.

## Implementation Design

### Interfaces

```go
// RuntimeCatalogue is what a runtime advertises before any selection is
// requested of it. Read from a disposable session ensured without an override.
type RuntimeCatalogue struct {
    Models   []string
    Efforts  []string
    // Contaminated records that the runtime echoed a later request back into
    // its own advertisement, so a passing proof rests on Roundfix's read
    // rather than on the runtime's refusal.
    Contaminated bool
}

func (runner ACPXRunner) readRuntimeCatalogue(
    ctx context.Context, runtime RuntimeSpec, sessionName, workDir string,
) (RuntimeCatalogue, error)

func (catalogue RuntimeCatalogue) AdvertisesModel(requested string) bool
```

`AdvertisesModel` reuses the canonical binding already used by capability
retention, so `opus` continues to bind `opus[1m]`.

### Data Models

`SelectionProof` gains the catalogue it was decided against, so readiness output
and evidence can name the source of the verdict rather than asserting it.

### API Contracts

`profiles validate` gains one refusal classification and one advisory:

- `model_not_advertised` — unchanged wording, now reachable on `claude` as it
  already is on `codex`.
- `catalogue_contaminated` — advisory, not a refusal. Recorded when the runtime
  echoed the request back into its advertisement, so a maintainer can tell that
  the pass rests on Roundfix's read.

## Coverage Map

| PRD goal | Component |
| --- | --- |
| 1 — an unoffered model is refused on every runtime | `readRuntimeCatalogue` plus the membership check in `ProveExactSelection` |
| 2 — the refusal names the advertised set | The catalogue is carried into `ModelNotAdvertisedError.Advertised` |
| 3 — proof stays token-free | The catalogue is an ACP read on the disposable session; no prompt |
| 4 — one actionable diagnosis | The session-close error is recorded, not appended, when the session was never created |
| 5 — codex and opencode unchanged | The stderr fast path is preserved and its tests are unchanged |

## Integration Points

- `internal/agent/selection_assignment.go` — catalogue read, membership,
  `SelectionProof` field.
- `internal/agent/acpx_runner.go` — the close-error composition on a session
  that was never opened.
- `internal/cli` — readiness output naming the catalogue source.

## Testing Approach

The corpus is captured first, against the live adapters, and it is the whole
point of this Spec: today `claude` / `opus-9-does-not-exist` / `high` proves
`passed`, and that fact is recorded as a declared break before anything changes.

The acceptance case no fixture substitutes for is a live proof on each of the
three runtimes: a nonexistent model must be refused on `codex`, on `claude`, and
on `opencode`, with the advertised set named. Fixtures cover the contaminated
read, because reproducing an adapter that rewrites its own advertisement is
exactly what a fixture is for.

## Build Order

1. **Characterization corpus.** Record that a nonexistent `claude` model proves
   `passed` today, and that a `codex` one is refused. Declares the break.
   Depends on nothing.
2. **The catalogue read.** `readRuntimeCatalogue` on the disposable session,
   with the contaminated-read fixture. Depends on step 1.
3. **Membership decides the verdict.** `ProveExactSelection` refuses a model the
   catalogue does not advertise, naming the set; the stderr fast path is
   preserved. Depends on step 2.
4. **A diagnosis that stops where it stops being useful.** The close error for a
   session that was never created is recorded rather than appended. Depends on
   nothing; runs parallel to 2 and 3.
5. **QA gate.** Depends on steps 3 and 4.

## Risks & Considerations

**One more round trip per distinct tuple.** Readiness already ensures a
disposable session per tuple; this adds one `sessions show` before the
selection is applied. Measured cost belongs in the QA gate rather than an
estimate here.

**An adapter that contaminates the catalogue read too.** If a runtime were to
rewrite its advertisement even without an override, the catalogue would be as
unreliable as the current read. That is why contamination is recorded rather
than assumed absent: the evidence names what the verdict rests on, so the next
maintainer to hit it starts from a fact.

**A model that exists but is not advertised.** A runtime may support a model it
does not list. Under this design that selection is refused. The PRD accepts it:
what a Run receives is what the runtime advertises now, and a selection that
depends on an unlisted model is exactly the kind that fails mid-Run.

## Decisions

- Membership is decided against a catalogue read before the request, per
  ADR-0112.
- The stderr fast path is preserved rather than replaced, because an adapter
  that refuses names its advertised set more precisely than a catalogue read.
- Contamination is recorded as an advisory on a passing proof, not as a refusal.
  Refusing every runtime that echoes would refuse `claude` entirely, which is
  not the defect being fixed.
