---
status: done
created_at: 2026-07-27
updated_at: 2026-09-08
absorbed_by: 0052-claude-adapter-standardization
---

# Agent Selection — standardize on the official Claude adapter and stop reading `[...]` as reasoning effort (2026-07-27)

The maintainer directs Roundfix to drop every reference to `@zed-industries/claude-code-acp` and support only `@agentclientprotocol/claude-agent-acp`, including Doctor and Setup. This finding records the empirical verification that the replacement works, the exact code sites that still name the deprecated package, and one **new blocking defect** that only becomes visible after the migration: Roundfix misreads the replacement's advertised model identifiers.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-07-27-claude-adapter-standardization.md`.
