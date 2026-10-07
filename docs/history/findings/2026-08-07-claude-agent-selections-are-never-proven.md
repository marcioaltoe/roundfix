---
status: done
absorbed_by: 0091-a-proof-that-can-refuse
created_at: 2026-08-07
updated_at: 2026-08-26
kind: finding
---

# Claude Agent Selections are never proven (2026-08-07)

`roundfix profiles configure --dry-run` proves `codex` tuples against the ACP adapter and accepts any `claude` tuple. A maintainer can pin a Claude model that does not exist and the preview reports success. The cause is now established: the Claude adapter refuses nothing.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-08-07-claude-agent-selections-are-never-proven.md`.
