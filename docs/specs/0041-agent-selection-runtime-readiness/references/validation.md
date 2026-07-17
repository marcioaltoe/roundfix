# Agent Selection Runtime Readiness Validation

Validated on 2026-07-17 against the local Secondbrain, primary web sources,
and disposable ACP Sessions. No Agent prompt was sent during runtime proof.

## Verdict

The proposed capability-driven proof boundary is correct, with one root-cause
correction: `gpt-5.6-sol / high` is supported by the current official Codex ACP
adapter. The observed rejection came from a setup-generated bare `codex-acp`
override that selected the legacy Zed adapter and masked ACPX's current built-in
adapter.

Roundfix must therefore prove both layers:

1. effective adapter command, official package lineage, and supported version;
2. exact runtime/model/reasoning state advertised and accepted by a disposable
   ACP Session.

Sol/high remains the generated Preferred Selection. GPT-5.5/xhigh remains the
generated fallback by explicit product policy. Terra and Luna are valid,
advertised official identifiers, but are not generated defaults.

## Secondbrain Evidence

The Secondbrain index was consulted first, followed by targeted `qmd query`
searches for runtime capability discovery, fail-before-mutation CLI design,
and Agent Selection proof. The following files informed the decision:

- `~/dev/secondbrain/wiki/index.md` — workspace routing and query procedure.
- `~/dev/secondbrain/projects/skills/mirror/skills/03-engineering-design/agentic-cli-design/references/principles.md`
  — machine-readable, safe-by-default, observable, and introspectable CLI
  contracts.
- `~/dev/secondbrain/projects/roundfix/mirror/docs/specs/0035-agent-selection-profiles/_techspec.md`
  — exact Agent Selection tuple, fallback, and validation boundaries.
- `~/dev/secondbrain/projects/roundfix/mirror/docs/adr/0039-model-availability-preflight-uses-a-disposable-agent-session.md`
  — disposable-session authority and cleanup requirements.
- `~/dev/secondbrain/projects/roundfix/mirror/docs/adr/0050-configured-fallbacks-activate-after-notification.md`
  — pre-proven fallback activation only after user-visible notification and
  before Agent work.

These sources support live, bounded, fail-before-mutation proof and do not
support static model-family assumptions or silent degradation of reasoning
effort.

## Primary Web Evidence

- [ACP Session Config Options](https://agentclientprotocol.com/protocol/v1/session-config-options)
  defines advertised `configOptions` as the selection authority. Select values
  must come from advertised options, and `session/set_config_option` returns the
  complete current configuration because dependent options may change.
- [ACPX agent configuration](https://github.com/openclaw/acpx/blob/main/docs/agents.md)
  documents the built-in Codex command as
  `npx -y @agentclientprotocol/codex-acp`.
- [ACPX 0.9.0 release](https://github.com/openclaw/acpx/releases/tag/v0.9.0)
  records the switch to the official Agent Client Protocol Codex adapter and
  advertised model IDs.
- [Official Codex ACP adapter](https://github.com/agentclientprotocol/codex-acp)
  documents model and reasoning-effort support.
- [Legacy adapter migration](https://github.com/zed-industries/codex-acp/pull/206)
  records that active development moved from the Zed repository to
  `agentclientprotocol/codex-acp`.

## Local Toolchain Evidence

The bounded environment inspection reported:

```text
ACPX:                                0.12.0
Codex CLI:                           0.144.5
legacy @zed-industries/codex-acp:    0.16.0
official @agentclientprotocol/codex-acp: 1.1.4
```

The effective ACPX User Config contained a bare `codex-acp` command. On this
machine, that command resolved to the legacy global package. ACPX's own current
built-in would instead invoke the official package.

The official adapter version was confirmed through the package registry and
its executable:

```text
npm view @agentclientprotocol/codex-acp version
1.1.4

npx -y @agentclientprotocol/codex-acp --version
@agentclientprotocol/codex-acp 1.1.4
```

## Disposable-Session Proof

A disposable session was created with the official adapter and Sol selected:

```text
acpx --cwd /Users/marcio/dev/roundfix \
  --agent "npx -y @agentclientprotocol/codex-acp" \
  --model gpt-5.6-sol sessions ensure \
  --name roundfix-spec0041-new-adapter
```

Setting `reasoning_effort` to `high` succeeded:

```text
acpx --cwd /Users/marcio/dev/roundfix \
  --agent "npx -y @agentclientprotocol/codex-acp" \
  --model gpt-5.6-sol set reasoning_effort high \
  -s roundfix-spec0041-new-adapter

config set: reasoning_effort=high (5 options)
```

The resulting bounded capability projection contained:

```text
models:
  gpt-5.6-sol
  gpt-5.6-terra
  gpt-5.6-luna
  gpt-5.5
  gpt-5.4
  gpt-5.4-mini
  gpt-5.3-codex-spark

reasoning_effort:
  low
  medium
  high
  xhigh
  max
  ultra

effective selection:
  model: gpt-5.6-sol
  reasoning_effort: high
```

The session was then closed successfully. Private ACPX persistence was used
only to cross-check this one research run after the public command response;
the product design explicitly forbids private files as production authority.

## Resulting Requirements

- Setup must not write a bare adapter override based only on executable
  presence.
- Setup, Doctor, profile commands, and operational Preflight Validation must
  share adapter identity and exact capability proof.
- Production proof must use public ACP/ACPX machine-readable responses, not
  `~/.acpx/sessions` or private Codex caches.
- An explicit non-empty reasoning effort must never be converted to an empty
  model-managed value after rejection.
- Legacy or unknown adapters fail before config, Run, worktree, or artifact
  mutation and receive one deterministic migration action.
