---
status: pending
created_at: 2026-08-06
updated_at: 2026-08-06
kind: rollup
members:
  - 2026-08-03-gate-and-review-rounds-need-a-convergence-rule.md
  - 2026-08-04-an-accepted-gap-has-no-terminal-state-so-the-loop-cannot-close.md
  - 2026-08-04-review-runs-halt-autonomous-delivery-on-unrelated-dirty-files.md
  - 2026-08-04-what-still-needs-a-supervisor-between-a-prd-and-a-merge.md
  - 2026-08-05-contract-seams-between-daemon-gate-and-archive.md
  - 2026-08-05-five-frictions-from-a-full-autonomous-spec-night.md
  - 2026-08-05-o-loop-devolve-controle-por-motivos-mecanicos.md
  - 2026-08-05-review-issues-have-no-identity-across-rounds.md
  - 2026-08-05-what-a-six-spec-autonomous-session-asks-roundfix-to-change.md
  - 2026-08-05-what-roundfix-should-do-differently-measured-over-one-queue-night.md
  - 2026-08-06-manual-thread-resolution-is-load-bearing-and-undocumented.md
---

# Review and delivery convergence — mechanical failures need a path back into the loop (2026-08-06)

The delivery findings describe a loop that returns control for recoverable
mechanical states: stale gates, missing review requests, commit-hook failures,
accepted gaps, unresolved conversations, and Review Issues that lose identity
across Rounds. Those states need typed recovery or terminal semantics, not a
Supervisor interpreting local artifacts by hand.

## Consolidated learning

- Review Issues need stable identity across Rounds so accepted, duplicated,
  failed, and still-open outcomes remain auditable.
- A Pull Request must request and observe its Review Source; manual thread
  resolution cannot remain both load-bearing and forbidden by the contract.
- Gate, Daemon settlement, commit, review, and Archive need one convergence
  model with a path from recoverable failure back to the responsible Agent.
- Preflight must report the actionable set instead of one defect per restart,
  and unrelated dirty files must not turn Review work into a delivery halt.

## Live edge

Reconciliation and newer Review Source contracts removed several manual steps.
The rollup remains `pending` until every non-decision interruption has a typed
recovery path and the loop stops returning control for mechanical reasons.
