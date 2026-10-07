---
status: done
created_at: 2026-07-27
updated_at: 2026-09-08
absorbed_by: 0054-tooling-task-and-verification-hygiene
---

# Verification — sandboxed Agents cannot reach the default Go build cache, and nothing tells them so (2026-07-27)

`make verify` is the authoritative gate for this repository, but an Agent running inside the ACP sandbox cannot always reach the default `GOCACHE` (`~/Library/Caches/go-build` on macOS). When the sandbox denies it, the gate fails **before compilation**, producing a failure that looks like a broken build rather than a denied path.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-07-27-sandboxed-agents-cannot-reach-the-default-go-cache.md`.
