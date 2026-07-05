---
task: task_02
spec: 0002-acpx-migration
status: pending
type: backend
complexity: medium
---

# Task 02: Implement the Agent Session lifecycle

## Overview

Give the acpx runner explicit session lifecycle per ADR-0018: ensure-once session creation named by the Run, verbatim full-access mode application, cooperative cancel, and `EndSession`. The `Runner` interface and `ExecuteRequest` grow their session surface and every existing fake extends trivially. Still not wired into Runs (that is task_04); verifiable alone through the helper-process rig from task_01.

## Requirements

1. MUST add `SessionRef` (name `roundfix-<run-id>`) to the agent package, extend `ExecuteRequest` with the session reference, and extend the `Runner` interface with `EndSession(ctx, runtime, session)`; all existing fakes and tests compile with trivial extensions.
2. MUST ensure the session exactly once per session name per runner instance before the first prompt: `acpx <agent> sessions ensure --name <session> --cwd <workdir>`; later prompts reuse it without re-ensuring; exit 4 from a prompt is classified as an infrastructure error (a Roundfix sequencing bug), per the TechSpec.
3. MUST apply the full-access opt-in after ensure and before the first prompt, per ADR-0011 verbatim ids: `acpx <agent> set-mode <FullAccessMode> -s <session>` (`full-access` codex, `bypassPermissions` claude; OpenCode untouched). For codex, apply the `danger-full-access` sandbox preset through acpx session config options, verifying the exact key against the pinned adapter and recording the verified invocation in this task's Result; if the preset is unavailable, degrade gracefully with a journaled warning while still applying the mode id.
4. MUST implement cooperative cancellation: when the context is canceled mid-prompt, send `acpx <agent> cancel -s <session>` and fall back to killing the prompt process — mirroring acpx's own SIGINT behavior; the classified outcome follows Stop Request semantics.
5. MUST implement `EndSession` as best-effort `acpx <agent> sessions close -s <session>`: failures are logged, never fatal (acpx's idle TTL is the backstop).
6. MUST keep mode/sandbox application failures fatal for session setup when full access was requested (fail fast rather than silently running sandboxed), matching the existing ADR-0011 behavior.

## Subtasks

- [ ] `SessionRef`, `ExecuteRequest.Session`, and the `EndSession` interface extension
- [ ] Ensure-once sequencing before the first prompt
- [ ] Full-access mode and codex sandbox mapping with fail-fast semantics
- [ ] Cooperative cancel with kill fallback
- [ ] Best-effort `EndSession`

## Acceptance Criteria

- [ ] Two prompts against one session produce exactly one `sessions ensure` invocation; a second session name re-ensures (asserted via captured invocations).
- [ ] Full-access tests assert the exact `set-mode` invocation per runtime, the codex sandbox application, the no-op for OpenCode, and that a failing `set-mode` fails session setup when full access was requested.
- [ ] Canceling the context mid-prompt issues the cancel invocation and classifies the outcome as a Stop Request; the fallback kill path is covered.
- [ ] `EndSession` issues the close invocation and swallows a non-zero exit with a logged warning.
- [ ] The full existing suite passes with only mechanical fake extensions (no behavioral edits to existing tests).

## Verification

- `rtk go test ./internal/agent/` — expected: all tests pass.
- `rtk go test ./...` — expected: full suite passes.

## References

`_prd.md` → User Stories 1, 3, 4; Core Features 2, 4, 7, 8. `_techspec.md` → Interfaces, acpx invocation mapping, Build Order 2, Risks (codex sandbox key). ADR-0011, ADR-0018.
