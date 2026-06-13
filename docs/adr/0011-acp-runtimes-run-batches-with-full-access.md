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
