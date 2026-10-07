---
type: perf # feat | fix | perf | refactor
status: declined # open | promoted | declined
created: 2026-08-10
spec: null # Spec slug when status: promoted
reason: resolved directly on 2026-08-10 — the coupling was .git, not the docs tree; test git reads stopped writing the index
---

# One reader in cli still couples verify to the docs tree

After the markdown contracts moved to `internal/docscontract`, touching a markdown file should leave every code package `(cached)`. `internal/speccheck` now behaves that way — proven on 2026-08-10: warm, append one line to a `docs/backlog` entry, re-run, `(cached)`.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/backlog/2026-08-10-one-reader-in-cli-still-couples-verify-to-the-docs-tree.md`.
