---
task: task_07
spec: 0052-claude-adapter-standardization
status: completed
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

## Result

Implementation:

- The Setup command reference now documents the official Codex and Claude
  adapter lineages, pins, deterministic install actions, and migration offers
  for legacy Codex plus both legacy Claude lineages. The Doctor reference now
  describes the deduplicated, runtime-sorted aggregate check and shows the
  shipped Claude-plus-Codex `adapter:` evidence line.
- The usage and configuration guides now render the built-in frontend
  Preferred Selection and working Claude examples as
  `claude / opus / xhigh`. Dated catalog labels remain identified as advisory,
  not working frontend selections.
- Both selection guides explain that advertised identifiers are opaque when
  the adapter exposes an independent reasoning control: `opus[1m]` is
  selectable as printed, `[1m]` is a context-window annotation rather than a
  reasoning effort, and reasoning remains a separate tuple field. Both guides
  link ADR-0079.

Focused checks:

- Before editing, focused repository searches found the stale
  `claude / claude-opus-5 / xhigh` frontend tuple in both guides, the
  `1.1.4` Codex pin in all three adapter-contract sections, and the
  single-runtime `adapter: ok (...)` Doctor placeholder.
- `rtk git -c core.fsmonitor=false diff --check` exited `0`.
- `rtk proxy rg -l "npm install -g @agentclientprotocol/claude-agent-acp@0\\.63\\.0" docs/user-guide/commands.md docs/user-guide/usage.md docs/user-guide/configuration.md`
  exited `0` and listed all three user-guide interfaces.
- `rtk proxy rg -n "claude / claude-opus-5|model: claude-opus-5|model: claude-fable-5|npx -y @agentclientprotocol/codex-acp@1\\.1\\.4|npm install -g @agentclientprotocol/codex-acp@1\\.1\\.4" docs/user-guide/commands.md docs/user-guide/usage.md docs/user-guide/configuration.md`
  exited `1`, the expected no-match result for stale working tuples and stale
  Setup/install commands.
- `rtk proxy rg -n "@zed-industries/claude-|claude-code-acp" docs/user-guide`
  listed only the short former `claude-code-acp` lineage name and the
  wrong-scope `@zed-industries/claude-agent-acp` lineage; it found no
  prohibited fully scoped former package name.
- `rtk proxy test -f docs/adr/0079-independent-reasoning-controls-make-model-identifiers-opaque.md`
  exited `0`.
- The Task's declared `## Verification` commands were not run; the Daemon owns
  that gate.

Acceptance evidence:

- Official Claude contract: commands, usage, and configuration each name
  `@agentclientprotocol/claude-agent-acp`, pin `0.63.0`, and
  `npm install -g @agentclientprotocol/claude-agent-acp@0.63.0` beside the
  corresponding Codex contract.
- Advertised working selections: every frontend or working Claude tuple in
  the user guide uses `opus`; remaining `claude-opus-5` and
  `claude-fable-5` mentions are explicitly dated advisory catalog labels.
- Doctor and Setup examples: the Doctor block contains sorted Claude and
  Codex evidence, and every current Codex proposal/install command uses
  `1.1.5`; `1.1.4` appears only as an earlier pin that Setup can migrate.
- Deprecated package exclusion: the focused lineage search found no user-guide
  occurrence of the prohibited fully scoped former Claude package.

Follow-up:

- `docs/agents/autonomous-work.md` still describes the built-in frontend
  profile as `claude / claude-opus-5 / xhigh`. That agent-facing guide is
  outside task_07's declared user-guide interfaces and was left unchanged.
