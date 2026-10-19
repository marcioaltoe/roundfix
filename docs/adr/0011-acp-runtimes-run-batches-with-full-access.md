---
status: accepted
created_at: 2026-06-13T00:36:33Z
updated_at: 2026-10-19T21:00:54Z
deprecated_at: null
superseded_by: null
---

# ACP Runtimes can run Batches with explicit full access

Roundfix keeps ACP runtime sandbox defaults unless the user explicitly opts in
with `--agent-full-access` or `defaults.agent_full_access: true`. When enabled,
Roundfix switches the ACP session to the runtime's full-access mode before
sending the Batch prompt: `full-access` for Codex through `codex-acp` and
`bypassPermissions` for Claude Code through `claude-agent-acp`. For Codex this
also selects the `danger-full-access` sandbox preset with approval prompts
disabled.

The opt-in exists because some verification commands need local services such
as databases. Codex's default workspace-write sandbox blocks all network
access, including localhost, which can make those commands fail with
`connect EPERM` even when the service is running on the developer's machine.

The mode is set through the Agent Client Protocol (`session/set_mode`), not
by editing the user's runtime configuration files, so Roundfix behavior does
not depend on per-machine `config.toml` state. If full access is enabled and
the runtime rejects the mode, session creation fails fast instead of silently
running a sandboxed Batch that cannot verify. OpenCode sessions keep that
runtime's defaults until it exposes an equivalent mode.

**Amended by [ADR-0293](0293-a-full-access-request-selects-the-mode-the-agent-session-advertises.md) (2026-10-19).** The runtime's full-access mode is now the first candidate the Agent Session advertises (`agent-full-access`, then `full-access`, for Codex). A session advertising none is reported degraded instead of failing, and the Codex sandbox preset is sent when the session advertises `sandbox_mode`.
