---
status: done
created_at: 2026-07-30
updated_at: 2026-09-08
absorbed_by: 0127-durable-unattended-spec-workflow
---

# The autonomous loop orders QA before its own preconditions (2026-07-30)

The loop discipline shipped in `rule.autonomous.loop` and in `docs/agents/autonomous-work.md` prescribes: implement the Task Graph, request the QA gate once when the graph closes, open the Pull Request, watch until Clean, merge. Spec 0053 needed four QA cycles to discover that this order cannot reach `pass` for any Spec whose acceptance observes its own Pull Request, and the correction arrived in two stages.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/findings/2026-07-30-the-autonomous-loop-orders-qa-before-its-own-preconditions.md`.
