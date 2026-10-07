---
status: done
created_at: 2026-08-04
updated_at: 2026-09-08
absorbed_by: 0125-repository-identity-and-run-branch-policy
---

# 2026-08-04 — Branch Integrity Preflight prescribes a remedy that reintroduces superseded work

A QA gate failed a Spec because the Supervisor had committed protected tooling (a Baseline profile) onto the Spec branch. The correct remediation, named by the gate itself, was *"remove the out-of-scope profile change from this Spec's ancestry"*. The branch was unpushed, so it was rebuilt: the offending file was dropped, the platform work moved to its own PR, and every Spec commit replayed.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-08-04-branch-integrity-preflight-prescribes-a-remedy-that-reintroduces-superseded-work.md`.
