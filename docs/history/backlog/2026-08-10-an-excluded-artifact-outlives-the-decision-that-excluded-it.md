---
type: fix # feat | fix | perf | refactor
status: deferred
created: 2026-08-10
spec: null # Spec slug when status: promoted
reason: null # required when status: declined
---

# An excluded artifact outlives the decision that excluded it

Answering a Baseline decision so that it excludes an artifact removes that artifact from the Setup Manifest but leaves its bytes on disk. Measured on 2026-08-10: setting `triage.external` to `false` and running `roundfix baseline update --yes` dropped the `external-triage` module, `root.external-triage`, and `guide.external-triage` from `docs/agents/setup-context.json`, and reported one file change. It did not delete `docs/agents/external-triage.md`, and it did not remove the `root.external-triage` region from `AGENTS.md`. Both survived with their `setup-context-driven` markers intact.

Full text in Git at `a30dc847a037fe584812b4471aa7da1701a9ec59`: `docs/history/backlog/2026-08-10-an-excluded-artifact-outlives-the-decision-that-excluded-it.md`.
