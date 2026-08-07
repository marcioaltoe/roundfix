---
status: pending
created_at: 2026-08-07
updated_at: 2026-08-07
kind: finding
---

# Claude Agent Selections are never proven (2026-08-07)

`roundfix profiles configure --dry-run` proves `codex` tuples against the ACP
adapter and accepts any `claude` tuple. A maintainer can pin a Claude model that
does not exist and the preview reports success. The cause is now established: the
Claude adapter refuses nothing.

## Reproduction

An invalid codex model is refused, and the adapter enumerates what it does
advertise:

```
codex / sonnet / xhigh
  → exit non-zero, classification: model_not_advertised
  → advertised Agent Models: gpt-5.6-sol, gpt-5.6-terra, gpt-5.6-luna,
    gpt-5.5, gpt-5.4, gpt-5.4-mini, gpt-5.3-codex-spark
```

An invalid claude model is accepted, in either position:

```
claude / __probe_invalid__ / high  as general.preferred    → exit 0, preview OK
claude / __probe_invalid__ / high  as general.fallbacks[0] → exit 0, preview OK
```

The difference is the runtime, not the position.

## Established cause

Opening a session against `@agentclientprotocol/claude-agent-acp` 0.63.0 and
reading the `model` entry of its advertised `configOptions` returns:

```
default              Opus 5 with 1M context
opus[1m]             Opus 5 with 1M context
claude-fable-5[1m]   Fable 5
sonnet               Sonnet 5
haiku                Haiku 4.5
__probe_invalid__    Custom model        ← the probe value, echoed back
```

The adapter appends the requested identifier to its own list as a **Custom
model** rather than rejecting it. There is no refusal for roundfix's proof to
observe, so the proof cannot distinguish a real model from a typo. The codex
adapter, by contrast, refuses and names its advertised set.

## A second, independent defect

`internal/agent/catalog.go` hardcodes the Claude models offered during
interactive configuration as `claude-opus-5`, `claude-fable-5`, and
`claude-opus-4-8`. Only the middle one matches the adapter, which advertises
`opus`, `claude-fable-5`, `sonnet`, `haiku`, and `default` after roundfix's
opaque-identifier parsing takes the part before `[`.

So the interactive picker offers two identifiers the adapter does not advertise,
and omits three it does — including `sonnet` and `haiku`. Because the adapter
accepts anything, choosing one of the two wrong entries produces a working
configuration that fails only inside a Run.

## Why it matters

Preflight exists to make a configured Selection trustworthy before a Run. For
the claude runtime that guarantee does not hold, and the interactive picker
actively steers toward identifiers that will fail. The cost lands inside a Run,
where it burns an Agent turn instead of a preview.

## Route

Not fixed here. Two separable repairs, and the first has a real design choice:

- **Proof.** The adapter cannot be made to refuse, so roundfix must either
  compare the requested identifier against the advertised list itself and refuse
  locally, or report honestly that a claude Selection cannot be proven. The
  second is acceptable only if the CLI says so instead of implying success.
- **Catalog.** The hardcoded Claude list should be read from the adapter rather
  than maintained by hand, or it will drift again. Correcting the current values
  without removing the hand-maintenance only resets the clock.
