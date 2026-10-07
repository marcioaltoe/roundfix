---
status: done
created_at: 2026-07-28
updated_at: 2026-09-08
absorbed_by: 0053-qa-gate-reachability-and-verdict-semantics
---

# QA gate — acceptance rows that need a live Pull Request are structurally unreachable, so their Spec can never be archived (2026-07-28)

The QA Agent runs in a Run Worktree checked out on `roundfix/run-<run-id>`. A Pull Request belongs to the user's feature branch, never to that Run Branch, so `gh pr view` from the QA surface always reports no Pull Request — no matter how many are open. Combined with the rule that Agents never push, any acceptance row whose journey requires a live Pull Request or reaches Final Push cannot be executed by the gate.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-07-28-qa-gate-cannot-reach-pull-request-journeys.md`.
