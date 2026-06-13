# ACP Runtimes run Batches with full access

Roundfix switches every ACP session it creates to the runtime's full-access
mode before sending the Batch prompt: `full-access` for Codex through
`codex-acp` and `bypassPermissions` for Claude Code through
`claude-agent-acp`. For Codex this selects the `danger-full-access` sandbox
preset with approval prompts disabled, because the default workspace-write
sandbox blocks all network access — including localhost — which makes
database-backed verification commands fail with `connect EPERM` even when
the service is running on the developer's machine.

The mode is set through the Agent Client Protocol (`session/set_mode`), not
by editing the user's runtime configuration files, so Roundfix behavior does
not depend on per-machine `config.toml` state. If the runtime rejects the
mode, session creation fails fast instead of silently running a sandboxed
Batch that cannot verify. OpenCode sessions keep that runtime's defaults
until it exposes an equivalent mode.
