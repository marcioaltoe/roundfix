---
status: done
absorbed_by: 0105-the-gates-own-economics
created_at: 2026-08-26
updated_at: 2026-09-08
kind: finding
---

# A QA Task's authored Verification is never executed, yet the authoring contract requires one

The authored terminal `qa` Task settles from the QA Report's `verdict:` — that is the contract and it works. What does not close: the authoring contract treats the `qa` Task like any other and requires a `## Verification` proving its own effect, while the runtime never executes it.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-08-26-a-qa-tasks-authored-verification-is-never-executed.md`.
