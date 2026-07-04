---
title: "PRD: Run Event seam"
type: PRD
category: enhancement
state: ready-for-agent
labels:
  - enhancement
  - ready-for-agent
blocked_by: []
---

# PRD: Run Event seam

## Problem Statement

Roundfix produces structured ACP stream updates, but they reach exactly one consumer through a hidden mechanism: the Agent runner writes to an output destination and type-asserts at publish time to discover whether that destination is "really" a Live Run View sink. Stream updates carry no Run or Batch identity, so no consumer can correlate output with the Run Database. The conversion from ACP to the local stream model also drops rendering data the product needs: diff old/new text, ACP tool kinds, and plan priorities are discarded before any consumer sees them.

This blocks the entire event-journal sequence. The Run Event Journal, the journaled Agent stream, attach/replay, and daemon event publication all need one publication path that carries Run identity and full-fidelity payloads. Today there is no seam to plug them into, and every fidelity improvement would require touching the runner again.

The user also expects the Agent Console to eventually render like a modern coding-agent cockpit: clear tool-use markers, file diffs, durations, and structured output. The current stream model cannot carry enough information for that, live or replayed.

## Solution

Introduce the Run Event seam: one product event type, one sink interface, and one fanout policy, owned by a new leaf module with no dependencies on other Roundfix modules.

Producers publish Run Events carrying Run identity, Batch when known, event source, event kind, a bounded text summary, a producer-stamped timestamp, and a structured payload. For Agent events, the payload is the raw ACP session update JSON exactly as the ACP Runtime sent it, per ADR 0008. The Agent runner converts its internal stream model into Run Events at publish time and stops knowing anything about terminal UI.

Consumers are sink adapters behind one interface: the Live Run View, the plain-text writer output for non-TTY use, and — in the next PRD — the Run Event Journal. A fanout multiplexer applies the failure policy: critical sinks fail the Run when they fail; best-effort sinks never block or fail producers.

The existing type-assert publication path and its single-consumer interface are removed entirely, not deprecated alongside the new seam.

## User Stories

1. As a developer, I want every Agent output record to carry its Run ID, so that output can be correlated with the Run Database.
2. As a developer, I want every Agent output record to carry its Batch number, so that multi-Batch Resolve Runs are understandable.
3. As a developer, I want Run Events to preserve the raw ACP session update payload, so that rendering fidelity is never limited by what an intermediate struct captured at write time.
4. As a developer, I want diff content, tool kinds, and plan priorities preserved in payloads, so that the Live Run View can later render tool markers, file diffs, and plans the way modern coding-agent cockpits do.
5. As a developer, I want Run Events to carry a bounded text summary, so that list views and plain-text output stay readable when tool output is huge.
6. As a developer, I want Run Events to carry producer-stamped timestamps, so that runtime delays and decision timing can be diagnosed later.
7. As a developer, I want non-TTY and headless output to keep working exactly as today, so that existing CLI output contracts do not regress.
8. As a developer, I want the Live Run View to keep streaming Agent output live, so that the cockpit experience does not regress while the seam is introduced.
9. As a developer, I want a stopped Agent to record a stopped status event before the runner returns, so that interrupted Runs explain themselves.
10. As a developer, I want the Agent runner to publish through one interface regardless of who consumes events, so that adding the Run Event Journal later does not modify runner code.
11. As a developer, I want Roundfix to track current ACP protocol stabilizations, so that journaled payloads use the current session update shape from day one rather than a stale one.
12. As a daemon, I want a critical-sink failure to surface as an error from publication, so that losing a durable consumer can fail the Run instead of being silently swallowed.
13. As a daemon, I want best-effort sink failures contained by the fanout, so that a slow or broken UI never blocks or fails event production.
14. As a TUI, I want to receive Run Events through a non-blocking adapter with a bounded buffer, so that rendering pressure never stalls the Agent runner.
15. As a future attach client, I want live consumers and durable consumers to receive identically structured Run Events, so that live and replay rendering can share one renderer.
16. As a test author, I want a fake sink to capture published Run Events, so that conversion fidelity and identity stamping can be proven without a terminal or a database.
17. As a test author, I want the fanout failure policy to be provable through the sink interface alone, so that tests assert external behavior rather than implementation branches.

## Implementation Decisions

- Create a new leaf module that owns the Run Event type, the event source and kind vocabularies, the sink interface, the fanout multiplexer, and the summary-bounding helper. The module depends only on the standard library. Producers and consumers all import it; it imports none of them.
- The Run Event type carries: Run ID, Batch number with zero meaning outside any Batch, event source, event kind, Review Issue reference when known, ACP tool ID and tool state when known, bounded text summary, producer-stamped time, and a raw JSON payload.
- Event sources are: agent, daemon, verification, git, and review source. Only the agent source is produced in this slice.
- Event kinds are namespaced typed strings. The initial set mirrors the existing stream model: agent message, agent thought, agent tool started, agent tool updated, agent plan, agent status, and agent raw. Daemon kinds arrive in later PRDs. Readers must treat unknown kinds as skippable, never as errors.
- The journal cursor is not part of the Run Event type. Cursor allocation belongs to the Run Event Journal at append time, so adapters without durable storage never see one.
- The sink interface is a single context-aware publish method returning an error. Context is the first argument because durable adapters perform IO.
- The fanout multiplexer is configured with critical sinks and best-effort sinks. Critical sink errors propagate to the producer. Best-effort sink errors are recorded and swallowed. Best-effort delivery must never block producers.
- The Run Event wraps rather than replaces the Agent stream model: the ACP stream update remains the Agent-internal protocol model, and the runner converts it to a Run Event at publish time, stamping Run ID and Batch number from the execute request it already receives. Timestamps come from an injectable clock.
- For agent-source events, the payload is the raw ACP session update JSON verbatim, per ADR 0008. Producers must not re-serialize or prune it. The summary, generated from the existing stream formatting and bounded by the new module's helper, serves list rendering.
- The runner output path is fully replaced: Agent runners accept a sink instead of a writer, the type-asserting publication function and the old single-consumer sink interface are deleted, and the non-TTY text path becomes a writer-sink adapter that formats events with the existing stream formatting. The non-ACP exec runner publishes raw-kind events.
- The Live Run View implements the sink interface directly with a bounded buffer and counted drops, so it remains non-blocking under rendering pressure. Fine-grained coalescing remains renderer work for the attach/replay PRD.
- Stop Request behavior: when the runner's context is canceled, it publishes an agent status event with the stopped state before returning, so later daemon slices can rely on the stop being visible in the event stream.
- Dependency currency is the first task, in an isolated commit: upgrade the ACP Go SDK from the pinned stale version to the current release so raw payloads carry the current session update shape, and upgrade Bubble Tea and its companion libraries to current patch releases. If the SDK upgrade surfaces breaking changes, they are resolved inside this first task before seam work begins.

## Testing Decisions

- Good tests prove externally visible behavior at the sink seam: what events a producer publishes, with what identity, kind, summary, and payload — never which internal functions ran.
- A fake sink captures published events to prove: raw ACP payload byte-equality with what the runtime sent, Run ID and Batch stamping from the execute request, kind mapping for every current stream update variant, and stopped-status publication on cancellation.
- Writer-sink tests prove the non-TTY text output contract is preserved. The existing stream formatting tests move to this adapter instead of asserting private formatting helpers.
- Fanout tests prove the failure policy through the interface: a critical sink error propagates to the publisher; a best-effort sink error does not; a blocked best-effort sink does not stall publication.
- Live Run View tests prove bounded-buffer behavior and counted drops without a real terminal, following the existing TUI test seams.
- Race tests are required for the fanout and the non-blocking TUI adapter, using explicit goroutine completion rather than sleeps.
- Prior art: existing Agent runner tests with fake runners, and existing TUI stream tests, both of which migrate to the new seam rather than being deleted.

## Out of Scope

- The SQLite Run Event Journal schema and store APIs.
- Journal persistence of Agent stream events.
- Attach, replay, or cursor-based reads.
- Daemon and Watch Run event kinds and publication.
- Timeline renderer redesign, chunk coalescing strategy, or rich diff rendering.
- Event retention, pruning, or export.
- Changing ACP Runtime selection, probing, or model behavior.

## Further Notes

This PRD is now the first slice of the event-journal sequence. Every later slice consumes this seam: the Run Event Journal becomes the first critical sink, the attach/replay Live Run View consumes the same event structure from the journal, and daemon orchestration publishes its own kinds through the same interface. The deliberate consequence is that the Agent runner is touched once, here, and never again as consumers are added.

The raw-payload rule and its trade-offs are recorded in ADR 0008. The glossary defines Run Event, Run Event Journal, and Attach.
