---
task: task_07
spec: 0052-claude-adapter-standardization
status: pending
type: docs
complexity: medium
---

# Task 07: Align adapter docs and user guide

## Overview

Bring the operator documentation in line with the shipped behavior: the
Claude adapter contract described like the Codex one, the multi-runtime
Doctor `adapter:` line, the Setup migration for Claude, the raised Codex pin,
the new frontend default tuple, and selection examples that use advertised
identifiers. Documentation only — no Go or Skill changes.

## Requirements

1. MUST document the supported Claude adapter (official package, `0.63.0`
   pin, official install command) wherever the Codex adapter contract is
   documented, including migration guidance from both legacy lineages.
2. MUST update the documented Doctor output examples to the multi-runtime
   `adapter:` line and the raised Codex pin in Setup examples.
3. MUST update configuration and usage examples that select Claude models to
   the advertised identifiers, replacing `claude-opus-5` tuples with the
   proven `opus` tuple where they describe the frontend default or working
   selections.
4. MUST describe the opaque-identifier rule where model selection is
   documented: bracketed advertised identifiers are selectable as printed,
   and a context-window annotation is not a reasoning effort.
5. MUST NOT edit the Roundfix Skill pair, baseline assets, or any protected
   tooling path — that is task_08's bounded scope.

## Subtasks

- [ ] Update the command reference for Doctor and Setup adapter behavior.
- [ ] Update usage and configuration guides for the frontend default and
      Claude selection examples.
- [ ] Add the Claude adapter migration note beside the existing Codex one.
- [ ] Cross-link ADR-0079 from the selection documentation.

## Acceptance Criteria

- [ ] The user guide names `@agentclientprotocol/claude-agent-acp` with its
      pin and install command wherever the supported Codex adapter is named.
- [ ] No user-guide example selects `claude-opus-5` as a working frontend
      tuple; documented working Claude selections use advertised
      identifiers.
- [ ] Doctor and Setup documented examples match the shipped multi-runtime
      output and raised Codex pin.
- [ ] No path under `docs/user-guide/` names `@zed-industries/claude-code-acp`.

## Context

- interface: `docs/user-guide/commands.md`
- interface: `docs/user-guide/usage.md`
- interface: `docs/user-guide/configuration.md`
- instruction: `docs/adr/0079-independent-reasoning-controls-make-model-identifiers-opaque.md`

## Verification

- `grep -rn 'zed-industries/claude-code-acp' docs/user-guide/ ; test $? -eq 1` — expected: no matches (exit 1).
- `grep -rln 'agentclientprotocol/claude-agent-acp' docs/user-guide/ | grep -q commands.md` — expected: the command reference documents the official adapter.
- `go test -count=1 ./internal/cli/ -run 'TestDocumentation|TestHelp'` — expected: documentation and help contract tests pass.

## References

`_prd.md` → Core Feature 9, Goals 1 and 4, User Experience; `_techspec.md` →
Build Order 7; ADR-0079.
