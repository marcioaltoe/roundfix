---
status: accepted
created_at: 2026-10-19T19:30:00Z
updated_at: 2026-10-19T19:30:00Z
deprecated_at: null
superseded_by: null
---

# A full-access request selects the mode the Agent Session advertises

On 2026-10-18 the Conexus QA gate ran on codex-acp 2.2.2 and acpx 0.19.6, and
`defaults.agent_full_access: true` could not help it. Roundfix sends the
fixed Codex mode `full-access`, and that adapter refuses it with ACP -32602.
It accepts `agent-full-access` and `read-only`. The Codex sandbox preset was
sent only after a mode literally named `full-access` succeeded. The adapter's
source names its modes `read-only`, `workspace-write`, `agent` and
`agent-full-access`, and advertises them in a `mode` config option described
as the "Approval and sandboxing preset for the session". Its
`agent-full-access` mode carries the `danger-full-access` sandbox and the
`never` approval policy. The Agent Client Protocol says a Client that
supports config options "SHOULD use them instead of the `modes` field". On
2026-10-19 the maintainer chose "Descobrir pelos availableModes
(Recomendado)".

## Decision

**Each runtime names its full-access candidates in preference order.** Codex
names `agent-full-access`, then `full-access`. Claude Code names
`bypassPermissions`. A runtime without candidates (OpenCode, Cursor, a command
override) keeps refusing a full-access request as unsupported. A runtime value
built with only a full-access mode name has that name as its single candidate.

**The Agent Session's advertised modes decide.** After Roundfix ensures an
Agent Session, it already reads `acpx sessions show` for the selection
capabilities. The advertised modes are the values of the select config option
whose identifier is `mode`. When no option has that identifier, they are the
values of the first option whose category is `mode`. The Full-Access Mode is
the first candidate, in the runtime's order, that those values contain,
compared exactly. Roundfix sends `set-mode` with that mode only, in the
disposable proof and in the Run alike. An adapter that refuses an advertised
mode still fails the session fast, as ADR-0011 decided.

**A session that advertises no candidate is reported degraded.** Roundfix
sends no `set-mode`. The effective access policy is `full-access (degraded: no
advertised mode grants full access; advertised modes: <values>)`, with the
values in advertised order, or `none` when the session advertises no mode
option. The Run event stream receives the status
`full_access_mode_not_advertised`, and a warning names the Agent Session and
the advertised modes. `roundfix profiles validate` and `roundfix doctor`
already print any degraded effective access policy.

**The Codex sandbox preset follows the session's own controls.** After the
Full-Access Mode is set on a Codex session, Roundfix sends `set sandbox_mode
danger-full-access` when the session advertises a `sandbox_mode` config
option. A failure there keeps today's degraded policy and status. When the
session advertises no `sandbox_mode` option, Roundfix sends no preset and the
policy is not degraded, because that adapter's mode is itself the approval
and sandbox preset.

## Considered Options

- Hard-coding `agent-full-access` for Codex was rejected. A fixed name is what
  broke: the adapter renamed its mode, and the next rename would break it
  again.
- Reading the older `modes.availableModes` field was rejected. The acpx 0.19.6
  session record does not keep it. The protocol tells a Client that supports
  config options to use them and ignore `modes`, and its version 2 removes
  `modes` and `session/set_mode`.
- Choosing a mode by its name or description, for example any value that
  contains `full`, was rejected: it would select a mode nobody vetted, so the
  preference stays an exact list.
- Refusing the session when no candidate is advertised was rejected for this
  change. A degraded policy that names the advertised modes is visible in the
  proof, the Doctor and the Run events, and it keeps a Run that can still
  work without full access.
- Always sending the sandbox preset was rejected. codex-acp 2.x advertises no
  `sandbox_mode` option, so every Codex session would be reported degraded
  while its mode already grants `danger-full-access`.
- Editing `~/.codex/config.toml` was rejected again, as in ADR-0011: the
  source measured that ACP sessions ignore its `sandbox_mode` and
  `default_permissions`.

## Consequences

- A Codex adapter that renames its full-access mode needs one new candidate,
  and until then Roundfix reports which modes it offered instead of failing
  every Task.
- This amends ADR-0011, which named `full-access` as the Codex mode and failed
  fast on any refused mode.
