---
status: deferred
created_at: 2026-09-08
updated_at: 2026-09-08
kind: finding
absorbed_by: 0125-repository-identity-and-run-branch-policy
---

# Run naming — internal branches ignore the repository's selected prefix (2026-09-08)

The Baseline tells Agents to use the repository's selected branch prefix, but Run and Task branches use a fixed namespace. A prefixed initial branch cannot prevent the executor from creating additional nonconforming branches.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-09-08-run-branches-ignore-the-selected-prefix.md`.
