---
spec: 0002-acpx-migration
status: active
created: 2026-07-05
surfaces: [cli, infra]
---

# ACPX Migration

Roundfix reaches its Agents through a hand-rolled ACP client layer that spawns a fresh runtime process for every Batch and re-implements session plumbing, permission answering, and filesystem jailing by hand. Every Work Item pays a runtime cold start, a crashed runtime ends the whole Run as Failed, and the client layer is the largest and most fragile subsystem Roundfix maintains. Migrating the agent layer to acpx — a headless ACP client with persistent named sessions, crash resume, and a raw protocol stream — makes multi-item Runs warm and crash-resilient while shrinking Roundfix to process orchestration it can actually own.

## Goals

- Every Run drives its ACP Runtime through acpx with behavior parity: commands, flags, stdout contracts, exit codes, and both Batch contracts are unchanged for users and Agents.
- One Agent Session per Run: consecutive Work Items reuse a warm session, and a runtime crash mid-Run resumes instead of failing the Run.
- The Run Event Journal keeps carrying raw ACP payloads with no schema change, so Attach and the Live Run View replay identically.
- The hand-rolled ACP client layer is deleted; Roundfix's agent layer becomes process orchestration of a pinned acpx version.

## User Stories

1. As a developer running resolve or implement over many Work Items, I want consecutive Batches and Tasks in one Run to reuse a warm Agent Session, so that the Run does not pay a runtime cold start per Work Item.
2. As a developer whose runtime crashes mid-Run, I want the next Work Item to resume the Agent Session transparently, so that a flaky runtime does not end the Run as Failed.
3. As a developer stopping a Run, I want cancellation to stay cooperative and immediate — Stop Request, cooperative cancel, exit 130 — so that stop semantics do not change under acpx.
4. As a developer using the full-access opt-in, I want the runtime modes applied exactly as before, so that Verification commands needing local services keep working.
5. As a developer with a custom stdio runtime command, I want the command override escape hatch to keep working, so that local testing setups survive the migration.
6. As a developer attaching to a Run, I want the Agent timeline rendered from the same raw events as today, so that Attach and replay behave identically.
7. As a developer on a machine without acpx or with the wrong version, I want Preflight Validation to name the missing dependency and the exact install command, so that the failure is actionable before any Run is created.

## Core Features

1. **Single agent layer.** acpx is the only path to ACP Runtimes for review and spec Runs; the version is pinned and verified before Runs start. See ADR-0017.
2. **Agent Session per Run.** A Run creates one named Agent Session at its first Agent work, routes every Work Item through it, and closes it at the Run's terminal outcome; sessions never survive a Run or cross Runs. See ADR-0018.
3. **Runtime selection parity.** Codex, Claude, and OpenCode remain the supported ACP Runtimes through acpx's adapters; the raw command override stays as the stdio escape hatch; the runtime is always named explicitly, never left to acpx's implicit default.
4. **Permission parity.** Blanket approval reproduces today's auto-approve behavior; the full-access opt-in maps to the runtime session modes verbatim per ADR-0011; a real permission policy engine remains deferred work.
5. **Journaling parity.** The raw ACP protocol stream is captured from acpx's machine-readable output into the Run Event Journal, preserving the ADR-0008 contract without schema changes.
6. **Failure mapping.** acpx's documented exit codes map onto existing outcomes: agent errors and permission denials fail the Batch under the ADR-0010 policy, timeouts keep today's timeout semantics, and environment problems (missing binary, wrong version, missing session) surface as Preflight Validation or infrastructure failures with actionable messages.
7. **Cooperative cancellation.** A Stop Request sends the cooperative cancel through acpx, then tears down; Run state and journaled history are preserved exactly as today.
8. **Clean teardown.** Every terminal outcome closes the Agent Session and its background owner process; no orphan processes outlive a Run.

## User Experience

No new commands and no new flags. The visible differences are confined to diagnostics: stderr may name the Agent Session, cold starts disappear after a Run's first Work Item, and a missing or mismatched acpx produces one actionable Preflight Validation message naming the install or upgrade command. Everything else — output lines, exit codes, Interactive Input, the Live Run View — reads exactly as before.

## Non-Goals / Out of Scope

- acpx flows: orchestration stays in the Daemon; Roundfix uses sessions and prompts only.
- A dual backend or SDK fallback — the cutover is total. See ADR-0017.
- Cross-Run session reuse or any warm state that outlives a Run. See ADR-0018.
- The long-lived daemon and work queue, parallel worktree execution, retry and escalation budgets, and the permission policy engine (work-plan items 3, 4, 7, 8).
- Windows support hardening beyond what acpx provides.
- Changes to prompts, Batch contracts, or the Roundfix skill semantics beyond documenting the new dependency.

## Success Metrics

- The full existing test suite passes with contract parity; no review-path or implement-path CLI test changes meaning.
- A multi-Task spec Run shows exactly one runtime spawn in its journal, with later Work Items reusing the warm Agent Session.
- An induced runtime kill mid-Run resumes on the next Work Item and the Run completes normally.
- Attach replays an acpx-era Run with an Agent timeline indistinguishable in structure from a pre-migration Run.
- The Go module no longer depends on the ACP SDK; the agent layer's owned code shrinks accordingly.

## Decisions

- Hard cutover to acpx as the only agent layer, pinned version, Node accepted as an agent-layer dependency. See ADR-0017.
- One Agent Session per Run, named by the Run, closed at the terminal outcome. See ADR-0018.
- Raw payloads are journaled from acpx's machine-readable stdout stream, not from acpx's rotating on-disk session files.
- acpx flows are out of scope; the Daemon remains the only orchestrator.
- Permission behavior stays at parity (blanket approval + ADR-0011 modes); the policy engine stays work-plan item 8.
- The glossary gains **Agent Session**.

## Open Questions

- The exact acpx version to pin at ship time (0.12.0 as of this writing) and the upgrade cadence policy. Default: pin the version current when the techspec lands; upgrades are deliberate commits. Owner: techspec.
