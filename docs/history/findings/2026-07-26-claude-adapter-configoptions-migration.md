---
status: done
created_at: 2026-07-26
updated_at: 2026-09-08
absorbed_by: 0052-claude-adapter-standardization
---

# Agent Selection — the configured Claude ACP adapter is deprecated and advertises no capability evidence (2026-07-26)

`claude/claude-opus-5/xhigh` fails Agent Selection preflight with `capability_evidence_invalid: missing_config_options`. The failure is not about the model: `claude-opus-5` is advertised and gets selected successfully. Roundfix is pointed at `@zed-industries/claude-code-acp`, which npm deprecated in favour of `@agentclientprotocol/claude-agent-acp` — the same rename Roundfix already followed for Codex. The replacement adapter advertises the exact two controls Roundfix requires, so the fix is an adapter migration, not a parser change.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-07-26-claude-adapter-configoptions-migration.md`.
